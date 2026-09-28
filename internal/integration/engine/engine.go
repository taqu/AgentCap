// Package engine provides the core AgentCap pipeline as a reusable function
// for coding-agent integration.
package engine

import (
	"context"
	"strings"
	"time"

	"github.com/taqu/agentcap/internal/delta"
	"github.com/taqu/agentcap/internal/exec"
	"github.com/taqu/agentcap/internal/integration/protocol"
	"github.com/taqu/agentcap/internal/project"
	"github.com/taqu/agentcap/internal/reduce"
	"github.com/taqu/agentcap/internal/session"
	"github.com/taqu/agentcap/internal/store"
)

// Execute runs a command through the full AgentCap pipeline and returns the
// compressed tool response. It is fail-open: if storage or session logic fails,
// it still returns the command's output.
func Execute(ctx context.Context, req *protocol.ToolRequest) (*protocol.ToolResponse, error) {
	out, err := Run(ctx, req)
	if out == nil {
		return nil, err
	}
	return out.Response, err
}

// Outcome is the result of one pass through the pipeline, together with the
// captured execution it was derived from and timings for observers such as
// the benchmark layer. Response is exactly what Execute returns.
type Outcome struct {
	Response *protocol.ToolResponse

	// StatelessPresentation is the fully rendered AgentCap presentation before
	// session comparison. It includes the same result-ID injection as Response
	// and is derived from the same captured execution. For a stateless request
	// it is identical to Response.Stdout.
	StatelessPresentation string

	// Exec is the captured child result the response was derived from.
	// Nil when the command could not be started or no command was given.
	Exec *exec.Result

	// ReduceDuration is the wall time spent in reducer.Reduce only.
	ReduceDuration time.Duration

	// ProcessingDuration is the wall time from child exit (exec.Run
	// returning) until Response is fully built. It covers reduction,
	// hashing, result/object storage, metadata updates, delta comparison,
	// result-ID injection, session history and stats recording.
	ProcessingDuration time.Duration
}

