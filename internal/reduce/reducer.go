// Package reduce provides command-output reducers for AI coding agents.
package reduce

import (
	"path/filepath"

	"github.com/taqu/agentcap/internal/exec"
)

// ReducedResult holds the compressed output and byte-count metadata.
type ReducedResult struct {
	Output   string
	RawBytes int
	RetBytes int
}

// Reducer transforms an exec.Result into a compact ReducedResult.
type Reducer interface {
	Reduce(r *exec.Result) *ReducedResult
}

// Select picks the most appropriate Reducer based on the command name.
func Select(args []string) Reducer {
	if len(args) == 0 {
		return &GenericReducer{}
	}
	name := filepath.Base(args[0])
	switch name {
	case "ls":
		return &LsReducer{}
	case "find":
		return &FindReducer{}
	case "grep", "rg", "ripgrep":
		return &GrepReducer{}
	case "cat":
		return &CatReducer{}
	case "head", "tail":
		return &HeadTailReducer{}
	case "tree":
		return &TreeReducer{}
	case "du":
		return &DuReducer{}
	case "wc":
		return &WcReducer{}
	default:
		return &GenericReducer{}
	}
}

// smallThresholdBytes is the byte threshold below which output is considered small.
const smallThresholdBytes = 4096

// smallThresholdLines is the line threshold below which output is considered small.
const smallThresholdLines = 50
