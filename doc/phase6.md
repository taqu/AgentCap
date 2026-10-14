# AgentCap Phase 6 — Coding Agent Integration and Transparent Command Interception

## Project

Project name: AgentCap  
CLI command: `acap`  
Language: Go

AgentCap is a command-output reduction and progressive-disclosure layer for AI coding agents.

Previous phases established:

- Phase 0: reliable command execution
- Phase 1: command-aware output reduction
- Phase 2: persistent results, capsules, and drill-down
- Phase 2.1: project-local SQLite storage under `.acap/store.db`
- Phase 3: session-aware deduplication and delta results
- Phase 4: Git-aware compression and structural Git deltas
- Phase 5: build/compiler/test-aware compression and diagnostic deltas

Phase 6 introduces coding-agent integration.

The central goal is:

> Allow coding agents to benefit from AgentCap automatically, without requiring the model to manually prefix every shell command with `acap run`.

Phase 6 must add integration adapters around the existing AgentCap core.

It must NOT redesign the execution, reduction, storage, session, Git, or build/test subsystems.

---

# 1. Phase 6 Goal

Implement a stable integration layer that allows supported coding agents to route shell/tool execution through AgentCap automatically.

Initial integration targets:

```text
Claude Code
OpenAI Codex
```

Secondary targets, only after the architecture is proven:

```text
Gemini CLI
OpenCode
generic hook-based agents
```

The primary user experience should move from:

```bash
acap run rg "Workspace" .
acap run git diff
acap run go test ./...
```

toward:

```bash
rg "Workspace" .
git diff
go test ./...
```

while AgentCap transparently captures and reduces the command result before it reaches the coding agent.

---

# 2. Core Architectural Rule

AgentCap core must remain independent of any coding-agent implementation.

The desired architecture is:

```text
Coding Agent
     |
     v
Agent-Specific Adapter
     |
     v
AgentCap Integration Protocol
     |
     v
Existing AgentCap Core
     |
     +--> executor
     +--> reducers
     +--> result store
     +--> sessions
     +--> delta engine
     +--> renderer
```

Do not introduce conditionals throughout core packages such as:

```go
if agent == "claude" { ... }

if agent == "codex" { ... }
```

Agent-specific behavior must remain isolated in adapters.

---

# 3. Product Principle

Phase 6 is not about adding more compression algorithms.

It is about making existing AgentCap behavior effortless to use.

The desired experience is:

```text
agent issues command
       |
       v
AgentCap intercepts execution
       |
       v
real command executes once
       |
       v
raw output stored locally
       |
       v
compact capsule/delta generated
       |
       v
agent receives compact output
```

The coding agent should not need special knowledge of every AgentCap reducer.

---

# 4. Strict Non-Goal: Do Not Become a Shell

AgentCap must not become a general-purpose shell.

Do NOT implement:

- full shell parsing
- shell expansion
- command pipelines
- variable expansion
- glob expansion
- job control
- shell scripting language
- shell history

If an agent executes:

```bash
sh -c 'foo | bar'
```

or:

```bash
bash -lc '...'
```

AgentCap may wrap that process as one command.

Do not reimplement shell semantics.

---

# 5. Integration Strategy

Prefer official, documented agent hooks or tool interception mechanisms where available.

The integration layer should support two broad modes:

```text
Mode A: pre/post tool hook integration

Mode B: explicit shell wrapper integration
```

Prefer Mode A when the agent exposes reliable hooks.

Use Mode B only where hook integration is unavailable or insufficient.

Do not patch coding-agent binaries.

Do not rely on undocumented memory injection or process modification.

---

# 6. Integration Boundary

The adapter is responsible for translating:

```text
agent tool request
```

into:

```text
AgentCap execution request
```

and translating:

```text
AgentCap presentation
```

back into:

```text
agent tool result
```

The adapter must not contain:

- Git parsing
- compiler parsing
- test parsing
- session delta logic
- storage logic
- token reduction heuristics

Those belong to existing AgentCap packages.

---

# 7. Introduce an Internal Integration Protocol

Define a small stable internal request/response model.

For example:

```go
type ToolRequest struct {
    Command    []string
    WorkingDir string
    Env        map[string]string

    SessionID string

    StdinMode string
    TTY       bool
}
```

And:

```go
type ToolResponse struct {
    ResultID string

    ExitCode int

    Stdout string
    Stderr string

    Presentation string
}
```

Exact fields may differ.

The purpose is to decouple agent adapters from AgentCap internals.

---

# 8. Do Not Use Public JSON Unless Necessary

The internal Go API should be preferred when the adapter executes inside the same process.

If adapters must communicate through subprocesses, introduce a compact machine-readable protocol.

JSON is acceptable for control messages.

Do not use verbose JSON as the normal agent-facing command result.

Remember:

```text
machine protocol != model-visible output
```

The model-visible response should remain AgentCap's compact capsule format.

---

# 9. New CLI Integration Surface

Add a dedicated machine-facing command if necessary.

Possible interface:

```bash
acap exec --protocol=json
```

