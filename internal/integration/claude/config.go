package claude

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	settingsFile = ".claude/settings.json"
	hookMarker   = "agentcap"
	hookCommand  = "acap hook claude"
	hookMatcher  = "Bash"
)

type HookEntry struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}
type HookMatcher struct {
	Matcher string      `json:"matcher"`
	Hooks   []HookEntry `json:"hooks"`
}

type object map[string]json.RawMessage

func readSettings(path string) (object, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return object{}, nil
	}
	if err != nil {
		return nil, err
	}
	var settings object
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("claude: malformed settings: %w", err)
	}
	if settings == nil {
		return nil, fmt.Errorf("claude: settings must be an object")
	}
	return settings, nil
}

func owned(entry object) bool {
	var kind, command string
	json.Unmarshal(entry["type"], &kind)
	json.Unmarshal(entry["command"], &command)
	return kind == "command" && command == hookCommand
}

// editHooks preserves all unknown fields on unrelated matchers and handlers.
func editHooks(settings object, remove bool) (bool, error) {
	hooks := object{}
	if raw, ok := settings["hooks"]; ok {
		if err := json.Unmarshal(raw, &hooks); err != nil || hooks == nil {
			return false, fmt.Errorf("claude: hooks must be an object")
		}
	}
	var matchers []object
	if raw, ok := hooks["PreToolUse"]; ok {
		if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
			return false, fmt.Errorf("claude: PreToolUse must be an array")
		}
		if err := json.Unmarshal(raw, &matchers); err != nil {
			return false, fmt.Errorf("claude: invalid PreToolUse: %w", err)
		}
	}
	found := false
	changed := false
	retained := make([]object, 0, len(matchers)+1)
	for _, matcher := range matchers {
		if matcher == nil {
			return false, fmt.Errorf("claude: invalid hook matcher")
		}
		var handlers []object
		raw, ok := matcher["hooks"]
		if !ok {
			return false, fmt.Errorf("claude: hook matcher missing hooks")
		}
		if err := json.Unmarshal(raw, &handlers); err != nil || handlers == nil {
			return false, fmt.Errorf("claude: invalid hook handlers")
		}
		kept := make([]object, 0, len(handlers))
		modified := false
		for _, handler := range handlers {
			if owned(handler) {
				if remove || found {
					modified = true
					changed = true
					continue
				}
				found = true
			}
			kept = append(kept, handler)
		}
		if modified {
			if len(kept) == 0 {
				continue
			}
			matcher["hooks"], _ = json.Marshal(kept)
		}
		retained = append(retained, matcher)
	}
	if !remove && !found {
		data, _ := json.Marshal(HookMatcher{Matcher: hookMatcher, Hooks: []HookEntry{{Type: "command", Command: hookCommand}}})
		var matcher object
		json.Unmarshal(data, &matcher)
		retained = append(retained, matcher)
		changed = true
	}
	if !changed {
		return false, nil
	}
	if len(retained) == 0 {
		delete(hooks, "PreToolUse")
	} else {
		hooks["PreToolUse"], _ = json.Marshal(retained)
	}
	if len(hooks) == 0 {
		delete(settings, "hooks")
	} else {
		settings["hooks"], _ = json.Marshal(hooks)
	}
	return true, nil
}

func IsInstalled(projectRoot string) bool {
	settings, err := readSettings(filepath.Join(projectRoot, settingsFile))
	if err != nil {
		return false
	}
	var hooks object
	if json.Unmarshal(settings["hooks"], &hooks) != nil {
		return false
	}
	var matchers []object
	if json.Unmarshal(hooks["PreToolUse"], &matchers) != nil {
		return false
	}
	for _, matcher := range matchers {
		var handlers []object
		json.Unmarshal(matcher["hooks"], &handlers)
		for _, handler := range handlers {
			if owned(handler) {
				return true
			}
		}
	}
	return false
}

func Install(projectRoot string, dryRun bool) error { return update(projectRoot, false, dryRun) }
func Remove(projectRoot string) error               { return update(projectRoot, true, false) }

func update(root string, remove, dryRun bool) error {
	path := filepath.Join(root, settingsFile)
	settings, err := readSettings(path)
	if err != nil {
		return err
	}
	changed, err := editHooks(settings, remove)
	if err != nil || !changed {
		return err
	}
	data, err := json.MarshalIndent(settings, "", "  ")
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
	mode := os.FileMode(0600)
	if info, err := os.Stat(path); err == nil {
		mode = info.Mode().Perm()
	}
	temp, err := os.CreateTemp(filepath.Dir(path), ".agentcap-settings-*")
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
