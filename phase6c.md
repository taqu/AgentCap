# Phase 6.5C — Antigravity Adapter

## Objective

Implement the Antigravity adapter on top of the common adapter contract established in Phase 6.5A.

The adapter must allow normal Antigravity `run_command` usage to benefit from AgentCap while preserving:

- command semantics
- Antigravity permission behavior
- exactly-once command execution
- exit status
- cwd
- environment behavior
- asynchronous/background execution behavior
- AgentCap session state
- progressive disclosure
- fail-open behavior

The adapter must remain thin.

Do not move Antigravity-specific behavior into AgentCap core.

---

# 1. Required Architecture

The intended architecture is:

```text
Antigravity
    |
    | run_command
    v
Antigravity Adapter
    |
    v
Phase 6.5A Common Adapter Boundary
    |
    v
AgentCap Core
```

Antigravity-specific responsibilities belong only in the Antigravity adapter.

The AgentCap core must not gain logic such as:

```text
if agent == "antigravity" { ... }
```

The adapter should translate Antigravity hook payloads into the normalized integration structures introduced in Phase 6.5A.

---

# 2. Supported Antigravity Surfaces

Current Antigravity hooks are available across:

```text
Antigravity 2.0
Antigravity CLI
Antigravity IDE
```

The same adapter design should be usable across these surfaces wherever their hook contracts are equivalent.

Do not build three independent adapters unless actual behavioral differences require it.

Conceptually:

```text
Antigravity 2.0
       |
Antigravity CLI
       |
Antigravity IDE
       |
       v
same Antigravity Adapter
       |
       v
AgentCap
```

Surface-specific differences should remain small and explicit.

---

# 3. Important Hook Constraint

Antigravity provides:

```text
PreToolUse
PostToolUse
PreInvocation
PostInvocation
Stop
```

For this phase, only the tool-execution lifecycle should normally matter.

`PreToolUse` can inspect and gate a proposed tool call.

`PostToolUse` runs after the tool completes.

However, the current PostToolUse output contract is:

```json
{}
```

It does not provide a documented mechanism for replacing the completed tool output with an AgentCap capsule.

Therefore, the following architecture is not sufficient:

```text
Antigravity executes run_command
        |
        v
PostToolUse
        |
        v
AgentCap compresses something
```

because the already-visible command output cannot simply be replaced through PostToolUse.

Do not assume undocumented output-replacement behavior.

---

# 4. Start With an Integration Feasibility Test

Before implementing substantial adapter infrastructure, build a minimal integration spike using a real Antigravity installation.

Test at least:

```bash
printf 'hello\n'
```

```bash
sh -c 'echo failure-output; exit 7'
```

and a realistic failing command such as:

```bash
go test ./...
```

in a fixture repository.

The spike must establish:

```text
1. Exactly when PreToolUse runs.

2. Exactly what run_command arguments are available.

3. Whether CommandLine can be safely redirected or wrapped.

4. Whether Cwd can be preserved exactly.

5. Whether asynchronous/background semantics survive wrapping.

6. Whether the original command can execute exactly once.

7. Whether Antigravity receives AgentCap's presentation
   instead of the original raw command output.

8. Whether non-zero exit status remains visible correctly.

9. Whether permission behavior remains equivalent.

10. Whether cancellation/termination propagates correctly.
```

Do not proceed with a design that cannot satisfy the critical invariants.

---

# 5. Critical Integration Strategy

Because PostToolUse cannot replace tool output, the likely integration point is before command execution.

Conceptually:

```text
Antigravity proposes:

    go test ./...

        |
        v

PreToolUse
        |
        v

Antigravity Adapter
        |
        v

execution transformed to AgentCap integration path
        |
        v

original command executes exactly once
        |
        v

AgentCap presentation becomes command output
```

The exact implementation may differ based on the existing Phase 6.5A integration boundary.

Do not introduce a new execution engine solely for Antigravity.

Reuse AgentCap's existing execution and capture path wherever possible.

---

# 6. Exactly-Once Execution Is Mandatory

The most important correctness invariant remains:

```text
Every intercepted command must execute exactly once.
```

Never produce:

```text
Antigravity executes command
+
AgentCap executes command again
```

This applies especially to:

```text
rm
mv
cp
git add
git commit
git checkout
go generate
database migrations
package installation
custom scripts
```

The Antigravity adapter must prove exactly-once behavior through tests.

---

# 7. Do Not Use PostToolUse as a Second Execution Path

Do not implement:

```text
run_command executes normally

PostToolUse sees command metadata

AgentCap reruns command to obtain output
```

This is strictly incorrect.

PostToolUse may be useful for:

```text
diagnostics
metrics
correlation
verification
cleanup
```

but it must never cause the command to execute again.

---

# 8. Inspect Phase 6.5A First

Before adding Antigravity-specific execution logic:

1. inspect the Phase 6.5A common adapter contract
2. identify the normalized invocation type
3. identify the shared execution boundary
4. identify the common recursion mechanism
5. identify the bypass mechanism
6. identify common session mapping helpers
7. identify shared metrics
8. identify fail-open behavior

Reuse these mechanisms.

Do not duplicate them in the Antigravity package.

---

# 9. Target Only `run_command`

The primary hook matcher should be:

```text
run_command
```

Do not match all Antigravity tools.

Do not initially intercept:

```text
view_file
write_to_file
replace_file_content
multi_replace_file_content
list_dir
browser tools
agent collaboration tools
```

AgentCap's responsibility is command-result mediation.

Keep the integration narrowly scoped.

---

# 10. Hook Configuration

Use workspace-local:

```text
.agents/hooks.json
```

as the initial integration mechanism.

Prefer workspace-local configuration for reproducible benchmark environments.

Conceptually, AgentCap should register a hook with:

```text
PreToolUse
matcher = run_command
```

Additional PostToolUse hooks may be used only where useful for metrics or verification.

Do not install unnecessary invocation-wide hooks.

---

# 11. Antigravity PreToolUse Input

The adapter should parse the documented PreToolUse payload.

Relevant fields include:

```text
toolCall
toolCall.name
toolCall.args
stepIdx

conversationId
workspacePaths
transcriptPath
artifactDirectoryPath
modelName
```

For `run_command`, current arguments include concepts such as:

```text
CommandLine
Cwd
WaitMsBeforeAsync
```

Only depend on fields actually required by AgentCap.

Unknown fields must be ignored safely.

Do not permanently mirror the entire Antigravity schema into AgentCap types.

---

# 12. Normalized Invocation

Translate Antigravity-specific input into the Phase 6.5A normalized invocation.

Conceptually:

```text
AgentInvocation
    agent              = antigravity
    externalSession    = conversationId
    command            = CommandLine
    cwd                = Cwd
    project context    = workspacePaths
    transport metadata = stepIdx / modelName / etc.
```

Do not let Antigravity types leak past the adapter boundary.

---

# 13. Session Mapping

Use:

```text
conversationId
```

as the primary external Antigravity session identifier.

Conceptually:

```text
agent            = antigravity
external session = conversationId
project          = AgentCap project root
```

Required behavior:

```text
same Antigravity conversation
    -> same AgentCap session

different Antigravity conversation
    -> different AgentCap session
```

Do not derive session identity from:

```text
PID
stepIdx
cwd alone
transcript filename alone
timestamp
```

when `conversationId` is available.

---

# 14. Workspace Paths Are Not the Execution Directory

Antigravity provides workspace-level path information.

Do not confuse:

```text
workspacePaths
```

with:

```text
run_command Cwd
```

The command must execute using the actual invocation cwd.

Maintain AgentCap's existing distinction:

```text
project root
!=
execution cwd
```

AgentCap project-root discovery must remain core-controlled.

---

# 15. Command String Preservation

`CommandLine` is shell command input.

Preserve it as shell syntax.

Do not naïvely split:

```bash
cd src && make test 2>&1 | tee build.log
```

into argv tokens.

Likewise preserve:

```bash
FOO="hello world" ./script.sh
```

```bash
printf 'a\nb\n' | grep b
```

```bash
false || echo recovered
```

```bash
for f in *.go; do echo "$f"; done
```

```bash
git diff -- '*.go'
```

The wrapper must not alter shell semantics.

---

# 16. Safe Command Transport

If `CommandLine` must be passed through an AgentCap wrapper, do not build the wrapper using unsafe string concatenation.

