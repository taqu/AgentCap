package workload

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestPhase7WorkloadsLoad keeps the generated Phase 7 suite loadable and
// checks that setup states stay invisible while setup edits stay uncommitted.
func TestPhase7WorkloadsLoad(t *testing.T) {
	paths, _ := filepath.Glob("../../benchmarks/workloads/phase7/*.yaml")
	if len(paths) < 8 {
		t.Fatalf("expected Phase 7 workloads, found %d", len(paths))
	}
	for _, p := range paths {
		d, err := Load(p)
		if err != nil {
			t.Fatalf("%s: %v", p, err)
		}
		if d.Category == "" || d.Language == "" {
			t.Errorf("%s: missing category/language", p)
		}
		if len(d.Setup) == 0 {
			continue
		}
		ws, err := PrepareWorkspace(context.Background(), d, t.TempDir())
		if err != nil {
			t.Fatalf("%s: prepare: %v", p, err)
		}
		if _, err := os.Stat(filepath.Join(ws, ".states")); !os.IsNotExist(err) {
			t.Errorf("%s: .states visible to agent", p)
		}
		out, _ := exec.Command("git", "-C", ws, "status", "--porcelain").Output()
		if len(out) == 0 {
			t.Errorf("%s: setup produced no uncommitted change", p)
		}
	}
}
