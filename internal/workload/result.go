package workload

// BenchmarkResultSchemaVersion identifies the JSON benchmark result contract.
// It is independent of the workload-definition and result-store schemas.
const BenchmarkResultSchemaVersion = 1

// BenchmarkResult is the canonical externally meaningful result of one
// workload run. Human and JSON formatters both consume this model.
//
// ShowBytes and RawRetrievalBytes are explicit zero-valued B4 fields because
// B0-B3 do not measure recovery operations. B5 may populate them without
// changing the meaning of the existing fields.
type BenchmarkResult struct {
	SchemaVersion int    `json:"schema_version"`
	Workload      string `json:"workload"`
	Commands      int    `json:"commands"`

	RawBytes       int64 `json:"raw_bytes"`
	StatelessBytes int64 `json:"stateless_bytes"`
	StatefulBytes  int64 `json:"stateful_bytes"`

	ShowBytes         int64 `json:"show_bytes"`
	RawRetrievalBytes int64 `json:"raw_retrieval_bytes"`

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
		r.ProcessingNS = run.Aggregate.ProcessingDuration.Nanoseconds()
	}
	return r
}
