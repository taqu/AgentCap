package agentbench

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/taqu/agentcap/internal/integration/antigravity"
	"github.com/taqu/agentcap/internal/integration/claude"
	"github.com/taqu/agentcap/internal/workload"
)

// Claude Code and Antigravity production adapters always map the agent session,
// so they expose only OFF (disabled, via the common bypass) and FULL
// (integrated). Stateless/stateful are not faked for them (§5, §6).
func productionModes(m Mode) bool { return m == ModeDisabled || m == ModeIntegrated }

// runHost starts an agent process in the workspace and classifies its exit.
func runHost(ctx context.Context, executable string, args []string, req AgentRunRequest, extraEnv ...string) (*AgentRunResult, error) {
	if _, err := exec.LookPath(executable); err != nil {
		return nil, fmt.Errorf("%s executable: %w", executable, err)
	}
	cmd := exec.CommandContext(ctx, executable, args...)
	cmd.Dir = req.Workspace
	env := benchmarkEnv(os.Environ(), req)
	if req.BinDir != "" {
		for i, item := range env {
			if strings.HasPrefix(item, "PATH=") {
				env[i] = "PATH=" + req.BinDir + string(os.PathListSeparator) + item[len("PATH="):]
			}
		}
	}
	cmd.Env = append(env, extraEnv...)
	isolateGroup(cmd)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start %s: %w", executable, err)
	}
	waitErr := cmd.Wait()
	result := &AgentRunResult{Stdout: stdout.Bytes(), Stderr: stderr.Bytes()}
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
	return nil, fmt.Errorf("wait for %s: %w", executable, waitErr)
}

// hostFailure recognizes service-side failures in a host's final error text.
var hostFailure = regexp.MustCompile(`(?i)quota|rate.?limit|overloaded|RESOURCE_EXHAUSTED|\b(429|500|502|503|529)\b|API Error|authentication`)

func finishTrace(r *AgentRunResult) {
	for _, e := range r.Trace {
		if e.Tool == "shell" {
			r.CommandCount++
			r.CommandVisibleBytes += e.Bytes
		}
	}
}

// originalCommand recovers the agent's command from a rewritten
// "acap <agent>-exec <payload>" command line, when a host reports the
// post-hook arguments.
var wrapperCommand = regexp.MustCompile(`(?:claude|antigravity)-exec '?([A-Za-z0-9_-]+)'?\s*$`)

func originalCommand(command string) string {
	m := wrapperCommand.FindStringSubmatch(command)
	if m == nil {
		return command
	}
	data, err := base64.RawURLEncoding.DecodeString(m[1])
	if err != nil {
		return command
	}
	var p struct {
		Request struct {
			ShellCommand *string `json:"shell_command"`
		} `json:"request"`
	}
	if json.Unmarshal(data, &p) != nil || p.Request.ShellCommand == nil {
		// Fall back to a generic search for the first string field.
		var generic map[string]any
		if json.Unmarshal(data, &generic) == nil {
			if s := findShellCommand(generic); s != "" {
				return s
			}
		}
		return command
	}
	return *p.Request.ShellCommand
}

func findShellCommand(v any) string {
	switch t := v.(type) {
	case map[string]any:
		for k, val := range t {
			if strings.EqualFold(k, "shellcommand") || strings.EqualFold(k, "shell_command") {
				if s, ok := val.(string); ok {
					return s
				}
			}
			if s := findShellCommand(val); s != "" {
				return s
			}
		}
	}
	return ""
}

// ---------------------------------------------------------------- Claude Code

// ClaudeAdapter drives `claude -p` in stream-json mode with the production
// Bash PreToolUse hook installed in the workspace in both modes.
type ClaudeAdapter struct {
	Executable string
	Model      string
}

func (a *ClaudeAdapter) Name() string             { return "claude" }
func (a *ClaudeAdapter) SupportsMode(m Mode) bool { return productionModes(m) }

