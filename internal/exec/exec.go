// Package exec provides a capturing command executor that buffers child
// stdout/stderr instead of streaming directly to the terminal.
package exec

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

// MaxOutputBytes is the maximum number of bytes captured from stdout or stderr.
const MaxOutputBytes = 10 * 1024 * 1024 // 10 MB

// Options holds optional configuration for Run.
type Options struct {
	// StdoutSink, if set, receives a tee of the raw stdout stream (unlimited).
	StdoutSink io.Writer
	// StderrSink, if set, receives a tee of the raw stderr stream (unlimited).
	StderrSink io.Writer
	// Dir, if non-empty, sets the working directory for the child process.
	Dir string
}

// Result holds the outcome of a completed child execution.
type Result struct {
	Args      []string
	ExitCode  int
	Stdout    []byte
	Stderr    []byte
	Duration  time.Duration
	Truncated bool
}

// Run executes the command specified by args, captures stdout/stderr (up to
// MaxOutputBytes each), and returns the result. ctx cancellation kills the
// child. stdin is forwarded from os.Stdin. opts may be nil.
func Run(ctx context.Context, args []string, opts *Options) (*Result, error) {
	if len(args) == 0 {
		return &Result{ExitCode: 1}, fmt.Errorf("no command specified")
	}

	name := args[0]
	rest := args[1:]

	c := exec.CommandContext(ctx, name, rest...)
	c.Stdin = os.Stdin
	if opts != nil && opts.Dir != "" {
		c.Dir = opts.Dir
	}

	// Limit stdout/stderr capture to MaxOutputBytes.
	stdoutBuf := &limitedBuffer{limit: MaxOutputBytes}
	stderrBuf := &limitedBuffer{limit: MaxOutputBytes}

	if opts != nil && opts.StdoutSink != nil {
		c.Stdout = io.MultiWriter(stdoutBuf, opts.StdoutSink)
	} else {
		c.Stdout = stdoutBuf
	}

	if opts != nil && opts.StderrSink != nil {
		c.Stderr = io.MultiWriter(stderrBuf, opts.StderrSink)
	} else {
		c.Stderr = stderrBuf
	}

	start := time.Now()

	if err := c.Start(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return &Result{Args: args, ExitCode: 127}, fmt.Errorf("acap: command not found: %s", name)
		}
		return &Result{Args: args, ExitCode: 1}, fmt.Errorf("acap: %w", err)
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

	return &Result{
		Args:      args,
		ExitCode:  exitCode,
		Stdout:    stdoutBuf.Bytes(),
		Stderr:    stderrBuf.Bytes(),
		Duration:  duration,
		Truncated: stdoutBuf.truncated || stderrBuf.truncated,
	}, nil
}

// limitedBuffer is an io.Writer that captures up to limit bytes.
type limitedBuffer struct {
	buf       bytes.Buffer
	limit     int
	truncated bool
}

func (lb *limitedBuffer) Write(p []byte) (int, error) {
	if lb.truncated {
		return len(p), nil // silently discard after cap
	}
	remaining := lb.limit - lb.buf.Len()
	if remaining <= 0 {
		lb.truncated = true
		return len(p), nil
	}
	if len(p) > remaining {
		lb.truncated = true
		p = p[:remaining]
	}
	return lb.buf.Write(p)
}

func (lb *limitedBuffer) Bytes() []byte {
	return lb.buf.Bytes()
}

// WriteTo satisfies io.WriterTo (unused but keeps the type usable).
var _ io.Writer = (*limitedBuffer)(nil)
