# Phase 7 — Agent Workflow Benchmark and Evaluation

## Objective

Build and run a reproducible benchmark that measures whether AgentCap materially improves coding-agent workflows.

The benchmark must evaluate AgentCap as a complete system, not merely as a terminal-output compressor.

The primary question is:

```text
Does AgentCap reduce total coding-agent context cost
without materially harming task success,
correctness,
or workflow efficiency?
```

This phase should produce evidence that can guide future AgentCap development.

Do not add major new AgentCap capabilities during this phase unless a benchmark-blocking correctness bug is discovered.

---

# 1. Starting Point

Phase 7 assumes the following are already complete:

```text
AgentCap execution core
result capture
persistent storage
progressive disclosure
session-aware compression
Git-aware reduction
build/test-aware reduction
transparent integration
benchmark CLI
Phase 6.5A common adapter contract
Phase 6.5B Claude Code adapter
Phase 6.5C Antigravity adapter
Phase 6.5D cross-agent conformance
```

The adapters must already be benchmark-ready.

Do not use Phase 7 to finish basic adapter correctness.

If Phase 6.5D conformance is failing, fix that first.

---

# 2. Primary Benchmark Question

The benchmark must measure the effect of AgentCap on complete coding workflows.

Do not optimize for:

```text
compression ratio of one command
```

The primary unit of evaluation is:

```text
one complete coding task
```

from:

```text
task start
```

through:

```text
final answer / task completion
```

including all intermediate shell/tool interaction.

---

# 3. Primary KPI

The primary AgentCap efficiency metric is:

```text
total agent-visible command-result tokens
over the complete coding workflow
```

This includes:

```text
initial AgentCap capsules
+
stateful delta results
+
acap show output
+
acap raw output
+
AgentCap-specific integration/instruction overhead
```

Do not report only initial compression ratio.

A run that returns tiny initial capsules but causes repeated raw retrieval may be worse than a moderately compressed run.

---

# 4. Task Success Is a Hard Constraint

Token reduction alone is not sufficient.

For every benchmark configuration measure whether the coding task was successfully completed.

AgentCap must not be considered beneficial if lower context consumption causes materially worse task completion.

Therefore always evaluate:

```text
efficiency
+
task outcome
```

together.

Do not optimize one while ignoring the other.

---

# 5. Benchmark Configurations

Support at least these AgentCap conditions:

```text
A. AgentCap OFF

B. AgentCap stateless

C. AgentCap stateful

D. AgentCap fully integrated
```

Conceptually:

```text
A:
normal agent command output

B:
AgentCap reduction
no cross-command state/delta

C:
AgentCap reduction
+
session state
+
delta behavior

D:
normal agent shell workflow
+
production adapter
+
full AgentCap integration
```

If the current implementation cannot cleanly expose all four modes, prioritize:

```text
OFF
vs
FULL
```

first.

Do not invent unreliable modes merely to complete the matrix.

---

# 6. Agent Matrix

Run Phase 7 against at least:

```text
Claude Code
Antigravity
```

using the production adapters from Phase 6.5.

Conceptually:

```text
                OFF    STATELESS    STATEFUL    FULL

Claude Code      ✓         ✓            ✓         ✓
Antigravity      ✓         ✓            ✓         ✓
```

If some intermediate AgentCap mode is not supported by one adapter, document that explicitly.

Do not create adapter-specific benchmark logic that changes core AgentCap behavior.

---

# 7. Model Configuration Control

Within a benchmark comparison, keep the agent/model configuration fixed.

For example:

```text
Claude configuration X
AgentCap OFF
vs
Claude configuration X
AgentCap ON
```

Do not compare:

```text
different model
+
different AgentCap mode
```

and attribute the difference to AgentCap.

Record the model/version information used for every run.

---

# 8. Agent Configuration Control

Hold relevant agent settings constant across ON/OFF comparisons.

Examples include:

```text
model
reasoning/thinking mode
permission mode
tool configuration
system instructions
project instructions
temperature or equivalent settings if exposed
```

Only AgentCap mode should intentionally differ.

If host-agent configuration must differ because integration requires it, record the difference explicitly.

---

# 9. Production Integration Only

FULL mode must use the actual production adapter.

Do not benchmark using:

```text
special benchmark wrappers
manual acap run prefixes
synthetic execution paths
mock hooks
```

unless the benchmark specifically evaluates those modes.

The main result should reflect real user-facing integration behavior.