Avoid constructions equivalent to:

```text
"acap-wrapper " + originalCommand
```

when that causes re-parsing or quoting ambiguity.

Use an encoding or argument transport mechanism that preserves the complete shell command exactly.

Possible implementation approaches may include:

```text
stdin transport
environment transport with safe encoding
temporary descriptor
opaque argument encoding
existing AgentCap shell-command entry point
```

Use whichever fits the existing architecture best.

Do not introduce shell injection vulnerabilities.

---

# 17. Permission Semantics Are Critical

PreToolUse is part of Antigravity's command permission flow.

The hook response may return decisions such as:

```text
allow
deny
ask
force_ask
deny_unless_prior_grant
```

The AgentCap adapter must not turn itself into a permission bypass.

Do not simply return:

```json
{"decision":"allow"}
```

for all commands.

The adapter's presence must not make a destructive command less restricted than it would be without AgentCap.

---

# 18. Preserve Existing Permission Decisions

Investigate how Antigravity combines:

```text
normal permission/autonomy behavior
+
PreToolUse hook decisions
```

The AgentCap hook should interfere as little as possible.

In particular, verify behavior under Antigravity permission modes such as:

```text
request-review
always-proceed
strict
```

where supported.

Do not override the user's configured autonomy level merely to enable compression.

---

# 19. Permission Overrides

Do not use `permissionOverrides` unless they are genuinely required.

AgentCap should not grant new command permissions.

If the integration architecture requires a wrapper command such as an AgentCap executable, carefully ensure that permission checking continues to apply to the original command's semantics.

The user should not be approving:

```text
acap internal-wrapper
```

while the real underlying command:

```text
rm -rf ...
```

is hidden from the permission system.

This is a critical acceptance requirement.

---

# 20. Wrapper Visibility Problem

A PreToolUse wrapper may change the command that Antigravity sees as executing.

Explicitly verify that wrapping:

```text
original command
```

into:

```text
AgentCap wrapper
```

does not break:

```text
permission prompts
auditability
UI command display
user understanding
command cancellation
background-task display
```

If it does, prefer a safer integration strategy.

Do not hide the underlying command from the user-facing execution UI where avoidable.

---

# 21. PreToolUse Decision Strategy

Keep the AgentCap hook's decision behavior minimal.

The adapter is not a security policy engine.

It should normally avoid making policy judgments such as:

```text
this command is safe
this command is unsafe
```

unless required to preserve existing Antigravity behavior.

AgentCap's responsibility is result mediation, not authorization.

---

# 22. Failure to Establish Safe Wrapping

If the feasibility spike shows that command wrapping necessarily weakens Antigravity permissions or changes command semantics, do not silently accept it.

Use one of the following outcomes:

```text
A. identify another supported integration path

B. restrict interception to commands for which semantics can
   be proven equivalent

C. expose the integration limitation explicitly

D. leave unsupported command modes bypassed
```

Do not use undocumented internal APIs.

Do not modify Antigravity binaries.

---

# 23. Successful Command Handling

For a successful command:

```text
exit = 0
```

Antigravity should receive the normal AgentCap presentation.

Example concept:

```text
@acap abc123 go-test PASS
packages=42 duration=8.1s
```

Do not create a separate Antigravity-specific capsule format.

---

# 24. Failed Command Handling

For:

```text
exit != 0
```

Antigravity must still receive AgentCap's reduced result.

This is mandatory.

Test realistic failures such as:

```text
Go test failure
GCC/Clang failure
cargo check failure
cargo test failure
linker failure
```

A failing command must not revert to full raw output merely because of adapter architecture.

---

# 25. Exit Status

The original command exit status must remain semantically visible.

Test at least:

```text
0
1
2
7
127
```

where practical.

Compression must never convert:

```text
exit 7
```

into apparent success.

Do not return wrapper success in place of command failure.

---

# 26. Wrapper Exit Semantics

If Antigravity actually executes an AgentCap wrapper process, that wrapper must terminate with the underlying command's effective exit status.

Conceptually:

```text
original exit = 0
wrapper exit  = 0

original exit = 7
wrapper exit  = 7
```

Internal AgentCap processing failures must be handled separately.

Do not collapse:

```text
command failure
```

and:

```text
AgentCap reduction failure
```

into one exit status.

---

# 27. stdout / stderr

AgentCap must still capture:

```text
full stdout
full stderr
```

before reduction.

The Antigravity-facing output should be AgentCap's presentation.

Raw output must remain recoverable via:

```bash
acap show <id>
acap raw <id>
```

Do not merge stdout/stderr destructively merely because Antigravity presents command output as one UI result.

---

# 28. Progressive Disclosure

Antigravity must receive normal AgentCap result IDs.

Existing retrieval commands should work naturally:

```bash
acap show <id>
acap raw <id>
acap show <id> --errors
acap show <id> --warnings
acap show <id> --test ...
acap show <id> --file ...
```

Do not introduce:

```text
antigravity-acap-show
```

or any parallel retrieval protocol.

---

# 29. Recursion Protection

Commands such as:

```bash
acap show ...
acap raw ...
acap stats
acap clean
```

must not recursively re-enter transparent interception.

Reuse the shared Phase 6.5A mechanism.

Do not implement recursion prevention using only fragile string-prefix matching if a stronger marker already exists.

---

# 30. Explicit Bypass

Support the Phase 6.5A bypass mechanism.

This is required for:

```text
benchmark OFF runs
debugging
compatibility testing
adapter recovery
```

When bypass is active:

```text
Antigravity
    -> normal run_command behavior
```

without AgentCap reduction.

Do not require hook uninstall/reinstall for every benchmark A/B run.

---

# 31. Asynchronous Execution

Antigravity's `run_command` may switch to asynchronous execution according to fields such as:

```text
WaitMsBeforeAsync
```

Do not assume every command remains synchronous.

Test commands that:

```text
finish before async threshold
```

and:

```text
continue beyond async threshold
```

The adapter must not change a command from asynchronous to synchronous solely to simplify capture.

---

# 32. Background/Long-Running Commands

If long-running command semantics cannot currently be preserved safely through AgentCap:

```text
bypass them conservatively
```

rather than blocking indefinitely or changing Antigravity behavior.

Record the bypass reason.

Examples may include:

```text
dev servers
watch mode
tail -f
long-lived language servers
interactive programs
```

Do not claim transparent support without tests.

---

# 33. Cancellation

Verify that stopping a command through Antigravity terminates the actual child command.

Expected path:

```text
Antigravity cancellation
        |
        v
AgentCap wrapper
        |
        v
underlying command terminates
```

Do not leave orphan processes.

Do not detach processes accidentally through the wrapper.

---

# 34. Signals

Preserve AgentCap's existing signal behavior.

If the wrapper receives:

```text
SIGINT
SIGTERM
```

or platform equivalents, the underlying command should receive appropriate cancellation.

Do not add Antigravity-specific signal logic to AgentCap core unless absolutely necessary.

---

# 35. Environment

Preserve the Antigravity command environment.

Do not reconstruct environment values from hook JSON.

The wrapper should inherit the command's normal environment unless existing AgentCap execution semantics explicitly require otherwise.

Any AgentCap-specific environment variables should be narrowly scoped to:

```text
recursion protection
session propagation
adapter identification
bypass
```

Do not leak sensitive values into result output.

---

# 36. Model Name

`modelName` may be retained as optional benchmark/debug metadata.

Do not include it in session identity.

A model change inside the same Antigravity conversation should not automatically create a new AgentCap session.

Conceptually:

```text
conversationId
    -> session identity

modelName
    -> metadata
```

---

# 37. stepIdx

`stepIdx` may be useful for correlation and diagnostics.

Do not use it as a persistent session identity.

Do not assume it is globally unique.

It may be retained as transport metadata only.

---

# 38. transcriptPath

Do not parse Antigravity transcripts as part of normal command interception.

The adapter already receives the required structured tool invocation data.

Transcript parsing would add:

```text
complexity
brittleness
version coupling
privacy surface
```

without being necessary for the primary adapter.

Use transcript information only if later benchmark work explicitly requires it.

---

# 39. artifactDirectoryPath

Do not store AgentCap result data inside Antigravity artifact directories.

AgentCap storage remains:

```text
<project-root>/.acap/
```

The Antigravity artifact directory may be retained as optional metadata if useful for debugging.

