package claude

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/taqu/agentcap/internal/integration/conformance"
	"github.com/taqu/agentcap/internal/integration/protocol"
	"github.com/taqu/agentcap/internal/project"
	"github.com/taqu/agentcap/internal/store"
)

func hookData(command, cwd, session, mode string, background bool) []byte {
	data, _ := json.Marshal(map[string]any{"session_id": session, "cwd": cwd, "permission_mode": mode, "hook_event_name": "PreToolUse", "tool_name": "Bash", "tool_use_id": "tool-1", "agent_id": "subagent-1", "tool_input": map[string]any{"command": command, "timeout": 30000, "run_in_background": background, "description": "preserve this", "future_field": "keep"}})
	return data
}

func preparedPayload(t *testing.T, command, cwd, session string) string {
	t.Helper()
	executable, _ := os.Executable()
	response, reason, err := Prepare(hookData(command, cwd, session, "bypassPermissions", false), executable)
	if err != nil || reason != "intercepted" {
		t.Fatalf("Prepare: %s %v", reason, err)
	}
	var output struct {
		Hook struct {
			Input struct {
				Command     string `json:"command"`
				Timeout     int    `json:"timeout"`
				Description string `json:"description"`
				Future      string `json:"future_field"`
			} `json:"updatedInput"`
			Permission string `json:"permissionDecision"`
		} `json:"hookSpecificOutput"`
	}
	if err := json.Unmarshal(response, &output); err != nil {
		t.Fatal(err)
	}
	if output.Hook.Permission != "" || output.Hook.Input.Timeout != 30000 || output.Hook.Input.Description != "preserve this" || output.Hook.Input.Future != "keep" {
		t.Fatalf("changed original input: %s", response)
	}
	_, payload, ok := strings.Cut(output.Hook.Input.Command, " claude-exec ")
	if !ok {
		t.Fatal(output.Hook.Input.Command)
	}
	return strings.Trim(payload, "'")
}

func TestClaudeConformance(t *testing.T) {
	if _, err := bashPath(); err != nil {
		t.Skip("Bash unavailable:", err)
	}
	conformance.Run(t, func(ctx context.Context, command, cwd, session string) *protocol.ToolResponse {
		payload := preparedPayload(t, command, cwd, session)
		var stdout, stderr bytes.Buffer
		code := ExecutePayload(ctx, payload, &stdout, &stderr)
		resp := &protocol.ToolResponse{ExitCode: code, Stdout: stdout.String(), Stderr: stderr.String()}
		fields := strings.Fields(resp.Stdout)
		if len(fields) > 1 && fields[0] == "@acap" {
			resp.ResultID = fields[1]
			// Load the persisted common response metadata, not Claude transport IDs.
			st, err := store.Open(project.FindRoot(cwd))
			if err != nil {
				t.Error(err)
				return resp
			}
			defer st.Close()
			entry, err := st.Open(resp.ResultID)
			if err != nil {
				t.Error(err)
				return resp
			}
			resp.Presentation = entry.Meta.Presentation
		}
		return resp
	})
}

func TestPermissionAndBypass(t *testing.T) {
	executable, _ := os.Executable()
	for _, mode := range []string{"default", "acceptEdits", "auto", "plan", "dontAsk", "", "future-mode"} {
		data := hookData("printf 'x' > marker", t.TempDir(), "session", mode, false)
		response, reason, err := Prepare(data, executable)
		if err != nil || reason != "permission-mode" || string(response) != "{}" {
			t.Fatalf("mode %s: %s %s %v", mode, response, reason, err)
		}
	}
	for _, command := range []string{"acap show abc", "acap raw abc", "cd nested && pwd", "export FOO=bar", "source setup.sh"} {
		response, reason, err := Prepare(hookData(command, t.TempDir(), "s", "bypassPermissions", false), executable)
		if err != nil || reason != "shell-state-or-recursion" || string(response) != "{}" {
			t.Fatalf("%s: %s %s %v", command, response, reason, err)
		}
	}
	for _, env := range []string{protocol.EnvBypass, protocol.EnvDepth} {
		t.Run(env, func(t *testing.T) {
			t.Setenv(env, "1")
			response, _, _ := Prepare(hookData("echo ok", t.TempDir(), "s", "bypassPermissions", false), executable)
			if string(response) != "{}" {
				t.Fatal(string(response))
			}
		})
	}
	response, reason, _ := Prepare(hookData("sleep 1", t.TempDir(), "s", "bypassPermissions", true), executable)
	if reason != "background" || string(response) != "{}" {
		t.Fatal(reason, string(response))
	}
}

func TestShellRoundTrip(t *testing.T) {
	if _, err := bashPath(); err != nil {
		t.Skip(err)
	}
	dir := t.TempDir()
	for _, command := range []string{
		`printf '%s\n' "a b"`,
		`FOO="hello world" sh -c 'printf "%s\n" "$FOO"'`,
		`printf 'a\nb\n' | grep b`,
		`false || echo recovered`,
		`printf 'x\n' > file.txt`,
		`printf '%s' '\"; touch injected; #'`,
	} {
		payload := preparedPayload(t, command, dir, "")
		var p execution
		decoded, _ := base64.RawURLEncoding.DecodeString(payload)
		json.Unmarshal(decoded, &p)
		if *p.Request.ShellCommand != command {
			t.Fatal("shell source changed")
		}
		var stdout, stderr bytes.Buffer
		if code := ExecutePayload(context.Background(), payload, &stdout, &stderr); code != 0 {
			t.Fatalf("%s: %d %s", command, code, stderr.String())
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "injected")); !os.IsNotExist(err) {
		t.Fatal("wrapper injected command syntax")
	}
	data, _ := os.ReadFile(filepath.Join(dir, "file.txt"))
	if string(data) != "x\n" {
		t.Fatalf("redirection: %q", data)
	}
}

func TestInvalidPayloadNeverExecutes(t *testing.T) {
	for _, payload := range []string{"invalid", base64.RawURLEncoding.EncodeToString([]byte(`{}`))} {
		var stdout, stderr bytes.Buffer
		if code := ExecutePayload(context.Background(), payload, &stdout, &stderr); code != 2 || stdout.Len() != 0 {
			t.Fatalf("%d %s", code, stdout.String())
		}
	}
}
