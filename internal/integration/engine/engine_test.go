package engine

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/taqu/agentcap/internal/integration/protocol"
)

func echoCmd() []string {
	if runtime.GOOS == "windows" {
		return []string{"cmd", "/c", "echo hello"}
	}
	return []string{"echo", "hello"}
}

func appendCmd(path string) []string {
	if runtime.GOOS == "windows" {
		return []string{"cmd", "/c", "echo x >> " + path}
	}
	return []string{"sh", "-c", "echo x >> " + path}
}

func exitCodeCmd(code int) []string {
	if runtime.GOOS == "windows" {
		return []string{"cmd", "/c", "exit " + itoa(code)}
	}
	return []string{"sh", "-c", "exit " + itoa(code)}
}

func itoa(n int) string {
	return strings.TrimSpace(strings.Join([]string{""}, "")) + func() string {
		if n == 0 {
			return "0"
		}
		s := ""
		for n > 0 {
			s = string(rune('0'+n%10)) + s
			n /= 10
		}
		return s
	}()
}

func TestEngineBasic(t *testing.T) {
	req := &protocol.ToolRequest{
		Protocol:   protocol.Version,
		Command:    echoCmd(),
		WorkingDir: t.TempDir(),
	}
	resp, err := Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Errorf("ExitCode = %d, want 0", resp.ExitCode)
	}
	if resp.Stdout == "" {
		t.Error("Stdout is empty, expected non-empty")
	}
}

func TestEngineExactlyOnce(t *testing.T) {
	dir := t.TempDir()
	counterFile := filepath.Join(dir, "counter.txt")

	cmd := appendCmd(counterFile)
	req := &protocol.ToolRequest{
		Protocol:   protocol.Version,
		Command:    cmd,
		WorkingDir: dir,
	}
	_, err := Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}

	data, err := os.ReadFile(counterFile)
	if err != nil {
		t.Fatalf("counter.txt not created: %v", err)
	}
	lines := strings.Split(strings.TrimRight(string(data), "\r\n"), "\n")
	// Filter empty lines.
	var nonEmpty []string
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			nonEmpty = append(nonEmpty, l)
		}
	}
	if len(nonEmpty) != 1 {
		t.Errorf("expected exactly 1 line in counter.txt, got %d: %q", len(nonEmpty), string(data))
	}
}

func TestEngineExitCodePreserved(t *testing.T) {
	req := &protocol.ToolRequest{
		Protocol:   protocol.Version,
		Command:    exitCodeCmd(7),
		WorkingDir: t.TempDir(),
	}
	resp, err := Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if resp.ExitCode != 7 {
		t.Errorf("ExitCode = %d, want 7", resp.ExitCode)
	}
}

func TestEngineFailOpen(t *testing.T) {
	// Use a non-existent working dir as root — store will fail to init,
	// but we should still get the command output.
	req := &protocol.ToolRequest{
		Protocol:   protocol.Version,
		Command:    echoCmd(),
		WorkingDir: filepath.Join(t.TempDir(), "nonexistent"),
	}
	// Even with a bad working dir for store, output should be returned.
	// The command itself will run in the current dir since the dir doesn't exist.
	// Actually, let's use a valid dir for the command but an invalid ACAP_ROOT.
	os.Setenv("ACAP_ROOT", "/nonexistent/path/that/does/not/exist/xyz")
	defer os.Unsetenv("ACAP_ROOT")

	req.WorkingDir = t.TempDir()
	resp, err := Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	// Even if store fails, we should get output.
	if resp.Stdout == "" {
		t.Error("Stdout is empty even with bad store root; expected fail-open")
	}
}
