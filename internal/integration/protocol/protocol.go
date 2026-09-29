// Package protocol defines the wire types for the AgentCap coding-agent
// integration protocol.
package protocol

const Version = 1

// EnvBypass is checked to skip AgentCap interception.
const EnvBypass = "ACAP_BYPASS"

// EnvDepth prevents recursive interception.
const EnvDepth = "ACAP_INTERCEPT_DEPTH"

// ToolRequest is the JSON payload sent to "acap exec --protocol=json" on stdin.
type ToolRequest struct {
	Protocol   int      `json:"protocol"`
	Command    []string `json:"command"`
	WorkingDir string   `json:"working_dir"`
	// StoreRoot optionally separates result persistence from command execution.
	// It is used by benchmark workloads whose working directory is disposable.
	StoreRoot string            `json:"store_root,omitempty"`
	Env       map[string]string `json:"env,omitempty"`
	SessionID string            `json:"session_id,omitempty"`
	MaxBytes  int               `json:"max_bytes,omitempty"`
}

// ToolResponse is written to stdout by "acap exec --protocol=json".
type ToolResponse struct {
	Protocol     int    `json:"protocol"`
	ResultID     string `json:"result_id,omitempty"`
	ExitCode     int    `json:"exit_code"`
	Stdout       string `json:"stdout"`
	Stderr       string `json:"stderr"`
	Presentation string `json:"presentation,omitempty"`
	Error        string `json:"error,omitempty"`
}

// PresentationBudget carries optional byte-budget hints to reducers.
type PresentationBudget struct {
	MaxBytes int
}
