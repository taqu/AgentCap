package workload

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/taqu/agentcap/internal/bench"
)

func sampleBenchmarkResult() (*Result, *BenchmarkResult) {
	run := &Result{
		Name:                       "git/repeated-diff",
		SessionID:                  "bench-test",
		TotalSteps:                 5,
		MutationCount:              2,
		ShowCount:                  2,
		RawRetrievalCount:          1,
		ShowBytes:                  50,
		RawRetrievalBytes:          25,
		RecoveryProcessingDuration: time.Millisecond,
		WallDuration:               9 * time.Millisecond,
		Aggregate: &bench.SessionMeasurement{
			CommandCount:       3,
			RawBytes:           1000,
			StatelessBytes:     400,
			StatefulBytes:      200,
			FullCount:          1,
			DeltaCount:         1,
			UnchangedCount:     1,
			ExecutionDuration:  5 * time.Millisecond,
			ReduceDuration:     time.Millisecond,
			ProcessingDuration: 2 * time.Millisecond,
		},
	}
	return run, NewBenchmarkResult(run)
}

func TestBenchmarkResultSchemaAndJSONFields(t *testing.T) {
	_, result := sampleBenchmarkResult()
	if result.SchemaVersion != BenchmarkResultSchemaVersion {
		t.Fatalf("schema version = %d, constant = %d", result.SchemaVersion, BenchmarkResultSchemaVersion)
	}
	if BenchmarkResultSchemaVersion != 2 {
		t.Fatalf("schema version constant = %d", BenchmarkResultSchemaVersion)
	}
	var buf bytes.Buffer
	if err := WriteJSON(&buf, result); err != nil {
		t.Fatal(err)
	}
	dec := json.NewDecoder(&buf)
	dec.UseNumber()
	var got map[string]any
	if err := dec.Decode(&got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	want := map[string]string{
		"schema_version":        "2",
		"commands":              "3",
		"raw_bytes":             "1000",
		"stateless_bytes":       "400",
		"stateful_bytes":        "200",
		"initial_visible_bytes": "200",
		"show_bytes":            "50",
		"raw_retrieval_bytes":   "25",
		"total_visible_bytes":   "275",
		"show_count":            "2",
		"raw_retrieval_count":   "1",
		"processing_ns":         "3000000",
	}
	if got["workload"] != "git/repeated-diff" {
		t.Errorf("workload = %v", got["workload"])
	}
	for field, value := range want {
		n, ok := got[field].(json.Number)
		if !ok || n.String() != value {
			t.Errorf("%s = %#v, want numeric %s", field, got[field], value)
		}
	}
	if len(got) != len(want)+1 {
		t.Errorf("unexpected schema fields: %#v", got)
	}
}

func TestHumanAndJSONUseCanonicalResult(t *testing.T) {
	run, result := sampleBenchmarkResult()
	// Changing internal byte totals after construction must not make either
	// formatter independently reconstruct the stable measurements.
	run.Aggregate.RawBytes = 9999
	var human, machine bytes.Buffer
	if err := WriteHuman(&human, result, false); err != nil {
		t.Fatal(err)
	}
	if err := WriteJSON(&machine, result); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(human.String(), "raw:                   1000 B") || !strings.Contains(human.String(), "total:                 275 B") {
		t.Fatalf("human output did not use canonical result:\n%s", human.String())
	}
	for _, want := range []string{"initial:               200 B", "show:                  50 B", "raw retrieval:         25 B", "effective vs raw:", "show calls:            2", "raw retrievals:        1"} {
		if !strings.Contains(human.String(), want) {
			t.Errorf("human output missing %q:\n%s", want, human.String())
		}
	}
	if !strings.Contains(machine.String(), `"raw_bytes": 1000`) {
		t.Fatalf("JSON output did not use canonical result:\n%s", machine.String())
	}
}
