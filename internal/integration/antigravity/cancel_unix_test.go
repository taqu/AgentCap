//go:build !windows

package antigravity

import (
	"os"
	osexec "os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

// Antigravity starts run_command as a process-group leader and cancels by
// SIGKILLing that group. The wrapped command must die with it.
func TestGroupKillTerminatesWrappedCommand(t *testing.T) {
	acap := filepath.Join(t.TempDir(), "acap")
	if out, err := osexec.Command("go", "build", "-o", acap, "github.com/taqu/agentcap/cmd/acap").CombinedOutput(); err != nil {
		t.Skip("cannot build acap:", err, string(out))
	}
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".git"), 0755)
	t.Setenv(EnvPermissions, permissionsDeclared)
	response, reason, err := Prepare(hookData("sh -c 'echo $$ > pid; exec sleep 30'; true", dir, "conv", 5000), acap)
	if reason != "intercepted" || err != nil {
		t.Fatal(reason, err)
	}
	wrapper := osexec.Command("sh", "-c", decode(t, response).Overwrite["CommandLine"])
	wrapper.Dir = dir
	wrapper.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := wrapper.Start(); err != nil {
		t.Fatal(err)
	}
	var pid int
	for i := 0; i < 200 && pid == 0; i++ {
		time.Sleep(50 * time.Millisecond)
		data, _ := os.ReadFile(filepath.Join(dir, "pid"))
		pid, _ = strconv.Atoi(strings.TrimSpace(string(data)))
	}
	if pid == 0 {
		syscall.Kill(-wrapper.Process.Pid, syscall.SIGKILL)
		t.Fatal("command did not start")
	}
	syscall.Kill(-wrapper.Process.Pid, syscall.SIGKILL)
	wrapper.Wait()
	for i := 0; i < 40; i++ {
		if syscall.Kill(pid, 0) != nil {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	syscall.Kill(pid, syscall.SIGKILL)
	t.Fatalf("command %d orphaned after group kill", pid)
}
