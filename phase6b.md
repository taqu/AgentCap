# Phase 6.5B — Claude Code Adapter

## Objective

Implement the Claude Code adapter on top of the common adapter contract established in Phase 6.5A.

The adapter must allow normal Claude Code shell usage to benefit from AgentCap while preserving:

- command semantics
- Claude Code permission behavior as far as technically possible
- exactly-once command execution
- exit status
- cwd
- environment behavior
- cancellation
- AgentCap session state
- progressive disclosure
- fail-open behavior

The adapter must remain thin.

Do not move Claude-specific behavior into AgentCap core.

---

# 1. Required Architecture

The intended architecture is:

```text
Claude Code
    |
    | Bash tool
    v
Claude Code Adapter
    |
    v
Phase 6.5A Common Adapter Boundary
    |
    v
AgentCap Core
```

Claude-specific responsibilities belong only in the Claude adapter.

The AgentCap core must not gain logic such as:

```text
if agent == "claude-code" { ... }
```

The adapter should translate Claude Code hook data into the normalized integration structures introduced in Phase 6.5A.

---

# 2. Important Claude Code Hook Constraints

Before implementing the final integration path, verify the behavior of the installed/supported Claude Code version against the current official hook contract.

The implementation must account for the following Claude Code behavior:

```text
PreToolUse
    -> runs before Bash execution
    -> receives command input
    -> may replace tool input

PostToolUse
    -> runs after successful Bash execution
    -> receives the Bash result
    -> may replace the tool output sent to Claude

PostToolUseFailure
    -> runs after failed tool execution
    -> receives failure information
    -> does not provide the same general Bash output-replacement mechanism
```

This asymmetry is important.

AgentCap must not accidentally compress only successful commands while leaking full compiler/test failures into Claude's context.

Build and test failure output is one of AgentCap's highest-value reduction cases.

---

# 3. Start With an Integration Feasibility Test

Before writing substantial adapter infrastructure, implement a minimal executable integration test/spike covering these cases:

```bash
printf 'hello\n'
```

```bash
sh -c 'echo failure-output; exit 7'
```

and one realistic failure such as:

```bash
go test ./...
```

with a known failing test fixture.

Verify exactly what Claude Code exposes to hooks and exactly what Claude receives after the hook response.

The spike must answer:

```text
1. Can successful command output be replaced safely?

2. Can failed command output be reduced before Claude sees it?

3. Can AgentCap participate without executing the original command twice?

4. Can the original command's exit status remain visible?

5. Can shell syntax remain unchanged?

6. Does integration change Claude Code's permission decision behavior?

7. Do timeout and background execution semantics remain correct?
```

Do not proceed with an architecture that cannot satisfy the critical invariants.

---

# 4. Do Not Accept Success-Only Integration

The following architecture is insufficient as the final implementation:

```text
Bash succeeds
    -> PostToolUse
    -> AgentCap compression

Bash fails
    -> raw Claude Code failure output
```

This would systematically bypass AgentCap for exactly the workflows where large diagnostic output is common:

```text
compiler errors
test failures
linker failures
cargo check failures
go test failures
large build failures
```

Do not mark Phase 6.5B complete with such behavior.

A success-only PostToolUse implementation may be used temporarily for investigation, but not as the final adapter.

---

# 5. Preferred Integration Property

The final adapter must satisfy:

```text
Claude requests shell command
        |
        v
command executes exactly once
        |
        v
AgentCap captures/processes result
        |
        v
Claude receives AgentCap presentation
```

for both:

```text
exit = 0
```

and:

```text
exit != 0
```

The underlying command must never execute once through Claude Code and then again through AgentCap.

---

# 6. Evaluate the Existing Phase 6 / Phase 6.5A Mechanism First

Before introducing a new execution strategy:

1. inspect the existing Phase 6 integration implementation
2. inspect the Phase 6.5A common adapter boundary
3. identify whether a completed-result ingestion path already exists
4. identify whether an execution-wrapping path already exists
5. reuse the existing mechanism when it satisfies the invariants

Do not create a second AgentCap integration architecture specifically for Claude Code.

If Phase 6.5A already defines an appropriate abstraction, use it directly.

---

# 7. Completed-Result Ingestion Path

If the common adapter boundary supports ingestion of an already executed command result, evaluate whether Claude Code hooks expose sufficient information to use that path.

