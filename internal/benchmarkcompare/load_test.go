package benchmarkcompare

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/taqu/agentcap/internal/workload"
)

func TestLoadSupportedSingleSchemas(t *testing.T) {
	for _, version := range []int{3, workload.BenchmarkResultSchemaVersion} {
		data, err := json.Marshal(&workload.BenchmarkResult{SchemaVersion: version, Workload: "deterministic/sample"})
		if err != nil {
			t.Fatal(err)
		}
		result, err := Load(data)
		if err != nil || result.Kind != KindSingle || result.Single.SchemaVersion != version {
			t.Fatalf("schema %d: result=%+v err=%v", version, result, err)
		}
	}
}

func TestLoadRepeatedAndValidateTrials(t *testing.T) {
	success := true
	zero := int64(0)
	trial := &workload.BenchmarkResult{
		SchemaVersion: workload.BenchmarkResultSchemaVersion, Trial: 1,
		Workload: "agent/sample", Agent: "codex", Mode: "integrated",
		TaskSuccess: &success, ExecutionStatus: "completed", WallTimeNS: 1,
	}
	value := &workload.RepeatedBenchmarkResult{
		SchemaVersion: workload.BenchmarkResultSchemaVersion,
		Workload:      "agent/sample", Agent: "codex", Mode: "integrated",
		RequestedTrialCount: 1, TrialCount: 1, SuccessCount: 1, RunStatus: "completed",
		Aggregate: workload.RepeatedBenchmarkAggregate{
			MedianTotalVisibleBytes: &zero, MedianCommandCount: &zero,
			MedianWallTimeNS: &zero, MedianProcessingNS: &zero,
		},
		Trials: []*workload.BenchmarkResult{trial},
	}
	data, _ := json.Marshal(value)
	result, err := Load(data)
	if err != nil || result.Kind != KindRepeated || result.Repeated.TrialCount != 1 {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	trial.Trial = 2
	data, _ = json.Marshal(value)
	if _, err := Load(data); err == nil || !strings.Contains(err.Error(), "index") {
		t.Fatalf("bad trial index error = %v", err)
	}
}

func TestLoadRejectsMalformedUnsupportedAndMissing(t *testing.T) {
	for _, tc := range []struct {
		name string
		data string
		want string
	}{
		{"malformed", `{`, "invalid benchmark JSON"},
		{"unsupported", `{"schema_version":2,"workload":"w"}`, "unsupported single"},
		{"missing", `{"schema_version":4,"workload":"w"}`, "missing required field"},
		{"negative", `{"schema_version":4,"workload":"w","commands":-1,"raw_bytes":0,"stateless_bytes":0,"stateful_bytes":0,"initial_visible_bytes":0,"show_bytes":0,"raw_retrieval_bytes":0,"total_visible_bytes":0,"show_count":0,"raw_retrieval_count":0,"processing_ns":0}`, "must not be negative"},
		{"unknown", `{"schema_version":4,"workload":"w","commands":0,"raw_bytes":0,"stateless_bytes":0,"stateful_bytes":0,"initial_visible_bytes":0,"show_bytes":0,"raw_retrieval_bytes":0,"total_visible_bytes":0,"show_count":0,"raw_retrieval_count":0,"processing_ns":0,"surprise":1}`, "unknown field"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Load([]byte(tc.data))
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want %q", err, tc.want)
			}
		})
	}
}
