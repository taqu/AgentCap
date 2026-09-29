package workload

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/taqu/agentcap/internal/store"
)

const helperArg = "-workload-helper"

func TestMain(m *testing.M) {
	if len(os.Args) > 1 && os.Args[1] == helperArg {
		os.Exit(runHelper(os.Args[2:]))
	}
	os.Exit(m.Run())
}

func runHelper(args []string) int {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "touch":
			i++
			f, err := os.OpenFile(args[i], os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
			if err != nil {
				return 98
			}
			_, _ = f.WriteString("once\n")
			_ = f.Close()
		case "out":
			i++
			fmt.Print(args[i])
		case "exit":
			i++
			n, _ := strconv.Atoi(args[i])
			return n
		}
	}
	return 0
}

func makeLayout(t *testing.T, body string, fixtureFiles map[string]string) string {
	t.Helper()
	root := t.TempDir()
	workloads := filepath.Join(root, "benchmarks", "workloads", "test")
	fixture := filepath.Join(root, "benchmarks", "fixtures", "sample")
	if err := os.MkdirAll(workloads, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(fixture, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, content := range fixtureFiles {
		path := filepath.Join(fixture, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	path := filepath.Join(workloads, "sample.yaml")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func validPrefix() string {
	return "version: 1\nname: test/sample\nfixture: ../../fixtures/sample\n"
}

func TestLoadValidWorkload(t *testing.T) {
	path := makeLayout(t, validPrefix()+`steps:
  - run:
      argv: ["go", "version"]
      cwd: .
      expect: {exit: 0}
  - copy: {from: state.txt, to: nested.txt}
  - write: {path: value.txt, content: value}
  - remove: {path: value.txt}
  - mkdir: {path: nested/dir}
  - show: {command: 1}
  - raw: {command: 1, stream: stderr}
`, map[string]string{"state.txt": "state"})
	d, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if d.Version != 1 || d.Name != "test/sample" || len(d.Steps) != 7 {
		t.Fatalf("definition = %+v", d)
	}
}

func TestLoadRejectsInvalidDefinitions(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{"version", "version: 2\nname: test/x\nfixture: ../../fixtures/sample\nsteps: [{run: {argv: [go]}}]\n", "unsupported workload schema version"},
		{"missing name", "version: 1\nfixture: ../../fixtures/sample\nsteps: [{run: {argv: [go]}}]\n", "name is required"},
		{"unknown step", validPrefix() + "steps: [{shell: echo}]\n", "field shell not found"},
		{"empty argv", validPrefix() + "steps: [{run: {argv: []}}]\n", "argv must not be empty"},
		{"ambiguous", validPrefix() + "steps: [{run: {argv: [go]}, mkdir: {path: x}}]\n", "exactly one"},
		{"future show", validPrefix() + "steps: [{show: {command: 1}}, {run: {argv: [go]}}]\n", "must reference a preceding run"},
		{"bad raw stream", validPrefix() + "steps: [{run: {argv: [go]}}, {raw: {command: 1, stream: both}}]\n", "stdout or stderr"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := makeLayout(t, tc.body, nil)
			_, err := Load(path)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want substring %q", err, tc.want)
			}
		})
	}
}

func TestLoadPathSafety(t *testing.T) {
	abs := filepath.Join(t.TempDir(), "escape.txt")
	tests := []struct {
		name string
		step string
	}{
		{"absolute destination", fmt.Sprintf("  - write: {path: %q, content: x}\n", abs)},
		{"destination traversal", "  - remove: {path: ../outside}\n"},
		{"source traversal", "  - copy: {from: ../outside, to: x}\n"},
		{"absolute source", fmt.Sprintf("  - copy: {from: %q, to: x}\n", abs)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := makeLayout(t, validPrefix()+"steps:\n"+tc.step, nil)
			if _, err := Load(path); err == nil {
				t.Fatal("unsafe path accepted")
			}
		})
	}
	t.Run("valid nested", func(t *testing.T) {
		path := makeLayout(t, validPrefix()+"steps:\n  - copy: {from: states/a.txt, to: nested/a.txt}\n", map[string]string{"states/a.txt": "a"})
		if _, err := Load(path); err != nil {
			t.Fatal(err)
		}
	})
}

func TestRunExactlyOnceIsolationFreshSessionsCleanupAndRecovery(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	stable := strings.Repeat("stable output line\n", 200)
	body := validPrefix() + fmt.Sprintf(`steps:
  - run:
      argv: [%q, %q, "touch", "counter.txt", "out", %q]
      expect: {exit: 0}
  - run:
      argv: [%q, %q, "touch", "counter.txt", "out", %q]
      expect: {exit: 0}
`, exe, helperArg, stable, exe, helperArg, stable)
	path := makeLayout(t, body, map[string]string{"source.txt": "unchanged"})
	d, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	fixtureSource := filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(path))), "fixtures", "sample", "source.txt")
	tempParent := t.TempDir()
	storeRoot := t.TempDir()
	t.Setenv("LOCALAPPDATA", filepath.Join(storeRoot, "cache"))

	first, err := Run(context.Background(), d, Options{StoreRoot: storeRoot, TempRoot: tempParent, KeepWorkspace: true})
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(first.RetainedWorkspace, "counter.txt"))
	if err != nil || strings.Count(string(data), "\n") != 2 {
		t.Fatalf("counter = %q, err=%v", data, err)
	}
	if first.Aggregate.CommandCount != 2 || first.Aggregate.FullCount != 1 || first.Aggregate.UnchangedCount != 1 {
		t.Fatalf("aggregate = %+v", first.Aggregate)
	}
	canonical := NewBenchmarkResult(first)
	if canonical.ShowCount != 0 || canonical.RawRetrievalCount != 0 || canonical.ShowBytes != 0 || canonical.RawRetrievalBytes != 0 || canonical.TotalVisibleBytes != canonical.InitialVisibleBytes {
		t.Fatalf("no-recovery accounting = %+v", canonical)
	}
	if first.Aggregate.Commands[1].StatefulVisibleBytes >= first.Aggregate.Commands[1].StatelessVisibleBytes {
		t.Error("unchanged presentation did not reduce stateful bytes")
	}
	st, err := store.Open(storeRoot)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.Open(first.Aggregate.Commands[0].ResultID); err != nil {
		t.Errorf("stored result is not recoverable: %v", err)
	}
	_ = st.Close()
	if got, _ := os.ReadFile(fixtureSource); string(got) != "unchanged" {
		t.Fatalf("fixture mutated: %q", got)
	}

	second, err := Run(context.Background(), d, Options{StoreRoot: storeRoot, TempRoot: tempParent})
	if err != nil {
		t.Fatal(err)
	}
	if second.SessionID == first.SessionID {
		t.Fatal("workload runs reused a benchmark session")
	}
	if second.RetainedWorkspace != "" {
		t.Fatalf("workspace unexpectedly retained: %s", second.RetainedWorkspace)
	}
	entries, err := os.ReadDir(tempParent)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 { // only the explicitly retained first workspace
		t.Fatalf("temporary workspace was not cleaned: %v", entries)
	}
}

