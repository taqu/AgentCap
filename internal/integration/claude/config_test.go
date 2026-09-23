package claude

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallCreatesSettings(t *testing.T) {
	dir := t.TempDir()

	if err := Install(dir, false); err != nil {
		t.Fatalf("Install error: %v", err)
	}

	path := filepath.Join(dir, settingsFile)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("settings.json not created: %v", err)
	}

	if !strings.Contains(string(data), hookCommand) {
		t.Errorf("settings.json does not contain hookCommand %q: %s", hookCommand, string(data))
	}
}

func TestIsInstalled(t *testing.T) {
	dir := t.TempDir()

	if IsInstalled(dir) {
		t.Error("IsInstalled: expected false before install")
	}

	if err := Install(dir, false); err != nil {
		t.Fatalf("Install error: %v", err)
	}

	if !IsInstalled(dir) {
		t.Error("IsInstalled: expected true after install")
	}
}

func TestRemoveRemovesOnlyAgentcapEntry(t *testing.T) {
	dir := t.TempDir()

	// Create settings with agentcap hook and a user hook.
	claudeDir := filepath.Join(dir, ".claude")
	if err := os.MkdirAll(claudeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	initial := map[string]interface{}{
		"hooks": map[string]interface{}{
			"PreToolUse": []interface{}{
				map[string]interface{}{
					"matcher": "Bash",
					"hooks": []interface{}{
						map[string]interface{}{"type": "command", "command": "acap hook claude"},
					},
				},
				map[string]interface{}{
					"matcher": "Edit",
					"hooks": []interface{}{
						map[string]interface{}{"type": "command", "command": "my-other-tool"},
					},
				},
			},
		},
	}

	data, _ := json.MarshalIndent(initial, "", "  ")
	if err := os.WriteFile(filepath.Join(claudeDir, "settings.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Remove(dir); err != nil {
		t.Fatalf("Remove error: %v", err)
	}

	result, err := os.ReadFile(filepath.Join(claudeDir, "settings.json"))
	if err != nil {
		t.Fatal(err)
	}

	content := string(result)
	if strings.Contains(content, hookMarker) {
		t.Errorf("Remove: agentcap hook still present in settings: %s", content)
	}
	if !strings.Contains(content, "my-other-tool") {
		t.Errorf("Remove: other hook was incorrectly removed: %s", content)
	}
}

func TestInstallIdempotent(t *testing.T) {
	dir := t.TempDir()

	// Install twice — should result in only one hook entry.
	if err := Install(dir, false); err != nil {
		t.Fatalf("first Install error: %v", err)
	}
	if err := Install(dir, false); err != nil {
		t.Fatalf("second Install error: %v", err)
	}

	path := filepath.Join(dir, settingsFile)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	// Count occurrences of hookMarker.
	count := strings.Count(string(data), hookCommand)
	if count != 1 {
		t.Errorf("expected exactly 1 hook entry, found %d in: %s", count, string(data))
	}
}

func TestInstallDryRun(t *testing.T) {
	dir := t.TempDir()

	if err := Install(dir, true); err != nil {
		t.Fatalf("Install dry-run error: %v", err)
	}

	// settings.json should NOT be created in dry-run mode.
	path := filepath.Join(dir, settingsFile)
	if _, err := os.Stat(path); err == nil {
		t.Error("Install dry-run: settings.json should not be created")
	}
}

func TestRemoveIdempotent(t *testing.T) {
	dir := t.TempDir()

	// Remove when not installed should not error.
	if err := Remove(dir); err != nil {
		t.Fatalf("Remove on empty dir error: %v", err)
	}
}
