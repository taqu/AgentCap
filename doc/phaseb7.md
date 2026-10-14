# Phase B7 — Repeated Trials and Statistical Aggregation

## Objective

Extend the Phase B6 coding-agent benchmark so nondeterministic agent behavior can be evaluated across multiple independent trials.

A single LLM-agent run is not sufficient evidence for comparing AgentCap modes because agent behavior may vary between runs.

Phase B7 introduces:

```text
repeated independent trials
        |
        v
per-trial results
        |
        v
statistical aggregation
```

The benchmark must preserve every individual trial while also producing useful aggregate statistics.

The primary goal is to support commands conceptually like:

```bash
acap bench agent \
    --workload bugfix-001 \
    --agent codex \
    --mode integrated \
    --repeat 5
```

and obtain both:

```text
individual trial results
```

and:

```text
aggregate statistics
```

Do not implement cross-result comparison or CI regression policy yet.

Those belong to B8 and B9.

---

# 1. Build Directly on the B6 Trial Model

Treat the existing B6 agent benchmark as the unit of repetition.

Conceptually:

```text
B6:
run one trial
    ->
TrialResult

B7:
repeat B6 trial N times
    ->
[]TrialResult
    ->
AggregateResult
```

Do not create a second agent-execution implementation for repeated benchmarks.

The repeated runner should invoke the same underlying single-trial execution path used by B6.

Prefer an architecture conceptually similar to:

```go
for i := 0; i < repeat; i++ {
    result, err := runner.RunTrial(...)
    ...
}
```

with appropriate workspace reset and result handling.

The exact API should follow the current implementation.

---

# 2. Add `--repeat`

Extend the agent benchmark command with:

```bash
--repeat <N>
```

For example:

```bash
acap bench agent \
    --workload bugfix-001 \
    --agent codex \
    --mode integrated \
    --repeat 5
```

Requirements:

```text
repeat >= 1
```

The default should preserve B6 behavior:

```text
--repeat 1
```

or equivalent.

A command without `--repeat` must continue to represent one trial.

Reject invalid values such as:

```text
--repeat 0
--repeat -1
```

using normal CLI validation conventions.

Do not silently reinterpret invalid values.

---

# 3. Every Trial Must Be Independent

Each trial must begin from the same canonical workload state.

For example:

```text
canonical fixture
      |
      +--> fresh workspace --> trial 1
      |
      +--> fresh workspace --> trial 2
      |
      +--> fresh workspace --> trial 3
      |
      +--> fresh workspace --> trial 4
      |
      +--> fresh workspace --> trial 5
```

Never reuse the mutated workspace from a previous trial.

This is a critical benchmark invariant.

The benchmark is intended to measure agent nondeterminism, not accumulated repository mutations.

Reuse the B6 fresh-workspace mechanism rather than creating a new reset implementation.

---

# 4. Keep Benchmark Configuration Constant Across Trials

Except for values inherently controlled by the external agent/provider, repeated trials must use the same benchmark configuration.

Keep constant:

```text
workload
agent
AgentCap mode
task instruction
fixture starting state
verification logic
integration instructions
timeout policy
benchmark configuration
```

Do not dynamically modify the prompt or AgentCap configuration based on earlier trial outcomes.

For example, never do:

```text
trial 1 fails
    ->
make trial 2 prompt more explicit
```

within the same repeated benchmark run.

That would invalidate the statistical comparison.

---

# 5. Preserve Every Individual Trial

Do not reduce repeated execution directly into aggregate counters.

Every trial must retain its complete B6 result.

Conceptually:

```go
type RepeatedBenchmarkResult struct {
    Trials    []AgentBenchmarkResult
    Aggregate BenchmarkAggregate
}
```

Adapt names and structure to the existing result model.

Each trial should retain the B6 metrics, including where applicable:

```text
outcome / task success
command count

raw bytes
initial visible bytes
show bytes
raw retrieval bytes
total visible bytes

show count
raw retrieval count

AgentCap processing time
wall time
```

and the relevant workload/agent/mode metadata.

The aggregate is a summary over those trial results.

It must never replace them.

---

# 6. Assign Stable Trial Indices

Each trial should have an explicit index or identifier.

For example:

```text
trial 1
trial 2
trial 3
trial 4
trial 5
```

Machine-readable results should expose this explicitly.

Conceptually:

```json
{
  "trial": 1,
  "task_success": true,
  ...
}
```

or an equivalent representation.

Do not rely solely on array position if the existing result architecture supports explicit metadata naturally.

This will make debugging and future B8 analysis easier.

---

# 7. Define Trial Outcomes Explicitly

B6 distinguishes task success from infrastructure failure.

B7 must preserve that distinction.

At minimum, repeated trials should be able to distinguish outcomes conceptually equivalent to:

```text
success
task_failure
timeout
```

and infrastructure-level benchmark failure where appropriate.

Use the outcome model established by B6 rather than introducing conflicting terminology.

Do not collapse:

```text
task failed verification
```

and:

```text
agent could not be launched
```

into the same category.

---

# 8. Define the Repetition Failure Policy

Repeated execution requires an explicit policy for infrastructure failures.

Prefer the following general rule:

```text
valid agent trial outcome
    -> retain trial
    -> continue remaining trials

benchmark infrastructure failure
    -> report clearly
    -> follow established benchmark error policy
```

For example:

```text
agent ran but verification failed
```

is a valid trial result and should not abort the remaining repetitions.

Likewise, if timeout is represented as a valid agent-run outcome by B6, retain it and continue.

However:

```text
fixture cannot be created
workload cannot be loaded
adapter executable is unavailable
benchmark configuration is invalid
```

may make further trials meaningless.

Use the B6 distinction between task outcome and infrastructure error.

Do not blindly continue after errors that invalidate the entire benchmark setup.

---

# 9. Aggregate Success Count

Report:

```text
success_count
trial_count
```

For example:

```text
success: 4/5
```

The numerator must be based on objective B6 task verification.

Do not infer success from agent exit status or agent-generated text.

Machine-readable output should preserve numeric values rather than only:

```text
"success": "4/5"
```

Prefer something conceptually like:

```json
{
  "trial_count": 5,
  "success_count": 4
}
```

A derived success rate may be displayed if useful, but primitive counts are authoritative.

---

# 10. Use Median for Core Continuous Metrics

For repeated LLM-agent benchmarks, use the median as the primary aggregate for skew-sensitive metrics.

At minimum aggregate:

```text
median agent-visible bytes
median wall time
median command count
AgentCap processing overhead
```

Where the B6 result model exposes the more specific metric:

```text
total_visible_bytes
```

use that as the primary AgentCap-side agent-visible cost.

Conceptually:

```text
median_total_visible_bytes
median_wall_time
median_command_count
median_processing_time
```

Use names consistent with the existing schema.

Do not use the arithmetic mean as the sole primary statistic.

LLM-agent workflows can contain outlier runs, so median should be the default summary.

---

# 11. Define Median Behavior Explicitly

Implement median calculation centrally and test it.

For odd sample counts:

```text
1, 3, 9
median = 3
```

For even sample counts, use a clearly defined policy.

Prefer the conventional midpoint:

```text
1, 3, 7, 9
median = 5
```

For integer byte/count metrics, decide how fractional medians are represented.

Prefer preserving precision where the result schema supports it, or use a documented deterministic integer policy.

Do not allow different formatters to calculate median differently.

All aggregate statistics should come from one canonical aggregation layer.

---

# 12. Preserve Raw Measurements Behind the Median

Never discard the distribution.

For example, if trials produce:

```text
total_visible_bytes:

82 KB
85 KB
88 KB
91 KB
410 KB
```

the median may be:

```text
88 KB
```

