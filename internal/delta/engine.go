package delta

import (
	"context"
	"fmt"
	"os"

	"github.com/taqu/agentcap/internal/store"
)

// Compare decides what to present given a baseline and current result.
// dr may be nil. Falls back to full capsule when no safe delta exists.
func Compare(
	ctx context.Context,
	baseline *store.Entry,
	current *store.Entry,
	currentCapsule string,
	dr Reducer,
	stdoutHash, stderrHash string,
) *Result {
	exitChanged := baseline.Meta.ExitCode != current.Meta.ExitCode

	// Exact equality check (requires stored hashes in baseline).
	if baseline.Meta.StdoutHash != "" && baseline.Meta.StderrHash != "" {
		stdoutEqual := stdoutHash == baseline.Meta.StdoutHash
		stderrEqual := stderrHash == baseline.Meta.StderrHash
		if !exitChanged && stdoutEqual && stderrEqual {
			return &Result{
				Output: fmt.Sprintf("@acap unchanged from %s\nexit=%d\n",
					baseline.Meta.ID, current.Meta.ExitCode),
				Presentation: PresentationUnchanged,
			}
		}
	}

	// Specialized delta.
	if dr != nil {
		if res, err := dr.Delta(ctx, baseline, current); err == nil && res != nil {
			if cheaper(res.Output, currentCapsule) {
				return res
			}
		}
	}

	// Generic text delta.
	stderrEqual := stderrHash == baseline.Meta.StderrHash && stderrHash != "" && baseline.Meta.StderrHash != ""
	if res := genericDelta(baseline, current, exitChanged, stderrEqual); res != nil {
		if cheaper(res.Output, currentCapsule) {
			return res
		}
	}

	// Full capsule fallback.
	return &Result{
		Output:       currentCapsule,
		Presentation: PresentationFull,
	}
}

func cheaper(delta, full string) bool {
	return len(full) > 0 && float64(len(delta)) < float64(len(full))*deltaThreshold
}

func readFile(path string) []byte {
	data, _ := os.ReadFile(path)
	return data
}