or:

```bash
acap hook execute
```

The exact naming may differ.

Do NOT overload normal:

```bash
acap run
```

with large amounts of agent-specific behavior.

`acap run` should remain usable manually.

---

# 10. Preserve Existing CLI

These existing commands must continue to work:

```text
acap run <command...>
acap show <id>
acap raw <id>
acap stats
acap clean
```

Phase 6 must not require an agent integration to use AgentCap.

The direct CLI remains the debugging and fallback interface.

---

# 11. Agent Session Mapping

Coding-agent sessions should map naturally onto:

```text
ACAP_SESSION_ID
```

or the equivalent Phase 3 session model.

Each adapter should attempt to obtain a stable session identity from the agent when available.

If the agent exposes a conversation/run/session identifier:

```text
agent session
     |
     v
AgentCap session
```

Use it after applying a safe local normalization if necessary.

Do not expose sensitive remote identifiers unnecessarily.

A local hash is acceptable.

---

# 12. Missing Agent Session ID

If the coding agent does not expose a stable session ID, the adapter may generate a local session ID for the lifetime of the integration process.

Example:

```text
acap-agent-7f31c2
```

Do not fall back to a global session shared by unrelated agent runs.

Session isolation is more important than maximizing deduplication.

---

# 13. Session Lifetime

The AgentCap session should approximately match the coding-agent task/conversation lifetime.

Do not create a new AgentCap session for every command.

That would defeat Phase 3 stateful compression.

Likewise, do not reuse one session indefinitely across unrelated tasks.

---

# 14. Working Directory Preservation

The adapter must preserve the exact working directory requested by the coding agent.

Do not automatically execute from the detected AgentCap project root.

These remain distinct:

```text
project root
execution cwd
```

For example:

```text
project_root=/repo
cwd=/repo/src/parser
```

The command must execute from:

```text
/repo/src/parser
```

while AgentCap storage remains:

```text
/repo/.acap/
```

---

# 15. Environment Preservation

The adapter must preserve the environment semantics expected by the coding agent.

Do not silently alter:

```text
PATH
LANG
LC_ALL
GIT_*
RUSTFLAGS
CFLAGS
GOFLAGS
```

except for integration-specific variables that are strictly necessary.

Any injected AgentCap environment variable should be minimal and documented.

---

# 16. Exit Code Preservation

This remains a strict requirement.

If the underlying command exits:

```text
7
```

the adapter must report exit code:

```text
7
```

to the coding agent wherever the agent API supports explicit exit status.

Do not convert a failed command into a successful tool call merely because AgentCap successfully compressed its output.

Distinguish:

```text
command status
```

from:

```text
AgentCap processing status
```

---

# 17. Tool Failure Semantics

There are two independent failure classes:

```text
child command failure
AgentCap integration failure
```

Examples:

```text
go test -> exit 1
```

is a valid command result.

But:

```text
AgentCap could not execute/capture command
```

is an integration failure.

Adapters must preserve this distinction.

---

# 18. Fail-Open Principle

AgentCap integration is an optimization.

If compression/reduction/session logic fails but the child command result is available:

> return the command's normal output rather than failing the coding-agent operation.

Examples:

```text
reducer panic/error
session DB unavailable
delta comparison failure
Git parser unsupported
test parser unsupported
```

should generally fall back to:

```text
captured raw or conservative generic result
```

Do not block coding work merely because compression failed.

---

# 19. Fail-Closed Cases

Do NOT fail open when command execution itself did not occur or its result is ambiguous.

Examples:

```text
could not spawn process
working directory invalid
permission denied before execution
```

Return a proper execution error.

Do not invent command output.

---

# 20. Command Execution Must Happen Exactly Once

This is a strict Phase 6 requirement.

Agent integration must not result in:

```text
agent executes command
+
AgentCap executes command again
```

The command must run once.

AgentCap must wrap the actual execution path.

Duplicate execution is especially dangerous for:

```text
rm
mv
git commit
formatters
generators
database migrations
build scripts with side effects
```

Add tests that detect accidental double execution.

---

# 21. Read-Only vs Side-Effecting Commands

Do not assume shell commands are read-only.

AgentCap must work with:

```text
cat
rg
git diff
```

and also:

```text
mkdir
mv
git add
go generate
npm install
custom scripts
```

Compression behavior may differ, but execution semantics must not.

Never rerun a command automatically just because parsing failed.

---

# 22. Tool Eligibility

The adapter should only intercept command/shell execution tools.

Do not attempt to intercept unrelated agent tools such as:

```text
file editor
browser
web search
MCP database tool
image generation
```

unless explicitly supported in a future phase.

Phase 6 focuses on shell/command execution.

---

# 23. Avoid Double Interception

A command may already explicitly use AgentCap:

```bash
acap run rg Foo .
```

If the agent adapter intercepts this command, it must avoid wrapping it again as:

```text
acap -> acap -> rg
```

Detect AgentCap invocations and bypass or handle them appropriately.

Likewise:

```bash
acap show ...
acap raw ...
acap stats
```