but the 410 KB trial is important evidence.

Therefore JSON results must retain all five trial results.

Human-readable output may summarize them compactly, but detailed trial data must remain available.

This is one of the core B7 requirements.

---

# 13. Aggregate Recovery Behavior

Report how often progressive-disclosure recovery occurs.

At minimum aggregate:

```text
show/raw rate
```

using explicit definitions.

Avoid an ambiguous single number.

Prefer metrics such as:

```text
trials_with_show
trials_with_raw_retrieval

total_show_count
total_raw_retrieval_count
```

and/or:

```text
show_calls_per_trial
raw_retrievals_per_trial
```

where useful.

For example:

```text
show used:      3/5 trials
raw used:       1/5 trials
show calls:     7 total
raw retrievals: 1 total
```

This is more informative than a vague:

```text
fallback rate: 20%
```

unless that rate has a precisely documented denominator.

---

# 14. Preserve Recovery Byte Metrics

In addition to recovery frequency, retain B5/B6 recovery byte measurements for each trial:

```text
show_bytes
raw_retrieval_bytes
```

Aggregate them where useful.

For example:

```text
median_show_bytes
median_raw_retrieval_bytes
```

may be included if they improve analysis.

However, keep:

```text
median_total_visible_bytes
```

as the primary workflow context-cost summary.

Do not let a large set of secondary aggregates obscure the core result.

---

# 15. Aggregate Command Count

Report:

```text
median_command_count
```

because command behavior may vary substantially between agent runs.

For example:

```text
trial 1: 18 commands
trial 2: 21 commands
trial 3: 19 commands
trial 4: 42 commands
trial 5: 20 commands
```

must not be represented only by total command count.

The workflow-level statistic should describe a typical trial.

Keep individual counts available in trial results.

---

# 16. Aggregate Wall Time

Report:

```text
median_wall_time
```

across trials.

Use complete B6 agent-run wall time.

Do not confuse this with:

```text
AgentCap processing time
```

which remains a separate metric.

Both should remain independently observable.

---

# 17. Aggregate AgentCap Processing Overhead

Aggregate the B6 AgentCap processing metric.

Prefer:

```text
median_processing_time
```

or the equivalent existing unit/field name.

Do not calculate AgentCap overhead by subtracting total wall times between modes.

That comparison would mix:

```text
agent nondeterminism
provider latency
command behavior
AgentCap processing
```

The authoritative AgentCap processing metric is the instrumentation introduced by earlier benchmark phases.

---

# 18. Do Not Aggregate Incompatible Trials Together

One repeated benchmark result must represent one fixed configuration:

```text
one workload
one agent
one AgentCap mode
one benchmark configuration
```

Do not aggregate:

```text
disabled trial
stateful trial
integrated trial
```

into one median.

Likewise, do not mix different workloads or different agents.

Those are separate repeated benchmark groups.

Cross-group comparison belongs to B8.

---

# 19. Avoid Automatic Mode Ranking

B7 should report statistics, not choose a winner.

Do not generate conclusions such as:

```text
integrated is best
stateful wins
disabled loses
```

Do not create a composite score combining:

```text
success
bytes
latency
commands
```

Keep those dimensions independent.

B8 will later display results side by side.

Even B8 should not require a synthetic overall score.

---

# 20. Machine-Readable Result Structure

Extend the existing JSON result format to represent repeated trials.

Conceptually:

```json
{
  "schema_version": 4,

  "workload": "bugfix-001",
  "agent": "codex",
  "mode": "integrated",

  "trial_count": 5,
  "success_count": 4,

  "aggregate": {
    "median_total_visible_bytes": 91020,
    "median_command_count": 24,
    "median_wall_time_ns": 82140000000,
    "median_processing_ns": 12840000,

    "trials_with_show": 3,
    "trials_with_raw_retrieval": 1,
    "total_show_count": 7,
    "total_raw_retrieval_count": 1
  },

  "trials": [
    {
      "trial": 1,
      "task_success": true,
      "command_count": 24,
      "total_visible_bytes": 91020
    },
    {
      "trial": 2,
      "task_success": true,
      "command_count": 21,
      "total_visible_bytes": 87410
    }
  ]
}
```

