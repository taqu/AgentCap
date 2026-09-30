# Phase 6.5D — Cross-Agent Adapter Conformance and Benchmark Readiness

## Objective

Validate that the Claude Code and Antigravity adapters conform to the same AgentCap integration contract and are ready for fair cross-agent benchmarking.

This phase is not about adding new AgentCap capabilities.

It is about proving that:

```text
Claude Code Adapter
```

and:

```text
Antigravity Adapter
```

produce equivalent AgentCap behavior for equivalent command executions.

The goal is to prevent upcoming benchmark results from being distorted by adapter-specific differences.

---

# 1. Core Question

Phase 6.5D must answer:

```text
If Claude Code and Antigravity execute the same command
under equivalent repository and session conditions,
does AgentCap behave equivalently?
```

Equivalent behavior includes:

```text
command semantics
execution count
cwd
project root
session mapping
raw capture
result storage
classification
reduction
delta behavior
exit status
progressive disclosure
bypass behavior
fail-open behavior
metrics
```

Agent-specific transport details may differ.

AgentCap core behavior should not.

---

# 2. Architecture Under Test

The target architecture is:

```text
                 AgentCap Core
                      ^
                      |
            Common Adapter Boundary
               ^              ^
               |              |
      Claude Adapter    Antigravity Adapter
               ^              ^
               |              |
          Claude Code      Antigravity
```

The test suite must validate behavior at two levels:

```text
adapter -> common boundary
```

and:

```text
common boundary -> AgentCap core
```

Do not test only end-user output.

The normalized invocation entering AgentCap should also be testable.

---

# 3. Scope

This phase covers:

```text
Claude Code Adapter
Antigravity Adapter
Phase 6.5A common adapter contract
cross-adapter equivalence
benchmark readiness
```

It does not introduce new coding-agent adapters.

Do not add:

```text
Gemini CLI
OpenCode
other agents
```

in this phase.

---

# 4. Primary Principle

The two adapters do not need to be internally identical.

They do need to be behaviorally equivalent where AgentCap semantics are concerned.

Acceptable differences include:

```text
hook payload format
hook configuration
external session identifier source
installation files
permission integration mechanism
agent-specific correlation metadata
async/bypass details required by the host agent
```

Unacceptable differences include:

```text
different reducer selection
different AgentCap result semantics
different raw storage rules
different delta rules
different project-root behavior
different execute-once guarantees
different drill-down behavior
```

unless explicitly required by host-agent constraints and documented.

---

# 5. Build a Shared Cross-Adapter Test Harness

Create a reusable test harness that can run the same logical test case through both adapters.

Conceptually:

```text
TestCase
    |
    +-- command
    +-- cwd
    +-- environment
    +-- session identity
    +-- project
    +-- expected exit status
    +-- expected filesystem effects
    +-- expected AgentCap semantics
```

Then:

```text
run via Claude adapter
run via Antigravity adapter
compare normalized/core behavior
```

Do not duplicate every test case manually in two separate suites.

Reuse common fixtures where possible.

---

# 6. Normalized Invocation Equivalence

For equivalent tool calls, compare the normalized invocation delivered to the Phase 6.5A boundary.

For example:

```bash
go test ./...
```

executed from:

```text
/repo/src
```

within one agent session should normalize to equivalent logical information.

Compare at least:

```text
agent-independent command representation
execution cwd
project root
session identity scope
environment behavior
execution mode
bypass state
```

Do not require agent identity itself to match.

Expected:

```text
agent = claude-code
```

versus:

```text
agent = antigravity
```

is correct.

But command semantics should remain equivalent.

---

# 7. Do Not Compare Raw External Session IDs Directly

Claude Code and Antigravity use different external session identifiers.

Do not assert:

```text
Claude session ID == Antigravity conversation ID
```

Instead verify the invariant:

```text
same logical session within one adapter
    -> same AgentCap session

new logical session within one adapter
    -> different AgentCap session
```

Then compare resulting session behavior.

The adapters should have equivalent session semantics, not identical external identifiers.

---

# 8. Exactly-Once Execution Across Both Adapters

Run the same side-effect fixture through both adapters.

Example:

```bash
./increment-counter.sh
```

For each adapter:

```text
before = 0
after  = 1
```

Never:

```text
after = 2
```

Repeat with a failing side-effect command:

```text
increment
exit 7
```

Expected for both:

```text
counter incremented once
exit = 7
```

Exactly-once behavior is mandatory.

---

# 9. Command Semantics Matrix

Use the same shell-semantic fixtures for both adapters.

Include:

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

```bash
for f in *.go; do echo "$f"; done
```

Compare:

```text
stdout
stderr
exit status
filesystem effects
cwd effects
```

AgentCap presentation may differ in incidental metadata.

Underlying semantics must match.

---

# 10. Success Path Equivalence

Run representative successful commands through both adapters.

Include:

```text
small stdout
large stdout
stderr with exit 0
stdout + stderr
no output
```

and realistic commands such as:

```bash
git status
```

```bash
git diff
```

```bash
rg Resolve .
```

```bash
go test ./...
```

for passing fixtures.

Verify equivalent AgentCap behavior.

---

# 11. Failure Path Equivalence

This is mandatory.

Run representative failures through both adapters:

```text
simple exit 1
simple exit 7
command not found
Go test failure
C/C++ compiler failure
Rust cargo failure
large stderr
mixed stdout/stderr failure
```

Verify that both adapters:

```text
capture full raw output
store the result
return AgentCap presentation
preserve failure status
expose result ID
support drill-down
```

Do not accept one adapter returning compressed failures while the other exposes raw failure logs.

---

# 12. Reducer Selection Equivalence

For equivalent commands, AgentCap core should select the same reducer regardless of adapter.

Examples:

```text
git diff
    -> Git reducer

go test
    -> Go test reducer

cargo check
    -> Rust build reducer

rg
    -> search reducer
```

Assert the reducer/classifier outcome where such metadata is available.

The adapter must not influence reducer selection except through genuine command semantics.

---

# 13. Raw Capture Equivalence

For equivalent executions, verify that AgentCap stores complete raw output under both adapters.

The stored bytes need not be byte-identical if the host execution environment itself introduces legitimate differences.

However, there must be no adapter-specific truncation.

Verify:

```text
stdout complete
stderr complete
exit status stored
command metadata stored
cwd stored
```

before reduction.

---

# 14. Capsule Equivalence

For deterministic fixtures, compare AgentCap capsules produced by both adapters.

Prefer semantic comparison over fragile byte-for-byte equality if capsules include fields such as:

```text
result ID
timestamp
duration
adapter metadata
```

Compare stable fields such as:

```text
result type
status
diagnostics
test failures
Git topology
counts
structured entries
reducer-selected fields
```

Do not require identical IDs or timestamps.

---

# 15. Presentation Equivalence

Compare model-visible AgentCap output.

The outputs should be semantically equivalent.

Ignore expected differences such as:

```text
result ID
duration
agent attribution
adapter attribution
```

Do not normalize away meaningful differences.

If one adapter includes important information that the other loses, treat that as a conformance failure.

---

# 16. Progressive Disclosure Equivalence

For results generated through both adapters, verify the same retrieval operations work.

Examples:

```bash
acap show <id>
```

```bash
acap raw <id>
```

```bash
acap show <id> --errors
```

```bash
acap show <id> --warnings
```

```bash
acap show <id> --test ...
```

```bash
acap show <id> --file ...
```

Equivalent result types should expose equivalent drill-down capability.

The adapters must not create different progressive-disclosure models.

---

# 17. Stateful Compression Equivalence

Create equivalent repeated-command workflows.

Example:

```text
run failing tests
edit source
run tests again
edit source
run tests again
```

For each adapter, verify expected transitions such as:

```text
full compact result
    ->
delta
    ->
unchanged or smaller delta
```

Compare:

```text
baseline selection behavior
unchanged detection
delta classification
resolved/new/remaining failures
```

The adapter must not alter AgentCap's stateful compression logic.

---

# 18. Session Isolation Equivalence

For each adapter independently verify:

```text
session A
session B
```

do not share state incorrectly.

Then verify equivalent behavior across:

```text
same project
different session
```

and:

```text
different project
same-looking external session identifier
```

Cross-project state mixing is never allowed.

---

# 19. Project Root Equivalence

Run equivalent commands from:

```text
/repo/
/repo/src/
/repo/src/parser/
/repo/tests/
```

for both adapters.

Verify that all resolve to the same AgentCap project root when appropriate.

Expected:

```text
/repo/.acap/
```

while execution cwd remains the actual nested directory.

This should be identical across adapters.

---

# 20. Environment Equivalence

Provide controlled environment fixtures.

Example:

```text
ACAP_TEST_VALUE=hello
```

Then execute:

```bash
printf '%s\n' "$ACAP_TEST_VALUE"
```

through both adapters.

Verify equivalent environment propagation.

Also verify AgentCap-specific adapter environment variables do not leak into model-visible output unless intentionally shown.

---

# 21. Exit Status Equivalence

Test:

```text
0
1
2
7
127
```

where practical.

Both adapters must preserve the same effective command status.

AgentCap must not report wrapper status in place of underlying command status.

---

# 22. stdout / stderr Equivalence

Test:

```text
stdout only
stderr only
stdout + stderr
large stdout
large stderr
interleaved output where practical
```

Verify both adapters preserve the same AgentCap storage semantics.

Do not require identical UI ordering if the host agent itself merges streams differently and AgentCap core cannot control that.

Document any unavoidable host-level difference.

---

# 23. Recursion Equivalence

From both agents run:

```bash
acap show <id>
```

```bash
acap raw <id>
```

```bash
acap stats
```

```bash
acap clean
```

Verify these follow the shared recursion/bypass mechanism.

Neither adapter should recursively wrap AgentCap retrieval commands.

---

# 24. Explicit Bypass Equivalence

Enable the common AgentCap bypass mode.

Run identical commands under both adapters.

Verify:

```text
AgentCap interception disabled
normal host-agent command behavior preserved
no unintended stateful reduction
```

The bypass mechanism should be the same conceptually for both adapters.

This mode will become the benchmark control path.

---

# 25. Fail-Open Equivalence

Inject equivalent AgentCap processing failures.

Where practical test failures in:

```text
reducer
presentation
adapter metadata handling
storage after safe capture
```

Verify both adapters follow the same fail-open principle:

```text
command already executed successfully
+
optional AgentCap processing failed
    ->
usable command output remains available
```

Do not rerun the command.

---

# 26. Adapter Failure Classification

Both adapters should classify integration failures consistently.

Examples:

```text
adapter unavailable
invalid hook payload
AgentCap binary unavailable
internal integration failure
unsupported invocation mode
```

These should remain distinguishable from:

```text
underlying command failed
```

Do not let one adapter map integration failure to command failure while the other fails open.

---

# 27. Unsupported Modes

Claude Code and Antigravity may differ in:

```text
background execution
async execution
timeout handling
interactive execution
permission model
```

Do not force identical internal handling where the hosts are fundamentally different.

Instead require equivalent policy:

```text
supported safely
```

or:

```text
conservatively bypassed
```

Each unsupported mode must have:

```text
explicit detection
documented behavior
bypass reason
tests
```

Do not silently degrade.

---

# 28. Permission Safety Verification

Permission models differ between Claude Code and Antigravity.

Do not attempt to make them identical.

Instead verify the common invariant:

```text
AgentCap integration must not weaken the host agent's permission model.
```

For each adapter, test representative:

```text
read-only command
write command
destructive command
compound command
pipeline
subshell
```

The cross-adapter test should verify that neither adapter relies on a universal allow decision.

---

# 29. Installer Conformance

Verify both integration installers satisfy equivalent ownership principles.

Expected:

```text
install
install
```

is idempotent.

Expected:

```text
existing config
+
AgentCap install
+
AgentCap remove
=
existing config
```

Agent-specific files differ:

```text
Claude settings/hooks
```

versus:

```text
.agents/hooks.json
```

but configuration ownership behavior should be equivalent.

---

# 30. Uninstaller Conformance

For both adapters verify:

```text
only AgentCap-owned entries removed
unrelated hooks preserved
unrelated permissions preserved
configuration remains valid
second removal is safe
```

