package main

import (
	"encoding/json"
	"errors"
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
	for _, args := range [][]string{{"bench", "--help"}, {"bench", "command", "--help"}, {"bench", "session", "--help"}} {
		out, code := acap(t, root, args...)
		if code != 0 || (!strings.Contains(out, "single command") && !strings.Contains(out, "benchmark session")) {
			t.Errorf("%v: exit %d, output:\n%s", args, code, out)
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