This is conceptual only.

Use the actual B4–B6 result schema and naming conventions.

Do not automatically use schema version 4 merely because this is Phase B7.

Follow the schema-versioning policy established in B4.

---

# 21. Avoid Duplicating Full Metadata Unnecessarily

The repeated-result structure should distinguish:

```text
group-level metadata
```

from:

```text
trial-specific data
```

For example:

```text
workload
agent
mode
```

normally belong at the repeated-run level if they are invariant.

Do not unnecessarily duplicate large static metadata inside every trial.

However, retain enough trial metadata to diagnose differences and preserve reproducibility.

Prefer clear normalization over aggressive deduplication.

---

# 22. Human-Readable Output

For repeated runs, provide a compact aggregate summary.

For example:

```text
Benchmark: bugfix-001
Agent:     codex
Mode:      integrated
Trials:    5

Task
  success                    4/5

Workflow median
  visible bytes           91,020
  commands                    24
  wall time                 82.1s

Recovery
  show used                3/5 trials
  raw used                 1/5 trials
  show calls                    7
  raw retrievals                1

AgentCap
  median processing        12.8ms
```

Then provide a compact per-trial section where appropriate:

```text
Trial   Result    Visible    Commands   Wall
1       success    91,020       24      82.1s
2       success    87,410       21      79.4s
3       failure   104,820       29      95.2s
4       success    89,300       23      80.7s
5       success   410,220       42     143.8s
```

Follow existing CLI formatting conventions.

The human output should make both:

```text
typical behavior
```

and:

```text
trial variability
```

visible without becoming excessively verbose.

---

# 23. Preserve `--repeat 1` Usability

Do not make the common single-trial case unnecessarily noisy.

If:

```bash
acap bench agent ... --repeat 1
```

is equivalent to normal B6 execution, it may retain the B6-style output.

Internally it may use the repeated-run infrastructure if that simplifies the implementation.

However, do not unnecessarily wrap every single trial in verbose statistical presentation.

Maintain good CLI ergonomics.

---

# 24. Result Persistence

If benchmark results are already persisted by B0–B6, repeated runs must preserve:

```text
the aggregate result
```

and:

```text
the individual trials
```

in a form that later B8 comparison can consume.

Do not persist only the formatted human summary.

The persisted representation must remain structured.

If the existing benchmark architecture distinguishes run IDs and result IDs, follow that model rather than inventing a second persistence namespace.

---

# 25. Trial Identity and Parent Run Identity

Where the persistence model supports it naturally, represent the relationship:

```text
repeated benchmark run
       |
       +--> trial 1
       +--> trial 2
       +--> trial 3
       +--> trial 4
       +--> trial 5
```

A parent run identifier plus trial index is preferable to unrelated anonymous result records.

Do not introduce complex relational infrastructure solely for this purpose if the current store has a simpler natural representation.

The important requirement is that B8 can later recover:

```text
which trials belong to this repeated benchmark run?
```

---

# 26. Handle Partial Repeated Runs

Consider interruption or infrastructure failure after some trials have completed.

For example:

```text
trial 1 complete
trial 2 complete
trial 3 complete
trial 4 infrastructure failure
```

Do not silently present this as:

```text
4 trials completed
```

or as a normal requested four-trial benchmark.

Preserve:

```text
requested trial count
completed/valid trial count
run completion status
```

where appropriate.

If the existing error model means the overall command fails, retain already completed results where the persistence architecture permits it.

Do not fabricate missing trials.

---

# 27. Timeouts Remain Trial Outcomes

If B6 treats agent timeout as a valid trial outcome, B7 must retain and count it explicitly.

