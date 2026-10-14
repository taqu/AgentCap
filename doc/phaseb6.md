# Phase B6 — Real Coding-Agent Workflow Benchmark

## Objective

Extend the benchmark system from deterministic AgentCap-only workloads to **real coding-agent workflow evaluation**.

Phase B0–B5 established deterministic measurement of:

```text
raw command cost
stateless AgentCap cost
stateful AgentCap cost
progressive-disclosure recovery cost
total agent-visible cost
```

Phase B6 must now answer the next question:

```text
Does AgentCap reduce the total context/output cost
of a real coding-agent workflow
without preventing the agent from completing the task?
```

This phase introduces:

- real coding-agent execution;
- an agent-adapter abstraction;
- coding-task workloads;
- AgentCap comparison modes;
- task success measurement;
- complete workflow measurement.

Do not implement statistical repeated-trial analysis yet. That belongs to Phase B7.

---

# 1. Preserve the Existing Benchmark Foundation

Build B6 on top of B0–B5.

Do not create a separate measurement system for agent benchmarks.

The architecture should conceptually become:

```text
Benchmark Runner
      |
      +---- workload
      |
      +---- agent adapter
      |        |
      |        +-- Claude Code
      |        +-- Codex
      |        +-- future agents
      |
      +---- AgentCap mode
      |
      +---- B0–B5 measurement
      |
      +---- task verifier
      |
      +---- BenchmarkResult
```

Reuse the existing:

```text
measurement model
workload infrastructure
session/stateful infrastructure
recovery accounting
result schema
human formatter
JSON formatter
```

wherever possible.

B6 should extend the benchmark system rather than fork it.

---

# 2. Introduce an Agent Adapter Abstraction

Agent-specific process invocation and configuration must be isolated behind an adapter interface.

The benchmark runner must not contain logic such as:

```text
if agent == "claude" ...
if agent == "codex" ...
```

throughout the execution path.

Use an abstraction appropriate for the existing Go architecture.

Conceptually:

```go
type AgentAdapter interface {
    Name() string

    Run(ctx context.Context, req AgentRunRequest) (AgentRunResult, error)
}
```

The exact interface should be derived from actual requirements and existing repository conventions.

Do not copy this interface literally if another shape fits better.

The important dependency direction is:

```text
benchmark runner
      |
      v
AgentAdapter
      |
      +--> Claude Code implementation
      +--> Codex implementation
```

Agent-specific behavior must remain below the adapter boundary.

---

# 3. Keep the Adapter Small

The adapter should be responsible only for agent-specific concerns such as:

```text
process invocation
CLI arguments
prompt/task injection
environment setup
working-directory setup
AgentCap integration hooks
exit/result capture
```

It should not own:

```text
benchmark metric aggregation
AgentCap measurement semantics
task-success policy
result formatting
statistical analysis
workload fixture management
```

Those remain benchmark-system responsibilities.

Avoid designing a large generic agent framework.

B6 only needs enough abstraction to run coding agents reproducibly.

---

# 4. Start With the Minimum Useful Adapter Set

Implement only the agent adapters that can be supported and tested reliably in the current development environment.

The roadmap identifies:

```text
Claude Code
Codex
```

as intended adapter examples.

Do not block the architecture on implementing every possible coding agent.

If only one adapter can be reliably exercised during B6 implementation, establish the adapter abstraction and implement that adapter first.

The design must make adding another adapter straightforward.

Do not introduce speculative adapters for tools that are not being tested.

---

# 5. Introduce Coding-Agent Workloads

B6 workloads represent actual coding tasks rather than only command sequences.

Target workload categories include:

```text
bug fixing
test failure repair
small feature
refactoring
compile-error repair
Git review
repository search
```

Do not require every category to be fully populated during the first implementation.

Implement enough representative fixtures to validate the architecture.

Prefer small, deterministic repositories with objective completion criteria.

---

# 6. Separate Task Setup, Agent Instruction, and Verification

A coding-agent workload should conceptually contain three distinct parts:

```text
setup
task
verification
```

For example:

```yaml
name: go-fix-001

fixture: fixtures/go/fix-001

task: |
  Fix the bug causing TestParseConfig to fail.
  Do not modify the test.

verify:
  - run: go test ./...
```

The exact workload syntax must follow or extend the B3 workload model.

