package workload

// BenchmarkResultSchemaVersion identifies the JSON benchmark result contract.
// It is independent of the workload-definition and result-store schemas.
const BenchmarkResultSchemaVersion = 2

// BenchmarkResult is the canonical externally meaningful result of one
// workload run. Human and JSON formatters both consume this model.
type BenchmarkResult struct {
	SchemaVersion int    `json:"schema_version"`
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

	// presentation contains B3 details used only by the human/verbose view.
	presentation *Result
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
