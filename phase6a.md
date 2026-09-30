# Phase 6.5A — Common Agent Adapter Contract

## Objective

Introduce and stabilize a common adapter contract for integrating AgentCap with multiple coding agents.

This phase must define the shared integration boundary before implementing or modifying agent-specific adapters such as Claude Code or Antigravity.

The goal is to ensure that all future adapters use the same execution, session, safety, and result-delivery semantics while keeping agent-specific behavior isolated from the AgentCap core.

Do not implement Claude Code-specific or Antigravity-specific behavior in this phase.

---

## Background

AgentCap already provides:

- command execution
- raw stdout/stderr capture
- result storage
- deterministic reduction
- progressive disclosure
- session-aware compression
- baseline selection
- delta presentation
- Git-aware reduction
- compiler/build/test-aware reduction
- transparent agent integration infrastructure

The intended architecture is:

```text
Coding Agent
    |
    v
Agent Adapter
    |
    v
AgentCap Integration Boundary
    |
    v
AgentCap Core
```

Agent-specific behavior must remain in adapters.

The AgentCap core must not accumulate logic such as:

```text
if claude ...
if antigravity ...
if codex ...
```

Phase 6.5A establishes the contract that all adapters must follow.

---

# 1. Core Design Principle

The adapter layer should be thin.

Its responsibilities are limited to translating between:

```text
agent-specific invocation/result format
```

and:

```text
AgentCap's agent-independent integration model
```

The adapter must not duplicate AgentCap core functionality.

In particular, adapters must not independently implement:

- command execution semantics
- result reduction
- Git parsing
- compiler/test parsing
- result persistence
- baseline selection
- delta computation
- raw-result storage
- progressive-disclosure logic

Those responsibilities remain in AgentCap core.

---

# 2. Critical Execution Invariant

The most important invariant is:

```text
Every intercepted command must execute exactly once.
```

The integration must never result in:

```text
agent executes command
+
AgentCap executes command again
```

This applies to all commands, including commands with side effects.

Examples include:

```text
rm
mv
cp
git add
git commit
git checkout
go generate
database migration commands
custom scripts
package installation commands
```

The common adapter contract must make it difficult or impossible for an adapter implementation to accidentally violate this rule.

Add tests that explicitly verify execution count.

---

# 3. Define a Common Adapter Boundary

Introduce or refine an agent-independent adapter boundary.

Use existing project naming and package conventions where possible.

Do not introduce unnecessary abstraction if an equivalent integration boundary already exists.

Conceptually, the common model should cover:

```text
AgentInvocation
    |
    +-- agent identity
    +-- agent session identity
    +-- tool identity
    +-- command
    +-- argv / shell representation
    +-- cwd
    +-- environment context when applicable
    +-- cancellation context
    +-- adapter metadata
```

and:

```text
AgentResult
    |
    +-- result ID
    +-- stdout / presentation output
    +-- stderr when required
    +-- exit status
    +-- AgentCap metadata
    +-- execution / reduction status
```

Exact type names are implementation-dependent.

Prefer small, explicit structures over a large generic framework.

---

# 4. Agent Identity

The common contract must provide an explicit agent identity.

Examples of future values may include:

```text
claude-code
antigravity
codex
```

Do not hard-code behavior based on these identities inside the core execution or reduction path.

Agent identity is primarily metadata for:

- diagnostics
- metrics
- benchmark attribution
- installation/configuration state
- adapter-specific handling outside the core

---

# 5. Session Mapping Contract

A coding-agent task/session must map to an AgentCap session.

Define the contract for providing a stable external session identifier to AgentCap.

The adapter should conceptually provide:

```text
agent
external session ID
project root
```

from which AgentCap can obtain or reuse the corresponding internal session.

Required behavior:

```text
same agent task/session
    -> same AgentCap session

different agent task/session
    -> different AgentCap session
```

Do not allow accidental state sharing between unrelated agent tasks.

Do not mix sessions across projects.

The project-local storage model must remain intact.

---

# 6. Missing Session IDs

Some integrations may not always provide a reliable agent session identifier.

Define conservative behavior for this case.

The fallback must not create unsafe cross-task deduplication.

Prefer reduced statefulness over incorrectly sharing state between unrelated agent sessions.

Document the fallback semantics clearly.

Do not invent unstable identifiers from values that may collide frequently.

---

# 7. Command Representation

The common boundary must preserve command semantics.

AgentCap must not reinterpret the user's command unnecessarily.

Where an agent provides direct argv execution, preserve argv.

