package reduce

import (
	"fmt"
	"strings"

	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/exec"
	"github.com/taqu/agentcap/internal/gitparse"
)

const maxShownCommits = 30

// GitLogReducer handles git log output.
type GitLogReducer struct{ Args []string }

func (r *GitLogReducer) Reduce(result *exec.Result) *ReducedResult {
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
	gl := gitparse.ParseLog(cleaned)
	if gl == nil || (gl.Total == 0 && len(gl.Commits) == 0) {
		return (&GenericReducer{}).Reduce(result)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap git-log\n")
	fmt.Fprintf(&sb, "total=%d\n", gl.Total)

	shown := gl.Commits
	omitted := 0
	if len(shown) > maxShownCommits {
		omitted = len(shown) - maxShownCommits
		shown = shown[:maxShownCommits]
	}

	sb.WriteByte('\n')
	for _, c := range shown {
		if c.Subject != "" {
			fmt.Fprintf(&sb, "%s %s\n", shortHash(c.Hash), c.Subject)
		} else {
			fmt.Fprintf(&sb, "%s\n", shortHash(c.Hash))
		}
	}

	if omitted > 0 {
		fmt.Fprintf(&sb, "\nomitted=%d\n", omitted)
	}

	sb.WriteString(buildStderr(result.Stderr))
	out := sb.String()
	return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
}

func shortHash(h string) string {
	if len(h) > 8 {
		return h[:8]
	}
	return h
}
