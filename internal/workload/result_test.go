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
		Name:          "git/repeated-diff",
		SessionID:     "bench-test",
		TotalSteps:    5,
		MutationCount: 2,
		WallDuration:  9 * time.Millisecond,
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
	if BenchmarkResultSchemaVersion != 1 {
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
		"schema_version":      "1",
		"commands":            "3",
		"raw_bytes":           "1000",
		"stateless_bytes":     "400",
		"stateful_bytes":      "200",
		"show_bytes":          "0",
		"raw_retrieval_bytes": "0",
		"processing_ns":       "2000000",
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
	if !strings.Contains(human.String(), "raw:                  1000 B") {
		t.Fatalf("human output did not use canonical result:\n%s", human.String())
	}
	if !strings.Contains(machine.String(), `"raw_bytes": 1000`) {
		t.Fatalf("JSON output did not use canonical result:\n%s", machine.String())
	}
}
