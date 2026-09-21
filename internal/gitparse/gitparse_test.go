package gitparse

import (
	"os"
	"testing"
)

// --- ParseStatus tests ---

func TestParseStatusClean(t *testing.T) {
	data, err := os.ReadFile("../../testdata/git/git-status-clean.txt")
	if err != nil {
		t.Fatal(err)
	}
	gs := ParseStatus(data)
	if gs == nil {
		t.Fatal("expected non-nil GitStatus")
	}
	if !gs.IsClean {
		t.Errorf("expected IsClean=true, got false")
	}
	if gs.Branch != "main" {
		t.Errorf("expected Branch=main, got %q", gs.Branch)
	}
	if len(gs.Files) != 0 {
		t.Errorf("expected 0 files, got %d", len(gs.Files))
	}
}

func TestParseStatusModified(t *testing.T) {
	data, err := os.ReadFile("../../testdata/git/git-status-modified.txt")
	if err != nil {
		t.Fatal(err)
	}
	gs := ParseStatus(data)
	if gs == nil {
		t.Fatal("expected non-nil GitStatus")
	}
	if gs.Branch != "features/phase4" {
		t.Errorf("expected Branch=features/phase4, got %q", gs.Branch)
	}
	if gs.Ahead != 2 {
		t.Errorf("expected Ahead=2, got %d", gs.Ahead)
	}
	if gs.IsClean {
		t.Errorf("expected IsClean=false")
	}
	if len(gs.Files) == 0 {
		t.Errorf("expected files, got none")
	}

	// Check for staged files.
	var staged []GitStatusFile
	for _, f := range gs.Files {
		if len(f.XY) == 2 && f.XY[0] != ' ' && !f.IsConflicted() && !f.IsUntracked() {
			staged = append(staged, f)
		}
	}
	if len(staged) != 2 {
		t.Errorf("expected 2 staged files, got %d", len(staged))
	}
}

func TestParseStatusShort(t *testing.T) {
	data, err := os.ReadFile("../../testdata/git/git-status-short.txt")
	if err != nil {
		t.Fatal(err)
	}
	gs := ParseStatus(data)
	if gs == nil {
		t.Fatal("expected non-nil GitStatus")
	}
	if gs.Branch != "main" {
		t.Errorf("expected Branch=main, got %q", gs.Branch)
	}
	if gs.Ahead != 1 {
		t.Errorf("expected Ahead=1, got %d", gs.Ahead)
	}
	if gs.Upstream != "origin/main" {
		t.Errorf("expected Upstream=origin/main, got %q", gs.Upstream)
	}
	if len(gs.Files) != 3 {
		t.Errorf("expected 3 files, got %d", len(gs.Files))
	}
}

func TestParseStatusDetached(t *testing.T) {
	raw := []byte("HEAD detached at 8fd31a2\nnothing to commit, working tree clean\n")
	gs := ParseStatus(raw)
	if gs == nil {
		t.Fatal("expected non-nil GitStatus")
	}
	if !gs.Detached {
		t.Errorf("expected Detached=true")
	}
	if gs.Branch != "DETACHED" {
		t.Errorf("expected Branch=DETACHED, got %q", gs.Branch)
	}
	if gs.HeadHash != "8fd31a2" {
		t.Errorf("expected HeadHash=8fd31a2, got %q", gs.HeadHash)
	}
}

func TestParseStatusEmpty(t *testing.T) {
	gs := ParseStatus([]byte{})
	if gs == nil {
		t.Fatal("expected non-nil GitStatus")
	}
	if !gs.IsClean {
		t.Errorf("expected IsClean=true for empty input")
	}
}

func TestGitStatusFileIsConflicted(t *testing.T) {
	cases := []struct {
		xy        string
		conflicts bool
	}{
		{"UU", true},
		{"AA", true},
		{"DD", true},
		{"AU", true},
		{"UD", true},
		{"M ", false},
		{" M", false},
		{"??", false},
	}
	for _, c := range cases {
		f := GitStatusFile{XY: c.xy}
		got := f.IsConflicted()
		if got != c.conflicts {
			t.Errorf("XY=%q IsConflicted()=%v, want %v", c.xy, got, c.conflicts)
		}
	}
}

