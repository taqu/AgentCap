package delta

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/store"
)

// GrepDelta implements structural delta for grep/rg output.
type GrepDelta struct{}

func (g *GrepDelta) Delta(ctx context.Context, baseline, current *store.Entry) (*Result, error) {
	prevLines := readStdoutLines(baseline)
	currLines := readStdoutLines(current)

	prevMatches := parseGrepMatches(prevLines)
	currMatches := parseGrepMatches(currLines)

	prevSet := toStringSet(prevMatches)
	currSet := toStringSet(currMatches)

	var added, removed []string
	for m := range currSet {
		if !prevSet[m] {
			added = append(added, m)
		}
	}
	for m := range prevSet {
		if !currSet[m] {
			removed = append(removed, m)
		}
	}
	sort.Strings(added)
	sort.Strings(removed)

	exitChanged := baseline.Meta.ExitCode != current.Meta.ExitCode

	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap delta from %s\n", baseline.Meta.ID)
	fmt.Fprintf(&sb, "matches: %d -> %d\n", len(prevMatches), len(currMatches))
	if exitChanged {
		fmt.Fprintf(&sb, "exit: %d -> %d\n", baseline.Meta.ExitCode, current.Meta.ExitCode)
	}
	if len(added) > 0 {
		sb.WriteString("\nadded:\n")
		for _, m := range added {
			sb.WriteString(m)
			sb.WriteByte('\n')
		}
	}
	if len(removed) > 0 {
		sb.WriteString("\nremoved:\n")
		for _, m := range removed {
			sb.WriteString(m)
			sb.WriteByte('\n')
		}
	}

	return &Result{Output: sb.String(), Presentation: PresentationDelta}, nil
}

func readStdoutLines(e *store.Entry) []string {
	data := readFile(e.StdoutPath())
	if data == nil {
		return nil
	}
	return clean.Lines(clean.StripANSI(data))
}

func parseGrepMatches(lines []string) []string {
	var matches []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.Contains(line, ":") {
			matches = append(matches, line)
		}
	}
	return matches
}

func toStringSet(lines []string) map[string]bool {
	s := make(map[string]bool, len(lines))
	for _, l := range lines {
		s[l] = true
	}
	return s
}
