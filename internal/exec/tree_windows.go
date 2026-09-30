package exec

import (
	"os/exec"
	"strconv"
	"time"
)

func configureProcessTree(c *exec.Cmd) {
	c.Cancel = func() error {
		err := exec.Command("taskkill.exe", "/PID", strconv.Itoa(c.Process.Pid), "/T", "/F").Run()
		if err != nil {
			return c.Process.Kill()
		}
		return nil
	}
	c.WaitDelay = 2 * time.Second
}
