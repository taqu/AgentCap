package benchmarkcompare

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

func WriteJSON(w io.Writer, comparison *Comparison) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(comparison)
}

func WriteHuman(w io.Writer, comparison *Comparison) error {
	var output strings.Builder
	output.WriteString("Benchmark comparison\n")
	fmt.Fprintf(&output, "Workload: %s\n", comparison.Baseline.Workload)
	if comparison.Baseline.Agent != "" {
		fmt.Fprintf(&output, "Agent:    %s\n", comparison.Baseline.Agent)
	}
	fmt.Fprintf(&output, "Statistic: %s\n\n", comparison.Statistic)
	writeRow(&output, "", "baseline", "candidate", "delta")
	writeRow(&output, "Mode", displayText(comparison.Baseline.Mode), displayText(comparison.Candidate.Mode), "-")
	writeRow(&output, "Trials", trialSummary(comparison.Baseline), trialSummary(comparison.Candidate), "-")
	writeRow(&output, "Task success", success(comparison.Baseline), success(comparison.Candidate), "-")
	writeRow(&output, "Task failures", optionalCount(comparison.Baseline.TaskFailureCount), optionalCount(comparison.Candidate.TaskFailureCount), "-")
	writeRow(&output, "Timeouts", optionalCount(comparison.Baseline.TimeoutCount), optionalCount(comparison.Candidate.TimeoutCount), "-")
	writeRow(&output, "Agent errors", optionalCount(comparison.Baseline.AgentErrorCount), optionalCount(comparison.Candidate.AgentErrorCount), "-")
	writeRow(&output, "Canceled", optionalCount(comparison.Baseline.CanceledCount), optionalCount(comparison.Candidate.CanceledCount), "-")
	writeRow(&output, metricLabel(comparison, "visible bytes"), metricValue(comparison.Metrics.TotalVisibleBytes, formatBytes), metricCandidate(comparison.Metrics.TotalVisibleBytes, formatBytes), metricDelta(comparison.Metrics.TotalVisibleBytes, formatSignedBytes))
	writeRow(&output, metricLabel(comparison, "commands"), metricValue(comparison.Metrics.CommandCount, formatInteger), metricCandidate(comparison.Metrics.CommandCount, formatInteger), metricDelta(comparison.Metrics.CommandCount, formatSignedInteger))
	writeRow(&output, metricLabel(comparison, "wall time"), metricValue(comparison.Metrics.WallTimeNS, formatDuration), metricCandidate(comparison.Metrics.WallTimeNS, formatDuration), metricDelta(comparison.Metrics.WallTimeNS, formatSignedDuration))
	writeRow(&output, "Show used", usage(comparison.Metrics.TrialsWithShow.Baseline, comparison.Baseline.TrialCount), usage(comparison.Metrics.TrialsWithShow.Candidate, comparison.Candidate.TrialCount), metricDelta(comparison.Metrics.TrialsWithShow, formatSignedInteger))
	writeRow(&output, "Raw used", usage(comparison.Metrics.TrialsWithRawRetrieval.Baseline, comparison.Baseline.TrialCount), usage(comparison.Metrics.TrialsWithRawRetrieval.Candidate, comparison.Candidate.TrialCount), metricDelta(comparison.Metrics.TrialsWithRawRetrieval, formatSignedInteger))
	writeRow(&output, "Show calls", metricValue(comparison.Metrics.TotalShowCount, formatInteger), metricCandidate(comparison.Metrics.TotalShowCount, formatInteger), metricDelta(comparison.Metrics.TotalShowCount, formatSignedInteger))
	writeRow(&output, "Raw retrievals", metricValue(comparison.Metrics.TotalRawRetrievalCount, formatInteger), metricCandidate(comparison.Metrics.TotalRawRetrievalCount, formatInteger), metricDelta(comparison.Metrics.TotalRawRetrievalCount, formatSignedInteger))
	writeRow(&output, metricLabel(comparison, "AgentCap processing"), metricValue(comparison.Metrics.ProcessingNS, formatDuration), metricCandidate(comparison.Metrics.ProcessingNS, formatDuration), metricDelta(comparison.Metrics.ProcessingNS, formatSignedDuration))
	_, err := io.WriteString(w, output.String())
	return err
}

func writeRow(output *strings.Builder, label, baseline, candidate, delta string) {
	fmt.Fprintf(output, "%-28s %14s %14s %20s\n", label, baseline, candidate, delta)
}

func metricLabel(comparison *Comparison, label string) string {
	if comparison.Statistic == "median" {
		return "Median " + label
	}
	return strings.ToUpper(label[:1]) + label[1:]
}

func success(side Side) string {
	if side.SuccessCount == nil {
		return "n/a"
	}
	return fmt.Sprintf("%d/%d", *side.SuccessCount, side.TrialCount)
}

func optionalCount(value *int) string {
	if value == nil {
		return "n/a"
	}
	return fmt.Sprint(*value)
}

func trialSummary(side Side) string {
	if side.RequestedTrialCount != side.TrialCount || side.RunStatus != "completed" {
		return fmt.Sprintf("%d/%d %s", side.TrialCount, side.RequestedTrialCount, side.RunStatus)
	}
	return fmt.Sprint(side.TrialCount)
}

func usage(value *int64, trials int) string {
	if value == nil {
		return "n/a"
	}
	return fmt.Sprintf("%d/%d", *value, trials)
}

func displayText(value string) string {
	if value == "" {
		return "n/a"
	}
	return value
}

func metricValue(metric NumericComparison, format func(int64) string) string {
	if metric.Baseline == nil {
		return "n/a"
	}
	return format(*metric.Baseline)
}

func metricCandidate(metric NumericComparison, format func(int64) string) string {
	if metric.Candidate == nil {
		return "n/a"
	}
	return format(*metric.Candidate)
}

func metricDelta(metric NumericComparison, format func(int64) string) string {
	if metric.Delta == nil {
		return "n/a"
	}
	value := format(*metric.Delta)
	if metric.RelativeDelta != nil {
		value += fmt.Sprintf(" (%+.1f%%)", *metric.RelativeDelta*100)
	}
	return value
}

func formatBytes(value int64) string         { return fmt.Sprintf("%d B", value) }
func formatInteger(value int64) string       { return fmt.Sprintf("%d", value) }
func formatSignedBytes(value int64) string   { return fmt.Sprintf("%+d B", value) }
func formatSignedInteger(value int64) string { return fmt.Sprintf("%+d", value) }

func formatDuration(value int64) string { return duration(time.Duration(value), false) }

func formatSignedDuration(value int64) string { return duration(time.Duration(value), true) }

func duration(value time.Duration, signed bool) string {
	sign := ""
	if signed {
		sign = "+"
		if value < 0 {
			sign = "-"
			value = -value
		}
	}
	switch {
	case value >= time.Second:
		return fmt.Sprintf("%s%.2fs", sign, value.Seconds())
	case value >= time.Millisecond:
		return fmt.Sprintf("%s%.1fms", sign, float64(value)/float64(time.Millisecond))
	default:
		return fmt.Sprintf("%s%.1fµs", sign, float64(value)/float64(time.Microsecond))
	}
}