func (a *ClaudeAdapter) Run(ctx context.Context, req AgentRunRequest) (*AgentRunResult, error) {
	if !productionModes(req.Mode) {
		return nil, fmt.Errorf("claude adapter supports modes disabled and integrated only")
	}
	// Hook stays installed for OFF; ACAP_BENCH_MODE=disabled is the bypass (§10).
	if err := claude.Install(req.Workspace, false); err != nil {
		return nil, err
	}
	model := req.Model
	if model == "" {
		model = a.Model
	}
	// The adapter only intercepts in bypassPermissions; use it in both modes so
	// only the AgentCap mode differs. Project settings only: user-level
	// settings must not vary the benchmark configuration.
	args := []string{"-p", req.Task, "--output-format", "stream-json", "--verbose",
		"--permission-mode", "bypassPermissions", "--setting-sources", "project", "--strict-mcp-config"}
	if model != "" {
		args = append(args, "--model", model)
	}
	exe := a.Executable
	if exe == "" {
		exe = "claude"
	}
	result, err := runHost(ctx, exe, args, req)
	if err != nil {
		return nil, err
	}
	parseClaudeStream(result)
	finishTrace(result)
	return result, nil
}

func parseClaudeStream(r *AgentRunResult) {
	type use struct{ name, command string }
	uses := map[string]use{}
	scanner := bufio.NewScanner(bytes.NewReader(r.Stdout))
	scanner.Buffer(make([]byte, 64*1024), 64*1024*1024)
	for scanner.Scan() {
		var ev struct {
			Type    string `json:"type"`
			Subtype string `json:"subtype"`
			Version string `json:"claude_code_version"`
			IsError bool   `json:"is_error"`
			Result  string `json:"result"`
			Message struct {
				Content json.RawMessage `json:"content"`
			} `json:"message"`
			Usage struct {
				Input      int64 `json:"input_tokens"`
				CacheWrite int64 `json:"cache_creation_input_tokens"`
				CacheRead  int64 `json:"cache_read_input_tokens"`
				Output     int64 `json:"output_tokens"`
			} `json:"usage"`
		}
		if json.Unmarshal(scanner.Bytes(), &ev) != nil {
			continue
		}
		switch ev.Type {
		case "system":
			if ev.Version != "" {
				r.AgentVersion = ev.Version
			}
		case "result":
			r.HostInput, r.HostCacheRead, r.HostCacheWrite, r.HostOutput = ev.Usage.Input, ev.Usage.CacheRead, ev.Usage.CacheWrite, ev.Usage.Output
			if ev.IsError && hostFailure.MatchString(ev.Result) {
				r.HostError = truncate(ev.Result, 200)
			}
		case "assistant", "user":
			var blocks []struct {
				Type      string          `json:"type"`
				ID        string          `json:"id"`
				Name      string          `json:"name"`
				Input     json.RawMessage `json:"input"`
				ToolUseID string          `json:"tool_use_id"`
				Content   json.RawMessage `json:"content"`
				IsError   bool            `json:"is_error"`
			}
			if json.Unmarshal(ev.Message.Content, &blocks) != nil {
				continue
			}
			for _, b := range blocks {
				switch b.Type {
				case "tool_use":
					var in struct {
						Command string `json:"command"`
					}
					json.Unmarshal(b.Input, &in)
					uses[b.ID] = use{b.Name, originalCommand(in.Command)}
				case "tool_result":
					u := uses[b.ToolUseID]
					text := contentText(b.Content)
					if u.name == "Bash" {
						r.Trace = append(r.Trace, shellEntry(u.command, text, b.IsError))
					} else {
						r.Trace = append(r.Trace, workload.TraceEntry{Tool: u.name, Family: "native", Bytes: int64(len(text)), IsError: b.IsError})
					}
				}
			}
		}
	}
}

// contentText flattens a tool_result content (string or text blocks).
func contentText(raw json.RawMessage) string {
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	var parts []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	if json.Unmarshal(raw, &parts) == nil {
		var sb strings.Builder
		for _, p := range parts {
			sb.WriteString(p.Text)
		}
		return sb.String()
	}
	return string(raw)
}