For example:

```text
trial 1 success
trial 2 timeout
trial 3 task failure
trial 4 success
trial 5 success
```

should not become merely:

```text
3/5 success
```

in machine-readable data.

The aggregate may report success count, but the individual outcomes must remain available.

Consider exposing aggregate counts such as:

```text
success_count
task_failure_count
timeout_count
```

if they fit the B6 outcome model cleanly.

---

# 28. Be Careful With Statistics Over Failed Trials

Define which trial populations are used for each aggregate metric.

For workflow-cost metrics, the default should generally be:

```text
all valid executed trials with that measurement
```

rather than only successful trials.

Otherwise a mode that fails after consuming large amounts of context could misleadingly appear cheap.

For example:

```text
trial 1 success:  90 KB
trial 2 success:  95 KB
trial 3 failure: 500 KB
```

must not silently report a cost statistic based only on:

```text
90 KB
95 KB
```

unless explicitly labeled as:

```text
successful-trials-only
```

Prefer all valid measured trials for the primary workflow metrics.

If successful-trial-only statistics are useful later, expose them separately and explicitly.

Do not mix the populations.

---

# 29. Do Not Treat Failed Trials as Zero

Never represent missing or failed measurements as:

```text
0 bytes
0 commands
0 seconds
```

unless zero is genuinely observed.

Missing measurement and measured zero are different.

Use the existing result/error representation to distinguish them.

This is especially important for medians.

Do not allow infrastructure failures with absent measurements to artificially lower aggregate statistics.

---

# 30. Centralize Aggregation Logic

Create one aggregation layer that consumes trial results and produces the aggregate.

Conceptually:

```text
[]TrialResult
      |
      v
Aggregator
      |
      v
AggregateResult
```

Both:

```text
human formatter
JSON formatter
```

must consume the same aggregate result.

Never calculate:

```text
median
success count
show/raw rate
```

independently in each formatter.

This preserves B4's separation between measurement and presentation.

---

# 31. Aggregation Must Be Pure Where Practical

Prefer aggregation logic that has no side effects.

Conceptually:

```go
aggregate := AggregateTrials(trials)
```

This makes the statistical behavior easy to unit test.

Do not couple median calculation to:

```text
CLI output
filesystem operations
agent execution
workspace setup
```

The aggregator should operate on completed trial data.

---

# 32. Deterministic Aggregation

Given the same set of trial results, aggregation must always produce the same result.

Do not introduce:

```text
sampling
randomized bootstrap analysis
approximate quantiles
```

in B7.

The initial statistical requirements are deliberately simple.

Exact medians and counts are sufficient.

More advanced statistics can be added later if real benchmark data demonstrates a need.

---

# 33. Do Not Add Statistical Significance Claims

B7 does not need:

```text
p-values
confidence intervals
hypothesis tests
effect-size significance thresholds
```

Do not claim that one mode is statistically significantly better than another.

B7's job is to collect repeated observations and summarize them robustly.

Cross-mode statistical methodology can be considered later if needed.

---

# 34. Sequential Execution First

Prefer running repeated trials sequentially initially:

```text
trial 1
   |
trial 2
   |
trial 3
   |
trial 4
   |
trial 5
```

Do not introduce parallel trial execution unless the existing architecture already supports it safely and there is a strong reason.

Parallel LLM-agent trials introduce additional variables:

```text
provider rate limits
CPU contention
disk contention
build-cache contention
network contention
```

which can distort wall-time measurements.

Sequential execution is simpler and more reproducible for B7.

---

# 35. Avoid Cross-Trial Cache Leakage Where Controllable

Fresh workspace isolation is mandatory.

Also inspect whether benchmark-controlled state can leak between trials, such as:

```text
AgentCap session state
benchmark result state
temporary session identifiers
```

Each trial must receive a fresh AgentCap session where the selected mode requires one.

Do not allow stateful compression history from trial 1 to affect trial 2.

