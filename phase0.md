# AgentCap Phase 0 — Execution Core

## Project

Project name: AgentCap  
CLI command: `acap`  
Language: Go

AgentCap is intended to become a command-output reduction layer for AI coding agents.

However, Phase 0 must NOT implement output compression.

The sole purpose of Phase 0 is to build a reliable command execution foundation that later phases can safely extend.

The implementation should be small, robust, testable, and idiomatic Go.

---

# 1. Phase 0 Goal

Implement the minimum execution core required for future AgentCap development.

Phase 0 must provide:

- a working `acap` CLI
- `acap run <command> [args...]`
- direct child-process execution
- stdout forwarding/capture
- stderr forwarding/capture
- exit-code preservation
- signal handling
- current working directory inheritance
- environment inheritance
- execution timing
- safe handling of large command output
- automated tests

The primary requirement is:

> Running a command through `acap run` should behave as closely as practical to running that command directly.

Phase 0 is infrastructure only.

Do not implement output compression or command-specific behavior.

---

# 2. Core Principle

AgentCap must not change the semantics of the command it executes.

Conceptually:

```text
argv
 |
 v
acap run
 |
 v
executor
 |
 +--> child stdout
 +--> child stderr
 +--> child exit status
 |
 v
caller
```

At this stage there is no reducer, summarizer, classifier, or output transformation layer.

The executor should be designed so those can be added later without rewriting the process-management core.

---

# 3. CLI

Implement:

```bash
acap run <command> [args...]
```

Examples:

```bash
acap run ls -la
acap run printf "hello\n"
acap run go test ./...
acap run rg "Workspace" .
```

Also implement basic help:

```bash
acap --help
acap run --help
```

Avoid introducing a large CLI framework unless clearly justified.

The Go standard library is preferred where practical.

A small CLI dependency is acceptable only if it materially improves maintainability.

---

# 4. Argument Handling

Everything after `run` must be treated as the child command argv.

For example:

```bash
acap run rg "foo bar" src
```

must execute an argv equivalent to:

```text
["rg", "foo bar", "src"]
```

Do not reconstruct a shell command string.

Do not use:

```bash
sh -c
```

or:

```bash
bash -c
```

for normal execution.

Use direct process execution, such as Go's `os/exec`.

This avoids:

- quoting bugs
- shell injection issues
- platform-specific parsing differences
- accidental wildcard expansion
- accidental environment expansion

Shell syntax support is not required in Phase 0.

For example:

```bash
acap run echo foo | grep foo
```

does NOT need to interpret the pipe.

Users or agents can explicitly run a shell themselves if necessary:

```bash
acap run sh -c 'echo foo | grep foo'
```

AgentCap itself should not parse shell syntax.

---

# 5. Execution Abstraction

Create a reusable internal execution layer.

A reasonable structure is:

```go
type Command struct {
    Args []string
}

type ExecutionResult struct {
    Args      []string
    ExitCode  int
    Duration  time.Duration
}
```

The exact API may differ.

Avoid prematurely storing all stdout/stderr inside `ExecutionResult` if streaming is a better design.

The architecture should make it possible for future phases to introduce captured output without redesigning the executor.

A possible package structure:

```text
cmd/
  acap/
    main.go

internal/
  command/
  executor/
```

Keep package boundaries simple.

---

# 6. Stdout and Stderr

Phase 0 should preserve stdout and stderr faithfully.

The preferred behavior is:

```text
child stdout -> acap stdout
child stderr -> acap stderr
```

Do not merge them.

Do not prefix lines.

Do not strip ANSI sequences.

Do not normalize whitespace.

Do not truncate output.

Do not reorder output intentionally.

AgentCap should act like a thin execution wrapper.

Later phases will introduce interception and reduction.

---

# 7. Exit Code Preservation

This is one of the most important Phase 0 requirements.

If:

```bash
some-command
```

exits with code:

```text
N
```

then:

```bash
acap run some-command
```

must normally exit with the same code:

```text
N
```

Examples:

```bash
acap run true
```

must exit `0`.

```bash
acap run false
```

must exit non-zero.

For:

```bash
acap run sh -c 'exit 42'
```

