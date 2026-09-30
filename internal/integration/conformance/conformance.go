// Package conformance supplies behavioral tests reusable by shell adapters.
package conformance

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/taqu/agentcap/internal/integration/protocol"
	"github.com/taqu/agentcap/internal/store"
)

type Invoke func(context.Context, string, string, string) *protocol.ToolResponse

func Run(t *testing.T, invoke Invoke) {
	t.Helper()
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, ".git"), 0755)
	nested := filepath.Join(root, "nested", "path")
	os.MkdirAll(nested, 0755)
	t.Run("exit-and-streams", func(t *testing.T) {
		for _, code := range []int{0, 1, 2, 7, 127} {
			resp := invoke(context.Background(), fmt.Sprintf("printf 'out'; printf 'err' >&2; exit %d", code), nested, "streams")
			if resp.ExitCode != code || !strings.Contains(resp.Stdout, "out") || !strings.Contains(resp.Stdout, "err") || resp.ResultID == "" {
				t.Fatalf("code %d: %+v", code, resp)
			}
			st, err := store.Open(root)
			if err != nil {
				t.Fatal(err)
			}
			entry, err := st.Open(resp.ResultID)
			if err != nil {
				t.Fatal(err)
			}
			stdout, _ := os.ReadFile(entry.StdoutPath())
			stderr, _ := os.ReadFile(entry.StderrPath())
			st.Close()
			if string(stdout) != "out" || string(stderr) != "err" {
				t.Fatalf("raw capture %q %q", stdout, stderr)
			}
		}
	})
	t.Run("exactly-once-success-and-failure", func(t *testing.T) {
		for _, code := range []int{0, 7} {
			name := fmt.Sprintf("marker-%d", code)
			resp := invoke(context.Background(), fmt.Sprintf("printf 'x' >> %s; exit %d", name, code), nested, "once")
			data, _ := os.ReadFile(filepath.Join(nested, name))
			if string(data) != "x" || resp.ExitCode != code {
				t.Fatalf("executed incorrectly: %q %+v", data, resp)
			}
		}
	})
	t.Run("cwd-and-environment", func(t *testing.T) {
		t.Setenv("ACAP_CONFORMANCE_VALUE", "hello world")
		resp := invoke(context.Background(), "printf '%s' \"$ACAP_CONFORMANCE_VALUE\"; printf 'ok' > cwd-marker", nested, "env")
		if !strings.Contains(resp.Stdout, "hello world") {
			t.Fatal(resp.Stdout)
		}
		if data, _ := os.ReadFile(filepath.Join(nested, "cwd-marker")); string(data) != "ok" {
			t.Fatal("cwd lost")
		}
	})
	t.Run("session-and-project-isolation", func(t *testing.T) {
		a := invoke(context.Background(), "printf 'session-output'", nested, "same")
		b := invoke(context.Background(), "printf 'session-output'", nested, "same")
		c := invoke(context.Background(), "printf 'session-output'", nested, "other")
		if a.Presentation == "unchanged" || b.Presentation != "unchanged" || c.Presentation == "unchanged" {
			t.Fatalf("session leak: %+v %+v %+v", a, b, c)
		}
		other := t.TempDir()
		d := invoke(context.Background(), "printf 'session-output'", other, "same")
		if d.Presentation == "unchanged" {
			t.Fatal("project leak")
		}
		st, _ := store.Open(root)
		defer st.Close()
		ae, _ := st.Open(a.ResultID)
		st2, _ := store.Open(other)
		defer st2.Close()
		de, _ := st2.Open(d.ResultID)
		if ae.Meta.SessionID == de.Meta.SessionID {
			t.Fatal("session mapping shared across projects")
		}
		for range 2 {
			resp := invoke(context.Background(), "printf 'stateless'", nested, "")
			if resp.Presentation == "unchanged" {
				t.Fatal("missing session used shared state")
			}
		}
	})
	t.Run("fail-open-without-reexecution", func(t *testing.T) {
		bad := t.TempDir()
		os.WriteFile(filepath.Join(bad, ".acap"), []byte("not a directory"), 0600)
		resp := invoke(context.Background(), "printf 'raw-out'; printf 'raw-err' >&2; printf 'x' >> marker; exit 7", bad, "failopen")
		data, _ := os.ReadFile(filepath.Join(bad, "marker"))
		if resp.ExitCode != 7 || resp.Stdout != "raw-out" || !strings.Contains(resp.Stderr, "raw-err") || string(data) != "x" {
			t.Fatalf("fail-open: %+v marker=%q", resp, data)
		}
	})
	t.Run("concurrent-invocations", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 4; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				token := fmt.Sprintf("unique-%d", i)
				resp := invoke(context.Background(), "printf '"+token+"'", nested, "parallel")
				if resp.ExitCode != 0 || !strings.Contains(resp.Stdout, token) || resp.ResultID == "" {
					t.Errorf("swapped or failed response: %+v", resp)
				}
			}(i)
		}
		wg.Wait()
	})
	t.Run("cancellation", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer cancel()
		started := time.Now()
		resp := invoke(ctx, "sleep 30", nested, "cancel")
		if time.Since(started) > 5*time.Second || resp.ExitCode == 0 {
			t.Fatalf("cancellation failed: %+v", resp)
		}
	})
}
