# Phase B8 — Benchmark Result Comparison

## Objective

Implement comparison of previously produced benchmark results.

Phase B0–B7 already provide:

```text
benchmark execution
measurement
stable result schemas
recovery-cost accounting
real coding-agent trials
repeated trials
statistical aggregation
result persistence
```

Phase B8 must add a read-only comparison layer that can answer:

```text
How did two benchmark configurations differ?
```

without rerunning the workload and without reducing multiple dimensions into a synthetic score.

The primary interface should be conceptually:

```bash
acap bench compare baseline.json candidate.json
```

or the equivalent using the benchmark result storage implemented by previous phases.

The comparison must preserve independent dimensions such as:

```text
task success
agent-visible context cost
command count
recovery behavior
wall time
AgentCap processing overhead
```

Do not automatically declare a winner.

Do not implement CI regression policy yet. That belongs to Phase B9.

---

# 1. B8 Is a Read-Only Analysis Phase

The comparison command must operate on existing benchmark results.

Conceptually:

```text
baseline result
       |
       +--------+
                |
                v
           Comparator
                |
                v
       ComparisonResult
                ^
                |
       +--------+
       |
candidate result
```

Do not:

```text
load baseline
    ->
rerun benchmark

load candidate
    ->
rerun benchmark
```

`acap bench compare` must not execute coding agents, benchmark commands, workload setup, verification, or AgentCap processing.

It compares already recorded observations.

---

# 2. Reuse the Existing Result Model

Do not create an independent parser for benchmark output.

Consume the structured result format introduced in B4 and extended through B5–B7.

Never compare human-readable CLI tables by parsing their text.

The intended dependency is:

```text
serialized BenchmarkResult
        |
        v
result loader
        |
        v
typed benchmark result
        |
        v
Comparator
```

not:

```text
human table
    |
string parsing
    |
comparison
```

Reuse existing serialization/deserialization and schema-version handling wherever possible.

---

# 3. Add `acap bench compare`

Implement a command conceptually similar to:

```bash
acap bench compare baseline.json candidate.json
```

If the existing benchmark persistence system uses run IDs rather than file paths, support the natural repository convention instead.

For example, the final interface may be:

```bash
acap bench compare <baseline-run> <candidate-run>
```

or both forms if existing abstractions make that simple.

Do not introduce multiple input mechanisms merely for convenience.

Choose the smallest interface that fits the B0–B7 architecture.

---

# 4. Use Explicit Baseline and Candidate Semantics

The two sides of the comparison must have stable meanings:

```text
baseline
candidate
```

Do not infer the baseline from timestamps or mode names.

For:

```bash
acap bench compare A B
```

interpret:

```text
A = baseline
B = candidate
```

unless existing CLI conventions strongly dictate otherwise.

This distinction matters for signed deltas.

For example:

```text
delta = candidate - baseline
```

must have one consistent definition throughout the implementation.

---

# 5. Introduce a Canonical Comparison Model

Do not let the human formatter independently calculate differences.

Create a structured comparison result.

Conceptually:

```go
type BenchmarkComparison struct {
    Baseline  BenchmarkResult
    Candidate BenchmarkResult

    Metrics ComparisonMetrics
}
```

or an equivalent architecture.

The exact type structure should match the existing codebase.

The important pipeline is:

```text
baseline result
candidate result
       |
       v
   Comparator
       |
       v
BenchmarkComparison
       |
       +--> human formatter
       |
       +--> JSON formatter
```

All comparison semantics must live in the comparator layer.

---

# 6. Compare Like With Like

Before computing meaningful deltas, validate that the results are comparable.

For coding-agent repeated benchmarks, compatibility should normally include:

```text
same workload
same agent
compatible benchmark/result semantics
```

Differences such as:

```text
mode = disabled
vs
mode = stateful
```

are expected and are often exactly what B8 is intended to compare.

However:

```text
workload A
vs
workload B
```

must not silently produce a normal comparison.

Likewise:

```text
Codex
vs
Claude Code
```

should not automatically be treated as an equivalent configuration comparison unless the user explicitly requests or the existing command semantics support cross-agent comparison.