Do not introduce a completely independent workload format unless the existing format fundamentally cannot represent agent tasks.

The separation matters because:

```text
task text
```

is what the coding agent sees, while:

```text
verification
```

is benchmark infrastructure and must not automatically be exposed to the agent.

---

# 7. Use Fresh Isolated Workspaces

Every agent benchmark run must execute in a fresh fixture workspace.

Conceptually:

```text
canonical fixture
      |
      v
temporary workspace
      |
      v
agent run
      |
      v
verification
      |
      v
discard workspace
```

Never let one benchmark run contaminate another.

Do not run coding-agent benchmarks directly against the AgentCap development repository unless a workload explicitly uses a disposable copy.

This becomes especially important in B7 when repeated trials are introduced.

---

# 8. Define the Four Comparison Modes

B6 introduces four formal benchmark modes:

```text
disabled
stateless
stateful
integrated
```

Their semantics must be explicitly defined in code and documentation.

## A. `disabled`

AgentCap is not used for command-output reduction.

The agent sees the normal command output.

This is the baseline.

## B. `stateless`

AgentCap reduces command output independently for each command.

No cross-command/session delta behavior is used.

## C. `stateful`

AgentCap uses the stateful/session behavior implemented in earlier phases.

Repeated or evolving command results may use unchanged/delta/full presentation behavior.

## D. `integrated`

AgentCap is exposed to the coding agent as the intended complete workflow integration.

This includes the normal progressive-disclosure mechanisms available to the agent, such as retrieving additional information through AgentCap when necessary.

The exact distinction between `stateful` and `integrated` must reflect the actual AgentCap architecture.

Do not invent artificial differences merely to preserve four names.

Before implementation, inspect the existing integration path and document the precise semantics of each mode.

---

# 9. Represent Modes Explicitly

Use a typed mode representation rather than scattered strings.

Conceptually:

```go
type AgentBenchmarkMode string

const (
    BenchmarkModeDisabled   AgentBenchmarkMode = "disabled"
    BenchmarkModeStateless  AgentBenchmarkMode = "stateless"
    BenchmarkModeStateful   AgentBenchmarkMode = "stateful"
    BenchmarkModeIntegrated AgentBenchmarkMode = "integrated"
)
```

Adapt naming to repository conventions.

Centralize validation and behavior selection.

Avoid mode-specific conditionals spread throughout unrelated packages.

---

# 10. Add an Agent Benchmark Command

Add a user-facing entry point conceptually similar to:

```bash
acap bench agent \
    --workload bugfix-001 \
    --agent <agent> \
    --mode disabled
```

and:

```bash
acap bench agent \
    --workload bugfix-001 \
    --agent <agent> \
    --mode integrated
```

Exact syntax should follow the existing B0–B5 CLI conventions.

Do not preserve roadmap example syntax blindly if the implemented benchmark CLI has evolved.

The command should make these dimensions explicit:

```text
workload
agent
AgentCap mode
```

Avoid hidden defaults that make benchmark results difficult to reproduce.

Reasonable defaults are acceptable if they are clearly represented in the resulting benchmark metadata.

---

# 11. Measure the Complete Agent Workflow

The benchmark boundary is the complete coding-agent task, not one command.

Measure all relevant command output presented to the agent during the run.

Conceptually:

```text
agent starts
   |
   +--> inspect files
   +--> search repository
   +--> run tests
   +--> edit
   +--> run tests
   +--> acap show
   +--> edit
   +--> run tests
   |
agent finishes
```

The benchmark should aggregate AgentCap-visible command/recovery output across this complete workflow.

Continue using the B5 principle:

```text
total agent-visible cost
=
initial command presentations
+
show recovery
+
raw recovery
```

for AgentCap-enabled modes.

For disabled mode, measure the corresponding raw command output visible to the agent.

---

# 12. Measure AgentCap-Relevant Output, Not Arbitrary Agent Tokens

Do not turn B6 into a general LLM token-accounting system unless the existing agent integration already exposes reliable token usage.

The primary KPI remains:

```text
Total agent-visible tokens/bytes
over the complete coding workflow
```

with B0–B5 currently providing byte-oriented measurement of AgentCap-controlled command output.

Be explicit about the measurement boundary.

For example, unless reliably available, do not pretend to measure:

```text
model internal reasoning tokens
provider-side cached tokens
agent framework hidden prompts
```

