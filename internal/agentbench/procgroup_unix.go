//go:build unix

package agentbench

import (
	"os/exec"
	"syscall"
)

// isolateGroup makes timeout cancellation kill the agent's whole process tree.
func isolateGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error { return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) }
}
