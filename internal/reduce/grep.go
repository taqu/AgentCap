package reduce

import (
	"fmt"
	"sort"
	"strings"

	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/exec"
)

// GrepReducer handles grep and rg output.
type GrepReducer struct{}

type grepMatch struct {
	file    string
	lineNum string
	content string
}

func (g *GrepReducer) Reduce(r *exec.Result) *ReducedResult {
	raw := clean.StripANSI(r.Stdout)
	rawBytes := len(r.Stdout)

	lines := clean.Lines(raw)
	stderr := buildStderr(r.Stderr)

	// Small output: pass through.
	if len(raw) <= smallThresholdBytes && len(lines) <= smallThresholdLines {
		out := string(raw) + stderr
		return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
	}

	// Parse matches.
	fileCount := make(map[string]int) // file -> match count
	var fileOrder []string            // preserve first-seen order
	var sample []grepMatch

	for _, line := range lines {
		if line == "" {
			continue
		}
		m, ok := parseGrepLine(line)
		if !ok {
			continue
		}
		if _, seen := fileCount[m.file]; !seen {
			fileOrder = append(fileOrder, m.file)
		}
		fileCount[m.file]++
		if len(sample) < 5 {
			sample = append(sample, m)
		}
	}

	total := len(lines)

	// If we couldn't parse anything meaningful, fall back.
	if len(fileCount) == 0 {
		gr := &GenericReducer{}
		return gr.Reduce(r)
	}

	// Sort files by match count descending.
	sortedFiles := make([]string, len(fileOrder))
	copy(sortedFiles, fileOrder)
	sort.Slice(sortedFiles, func(i, j int) bool {
		if fileCount[sortedFiles[i]] != fileCount[sortedFiles[j]] {
			return fileCount[sortedFiles[i]] > fileCount[sortedFiles[j]]
		}
		return sortedFiles[i] < sortedFiles[j]
	})

	omitted := total - len(sample)

	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap rg matches=%d files=%d\n", total, len(fileCount))

	sb.WriteString("\nfiles:\n")
	shown := 0
	for _, f := range sortedFiles {
		if shown >= 20 {
			break
		}
		fmt.Fprintf(&sb, "%s %d\n", f, fileCount[f])
		shown++
	}

	if len(sample) > 0 {
		sb.WriteString("\nsample:\n")
		for _, m := range sample {
			if m.lineNum != "" {
				fmt.Fprintf(&sb, "%s:%s:%s\n", m.file, m.lineNum, m.content)
			} else {
				fmt.Fprintf(&sb, "%s:%s\n", m.file, m.content)
			}
		}
	}

	if omitted > 0 {
		fmt.Fprintf(&sb, "\nomitted=%d\n", omitted)
	}

	sb.WriteString(stderr)
	out := sb.String()
	return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
}

// parseGrepLine parses a line of the form "file:line:content" or "file:content".
func parseGrepLine(line string) (grepMatch, bool) {
	// Find first colon.
	idx1 := strings.IndexByte(line, ':')
	if idx1 < 0 {
		return grepMatch{}, false
	}
	file := line[:idx1]
	rest := line[idx1+1:]

	if file == "" {
		return grepMatch{}, false
	}

	// Try to find a line number (numeric field before second colon).
	idx2 := strings.IndexByte(rest, ':')
	if idx2 >= 0 {
		possibleNum := rest[:idx2]
		if isAllDigits(possibleNum) {
			return grepMatch{file: file, lineNum: possibleNum, content: rest[idx2+1:]}, true
		}
	}

	return grepMatch{file: file, lineNum: "", content: rest}, true
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