---

# 10. OFF Mode

Prefer the common Phase 6.5 bypass mechanism for AgentCap OFF.

Do not repeatedly uninstall and reinstall hooks between trials if bypass can provide equivalent host-agent behavior.

OFF mode must disable:

```text
AgentCap presentation
stateful delta
result interception
AgentCap-specific command rewriting
```

while leaving unrelated agent configuration unchanged.

---

# 11. Clean Trial State

Every independent benchmark trial must start from a known state.

Reset at least:

```text
repository contents
Git working tree
generated files
.acap state
AgentCap session state
benchmark metadata
temporary files
```

unless the benchmark explicitly tests persistent-session behavior.

Do not let previous runs influence later independent trials.

---

# 12. Benchmark Task Categories

Include representative coding-agent workflows.

At minimum cover:

```text
bug fixing
test failure repair
compile-error repair
small feature implementation
refactoring
Git review / diff inspection
large repository search
```

These categories should exercise different AgentCap reducers and interaction patterns.

---

# 13. Task Selection Principle

Prefer tasks that generate meaningful command interaction.

A benchmark task should normally require some combination of:

```text
repository exploration
search
build
test
Git inspection
repeated command execution
diagnostic interpretation
```

Avoid tasks that can be completed almost entirely from one small source file without shell interaction.

Such tasks provide little information about AgentCap.

---

# 14. Task Realism

Benchmark tasks should resemble real coding-agent work.

Avoid artificial tasks whose only purpose is to generate huge output.

For example, prefer:

```text
fix this failing parser test
```

over:

```text
print 100,000 lines and compress them
```

Synthetic output tests remain useful for microbenchmarks but should not dominate Phase 7 conclusions.

---

# 15. Repository Scale

Use more than one repository size where practical.

Useful categories include:

```text
small fixture repository
medium real-world-style repository
large repository
```

The small fixture is useful for determinism.

The larger repositories are needed to expose:

```text
search cost
Git diff cost
test-suite output
repeated exploration
```

---

# 16. Language Coverage

Prioritize languages/toolchains already well supported by AgentCap.

Recommended core set:

```text
Go
C/C++
Rust
```

These directly exercise existing build/test reducers.

Do not add new reducer implementations merely to expand benchmark language coverage.

Additional languages may be used only if AgentCap generic reduction is intentionally being evaluated.

---

# 17. Deterministic Task Fixtures

Where possible, create deterministic task fixtures.

Each task should specify:

```text
starting revision
task description
expected correct outcome
validation command
allowed file scope if needed
timeout
```

Example:

```text
repository revision: <commit>
task:
    fix failing parser handling for empty input

validation:
    go test ./...

expected:
    tests pass
```

Keep task definitions version-controlled.

---

# 18. Ground Truth

Every benchmark task must have a reliable success oracle.

Prefer:

```text
automated tests
build success
structured validation script
known Git diff constraints
```

Do not rely solely on the agent claiming:

```text
done
```

The benchmark runner should independently verify success.

---

# 19. Hidden Validation

Where useful, include validation not directly exposed to the coding agent.

This helps detect:

```text
overfitting to visible tests
partial fixes
incorrect implementation
```

However, hidden validation must evaluate task correctness, not arbitrary style preferences.

---

# 20. Run Repetition

Agent behavior is nondeterministic.

Do not draw conclusions from a single run.

Run each important condition multiple times.

At minimum, structure the benchmark so repetition count is configurable.

Conceptually:

```text
task
x
agent
x
AgentCap mode
x
N repetitions
```

Use enough repetitions to reveal large variance.

Do not hard-code conclusions from one lucky or unlucky trajectory.

---

# 21. Randomization

Where practical, randomize or rotate run order.

Avoid always running:

```text
OFF first
ON second
```

because:

```text
machine cache
repository cache
agent service variation
environment warmup
```

could bias results.

Record run order.

---

# 22. Reproducibility

Every run should record enough information to reproduce it.

At minimum:

```text
benchmark task ID
repository revision
agent
agent version
model
model configuration
AgentCap version
adapter version
AgentCap mode
OS/platform
start time
run identifier
```

Do not depend on manually remembered environment state.

---

# 23. Primary Metrics

For every task run capture:

```text
task success
total agent-visible command-result tokens
total agent-visible command-result bytes
```

These are the primary AgentCap metrics.

If reliable model input token counts are available from the agent, record them as well.

