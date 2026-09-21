package reduce

import (
	"fmt"
	"strings"

	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/exec"
	"github.com/taqu/agentcap/internal/gitparse"
)

// GitStatusReducer handles git status output.
type GitStatusReducer struct{ Args []string }

func (r *GitStatusReducer) Reduce(result *exec.Result) *ReducedResult {
	raw := result.Stdout
	rawBytes := len(raw)

	// On failure, pass through with stderr.
	if result.ExitCode != 0 {
		out := string(clean.StripANSI(raw)) + buildStderr(result.Stderr)
		return &ReducedResult{
			Output:   out,
			RawBytes: rawBytes,
			RetBytes: rawBytes,
		}
	}

	cleaned := clean.StripANSI(raw)
	gs := gitparse.ParseStatus(cleaned)
	if gs == nil {
		return (&GenericReducer{}).Reduce(result)
	}

	var sb strings.Builder
	renderGitStatus(&sb, gs)
	sb.WriteString(buildStderr(result.Stderr))
	out := sb.String()
	return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
}

func renderGitStatus(sb *strings.Builder, gs *gitparse.GitStatus) {
	fmt.Fprintf(sb, "@acap git-status\n")
	if gs.Detached {
		fmt.Fprintf(sb, "branch=DETACHED head=%s\n", gs.HeadHash)
	} else if gs.Branch != "" {
		if gs.Ahead > 0 && gs.Behind > 0 {
			fmt.Fprintf(sb, "branch=%s ahead=%d behind=%d\n", gs.Branch, gs.Ahead, gs.Behind)
		} else if gs.Ahead > 0 {
			fmt.Fprintf(sb, "branch=%s ahead=%d\n", gs.Branch, gs.Ahead)
		} else if gs.Behind > 0 {
			fmt.Fprintf(sb, "branch=%s behind=%d\n", gs.Branch, gs.Behind)
		} else {
			fmt.Fprintf(sb, "branch=%s\n", gs.Branch)
		}
	}

	if gs.IsClean && len(gs.Files) == 0 {
		sb.WriteString("clean\n")
		return
	}

	// Collect file categories
	var conflicts, staged, unstaged, untracked []gitparse.GitStatusFile
	for _, f := range gs.Files {
		if f.IsConflicted() {
			conflicts = append(conflicts, f)
		} else if f.IsUntracked() {
			untracked = append(untracked, f)
		} else if len(f.XY) == 2 && f.XY[0] != ' ' {
			staged = append(staged, f)
		} else {
			unstaged = append(unstaged, f)
		}
	}

	// Conflicts first (most important)
	if len(conflicts) > 0 {
		fmt.Fprintf(sb, "\nconflicts=%d\n", len(conflicts))
		for _, f := range conflicts {
			fmt.Fprintf(sb, "%s %s\n", f.XY, f.Path)
		}
	}
	if len(staged) > 0 {
		fmt.Fprintf(sb, "\nstaged=%d\n", len(staged))
		for _, f := range staged {
			printFileEntry(sb, f)
		}
	}
	if len(unstaged) > 0 {
		fmt.Fprintf(sb, "\nunstaged=%d\n", len(unstaged))
		for _, f := range unstaged {
			printFileEntry(sb, f)
		}
	}
	if len(untracked) > 0 {
		fmt.Fprintf(sb, "\nuntracked=%d\n", len(untracked))
		for _, f := range untracked {
			fmt.Fprintf(sb, "? %s\n", f.Path)
		}
	}
}

func printFileEntry(sb *strings.Builder, f gitparse.GitStatusFile) {
	xy := strings.TrimSpace(f.XY)
	if f.OldPath != "" {
		fmt.Fprintf(sb, "%s %s -> %s\n", xy, f.OldPath, f.Path)
	} else {
		fmt.Fprintf(sb, "%s %s\n", xy, f.Path)
	}
}
