package benchcompare

import (
	"bytes"
	"encoding/json"
	"errors"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/taqu/agentcap/internal/workload"
)

func ptr(v int64) *int64 { return &v }

func repeated(mode string, trials int, visible, commands int64) *workload.StoredResult {
	return &workload.StoredResult{Kind: workload.ResultKindRepeated, Repeated: &workload.RepeatedBenchmarkResult{
		SchemaVersion: workload.BenchmarkResultSchemaVersion,
		Workload:      "bugfix-001", Agent: "codex", Mode: mode,
		RequestedTrialCount: trials, TrialCount: trials, SuccessCount: trials, RunStatus: "completed",
		Aggregate: workload.RepeatedBenchmarkAggregate{
			MedianTotalVisibleBytes: ptr(visible),
			MedianCommandCount:      ptr(commands),
			MedianWallTimeNS:        ptr(84_200_000_000),
			MedianProcessingNS:      ptr(0),
		},
	}}
}

func single(agent, mode string, visible int64) *workload.StoredResult {
	r := &workload.BenchmarkResult{
		SchemaVersion: workload.BenchmarkResultSchemaVersion,
		Workload:      "bugfix-001", Agent: agent, Mode: mode,
		Commands: 3, TotalVisibleBytes: visible, ProcessingNS: 5_000,
	}
	if agent != "" {
		ok := true
		r.TaskSuccess, r.ExecutionStatus, r.WallTimeNS = &ok, "completed", 2_000_000_000
	}
	return &workload.StoredResult{Kind: workload.ResultKindSingle, Single: r}
}

func mustCompare(t *testing.T, b, k *workload.StoredResult) *Comparison {
	t.Helper()
	c, err := Compare(b, k)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestNumericDeltaIsCandidateMinusBaseline(t *testing.T) {
	c := mustCompare(t, repeated("stateful", 5, 100, 20), repeated("integrated", 5, 60, 25))
	visible := c.Metric("median_total_visible_bytes")
	if *visible.Baseline != 100 || *visible.Candidate != 60 || *visible.Delta != -40 || *visible.RelativeDelta != -0.4 {
		t.Fatalf("visible = %+v", visible)
	}
	commands := c.Metric("median_command_count")
	if *commands.Delta != 5 || *commands.RelativeDelta != 0.25 {
		t.Fatalf("commands = %+v", commands)
	}
}

func TestZeroBaselineHasNoRelativeDelta(t *testing.T) {
	b, k := repeated("stateful", 5, 0, 20), repeated("integrated", 5, 10, 20)
	k.Repeated.Aggregate.TotalRawRetrievalCount = 1
	c := mustCompare(t, b, k)
	for _, name := range []string{"median_total_visible_bytes", "total_raw_retrieval_count"} {
		m := c.Metric(name)
		if m.Delta == nil || *m.Delta <= 0 || m.RelativeDelta != nil {
			t.Fatalf("%s = %+v", name, m)
		}
	}
	out, err := json.Marshal(c)
	if err != nil || bytes.Contains(out, []byte("Inf")) {
		t.Fatalf("json = %s, %v", out, err)
	}
}

func TestSuccessAndDifferentTrialCountsStayStructured(t *testing.T) {
	b, k := repeated("disabled", 5, 400, 31), repeated("stateful", 10, 90, 32)
	k.Repeated.SuccessCount, k.Repeated.TaskFailureCount = 8, 2
	c := mustCompare(t, b, k)
	success := c.Outcome("task_success")
	if *success.Baseline != (Fraction{5, 5}) || *success.Candidate != (Fraction{8, 10}) {
		t.Fatalf("success = %+v", success)
	}
	if c.Baseline.TrialCount != 5 || c.Candidate.TrialCount != 10 {
		t.Fatalf("sides = %+v %+v", c.Baseline, c.Candidate)
	}
	// Medians remain comparable; per-group totals are not subtracted.
	if c.Metric("median_command_count").Delta == nil {
		t.Fatal("median delta missing")
	}
	if m := c.Metric("total_show_count"); m.Delta != nil || m.Candidate == nil {
		t.Fatalf("total_show_count = %+v", m)
	}
}

func TestRecoveryComparison(t *testing.T) {
	b, k := repeated("stateful", 5, 100, 10), repeated("integrated", 5, 100, 10)
	b.Repeated.Aggregate.TrialsWithShow, b.Repeated.Aggregate.TotalShowCount = 1, 1
	k.Repeated.Aggregate = workload.RepeatedBenchmarkAggregate{
		MedianTotalVisibleBytes: ptr(100), MedianCommandCount: ptr(10), MedianWallTimeNS: ptr(1), MedianProcessingNS: ptr(1),
		TrialsWithShow: 3, TrialsWithRawRetrieval: 1, TotalShowCount: 4, TotalRawRetrievalCount: 2,
	}
	c := mustCompare(t, b, k)
	if o := c.Outcome("trials_with_show"); *o.Baseline != (Fraction{1, 5}) || *o.Candidate != (Fraction{3, 5}) {
		t.Fatalf("trials_with_show = %+v", o)
	}
	if o := c.Outcome("trials_with_raw_retrieval"); *o.Baseline != (Fraction{0, 5}) || *o.Candidate != (Fraction{1, 5}) {
		t.Fatalf("trials_with_raw_retrieval = %+v", o)
	}
	if m := c.Metric("total_show_count"); *m.Delta != 3 {
		t.Fatalf("total_show_count = %+v", m)
	}
	if m := c.Metric("total_raw_retrieval_count"); *m.Delta != 2 || m.RelativeDelta != nil {
		t.Fatalf("total_raw_retrieval_count = %+v", m)
	}
}

func TestMissingAndNotApplicableAreNotZero(t *testing.T) {
	b, k := repeated("disabled", 0, 0, 0), repeated("stateful", 5, 100, 10)
	b.Repeated.RequestedTrialCount, b.Repeated.SuccessCount, b.Repeated.RunStatus = 5, 0, "incomplete"
	b.Repeated.Aggregate = workload.RepeatedBenchmarkAggregate{}
	c := mustCompare(t, b, k)
	visible := c.Metric("median_total_visible_bytes")
	if visible.Baseline != nil || visible.BaselineStatus != Unavailable || visible.Delta != nil {
		t.Fatalf("visible = %+v", visible)
	}
	processing := c.Metric("median_processing_ns")
	if processing.Baseline != nil || processing.BaselineStatus != NotApplicable || processing.CandidateStatus != Measured {
		t.Fatalf("processing = %+v", processing)
	}
	if o := c.Outcome("task_success"); o.Baseline != nil || o.BaselineStatus != Unavailable {
		t.Fatalf("success = %+v", o)
	}
	var out bytes.Buffer
	if err := WriteHuman(&out, c); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "n/a") || !strings.Contains(out.String(), "0 of 5") {
		t.Fatalf("human output:\n%s", out.String())
	}
}

