package antigravity

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
	"strconv"
	"strings"
	"time"

	"github.com/taqu/agentcap/internal/integration/common"
	"github.com/taqu/agentcap/internal/integration/engine"
	"github.com/taqu/agentcap/internal/integration/protocol"
)

const (
	AdapterVersion = "1"
	agentName      = "antigravity"

	// EnvPermissions declares that the caller's Antigravity session already
	// runs without command-level permission checks (always-proceed /
	// --dangerously-skip-permissions). Antigravity evaluates permissions
	// against the overwritten CommandLine and does not expose its mode to
	// hooks, so rewriting is only equivalent when this is explicitly declared.
	EnvPermissions      = "ACAP_ANTIGRAVITY_PERMISSIONS"
	permissionsDeclared = "unrestricted"

	// minWaitMs: smaller WaitMsBeforeAsync values signal that the agent
	// expects the command to be promoted to a background task. The wrapper
	// only emits output on completion, so those commands are bypassed.
	minWaitMs = 2000
)

type execution struct {
	Request protocol.ToolRequest `json:"request"`
}

// Retrieval commands must never re-enter interception; quoted occurrences are
// over-matched on purpose.
var acapCommand = regexp.MustCompile(`(^|[^a-zA-Z0-9_./-])(\S*/)?acap(\.exe)?([^a-zA-Z0-9_.-]|$)`)

// Long-lived programs whose output must stream through Antigravity.
var longLived = regexp.MustCompile(`(^|[\s;&|(])(tail\s+(-\S*\s+)*-[a-zA-Z]*[fF]|watch\s|top($|\s)|htop|less($|\s)|more($|\s)|vim?($|\s)|nano\s|ssh\s)|--watch\b|\b(npm|pnpm|yarn|bun)\s+(run\s+)?(dev|start|serve|watch)\b|\b(serve|nodemon|http-server|live-server)\b|runserver|flask\s+run|uvicorn|vite($|\s)|next\s+dev`)