Start conservatively.

Prefer rejecting or clearly marking incompatible comparisons rather than producing misleading deltas.

---

# 7. Do Not Require Identical Trial Counts

Repeated benchmark results may contain different numbers of trials.

For example:

```text
baseline:
    5 trials

candidate:
    10 trials
```

This does not automatically make the results incomparable.

Display the trial counts explicitly.

Use each result's already computed B7 aggregates.

Do not discard candidate trials merely to force equal sample sizes.

Do not resample or randomly select trials.

B8 is comparison, not statistical normalization.

---

# 8. Preserve B7 Aggregation Semantics

Do not recalculate benchmark statistics using different rules.

If B7 defines:

```text
median_total_visible_bytes
median_command_count
median_wall_time
median_processing_time
```

then B8 should compare those canonical B7 aggregate values.

Do not independently recompute medians from raw trials unless necessary for schema migration or validation.

Likewise, preserve B7's population rules regarding:

```text
successful trials
task failures
timeouts
missing measurements
```

B8 must not silently change which trials contribute to a metric.

---

# 9. Compare Task Success Directly

Display success counts independently.

For example:

```text
                       baseline     candidate
Task success               5/5           4/5
```

Do not convert this into a winner label.

Do not say:

```text
candidate loses
baseline is better
```

The user should be able to see the success behavior directly.

If trial counts differ, retain both numerator and denominator:

```text
baseline:   5/5
candidate:  8/10
```

A percentage may be shown as supplemental information, but counts remain important.

---

# 10. Compare Agent-Visible Context Cost

The primary context-cost metric for B5+ results is:

```text
total_visible_bytes
```

For repeated B7 results, compare:

```text
median_total_visible_bytes
```

Conceptually:

```text
                       baseline     candidate      delta
Visible bytes             412k           91k       -321k
```

Optionally display a relative delta:

```text
-77.9%
```

if the baseline is non-zero and the calculation is well-defined.

Do not use initial capsule size as the primary comparison when total recovery-aware cost is available.

The workflow-level KPI must include progressive-disclosure recovery.

---

# 11. Compare Command Count

Display:

```text
median_command_count
```

for repeated coding-agent results.

For example:

```text
                       baseline     candidate
Median commands              31            32
```

This prevents a context reduction from being interpreted without seeing changes in agent behavior.

Do not hide increased command count merely because visible bytes decreased.

---

# 12. Compare Wall Time

Display:

```text
median_wall_time
```

separately.

For example:

```text
                       baseline     candidate
Median wall time          84.2s         86.1s
```

Do not infer AgentCap processing overhead by subtracting these values.

Wall time includes many sources of variability:

```text
LLM latency
provider latency
command execution
agent decisions
filesystem operations
AgentCap
```

Keep wall time observational.

---

# 13. Compare AgentCap Processing Time Separately

Where applicable, compare:

```text
median_processing_time
```

or the equivalent canonical B7 metric.

For disabled mode, AgentCap processing may be:

```text
0
not applicable
```

depending on the established result schema.

Preserve that semantic distinction.

Do not fabricate zero if the value is actually unavailable.

---

# 14. Compare Recovery Behavior

Progressive-disclosure behavior must remain visible.

At minimum compare metrics such as:

```text
trials_with_show
trials_with_raw_retrieval
total_show_count
total_raw_retrieval_count
```

or the canonical equivalents established by B7.

A useful human representation may include:

```text
                       baseline     candidate
Show used                  -           3/5
Raw fallback               -           1/5
```

or:

```text
Raw retrievals/run         -           0.2
```

if that derived metric is clearly defined.

Do not hide raw fallback behind total visible bytes.

A candidate with low median bytes but frequent raw fallback may behave very differently from one that rarely needs recovery.

---

# 15. Keep Recovery Bytes Observable

Where useful, also expose:

```text
median_show_bytes
median_raw_retrieval_bytes
```

or other B7 recovery-byte aggregates if they already exist.

Do not introduce a large number of redundant derived metrics solely for B8.

The primary goal is to preserve visibility into:

```text
how much context was consumed
and
how that context was recovered
```

---