Where an agent exposes a shell command string, preserve the original shell semantics and route it through the appropriate existing execution path.

Do not convert:

```text
shell command string
```

into:

```text
naively split argv
```

if doing so would change shell behavior.

Likewise, do not reconstruct shell syntax from argv unless explicitly required.

Preserve the distinction between:

```text
direct argv execution
```

and:

```text
shell-mediated execution
```

if AgentCap currently supports both concepts.

---

# 8. Working Directory

The adapter contract must explicitly preserve execution cwd.

Remember that:

```text
project root != execution cwd
```

A command invoked from:

```text
/repo/src/parser
```

must execute from that directory even if AgentCap storage is located at:

```text
/repo/.acap/
```

Add tests covering nested working directories.

---

# 9. Environment Handling

Preserve the existing environment semantics.

The adapter must not silently remove, rewrite, or globally mutate environment variables required by the command.

AgentCap-specific integration variables may be added where necessary for:

- recursion protection
- session propagation
- debugging
- adapter identification

Keep such variables narrowly scoped and documented.

Avoid using global persistent environment state when per-process state is sufficient.

---

# 10. Exit Status Preservation

The command's real exit status must remain observable by the coding agent.

For example:

```text
command exit = 0
-> adapter reports success

command exit = 1
-> adapter preserves exit = 1
```

Compression must never turn command failure into command success.

Likewise, an internal reduction failure must not replace the actual command exit status.

---

# 11. stdout / stderr Semantics

Define the adapter contract for command output carefully.

The core should continue capturing complete raw stdout and stderr.

The agent-facing response may contain AgentCap's reduced presentation.

However:

- stderr must not be discarded unsafely
- diagnostic information must remain recoverable
- raw output must remain available through stored results
- result IDs must remain available for drill-down

Do not make adapter-specific compression decisions.

Adapters should present the result produced by the core integration layer.

---

# 12. Progressive Disclosure Compatibility

All adapters must preserve AgentCap's existing progressive-disclosure workflow.

An agent receiving:

```text
@acap <result-id> ...
```

must still be able to use existing retrieval mechanisms such as:

```bash
acap show <result-id>
acap show <result-id> --file ...
acap show <result-id> --errors
acap show <result-id> --warnings
acap show <result-id> --test ...
acap raw <result-id>
```

Do not create adapter-specific result stores.

Do not create adapter-specific drill-down formats unless unavoidable.

---

# 13. Recursion Protection

Define common recursion protection at the integration boundary.

AgentCap must not intercept itself recursively.

Examples that should not recursively re-enter the adapter:

```bash
acap show ...
acap raw ...
acap stats
acap clean
```

Likewise, internal helper processes launched by the adapter or AgentCap must not accidentally pass through the integration layer repeatedly.

Use an explicit and testable recursion/bypass mechanism.

Do not depend on fragile command-name heuristics alone if a stronger mechanism already exists.

---

# 14. Explicit Bypass

Provide a shared mechanism for intentionally bypassing AgentCap interception.

This is required for:

- debugging
- integration recovery
- compatibility testing
- benchmark control groups
- diagnosing adapter failures

The bypass mechanism should behave consistently across adapters where possible.

Examples of acceptable approaches include an existing environment flag or integration configuration.

Prefer one common mechanism over separate flags for each adapter.

---

# 15. Fail-Open Behavior

Preserve the existing integration safety principle:

```text
If command execution succeeds but AgentCap reduction/presentation fails,
prefer returning usable command output over breaking the coding agent.
```

The common contract must distinguish between at least:

```text
command execution failure
```

and:

```text
AgentCap processing/integration failure
```

Do not hide a real command failure.

Do not report a command as failed solely because optional compression failed if complete command output can still be safely returned.

Any fail-open behavior must remain deterministic and testable.

---

# 16. Cancellation and Signals

The common boundary must preserve cancellation semantics.

If the coding agent cancels a running command:

```text
agent cancellation
    ->
adapter
    ->
AgentCap
    ->
child process
```

Cancellation must propagate to the actual command.

Avoid orphaned child processes.

Preserve the existing AgentCap signal-handling behavior.

Add at least one integration-level test where practical.

---

# 17. Result Metadata

Add enough adapter metadata to support future benchmarking and debugging.

At minimum, result or integration metadata should make it possible to determine:

```text
agent
adapter
adapter version if applicable
session mapping
intercepted vs bypassed
integration success/failure
adapter processing latency
```

Do not overload normal agent-facing output with this metadata.

