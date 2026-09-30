//go:build !windows

package exec

import (
	"context"
	"os"
	osexec "os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

// The helper runs a process tree through Run, the way adapter wrappers do.
func TestMain(m *testing.M) {
	if pidFile := os.Getenv("ACAP_EXEC_TREE_HELPER"); pidFile != "" {
		Run(context.Background(), []string{"sh", "-c", "sh -c 'echo $$ > " + pidFile + "; exec sleep 30'; true"}, &Options{ProcessTree: true})
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func TestTerminatedWrapperDoesNotOrphanTree(t *testing.T) {
	for _, sig := range []syscall.Signal{syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT} {
		t.Run(sig.String(), func(t *testing.T) {
			pidFile := filepath.Join(t.TempDir(), "pid")
			helper := osexec.Command(os.Args[0], "-test.run=^$")
			helper.Env = append(os.Environ(), "ACAP_EXEC_TREE_HELPER="+pidFile)
			if err := helper.Start(); err != nil {
				t.Fatal(err)
			}
			var pid int
			for i := 0; i < 100 && pid == 0; i++ {
				time.Sleep(50 * time.Millisecond)
				data, _ := os.ReadFile(pidFile)
				pid, _ = strconv.Atoi(strings.TrimSpace(string(data)))
			}
			if pid == 0 {
				helper.Process.Kill()
				t.Fatal("tree did not start")
			}
			helper.Process.Signal(sig)
			done := make(chan struct{})
			go func() { helper.Wait(); close(done) }()
			select {
			case <-done:
			case <-time.After(10 * time.Second):
				helper.Process.Kill()
				t.Fatal("wrapper did not finish after signal")
			}
			for i := 0; i < 40; i++ {
				if syscall.Kill(pid, 0) != nil {
					return
				}
				time.Sleep(50 * time.Millisecond)
			}
			syscall.Kill(pid, syscall.SIGKILL)
			t.Fatalf("grandchild %d orphaned after %s", pid, sig)
		})
	}
}
