# Task: Implement AgentCap Benchmark CLI Phase B3 — Reproducible Workloads

Implement Phase B3 of the AgentCap benchmark CLI.

Phase B0+B1 established reliable single-command measurement.

Phase B2 added stateful benchmark sessions capable of measuring:

```text
raw workflow cost
vs
stateless AgentCap cost
vs
stateful AgentCap cost
```

Phase B3 adds a reproducible workload layer on top of the B2 benchmark-session API.

The goal is to define realistic multi-step command workflows that can be executed repeatedly from a known initial workspace state.

Do not implement coding-agent orchestration, LLM evaluation, statistical repeated trials, historical benchmark comparison, or CI regression thresholds in this phase.

---

# 1. Inspect B0-B2 Before Designing B3

Before implementing anything, inspect the actual benchmark infrastructure currently present in the repository.

Determine:

- B0 measurement structures
- B1 single-command benchmark execution path
- B2 benchmark-session API
- how benchmark sessions map to AgentCap sessions
- how raw/stateless/stateful bytes are measured
- how full/delta/unchanged presentations are classified
- how benchmark data survives separate CLI invocations
- current CLI conventions
- existing JSON support, if any
- current result-store schema
- existing test fixture conventions
- existing temporary-workspace helpers

Build B3 on top of these mechanisms.

Do not introduce a parallel execution or measurement architecture.

---

# 2. Core Goal

A developer should be able to define a benchmark workload such as:

```text
initial fixture
    |
    v
git diff
    |
    v
deterministic workspace mutation
    |
    v
git diff
    |
    v
deterministic workspace mutation
    |
    v
go test ./...
    |
    v
benchmark summary
```

and execute it reproducibly through a command conceptually similar to:

```bash
acap bench run <workload>
```

The exact syntax should follow the existing CLI conventions.

The workload runner must reuse the B2 benchmark-session API so the final report includes:

```text
raw bytes
stateless bytes
stateful bytes

full presentations
delta presentations
unchanged presentations

execution duration
AgentCap processing duration
```

as already defined by B0-B2.

---

# 3. Why Workloads Need State Transitions

Do not model B3 as only:

```yaml
commands:
  - git diff
  - git diff
  - git diff
```

That can test unchanged-result deduplication, but it does not adequately test AgentCap's main stateful behavior.

Realistic workflows include transitions such as:

```text
compile
  -> 7 errors

fix source

compile
  -> 3 errors

fix source

compile
  -> PASS
```

or:

```text
git diff
  -> files A, B changed

edit A

git diff
  -> A changed further, B unchanged
```

The workload format therefore needs a small, deterministic way to mutate workspace state between command executions.

Keep this mechanism intentionally narrow.

---

# 4. Workload Format

Introduce a simple declarative workload format.

Prefer YAML if the repository does not already have a stronger convention.

A workload may conceptually look like:

```yaml
version: 1

name: git-repeated-diff

fixture: fixtures/git-repeated-diff

steps:
  - run:
      argv: ["git", "diff"]

  - copy:
      from: states/change-2/main.go
      to: main.go

  - run:
      argv: ["git", "diff"]

  - run:
      argv: ["git", "diff"]
```

This is only a conceptual example.

Design the actual schema carefully based on repository needs.

Do not implement a general-purpose scripting language.

---

# 5. Version the Workload Schema

Every workload definition must contain an explicit schema version.

For example:

```yaml
version: 1
```

Reject unsupported versions clearly.

Do not silently interpret unknown versions using the newest parser.

The workload format will likely evolve, so schema versioning must exist from the beginning.

---

# 6. Keep the Initial DSL Small

B3 should support only the operations required for deterministic AgentCap benchmark scenarios.

A reasonable initial set is:

```text
run
copy
write
remove
```

Potentially add:

```text
mkdir
```

only if actual benchmark fixtures require it.

Do not add:

```text
arbitrary shell mutation hooks
conditionals
loops
variables
templating
network operations
sleep-based synchronization
embedded scripts
```

unless there is a compelling existing requirement.

The benchmark format should describe controlled state transitions, not become another shell.

---

# 7. `run` Step

The `run` operation executes a command through the normal B2 benchmark-session path.

Prefer structured argv:

```yaml
- run:
    argv: ["go", "test", "./..."]
```

rather than:

```yaml
- run: "go test ./..."
```

