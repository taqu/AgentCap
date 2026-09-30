package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// The test binary doubles as the acap CLI (ACAP_TEST_AS_CLI=1) and as a
// deterministic target command (first arg helperArg). The helper check runs
// first so a target launched from the CLI role does not re-enter main().
const helperArg = "-cli-helper"

func TestMain(m *testing.M) {
	if len(os.Args) > 1 && os.Args[1] == helperArg {
		os.Exit(runHelper(os.Args[2:]))
	}
	if os.Getenv("ACAP_TEST_AS_CLI") == "1" {
		os.Args = append([]string{"acap"}, os.Args[1:]...)
		main()
		os.Exit(0)
	}
	os.Exit(m.Run())
}

// runHelper: "touch <path>" appends a line, "out <s>" / "err <s>" write,
// "exit <n>" exits with n.
func runHelper(args []string) int {
	for i := 0; i+1 < len(args); i += 2 {
		switch args[i] {
		case "touch":
			f, err := os.OpenFile(args[i+1], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
			if err != nil {
				return 99
			}
			f.WriteString("x\n")
			f.Close()
		case "out":
			os.Stdout.WriteString(args[i+1])
		case "err":
			os.Stderr.WriteString(args[i+1])
		case "exit":
			n, _ := strconv.Atoi(args[i+1])
			return n
		}
	}
	return 0
}

// acap runs the CLI with args in an isolated store root and returns stdout
// and the exit code.
func acap(t *testing.T, root string, args ...string) (string, int) {
	t.Helper()
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(exe, args...)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "ACAP_TEST_AS_CLI=1", "ACAP_ROOT="+root, "ACAP_SESSION_ID=", "LOCALAPPDATA="+root, "XDG_CACHE_HOME="+root)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	code := 0
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		code = ee.ExitCode()
	} else if err != nil {
		t.Fatal(err)
	}
	if stderr.Len() > 0 {
		t.Logf("stderr: %s", stderr.String())
	}
	return string(out), code
}

func TestBenchCommandExecutesOnce(t *testing.T) {
	root := t.TempDir()
	counter := filepath.Join(root, "counter.txt")
	exe, _ := os.Executable()

	out, code := acap(t, root, "bench", "command", "--",
		exe, helperArg, "touch", counter, "out", "hello\n")
	if code != 0 {
		t.Fatalf("exit = %d, output:\n%s", code, out)
	}
	data, err := os.ReadFile(counter)
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(string(data), "\n"); n != 1 {
		t.Fatalf("target executed %d times, want exactly 1", n)
	}
	if !strings.Contains(out, "Benchmark: ") || !strings.Contains(out, "total:      6 B") {
		t.Errorf("unexpected report:\n%s", out)
	}
}

func TestBenchCommandFailedExitAndRecovery(t *testing.T) {
	root := t.TempDir()
	exe, _ := os.Executable()

	out, code := acap(t, root, "bench", "command", "--json", "--",
		exe, helperArg, "out", "FAIL\n", "err", "boom\n", "exit", "4")
	if code != 4 {
		t.Fatalf("exit = %d, want target exit 4", code)
	}
	var rep struct {
		ExitCode          int      `json:"exit_code"`
		RawStdoutBytes    int64    `json:"raw_stdout_bytes"`
		RawStderrBytes    int64    `json:"raw_stderr_bytes"`
		RawBytes          int64    `json:"raw_bytes"`
		AgentVisibleBytes int64    `json:"agent_visible_bytes"`
		ReductionRatio    *float64 `json:"reduction_ratio"`
		ResultID          string   `json:"result_id"`
	}
	if err := json.Unmarshal([]byte(out), &rep); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	if rep.ExitCode != 4 || rep.RawStdoutBytes != 5 || rep.RawStderrBytes != 5 || rep.RawBytes != 10 {
		t.Errorf("report = %+v", rep)
	}
	if rep.ResultID == "" || rep.ReductionRatio == nil {
		t.Fatalf("report = %+v", rep)
	}

	raw, code := acap(t, root, "raw", rep.ResultID)
	if code != 0 || raw != "FAIL\n" {
		t.Errorf("acap raw = %q (exit %d)", raw, code)
	}
	rawErr, code := acap(t, root, "raw", "--stderr", rep.ResultID)
	if code != 0 || rawErr != "boom\n" {
		t.Errorf("acap raw --stderr = %q (exit %d)", rawErr, code)
	}
	shown, code := acap(t, root, "show", rep.ResultID)
	if code != 0 || !strings.Contains(shown, "FAIL") || !strings.Contains(shown, "boom") {
		t.Errorf("acap show = %q (exit %d)", shown, code)
	}
}