Persist or expose it through existing statistics/debug facilities where appropriate.

---

# 18. Adapter Metrics

Prepare the common integration layer for the upcoming benchmark phase.

Support measurement of:

```text
commands intercepted
commands bypassed
adapter failures
adapter latency
AgentCap processing latency
```

Keep:

```text
adapter latency
```

separate from:

```text
core AgentCap processing latency
```

where feasible.

This distinction is important because benchmark results must be able to distinguish:

```text
cost of AgentCap
```

from:

```text
cost of a particular agent integration
```

Do not build the full benchmark framework in this phase.

Only expose the necessary adapter-level measurements.

---

# 19. Adapter Versioning

Provide a lightweight way to identify adapter implementation/version in diagnostics or benchmark records.

Do not create a complex compatibility negotiation protocol.

The purpose is simply to make benchmark and bug reports attributable to a concrete adapter implementation.

If repository/build version information already exists and is sufficient, reuse it.

---

# 20. Installation Is Out of Scope

Phase 6.5A defines runtime contracts.

Do not implement agent-specific installation workflows here.

The following belong to later phases:

```text
Claude Code hook installation
Claude Code hook configuration

Antigravity .agents/hooks.json generation
Antigravity plugin packaging
```

This phase may define shared installer-facing interfaces only if clearly necessary.

Do not build installers prematurely.

---

# 21. Agent-Specific Parsing Is Out of Scope

Do not add parsers for:

```text
Claude hook payloads
Antigravity hook payloads
```

in the common package.

Those belong to their respective adapters.

The common layer should receive already-normalized invocation data.

Preferred flow:

```text
agent-specific payload
    |
    v
agent adapter
    |
    v
normalized AgentInvocation
    |
    v
AgentCap integration boundary
```

---

# 22. Concurrency

The common contract must tolerate concurrent agent tool executions if the surrounding coding agent permits them.

Ensure that:

- result storage remains safe
- session association remains correct
- one command does not inherit another command's transient adapter state
- recursion/bypass state is process-safe
- adapter metrics do not race

Do not introduce global mutable state unless necessary and synchronized.

---

# 23. Security and Trust Boundary

Treat agent-provided command input as command input, not configuration.

Do not evaluate or reinterpret command strings inside configuration templates.

Avoid introducing command-injection opportunities through adapter metadata.

Do not construct shell commands from unescaped metadata.

Preserve AgentCap's existing local-only model.

Do not add:

- telemetry
- remote services
- network dependencies

---

# 24. Common Conformance Tests

Create a shared adapter conformance test suite or equivalent reusable test helpers.

Future adapters should be able to run the same behavioral tests.

At minimum, cover:

### Basic execution

```text
successful command
failing command
stdout-only command
stderr-only command
stdout + stderr command
```

### Execution count

Verify that an intercepted command executes exactly once.

Use a side-effect fixture such as incrementing/writing a local marker.

Expected result:

```text
marker changed once
not twice
```

### Exit status

Verify multiple non-zero exit codes.

### cwd

Execute from nested directories and verify the actual working directory.

### Environment

Verify normal environment propagation.

### Session mapping

Verify:

```text
same external session
    -> same internal session

different external session
    -> different internal session
```

### Project isolation

Verify that identical session identifiers in separate repositories do not share AgentCap project state.

### Recursion

Verify that AgentCap retrieval commands do not recursively re-enter interception.

### Bypass

Verify that explicit bypass executes normally without AgentCap processing.

### Fail-open

Simulate AgentCap presentation/reduction failure where practical and verify usable command output remains available.

### Cancellation

Verify cancellation propagation if the existing test infrastructure makes this reliable.

---

# 25. Cross-Adapter Consistency Contract

Future Claude Code and Antigravity adapters must produce equivalent normalized requests for equivalent command executions.

For an equivalent command such as:

```bash
go test ./...
```

the AgentCap core should see equivalent information for:

```text
command semantics
cwd
project
session
environment behavior
```

The resulting AgentCap behavior should therefore be equivalent with respect to:

```text
classification
raw capture
storage
reduction
delta handling
result IDs
drill-down support
exit semantics
```

Only agent-specific transport details should differ.

---

# 26. Keep the API Small

Do not create a large plugin SDK.

This phase is intended to support a small number of coding-agent adapters.

Prefer:

```text
small normalized request/result contract
+
shared execution integration helpers
+
shared conformance tests
```

over:

```text
generic extensible agent framework
```

Avoid speculative abstractions for agents that are not currently being integrated.

---

# 27. Backward Compatibility