If an adapter exposes trustworthy token metadata, it may be recorded as optional metadata, but it must not redefine the B6 core measurement model.

---

# 13. Capture Command Count

Record the number of benchmark-relevant commands executed during the agent workflow.

This matters because AgentCap must not appear more efficient merely because the agent behaves very differently.

At minimum expose:

```text
command_count
```

using the existing measurement semantics where possible.

Later comparison should make it possible to observe situations such as:

```text
disabled:
    20 commands
    300 KB visible

integrated:
    45 commands
    100 KB visible
```

without collapsing those differences into a single score.

---

# 14. Preserve Recovery Metrics

For AgentCap-enabled modes, preserve the B5 metrics:

```text
initial_visible_bytes
show_bytes
raw_retrieval_bytes
total_visible_bytes
show_count
raw_retrieval_count
```

A real coding agent may choose to retrieve additional information.

Those retrievals are part of the workflow cost.

Do not hide them inside a generic total.

B6 should make it possible to answer:

```text
Did the agent need show?
Did it fall back to raw?
How often?
How expensive was recovery?
```

---

# 15. Observe Real Recovery Behavior

Unlike B5, B6 should not prescribe recovery actions merely to exercise accounting.

In `integrated` mode, the coding agent should decide whether it needs:

```text
acap show
acap raw
```

through the normal AgentCap integration available to that agent.

Measure what actually happens.

This is the key transition:

```text
B5:
deterministic benchmark says:
"perform show now"

B6:
coding agent decides:
"I need more information"
```

Do not add benchmark logic that forces recovery merely to improve coverage.

Recovery accounting itself has already been validated in B5.

---

# 16. Define Objective Task Success

B6 introduces task success as a first-class benchmark dimension.

Task success must be determined by benchmark-controlled verification, not by the coding agent claiming that it succeeded.

For example:

```text
agent says:
"Fixed."
```

is not sufficient.

Instead use objective verification such as:

```text
tests pass
expected files changed
forbidden files unchanged
build succeeds
expected output produced
```

Prefer executable verification.

Conceptually:

```text
agent run
   |
   v
task verifier
   |
   +--> success
   |
   +--> failure
```

---

# 17. Keep Verification Outside the Agent Measurement Boundary

Verification commands executed by the benchmark runner after the agent has finished must be distinguishable from commands the agent itself executed.

Do not accidentally count benchmark-only verification output as agent-visible workflow cost.

Conceptually:

```text
agent command
    -> agent-visible
    -> count in workflow cost

benchmark verifier command
    -> benchmark infrastructure
    -> do NOT count as agent-visible
```

Unless the verifier output is explicitly fed back to the agent as part of the task, it is not agent-visible context.

This distinction must be explicit in the implementation.

---

# 18. Record Success Without Creating a Composite Score

The benchmark result should expose task success directly.

Conceptually:

```json
{
  "task_success": true
}
```

or an equivalent structured result.

Do not create a synthetic score such as:

```text
score =
    compression * success / latency
```

B6 should expose independent dimensions:

```text
task success
agent-visible cost
command count
recovery behavior
processing overhead
wall time
```

Later analysis can compare those dimensions.

Do not declare one mode the "winner."

---

# 19. Measure Wall Time

Measure complete agent-run wall time.

Conceptually:

```text
agent start
    |
    | wall time
    |
agent termination
```

Keep this separate from:

```text
AgentCap processing time
```

Both are useful.

AgentCap overhead may be tiny while total workflow duration changes because the agent executes additional commands.

Do not conflate them.

---

# 20. Define Agent Run Termination

The benchmark runner must robustly detect agent completion.

Handle at least:

```text
normal completion
non-zero process exit
timeout
context cancellation
adapter failure
```

Use explicit timeouts for agent workloads.

Do not allow a benchmark to hang indefinitely because an agent waits for input.

The adapter should run agents in non-interactive or automation-compatible mode where supported.

---

# 21. Prevent Interactive Prompts

Agent benchmarks must be reproducible and unattended.

Configure adapters so they do not stop for interactive questions such as:

```text
Approve this command?
Continue?
Select an option:
```

Use supported non-interactive configuration appropriate for each agent.

Do not implement fragile stdin automation to answer arbitrary prompts unless unavoidable.

If an agent cannot reliably run unattended, report that limitation rather than hiding it in benchmark-runner hacks.