Do not compare file structures directly.

Compare behavior and ownership safety.

---

# 31. Adapter Metrics Conformance

Both adapters must emit the same logical metrics.

At minimum:

```text
agent
adapter
adapter version
commands intercepted
commands bypassed
adapter failures
adapter latency
```

Where bypass reasons are recorded, use common reason categories where possible.

For example:

```text
explicit-bypass
recursion
unsupported-async
unsupported-interactive
integration-failure
```

Avoid agent-specific reason names when a shared semantic category exists.

---

# 32. Core vs Adapter Latency

Ensure the metrics distinguish:

```text
adapter latency
```

from:

```text
AgentCap core processing latency
```

The upcoming benchmark must be able to answer:

```text
Was overhead caused by AgentCap itself
or by one specific adapter?
```

Do not collapse these into one timing field.

---

# 33. Benchmark Control Condition

Verify that each agent can run a clean AgentCap OFF condition.

Preferred:

```text
same installed adapter
+
explicit bypass
```

rather than:

```text
uninstall adapter
run task
reinstall adapter
```

This reduces configuration drift.

Confirm that OFF mode does not accidentally leave:

```text
session deltas
presentation rewriting
result interception
```

enabled.

---

# 34. Benchmark ON Condition

Verify that AgentCap ON uses the same production adapter path that normal users would use.

Do not create a special benchmark-only execution path.

The upcoming benchmark must measure:

```text
real Claude adapter
real Antigravity adapter
```

not synthetic wrappers.

---

# 35. Clean-State Benchmark Setup

Provide a reproducible helper or documented procedure for resetting local AgentCap state between benchmark trials.

At minimum consider:

```text
.acap store state
session state
previous result baselines
fixture repository modifications
generated files
adapter metrics
```

Do not accidentally carry state from one trial into the next unless the benchmark explicitly tests workflow continuation.

---

# 36. Deterministic Fixture Repository

Create or reuse a deterministic fixture repository for cross-adapter validation.

It should contain examples capable of exercising:

```text
Git diff
repository search
passing tests
failing tests
compiler errors
stateful test repair
large output
nested cwd
side effects
```

Keep the fixture small enough to run repeatedly.

Do not use a rapidly changing external repository.

---

# 37. Suggested Fixture Categories

A useful fixture set includes:

```text
fixture/simple-output
fixture/shell-semantics
fixture/git
fixture/go-tests
fixture/cpp-build
fixture/rust-build
fixture/stateful-delta
fixture/large-output
fixture/session-isolation
fixture/side-effect
```

This layout is optional.

Follow existing test organization where practical.

---

# 38. Golden Comparison Strategy

Avoid brittle full-text golden files containing:

```text
result IDs
timestamps
durations
temporary paths
agent session IDs
```

Prefer structured or normalized comparisons.

For example:

```text
normalize result ID
normalize duration
normalize temp directory
compare semantic presentation
```

Do not normalize away actual diagnostic differences.

---

# 39. Differential Testing

Where useful, implement a differential test model:

```text
Claude normalized/core result
        vs
Antigravity normalized/core result
```

Report semantic differences clearly.

Example:

```text
Mismatch:
  Claude reducer: go-test
  Antigravity reducer: generic
```

or:

```text
Mismatch:
  Claude exit: 7
  Antigravity exit: 0
```

Make conformance failures easy to diagnose.

---

# 40. Baseline Without AgentCap

For selected fixtures, run the underlying commands directly outside both adapters.

Use this only to establish command semantics.

Conceptually:

```text
direct command execution
```

serves as a semantic reference for:

```text
filesystem changes
exit status
stdout/stderr
```

Then verify both adapters preserve those semantics.

Do not use direct execution as a substitute for real adapter tests.

---

# 41. Three-Way Semantic Comparison

For critical command fixtures, compare:

```text
A. direct command

B. Claude + AgentCap

C. Antigravity + AgentCap
```

The command-level semantics should satisfy:

```text
A ~= B ~= C
```

except for intentional AgentCap presentation reduction.

This is particularly important for:

```text
exit status
filesystem effects
Git effects
cwd
environment
```

---

# 42. Stateful Workflow Comparison

