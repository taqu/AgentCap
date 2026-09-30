package antigravity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	hooksFile   = ".agents/hooks.json"
	hookName    = "agentcap"
	hookCommand = "acap hook antigravity"
)

type object map[string]json.RawMessage

// ownedSpec is the complete AgentCap-owned named hook. Only PreToolUse on
// run_command is registered; PostToolUse cannot replace output.
func ownedSpec() json.RawMessage {
	data, _ := json.Marshal(map[string]any{
		"PreToolUse": []any{map[string]any{
			"matcher": "run_command",
			"hooks":   []any{map[string]any{"type": "command", "command": hookCommand, "timeout": 10}},
		}},
	})
	return data
}

func readHooks(path string) (object, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return object{}, nil
	}
	if err != nil {
		return nil, err
	}
	var hooks object
	if err := json.Unmarshal(data, &hooks); err != nil {
		return nil, fmt.Errorf("antigravity: malformed %s: %w", hooksFile, err)
	}
	if hooks == nil {
		return nil, fmt.Errorf("antigravity: %s must be an object", hooksFile)
	}
	return hooks, nil
}

// IsInstalled reports whether the AgentCap-owned named hook exists.
func IsInstalled(projectRoot string) bool {
	hooks, err := readHooks(filepath.Join(projectRoot, hooksFile))
	if err != nil {
		return false
	}
	_, ok := hooks[hookName]
	return ok
}

// Install adds or refreshes the "agentcap" named hook, preserving all others.
func Install(projectRoot string, dryRun bool) error { return update(projectRoot, false, dryRun) }

// Remove deletes only the "agentcap" named hook.
func Remove(projectRoot string) error { return update(projectRoot, true, false) }

func update(root string, remove, dryRun bool) error {
	path := filepath.Join(root, hooksFile)
	hooks, err := readHooks(path)
	if err != nil {
		return err
	}
	current, present := hooks[hookName]
	switch {
	case remove && !present:
		return nil
	case remove:
		delete(hooks, hookName)
	default:
		var a, b bytes.Buffer
		if present && json.Compact(&a, current) == nil && json.Compact(&b, ownedSpec()) == nil && a.String() == b.String() {
			return nil
		}
		hooks[hookName] = ownedSpec()
	}
	data, err := json.MarshalIndent(hooks, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if dryRun {
		fmt.Printf("[dry-run] Would update %s:\n%s", path, data)
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	mode := os.FileMode(0644)
	if info, err := os.Stat(path); err == nil {
		mode = info.Mode().Perm()
	}
	temp, err := os.CreateTemp(filepath.Dir(path), ".agentcap-hooks-*")
	if err != nil {
		return err
	}
	defer os.Remove(temp.Name())
	if err = temp.Chmod(mode); err == nil {
		_, err = temp.Write(data)
	}
	closeErr := temp.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	return os.Rename(temp.Name(), path)
}