Conceptually:

```text
trial 1
    -> fresh AgentCap session A

trial 2
    -> fresh AgentCap session B
```

Repeated trials test independent workflows.

They are not one long stateful workflow.

---

# 36. Do Not Pretend External Provider State Is Fully Controlled

Some external nondeterminism may remain outside AgentCap's control.

Examples may include:

```text
model updates
provider-side caching
service load
agent implementation changes
```

Record reliable version/configuration metadata where available.

Do not claim bit-for-bit reproducibility of LLM-agent trials.

The purpose of repetition is precisely to observe variability that cannot be eliminated.

---

# 37. Tests

Add focused tests for repeated execution and aggregation.

## Repeat validation

Verify:

```text
repeat = 1
repeat > 1
repeat = 0 rejected
negative repeat rejected
```

---

## Trial count

Using a fake agent, verify:

```text
--repeat 5
```

causes exactly five independent B6 trial executions.

---

## Fresh workspace

Verify mutations from trial 1 are absent when trial 2 starts.

---

## Fresh AgentCap state

For stateful/integrated mode, verify state from one trial does not become the baseline for another trial.

---

## Configuration consistency

Verify each repeated trial receives the same:

```text
workload
agent
mode
task configuration
```

---

## Success aggregation

For outcomes:

```text
success
failure
success
success
failure
```

verify:

```text
trial_count = 5
success_count = 3
```

---

## Median: odd count

For:

```text
10
20
100
```

verify:

```text
median = 20
```

---

## Median: even count

Test the chosen even-sample median policy explicitly.

---

## Outlier behavior

For:

```text
10
11
12
13
1000
```

verify the primary median remains:

```text
12
```

while the `1000` trial remains present in individual results.

---

## Recovery aggregation

Given trial-level recovery behavior, verify:

```text
trials_with_show
trials_with_raw_retrieval
total_show_count
total_raw_retrieval_count
```

are calculated correctly.

---

## Failed-trial accounting

Verify a task-failure trial with valid measurements remains part of the primary workflow-cost aggregation.

Do not silently aggregate only successful trials.

---

## Missing measurements

Verify missing measurements are not interpreted as zero.

---

## Infrastructure failure

Verify the repeated runner follows the defined policy when a trial cannot execute because of benchmark infrastructure failure.

---

## JSON preservation

Verify JSON contains both:

```text
aggregate
trials
```

and that individual trial data is not discarded.

---

## Single-trial compatibility

Verify the normal B6 single-trial path remains usable.

---

# 38. Documentation

Document:

```bash
acap bench agent ... --repeat N
```

Explain why repeated trials are necessary for coding-agent evaluation.

Document the primary aggregate metrics:

```text
success count
median agent-visible bytes
median wall time
median command count
show/raw usage
AgentCap processing overhead
```

Also explicitly document that:

```text
individual trials are preserved
```

and that aggregate medians do not replace raw trial observations.

Document which trials are included in cost medians.

In particular, make clear whether failed but valid trials remain in the primary metric population.

---

# 39. Do Not Implement B8 Yet

Do not implement:

```bash
acap bench compare baseline.json candidate.json
```

B7 operates on repeated trials of one benchmark configuration.

It does not compare:

```text
disabled vs integrated
baseline vs candidate
commit A vs commit B
```

Those belong to B8.

The output structure should make B8 straightforward, but do not implement comparison logic yet.

---

# 40. Do Not Implement B9 Yet

Do not introduce:

```text
regression thresholds
CI pass/fail decisions
automatic benchmark gates
scheduled benchmark execution
```

B7 produces measurements and statistical summaries.

B9 will decide how deterministic benchmark regressions are used in CI.

---

# Implementation Guidance

Before changing code:

1. inspect the B6 single-trial runner;
2. identify the exact result type representing one agent trial;
3. reuse the B6 fresh-workspace setup for every repetition;
4. verify that AgentCap session state is fresh per trial;
5. inspect the benchmark result persistence model;
6. inspect the B4–B6 schema-versioning policy;
7. add a thin repeated-run orchestration layer;
8. implement a pure aggregation layer;
9. extend human/JSON formatting from the canonical aggregate;
10. validate everything with fake-agent tests before running expensive real-agent repetitions.

Prefer:

```text
existing B6 RunTrial()
        |
        v
small RepeatRunner
        |
        v
[]TrialResult
        |
        v
pure Aggregator
        |
        v
RepeatedBenchmarkResult
```

over a broad rewrite of the benchmark runner.

---

# Architectural Invariants

After B7, the benchmark architecture should conceptually be:

```text
                 Repeated Benchmark
                         |
                         v
                    RepeatRunner
                         |
        +----------------+----------------+
        |                |                |
        v                v                v
     Trial 1          Trial 2          Trial N
        |                |                |
        v                v                v
 B6 RunTrial()     B6 RunTrial()     B6 RunTrial()
        |                |                |
        v                v                v
 fresh workspace    fresh workspace    fresh workspace
 fresh session      fresh session      fresh session
        |                |                |
        +----------------+----------------+
                         |
                         v
                   []TrialResult
                         |
                         v
                     Aggregator
                         |
                         v
                AggregateResult
                         |
              +----------+----------+
              |                     |
              v                     v
        human formatter       JSON formatter
```

The central invariants are:

```text
one repetition
=
one independent B6 trial
```

and:

```text
aggregate result
does not replace
individual trial results
```

and:

```text
same benchmark configuration
+
fresh state per trial
```

---

# Primary Statistical Model

For one fixed:

```text
workload × agent × mode
```

B7 should produce:

```text
Trials
  trial 1
  trial 2
  ...
  trial N

Outcomes
  success count
  failure/timeout counts where applicable

Typical workflow
  median total agent-visible bytes
  median command count
  median wall time
  median AgentCap processing time

Recovery
  trials using show
  trials using raw
  total show calls
  total raw retrievals
```

while preserving the complete individual observations behind those summaries.

---

# Exit Criteria

Phase B7 is complete when all of the following are true:

1. `acap bench agent` supports repeated trials through `--repeat N`.
2. The default single-trial behavior remains compatible with B6.
3. Every repetition executes through the existing B6 single-trial path.
4. Every trial starts from a fresh fixture workspace.
5. Every stateful/integrated trial starts with fresh AgentCap state.
6. Benchmark configuration remains constant across repetitions.
7. Every individual trial result is preserved.
8. Trial identity/index is available.
9. Objective B6 task-success results are preserved per trial.
10. Success count is aggregated correctly.
11. Median total agent-visible bytes is available.
12. Median wall time is available.
13. Median command count is available.
14. AgentCap processing overhead is aggregated.
15. `show` and `raw` usage across trials is observable.
16. Failed but valid trials are not silently discarded from workflow-cost statistics.
17. Missing measurements are not treated as zero.
18. Infrastructure failure remains distinguishable from task failure.
19. Aggregate statistics are produced by one canonical aggregation layer.
20. Human and JSON output consume the same aggregate result.
21. Machine-readable output preserves both aggregate statistics and individual trials.
22. Result persistence retains enough information for later B8 comparison.
23. Repeated trials run sequentially unless existing architecture provides a compelling safe alternative.
24. No B8 cross-result comparison is implemented.
25. No B9 CI regression policy is implemented.

The completed Phase B7 should turn:

```text
"the agent used 91 KB in one run"
```

into evidence of the form:

```text
5 independent trials

task success:          4/5
median visible bytes:  91 KB
median commands:       24
median wall time:      82 s

show used:             3/5 trials
raw used:              1/5 trials

individual trials:
    retained
```

so later phases can compare benchmark configurations without pretending that one nondeterministic LLM-agent run is representative.
