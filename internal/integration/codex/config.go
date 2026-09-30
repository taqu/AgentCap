package codex

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	hookMarker         = "agentcap"
	instructionsFile   = ".codex/instructions.md"
	agentsFile         = "AGENTS.md"
	configFile         = ".codex/config.toml"
	defaultHookCommand = "acap hook codex"
)

const instructionsContent = `
<!-- agentcap -->
## AgentCap

AgentCap is active in this project. Shell command output is automatically compressed.

Compact results include an ID like ` + "`@acap 81bc2f`" + `. Use targeted drill-down:

- ` + "`acap show <id>`" + ` — full stored capsule
- ` + "`acap show <id> --errors`" + ` — build errors only
- ` + "`acap raw <id>`" + ` — full raw output (use as last resort)
`

const instructionsMarker = "<!-- agentcap -->"

// IsInstalled checks whether the AgentCap hook is present for Codex.
func IsInstalled(projectRoot string) bool {
	// Check instructions.md
	for _, name := range []string{instructionsFile, agentsFile} {
		instrPath := filepath.Join(projectRoot, name)
		if data, err := os.ReadFile(instrPath); err == nil {
			if strings.Contains(string(data), hookMarker) {
				return true
			}
		}
	}
	// Check config.toml
	cfgPath := filepath.Join(projectRoot, configFile)
	if data, err := os.ReadFile(cfgPath); err == nil {
		if strings.Contains(string(data), hookMarker) {
			return true
		}
	}
	return false
}

// Install adds AgentCap configuration for Codex CLI.
func Install(projectRoot string, dryRun bool) error {
	if dryRun {
		fmt.Printf("[dry-run] Would install AgentCap for Codex in %s\n", projectRoot)
		return nil
	}
	if err := InstallInstructions(projectRoot); err != nil {
		return err
	}
	return InstallHook(projectRoot, defaultHookCommand)
}

// InstallInstructions exposes AgentCap recovery guidance to Codex.
func InstallInstructions(projectRoot string) error {
	if err := os.MkdirAll(filepath.Join(projectRoot, ".codex"), 0o755); err != nil {
		return fmt.Errorf("codex: mkdir: %w", err)
	}
	if err := appendIfMissing(filepath.Join(projectRoot, instructionsFile), instructionsMarker, instructionsContent); err != nil {
		return fmt.Errorf("codex: instructions: %w", err)
	}
	if err := appendIfMissing(filepath.Join(projectRoot, agentsFile), instructionsMarker, instructionsContent); err != nil {
		return fmt.Errorf("codex: AGENTS.md: %w", err)
	}
	return nil
}

// InstallHook installs command interception without implicitly adding recovery
// instructions. hookCommand is quoted as a TOML string.
func InstallHook(projectRoot, hookCommand string) error {
	if hookCommand == "" {
		hookCommand = defaultHookCommand
	}
	if err := os.MkdirAll(filepath.Join(projectRoot, ".codex"), 0o755); err != nil {
		return fmt.Errorf("codex: mkdir: %w", err)
	}
	hookTOML := "\n# managed by agentcap\n[hooks]\npre_tool_use = [" + strconv.Quote(hookCommand) + "]\n"
	if err := appendIfMissing(filepath.Join(projectRoot, configFile), hookMarker, hookTOML); err != nil {
		return fmt.Errorf("codex: config: %w", err)
	}
	return nil
}

// Remove removes AgentCap-owned entries from Codex configuration.
func Remove(projectRoot string) error {
	if err := removeSection(filepath.Join(projectRoot, instructionsFile), instructionsMarker); err != nil {
		return err
	}
	if err := removeSection(filepath.Join(projectRoot, agentsFile), instructionsMarker); err != nil {
		return err
	}
	if err := removeSection(filepath.Join(projectRoot, configFile), hookMarker); err != nil {
		return err
	}
	return nil
}

// appendIfMissing appends content to path if marker is not already present.
func appendIfMissing(path, marker, content string) error {
	data, err := os.ReadFile(path)
	if err == nil && strings.Contains(string(data), marker) {
		return nil // already present
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}

// removeSection removes lines between the marker and the next section (empty line after marker block).
// Simple approach: remove all lines containing hookMarker and the lines following it until an empty line.
func removeSection(path, marker string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	content := string(data)
	if !strings.Contains(content, marker) {
		return nil
	}

	lines := strings.Split(content, "\n")
	var result []string
	skip := false
	for _, line := range lines {
		if strings.Contains(line, marker) {
			skip = true
			continue
		}
		if skip && strings.TrimSpace(line) == "" {
			skip = false
			continue
		}
		if !skip {
			result = append(result, line)
		}
	}

	out := strings.Join(result, "\n")
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(out), 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}
