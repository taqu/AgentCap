package benchcompare

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/taqu/agentcap/internal/stats"
	"github.com/taqu/agentcap/internal/workload"
)

// WriteJSON writes the canonical comparison as one JSON object.
func WriteJSON(w io.Writer, c *Comparison) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(c)
}

// label names for the human view; metrics without a label are JSON-only.
var metricLabels = map[string]string{
	"median_total_visible_bytes": "Median visible bytes",
	"median_command_count":       "Median commands",
	"median_wall_time_ns":        "Median wall time",
	"median_processing_ns":       "AgentCap processing",
	"total_show_count":           "Show calls (total)",
	"total_raw_retrieval_count":  "Raw retrievals (total)",
	"total_visible_bytes":        "Visible bytes",
	"command_count":              "Commands",
	"wall_time_ns":               "Wall time",
	"processing_ns":              "AgentCap processing",
	"show_count":                 "Show calls",
	"raw_retrieval_count":        "Raw retrievals",
	"show_bytes":                 "Show bytes",
	"raw_retrieval_bytes":        "Raw retrieval bytes",
}

var outcomeLabels = map[string]string{
	"task_success":              "Task success",
	"timeout":                   "Timeouts",
	"agent_error":               "Agent errors",
	"canceled":                  "Canceled",
	"trials_with_show":          "Show used",
	"trials_with_raw_retrieval": "Raw fallback",
}

// WriteHuman renders a compact side-by-side table from the canonical
// comparison. It performs no arithmetic of its own beyond formatting.
func WriteHuman(w io.Writer, c *Comparison) error {
	var sb strings.Builder
	sb.WriteString("Benchmark comparison\n")
	fmt.Fprintf(&sb, "Workload: %s\n", c.Baseline.Workload)
	if c.Baseline.Agent != "" {
		fmt.Fprintf(&sb, "Agent:    %s\n", c.Baseline.Agent)
	}
	if c.Kind == workload.ResultKindRepeated {
		sb.WriteString("Results:  repeated trials (medians)\n")
	} else {
		sb.WriteString("Results:  single observation\n")
	}

	t := &table{}
	t.row("", "baseline", "candidate", "delta")
	if c.Baseline.Agent != "" {
		t.row("Mode", c.Baseline.Mode, c.Candidate.Mode, "")
	}
	if c.Kind == workload.ResultKindRepeated {
		t.row("Trials", trials(c.Baseline), trials(c.Candidate), "")
		if c.Baseline.RunStatus != "completed" || c.Candidate.RunStatus != "completed" {
			t.row("Run status", c.Baseline.RunStatus, c.Candidate.RunStatus, "")
		}
	} else if c.Baseline.ExecutionStatus != "" || c.Candidate.ExecutionStatus != "" {
		t.row("Execution", orNA(c.Baseline.ExecutionStatus), orNA(c.Candidate.ExecutionStatus), "")
	}
	t.blank()

	t.outcome(c.Outcome("task_success"))
	for _, name := range []string{"timeout", "agent_error", "canceled"} {
		if o := c.Outcome(name); o != nil && (nonZero(o.Baseline) || nonZero(o.Candidate)) {
			t.outcome(o)
		}
	}
	t.blank()

	if c.Kind == workload.ResultKindRepeated {
		t.metrics(c, "median_total_visible_bytes", "median_command_count", "median_wall_time_ns")
		t.blank()
		t.outcome(c.Outcome("trials_with_show"))
		t.outcome(c.Outcome("trials_with_raw_retrieval"))
		t.metrics(c, "total_show_count", "total_raw_retrieval_count", "median_processing_ns")
	} else {
		t.metrics(c, "total_visible_bytes", "command_count", "wall_time_ns")
		t.blank()
		t.metrics(c, "show_count", "raw_retrieval_count", "show_bytes", "raw_retrieval_bytes", "processing_ns")
	}
	sb.WriteString("\n")
	t.write(&sb)
	sb.WriteString("\ndelta = candidate - baseline.  - = not applicable, n/a = not recorded.\n")
	_, err := io.WriteString(w, sb.String())
	return err
}

type table struct{ rows [][4]string }

func (t *table) row(label, b, k, d string) { t.rows = append(t.rows, [4]string{label, b, k, d}) }
func (t *table) blank()                    { t.rows = append(t.rows, [4]string{}) }

func (t *table) outcome(o *OutcomeMetric) {
	if o == nil {
		return
	}
	t.row(outcomeLabels[o.Name], formatFraction(o.Baseline, o.BaselineStatus), formatFraction(o.Candidate, o.CandidateStatus), "")
}

func (t *table) metrics(c *Comparison, names ...string) {
	for _, name := range names {
		m := c.Metric(name)
		if m == nil {
			continue
		}
		t.row(metricLabels[name], formatValue(m.Unit, m.Baseline, m.BaselineStatus), formatValue(m.Unit, m.Candidate, m.CandidateStatus), formatDelta(m))
	}
}

func (t *table) write(sb *strings.Builder) {
	var width [4]int
	for _, r := range t.rows {
		for i, cell := range r {
			width[i] = max(width[i], len([]rune(cell)))
		}
	}
	for _, r := range t.rows {
		if r == [4]string{} {
			sb.WriteString("\n")
			continue
		}
		line := fmt.Sprintf("%-*s  %*s  %*s  %*s", width[0], r[0], width[1], r[1], width[2], r[2], width[3], r[3])
		sb.WriteString(strings.TrimRight(line, " "))
		sb.WriteString("\n")
	}
}

func trials(s Side) string {
	if s.RequestedTrialCount != 0 && s.RequestedTrialCount != s.TrialCount {
		return fmt.Sprintf("%d of %d", s.TrialCount, s.RequestedTrialCount)
	}
	return fmt.Sprintf("%d", s.TrialCount)
}

func orNA(s string) string {
	if s == "" {
		return "n/a"
	}
	return s
}

func nonZero(f *Fraction) bool { return f != nil && f.Count != 0 }

func unavailable(status Availability) string {
	if status == NotApplicable {
		return "-"
	}
	return "n/a"
}

func formatFraction(f *Fraction, status Availability) string {
	if status != Measured || f == nil {
		return unavailable(status)
	}
	return fmt.Sprintf("%d/%d", f.Count, f.Total)
}

func formatValue(unit Unit, v *int64, status Availability) string {
	if status != Measured || v == nil {
		return unavailable(status)
	}
	return formatAmount(unit, *v)
}

func formatAmount(unit Unit, v int64) string {
	switch unit {
	case UnitBytes:
		return stats.FormatBytes(v)
	case UnitNanoseconds:
		return formatDuration(time.Duration(v))
	default:
		return fmt.Sprintf("%d", v)
	}
}

func formatDelta(m *NumericMetric) string {
	if m.Delta == nil {
		return ""
	}
	d := *m.Delta
	sign := "+"
	if d < 0 {
		sign, d = "-", -d
	} else if d == 0 {
		sign = ""
	}
	s := sign + formatAmount(m.Unit, d)
	if m.RelativeDelta != nil {
		s += fmt.Sprintf(" (%+.1f%%)", *m.RelativeDelta*100)
	}
	return s
}

func formatDuration(d time.Duration) string {
	switch {
	case d >= time.Second:
		return fmt.Sprintf("%.2fs", d.Seconds())
	case d >= time.Millisecond:
		return fmt.Sprintf("%.1fms", float64(d)/float64(time.Millisecond))
	default:
		return fmt.Sprintf("%.1fµs", float64(d)/float64(time.Microsecond))
	}
}
