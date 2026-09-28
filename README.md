# AgentCap

AgentCap (`acap`) is a command execution wrapper for AI coding agents.

**Phase 0** — execution core only. Output compression is not yet implemented.

## Build

```bash
go build ./cmd/acap
```

## Usage

```bash
acap run <command> [args...]

acap run ls -la
acap run go test ./...
acap run sh -c 'echo foo | grep foo'
```

Arguments are passed directly to the child process — no shell interpretation.

## Benchmark

```bash
acap bench command -- git diff
acap bench command --json -- go test ./...
```

Runs the command **exactly once** through the normal AgentCap pipeline
(stateless, no session delta) and reports:

- `raw` — bytes captured from the command's stdout and stderr (total = stdout + stderr)
- `visible` — bytes of the AgentCap result presentation the agent receives (`@acap <id>` header, capsule, folded-in stderr)
- `reduction` — `1 - visible/raw` for this command; `n/a` when raw is 0; can be negative for tiny outputs
- `duration` — command execution time (process start to exit)
- `processing` — AgentCap post-execution time (reduction, hashing, storage, stats), with reducer-only time in parentheses
- `id` — the stored result, inspectable with `acap show` / `acap raw`

Reduction is single-command agent-visible output reduction, not workflow or
end-to-end agent savings. `acap bench` exits with the command's exit code.

### Stateful benchmark sessions

```bash
session=$(acap bench session start)
acap bench command --session "$session" -- git diff
acap bench command --session "$session" -- go test ./...
acap bench command --session "$session" -- git diff
acap bench session show "$session"
acap bench session show --json "$session"
```

A benchmark session persists across CLI invocations and maps one-to-one to a
fresh normal AgentCap session, so its delta baseline is isolated from regular
commands and other benchmarks. Each target command executes exactly once. The
same captured result supplies all three byte counts:

- `raw` = captured stdout bytes + captured stderr bytes
- `stateless` = the real full reduced presentation, including its result ID,
  before session comparison
- `stateful` = the actual full, delta, or unchanged presentation after the
  normal AgentCap session/delta engine

The reported reductions are `1 - stateless/raw`, `1 - stateful/raw`, and
`1 - stateful/stateless`; a zero denominator is shown as `n/a`. Processing time
keeps the single-command definition: child exit until the stateful response is
ready, including reduction, hashing, storage, delta comparison, session history,
and stats. Every result remains complete and inspectable with `acap show` and
`acap raw`.

This deterministic benchmark measures AgentCap output bytes and processing
time. It does not measure tokens, task success, model quality, or coding-agent
accuracy.

## Debug

```bash
ACAP_DEBUG=1 acap run go test ./...
```

Prints argv, exit code, and duration to stderr.

## Test

```bash
go test ./...
```

## Notes

- stdout and stderr are forwarded separately; relative ordering between them is not guaranteed (OS scheduling).
- TTY-dependent commands may behave differently when wrapped (no PTY in Phase 0).
- Signal forwarding is implemented for Unix (SIGINT, SIGTERM) and Windows (os.Interrupt).
