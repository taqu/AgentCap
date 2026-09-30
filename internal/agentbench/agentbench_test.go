package agentbench

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/taqu/agentcap/internal/store"
	"github.com/taqu/agentcap/internal/workload"
)

func TestMain(m *testing.M) {
	if os.Getenv("AGENTBENCH_PROCESS") == "1" {
		os.Exit(runAdapterProcess())
	}
	if len(os.Args) > 1 && os.Args[1] == "-verify-file" {
		data, err := os.ReadFile(os.Args[2])
		fmt.Print(strings.Repeat("verification output is not agent-visible\n", 100))
		if err == nil && string(data) == os.Args[3] {
			os.Exit(0)
		}
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func runAdapterProcess() int {
	if os.Getenv("AGENTBENCH_SLEEP") == "1" {
		time.Sleep(5 * time.Second)
	}
	fmt.Println(`{"type":"item.completed","item":{"type":"command_execution","aggregated_output":"abc"}}`)
	fmt.Println(`{"type":"item.completed","item":{"type":"command_execution","output":"12345"}}`)
	cwd, _ := os.Getwd()
	fmt.Fprintf(os.Stderr, "%s|%s|cwd=%s|args=%s", os.Getenv("ACAP_ROOT"), os.Getenv("ACAP_BENCH_MODE"), cwd, strings.Join(os.Args[1:], " "))
	if os.Getenv("AGENTBENCH_EXIT") == "1" {
		return 7
	}
	return 0
}

type fakeAdapter struct {
	correct bool
	timeout bool
	calls   int
	initial []string
}

func (a *fakeAdapter) Name() string { return "fake" }

func (a *fakeAdapter) Run(_ context.Context, req AgentRunRequest) (*AgentRunResult, error) {
	a.calls++
	path := filepath.Join(req.Workspace, "solution.txt")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	a.initial = append(a.initial, string(data))
	if a.correct {
		if err := os.WriteFile(path, []byte("fixed"), 0o644); err != nil {
			return nil, err
		}
	}
	if req.Mode == ModeDisabled {
		return &AgentRunResult{CommandCount: 2, CommandVisibleBytes: 700, TimedOut: a.timeout}, nil
	}
	st, err := store.Open(req.StoreRoot)
	if err != nil {
		return nil, err
	}
	defer st.Close()
	stateful := 200
	if req.Mode != ModeStateless {
		stateful = 120
	}
	_ = st.RecordRun(1000, 200, stateful, "full")
	_ = st.RecordProcessing(50)
	if req.Mode == ModeIntegrated {
		_ = st.RecordShow(30)
		_ = st.RecordRaw(40)
		_ = st.RecordProcessing(10)
	}
	return &AgentRunResult{ExitCode: 0, CommandCount: 1, CommandVisibleBytes: int64(stateful), TimedOut: a.timeout}, nil
}

func makeAgentWorkload(t *testing.T) *workload.Definition {
	t.Helper()
	root := t.TempDir()
	workloads := filepath.Join(root, "benchmarks", "workloads", "agent")
	fixture := filepath.Join(root, "benchmarks", "fixtures", "agent-sample")
	if err := os.MkdirAll(workloads, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(fixture, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(fixture, "solution.txt"), []byte("broken"), 0o644); err != nil {
		t.Fatal(err)
	}
	exe, _ := os.Executable()
	body := fmt.Sprintf(`version: 1
name: agent/sample
fixture: ../../fixtures/agent-sample
task: Fix solution.txt.
timeout: 2s
verify:
  - run:
      argv: [%q, "-verify-file", "solution.txt", "fixed"]
`, exe)
	path := filepath.Join(workloads, "sample.yaml")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	d, err := workload.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func TestRunnerModesFreshWorkspaceVerificationAndAccounting(t *testing.T) {
	definition := makeAgentWorkload(t)
	for _, mode := range []Mode{ModeDisabled, ModeStateless, ModeStateful, ModeIntegrated} {
		t.Run(string(mode), func(t *testing.T) {
			adapter := &fakeAdapter{correct: true}
			tempRoot := t.TempDir()
			trial, err := Run(context.Background(), definition, Options{
				Adapter: adapter, Mode: mode, StoreRoot: t.TempDir(), TempRoot: tempRoot,
			})
			if err != nil {
				t.Fatal(err)
			}
			if adapter.calls != 1 || len(adapter.initial) != 1 || adapter.initial[0] != "broken" {
				t.Fatalf("adapter calls/initial = %d/%q", adapter.calls, adapter.initial)
			}
			if trial.Benchmark.TaskSuccess == nil || !*trial.Benchmark.TaskSuccess || trial.Benchmark.Agent != "fake" || trial.Benchmark.Mode != string(mode) {
				t.Fatalf("benchmark metadata = %+v", trial.Benchmark)
			}
			if mode == ModeDisabled {
				if trial.Benchmark.Commands != 2 || trial.Benchmark.RawBytes != 700 || trial.Benchmark.TotalVisibleBytes != 700 {
					t.Fatalf("disabled metrics = %+v", trial.Benchmark)
				}
			} else {
				wantInitial := int64(120)
				if mode == ModeStateless {
					wantInitial = 200
				}
				if trial.Benchmark.Commands != 1 || trial.Benchmark.RawBytes != 1000 || trial.Benchmark.InitialVisibleBytes != wantInitial {
					t.Fatalf("enabled metrics = %+v", trial.Benchmark)
				}
				if mode == ModeIntegrated && (trial.Benchmark.ShowCount != 1 || trial.Benchmark.RawRetrievalCount != 1 || trial.Benchmark.TotalVisibleBytes != 190) {
					t.Fatalf("integrated recovery = %+v", trial.Benchmark)
				}
			}
			entries, _ := os.ReadDir(tempRoot)
			if len(entries) != 0 {
				t.Fatalf("workspace not cleaned: %v", entries)
			}
		})
	}
}

func TestRunnerTaskFailureAndTimeoutAreValidOutcomes(t *testing.T) {
	definition := makeAgentWorkload(t)
	for _, tc := range []struct {
		name    string
		adapter *fakeAdapter
		status  string
	}{
		{"task failure", &fakeAdapter{}, "completed"},
		{"timeout", &fakeAdapter{timeout: true}, "timeout"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			trial, err := Run(context.Background(), definition, Options{Adapter: tc.adapter, Mode: ModeDisabled, StoreRoot: t.TempDir()})
			if err != nil {
				t.Fatal(err)
			}
			if trial.Benchmark.TaskSuccess == nil || *trial.Benchmark.TaskSuccess || trial.Benchmark.ExecutionStatus != tc.status {
				t.Fatalf("outcome = %+v", trial.Benchmark)
			}
		})
	}
}

func TestCodexAdapterInvocationModesExitAndTimeout(t *testing.T) {
	exe, _ := os.Executable()
	t.Setenv("AGENTBENCH_PROCESS", "1")
	adapter := &CodexAdapter{Executable: exe}
	root := t.TempDir()
	request := AgentRunRequest{Workspace: t.TempDir(), Task: "test task", Mode: ModeDisabled, StoreRoot: root}
	result, err := adapter.Run(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if result.CommandCount != 2 || result.CommandVisibleBytes != 8 || !strings.Contains(string(result.Stderr), root+"|disabled") || !strings.Contains(string(result.Stderr), "cwd="+request.Workspace) || !strings.Contains(string(result.Stderr), "exec --json --ephemeral") {
		t.Fatalf("adapter result = %+v stderr=%q", result, result.Stderr)
	}

	t.Setenv("AGENTBENCH_EXIT", "1")
	result, err = adapter.Run(context.Background(), request)
	if err != nil || result.ExitCode != 7 {
		t.Fatalf("nonzero result=%+v err=%v", result, err)
	}
	t.Setenv("AGENTBENCH_EXIT", "")
	t.Setenv("AGENTBENCH_SLEEP", "1")
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	result, err = adapter.Run(ctx, request)
	if err != nil || !result.TimedOut {
		t.Fatalf("timeout result=%+v err=%v", result, err)
	}
}

func TestCodexAdapterConfiguresStatefulAndIntegratedDifferently(t *testing.T) {
	exe, _ := os.Executable()
	t.Setenv("AGENTBENCH_PROCESS", "1")
	adapter := &CodexAdapter{Executable: exe}
	for _, tc := range []struct {
		mode             Mode
		wantInstructions bool
	}{
		{ModeStateful, false},
		{ModeIntegrated, true},
	} {
		workspace := t.TempDir()
		_, err := adapter.Run(context.Background(), AgentRunRequest{
			Workspace: workspace, Task: "task", Mode: tc.mode, StoreRoot: t.TempDir(), HookCommand: `"C:\\tools\\acap.exe" hook codex`,
		})
		if err != nil {
			t.Fatal(err)
		}
		config, err := os.ReadFile(filepath.Join(workspace, ".codex", "config.toml"))
		if err != nil || !strings.Contains(string(config), "pre_tool_use") {
			t.Fatalf("hook config missing: %q err=%v", config, err)
		}
		_, instructionErr := os.Stat(filepath.Join(workspace, "AGENTS.md"))
		if (instructionErr == nil) != tc.wantInstructions {
			t.Fatalf("mode %s instructions error = %v", tc.mode, instructionErr)
		}
	}
}

func TestParseMode(t *testing.T) {
	for _, value := range []string{"disabled", "stateless", "stateful", "integrated"} {
		if mode, err := ParseMode(value); err != nil || string(mode) != value {
			t.Errorf("ParseMode(%q) = %q, %v", value, mode, err)
		}
	}
	if _, err := ParseMode("other"); err == nil {
		t.Fatal("invalid mode accepted")
	}
}
