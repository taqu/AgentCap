package benchmarkcompare

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/taqu/agentcap/internal/workload"
)

// LoadFile reads a benchmark JSON result without executing any benchmark code.
func LoadFile(path string) (*Result, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	result, err := Load(data)
	if err != nil {
		return nil, fmt.Errorf("load %s: %w", path, err)
	}
	return result, nil
}

// Load decodes the stable B4-B7 JSON result types supported by comparison.
func Load(data []byte) (*Result, error) {
	var fields map[string]json.RawMessage
	if err := decodeStrict(data, &fields); err != nil {
		return nil, fmt.Errorf("invalid benchmark JSON: %w", err)
	}
	if err := require(fields, "schema_version", "workload"); err != nil {
		return nil, err
	}
	var version int
	if err := json.Unmarshal(fields["schema_version"], &version); err != nil {
		return nil, errors.New("schema_version must be an integer")
	}
	_, hasTrials := fields["trials"]
	_, hasAggregate := fields["aggregate"]
	_, hasTrialCount := fields["trial_count"]
	if hasTrials || hasAggregate || hasTrialCount {
		if version != 4 {
			return nil, fmt.Errorf("unsupported repeated benchmark schema version %d (supported: 4)", version)
		}
		return loadRepeated(data, fields)
	}
	if version != 3 && version != 4 {
		return nil, fmt.Errorf("unsupported single benchmark schema version %d (supported: 3 and 4)", version)
	}
	return loadSingle(data, fields)
}

func loadSingle(data []byte, fields map[string]json.RawMessage) (*Result, error) {
	if err := require(fields, "commands", "raw_bytes", "stateless_bytes", "stateful_bytes", "initial_visible_bytes", "show_bytes", "raw_retrieval_bytes", "total_visible_bytes", "show_count", "raw_retrieval_count", "processing_ns"); err != nil {
		return nil, err
	}
	var value workload.BenchmarkResult
	if err := decodeStrict(data, &value); err != nil {
		return nil, fmt.Errorf("invalid single benchmark result: %w", err)
	}
	if value.Workload == "" {
		return nil, errors.New("workload must not be empty")
	}
	if value.Agent != "" {
		if err := require(fields, "mode", "task_success", "execution_status", "wall_time_ns"); err != nil {
			return nil, err
		}
		if value.Mode == "" || value.TaskSuccess == nil || value.ExecutionStatus == "" {
			return nil, errors.New("coding-agent result has incomplete outcome metadata")
		}
	}
	if value.Commands < 0 || value.RawBytes < 0 || value.StatelessBytes < 0 || value.StatefulBytes < 0 || value.InitialVisibleBytes < 0 || value.ShowBytes < 0 || value.RawRetrievalBytes < 0 || value.TotalVisibleBytes < 0 || value.ShowCount < 0 || value.RawRetrievalCount < 0 || value.ProcessingNS < 0 || value.WallTimeNS < 0 {
		return nil, errors.New("benchmark measurements must not be negative")
	}
	_, wallTimeAvailable := fields["wall_time_ns"]
	return &Result{Kind: KindSingle, Single: &value, singleWallTimeAvailable: wallTimeAvailable}, nil
}

func loadRepeated(data []byte, fields map[string]json.RawMessage) (*Result, error) {
	if err := require(fields, "agent", "mode", "requested_trial_count", "trial_count", "success_count", "task_failure_count", "timeout_count", "agent_error_count", "canceled_count", "run_status", "aggregate", "trials"); err != nil {
		return nil, err
	}
	var value workload.RepeatedBenchmarkResult
	if err := decodeStrict(data, &value); err != nil {
		return nil, fmt.Errorf("invalid repeated benchmark result: %w", err)
	}
	if value.Workload == "" || value.Agent == "" || value.Mode == "" || value.RunStatus == "" {
		return nil, errors.New("repeated benchmark has incomplete metadata")
	}
	if value.RequestedTrialCount < 1 || value.TrialCount < 0 || value.TrialCount > value.RequestedTrialCount || value.TrialCount != len(value.Trials) {
		return nil, errors.New("repeated benchmark has inconsistent trial counts")
	}
	if value.SuccessCount < 0 || value.SuccessCount > value.TrialCount || value.TaskFailureCount < 0 || value.TaskFailureCount > value.TrialCount || value.TimeoutCount < 0 || value.TimeoutCount > value.TrialCount || value.AgentErrorCount < 0 || value.AgentErrorCount > value.TrialCount || value.CanceledCount < 0 || value.CanceledCount > value.TrialCount {
		return nil, errors.New("repeated benchmark has invalid outcome counts")
	}
	if value.SuccessCount+value.TaskFailureCount != value.TrialCount {
		return nil, errors.New("repeated benchmark task outcomes do not match trial count")
	}
	if value.RunStatus != "completed" && value.RunStatus != "incomplete" {
		return nil, fmt.Errorf("unsupported repeated run status %q", value.RunStatus)
	}
	if value.RunStatus == "completed" && value.TrialCount != value.RequestedTrialCount {
		return nil, errors.New("completed repeated benchmark is missing trials")
	}
	var aggregate map[string]json.RawMessage
	if err := json.Unmarshal(fields["aggregate"], &aggregate); err != nil {
		return nil, errors.New("aggregate must be an object")
	}
	if err := require(aggregate, "median_total_visible_bytes", "median_command_count", "median_wall_time_ns", "median_processing_ns", "trials_with_show", "trials_with_raw_retrieval", "total_show_count", "total_raw_retrieval_count"); err != nil {
		return nil, fmt.Errorf("aggregate: %w", err)
	}
	if negative(value.Aggregate.MedianTotalVisibleBytes) || negative(value.Aggregate.MedianCommandCount) || negative(value.Aggregate.MedianWallTimeNS) || negative(value.Aggregate.MedianProcessingNS) || value.Aggregate.TrialsWithShow < 0 || value.Aggregate.TrialsWithShow > value.TrialCount || value.Aggregate.TrialsWithRawRetrieval < 0 || value.Aggregate.TrialsWithRawRetrieval > value.TrialCount || value.Aggregate.TotalShowCount < 0 || value.Aggregate.TotalRawRetrievalCount < 0 {
		return nil, errors.New("aggregate has invalid measurements or recovery trial counts")
	}
	var rawTrials []json.RawMessage
	if err := json.Unmarshal(fields["trials"], &rawTrials); err != nil {
		return nil, errors.New("trials must be an array")
	}
	for i, raw := range rawTrials {
		loaded, err := Load(raw)
		if err != nil {
			return nil, fmt.Errorf("trial %d: %w", i+1, err)
		}
		if loaded.Kind != KindSingle || loaded.Single.SchemaVersion != 4 || loaded.Single.Workload != value.Workload || loaded.Single.Agent != value.Agent || loaded.Single.Mode != value.Mode {
			return nil, fmt.Errorf("trial %d metadata does not match repeated benchmark", i+1)
		}
		if loaded.Single.Trial != i+1 {
			return nil, fmt.Errorf("trial %d has index %d", i+1, loaded.Single.Trial)
		}
	}
	return &Result{Kind: KindRepeated, Repeated: &value}, nil
}

func decodeStrict(data []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("multiple JSON values")
		}
		return err
	}
	return nil
}

func require(fields map[string]json.RawMessage, names ...string) error {
	for _, name := range names {
		if _, ok := fields[name]; !ok {
			return fmt.Errorf("missing required field %q", name)
		}
	}
	return nil
}

func negative(value *int64) bool { return value != nil && *value < 0 }
