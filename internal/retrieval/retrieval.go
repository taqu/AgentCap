// Package retrieval implements the normal agent-visible result recovery paths.
package retrieval

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/taqu/agentcap/internal/store"
)

// Stream identifies a captured raw command stream.
type Stream string

const (
	Stdout Stream = "stdout"
	Stderr Stream = "stderr"
)

// Show writes the stored reduced capsule exactly as the default `acap show`
// operation does and records the visible byte count in normal AgentCap stats.
func Show(w io.Writer, st *store.Store, id string) (int64, error) {
	started := time.Now()
	entry, err := st.Open(id)
	if err != nil {
		return 0, err
	}
	capsule, err := entry.Capsule()
	if err != nil {
		return 0, err
	}
	n, err := io.WriteString(w, capsule)
	if err != nil {
		return int64(n), err
	}
	_ = st.RecordShow(n)
	_ = st.RecordProcessing(time.Since(started).Nanoseconds())
	return int64(n), nil
}

// Raw writes a captured stdout or stderr stream without rerunning the command
// and records the visible byte count in normal AgentCap stats.
func Raw(w io.Writer, st *store.Store, id string, stream Stream) (int64, error) {
	started := time.Now()
	entry, err := st.Open(id)
	if err != nil {
		return 0, err
	}
	var path string
	switch stream {
	case "", Stdout:
		path = entry.StdoutPath()
	case Stderr:
		path = entry.StderrPath()
	default:
		return 0, fmt.Errorf("unsupported raw stream %q", stream)
	}
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	n, err := io.Copy(w, f)
	if err != nil {
		return n, err
	}
	_ = st.RecordRaw(int(n))
	_ = st.RecordProcessing(time.Since(started).Nanoseconds())
	return n, nil
}
