package delta

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/store"
)

// FindDelta implements structural delta for find output.
type FindDelta struct{}

func (f *FindDelta) Delta(ctx context.Context, baseline, current *store.Entry) (*Result, error) {
	prevPaths := readPaths(baseline)
	currPaths := readPaths(current)

	prevSet := toStringSet(prevPaths)
	currSet := toStringSet(currPaths)

	var added, removed []string
	for p := range currSet {
		if !prevSet[p] {
			added = append(added, p)
		}
	}
	for p := range prevSet {
		if !currSet[p] {
			removed = append(removed, p)
		}
	}
	sort.Strings(added)
	sort.Strings(removed)

	exitChanged := baseline.Meta.ExitCode != current.Meta.ExitCode

	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap delta from %s\n", baseline.Meta.ID)
	fmt.Fprintf(&sb, "paths: %d -> %d\n", len(prevPaths), len(currPaths))
	if exitChanged {
		fmt.Fprintf(&sb, "exit: %d -> %d\n", baseline.Meta.ExitCode, current.Meta.ExitCode)
	}
	if len(added) > 0 {
		sb.WriteString("\nadded:\n")
		for _, p := range added {
			sb.WriteString(p)
			sb.WriteByte('\n')
		}
	}
	if len(removed) > 0 {
		sb.WriteString("\nremoved:\n")
		for _, p := range removed {
			sb.WriteString(p)
			sb.WriteByte('\n')
		}
	}

	return &Result{Output: sb.String(), Presentation: PresentationDelta}, nil
}

func readPaths(e *store.Entry) []string {
	data := readFile(e.StdoutPath())
	if data == nil {
		return nil
	}
	lines := clean.Lines(clean.StripANSI(data))
	var paths []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			paths = append(paths, l)
		}
	}
	return paths
}
