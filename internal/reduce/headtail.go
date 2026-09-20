package reduce

import (
	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/exec"
)

// HeadTailReducer handles head and tail output.
// These commands already constrain output, so we only strip ANSI.
type HeadTailReducer struct{}

func (h *HeadTailReducer) Reduce(r *exec.Result) *ReducedResult {
	raw := clean.StripANSI(r.Stdout)
	rawBytes := len(r.Stdout)
	stderr := buildStderr(r.Stderr)

	// head/tail already limit output — only strip ANSI noise.
	// Only add a header if the output is unusually large (> 4KB).
	out := string(raw) + stderr
	if len(raw) > smallThresholdBytes {
		// Still large — pass through generic.
		g := &GenericReducer{}
		g2 := g.Reduce(r)
		return g2
	}

	return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
}