Do not reconstruct structured argv into:

```text
sh -c
```

unless the existing AgentCap execution semantics explicitly require it.

Each `run` step must preserve the existing exactly-once invariant.

A workload containing five `run` steps must execute exactly five target commands.

---

# 8. Mutation Steps Must Not Be AgentCap Commands

Operations such as:

```text
copy
write
remove
mkdir
```

are benchmark harness operations.

They prepare workspace state.

They are NOT benchmarked command executions and must not:

- create AgentCap result IDs
- count toward benchmark command count
- contribute command output bytes
- enter AgentCap session baseline selection

Conceptually:

```text
run
  -> measured

mutation
  -> not measured

run
  -> measured
```

Keep this boundary explicit in the implementation.

---

# 9. Prefer Fixture-State Replacement Over Patch Engines

Do not build a patch/diff application engine unless the repository already has a safe implementation that can be reused.

For deterministic benchmark transitions, prefer simple operations such as:

```text
copy fixture state N
    ->
workspace path
```

For example:

```text
states/01/main.go
states/02/main.go
states/03/main.go
```

can represent:

```text
7 errors
3 errors
PASS
```

This is easier to understand and reproduce than dynamically editing source text.

---

# 10. Isolated Temporary Workspace

Never run destructive benchmark mutations directly against the user's working repository.

Each workload run should execute in an isolated temporary workspace.

Conceptually:

```text
fixture
   |
   v
temporary workspace
   |
   +--> benchmark commands
   |
   +--> deterministic mutations
   |
   v
discard
```

The original fixture must remain unchanged.

The user's current source tree must remain unchanged.

---

# 11. Fixture Copying

Define fixture-copy semantics carefully.

The runner should create a fresh workspace from the fixture for every workload invocation.

For example:

```text
benchmarks/fixtures/git-repeated-diff/
                   |
                   v
/tmp/acap-bench-.../
```

Do not reuse a dirty workspace from a previous run.

The benchmark must begin from a known initial state every time.

---

# 12. Git Fixtures

Some important AgentCap benchmarks require real Git repository state.

For Git workloads, support fixtures that can deterministically establish:

```text
repository
baseline commit
working-tree changes
```

Do not depend on the developer's current repository Git state.

A Git workload should run inside its temporary fixture repository.

If fixture storage should avoid embedding `.git/` directories, initialize the temporary repository deterministically during fixture setup.

Prefer the simplest reproducible approach.

Document the chosen strategy.

---

# 13. Deterministic Git Setup

If B3 initializes Git repositories dynamically, make them independent from user-level Git configuration where practical.

For example, tests should not fail merely because the developer has no global:

```text
user.name
user.email
```

configured.

If commits are required, configure the temporary repository locally.

Avoid dependencies on:

```text
global Git config
Git hooks
network remotes
user signing configuration
```

Benchmark fixtures should be self-contained.

---

# 14. Working Directory

Allow workloads or `run` steps to specify a working directory relative to the temporary workspace when necessary.

For example:

```yaml
- run:
    argv: ["go", "test", "./..."]
    cwd: "project"
```

All workload paths must remain relative to the isolated benchmark workspace.

Do not allow a workload to escape into arbitrary host filesystem locations.

---

# 15. Path Safety

Treat workload definitions as potentially unsafe input.

For all fixture and mutation paths:

- reject absolute destination paths
- reject `..` traversal that escapes the benchmark workspace
- resolve symlinks conservatively
- prevent mutation operations from escaping the temporary workspace

A workload must not be able to write:

```text
~/.ssh/
../real-project/
 /etc/
```

through path manipulation.

Add tests for path traversal rejection.

---

# 16. Source Paths

Source fixture paths should also be constrained to an explicitly allowed benchmark/fixture root.

Do not allow:

```yaml
copy:
  from: /some/arbitrary/host/file
```

unless there is an explicit future design for external inputs.

B3 workloads should be self-contained and portable.

---

# 17. Environment

Use the existing AgentCap environment inheritance semantics for `run` commands unless reproducibility requires explicit overrides.

If environment overrides are supported, keep them declarative and minimal.

For example:

```yaml
env:
  GOFLAGS: "-count=1"
```

or per step if necessary.

Do not capture or serialize the user's entire environment into workload definitions.

Avoid environment-dependent benchmark behavior where possible.

---

# 18. Network Independence

