package delta

import (
	"fmt"
	"strings"

	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/store"
)

const maxGenericDeltaLines = 200

func genericDelta(baseline, current *store.Entry, exitChanged, stderrEqual bool) *Result {
	prevRaw := readFile(baseline.StdoutPath())
	currRaw := readFile(current.StdoutPath())
	if prevRaw == nil || currRaw == nil {
		return nil
	}

	prevLines := clean.Lines(clean.StripANSI(prevRaw))
	currLines := clean.Lines(clean.StripANSI(currRaw))

	added, removed := diffLineSet(prevLines, currLines)
	if len(added)+len(removed) > maxGenericDeltaLines {
		return nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap changed from %s\n", baseline.Meta.ID)
	if exitChanged {
		fmt.Fprintf(&sb, "exit: %d -> %d\n", baseline.Meta.ExitCode, current.Meta.ExitCode)
	}

	if len(added) > 0 {
		sb.WriteString("\nadded:\n")
		for _, l := range added {
			fmt.Fprintf(&sb, "+ %s\n", l)
		}
	}
	if len(removed) > 0 {
		sb.WriteString("\nremoved:\n")
		for _, l := range removed {
			fmt.Fprintf(&sb, "- %s\n", l)
		}
	}

	// Surface new stderr.
	if !stderrEqual {
		prevErr := readFile(baseline.StderrPath())
		currErr := readFile(current.StderrPath())
		errAdded, _ := diffLineSet(
			clean.Lines(clean.StripANSI(prevErr)),
			clean.Lines(clean.StripANSI(currErr)),
		)
		if len(errAdded) > 0 {
			sb.WriteString("\nstderr added:\n")
			for _, l := range errAdded {
				fmt.Fprintf(&sb, "%s\n", l)
			}
		}
	}

	if len(added) == 0 && len(removed) == 0 && !exitChanged {
		return nil
	}

	return &Result{
		Output:       sb.String(),
		Presentation: PresentationDelta,
	}
}

// diffLineSet returns (added, removed) as multiset difference.
func diffLineSet(prev, curr []string) (added, removed []string) {
	prevCount := make(map[string]int, len(prev))
	for _, l := range prev {
		prevCount[l]++
	}
	currCount := make(map[string]int, len(curr))
	for _, l := range curr {
		currCount[l]++
	}

	seen := make(map[string]bool)
	for _, l := range curr {
		if !seen[l] {
			seen[l] = true
			diff := currCount[l] - prevCount[l]
			for i := 0; i < diff; i++ {
				added = append(added, l)
			}
		}
	}

	seen = make(map[string]bool)
	for _, l := range prev {
		if !seen[l] {
			seen[l] = true
			diff := prevCount[l] - currCount[l]
			for i := 0; i < diff; i++ {
				removed = append(removed, l)
			}
		}
	}
	return added, removed
}
