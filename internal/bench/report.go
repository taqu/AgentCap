package bench

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// WriteText writes the human-readable benchmark report for m.
//
// The report is benchmark UI only; none of it is part of the agent-visible
// result measured by AgentVisibleBytes.
func WriteText(w io.Writer, m *Measurement) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Benchmark: %s\n\n", strings.Join(m.Command, " "))

	sb.WriteString("command:\n")
	fmt.Fprintf(&sb, "  exit:       %d\n", m.ExitCode)
	fmt.Fprintf(&sb, "  duration:   %s\n", formatDuration(m.ExecutionDuration))

	sb.WriteString("\nraw:\n")
	fmt.Fprintf(&sb, "  stdout:     %d B\n", m.RawStdoutBytes)
	fmt.Fprintf(&sb, "  stderr:     %d B\n", m.RawStderrBytes)
	fmt.Fprintf(&sb, "  total:      %d B\n", m.RawBytes)
	if m.Truncated {
		sb.WriteString("  truncated:  yes (capture limit reached; raw sizes are lower bounds)\n")
	}

	sb.WriteString("\nagentcap:\n")
	fmt.Fprintf(&sb, "  visible:    %d B\n", m.AgentVisibleBytes)
	fmt.Fprintf(&sb, "  reduction:  %s\n", formatReduction(m))
	fmt.Fprintf(&sb, "  processing: %s (reduce %s)\n",
		formatDuration(m.ProcessingDuration), formatDuration(m.ReduceDuration))

	sb.WriteString("\nresult:\n")
	if m.ResultID != "" {
		fmt.Fprintf(&sb, "  id:         %s\n", m.ResultID)
	} else {
		sb.WriteString("  id:         (not stored)\n")
	}

	_, err := io.WriteString(w, sb.String())
	return err
}

// jsonReport is the stable --json schema for a single-command benchmark.
type jsonReport struct {
	Command              []string `json:"command"`
	ExitCode             int      `json:"exit_code"`
	RawStdoutBytes       int64    `json:"raw_stdout_bytes"`
	RawStderrBytes       int64    `json:"raw_stderr_bytes"`
	RawBytes             int64    `json:"raw_bytes"`
	Truncated            bool     `json:"truncated"`
	AgentVisibleBytes    int64    `json:"agent_visible_bytes"`
	ReductionRatio       *float64 `json:"reduction_ratio"` // null when raw_bytes == 0
	ExecutionDurationNs  int64    `json:"execution_duration_ns"`
	ReduceDurationNs     int64    `json:"reduce_duration_ns"`
	ProcessingDurationNs int64    `json:"processing_duration_ns"`
	ResultID             string   `json:"result_id"`
}

// WriteJSON writes m as a single JSON object followed by a newline.
func WriteJSON(w io.Writer, m *Measurement) error {
	r := jsonReport{
		Command:              m.Command,
		ExitCode:             m.ExitCode,
		RawStdoutBytes:       m.RawStdoutBytes,
		RawStderrBytes:       m.RawStderrBytes,
		RawBytes:             m.RawBytes,
		Truncated:            m.Truncated,
		AgentVisibleBytes:    m.AgentVisibleBytes,
		ExecutionDurationNs:  m.ExecutionDuration.Nanoseconds(),
		ReduceDurationNs:     m.ReduceDuration.Nanoseconds(),
		ProcessingDurationNs: m.ProcessingDuration.Nanoseconds(),
		ResultID:             m.ResultID,
	}
	if ratio, ok := m.ReductionRatio(); ok {
		r.ReductionRatio = &ratio
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func formatReduction(m *Measurement) string {
	ratio, ok := m.ReductionRatio()
	if !ok {
		return "n/a (no raw output)"
	}
	return fmt.Sprintf("%.1f%%", ratio*100)
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
