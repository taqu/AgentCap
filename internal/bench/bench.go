// Package bench measures AgentCap's behavior on real command executions.
//
// The benchmark layer is an observer: it runs a command through the normal
// AgentCap pipeline (engine.Run) exactly once and derives every measurement
// from that single captured execution. It never re-runs the command and never
// uses a benchmark-specific reducer.
package bench

import (
	"context"
	"errors"
	"time"

	"github.com/taqu/agentcap/internal/integration/engine"
	"github.com/taqu/agentcap/internal/integration/protocol"
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
	// Presentation is the pipeline presentation kind (always "full" here).
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
	return measure(ctx, engine.Run, args, dir)
}

// runFunc matches engine.Run; it is a parameter only so tests can observe
// how many times the pipeline is invoked.
type runFunc func(context.Context, *protocol.ToolRequest) (*engine.Outcome, error)

func measure(ctx context.Context, run runFunc, args []string, dir string) (*Measurement, error) {
	if len(args) == 0 {
		return nil, errors.New("no command specified")
	}
	req := &protocol.ToolRequest{
		Protocol:   protocol.Version,
		Command:    args,
		WorkingDir: dir,
		// SessionID intentionally empty: stateless measurement.
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
		Command:            args,
		ExitCode:           resp.ExitCode,
		RawStdoutBytes:     int64(len(r.Stdout)),
		RawStderrBytes:     int64(len(r.Stderr)),
		Truncated:          r.Truncated,
		AgentVisibleBytes:  int64(len(resp.Stdout)),
		ExecutionDuration:  r.Duration,
		ReduceDuration:     out.ReduceDuration,
		ProcessingDuration: out.ProcessingDuration,
		ResultID:           resp.ResultID,
		Presentation:       resp.Presentation,
	}
	m.RawBytes = m.RawStdoutBytes + m.RawStderrBytes
	return m, nil
}

// ExecError reports that the target command could not be executed
// (e.g. not found). ExitCode follows the `acap run` convention (127 for
// not found).
type ExecError struct {
	ExitCode int
	Msg      string
}

func (e *ExecError) Error() string { return e.Msg }