must not be recursively intercepted.

---

# 24. Recursion Protection

Introduce explicit recursion protection.

For example:

```text
ACAP_INTERCEPT_DEPTH
```

or an internal marker.

A child process launched by AgentCap must not accidentally be intercepted again by the same integration path.

Keep recursion prevention simple and deterministic.

---

# 25. Transparent Agent-Facing Output

The coding agent should receive the same compact presentation it would receive from:

```bash
acap run ...
```

Examples:

```text
@acap 81bc2f rg matches=317 files=42
...
```

or:

```text
@acap c81a20 go-test FAIL
failed=3
...
```

or:

```text
@acap f91c20 unchanged from d811aa
exit=0
```

Do not wrap these again in verbose integration prose.

---

# 26. Result IDs Must Remain Visible

Result IDs are critical for progressive disclosure.

The agent-facing result must expose:

```text
@acap <result-id>
```

so the coding agent can later execute:

```bash
acap show <id>
acap raw <id>
```

This behavior must survive agent integration.

---

# 27. Teach the Agent the Drill-Down Contract

The integration should provide a concise instruction or skill file explaining:

```text
When AgentCap returns a compact result:
- use acap show <id> for the full capsule
- use selectors for targeted detail
- use acap raw <id> only when necessary
```

Do not tell the agent to retrieve raw output routinely.

The desired behavior is:

```text
capsule
  -> targeted show
  -> raw only as last resort
```

---

# 28. Agent Instruction Must Be Short

Do not inject a multi-page AgentCap manual into every coding-agent context.

Use a short stable instruction.

For example, the conceptual guidance should fit roughly into:

```text
AgentCap compresses command output.
Results include an ID.
Use `acap show <id> ...` for targeted drill-down.
Use `acap raw <id>` only when the compact result is insufficient.
```

Every token spent explaining the compression tool reduces its benefit.

---

# 29. Avoid Repeated Instruction Injection

If the agent integration mechanism supports persistent project instructions, install AgentCap guidance once.

Do not append the same usage text after every command result.

Phase 6 must measure integration overhead as part of total token savings.

---

# 30. Claude Code Adapter

Implement Claude Code support as an isolated adapter.

Possible responsibilities:

```text
detect shell/tool execution
capture command request
map agent session -> AgentCap session
route execution through AgentCap
return compressed result
```

Do not place Claude-specific event names or payload types outside the adapter package.

Suggested layout:

```text
internal/integration/claude/
```

or an external integration package if cleaner.

---

# 31. Codex Adapter

Implement OpenAI Codex support through its supported command/tool integration mechanism.

Responsibilities should mirror the Claude adapter:

```text
command extraction
cwd extraction
environment handling
session mapping
AgentCap execution
result translation
```

Suggested layout:

```text
internal/integration/codex/
```

Do not share Codex-specific payload structures with core packages.

---

# 32. Common Adapter Interface

Where practical, define a small common interface.

For example:

```go
type Adapter interface {
    Name() string
    Handle(ctx context.Context, req ToolRequest) (ToolResponse, error)
}
```

Or separate parsing and rendering interfaces.

Do not force adapters with fundamentally different event mechanisms into an unnatural abstraction.

Shared semantics matter more than identical plumbing.

---

# 33. Generic Integration Adapter

After at least two real agent adapters work, consider a generic integration mode.

For example:

```text
stdin JSON request
stdout JSON response
```

or another simple process protocol.

This can support agents or wrappers not built directly into AgentCap.

Do not build the generic protocol first.

Prove the requirements using real adapters.

---

# 34. Adapter Configuration

Add minimal explicit configuration.

Possible commands:

```bash
acap integrate claude
acap integrate codex
```

or:

```bash
acap setup claude
acap setup codex
```

Choose one naming scheme and keep it consistent.

The command should configure integration where safe and practical.

---

# 35. Prefer Dry-Run Inspection

Before modifying external agent configuration, support a dry-run or inspect mode.

For example:

```bash
acap integrate claude --dry-run
```

Output:

```text
would add hook:
...
```

This makes integration safer and easier to debug.

---

# 36. Configuration Mutation Safety

Do not overwrite an agent's entire configuration file.

If configuration changes are required:

- preserve unrelated settings
- add only AgentCap-owned entries
- make changes idempotent
- detect conflicts
- create a backup if appropriate
- support uninstall/removal

Do not use fragile regex replacement over structured config formats.

Parse and modify them properly.

---

# 37. Integration Idempotency

Running:

```bash
acap integrate claude
```

twice must not install duplicate hooks.

Likewise:

```bash
acap integrate codex
```

must be idempotent.

Test this.

---

# 38. Uninstall / Disable

Provide a way to remove AgentCap-owned integration configuration.

For example:

```bash
acap integrate claude --remove
acap integrate codex --remove
```

or equivalent.

Only remove entries that AgentCap owns.

Do not remove unrelated hooks.

---

# 39. Integration Status

A command such as:

```bash
acap integrate status
```

may report:

```text
claude configured
codex not configured
```

This is useful but secondary.

