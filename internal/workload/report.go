package workload

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// WriteHuman writes the concise human representation of the canonical result.
func WriteHuman(w io.Writer, result *BenchmarkResult, verbose bool) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Benchmark Workload: %s\n\n", result.Workload)
	run := result.presentation
	if run != nil && run.SessionID != "" {
		fmt.Fprintf(&sb, "session: %s\n\n", run.SessionID)
	}
	if run != nil {
		fmt.Fprintf(&sb, "steps:\n  total:       %d\n  commands:    %d\n  mutations:   %d\n", run.TotalSteps, result.Commands, run.MutationCount)
	} else {
		fmt.Fprintf(&sb, "commands: %d\n", result.Commands)
	}
	fmt.Fprintf(&sb, "\noutput:\n  raw:                  %d B\n  stateless:            %d B\n  stateful:             %d B\n", result.RawBytes, result.StatelessBytes, result.StatefulBytes)
	fmt.Fprintf(&sb, "\nreduction:\n  stateless vs raw:      %s\n  stateful vs raw:       %s\n  stateful vs stateless: %s\n", ratio(result.StatelessBytes, result.RawBytes), ratio(result.StatefulBytes, result.RawBytes), ratio(result.StatefulBytes, result.StatelessBytes))
	if run != nil && run.Aggregate != nil {
		a := run.Aggregate
		fmt.Fprintf(&sb, "\npresentations:\n  full:                  %d\n  delta:                 %d\n  unchanged:             %d\n", a.FullCount, a.DeltaCount, a.UnchangedCount)
		fmt.Fprintf(&sb, "\ntime:\n  command execution:     %s\n  AgentCap processing:   %s (reduce %s)\n  workload wall clock:   %s\n", duration(a.ExecutionDuration), duration(a.ProcessingDuration), duration(a.ReduceDuration), duration(run.WallDuration))
	}
	if verbose && run != nil && len(run.Steps) > 0 {
		sb.WriteString("\nsteps:\n")
		for _, step := range run.Steps {
			if step.Measurement == nil {
				fmt.Fprintf(&sb, "  %d  %s\n", step.Number, step.Type)
				continue
			}
			m := step.Measurement
			fmt.Fprintf(&sb, "  %d  run  exit=%d raw=%d B stateless=%d B stateful=%d B %s id=%s\n", step.Number, m.ExitCode, m.RawBytes, m.StatelessVisibleBytes, m.StatefulVisibleBytes, m.Presentation, m.ResultID)
		}
	}
	if run != nil && run.RetainedWorkspace != "" {
		fmt.Fprintf(&sb, "\nworkspace retained: %s\n", run.RetainedWorkspace)
	}
	_, err := io.WriteString(w, sb.String())
	return err
}

// WriteJSON writes only the stable schema as one JSON object. Internal B3
// presentation details are deliberately excluded.
func WriteJSON(w io.Writer, result *BenchmarkResult) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}

// WriteText preserves the B3 programmatic API while routing it through the
// canonical B4 result model.
func WriteText(w io.Writer, run *Result, verbose bool) error {
	return WriteHuman(w, NewBenchmarkResult(run), verbose)
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
