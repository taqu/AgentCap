package reduce

import (
	"fmt"
	"strings"

	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/exec"
	"github.com/taqu/agentcap/internal/gitparse"
)

// GitBranchReducer handles git branch output.
type GitBranchReducer struct{ Args []string }

func (r *GitBranchReducer) Reduce(result *exec.Result) *ReducedResult {
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
	bl := gitparse.ParseBranch(cleaned)
	if bl == nil {
		return (&GenericReducer{}).Reduce(result)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap git-branch\n")
	if bl.Current != "" {
		fmt.Fprintf(&sb, "current=%s\n", bl.Current)
	}
	fmt.Fprintf(&sb, "local=%d remote=%d\n", len(bl.Local), len(bl.Remote))

	if len(bl.Local) > 0 {
		sb.WriteString("\nlocal:\n")
		for _, b := range bl.Local {
			if b.Current {
				fmt.Fprintf(&sb, "* %s\n", b.Name)
			} else {
				fmt.Fprintf(&sb, "  %s\n", b.Name)
			}
		}
	}

	if len(bl.Remote) > 0 {
		sb.WriteString("\nremote:\n")
		for _, b := range bl.Remote {
			fmt.Fprintf(&sb, "  %s\n", b.Name)
		}
	}

	sb.WriteString(buildStderr(result.Stderr))
	out := sb.String()
	return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
}
