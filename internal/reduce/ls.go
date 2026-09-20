package reduce

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/exec"
)

// LsReducer handles ls output.
type LsReducer struct{}

func (l *LsReducer) Reduce(r *exec.Result) *ReducedResult {
	raw := clean.StripANSI(r.Stdout)
	rawBytes := len(r.Stdout)

	lines := clean.Lines(raw)
	stderr := buildStderr(r.Stderr)

	// Small output: pass through.
	if len(raw) <= smallThresholdBytes && len(lines) <= smallThresholdLines {
		out := string(raw) + stderr
		return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
	}

	// Detect -l format by checking if any non-total line starts with permission chars.
	longFormat := isLongFormat(lines)

	var out string
	if longFormat {
		out = reduceLsLong(lines, stderr)
	} else {
		out = reduceLsPlain(lines, stderr)
	}

	return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
}

func isLongFormat(lines []string) bool {
	for _, line := range lines {
		if len(line) < 4 {
			continue
		}
		c := line[0]
		if c == '-' || c == 'd' || c == 'l' || c == 'c' || c == 'b' || c == 'p' || c == 's' {
			// Likely a permissions line.
			return true
		}
	}
	return false
}

type lsEntry struct {
	name  string
	isDir bool
	size  int64 // bytes, 0 if unknown
}

func parseLsLongLine(line string) (lsEntry, bool) {
	// Format: permissions links owner group size month day time name
	// e.g.:  drwxr-xr-x  2 user group 4096 Jan  1 12:00 dirname
	fields := strings.Fields(line)
	if len(fields) < 9 {
		return lsEntry{}, false
	}
	mode := fields[0]
	if len(mode) < 1 {
		return lsEntry{}, false
	}
	isDir := mode[0] == 'd'
	// Size is field index 4
	size, err := strconv.ParseInt(fields[4], 10, 64)
	if err != nil {
		size = 0
	}
	// Name is everything after the 8th field (handles spaces)
	name := strings.Join(fields[8:], " ")
	return lsEntry{name: name, isDir: isDir, size: size}, true
}

func reduceLsLong(lines []string, stderr string) string {
	var dirs []string
	var files []lsEntry

	for _, line := range lines {
		if strings.HasPrefix(line, "total ") || line == "" {
			continue
		}
		entry, ok := parseLsLongLine(line)
		if !ok {
			continue
		}
		if entry.isDir {
			dirs = append(dirs, entry.name+"/")
		} else {
			files = append(files, entry)
		}
	}

	// Sort files by size descending.
	sort.Slice(files, func(i, j int) bool {
		return files[i].size > files[j].size
	})

	total := len(dirs) + len(files)
	shown := len(dirs) + min(5, len(files))
	omitted := total - shown

	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap ls entries=%d files=%d dirs=%d\n", total, len(files), len(dirs))

	if len(dirs) > 0 {
		sb.WriteString("\ndirs:\n")
		for _, d := range dirs {
			sb.WriteString(d)
			sb.WriteByte('\n')
		}
	}

	if len(files) > 0 {
		sb.WriteString("\nlargest:\n")
		top := files
		if len(top) > 5 {
			top = top[:5]
		}
		for _, f := range top {
			fmt.Fprintf(&sb, "%s %s\n", f.name, formatSize(f.size))
		}
	}

	if omitted > 0 {
		fmt.Fprintf(&sb, "\nomitted=%d\n", omitted)
	}

	sb.WriteString(stderr)
	return sb.String()
}

func reduceLsPlain(lines []string, stderr string) string {
	var dirs []string
	var files []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// In plain ls, we can't easily distinguish dirs from files without stat.
		// Show all entries, count them.
		if strings.HasSuffix(line, "/") {
			dirs = append(dirs, line)
		} else {
			files = append(files, line)
		}
	}

	total := len(dirs) + len(files)
	shown := len(dirs)
	omitted := total - shown

	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap ls entries=%d\n", total)

	if len(dirs) > 0 {
		sb.WriteString("\ndirs:\n")
		for _, d := range dirs {
			sb.WriteString(d)
			sb.WriteByte('\n')
		}
	}

	if omitted > 0 {
		fmt.Fprintf(&sb, "\nomitted=%d\n", omitted)
	}

	sb.WriteString(stderr)
	return sb.String()
}

func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1fGB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1fMB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1fKB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%dB", bytes)
	}
}