---

# 22. Preserve AgentCap's Single-Execution Invariant

AgentCap instrumentation must continue to preserve:

```text
Every intercepted command executes exactly once.
```

B6 must not execute an agent command once for the agent and again merely to collect benchmark measurements.

The intended path is:

```text
agent requests command
       |
       v
single execution
       |
       v
capture
       |
       +--> raw measurement
       |
       +--> AgentCap processing
       |
       +--> agent-visible presentation
       |
       +--> benchmark accounting
```

Add integration tests where practical.

---

# 23. Do Not Compare Modes by Reusing a Mutated Workspace

Each mode must start from the same fixture state.

Never:

```text
disabled run
    |
mutated repository
    |
stateful run
```

Instead:

```text
fixture
  |
  +--> fresh workspace --> disabled
  |
  +--> fresh workspace --> stateless
  |
  +--> fresh workspace --> stateful
  |
  +--> fresh workspace --> integrated
```

This is required for meaningful comparison.

---

# 24. Result Metadata

Extend the benchmark result with enough metadata to identify the run.

At minimum, agent benchmark results should make it possible to determine:

```text
workload
agent
mode
task success
command count
agent-visible bytes
recovery metrics
wall time
AgentCap processing time
```

Include other reproducibility metadata only where it is reliable and useful.

Potential examples include:

```text
agent version
AgentCap version/commit
fixture identifier
workload version
```

Do not add large environment dumps to every result.

---

# 25. JSON Output

Support machine-readable agent benchmark output using the B4/B5 result-format infrastructure.

Conceptually:

```bash
acap bench agent \
    --workload bugfix-001 \
    --agent codex \
    --mode integrated \
    --json
```

should produce structured data similar to:

```json
{
  "schema_version": 3,
  "workload": "bugfix-001",
  "agent": "codex",
  "mode": "integrated",

  "task_success": true,

  "command_count": 24,

  "raw_bytes": 412840,
  "initial_visible_bytes": 84210,
  "show_bytes": 6810,
  "raw_retrieval_bytes": 0,
  "total_visible_bytes": 91020,

  "show_count": 2,
  "raw_retrieval_count": 0,

  "processing_ns": 12840000,
  "wall_time_ns": 82140000000
}
```

Treat this only as a conceptual example.

Use the actual B4/B5 schema and follow its versioning policy.

Do not automatically choose schema version 3 merely because this is Phase B6.

Increment the schema version only according to the established schema compatibility rules.

---

# 26. Human-Readable Output

Provide a compact summary.

For example:

```text
Benchmark: bugfix-001
Agent:     codex
Mode:      integrated

Task
  success              yes

Workflow
  commands               24
  wall time            82.1s

Agent-visible output
  initial             84,210
  show                 6,810
  raw retrieval            0
                     -------
  total               91,020

Raw baseline          412,840
Effective reduction     78.0%

Recovery
  show calls                2
  raw retrievals            0

AgentCap processing     12.8ms
```

Follow existing formatter conventions rather than reproducing this spacing exactly.

Do not automatically print claims such as:

```text
integrated wins
AgentCap is better
benchmark passed
```

Report the measured dimensions.

---

# 27. Disabled-Mode Measurement

Disabled mode is important because it establishes the real coding-agent baseline.

Ensure command output presented to the agent is measured even though AgentCap reduction is disabled.

Conceptually:

```text
disabled total visible bytes
=
normal command output visible to agent
```

Do not populate AgentCap recovery metrics with fabricated values.

They should be zero or not applicable according to the established schema conventions.

The task verifier must remain identical across modes.

---

# 28. Avoid Prompt Differences Between Modes

As much as possible, the coding task presented to the agent should remain identical across:

```text
disabled
stateless
stateful
integrated
```

Any mode-specific instructions required to use AgentCap must be minimal and explicitly attributable to integration.

Do not improve the task prompt for one mode merely to help it succeed.

Otherwise the benchmark would measure prompt differences rather than AgentCap behavior.

---

# 29. Integrated-Mode Instructions Are Part of the Experiment

If integrated mode requires instructions telling the agent how to use AgentCap, treat those instructions as part of the benchmark configuration.

Keep them stable and inspectable.

Do not dynamically rewrite them based on benchmark results.

If practical, record the integration instruction version or identifier in benchmark metadata.