Do not build an elaborate integration manager.

---

# 40. Project-Local vs Global Integration

Separate:

```text
AgentCap result storage
```

from:

```text
agent integration configuration
```

Result storage remains project-local:

```text
<project>/.acap/store.db
```

Agent hook configuration may be:

```text
project-local
or
user-global
```

depending on the coding agent's supported configuration model.

Do not force `.acap` to contain third-party agent configuration unless appropriate.

---

# 41. Prefer Project-Local Integration When Practical

For coding-agent tooling that supports per-project hooks/instructions, project-local integration is preferable.

Advantages:

```text
repository-specific opt-in
easy removal
no effect on unrelated projects
```

However, do not create unsupported configuration structures.

Use the agent's documented configuration model.

---

# 42. Global Integration Must Still Respect Project Storage

Even if an adapter is installed globally:

```text
every intercepted command
   |
   v
resolve project root
   |
   v
use that project's .acap/store.db
```

Do not use one global result database merely because the hook itself is global.

---

# 43. Non-Project Commands

If an intercepted command runs outside Git or an existing `.acap` project, preserve Phase 2.1 behavior:

```text
project root = cwd
```

if persistence is needed.

Be cautious about littering arbitrary directories with `.acap`.

If Phase 2.1 already has a policy for transient/non-project directories, reuse it.

Do not invent different storage semantics in the adapter.

---

# 44. AgentCap Self-Commands

Allow coding agents to execute:

```text
acap show
acap raw
acap stats
```

normally.

These commands are control-plane operations and should bypass normal compression interception.

Otherwise the agent may be unable to drill down.

---

# 45. Tool Result Size Limits

Some coding agents impose their own output limits.

AgentCap should produce output well below typical limits where possible.

Adapters should not rely on the agent truncating results.

If AgentCap output itself exceeds an integration API limit:

```text
further reduce safely
or
store the result and return a smaller retrieval pointer
```

Never silently lose all diagnostics.

---

# 46. Integration-Level Emergency Truncation

If an agent API imposes a hard output cap and the normal AgentCap presentation exceeds it:

1. preserve result ID
2. preserve exit code/status
3. preserve highest-priority diagnostics
4. indicate omitted output
5. instruct through the compact format that detail remains available via `acap show`

Example:

```text
@acap 91ab22 clang-build FAIL
errors=184 shown=20 omitted=164

...
detail available via acap show 91ab22 --errors
```

Do not truncate without recoverability.

---

# 47. Stdout and Stderr Mapping

Different agent APIs may model shell results differently.

Where possible preserve:

```text
stdout
stderr
exit code
```

separately.

If an agent only accepts one textual result field, render AgentCap's normal compact presentation there.

Do not fabricate original stdout/stderr ordering.

---

# 48. Streaming Tool Output

Some coding agents may stream command output incrementally.

Phase 6 does not need to support live compressed streaming initially.

A valid initial approach is:

```text
execute
capture
reduce
return final result
```

Do not compromise correctness merely to support streaming.

Live progressive compression can be explored later.

---

# 49. Interactive Commands

Coding agents occasionally invoke commands requiring a TTY.

AgentCap integration must not pretend these are ordinary non-interactive commands if that breaks semantics.

If the adapter can detect interactive/PTY mode:

```text
bypass reduction or use compatible capture path
```

is acceptable.

Document limitations.

---

# 50. Long-Running Commands

Do not impose arbitrary new timeouts.

The agent's cancellation/timeout semantics should remain authoritative.

If the agent cancels a command:

```text
adapter
 -> AgentCap context cancellation
 -> child process termination
```

Use the execution infrastructure from Phase 0.

---

# 51. Cancellation Propagation

Test cancellation explicitly.

If the coding-agent tool request is canceled:

- AgentCap must receive cancellation
- child process should terminate
- result state should not be falsely marked complete
- temporary storage should be cleaned safely

Do not leave orphan processes.

---

# 52. Partial Results on Cancellation

If a canceled command produced partial output, choose a clear policy.

Recommended:

```text
store partial result only if it can be clearly marked incomplete
```

or discard it.

Do not store a partial result as though the command completed normally.

If stored, metadata should indicate:

```text
interrupted=true
```

or equivalent.

---

# 53. Side-Effect Safety Test

Add a test command that increments a file counter:

```text
counter = counter + 1
```

Route it through each adapter.

Verify the counter increments exactly once.

This test protects against one of the most dangerous integration bugs.

---

# 54. Integration Recursion Test

Execute:

```bash
acap run printf hello
```

through an installed agent adapter.

Verify:

```text
only one AgentCap execution layer
```

and no recursion.

Also test:

```bash
acap show <id>
```

through the adapter.

---

# 55. Session Test

Within one simulated agent session:

```text
rg Workspace .
rg Workspace .
```

must map to the same AgentCap session.

The second result should benefit from unchanged/delta logic.

Start a second simulated coding-agent session and verify it does not inherit the first session's stateful presentation.

---

# 56. Project Isolation Test

Simulate one coding agent operating in:

```text
repo-a
```

and another in:

```text
repo-b
```

Verify:

```text
repo-a/.acap/store.db
repo-b/.acap/store.db
```

remain isolated even if the integration adapter itself is globally installed.

---

# 57. CWD Transition Test

Within one agent session:

```text
/repo
/repo/src
/repo/tests
```

commands should still resolve to the same project store when appropriate.

But Phase 3 command identity must continue to include exact cwd.

Do not accidentally normalize cwd to project root for baseline matching.

---

# 58. Agent Instruction / Skill File

Provide an optional concise AgentCap usage instruction suitable for coding agents.

It should explain:

```text
AgentCap may return compact command results with an ID.

Use:
acap show <id>
for expanded stored output.

Use selectors such as:
--file
--hunk
--errors
--warnings
--test
--match
--lines

Use:
acap raw <id>
only when targeted drill-down is insufficient.
```

Keep it concise.

---

# 59. Do Not Encourage `acap raw` by Default

This is important.

The integration instruction should NOT say:

```text
If output is compressed, run acap raw.
```

Instead:

```text
use the smallest relevant drill-down operation
```

Otherwise the coding agent will immediately defeat AgentCap's token-saving model.

---

# 60. Drill-Down Preference

Teach the agent a preference order:

```text
1. trust sufficient capsule
2. targeted acap show selector
3. broader acap show
4. acap raw only if necessary
```

Do not implement this as hard enforcement.

It is behavioral guidance for the coding agent.

---

# 61. Phase 6 Telemetry Policy

Do not add remote telemetry.

Local integration metrics are allowed.

For example:

```text
intercepted commands
bypassed commands
integration failures
raw bytes
returned bytes
show/raw follow-up calls
```

Store them locally if useful.

Do not transmit them anywhere.

---

# 62. Integration Overhead Metrics

Measure:

```text
adapter invocation latency
AgentCap processing latency
total command latency
```

The adapter should add very little overhead beyond existing AgentCap processing.

Avoid:

```text
spawn acap multiple times
parse config every command
scan large directories
```

where possible.

---

# 63. Token-Saving Metrics

Phase 6 should measure the complete integrated workflow.

Important:

```text
compression savings
-
integration instruction overhead
-
adapter wrapper overhead
-
drill-down overhead
=
net context savings
```

Do not claim savings based only on reducer output sizes.

---

# 64. Phase 6 Benchmark

Create realistic scripted agent workflows.

Example:

```text
1. rg Workspace .
2. cat relevant file
3. git status
4. git diff
5. go test ./...
6. edit fixture
7. go test ./...
8. git diff
9. rg Workspace .
```

Run the workflow:

```text
without AgentCap
with explicit acap run
with transparent Phase 6 integration
```

Compare total agent-visible bytes.

---

# 65. Baseline Correctness Benchmark

Verify transparent integration produces the same command-side effects and exit statuses as direct execution.

Compare:

```text
direct
vs
acap run
vs
agent adapter
```

for representative commands.

Any semantic difference is a correctness bug.

---

# 66. Integration Logging

Debug mode may log:

```text
adapter name
agent session ID hash
command argv
cwd
intercept/bypass decision
AgentCap result ID
exit code
presentation type
latency
```

Never log secret environment variables by default.

Debug output must remain local.

---

# 67. Sensitive Environment Data

Do not persist complete environment maps into `.acap/store.db`.

The adapter may need environment variables to execute a process, but AgentCap should not store them by default.

Environment variables may contain:

```text
tokens
API keys
credentials
secrets
```

Persist only specifically needed non-sensitive metadata.

---

# 68. Command Metadata Privacy

Existing result metadata may store argv.

Continue doing so according to current AgentCap semantics.

Be aware that argv itself can occasionally contain secrets.

Do not introduce additional copies of argv in integration logs or config files unnecessarily.

---

# 69. Adapter Crash Isolation

A bug in one adapter must not corrupt AgentCap result storage.

Keep adapter parsing/configuration separate from:

```text
store
executor
reducers
```

Where practical, validate requests before invoking the core.

---

# 70. Version Compatibility

Coding-agent configuration formats can evolve.

Keep adapter code versioned and isolated.

Do not make the core storage schema depend directly on third-party agent configuration versions.

If an integration format is unsupported:

```text
fail that adapter clearly
```

while normal:

```bash
acap run
```

continues to work.

---

# 71. Detect Unsupported Integration Version

If practical, validate the target coding-agent version/configuration format during setup.

Do not install a configuration known to be incompatible.

Avoid brittle assumptions about unspecified third-party internals.

---

# 72. Configuration Backup

When modifying an existing third-party config file:

```text
parse
validate
modify minimally
write atomically
```

Consider creating a backup when the file format or environment makes that appropriate.

Never truncate the original config before a valid replacement is ready.

---

# 73. Atomic Config Writes

Use:

```text
temporary file
validate serialized output
atomic rename
```

where appropriate.

Do not leave half-written coding-agent configuration files.

---

# 74. Preserve Formatting Where Practical

