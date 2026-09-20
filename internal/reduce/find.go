package reduce

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/exec"
)

// FindReducer handles find output (one path per line).
type FindReducer struct{}

func (f *FindReducer) Reduce(r *exec.Result) *ReducedResult {
	raw := clean.StripANSI(r.Stdout)
	rawBytes := len(r.Stdout)

	lines := clean.Lines(raw)
	stderr := buildStderr(r.Stderr)

	// Small output: pass through.
	if len(raw) <= smallThresholdBytes && len(lines) <= smallThresholdLines {
		out := string(raw) + stderr
		return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
	}

	// Validate that output looks like paths.
	if !looksLikePaths(lines) {
		g := &GenericReducer{}
		return g.Reduce(r)
	}

	byRoot := make(map[string]int)
	byExt := make(map[string]int)
	var sample []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// First path component (root).
		root := firstComponent(line)
		byRoot[root]++

		// Extension.
		ext := filepath.Ext(line)
		if ext != "" {
			byExt[ext]++
		}

		if len(sample) < 5 {
			sample = append(sample, line)
		}
	}

	total := len(lines)
	omitted := total - len(sample)

	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap find paths=%d\n", total)

	// Top roots by count.
	sb.WriteString("\nby-root:\n")
	roots := sortedMapKeys(byRoot)
	shown := 0
	for _, k := range roots {
		if shown >= 10 {
			break
		}
		fmt.Fprintf(&sb, "%s %d\n", k, byRoot[k])
		shown++
	}

	// Top extensions.
	if len(byExt) > 0 {
		sb.WriteString("\nextensions:\n")
		exts := sortedMapKeys(byExt)
		shown = 0
		for _, k := range exts {
			if shown >= 10 {
				break
			}
			fmt.Fprintf(&sb, "%s %d\n", k, byExt[k])
			shown++
		}
	}

	// Sample.
	if len(sample) > 0 {
		sb.WriteString("\nsample:\n")
		for _, s := range sample {
			sb.WriteString(s)
			sb.WriteByte('\n')
		}
	}

	if omitted > 0 {
		fmt.Fprintf(&sb, "\nomitted=%d\n", omitted)
	}

	sb.WriteString(stderr)
	out := sb.String()
	return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
}

// looksLikePaths returns true if the majority of lines look like filesystem paths.
func looksLikePaths(lines []string) bool {
	if len(lines) == 0 {
		return false
	}
	pathLike := 0
	check := min(20, len(lines))
	for _, line := range lines[:check] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// A path should contain no spaces (common for find output) or start with . or /
		if strings.HasPrefix(line, "./") || strings.HasPrefix(line, "/") || strings.HasPrefix(line, "../") {
			pathLike++
		} else if !strings.ContainsAny(line, " \t") {
			pathLike++
		}
	}
	return pathLike*2 >= check
}

// firstComponent returns the first path component of p.
func firstComponent(p string) string {
	// Normalize separators.
	p = filepath.ToSlash(p)
	// Strip leading "./"
	p = strings.TrimPrefix(p, "./")
	p = strings.TrimPrefix(p, "/")
	idx := strings.IndexByte(p, '/')
	if idx < 0 {
		return p
	}
	return p[:idx]
}

// sortedMapKeys returns keys of a map[string]int sorted by value descending,
// then alphabetically.
func sortedMapKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if m[keys[i]] != m[keys[j]] {
			return m[keys[i]] > m[keys[j]]
		}
		return keys[i] < keys[j]
	})
	return keys
}
