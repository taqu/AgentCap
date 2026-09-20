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
