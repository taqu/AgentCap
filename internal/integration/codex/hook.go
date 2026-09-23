// Package codex provides OpenAI Codex CLI integration helpers for AgentCap.
package codex

import "encoding/json"

// HookInput is the JSON that Codex CLI sends to a pre_tool_use hook on stdin.
type HookInput struct {
	SessionID string          `json:"session_id"`
	Tool      string          `json:"tool"`
	ToolInput json.RawMessage `json:"tool_input"`
	HookEvent string          `json:"hook_event"`
}

// ShellInput is the tool_input for Codex's shell tool.
type ShellInput struct {
	Cmd     string   `json:"cmd"`
	Cmdline []string `json:"cmdline,omitempty"`
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

// ParseHookInput parses Codex CLI's pre_tool_use hook stdin.
func ParseHookInput(data []byte) (*HookInput, error) {
	var h HookInput
	if err := json.Unmarshal(data, &h); err != nil {
		return nil, err
	}
	return &h, nil
}

// ParseShellInput parses the tool_input for a shell tool call.
func ParseShellInput(raw json.RawMessage) (*ShellInput, error) {
	var s ShellInput
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// MakeBlockResponse creates the JSON response that blocks a tool and provides replacement output.
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
