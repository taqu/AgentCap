package executor

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

// Command holds the argv for a child process.
type Command struct {
	Args []string
}

// Result holds the outcome of a completed child execution.
type Result struct {
	ExitCode int
	Duration time.Duration
}

var debugMode = os.Getenv("ACAP_DEBUG") != ""

// Run executes cmd as a direct child process, forwarding stdin/stdout/stderr,
// and returns the exit code and duration. ctx cancellation kills the child.
func Run(ctx context.Context, cmd Command) (Result, error) {
	if len(cmd.Args) == 0 {
		return Result{ExitCode: 1}, fmt.Errorf("no command specified")
	}

	name := cmd.Args[0]
	args := cmd.Args[1:]

	if debugMode {
		fmt.Fprintf(os.Stderr, "acap: args=%v\n", cmd.Args)
	}

	c := exec.CommandContext(ctx, name, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	start := time.Now()

	if err := c.Start(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return Result{ExitCode: 127}, fmt.Errorf("acap: command not found: %s", name)
		}
		return Result{ExitCode: 1}, fmt.Errorf("acap: %w", err)
	}

	// Forward termination signals to the child.
	sigs := make(chan os.Signal, 1)
	notifySignals(sigs)
	sigDone := make(chan struct{})
	go func() {
		defer signal.Stop(sigs)
		select {
		case sig := <-sigs:
			if c.Process != nil {
				_ = c.Process.Signal(sig)
			}
		case <-sigDone:
		}
	}()

	waitErr := c.Wait()
	close(sigDone)

	duration := time.Since(start)

	exitCode := 0
	if waitErr != nil {
		var exitErr *exec.ExitError
		if errors.As(waitErr, &exitErr) {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}

	result := Result{ExitCode: exitCode, Duration: duration}

	if debugMode {
		fmt.Fprintf(os.Stderr, "acap: exit_code=%d duration=%s\n", exitCode, duration)
	}

	return result, nil
}
