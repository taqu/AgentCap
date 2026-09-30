package benchmarkcompare

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/taqu/agentcap/internal/workload"
)

func pointer(value int64) *int64 { return &value }

func repeated(mode string, trials, successes int, visible, commands, wall, processing *int64) *Result {
	return &Result{Kind: KindRepeated, Repeated: &workload.RepeatedBenchmarkResult{
		SchemaVersion: workload.BenchmarkResultSchemaVersion,
		Workload:      "agent/sample", Agent: "codex", Mode: mode,
		RequestedTrialCount: trials, TrialCount: trials, SuccessCount: successes,
		TaskFailureCount: trials - successes, RunStatus: "completed",
		Aggregate: workload.RepeatedBenchmarkAggregate{
			MedianTotalVisibleBytes: visible, MedianCommandCount: commands,
			MedianWallTimeNS: wall, MedianProcessingNS: processing,
			TrialsWithShow: 2, TrialsWithRawRetrieval: 1,
			TotalShowCount: 4, TotalRawRetrievalCount: 1,
		},
	}}
}

func TestCompareRepeatedNumericSuccessAndDifferentTrialCounts(t *testing.T) {
	baseline := repeated("disabled", 5, 5, pointer(100), pointer(20), pointer(80_000_000_000), pointer(0))
	candidate := repeated("integrated", 10, 8, pointer(60), pointer(25), pointer(82_000_000_000), pointer(12_000_000))
	before := *baseline.Repeated
	comparison, err := Compare(baseline, candidate)
	if err != nil {
		t.Fatal(err)
	}
	if comparison.SchemaVersion != 1 || comparison.Statistic != "median" || comparison.Baseline.TrialCount != 5 || comparison.Candidate.TrialCount != 10 || comparison.Baseline.SuccessCount == nil || *comparison.Baseline.SuccessCount != 5 || comparison.Candidate.SuccessCount == nil || *comparison.Candidate.SuccessCount != 8 {
		t.Fatalf("comparison metadata = %+v", comparison)
	}
	visible := comparison.Metrics.TotalVisibleBytes
	if visible.Delta == nil || *visible.Delta != -40 || visible.RelativeDelta == nil || *visible.RelativeDelta != -0.4 {
		t.Fatalf("visible comparison = %+v", visible)
	}
	commands := comparison.Metrics.CommandCount
	if commands.Delta == nil || *commands.Delta != 5 {
		t.Fatalf("command comparison = %+v", commands)
	}
	if !reflect.DeepEqual(before, *baseline.Repeated) {
		t.Fatal("comparison modified source result")
	}
}

func TestCompareZeroBaselineAndMissingMetrics(t *testing.T) {
	baseline := repeated("disabled", 2, 2, pointer(0), pointer(0), nil, nil)
	candidate := repeated("stateful", 2, 2, pointer(10), pointer(1), pointer(5), pointer(0))
	comparison, err := Compare(baseline, candidate)
	if err != nil {
		t.Fatal(err)
	}
	visible := comparison.Metrics.TotalVisibleBytes
	if visible.Delta == nil || *visible.Delta != 10 || visible.RelativeDelta != nil {
		t.Fatalf("zero-baseline comparison = %+v", visible)
	}
	if comparison.Metrics.WallTimeNS.Delta != nil || comparison.Metrics.WallTimeNS.Baseline != nil || comparison.Metrics.ProcessingNS.Delta != nil {
		t.Fatalf("missing metric was treated as zero: wall=%+v processing=%+v", comparison.Metrics.WallTimeNS, comparison.Metrics.ProcessingNS)
	}
}