Do not couple AgentCap persistence to Antigravity's internal conversation storage.

---

# 40. Multi-Workspace Sessions

Antigravity may provide multiple:

```text
workspacePaths
```

Do not simply select the first path and assume it is the project root.

Use actual execution cwd plus AgentCap's existing project-root discovery.

If a command runs outside all mounted workspace paths, preserve normal AgentCap root-resolution behavior.

Add tests where practical.

---

# 41. Concurrent Tool Calls

The adapter must tolerate concurrent `run_command` executions if Antigravity allows them.

Do not store invocation-specific data in unsynchronized global variables.

Verify that concurrent commands cannot exchange:

```text
CommandLine
Cwd
conversationId
result ID
stdout/stderr
exit status
```

Reuse the common Phase 6.5A concurrency-safe mechanisms.

---

# 42. PostToolUse Usage

PostToolUse may be used for lightweight purposes such as:

```text
confirmation that a tool completed
metrics correlation
diagnostic validation
cleanup
```

but only if useful.

Do not add PostToolUse merely because the hook exists.

Do not expect it to replace tool output.

Do not rerun commands from PostToolUse.

---

# 43. PostToolUse Error Field

Where PostToolUse is used, the documented:

```text
error
```

field may help correlate whether Antigravity considered the invocation failed.

Treat it as transport information.

Do not use it instead of AgentCap's captured exit status when AgentCap executed the actual command.

AgentCap execution remains authoritative for its stored result.

---

# 44. Fail-Open Behavior

If:

```text
original command executes
+
AgentCap reduction fails
```

prefer returning usable original command output rather than breaking the task.

Conceptually:

```text
execute
capture
try reduce

if reduce succeeds:
    return capsule
else:
    return safe original output
```

Preserve the actual command exit status.

Do not fabricate success or failure.

---

# 45. Wrapper Failure Before Command Execution

Distinguish:

```text
AgentCap wrapper failed before command started
```

from:

```text
command executed and failed
```

Do not claim exactly-once execution if the adapter cannot determine whether the command started.

Make the execution state explicit internally.

Where possible use an architecture in which command ownership is unambiguous.

---

# 46. AgentCap Binary Missing

Test behavior when the AgentCap executable cannot be found.

The integration should fail safely.

Depending on the selected architecture, either:

```text
bypass and execute normally
```

or:

```text
report a concise integration failure before execution
```

may be appropriate.

Do not accidentally execute the command twice while attempting recovery.

---

# 47. Project-Local Installer

Provide an idempotent installation command following existing AgentCap CLI conventions.

Conceptually:

```text
acap integrate antigravity
```

Installation should create or update:

```text
.agents/hooks.json
```

in the project/workspace.

Do not invent a separate management utility if Phase 6 already has integration commands.

---

# 48. Preserve Existing hooks.json

`.agents/hooks.json` may contain unrelated user hooks.

Installation must:

```text
read existing file
preserve unrelated hook definitions
add only AgentCap-owned configuration
avoid duplicate entries
write valid JSON
```

Do not replace the entire file.

---

# 49. Safe Uninstall

Support removal using existing AgentCap CLI conventions.

Conceptually:

```text
acap integrate antigravity --remove
```

Removal must delete only AgentCap-owned hook entries.

Verify:

```text
existing hooks
+
AgentCap install
+
AgentCap remove
=
existing hooks
```

Do not delete another hook that also matches `run_command`.

---

# 50. Idempotency

Verify:

```text
install
install
```

creates only one logical AgentCap integration.

Likewise:

```text
remove
remove
```

must be safe.

Use stable ownership markers or hook names.

Do not rely on array position.

---

# 51. Hook Name

Use a stable AgentCap-owned hook key.

Conceptually:

```text
agentcap
```

or an existing project naming convention.

It should be:

```text
recognizable
stable
safe to update
safe to remove
```

Do not use a generated random name on every installation.

---

# 52. Plugin Packaging Is Not Required Yet

Antigravity plugins may package:

```text
plugin.json
hooks.json
skills
agents
rules
MCP configuration
```

However, Phase 6.5C should initially prefer direct workspace hook installation.

Plugin packaging is optional future distribution work.

Do not block benchmark readiness on marketplace/plugin packaging.

---