# 16. Define Numeric Delta Semantics Centrally

For comparable numeric metrics, use one consistent definition:

```text
absolute_delta = candidate - baseline
```

For example:

```text
baseline  = 412000
candidate =  91000

delta = -321000
```

For relative change:

```text
relative_delta =
    (candidate - baseline) / baseline
```

when:

```text
baseline != 0
```

Centralize this calculation.

Do not let individual formatters invent their own sign conventions.

---

# 17. Handle Zero Baselines Correctly

Relative percentage change is undefined when the baseline is zero.

For example:

```text
baseline raw retrieval count = 0
candidate raw retrieval count = 1
```

Do not emit:

```text
+Inf%
```

or an arbitrary percentage.

Represent relative change as:

```text
n/a
undefined
```

or the structured equivalent.

Absolute delta remains valid:

```text
+1
```

Use existing project conventions for optional values.

---

# 18. Do Not Calculate Meaningless Deltas

Some fields should be displayed side by side rather than numerically subtracted.

Examples include:

```text
agent name
mode
workload
outcome labels
schema version
```

Likewise, task success counts with different denominators should not be reduced to a simplistic integer delta without context.

Prefer:

```text
5/5 -> 8/10
```

over:

```text
+3 successes
```

when trial counts differ.

The comparator should distinguish:

```text
numeric comparable metrics
categorical metadata
structured outcome metrics
```

---

# 19. No Composite Score

Do not introduce:

```text
benchmark score
efficiency score
AgentCap score
overall score
```

Do not combine:

```text
task success
visible bytes
latency
command count
recovery
```

using arbitrary weights.

These dimensions answer different questions.

B8 should expose the tradeoffs rather than hide them behind one number.

---

# 20. No Automatic Winner

Do not print:

```text
Winner: candidate
```

or:

```text
PASS
FAIL
```

based on benchmark comparison.

Likewise, do not implement logic such as:

```text
if visible_bytes improved:
    candidate wins
```

because task success or recovery behavior may have changed.

The comparison layer reports observations.

B9 may later introduce explicit regression policy for CI, but that policy must remain separate from B8's descriptive comparison.

---

# 21. Human-Readable Comparison Output

Produce a compact side-by-side table.

For example:

```text
Benchmark comparison
Workload: bugfix-001
Agent:    codex

                         baseline     candidate
Mode                     disabled      stateful
Trials                        5             5

Task success                5/5           5/5

Median visible bytes       412k           91k
Median commands              31            32
Median wall time           84.2s         86.1s

Show used                     -           3/5
Raw fallback                  -           1/5
AgentCap processing           -          12.8ms
```

Where useful, include a delta column:

```text
                         baseline   candidate      delta
Visible bytes               412k         91k      -321k
Commands                      31          32          +1
Wall time                  84.2s       86.1s       +1.9s
```

Keep the output compact.

Do not overload the default view with every field present in the result schema.

Prioritize the metrics needed to understand workflow tradeoffs.

---

# 22. Make Direction Explicit Without Assigning Value

A negative byte delta means:

```text
candidate used fewer bytes
```

A positive wall-time delta means:

```text
candidate took longer
```

Those are descriptive facts.

Do not automatically attach semantic labels such as:

```text
GOOD
BAD
REGRESSION
IMPROVEMENT
```

because whether a tradeoff is acceptable depends on other dimensions and later policy.

B8 should report signed values clearly enough that the user does not need those labels.

---

# 23. JSON Comparison Output

Support:

```bash
acap bench compare baseline.json candidate.json --json
```

or the equivalent argument order established by the CLI.

The JSON result should contain structured comparison information.

Conceptually:

```json
{
  "schema_version": 1,

  "baseline": {
    "workload": "bugfix-001",
    "agent": "codex",
    "mode": "disabled",
    "trial_count": 5
  },

  "candidate": {
    "workload": "bugfix-001",
    "agent": "codex",
    "mode": "stateful",
    "trial_count": 5
  },

  "metrics": {
    "median_total_visible_bytes": {
      "baseline": 412000,
      "candidate": 91000,
      "delta": -321000,
      "relative_delta": -0.7791
    },

    "median_command_count": {
      "baseline": 31,
      "candidate": 32,
      "delta": 1
    },

    "median_wall_time_ns": {
      "baseline": 84200000000,
      "candidate": 86100000000,
      "delta": 1900000000
    }
  }
}
```

