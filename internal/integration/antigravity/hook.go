// Package antigravity provides the Antigravity run_command adapter for AgentCap.
// The same hook contract is shared by Antigravity 2.0, the CLI and the IDE.
package antigravity

import (
	"encoding/json"
	"fmt"
	"strings"
)

// HookInput holds only the PreToolUse fields AgentCap needs. Unknown fields
// are ignored; transcriptPath and artifactDirectoryPath are deliberately unused.
type HookInput struct {
	ConversationID string   `json:"conversationId"`
	WorkspacePaths []string `json:"workspacePaths"`
	ModelName      string   `json:"modelName"`
	StepIdx        *int64   `json:"stepIdx"`
	ToolCall       *struct {
		Name string          `json:"name"`
		Args json.RawMessage `json:"args"`
	} `json:"toolCall"`
}

// RunCommandArgs is the subset of run_command arguments AgentCap reads.
type RunCommandArgs struct {
	CommandLine       string `json:"CommandLine"`
	Cwd               string `json:"Cwd"`
	WaitMsBeforeAsync *int64 `json:"WaitMsBeforeAsync"`
}

// ParseHookInput parses Antigravity's PreToolUse stdin.
func ParseHookInput(data []byte) (*HookInput, error) {
	var h HookInput
	if err := json.Unmarshal(data, &h); err != nil {
		return nil, fmt.Errorf("antigravity: malformed hook payload: %w", err)
	}
	if h.ToolCall == nil || h.ToolCall.Name == "" {
		return nil, fmt.Errorf("antigravity: hook payload missing toolCall.name")
	}
	return &h, nil
}

// ParseRunCommandArgs parses run_command arguments. A missing CommandLine
// means the tool contract changed incompatibly.
func ParseRunCommandArgs(raw json.RawMessage) (*RunCommandArgs, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("antigravity: run_command missing args")
	}
	var a RunCommandArgs
	if err := json.Unmarshal(raw, &a); err != nil {
		return nil, fmt.Errorf("antigravity: invalid run_command args: %w", err)
	}
	if strings.TrimSpace(a.CommandLine) == "" {
		return nil, fmt.Errorf("antigravity: run_command missing CommandLine")
	}
	return &a, nil
}

type decision struct {
	Decision  string            `json:"decision"`
	Overwrite map[string]string `json:"overwrite,omitempty"`
}

// MakePassResponse defers to Antigravity's normal permission flow. An empty
// object is not neutral: Antigravity treats a missing decision as a denial.
// "ask" respects configured allow rules, cached grants and the autonomy mode.
func MakePassResponse() []byte {
	b, _ := json.Marshal(decision{Decision: "ask"})
	return b
}

func makeOverwriteResponse(commandLine string) ([]byte, error) {
	return json.Marshal(decision{Decision: "ask", Overwrite: map[string]string{"CommandLine": commandLine}})
}