Do not fabricate model-token estimates if the host does not expose them reliably.

---

# 24. Secondary Metrics

Also capture:

```text
wall-clock time
number of shell commands
number of repeated commands
number of AgentCap results
number of acap show calls
number of acap raw calls
raw fallback rate
unchanged-result rate
delta-result rate
full-result rate
commands intercepted
commands bypassed
adapter failures
AgentCap processing latency
adapter latency
```

Where available also capture:

```text
model input tokens
model output tokens
```

Keep host-reported token metrics separate from locally estimated metrics.

---

# 25. Raw Fallback Rate

A critical metric is:

```text
how often does the agent request acap raw?
```

Measure:

```text
raw calls / AgentCap results
```

and:

```text
tasks containing at least one raw fallback
```

Also measure immediate raw fallback.

For example:

```text
AgentCap returns result
next relevant command is acap raw <same-id>
```

This is a strong signal that the capsule may have been too aggressive.

---

# 26. Drill-Down Rate

Measure targeted recovery separately from full raw recovery.

Examples:

```text
acap show <id> --file
acap show <id> --test
acap show <id> --errors
acap show <id> --warnings
```

Targeted drill-down is expected behavior.

Do not treat all `show` calls as failures.

The useful distinction is:

```text
progressive targeted recovery
vs
full raw fallback
```

---

# 27. Information Recovery Cost

Calculate:

```text
initial AgentCap output
+
all drill-down output
+
all raw output
```

for each result/workflow.

This gives actual AgentCap context cost.

Do not measure only the first response.

---

# 28. Stateful Value

Explicitly measure the benefit of stateful compression.

For repeated commands, compare:

```text
stateless returned bytes/tokens
```

against:

```text
stateful returned bytes/tokens
```

Measure especially:

```text
repeated tests
repeated builds
repeated git diff
repeated repository search
```

This should quantify whether session/delta logic materially contributes beyond basic reduction.

---

# 29. Unchanged Result Rate

Track:

```text
number of unchanged results
```

and:

```text
bytes/tokens avoided through unchanged representation
```

Repeated no-change commands are one of the clearest stateful-compression opportunities.

---

# 30. Delta Result Rate

Track:

```text
results presented as delta
```

and compare their returned size to:

```text
full AgentCap capsule
raw result
```

Do not assume all deltas are useful.

Also measure whether delta results cause increased drill-down.

---

# 31. Full Fallback Rate

Track how often AgentCap cannot apply a specialized/stateful representation and returns a full or generic result.

This can reveal:

```text
unsupported command patterns
classification failures
unsafe delta situations
```

Do not automatically treat full fallback as a bug.

Conservative fallback may be correct.

---

# 32. Command Distribution

Record which command families dominate model-visible context.

At minimum categorize:

```text
search
file listing
Git
build
test
compiler
generic shell
AgentCap drill-down
other
```

This is important for determining future optimization priorities.

---

# 33. Per-Command Cost

For each command/result record:

```text
raw bytes
raw estimated tokens
AgentCap returned bytes
AgentCap returned estimated tokens
drill-down bytes/tokens
raw recovery bytes/tokens
processing latency
```

This enables workflow-level attribution later.

---

# 34. Token Estimation

If exact tool-output token counts are unavailable, use one consistent tokenizer/estimation method across all configurations.

Do not compare:

```text
exact host token count
```

against:

```text
character/4 estimate
```

as though they are the same metric.

Clearly label:

```text
measured tokens
```

versus:

```text
estimated tokens
```

---

# 35. Byte Metrics

Always retain raw byte counts even if token estimates are available.

Bytes provide:

```text
deterministic
model-independent
easy-to-reproduce
```

measurements.

Token metrics are useful but may change with model/tokenizer.

---

# 36. Wall-Clock Time

Measure complete task duration:

```text
task start
->
task completion/failure
```

Also retain command-level timing.

Distinguish:

```text
agent reasoning time
command execution time
AgentCap processing latency
adapter latency
```

where technically possible.

Do not attribute total wall-clock differences directly to AgentCap without examining these components.

---

# 37. AgentCap Processing Latency

Measure AgentCap processing separately.

For each intercepted result, record:

```text
capture overhead
storage overhead
reduction overhead
state/delta overhead
presentation overhead
```

where existing instrumentation makes this practical.

Do not add excessive instrumentation complexity if only total core processing latency is currently available.

---

# 38. Adapter Latency