Do not attempt full prompt-token accounting in B6 unless that infrastructure already exists.

Later phases can refine instruction-overhead measurement if necessary.

---

# 30. Agent Output and Logs

Keep benchmark result output separate from verbose agent logs.

The CLI should make it possible to obtain a clean benchmark result without mixing it with arbitrary agent stdout/stderr.

In particular:

```bash
acap bench agent ... --json
```

must still emit valid JSON on its result stdout channel.

Agent logs may be:

```text
captured
stored
sent to stderr
```

according to existing CLI conventions.

Do not contaminate machine-readable output.

---

# 31. Failure Handling

Distinguish infrastructure failure from task failure.

These are not equivalent:

```text
agent ran successfully but failed the task
```

and:

```text
agent process could not start
```

Conceptually represent:

```text
benchmark execution status
task success
```

separately where needed.

Examples:

```text
agent exits normally
tests fail
    -> valid benchmark trial
    -> task_success = false
```

versus:

```text
agent executable missing
    -> benchmark infrastructure error
```

Do not turn every task failure into a CLI infrastructure error.

This distinction becomes important for B7 repeated trials.

---

# 32. Timeout Is a Valid Trial Outcome When Appropriate

If the agent starts correctly but fails to finish the task before the workload timeout, preserve enough information to distinguish that outcome.

Do not fabricate:

```text
task_success = false
wall_time = 0
```

Record the actual observed workflow measurements where reliable.

Design the result so B7 can later aggregate:

```text
success
failure
timeout
```

without parsing error strings.

Do not implement B7 aggregation yet.

---

# 33. Initial Benchmark Fixtures

Create a small representative fixture set.

Prefer tasks that are:

```text
small
objective
fast to verify
easy to reset
unlikely to depend on network access
```

Good initial candidates include:

```text
compile-error repair
single failing unit-test repair
small deterministic bug fix
simple repository search task
```

Avoid starting with large real-world repositories.

The goal of B6 is validating the benchmark architecture, not establishing a comprehensive benchmark suite.

---

# 34. Tests

Add tests at multiple levels.

## Adapter tests

Verify:

```text
request construction
working directory
environment propagation
mode configuration
timeout/cancellation
exit handling
stdout/stderr capture
```

Use fake/stub agent processes where appropriate.

Unit tests must not require paid external LLM calls.

---

## Runner tests

Using a fake agent adapter, verify:

```text
fresh workspace creation
agent invocation
verification
result construction
measurement aggregation
cleanup
```

---

## Task-success tests

Test at least:

```text
agent makes correct change
    -> task_success = true

agent completes without correct change
    -> task_success = false
```

Do not use the agent's own completion message as the success signal.

---

## Measurement-boundary tests

Verify:

```text
agent-visible command output
    -> counted

benchmark-only verification output
    -> not counted
```

This distinction is critical.

---

## Mode tests

Verify that:

```text
disabled
stateless
stateful
integrated
```

select the intended execution/integration behavior.

Avoid duplicating all B0–B5 reducer tests.

Test the B6 wiring and boundaries.

---

## Recovery tests

With a fake/instrumented agent, verify that integrated-mode:

```text
show
raw
```

operations contribute to B5 recovery accounting.

---

## Single-execution tests

Verify benchmark instrumentation does not cause an intercepted command to execute more than once.

---

## JSON tests

Verify agent-specific metadata and outcome fields serialize correctly under the established schema-versioning policy.

---

# 35. Documentation

Document:

```text
acap bench agent
```

including:

```text
how to select a workload
how to select an agent
how to select a mode
how task success is verified
what measurements are included
what is not measured
```

Clearly define the four modes.

Document the distinction between:

```text
AgentCap processing time
```

and:

```text
complete agent wall time
```

Also document that B6 executes one trial per invocation unless otherwise already supported for unrelated reasons.

Repeated statistical trials belong to B7.

---

# 36. Do Not Implement B7 Yet

Do not add:

```text
--repeat 5
median calculations
mean calculations
percentiles
success-rate aggregation
multi-trial summaries
```

A B6 run represents one agent trial.

The result format should be suitable for later aggregation, but B6 must not implement that aggregation.

---

# 37. Do Not Implement B8 Yet

Do not implement:

```bash
acap bench compare ...
```

Do not automatically compare stored results.

Do not select a winning mode.

B6 should produce trustworthy individual results that B8 can compare later.

