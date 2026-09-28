// Package bench measures AgentCap's behavior on real command executions.
//
// The benchmark layer is an observer: it runs a command through the normal
// AgentCap pipeline (engine.Run) exactly once and derives every measurement
// from that single captured execution. It never re-runs the command and never
// uses a benchmark-specific reducer.
package bench

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/taqu/agentcap/internal/integration/engine"
	"github.com/taqu/agentcap/internal/integration/protocol"
	"github.com/taqu/agentcap/internal/project"
	"github.com/taqu/agentcap/internal/session"
	"github.com/taqu/agentcap/internal/store"
)

// Measurement describes one command execution observed through AgentCap.
//
// Byte semantics:
//   - RawStdoutBytes / RawStderrBytes are the bytes captured from the child's
//     stdout / stderr before any AgentCap processing. They are the same bytes
//     that are stored and returned by `acap raw <id>`. Capture is capped at
//     exec.MaxOutputBytes per stream; Truncated reports when the cap was hit,
//     in which case the raw counts understate the true output.
//   - RawBytes = RawStdoutBytes + RawStderrBytes.
//   - AgentVisibleBytes is the length of the rendered AgentCap command-result
//     presentation returned to the coding agent (the "@acap <id> ..." header,
//     capsule and any folded-in stderr section), i.e. the text the Claude and
//     Codex integrations hand back to the agent. It does not include the
//     benchmark report itself, storage metadata, or on-disk store size.
type Measurement struct {
	Command []string

	ExitCode int

	RawStdoutBytes int64
	RawStderrBytes int64
	RawBytes       int64
	Truncated      bool

	AgentVisibleBytes int64
	// StatelessVisibleBytes is the full rendered reduction before session
	// comparison. StatefulVisibleBytes is the actual post-session presentation.
	// AgentVisibleBytes remains the B1-compatible alias of the actual visible
	// presentation (and therefore equals StatefulVisibleBytes).
	StatelessVisibleBytes int64
	StatefulVisibleBytes  int64

	// ExecutionDuration is child process start to child process exit.
	ExecutionDuration time.Duration
	// ReduceDuration is time spent in the reducer only.
	ReduceDuration time.Duration
	// ProcessingDuration is AgentCap's post-execution work, from child exit
	// until the agent-visible result is ready. It includes reduction,
	// hashing, writing the result to .acap/store.db and the object store,
	// and stats recording. See engine.Outcome.ProcessingDuration.
	ProcessingDuration time.Duration

	// ResultID is the stored result ID, or "" if the result was not stored.
	ResultID string
	// Presentation is the normal pipeline kind: full, delta, or unchanged.
	// Stateless B1 measurements are always full.
	Presentation string
}

// ReductionRatio returns 1 - AgentVisibleBytes/RawBytes, the single-command
// reduction in agent-visible bytes. ok is false when RawBytes is zero, where
// the ratio is undefined. The ratio is negative when AgentCap's presentation
// is larger than the raw output (e.g. the result-ID header on tiny output).
func (m *Measurement) ReductionRatio() (ratio float64, ok bool) {
	if m.RawBytes == 0 {
		return 0, false
	}
	return 1 - float64(m.AgentVisibleBytes)/float64(m.RawBytes), true
}

// MeasureCommand executes args once through the normal AgentCap pipeline
// in dir and returns the resulting measurement.
//
// The execution is always stateless: no session ID is passed to the
// pipeline, so the presentation is never a delta against an earlier result,
// regardless of ACAP_SESSION_ID. The result is stored normally and is
// inspectable with `acap show` / `acap raw`.
//
// A non-zero exit status of the target command is not an error. An error is
// returned only when the command could not be executed at all.
func MeasureCommand(ctx context.Context, args []string, dir string) (*Measurement, error) {
	return measure(ctx, engine.Run, args, dir, "")
}

// StartSession creates a persistent benchmark session. Its ID is also used as
// the normal AgentCap session ID, giving the benchmark an isolated baseline.
func StartSession(dir string) (string, error) {
	a, err := session.GenerateID()
	if err != nil {
		return "", err
	}
	b, err := session.GenerateID()
	if err != nil {
		return "", err
	}
	id := "bench-" + a + b
	if _, err := session.Open(id); err != nil {
		return "", err
	}
	st, err := store.Open(project.FindRoot(dir))
	if err != nil {
		return "", err
	}
	defer st.Close()
	if err := st.CreateBenchmarkSession(id); err != nil {
		return "", err
	}
	return id, nil
}

