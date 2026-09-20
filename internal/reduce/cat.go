package reduce

import (
	"fmt"
	"strings"

	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/exec"
)

// CatReducer handles cat output.
type CatReducer struct{}

func (c *CatReducer) Reduce(r *exec.Result) *ReducedResult {
	raw := r.Stdout
	rawBytes := len(raw)
	stderr := buildStderr(r.Stderr)

	// Detect binary content (>15% non-printable bytes excluding newlines/tabs).
	if isBinary(raw) {
		out := fmt.Sprintf("@acap cat binary bytes=%d\noutput omitted\n", len(raw)) + stderr
		return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
	}

	cleaned := clean.StripANSI(raw)
	lines := clean.Lines(cleaned)

	// Extract filename from args if possible.
	filename := catFilename(r.Args)

	// Small output: pass through.
	if len(cleaned) <= smallThresholdBytes && len(lines) <= smallThresholdLines {
		out := string(cleaned) + stderr
		return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
	}

	const head = 80
	const tail = 80

	total := len(lines)
	omitted := max(0, total-head-tail)

	var sb strings.Builder
	if filename != "" {
		fmt.Fprintf(&sb, "@acap cat file=%s lines=%d bytes=%d\n\n", filename, total, len(cleaned))
	} else {
		fmt.Fprintf(&sb, "@acap cat lines=%d bytes=%d\n\n", total, len(cleaned))
	}

	if total <= head+tail {
		sb.WriteString(strings.Join(lines, "\n"))
		sb.WriteByte('\n')
	} else {
		sb.WriteString(strings.Join(lines[:head], "\n"))
		sb.WriteByte('\n')
		fmt.Fprintf(&sb, "\n...\n\n")
		sb.WriteString(strings.Join(lines[total-tail:], "\n"))
		sb.WriteByte('\n')
		fmt.Fprintf(&sb, "\nomitted_lines=%d\n", omitted)
	}

	sb.WriteString(stderr)
	out := sb.String()
	return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
}

// isBinary returns true if more than 15% of non-whitespace bytes are non-printable.
func isBinary(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	nonPrintable := 0
	total := 0
	for _, b := range data {
		if b == '\n' || b == '\r' || b == '\t' {
			continue
		}
		total++
		if b < 0x20 || b == 0x7f {
			nonPrintable++
		}
	}
	if total == 0 {
		return false
	}
	return nonPrintable*100/total > 15
}

// catFilename extracts the filename from cat's args (last non-flag argument).
func catFilename(args []string) string {
	if len(args) < 2 {
		return ""
	}
	// Return first non-flag argument after "cat".
	for _, a := range args[1:] {
		if !strings.HasPrefix(a, "-") {
			return a
		}
	}
	return ""
}