---

# 38. Do Not Implement B9 Yet

Do not add:

```text
CI regression thresholds
automatic benchmark gates
nightly benchmark orchestration
regression pass/fail policy
```

Those belong to B9.

---

# Implementation Guidance

Before modifying code:

1. inspect the complete B0–B5 benchmark architecture;
2. inspect how B3 workloads create/reset fixture workspaces;
3. inspect the B5 measurement boundary for agent-visible output;
4. inspect existing AgentCap integration mechanisms;
5. identify how coding agents can be routed through AgentCap without duplicating command execution;
6. inspect the current result-schema versioning policy;
7. inspect existing subprocess/context/timeout utilities;
8. implement the smallest viable agent adapter;
9. validate the architecture with a small deterministic fixture before adding more workloads.

Do not begin by creating a large abstraction hierarchy.

Prefer:

```text
one runner
one small adapter interface
one real adapter
one fake adapter for tests
a few deterministic fixtures
```

Then add another real adapter once the boundary has proven sufficient.

---

# Architectural Invariants

After B6, the system should conceptually look like:

```text
                    Coding Task
                        |
                        v
                Benchmark Runner
                        |
              +---------+---------+
              |                   |
              v                   v
        Agent Adapter         Task Verifier
              |
              v
          Coding Agent
              |
              v
        command/tool use
              |
              v
      AgentCap integration
              |
              v
       B0–B5 measurement
              |
              v
       agent-visible output
              |
              +--------------------+
                                   |
                                   v
                           BenchmarkResult
```

Agent-specific dependencies must point inward through the adapter:

```text
Benchmark Runner
      |
      v
AgentAdapter
      |
      +--> Codex
      +--> Claude Code
```

Never:

```text
Benchmark Runner
      |
      +--> Codex-specific parsing everywhere
      +--> Claude-specific parsing everywhere
```

And benchmark verification must remain outside the agent-visible measurement boundary:

```text
agent workflow
     |
     +--> measured agent-visible output

task verification
     |
     +--> benchmark-only output
         NOT agent-visible
```

---

# Core Comparison Model

B6 should make it possible to collect four independent results for the same task:

```text
same workload
same starting fixture
same task instruction
same verifier

        |
        +--> disabled
        |
        +--> stateless
        |
        +--> stateful
        |
        +--> integrated
```

Each result exposes:

```text
task success
agent-visible cost
command count
recovery behavior
wall time
AgentCap processing overhead
```

B6 does not combine these dimensions into a single score.

---

# Exit Criteria

Phase B6 is complete when all of the following are true:

1. A real coding agent can be executed through `acap bench agent`.
2. Agent-specific behavior is isolated behind an adapter abstraction.
3. At least one real adapter is implemented and usable.
4. Tests can use a fake adapter without external LLM calls.
5. Coding-agent workloads run in fresh isolated fixture workspaces.
6. Workloads clearly separate task instructions from benchmark verification.
7. The benchmark supports the defined `disabled`, `stateless`, `stateful`, and `integrated` modes with documented semantics.
8. Each mode starts from the same canonical fixture state.
9. The complete agent workflow is measured rather than a single command.
10. B0–B5 agent-visible and recovery accounting is reused rather than reimplemented.
11. Real agent `show` / `raw` behavior is measured in integrated mode.
12. Task success is determined by benchmark-controlled verification.
13. Benchmark-only verification output is excluded from agent-visible cost.
14. Command count is recorded.
15. Complete agent wall time is recorded separately from AgentCap processing time.
16. Infrastructure failure and task failure are distinguishable.
17. Agent benchmark results are available in both human-readable and machine-readable forms.
18. JSON mode remains clean machine-readable output.
19. AgentCap instrumentation preserves the single-execution invariant.
20. No B7 repeated-trial/statistical aggregation is introduced.
21. No B8 comparison command is introduced.
22. No B9 CI regression policy is introduced.

The completed Phase B6 should allow us to run:

```text
real coding task
        ×
real coding agent
        ×
AgentCap mode
```

and obtain one trustworthy trial answering:

```text
Did the agent complete the task?

How much command/recovery output
was actually visible to the agent?

How many commands and recovery operations
were required?

How long did the workflow take?

How much AgentCap processing overhead
was introduced?
```

without yet trying to draw statistical conclusions from a single nondeterministic LLM run.
