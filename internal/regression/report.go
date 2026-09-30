package regression

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/taqu/agentcap/internal/benchcompare"
)

// WriteJSON writes the canonical evaluation as one JSON object.
func WriteJSON(w io.Writer, e *Evaluation) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(e)
}

// WriteHuman renders the evaluation for local use and CI logs. Every check
// shows its baseline, candidate, change, configured limit, and status.
func WriteHuman(w io.Writer, e *Evaluation) error {
	var sb strings.Builder
	sb.WriteString("Benchmark regression check\n\n")
	if e.Status == StatusError {
		fmt.Fprintf(&sb, "Result: ERROR (evaluation could not be completed)\n%s\n", e.Error)
		_, err := io.WriteString(w, sb.String())
		return err
	}

	tw := tabwriter.NewWriter(&sb, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "Workload\tMetric\tBaseline\tCandidate\tChange\tAllowed increase\tStatus")
	var notes []string
	for _, we := range e.Workloads {
		if we.Status != WorkloadEvaluated {
			fmt.Fprintf(tw, "%s\t-\t-\t-\t-\t-\t%s\n", we.Workload, strings.ToUpper(string(we.Status)))
			notes = append(notes, fmt.Sprintf("%s: %s", we.Workload, we.Reason))
			continue
		}
		for _, c := range we.Checks {
			change := benchcompare.FormatDelta(&c.NumericMetric)
			if change == "" {
				change = "n/a"
			}
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", we.Workload, c.Name,
				benchcompare.FormatValue(c.Unit, c.Baseline, c.BaselineStatus),
				benchcompare.FormatValue(c.Unit, c.Candidate, c.CandidateStatus),
				change, formatLimit(c), statusLabel(c.Status))
			if c.Reason != "" {
				notes = append(notes, fmt.Sprintf("%s %s: %s", we.Workload, c.Name, c.Reason))
			}
		}
		if raw := we.RawBytes; raw != nil && raw.Delta != nil && *raw.Delta != 0 {
			notes = append(notes, fmt.Sprintf("%s raw_bytes (reference, not gated): %s -> %s (%s); fixture or tool output differs",
				we.Workload,
				benchcompare.FormatValue(raw.Unit, raw.Baseline, raw.BaselineStatus),
				benchcompare.FormatValue(raw.Unit, raw.Candidate, raw.CandidateStatus),
				benchcompare.FormatDelta(raw)))
		}
	}
	if err := tw.Flush(); err != nil {
		return err
	}
	if len(notes) > 0 {
		sb.WriteString("\nNotes:\n")
		for _, note := range notes {
			fmt.Fprintf(&sb, "  %s\n", note)
		}
	}
	s := e.Summary
	fmt.Fprintf(&sb, "\nChecks: %d passed, %d regression, %d not evaluable; workloads: %d of %d failed\n",
		s.Passed, s.Regressions, s.NotEvaluable, s.FailedWorkloads, s.Workloads)
	if e.Passed {
		sb.WriteString("Result: PASSED\n")
	} else {
		sb.WriteString("Result: FAILED\n")
	}
	_, err := io.WriteString(w, sb.String())
	return err
}

func formatLimit(c Check) string {
	var parts []string
	if r := c.MaxRelativeIncrease; r != nil {
		parts = append(parts, fmt.Sprintf("+%.1f%%", *r*100))
	}
	if a := c.MaxAbsoluteIncrease; a != nil {
		unit := c.Unit
		if unit == "" {
			unit = benchcompare.UnitCount
		}
		value := *a
		parts = append(parts, "+"+benchcompare.FormatValue(unit, &value, benchcompare.Measured))
	}
	return strings.Join(parts, " or ")
}

func statusLabel(s Status) string {
	switch s {
	case Pass:
		return "PASS"
	case Regression:
		return "REGRESSION"
	default:
		return "NOT EVALUABLE"
	}
}