Conceptually:

```text
Claude Bash executes command
        |
        v
Claude hook receives completed result
        |
        v
Claude Adapter
        |
        v
AgentCap ingest/capture/reduce/store
        |
        v
replacement tool output
```

This has an important advantage:

```text
Claude Code itself remains responsible for command execution
```

and therefore naturally preserves much of Claude Code's:

- sandbox behavior
- permission handling
- shell behavior
- timeout handling

However, this approach is acceptable only if both successful and failed shell executions can be handled correctly.

If failed Bash output cannot be replaced before it reaches Claude, this cannot be the sole integration path.

Do not silently accept that limitation.

---

# 8. Pre-Execution Wrapping Path

If completed-result ingestion cannot cover failed Bash commands, evaluate a PreToolUse wrapping strategy using the Phase 6.5A execution boundary.

Conceptually:

```text
Claude proposes:

    go test ./...

PreToolUse adapter transforms execution path:

    AgentCap integration wrapper
        -> original shell command
```

The original command must still execute exactly once.

The adapter must preserve the original command as opaque shell syntax.

Do not parse it using naïve whitespace splitting.

---

# 9. Permission Semantics Are a Critical Risk

Claude Code evaluates permissions against tool input.

If PreToolUse rewrites:

```text
go test ./...
```

into something resembling:

```text
acap ...
```

the adapter may accidentally change what Claude Code's permission system evaluates.

This is a correctness and security issue.

Therefore explicitly test:

```text
ordinary allowed command
command requiring confirmation
command covered by an allow rule
command covered by a deny rule
destructive command
compound shell command
pipeline
redirection
subshell
```

Examples:

```bash
git status
```

```bash
rm temporary-file
```

```bash
echo foo | grep foo
```

```bash
cd subdir && go test ./...
```

```bash
FOO=bar ./script.sh
```

The adapter must not make a command effectively less restricted merely because it is wrapped by AgentCap.

Do not auto-return a broad:

```text
permissionDecision = allow
```

for intercepted Bash calls.

AgentCap is not a permission bypass.

---

# 10. If Permission Semantics Cannot Be Preserved

If testing demonstrates that pre-execution wrapping materially bypasses or weakens Claude Code's command permission semantics, do not hide the issue.

Prefer one of these outcomes:

```text
A. use a safer supported integration mechanism

B. restrict transparent interception to cases whose permission behavior
   can be proven equivalent

C. expose the unsupported limitation explicitly
```

Do not trade permission correctness for compression.

Do not introduce undocumented hacks into Claude Code internals.

---

# 11. Claude Hook Input Parser

Implement a small Claude-specific parser for the hook JSON.

The adapter should extract only fields it actually needs.

Relevant fields include conceptually:

```text
session_id
cwd
hook_event_name
tool_name
tool_input
tool_use_id
duration_ms when available
```

For Bash:

```text
tool_input.command
tool_input.timeout
tool_input.run_in_background
```

Do not deserialize the entire Claude hook schema into a large permanent compatibility API unless necessary.

Unknown fields should normally be ignored safely.

---

# 12. Tool Scope

Phase 6.5B should initially target:

```text
Bash
```

Do not intercept unrelated Claude tools such as:

```text
Read
Write
Edit
Glob
Grep
WebFetch
WebSearch
Agent
```

unless required by the existing Phase 6 design.

AgentCap's responsibility is command-result mediation.

Do not turn the Claude adapter into a general Claude Code tool proxy.

---

# 13. PowerShell

Do not automatically expand the scope to PowerShell unless AgentCap already has a tested shell-execution model that supports it correctly.

If PowerShell support already exists and requires little additional adapter logic, it may be included.

Otherwise:

```text
Bash first
PowerShell later
```

is acceptable.

Document the scope explicitly.

Do not claim PowerShell compatibility without tests.

---

# 14. Shell Command Preservation

Claude Code's Bash tool supplies the command as a shell command string.

Preserve it as a shell command string.

For example:

```bash
cd src && make test 2>&1 | tee build.log
```

must not become:

```text
["cd", "src", "&&", "make", ...]
```

through naïve splitting.

Likewise preserve:

```bash
FOO="a b" ./script.sh
```

```bash
for f in *.go; do echo "$f"; done
```

