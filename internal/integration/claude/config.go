package claude

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	settingsFile = ".claude/settings.json"
	hookMarker   = "agentcap"
	hookCommand  = "acap hook claude"
	hookMatcher  = "Bash"
)

// HookEntry represents one hook in Claude Code settings.
type HookEntry struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}

// HookMatcher is one entry in the PreToolUse array.
type HookMatcher struct {
	Matcher string      `json:"matcher"`
	Hooks   []HookEntry `json:"hooks"`
}

// IsInstalled checks whether the AgentCap hook is present in the project's Claude Code settings.
func IsInstalled(projectRoot string) bool {
	path := filepath.Join(projectRoot, settingsFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	s := string(data)
	return strings.Contains(s, hookMarker) || strings.Contains(s, hookCommand)
}

// Install adds the AgentCap PreToolUse hook to .claude/settings.json.
// If dryRun is true, prints what would be done instead.
// Idempotent: if already installed, does nothing.
func Install(projectRoot string, dryRun bool) error {
	if IsInstalled(projectRoot) {
		if dryRun {
			fmt.Printf("[dry-run] AgentCap hook already installed in %s\n", filepath.Join(projectRoot, settingsFile))
		}
		return nil
	}

	path := filepath.Join(projectRoot, settingsFile)

	// Read existing settings or start with empty object.
	var rawSettings map[string]json.RawMessage
	data, err := os.ReadFile(path)
	if err == nil {
		if err2 := json.Unmarshal(data, &rawSettings); err2 != nil {
			rawSettings = make(map[string]json.RawMessage)
		}
	} else {
		rawSettings = make(map[string]json.RawMessage)
	}

	// Parse existing hooks section.
	var hooksMap map[string][]HookMatcher
	if raw, ok := rawSettings["hooks"]; ok {
		if err := json.Unmarshal(raw, &hooksMap); err != nil {
			hooksMap = make(map[string][]HookMatcher)
		}
	} else {
		hooksMap = make(map[string][]HookMatcher)
	}

	// Build the new hook entry.
	newEntry := HookEntry{
		Type:    "command",
		Command: hookCommand,
	}
	newMatcher := HookMatcher{
		Matcher: hookMatcher,
		Hooks:   []HookEntry{newEntry},
	}

	// Add to PreToolUse, avoiding duplicates.
	preToolUse := hooksMap["PreToolUse"]
	for _, m := range preToolUse {
		for _, h := range m.Hooks {
			if strings.Contains(h.Command, hookMarker) || h.Command == hookCommand {
				// Already present.
				return nil
			}
		}
	}
	preToolUse = append(preToolUse, newMatcher)
	hooksMap["PreToolUse"] = preToolUse

	// Serialize hooks back.
	hooksJSON, err := json.MarshalIndent(hooksMap, "", "  ")
	if err != nil {
		return fmt.Errorf("claude: marshal hooks: %w", err)
	}
	rawSettings["hooks"] = json.RawMessage(hooksJSON)

	// Serialize full settings.
	out, err := json.MarshalIndent(rawSettings, "", "  ")
	if err != nil {
		return fmt.Errorf("claude: marshal settings: %w", err)
	}

	if dryRun {
		fmt.Printf("[dry-run] Would write to %s:\n%s\n", path, string(out))
		return nil
	}

	// Ensure parent dir exists.
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("claude: mkdir: %w", err)
	}

	// Atomic write: temp file + rename.
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, out, 0o644); err != nil {
		return fmt.Errorf("claude: write temp: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("claude: rename: %w", err)
	}

	// Install CLAUDE.md instruction.
	if err := installClaudeMD(projectRoot, dryRun); err != nil {
		// Non-fatal.
		fmt.Fprintf(os.Stderr, "acap: warning: CLAUDE.md: %v\n", err)
	}

	return nil
}

// Remove removes the AgentCap hook from .claude/settings.json.
// Only removes entries it owns (containing hookMarker).
// Idempotent: if not installed, does nothing.
func Remove(projectRoot string) error {
	path := filepath.Join(projectRoot, settingsFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("claude: read settings: %w", err)
	}

	var rawSettings map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawSettings); err != nil {
		return fmt.Errorf("claude: parse settings: %w", err)
	}

	hooksRaw, ok := rawSettings["hooks"]
	if !ok {
		return nil
	}

	var hooksMap map[string][]HookMatcher
	if err := json.Unmarshal(hooksRaw, &hooksMap); err != nil {
		return nil
	}

	// Remove agentcap-owned entries from PreToolUse.
	preToolUse := hooksMap["PreToolUse"]
	var filtered []HookMatcher
	for _, m := range preToolUse {
		var filteredHooks []HookEntry
		for _, h := range m.Hooks {
			if !strings.Contains(h.Command, hookMarker) {
				filteredHooks = append(filteredHooks, h)
			}
		}
		if len(filteredHooks) > 0 {
			filtered = append(filtered, HookMatcher{
				Matcher: m.Matcher,
				Hooks:   filteredHooks,
			})
		}
	}

	if len(filtered) == 0 {
		delete(hooksMap, "PreToolUse")
	} else {
		hooksMap["PreToolUse"] = filtered
	}

	if len(hooksMap) == 0 {
		delete(rawSettings, "hooks")
	} else {
		hooksJSON, err := json.MarshalIndent(hooksMap, "", "  ")
		if err != nil {
			return fmt.Errorf("claude: marshal hooks: %w", err)
		}
		rawSettings["hooks"] = json.RawMessage(hooksJSON)
	}

	out, err := json.MarshalIndent(rawSettings, "", "  ")
	if err != nil {
		return fmt.Errorf("claude: marshal settings: %w", err)
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, out, 0o644); err != nil {
		return fmt.Errorf("claude: write temp: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("claude: rename: %w", err)
	}

	return nil
}

const claudeMDMarker = "<!-- agentcap -->"

const claudeMDContent = `
<!-- agentcap -->
## AgentCap

AgentCap is active in this project. Shell command output is automatically compressed.

Compact results include an ID like ` + "`@acap 81bc2f`" + `. Use targeted drill-down before fetching raw output:

- ` + "`acap show <id>`" + ` — full stored capsule
- ` + "`acap show <id> --errors`" + ` — build errors only
- ` + "`acap show <id> --warnings`" + ` — build warnings only
- ` + "`acap show <id> --test <name>`" + ` — specific test failure
- ` + "`acap show <id> --file <path>`" + ` — specific file diff
- ` + "`acap show <id> --hunk <N>`" + ` — specific diff hunk
- ` + "`acap show <id> --match <text>`" + ` — search stored output
- ` + "`acap show <id> --lines X:Y`" + ` — line range
- ` + "`acap raw <id>`" + ` — full raw output (use as last resort)

Prefer ` + "`acap show`" + ` with a selector over ` + "`acap raw`" + `.
`

func installClaudeMD(projectRoot string, dryRun bool) error {
	path := filepath.Join(projectRoot, "CLAUDE.md")

	data, err := os.ReadFile(path)
	if err == nil && strings.Contains(string(data), claudeMDMarker) {
		return nil // already present
	}

	if dryRun {
		fmt.Printf("[dry-run] Would append AgentCap instructions to %s\n", path)
		return nil
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open CLAUDE.md: %w", err)
	}
	defer f.Close()

	_, err = f.WriteString(claudeMDContent)
	return err
}