func TestBenchCommandNotFound(t *testing.T) {
	root := t.TempDir()
	_, code := acap(t, root, "bench", "command", "--", "acap-definitely-not-a-command-xyz")
	if code != 127 {
		t.Errorf("exit = %d, want 127", code)
	}
}

func TestBenchHelp(t *testing.T) {
	root := t.TempDir()
	for _, args := range [][]string{{"bench", "--help"}, {"bench", "command", "--help"}, {"bench", "session", "--help"}, {"bench", "run", "--help"}, {"bench", "agent", "--help"}} {
		out, code := acap(t, root, args...)
		if code != 0 || (!strings.Contains(out, "single command") && !strings.Contains(out, "benchmark session") && !strings.Contains(out, "versioned workload") && !strings.Contains(out, "coding-agent trial")) {
			t.Errorf("%v: exit %d, output:\n%s", args, code, out)
		}
	}
}

func TestBenchRunWorkloadCLI(t *testing.T) {
	root := t.TempDir()
	workloads := filepath.Join(root, "benchmarks", "workloads", "cli")
	fixture := filepath.Join(root, "benchmarks", "fixtures", "cli")
	if err := os.MkdirAll(workloads, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(fixture, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(fixture, "marker.txt"), []byte("original"), 0o644); err != nil {
		t.Fatal(err)
	}
	exe, _ := os.Executable()
	definition := fmt.Sprintf(`version: 1
name: cli/smoke
fixture: ../../fixtures/cli
steps:
  - write: {path: marker.txt, content: changed}
  - run:
      argv: [%q, %q, out, %q]
      expect: {exit: 0}
  - run:
      argv: [%q, %q, out, %q]
      expect: {exit: 0}
`, exe, helperArg, strings.Repeat("same line\n", 100), exe, helperArg, strings.Repeat("same line\n", 100))
	path := filepath.Join(workloads, "smoke.yaml")
	if err := os.WriteFile(path, []byte(definition), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := acap(t, root, "bench", "run", "--verbose", path)
	if code != 0 {
		t.Fatalf("exit=%d output:\n%s", code, out)
	}
	for _, want := range []string{"Benchmark Workload: cli/smoke", "commands:    2", "mutations:   1", "unchanged:             1", "id="} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in output:\n%s", want, out)
		}
	}
	if got, _ := os.ReadFile(filepath.Join(fixture, "marker.txt")); string(got) != "original" {
		t.Fatalf("source fixture changed: %q", got)
	}
}