Keep adapter latency separate from AgentCap core latency.

For example:

```text
Claude adapter overhead
Antigravity adapter overhead
AgentCap core overhead
```

This is necessary for diagnosing integration-specific regressions.

---

# 39. Command Count

Compare the total number of shell commands per workflow.

AgentCap may reduce output yet accidentally cause agents to issue more commands.

Therefore evaluate:

```text
command count
+
context cost
```

together.

---

# 40. Repeated Command Count

Track repeated equivalent commands.

Examples:

```text
go test ./...
go test ./...
go test ./...
```

or:

```text
git diff
git diff
```

This is where AgentCap stateful behavior should provide disproportionate benefit.

---

# 41. AgentCap Retrieval Commands

Track:

```text
acap show
acap raw
```

as their own command category.

They contribute to:

```text
command count
wall-clock time
agent-visible output
```

Do not exclude them from AgentCap cost.

---

# 42. Task Success Categories

Use explicit task outcomes.

Recommended:

```text
PASS
FAIL
TIMEOUT
INFRA_ERROR
```

Optionally:

```text
PARTIAL
```

only if the task has a deterministic partial-success definition.

Do not classify based on subjective impressions.

---

# 43. Infrastructure Failure

Separate benchmark infrastructure failures from agent failures.

Examples:

```text
agent API/service unavailable
hook startup failure
benchmark runner crash
repository setup failure
machine resource failure
```

These should not count as task failures.

Record them as:

```text
INFRA_ERROR
```

and rerun where appropriate.

---

# 44. Adapter Failure

If the benchmark detects:

```text
adapter crash
incorrect interception
duplicate execution
permission failure
```

classify it explicitly.

Do not silently include adapter-corrupted runs in performance aggregates.

Phase 6.5D should make these rare.

---

# 45. Timeout Policy

Define a per-task timeout.

Use the same timeout for AgentCap ON and OFF comparisons.

A timeout counts as a task outcome, not simply missing data.

Do not allow AgentCap ON runs extra time merely because integration adds overhead.

---

# 46. Correctness Incidents

Track cases where AgentCap appears to hide required information.

Examples:

```text
agent makes incorrect inference from compressed output
agent repeatedly searches for omitted error
agent requests raw immediately
agent misses a diagnostic available in raw output
```

These should be recorded as benchmark observations.

Do not infer causality automatically.

Flag them for manual review.

---

# 47. Over-Compression Review

For failed or anomalous AgentCap runs, inspect stored raw results.

Determine whether:

```text
necessary information was omitted
```

or:

```text
agent simply chose a different strategy
```

Do not label every failed ON run as over-compression.

Use concrete evidence from stored results.

---

# 48. Compare Paired Runs Carefully

The strongest comparison is:

```text
same task
same agent
same model/settings
same repository state
AgentCap OFF
vs
AgentCap ON
```

Use this as the primary analytical unit.

Cross-agent comparisons are secondary.

Do not infer AgentCap quality by comparing Claude ON against Antigravity OFF.

---

# 49. Claude and Antigravity Are Separate Cohorts

Report results separately for:

```text
Claude Code
```

and:

```text
Antigravity
```

before aggregating.

Agent behavior may differ substantially.

A feature useful for one agent may be neutral for another.

---

# 50. Do Not Rank Agents

Phase 7 evaluates AgentCap.

Do not turn the benchmark into:

```text
Claude vs Antigravity leaderboard
```

unless a future project explicitly requires that.

Cross-agent differences are useful for understanding integration behavior, not for selecting a winner.

---

# 51. Workflow-Level Aggregation

For each run compute:

```text
raw command-result cost
AgentCap-visible cost
recovery cost
total AgentCap context cost
```

Conceptually:

```text
AgentCap total cost =
    capsules
  + deltas
  + show output
  + raw output
  + AgentCap instruction overhead
```

This is the main workflow metric.

---

# 52. Savings Metric

A useful derived metric is:

```text
context savings =
1 - AgentCap total context cost / baseline raw context cost
```

Report this only alongside:

```text
task success
```

Do not present savings alone as product success.

---

# 53. Success-Normalized Analysis

When comparing modes, separate:

```text
all runs
```

from:

```text
successful runs only
```

Both views are useful.

Successful-only analysis measures efficiency when the task succeeds.

All-run analysis reveals whether compression changes success rates.

Do not hide failures by reporting successful runs only.

---

# 54. Latency Analysis

Report:

```text
median
distribution
```

rather than only mean when enough runs exist.

Coding-agent timings often contain large outliers.

Do the same for:

```text
total tokens
command count
raw fallback rate
```

where useful.

---

# 55. Avoid Premature Statistical Complexity

Do not begin with sophisticated statistical modeling.

Start with:

```text
counts
medians
means where useful
percentiles
paired differences
success rates
```

Only add advanced significance analysis if the dataset becomes large enough to justify it.

---

# 56. Benchmark CLI Integration

Use the existing benchmark CLI as the central orchestration and reporting interface.

Do not build a second benchmark runner.

Extend it only where Phase 7 requires missing capabilities.

Keep:

```text
task definition
run orchestration
metrics collection
result storage
report generation
```

coherent within the existing benchmark architecture.

---

# 57. Benchmark Result Storage

Store raw benchmark records in a machine-readable format.

Each run should have a stable identifier.

A record should contain conceptually:

```text
run_id
task_id

agent
agent_version
model
model_settings

agentcap_mode
agentcap_version
adapter
adapter_version

repository
revision

task_outcome

timing metrics
token/byte metrics
command metrics
AgentCap metrics

artifact/result references
```

Use existing benchmark schemas where available.

Do not introduce duplicate representations unnecessarily.

---

# 58. Preserve Raw Benchmark Evidence

Keep enough raw evidence to investigate anomalous runs.

This may include:

```text
agent transcript or available tool log
command sequence
AgentCap result IDs
benchmark runner logs
final repository diff
validation output
```

Respect host-agent storage constraints.

Do not rely only on aggregated CSV rows.

---

# 59. Command Trace

Record an ordered command trace for each benchmark run.

Example:

```text
1. rg Resolve .
2. sed ...
3. go test ./...
4. acap show ...
5. edit
6. go test ./...
7. git diff
```

Include per-command AgentCap mode/result metadata where available.

This is essential for understanding why a run consumed context.

---

# 60. AgentCap Result Trace

For intercepted commands record:

```text
result ID
reducer
presentation type
raw size
returned size
baseline result if any
full/delta/unchanged
drill-down relationship
```

This allows workflow reconstruction without parsing model text.

---

# 61. Benchmark Report — Overview

Generate a summary report containing at least:

```text
number of tasks
number of runs
agents
AgentCap modes
success rates
total context cost
context savings
wall-clock time
command count
raw fallback rate
```

Keep raw detailed data separately available.

---

# 62. Benchmark Report — Per Agent

Provide separate sections for:

```text
Claude Code
Antigravity
```

For each report:

```text
OFF vs ON task success
context cost
wall time
command count
drill-down rate
raw fallback rate
AgentCap overhead
```

Do not hide agent-specific differences inside one aggregate.

---

# 63. Benchmark Report — Per Task Category

Break down results by categories such as:

```text
bug fix
test repair
compile repair
feature
refactor
Git inspection
search
```

This may reveal that AgentCap is highly useful in some workflows and neutral in others.

---

# 64. Benchmark Report — Command Families

Report which command families account for the largest amount of:

```text
raw output
AgentCap output
remaining context cost
raw fallback
```

This should directly inform the next optimization phase.

---

# 65. Benchmark Report — Stateful Value

Include a specific section comparing:

```text
stateless
vs
stateful
```

where data is available.

Measure:

```text
total returned context
unchanged savings
delta savings
drill-down change
task success change
```

This determines whether Phase 3-style session state is actually valuable in realistic workflows.

---

# 66. Benchmark Report — Drill-Down

Report:

```text
show calls
raw calls
targeted drill-down rate
immediate raw fallback rate
```

Identify which reducer/result types generate the most recovery requests.

This is a key signal for over-compression.

---

# 67. Benchmark Report — Latency

Separate:

```text
task wall time
AgentCap processing overhead
adapter overhead
```

Do not report one combined number only.

Include enough detail to determine whether latency is noticeable in interactive workflows.

---

# 68. Benchmark Report — Failure Review

Include a small manually reviewed section for anomalous cases.

Examples:

```text
AgentCap ON failed while OFF succeeded
ON required immediate raw recovery
adapter bypassed unexpectedly
compression omitted critical diagnostic
AgentCap produced misleading delta
```

Provide concrete run IDs.

Do not generalize from one incident without broader evidence.

---

# 69. No Product Claims Before Data

Do not encode expectations such as:

```text
AgentCap should save 80%
```