This is conceptual only.

Use naming consistent with the existing result architecture.

---

# 24. Version the Comparison Schema Separately Where Appropriate

Do not automatically assume:

```text
benchmark result schema version
=
comparison result schema version
```

They represent different serialized contracts.

If B8 introduces a new persisted/exported comparison JSON format, give that format its own clearly defined schema-version constant.

Conceptually:

```go
const BenchmarkComparisonSchemaVersion = 1
```

This avoids coupling future changes to:

```text
individual benchmark result schema
```

and:

```text
comparison output schema
```

unnecessarily.

Follow existing repository versioning conventions if they already solve this problem.

---

# 25. Preserve Source Results

A comparison result must not destroy or rewrite its source benchmark results.

The comparator is read-only.

Conceptually:

```text
baseline result ----+
                    |
                    v
                comparison
                    ^
                    |
candidate result ---+
```

The source results remain independently inspectable.

Do not normalize them by modifying persisted data in place.

---

# 26. Handle Result Schema Compatibility

Benchmark results may have been produced using different result-schema versions.

Use the B4+ schema-version infrastructure.

Possible approaches include:

```text
load old schema
    ->
upgrade into current in-memory representation
```

if migration support already exists.

Otherwise, reject unsupported combinations clearly.

Do not silently interpret fields whose semantics changed between schema versions.

The comparator must never compare values merely because their JSON field names happen to match.

---

# 27. Distinguish Missing From Zero

Preserve the rule introduced in earlier phases:

```text
missing measurement != measured zero
```

For example:

```text
baseline:
    raw_retrieval_count unavailable

candidate:
    raw_retrieval_count = 0
```

must not be displayed as:

```text
0 vs 0
```

Likewise, do not calculate deltas involving unavailable values.

Use:

```text
n/a
```

or structured null/optional representation.

---

# 28. Support Single-Trial Results Where Reasonable

B8 should primarily support B7 repeated results, but comparison architecture should not unnecessarily exclude B6 single-trial results if they share compatible semantics.

For single-trial comparisons, compare the direct trial metrics:

```text
task success
total_visible_bytes
command_count
wall_time
processing time
recovery
```

Do not invent:

```text
median
```

for a result that is not represented as an aggregate unless the existing architecture naturally treats a single trial as a one-element aggregate.

Prefer reuse over special-case duplication.

---

# 29. Do Not Mix Single-Trial and Repeated Results Silently

If comparing:

```text
one single B6 trial
```

against:

```text
a B7 five-trial aggregate
```

the comparison must make the mismatch obvious or reject it.

Do not present:

```text
91 KB vs 88 KB
```

as if both numbers had identical statistical meaning when one is a single observation and the other is a median.

A conservative initial implementation may simply require both sides to have the same result kind.

That is acceptable for B8.

---

# 30. Optional Directory/Suite Comparison

The roadmap allows future comparison at benchmark-directory granularity.

Do not make suite comparison mandatory if it substantially expands B8.

Implement result-to-result comparison first.

If the existing B7 persistence model naturally groups multiple workloads, a minimal suite comparison may be added, but it must preserve per-workload results.

Never aggregate unrelated workloads into one global byte total and treat that as the sole result.

For example:

```text
bugfix-001
bugfix-002
compile-001
```

should remain individually observable.

---

# 31. If Suite Comparison Is Implemented, Match by Stable Workload Identity

Do not compare suite entries by array position.

Match:

```text
baseline workload ID
```

to:

```text
candidate workload ID
```

using the stable workload identity established in previous phases.

Explicitly report:

```text
missing from baseline
missing from candidate
```

rather than silently dropping unmatched workloads.

Do not treat missing workloads as zero-valued results.

---

# 32. Comparison Should Be Deterministic

Given the same two benchmark results, B8 must always produce the same comparison.

Do not introduce:

```text
sampling
randomized statistics
LLM-generated interpretation
heuristic winner selection
```