```bash
git diff -- '*.go'
```

AgentCap must preserve shell semantics exactly.

---

# 15. Session Mapping

Use Claude Code's supplied session identifier as the external adapter session identity.

Conceptually:

```text
agent             = claude-code
external_session  = Claude session_id
project           = AgentCap project root
```

Map this through the Phase 6.5A session contract.

Required:

```text
same Claude session
    -> same AgentCap session

different Claude session
    -> different AgentCap session
```

A resumed Claude session should continue using the same AgentCap session when Claude preserves the same session identity.

Do not derive session identity from:

```text
PID
cwd alone
transcript filename alone
current timestamp
```

when Claude's explicit session identifier is available.

---

# 16. Subagents

Claude hooks may also contain subagent identity information.

Do not automatically create separate AgentCap sessions for every subagent unless measurement or existing Phase 6.5A semantics require it.

Define the behavior deliberately.

A reasonable initial model is:

```text
Claude top-level session_id
    -> AgentCap session
```

with optional subagent identity retained as metadata.

If concurrent subagent commands can corrupt stateful baseline selection, isolate them using the common contract rather than adding ad hoc Claude logic.

Whichever behavior is chosen, add tests.

---

# 17. Working Directory

Use the cwd supplied by Claude Code for the actual tool invocation.

Do not assume:

```text
cwd == CLAUDE_PROJECT_DIR
```

Claude may execute commands after changing directories.

AgentCap must retain the distinction:

```text
project root
!=
execution cwd
```

Test:

```bash
cd nested/path
```

followed by:

```bash
pwd
```

and commands invoked from nested directories.

---

# 18. Project Root Resolution

Do not let Claude-specific project metadata replace AgentCap's existing root resolution.

Continue using AgentCap's normal project-root rules.

Claude's cwd is execution context.

AgentCap's project root remains determined by AgentCap.

This preserves:

```text
nearest .acap/
nearest Git root
cwd fallback
```

according to existing core behavior.

---

# 19. Environment

Do not reconstruct Claude Code's environment from hook JSON.

Hook processes inherit their environment from Claude Code subject to Claude Code's own subprocess behavior.

Use the existing Phase 6.5A environment mechanism.

Any adapter-specific environment variables must be:

- narrowly scoped
- documented
- collision-resistant
- removed or ignored outside adapter execution

Possible uses include:

```text
recursion protection
adapter identification
session propagation
bypass
```

Do not expose sensitive environment values in AgentCap output or logs.

---

# 20. Exactly-Once Verification

This phase must include a real exactly-once test through the Claude integration path.

Use a fixture command that performs an observable side effect.

For example:

```bash
./increment-counter.sh
```

Expected:

```text
before = 0
after  = 1
```

Never:

```text
after = 2
```

Also test a failing side-effect command:

```text
perform side effect
then exit 7
```

It must still execute exactly once.

---

# 21. Exit Status Preservation

Verify at least:

```text
0
1
2
7
127
```

where practical.

The agent must still receive the semantic fact that the command failed.

AgentCap compression must not transform:

```text
exit 7
```

into apparent success.

Likewise, adapter failure must not overwrite the original command's status.

---

# 22. Failure Output

Pay special attention to non-zero command exits.

Test realistic large failure output:

```text
Go test failure
C/C++ compiler failure
Rust cargo check/test failure
```

The result Claude receives should be the normal AgentCap reduced presentation, including a usable result ID for drill-down.

For example, conceptually:

```text
@acap abc123 go-test FAIL
failed=2

TestParserInvalidInput
TestResolveMissingSymbol
```

rather than the complete raw log.

Claude must still be able to request:

```bash
acap show abc123 --test TestParserInvalidInput
```

or:

```bash
acap raw abc123
```

---

# 23. Success Output

Verify successful noisy commands as well:

```bash
go test ./...
```

```bash
cargo test
```

```bash
git diff
```

```bash
rg pattern .
```

Claude should receive AgentCap's normal presentation, not a Claude-specific summary format.

Do not duplicate reducers in the adapter.

---

# 24. Stateful Compression

Repeated commands within one Claude session must reuse the AgentCap session established from Claude's session ID.

Test:

```text
run #1 -> full compact result

edit

run #2 -> delta result

run #3 with no relevant change -> unchanged result
```

