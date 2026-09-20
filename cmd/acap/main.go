package main

import (
	"context"
	"fmt"
	"os"

	"github.com/taqu/agentcap/internal/executor"
)

const usage = `acap — command execution wrapper for AI coding agents

Usage:
  acap run <command> [args...]

Examples:
  acap run ls -la
  acap run go test ./...
  acap run sh -c 'echo foo | grep foo'

Environment:
  ACAP_DEBUG=1   Print argv, exit code, and duration to stderr.
`

const runUsage = `Usage: acap run <command> [args...]

Executes <command> directly (no shell interpretation) and forwards
stdin, stdout, and stderr. The child's exit code is preserved.
`

func main() {
	args := os.Args[1:]

	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprint(os.Stdout, usage)
		os.Exit(0)
	}

	switch args[0] {
	case "run":
		runCmd(args[1:])
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

	cmd := executor.Command{Args: args}
	result, err := executor.Run(context.Background(), cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(result.ExitCode)
	}
	os.Exit(result.ExitCode)
}