Also perform workflow-level comparison.

Example:

```text
Initial failing tests
    |
fix one issue
    |
rerun
    |
fix final issue
    |
rerun
```

Collect for both adapters:

```text
number of AgentCap results
full/delta/unchanged classification
raw fallback calls
show calls
total AgentCap-visible output
```

Do not optimize these numbers during this phase.

The purpose is to detect integration inconsistency.

---

# 43. Do Not Require Identical Agent Behavior

Claude Code and Antigravity may decide to issue different commands during a real coding task.

That is expected.

Phase 6.5D conformance tests should therefore primarily use controlled commands and controlled workflows.

Do not treat:

```text
Claude chose rg
Antigravity chose grep
```

as an adapter conformance failure.

That belongs to later agent benchmark analysis.

---

# 44. Separate Adapter Conformance From Agent Quality

This phase must not score:

```text
which agent is better
which agent uses fewer tokens
which agent solves tasks faster
```

Those are benchmark questions.

Phase 6.5D only establishes that both adapters provide a fair AgentCap integration surface.

---

# 45. No Reducer Tuning

If a test reveals that both adapters produce the same undesirable AgentCap output:

```text
that is not an adapter conformance bug
```

unless it reflects an existing correctness defect.

Record reducer-quality issues for the benchmark/optimization phase.

Do not tune reducers during Phase 6.5D merely to make tests look better.

---

# 46. Detect Adapter-Specific Divergence

If only one adapter produces problematic behavior, classify it as an adapter issue.

Examples:

```text
Claude preserves exit=7
Antigravity returns exit=0
```

```text
Claude stores complete stderr
Antigravity truncates stderr
```

```text
Claude maps session correctly
Antigravity creates a new session every command
```

These must be fixed before benchmark work begins.

---

# 47. Detect Core Divergence Caused by Adapter Input

A mismatch may originate from different normalized input.

For every differential failure determine whether the difference starts at:

```text
host agent
adapter parsing
normalization
common integration boundary
AgentCap core
presentation
```

Add diagnostics sufficient to identify this boundary.

Do not patch symptoms at the presentation layer.

---

# 48. Cross-Agent Conformance Report

Add a machine-readable or test-generated conformance summary.

It should include categories such as:

```text
execution
shell semantics
exit status
cwd
environment
raw capture
reducer selection
storage
session mapping
delta
progressive disclosure
recursion
bypass
fail-open
installation
metrics
```

Each category should be reportable as:

```text
PASS
FAIL
SKIP with explicit reason
```

Avoid silent skips.

---

# 49. Skip Policy

A test may be skipped only when the underlying host agent genuinely cannot support that mode.

Examples might include unsupported platform behavior.

Every skip must state:

```text
adapter
test
reason
whether benchmark scope is affected
```

Do not use skips to hide flaky adapter behavior.

---

# 50. Flakiness Policy

Cross-agent conformance must be stable before benchmarking.

Run critical tests multiple times where timing is involved.

Especially watch:

```text
async execution
timeouts
cancellation
concurrent commands
session reuse
```

Fix deterministic adapter races.

Do not accept recurring flakes as normal.

---

# 51. Concurrency Test

Where both agents support concurrent command execution, run two controlled commands simultaneously.

Verify no cross-talk in:

```text
stdout
stderr
cwd
session
result ID
storage record
adapter metadata
```

If one host does not expose concurrent shell calls, document that host limitation rather than fabricating concurrency.

---

# 52. Large Output Test

Generate deterministic large output.

Test through both adapters.

Verify:

```text
command executes once
raw output captured completely
no pipe deadlock
AgentCap reduces output
result is recoverable
adapter memory behavior remains reasonable
```

Compare reduction semantics.

---

# 53. Large Failure Test

Generate deterministic large stderr or compiler/test failure output.

This is especially important because AgentCap exists to reduce such context.

Both adapters must produce equivalent compact failure behavior.

A raw-output leak in one adapter is a blocker for benchmark readiness.

---

# 54. Command Families to Cover

At minimum include representative commands from:

```text
generic shell
search
Git
build
test
failure
large output
stateful repeated execution
```

Examples:

```bash
ls
rg
git status
git diff
go test
cargo check
clang++
```

Use the toolchains already available in CI/test environments.

Do not introduce large dependency requirements solely for coverage.

---

# 55. Bypass Reason Consistency

Use consistent semantic bypass categories.

For example:

```text
explicit
recursion
unsupported-background
unsupported-interactive
unsupported-host-mode
integration-failure
```

The exact enum names may differ from these examples.

Do not create:

```text
claude-background-thing
```

and:

```text
antigravity-async-case
```

if both represent the same logical limitation.

This matters for benchmark aggregation.

---

# 56. Metrics Validation

Verify metric counters match actual test behavior.

For a controlled run with:

```text
10 intercepted
2 bypassed
1 simulated adapter failure
```

metrics should report the same values.

Test both adapters.

Do not assume instrumentation is correct merely because commands work.

---

# 57. Adapter Latency Measurement Validation

Verify timing measurement boundaries.

For both adapters, clearly define:

```text
adapter start
adapter end
core processing start
core processing end
```

Avoid double-counting core time inside adapter-only latency if the benchmark expects these values separately.

Document the timing definitions.

---

# 58. Result Attribution

Every stored result generated through an adapter should be attributable to the correct agent integration where Phase 6.5A metadata supports it.

Example:

```text
agent = claude-code
```

or:

```text
agent = antigravity
```

Do not let agent attribution affect reducer/core behavior.

It is metadata only.

---

# 59. Adapter Version Attribution

Verify benchmark/debug metadata includes enough information to distinguish adapter revisions.

This may reuse AgentCap build version.

The benchmark should not accidentally combine results from:

```text
adapter version X
adapter version Y
```

without being able to identify them.

---

# 60. Installation Verification in Clean Repositories

Test installation into a clean project for both adapters.

Then verify:

```text
agent starts
normal command works
AgentCap intercepts
result stored
uninstall works
```

Also test installation into repositories that already contain unrelated hook configuration.

---

# 61. Do Not Modify Agent Configuration for Benchmark Advantage

Do not tune:

```text
agent prompts
permission defaults
model settings
system instructions
```

differently merely to make AgentCap integration perform better.

Cross-agent benchmarks may require different native configuration formats, but adapter conformance must not depend on special model prompting.

---

# 62. Minimal Adapter Instructions

If either adapter requires model-facing instructions for progressive disclosure, document them explicitly.

Keep them minimal.

Measure their size.

The benchmark must later include such instruction overhead in AgentCap context cost where appropriate.

Do not hide instruction overhead outside benchmark accounting.

---

# 63. Progressive Disclosure Instruction Parity

If one adapter requires additional instruction text and the other does not, record that as an adapter difference.

Do not artificially add unnecessary instructions to the other adapter just to make them textually equal.

Fairness means measuring real integration cost.

---

# 64. Benchmark Metadata Contract

Before leaving Phase 6.5D, define the stable adapter metadata fields that the benchmark CLI will consume.

At minimum:

```text
agent
adapter
adapter_version

agentcap_enabled

commands_intercepted
commands_bypassed
adapter_failures

adapter_latency_ms
agentcap_processing_latency_ms
```

Reuse existing field names if Phase 6.5A already established them.

Do not create duplicate metrics schemas.

---

# 65. Do Not Build the Full Benchmark Yet

Phase 6.5D may expose and validate benchmark-facing metadata.

It should not yet implement:

```text
task scoring
agent quality scoring
benchmark datasets
statistical aggregation
leaderboards
cross-model evaluation
```

Those belong to the benchmark phase.

---

# 66. Benchmark Readiness Gate

Create a single readiness check or documented checklist.

Benchmarking must not begin unless:

```text
Claude adapter conformance passes
Antigravity adapter conformance passes
cross-adapter differential tests pass
critical exactly-once tests pass
failure compression tests pass
bypass works
metrics are validated
```

Make benchmark readiness explicit.

---

# 67. Blockers

Treat the following as benchmark blockers:

```text
duplicate command execution
different command semantics
exit-code corruption
raw output loss
failed-command compression bypass
incorrect project/session isolation
permission weakening
broken progressive disclosure
adapter-specific reducer behavior
unstable/flaky interception
metrics that cannot distinguish adapter overhead
```

