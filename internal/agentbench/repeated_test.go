package agentbench

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/taqu/agentcap/internal/workload"
)

type repeatedAdapter struct {
	successes []bool
	visible   []int64
	errorAt   int
	calls     int
	initial   []string
	roots     []string
	scopes    []string
	tasks     []string
	modes     []Mode
}

func (a *repeatedAdapter) Name() string { return "repeated-fake" }

func (a *repeatedAdapter) Run(_ context.Context, req AgentRunRequest) (*AgentRunResult, error) {
	a.calls++
	if a.errorAt == a.calls {
		return nil, errors.New("adapter infrastructure failure")
	}
	data, err := os.ReadFile(filepath.Join(req.Workspace, "solution.txt"))
	if err != nil {
		return nil, err
	}
	a.initial = append(a.initial, string(data))
	a.roots = append(a.roots, req.StoreRoot)
	a.scopes = append(a.scopes, req.SessionScope)
	a.tasks = append(a.tasks, req.Task)
	a.modes = append(a.modes, req.Mode)
	index := a.calls - 1
	if index < len(a.successes) && a.successes[index] {
		if err := os.WriteFile(filepath.Join(req.Workspace, "solution.txt"), []byte("fixed"), 0o644); err != nil {
			return nil, err
		}
	}
	visible := int64(a.calls * 10)
	if index < len(a.visible) {
		visible = a.visible[index]
	}
	return &AgentRunResult{ExitCode: 0, CommandCount: a.calls, CommandVisibleBytes: visible}, nil
}

func TestRunRepeatedIndependentTrialsAndAggregation(t *testing.T) {
	definition := makeAgentWorkload(t)
	// A caller-level root override must not collapse isolated trial stores.
	t.Setenv("ACAP_ROOT", t.TempDir())
	adapter := &repeatedAdapter{
		successes: []bool{true, false, true, true, false},
		visible:   []int64{10, 11, 12, 13, 1000},
	}
	tempRoot := t.TempDir()
	run, err := RunRepeated(context.Background(), definition, Options{
		Adapter: adapter, Mode: ModeDisabled, TempRoot: tempRoot,
	}, 5)
	if err != nil {
		t.Fatal(err)
	}
	if adapter.calls != 5 || len(run.Trials) != 5 || run.Benchmark.TrialCount != 5 || run.Benchmark.SuccessCount != 3 {
		t.Fatalf("calls/trials/success = %d/%d/%d/%d", adapter.calls, len(run.Trials), run.Benchmark.TrialCount, run.Benchmark.SuccessCount)
	}
	if run.Benchmark.Aggregate.MedianTotalVisibleBytes == nil || *run.Benchmark.Aggregate.MedianTotalVisibleBytes != 12 || run.Benchmark.Aggregate.MedianCommandCount == nil || *run.Benchmark.Aggregate.MedianCommandCount != 3 {
		t.Fatalf("aggregate = %+v", run.Benchmark.Aggregate)
	}
	seenRoots, seenScopes := map[string]bool{}, map[string]bool{}
	for i := range adapter.initial {
		if adapter.initial[i] != "broken" || adapter.tasks[i] != definition.Task || adapter.modes[i] != ModeDisabled {
			t.Fatalf("trial %d configuration/fixture = %q/%q/%q", i+1, adapter.initial[i], adapter.tasks[i], adapter.modes[i])
		}
		if adapter.roots[i] == "" || seenRoots[adapter.roots[i]] {
			t.Fatalf("trial %d reused state root %q", i+1, adapter.roots[i])
		}
		if adapter.scopes[i] == "" || seenScopes[adapter.scopes[i]] {
			t.Fatalf("trial %d reused session scope %q", i+1, adapter.scopes[i])
		}
		seenRoots[adapter.roots[i]], seenScopes[adapter.scopes[i]] = true, true
		if run.Trials[i].Benchmark.Trial != i+1 {
			t.Fatalf("trial index %d = %d", i+1, run.Trials[i].Benchmark.Trial)
		}
	}
	entries, err := os.ReadDir(tempRoot)
	if err != nil || len(entries) != 0 {
		t.Fatalf("trial roots not cleaned: %v, %v", entries, err)
	}
	if run.Benchmark.Trials[4].TotalVisibleBytes != 1000 {
		t.Fatal("outlier trial was not preserved")
	}
}

func TestMedianOddEvenAndMissing(t *testing.T) {
	if got := median([]int64{100, 10, 20}); got != 20 {
		t.Fatalf("odd median = %d", got)
	}
	if got := median([]int64{9, 1, 7, 3}); got != 5 {
		t.Fatalf("even median = %d", got)
	}
	if got := median([]int64{2, 3}); got != 2 {
		t.Fatalf("fractional integer median = %d, want rounded-down 2", got)
	}
	if got := median(nil); got != 0 {
		t.Fatalf("empty median = %d", got)
	}
}

