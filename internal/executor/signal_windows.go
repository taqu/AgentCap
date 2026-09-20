//go:build windows

package executor

import (
	"os"
	"os/signal"
)

func notifySignals(ch chan os.Signal) {
	signal.Notify(ch, os.Interrupt)
}