Do not proceed to benchmark collection with these unresolved.

---

# 68. Non-Blocking Differences

Differences that may be documented rather than blocked include:

```text
different hook configuration format
different session ID syntax
different host permission UI
different host async mechanism
different installation file
different unavoidable host metadata
```

provided AgentCap semantics remain correct.

---

# 69. Suggested Test Matrix

Use a matrix similar to:

```text
                               Claude   Antigravity

basic success                    PASS       PASS
basic failure                    PASS       PASS
execute exactly once             PASS       PASS
shell semantics                  PASS       PASS
cwd                              PASS       PASS
environment                      PASS       PASS
exit status                      PASS       PASS
Git reducer                      PASS       PASS
test reducer                     PASS       PASS
compiler failure                 PASS       PASS
raw recovery                     PASS       PASS
stateful delta                   PASS       PASS
session isolation                PASS       PASS
project isolation                PASS       PASS
recursion                        PASS       PASS
bypass                           PASS       PASS
fail-open                        PASS       PASS
large output                     PASS       PASS
permission safety                PASS       PASS
installer idempotency            PASS       PASS
uninstall safety                 PASS       PASS
metrics                          PASS       PASS
```

Generate this from tests where practical rather than maintaining it manually.

---

# 70. Cross-Adapter Semantic Diff

For deterministic fixtures, produce a concise semantic diff when adapters disagree.

Example:

```text
fixture: go-test-failure

Claude:
  exit: 1
  reducer: go-test
  failures: 3
  raw_stored: true

Antigravity:
  exit: 1
  reducer: generic
  failures: unknown
  raw_stored: true
```

This should immediately identify the mismatch.

Avoid dumping entire raw command logs unless requested.

---

# 71. Test Isolation

Each conformance test must start from known state.

Reset as needed:

```text
repository working tree
.acap state
environment
temporary files
session mappings
side-effect counters
```

Do not allow test order to determine outcomes.

---

# 72. Cross-Platform Scope

Run on the operating systems AgentCap currently supports.

Do not expand platform support during this phase.

If Claude and Antigravity support different subsets of AgentCap-supported platforms, document that clearly.

Do not call a platform supported unless the adapter test passes there.

---

# 73. CI Integration

Add automated conformance tests to CI where practical.

Separate:

```text
pure unit/conformance tests
```

from:

```text
real installed Claude/Antigravity end-to-end tests
```

if the latter cannot run in normal CI.

The core conformance suite should remain runnable without requiring live agent sessions whenever possible.

---

# 74. Real-Agent Smoke Tests

In addition to mocks, perform at least one real end-to-end smoke test for each adapter before declaring readiness.

Mocks alone are insufficient because:

```text
hook behavior
permission flow
async behavior
host execution details
```

may differ from assumptions.

---

# 75. Mock Contract Validation

Where adapters use mocked hook payloads in tests, periodically validate those fixtures against actual current host payloads.

Do not let mock JSON drift indefinitely from Claude Code or Antigravity.

Keep fixtures minimal and version-aware where necessary.

---

# 76. Documentation

Add concise developer documentation covering:

```text
what cross-adapter conformance means
what is expected to be equal
what is allowed to differ
how to run the conformance suite
how to run adapter-specific tests
how to interpret failures
how to reset test state
how benchmark readiness is determined
```

Do not duplicate the full Phase 6.5A/B/C implementation documents.

---

# 77. Final Cross-Agent Workflow Test

Run one controlled coding workflow through both adapters.

Use the same repository state and same scripted command sequence.

For example:

```text
1. repository search
2. git status
3. failing tests
4. inspect targeted failure
5. modify fixture
6. rerun tests
7. git diff
8. rerun tests unchanged
```

Compare AgentCap behavior at every command boundary.

Do not compare model reasoning quality in this test.

---

# 78. Expected Workflow-Level Equivalence

For an identical scripted sequence, the two adapters should produce equivalent AgentCap state transitions such as:

```text
search        -> compact result
git status    -> compact result
test failure  -> structured failure
rerun         -> delta
git diff      -> structured Git result
rerun         -> unchanged
```

