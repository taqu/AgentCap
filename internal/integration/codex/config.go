package codex

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	hookMarker       = "agentcap"
	instructionsFile = ".codex/instructions.md"
	configFile       = ".codex/config.toml"
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

const hookTOML = `
# managed by agentcap
[hooks]
pre_tool_use = ["acap hook codex"]
`

// IsInstalled checks whether the AgentCap hook is present for Codex.
func IsInstalled(projectRoot string) bool {
	// Check instructions.md
	instrPath := filepath.Join(projectRoot, instructionsFile)
	if data, err := os.ReadFile(instrPath); err == nil {
		if strings.Contains(string(data), hookMarker) {
			return true
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
	if IsInstalled(projectRoot) {
		if dryRun {
			fmt.Printf("[dry-run] AgentCap already installed for Codex in %s\n", projectRoot)
		}
		return nil
	}

	if dryRun {
		fmt.Printf("[dry-run] Would install AgentCap for Codex in %s\n", projectRoot)
		return nil
	}

	codexDir := filepath.Join(projectRoot, ".codex")
	if err := os.MkdirAll(codexDir, 0o755); err != nil {
		return fmt.Errorf("codex: mkdir: %w", err)
	}

	// Write instructions.md
	if err := appendIfMissing(filepath.Join(projectRoot, instructionsFile), instructionsMarker, instructionsContent); err != nil {
		return fmt.Errorf("codex: instructions: %w", err)
	}

	// Write config.toml
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
