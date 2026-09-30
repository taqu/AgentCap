package agentbench

import (
	"context"
	"fmt"

	"github.com/taqu/agentcap/internal/workload"
)

// Adapter isolates agent-specific invocation and event parsing.
type Adapter interface {
	Name() string
	Run(context.Context, AgentRunRequest) (*AgentRunResult, error)
}

type AgentRunRequest struct {
	Workspace    string
	Task         string
	Mode         Mode
	StoreRoot    string
	HookCommand  string
	Model        string
	SessionScope string
	// BinDir is prepended to PATH so production hooks resolve "acap" to the
	// benchmarked binary.
	BinDir string
}

type AgentRunResult struct {
	ExitCode            int
	TimedOut            bool
	Canceled            bool
	CommandCount        int
	CommandVisibleBytes int64
	Stdout              []byte
	Stderr              []byte

	// Phase 7 workflow evidence; nil/zero when the adapter cannot provide it.
	Trace            []workload.TraceEntry
	AgentVersion     string
	InstructionBytes int64
	HostInput        int64
	HostCacheRead    int64
	HostCacheWrite   int64
	HostOutput       int64
	// HostError is a host/service failure reported by the agent (quota, API
	// outage). It makes the trial INFRA_ERROR regardless of verification.
	HostError string
}

// ModeSupporter is implemented by adapters whose production integration does
// not expose every benchmark mode.
type ModeSupporter interface {
	SupportsMode(Mode) bool
}

// NewAdapter resolves a supported real adapter without leaking agent-specific
// selection logic into the benchmark runner.
func NewAdapter(name, model string) (Adapter, error) {
	switch name {
	case "codex":
		return &CodexAdapter{Executable: "codex", Model: model}, nil
	case "claude":
		return &ClaudeAdapter{Executable: "claude", Model: model}, nil
	case "antigravity":
		return &AntigravityAdapter{Executable: "agy", Model: model}, nil
	default:
		return nil, fmt.Errorf("unsupported coding agent %q (supported: codex, claude, antigravity)", name)
	}
}