// ---------------------------------------------------------------- Antigravity

// AntigravityAdapter drives `agy -p` in stream-json mode with the production
// run_command PreToolUse hook installed in the workspace in both modes.
type AntigravityAdapter struct {
	Executable string
	Model      string
}

func (a *AntigravityAdapter) Name() string             { return "antigravity" }
func (a *AntigravityAdapter) SupportsMode(m Mode) bool { return productionModes(m) }

func (a *AntigravityAdapter) Run(ctx context.Context, req AgentRunRequest) (*AgentRunResult, error) {
	if !productionModes(req.Mode) {
		return nil, fmt.Errorf("antigravity adapter supports modes disabled and integrated only")
	}
	if err := antigravity.Install(req.Workspace, false); err != nil {
		return nil, err
	}
	model := req.Model
	if model == "" {
		model = a.Model
	}
	args := []string{"-p", req.Task, "--output-format", "stream-json", "--dangerously-skip-permissions"}
	if model != "" {
		args = append(args, "--model", model)
	}
	exe := a.Executable
	if exe == "" {
		exe = "agy"
	}
	// Declared unrestricted in both modes; OFF is bypassed by ACAP_BENCH_MODE.
	result, err := runHost(ctx, exe, args, req, "ACAP_ANTIGRAVITY_PERMISSIONS=unrestricted")
	if err != nil {
		return nil, err
	}
	if out, err := exec.Command(exe, "--version").Output(); err == nil {
		result.AgentVersion = strings.TrimSpace(string(out))
	}
	parseAntigravityStream(result)
	finishTrace(result)
	return result, nil
}

func parseAntigravityStream(r *AgentRunResult) {
	scanner := bufio.NewScanner(bytes.NewReader(r.Stdout))
	scanner.Buffer(make([]byte, 64*1024), 64*1024*1024)
	for scanner.Scan() {
		var ev struct {
			Event string `json:"event"`
			Step  struct {
				State    string `json:"state"`
				StepType string `json:"step_type"`
				ToolName string `json:"tool_name"`
				ToolInfo struct {
					Parameters map[string]any `json:"parameters"`
					Output     string         `json:"output"`
					Error      string         `json:"error"`
				} `json:"tool_info"`
			} `json:"step_update"`
			Result struct {
				Status string `json:"status"`
				Error  string `json:"error"`
				Usage  struct {
					Input     int64 `json:"input_tokens"`
					Output    int64 `json:"output_tokens"`
					Thinking  int64 `json:"thinking_tokens"`
					CacheRead int64 `json:"cache_read_tokens"`
				} `json:"usage"`
			} `json:"result"`
		}
		if json.Unmarshal(scanner.Bytes(), &ev) != nil {
			continue
		}
		switch ev.Event {
		case "result":
			if ev.Result.Status == "ERROR" && ev.Result.Error != "" {
				r.HostError = truncate(ev.Result.Error, 200)
			}
			u := ev.Result.Usage
			r.HostInput, r.HostCacheRead, r.HostOutput = u.Input, u.CacheRead, u.Output+u.Thinking
		case "step_update":
			s := ev.Step
			if s.StepType != "tool" || (s.State != "DONE" && s.State != "ERROR") {
				continue
			}
			text := strings.ReplaceAll(s.ToolInfo.Output+s.ToolInfo.Error, "\r\n", "\n")
			if s.ToolName == "run_command" {
				command, _ := s.ToolInfo.Parameters["CommandLine"].(string)
				r.Trace = append(r.Trace, shellEntry(originalCommand(command), text, s.State == "ERROR"))
			} else {
				r.Trace = append(r.Trace, workload.TraceEntry{Tool: s.ToolName, Family: "native", Bytes: int64(len(text)), IsError: s.State == "ERROR"})
			}
		}
	}
}
