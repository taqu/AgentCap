package workload

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ResultKind distinguishes the two serialized benchmark result shapes.
type ResultKind string

const (
	// ResultKindSingle is one BenchmarkResult: a deterministic workload run or
	// one coding-agent trial.
	ResultKindSingle ResultKind = "single"
	// ResultKindRepeated is a RepeatedBenchmarkResult with a B7 aggregate.
	ResultKindRepeated ResultKind = "repeated"
)

// Oldest schema versions whose fields keep the current meaning. Version 1
// predates recovery-aware total_visible_bytes; repeated results first
// appeared in version 4. Versions 2-4 only added fields to single results.
const (
	minSingleResultSchemaVersion   = 2
	minRepeatedResultSchemaVersion = 4
)

// StoredResult is one previously serialized benchmark result decoded into the
// canonical in-memory types. Exactly one of Single and Repeated is set.
type StoredResult struct {
	Kind     ResultKind
	Single   *BenchmarkResult
	Repeated *RepeatedBenchmarkResult
}

// SchemaVersion returns the result-schema version the result was written with.
func (r *StoredResult) SchemaVersion() int {
	if r.Repeated != nil {
		return r.Repeated.SchemaVersion
	}
	return r.Single.SchemaVersion
}

// LoadResult reads a result previously written by "acap bench run --json" or
// "acap bench agent --json". It never modifies the file.
func LoadResult(path string) (*StoredResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	result, err := DecodeResult(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return result, nil
}

// DecodeResult decodes serialized single or repeated benchmark JSON and
// rejects schema versions and shapes whose semantics are not supported.
func DecodeResult(data []byte) (*StoredResult, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, fmt.Errorf("invalid benchmark result JSON: %w", err)
	}
	var version int
	if raw, ok := fields["schema_version"]; !ok {
		return nil, errors.New("missing required field \"schema_version\"")
	} else if err := json.Unmarshal(raw, &version); err != nil {
		return nil, fmt.Errorf("invalid schema_version: %w", err)
	}
	if version > BenchmarkResultSchemaVersion {
		return nil, fmt.Errorf("unsupported benchmark result schema version %d (newest supported: %d)", version, BenchmarkResultSchemaVersion)
	}

	_, hasAggregate := fields["aggregate"]
	_, hasTrials := fields["trials"]
	if hasAggregate || hasTrials {
		if version < minRepeatedResultSchemaVersion {
			return nil, fmt.Errorf("unsupported repeated benchmark result schema version %d (supported: %d-%d)", version, minRepeatedResultSchemaVersion, BenchmarkResultSchemaVersion)
		}
		if err := requireFields(fields, "workload", "agent", "mode", "requested_trial_count", "trial_count",
			"success_count", "task_failure_count", "timeout_count", "agent_error_count", "canceled_count",
			"run_status", "aggregate", "trials"); err != nil {
			return nil, err
		}
		var aggregate map[string]json.RawMessage
		if err := json.Unmarshal(fields["aggregate"], &aggregate); err != nil {
			return nil, fmt.Errorf("invalid aggregate: %w", err)
		}
		if err := requireFields(aggregate, "median_total_visible_bytes", "median_command_count",
			"median_wall_time_ns", "median_processing_ns", "trials_with_show", "trials_with_raw_retrieval",
			"total_show_count", "total_raw_retrieval_count"); err != nil {
			return nil, fmt.Errorf("aggregate: %w", err)
		}
		var result RepeatedBenchmarkResult
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("invalid repeated benchmark result: %w", err)
		}
		return &StoredResult{Kind: ResultKindRepeated, Repeated: &result}, nil
	}

	if version < minSingleResultSchemaVersion {
		return nil, fmt.Errorf("unsupported benchmark result schema version %d (supported: %d-%d)", version, minSingleResultSchemaVersion, BenchmarkResultSchemaVersion)
	}
	if err := requireFields(fields, "workload", "commands", "initial_visible_bytes", "show_bytes",
		"raw_retrieval_bytes", "total_visible_bytes", "show_count", "raw_retrieval_count", "processing_ns"); err != nil {
		return nil, err
	}
	var result BenchmarkResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("invalid benchmark result: %w", err)
	}
	if result.Agent != "" {
		// Agent-trial fields are omitempty in the serialized form, so their
		// presence is checked explicitly rather than defaulted to zero.
		if err := requireFields(fields, "mode", "task_success", "execution_status", "wall_time_ns"); err != nil {
			return nil, fmt.Errorf("agent result: %w", err)
		}
	}
	return &StoredResult{Kind: ResultKindSingle, Single: &result}, nil
}

// LoadResultSet loads the results at path, keyed by workload name. path is
// either one result file or a directory whose *.json files (not recursive)
// are each one result, as written by "acap bench suite". An empty set and two
// results for the same workload are errors.
func LoadResultSet(path string) (map[string]*StoredResult, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	files := []string{path}
	if info.IsDir() {
		if files, err = filepath.Glob(filepath.Join(path, "*.json")); err != nil {
			return nil, err
		}
		sort.Strings(files)
	}
	set := make(map[string]*StoredResult, len(files))
	for _, file := range files {
		result, err := LoadResult(file)
		if err != nil {
			return nil, err
		}
		name := result.WorkloadName()
		if _, dup := set[name]; dup {
			return nil, fmt.Errorf("%s: more than one result for workload %q", path, name)
		}
		set[name] = result
	}
	if len(set) == 0 {
		return nil, fmt.Errorf("%s: no benchmark result files", path)
	}
	return set, nil
}

// WorkloadName returns the stable workload identity of the result.
func (r *StoredResult) WorkloadName() string {
	if r.Repeated != nil {
		return r.Repeated.Workload
	}
	return r.Single.Workload
}

// ResultFileName is the file name "acap bench suite" uses for a workload's
// result. Workload names may contain "/", which is replaced by "_".
func ResultFileName(workloadName string) string {
	return strings.ReplaceAll(workloadName, "/", "_") + ".json"
}

func requireFields(fields map[string]json.RawMessage, names ...string) error {
	for _, name := range names {
		if _, ok := fields[name]; !ok {
			return fmt.Errorf("missing required field %q", name)
		}
	}
	return nil
}