func TestAggregateRecoveryFailurePopulationAndJSON(t *testing.T) {
	values := []struct {
		success   bool
		visible   int64
		show, raw int
		status    string
	}{
		{true, 90, 2, 0, "completed"},
		{true, 95, 0, 1, "completed"},
		{false, 500, 3, 2, "timeout"},
	}
	trials := make([]*Trial, 0, len(values)+1)
	for i, value := range values {
		success := value.success
		trials = append(trials, &Trial{Benchmark: &workload.BenchmarkResult{
			SchemaVersion: workload.BenchmarkResultSchemaVersion, Trial: i + 1,
			Workload: "w", Agent: "a", Mode: "integrated", TaskSuccess: &success,
			ExecutionStatus: value.status, Commands: i + 1, TotalVisibleBytes: value.visible,
			ShowCount: value.show, RawRetrievalCount: value.raw,
		}})
	}
	trials = append(trials, nil) // missing measurement is omitted, never a zero sample
	result, err := AggregateTrials("w", "a", ModeIntegrated, 4, trials, "incomplete")
	if err != nil {
		t.Fatal(err)
	}
	if result.TrialCount != 3 || result.SuccessCount != 2 || result.TaskFailureCount != 1 || result.TimeoutCount != 1 || result.Aggregate.MedianTotalVisibleBytes == nil || *result.Aggregate.MedianTotalVisibleBytes != 95 {
		t.Fatalf("result = %+v aggregate=%+v", result, result.Aggregate)
	}
	if result.Aggregate.TrialsWithShow != 2 || result.Aggregate.TrialsWithRawRetrieval != 2 || result.Aggregate.TotalShowCount != 5 || result.Aggregate.TotalRawRetrievalCount != 3 {
		t.Fatalf("recovery aggregate = %+v", result.Aggregate)
	}
	var machine, human bytes.Buffer
	if err := workload.WriteRepeatedJSON(&machine, result); err != nil {
		t.Fatal(err)
	}
	if err := workload.WriteRepeatedHuman(&human, result); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`"aggregate"`, `"trials"`, `"trial": 3`, `"total_visible_bytes": 500`} {
		if !bytes.Contains(machine.Bytes(), []byte(want)) {
			t.Errorf("JSON missing %q:\n%s", want, machine.String())
		}
	}
	if !bytes.Contains(human.Bytes(), []byte("success:               2/3")) || !bytes.Contains(human.Bytes(), []byte("500")) {
		t.Fatalf("human output missing aggregate/trials:\n%s", human.String())
	}
}

func TestRunRepeatedInfrastructureFailureReturnsPartialResult(t *testing.T) {
	definition := makeAgentWorkload(t)
	adapter := &repeatedAdapter{successes: []bool{true, true}, errorAt: 3}
	run, err := RunRepeated(context.Background(), definition, Options{Adapter: adapter, Mode: ModeDisabled}, 5)
	if err == nil || run == nil || run.Benchmark == nil {
		t.Fatalf("run=%+v err=%v", run, err)
	}
	if run.Benchmark.RunStatus != "incomplete" || run.Benchmark.RequestedTrialCount != 5 || run.Benchmark.TrialCount != 2 {
		t.Fatalf("partial result = %+v", run.Benchmark)
	}
}

func TestAggregateNoMeasurementsUsesNullMedians(t *testing.T) {
	result, err := AggregateTrials("w", "a", ModeDisabled, 2, nil, "incomplete")
	if err != nil {
		t.Fatal(err)
	}
	if result.Aggregate.MedianTotalVisibleBytes != nil || result.Aggregate.MedianCommandCount != nil || result.Aggregate.MedianWallTimeNS != nil || result.Aggregate.MedianProcessingNS != nil {
		t.Fatalf("missing medians were represented as measurements: %+v", result.Aggregate)
	}
	var machine bytes.Buffer
	if err := workload.WriteRepeatedJSON(&machine, result); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(machine.Bytes(), []byte(`"median_total_visible_bytes": null`)) {
		t.Fatalf("missing median is not JSON null:\n%s", machine.String())
	}
}

func TestRunRepeatedSingleTrialPreservesB6Result(t *testing.T) {
	adapter := &repeatedAdapter{successes: []bool{true}, visible: []int64{77}}
	run, err := RunRepeated(context.Background(), makeAgentWorkload(t), Options{Adapter: adapter, Mode: ModeDisabled}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(run.Trials) != 1 || run.Trials[0].Benchmark == nil || run.Trials[0].Benchmark.TotalVisibleBytes != 77 || run.Trials[0].Benchmark.Trial != 1 {
		t.Fatalf("single trial = %+v", run)
	}
}

func TestRunRepeatedRejectsInvalidCountAndIncompatibleTrials(t *testing.T) {
	if _, err := RunRepeated(context.Background(), makeAgentWorkload(t), Options{}, 0); err == nil {
		t.Fatal("repeat zero accepted")
	}
	success := true
	_, err := AggregateTrials("w", "a", ModeDisabled, 1, []*Trial{{Benchmark: &workload.BenchmarkResult{
		Trial: 1, Workload: "other", Agent: "a", Mode: "disabled", TaskSuccess: &success,
	}}}, "completed")
	if err == nil {
		t.Fatal("incompatible trial accepted")
	}
}