func TestBenchRunJSONIsStableAndExecutesOnce(t *testing.T) {
	root := t.TempDir()
	workloads := filepath.Join(root, "benchmarks", "workloads", "cli")
	fixture := filepath.Join(root, "benchmarks", "fixtures", "cli-json")
	if err := os.MkdirAll(workloads, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(fixture, 0o755); err != nil {
		t.Fatal(err)
	}
	counter := filepath.Join(root, "json-counter.txt")
	exe, _ := os.Executable()
	definition := fmt.Sprintf(`version: 1
name: cli/json
fixture: ../../fixtures/cli-json
steps:
  - run:
      argv: [%q, %q, touch, %q, out, %q]
      expect: {exit: 0}
  - show: {command: 1}
  - show: {command: 1}
  - raw: {command: 1}
`, exe, helperArg, counter, strings.Repeat("JSON output line\n", 100))
	path := filepath.Join(workloads, "json.yaml")
	if err := os.WriteFile(path, []byte(definition), 0o644); err != nil {
		t.Fatal(err)
	}
	// Deliberately put --json after the workload, as documented.
	out, code := acap(t, root, "bench", "run", path, "--json")
	if code != 0 {
		t.Fatalf("exit=%d output=%q", code, out)
	}
	var result struct {
		SchemaVersion       int    `json:"schema_version"`
		Workload            string `json:"workload"`
		Commands            int    `json:"commands"`
		RawBytes            int64  `json:"raw_bytes"`
		StatelessBytes      int64  `json:"stateless_bytes"`
		StatefulBytes       int64  `json:"stateful_bytes"`
		InitialVisibleBytes int64  `json:"initial_visible_bytes"`
		ShowBytes           int64  `json:"show_bytes"`
		RawRetrievalBytes   int64  `json:"raw_retrieval_bytes"`
		TotalVisibleBytes   int64  `json:"total_visible_bytes"`
		ShowCount           int    `json:"show_count"`
		RawRetrievalCount   int    `json:"raw_retrieval_count"`
		ProcessingNS        int64  `json:"processing_ns"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("stdout is not JSON-only: %v\n%s", err, out)
	}
	if result.SchemaVersion != 4 || result.Workload != "cli/json" || result.Commands != 1 || result.RawBytes == 0 || result.ProcessingNS <= 0 || result.InitialVisibleBytes != result.StatefulBytes || result.ShowCount != 2 || result.RawRetrievalCount != 1 || result.ShowBytes <= 0 || result.RawRetrievalBytes <= 0 || result.TotalVisibleBytes != result.InitialVisibleBytes+result.ShowBytes+result.RawRetrievalBytes {
		t.Fatalf("result = %+v", result)
	}
	data, err := os.ReadFile(counter)
	if err != nil || strings.Count(string(data), "\n") != 1 {
		t.Fatalf("target executed other than once: counter=%q err=%v", data, err)
	}
}

func TestBenchmarkSessionIDModes(t *testing.T) {
	t.Setenv("ACAP_BENCH_MODE", "stateless")
	if got := benchmarkSessionID("agent-session"); got != "" {
		t.Fatalf("stateless session = %q", got)
	}
	t.Setenv("ACAP_BENCH_MODE", "stateful")
	if got := benchmarkSessionID("agent-session"); got == "" {
		t.Fatal("stateful mode did not map the agent session")
	}
	withoutScope := benchmarkSessionID("agent-session")
	t.Setenv("ACAP_BENCH_SESSION_SCOPE", "trial-one")
	first := benchmarkSessionID("agent-session")
	t.Setenv("ACAP_BENCH_SESSION_SCOPE", "trial-two")
	second := benchmarkSessionID("agent-session")
	if first == withoutScope || second == withoutScope || first == second {
		t.Fatalf("trial scopes not isolated: unscoped=%q first=%q second=%q", withoutScope, first, second)
	}
}

func TestBenchAgentRejectsInvalidRepeat(t *testing.T) {
	root := t.TempDir()
	for _, repeat := range []string{"0", "-1"} {
		_, code := acap(t, root, "bench", "agent", "--workload", "unused.yaml", "--agent", "codex", "--mode", "disabled", "--repeat", repeat)
		if code == 0 {
			t.Fatalf("repeat %s accepted", repeat)
		}
	}
}

func TestBenchSessionPersistsAcrossCLIInvocations(t *testing.T) {
	root := t.TempDir()
	started, code := acap(t, root, "bench", "session", "start")
	if code != 0 {
		t.Fatalf("start exit=%d output=%q", code, started)
	}
	id := strings.TrimSpace(started)
	if !strings.HasPrefix(id, "bench-") {
		t.Fatalf("session ID = %q", id)
	}
	exe, _ := os.Executable()
	output := strings.Repeat("cross process output line\n", 100)
	var second struct {
		Presentation          string `json:"presentation"`
		StatelessVisibleBytes int64  `json:"stateless_visible_bytes"`
		StatefulVisibleBytes  int64  `json:"stateful_visible_bytes"`
		ResultID              string `json:"result_id"`
	}
	for i := 0; i < 2; i++ {
		out, commandCode := acap(t, root, "bench", "command", "--json", "--session", id, "--",
			exe, helperArg, "out", output)
		if commandCode != 0 {
			t.Fatalf("command %d exit=%d output=%s", i+1, commandCode, out)
		}
		if i == 1 {
			if err := json.Unmarshal([]byte(out), &second); err != nil {
				t.Fatal(err)
			}
		}
	}
	if second.Presentation != "unchanged" || second.StatefulVisibleBytes >= second.StatelessVisibleBytes {
		t.Errorf("second measurement = %+v", second)
	}
	shown, code := acap(t, root, "bench", "session", "show", "--json", id)
	if code != 0 {
		t.Fatalf("show exit=%d output=%s", code, shown)
	}
	var summary struct {
		CommandCount   int `json:"command_count"`
		FullCount      int `json:"full_count"`
		UnchangedCount int `json:"unchanged_count"`
	}
	if err := json.Unmarshal([]byte(shown), &summary); err != nil {
		t.Fatal(err)
	}
	if summary.CommandCount != 2 || summary.FullCount != 1 || summary.UnchangedCount != 1 {
		t.Errorf("summary = %+v", summary)
	}
	raw, code := acap(t, root, "raw", second.ResultID)
	if code != 0 || raw != output {
		t.Errorf("raw recovery exit=%d len=%d", code, len(raw))
	}
}

func TestBenchSessionCommandExecutesOnce(t *testing.T) {
	root := t.TempDir()
	started, code := acap(t, root, "bench", "session", "start")
	if code != 0 {
		t.Fatal(started)
	}
	counter := filepath.Join(root, "session-cli-counter.txt")
	exe, _ := os.Executable()
	_, code = acap(t, root, "bench", "command", "--session", strings.TrimSpace(started), "--",
		exe, helperArg, "touch", counter, "out", "once\n")
	if code != 0 {
		t.Fatalf("command exit=%d", code)
	}
	data, err := os.ReadFile(counter)
	if err != nil || strings.Count(string(data), "\n") != 1 {
		t.Fatalf("counter=%q err=%v", data, err)
	}
}

func writeRepeatedResult(t *testing.T, path, workloadName, mode string, trials, success int, visible int64) {
	t.Helper()
	processing := `0`
	if mode != "disabled" {
		processing = `12800000`
	}
	data := fmt.Sprintf(`{"schema_version":4,"workload":%q,"agent":"codex","mode":%q,
		"requested_trial_count":%d,"trial_count":%d,"success_count":%d,"task_failure_count":%d,
		"timeout_count":0,"agent_error_count":0,"canceled_count":0,"run_status":"completed",
		"aggregate":{"median_total_visible_bytes":%d,"median_command_count":31,"median_wall_time_ns":84200000000,
		"median_processing_ns":%s,"trials_with_show":1,"trials_with_raw_retrieval":0,"total_show_count":1,"total_raw_retrieval_count":0},
		"trials":[]}`, workloadName, mode, trials, trials, success, trials-success, visible, processing)
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestBenchCompareReadsResultsWithoutExecution(t *testing.T) {
	root := t.TempDir()
	baseline := filepath.Join(root, "baseline.json")
	candidate := filepath.Join(root, "candidate.json")
	writeRepeatedResult(t, baseline, "bugfix-001", "disabled", 5, 5, 412000)
	writeRepeatedResult(t, candidate, "bugfix-001", "stateful", 10, 8, 91000)
	before := map[string][]byte{}
	for _, p := range []string{baseline, candidate} {
		before[p], _ = os.ReadFile(p)
	}

	out, code := acap(t, root, "bench", "compare", "--json", baseline, candidate)
	if code != 0 {
		t.Fatalf("exit %d: %s", code, out)
	}
	var decoded struct {
		SchemaVersion int `json:"schema_version"`
		Candidate     struct {
			TrialCount int `json:"trial_count"`
		} `json:"candidate"`
		Metrics []struct {
			Name  string `json:"name"`
			Delta *int64 `json:"delta"`
		} `json:"metrics"`
	}
	if err := json.Unmarshal([]byte(out), &decoded); err != nil {
		t.Fatalf("stdout is not clean JSON: %v\n%s", err, out)
	}
	if decoded.SchemaVersion != 2 || decoded.Candidate.TrialCount != 10 || decoded.Metrics[0].Delta == nil || *decoded.Metrics[0].Delta != -321000 {
		t.Fatalf("decoded = %+v", decoded)
	}

	out, code = acap(t, root, "bench", "compare", baseline, candidate)
	if code != 0 {
		t.Fatalf("exit %d: %s", code, out)
	}
	for _, want := range []string{"Benchmark comparison", "Task success", "5/5", "8/10", "Median visible bytes", "Show used", "AgentCap processing"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}

	// Comparison is read-only: inputs are untouched and no store is created.
	entries, _ := os.ReadDir(root)
	if len(entries) != 2 {
		t.Fatalf("compare created files: %v", entries)
	}
	for p, data := range before {
		if now, _ := os.ReadFile(p); string(now) != string(data) {
			t.Fatalf("%s was modified", p)
		}
	}
}

func TestBenchCompareRejectsIncompatibleAndInvalidInputs(t *testing.T) {
	root := t.TempDir()
	baseline := filepath.Join(root, "baseline.json")
	other := filepath.Join(root, "other.json")
	invalid := filepath.Join(root, "invalid.json")
	writeRepeatedResult(t, baseline, "bugfix-001", "disabled", 5, 5, 1)
	writeRepeatedResult(t, other, "bugfix-002", "stateful", 5, 5, 1)
	if err := os.WriteFile(invalid, []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{
		{baseline, other},
		{baseline, invalid},
		{baseline, filepath.Join(root, "missing.json")},
		{baseline},
	} {
		out, code := acap(t, root, append([]string{"bench", "compare", "--json"}, args...)...)
		if code == 0 || out != "" {
			t.Errorf("%v: exit %d, stdout %q", args, code, out)
		}
	}
}

// writeRegressionPolicy writes a one-workload regression suite gating only
// deterministic byte metrics, and returns its path and workload name.
func writeRegressionPolicy(t *testing.T, dir, metrics string) (string, string) {
	t.Helper()
	workloadPath, err := filepath.Abs("../../benchmarks/workloads/git/repeated-diff.yaml")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "regression.yaml")
	content := "version: 1\nworkloads:\n  - " + strconv.Quote(workloadPath) + "\nmetrics:\n" + metrics
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path, "git/repeated-diff"
}

func writeSingleResult(t *testing.T, dir, workloadName string, stateful int64) string {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "result.json")
	data := fmt.Sprintf(`{"schema_version":4,"workload":%q,"commands":3,"raw_bytes":716,"stateless_bytes":285,
		"stateful_bytes":%d,"initial_visible_bytes":%d,"show_bytes":0,"raw_retrieval_bytes":0,
		"total_visible_bytes":%d,"show_count":0,"raw_retrieval_count":0,"processing_ns":1000}`,
		workloadName, stateful, stateful, stateful)
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestBenchCheckExitStatus(t *testing.T) {
	root := t.TempDir()
	policy, name := writeRegressionPolicy(t, root, "  stateful_bytes:\n    max_relative_increase: 0.10\n")
	baseline := filepath.Join(root, "baseline")
	same := filepath.Join(root, "same")
	regressed := filepath.Join(root, "regressed")
	writeSingleResult(t, baseline, name, 200)
	writeSingleResult(t, same, name, 200)
	regressedFile := writeSingleResult(t, regressed, name, 300)

	out, code := acap(t, root, "bench", "check", "--baseline", baseline, "--candidate", same, "--policy", policy)
	if code != 0 || !strings.Contains(out, "Result: PASSED") {
		t.Fatalf("pass: exit %d\n%s", code, out)
	}

	out, code = acap(t, root, "bench", "check", "--baseline", baseline, "--candidate", regressed, "--policy", policy)
	if code != 1 || !strings.Contains(out, "REGRESSION") || !strings.Contains(out, "+50.0%") {
		t.Fatalf("regression: exit %d\n%s", code, out)
	}
	out, code = acap(t, root, "bench", "check", "--json", "--baseline", baseline, "--candidate", regressed, "--policy", policy)
	var evaluation struct {
		Status string `json:"status"`
		Passed bool   `json:"passed"`
	}
	if err := json.Unmarshal([]byte(out), &evaluation); err != nil || code != 1 || evaluation.Status != "failed" || evaluation.Passed {
		t.Fatalf("regression json: exit %d err %v\n%s", code, err, out)
	}

	// Infrastructure errors are non-zero and distinguishable from regressions.
	out, code = acap(t, root, "bench", "check", "--json", "--baseline", filepath.Join(root, "missing"), "--candidate", regressed, "--policy", policy)
	if err := json.Unmarshal([]byte(out), &evaluation); err != nil || code != 2 || evaluation.Status != "error" || evaluation.Passed {
		t.Fatalf("infrastructure error: exit %d err %v\n%s", code, err, out)
	}
	badPolicy, _ := writeRegressionPolicy(t, t.TempDir(), "  statefull_bytes:\n    max_relative_increase: 0.10\n")
	if _, code = acap(t, root, "bench", "check", "--baseline", baseline, "--candidate", same, "--policy", badPolicy); code != 2 {
		t.Fatalf("invalid policy: exit %d", code)
	}

	// B8 comparison stays descriptive for the same regressed pair.
	if out, code = acap(t, root, "bench", "compare", filepath.Join(baseline, "result.json"), regressedFile); code != 0 {
		t.Fatalf("compare: exit %d\n%s", code, out)
	}
}

// The deterministic path end to end: benchmark -> compare -> policy, with no
// coding agent involved.
func TestBenchSuiteAndCheckDeterministicWorkload(t *testing.T) {
	root := t.TempDir()
	policy, _ := writeRegressionPolicy(t, root,
		"  stateless_bytes:\n    max_relative_increase: 0\n  stateful_bytes:\n    max_relative_increase: 0\n")
	baseline := filepath.Join(root, "results", "baseline")
	candidate := filepath.Join(root, "results", "candidate")
	for _, dir := range []string{baseline, candidate} {
		if out, code := acap(t, root, "bench", "suite", "--out", dir, policy); code != 0 {
			t.Fatalf("suite: exit %d\n%s", code, out)
		}
	}
	if _, err := os.Stat(filepath.Join(baseline, "git_repeated-diff.json")); err != nil {
		t.Fatal(err)
	}
	out, code := acap(t, root, "bench", "check", "--baseline", baseline, "--candidate", candidate, "--policy", policy)
	if code != 0 || !strings.Contains(out, "Result: PASSED") {
		t.Fatalf("check: exit %d\n%s", code, out)
	}
	// Existing results are never overwritten.
	if _, code := acap(t, root, "bench", "suite", "--out", baseline, policy); code == 0 {
		t.Fatal("suite wrote into a non-empty directory")
	}
}
