package agentbench

import (
	"context"
	"fmt"
)

// Adapter isolates agent-specific invocation and event parsing.
type Adapter interface {
	Name() string
	Run(context.Context, AgentRunRequest) (*AgentRunResult, error)
}

type AgentRunRequest struct {
	Workspace   string
	Task        string
	Mode        Mode
	StoreRoot   string
	HookCommand string
	Model       string
}

type AgentRunResult struct {
	ExitCode            int
	TimedOut            bool
	Canceled            bool
	CommandCount        int
	CommandVisibleBytes int64
	Stdout              []byte
	Stderr              []byte
}

// NewAdapter resolves a supported real adapter without leaking agent-specific
// selection logic into the benchmark runner.
func NewAdapter(name, model string) (Adapter, error) {
	switch name {
	case "codex":
		return &CodexAdapter{Executable: "codex", Model: model}, nil
	default:
		return nil, fmt.Errorf("unsupported coding agent %q (supported: codex)", name)
	}
}