// Run is Execute, but also returns the captured execution and timings.
// The command is executed exactly once.
func Run(ctx context.Context, req *protocol.ToolRequest) (*Outcome, error) {
	if len(req.Command) == 0 {
		return &Outcome{Response: &protocol.ToolResponse{
			Protocol: protocol.Version,
			ExitCode: 1,
			Error:    "no command specified",
		}}, nil
	}

	startTime := time.Now().UTC()

	// Execute the command.
	opts := &exec.Options{}
	if req.WorkingDir != "" {
		opts.Dir = req.WorkingDir
	}
	result, execErr := exec.Run(ctx, req.Command, opts)
	if execErr != nil {
		// Command not found or permission denied — return error response.
		exitCode := 1
		if result != nil {
			exitCode = result.ExitCode
		}
		stdout := ""
		stderr := ""
		if result != nil {
			stdout = string(result.Stdout)
			stderr = string(result.Stderr)
		}
		return &Outcome{Response: &protocol.ToolResponse{
			Protocol: protocol.Version,
			ExitCode: exitCode,
			Stdout:   stdout,
			Stderr:   stderr,
			Error:    execErr.Error(),
		}}, nil
	}

	processStart := time.Now()

	// Reduce output.
	reducer := reduce.Select(req.Command)
	reduceStart := time.Now()
	reduced := func() (r *reduce.ReducedResult) {
		defer func() {
			if rec := recover(); rec != nil {
				// Fail-open: reducer panicked; return raw output.
				r = &reduce.ReducedResult{
					Output:   string(result.Stdout),
					RawBytes: len(result.Stdout),
					RetBytes: len(result.Stdout),
				}
			}
		}()
		return reducer.Reduce(result)
	}()
	reduceDuration := time.Since(reduceStart)

	// Compute hashes.
	stdoutHash := delta.HashBytes(result.Stdout)
	stderrHash := delta.HashBytes(result.Stderr)

	// Determine working directory.
	workDir := req.WorkingDir

	// Build metadata.
	reducerName := reducerTypeName(reducer)
	meta := store.Meta{
		Command:     req.Command,
		ExitCode:    result.ExitCode,
		StartedAt:   startTime,
		DurationMs:  result.Duration.Milliseconds(),
		StdoutBytes: int64(len(result.Stdout)),
		StderrBytes: int64(len(result.Stderr)),
		Reducer:     reducerName,
		CreatedAt:   time.Now().UTC(),
		Truncated:   result.Truncated,
		StdoutHash:  stdoutHash,
		StderrHash:  stderrHash,
		WorkDir:     workDir,
	}

	// Open session.
	var sess *session.Session
	if req.SessionID != "" {
		if s, err := session.Open(req.SessionID); err == nil {
			sess = s
		}
	}

	var baselineEntry *store.Entry
	var baselineRec *session.HistoryRecord

	if sess != nil {
		meta.SessionID = sess.ID
		key := session.NewKey(req.Command, workDir)
		if rec, err := sess.LatestBaseline(key); err == nil && rec != nil {
			baselineRec = rec
		}
	}

	// Open store — fail-open if it fails.
	root := project.FindRoot(workDir)
	st, storeInitErr := store.Open(root)
	if storeInitErr == nil {
		defer st.Close()
	}

	var entry *store.Entry
	if storeInitErr == nil && reduced != nil {
		var saveErr error
		entry, saveErr = st.Save(meta, result.Stdout, result.Stderr, reduced.Output)
		if saveErr != nil {
			entry = nil
		}
	}

	// Load baseline entry if we have a reference.
	if sess != nil && baselineRec != nil && st != nil {
		if be, err := st.Open(baselineRec.ResultID); err == nil {
			baselineEntry = be
		}
	}

	// Save git diff metadata if applicable.
	if entry != nil && storeInitErr == nil {
		if gdr, ok := reducer.(*reduce.GitDiffReducer); ok && gdr.ParsedDiff != nil {
			_ = st.SaveGitDiff(entry.ID, gdr.ParsedDiff.Files)
		}
	}

	// Save build diagnostics / test failures if applicable.
	if entry != nil && st != nil {
		switch rv := reducer.(type) {
		case *reduce.GoBuildReducer:
			if rv.ParsedBuild != nil && len(rv.ParsedBuild.Diagnostics) > 0 {
				_ = st.SaveDiagnostics(entry.ID, rv.ParsedBuild.Diagnostics)
			}
		case *reduce.GoTestReducer:
			if rv.ParsedRun != nil && len(rv.ParsedRun.Failures) > 0 {
				_ = st.SaveTestFailures(entry.ID, rv.ParsedRun.Failures)
			}
		case *reduce.GccReducer:
			if rv.ParsedBuild != nil && len(rv.ParsedBuild.Diagnostics) > 0 {
				_ = st.SaveDiagnostics(entry.ID, rv.ParsedBuild.Diagnostics)
			}
		case *reduce.CargoBuildReducer:
			if rv.ParsedBuild != nil && len(rv.ParsedBuild.Diagnostics) > 0 {
				_ = st.SaveDiagnostics(entry.ID, rv.ParsedBuild.Diagnostics)
			}
		case *reduce.CargoTestReducer:
			if rv.ParsedRun != nil && len(rv.ParsedRun.Failures) > 0 {
				_ = st.SaveTestFailures(entry.ID, rv.ParsedRun.Failures)
			}
		}
	}

	// Determine presentation. Keep the full reduced form so observers can
	// measure the real pre-delta presentation without executing or reducing the
	// command a second time.
	output := ""
	if reduced != nil {
		output = reduced.Output
	}
	statelessOutput := output
	presentation := string(delta.PresentationFull)

	if sess != nil && baselineEntry != nil && entry != nil && reduced != nil {
		dr := selectDeltaReducer(reducer, st)
		deltaResult := delta.Compare(
			ctx,
			baselineEntry,
			entry,
			reduced.Output,
			dr,
			stdoutHash,
			stderrHash,
		)
		output = deltaResult.Output
		presentation = string(deltaResult.Presentation)

		updatedMeta := entry.Meta
		updatedMeta.BaselineID = baselineEntry.Meta.ID
		updatedMeta.Presentation = presentation
		_ = st.UpdateMeta(entry, updatedMeta)
	} else if entry != nil {
		updatedMeta := entry.Meta
		updatedMeta.Presentation = "full"
		_ = st.UpdateMeta(entry, updatedMeta)
	}

	// Inject result ID.
	resultID := ""
	if entry != nil {
		resultID = entry.ID
		output = injectResultID(output, entry.ID)
		statelessOutput = injectResultID(statelessOutput, entry.ID)
	}

	// Record in session history.
	if sess != nil && entry != nil {
		key := session.NewKey(req.Command, workDir)
		_ = sess.Record(entry.ID, key, result.ExitCode, stdoutHash, stderrHash, presentation)
	}

	// Record stats.
	if st != nil && reduced != nil {
		_ = st.RecordRun(reduced.RawBytes, reduced.RetBytes, len(output), presentation)
	}

	resp := &protocol.ToolResponse{
		Protocol:     protocol.Version,
		ResultID:     resultID,
		ExitCode:     result.ExitCode,
		Stdout:       output,
		Stderr:       string(result.Stderr),
		Presentation: presentation,
	}
	return &Outcome{
		Response:              resp,
		StatelessPresentation: statelessOutput,
		Exec:                  result,
		ReduceDuration:        reduceDuration,
		ProcessingDuration:    time.Since(processStart),
	}, nil
}

