# Claude Code shell adapter

Status: **restricted integration; Phase 6.5B is not complete for normal permission modes.**

The public hook API cannot both preserve normal command permission rules and
replace failed Bash results. This adapter compresses foreground Bash calls only
when the caller has already explicitly selected `bypassPermissions`. It never
selects that mode or returns a permission `allow` decision. In `default`, `auto`,
`acceptEdits`, `plan`, `dontAsk`, missing, and unknown modes, the hook returns `{}`
and ordinary Claude Bash execution continues. Do not disable permissions just
to enable compression in a trusted working environment. A full transparent
adapter needs a supported failed-result replacement API or another execution
mechanism that preserves permission decisions on the original command.

## Feasibility evidence

Verified with the installed Claude Code **2.1.285** on Windows/Git Bash on
2026-09-30. Official contract: <https://code.claude.com/docs/en/hooks>.

| Probe | Hook event | What Claude received |
| --- | --- | --- |
| `printf 'hello\n'` | PostToolUse | `PROBE_REPLACEMENT`, supplied through `updatedToolOutput` |
| `sh -c 'echo failure-output; exit 7'` | PostToolUseFailure | Original `Exit code 7` and `failure-output` |
| Known failing `go test ./...` fixture | PostToolUseFailure | Original `Exit code 1` and test diagnostics |

`PostToolUseFailure` exposes display text that can already be truncated and
interleaves stdout/stderr. It supports additional context, not replacement.
`PreToolUse.updatedInput` changes the input used for permission checks. A
post-execution success-only adapter and unrestricted pre-execution wrapping are
therefore unsuitable. The old execute-inside-PreToolUse-and-block implementation
has been removed because it ran before native tool permission checks.

The executable probe source is `testdata/probe/main.go`. Build it using
`go build -o <probe> ./internal/integration/claude/testdata/probe`; configure a
disposable project's PreToolUse, PostToolUse, and PostToolUseFailure Bash hooks
to run `<probe> <absolute-log-path>`. The probe records JSON stdin and replaces
only successful output. Run the three commands above in a real Claude session.
Use a disposable fixture: logs can contain sensitive input and output.

## Runtime boundary

`Prepare` parses the small Claude-specific input, checks the safety scope, and
returns `hookSpecificOutput.updatedInput`, preserving all other input fields.
It **never executes** the command. Source, cwd, timeout, session, and attribution
travel as a base64url JSON argument to the hidden `acap claude-exec` entry point.
Shell quoting protects the executable path; metadata never becomes shell syntax.

Claude launches the wrapper in its Bash execution context. `ExecutePayload`
calls `engine.Run` exactly once with `ToolRequest.ShellCommand`; the engine runs
the original source opaquely with Bash `-c`. Direct argv callers keep the original
`ToolRequest.Command` path. A conservative classification of plain commands
selects existing reducers; classification is never used to execute commands.
Complex syntax uses the existing generic reducer, with its existing reduction
limitations, including conservative stderr handling. No reducers are modified.

The wrapper returns the real command exit status and the core presentation for
both success and failure. Stored raw stdout and stderr remain separate and
complete; adapter executions do not use the legacy CLI 10 MiB capture cap.
Complete capture uses memory proportional to raw output. `acap show`/`raw`
read the normal project store. After storage/reduction failure the wrapper
returns captured raw streams and the command status, with a concise processing
diagnostic; it never reruns the command. A wrapper crash after execution cannot
be recovered by retrying and is not automatically retried.

Project root follows the existing `project.FindRoot` rules, independently of
execution cwd. Adapter session history is project-local under
`.acap/sessions/<hash>/history.jsonl`. The full SHA-256 identity includes agent,
external session ID, and project root. Missing external IDs are stateless.
Resumed sessions reuse the mapping. Subagents share the top-level `session_id`;
`agent_id` and `tool_use_id` are optional diagnostic attribution.

The common `ACAP_INTERCEPT_DEPTH` guard is incremented only in the child process
environment. `ACAP_BYPASS=1` preserves normal Bash behavior for OFF benchmarks;
`ACAP_BENCH_MODE=disabled` also bypasses. No persistent environment changes are
made. Inherited environment is preserved, and normalized per-process overrides
are applied by the common engine.

Explicit background calls are bypassed. Commands mentioning `cd`, shell state
mutation, or AgentCap are conservatively bypassed, including quoted occurrences.
This preserves Claude's persistent cwd and retrieval behavior. The wrapper uses
a child Bash rather than Claude's non-exported functions/aliases/options;
commands depending on those features are outside the supported scope.
PowerShell is not intercepted. Windows uses `CLAUDE_CODE_GIT_BASH_PATH` when set,
otherwise a discoverable Git Bash; the Windows WSL launcher is rejected.

The original timeout field is retained and also applied as a millisecond context
deadline to execution. Timeout/cancellation returns non-success (124) and a
concise interruption message. Context cancellation terminates the child process
tree using process groups on Unix and taskkill on Windows. Native Claude
background promotion, interactive shell behavior, and real user cancellation
have not been proven equivalent; they remain acceptance limitations.

## Configuration and diagnostics

Use `acap integrate claude [--dry-run]` and `acap integrate claude --remove`.
The installer merges one Bash PreToolUse hook into `.claude/settings.json`.
`acap` must be on the PATH available to Claude hooks. A missing hook executable
leaves normal native tool execution available through Claude's hook failure
behavior; a missing wrapper after rewriting fails before command execution.
No permission settings or `CLAUDE.md` instructions are installed.

Ownership is the exact command handler `acap hook claude`, never an `agentcap`
substring. Unknown fields on other matchers/handlers, other lifecycle events,
permissions, and environment settings survive install/remove. Malformed settings
are rejected without writing. Writes use a unique adjacent temporary file and
rename; repeated installation/removal makes no changes. Concurrent independent
settings editors are not coordinated; serialize installation with other editors.

Best-effort `.acap/adapter-metrics.jsonl` records agent, adapter version, mapped
session, tool/subagent IDs, bypass reason, intercepted/bypassed/failure counters,
adapter latency and core processing latency. Command text, external session IDs,
and environment values are excluded from this log. Intercepted result metadata
also carries attribution in the existing store. Counters are per-record (sum
them); metrics failures cannot break execution. Latency is pre-execution adapter
preparation time, distinct from execution and core processing time.

## Verification

Run `go test ./...`. `TestClaudeConformance` runs the reusable
`internal/integration/conformance.Run` suite through `Prepare` and the wrapper,
covering success/failure, exit statuses 0/1/2/7/127, raw streams, exactly-once
side effects, nested cwd, environment, session/project isolation, missing IDs,
fail-open, concurrency, and cancellation. Future adapters supply an `Invoke`
function to the same suite. Claude-specific tests cover permission-mode guards,
background/recursion/bypass, shell transport, malformed payloads, settings
ownership, unknown fields and configuration preservation.

A real restricted-mode Claude session has verified search, compact Go failure,
targeted `show`, editing, compact successful reruns, and success/failure markers
written once. Raw failure output retained all 500 diagnostic lines and `FIX_ME`.
The existing Go test reducer and `show --test` display only the first diagnostic
lines in this fixture; that existing drill-down limitation was not changed.

Unix execution is covered by portable code/tests but was not run in this Windows
session. Normal permission-mode transparent compression, real cancellation,
auto-background equivalence and version compatibility beyond 2.1.285 remain
unverified/unsupported. Do not label Phase 6.5B complete or run ON/OFF benchmarks
in normal permission modes: ON would correctly bypass too.