The Claude adapter must not interfere with:

```text
baseline selection
delta generation
unchanged detection
```

These remain core responsibilities.

---

# 25. Progressive Disclosure

Claude must receive normal AgentCap result IDs.

Existing commands must continue to work naturally from Claude Code:

```bash
acap show <id>
acap raw <id>
acap show <id> --errors
acap show <id> --warnings
acap show <id> --test ...
acap show <id> --file ...
```

Do not intercept these recursively.

They are part of AgentCap's intended model-facing API.

---

# 26. Recursion Protection

The Claude adapter must detect when a command is already inside AgentCap integration.

At minimum ensure commands such as:

```bash
acap show ...
acap raw ...
acap stats
acap clean
```

do not become:

```text
Claude hook
 -> AgentCap
 -> shell
 -> Claude hook equivalent
 -> AgentCap
 -> ...
```

Use the shared Phase 6.5A recursion mechanism.

Do not create a second Claude-only recursion scheme unless unavoidable.

---

# 27. Bypass

Support the common explicit bypass mechanism from Phase 6.5A.

The bypass is required for:

```text
benchmark AgentCap OFF runs
debugging
adapter diagnosis
compatibility testing
emergency recovery
```

The Claude adapter must not invent an unrelated bypass flag when a common one exists.

Test that bypass returns normal Claude Code Bash behavior.

---

# 28. Background Commands

Claude Bash supports background execution.

Do not assume synchronous behavior.

Test or explicitly exclude:

```text
run_in_background = true
```

If AgentCap cannot safely capture/reduce background command results through the current integration contract, bypass them conservatively.

Do not change:

```text
background
```

into:

```text
foreground
```

just to make compression easier.

Record the bypass reason where adapter metrics support it.

---

# 29. Timeouts

Preserve Claude Code's timeout semantics.

If `tool_input.timeout` is present:

- do not silently remove it
- do not replace it with AgentCap's unrelated default
- do not allow the wrapper to outlive the intended command timeout

Test at least one command that exceeds a short timeout.

The agent-visible result must remain semantically consistent with the command being interrupted/timed out.

---

# 30. Cancellation

Test cancellation through Claude Code when practical.

Required behavior:

```text
Claude cancellation
    ->
integration
    ->
underlying process terminates
```

Do not leave orphaned build/test processes.

If Claude's hook mechanism prevents perfect propagation in a particular integration architecture, document that constraint and choose the safest behavior.

---

# 31. Fail-Open

If AgentCap processing fails after the command has already executed, Claude should receive usable original command output whenever technically possible.

Conceptually:

```text
command execution succeeded
AgentCap reducer failed
    ->
return original usable result
```

Do not:

```text
command succeeded
AgentCap compression failed
    ->
fabricate command failure
```

For an architecture that wraps execution before Claude Code receives the result, ensure the wrapper itself implements equivalent fail-open behavior.

---

# 32. Hook Failure

A broken Claude hook must not silently consume command output.

Test:

```text
AgentCap binary missing
malformed adapter configuration
storage unavailable
internal reducer failure
```

Distinguish:

```text
adapter unavailable
```

from:

```text
original command failed
```

Keep diagnostics concise.

Avoid flooding Claude's model context with adapter debugging output.

---

# 33. Hook Configuration

Implement project-local Claude Code integration using the normal Claude hook configuration mechanism.

Prefer project-scoped configuration for reproducible benchmarks.

The installed configuration should target only the required Bash lifecycle events.

Do not match every Claude tool unless necessary.

Keep hook configuration minimal.

---

# 34. Installer

Provide an idempotent installation path.

Conceptually, support something equivalent to:

```text
acap integrate claude
```

or the naming convention already established by Phase 6.

Do not invent a parallel CLI hierarchy if an integration command already exists.

Installation should:

```text
detect existing Claude settings
preserve unrelated settings
add AgentCap hooks
avoid duplicate entries
write valid JSON
report what changed
```

---

# 35. Do Not Overwrite User Configuration

Claude Code project settings may already contain:

```text
permissions
hooks
environment configuration
other integrations
```

The installer must merge only the AgentCap-owned configuration.

Do not replace the whole settings file.

Do not remove unrelated hooks.

Do not reorder or rewrite large portions of user configuration unnecessarily.

---

# 36. Uninstaller

