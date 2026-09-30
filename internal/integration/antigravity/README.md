# Antigravity run_command adapter

Status: **working, restricted to sessions that declare unrestricted command
permissions.** In every other session the hook is a verified no-op, and
Antigravity's normal behavior is unchanged.

One adapter serves Antigravity 2.0, the CLI and the IDE. They share the
`hooks.json` contract. Verified with **Antigravity CLI 1.2.14** (Linux/WSL2,
2026-10-01). The official contract ships with the product at
`~/.gemini/antigravity-cli/builtin/skills/agy-customizations/docs/hooks.md`.

## Why PreToolUse, and why PostToolUse alone is insufficient

`PostToolUse` runs after execution, and its output contract is `{}`. It cannot
replace what the agent sees, and rerunning the command there would break
exactly-once execution. `PreToolUse` supports `overwrite`, a shallow merge into
the tool arguments before execution. The adapter overwrites only `CommandLine`,
changing it to:

```
acap antigravity-exec '<base64url JSON payload>'
```

Antigravity then executes the wrapper in the original `Cwd`, with its own shell,
PTY, environment and async policy, as the real `run_command` step. `Cwd` and
`WaitMsBeforeAsync` are left untouched. The agent sees a one-line notice
("A pre-tool hook changed the arguments…") followed by the AgentCap capsule and
the real exit code. No PostToolUse hook is installed.

## Feasibility findings (real CLI, headless `agy -p`)

| Question | Finding |
| --- | --- |
| Payload | `toolCall.args` has `CommandLine`, `Cwd`, `WaitMsBeforeAsync`, `toolAction`, `toolSummary`; common fields as documented. The hook cwd is the `.agents` directory. |
| `{}` response | **Denies** the call (`tool call denied by pre-tool hook`). `decision` is mandatory. |
| `{"decision":"ask"}` | Same as having no hook: allow rules and grants apply, unapproved commands are denied or prompted, safe commands run. |
| Permission target with `overwrite` | Rules are evaluated against the **overwritten** CommandLine. Allow `touch` (original) plus rewrite to `mkdir` ran nothing; allow `mkdir` ran the rewrite. |
| Permission mode visible to hook? | **No.** Payload and environment (`ANTIGRAVITY_CONVERSATION_ID` only) are identical with and without `--dangerously-skip-permissions`. |
| Exit status / output | The model receives `The command exited with code N` plus the capsule. |
| Async promotion | A wrapper exceeding `WaitMsBeforeAsync` is backgrounded normally. The capsule and exit code arrive on completion, and the command runs once. |
| Cancellation | Antigravity makes each command a process-group leader and SIGKILLs the group. See "Cancellation" below. |
| Missing `acap` on PATH | The hook exits 127, and Antigravity fails the step before execution with `acap: not found`. The command does not run. |

## Permission preservation

Permissions are checked against the overwritten command, and the hook cannot see
the autonomy mode. Unconditional rewriting would therefore hide the real command
from allow/deny rules and approval prompts. So:

- The adapter **always** answers `decision: "ask"`. It never returns `allow`,
  `deny` or `permissionOverrides`.
- It rewrites **only** when the launching environment sets
  `ACAP_ANTIGRAVITY_PERMISSIONS=unrestricted`. This is the caller's declaration
  that the session already runs without command-level permission checks
  (always-proceed / `--dangerously-skip-permissions`), for example a benchmark
  sandbox. Hooks inherit Antigravity's environment.
- Without that declaration, every call passes through unchanged (reason
  `permission-mode`). That includes `request-review`, `strict`, and unknown
  modes.

In declared mode, command-level allow/deny rules would see
`acap antigravity-exec …`. Do not set the variable in sessions that rely on
command rules or interactive approval. Full transparent support in normal modes
needs Antigravity to evaluate permissions on the pre-overwrite command, or to
expose the mode to hooks.

## Normalization and session mapping

`ParseHookInput` reads only `conversationId`, `workspacePaths`, `modelName`,
`stepIdx` and `toolCall`. Unknown fields are ignored. A missing
`toolCall.name`, or a `run_command` without `CommandLine`, produces a stderr
diagnostic and a pass-through. Transcript and artifact paths are never read.

`protocol.ToolRequest` carries: `ShellCommand` = `CommandLine` (opaque shell
source, run with `bash -c` or `sh -c`), `WorkingDir` = `Cwd`, and
`Integration{Agent/Adapter: "antigravity", ExternalSession: conversationId,
ToolUseID: stepIdx, Model: modelName}`. The core's `common.MapSession` hashes
(agent, conversationId, project root). The same conversation therefore gets one
session, and different conversations or projects are isolated. `modelName` and
`stepIdx` are attribution only. The project root comes from `Cwd` via
`project.FindRoot`, never from `workspacePaths`. Results live in
`<root>/.acap/`.

