package reduce

import (
	"fmt"
	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/exec"
	"strings"
)

// TreeReducer handles tree output.
type TreeReducer struct{}

func (t *TreeReducer) Reduce(r *exec.Result) *ReducedResult {
	raw := clean.StripANSI(r.Stdout)
	rawBytes := len(r.Stdout)

	lines := clean.Lines(raw)
	stderr := buildStderr(r.Stderr)

	// Small output: pass through.
	if len(raw) <= smallThresholdBytes && len(lines) <= smallThresholdLines {
		out := string(raw) + stderr
		return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
	}

	// Parse tree summary line (last non-empty line).
	dirs, files := parseTreeSummary(lines)

	// Collect lines up to depth 3.
	var kept []string
	truncated := false
	for _, line := range lines {
		if isTreeSummaryLine(line) {
			continue
		}
		depth := treeLineDepth(line)
		if depth <= 3 {
			kept = append(kept, line)
		} else {
			truncated = true
		}
	}

	var sb strings.Builder
	if dirs > 0 || files > 0 {
		fmt.Fprintf(&sb, "@acap tree dirs=%d files=%d\n\n", dirs, files)
	} else {
		fmt.Fprintf(&sb, "@acap tree lines=%d\n\n", len(lines))
	}

	sb.WriteString(strings.Join(kept, "\n"))
	if len(kept) > 0 {
		sb.WriteByte('\n')
	}

	if truncated {
		sb.WriteString("\ndepth_truncated=true\n")
	}

	sb.WriteString(stderr)
	out := sb.String()
	return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
}

// parseTreeSummary parses the summary line like "3 directories, 5 files".
func parseTreeSummary(lines []string) (dirs, files int) {
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		d, f, ok := parseTreeSummaryLine(line)
		if ok {
			return d, f
		}
		break
	}
	return 0, 0
}

func isTreeSummaryLine(line string) bool {
	_, _, ok := parseTreeSummaryLine(strings.TrimSpace(line))
	return ok
}

func parseTreeSummaryLine(line string) (dirs, files int, ok bool) {
	// Matches: "N director(y|ies), M file(s)"
	var d, f int
	n, err := fmt.Sscanf(line, "%d director", &d)
	if n != 1 || err != nil {
		return 0, 0, false
	}
	// Find ", " separator
	idx := strings.Index(line, ", ")
	if idx < 0 {
		return 0, 0, false
	}
	rest := line[idx+2:]
	n, err = fmt.Sscanf(rest, "%d file", &f)
	if n != 1 || err != nil {
		return 0, 0, false
	}
	return d, f, true
}

// treeLineDepth estimates the depth of a tree output line by counting
// the indentation levels. Tree uses UTF-8 box-drawing characters:
//
//	│   (U+2502 + 3 spaces)  — continuation
//	├── (U+251C + U+2500 + U+2500 + space) — branch
//	└── (U+2514 + U+2500 + U+2500 + space) — last branch
//	    (4 spaces) — empty continuation
//
// Each level contributes exactly 4 visible characters / rune-groups.
func treeLineDepth(line string) int {
	if line == "" || line == "." {
		return 0
	}
	depth := 0
	i := 0
	runes := []rune(line)
	for i < len(runes) {
		r := runes[i]
		if r == '│' { // U+2502 vertical bar
			// Expect 3 spaces after it
			if i+3 < len(runes) && runes[i+1] == ' ' && runes[i+2] == ' ' && runes[i+3] == ' ' {
				depth++
				i += 4
				continue
			}
			return depth
		} else if r == '├' || r == '└' { // branch chars U+251C / U+2514
			// Expect ── (two U+2500) + space
			if i+3 < len(runes) &&
				runes[i+1] == '─' &&
				runes[i+2] == '─' &&
				runes[i+3] == ' ' {
				depth++ // this char counts as one depth level
				i += 4
			}
			return depth
		} else if r == ' ' {
			// 4 spaces = one empty continuation level
			if i+3 < len(runes) && runes[i+1] == ' ' && runes[i+2] == ' ' && runes[i+3] == ' ' {
				depth++
				i += 4
				continue
			}
			return depth
		} else {
			return depth
		}
	}
	return depth
}