If the third-party config format supports structured parsing but rewriting causes massive unrelated formatting changes, minimize them if practical.

However, correctness is more important than preserving whitespace.

Do not use unsafe text replacement solely to preserve formatting.

---

# 75. Integration Ownership Marker

Where supported, mark AgentCap-owned configuration entries clearly.

For example conceptually:

```text
managed_by = "agentcap"
```

or a unique hook name.

This allows safe idempotent updates and removal.

Do not mark unrelated user configuration as AgentCap-owned.

---

# 76. No MCP Yet Unless Required

Do not add an MCP server merely because Phase 6 integrates coding agents.

Transparent shell interception is the primary objective.

MCP remains optional and separate.

Only introduce MCP if a target coding agent absolutely requires it for the requested integration path.

Otherwise leave it for a later phase.

---

# 77. Why MCP Is Not the Default

AgentCap's main value operates on:

```text
existing shell/tool execution
```

Turning every command into a custom MCP tool would require the model to explicitly choose AgentCap tools.

That weakens transparency.

Prefer:

```text
agent thinks it is running normal shell commands
AgentCap transparently compresses the result
```

where supported.

---

# 78. No AST Integration

Do not add source semantic analysis in Phase 6.

Phase 6 is integration infrastructure only.

No:

```text
Tree-sitter
AST
symbol index
semantic diff
repository analysis
```

unless already provided by prior completed phases.

---

# 79. No New Reducers

Do not expand Phase 6 scope by adding support for dozens of new commands.

If a command is unsupported:

```text
generic reducer
```

already exists.

Phase 6 success depends on transparent adoption, not command breadth.

---

# 80. No Agent-Specific Compression Rules

Do not make the same command compress differently merely because it came from Claude or Codex.

For example:

```text
git diff
```

should use the same core Git reducer regardless of adapter.

Agent-specific differences should be limited to transport/integration constraints.

---

# 81. Adapter Result Limit Handling

If one coding agent permits less output than another, the adapter may request a stricter presentation budget.

Prefer a generic concept such as:

```go
type PresentationBudget struct {
    MaxBytes int
}
```

rather than hard-coded:

```go
if claude ...
if codex ...
```

This keeps presentation policy reusable.

---

# 82. Do Not Optimize for Exact Token Counts Yet

Agent models tokenize differently.

Use:

```text
bytes
lines
rough estimated tokens
```

for integration budgeting.

Do not embed multiple tokenizer libraries merely to target every agent model exactly.

The main reduction gains should be large enough that exact tokenization is unnecessary.

---

# 83. Bypass Mechanism

Provide a way to bypass AgentCap interception explicitly.

For example:

```text
ACAP_BYPASS=1
```

or:

```bash
acap bypass <command...>
```

A simple environment mechanism is preferable.

This is important for debugging integration issues.

---

# 84. Bypass Semantics

When bypass is enabled:

```text
execute command normally
return normal raw result through agent mechanism
```

Do not create a compressed AgentCap result unless explicitly desired.

Avoid creating confusing half-bypassed states.

---

# 85. Per-Command Bypass

The adapter may need to bypass certain commands automatically.

Examples:

```text
acap itself
interactive TTY command unsupported by capture
agent-internal bootstrap commands
```

Keep the automatic bypass list extremely small.

Do not create a giant blacklist.

---

# 86. Installation Discovery

The adapter must locate the `acap` executable reliably if integration uses subprocess execution.

Prefer:

```text
explicit configured path
or
PATH resolution
```

Do not assume:

```text
/usr/local/bin/acap
```

or another hard-coded installation path.

---

# 87. Version Handshake

If adapters invoke an external `acap` process, provide a lightweight version/protocol handshake.

For example:

```text
protocol_version
agentcap_version
```

Do not rely on parsing human-readable `--version` strings to determine protocol compatibility.

---

# 88. Internal Protocol Version

If a machine-readable adapter protocol is introduced, version it from the beginning.

Example:

```json
{
  "protocol": 1
}
```

Keep protocol changes backward-compatible where practical.

Do not conflate:

```text
database schema version
integration protocol version
AgentCap release version
```

These are separate concepts.

---

# 89. Suggested Package Structure

A possible layout:

```text
internal/
├── executor/
├── reduce/
├── result/
├── store/
├── session/
├── delta/
├── integration/
│   ├── protocol/
│   ├── common/
│   ├── claude/
│   └── codex/
├── render/
└── stats/
```

CLI integration setup code may live under:

```text
internal/setup/
```

if that keeps third-party config mutation separate.

Do not force this exact structure if the repository already has cleaner boundaries.

---

# 90. Integration Interface Test Harness

Build a generic adapter test harness.

It should be able to simulate:

```text
agent session
tool command
cwd
env
cancellation
expected exit code
expected returned text
```

Run the same semantic tests against each adapter where possible.

This prevents behavioral divergence.

---

# 91. Golden Tests

For adapter request/response translation, use deterministic golden fixtures where appropriate.

Test:

```text
normal success
normal failure
large compressed result
unchanged result
Git delta
test failure delta
AgentCap internal fallback
```