func TestDisabledModeRecoveryIsNotApplicable(t *testing.T) {
	c := mustCompare(t, repeated("disabled", 5, 400, 31), repeated("integrated", 5, 90, 32))
	for _, name := range []string{"trials_with_show", "trials_with_raw_retrieval"} {
		if o := c.Outcome(name); o.Baseline != nil || o.BaselineStatus != NotApplicable || o.Candidate == nil {
			t.Fatalf("%s = %+v", name, o)
		}
	}
	if m := c.Metric("total_show_count"); m.BaselineStatus != NotApplicable || m.Delta != nil {
		t.Fatalf("total_show_count = %+v", m)
	}
}

func TestIncompatibleResultsAreRejected(t *testing.T) {
	otherWorkload := repeated("stateful", 5, 1, 1)
	otherWorkload.Repeated.Workload = "bugfix-002"
	otherAgent := repeated("stateful", 5, 1, 1)
	otherAgent.Repeated.Agent = "claude"
	cases := map[string][2]*workload.StoredResult{
		"workload":      {repeated("disabled", 5, 1, 1), otherWorkload},
		"agent":         {repeated("disabled", 5, 1, 1), otherAgent},
		"kind":          {single("codex", "disabled", 1), repeated("stateful", 5, 1, 1)},
		"deterministic": {single("", "", 1), single("codex", "stateful", 1)},
	}
	for name, pair := range cases {
		_, err := Compare(pair[0], pair[1])
		var incompatible *IncompatibleError
		if !errors.As(err, &incompatible) {
			t.Errorf("%s: err = %v", name, err)
		}
	}
}

func TestDifferentModesAreCompatible(t *testing.T) {
	c := mustCompare(t, repeated("disabled", 5, 1, 1), repeated("stateful", 5, 1, 1))
	if c.Baseline.Mode != "disabled" || c.Candidate.Mode != "stateful" {
		t.Fatalf("modes = %q %q", c.Baseline.Mode, c.Candidate.Mode)
	}
}

