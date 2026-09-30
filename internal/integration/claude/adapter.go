package claude

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	osexec "os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/taqu/agentcap/internal/integration/common"
	"github.com/taqu/agentcap/internal/integration/engine"
	"github.com/taqu/agentcap/internal/integration/protocol"
	"github.com/taqu/agentcap/internal/project"
)

const AdapterVersion = "1"

type execution struct {
	Request protocol.ToolRequest `json:"request"`
	Timeout int                  `json:"timeout,omitempty"`
}

// Shell state changes must remain in Claude's persistent Bash execution context.
// This intentionally over-bypasses quoted occurrences rather than guessing ASTs.
var statefulShell = regexp.MustCompile(`(^|[^a-zA-Z0-9_])(cd|pushd|popd|export|unset|source|alias|unalias|umask|set|trap|eval|exec|acap|acap\.exe)([^a-zA-Z0-9_]|$)`)

func bashPath() (string, error) {
	if runtime.GOOS == "windows" {
		if path := os.Getenv("CLAUDE_CODE_GIT_BASH_PATH"); path != "" {
			if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
				return path, nil
			}
			return "", fmt.Errorf("configured Git Bash is unavailable")
		}
	}
	path, err := osexec.LookPath("bash")
	if err != nil {
		return "", err
	}
	if runtime.GOOS == "windows" && strings.Contains(strings.ToLower(path), `\windows\`) {
		return "", fmt.Errorf("WSL Bash is unsupported; configure CLAUDE_CODE_GIT_BASH_PATH")
	}
	return path, nil
}

func shellQuote(value string) string { return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'" }

// Prepare never executes commands. Only an already explicitly unpermissioned
// Claude session may use rewriting: original allow/deny rules cannot be preserved
// by the public hook API. Other modes safely retain ordinary Bash behavior.
func Prepare(data []byte, executable string) ([]byte, string, error) {
	started := time.Now()
	h, err := ParseHookInput(data)
	if err != nil {
		return MakeAllowResponse(), "invalid-input", err
	}
	if h.ToolName != "Bash" || h.HookEventName != "PreToolUse" {
		return MakeAllowResponse(), "unsupported-tool", nil
	}
	b, err := ParseBashInput(h.ToolInput)
	if err != nil {
		return MakeAllowResponse(), "invalid-input", err
	}
	reason := ""
	switch {
	case common.IsBypassed() || os.Getenv(protocol.EnvBenchmarkMode) == "disabled":
		reason = "explicit-bypass"
	case common.IsRecursive():
		reason = "recursion"
	case b.RunInBackground:
		reason = "background"
	case h.PermissionMode != "bypassPermissions":
		reason = "permission-mode"
	case h.Cwd == "" || !filepath.IsAbs(h.Cwd):
		reason = "missing-cwd"
	case statefulShell.MatchString(b.Command):
		reason = "shell-state-or-recursion"
	}
	if reason != "" {
		recordMetric(h.Cwd, h, reason, false, time.Since(started), 0)
		return MakeAllowResponse(), reason, nil
	}
	shell, err := bashPath()
	if err != nil {
		recordMetric(h.Cwd, h, "shell-unavailable", true, time.Since(started), 0)
		return MakeAllowResponse(), "shell-unavailable", err
	}
	if fi, err := os.Stat(executable); err != nil || fi.IsDir() {
		return MakeAllowResponse(), "binary-unavailable", fmt.Errorf("AgentCap binary unavailable")
	}
	req := protocol.ToolRequest{Protocol: protocol.Version, Shell: shell, ShellCommand: &b.Command, WorkingDir: h.Cwd,
		Integration: &protocol.Metadata{Agent: "claude-code", Adapter: "claude", Version: AdapterVersion, ExternalSession: h.SessionID, ToolUseID: h.ToolUseID, SubagentID: h.AgentID}}
	req.Integration.AdapterLatencyNs = time.Since(started).Nanoseconds()
	encoded, err := json.Marshal(execution{Request: req, Timeout: b.Timeout})
	if err != nil {
		return MakeAllowResponse(), "encoding-failure", err
	}
	// Carry source solely as base64 argument data. All other Bash input fields,
	// including future fields, are preserved rather than reconstructed.
	var input map[string]json.RawMessage
	if err := json.Unmarshal(h.ToolInput, &input); err != nil {
		return MakeAllowResponse(), "invalid-input", err
	}
	binary := filepath.ToSlash(executable)
	command := shellQuote(binary) + " claude-exec " + shellQuote(base64.RawURLEncoding.EncodeToString(encoded))
	input["command"], _ = json.Marshal(command)
	response, err := json.Marshal(map[string]any{"hookSpecificOutput": map[string]any{"hookEventName": "PreToolUse", "updatedInput": input}})
	return response, "intercepted", err
}

// ExecutePayload is the wrapper process entry point. Never retry execution after
// processing failure: the common engine returns original capture fail-open.
func ExecutePayload(ctx context.Context, encoded string, stdout, stderr io.Writer) int {
	data, err := base64.RawURLEncoding.DecodeString(encoded)
	var p execution
	if err != nil || json.Unmarshal(data, &p) != nil || p.Request.ShellCommand == nil || p.Request.Integration == nil || p.Timeout < 0 {
		fmt.Fprintln(stderr, "acap: invalid Claude execution payload")
		return 2
	}
	if p.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(p.Timeout)*time.Millisecond)
		defer cancel()
	}
	out, err := engine.Run(ctx, &p.Request)
	if err != nil || out == nil || out.Response == nil {
		fmt.Fprintln(stderr, "acap: Claude execution unavailable:", err)
		return 1
	}
	resp := out.Response
	// Successful persistence means core presentation includes relevant diagnostics;
	// do not also leak the complete original stderr into Claude's failure context.
	io.WriteString(stdout, resp.Stdout)
	if resp.ResultID == "" || resp.Error != "" {
		io.WriteString(stderr, resp.Stderr)
	}
	if resp.Error != "" {
		fmt.Fprintln(stderr, "acap:", resp.Error)
	}
	reason := "intercepted"
	if resp.Error != "" {
		reason = "processing-failure"
	}
	h := &HookInput{SessionID: p.Request.Integration.ExternalSession, ToolUseID: p.Request.Integration.ToolUseID, AgentID: p.Request.Integration.SubagentID}
	recordMetric(p.Request.WorkingDir, h, reason, resp.Error != "", time.Duration(p.Request.Integration.AdapterLatencyNs), out.ProcessingDuration)
	if ctx.Err() != nil {
		fmt.Fprintln(stderr, "acap: command interrupted or timed out")
		return 124
	}
	if resp.ExitCode < 0 {
		return 128
	}
	return resp.ExitCode
}

// One bounded append per invocation; no command text or environment values.
func recordMetric(cwd string, h *HookInput, reason string, failed bool, adapter, processing time.Duration) {
	if cwd == "" || !filepath.IsAbs(cwd) {
		return
	}
	root := project.FindRoot(cwd)
	dir := filepath.Join(root, ".acap")
	if os.MkdirAll(dir, 0755) != nil {
		return
	}
	metric := struct {
		Agent        string `json:"agent"`
		Adapter      string `json:"adapter"`
		Version      string `json:"adapter_version"`
		Session      string `json:"session_mapping"`
		ToolUseID    string `json:"tool_use_id,omitempty"`
		SubagentID   string `json:"subagent_id,omitempty"`
		Reason       string `json:"reason"`
		Intercepted  int    `json:"commands_intercepted"`
		Bypassed     int    `json:"commands_bypassed"`
		Failures     int    `json:"adapter_failures"`
		AdapterNs    int64  `json:"adapter_latency_ns"`
		ProcessingNs int64  `json:"processing_latency_ns"`
	}{Agent: "claude-code", Adapter: "claude", Version: AdapterVersion, Session: common.MapSession("claude-code", h.SessionID, root), ToolUseID: h.ToolUseID, SubagentID: h.AgentID, Reason: reason, AdapterNs: adapter.Nanoseconds(), ProcessingNs: processing.Nanoseconds()}
	if reason == "intercepted" || reason == "processing-failure" {
		metric.Intercepted = 1
	} else {
		metric.Bypassed = 1
	}
	if failed {
		metric.Failures = 1
	}
	data, _ := json.Marshal(metric)
	f, err := os.OpenFile(filepath.Join(dir, "adapter-metrics.jsonl"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer f.Close()
	f.Write(append(data, '\n'))
}