Provide safe integration removal.

Conceptually:

```text
acap integrate claude --remove
```

or the existing project convention.

Removal must delete only AgentCap-owned entries.

If the settings file contained unrelated configuration before installation, it must remain intact afterward.

Test:

```text
existing hooks
+
install AgentCap
+
remove AgentCap
=
existing hooks
```

---

# 37. Idempotency

Running installation twice must not produce duplicated hook entries.

Required:

```text
install
install
```

results in one logical AgentCap integration.

Likewise:

```text
remove
remove
```

should be safe.

---

# 38. Configuration Ownership

Make AgentCap-owned hook entries identifiable.

Use a stable command path, marker, or structure sufficient for safe update/removal.

Do not depend on fragile array indexes.

Do not delete another hook merely because it also matches Bash.

---

# 39. Version Compatibility

Do not assume every installed Claude Code release has identical hook capabilities.

At installation or diagnostic time, provide enough information to identify unsupported versions where necessary.

Avoid a large version-compatibility framework.

The goal is only:

```text
supported
unsupported
possibly degraded
```

with useful diagnostics.

Do not silently select a degraded success-only integration if it violates Phase 6.5B requirements.

---

# 40. Adapter Metrics

Populate the Phase 6.5A adapter metrics for Claude Code.

At minimum:

```text
agent = claude-code

commands_intercepted
commands_bypassed
adapter_failures
adapter_latency
```

When useful, record bypass reason:

```text
background
explicit bypass
unsupported invocation
recursion
integration failure
```

Keep adapter overhead distinct from AgentCap core processing latency.

---

# 41. Tool-Use Identity

Retain Claude's tool-use identifier as adapter/debug metadata if useful.

It must not replace AgentCap result IDs.

Conceptually:

```text
Claude tool_use_id
    -> transport/debug correlation

AgentCap result ID
    -> stored command result / progressive disclosure
```

Do not expose Claude internal IDs unnecessarily in normal capsule output.

---

# 42. Concurrency

Claude Code may issue tool calls concurrently.

The adapter must not keep invocation state in unsafe process-global mutable variables.

Verify that two Bash calls from one Claude session cannot:

- swap outputs
- swap cwd
- swap result IDs
- corrupt session state
- overwrite temporary files
- share transient recursion state incorrectly

Reuse Phase 6.5A concurrency-safe mechanisms.

---

# 43. Large Output

Test at least one large stdout result and one large failure result.

Verify:

```text
raw result stored completely
AgentCap capsule returned
no hook pipe deadlock
no accidental truncation before AgentCap storage
no excessive in-memory duplication where avoidable
```

Be especially careful if hook payloads themselves contain large command output.

---

# 44. Benchmark Compatibility

The upcoming benchmark must be able to compare:

```text
Claude Code + AgentCap OFF

vs

Claude Code + AgentCap ON
```

without changing unrelated Claude configuration between runs.

The OFF condition should use the shared bypass/disable mechanism rather than uninstalling and reinstalling integration for every benchmark task if possible.

Benchmark records must be attributable to:

```text
agent = claude-code
adapter version
AgentCap enabled/disabled
```

---

# 45. Do Not Add Claude-Specific Prompt Instructions Unless Required

Do not solve integration problems by adding large instructions to `CLAUDE.md`.

Transparent integration is preferred.

The model should continue issuing ordinary commands such as:

```bash
git diff
go test ./...
cargo check
rg foo .
```

rather than being instructed to manually write:

```bash
acap run ...
```

for every command.

A very small instruction explaining progressive disclosure may be acceptable only if existing Phase 6 integration already requires it.

Measure its token overhead if added.

---

# 46. No Reducer Changes

Do not tune AgentCap reducers during this phase merely because Claude behaves differently.

If Claude requests raw output too frequently, record that behavior for benchmark/evaluation work.

Do not hide adapter problems by changing unrelated compression logic.

Reducer optimization belongs after measurement unless a clear correctness bug is found.

---

# 47. No Benchmark Optimization

Do not special-case benchmark commands.

The adapter must work for ordinary Claude Code usage.

Do not detect fixture names, benchmark repositories, or benchmark task IDs to change compression behavior.

The benchmark must measure the real adapter.

---

# 48. Security

Treat all hook input as untrusted structured input.

Do not construct shell commands by unsafe concatenation.