// --- ParseDiff tests ---

func TestParseDiffSmall(t *testing.T) {
	data, err := os.ReadFile("../../testdata/git/git-diff-small.txt")
	if err != nil {
		t.Fatal(err)
	}
	d := ParseDiff(data)
	if d == nil {
		t.Fatal("expected non-nil GitDiff")
	}
	if len(d.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(d.Files))
	}
	// Check additions/deletions.
	add, del := d.Totals()
	if add == 0 {
		t.Errorf("expected additions > 0")
	}
	_ = del
}

func TestParseDiffBinary(t *testing.T) {
	data, err := os.ReadFile("../../testdata/git/git-diff-binary.txt")
	if err != nil {
		t.Fatal(err)
	}
	d := ParseDiff(data)
	if d == nil {
		t.Fatal("expected non-nil GitDiff")
	}
	if len(d.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(d.Files))
	}
	// First file should be binary.
	if !d.Files[0].Binary {
		t.Errorf("expected first file to be binary")
	}
	// Second file should not be binary.
	if d.Files[1].Binary {
		t.Errorf("expected second file to not be binary")
	}
}

func TestParseDiffRename(t *testing.T) {
	data, err := os.ReadFile("../../testdata/git/git-diff-rename.txt")
	if err != nil {
		t.Fatal(err)
	}
	d := ParseDiff(data)
	if d == nil {
		t.Fatal("expected non-nil GitDiff")
	}
	if len(d.Files) != 1 {
		t.Errorf("expected 1 file, got %d", len(d.Files))
	}
	f := d.Files[0]
	if f.Status != "R" {
		t.Errorf("expected Status=R, got %q", f.Status)
	}
	if f.Similarity != 85 {
		t.Errorf("expected Similarity=85, got %d", f.Similarity)
	}
	if f.OldPath != "internal/old/package.go" {
		t.Errorf("expected OldPath=internal/old/package.go, got %q", f.OldPath)
	}
	if f.NewPath != "internal/new/package.go" {
		t.Errorf("expected NewPath=internal/new/package.go, got %q", f.NewPath)
	}
}

func TestParseDiffEmpty(t *testing.T) {
	d := ParseDiff([]byte{})
	if d == nil {
		t.Fatal("expected non-nil GitDiff for empty input")
	}
	if len(d.Files) != 0 {
		t.Errorf("expected 0 files, got %d", len(d.Files))
	}
}

func TestParseDiffTooLarge(t *testing.T) {
	large := make([]byte, maxDiffBytes+1)
	d := ParseDiff(large)
	if d != nil {
		t.Errorf("expected nil for input exceeding size limit")
	}
}

func TestParseHunkHeader(t *testing.T) {
	cases := []struct {
		line     string
		oldStart int
		oldLines int
		newStart int
		newLines int
		header   string
	}{
		{"@@ -10,5 +10,7 @@ func Foo() {", 10, 5, 10, 7, "func Foo() {"},
		{"@@ -1 +1 @@", 1, 1, 1, 1, ""},
		{"@@ -0,0 +1,3 @@", 0, 0, 1, 3, ""},
		{"@@ -28,6 +28,9 @@ func Select(args []string) Reducer {", 28, 6, 28, 9, "func Select(args []string) Reducer {"},
	}
	for _, c := range cases {
		os, ol, ns, nl, hdr := parseHunkHeader(c.line)
		if os != c.oldStart || ol != c.oldLines || ns != c.newStart || nl != c.newLines || hdr != c.header {
			t.Errorf("parseHunkHeader(%q): got (%d,%d,%d,%d,%q), want (%d,%d,%d,%d,%q)",
				c.line, os, ol, ns, nl, hdr,
				c.oldStart, c.oldLines, c.newStart, c.newLines, c.header)
		}
	}
}
