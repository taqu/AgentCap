package codex

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallCreatesFiles(t *testing.T) {
	dir := t.TempDir()

	if err := Install(dir, false); err != nil {
		t.Fatalf("Install error: %v", err)
	}

	// Check instructions.md was created.
	instrPath := filepath.Join(dir, instructionsFile)
	data, err := os.ReadFile(instrPath)
	if err != nil {
		t.Fatalf("instructions.md not created: %v", err)
	}
	if !strings.Contains(string(data), hookMarker) {
		t.Errorf("instructions.md does not contain hookMarker: %s", string(data))
	}

	// Check config.toml was created.
	cfgPath := filepath.Join(dir, configFile)
	cfgData, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("config.toml not created: %v", err)
	}
	if !strings.Contains(string(cfgData), hookMarker) {
		t.Errorf("config.toml does not contain hookMarker: %s", string(cfgData))
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

func TestRemoveRemovesAgentcapSection(t *testing.T) {
	dir := t.TempDir()

	if err := Install(dir, false); err != nil {
		t.Fatalf("Install error: %v", err)
	}

	if err := Remove(dir); err != nil {
		t.Fatalf("Remove error: %v", err)
	}

	if IsInstalled(dir) {
		t.Error("IsInstalled: expected false after remove")
	}
}

func TestInstallIdempotent(t *testing.T) {
	dir := t.TempDir()

	if err := Install(dir, false); err != nil {
		t.Fatalf("first Install error: %v", err)
	}
	if err := Install(dir, false); err != nil {
		t.Fatalf("second Install error: %v", err)
	}

	// Count marker occurrences — should be 1.
	instrPath := filepath.Join(dir, instructionsFile)
	data, err := os.ReadFile(instrPath)
	if err != nil {
		t.Fatal(err)
	}
	count := strings.Count(string(data), instructionsMarker)
	if count != 1 {
		t.Errorf("expected 1 marker in instructions.md, found %d", count)
	}
}

func TestInstallDryRun(t *testing.T) {
	dir := t.TempDir()

	if err := Install(dir, true); err != nil {
		t.Fatalf("Install dry-run error: %v", err)
	}

	// Files should NOT be created in dry-run.
	if _, err := os.Stat(filepath.Join(dir, instructionsFile)); err == nil {
		t.Error("Install dry-run: instructions.md should not be created")
	}
}

func TestRemoveIdempotent(t *testing.T) {
	dir := t.TempDir()

	// Remove when not installed should not error.
	if err := Remove(dir); err != nil {
		t.Fatalf("Remove on empty dir error: %v", err)
	}
}
