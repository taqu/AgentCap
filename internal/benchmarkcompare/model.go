// Package benchmarkcompare loads and compares existing benchmark results.
// It is deliberately independent of benchmark and agent execution.
package benchmarkcompare

import "github.com/taqu/agentcap/internal/workload"

const SchemaVersion = 1

type ResultKind string

const (
	KindSingle   ResultKind = "single"
	KindRepeated ResultKind = "repeated"
)

// Result is one typed, validated benchmark result loaded from JSON.
type Result struct {
	Kind     ResultKind
	Single   *workload.BenchmarkResult
	Repeated *workload.RepeatedBenchmarkResult

	singleWallTimeAvailable bool
}

type Side struct {
	Kind                ResultKind `json:"kind"`
	ResultSchemaVersion int        `json:"result_schema_version"`
	Workload            string     `json:"workload"`
	Agent               string     `json:"agent,omitempty"`
	Mode                string     `json:"mode,omitempty"`
	TrialCount          int        `json:"trial_count"`
	RequestedTrialCount int        `json:"requested_trial_count"`
	RunStatus           string     `json:"run_status"`
	SuccessCount        *int       `json:"success_count"`
	TaskFailureCount    *int       `json:"task_failure_count"`
	TimeoutCount        *int       `json:"timeout_count"`
	AgentErrorCount     *int       `json:"agent_error_count"`
	CanceledCount       *int       `json:"canceled_count"`
}

// NumericComparison uses candidate - baseline for every delta. RelativeDelta
// is nil if either input is unavailable or the baseline is zero.
type NumericComparison struct {
	Baseline      *int64   `json:"baseline"`
	Candidate     *int64   `json:"candidate"`
	Delta         *int64   `json:"delta"`
	RelativeDelta *float64 `json:"relative_delta"`
}

type Metrics struct {
	TotalVisibleBytes NumericComparison `json:"total_visible_bytes"`
	CommandCount      NumericComparison `json:"command_count"`
	WallTimeNS        NumericComparison `json:"wall_time_ns"`
	ProcessingNS      NumericComparison `json:"processing_ns"`

	TrialsWithShow         NumericComparison `json:"trials_with_show"`
	TrialsWithRawRetrieval NumericComparison `json:"trials_with_raw_retrieval"`
	TotalShowCount         NumericComparison `json:"total_show_count"`
	TotalRawRetrievalCount NumericComparison `json:"total_raw_retrieval_count"`
}

// Comparison is the canonical read-only analysis consumed by both formatters.
type Comparison struct {
	SchemaVersion int     `json:"schema_version"`
	Statistic     string  `json:"statistic"`
	Baseline      Side    `json:"baseline"`
	Candidate     Side    `json:"candidate"`
	Metrics       Metrics `json:"metrics"`
}