In particular, a command such as:

```bash
echo '"; rm -rf /; #'
```

must remain command data, not become syntax accidentally introduced by the adapter's wrapper construction.

If wrapping requires transporting a shell command through another process, use an encoding/argument mechanism that preserves it exactly.

Avoid temporary shell scripts containing unsafely interpolated command text.

---

# 49. Cross-Platform Handling

Support the operating systems already supported by AgentCap.

Do not claim additional platform compatibility in this phase.

Pay attention to:

```text
path separators
executable lookup
settings paths
shell invocation
temporary files
quoting
```

If Windows Claude integration requires materially different behavior, isolate that behavior in the Claude adapter rather than the core.

---

# 50. Tests — Hook Parsing

Add unit tests for Claude hook payload parsing.

Cover at least:

```text
valid Bash input
optional description
optional timeout
background flag
unknown additional fields
missing required fields
wrong tool name
malformed JSON
session ID
cwd
tool_use_id
```

Unknown future Claude fields must not break parsing.

---

# 51. Tests — Successful Bash

Test successful commands through the adapter:

```text
small stdout
large stdout
stderr with exit 0
stdout + stderr
no output
Git output
successful tests
```

Verify Claude-facing output is AgentCap's presentation where interception applies.

---

# 52. Tests — Failed Bash

This is mandatory.

Test:

```text
exit 1
exit 7
compiler failure
test failure
large stderr
stdout + stderr before failure
command-not-found
```

Verify that failure output follows the selected safe AgentCap integration path.

Do not rely only on successful command tests.

---

# 53. Tests — Shell Semantics

Include shell-sensitive cases:

```bash
printf '%s\n' "a b"
```

```bash
FOO="hello world" sh -c 'printf "%s\n" "$FOO"'
```

```bash
printf 'a\nb\n' | grep b
```

```bash
false || echo recovered
```

```bash
mkdir -p nested && cd nested && pwd
```

```bash
printf 'x\n' > file.txt
```

The observable behavior with AgentCap enabled must match normal Claude Bash execution.

---

# 54. Tests — Permission Behavior

Where automated testing is practical, verify that integration does not convert:

```text
denied command
```

into:

```text
allowed command
```

or otherwise bypass Claude Code's normal permission system.

At minimum, document and manually verify the effect of any PreToolUse command rewrite on:

```text
allow rules
deny rules
ask behavior
auto mode
bypassPermissions mode
```

Permission equivalence is part of acceptance, not an optional polish item.

---

# 55. Tests — Session State

Use one Claude session equivalent and run a repeated command sequence.

Verify stateful compression.

Then create a second Claude session with the same project.

Verify that state does not leak incorrectly.

Also verify same session identifier in two different repositories does not mix project state.

---

# 56. Tests — Installation

Cover:

```text
empty settings
existing settings
existing unrelated hooks
existing Bash hooks
AgentCap already installed
remove after install
remove when absent
malformed settings
```

Never destroy unrelated settings.

---

# 57. Shared Conformance Suite

Run the Phase 6.5A adapter conformance suite against the Claude adapter.

The Claude implementation must pass the common requirements for:

```text
execution correctness
exactly once
cwd
environment
exit status
session mapping
project isolation
recursion
bypass
fail-open
```

If a Claude platform constraint prevents one of these guarantees, do not silently skip it.

Document the exact limitation and keep Phase 6.5B incomplete unless the limitation has an intentionally accepted safe behavior.

---

# 58. Manual End-to-End Verification

Before completion, run a real Claude Code session in a fixture repository.

Verify this workflow:

```text
1. Claude runs repository search.
2. Claude runs tests.
3. Tests fail.
4. Claude sees compact AgentCap failure output.
5. Claude requests targeted detail if needed.
6. Claude edits code.
7. Claude reruns tests.
8. AgentCap returns a smaller stateful delta.
9. Tests eventually pass.
10. Claude receives compact success output.
```

Inspect `.acap` afterward and verify raw results are recoverable.

---

# 59. Compare Enabled vs Bypassed Behavior

For representative commands, compare:

```text
AgentCap bypassed
```

against:

```text
Claude adapter enabled
```

Confirm equivalent command semantics for:

```text
filesystem effects
Git effects
exit status
cwd
environment
timeout
```