into pass/fail logic.

Phase 7 exists to discover the actual effect.

The benchmark must be capable of showing:

```text
large improvement
small improvement
neutral result
regression
```

without bias.

---

# 70. Do Not Optimize During Initial Baseline

Run an initial benchmark baseline before tuning reducers.

Otherwise there is no stable starting point.

Recommended sequence:

```text
baseline benchmark
    ->
analyze
    ->
identify major losses
    ->
future optimization phase
```

Do not repeatedly modify AgentCap halfway through one benchmark dataset and aggregate the results together.

---

# 71. Freeze Versions During a Benchmark Batch

For one benchmark batch, freeze:

```text
AgentCap commit/version
adapter version
benchmark task version
agent configuration
model configuration
```

If code changes, start a new benchmark batch.

Do not mix incompatible revisions in one aggregate without explicit labeling.

---

# 72. Benchmark Batch Identity

Assign a benchmark batch ID.

Example concept:

```text
2026-10-phase7-baseline-01
```

The exact format is implementation-specific.

Every run should reference its batch.

---

# 73. Warmup

If the host or repository has meaningful one-time startup costs, consider a non-measured warmup step.

However, keep warmup identical for ON and OFF.

Do not warm AgentCap ON more aggressively than OFF.

Document the warmup policy.

---

# 74. Build Cache Policy

Decide explicitly how to handle:

```text
Go build cache
Rust cargo cache
C/C++ build artifacts
```

Two reasonable benchmark types exist:

```text
cold build
warm incremental workflow
```

Do not accidentally mix them.

For coding-agent workflows, warm incremental behavior may often be more realistic.

Record the policy per task.

---

# 75. Repository Reset

Use a reliable reset mechanism after every independent trial.

Conceptually:

```text
git reset --hard <revision>
git clean ...
```

or isolated worktrees/clones.

Be careful not to remove benchmark infrastructure or AgentCap integration configuration accidentally.

Prefer isolated worktrees or disposable copies if they improve safety.

---

# 76. AgentCap State Reset

Independent runs must start without stale:

```text
.acap/store.db
.acap/objects
session relationships
baseline relationships
```

unless state persistence itself is the benchmark target.

Provide a deterministic reset command/helper.

---

# 77. Session Scope

For one coding task, preserve one AgentCap session across the workflow.

Do not reset AgentCap state between commands.

That would invalidate stateful evaluation.

Reset only between independent benchmark trials.

---

# 78. Progressive Disclosure Is Part of the Benchmark

Do not prevent the agent from using:

```text
acap show
acap raw
```

AgentCap's intended design includes recovery.

The benchmark must measure natural drill-down behavior.

Artificially forbidding raw retrieval would exaggerate compression gains.

---

# 79. Do Not Force Agents to Use Drill-Down

Likewise, do not instruct agents:

```text
always use acap show before raw
```

unless this is part of actual production integration instructions.

Observe natural behavior.

The goal is to measure AgentCap as deployed.

---

# 80. Integration Instruction Overhead

If Claude or Antigravity receives AgentCap-specific model instructions, count those tokens as AgentCap cost.

Do not exclude them merely because they are not shell output.

Conceptually:

```text
AgentCap context cost =
tool-result cost
+
required integration instruction cost
```

---

# 81. Agent Behavior Changes

Record whether AgentCap changes command behavior.

Examples:

```text
more commands
fewer commands
more repeated commands
more targeted commands
more raw recovery
```

This may be more important than per-command compression ratio.

---

# 82. Information-Loss Investigation

When an AgentCap run performs worse, inspect:

```text
raw command result
AgentCap presentation
agent next action
```

Determine whether the omitted information plausibly contributed.

Classify cautiously.

Suggested labels:

```text
no evidence of information loss
possible information loss
probable information loss
confirmed reducer correctness bug
```

Keep these as analysis labels, not automatic benchmark metrics.

---

# 83. Reducer Attribution

When a problematic result is identified, record which reducer produced it.

For example:

```text
generic
Git
Go test
Rust build
compiler diagnostics
search
```

This allows future optimization to target the highest-impact reducer.

---

# 84. Benchmark Exit Criteria

Phase 7 does not need to prove AgentCap is beneficial.

It needs to produce reliable evidence.

Phase 7 is complete when:

```text
benchmark runs are reproducible
metrics are trustworthy
ON/OFF comparison works
success is independently validated
full workflow cost is measured
drill-down cost is included
adapter/core overhead is separated
results can identify major context consumers
```

