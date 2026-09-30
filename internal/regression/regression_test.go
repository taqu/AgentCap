package regression

import (
	"bytes"
	"encoding/json"
	"go/parser"
	"go/token"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/taqu/agentcap/internal/benchcompare"
	"github.com/taqu/agentcap/internal/workload"
)

func f64(v float64) *float64 { return &v }
func i64(v int64) *int64     { return &v }

func result(name string, stateful int64) *workload.StoredResult {
	return &workload.StoredResult{Kind: workload.ResultKindSingle, Single: &workload.BenchmarkResult{
		SchemaVersion: workload.BenchmarkResultSchemaVersion, Workload: name, Commands: 3,
		RawBytes: 1000, StatelessBytes: stateful, StatefulBytes: stateful, InitialVisibleBytes: stateful,
		TotalVisibleBytes: stateful, ProcessingNS: 5_000_000,
	}}
}

func policy(t *testing.T, rules map[string]Rule) *Policy {
	t.Helper()
	p, err := NewPolicy(rules)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func checkOne(t *testing.T, baseline, candidate int64, rule Rule) Check {
	t.Helper()
	c, err := benchcompare.Compare(result("w", baseline), result("w", candidate))
	if err != nil {
		t.Fatal(err)
	}
	checks := EvaluateComparison(c, policy(t, map[string]Rule{"stateful_bytes": rule}))
	if len(checks) != 1 {
		t.Fatalf("checks = %+v", checks)
	}
	return checks[0]
}

func TestRelativeThreshold(t *testing.T) {
	rule := Rule{MaxRelativeIncrease: f64(0.10)}
	cases := []struct {
		candidate int64
		want      Status
	}{
		{105, Pass},       // +5%
		{115, Regression}, // +15%
		{110, Pass},       // exactly +10%: equal to the limit passes
		{90, Pass},        // decrease
	}
	for _, tc := range cases {
		c := checkOne(t, 100, tc.candidate, rule)
		if c.Status != tc.want {
			t.Errorf("100 -> %d: status %s, want %s", tc.candidate, c.Status, tc.want)
		}
		// The check carries the B8 comparison values unchanged.
		if *c.Delta != tc.candidate-100 || *c.Baseline != 100 || *c.Candidate != tc.candidate {
			t.Errorf("100 -> %d: metric %+v", tc.candidate, c.NumericMetric)
		}
	}
}

func TestAbsoluteThreshold(t *testing.T) {
	if s := checkOne(t, 100, 150, Rule{MaxAbsoluteIncrease: i64(50)}).Status; s != Pass {
		t.Errorf("+50 with limit 50: %s", s)
	}
	if s := checkOne(t, 100, 151, Rule{MaxAbsoluteIncrease: i64(50)}).Status; s != Regression {
		t.Errorf("+51 with limit 50: %s", s)
	}
}

func TestAbsoluteLimitIsNoiseFloorForRelative(t *testing.T) {
	both := Rule{MaxRelativeIncrease: f64(0.10), MaxAbsoluteIncrease: i64(32)}
	// +100% but only +10 bytes: within the absolute limit.
	if s := checkOne(t, 10, 20, both).Status; s != Pass {
		t.Errorf("10 -> 20: %s", s)
	}
	// +5% but +5000 bytes: within the relative limit.
	if s := checkOne(t, 100000, 105000, both).Status; s != Pass {
		t.Errorf("100000 -> 105000: %s", s)
	}
	// Exceeds both limits.
	if s := checkOne(t, 100, 200, both).Status; s != Regression {
		t.Errorf("100 -> 200: %s", s)
	}
}

func TestZeroBaseline(t *testing.T) {
	c := checkOne(t, 0, 10, Rule{MaxRelativeIncrease: f64(0.10)})
	if c.Status != NotEvaluable || c.RelativeDelta != nil || c.Reason == "" {
		t.Fatalf("relative-only zero baseline: %+v", c)
	}
	if s := checkOne(t, 0, 10, Rule{MaxRelativeIncrease: f64(0.10), MaxAbsoluteIncrease: i64(5)}).Status; s != Regression {
		t.Errorf("absolute limit should decide: %s", s)
	}
	if s := checkOne(t, 0, 0, Rule{MaxRelativeIncrease: f64(0.10)}).Status; s != Pass {
		t.Errorf("no increase: %s", s)
	}
	out, _ := json.Marshal(c)
	if bytes.Contains(out, []byte("Inf")) {
		t.Fatalf("json = %s", out)
	}
}

func TestMissingMetricIsNotEvaluable(t *testing.T) {
	// Deterministic workload results do not record wall time.
	c, err := benchcompare.Compare(result("w", 1), result("w", 1))
	if err != nil {
		t.Fatal(err)
	}
	checks := EvaluateComparison(c, policy(t, map[string]Rule{
		"wall_time_ns":               {MaxRelativeIncrease: f64(1)},
		"median_total_visible_bytes": {MaxRelativeIncrease: f64(1)},
	}))
	for _, check := range checks {
		if check.Status != NotEvaluable || check.Reason == "" {
			t.Errorf("%s: %+v", check.Name, check)
		}
	}
	e := EvaluateSuite([]string{"w"}, policy(t, map[string]Rule{
		"stateful_bytes": {MaxRelativeIncrease: f64(1)},
		"wall_time_ns":   {MaxRelativeIncrease: f64(1)},
	}), map[string]*workload.StoredResult{"w": result("w", 1)}, map[string]*workload.StoredResult{"w": result("w", 1)})
	if e.Passed || e.Status != StatusFailed || e.Summary.NotEvaluable != 1 || e.Summary.Passed != 1 {
		t.Fatalf("evaluation = %+v", e.Summary)
	}
}

func TestPolicyValidation(t *testing.T) {
	cases := map[string]map[string]Rule{
		"unknown metric": {"statefull_bytes": {MaxRelativeIncrease: f64(0.1)}},
		"raw bytes":      {"raw_bytes": {MaxRelativeIncrease: f64(0.1)}},
		"no limits":      {"stateful_bytes": {}},
		"negative rel":   {"stateful_bytes": {MaxRelativeIncrease: f64(-0.1)}},
		"negative abs":   {"stateful_bytes": {MaxAbsoluteIncrease: i64(-1)}},
		"nan":            {"stateful_bytes": {MaxRelativeIncrease: f64(math.NaN())}},
		"infinite":       {"stateful_bytes": {MaxRelativeIncrease: f64(math.Inf(1))}},
		"empty policy":   {},
	}
	for name, rules := range cases {
		if _, err := NewPolicy(rules); err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
}

func TestSuiteFileValidation(t *testing.T) {
	workloadPath, err := filepath.Abs("../../benchmarks/workloads/git/repeated-diff.yaml")
	if err != nil {
		t.Fatal(err)
	}
	agentPath, _ := filepath.Abs("../../benchmarks/workloads/agent/go-bugfix.yaml")
	metric := "metrics:\n  stateful_bytes:\n    max_relative_increase: 0.1\n"
	cases := map[string]string{
		"typo metric":        "version: 1\nworkloads: [" + strconv.Quote(workloadPath) + "]\nmetrics:\n  statefull_bytes:\n    max_relative_increase: 0.1\n",
		"typo rule key":      "version: 1\nworkloads: [" + strconv.Quote(workloadPath) + "]\nmetrics:\n  stateful_bytes:\n    max_relative_increse: 0.1\n",
		"duplicate metric":   "version: 1\nworkloads: [" + strconv.Quote(workloadPath) + "]\n" + metric + "  stateful_bytes:\n    max_relative_increase: 0.2\n",
		"bad percentage":     "version: 1\nworkloads: [" + strconv.Quote(workloadPath) + "]\nmetrics:\n  stateful_bytes:\n    max_relative_increase: ten\n",
		"agent workload":     "version: 1\nworkloads: [" + strconv.Quote(agentPath) + "]\n" + metric,
		"duplicate workload": "version: 1\nworkloads: [" + strconv.Quote(workloadPath) + ", " + strconv.Quote(workloadPath) + "]\n" + metric,
		"no workloads":       "version: 1\n" + metric,
		"version":            "version: 2\nworkloads: [" + strconv.Quote(workloadPath) + "]\n" + metric,
		"unknown top key":    "version: 1\nthreshold: 0.1\nworkloads: [" + strconv.Quote(workloadPath) + "]\n" + metric,
	}
	for name, content := range cases {
		path := filepath.Join(t.TempDir(), "regression.yaml")
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := LoadSuite(path); err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
}

func TestCheckedInSuiteIsValidAndDeterministic(t *testing.T) {
	suite, err := LoadSuite("../../benchmarks/regression.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if len(suite.Workloads) == 0 || len(suite.Policy.Rules) == 0 {
		t.Fatalf("suite = %+v", suite)
	}
	for _, w := range suite.Workloads {
		if w.Definition.Task != "" {
			t.Errorf("%s is an agent workload", w.Name)
		}
	}
}

func TestMultipleChecksAndPerWorkloadVisibility(t *testing.T) {
	p := policy(t, map[string]Rule{
		"stateless_bytes": {MaxRelativeIncrease: f64(0.10)},
		"stateful_bytes":  {MaxRelativeIncrease: f64(0.10)},
		"processing_ns":   {MaxRelativeIncrease: f64(0.50)},
	})
	baseline := map[string]*workload.StoredResult{"a": result("a", 100_000), "b": result("b", 100_000)}
	candidate := map[string]*workload.StoredResult{"a": result("a", 50_000), "b": result("b", 180_000)}
	// Within workload b: stateless passes, stateful regresses, processing passes.
	candidate["b"].Single.StatelessBytes = 100_000
	e := EvaluateSuite([]string{"a", "b"}, p, baseline, candidate)
	if e.Passed || e.Status != StatusFailed {
		t.Fatalf("status = %s", e.Status)
	}
	// Workload a improving by 50 KB does not offset workload b's +80 KB.
	b := e.Workloads[1]
	statuses := []Status{b.Checks[0].Status, b.Checks[1].Status, b.Checks[2].Status}
	if statuses[0] != Pass || statuses[1] != Regression || statuses[2] != Pass {
		t.Fatalf("b statuses = %v", statuses)
	}
	if e.Summary.Regressions != 1 || e.Summary.FailedWorkloads != 1 || e.Summary.Passed != 5 {
		t.Fatalf("summary = %+v", e.Summary)
	}
}

func TestSuiteWorkloadOutcomes(t *testing.T) {
	p := policy(t, map[string]Rule{"stateful_bytes": {MaxRelativeIncrease: f64(0.10)}})
	same := map[string]*workload.StoredResult{"a": result("a", 1)}
	if e := EvaluateSuite([]string{"a"}, p, same, same); !e.Passed || e.Status != StatusPassed {
		t.Fatalf("identical results: %+v", e)
	}
	missing := EvaluateSuite([]string{"a", "b"}, p, same, same)
	if missing.Passed || missing.Workloads[1].Status != WorkloadMissing {
		t.Fatalf("missing workload: %+v", missing.Workloads)
	}
	agent := result("a", 1)
	agent.Single.Agent = "codex"
	if e := EvaluateSuite([]string{"a"}, p, same, map[string]*workload.StoredResult{"a": agent}); e.Passed || e.Workloads[0].Status != WorkloadIncompatible {
		t.Fatalf("incompatible: %+v", e.Workloads)
	}
	// Results outside the suite are ignored; a suite with nothing evaluated
	// never passes.
	if e := EvaluateSuite(nil, p, same, same); e.Passed {
		t.Fatal("empty evaluation passed")
	}
}

func TestEvaluationDoesNotModifyInputs(t *testing.T) {
	p := policy(t, map[string]Rule{"stateful_bytes": {MaxRelativeIncrease: f64(0.10)}})
	b, k := result("a", 100), result("a", 500)
	before := [2]workload.BenchmarkResult{*b.Single, *k.Single}
	first, _ := json.Marshal(EvaluateSuite([]string{"a"}, p, map[string]*workload.StoredResult{"a": b}, map[string]*workload.StoredResult{"a": k}))
	second, _ := json.Marshal(EvaluateSuite([]string{"a"}, p, map[string]*workload.StoredResult{"a": b}, map[string]*workload.StoredResult{"a": k}))
	if !bytes.Equal(first, second) {
		t.Fatal("evaluation is not deterministic")
	}
	if before != [2]workload.BenchmarkResult{*b.Single, *k.Single} {
		t.Fatal("inputs modified")
	}
}

func TestHumanOutput(t *testing.T) {
	p := policy(t, map[string]Rule{"stateful_bytes": {MaxRelativeIncrease: f64(0.10), MaxAbsoluteIncrease: i64(1000)}})
	b, k := result("git/diff", 20_000), result("git/diff", 23_000)
	k.Single.RawBytes = 1500
	e := EvaluateSuite([]string{"git/diff", "gone"}, p, map[string]*workload.StoredResult{"git/diff": b}, map[string]*workload.StoredResult{"git/diff": k})
	var out bytes.Buffer
	if err := WriteHuman(&out, e); err != nil {
		t.Fatal(err)
	}
	s := out.String()
	for _, want := range []string{"git/diff", "stateful_bytes", "19.5KB", "22.5KB", "+2.9KB (+15.0%)", "+10.0% or +1000B",
		"REGRESSION", "MISSING", "no baseline or candidate result", "raw_bytes (reference, not gated)", "Result: FAILED"} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q in:\n%s", want, s)
		}
	}
	out.Reset()
	if err := WriteHuman(&out, ErrorEvaluation(os.ErrNotExist)); err != nil || !strings.Contains(out.String(), "Result: ERROR") {
		t.Fatalf("error output: %q %v", out.String(), err)
	}
}

func TestJSONOutput(t *testing.T) {
	p := policy(t, map[string]Rule{"stateful_bytes": {MaxRelativeIncrease: f64(0.10)}})
	e := EvaluateSuite([]string{"a"}, p, map[string]*workload.StoredResult{"a": result("a", 20000)}, map[string]*workload.StoredResult{"a": result("a", 23000)})
	var out bytes.Buffer
	if err := WriteJSON(&out, e); err != nil {
		t.Fatal(err)
	}
	var decoded struct {
		SchemaVersion int    `json:"schema_version"`
		Status        string `json:"status"`
		Passed        bool   `json:"passed"`
		Workloads     []struct {
			Workload string `json:"workload"`
			Checks   []struct {
				Name                string   `json:"name"`
				Baseline            int64    `json:"baseline"`
				Candidate           int64    `json:"candidate"`
				Delta               int64    `json:"delta"`
				RelativeDelta       float64  `json:"relative_delta"`
				MaxRelativeIncrease *float64 `json:"max_relative_increase"`
				MaxAbsoluteIncrease *int64   `json:"max_absolute_increase"`
				Status              string   `json:"status"`
			} `json:"checks"`
		} `json:"workloads"`
	}
	if err := json.Unmarshal(out.Bytes(), &decoded); err != nil {
		t.Fatal(err)
	}
	c := decoded.Workloads[0].Checks[0]
	if decoded.SchemaVersion != EvaluationSchemaVersion || decoded.Status != "failed" || decoded.Passed ||
		c.Name != "stateful_bytes" || c.Baseline != 20000 || c.Candidate != 23000 || c.Delta != 3000 ||
		c.RelativeDelta != 0.15 || *c.MaxRelativeIncrease != 0.10 || c.MaxAbsoluteIncrease != nil || c.Status != "regression" {
		t.Fatalf("decoded = %+v", decoded)
	}
	out.Reset()
	WriteJSON(&out, ErrorEvaluation(os.ErrNotExist))
	if !strings.Contains(out.String(), `"status": "error"`) || !strings.Contains(out.String(), `"passed": false`) {
		t.Fatalf("error json = %s", out.String())
	}
}

// Every policy metric must be a metric the B8 comparator actually reports.
func TestPolicyMetricsExistInComparison(t *testing.T) {
	single, _ := benchcompare.Compare(result("w", 1), result("w", 1))
	visible := int64(1)
	rep := &workload.StoredResult{Kind: workload.ResultKindRepeated, Repeated: &workload.RepeatedBenchmarkResult{
		Workload: "w", Agent: "codex", Mode: "stateful", TrialCount: 1,
		Aggregate: workload.RepeatedBenchmarkAggregate{MedianTotalVisibleBytes: &visible},
	}}
	repeated, _ := benchcompare.Compare(rep, rep)
	for _, name := range Metrics {
		if single.Metric(name) == nil && repeated.Metric(name) == nil {
			t.Errorf("policy metric %s is not produced by the comparator", name)
		}
	}
}

// The policy layer consumes comparisons only: it must not depend on agent
// adapters, command execution, integrations, or workload execution.
func TestEvaluatorDoesNotDependOnExecution(t *testing.T) {
	allowed := map[string]bool{
		"github.com/taqu/agentcap/internal/workload":     true,
		"github.com/taqu/agentcap/internal/benchcompare": true,
	}
	files, _ := filepath.Glob("*.go")
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
			if (strings.HasPrefix(path, "github.com/taqu/agentcap/") && !allowed[path]) || path == "os/exec" {
				t.Errorf("%s imports %s", file, path)
			}
		}
		if bytes.Contains(src, []byte("workload.Run(")) {
			t.Errorf("%s executes workloads", file)
		}
	}
}