AgentCap should exit with code `42` where the operating system allows that status to be represented directly.

Do not convert all child failures into:

```text
exit 1
```

---

# 8. AgentCap Internal Errors

Distinguish between:

1. child command failure
2. AgentCap execution failure

For example:

```bash
acap run go test ./...
```

returning exit code `1` because tests failed is a normal child result.

But:

```bash
acap run definitely-not-a-real-command
```

is an AgentCap/process-launch failure.

Handle these cases clearly.

Print internal errors to stderr.

Keep error messages concise.

Do not print Go stack traces during normal CLI failures.

---

# 9. Command Not Found

Handle missing executables cleanly.

Example:

```bash
acap run this-command-does-not-exist
```

Expected properties:

- no panic
- useful error on stderr
- non-zero exit status
- behavior documented and tested

Do not silently fall back to a shell.

---

# 10. Working Directory

The child process must inherit AgentCap's current working directory.

Example:

```bash
cd /repo
acap run pwd
```

should behave like:

```bash
cd /repo
pwd
```

Do not introduce AgentCap-specific working-directory behavior in Phase 0.

---

# 11. Environment Variables

The child process must inherit the current environment.

Example:

```bash
FOO=bar acap run sh -c 'echo "$FOO"'
```

should expose `FOO=bar` to the child.

Do not sanitize or rewrite the environment unless required for basic correctness.

Do not add large sets of AgentCap-specific environment variables.

---

# 12. stdin

Child processes should receive stdin appropriately.

This is important for commands such as:

```bash
echo hello | acap run cat
```

where practical.

The preferred initial behavior is:

```text
acap stdin -> child stdin
```

Do not buffer stdin unnecessarily.

Interactive command behavior should remain usable where possible.

---

# 13. Signals

Handle termination signals correctly.

At minimum consider:

- SIGINT
- SIGTERM

On Unix-like systems, when AgentCap receives an interrupt, the child process should not be left running unintentionally.

For example, pressing Ctrl+C during:

```bash
acap run long-running-command
```

should terminate the command in a predictable way.

Avoid orphaning child processes.

Do not over-engineer cross-platform process-group handling in Phase 0, but implement sane behavior for the primary supported platform.

Document any platform limitations.

---

# 14. Execution Duration

Measure child execution duration internally.

For example:

```go
start := time.Now()
...
duration := time.Since(start)
```

Do not print timing information during normal execution.

It should be available internally for future phases and tests.

Optional debug output may expose it.

---

# 15. Large Output

Phase 0 must safely support commands that emit large amounts of output.

Since Phase 0 forwards output directly, do not read the complete output into a single in-memory byte slice unless there is a strong reason.

Prefer streaming using:

```text
child stdout -> parent stdout
child stderr -> parent stderr
```

This avoids excessive memory usage.

The execution architecture should make it possible for Phase 1 or Phase 2 to add interception/capture later.

Do not implement spill-to-disk result storage yet unless required by the chosen internal architecture.

---

# 16. Streaming Requirement

Avoid an architecture equivalent to:

```go
output, err := cmd.CombinedOutput()
```

for the main execution path.

Reasons:

- stdout and stderr become merged
- output must be fully buffered
- large output consumes excessive memory
- streaming behavior is lost

Prefer connecting streams directly.

For example:

```go
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
cmd.Stdin = os.Stdin
```

or an abstraction that preserves equivalent behavior.

---

# 17. stdout / stderr Ordering

Do not promise perfect global ordering between stdout and stderr.

They are separate streams and operating-system scheduling may interleave them differently.

The important requirement is:

- preserve stdout as stdout
- preserve stderr as stderr
- do not intentionally reorder either stream internally

Document this if necessary.

---

# 18. TTY Behavior

Investigate whether direct stream attachment is sufficient for common interactive commands.

Do not implement a PTY subsystem in Phase 0 unless it is genuinely necessary.

Commands requiring a real terminal may behave differently when wrapped.

That limitation is acceptable for Phase 0 if documented.

The Phase 0 goal is reliable non-interactive coding-agent command execution first.

---

# 19. Platform Scope

Primary target:

```text
Linux / Unix-like environments
```