Do not make tests depend on live coding-agent services.

---

# 92. Do Not Require Network Access in Tests

Phase 6 integration tests should run locally.

Do not require:

```text
Claude API
OpenAI API
Gemini API
```

to execute the test suite.

Test hook/config/protocol behavior with fixtures and local subprocesses.

---

# 93. End-to-End Manual Tests

Actual installed coding agents may be used for manual verification, but not as the only correctness test.

The repository test suite must remain deterministic.

---

# 94. Claude Integration Verification

Manually verify:

```text
agent executes rg
agent receives AgentCap capsule

agent executes git diff
agent receives Git capsule

agent executes go test
agent receives diagnostic/test capsule

agent requests acap show
drill-down works

agent repeats unchanged command
session-aware compression works
```

Verify the actual command executes only once.

---

# 95. Codex Integration Verification

Perform the same semantic verification for Codex.

Do not accept adapter behavior that differs materially without a documented platform limitation.

---

# 96. Cross-Agent Consistency

Given:

```text
same repository
same command
same repository state
```

the core AgentCap capsule should be substantially identical regardless of which supported coding agent issued it.

Transport wrappers may differ.

Compression semantics should not.

---

# 97. Integration Setup Documentation

Add concise documentation for each supported coding agent.

Include:

```text
requirements
installation/setup command
how to verify
how to disable
how to remove
how to debug
```

Do not duplicate the entire AgentCap README for each adapter.

---

# 98. Diagnostic Command

Consider:

```bash
acap doctor
```

or:

```bash
acap integrate doctor
```

to check:

```text
acap binary available
project root detected
.acap writable
SQLite store usable
supported integration config detected
hook installed
```

This is useful but secondary.

Do not turn it into a large environment scanner.

---

# 99. `acap doctor` Must Be Read-Only by Default

If implemented, it should not silently repair or modify configuration.

It may recommend:

```text
run acap integrate claude
```

but should not do so automatically.

---

# 100. Phase 6 Statistics

Add local integration-level statistics such as:

```text
commands intercepted
commands bypassed
integration failures
raw bytes
stateless capsule bytes
stateful returned bytes
drill-down bytes
raw retrieval bytes
```

Optionally split by adapter:

```text
claude
codex
```

Do not expose user prompts or source contents in aggregate stats.

---

# 101. Net Savings

The critical Phase 6 metric is:

```text
agent-visible bytes without AgentCap
-
agent-visible bytes with transparent AgentCap
```

across a complete realistic workflow.

Also measure:

```text
number of extra drill-down commands
```

A system that saves 95% initially but causes excessive `acap raw` calls is not successful.

---

# 102. Integration Latency Budget

Measure:

```text
direct command latency
acap run latency
integrated agent latency
```

The Phase 6 adapter overhead should be small relative to existing AgentCap processing.

Do not add expensive per-command configuration discovery.

Cache safe immutable integration metadata within a long-lived adapter process when available.

---

# 103. No Background Daemon by Default

Do not introduce a required AgentCap daemon in Phase 6.

Prefer:

```text
hook/subprocess
or
in-process adapter
```

depending on the coding agent.

A daemon adds:

```text
lifecycle complexity
socket security
version coordination
installation complexity
```

Only introduce one if measurements prove it necessary.

---

# 104. No Remote Service

AgentCap remains local.

Phase 6 must not add:

```text
cloud relay
hosted compression API
remote command execution
remote result storage
telemetry backend
```

The coding agent may itself be remote, but AgentCap command interception and storage remain local to the execution environment.

---

# 105. Explicit Non-Goals

Do NOT implement the following in Phase 6 unless strictly required for the two primary adapters:

- MCP server
- AST analysis
- Tree-sitter
- semantic code indexing
- repository-wide symbol database
- new Git features
- new compiler/test reducers
- remote storage
- telemetry
- cloud synchronization
- background daemon
- shell implementation
- command rewriting
- autonomous command selection
- automatic code modifications
- automatic test reruns
- agent prompt rewriting beyond minimal AgentCap usage guidance

Keep Phase 6 focused on integration.

---

# 106. Suggested Implementation Order

Implement in this order:

1. define common integration request/response model
2. define recursion/bypass semantics
3. expose existing AgentCap execution pipeline through a reusable internal API
4. ensure the core can be invoked without CLI parsing
5. build an adapter test harness
6. implement stable session mapping
7. implement the first primary agent adapter
8. verify command executes exactly once
9. verify result IDs and drill-down
10. verify cancellation and cwd/environment preservation
11. add integration setup/removal support
12. add idempotency/config safety tests
13. implement the second primary agent adapter
14. run cross-agent semantic consistency tests
15. add compact agent usage instruction
16. add integration-level statistics
17. add realistic end-to-end workflow benchmarks
18. add bypass/debug/doctor support where justified
19. update documentation
20. only then evaluate secondary adapters

Do not build all adapters simultaneously.

Make one adapter reliable first, then reuse the proven architecture.

---

# 107. Definition of Done