// injectResultID inserts the result ID into the output header.
func injectResultID(output, id string) string {
	if strings.HasPrefix(output, "@acap ") {
		rest := output[len("@acap "):]
		return "@acap " + id + " " + rest
	}
	return "@acap " + id + "\n" + output
}

// reducerTypeName returns a short string identifier for the reducer.
func reducerTypeName(r reduce.Reducer) string {
	switch rv := r.(type) {
	case *reduce.LsReducer:
		return "ls"
	case *reduce.FindReducer:
		return "find"
	case *reduce.GrepReducer:
		return "grep"
	case *reduce.CatReducer:
		return "cat"
	case *reduce.HeadTailReducer:
		return "headtail"
	case *reduce.TreeReducer:
		return "tree"
	case *reduce.DuReducer:
		return "du"
	case *reduce.WcReducer:
		return "wc"
	case *reduce.GitStatusReducer:
		return "git-status"
	case *reduce.GitDiffReducer:
		return "git-diff"
	case *reduce.GitShowReducer:
		return "git-show"
	case *reduce.GitLogReducer:
		return "git-log"
	case *reduce.GitBranchReducer:
		return "git-branch"
	case *reduce.GoBuildReducer:
		return "go-build"
	case *reduce.GoTestReducer:
		return "go-test"
	case *reduce.GccReducer:
		if rv.Tool != "" {
			return rv.Tool
		}
		return "gcc"
	case *reduce.CargoBuildReducer:
		return "cargo-build"
	case *reduce.CargoTestReducer:
		return "cargo-test"
	case *reduce.MakeReducer:
		return "make"
	case *reduce.NinjaReducer:
		return "ninja"
	case *reduce.GenericReducer:
		return "generic"
	default:
		return "unknown"
	}
}

// selectDeltaReducer maps a reduce.Reducer to its corresponding delta.Reducer.
func selectDeltaReducer(r reduce.Reducer, st *store.Store) delta.Reducer {
	switch r.(type) {
	case *reduce.GrepReducer:
		return &delta.GrepDelta{}
	case *reduce.FindReducer:
		return &delta.FindDelta{}
	case *reduce.LsReducer:
		return &delta.LsDelta{}
	case *reduce.GitStatusReducer:
		return &delta.GitStatusDelta{}
	case *reduce.GitDiffReducer:
		return &delta.GitDiffDelta{Store: st}
	case *reduce.GoBuildReducer:
		return &delta.DiagnosticsDelta{Store: st}
	case *reduce.GoTestReducer:
		return &delta.TestResultsDelta{Store: st}
	case *reduce.GccReducer:
		return &delta.DiagnosticsDelta{Store: st}
	case *reduce.CargoBuildReducer:
		return &delta.DiagnosticsDelta{Store: st}
	case *reduce.CargoTestReducer:
		return &delta.TestResultsDelta{Store: st}
	default:
		return nil
	}
}
