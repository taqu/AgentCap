package workload

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// WriteText writes the concise B3 aggregate without replaying command output.
func WriteText(w io.Writer, r *Result, verbose bool) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Benchmark Workload: %s\n\n", r.Name)
	fmt.Fprintf(&sb, "session: %s\n\n", r.SessionID)
	commands := 0
	if r.Aggregate != nil {
		commands = r.Aggregate.CommandCount
	}
	fmt.Fprintf(&sb, "steps:\n  total:       %d\n  commands:    %d\n  mutations:   %d\n", r.TotalSteps, commands, r.MutationCount)
	if a := r.Aggregate; a != nil {
		fmt.Fprintf(&sb, "\noutput:\n  raw:                  %d B\n  stateless:            %d B\n  stateful:             %d B\n", a.RawBytes, a.StatelessBytes, a.StatefulBytes)
		fmt.Fprintf(&sb, "\nreduction:\n  stateless vs raw:      %s\n  stateful vs raw:       %s\n  stateful vs stateless: %s\n", ratio(a.StatelessBytes, a.RawBytes), ratio(a.StatefulBytes, a.RawBytes), ratio(a.StatefulBytes, a.StatelessBytes))
		fmt.Fprintf(&sb, "\npresentations:\n  full:                  %d\n  delta:                 %d\n  unchanged:             %d\n", a.FullCount, a.DeltaCount, a.UnchangedCount)
		fmt.Fprintf(&sb, "\ntime:\n  command execution:     %s\n  AgentCap processing:   %s (reduce %s)\n  workload wall clock:   %s\n", duration(a.ExecutionDuration), duration(a.ProcessingDuration), duration(a.ReduceDuration), duration(r.WallDuration))
	}
	if verbose && len(r.Steps) > 0 {
		sb.WriteString("\nsteps:\n")
		for _, step := range r.Steps {
			if step.Measurement == nil {
				fmt.Fprintf(&sb, "  %d  %s\n", step.Number, step.Type)
				continue
			}
			m := step.Measurement
			fmt.Fprintf(&sb, "  %d  run  exit=%d raw=%d B stateless=%d B stateful=%d B %s id=%s\n", step.Number, m.ExitCode, m.RawBytes, m.StatelessVisibleBytes, m.StatefulVisibleBytes, m.Presentation, m.ResultID)
		}
	}
	if r.RetainedWorkspace != "" {
		fmt.Fprintf(&sb, "\nworkspace retained: %s\n", r.RetainedWorkspace)
	}
	_, err := io.WriteString(w, sb.String())
	return err
}

func ratio(n, d int64) string {
	if d == 0 {
		return "n/a"
	}
	return fmt.Sprintf("%.1f%%", (1-float64(n)/float64(d))*100)
}

func duration(d time.Duration) string {
	switch {
	case d >= time.Second:
		return fmt.Sprintf("%.2fs", d.Seconds())
	case d >= time.Millisecond:
		return fmt.Sprintf("%.1fms", float64(d)/float64(time.Millisecond))
	default:
		return fmt.Sprintf("%.1fµs", float64(d)/float64(time.Microsecond))
	}
}