func shellPath() (string, error) {
	for _, name := range []string{"bash", "sh"} {
		if path, err := osexec.LookPath(name); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("no POSIX shell found")
}

// wrapperBinary prefers the bare name when PATH resolves to this binary, so
// Antigravity's UI and command(acap ...) permission rules see "acap".
func wrapperBinary(executable string) string {
	if found, err := osexec.LookPath("acap"); err == nil {
		a, errA := filepath.EvalSymlinks(found)
		b, errB := filepath.EvalSymlinks(executable)
		if errA == nil && errB == nil && a == b {
			return "acap"
		}
	}
	return shellQuote(filepath.ToSlash(executable))
}

func shellQuote(value string) string { return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'" }

// Prepare never executes commands. It returns the PreToolUse response, a
// reason for metrics, and a diagnostic error. Every path returns decision
// "ask": the adapter never grants permission.
func Prepare(data []byte, executable string) ([]byte, string, error) {
	started := time.Now()
	h, err := ParseHookInput(data)
	if err != nil {
		return MakePassResponse(), "invalid-input", err
	}
	if h.ToolCall.Name != "run_command" {
		return MakePassResponse(), "unsupported-tool", nil
	}
	args, err := ParseRunCommandArgs(h.ToolCall.Args)
	if err != nil {
		return MakePassResponse(), "invalid-input", err
	}
	reason := ""
	switch {
	case common.IsBypassed() || os.Getenv(protocol.EnvBenchmarkMode) == "disabled":
		reason = "explicit-bypass"
	case common.IsRecursive() || acapCommand.MatchString(args.CommandLine):
		reason = "recursion"
	case os.Getenv(EnvPermissions) != permissionsDeclared:
		reason = "permission-mode"
	case args.Cwd == "" || !filepath.IsAbs(args.Cwd):
		reason = "missing-cwd"
	case args.WaitMsBeforeAsync != nil && *args.WaitMsBeforeAsync < minWaitMs:
		reason = "async"
	case longLived.MatchString(args.CommandLine):
		reason = "long-running-or-interactive"
	}
	if reason != "" {
		recordMetric(args.Cwd, h, reason, false, time.Since(started), 0)
		return MakePassResponse(), reason, nil
	}
	shell, err := shellPath()
	if err != nil {
		recordMetric(args.Cwd, h, "shell-unavailable", true, time.Since(started), 0)
		return MakePassResponse(), "shell-unavailable", err
	}
	if fi, err := os.Stat(executable); err != nil || fi.IsDir() {
		recordMetric(args.Cwd, h, "binary-unavailable", true, time.Since(started), 0)
		return MakePassResponse(), "binary-unavailable", fmt.Errorf("AgentCap binary unavailable")
	}
	meta := &protocol.Metadata{Agent: agentName, Adapter: agentName, Version: AdapterVersion, ExternalSession: h.ConversationID, Model: h.ModelName}
	if h.StepIdx != nil {
		meta.ToolUseID = strconv.FormatInt(*h.StepIdx, 10)
	}
	// workspacePaths are not the execution directory; the core resolves the
	// project root from the actual cwd.
	// Antigravity runs each command as its own process-group leader and
	// cancels by killing that group; the command must stay inside it.
	req := protocol.ToolRequest{Protocol: protocol.Version, Shell: shell, ShellCommand: &args.CommandLine, WorkingDir: args.Cwd, InheritProcessGroup: true, Integration: meta}
	meta.AdapterLatencyNs = time.Since(started).Nanoseconds()
	encoded, err := json.Marshal(execution{Request: req})
	if err != nil {
		return MakePassResponse(), "encoding-failure", err
	}
	// Source travels only as base64url argument data, never as shell syntax.
	// overwrite is a shallow merge, so Cwd and WaitMsBeforeAsync are unchanged.
	command := wrapperBinary(executable) + " antigravity-exec " + shellQuote(base64.RawURLEncoding.EncodeToString(encoded))
	response, err := makeOverwriteResponse(command)
	if err != nil {
		return MakePassResponse(), "encoding-failure", err
	}
	return response, "intercepted", nil
}

// ExecutePayload is the wrapper entry point run by Antigravity in place of the
// original CommandLine. It runs the command exactly once and exits with its
// status. Exit 2 with no execution means the payload was rejected before start.
func ExecutePayload(ctx context.Context, encoded string, stdout, stderr io.Writer) int {
	data, err := base64.RawURLEncoding.DecodeString(encoded)
	var p execution
	if err != nil || json.Unmarshal(data, &p) != nil || p.Request.ShellCommand == nil || p.Request.Integration == nil || p.Request.Shell == "" {
		fmt.Fprintln(stderr, "acap: invalid Antigravity execution payload; command not started")
		return 2
	}
	out, err := engine.Run(ctx, &p.Request)
	if err != nil || out == nil || out.Response == nil {
		fmt.Fprintln(stderr, "acap: Antigravity execution unavailable; command not started:", err)
		return 1
	}
	resp := out.Response
	io.WriteString(stdout, resp.Stdout)
	// On success the capsule carries relevant diagnostics; raw stderr stays in
	// the store. Fail-open paths return the captured streams unchanged.
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
	h := &HookInput{ConversationID: p.Request.Integration.ExternalSession}
	if id, err := strconv.ParseInt(p.Request.Integration.ToolUseID, 10, 64); err == nil {
		h.StepIdx = &id
	}
	recordMetric(p.Request.WorkingDir, h, reason, resp.Error != "", time.Duration(p.Request.Integration.AdapterLatencyNs), out.ProcessingDuration)
	if ctx.Err() != nil {
		fmt.Fprintln(stderr, "acap: command interrupted")
		return 124
	}
	if resp.ExitCode < 0 {
		return 128
	}
	return resp.ExitCode
}

func recordMetric(cwd string, h *HookInput, reason string, failed bool, adapter, processing time.Duration) {
	m := common.AdapterMetric{Agent: agentName, Adapter: agentName, Version: AdapterVersion, Session: h.ConversationID, Reason: reason, Failed: failed, AdapterLatency: adapter, ProcessingLatency: processing}
	if h.StepIdx != nil {
		m.ToolUseID = strconv.FormatInt(*h.StepIdx, 10)
	}
	common.RecordMetric(cwd, m)
}
