// Package claude provides Claude Code integration helpers for AgentCap.
package claude

import (
	"encoding/json"
	"fmt"
	"strings"
)

// HookInput is the JSON that Claude Code sends to a PreToolUse hook on stdin.
type HookInput struct {
	SessionID      string          `json:"session_id"`
	ToolName       string          `json:"tool_name"`
	ToolInput      json.RawMessage `json:"tool_input"`
	HookEventName  string          `json:"hook_event_name"`
	Cwd            string          `json:"cwd"`
	PermissionMode string          `json:"permission_mode"`
	ToolUseID      string          `json:"tool_use_id"`
	AgentID        string          `json:"agent_id"`
}

// BashInput is the tool_input payload for the Bash tool.
type BashInput struct {
	Command         string `json:"command"`
	Timeout         int    `json:"timeout,omitempty"`
	RunInBackground bool   `json:"run_in_background,omitempty"`
}

// ParseHookInput parses Claude Code's PreToolUse hook stdin.
func ParseHookInput(data []byte) (*HookInput, error) {
	var h HookInput
	if err := json.Unmarshal(data, &h); err != nil {
		return nil, err
	}
	if h.ToolName == "" || h.HookEventName == "" || len(h.ToolInput) == 0 {
		return nil, fmt.Errorf("missing hook fields")
	}
	return &h, nil
}

// ParseBashInput parses the tool_input for a Bash tool call.
func ParseBashInput(raw json.RawMessage) (*BashInput, error) {
	var b BashInput
	if err := json.Unmarshal(raw, &b); err != nil {
		return nil, err
	}
	if strings.TrimSpace(b.Command) == "" || b.Timeout < 0 {
		return nil, fmt.Errorf("invalid Bash input")
	}
	return &b, nil
}

// MakeAllowResponse returns no decision, leaving permission checks to Claude.
func MakeAllowResponse() []byte {
	// No decision: preserve Claude's normal permission flow, never grant allow.
	return []byte(`{}`)
}