B3 benchmark fixtures should preferably run without network access.

Do not make initial benchmark workloads depend on:

```text
package downloads
remote Git repositories
HTTP services
external APIs
```

If a language tool normally downloads dependencies, design the fixture so tests can run from already-contained or standard-library-only code where practical.

The goal is deterministic AgentCap evaluation, not dependency-manager benchmarking.

---

# 19. Workload Runner

Introduce a reusable workload runner below the CLI layer.

Conceptually:

```text
Workload Definition
       |
       v
Workload Parser
       |
       v
Validator
       |
       v
Workspace Setup
       |
       v
B2 Benchmark Session
       |
       v
Step Runner
       |
       +--> run ------> B2 command measurement
       |
       +--> mutation -> workspace only
       |
       v
Session Aggregate
```

Do not place all workload logic directly inside the CLI command handler.

Future benchmark phases should be able to invoke the workload runner programmatically.

---

# 20. Fail Fast on Invalid Workloads

Validate the workload before executing target commands wherever possible.

Examples of validation errors:

```text
unsupported schema version
missing name
missing fixture
unknown step type
empty argv
invalid cwd
absolute destination path
path traversal
missing fixture file
duplicate/ambiguous fields
```

Do not partially execute a workload and only later discover that its definition was structurally invalid.

---

# 21. Runtime Failure Semantics

Distinguish:

```text
benchmark infrastructure failure
```

from:

```text
target command returned non-zero
```

A non-zero target command is often expected.

For example:

```text
go test ./...
```

may intentionally fail during the first step of a fail-fix-pass workload.

Do not abort the workload merely because a benchmarked command exits non-zero unless the workload explicitly requires a particular exit status.

---

# 22. Expected Exit Status

Support a simple optional assertion for `run` steps.

For example:

```yaml
- run:
    argv: ["go", "test", "./..."]
    expect:
      exit: 1
```

and later:

```yaml
- run:
    argv: ["go", "test", "./..."]
    expect:
      exit: 0
```

This helps verify that the workload itself remains valid.

If no expected exit status is specified, record the actual status without treating non-zero as a harness failure.

Keep assertions small in B3.

---

# 23. Minimal Assertions

In addition to exit status, consider only deterministic assertions that materially improve fixture correctness.

Potentially useful:

```text
stdout contains
stderr contains
```

but add them only if needed by real fixtures.

Do not build a full assertion language.

The purpose of B3 assertions is to detect broken benchmark fixtures, not to replace a testing framework.

---

# 24. Benchmark Session Boundary

Every workload invocation should create a fresh B2 benchmark session unless the CLI explicitly supports another well-defined mode.

Conceptually:

```text
acap bench run workload.yaml
        |
        v
new benchmark session
        |
        +--> step #1
        +--> step #2
        +--> step #3
        |
        v
aggregate
```

This guarantees that previous AgentCap history does not artificially improve benchmark results.

---

# 25. Result Persistence

Every `run` step remains a normal independently complete AgentCap result.

Preserve:

```text
raw stdout
raw stderr
metadata
capsule
indexes
result ID
```

as appropriate.

A user should be able to inspect a result from a completed workload using existing commands such as:

```bash
acap show <result-id>
acap raw <result-id>
```

Do not store only benchmark aggregates.

---

# 26. Workload Result

The workload result should contain at least:

```text
workload name
benchmark session ID
number of run steps
raw bytes
stateless bytes
stateful bytes
full count
delta count
unchanged count
execution duration
AgentCap processing duration
overall workload duration
```

Reuse B2 aggregate semantics.

Do not redefine byte metrics in B3.

---

# 27. Wall-Clock Duration

B3 may introduce an additional:

```text
workload wall-clock duration
```

covering:

```text
workspace setup
+
mutations
+
command execution
+
AgentCap processing
+
harness overhead
```

If added, clearly distinguish it from:

```text
command execution duration
AgentCap processing duration
```

Do not silently change B0-B2 timing definitions.

---

# 28. Human-Readable Output

A successful run should produce a concise summary.

Conceptually:

```text
Benchmark Workload: git/repeated-diff

session: 8a13c91d...

steps:
  total:       5
  commands:    3
  mutations:   2

output:
  raw:                  184220 B
  stateless:             28140 B
  stateful:              13720 B

reduction:
  stateless vs raw:       84.7%
  stateful vs raw:        92.6%
  stateful vs stateless:  51.2%

presentations:
  full:                       1
  delta:                      1
  unchanged:                  1

time:
  command execution:        ...
  AgentCap processing:      ...
  workload wall clock:      ...
```

