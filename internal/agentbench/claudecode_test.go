package agentbench

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const sampleClaudeOutput = `{"type":"system","subtype":"init","model":"claude-test","session_id":"s1","tools":[]}
{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Bash","id":"t1","input":{"command":"echo hi"}}]}}
{"type":"user","tool_use_result":{"stdout":"hi\n","stderr":""}}
{"type":"assistant","message":{"content":[{"type":"tool_use","name":"Bash","id":"t2","input":{"command":"ls"}}]}}
{"type":"user","tool_use_result":{"stdout":"file.txt\n","stderr":"warn\n"}}
{"type":"result","subtype":"success","usage":{"input_tokens":10,"output_tokens":5,"cache_creation_input_tokens":100,"cache_read_input_tokens":50}}
`

func TestParseClaudeCodeEvents(t *testing.T) {
	result := &AgentRunResult{}
	parseClaudeCodeEvents([]byte(sampleClaudeOutput), result)

	if result.Model != "claude-test" {
		t.Errorf("model = %q, want %q", result.Model, "claude-test")
	}
	if result.CommandCount != 2 {
		t.Errorf("CommandCount = %d, want 2", result.CommandCount)
	}
	// stdout "hi\n" + stdout "file.txt\n" + stderr "warn\n"
	wantBytes := int64(len("hi\n") + len("file.txt\n") + len("warn\n"))
	if result.CommandVisibleBytes != wantBytes {
		t.Errorf("CommandVisibleBytes = %d, want %d", result.CommandVisibleBytes, wantBytes)
	}
	if result.InputTokens != 10 {
		t.Errorf("InputTokens = %d, want 10", result.InputTokens)
	}
	if result.OutputTokens != 5 {
		t.Errorf("OutputTokens = %d, want 5", result.OutputTokens)
	}
	if result.CacheCreationTokens != 100 {
		t.Errorf("CacheCreationTokens = %d, want 100", result.CacheCreationTokens)
	}
	if result.CacheReadTokens != 50 {
		t.Errorf("CacheReadTokens = %d, want 50", result.CacheReadTokens)
	}
}

func TestParseClaudeCodeEventsEmptyAndMalformed(t *testing.T) {
	result := &AgentRunResult{}
	parseClaudeCodeEvents([]byte(""), result)
	if result.CommandCount != 0 || result.InputTokens != 0 {
		t.Error("empty input should produce zero metrics")
	}

	result2 := &AgentRunResult{}
	parseClaudeCodeEvents([]byte("{not valid json}\n{\"type\":\"result\",\"usage\":{\"input_tokens\":3,\"output_tokens\":1}}\n"), result2)
	if result2.InputTokens != 3 || result2.OutputTokens != 1 {
		t.Errorf("malformed line followed by valid: tokens = %d/%d", result2.InputTokens, result2.OutputTokens)
	}
}

