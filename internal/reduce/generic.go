package reduce

import (
	"fmt"
	"strings"

	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/exec"
)

// GenericReducer handles unknown commands with conservative reduction.
type GenericReducer struct{}

func (g *GenericReducer) Reduce(r *exec.Result) *ReducedResult {
	raw := clean.StripANSI(r.Stdout)
	raw = clean.StripProgress(raw)
	rawBytes := len(r.Stdout)

	lines := clean.Lines(raw)
	lines = clean.CollapseRepeated(lines)

	stderr := buildStderr(r.Stderr)

	// Small output: pass through unchanged.
	if len(raw) <= smallThresholdBytes && len(lines) <= smallThresholdLines {
		out := string(raw) + stderr
		return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
	}

	const head = 60
	const tail = 60

	total := len(lines)
	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap generic lines=%d shown=%d omitted=%d\n\n", total, min(head+tail, total), max(0, total-head-tail))

	if total <= head+tail {
		sb.WriteString(strings.Join(lines, "\n"))
		if total > 0 {
			sb.WriteByte('\n')
		}
	} else {
		sb.WriteString(strings.Join(lines[:head], "\n"))
		sb.WriteByte('\n')
		fmt.Fprintf(&sb, "\n...\n\n")
		sb.WriteString(strings.Join(lines[total-tail:], "\n"))
		sb.WriteByte('\n')
	}

	sb.WriteString(stderr)
	out := sb.String()
	return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
}

// buildStderr strips ANSI from stderr and formats it as a section.
func buildStderr(stderr []byte) string {
	if len(stderr) == 0 {
		return ""
	}
	cleaned := clean.StripANSI(stderr)
	cleaned = clean.StripProgress(cleaned)
	if len(cleaned) == 0 {
		return ""
	}
	return "\nstderr:\n" + string(cleaned)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
