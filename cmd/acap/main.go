package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/taqu/agentcap/internal/exec"
	"github.com/taqu/agentcap/internal/reduce"
	"github.com/taqu/agentcap/internal/stats"
)

const usage = `acap — command output compression layer for AI coding agents

Usage:
  acap run <command> [args...]
  acap stats

Examples:
  acap run ls -la
  acap run find . -type f
  acap run rg "Workspace" .
  acap run cat src/main.go

Environment:
  ACAP_DEBUG=1   Print reducer debug info to stderr.
`

const runUsage = `Usage: acap run <command> [args...]

Executes <command> directly (no shell interpretation), reduces the output
for LLM consumption, and writes it to stdout.  The child's exit code is
preserved.
`

var debugMode = os.Getenv("ACAP_DEBUG") == "1"

func main() {
	args := os.Args[1:]

	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprint(os.Stdout, usage)
		os.Exit(0)
	}

	switch args[0] {
	case "run":
		runCmd(args[1:])
	case "stats":
		statsCmd()
	default:
		fmt.Fprintf(os.Stderr, "acap: unknown command %q\n\n%s", args[0], usage)
		os.Exit(1)
	}
}

func runCmd(args []string) {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprint(os.Stdout, runUsage)
		os.Exit(0)
	}

	result, err := exec.Run(context.Background(), args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(result.ExitCode)
	}

	execDuration := result.Duration

	t0 := time.Now()
	reducer := reduce.Select(args)
	reduced := reducer.Reduce(result)
	reduceDuration := time.Since(t0)

	// Write reduced output to stdout.
	fmt.Fprint(os.Stdout, reduced.Output)

	// Persist statistics.
	_ = stats.Record(reduced.RawBytes, reduced.RetBytes)

	if debugMode {
		reducerName := reducerTypeName(reducer)
		fmt.Fprintf(os.Stderr, "acap: cmd=%s reducer=%s raw=%d ret=%d exec=%s reduce=%s\n",
			args[0], reducerName, reduced.RawBytes, reduced.RetBytes,
			execDuration.Round(time.Millisecond),
			reduceDuration.Round(time.Millisecond),
		)
	}

	os.Exit(result.ExitCode)
}

func statsCmd() {
	s, err := stats.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: stats: %v\n", err)
		os.Exit(1)
	}

	saved := s.RawBytes - s.RetBytes
	var reduction float64
	if s.RawBytes > 0 {
		reduction = float64(saved) / float64(s.RawBytes) * 100
	}

	fmt.Printf("commands: %d\n", s.Commands)
	fmt.Printf("raw:      %s\n", stats.FormatBytes(s.RawBytes))
	fmt.Printf("returned: %s\n", stats.FormatBytes(s.RetBytes))
	fmt.Printf("saved:    %s\n", stats.FormatBytes(saved))
	fmt.Printf("reduction: %.1f%%\n", reduction)
}

func reducerTypeName(r reduce.Reducer) string {
	switch r.(type) {
	case *reduce.LsReducer:
		return "ls"
	case *reduce.FindReducer:
		return "find"
	case *reduce.GrepReducer:
		return "grep"
	case *reduce.CatReducer:
		return "cat"
	case *reduce.HeadTailReducer:
		return "headtail"
	case *reduce.TreeReducer:
		return "tree"
	case *reduce.DuReducer:
		return "du"
	case *reduce.WcReducer:
		return "wc"
	case *reduce.GenericReducer:
		return "generic"
	default:
		return "unknown"
	}
}
