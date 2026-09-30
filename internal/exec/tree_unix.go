//go:build !windows

package exec

import (
	"os/exec"
	"syscall"
	"time"
)

func configureProcessTree(c *exec.Cmd) {
	c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	c.Cancel = func() error { return syscall.Kill(-c.Process.Pid, syscall.SIGKILL) }
	c.WaitDelay = 2 * time.Second
}
