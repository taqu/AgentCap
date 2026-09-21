package delta

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/store"
)

// LsDelta implements structural delta for ls output.
type LsDelta struct{}

func (l *LsDelta) Delta(ctx context.Context, baseline, current *store.Entry) (*Result, error) {
	prevEntries := readLsEntries(baseline)
	currEntries := readLsEntries(current)

	prevSet := toStringSet(prevEntries)
	currSet := toStringSet(currEntries)

	var added, removed []string
	for e := range currSet {
		if !prevSet[e] {
			added = append(added, e)
		}
	}
	for e := range prevSet {
		if !currSet[e] {
			removed = append(removed, e)
		}
	}
	sort.Strings(added)
	sort.Strings(removed)

	exitChanged := baseline.Meta.ExitCode != current.Meta.ExitCode

	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap delta from %s\n", baseline.Meta.ID)
	if exitChanged {
		fmt.Fprintf(&sb, "exit: %d -> %d\n", baseline.Meta.ExitCode, current.Meta.ExitCode)
	}
	if len(added) > 0 {
		sb.WriteString("\nadded:\n")
		for _, e := range added {
			sb.WriteString(e)
			sb.WriteByte('\n')
		}
	}
	if len(removed) > 0 {
		sb.WriteString("\nremoved:\n")
		for _, e := range removed {
			sb.WriteString(e)
			sb.WriteByte('\n')
		}
	}

	return &Result{Output: sb.String(), Presentation: PresentationDelta}, nil
}

func readLsEntries(e *store.Entry) []string {
	data := readFile(e.StdoutPath())
	if data == nil {
		return nil
	}
	lines := clean.Lines(clean.StripANSI(data))

	// If long format, extract just the filename (last field).
	if isLongFormat(lines) {
		var names []string
		for _, line := range lines {
			if strings.HasPrefix(line, "total ") || line == "" {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) >= 9 {
				names = append(names, strings.Join(fields[8:], " "))
			}
		}
		return names
	}

	// Plain format.
	var entries []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			entries = append(entries, l)
		}
	}
	return entries
}

func isLongFormat(lines []string) bool {
	for _, line := range lines {
		if len(line) < 4 {
			continue
		}
		c := line[0]
		if c == '-' || c == 'd' || c == 'l' || c == 'c' || c == 'b' || c == 'p' || c == 's' {
			return true
		}
	}
	return false
}
