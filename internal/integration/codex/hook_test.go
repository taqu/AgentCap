package codex

import (
	"encoding/json"
	"testing"
)

func TestParseHookInput(t *testing.T) {
	raw := `{
		"session_id": "xyz789",
		"tool": "shell",
		"tool_input": {"cmd": "git diff"},
		"hook_event": "pre_tool_use"
	}`
	h, err := ParseHookInput([]byte(raw))
	if err != nil {
		t.Fatalf("ParseHookInput error: %v", err)
	}
	if h.SessionID != "xyz789" {
		t.Errorf("SessionID = %q, want %q", h.SessionID, "xyz789")
	}
	if h.Tool != "shell" {
		t.Errorf("Tool = %q, want %q", h.Tool, "shell")
	}
	if h.HookEvent != "pre_tool_use" {
		t.Errorf("HookEvent = %q, want %q", h.HookEvent, "pre_tool_use")
	}
}

func TestParseHookInputInvalid(t *testing.T) {
	_, err := ParseHookInput([]byte("not json"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestParseShellInput(t *testing.T) {
	raw := json.RawMessage(`{"cmd": "ls -la", "cmdline": ["ls", "-la"]}`)
	s, err := ParseShellInput(raw)
	if err != nil {
		t.Fatalf("ParseShellInput error: %v", err)
	}
	if s.Cmd != "ls -la" {
		t.Errorf("Cmd = %q, want %q", s.Cmd, "ls -la")
	}
	if len(s.Cmdline) != 2 || s.Cmdline[0] != "ls" || s.Cmdline[1] != "-la" {
		t.Errorf("Cmdline = %v, want [ls -la]", s.Cmdline)
	}
}

func TestMakeBlockResponse(t *testing.T) {
	out := "compressed output"
	data, err := MakeBlockResponse(out)
	if err != nil {
		t.Fatalf("MakeBlockResponse error: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if m["decision"] != "block" {
		t.Errorf("decision = %v, want block", m["decision"])
	}
	if m["reason"] != out {
		t.Errorf("reason = %v, want %q", m["reason"], out)
	}
}

func TestMakeAllowResponse(t *testing.T) {
	data := MakeAllowResponse()
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if m["decision"] != "allow" {
		t.Errorf("decision = %v, want allow", m["decision"])
	}
}
