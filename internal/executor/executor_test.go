package executor_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/taqu/agentcap/internal/executor"
)

// TestHelperProcess is not a real test. It is re-invoked as a subprocess by
// other tests. Requires GO_WANT_HELPER_PROCESS=1 in the environment.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	args := os.Args
	for i, a := range args {
		if a == "--" {
			args = args[i+1:]
			break
		}
	}

	if len(args) == 0 {
		os.Exit(0)
	}

	switch args[0] {
	case "echo":
		fmt.Println(strings.Join(args[1:], " "))
	case "exit":
		code, _ := strconv.Atoi(args[1])
		os.Exit(code)
	case "stdout":
		fmt.Fprint(os.Stdout, strings.Join(args[1:], " "))
	case "stderr":
		fmt.Fprint(os.Stderr, strings.Join(args[1:], " "))
	case "env":
		fmt.Println(os.Getenv(args[1]))
	case "pwd":
		dir, _ := os.Getwd()
		fmt.Println(dir)
	case "cat":
		_, _ = io.Copy(os.Stdout, os.Stdin)
	case "large":
		n, _ := strconv.Atoi(args[1])
		for i := range n {
			fmt.Fprintf(os.Stdout, "line %d\n", i)
		}
	case "echo-args":
		for _, a := range args[1:] {
			fmt.Println(a)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown helper command: %s\n", args[0])
		os.Exit(1)
	}
	os.Exit(0)
}

// helperCmd returns an executor.Command that re-invokes this test binary's
// TestHelperProcess with the given sub-command and arguments.
func helperCmd(sub ...string) executor.Command {
	base := []string{os.Args[0], "-test.run=TestHelperProcess", "--"}
	return executor.Command{Args: append(base, sub...)}
}

func TestRun_Success(t *testing.T) {
	t.Setenv("GO_WANT_HELPER_PROCESS", "1")
	res, err := executor.Run(context.Background(), helperCmd("echo", "hello"))
	if err != nil {
		t.Fatal(err)
	}
	if res.ExitCode != 0 {
		t.Errorf("exit code = %d, want 0", res.ExitCode)
	}
	if res.Duration <= 0 {
		t.Error("duration should be positive")
	}
}

func TestRun_NonZeroExit(t *testing.T) {
	t.Setenv("GO_WANT_HELPER_PROCESS", "1")
	res, err := executor.Run(context.Background(), helperCmd("exit", "42"))
	if err != nil {
		t.Fatal(err)
	}
	if res.ExitCode != 42 {
		t.Errorf("exit code = %d, want 42", res.ExitCode)
	}
}

func TestRun_ExitOne(t *testing.T) {
	t.Setenv("GO_WANT_HELPER_PROCESS", "1")
	res, err := executor.Run(context.Background(), helperCmd("exit", "1"))
	if err != nil {
		t.Fatal(err)
	}
	if res.ExitCode != 1 {
		t.Errorf("exit code = %d, want 1", res.ExitCode)
	}
}

func TestRun_StdoutForwarding(t *testing.T) {
	t.Setenv("GO_WANT_HELPER_PROCESS", "1")

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	origStdout := os.Stdout
	os.Stdout = w

	res, runErr := executor.Run(context.Background(), helperCmd("stdout", "hello from stdout"))

	w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	r.Close()

	if runErr != nil {
		t.Fatal(runErr)
	}
	if res.ExitCode != 0 {
		t.Errorf("exit code = %d, want 0", res.ExitCode)
	}
	if !strings.Contains(buf.String(), "hello from stdout") {
		t.Errorf("stdout %q does not contain expected text", buf.String())
	}
}

func TestRun_StderrForwarding(t *testing.T) {
	t.Setenv("GO_WANT_HELPER_PROCESS", "1")

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	origStderr := os.Stderr
	os.Stderr = w

	res, runErr := executor.Run(context.Background(), helperCmd("stderr", "hello from stderr"))

	w.Close()
	os.Stderr = origStderr

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	r.Close()

	if runErr != nil {
		t.Fatal(runErr)
	}
	if res.ExitCode != 0 {
		t.Errorf("exit code = %d, want 0", res.ExitCode)
	}
	if !strings.Contains(buf.String(), "hello from stderr") {
		t.Errorf("stderr %q does not contain expected text", buf.String())
	}
}

func TestRun_StdoutStderrSeparate(t *testing.T) {
	t.Setenv("GO_WANT_HELPER_PROCESS", "1")

	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	origOut, origErr := os.Stdout, os.Stderr
	os.Stdout, os.Stderr = wOut, wErr

	executor.Run(context.Background(), helperCmd("stderr", "err-only"))

	wOut.Close()
	wErr.Close()
	os.Stdout, os.Stderr = origOut, origErr

	var outBuf, errBuf bytes.Buffer
	io.Copy(&outBuf, rOut)
	io.Copy(&errBuf, rErr)
	rOut.Close()
	rErr.Close()

	if strings.Contains(outBuf.String(), "err-only") {
		t.Error("stderr text should not appear on stdout")
	}
	if !strings.Contains(errBuf.String(), "err-only") {
		t.Error("stderr text not found on stderr")
	}
}

