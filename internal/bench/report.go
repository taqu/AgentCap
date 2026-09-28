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
	Command               []string `json:"command"`
	ExitCode              int      `json:"exit_code"`
	RawStdoutBytes        int64    `json:"raw_stdout_bytes"`
	RawStderrBytes        int64    `json:"raw_stderr_bytes"`
	RawBytes              int64    `json:"raw_bytes"`
	Truncated             bool     `json:"truncated"`
	AgentVisibleBytes     int64    `json:"agent_visible_bytes"`
	StatelessVisibleBytes int64    `json:"stateless_visible_bytes"`
	StatefulVisibleBytes  int64    `json:"stateful_visible_bytes"`
	Presentation          string   `json:"presentation"`
	ReductionRatio        *float64 `json:"reduction_ratio"` // null when raw_bytes == 0
	ExecutionDurationNs   int64    `json:"execution_duration_ns"`
	ReduceDurationNs      int64    `json:"reduce_duration_ns"`
	ProcessingDurationNs  int64    `json:"processing_duration_ns"`
	ResultID              string   `json:"result_id"`
}

// WriteJSON writes m as a single JSON object followed by a newline.
func WriteJSON(w io.Writer, m *Measurement) error {
	r := jsonReport{
		Command:               m.Command,
		ExitCode:              m.ExitCode,
		RawStdoutBytes:        m.RawStdoutBytes,
		RawStderrBytes:        m.RawStderrBytes,
		RawBytes:              m.RawBytes,
		Truncated:             m.Truncated,
		AgentVisibleBytes:     m.AgentVisibleBytes,
		StatelessVisibleBytes: m.StatelessVisibleBytes,
		StatefulVisibleBytes:  m.StatefulVisibleBytes,
		Presentation:          m.Presentation,
		ExecutionDurationNs:   m.ExecutionDuration.Nanoseconds(),
		ReduceDurationNs:      m.ReduceDuration.Nanoseconds(),
		ProcessingDurationNs:  m.ProcessingDuration.Nanoseconds(),
		ResultID:              m.ResultID,
	}
	if ratio, ok := m.ReductionRatio(); ok {
		r.ReductionRatio = &ratio
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

// WriteSessionText writes a compact workflow-level benchmark report.
func WriteSessionText(w io.Writer, m *SessionMeasurement) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Benchmark Session: %s\n\n", m.SessionID)
	fmt.Fprintf(&sb, "commands: %d\n\n", m.CommandCount)
	sb.WriteString("output:\n")
	fmt.Fprintf(&sb, "  raw:                  %d B\n", m.RawBytes)
	fmt.Fprintf(&sb, "  stateless:            %d B\n", m.StatelessBytes)
	fmt.Fprintf(&sb, "  stateful:             %d B\n", m.StatefulBytes)
	sb.WriteString("\nreduction:\n")
	fmt.Fprintf(&sb, "  stateless vs raw:      %s\n", formatRatio(m.StatelessBytes, m.RawBytes))
	fmt.Fprintf(&sb, "  stateful vs raw:       %s\n", formatRatio(m.StatefulBytes, m.RawBytes))
	fmt.Fprintf(&sb, "  stateful vs stateless: %s\n", formatRatio(m.StatefulBytes, m.StatelessBytes))
	sb.WriteString("\npresentations:\n")
	fmt.Fprintf(&sb, "  full:                  %d\n", m.FullCount)
	fmt.Fprintf(&sb, "  delta:                 %d\n", m.DeltaCount)
	fmt.Fprintf(&sb, "  unchanged:             %d\n", m.UnchangedCount)
	sb.WriteString("\ntime:\n")
	fmt.Fprintf(&sb, "  command execution:     %s\n", formatDuration(m.ExecutionDuration))
	fmt.Fprintf(&sb, "  AgentCap processing:   %s (reduce %s)\n", formatDuration(m.ProcessingDuration), formatDuration(m.ReduceDuration))
	if len(m.Commands) > 0 {
		sb.WriteString("\ncommands:\n")
		for i, command := range m.Commands {
			fmt.Fprintf(&sb, "  %d  %s  raw=%d B stateless=%d B stateful=%d B %s id=%s\n",
				i+1, strings.Join(command.Command, " "), command.RawBytes,
				command.StatelessVisibleBytes, command.StatefulVisibleBytes,
				command.Presentation, command.ResultID)
		}
	}
	_, err := io.WriteString(w, sb.String())
	return err
}

type sessionJSONReport struct {
	SessionID            string   `json:"session_id"`
	CommandCount         int      `json:"command_count"`
	RawBytes             int64    `json:"raw_bytes"`
	StatelessBytes       int64    `json:"stateless_bytes"`
	StatefulBytes        int64    `json:"stateful_bytes"`
	StatelessVsRaw       *float64 `json:"stateless_vs_raw"`
	StatefulVsRaw        *float64 `json:"stateful_vs_raw"`
	StatefulVsStateless  *float64 `json:"stateful_vs_stateless"`
	FullCount            int      `json:"full_count"`
	DeltaCount           int      `json:"delta_count"`
	UnchangedCount       int      `json:"unchanged_count"`
	ExecutionDurationNs  int64    `json:"execution_duration_ns"`
	ReduceDurationNs     int64    `json:"reduce_duration_ns"`
	ProcessingDurationNs int64    `json:"processing_duration_ns"`
}

// WriteSessionJSON writes the aggregate session schema as JSON.
func WriteSessionJSON(w io.Writer, m *SessionMeasurement) error {
	r := sessionJSONReport{
		SessionID: m.SessionID, CommandCount: m.CommandCount, RawBytes: m.RawBytes,
		StatelessBytes: m.StatelessBytes, StatefulBytes: m.StatefulBytes,
		FullCount: m.FullCount, DeltaCount: m.DeltaCount, UnchangedCount: m.UnchangedCount,
		ExecutionDurationNs: m.ExecutionDuration.Nanoseconds(), ReduceDurationNs: m.ReduceDuration.Nanoseconds(),
		ProcessingDurationNs: m.ProcessingDuration.Nanoseconds(),
	}
	r.StatelessVsRaw = ratioPtr(m.StatelessBytes, m.RawBytes)
	r.StatefulVsRaw = ratioPtr(m.StatefulBytes, m.RawBytes)
	r.StatefulVsStateless = ratioPtr(m.StatefulBytes, m.StatelessBytes)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func ratioPtr(numerator, denominator int64) *float64 {
	if denominator == 0 {
		return nil
	}
	ratio := 1 - float64(numerator)/float64(denominator)
	return &ratio
}

func formatRatio(numerator, denominator int64) string {
	ratio := ratioPtr(numerator, denominator)
	if ratio == nil {
		return "n/a"
	}
	return fmt.Sprintf("%.1f%%", *ratio*100)
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