The comparison layer should be deterministic and suitable for scripting.

---

# 33. Keep Comparison Independent From Agent Execution

The comparator package should not depend on:

```text
Claude adapter
Codex adapter
agent process execution
fixture workspace management
task verification
```

It should depend on benchmark result types.

Conceptually:

```text
agent execution
      |
      v
BenchmarkResult
      |
      v
Comparator
```

Never:

```text
Comparator
      |
      v
AgentAdapter
```

This separation will make B8 useful for both coding-agent results and deterministic B0–B5 benchmark results where applicable.

---

# 34. Keep Comparison Independent From CI Policy

Similarly:

```text
Comparator
```

must not know about:

```text
allowed regression percentages
CI pass/fail thresholds
required success rates
release gating
```

Those belong to B9.

B8 produces facts.

B9 applies policy.

Maintain that dependency direction.

---

# 35. Result Loading Errors

Handle malformed or unsupported inputs explicitly.

Examples:

```text
file does not exist
invalid JSON
unsupported schema version
missing required fields
incompatible result type
incompatible workload
```

These are comparison-command errors.

Do not produce a partially valid-looking comparison using fabricated defaults.

---

# 36. Tests

Add focused tests for the comparison model and CLI.

## Basic numeric comparison

Given:

```text
baseline visible bytes  = 100
candidate visible bytes = 60
```

verify:

```text
absolute delta = -40
relative delta = -0.4
```

---

## Positive delta

Given:

```text
baseline commands  = 20
candidate commands = 25
```

verify:

```text
delta = +5
```

---

## Zero baseline

Verify relative delta is unavailable rather than infinity or a fabricated percentage.

---

## Success comparison

Verify:

```text
5/5
vs
4/5
```

is preserved as structured outcome data.

Do not collapse it into a generic numeric score.

---

## Different trial counts

Verify:

```text
5 trials
vs
10 trials
```

can be represented without discarding observations.

---

## Recovery comparison

Verify differences in:

```text
trials_with_show
trials_with_raw_retrieval
total_show_count
total_raw_retrieval_count
```

are represented correctly.

---

## Missing metrics

Verify unavailable measurements remain unavailable and are not interpreted as zero.

---

## Incompatible workloads

Verify comparing:

```text
bugfix-001
vs
bugfix-002
```

does not silently produce a normal comparison.

---

## Different modes

Verify:

```text
disabled
vs
stateful
```

is accepted when the workload/agent/result semantics are otherwise compatible.

---

## Schema compatibility

Test supported schema-version combinations.

Reject unsupported versions clearly.

---

## Human output

Verify the human table contains the core independent dimensions:

```text
task success
visible bytes
commands
wall time
recovery
processing overhead
```

Avoid brittle exact-whitespace tests unless the repository already uses golden CLI tests.

---

## JSON output

Verify:

```bash
acap bench compare ... --json
```

produces clean valid JSON with no human-oriented text on stdout.

---

## No execution

Use mocks/fakes or architecture tests where appropriate to verify the compare path never invokes:

```text
agent execution
workload execution
fixture mutation
AgentCap command processing
```

This is an important B8 invariant.

---

# 37. Documentation

Document:

```bash
acap bench compare <baseline> <candidate>
```

Explain:

```text
first input = baseline
second input = candidate
delta = candidate - baseline
```

Document the primary comparison dimensions:

```text
task success
agent-visible bytes
command count
wall time
show/raw behavior
AgentCap processing overhead
```

Explain that the command intentionally does not produce an overall winner or composite score.

Document compatibility requirements.

If different trial counts are supported, explicitly state that each side retains its own sample count.

---

# 38. Do Not Implement B9 Yet

Do not add:

```text
regression threshold configuration
PASS / FAIL comparison
CI exit status based on metric changes
allowed percentage regression
GitHub Actions policy
automatic baseline selection
PR benchmark gating
```

For example, do not implement:

```text
fail if visible bytes increase by 10%
```

or:

```text
fail if median wall time increases
```

Those are B9 policy decisions.

B8 should provide the structured comparison data B9 will later consume.

---

