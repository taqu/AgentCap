//go:build !windows

package exec

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// SIGHUP is included: a terminal owner hangs up its foreground group, which
// no longer contains a child placed in its own process group.
func notifySignals(ch chan os.Signal) {
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
}

func forwardSignal(c *exec.Cmd, sig os.Signal, tree bool) {
	if c.Process == nil {
		return
	}
	if s, ok := sig.(syscall.Signal); ok && tree {
		_ = syscall.Kill(-c.Process.Pid, s)
		return
	}
	_ = c.Process.Signal(sig)
}