Adapt formatting to existing CLI conventions.

Do not print all command output by default.

---

# 29. Failure Output

If the workload infrastructure fails, report:

- workload name
- failing step number
- step type
- concise reason

For example:

```text
workload failed at step 4 (copy):
source fixture file not found: states/03/main.go
```

Do not bury harness failures inside ordinary command output.

---

# 30. Optional Per-Step Reporting

If useful, provide a concise per-step mode or verbose flag.

For example:

```text
#  step       exit   raw     stateless   stateful   kind
1  git diff      0   41 KB    7 KB        7 KB      full
2  copy          -     -       -           -        -
3  git diff      0   46 KB    8 KB        3 KB      delta
4  git diff      0   46 KB    8 KB       40 B       unchanged
```

Do not make verbose per-step output the default if it obscures the aggregate result.

---

# 31. Workload Discovery

Use a predictable repository layout.

A reasonable structure is:

```text
benchmarks/
├── workloads/
│   ├── git/
│   │   ├── repeated-diff.yaml
│   │   └── status-change-status.yaml
│   ├── build/
│   │   └── compile-fix.yaml
│   └── test/
│       └── fail-fix-pass.yaml
│
└── fixtures/
    ├── git-repeated-diff/
    ├── compile-fix/
    └── fail-fix-pass/
```

Adapt this if the repository already has benchmark/test conventions.

Keep workload definitions separate from fixture content.

---

# 32. Workload Naming

Give every workload a stable logical name independent of its absolute filesystem location.

For example:

```yaml
name: git/repeated-diff
```

This will be useful in later phases for:

```text
result storage
comparison
repeated trials
CI
agent benchmarks
```

Reject obviously invalid or duplicate names where relevant.

Do not design a global benchmark registry yet.

---

# 33. Initial Benchmark Workloads

Add a small set of high-value workloads that exercise existing AgentCap strengths.

Do not attempt broad command coverage.

Prefer approximately three workloads.

## A. Git repeated diff

Exercise:

```text
git diff
mutation
git diff
git diff unchanged
```

This should demonstrate:

```text
full
delta
unchanged
```

where supported by the existing Git delta implementation.

## B. Build/compile transition

Exercise a deterministic transition such as:

```text
compile failure
mutation
fewer/different errors
mutation
success
```

Choose a toolchain already supported by AgentCap and readily available in the repository's test environment.

## C. Test transition

Exercise:

```text
failing test
mutation
passing test
repeated passing test
```

Again, use an existing supported toolchain.

If environment portability makes one of these impractical, implement fewer fixtures rather than introducing network/toolchain fragility.

---

# 34. Fixture Size

Keep B3 fixtures intentionally small.

They need to generate representative command output, not become realistic production repositories.

Optimize for:

```text
determinism
clarity
fast execution
easy debugging
```

Large-repository benchmarking can be added later after the framework is trustworthy.

---

# 35. Determinism

A workload should produce semantically equivalent command/output transitions when run repeatedly from the same code revision and toolchain environment.

Avoid embedding:

```text
current timestamps
random values
temporary absolute paths in expected output
host-specific usernames
network results
machine-specific state
```

where possible.

If a tool naturally emits nondeterministic fields such as durations, do not assert those exact fields.

---

# 36. Exactly-Once Invariant

The B0/B1/B2 exactly-once invariant remains mandatory.

If a workload contains:

```text
run A
mutation
run B
run C
```

then exactly:

```text
A once
B once
C once
```

must execute.

Do not rerun commands to obtain:

```text
raw cost
stateless cost
stateful cost
assertions
```

All measurements and assertions must derive from the single captured execution.

Add an end-to-end workload test that proves this.

---

# 37. Do Not Benchmark Mutation Output

Mutation operations may produce internal errors or diagnostic information for the harness, but this is not AgentCap-visible command output.

Do not include mutation operation bytes in:

```text
raw_bytes
stateless_bytes
stateful_bytes
```

They may contribute to overall workload wall-clock time if that metric is implemented.

---

# 38. Cleanup

Temporary benchmark workspaces should normally be removed after completion.