The code should remain reasonably portable to macOS.

Do not spend significant implementation effort on full Windows parity during Phase 0 unless it comes naturally.

Avoid unnecessary Unix-only assumptions in general code.

If Unix-specific signal handling is required, isolate it behind platform-specific files where appropriate.

For example:

```text
signal_unix.go
signal_windows.go
```

Only add this split if needed.

---

# 20. Debug Mode

Implement a minimal opt-in debug mode if it is useful for development.

For example:

```bash
ACAP_DEBUG=1 acap run ...
```

Debug information may include:

```text
argv
resolved executable
exit code
duration
```

All debug output must go to stderr.

Normal execution must remain clean.

Do not add timestamps, logging frameworks, structured JSON logs, or log files unless needed.

---

# 21. No Output Transformation

This is a strict Phase 0 rule.

Do NOT:

- strip ANSI escape codes
- remove duplicate lines
- summarize output
- truncate output
- group lines
- detect filenames
- parse diagnostics
- parse search results
- rewrite progress output
- estimate tokens
- classify commands

Phase 0 must provide a trustworthy baseline against which future transformed behavior can be compared.

---

# 22. No Persistent State

Do not create:

```text
~/.cache/agentcap
```

or any persistent result database in Phase 0 unless absolutely required.

Do not save:

- command history
- stdout
- stderr
- result IDs
- sessions
- statistics

Persistent state belongs to later phases.

---

# 23. Testing

Add unit and integration tests.

At minimum test the following.

## Successful execution

Verify that:

```bash
acap run <successful-command>
```

returns exit code `0`.

## Non-zero exit

Verify that a child exit code is propagated.

For example, using a test helper or:

```bash
sh -c 'exit 42'
```

on supported Unix environments.

## Stdout forwarding

Child stdout must reach AgentCap stdout.

## Stderr forwarding

Child stderr must reach AgentCap stderr.

Test them independently.

## argv preservation

Verify arguments containing:

- spaces
- quotes as literal argument content
- Unicode
- empty strings where practical

are passed correctly.

## Environment inheritance

Verify a child can read an environment variable inherited through AgentCap.

## Working directory inheritance

Verify the child observes the expected current directory.

## stdin forwarding

Verify input piped to AgentCap reaches the child.

## Command not found

Verify:

- no panic
- useful stderr error
- non-zero exit status

## Large output

Generate a reasonably large stream and verify AgentCap handles it without truncating or deadlocking.

Do not depend on large external repositories for tests.

---

# 24. Test Helpers

Prefer deterministic test helpers over relying on arbitrary system commands.

One good approach is to use the Go test binary itself as a helper subprocess.

Alternatively, use small cross-platform helper programs under testdata where appropriate.

External commands such as:

```text
sh
printf
cat
```

may be used for Unix-specific integration tests, but the core test suite should remain as portable as practical.

---

# 25. Avoid Deadlocks

Be careful when dealing with:

- stdin
- stdout
- stderr
- process waiting
- pipes

Do not implement custom pipe-reading goroutines unless necessary.

If custom stream handling is added, ensure both stdout and stderr can be drained concurrently.

A child producing large stderr output must not block because only stdout is being read.

Prefer direct file descriptor/stream attachment for Phase 0.

---

# 26. Context Support

The executor API should accept a Go `context.Context`.

For example:

```go
func Run(ctx context.Context, command Command) Result
```

or equivalent.

This enables future:

- cancellation
- timeouts
- agent request cancellation

Do not add a default command timeout in Phase 0.

Commands should run until:

- completion
- user cancellation
- external signal

unless explicitly configured otherwise in future phases.

---

# 27. Data Model

Keep internal models minimal.

Avoid adding Phase 1 concepts such as:

```text
ReducedResult
Capsule
Reducer
CommandClassifier
TokenStats
```

Phase 0 may define something similar to:

```go
type Result struct {
    ExitCode int
    Duration time.Duration
}
```

Only add fields that Phase 0 actually uses or clearly needs for the execution boundary.

---

# 28. Suggested Project Layout

Keep the repository small.

For example:

```text
agentcap/
├── cmd/
│   └── acap/
│       └── main.go
├── internal/
│   └── executor/
│       ├── executor.go
│       └── executor_test.go
├── go.mod
├── go.sum
├── README.md
└── LICENSE
```

Additional files are fine if justified.

Do not create deep package hierarchies for hypothetical future features.

---

# 29. README

Add a minimal README describing:

- what AgentCap is
- that Phase 0 currently provides an execution wrapper only
- build instructions
- basic usage

Example:

```bash
go build ./cmd/acap

./acap run ls -la
./acap run go test ./...
```

Clearly state that output compression is not implemented yet.

---

# 30. Build Requirements

The project must build with standard Go tooling:

```bash
go build ./...
```

Tests:

```bash
go test ./...
```

Formatting:

```bash
gofmt
```

Prefer no external runtime dependencies.

The resulting `acap` CLI should be distributable as a single Go executable.

---

# 31. Explicit Non-Goals

Do NOT implement any of the following in Phase 0:

- output compression
- output truncation
- output summarization
- ANSI stripping
- progress-line removal
- command classification
- reducers
- `ls` special handling
- `find` special handling
- `grep` or `rg` special handling
- `cat` special handling
- token counting
- compression statistics
- persistent state
- command result storage
- result IDs
- `acap show`
- `acap raw`
- session management
- Git integration
- build/test parsing
- AST parsing
- Tree-sitter
- semantic source analysis
- MCP
- coding-agent-specific hooks
- Claude Code integration
- Codex integration
- remote services
- telemetry
- plugin architecture

Do not build abstractions solely for these future features.

Phase 0 should remain deliberately small.

---

# 32. Suggested Implementation Order

Implement in this order:

1. initialize Go module
2. create `acap` CLI entry point
3. parse `run` command and argv
4. implement direct process execution
5. connect stdin/stdout/stderr
6. preserve child exit code
7. handle command-not-found errors
8. add execution timing
9. add context cancellation
10. add signal handling
11. add integration tests
12. add large-output tests
13. add minimal debug mode if useful
14. add README
15. run formatting, tests, and vet/static checks where appropriate

Do not start Phase 1 work during this implementation.

---

# 33. Definition of Done

Phase 0 is complete when all of the following are true.

1. AgentCap builds as a single Go executable.

2. This works:

```bash
acap run <command> [args...]
```

3. argv is passed directly without implicit shell interpretation.

4. stdin is forwarded to the child.

5. child stdout is forwarded to AgentCap stdout.

6. child stderr is forwarded to AgentCap stderr.

7. stdout and stderr are not merged intentionally.

8. the child's exit code is preserved.

9. command-not-found errors are handled cleanly.

10. the child inherits cwd.

11. the child inherits environment variables.

12. Ctrl+C / termination does not normally leave the child orphaned.

13. large output can pass through without being fully buffered in memory.

14. execution duration is available internally.

15. context cancellation is supported by the executor.

16. automated tests cover the primary execution semantics.

17. `go test ./...` passes.

18. `go build ./...` passes.

19. normal command output is not transformed in any way.

20. no Phase 1 or later features have been implemented.

---

# 34. Engineering Priorities

When implementation tradeoffs are required, use this order:

1. Command semantic correctness
2. Exit-code correctness
3. Process lifecycle correctness
4. Stream correctness
5. Simplicity
6. Low overhead
7. Portability
8. Future extensibility

Do not sacrifice execution correctness for abstractions intended for later phases.

---

# 35. Final Verification

Before finishing Phase 0, manually verify representative cases such as:

```bash
acap run echo hello

acap run sh -c 'echo stdout; echo stderr >&2'

acap run sh -c 'exit 7'

printf 'hello\n' | acap run cat

FOO=bar acap run sh -c 'printf "%s\n" "$FOO"'

acap run command-that-does-not-exist
```

Also verify interrupting a long-running process with Ctrl+C.

Run:

```bash
gofmt
go test ./...
go build ./...
```

Then report:

- final project structure
- implemented execution behavior
- test coverage
- platform-specific limitations
- known issues

Do not proceed into Phase 1.

Stop once the Phase 0 Definition of Done is satisfied.
