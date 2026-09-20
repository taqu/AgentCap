package reduce

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/exec"
)

// DuReducer handles du output.
type DuReducer struct{}

type duEntry struct {
	size    int64  // in bytes (or blocks if not -h)
	sizeStr string // original size string (may be human-readable)
	path    string
}

func (d *DuReducer) Reduce(r *exec.Result) *ReducedResult {
	raw := clean.StripANSI(r.Stdout)
	rawBytes := len(r.Stdout)

	lines := clean.Lines(raw)
	stderr := buildStderr(r.Stderr)

	// Small output: pass through.
	if len(raw) <= smallThresholdBytes && len(lines) <= smallThresholdLines {
		out := string(raw) + stderr
		return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
	}

	entries := parseDuLines(lines)
	if len(entries) == 0 {
		g := &GenericReducer{}
		return g.Reduce(r)
	}

	// Sort by size descending.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].size > entries[j].size
	})

	// Find total (usually the last entry is the top-level total, or max).
	var total int64
	for _, e := range entries {
		if e.size > total {
			total = e.size
		}
	}

	totalStr := formatSize(total)
	omitted := max(0, len(entries)-10)

	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap du entries=%d total=%s\n", len(entries), totalStr)

	sb.WriteString("\nlargest:\n")
	top := entries
	if len(top) > 10 {
		top = top[:10]
	}
	for _, e := range top {
		sizeDisplay := e.sizeStr
		if sizeDisplay == "" {
			sizeDisplay = formatSize(e.size)
		}
		fmt.Fprintf(&sb, "%s %s\n", sizeDisplay, e.path)
	}

	if omitted > 0 {
		fmt.Fprintf(&sb, "\nomitted=%d\n", omitted)
	}

	sb.WriteString(stderr)
	out := sb.String()
	return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
}

func parseDuLines(lines []string) []duEntry {
	var entries []duEntry
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// du output: "SIZE\tPATH"
		idx := strings.IndexByte(line, '\t')
		if idx < 0 {
			// Fall back: split on whitespace.
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}
			size, sizeStr := parseDuSize(fields[0])
			path := strings.Join(fields[1:], " ")
			entries = append(entries, duEntry{size: size, sizeStr: sizeStr, path: path})
			continue
		}
		sizeField := strings.TrimSpace(line[:idx])
		path := strings.TrimSpace(line[idx+1:])
		size, sizeStr := parseDuSize(sizeField)
		entries = append(entries, duEntry{size: size, sizeStr: sizeStr, path: path})
	}
	return entries
}

// parseDuSize parses a du size field which may be a plain number (blocks) or
// a human-readable string like "4.8G", "1.1M", "420K".
func parseDuSize(s string) (int64, string) {
	if s == "" {
		return 0, s
	}
	// Try plain integer.
	n, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return n * 1024, "" // du default is 1K blocks
	}

	// Try human-readable suffix.
	if len(s) < 2 {
		return 0, s
	}
	suffix := s[len(s)-1]
	num, err := strconv.ParseFloat(s[:len(s)-1], 64)
	if err != nil {
		return 0, s
	}
	var mult float64
	switch suffix {
	case 'K', 'k':
		mult = 1024
	case 'M', 'm':
		mult = 1024 * 1024
	case 'G', 'g':
		mult = 1024 * 1024 * 1024
	case 'T', 't':
		mult = 1024 * 1024 * 1024 * 1024
	default:
		return 0, s
	}
	return int64(num * mult), s
}