func TestSingleResultsCompareDirectMetrics(t *testing.T) {
	c := mustCompare(t, single("codex", "disabled", 1000), single("codex", "integrated", 250))
	if c.Kind != workload.ResultKindSingle || c.Metric("median_total_visible_bytes") != nil {
		t.Fatalf("kind = %s", c.Kind)
	}
	if m := c.Metric("total_visible_bytes"); *m.Delta != -750 {
		t.Fatalf("visible = %+v", m)
	}
	if o := c.Outcome("task_success"); *o.Baseline != (Fraction{1, 1}) {
		t.Fatalf("success = %+v", o)
	}
	if m := c.Metric("processing_ns"); m.BaselineStatus != NotApplicable || m.CandidateStatus != Measured {
		t.Fatalf("processing = %+v", m)
	}

	deterministic := mustCompare(t, single("", "", 1000), single("", "", 900))
	if m := deterministic.Metric("wall_time_ns"); m.BaselineStatus != Unavailable || m.Delta != nil {
		t.Fatalf("wall = %+v", m)
	}
	if o := deterministic.Outcome("task_success"); o.BaselineStatus != Unavailable {
		t.Fatalf("success = %+v", o)
	}
}

func TestCompareDoesNotModifyInputsAndIsDeterministic(t *testing.T) {
	b, k := repeated("disabled", 5, 412000, 31), repeated("stateful", 5, 91000, 32)
	before := [2]workload.RepeatedBenchmarkResult{*b.Repeated, *k.Repeated}
	first, _ := json.Marshal(mustCompare(t, b, k))
	second, _ := json.Marshal(mustCompare(t, b, k))
	if !bytes.Equal(first, second) {
		t.Fatal("comparison is not deterministic")
	}
	if !reflect.DeepEqual(before, [2]workload.RepeatedBenchmarkResult{*b.Repeated, *k.Repeated}) {
		t.Fatal("inputs were modified")
	}
}

func TestHumanOutputShowsIndependentDimensions(t *testing.T) {
	b, k := repeated("disabled", 5, 412000, 31), repeated("stateful", 5, 91000, 32)
	k.Repeated.Aggregate.MedianProcessingNS = ptr(12_800_000)
	k.Repeated.Aggregate.TrialsWithShow = 3
	var out bytes.Buffer
	if err := WriteHuman(&out, mustCompare(t, b, k)); err != nil {
		t.Fatal(err)
	}
	s := out.String()
	for _, want := range []string{"Task success", "5/5", "Median visible bytes", "Median commands", "+1", "Median wall time",
		"Show used", "3/5", "Raw fallback", "AgentCap processing", "12.8ms", "-313.5KB (-77.9%)", "disabled", "stateful"} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q in:\n%s", want, s)
		}
	}
	for _, banned := range []string{"winner", "Winner", "score", "PASS", "FAIL", "better", "REGRESSION", "IMPROVEMENT"} {
		if strings.Contains(s, banned) {
			t.Errorf("human output contains %q", banned)
		}
	}
}

func TestJSONOutput(t *testing.T) {
	var out bytes.Buffer
	if err := WriteJSON(&out, mustCompare(t, repeated("disabled", 5, 100, 20), repeated("stateful", 10, 60, 25))); err != nil {
		t.Fatal(err)
	}
	var decoded struct {
		SchemaVersion int `json:"schema_version"`
		Baseline      struct {
			Mode       string `json:"mode"`
			TrialCount int    `json:"trial_count"`
		} `json:"baseline"`
		Metrics []struct {
			Name          string   `json:"name"`
			Delta         *int64   `json:"delta"`
			RelativeDelta *float64 `json:"relative_delta"`
		} `json:"metrics"`
	}
	if err := json.Unmarshal(out.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if decoded.SchemaVersion != ComparisonSchemaVersion || decoded.Baseline.Mode != "disabled" || decoded.Baseline.TrialCount != 5 {
		t.Fatalf("decoded = %+v", decoded)
	}
	if decoded.Metrics[0].Name != "median_total_visible_bytes" || *decoded.Metrics[0].Delta != -40 {
		t.Fatalf("metrics = %+v", decoded.Metrics)
	}
}

// The analysis path must never depend on execution: no agent adapters,
// command execution, integrations, or AgentCap processing packages.
func TestComparatorDoesNotImportExecution(t *testing.T) {
	allowed := map[string]bool{
		"github.com/taqu/agentcap/internal/workload": true,
		"github.com/taqu/agentcap/internal/stats":    true,
	}
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		src, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		parsed, err := parser.ParseFile(token.NewFileSet(), file, src, parser.ImportsOnly)
		if err != nil {
			t.Fatal(err)
		}
		for _, spec := range parsed.Imports {
			path, _ := strconv.Unquote(spec.Path.Value)
			if strings.HasPrefix(path, "github.com/taqu/agentcap/") && !allowed[path] {
				t.Errorf("%s imports %s", file, path)
			}
			if path == "os/exec" {
				t.Errorf("%s imports os/exec", file)
			}
		}
	}
}