func TestRunExpectedExitAndAssertionFailure(t *testing.T) {
	exe, _ := os.Executable()
	for _, tc := range []struct {
		name     string
		expected int
		wantErr  bool
	}{
		{"expected failure", 7, false},
		{"unexpected status", 0, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("LOCALAPPDATA", filepath.Join(t.TempDir(), "cache"))
			body := validPrefix() + fmt.Sprintf("steps:\n  - run:\n      argv: [%q, %q, exit, \"7\"]\n      expect: {exit: %d}\n", exe, helperArg, tc.expected)
			path := makeLayout(t, body, nil)
			d, err := Load(path)
			if err != nil {
				t.Fatal(err)
			}
			tempParent := t.TempDir()
			r, err := Run(context.Background(), d, Options{StoreRoot: t.TempDir(), TempRoot: tempParent})
			if (err != nil) != tc.wantErr {
				t.Fatalf("error = %v", err)
			}
			if r.Aggregate == nil || r.Aggregate.CommandCount != 1 || r.Aggregate.Commands[0].ExitCode != 7 {
				t.Fatalf("non-zero command was not measured: %+v", r.Aggregate)
			}
			if tc.wantErr && (!strings.Contains(err.Error(), "step 1 (run)") || !strings.Contains(err.Error(), "expected 0")) {
				t.Fatalf("unclear assertion error: %v", err)
			}
			entries, readErr := os.ReadDir(tempParent)
			if readErr != nil || len(entries) != 0 {
				t.Fatalf("workspace not cleaned after run (err=%v, entries=%v)", readErr, entries)
			}
		})
	}
}