# 53. Do Not Add Skills Unless Required

Do not add an Antigravity skill merely to teach the model to use AgentCap.

Transparent integration is preferred.

The agent should continue issuing:

```bash
git diff
go test ./...
cargo check
rg pattern .
```

not:

```bash
acap run ...
```

for every command.

Progressive disclosure should remain usable through normal shell commands when needed.

---

# 54. Adapter Metrics

Populate Phase 6.5A metrics with:

```text
agent = antigravity
```

and at least:

```text
commands_intercepted
commands_bypassed
adapter_failures
adapter_latency
```

Useful bypass reasons may include:

```text
explicit bypass
recursive AgentCap command
unsupported async mode
unsupported interactive command
integration failure
```

Keep adapter latency separate from core AgentCap processing latency.

---

# 55. Adapter Version

Expose enough version information for benchmark attribution.

Conceptually:

```text
agent         = antigravity
adapter       = antigravity
adapterVersion
```

Reuse AgentCap build/version metadata where sufficient.

Do not build a complex compatibility negotiation system.

---

# 56. Antigravity Version Compatibility

The adapter should tolerate additional unknown fields in hook payloads.

Do not fail because Antigravity adds unrelated fields.

At the same time, do not silently continue if required fields change incompatibly.

Required-field failures should produce clear diagnostics.

---

# 57. Hook Parser Tests

Add unit tests covering:

```text
valid run_command
CommandLine
Cwd
WaitMsBeforeAsync
conversationId
workspacePaths
modelName
stepIdx
unknown extra fields
missing fields
wrong tool name
malformed JSON
```

Unknown optional fields should not fail parsing.

---

# 58. Command Semantics Tests

Test shell-sensitive commands:

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

```bash
git diff -- '*.go'
```

With AgentCap enabled, behavior must match normal Antigravity command execution except for intentional presentation reduction.

---

# 59. Exactly-Once Tests

Use a side-effect fixture.

Example:

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

Also test:

```text
increment counter
exit 7
```

The counter must still increase once.

---

# 60. Build/Test Failure Tests

Test large failures through the real adapter:

```text
Go test failure
C/C++ compiler failure
Rust cargo failure
```

Verify:

```text
AgentCap captures full output
AgentCap stores raw output
Antigravity receives compact presentation
exit status remains non-zero
result ID is available
drill-down works
```

These tests are mandatory.

---

# 61. Stateful Compression Tests

Within one `conversationId`:

```text
run failing tests
edit
run failing tests again
edit
run tests again
```

Verify:

```text
first result  -> normal compact result
second result -> delta where appropriate
third result  -> delta/unchanged where appropriate
```

The adapter must not implement delta logic.

It must only preserve session identity correctly.

---

# 62. Session Isolation Tests

Verify:

```text
conversation A
conversation B
```

do not accidentally share AgentCap session state.

Also verify:

```text
same conversation ID
different project
```

does not cause cross-project mixing.

Project scope remains authoritative.

---

# 63. Async Tests

Test a command that returns quickly.

Test another that exceeds:

```text
WaitMsBeforeAsync
```

Verify the adapter does not:

```text
block incorrectly
lose output
duplicate execution
break cancellation
```

If async execution is intentionally bypassed initially, assert that bypass explicitly.

---

# 64. Permission Tests

Test representative commands under supported Antigravity permission/autonomy settings.

Include:

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

Verify that AgentCap does not weaken the permission policy.

Particularly verify that the wrapper does not obscure the actual underlying command from approval logic.

---

# 65. Bypass Tests

With AgentCap bypass enabled:

```text
run_command
    -> normal Antigravity execution
```

No AgentCap result should be produced for that intercepted execution unless existing bypass semantics intentionally record one.

This mode will be used by the benchmark control condition.

---

# 66. Recursion Tests

Run:

```bash
acap show <id>
```

```bash
acap raw <id>
```

```bash
acap stats
```

from Antigravity.

Verify these commands do not recursively wrap themselves.

---

# 67. Fail-Open Tests

Inject failures in:

```text
reducer
storage if safely testable
presentation
adapter metadata
```

and verify the safest fallback.

Do not lose original command output when a recoverable AgentCap processing failure occurs.

Do not execute the command again.

---

# 68. Large Output Tests

