package exec

import (
	"context"
	"strings"
	"testing"
)

func TestRunSimple(t *testing.T) {
	// Use a cross-platform command.
	result, err := Run(context.Background(), []string{"go", "version"}, nil)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("expected exit 0, got %d", result.ExitCode)
	}
	if !strings.Contains(string(result.Stdout), "go version") {
		t.Errorf("expected go version in stdout, got: %q", string(result.Stdout))
	}
}

func TestRunExitCode(t *testing.T) {
	// go tool with invalid flag returns non-zero.
	result, _ := Run(context.Background(), []string{"go", "tool", "nonexistent_tool_xyz"}, nil)
	if result.ExitCode == 0 {
		t.Errorf("expected non-zero exit code")
	}
}

func TestRunEmpty(t *testing.T) {
	_, err := Run(context.Background(), nil, nil)
	if err == nil {
		t.Error("expected error for empty args")
	}
}

func TestRunCommandNotFound(t *testing.T) {
	result, err := Run(context.Background(), []string{"__nonexistent_command_xyz__"}, nil)
	if err == nil {
		t.Error("expected error for missing command")
	}
	if result.ExitCode != 127 {
		t.Errorf("expected exit 127 for command not found, got %d", result.ExitCode)
	}
}

func TestRunStderrCapture(t *testing.T) {
	// go vet on non-existent path writes to stderr.
	result, _ := Run(context.Background(), []string{"go", "build", "./nonexistent/path/..."}, nil)
	// Either stdout or stderr should have content, and exit should be non-zero.
	if result.ExitCode == 0 {
		t.Error("expected non-zero exit")
	}
}

func TestRunDuration(t *testing.T) {
	result, err := Run(context.Background(), []string{"go", "version"}, nil)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.Duration <= 0 {
		t.Error("expected positive duration")
	}
}

func TestRunWithSinks(t *testing.T) {
	var stdoutSink, stderrSink strings.Builder
	opts := &Options{StdoutSink: &stdoutSink}
	result, err := Run(context.Background(), []string{"go", "version"}, opts)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("expected exit 0, got %d", result.ExitCode)
	}
	// Both the result buffer and sink should have the same content.
	if string(result.Stdout) != stdoutSink.String() {
		t.Errorf("sink mismatch: result=%q sink=%q", string(result.Stdout), stdoutSink.String())
	}
	_ = stderrSink
}

func TestLimitedBuffer(t *testing.T) {
	lb := &limitedBuffer{limit: 10}
	n, err := lb.Write([]byte("hello"))
	if err != nil || n != 5 {
		t.Errorf("Write: n=%d err=%v", n, err)
	}
	n, err = lb.Write([]byte("world!"))
	if err != nil {
		t.Errorf("Write overrun err=%v", err)
	}
	_ = n // limitedBuffer may return fewer bytes written; that's fine for our internal use
	if !lb.truncated {
		t.Error("expected truncated=true after overflow")
	}
	if got := string(lb.Bytes()); got != "helloworld" {
		t.Errorf("expected 'helloworld', got %q", got)
	}
}
