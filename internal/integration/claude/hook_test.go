package claude

import (
	"encoding/json"
	"testing"
)

func TestParseHookInput(t *testing.T) {
	raw := `{
		"session_id": "abc123",
		"tool_name": "Bash",
		"tool_input": {"command": "echo hello"},
		"hook_event_name": "PreToolUse"
	}`
	h, err := ParseHookInput([]byte(raw))
	if err != nil {
		t.Fatalf("ParseHookInput error: %v", err)
	}
	if h.SessionID != "abc123" {
		t.Errorf("SessionID = %q, want %q", h.SessionID, "abc123")
	}
	if h.ToolName != "Bash" {
		t.Errorf("ToolName = %q, want %q", h.ToolName, "Bash")
	}
	if h.HookEventName != "PreToolUse" {
		t.Errorf("HookEventName = %q, want %q", h.HookEventName, "PreToolUse")
	}
}

func TestParseHookInputInvalid(t *testing.T) {
	_, err := ParseHookInput([]byte("not json"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestParseBashInput(t *testing.T) {
	raw := json.RawMessage(`{"command": "ls -la", "timeout": 30}`)
	b, err := ParseBashInput(raw)
	if err != nil {
		t.Fatalf("ParseBashInput error: %v", err)
	}
	if b.Command != "ls -la" {
		t.Errorf("Command = %q, want %q", b.Command, "ls -la")
	}
	if b.Timeout != 30 {
		t.Errorf("Timeout = %d, want %d", b.Timeout, 30)
	}
}

func TestMakeAllowResponse(t *testing.T) {
	data := MakeAllowResponse()
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(m) != 0 {
		t.Errorf("bypass must not grant permission: %v", m)
	}
}