Test:

```text
large stdout
large stderr
large mixed output
large compiler failure
```

Verify:

```text
no pipe deadlock
no unexpected truncation
raw output stored completely
presentation reduced
memory use remains reasonable
```

---

# 69. Shared Phase 6.5A Conformance Suite

Run the common adapter conformance suite against the Antigravity adapter.

The adapter must satisfy:

```text
execution correctness
exactly-once execution
cwd
environment
exit status
session mapping
project isolation
recursion protection
bypass
fail-open
```

Do not silently skip failed conformance requirements because of Antigravity-specific behavior.

---

# 70. Real End-to-End Antigravity Verification

Before completion, run a real Antigravity coding workflow in a fixture repository.

Verify:

```text
1. Antigravity runs repository search.

2. Antigravity runs tests.

3. Tests fail.

4. Antigravity receives compact AgentCap failure output.

5. Antigravity requests targeted detail if needed.

6. Antigravity edits code.

7. Antigravity reruns tests.

8. AgentCap returns stateful delta information.

9. Tests eventually pass.

10. Antigravity receives compact success output.
```

Inspect `.acap` afterward.

Verify that full raw command results remain recoverable.

---

# 71. Compare Enabled vs Bypassed Behavior

For representative commands compare:

```text
AgentCap OFF
```

and:

```text
Antigravity Adapter ON
```

Verify equivalent:

```text
filesystem effects
Git effects
exit status
cwd
environment
async behavior
cancellation
```

Only model-visible result presentation should intentionally differ.

---

# 72. Benchmark Compatibility

The upcoming benchmark must be able to run:

```text
Antigravity + AgentCap OFF
```

versus:

```text
Antigravity + AgentCap ON
```

without changing unrelated Antigravity configuration.

Use bypass/enable state rather than reinstalling hooks for every task where possible.

Benchmark records must be attributable to:

```text
agent = antigravity
adapter version
AgentCap enabled/disabled
```

---

# 73. Do Not Tune Reducers in This Phase

Do not modify reducer behavior simply because Antigravity reacts differently from Claude Code.

If Antigravity frequently requests:

```bash
acap raw ...
```

record that for the benchmark.

Do not optimize against anecdotal adapter testing.

Reducer tuning should be measurement-driven.

---

# 74. No Benchmark-Specific Behavior

Do not detect:

```text
benchmark fixture repository
benchmark task ID
benchmark command pattern
```

and change adapter behavior.

The benchmark must measure the real production adapter.

---

# 75. No Transcript-Based Token Optimization

Do not analyze or rewrite Antigravity conversation transcripts in this phase.

AgentCap optimizes command-result presentation.

It does not become:

```text
conversation compressor
prompt compressor
model-context rewriter
```

Keep the project boundary narrow.

---

# 76. No General Antigravity Integration Framework

Do not build a framework for every Antigravity tool.

This phase needs:

```text
run_command
```

only.

Avoid speculative abstractions for:

```text
browser
file editing
MCP
subagents
skills
artifact tools
```

---

# 77. No Plugin Marketplace Work

Do not:

```text
publish plugin
build marketplace metadata
implement plugin discovery
implement remote installation
```

during Phase 6.5C.

A later packaging phase can convert the adapter into an Antigravity plugin if useful.

Benchmark readiness is the priority.

---

# 78. Documentation

Add concise developer documentation covering:

```text
how Antigravity hooks invoke the adapter
why PreToolUse is required
why PostToolUse alone is insufficient
how run_command is normalized
how conversationId maps to AgentCap sessions
how exactly-once execution is preserved
how permissions are preserved
how async commands are handled
how bypass works
how hooks are installed
how hooks are removed
```

Keep documentation implementation-oriented.

---

# 79. Non-Goals

Do not implement:

- Claude Code changes
- Gemini CLI adapter
- OpenCode adapter
- general Antigravity tool interception
- new reducers
- semantic source navigation
- AST integration
- transcript analysis
- prompt compression
- benchmark scoring
- marketplace publishing
- remote telemetry
- daemon infrastructure

This phase is specifically the Antigravity `run_command` adapter.

---

# 80. Acceptance Criteria

Phase 6.5C is complete only when all of the following are true.

## Architecture

