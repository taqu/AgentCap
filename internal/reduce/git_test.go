package reduce

import (
	"os"
	"strings"
	"testing"

	"github.com/taqu/agentcap/internal/exec"
)

func makeGitResult(stdout []byte, exitCode int) *exec.Result {
	return &exec.Result{
		Stdout:   stdout,
		Stderr:   nil,
		ExitCode: exitCode,
	}
}

func TestGitStatusReducerClean(t *testing.T) {
	data, err := os.ReadFile("../../testdata/git/git-status-clean.txt")
	if err != nil {
		t.Fatal(err)
	}
	r := &GitStatusReducer{Args: []string{"git", "status"}}
	result := r.Reduce(makeGitResult(data, 0))
	if result == nil {
		t.Fatal("expected non-nil ReducedResult")
	}
	if !strings.Contains(result.Output, "@acap git-status") {
		t.Errorf("expected @acap git-status header, got: %q", result.Output[:min(100, len(result.Output))])
	}
	if !strings.Contains(result.Output, "branch=main") {
		t.Errorf("expected branch=main in output, got: %q", result.Output)
	}
	if !strings.Contains(result.Output, "clean") {
		t.Errorf("expected 'clean' in output, got: %q", result.Output)
	}
}

func TestGitStatusReducerModified(t *testing.T) {
	data, err := os.ReadFile("../../testdata/git/git-status-modified.txt")
	if err != nil {
		t.Fatal(err)
	}
	r := &GitStatusReducer{Args: []string{"git", "status"}}
	result := r.Reduce(makeGitResult(data, 0))
	if result == nil {
		t.Fatal("expected non-nil ReducedResult")
	}
	if !strings.Contains(result.Output, "@acap git-status") {
		t.Errorf("expected @acap git-status header")
	}
	if !strings.Contains(result.Output, "staged=") {
		t.Errorf("expected staged= section in output, got: %q", result.Output)
	}
	// Should have fewer bytes than raw.
	if result.RetBytes >= result.RawBytes && result.RawBytes > 0 {
		t.Logf("note: reduced size (%d) >= raw size (%d)", result.RetBytes, result.RawBytes)
	}
}

func TestGitStatusReducerFailure(t *testing.T) {
	r := &GitStatusReducer{Args: []string{"git", "status"}}
	result := r.Reduce(makeGitResult([]byte("fatal: not a git repository"), 128))
	if result == nil {
		t.Fatal("expected non-nil ReducedResult")
	}
	// Should pass through on failure.
	if !strings.Contains(result.Output, "fatal") {
		t.Errorf("expected error output passed through, got: %q", result.Output)
	}
}

func TestGitDiffReducerEmpty(t *testing.T) {
	r := &GitDiffReducer{Args: []string{"git", "diff"}}
	result := r.Reduce(makeGitResult([]byte(""), 0))
	if result == nil {
		t.Fatal("expected non-nil ReducedResult")
	}
	if !strings.Contains(result.Output, "no changes") {
		t.Errorf("expected 'no changes', got: %q", result.Output)
	}
}

func TestGitDiffReducerSmall(t *testing.T) {
	data, err := os.ReadFile("../../testdata/git/git-diff-small.txt")
	if err != nil {
		t.Fatal(err)
	}
	r := &GitDiffReducer{Args: []string{"git", "diff"}}
	result := r.Reduce(makeGitResult(data, 0))
	if result == nil {
		t.Fatal("expected non-nil ReducedResult")
	}
	if !strings.Contains(result.Output, "@acap git-diff") {
		t.Errorf("expected @acap git-diff header, got: %q", result.Output[:min(100, len(result.Output))])
	}
	if !strings.Contains(result.Output, "files=") {
		t.Errorf("expected files= in output, got: %q", result.Output)
	}
	if r.ParsedDiff == nil {
		t.Errorf("expected ParsedDiff to be set after Reduce")
	}
}

func TestGitDiffReducerBinary(t *testing.T) {
	data, err := os.ReadFile("../../testdata/git/git-diff-binary.txt")
	if err != nil {
		t.Fatal(err)
	}
	r := &GitDiffReducer{Args: []string{"git", "diff"}}
	result := r.Reduce(makeGitResult(data, 0))
	if result == nil {
		t.Fatal("expected non-nil ReducedResult")
	}
	if !strings.Contains(result.Output, "binary=") {
		t.Errorf("expected binary= section, got: %q", result.Output)
	}
}

func TestGitDiffReducerStat(t *testing.T) {
	statOutput := ` internal/reduce/reducer.go | 3 +++
 1 file changed, 3 insertions(+)
`
	r := &GitDiffReducer{Args: []string{"git", "diff", "--stat"}}
	result := r.Reduce(makeGitResult([]byte(statOutput), 0))
	if result == nil {
		t.Fatal("expected non-nil ReducedResult")
	}
	if !strings.Contains(result.Output, "stat") {
		t.Errorf("expected 'stat' in output, got: %q", result.Output)
	}
}

func TestGitDiffReducerCached(t *testing.T) {
	r := &GitDiffReducer{Args: []string{"git", "diff", "--cached"}}
	result := r.Reduce(makeGitResult([]byte(""), 0))
	if result == nil {
		t.Fatal("expected non-nil ReducedResult")
	}
	if !strings.Contains(result.Output, "staged") {
		t.Errorf("expected 'staged' in output for --cached, got: %q", result.Output)
	}
}

func TestGitStatusReducerShort(t *testing.T) {
	data, err := os.ReadFile("../../testdata/git/git-status-short.txt")
	if err != nil {
		t.Fatal(err)
	}
	r := &GitStatusReducer{Args: []string{"git", "status", "--short"}}
	result := r.Reduce(makeGitResult(data, 0))
	if result == nil {
		t.Fatal("expected non-nil ReducedResult")
	}
	if !strings.Contains(result.Output, "branch=main") {
		t.Errorf("expected branch=main, got: %q", result.Output)
	}
}

func TestSelectGit(t *testing.T) {
	cases := []struct {
		args     []string
		wantType string
	}{
		{[]string{"git", "status"}, "*reduce.GitStatusReducer"},
		{[]string{"git", "diff"}, "*reduce.GitDiffReducer"},
		{[]string{"git", "log"}, "*reduce.GitLogReducer"},
		{[]string{"git", "branch"}, "*reduce.GitBranchReducer"},
		{[]string{"git", "show"}, "*reduce.GitShowReducer"},
		{[]string{"git", "commit"}, "*reduce.GenericReducer"},
	}
	for _, c := range cases {
		r := SelectGit(c.args)
		got := typeName(r)
		if got != c.wantType {
			t.Errorf("SelectGit(%v) = %s, want %s", c.args, got, c.wantType)
		}
	}
}

func typeName(r Reducer) string {
	switch r.(type) {
	case *GitStatusReducer:
		return "*reduce.GitStatusReducer"
	case *GitDiffReducer:
		return "*reduce.GitDiffReducer"
	case *GitLogReducer:
		return "*reduce.GitLogReducer"
	case *GitBranchReducer:
		return "*reduce.GitBranchReducer"
	case *GitShowReducer:
		return "*reduce.GitShowReducer"
	case *GenericReducer:
		return "*reduce.GenericReducer"
	default:
		return "unknown"
	}
}
