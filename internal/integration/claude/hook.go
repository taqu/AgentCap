// Package claude provides Claude Code integration helpers for AgentCap.
package claude

import "encoding/json"

// HookInput is the JSON that Claude Code sends to a PreToolUse hook on stdin.
type HookInput struct {
	SessionID     string          `json:"session_id"`
	ToolName      string          `json:"tool_name"`
	ToolInput     json.RawMessage `json:"tool_input"`
	HookEventName string          `json:"hook_event_name"`
}

// BashInput is the tool_input payload for the Bash tool.
type BashInput struct {
	Command string `json:"command"`
	Timeout int    `json:"timeout,omitempty"`
}

// blockResponse is the JSON written to stdout when blocking a tool.
type blockResponse struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason"`
}

// allowResponse is the JSON written to stdout when allowing a tool.
type allowResponse struct {
	Decision string `json:"decision"`
}

// ParseHookInput parses Claude Code's PreToolUse hook stdin.
func ParseHookInput(data []byte) (*HookInput, error) {
	var h HookInput
	if err := json.Unmarshal(data, &h); err != nil {
		return nil, err
	}
	return &h, nil
}

// ParseBashInput parses the tool_input for a Bash tool call.
func ParseBashInput(raw json.RawMessage) (*BashInput, error) {
	var b BashInput
	if err := json.Unmarshal(raw, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// MakeBlockResponse creates the JSON response that blocks a tool and provides replacement output.
// Claude Code hook exit code 2 + this JSON = tool result shown to model.
func MakeBlockResponse(output string) ([]byte, error) {
	return json.Marshal(blockResponse{
		Decision: "block",
		Reason:   output,
	})
}

// MakeAllowResponse creates the JSON response that allows a tool to proceed.
func MakeAllowResponse() []byte {
	b, _ := json.Marshal(allowResponse{Decision: "allow"})
	return b
}
