//go:build !windows

package exec

import (
	"os"
	"os/signal"
	"syscall"
)

func notifySignals(ch chan os.Signal) {
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
}