func TestCompareRecoveryAndFormatting(t *testing.T) {
	baseline := repeated("disabled", 5, 5, pointer(100), pointer(20), pointer(1_000_000_000), pointer(0))
	baseline.Repeated.Aggregate.TrialsWithShow = 0
	baseline.Repeated.Aggregate.TrialsWithRawRetrieval = 0
	baseline.Repeated.Aggregate.TotalShowCount = 0
	baseline.Repeated.Aggregate.TotalRawRetrievalCount = 0
	candidate := repeated("integrated", 5, 4, pointer(60), pointer(25), pointer(2_000_000_000), pointer(2_000_000))
	comparison, err := Compare(baseline, candidate)
	if err != nil {
		t.Fatal(err)
	}
	if delta := comparison.Metrics.TrialsWithShow.Delta; delta == nil || *delta != 2 {
		t.Fatalf("show comparison = %+v", comparison.Metrics.TrialsWithShow)
	}
	if delta := comparison.Metrics.TotalRawRetrievalCount.Delta; delta == nil || *delta != 1 {
		t.Fatalf("raw comparison = %+v", comparison.Metrics.TotalRawRetrievalCount)
	}
	var human, machine bytes.Buffer
	if err := WriteHuman(&human, comparison); err != nil {
		t.Fatal(err)
	}
	if err := WriteJSON(&machine, comparison); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Task success", "Median visible bytes", "Median commands", "Median wall time", "Show used", "Raw used", "AgentCap processing", "-40 B"} {
		if !strings.Contains(human.String(), want) {
			t.Errorf("human output missing %q:\n%s", want, human.String())
		}
	}
	if strings.Contains(strings.ToLower(human.String()), "winner") || strings.Contains(human.String(), "PASS") || strings.Contains(human.String(), "FAIL") {
		t.Fatalf("human output assigns a judgment:\n%s", human.String())
	}
	var decoded map[string]any
	if err := json.Unmarshal(machine.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid comparison JSON: %v\n%s", err, machine.String())
	}
	if decoded["schema_version"] != float64(1) || decoded["statistic"] != "median" {
		t.Fatalf("JSON metadata = %#v", decoded)
	}
}

func TestCompareCompatibility(t *testing.T) {
	baseline := repeated("disabled", 1, 1, pointer(1), pointer(1), pointer(1), pointer(0))
	candidate := repeated("stateful", 1, 1, pointer(1), pointer(1), pointer(1), pointer(1))
	if _, err := Compare(baseline, candidate); err != nil {
		t.Fatalf("different modes should compare: %v", err)
	}
	candidate.Repeated.Workload = "other"
	if _, err := Compare(baseline, candidate); err == nil || !strings.Contains(err.Error(), "workloads") {
		t.Fatalf("incompatible workloads error = %v", err)
	}
	candidate.Repeated.Workload = baseline.Repeated.Workload
	candidate.Repeated.Agent = "claude"
	if _, err := Compare(baseline, candidate); err == nil || !strings.Contains(err.Error(), "agents") {
		t.Fatalf("incompatible agents error = %v", err)
	}
	single := &Result{Kind: KindSingle, Single: &workload.BenchmarkResult{Workload: baseline.Repeated.Workload, Agent: "codex"}}
	if _, err := Compare(baseline, single); err == nil || !strings.Contains(err.Error(), "kinds") {
		t.Fatalf("mixed result kinds error = %v", err)
	}
}

func TestSingleDeterministicSuccessAndWallTimeAreUnavailable(t *testing.T) {
	baseline := &Result{Kind: KindSingle, Single: &workload.BenchmarkResult{SchemaVersion: 4, Workload: "deterministic/sample", TotalVisibleBytes: 10}}
	candidate := &Result{Kind: KindSingle, Single: &workload.BenchmarkResult{SchemaVersion: 4, Workload: "deterministic/sample", TotalVisibleBytes: 9}}
	comparison, err := Compare(baseline, candidate)
	if err != nil {
		t.Fatal(err)
	}
	if comparison.Statistic != "observation" || comparison.Baseline.SuccessCount != nil || comparison.Metrics.WallTimeNS.Baseline != nil || comparison.Metrics.WallTimeNS.Delta != nil {
		t.Fatalf("unavailable fields were fabricated: %+v", comparison)
	}
	var human bytes.Buffer
	if err := WriteHuman(&human, comparison); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(human.String(), "Task success") || !strings.Contains(human.String(), "n/a") {
		t.Fatalf("unavailable fields not shown clearly:\n%s", human.String())
	}
}
