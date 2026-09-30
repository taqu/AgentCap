package agentbench

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	acapexec "github.com/taqu/agentcap/internal/exec"
	"github.com/taqu/agentcap/internal/session"
	"github.com/taqu/agentcap/internal/stats"
	"github.com/taqu/agentcap/internal/store"
	"github.com/taqu/agentcap/internal/workload"
)

type Options struct {
	Adapter       Adapter
	Mode          Mode
	StoreRoot     string
	TempRoot      string
	Timeout       time.Duration
	KeepWorkspace bool
	HookCommand   string
	Model         string

	// Phase 7 batch metadata and evidence retention.
	Batch       string
	AcapVersion string
	EvidenceDir string
}

type Trial struct {
	Benchmark         *workload.BenchmarkResult
	AgentStdout       []byte
	AgentStderr       []byte
	RetainedWorkspace string
	RetainedStateRoot string
}

// Run executes one agent trial in a fresh fixture workspace and verifies it
// outside the agent-visible measurement boundary.
func Run(ctx context.Context, definition *workload.Definition, opts Options) (trial *Trial, err error) {
	if opts.Adapter == nil {
		return nil, errors.New("agent adapter is required")
	}
	if definition.Task == "" || len(definition.Verify) == 0 {
		return nil, errors.New("coding-agent workload with verification is required")
	}
	if opts.StoreRoot == "" {
		return nil, errors.New("persistent store root is required")
	}
	if ms, ok := opts.Adapter.(ModeSupporter); ok && !ms.SupportsMode(opts.Mode) {
		return nil, fmt.Errorf("agent %s does not support mode %s through its production adapter", opts.Adapter.Name(), opts.Mode)
	}
	binDir, err := linkAcapBinary(opts.TempRoot)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(binDir)
	workspace, err := workload.PrepareWorkspace(ctx, definition, opts.TempRoot)
	if err != nil {
		return nil, err
	}
	trial = &Trial{}
	defer func() {
		if opts.KeepWorkspace {
			trial.RetainedWorkspace = workspace
			return
		}
		if removeErr := os.RemoveAll(workspace); err == nil && removeErr != nil {
			err = fmt.Errorf("clean temporary workspace: %w", removeErr)
		}
	}()

	// StoreRoot is explicit benchmark state. Do not apply normal project-root
	// discovery here because ACAP_ROOT from the caller could defeat per-trial
	// isolation established by RunRepeated.
	root := filepath.Clean(opts.StoreRoot)
	sessionID, err := session.GenerateID()
	if err != nil {
		return trial, fmt.Errorf("create benchmark session scope: %w", err)
	}
	// Including the per-trial store root makes repeated-run scopes distinct even
	// in the extraordinarily unlikely event that two random session IDs collide.
	sessionScope := root + ":" + sessionID
	before, err := loadStats(root)
	if err != nil {
		return trial, err
	}
	timeout := opts.Timeout
	if timeout <= 0 && definition.Timeout != "" {
		timeout, _ = time.ParseDuration(definition.Timeout)
	}
	if timeout <= 0 {
		timeout = 10 * time.Minute
	}
	agentCtx, cancel := context.WithTimeout(ctx, timeout)
	started := time.Now()
	runID := fmt.Sprintf("%s-%d", started.UTC().Format("20060102T150405.000"), os.Getpid())
	agentResult, runErr := opts.Adapter.Run(agentCtx, AgentRunRequest{
		Workspace: workspace, Task: definition.Task, Mode: opts.Mode,
		StoreRoot: root, HookCommand: opts.HookCommand, Model: opts.Model,
		SessionScope: sessionScope, BinDir: binDir,
	})
	wall := time.Since(started)
	cancel()
	if runErr != nil {
		return trial, runErr
	}
	trial.AgentStdout = agentResult.Stdout
	trial.AgentStderr = agentResult.Stderr
	after, err := loadStats(root)
	if err != nil {
		return trial, err
	}
	delta := subtractStats(after, before)
	success, verifyLog, err := verify(ctx, workspace, definition.Verify)
	if err != nil {
		return trial, err
	}
	status := "completed"
	if agentResult.TimedOut {
		status = "timeout"
	} else if agentResult.Canceled {
		status = "canceled"
	} else if agentResult.ExitCode != 0 {
		status = "agent_error"
	}
	exitCode := agentResult.ExitCode
	benchmark := &workload.BenchmarkResult{
		SchemaVersion: workload.BenchmarkResultSchemaVersion,
		Workload:      definition.Name, Agent: opts.Adapter.Name(), Mode: string(opts.Mode),
		TaskSuccess: &success, ExecutionStatus: status, AgentExitCode: &exitCode,
		WallTimeNS: wall.Nanoseconds(),
	}
	if opts.Mode == ModeDisabled {
		benchmark.Commands = agentResult.CommandCount
		benchmark.RawBytes = agentResult.CommandVisibleBytes
		benchmark.StatelessBytes = agentResult.CommandVisibleBytes
		benchmark.StatefulBytes = agentResult.CommandVisibleBytes
		benchmark.InitialVisibleBytes = agentResult.CommandVisibleBytes
		benchmark.TotalVisibleBytes = agentResult.CommandVisibleBytes
	} else {
		benchmark.Commands = int(delta.Commands)
		benchmark.RawBytes = delta.RawBytes
		benchmark.StatelessBytes = delta.StatelessBytes
		benchmark.StatefulBytes = delta.StatefulBytes
		benchmark.InitialVisibleBytes = delta.StatefulBytes
		benchmark.ShowBytes = delta.ShowRetBytes
		benchmark.RawRetrievalBytes = delta.RawRetBytes
		benchmark.TotalVisibleBytes = benchmark.InitialVisibleBytes + benchmark.ShowBytes + benchmark.RawRetrievalBytes
		benchmark.ShowCount = int(delta.ShowCalls)
		benchmark.RawRetrievalCount = int(delta.RawCalls)
		benchmark.ProcessingNS = delta.ProcessingNS
	}
	benchmark.Workflow = workflowMetrics(definition, opts, agentResult, &delta, root, runID, started, success, status)
	trial.Benchmark = benchmark
	if opts.EvidenceDir != "" {
		if err := writeEvidence(ctx, opts.EvidenceDir, runID, workspace, agentResult, verifyLog, benchmark); err != nil {
			return trial, err
		}
	}
	return trial, nil
}

