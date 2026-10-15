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

	"github.com/taqu/agentcap/internal/integration/claude"
)

// ClaudeCodeAdapter invokes the Claude Code CLI in non-interactive print mode
// with structured JSONL streaming output.
type ClaudeCodeAdapter struct {
	Executable string
	Model      string
}

func (a *ClaudeCodeAdapter) Name() string { return "claude" }

func (a *ClaudeCodeAdapter) Run(ctx context.Context, req AgentRunRequest) (*AgentRunResult, error) {
	executable := a.Executable
	if executable == "" {
		executable = "claude"
	}
	if _, err := exec.LookPath(executable); err != nil && !filepath.IsAbs(executable) {
		return nil, fmt.Errorf("claude-code executable: %w", err)
	}
	if req.Mode.AgentCapEnabled() {
		if err := claude.InstallHook(req.Workspace, req.HookCommand); err != nil {
			return nil, err
		}
		if req.Mode == ModeIntegrated {
			if err := claude.InstallInstructions(req.Workspace); err != nil {
				return nil, err
			}
		}
	}

	args := []string{
		"-p",
		"--output-format", "stream-json",
		"--verbose",
		"--no-session-persistence",
		"--dangerously-skip-permissions",
		"--setting-sources", "project",
		"--strict-mcp-config",
		"--tools", "Bash,Read,Edit,Write,Glob,Grep",
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
		return nil, fmt.Errorf("start claude-code: %w", err)
	}
	waitErr := cmd.Wait()
	result := &AgentRunResult{Stdout: stdout.Bytes(), Stderr: stderr.Bytes()}
	parseClaudeCodeEvents(result.Stdout, result)
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
	return nil, fmt.Errorf("wait for claude-code: %w", waitErr)
}

// parseClaudeCodeEvents extracts metrics from the Claude Code stream-json JSONL
// output and populates result in place.
func parseClaudeCodeEvents(data []byte, result *AgentRunResult) {
	shellTools := map[string]bool{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		var event struct {
			Type    string `json:"type"`
			Subtype string `json:"subtype"`
			Model   string `json:"model"`
			IsError bool   `json:"is_error"`
			Message *struct {
				Content []struct {
					Type      string          `json:"type"`
					Name      string          `json:"name"`
					ID        string          `json:"id"`
					ToolUseID string          `json:"tool_use_id"`
					Content   json.RawMessage `json:"content"`
				} `json:"content"`
			} `json:"message"`
			ToolUseResult json.RawMessage `json:"tool_use_result"`
			Usage         *struct {
				InputTokens              int64 `json:"input_tokens"`
				OutputTokens             int64 `json:"output_tokens"`
				CacheCreationInputTokens int64 `json:"cache_creation_input_tokens"`
				CacheReadInputTokens     int64 `json:"cache_read_input_tokens"`
			} `json:"usage"`
			ModelUsage map[string]struct {
				Input  int64 `json:"inputTokens"`
				Output int64 `json:"outputTokens"`
				Create int64 `json:"cacheCreationInputTokens"`
				Read   int64 `json:"cacheReadInputTokens"`
			} `json:"modelUsage"`
		}
		if json.Unmarshal(scanner.Bytes(), &event) != nil {
			continue
		}
		switch event.Type {
		case "system":
			if event.Subtype == "init" && event.Model != "" {
				result.Model = event.Model
			}
		case "assistant":
			if event.Message != nil {
				for _, item := range event.Message.Content {
					if item.Type == "tool_use" {
						result.ToolCount++
					}
					if item.Type == "tool_use" && item.Name == "Bash" {
						result.CommandCount++
						shellTools[item.ID] = true
					}
				}
			}
		case "user":
			var shellVisible int64
			matchedShell := false
			if event.Message != nil {
				for _, item := range event.Message.Content {
					if item.Type == "tool_result" && item.ToolUseID != "" && shellTools[item.ToolUseID] {
						matchedShell = true
						shellVisible += claudeContentBytes(item.Content)
					}
				}
			}
			var shell struct {
				Stdout string `json:"stdout"`
				Stderr string `json:"stderr"`
			}
			if matchedShell {
				result.CommandVisibleBytes += shellVisible
			} else if json.Unmarshal(event.ToolUseResult, &shell) == nil {
				result.CommandVisibleBytes += int64(len(shell.Stdout) + len(shell.Stderr))
			}
			if event.Message != nil {
				for _, item := range event.Message.Content {
					if item.Type == "tool_result" {
						result.ToolVisibleBytes += claudeContentBytes(item.Content)
					}
				}
			}
		case "result":
			if event.Usage != nil {
				result.UsageAvailable = !event.IsError || event.Usage.InputTokens+event.Usage.OutputTokens+event.Usage.CacheCreationInputTokens+event.Usage.CacheReadInputTokens > 0
				result.InputTokens = event.Usage.InputTokens
				result.OutputTokens = event.Usage.OutputTokens
				result.CacheCreationTokens = event.Usage.CacheCreationInputTokens
				result.CacheReadTokens = event.Usage.CacheReadInputTokens
			}
			if len(event.ModelUsage) > 0 {
				result.InputTokens, result.OutputTokens, result.CacheCreationTokens, result.CacheReadTokens = 0, 0, 0, 0
				for _, usage := range event.ModelUsage {
					result.InputTokens += usage.Input
					result.OutputTokens += usage.Output
					result.CacheCreationTokens += usage.Create
					result.CacheReadTokens += usage.Read
				}
				result.UsageAvailable = true
			}
		}
	}
}

func claudeContentBytes(data json.RawMessage) int64 {
	var value string
	if json.Unmarshal(data, &value) == nil {
		return int64(len(value))
	}
	var blocks []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	if json.Unmarshal(data, &blocks) != nil {
		return 0
	}
	var total int64
	for _, block := range blocks {
		if block.Type == "text" {
			total += int64(len(block.Text))
		}
	}
	return total
}