func TestRecoveryAccountingAndSingleExecution(t *testing.T) {
	exe, _ := os.Executable()
	outputA := strings.Repeat("alpha recovery line\n", 40)
	outputB := strings.Repeat("beta recovery line\n", 30)
	outputC := strings.Repeat("gamma recovery line\n", 20)
	tests := []struct {
		name     string
		steps    func(string) string
		runs     int
		shows    int
		raws     int
		rawBytes int64
	}{
		{
			name: "show",
			steps: func(counter string) string {
				return fmt.Sprintf(`
  - run: {argv: [%q, %q, touch, %q, out, %q]}
  - show: {command: 1}
`, exe, helperArg, counter, outputA)
			},
			runs: 1, shows: 1,
		},
		{
			name: "raw",
			steps: func(counter string) string {
				return fmt.Sprintf(`
  - run: {argv: [%q, %q, touch, %q, out, %q]}
  - raw: {command: 1}
`, exe, helperArg, counter, outputA)
			},
			runs: 1, raws: 1, rawBytes: int64(len(outputA)),
		},
		{
			name: "repeated",
			steps: func(counter string) string {
				return fmt.Sprintf(`
  - run: {argv: [%q, %q, touch, %q, out, %q]}
  - show: {command: 1}
  - show: {command: 1}
  - raw: {command: 1, stream: stdout}
`, exe, helperArg, counter, outputA)
			},
			runs: 1, shows: 2, raws: 1, rawBytes: int64(len(outputA)),
		},
		{
			name: "mixed",
			steps: func(counter string) string {
				return fmt.Sprintf(`
  - run: {argv: [%q, %q, touch, %q, out, %q]}
  - show: {command: 1}
  - run: {argv: [%q, %q, touch, %q, out, %q]}
  - run: {argv: [%q, %q, touch, %q, out, %q]}
  - show: {command: 3}
  - raw: {command: 3}
`, exe, helperArg, counter, outputA, exe, helperArg, counter, outputB, exe, helperArg, counter, outputC)
			},
			runs: 3, shows: 2, raws: 1, rawBytes: int64(len(outputC)),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			t.Setenv("LOCALAPPDATA", filepath.Join(root, "cache"))
			counter := filepath.Join(root, "counter.txt")
			path := makeLayout(t, validPrefix()+"steps:"+tc.steps(counter), nil)
			d, err := Load(path)
			if err != nil {
				t.Fatal(err)
			}
			storeRoot := t.TempDir()
			run, err := Run(context.Background(), d, Options{StoreRoot: storeRoot})
			if err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(counter)
			if err != nil || strings.Count(string(data), "\n") != tc.runs {
				t.Fatalf("commands executed %d times, want %d (err=%v)", strings.Count(string(data), "\n"), tc.runs, err)
			}
			result := NewBenchmarkResult(run)
			if result.Commands != tc.runs || result.ShowCount != tc.shows || result.RawRetrievalCount != tc.raws {
				t.Fatalf("counts = commands:%d show:%d raw:%d", result.Commands, result.ShowCount, result.RawRetrievalCount)
			}
			if tc.shows > 0 && result.ShowBytes <= 0 {
				t.Fatal("show retrieval produced no visible bytes")
			}
			if result.RawRetrievalBytes != tc.rawBytes {
				t.Fatalf("raw retrieval bytes = %d, want %d", result.RawRetrievalBytes, tc.rawBytes)
			}
			if result.TotalVisibleBytes != result.InitialVisibleBytes+result.ShowBytes+result.RawRetrievalBytes {
				t.Fatalf("total invariant failed: %+v", result)
			}
			if tc.name == "repeated" && run.Steps[1].RecoveryBytes*2 != result.ShowBytes {
				t.Fatalf("repeated show was deduplicated: first=%d total=%d", run.Steps[1].RecoveryBytes, result.ShowBytes)
			}
			st, err := store.Open(storeRoot)
			if err != nil {
				t.Fatal(err)
			}
			stats, err := st.LoadStats()
			_ = st.Close()
			if err != nil || stats.ShowCalls != int64(tc.shows) || stats.RawCalls != int64(tc.raws) {
				t.Fatalf("normal retrieval stats = %+v, err=%v", stats, err)
			}
		})
	}
}