// MeasureSessionCommand executes args once through a benchmark session and
// persists the resulting measurement for cross-process aggregation.
func MeasureSessionCommand(ctx context.Context, args []string, dir, sessionID string) (*Measurement, error) {
	if sessionID == "" {
		return nil, errors.New("benchmark session ID is required")
	}
	root := project.FindRoot(dir)
	st, err := store.Open(root)
	if err != nil {
		return nil, err
	}
	exists, err := st.BenchmarkSessionExists(sessionID)
	_ = st.Close()
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("benchmark session %q not found", sessionID)
	}

	m, err := measure(ctx, engine.Run, args, dir, sessionID)
	if err != nil {
		return nil, err
	}
	st, err = store.Open(root)
	if err != nil {
		return m, err
	}
	defer st.Close()
	err = st.SaveBenchmarkRecord(&store.BenchmarkRecord{
		SessionID: sessionID, ResultID: m.ResultID, Command: m.Command, ExitCode: m.ExitCode,
		RawStdoutBytes: m.RawStdoutBytes, RawStderrBytes: m.RawStderrBytes, RawBytes: m.RawBytes,
		StatelessBytes: m.StatelessVisibleBytes, StatefulBytes: m.StatefulVisibleBytes,
		Presentation: m.Presentation, Truncated: m.Truncated,
		ExecutionDuration: m.ExecutionDuration, ReduceDuration: m.ReduceDuration,
		ProcessingDuration: m.ProcessingDuration,
	})
	return m, err
}

// runFunc matches engine.Run; it is a parameter only so tests can observe
// how many times the pipeline is invoked.
type runFunc func(context.Context, *protocol.ToolRequest) (*engine.Outcome, error)

func measure(ctx context.Context, run runFunc, args []string, dir, sessionID string) (*Measurement, error) {
	if len(args) == 0 {
		return nil, errors.New("no command specified")
	}
	req := &protocol.ToolRequest{
		Protocol:   protocol.Version,
		Command:    args,
		WorkingDir: dir,
		SessionID:  sessionID,
	}
	out, err := run(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := out.Response
	if resp.Error != "" || out.Exec == nil {
		msg := resp.Error
		if msg == "" {
			msg = "command did not execute"
		}
		return nil, &ExecError{ExitCode: resp.ExitCode, Msg: msg}
	}

	r := out.Exec
	m := &Measurement{
		Command:               args,
		ExitCode:              resp.ExitCode,
		RawStdoutBytes:        int64(len(r.Stdout)),
		RawStderrBytes:        int64(len(r.Stderr)),
		Truncated:             r.Truncated,
		AgentVisibleBytes:     int64(len(resp.Stdout)),
		StatelessVisibleBytes: int64(len(out.StatelessPresentation)),
		StatefulVisibleBytes:  int64(len(resp.Stdout)),
		ExecutionDuration:     r.Duration,
		ReduceDuration:        out.ReduceDuration,
		ProcessingDuration:    out.ProcessingDuration,
		ResultID:              resp.ResultID,
		Presentation:          resp.Presentation,
	}
	m.RawBytes = m.RawStdoutBytes + m.RawStderrBytes
	return m, nil
}

// SessionMeasurement is a reusable aggregate of persisted command measurements.
type SessionMeasurement struct {
	SessionID string
	Commands  []Measurement

	CommandCount       int
	RawBytes           int64
	StatelessBytes     int64
	StatefulBytes      int64
	FullCount          int
	DeltaCount         int
	UnchangedCount     int
	ExecutionDuration  time.Duration
	ReduceDuration     time.Duration
	ProcessingDuration time.Duration
}

// LoadSession reconstructs a benchmark aggregate from its individual records.
func LoadSession(dir, sessionID string) (*SessionMeasurement, error) {
	st, err := store.Open(project.FindRoot(dir))
	if err != nil {
		return nil, err
	}
	defer st.Close()
	records, err := st.LoadBenchmarkRecords(sessionID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("benchmark session %q not found", sessionID)
	}
	if err != nil {
		return nil, err
	}
	agg := &SessionMeasurement{SessionID: sessionID}
	for _, rec := range records {
		m := Measurement{
			Command: rec.Command, ExitCode: rec.ExitCode, ResultID: rec.ResultID,
			RawStdoutBytes: rec.RawStdoutBytes, RawStderrBytes: rec.RawStderrBytes, RawBytes: rec.RawBytes,
			Truncated: rec.Truncated, AgentVisibleBytes: rec.StatefulBytes,
			StatelessVisibleBytes: rec.StatelessBytes, StatefulVisibleBytes: rec.StatefulBytes,
			Presentation: rec.Presentation, ExecutionDuration: rec.ExecutionDuration,
			ReduceDuration: rec.ReduceDuration, ProcessingDuration: rec.ProcessingDuration,
		}
		agg.Commands = append(agg.Commands, m)
		agg.RawBytes += m.RawBytes
		agg.StatelessBytes += m.StatelessVisibleBytes
		agg.StatefulBytes += m.StatefulVisibleBytes
		agg.ExecutionDuration += m.ExecutionDuration
		agg.ReduceDuration += m.ReduceDuration
		agg.ProcessingDuration += m.ProcessingDuration
		switch m.Presentation {
		case "unchanged":
			agg.UnchangedCount++
		case "delta":
			agg.DeltaCount++
		default:
			agg.FullCount++
		}
	}
	agg.CommandCount = len(agg.Commands)
	return agg, nil
}

// ExecError reports that the target command could not be executed
// (e.g. not found). ExitCode follows the `acap run` convention (127 for
// not found).
type ExecError struct {
	ExitCode int
	Msg      string
}

func (e *ExecError) Error() string { return e.Msg }