Provide a reasonable debugging mechanism for preserving the workspace when a workload fails, only if useful.

For example, a flag conceptually similar to:

```text
--keep-workspace
```

may be acceptable.

If implemented, clearly print the retained path.

Do not retain temporary benchmark repositories by default.

---

# 39. Interruptions and Cancellation

Preserve existing AgentCap cancellation semantics.

If the benchmark process is interrupted:

- stop the currently running command through the normal cancellation path
- do not start subsequent workload steps
- clean up the temporary workspace where practical
- do not corrupt normal AgentCap state

Do not invent a separate process-management mechanism for workloads.

---

# 40. JSON Output

If B0-B2 already support stable JSON output, extend it naturally to workload results.

If JSON was intentionally deferred to Phase B4, continue to defer it.

Do not prematurely design the full long-term benchmark result schema in B3.

The workload definition schema and benchmark result schema are separate concerns.

The workload definition MUST be versioned even if result JSON remains deferred.

---

# 41. Security Boundary

The workload format executes commands.

Therefore it is not a safe format for untrusted workloads.

Document that benchmark workloads should be treated similarly to scripts or test code from the repository.

However, still enforce filesystem path boundaries for built-in mutation operations.

Do not imply that the workload runner is a general sandbox.

---

# 42. Tests

Add focused unit and integration tests.

At minimum cover:

## Parsing

- valid workload
- unsupported schema version
- unknown step
- missing required fields
- empty command argv

## Path safety

- absolute mutation destination rejected
- `..` traversal rejected
- fixture source escape rejected
- valid nested relative paths accepted

## Workspace isolation

Verify that:

```text
original fixture
```

is unchanged after execution.

Verify that:

```text
user repository
```

is not modified.

## Fresh-run behavior

Run the same workload twice.

Each invocation must begin from the same fixture state and a fresh benchmark session.

## Stateful behavior

Use an actual workload to exercise the existing B2 session/delta pipeline.

Do not fake presentation kinds.

## Exactly once

Use an observable side effect to verify every `run` step executes once.

## Expected exit status

Test expected failure and expected success.

## Unexpected assertion result

Verify a broken fixture produces a clear workload failure.

## Cleanup

Verify temporary workspace cleanup on success.

Test failure cleanup behavior according to the chosen policy.

## Result recoverability

Verify results produced by workload `run` steps remain available through normal AgentCap result retrieval.

---

# 43. Do Not Overfit Tests to Byte Counts

Avoid tests that require exact compressed byte totals for large real tool output unless the rendering format is intentionally stable.

Prefer invariants such as:

```text
raw > 0

stateless <= raw
where expected

stateful < stateless
for a known unchanged/delta scenario

command_count == expected

full/delta/unchanged counts == expected
```

Exact byte-count tests are appropriate for small controlled renderer unit tests.

Do not make the benchmark suite brittle to harmless formatting changes.

---

# 44. No Automatic Benchmark Repetition

Do not add:

```text
--repeat
```

in B3.

A workload invocation executes one trial.

Repeated trials and statistical aggregation belong to a later phase.

This also keeps the semantics around side effects and workspace lifecycle simple.

---

# 45. No Agent Execution

Do not invoke coding agents in B3.

The workload runner itself determines the command sequence.

This phase is deterministic infrastructure for measuring AgentCap.

Later phases can allow a coding agent to operate inside a benchmark task environment.

---

# 46. No Task-Success Evaluation Yet

B3 assertions validate that the benchmark workload behaved as designed.

They do not measure coding-agent task success.

Do not introduce:

```text
LLM judges
solution scoring
patch correctness evaluation
agent success rates
```

Those belong to the agent-workflow benchmark phase.

---

# 47. No Historical Comparison Yet

Do not implement:

```bash
acap bench compare ...
```

in B3.

Do not introduce baseline benchmark databases or regression thresholds.

B3 creates reproducible runs.

B4 and later phases can define stable result serialization and comparison.

---

# 48. Architecture Boundary

Keep the architecture approximately:

```text
                    Workload Definition
                           |
                           v
                       Parser
                           |
                           v
                      Validator
                           |
                           v
                    Workspace Setup
                           |
                           v
                   Workload Runner
                           |
              +------------+------------+
              |                         |
              v                         v
         run step                   mutation step
              |                         |
              v                         v
     B2 Benchmark API               workspace
              |
              v
       AgentCap Core
              |
              v
     B2 Measurements
              |
              +------------+
                           |
                           v
                    B2 Aggregator
                           |
                           v
                    Workload Result
```