Phase 6 is complete when all of the following are true:

1. AgentCap core remains agent-independent.

2. A reusable integration request/response boundary exists.

3. At least Claude Code and Codex have working isolated adapters.

4. Supported shell/tool commands can be routed transparently through AgentCap.

5. The coding agent does not need to manually prefix normal commands with `acap run`.

6. Each intercepted command executes exactly once.

7. Exact argv semantics are preserved.

8. Exact cwd semantics are preserved.

9. Required environment semantics are preserved.

10. Child exit codes remain authoritative.

11. Command failures are distinguished from AgentCap integration failures.

12. Session IDs are stable within one coding-agent task.

13. Separate coding-agent tasks do not accidentally share stateful presentation history.

14. Project-local `.acap/store.db` isolation remains intact.

15. Existing Phase 3 unchanged/delta behavior works through adapters.

16. Phase 4 Git capsules and Git deltas work through adapters.

17. Phase 5 build/test capsules and deltas work through adapters.

18. Result IDs remain visible to the coding agent.

19. `acap show` and `acap raw` remain usable from the coding agent.

20. AgentCap commands do not recursively intercept themselves.

21. A bypass mechanism exists for debugging.

22. Cancellation propagates correctly to child processes.

23. Interactive/unsupported command behavior has a safe documented fallback.

24. Adapter configuration installation is idempotent.

25. Adapter removal does not delete unrelated user configuration.

26. External config modifications are performed safely and atomically where practical.

27. Integration tests do not require remote APIs.

28. Side-effect tests confirm commands execute exactly once.

29. Cross-project isolation tests pass.

30. Cross-agent behavior is semantically consistent.

31. Integration-level benchmarks measure net agent-visible output savings.

32. Adapter overhead is measured.

33. No remote telemetry is added.

34. No mandatory background daemon is added.

35. `go test ./...` passes.

36. `go build ./...` passes.

---

# 108. Engineering Priorities

When tradeoffs are necessary, use this order:

1. Never execute a command twice
2. Preserve command semantics
3. Preserve exit/cancellation behavior
4. Never hide new actionable output
5. Fail open to useful command output when compression fails
6. Maintain project/session isolation
7. Preserve drill-down and raw recoverability
8. Keep adapters isolated from core
9. Reduce total agent-visible context
10. Minimize integration latency
11. Add more agent integrations

Transparent integration is valuable only if it is trustworthy.

---

# 109. Phase 6 Success Criterion

Phase 6 succeeds when a coding agent can work normally:

```text
rg ...
cat ...
git diff
edit
go test ./...
edit
go test ./...
git diff
```

without being instructed to manually wrap each command in AgentCap.

The effective flow should be:

```text
Coding Agent
    |
    | normal shell command
    v
AgentCap interception
    |
    v
real execution exactly once
    |
    v
local raw capture
    |
    v
stateful structured compression
    |
    v
compact tool result
    |
    v
Coding Agent
```

And when more detail is required:

```text
agent
  |
  +--> acap show <id> --file ...
  +--> acap show <id> --errors
  +--> acap show <id> --test ...
  |
  +--> acap raw <id> only when necessary
```

The coding agent should gain AgentCap's token-saving behavior without substantially changing how it normally uses shell tools.

---

# 110. Final Verification

Before finishing Phase 6, perform an end-to-end workflow using each primary coding-agent adapter.

Use a repository with:

- searchable source
- Git modifications
- failing tests
- successful tests after fixes

Use one stable AgentCap session.

Have the agent execute normal commands such as:

```bash
rg "Workspace" .
git status
git diff
go test ./...
```

Verify the agent did NOT explicitly need:

```bash
acap run ...
```

Verify:

```text
rg
 -> Phase 1/3 capsule

git diff
 -> Phase 4 capsule

go test
 -> Phase 5 failure capsule
```

Then modify repository state and repeat.

Verify:

```text
repeated search
 -> unchanged/delta

repeated git diff
 -> Git structural delta

repeated tests
 -> resolved/new/remaining failure delta
```

Have the agent request targeted detail:

```bash
acap show <id> --file src/foo.go
acap show <id> --hunk 2
acap show <id> --errors
acap show <id> --test TestFoo
```

Verify no command is re-executed.

Verify:

```bash
acap raw <id>
```

still returns the stored raw result.

Run a side-effecting command and confirm it executes exactly once.

Cancel a long-running command and confirm the child terminates correctly.

Start a second coding-agent session and verify it does not inherit the first session's delta state.

Run the same workflow in a second repository and verify storage isolation.

Test integration removal and confirm unrelated agent configuration remains intact.

Finally run:

```bash
gofmt
go test ./...
go build ./...
```

Report:

- supported coding-agent adapters
- integration architecture
- internal integration protocol
- session mapping rules
- recursion/bypass strategy
- cancellation behavior
- configuration/setup strategy
- known agent-specific limitations
- direct vs integrated latency
- total workflow output savings
- drill-down frequency
- known limitations

Stop once the Phase 6 Definition of Done is satisfied.
