package agentbench

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/taqu/agentcap/internal/stats"
	"github.com/taqu/agentcap/internal/workload"
)

// linkAcapBinary exposes the running acap binary as "acap" on a private PATH
// directory so production hooks ("acap hook <agent>") use the benchmarked build.
func linkAcapBinary(tempRoot string) (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("locate acap binary: %w", err)
	}
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		exe = resolved
	}
	dir, err := os.MkdirTemp(tempRoot, "acap-bin-")
	if err != nil {
		return "", err
	}
	name := "acap"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	if err := os.Symlink(exe, filepath.Join(dir, name)); err != nil {
		os.RemoveAll(dir)
		return "", fmt.Errorf("link acap binary: %w", err)
	}
	return dir, nil
}

// workflowMetrics assembles the Phase 7 record and flags implausible data (§96).
func workflowMetrics(def *workload.Definition, opts Options, r *AgentRunResult, delta *stats.Stats, root, runID string, started time.Time, success bool, status string) *workload.WorkflowMetrics {
	w := buildWorkflow(r.Trace, r.InstructionBytes)
	w.Trace = r.Trace
	w.RunID = runID
	w.Batch = opts.Batch
	w.StartedAt = started.UTC().Format(time.RFC3339)
	w.Platform = runtime.GOOS + "/" + runtime.GOARCH
	w.Category, w.Language = def.Category, def.Language
	w.AgentVer = r.AgentVersion
	w.Model = opts.Model
	w.AcapCommit = opts.AcapVersion
	w.HostInputTokens, w.HostCacheReadTokens, w.HostCacheWriteTokens, w.HostOutputTokens = r.HostInput, r.HostCacheRead, r.HostCacheWrite, r.HostOutput
	// Read in OFF too: explicit-bypass records verify the OFF switch (§10).
	applyAdapterMetrics(root, w)
	switch {
	case r.HostError != "":
		w.Outcome = "INFRA_ERROR"
		w.Anomalies = append(w.Anomalies, "host error: "+r.HostError)
	case status == "timeout":
		w.Outcome = "TIMEOUT"
	case status == "agent_error" && len(r.Trace) == 0:
		w.Outcome = "INFRA_ERROR" // agent never reached a tool call (API/startup failure)
	case success:
		w.Outcome = "PASS"
	default:
		w.Outcome = "FAIL"
	}
	// Sanity checks: flag, never silently fix.
	if w.ShellCommands+w.NativeToolCalls != len(r.Trace) {
		w.Anomalies = append(w.Anomalies, "trace length mismatch")
	}
	if opts.Mode == ModeDisabled && w.AcapResults > 0 {
		w.Anomalies = append(w.Anomalies, "ADAPTER: AgentCap capsules observed in OFF mode")
	}
	if opts.Mode.AgentCapEnabled() {
		nonRetrieval := w.ShellCommands - w.ShowCalls - w.RawCalls
		if nonRetrieval > 0 && w.Intercepted+w.Bypassed == 0 {
			w.Anomalies = append(w.Anomalies, "ADAPTER: hook never ran")
		}
		if int64(w.AcapResults) > delta.Commands {
			w.Anomalies = append(w.Anomalies, "capsules in trace exceed stored results")
		}
		if int64(w.ShowCalls) != delta.ShowCalls || int64(w.RawCalls) != delta.RawCalls {
			w.Anomalies = append(w.Anomalies, fmt.Sprintf("show/raw trace counts (%d/%d) differ from store (%d/%d)", w.ShowCalls, w.RawCalls, delta.ShowCalls, delta.RawCalls))
		}
		if w.AdapterFailures > 0 {
			w.Anomalies = append(w.Anomalies, "ADAPTER: adapter failures recorded")
		}
	}
	return w
}

// writeEvidence retains the raw material needed to review one run (§58).
func writeEvidence(ctx context.Context, dir, runID, workspace string, r *AgentRunResult, verifyLog []byte, result *workload.BenchmarkResult) error {
	runDir := filepath.Join(dir, runID)
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		return err
	}
	diff, _ := exec.CommandContext(ctx, "git", "-C", workspace, "diff", "HEAD").Output()
	status, _ := exec.CommandContext(ctx, "git", "-C", workspace, "status", "--porcelain").Output()
	data, _ := json.MarshalIndent(result, "", "  ")
	files := map[string][]byte{
		"agent.stdout.jsonl": r.Stdout, "agent.stderr.txt": r.Stderr,
		"final.diff": diff, "final.status": status, "verify.txt": verifyLog, "result.json": data,
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(runDir, name), content, 0o644); err != nil {
			return err
		}
	}
	return nil
}
