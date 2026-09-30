//go:build windows

package exec

import (
	"os"
	"os/exec"
	"os/signal"
)

func notifySignals(ch chan os.Signal) {
	signal.Notify(ch, os.Interrupt)
}

func forwardSignal(c *exec.Cmd, sig os.Signal, tree bool) {
	if c.Process != nil {
		_ = c.Process.Signal(sig)
	}
}