- Antigravity-specific logic is isolated in its adapter.
- The Phase 6.5A common adapter contract is used.
- AgentCap core contains no Antigravity-specific execution branches.

## Hook integration

- Workspace-local `.agents/hooks.json` integration works.
- `run_command` is targeted explicitly.
- unrelated Antigravity tools are not intercepted.
- PostToolUse is not incorrectly relied upon for output replacement.

## Exactly-once execution

- successful commands execute once.
- failed commands execute once.
- side-effect tests prove no duplicate execution.

## Output behavior

- successful output goes through AgentCap.
- failed output goes through AgentCap.
- full raw stdout/stderr remain recoverable.
- result IDs are available.

## Semantics

- CommandLine shell behavior is preserved.
- cwd is preserved.
- environment is preserved.
- exit status is preserved.
- cancellation is preserved.
- supported async behavior is preserved.
- unsupported execution modes are conservatively bypassed.

## Permission safety

- AgentCap does not broadly return allow for all commands.
- AgentCap does not intentionally weaken Antigravity permission behavior.
- original command semantics remain visible to permission handling as far as technically possible.
- wrapper execution does not become an authorization bypass.

## Session behavior

- `conversationId` maps to AgentCap session identity.
- separate conversations remain isolated.
- separate projects remain isolated.
- model changes do not accidentally reset sessions.

## Progressive disclosure

- `acap show` works.
- `acap raw` works.
- retrieval commands do not recurse.
- no Antigravity-specific result store exists.

## Installation

- installation safely updates `.agents/hooks.json`.
- existing hooks are preserved.
- installation is idempotent.
- removal deletes only AgentCap-owned configuration.

## Reliability

- explicit bypass works.
- fail-open works.
- wrapper failures are distinguishable from command failures.
- concurrent invocations do not corrupt state.

## Benchmark readiness

Metrics expose or record:

```text
agent = antigravity
adapter version
commands intercepted
commands bypassed
adapter failures
adapter latency
```

and AgentCap can be enabled/bypassed cleanly for benchmark A/B runs.

## Testing

- Antigravity-specific unit tests pass.
- shell-semantic tests pass.
- failure-output tests pass.
- permission tests pass.
- async tests pass or explicitly prove safe bypass behavior.
- Phase 6.5A conformance tests pass.
- all existing AgentCap tests continue to pass.
- real Antigravity end-to-end verification succeeds.

---

# 81. Implementation Order

Use this order:

```text
1. Inspect Phase 6.5A common adapter contract.

2. Build minimal real Antigravity PreToolUse feasibility tests.

3. Prove exactly-once command wrapping.

4. Prove failed commands can return AgentCap presentation.

5. Verify permission behavior before committing to the architecture.

6. Implement Antigravity hook payload parsing.

7. Implement conversationId -> AgentCap session mapping.

8. Implement run_command interception.

9. Preserve cwd/environment/exit semantics.

10. Implement async-safe handling or conservative bypass.

11. Implement recursion and explicit bypass.

12. Implement fail-open behavior.

13. Add adapter metrics.

14. Implement safe .agents/hooks.json installer.

15. Implement safe uninstaller.

16. Run the Phase 6.5A conformance suite.

17. Run Antigravity-specific permission and shell tests.

18. Run real end-to-end Antigravity workflow verification.
```

Do not begin with plugin packaging or installer polish before command interception has been proven safe.

---

# 82. Design Priority

When tradeoffs arise, use this priority order:

```text
1. Preserve command semantics.

2. Preserve permission/security semantics.

3. Execute exactly once.

4. Preserve failure information.

5. Preserve exit semantics.

6. Preserve raw-result recoverability.

7. Preserve session correctness.

8. Preserve asynchronous/cancellation behavior.

9. Reduce model-visible output.

10. Minimize adapter latency.
```

A larger correct result is preferable to a smaller result produced through unsafe interception.

---

# 83. Final Deliverable

At the end of Phase 6.5C, Antigravity should be a first-class AgentCap integration.

The agent should continue issuing ordinary commands such as:

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

without duplicate execution, without weakening Antigravity permissions, and without requiring the model to manually prefix commands with:

```bash
acap run
```

The implementation must be suitable for immediate use in the upcoming cross-agent benchmark phase.