func verify(ctx context.Context, workspace string, steps []workload.VerifyStep) (bool, []byte, error) {
	var log []byte
	for i, step := range steps {
		cwd, err := workload.ResolveWorkspaceDir(workspace, step.Run.Cwd)
		if err != nil {
			return false, log, fmt.Errorf("verify step %d cwd: %w", i+1, err)
		}
		result, err := acapexec.Run(ctx, step.Run.Argv, &acapexec.Options{Dir: cwd})
		if err != nil {
			return false, log, fmt.Errorf("verify step %d: %w", i+1, err)
		}
		expected := 0
		if step.Run.Expect != nil && step.Run.Expect.Exit != nil {
			expected = *step.Run.Expect.Exit
		}
		log = fmt.Appendf(log, "$ %v (exit %d, expected %d)\n%s%s\n", step.Run.Argv, result.ExitCode, expected, result.Stdout, result.Stderr)
		if result.ExitCode != expected {
			return false, log, nil
		}
	}
	return true, log, nil
}

func loadStats(root string) (*stats.Stats, error) {
	st, err := store.Open(root)
	if err != nil {
		return nil, err
	}
	defer st.Close()
	return st.LoadStats()
}

func subtractStats(after, before *stats.Stats) stats.Stats {
	return stats.Stats{
		Commands:       after.Commands - before.Commands,
		RawBytes:       after.RawBytes - before.RawBytes,
		ShowCalls:      after.ShowCalls - before.ShowCalls,
		ShowRetBytes:   after.ShowRetBytes - before.ShowRetBytes,
		RawCalls:       after.RawCalls - before.RawCalls,
		RawRetBytes:    after.RawRetBytes - before.RawRetBytes,
		StatelessBytes: after.StatelessBytes - before.StatelessBytes,
		StatefulBytes:  after.StatefulBytes - before.StatefulBytes,
		ProcessingNS:   after.ProcessingNS - before.ProcessingNS,
		UnchangedCount: after.UnchangedCount - before.UnchangedCount,
		DeltaCount:     after.DeltaCount - before.DeltaCount,
	}
}