# Implementation Guidance

Before modifying code:

1. inspect the B4–B7 serialized result structures;
2. inspect how benchmark results are persisted and loaded;
3. identify the canonical B7 aggregate fields;
4. inspect the result-schema compatibility/versioning mechanism;
5. define comparison compatibility rules;
6. implement a small pure comparator;
7. create a canonical `BenchmarkComparison` result;
8. implement human and JSON formatters from that result;
9. add the CLI command only after comparator behavior is well tested;
10. avoid touching agent execution or benchmark measurement unless strictly necessary.

Prefer an architecture like:

```text
ResultLoader
     |
     +--> baseline
     |
     +--> candidate
              |
              v
          Comparator
              |
              v
     BenchmarkComparison
              |
        +-----+-----+
        |           |
        v           v
      human        JSON
```

The comparator should ideally be close to a pure function:

```go
comparison, err := Compare(baseline, candidate)
```

It should not know where the results came from.

---

# Architectural Invariants

After B8, the benchmark pipeline should conceptually be divided into two paths:

```text
EXECUTION PATH

workload
   |
agent
   |
AgentCap
   |
measurement
   |
BenchmarkResult
   |
persistence
```

and:

```text
ANALYSIS PATH

stored baseline result
          |
          +----------+
                     |
                     v
                 Comparator
                     ^
                     |
          +----------+
          |
stored candidate result
                     |
                     v
             BenchmarkComparison
                     |
               +-----+-----+
               |           |
               v           v
             human        JSON
```

The analysis path must never loop back into execution:

```text
Comparator
    X
    |
    +--> workload runner
    +--> agent adapter
    +--> command execution
```

The central B8 invariant is:

```text
comparison observes benchmark results;
comparison does not create new benchmark observations.
```

---

# Core Comparison Model

For the same:

```text
workload × agent
```

B8 should make comparisons such as:

```text
disabled
vs
stateful
```

or:

```text
stateful
vs
integrated
```

observable across independent dimensions:

```text
                         baseline     candidate

Task success                5/5           5/5

Median visible bytes       412k           91k
Median commands              31            32
Median wall time           84.2s         86.1s

Show usage                    -           3/5
Raw usage                     -           1/5

AgentCap processing           -          12.8ms
```

The command may additionally expose signed numeric deltas:

```text
visible bytes:  -321k
commands:       +1
wall time:      +1.9s
```

but must leave interpretation of those tradeoffs to the consumer.

---

# Exit Criteria

Phase B8 is complete when all of the following are true:

1. `acap bench compare` can compare two previously produced benchmark results.
2. Comparison never reruns workloads, commands, or coding agents.
3. Baseline and candidate semantics are explicit and stable.
4. Numeric delta semantics are consistently defined as `candidate - baseline`.
5. Compatible B7 repeated benchmark results can be compared.
6. Trial counts remain visible and do not need to be equal.
7. Task-success counts are displayed independently.
8. Median total agent-visible bytes are compared.
9. Median command counts are compared.
10. Median wall times are compared.
11. AgentCap processing overhead is independently observable.
12. `show` / `raw` recovery behavior remains observable.
13. Missing measurements are not interpreted as zero.
14. Relative deltas with a zero baseline are handled safely.
15. Incompatible workloads/results are rejected or clearly identified.
16. Supported result-schema versions are handled explicitly.
17. Comparison uses one canonical comparison model.
18. Human and JSON output consume the same comparison model.
19. JSON comparison output is machine-readable and clean on stdout.
20. Individual source benchmark results remain unchanged.
21. No composite benchmark score is introduced.
22. No automatic winner is selected.
23. No B9 regression threshold or CI pass/fail policy is introduced.

The completed B8 should turn two stored benchmark observations such as:

```text
baseline:
    task success = 5/5
    median visible = 412 KB
    median commands = 31
    median wall = 84.2 s

candidate:
    task success = 5/5
    median visible = 91 KB
    median commands = 32
    median wall = 86.1 s
```

into a deterministic, structured comparison that exposes:

```text
what changed
by how much
and on which independent metric
```

without converting those measurements into an automatic judgment about which configuration is "better."
