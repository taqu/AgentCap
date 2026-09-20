package reduce

import (
	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/exec"
)

// WcReducer handles wc output. wc output is already compact so we only
// strip ANSI sequences and return as-is.
type WcReducer struct{}

func (w *WcReducer) Reduce(r *exec.Result) *ReducedResult {
	raw := clean.StripANSI(r.Stdout)
	rawBytes := len(r.Stdout)
	stderr := buildStderr(r.Stderr)

	out := string(raw) + stderr
	return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
}