Differences in result IDs and timings are expected.

Differences in AgentCap semantic state are not.

---

# 79. Final Acceptance Criteria

Phase 6.5D is complete only when all of the following are true.

## Common Contract

- Claude Code and Antigravity both pass the Phase 6.5A conformance suite.
- Both use the same agent-independent integration boundary.
- No new agent-specific branches have been added to AgentCap core.

## Execution

- Equivalent commands preserve equivalent semantics.
- Successful commands execute exactly once.
- Failed commands execute exactly once.
- Side-effect tests pass for both adapters.

## Output

- Full stdout/stderr remain recoverable.
- Failure output is compressed for both adapters.
- Successful output is processed consistently.
- Result IDs remain usable for progressive disclosure.

## Core Behavior

- Equivalent commands select equivalent reducers.
- Equivalent structured results produce equivalent capsules.
- Stateful delta behavior is equivalent.
- Project-root behavior is equivalent.
- Session semantics are equivalent.

## Safety

- Neither adapter weakens host permission behavior.
- recursion protection works.
- bypass works.
- fail-open behavior works.
- unsupported modes are explicitly and conservatively handled.

## Configuration

- both installers are idempotent.
- both uninstallers preserve unrelated configuration.
- AgentCap-owned configuration can be identified reliably.

## Metrics

Both adapters correctly expose:

```text
agent
adapter
adapter version
commands intercepted
commands bypassed
adapter failures
adapter latency
AgentCap processing latency
```

## Benchmark Control

- AgentCap ON works through the production adapter path.
- AgentCap OFF works through explicit bypass.
- switching ON/OFF does not require destructive configuration changes.
- test state can be reset reproducibly.

## Differential Tests

- no unexplained semantic differences remain between adapters.
- all benchmark-blocking differences are resolved.
- skips are documented explicitly.
- critical tests are stable and non-flaky.

## End-to-End

- a real Claude Code smoke test succeeds.
- a real Antigravity smoke test succeeds.
- a controlled equivalent workflow produces equivalent AgentCap behavior.

---

# 80. Implementation Order

Use this order:

```text
1. Inspect Phase 6.5A/B/C implementations.

2. Define the semantic equivalence model.

3. Build the shared cross-adapter harness.

4. Add normalized invocation comparison.

5. Add exactly-once differential tests.

6. Add shell-semantics tests.

7. Add success/failure comparison.

8. Add reducer/capsule comparison.

9. Add project/session isolation tests.

10. Add stateful delta comparison.

11. Add progressive-disclosure comparison.

12. Add bypass/recursion/fail-open comparison.

13. Add adapter metrics validation.

14. Add installer/uninstaller conformance tests.

15. Add deterministic fixture repository/workflows.

16. Run real Claude Code smoke test.

17. Run real Antigravity smoke test.

18. Produce conformance summary.

19. Fix all benchmark-blocking differences.

20. Mark the adapters benchmark-ready.
```

Do not begin benchmark result collection until the readiness gate passes.

---

# 81. Design Priority

When deciding whether a difference is acceptable, use this priority order:

```text
1. Preserve command semantics.
2. Preserve permission/security semantics.
3. Execute exactly once.
4. Preserve exit and failure information.
5. Preserve complete raw recoverability.
6. Preserve project/session correctness.
7. Preserve AgentCap core equivalence.
8. Preserve progressive disclosure.
9. Preserve benchmark measurability.
10. Minimize adapter overhead.
```

Do not sacrifice correctness merely to make adapter outputs look identical.

---

# 82. Final Deliverable

At the end of Phase 6.5D, AgentCap should have two independently implemented but behaviorally conformant production adapters:

```text
Claude Code Adapter
Antigravity Adapter
```

Both should feed the same AgentCap core semantics.

The project should be able to state confidently:

```text
Differences observed in the upcoming benchmark are not caused
by basic adapter correctness or inconsistent AgentCap integration.
```

The next phase can then begin measuring:

```text
AgentCap OFF
vs
AgentCap ON
```

across real coding-agent workflows without first questioning whether the two adapters are behaving fundamentally differently.
