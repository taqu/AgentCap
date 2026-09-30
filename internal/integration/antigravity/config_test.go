package antigravity

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

const existingHooks = `{
  "lint-checker": {
    "PostToolUse": [{"matcher": "run_command", "hooks": [{"type": "command", "command": "./lint.sh", "timeout": 10}]}]
  },
  "safety-gate": {
    "enabled": false,
    "PreToolUse": [{"matcher": "run_command", "hooks": [{"command": "./safety-check.sh", "futureField": 1}]}]
  }
}`

func semantic(t *testing.T, data []byte) map[string]any {
	t.Helper()
	var v map[string]any
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatal(err)
	}
	return v
}

func TestInstallPreservesAndRemoves(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, hooksFile)
	os.MkdirAll(filepath.Dir(path), 0755)
	os.WriteFile(path, []byte(existingHooks), 0640)

	for range 2 {
		if err := Install(root, false); err != nil {
			t.Fatal(err)
		}
	}
	data, _ := os.ReadFile(path)
	got := semantic(t, data)
	if len(got) != 3 || !IsInstalled(root) {
		t.Fatalf("expected one AgentCap hook plus existing: %s", data)
	}
	spec := got[hookName].(map[string]any)["PreToolUse"].([]any)[0].(map[string]any)
	if spec["matcher"] != "run_command" || len(got[hookName].(map[string]any)) != 1 {
		t.Fatalf("only PreToolUse on run_command: %s", data)
	}
	if info, _ := os.Stat(path); info.Mode().Perm() != 0640 {
		t.Fatal("mode changed")
	}
	before, _ := os.ReadFile(path)
	Install(root, false)
	if after, _ := os.ReadFile(path); string(after) != string(before) {
		t.Fatal("repeat install rewrote file")
	}

	for range 2 {
		if err := Remove(root); err != nil {
			t.Fatal(err)
		}
	}
	data, _ = os.ReadFile(path)
	want, _ := json.Marshal(semantic(t, []byte(existingHooks)))
	gotJSON, _ := json.Marshal(semantic(t, data))
	if string(want) != string(gotJSON) || IsInstalled(root) {
		t.Fatalf("existing + install + remove != existing:\n%s", data)
	}
}

func TestInstallRefreshesStaleOwnedHook(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, hooksFile)
	os.MkdirAll(filepath.Dir(path), 0755)
	os.WriteFile(path, []byte(`{"agentcap":{"PreToolUse":[{"matcher":"*","hooks":[{"command":"old"}]}]}}`), 0644)
	Install(root, false)
	data, _ := os.ReadFile(path)
	if string(semantic(t, data)[hookName].(map[string]any)["PreToolUse"].([]any)[0].(map[string]any)["matcher"].(string)) != "run_command" {
		t.Fatal(string(data))
	}
}

func TestInstallCreatesAndRejectsMalformed(t *testing.T) {
	root := t.TempDir()
	if err := Install(root, false); err != nil || !IsInstalled(root) {
		t.Fatal(err)
	}
	if err := Remove(t.TempDir()); err != nil {
		t.Fatal("remove without file:", err)
	}
	bad := t.TempDir()
	path := filepath.Join(bad, hooksFile)
	os.MkdirAll(filepath.Dir(path), 0755)
	for _, content := range []string{`{`, `null`, `[]`} {
		os.WriteFile(path, []byte(content), 0644)
		if err := Install(bad, false); err == nil {
			t.Fatalf("%s: expected error", content)
		}
		if data, _ := os.ReadFile(path); string(data) != content {
			t.Fatal("malformed file overwritten")
		}
	}
}
