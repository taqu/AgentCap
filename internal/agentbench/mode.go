// Package agentbench runs one real coding-agent benchmark trial.
package agentbench

import "fmt"

// Mode defines how AgentCap is exposed to the coding agent.
type Mode string

const (
	ModeDisabled   Mode = "disabled"
	ModeStateless  Mode = "stateless"
	ModeStateful   Mode = "stateful"
	ModeIntegrated Mode = "integrated"
)

func ParseMode(value string) (Mode, error) {
	mode := Mode(value)
	switch mode {
	case ModeDisabled, ModeStateless, ModeStateful, ModeIntegrated:
		return mode, nil
	default:
		return "", fmt.Errorf("unsupported agent benchmark mode %q", value)
	}
}

func (m Mode) AgentCapEnabled() bool { return m != ModeDisabled }
