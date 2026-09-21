package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindRoot_EnvOverride(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(EnvRoot, dir)
	got := FindRoot("/some/other/path")
	if got != filepath.Clean(dir) {
		t.Errorf("expected %s, got %s", dir, got)
	}
}

func TestFindRoot_AcapAncestor(t *testing.T) {
	root := t.TempDir()
	// Create .acap dir in root.
	if err := os.MkdirAll(filepath.Join(root, ".acap"), 0o700); err != nil {
		t.Fatal(err)
	}
	// Create a nested directory.
	nested := filepath.Join(root, "a", "b", "c")
	if err := os.MkdirAll(nested, 0o700); err != nil {
		t.Fatal(err)
	}

	got := FindRoot(nested)
	if got != root {
		t.Errorf("expected %s, got %s", root, got)
	}
}

func TestFindRoot_GitAncestor(t *testing.T) {
	root := t.TempDir()
	// Create .git dir in root.
	if err := os.MkdirAll(filepath.Join(root, ".git"), 0o700); err != nil {
		t.Fatal(err)
	}
	// Create nested dir (no .acap anywhere).
	nested := filepath.Join(root, "src", "pkg")
	if err := os.MkdirAll(nested, 0o700); err != nil {
		t.Fatal(err)
	}

	got := FindRoot(nested)
	if got != root {
		t.Errorf("expected %s, got %s", root, got)
	}
}

func TestFindRoot_NoGitFallback(t *testing.T) {
	dir := t.TempDir()
	// No .acap, no .git anywhere — should fall back to cwd.
	got := FindRoot(dir)
	if got != filepath.Clean(dir) {
		t.Errorf("expected %s, got %s", dir, got)
	}
}

func TestFindRoot_AcapBeatsGit(t *testing.T) {
	root := t.TempDir()
	// .git at top level.
	if err := os.MkdirAll(filepath.Join(root, ".git"), 0o700); err != nil {
		t.Fatal(err)
	}
	// .acap in subdirectory — should win because it's closer.
	sub := filepath.Join(root, "project")
	if err := os.MkdirAll(filepath.Join(sub, ".acap"), 0o700); err != nil {
		t.Fatal(err)
	}
	// Start from a dir nested inside sub.
	nested := filepath.Join(sub, "src")
	if err := os.MkdirAll(nested, 0o700); err != nil {
		t.Fatal(err)
	}

	got := FindRoot(nested)
	if got != sub {
		t.Errorf("expected %s (acap wins over git), got %s", sub, got)
	}
}

func TestFindRoot_Nested(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".acap"), 0o700); err != nil {
		t.Fatal(err)
	}
	deep := filepath.Join(root, "a", "b", "c", "d", "e")
	if err := os.MkdirAll(deep, 0o700); err != nil {
		t.Fatal(err)
	}
	got := FindRoot(deep)
	if got != root {
		t.Errorf("expected %s, got %s", root, got)
	}
}
