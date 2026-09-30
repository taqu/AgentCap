package workload

// BenchmarkResultSchemaVersion identifies the JSON benchmark result contract.
// It is independent of the workload-definition and result-store schemas.
const BenchmarkResultSchemaVersion = 4

// BenchmarkResult is the canonical externally meaningful result of one
// workload run. Human and JSON formatters both consume this model.
type BenchmarkResult struct {
	SchemaVersion int    `json:"schema_version"`
	Trial         int    `json:"trial,omitempty"`
	Workload      string `json:"workload"`
	Commands      int    `json:"commands"`

	RawBytes       int64 `json:"raw_bytes"`
	StatelessBytes int64 `json:"stateless_bytes"`
	StatefulBytes  int64 `json:"stateful_bytes"`

	InitialVisibleBytes int64 `json:"initial_visible_bytes"`
	ShowBytes           int64 `json:"show_bytes"`
	RawRetrievalBytes   int64 `json:"raw_retrieval_bytes"`
	TotalVisibleBytes   int64 `json:"total_visible_bytes"`
	ShowCount           int   `json:"show_count"`
	RawRetrievalCount   int   `json:"raw_retrieval_count"`

	ProcessingNS int64 `json:"processing_ns"`

	Agent           string `json:"agent,omitempty"`
	Mode            string `json:"mode,omitempty"`
	TaskSuccess     *bool  `json:"task_success,omitempty"`
	ExecutionStatus string `json:"execution_status,omitempty"`
	AgentExitCode   *int   `json:"agent_exit_code,omitempty"`
	WallTimeNS      int64  `json:"wall_time_ns,omitempty"`

	Workflow *WorkflowMetrics `json:"workflow,omitempty"`

	// presentation contains B3 details used only by the human/verbose view.
	presentation *Result
}

// RepeatedBenchmarkAggregate is the canonical statistical summary of all
// valid executed trials in one fixed workload/agent/mode configuration.
// Integer medians use the conventional midpoint, rounded down when it falls
// between two integer values.
type RepeatedBenchmarkAggregate struct {
	MedianTotalVisibleBytes *int64 `json:"median_total_visible_bytes"`
	MedianCommandCount      *int64 `json:"median_command_count"`
	MedianWallTimeNS        *int64 `json:"median_wall_time_ns"`
	MedianProcessingNS      *int64 `json:"median_processing_ns"`

	TrialsWithShow         int `json:"trials_with_show"`
	TrialsWithRawRetrieval int `json:"trials_with_raw_retrieval"`
	TotalShowCount         int `json:"total_show_count"`
	TotalRawRetrievalCount int `json:"total_raw_retrieval_count"`
}

// RepeatedBenchmarkResult preserves every B6 trial together with the B7
// aggregate for one fixed benchmark configuration.
type RepeatedBenchmarkResult struct {
	SchemaVersion int    `json:"schema_version"`
	Workload      string `json:"workload"`
	Agent         string `json:"agent"`
	Mode          string `json:"mode"`

	RequestedTrialCount int    `json:"requested_trial_count"`
	TrialCount          int    `json:"trial_count"`
	SuccessCount        int    `json:"success_count"`
	TaskFailureCount    int    `json:"task_failure_count"`
	TimeoutCount        int    `json:"timeout_count"`
	AgentErrorCount     int    `json:"agent_error_count"`
	CanceledCount       int    `json:"canceled_count"`
	RunStatus           string `json:"run_status"`

	Aggregate RepeatedBenchmarkAggregate `json:"aggregate"`
	Trials    []*BenchmarkResult         `json:"trials"`
}

// NewBenchmarkResult transforms internal B3 measurements into the stable B4
// result contract. It performs no execution or measurement.
func NewBenchmarkResult(run *Result) *BenchmarkResult {
	r := &BenchmarkResult{
		SchemaVersion: BenchmarkResultSchemaVersion,
		presentation:  run,
	}
	if run == nil {
		return r
	}
	r.Workload = run.Name
	if run.Aggregate != nil {
		r.Commands = run.Aggregate.CommandCount
		r.RawBytes = run.Aggregate.RawBytes
		r.StatelessBytes = run.Aggregate.StatelessBytes
		r.StatefulBytes = run.Aggregate.StatefulBytes
		r.InitialVisibleBytes = run.Aggregate.StatefulBytes
		r.ShowBytes = run.ShowBytes
		r.RawRetrievalBytes = run.RawRetrievalBytes
		r.TotalVisibleBytes = r.InitialVisibleBytes + r.ShowBytes + r.RawRetrievalBytes
		r.ShowCount = run.ShowCount
		r.RawRetrievalCount = run.RawRetrievalCount
		r.ProcessingNS = (run.Aggregate.ProcessingDuration + run.RecoveryProcessingDuration).Nanoseconds()
	}
	return r
}