func TestClaudeCodeAdapterInvocationModesExitAndTimeout(t *testing.T) {
	exe, _ := os.Executable()
	t.Setenv("AGENTBENCH_CLAUDE_PROCESS", "1")
	adapter := &ClaudeCodeAdapter{Executable: exe}
	root := t.TempDir()
	request := AgentRunRequest{Workspace: t.TempDir(), Task: "test task", Mode: ModeDisabled, StoreRoot: root}
	result, err := adapter.Run(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	// mock emits two tool_use_result events with empty stderr
	wantVisibleBytes := int64(len("hi\n") + len("file.txt\n"))
	if result.CommandCount != 2 || result.CommandVisibleBytes != wantVisibleBytes {
		t.Errorf("unexpected result: CommandCount=%d CommandVisibleBytes=%d (want %d)", result.CommandCount, result.CommandVisibleBytes, wantVisibleBytes)
	}
	if result.InputTokens != 10 || result.OutputTokens != 5 {
		t.Errorf("tokens: input=%d output=%d", result.InputTokens, result.OutputTokens)
	}
	if result.Model != "claude-test" {
		t.Errorf("model = %q", result.Model)
	}
	if !strings.Contains(string(result.Stderr), root+"|disabled") {
		t.Errorf("stderr missing env markers: %q", result.Stderr)
	}
	if !strings.Contains(string(result.Stderr), "-p") {
		t.Errorf("stderr missing -p flag: %q", result.Stderr)
	}

	t.Setenv("AGENTBENCH_EXIT", "1")
	result, err = adapter.Run(context.Background(), request)
	if err != nil || result.ExitCode != 7 {
		t.Fatalf("nonzero exit: result=%+v err=%v", result, err)
	}

	t.Setenv("AGENTBENCH_EXIT", "")
	t.Setenv("AGENTBENCH_SLEEP", "1")
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	result, err = adapter.Run(ctx, request)
	if err != nil || !result.TimedOut {
		t.Fatalf("timeout: result=%+v err=%v", result, err)
	}
}

func TestClaudeCodeAdapterInstallsHookForAgentCapModes(t *testing.T) {
	exe, _ := os.Executable()
	t.Setenv("AGENTBENCH_CLAUDE_PROCESS", "1")
	adapter := &ClaudeCodeAdapter{Executable: exe}

	for _, tc := range []struct {
		mode         Mode
		wantHook     bool
		wantCLAUDEmd bool
	}{
		{ModeDisabled, false, false},
		{ModeStateless, true, false},
		{ModeStateful, true, false},
		{ModeIntegrated, true, true},
	} {
		t.Run(string(tc.mode), func(t *testing.T) {
			workspace := t.TempDir()
			_, err := adapter.Run(context.Background(), AgentRunRequest{
				Workspace:   workspace,
				Task:        "task",
				Mode:        tc.mode,
				StoreRoot:   t.TempDir(),
				HookCommand: `"` + filepath.ToSlash(exe) + `" hook claude`,
			})
			if err != nil {
				t.Fatal(err)
			}
			settingsPath := filepath.Join(workspace, ".claude", "settings.json")
			_, settingsErr := os.Stat(settingsPath)
			if tc.wantHook && settingsErr != nil {
				t.Fatalf("settings.json missing for mode %s: %v", tc.mode, settingsErr)
			}
			if !tc.wantHook && settingsErr == nil {
				t.Fatalf("settings.json unexpectedly present for mode %s", tc.mode)
			}
			claudeMDPath := filepath.Join(workspace, "CLAUDE.md")
			_, mdErr := os.Stat(claudeMDPath)
			if tc.wantCLAUDEmd && mdErr != nil {
				t.Fatalf("CLAUDE.md missing for mode %s: %v", tc.mode, mdErr)
			}
			if !tc.wantCLAUDEmd && mdErr == nil {
				t.Fatalf("CLAUDE.md unexpectedly present for mode %s", tc.mode)
			}
		})
	}
}

func TestClaudeCodeAdapterPassesDangerouslySkipPermissionsOnlyWhenEnabled(t *testing.T) {
	exe, _ := os.Executable()
	t.Setenv("AGENTBENCH_CLAUDE_PROCESS", "1")
	adapter := &ClaudeCodeAdapter{Executable: exe}

	for _, tc := range []struct {
		mode    Mode
		wantArg string
	}{
		{ModeDisabled, ""},
		{ModeStateless, "--dangerously-skip-permissions"},
		{ModeIntegrated, "--dangerously-skip-permissions"},
	} {
		t.Run(string(tc.mode), func(t *testing.T) {
			result, err := adapter.Run(context.Background(), AgentRunRequest{
				Workspace: t.TempDir(),
				Task:      "task",
				Mode:      tc.mode,
				StoreRoot: t.TempDir(),
			})
			if err != nil {
				t.Fatal(err)
			}
			stderr := string(result.Stderr)
			if tc.wantArg != "" && !strings.Contains(stderr, tc.wantArg) {
				t.Errorf("mode %s: expected %q in args, stderr=%q", tc.mode, tc.wantArg, stderr)
			}
			if tc.wantArg == "" && strings.Contains(stderr, "--dangerously-skip-permissions") {
				t.Errorf("mode %s: unexpected --dangerously-skip-permissions in args, stderr=%q", tc.mode, stderr)
			}
		})
	}
}

func TestClaudeCodeAdapterModelFlag(t *testing.T) {
	exe, _ := os.Executable()
	t.Setenv("AGENTBENCH_CLAUDE_PROCESS", "1")

	adapter := &ClaudeCodeAdapter{Executable: exe, Model: "claude-opus-4-7"}
	result, err := adapter.Run(context.Background(), AgentRunRequest{
		Workspace: t.TempDir(), Task: "task", Mode: ModeDisabled, StoreRoot: t.TempDir(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(result.Stderr), "--model") || !strings.Contains(string(result.Stderr), "claude-opus-4-7") {
		t.Errorf("model flag missing: %q", result.Stderr)
	}

	// Request model overrides adapter model.
	result2, err := adapter.Run(context.Background(), AgentRunRequest{
		Workspace: t.TempDir(), Task: "task", Mode: ModeDisabled, StoreRoot: t.TempDir(), Model: "claude-haiku-4-5",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(result2.Stderr), "claude-haiku-4-5") {
		t.Errorf("request model not used: %q", result2.Stderr)
	}
	if strings.Contains(string(result2.Stderr), "claude-opus-4-7") {
		t.Errorf("adapter model leaked when request model set: %q", result2.Stderr)
	}
}