---

# 85. Required Initial Benchmark Set

Before expanding the benchmark suite, complete a smaller baseline set.

Recommended initial matrix:

```text
Agents:
    Claude Code
    Antigravity

Modes:
    OFF
    FULL

Task categories:
    test repair
    compile-error repair
    bug fix
    Git/diff inspection
    repository search

Languages:
    Go
    C/C++
    Rust
```

Use a manageable number of deterministic tasks first.

Only expand after the measurement pipeline is validated.

---

# 86. Phase 7A — Benchmark Harness Validation

Treat the first subphase as measurement validation.

Goals:

```text
verify clean reset
verify ON/OFF switching
verify metrics
verify task validation
verify command traces
verify AgentCap trace
verify reproducibility
```

Run a small number of tasks.

Do not interpret product performance yet.

---

# 87. Phase 7B — Baseline Benchmark

Once the harness is trustworthy, run the first fixed benchmark batch.

Use frozen:

```text
AgentCap version
adapters
tasks
agent/model settings
```

Do not modify reducers during the batch.

This establishes the reference dataset.

---

# 88. Phase 7C — Analysis

Analyze:

```text
success
context savings
command count
latency
drill-down
raw fallback
stateful behavior
command-family costs
```

Identify the largest remaining context sources.

Also identify any task categories where AgentCap appears harmful.

---

# 89. Phase 7D — Prioritization Output

The final Phase 7 output should not immediately implement fixes.

Produce a prioritized evidence list.

Example:

```text
High impact:
    repeated C++ diagnostics still dominate context

High risk:
    Git delta causes frequent raw recovery

Low value:
    ls compression contributes negligible total savings
```

Priorities must be driven by measured benchmark data.

---

# 90. Future Work Gate

Do not begin semantic source navigation, AST integration, or additional reducers solely because they appear attractive.

Require benchmark evidence such as:

```text
large source-file output is a dominant remaining context cost
```

before adding major complexity.

This preserves AgentCap's narrow architecture.

---

# 91. Potential Decision Outcomes

Phase 7 should be able to support decisions such as:

```text
keep current architecture
tune specific reducer
change delta presentation
reduce compression aggressiveness
improve progressive disclosure
optimize adapter latency
add support for one high-cost command family
investigate semantic source integration
```

Do not predetermine which outcome is correct.

---

# 92. Regression Benchmark

Once the baseline is stable, retain it as a regression suite.

Future AgentCap changes should be comparable against:

```text
Phase 7 baseline
```

where practical.

At minimum preserve:

```text
task definitions
result schema
key metrics
benchmark batch metadata
```

---

# 93. No Benchmark Gaming

Do not special-case:

```text
benchmark repositories
task IDs
fixture commands
expected edits
```

inside AgentCap.

Do not modify model instructions to reveal solutions.

Do not preload task answers.

The benchmark should exercise production behavior.

---

# 94. No Hidden Compression Exclusions

Do not exclude inconvenient AgentCap costs such as:

```text
show calls
raw calls
hook instructions
adapter overhead
failed reductions
```

from the final totals.

If the agent sees or pays the cost, it belongs in the evaluation.

---

# 95. No Double Counting

Be careful not to double count:

```text
AgentCap capsule bytes
```

both as:

```text
tool output
```

and:

```text
AgentCap metrics output
```

Define metric boundaries clearly.

Benchmark accounting should be auditable.

---

# 96. Data Validation

Add sanity checks for benchmark records.

Examples:

```text
returned bytes <= impossible limits
command count matches trace length
show/raw counts match command trace
intercepted + bypassed counts are plausible
success validation exists
AgentCap mode recorded
```

Reject or flag corrupted records before aggregation.

---

# 97. Benchmark Diagnostics

Provide a way to inspect one run deeply.

Conceptually:

```text
benchmark show <run-id>
```

or equivalent existing CLI behavior.

It should expose:

```text
task metadata
command timeline
AgentCap results
metrics
validation result
final outcome
```

Do not force manual database inspection for ordinary analysis.

---

# 98. Comparison Command

Provide or extend a benchmark comparison operation.

Conceptually:

```text
benchmark compare <run-a> <run-b>
```

or:

```text
benchmark compare --task <task> --mode off --mode full
```

Use existing benchmark CLI conventions.

The output should emphasize:

```text
success difference
context-cost difference
command-count difference
time difference
drill-down difference
```

---

# 99. Batch Summary

Provide a batch-level summary.

Conceptually:

```text
benchmark report <batch>
```

The report should make it easy to answer:

```text
Did AgentCap reduce context?

Did success change?

Where did savings come from?

Where was context still spent?

How often did agents recover raw output?

What was AgentCap's latency cost?
```

---

# 100. Acceptance Criteria

Phase 7 is complete only when all of the following are true.

## Harness

- benchmark CLI can run controlled AgentCap OFF/ON trials.
- Claude Code and Antigravity production adapters are used.
- each run starts from deterministic repository and AgentCap state.
- benchmark batches are reproducible.

## Correctness

- task success is independently validated.
- infrastructure failures are separated from task failures.
- adapter failures are identifiable.
- timeouts are handled consistently.

## Metrics

Every run records at least:

```text
task success
agent
model/configuration
AgentCap mode
total agent-visible command-result bytes/tokens
wall-clock time
shell command count
repeated command count
show calls
raw calls
commands intercepted
commands bypassed
adapter failures
AgentCap processing latency
adapter latency
```

where technically measurable.

## AgentCap accounting

- initial capsules are counted.
- delta/unchanged output is counted.
- targeted drill-down is counted.
- raw retrieval is counted.
- required AgentCap instruction overhead is counted.
- costs are not double-counted.

## Stateful evaluation

- repeated-command workflows are represented.
- unchanged rate is measurable.
- delta rate is measurable.
- stateful vs stateless cost can be compared where supported.

## Reports

The benchmark can report:

```text
per-run metrics
per-task comparisons
per-agent results
per-task-category results
command-family costs
raw fallback behavior
latency overhead
```

## Evidence quality

- benchmark versions/configuration are frozen per batch.
- anomalous runs can be inspected.
- raw benchmark evidence is retained.
- results do not depend on one trial only.
- AgentCap is not benchmark-special-cased.

## Initial baseline

A fixed baseline batch has been completed for:

```text
Claude Code
Antigravity
```

with at least:

```text
AgentCap OFF
AgentCap FULL
```

across multiple realistic coding-task categories.

## Final output

Phase 7 produces a data-driven list of:

```text
measured strengths
measured regressions
largest remaining context costs
highest raw-fallback sources
highest-latency components
highest-value next optimization candidates
```

without implementing those optimizations yet.

---

# 101. Implementation Order

Use this order:

```text
1. Inspect the existing benchmark CLI.

2. Inspect Phase 6.5 adapter metrics and ON/OFF controls.

3. Define the benchmark run schema.

4. Define deterministic task specifications.

5. Implement clean repository/state reset.

6. Implement task validation.

7. Capture command traces.

8. Capture AgentCap result traces.

9. Capture context-cost metrics.

10. Capture drill-down/raw-recovery metrics.

11. Capture adapter/core latency separately.

12. Add per-run inspection.

13. Add paired OFF/ON comparison.

14. Add batch summaries.

15. Validate the harness on a very small fixture set.

16. Fix measurement bugs.

17. Freeze the Phase 7 baseline configuration.

18. Run repeated Claude Code baseline trials.

19. Run repeated Antigravity baseline trials.

20. Analyze results.

21. Produce Phase 7 findings and next-phase priorities.
```

Do not tune AgentCap between steps 17 and 20.

---

# 102. Design Priority

When benchmark design tradeoffs arise, use this order:

```text
1. Correct task validation.
2. Reproducible runs.
3. Correct AgentCap ON/OFF isolation.
4. Complete workflow-level cost accounting.
5. Accurate command/result traces.
6. Separation of adapter and core overhead.
7. Natural agent behavior.
8. Statistical usefulness.
9. Benchmark throughput.
```

A smaller trustworthy benchmark is preferable to a large noisy benchmark.

---

# 103. Final Deliverable

At the end of Phase 7, AgentCap should have a reproducible evaluation system capable of answering:

```text
Does AgentCap reduce the amount of command-result context
that coding agents consume across realistic workflows?

Does that reduction preserve task success?

How much additional latency does AgentCap introduce?

How much of the theoretical compression benefit is lost
through drill-down or raw recovery?

How much additional value comes from stateful compression?

Which command families still dominate context consumption?

Which reducers or presentation strategies should be improved next?
```

The result of Phase 7 should be measurement, not another feature expansion.

Future development should be driven by those measurements.
