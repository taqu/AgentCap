package reduce

import (
	"fmt"
	"strings"

	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/exec"
	"github.com/taqu/agentcap/internal/gitparse"
)

// GitShowReducer handles git show output.
type GitShowReducer struct{ Args []string }

func (r *GitShowReducer) Reduce(result *exec.Result) *ReducedResult {
	raw := result.Stdout
	rawBytes := len(raw)

	if result.ExitCode != 0 {
		out := string(clean.StripANSI(raw)) + buildStderr(result.Stderr)
		return &ReducedResult{
			Output:   out,
			RawBytes: rawBytes,
			RetBytes: rawBytes,
		}
	}

	cleaned := clean.StripANSI(raw)
	gs := gitparse.ParseShow(cleaned)
	if gs == nil {
		return (&GenericReducer{}).Reduce(result)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap git-show\n")

	c := gs.Commit
	if c.Hash != "" {
		fmt.Fprintf(&sb, "commit=%s\n", c.Hash)
	}
	if c.Author != "" {
		fmt.Fprintf(&sb, "author=%s\n", c.Author)
	}
	if c.Date != "" {
		fmt.Fprintf(&sb, "date=%s\n", c.Date)
	}
	if c.Subject != "" {
		fmt.Fprintf(&sb, "subject=%s\n", c.Subject)
	}

	if gs.Diff != nil && len(gs.Diff.Files) > 0 {
		totalAdd, totalDel := gs.Diff.Totals()
		fmt.Fprintf(&sb, "\nfiles=%d +%d -%d\n", len(gs.Diff.Files), totalAdd, totalDel)
		sb.WriteByte('\n')
		for _, f := range gs.Diff.Files {
			if f.Binary {
				fmt.Fprintf(&sb, "B %s\n", f.Path())
				continue
			}
			renderDiffFileLine(&sb, f)
		}
	}

	sb.WriteString(buildStderr(result.Stderr))
	out := sb.String()
	return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
}