Only the agent-visible presentation should intentionally differ.

---

# 60. Non-Goals

Do not implement:

- Antigravity adapter
- Gemini CLI adapter
- OpenCode adapter
- new AgentCap reducers
- benchmark scoring
- benchmark task generation
- semantic source analysis
- AST integration
- general Claude tool interception
- Claude transcript analysis
- prompt rewriting
- model-response rewriting
- telemetry upload
- remote services

This phase is specifically the Claude Code shell adapter.

---

# 61. Acceptance Criteria

Phase 6.5B is complete only when all of the following are true.

## Integration

- Claude Code Bash commands pass through the Phase 6.5A adapter boundary.
- Claude-specific code remains isolated from AgentCap core.
- Normal Claude command syntax remains unchanged from the model's perspective.

## Exactly-once execution

- Successful commands execute exactly once.
- Failed commands execute exactly once.
- Side-effect tests prove there is no duplicate execution.

## Output

- Successful command output can be processed by AgentCap.
- Failed command output can be processed by AgentCap.
- Large build/test failures do not bypass compression merely because the command exited non-zero.
- Raw stdout/stderr remain recoverable.

## Command semantics

- shell syntax is preserved.
- cwd is preserved.
- environment behavior is preserved.
- exit status is preserved.
- timeout behavior is preserved.
- background commands are either correctly supported or conservatively bypassed.

## Permission safety

- integration does not intentionally bypass Claude Code permission checks.
- any command rewriting has been tested against permission behavior.
- destructive commands do not become implicitly approved by the adapter.

## Session behavior

- Claude `session_id` maps correctly to AgentCap session state.
- resumed sessions behave consistently.
- separate Claude sessions remain isolated.
- separate projects remain isolated.

## Progressive disclosure

- AgentCap result IDs are visible.
- `acap show` works.
- `acap raw` works.
- retrieval commands do not recurse.

## Installation

- project-level installation works.
- installation is idempotent.
- unrelated Claude settings are preserved.
- unrelated hooks are preserved.
- removal deletes only AgentCap-owned configuration.

## Reliability

- bypass works.
- fail-open works.
- adapter failures are distinguishable from command failures.
- concurrent calls do not corrupt results.

## Benchmark readiness

Metrics can identify:

```text
agent = claude-code
adapter version
commands intercepted
commands bypassed
adapter failures
adapter latency
```

and the adapter can be cleanly enabled/bypassed for A/B benchmark runs.

## Testing

- Claude-specific tests pass.
- Phase 6.5A shared conformance tests pass.
- all existing AgentCap tests continue to pass.
- manual Claude Code end-to-end verification succeeds.

---

# 62. Implementation Order

Use this order:

```text
1. Inspect Phase 6.5A common adapter contract.

2. Build minimal Claude hook feasibility tests.

3. Prove successful and failed Bash integration.

4. Choose the safest exactly-once architecture.

5. Implement Claude hook parsing.

6. Implement session mapping.

7. Implement Bash integration.

8. Implement recursion and bypass behavior.

9. Implement fail-open handling.

10. Add adapter metrics.

11. Add project-local installer.

12. Add safe uninstaller.

13. Run common conformance tests.

14. Run Claude-specific shell/permission tests.

15. Run real Claude Code end-to-end verification.
```

Do not start with installer work before the execution architecture has been proven.

---

# 63. Design Priority

When tradeoffs arise, use this priority order:

```text
1. Preserve command semantics.
2. Preserve permission/security semantics.
3. Execute exactly once.
4. Preserve failure information.
5. Preserve raw-result recoverability.
6. Preserve session correctness.
7. Preserve fail-open behavior.
8. Reduce model-visible output.
9. Minimize adapter latency.
```

A larger but correct result is preferable to a smaller result produced by unsafe or semantically incorrect interception.

---

# 64. Final Deliverable

At the end of Phase 6.5B, Claude Code should be a first-class AgentCap integration.

Claude should continue issuing ordinary commands:

```bash
git diff
go test ./...
cargo check
rg Resolve .
```

while AgentCap transparently provides:

```text
capture
storage
reduction
session state
deltas
result IDs
progressive disclosure
```

without duplicate execution and without requiring Claude to manually prefix every command with `acap run`.

The implementation must be suitable for immediate use in the upcoming AgentCap benchmark phase.