## Exactly-once, exit status, fail-open

`Prepare` never executes anything. `ExecutePayload` calls `engine.Run` once, and
its exit code is the command's. After a reduction or storage failure the wrapper
prints the captured raw streams and a one-line `acap:` diagnostic, and keeps the
real exit status. It never retries. A rejected payload exits 2 with
"command not started". Raw stdout and stderr are stored separately and in full,
and `acap show` / `acap raw` work unchanged.

## Bypass and conservative scope

Pass-through (`{"decision":"ask"}`, no overwrite) with a recorded reason:

- `explicit-bypass`: `ACAP_BYPASS=1` or `ACAP_BENCH_MODE=disabled`. This is the
  benchmark OFF condition and needs no hook reinstall.
- `recursion`: `ACAP_INTERCEPT_DEPTH>0`, or any `acap` word in the command.
  Quoted occurrences are intentionally over-matched.
- `async`: `WaitMsBeforeAsync < 2000`. The agent expects early backgrounding,
  and the wrapper emits output only on completion.
- `long-running-or-interactive`: `tail -f`, `watch`, dev servers, editors,
  pagers, `ssh`, `--watch`, …
- `missing-cwd`, `unsupported-tool`, `invalid-input`, `binary-unavailable`,
  `shell-unavailable`.

Known differences in declared mode: the child's stdout/stderr are pipes, not
Antigravity's PTY, so TTY-dependent formatting can differ. Output of
backgrounded commands appears only when they complete.

## Cancellation

Antigravity SIGKILLs the command's process group, and SIGKILL cannot be
forwarded. The adapter therefore sets `ToolRequest.InheritProcessGroup`, which
keeps the child in the wrapper's group so the kill reaches the whole tree
(verified: no orphan after `--print-timeout`). Independently, `internal/exec`
now forwards SIGINT, SIGTERM and SIGHUP to the child's whole process group until
it exits.

## Install / remove

```
acap integrate antigravity [--dry-run]
acap integrate antigravity --remove
```

The installer writes one named hook, `"agentcap"`, to `.agents/hooks.json`:
`PreToolUse` with matcher `run_command` and command `acap hook antigravity`,
timeout 10 s. Other named hooks are preserved byte-for-byte in meaning.
Reinstalling is a no-op, and a stale `agentcap` entry is refreshed. Removal
deletes only that key, and existing + install + remove = existing. Malformed
files are rejected without writing. Writes go through a temp file and rename,
keeping the file mode. `acap` must be on the PATH Antigravity gives hooks. When
PATH resolves to the hook binary, the overwrite uses a bare `acap` (readable in
the UI and matchable by `command(acap)` rules). Otherwise it uses the quoted
absolute path.

## Metrics

`.acap/adapter-metrics.jsonl` gets one record per hook decision or wrapper run,
with: `agent`/`adapter` = `antigravity`, `adapter_version`, the hashed session
mapping, `tool_use_id` (stepIdx), `reason`, `commands_intercepted`,
`commands_bypassed`, `adapter_failures`, `adapter_latency_ns` (hook preparation)
and `processing_latency_ns` (core). Command text and conversation IDs are
excluded. Stored results also carry `integration` attribution, including
`model`.

## Verification

`go test ./internal/integration/antigravity/` covers:

- the shared conformance suite: exit codes 0/1/2/7/127, streams, exactly-once,
  cwd/env, session/project isolation, fail-open, concurrency, cancellation
- parser cases, and pass-through for every permission and bypass reason
- shell-semantics round trips and injection safety
- running the overwritten CommandLine under `sh` (counter increments once,
  exit 7 preserved, `acap raw` works)
- Go test failure compaction and deltas, and model-change session stability
- 40k-line mixed output, group-kill cancellation, installer preservation,
  idempotency and removal

Real end-to-end run (CLI 1.2.14, declared mode, fixture Go repo): grep, then a
compact `go-test FAIL`, then `acap show` / `acap raw` drill-down (bypassed as
recursion), then edits, then `delta … failed 2 -> 1`, then `go-test PASS`, all
in one mapped session. OFF (`ACAP_BYPASS=1`) returned native raw output with no
rewrite notice.

Not verified here: the Antigravity 2.0 and IDE surfaces, interactive approval
prompts, Windows, and deny-rule behavior under `--dangerously-skip-permissions`
(if deny rules still apply in that mode, declared mode would bypass them).