The workload runner orchestrates.

B2 measures.

AgentCap core executes/reduces/stores.

Do not merge these responsibilities.

---

# 49. Future Compatibility

Design B3 so later phases can use the workload runner programmatically.

Future callers may include:

```text
repeated-trial runner
benchmark comparison tool
CI benchmark runner
coding-agent benchmark harness
```

Avoid making the workload engine dependent on terminal rendering or global CLI state.

At the same time, do not implement these future systems now.

---

# 50. Documentation

Document:

- workload directory layout
- workload schema
- supported step types
- fixture semantics
- path restrictions
- benchmark-session behavior
- how to run one workload
- how to inspect results
- how to preserve a failed workspace if supported
- how non-zero target exits differ from harness failures

Include at least one complete small workload example.

Do not document future DSL operations as implemented.

---

# 51. Out of Scope

Do NOT implement:

- general-purpose benchmark scripting
- arbitrary shell hooks for mutations
- loops or conditionals
- workload variables/templates
- remote fixtures
- network-dependent workloads
- large repository downloads
- automatic repeated trials
- statistical aggregation
- coding-agent orchestration
- task-success scoring
- LLM judges
- tokenizer integration
- benchmark history
- historical comparison
- CI regression thresholds
- automatic reducer tuning
- semantic source-code benchmarks
- B4+ functionality

Keep B3 focused on deterministic reproducible workloads.

---

# 52. Validation

Before finishing:

1. Build AgentCap.
2. Run all existing tests.
3. Run all B0-B2 benchmark tests.
4. Run all new B3 tests.
5. Validate every bundled workload.
6. Execute each bundled workload at least once.
7. Execute at least one workload twice and verify fresh initial state.
8. Verify benchmark-session isolation between runs.
9. Verify raw/stateless/stateful totals are produced through B2.
10. Verify full/delta/unchanged behavior on a suitable workload.
11. Verify exactly-once execution.
12. Verify non-zero target commands can be expected workflow steps.
13. Verify invalid path traversal is rejected.
14. Verify the source fixture remains unchanged.
15. Verify the current user repository remains unchanged.
16. Verify temporary workspace cleanup.
17. Verify stored command results remain accessible with `acap show`.
18. Verify raw command results remain accessible with `acap raw`.
19. Confirm normal non-benchmark AgentCap behavior remains unchanged.

---

# 53. Acceptance Criteria

Phase B3 is complete when:

```text
[ ] A versioned declarative workload format exists.

[ ] Workloads can contain multiple run steps.

[ ] Workloads can perform minimal deterministic workspace mutations.

[ ] Commands use structured argv semantics.

[ ] Every run step uses the B2 benchmark-session path.

[ ] Every target command executes exactly once.

[ ] Raw/stateless/stateful measurements come from B2.

[ ] A fresh isolated workspace is created for every workload run.

[ ] A fresh benchmark session is created for every workload run.

[ ] User repositories and source fixtures are not mutated.

[ ] Mutation paths cannot escape the benchmark workspace.

[ ] Non-zero target command exits are supported.

[ ] Optional expected exit-status assertions are supported.

[ ] Mutation operations are not counted as AgentCap commands.

[ ] Every run result remains independently stored and recoverable.

[ ] At least one workload exercises stateful delta/unchanged behavior.

[ ] At least one realistic Git/build/test transition is represented.

[ ] Workload definitions are validated before execution where possible.

[ ] Temporary workspaces are cleaned up according to documented semantics.

[ ] No coding-agent or repeated-trial functionality has leaked into B3.
```

---

# 54. Final Report

When implementation is complete, report:

1. files added or modified
2. workload schema introduced
3. schema versioning strategy
4. supported step types
5. CLI syntax implemented
6. workload directory/fixture layout
7. temporary-workspace strategy
8. Git fixture setup strategy, if applicable
9. filesystem safety rules
10. how B3 invokes the B2 benchmark API
11. how benchmark sessions are isolated
12. how exactly-once execution is preserved
13. bundled workloads added
14. tests added
15. example output from each bundled workload
16. existing tests executed and their results
17. any portability/toolchain limitations discovered
18. any architectural issues discovered

Do not make unrelated refactors.
