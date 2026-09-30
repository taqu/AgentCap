package agentbench

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/taqu/agentcap/internal/integration/codex"
	"github.com/taqu/agentcap/internal/integration/protocol"
)

// CodexAdapter invokes `codex exec` in non-interactive JSONL mode.
type CodexAdapter struct {
	Executable string
	Model      string
}

func (a *CodexAdapter) Name() string { return "codex" }

func (a *CodexAdapter) Run(ctx context.Context, req AgentRunRequest) (*AgentRunResult, error) {
	executable := a.Executable
	if executable == "" {
		executable = "codex"
	}
	if _, err := exec.LookPath(executable); err != nil && !filepath.IsAbs(executable) {
		return nil, fmt.Errorf("codex executable: %w", err)
	}
	if req.Mode.AgentCapEnabled() {
		if err := codex.InstallHook(req.Workspace, req.HookCommand); err != nil {
			return nil, err
		}
		if req.Mode == ModeIntegrated {
			if err := codex.InstallInstructions(req.Workspace); err != nil {
				return nil, err
			}
		}
	}

	args := []string{"exec", "--json", "--ephemeral", "--ignore-user-config", "--ignore-rules", "--skip-git-repo-check", "--sandbox", "workspace-write", "--approve-for-me", "--cd", req.Workspace}
	if req.Mode.AgentCapEnabled() {
		args = append(args, "--dangerously-bypass-hook-trust")
	}
	model := req.Model
	if model == "" {
		model = a.Model
	}
	if model != "" {
		args = append(args, "--model", model)
	}
	args = append(args, req.Task)

	cmd := exec.CommandContext(ctx, executable, args...)
	cmd.Dir = req.Workspace
	cmd.Env = benchmarkEnv(os.Environ(), req)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start codex: %w", err)
	}
	waitErr := cmd.Wait()
	result := &AgentRunResult{Stdout: stdout.Bytes(), Stderr: stderr.Bytes()}
	result.CommandCount, result.CommandVisibleBytes = parseCodexEvents(result.Stdout)
	if waitErr == nil {
		return result, nil
	}
	if ctx.Err() != nil {
		result.TimedOut = errors.Is(ctx.Err(), context.DeadlineExceeded)
		result.Canceled = errors.Is(ctx.Err(), context.Canceled)
		result.ExitCode = -1
		return result, nil
	}
	var exitErr *exec.ExitError
	if errors.As(waitErr, &exitErr) {
		result.ExitCode = exitErr.ExitCode()
		return result, nil
	}
	return nil, fmt.Errorf("wait for codex: %w", waitErr)
}

func benchmarkEnv(env []string, req AgentRunRequest) []string {
	values := map[string]string{
		"ACAP_ROOT":               req.StoreRoot,
		protocol.EnvBenchmarkMode: string(req.Mode),
		protocol.EnvBypass:        "",
		protocol.EnvDepth:         "",
		"ACAP_SESSION_ID":         "",
	}
	out := make([]string, 0, len(env)+len(values))
	for _, item := range env {
		key, _, ok := strings.Cut(item, "=")
		if _, replace := values[key]; ok && replace {
			continue
		}
		out = append(out, item)
	}
	for key, value := range values {
		out = append(out, key+"="+value)
	}
	return out
}

func parseCodexEvents(data []byte) (commands int, visibleBytes int64) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		var event struct {
			Type string `json:"type"`
			Item struct {
				Type             string `json:"type"`
				AggregatedOutput string `json:"aggregated_output"`
				Output           string `json:"output"`
			} `json:"item"`
		}
		if json.Unmarshal(scanner.Bytes(), &event) != nil || event.Type != "item.completed" || event.Item.Type != "command_execution" {
			continue
		}
		commands++
		output := event.Item.AggregatedOutput
		if output == "" {
			output = event.Item.Output
		}
		visibleBytes += int64(len(output))
	}
	return commands, visibleBytes
}