Existing supported AgentCap integration behavior must continue to work.

Do not break:

```text
acap run
acap show
acap raw
existing storage
existing sessions
existing reducers
existing result IDs
existing CLI behavior
```

Any migration of existing Phase 6 integration code into the new adapter boundary should preserve observable behavior.

Refactor incrementally.

---

# 28. Code Organization

Keep the architecture visibly separated.

A reasonable conceptual structure is:

```text
internal/
    integration/
        common agent-independent integration logic

    adapters/
        shared adapter helpers if necessary

        claude/
            # later phase

        antigravity/
            # later phase
```

Do not force this exact directory layout if the repository already has a better established organization.

The architectural rule matters more than directory names:

```text
agent-specific transport
    !=
AgentCap core
```

---

# 29. Documentation

Add concise developer documentation describing:

```text
what an AgentCap adapter is
what it is allowed to do
what it must not do
how session mapping works
how commands are executed exactly once
how bypass/recursion protection works
how failures are handled
how a new adapter runs conformance tests
```

Keep documentation implementation-oriented.

Do not write user-facing Claude or Antigravity setup instructions yet.

---

# 30. Non-Goals

Do not implement:

- Claude Code adapter behavior
- Antigravity adapter behavior
- Gemini CLI integration
- OpenCode integration
- new reducers
- new Git parsing
- new compiler/test parsers
- semantic source-code processing
- AST integration
- benchmark task orchestration
- benchmark scoring
- automatic adapter discovery
- generic third-party plugin SDK
- remote telemetry
- daemon infrastructure

This phase is strictly about stabilizing the common adapter contract.

---

# 31. Acceptance Criteria

Phase 6.5A is complete when all of the following are true.

### Architecture

- There is a clear agent-independent integration boundary.
- Agent-specific transport behavior is outside the AgentCap core.
- The API is small and implementation-oriented.

### Execution correctness

- Every intercepted command executes exactly once.
- cwd is preserved.
- environment semantics are preserved.
- exit status is preserved.
- stdout/stderr semantics remain correct.
- cancellation propagates correctly where supported.

### Session correctness

- Adapter sessions map deterministically to AgentCap sessions.
- Different agent sessions remain isolated.
- Projects remain isolated.
- Missing external session IDs use a conservative fallback.

### Safety

- recursion protection exists and is tested.
- explicit bypass exists and is tested.
- fail-open behavior is defined and tested.
- existing local-only behavior is preserved.

### Progressive disclosure

- normal AgentCap result IDs remain available.
- `acap show` continues to work.
- `acap raw` continues to work.
- adapters do not create parallel result stores.

### Benchmark readiness

The integration layer can expose or record:

```text
agent
adapter
adapter version
commands intercepted
commands bypassed
adapter failures
adapter latency
```

without implementing the benchmark itself.

### Tests

Shared conformance tests cover at least:

```text
success
failure
stdout/stderr
exit status
cwd
environment
execute exactly once
session mapping
project isolation
recursion protection
bypass
fail-open
```

All existing AgentCap tests continue to pass.

---

# 32. Implementation Guidance

Before changing code:

1. Inspect the existing Phase 6 integration implementation.
2. Identify which parts are already agent-independent.
3. Identify any existing agent-specific assumptions embedded in shared code.
4. Preserve working code where possible.
5. Extract the smallest useful common adapter contract.

Do not rewrite the integration subsystem solely for architectural cleanliness.

Prefer incremental refactoring backed by tests.

When there is a choice between:

```text
more abstraction
```

and:

```text
a smaller explicit contract sufficient for Claude Code + Antigravity
```

choose the smaller explicit contract.

---

# 33. Final Verification

Before declaring the phase complete, verify the following manually or through automated tests:

```text
1. A normal command executes once.
2. A side-effect command executes once.
3. A failing command preserves its exit code.
4. stdout and stderr remain recoverable.
5. A command from a nested cwd executes from that cwd.
6. Two invocations in one external agent session share AgentCap session state.
7. Two external agent sessions do not share session state.
8. Two repositories do not share project state.
9. AgentCap's own retrieval commands do not recurse.
10. Explicit bypass works.
11. A simulated AgentCap processing failure fails open safely.
12. Existing acap run/show/raw behavior is unchanged.
13. Existing tests pass.
14. The new conformance suite passes.
```

Do not begin Claude Code- or Antigravity-specific adapter work as part of this phase.

The output of Phase 6.5A should be a stable integration contract that Phase 6.5B and Phase 6.5C can implement independently.
