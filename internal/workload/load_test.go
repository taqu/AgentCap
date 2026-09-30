package workload

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDecodeResultRoundTripsWrittenResults(t *testing.T) {
	ok := true
	exit := 0
	single := &BenchmarkResult{SchemaVersion: BenchmarkResultSchemaVersion, Workload: "w", Commands: 2,
		TotalVisibleBytes: 10, Agent: "codex", Mode: "stateful", TaskSuccess: &ok,
		ExecutionStatus: "completed", AgentExitCode: &exit, WallTimeNS: 5}
	var out bytes.Buffer
	if err := WriteJSON(&out, single); err != nil {
		t.Fatal(err)
	}
	got, err := DecodeResult(out.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if got.Kind != ResultKindSingle || got.Single.TotalVisibleBytes != 10 || *got.Single.TaskSuccess != true {
		t.Fatalf("got %+v", got.Single)
	}

	visible := int64(7)
	repeated := &RepeatedBenchmarkResult{SchemaVersion: BenchmarkResultSchemaVersion, Workload: "w", Agent: "codex",
		Mode: "stateful", RequestedTrialCount: 1, TrialCount: 1, RunStatus: "completed",
		Aggregate: RepeatedBenchmarkAggregate{MedianTotalVisibleBytes: &visible}, Trials: []*BenchmarkResult{single}}
	out.Reset()
	if err := WriteRepeatedJSON(&out, repeated); err != nil {
		t.Fatal(err)
	}
	got, err = DecodeResult(out.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if got.Kind != ResultKindRepeated || *got.Repeated.Aggregate.MedianTotalVisibleBytes != 7 || got.Repeated.Aggregate.MedianCommandCount != nil {
		t.Fatalf("got %+v", got.Repeated)
	}
}

func TestDecodeResultSchemaVersions(t *testing.T) {
	singleFields := `"workload":"w","commands":1,"initial_visible_bytes":1,"show_bytes":0,"raw_retrieval_bytes":0,"total_visible_bytes":1,"show_count":0,"raw_retrieval_count":0,"processing_ns":1`
	repeatedFields := `"workload":"w","agent":"codex","mode":"disabled","requested_trial_count":1,"trial_count":1,"success_count":1,"task_failure_count":0,"timeout_count":0,"agent_error_count":0,"canceled_count":0,"run_status":"completed","trials":[],"aggregate":{"median_total_visible_bytes":1,"median_command_count":1,"median_wall_time_ns":1,"median_processing_ns":0,"trials_with_show":0,"trials_with_raw_retrieval":0,"total_show_count":0,"total_raw_retrieval_count":0}`
	cases := []struct {
		name, json, wantErr string
	}{
		{"single v2", `{"schema_version":2,` + singleFields + `}`, ""},
		{"single v3", `{"schema_version":3,` + singleFields + `}`, ""},
		{"single v4", `{"schema_version":4,` + singleFields + `}`, ""},
		{"single v1", `{"schema_version":1,"workload":"w","commands":1}`, "unsupported benchmark result schema version 1"},
		{"future", `{"schema_version":99,` + singleFields + `}`, "unsupported benchmark result schema version 99"},
		{"repeated v4", `{"schema_version":4,` + repeatedFields + `}`, ""},
		{"repeated v3", `{"schema_version":3,` + repeatedFields + `}`, "unsupported repeated"},
		{"no version", `{` + singleFields + `}`, "schema_version"},
		{"invalid", `{"schema_version":`, "invalid benchmark result JSON"},
		{"missing field", `{"schema_version":4,"workload":"w"}`, "missing required field"},
		{"agent missing wall", `{"schema_version":4,"agent":"codex","mode":"disabled","task_success":true,"execution_status":"completed",` + singleFields + `}`, "wall_time_ns"},
		{"repeated missing aggregate field", `{"schema_version":4,` + strings.Replace(repeatedFields, `"median_wall_time_ns":1,`, "", 1) + `}`, "median_wall_time_ns"},
	}
	for _, tc := range cases {
		_, err := DecodeResult([]byte(tc.json))
		if tc.wantErr == "" && err != nil {
			t.Errorf("%s: unexpected error %v", tc.name, err)
		}
		if tc.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tc.wantErr)) {
			t.Errorf("%s: err = %v, want %q", tc.name, err, tc.wantErr)
		}
	}
}

func TestLoadResultMissingFile(t *testing.T) {
	if _, err := LoadResult(filepath.Join(t.TempDir(), "missing.json")); !os.IsNotExist(err) {
		t.Fatalf("err = %v", err)
	}
}