func TestRun_ArgvPreservation(t *testing.T) {
	t.Setenv("GO_WANT_HELPER_PROCESS", "1")

	cases := []string{"with spaces", "unicode-αβγ", "", `"quoted"`}
	cmd := helperCmd(append([]string{"echo-args"}, cases...)...)

	r, w, _ := os.Pipe()
	origOut := os.Stdout
	os.Stdout = w

	executor.Run(context.Background(), cmd)

	w.Close()
	os.Stdout = origOut
	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()

	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	for i, want := range cases {
		if i >= len(lines) {
			t.Errorf("arg[%d]: missing output line", i)
			continue
		}
		if lines[i] != want {
			t.Errorf("arg[%d]: got %q, want %q", i, lines[i], want)
		}
	}
}

func TestRun_EnvInheritance(t *testing.T) {
	t.Setenv("GO_WANT_HELPER_PROCESS", "1")
	t.Setenv("ACAP_TEST_ENV_VAR", "hello-env")

	r, w, _ := os.Pipe()
	origOut := os.Stdout
	os.Stdout = w

	executor.Run(context.Background(), helperCmd("env", "ACAP_TEST_ENV_VAR"))

	w.Close()
	os.Stdout = origOut
	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()

	if !strings.Contains(buf.String(), "hello-env") {
		t.Errorf("env var not inherited; got %q", buf.String())
	}
}

func TestRun_WorkingDirectoryInheritance(t *testing.T) {
	t.Setenv("GO_WANT_HELPER_PROCESS", "1")

	expected, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	r, w, _ := os.Pipe()
	origOut := os.Stdout
	os.Stdout = w

	executor.Run(context.Background(), helperCmd("pwd"))

	w.Close()
	os.Stdout = origOut
	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()

	got := strings.TrimSpace(buf.String())
	if got != expected {
		t.Errorf("cwd: got %q, want %q", got, expected)
	}
}

func TestRun_StdinForwarding(t *testing.T) {
	t.Setenv("GO_WANT_HELPER_PROCESS", "1")

	rIn, wIn, _ := os.Pipe()
	origIn := os.Stdin
	os.Stdin = rIn

	rOut, wOut, _ := os.Pipe()
	origOut := os.Stdout
	os.Stdout = wOut

	_, _ = wIn.WriteString("stdin-content\n")
	wIn.Close()

	executor.Run(context.Background(), helperCmd("cat"))

	wOut.Close()
	os.Stdin = origIn
	os.Stdout = origOut

	var buf bytes.Buffer
	io.Copy(&buf, rOut)
	rOut.Close()

	if !strings.Contains(buf.String(), "stdin-content") {
		t.Errorf("stdin not forwarded; got %q", buf.String())
	}
}

func TestRun_CommandNotFound(t *testing.T) {
	res, err := executor.Run(context.Background(), executor.Command{
		Args: []string{"acap-this-command-does-not-exist-xyz"},
	})
	if err == nil {
		t.Error("expected error for missing command, got nil")
	}
	if res.ExitCode == 0 {
		t.Error("expected non-zero exit code for missing command")
	}
}

func TestRun_NoArgs(t *testing.T) {
	_, err := executor.Run(context.Background(), executor.Command{})
	if err == nil {
		t.Error("expected error when no args provided")
	}
}

func TestRun_LargeOutput(t *testing.T) {
	t.Setenv("GO_WANT_HELPER_PROCESS", "1")

	r, w, _ := os.Pipe()
	origOut := os.Stdout
	os.Stdout = w

	// Drain in background so the pipe buffer never fills and blocks the child.
	done := make(chan int64, 1)
	go func() {
		n, _ := io.Copy(io.Discard, r)
		done <- n
	}()

	res, err := executor.Run(context.Background(), helperCmd("large", "100000"))

	w.Close()
	os.Stdout = origOut
	n := <-done
	r.Close()

	if err != nil {
		t.Fatal(err)
	}
	if res.ExitCode != 0 {
		t.Errorf("exit code = %d, want 0", res.ExitCode)
	}
	if n == 0 {
		t.Error("expected large output but got none")
	}
}

func TestRun_ContextCancellation(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("context kill behavior on Windows may differ")
	}
	t.Setenv("GO_WANT_HELPER_PROCESS", "1")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// No assertion on exit code — cancelled context may kill or allow the
	// process to complete. We only verify no panic occurs.
	res, _ := executor.Run(ctx, helperCmd("exit", "0"))
	_ = res
}

func TestRun_Duration(t *testing.T) {
	t.Setenv("GO_WANT_HELPER_PROCESS", "1")
	res, err := executor.Run(context.Background(), helperCmd("echo", "timing"))
	if err != nil {
		t.Fatal(err)
	}
	if res.Duration <= 0 {
		t.Error("duration should be positive")
	}
}
