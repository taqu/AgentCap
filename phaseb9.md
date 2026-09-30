# Phase B9 — CI Regression Benchmark

## Objective

Integrate the benchmark system into the AgentCap development feedback loop so deterministic regressions can be detected automatically in CI.

Phase B0–B8 already provide:

```text
measurement
single-command benchmarks
stateful workflow benchmarks
reproducible workloads
stable result formats
recovery-cost accounting
coding-agent benchmarks
repeated trials
benchmark comparison
```

Phase B9 adds a separate policy layer that can answer:

```text
Did this change cause a benchmark regression
large enough that CI should fail?
```

The primary CI target is the lightweight deterministic benchmark suite.

The initial regression dimensions should focus on measurements such as:

```text
raw bytes
capsule/stateless bytes
stateful bytes
AgentCap processing latency
```

and other deterministic B0–B5 metrics where appropriate.

Do not make nondeterministic LLM-agent benchmarks from B6/B7 a required normal CI gate.

Those benchmarks should remain available as a separate evaluation workflow.

---

# 1. Preserve the Measurement / Comparison / Policy Separation

B9 must not embed regression thresholds into benchmark measurement or B8 comparison logic.

Maintain the architecture:

```text
benchmark execution
        |
        v
BenchmarkResult
        |
        v
Comparator
        |
        v
BenchmarkComparison
        |
        v
RegressionPolicy
        |
        v
RegressionEvaluation
        |
        v
CI exit status
```

Each layer has a separate responsibility.

### Measurement

Answers:

```text
What happened?
```

### Comparison

Answers:

```text
What changed?
```

### Regression policy

Answers:

```text
Does this change violate an explicitly configured threshold?
```

Do not merge these responsibilities.

---

# 2. Reuse B8 Comparison Semantics

Use the comparison model implemented in B8.

In particular, preserve:

```text
baseline
candidate
delta = candidate - baseline
```

Do not create a second implementation of percentage or absolute delta calculations inside the CI layer.

For example:

```text
baseline stateful bytes  = 20,000
candidate stateful bytes = 23,000

delta          = +3,000
relative delta = +15%
```

must use the same comparison semantics in:

```text
acap bench compare
```

and:

```text
CI regression evaluation
```

The difference is that B9 applies an explicit policy to those observations.

---

# 3. Introduce an Explicit Regression Policy Model

Regression thresholds must be explicit configuration, not hard-coded assumptions scattered through the CLI.

Conceptually:

```yaml
metrics:
  stateful_bytes:
    max_relative_increase: 0.10

  stateless_bytes:
    max_relative_increase: 0.10

  processing_ns:
    max_relative_increase: 0.20
```

The exact configuration syntax should fit the existing benchmark/workload architecture.

Possible locations include:

```text
benchmark suite configuration
dedicated regression policy file
CI benchmark configuration
```

Choose one clear model.

Do not duplicate thresholds across workload files, CLI code, and CI scripts.

---

# 4. Start With a Small Policy Surface

Do not build a generic policy language.

The initial policy model only needs to express useful benchmark regression constraints.

At minimum consider support for:

```text
max absolute increase
max relative increase
```

where appropriate.

For example:

```yaml
stateful_bytes:
  max_relative_increase: 0.10
```

means conceptually:

```text
candidate may be at most 10% larger than baseline
```

For latency:

```yaml
processing_ns:
  max_relative_increase: 0.20
```

Do not add:

```text
arbitrary expressions
boolean expression trees
embedded scripting
custom user code
```

B9 needs regression thresholds, not a general rules engine.

---

# 5. Prefer Upper-Bound Regression Rules

For the initial deterministic metrics, regression generally means an undesirable increase.

Examples:

```text
capsule bytes increase
stateful bytes increase
processing latency increase
recovery bytes increase
```

Represent this directly.

Do not over-generalize the initial implementation merely to support every possible metric direction.

If the policy abstraction can cleanly represent direction, that is fine, but keep the user-facing configuration simple.

---

# 6. Do Not Treat Raw Baseline Changes as AgentCap Regressions Automatically

Be careful with:

```text
raw_bytes
```

A change in fixture/compiler/tool output may legitimately change the raw output size.

For example:

```text
baseline raw bytes  = 100 KB
candidate raw bytes = 150 KB
```

does not by itself prove AgentCap regressed.

Raw output is primarily a reference measurement.

AgentCap-controlled metrics such as:

```text
stateless/capsule bytes
stateful bytes
processing latency
recovery behavior
```

are more direct regression targets.

If raw-byte changes are checked, give them explicitly defined semantics rather than automatically treating every increase as a failure.

---

# 7. Support Absolute and Relative Thresholds Carefully

Relative thresholds are useful for large metrics.

For example:

```text
baseline = 100,000 bytes
candidate = 108,000 bytes

relative increase = 8%
```

An allowed threshold of:

```text
10%
```

would not trigger a regression.

However, relative percentages become unstable for very small baselines.

For example:

```text
baseline = 10 bytes
candidate = 20 bytes
```

is:

```text
+100%
```

but only:

```text
+10 bytes
```

Therefore allow absolute thresholds where appropriate.

Do not assume relative percentage is always the correct regression criterion.

---

# 8. Handle Zero Baselines Correctly

Reuse B8's zero-baseline semantics.

For:

```text
baseline = 0
candidate > 0
```

relative change is undefined.

Do not produce:

```text
+Inf%
```

and do not silently treat it as zero.

If a policy requires relative comparison and the baseline is zero, evaluation should either:

```text
use an explicitly configured absolute rule
```

or:

```text
report that the relative policy cannot be evaluated
```

according to the chosen policy semantics.

Do not invent an arbitrary percentage.

---

# 9. Introduce a Canonical Regression Evaluation Result

Do not let CLI code independently decide whether CI should pass.

Create a structured result.

Conceptually:

```go
type RegressionEvaluation struct {
    Passed bool

    Checks []RegressionCheck
}
```

with checks conceptually containing:

```go
type RegressionCheck struct {
    Metric string

    Baseline  ...
    Candidate ...

    Delta         ...
    RelativeDelta ...

    Threshold ...

    Status ...
}
```

Adapt types to the current architecture.

The important pipeline is:

```text
BenchmarkComparison
        +
RegressionPolicy
        |
        v
RegressionEvaluator
        |
        v
RegressionEvaluation
```

Human output, JSON output, and process exit status should all consume this same evaluation result.

---

# 10. Keep Metric Status Structured

Each policy check should have a structured outcome.

Conceptually:

```text
pass
regression
not_evaluable
```

or equivalent naming consistent with the codebase.

Do not encode policy results only as arbitrary error strings.

This makes the result usable by:

```text
CLI
CI
JSON consumers
future reporting tools
```

---

# 11. Define Overall CI Status Explicitly

The overall regression evaluation should fail when at least one required metric violates its configured threshold.

Conceptually:

```text
all required checks pass
    ->
overall pass

one or more required checks regress
    ->
overall fail
```

Define behavior for:

```text
not_evaluable
```

explicitly.

For required checks, conservative behavior is generally preferable:

```text
required metric cannot be evaluated
    ->
CI evaluation does not silently pass
```

However, follow existing project conventions if there is already a policy for missing benchmark data.

Do not silently ignore unavailable required metrics.

---

# 12. Separate Regression Failure From Infrastructure Failure

These are different outcomes:

```text
benchmark completed
candidate exceeded configured threshold
```

versus:

```text
benchmark could not execute
```

and:

```text
baseline result could not be loaded
```

Preserve this distinction.

Conceptually:

```text
regression
    -> valid benchmark evaluation
    -> CI policy failed

infrastructure error
    -> evaluation could not be completed
```

Do not report infrastructure failure as:

```text
0 regressions
```

or as a successful policy evaluation.

---

# 13. Provide CI-Friendly Exit Codes

The command used by CI must return a non-zero exit status when configured regression policy fails.

Conceptually:

```text
0
    benchmark evaluation completed
    no configured regression detected

non-zero
    regression detected
    OR
    benchmark/evaluation infrastructure error
```

If the project already distinguishes exit codes by error class, reuse that convention.

Do not introduce complicated exit-code taxonomy solely for B9 unless it already fits the CLI architecture.

Machine-readable JSON should provide enough structured information to distinguish regression from infrastructure error even if both produce non-zero process status.

---

# 14. Add a Regression Check Command

Add a command conceptually similar to:

```bash
acap bench check \
    --baseline baseline.json \
    --candidate candidate.json \
    --policy benchmarks/regression.yaml
```

or an equivalent interface consistent with B8.

Another reasonable shape is:

```bash
acap bench compare baseline.json candidate.json \
    --policy benchmarks/regression.yaml \
    --check
```

Prefer a separate `check` concept if it keeps the architecture clear:

```text
compare
    -> descriptive

check
    -> policy evaluation
```

Do not force B8's normal descriptive comparison command to start returning failure merely because a metric changed.

Preserve the distinction between:

```text
acap bench compare
```

and:

```text
acap bench check
```

unless existing CLI architecture strongly favors another clean design.

---

# 15. Keep `acap bench compare` Descriptive

After B9, normal B8 comparison should still work as:

```bash
acap bench compare baseline.json candidate.json
```

and simply display differences.

It must not unexpectedly return failure because:

```text
visible bytes increased
latency increased
```

Regression policy should only apply when explicitly requested through the B9 check path.

This preserves B8 as a general analysis tool.

---

# 16. Add JSON Output for Regression Evaluation

Support machine-readable policy results.

Conceptually:

```bash
acap bench check \
    --baseline baseline.json \
    --candidate candidate.json \
    --policy regression.yaml \
    --json
```

may produce:

```json
{
  "schema_version": 1,
  "passed": false,

  "checks": [
    {
      "metric": "stateful_bytes",
      "baseline": 20000,
      "candidate": 23000,
      "delta": 3000,
      "relative_delta": 0.15,
      "max_relative_increase": 0.10,
      "status": "regression"
    },
    {
      "metric": "processing_ns",
      "baseline": 5000000,
      "candidate": 5400000,
      "delta": 400000,
      "relative_delta": 0.08,
      "max_relative_increase": 0.20,
      "status": "pass"
    }
  ]
}
```

This is conceptual only.

Use repository naming and serialization conventions.

---

# 17. Version the Regression-Evaluation Schema

If regression evaluation is serialized as its own JSON contract, version it independently.

Conceptually:

```go
const RegressionEvaluationSchemaVersion = 1
```

Do not automatically couple it to:

```text
BenchmarkResult schema
BenchmarkComparison schema
```

They represent different contracts.

Follow existing schema-versioning infrastructure where possible.

---

# 18. Human-Readable CI Output

Produce compact output suitable for local development and CI logs.

For example:

```text
Benchmark regression check

Metric               Baseline   Candidate    Change    Limit    Status
stateful bytes         20,000      23,000    +15.0%   +10.0%    REGRESSION
stateless bytes        42,000      43,500     +3.6%   +10.0%    PASS
processing              5.0ms       5.4ms     +8.0%   +20.0%    PASS

Regression checks: 1 failed
```

Exact formatting should follow existing CLI conventions.

Keep enough information in the output to answer:

```text
which metric failed?
what was the baseline?
what was the candidate?
how much did it change?
what threshold was configured?
```

Do not require developers to rerun the comparison just to understand the CI failure.

---

# 19. Use Deterministic Benchmarks for Normal CI

The default CI suite should focus on deterministic AgentCap-only workloads from B0–B5.

Good CI candidates include:

```text
single-command reducer fixtures
repeated-command stateful fixtures
deterministic state-transition workloads
deterministic recovery-cost workloads
```

These can reliably detect regressions in:

```text
output size
stateful compression
recovery accounting
processing latency
```

without introducing LLM nondeterminism.

---

# 20. Keep B6/B7 LLM Benchmarks Out of Required Normal CI

Do not make normal pull-request CI depend on:

```text
Claude Code
Codex
external LLM API availability
LLM nondeterminism
paid model usage
provider rate limits
```

The B6/B7 benchmark suite should remain separate.

Possible uses include:

```text
manual evaluation
scheduled evaluation
release evaluation
nightly experimentation
```

but B9 does not need to implement all of those orchestration mechanisms.

The important architectural boundary is:

```text
deterministic AgentCap benchmark
    -> suitable for CI regression gating

real LLM-agent benchmark
    -> separate evaluation
```

---

# 21. Do Not Hide Agent Benchmarks

Although full agent benchmarks should not gate normal CI, preserve their usefulness.

Documentation may recommend a separate command or workflow such as:

```text
full agent evaluation
```

for significant AgentCap changes.

Do not delete, weaken, or special-case B6/B7 merely because B9 focuses on deterministic CI.

B9 adds a fast feedback loop; it does not replace full agent evaluation.

---

# 22. Define a Small Deterministic Regression Suite

Create or identify a lightweight benchmark suite suitable for routine CI.

It should exercise representative AgentCap behavior.

For example:

```text
git diff reduction
repeated git diff
repeated test output
error-count transition
unchanged stateful result
delta stateful result
show recovery
raw recovery
```

Use the B3 workload/fixture infrastructure.

Do not create special CI-only execution semantics if normal benchmark workloads can express these cases.

---

# 23. Keep CI Runtime Reasonable

The regression suite should be small enough to run routinely.

Do not simply execute every benchmark fixture created during development.

Select workloads that provide coverage of major behavior classes.

Prefer:

```text
small deterministic fixtures
fast commands
stable output
representative reducers
stateful transitions
```

Avoid:

```text
large repositories
network-dependent commands
long builds
environment-sensitive tools
```

unless specifically isolated from the default CI suite.

---

# 24. Baselines Must Be Explicit

Do not silently choose an arbitrary previous benchmark result as the baseline.

The CI system should have an explicit baseline source.

Possible approaches include:

```text
checked-in baseline result
baseline generated from the target/base branch
known benchmark artifact
```

Choose the approach that best fits the repository's existing CI architecture.

The key requirement is that developers can determine:

```text
what result was used as the baseline?
```

Do not hide baseline selection behind implicit timestamp ordering.

---

# 25. Prefer Base-Branch Baselines Where Infrastructure Allows

If CI infrastructure can reliably execute the deterministic benchmark on both:

```text
base revision
candidate revision
```

that can reduce stale-baseline maintenance.

Conceptually:

```text
base commit
    ->
deterministic benchmark
    ->
baseline

candidate commit
    ->
same deterministic benchmark
    ->
candidate

baseline + candidate
    ->
B8 comparison
    ->
B9 policy
```

However, do not implement complex Git/CI orchestration inside the benchmark library itself.

Keep repository/CI workflow mechanics outside the core evaluator where possible.

The B9 core should consume explicit results regardless of how CI obtained them.

---

# 26. Avoid Benchmark Self-Modification

CI benchmark execution must use disposable fixtures/workspaces exactly as established in B3+.

Do not mutate the source repository to produce benchmark state transitions.

The same workload should be runnable:

```text
locally
in CI
on baseline revision
on candidate revision
```

without contaminating subsequent runs.

---

# 27. Stabilize the Environment Where Practical

Processing-latency thresholds can be noisy in CI.

Do not assume nanosecond-level benchmark stability.

Use tolerances appropriate for the metric and environment.

For latency regression checks, consider:

```text
larger relative tolerance
multiple deterministic repetitions
median deterministic processing time
```

if the existing benchmark infrastructure supports this cleanly.

Do not introduce heavy statistical infrastructure merely for CI.

The goal is detecting substantial regressions, not microbenchmark precision.

---

# 28. Be Conservative With Wall-Time CI Gates

AgentCap processing time is generally a better deterministic regression metric than complete command wall time.

Complete wall time can include:

```text
OS scheduling
filesystem cache
compiler variability
process startup
CI host contention
```

Prefer directly instrumented:

```text
AgentCap processing latency
```

for initial CI policy.

If wall time is exposed, consider reporting it without gating initially.

Do not create fragile CI failures from noisy timing data.

---

# 29. Detect Large Regressions, Not Tiny Noise

The roadmap goal is to detect meaningful regressions.

Do not configure thresholds so tightly that normal environmental variation constantly fails CI.

For example, a policy may allow modest variation while catching:

```text
stateful output unexpectedly doubles
reducer stops reducing a fixture
processing overhead increases dramatically
recovery output unexpectedly explodes
```

Exact thresholds should be configurable rather than hard-coded into implementation.

Do not choose permanent product thresholds in code during B9 implementation.

---

# 30. Allow Per-Metric Thresholds

Different metrics need different tolerances.

For example:

```yaml
metrics:
  stateless_bytes:
    max_relative_increase: 0.05

  stateful_bytes:
    max_relative_increase: 0.05

  processing_ns:
    max_relative_increase: 0.25
```

Do not use one global:

```text
max_regression = 10%
```

for every metric.

Output size and processing latency have different noise characteristics.

---

# 31. Consider Per-Workload Policy Only If Needed

Initially prefer a simple suite-wide metric policy.

Do not immediately build:

```text
workload A has threshold X
workload B has threshold Y
workload C has threshold Z
```

unless actual benchmark data demonstrates that per-workload thresholds are necessary.

If the existing suite architecture naturally supports overrides, they may be included carefully.

But avoid turning B9 configuration into a large hierarchy before there is evidence it is needed.

---

# 32. Preserve Per-Workload Visibility

If the CI suite contains multiple workloads, do not hide a severe regression in one workload behind aggregate improvement elsewhere.

For example:

```text
workload A: -50 KB
workload B: +80 KB
```

must not automatically become:

```text
suite total: +30 KB
```

as the only regression signal.

Evaluate or at least expose workload-level results.

A reducer can regress badly for one output class while improving another.

Per-workload observability is important.

---

# 33. Do Not Average Away Regressions

Similarly, avoid policies based only on suite-wide average reduction.

For example:

```text
9 workloads unchanged
1 workload regresses by 200%
```

should remain visible.

The CI policy should be able to detect that individual regression.

Do not let unrelated benchmark cases compensate for one another unless a future explicit policy intentionally allows it.

---

# 34. Keep Raw Benchmark Results as CI Artifacts Where Practical

When CI produces baseline/candidate benchmark results, preserve the structured results where the existing CI environment makes that straightforward.

This makes failures diagnosable and allows later analysis.

Useful artifacts include:

```text
baseline result JSON
candidate result JSON
comparison JSON
regression evaluation JSON
```

Do not make artifact upload a hard dependency of the core benchmark evaluator.

The evaluator should work locally without CI-specific infrastructure.

---

# 35. Make Local Reproduction Easy

A developer who sees a CI regression should be able to reproduce the same policy check locally.

For example:

```bash
acap bench check \
    --baseline baseline.json \
    --candidate candidate.json \
    --policy benchmarks/regression.yaml
```

Avoid putting essential regression logic only inside:

```text
GitHub Actions YAML
shell scripts
CI vendor configuration
```

CI should invoke AgentCap's benchmark/regression functionality, not reimplement it.

---

# 36. Keep CI Integration Thin

The CI configuration should conceptually do:

```text
produce/load baseline
        |
produce candidate
        |
acap bench check
        |
use exit status
```

Do not implement comparison math or threshold checks in shell scripts.

For example, avoid:

```bash
if [ "$candidate" -gt ... ]; then
    exit 1
fi
```

for benchmark semantics.

That logic belongs in the tested B9 regression evaluator.

---

# 37. Policy Validation

Validate regression policy files before evaluation.

Reject cases such as:

```text
unknown metric
negative max increase where unsupported
invalid percentage
duplicate conflicting rules
malformed configuration
```

Do not silently ignore unknown metrics.

A typo such as:

```yaml
statefull_bytes:
```

must not cause the intended regression check to disappear.

Fail clearly.

---

# 38. Missing Metric Policy

If a configured metric is missing from either result, do not silently pass.

For example:

```text
policy requires:
    stateful_bytes

candidate result:
    stateful_bytes unavailable
```

must produce an explicit evaluation outcome.

For required CI checks, this should normally prevent a successful evaluation.

Do not convert missing measurements to zero.

---

# 39. Unsupported Schema Handling

Reuse B8's result loading/schema compatibility infrastructure.

If a baseline or candidate result uses an unsupported schema:

```text
reject clearly
```

or:

```text
upgrade through an existing supported migration path
```

Do not compare semantically incompatible metrics merely because their field names match.

---

# 40. Regression Evaluation Must Be Deterministic

Given:

```text
same baseline
same candidate
same policy
```

the evaluator must always produce the same result.

Do not use:

```text
LLM interpretation
random sampling
heuristic severity classification
```

in the core CI decision path.

B9 must be reliable enough to control process exit status.

---

# 41. Do Not Introduce Automatic Baseline Learning

Do not automatically update regression baselines because the candidate passed.

Do not implement:

```text
candidate passed
    ->
replace baseline
```

inside the benchmark command.

Baseline lifecycle is a repository/CI policy decision.

Keep baseline evaluation separate from baseline mutation.

This prevents accidental normalization of regressions.

---

# 42. Do Not Automatically Rewrite Thresholds

Similarly, do not adapt thresholds based on observed benchmark variance.

For example, do not:

```text
CI failed three times
    ->
increase threshold automatically
```

Thresholds must remain explicit configuration controlled by developers.

---

# 43. Tests

Add focused tests for regression policy and CI behavior.

## Relative threshold pass

Given:

```text
baseline  = 100
candidate = 105
limit     = +10%
```

verify:

```text
status = pass
```

---

## Relative threshold regression

Given:

```text
baseline  = 100
candidate = 115
limit     = +10%
```

verify:

```text
status = regression
```

---

## Exact threshold boundary

Define and test whether:

```text
candidate increase == allowed threshold
```

passes.

Prefer:

```text
increase <= threshold
    -> pass
```

unless existing conventions dictate otherwise.

---

## Absolute threshold

Verify absolute-increase rules independently from relative rules.

---

## Zero baseline

Verify relative policy does not produce infinity or fabricated percentages.

---

## Missing metric

Verify a required missing metric does not silently pass.

---

## Unknown policy metric

Verify configuration validation rejects unknown metric names.

---

## Multiple checks

Given:

```text
check A passes
check B regresses
check C passes
```

verify:

```text
overall evaluation = failed
```

---

## Infrastructure failure

Verify invalid benchmark input remains distinguishable from a valid regression result.

---

## Human output

Verify failed metrics clearly display:

```text
baseline
candidate
change
threshold
status
```

---

## JSON output

Verify structured policy evaluation can be consumed without parsing human output.

---

## Exit status

Add CLI tests verifying:

```text
all checks pass
    -> success exit status

regression detected
    -> non-zero exit status

evaluation infrastructure error
    -> non-zero exit status
```

---

## B8 compatibility

Verify normal:

```bash
acap bench compare ...
```

remains descriptive and does not start failing merely because candidate metrics increased.

---

## Deterministic suite

Where practical, add an integration test running a small deterministic workload through:

```text
benchmark
    ->
compare
    ->
policy evaluation
```

without any LLM dependency.

---

# 44. CI Workflow

Add a lightweight CI integration using the repository's existing CI system.

The workflow should conceptually:

```text
prepare baseline
      |
run deterministic benchmark
      |
produce candidate result
      |
compare/evaluate
      |
fail CI if configured regression exists
```

Keep CI-specific logic thin.

Do not duplicate regression semantics in workflow configuration.

If baseline generation from the base revision is too complex for the first implementation, use a clearly documented stable baseline artifact/result.

Prefer a simple reliable implementation over a sophisticated but fragile CI pipeline.

---

# 45. Documentation

Document:

```text
how to run the deterministic regression suite
how to generate benchmark results
how to compare results
how to apply regression policy
how CI interprets the result
how thresholds are configured
how to reproduce CI failures locally
```

Clearly distinguish:

```text
B8 comparison:
    descriptive

B9 regression check:
    policy-driven
```

Also document that normal CI intentionally does not use full LLM-agent benchmarks as a required gate.

Explain that B6/B7 remain important for higher-level evaluation of:

```text
agent task performance
real recovery behavior
complete coding-agent efficiency
```

but have different reproducibility/cost characteristics from deterministic CI benchmarks.

---

# 46. Do Not Expand B9 Into Benchmark Infrastructure Redesign

B9 is the final integration phase.

Do not use it as an opportunity to rewrite:

```text
measurement
workload execution
session handling
AgentCap reducer architecture
agent adapters
result persistence
comparison semantics
```

Fix small issues if required for CI integration, but preserve the boundaries established in B0–B8.

The expected implementation should mostly consist of:

```text
regression policy model
regression evaluator
CLI check command
deterministic suite selection
CI integration
tests
documentation
```

---

# Implementation Guidance

Before modifying code:

1. inspect the B8 comparison model;
2. inspect deterministic B0–B5 benchmark results;
3. identify metrics stable enough for CI gating;
4. inspect the benchmark result/schema loader;
5. inspect the repository's existing configuration conventions;
6. inspect the existing CI system;
7. define a minimal regression policy model;
8. implement a pure regression evaluator;
9. add human/JSON presentation;
10. add CLI exit-status behavior;
11. validate locally using deterministic fixtures;
12. only then add the thin CI workflow.

Prefer:

```text
BenchmarkComparison
        +
RegressionPolicy
        |
        v
pure RegressionEvaluator
        |
        v
RegressionEvaluation
        |
        +--> human
        +--> JSON
        +--> exit status
```

The evaluator should ideally be usable independently of the CLI:

```go
evaluation, err := EvaluateRegression(comparison, policy)
```

The core evaluator must not know about:

```text
GitHub Actions
shell scripts
agent adapters
LLM providers
fixture creation
```

---

# Architectural Invariants

After B9, the full benchmark architecture should conceptually be:

```text
                 BENCHMARK EXECUTION

workload
   |
   v
AgentCap / agent execution
   |
   v
measurement
   |
   v
BenchmarkResult
   |
   v
persistence


                 BENCHMARK ANALYSIS

baseline result          candidate result
       |                       |
       +-----------+-----------+
                   |
                   v
                B8 Comparator
                   |
                   v
          BenchmarkComparison


                 CI POLICY

          BenchmarkComparison
                   +
           RegressionPolicy
                   |
                   v
           B9 Evaluator
                   |
                   v
       RegressionEvaluation
                   |
          +--------+--------+
          |        |        |
          v        v        v
        human     JSON    exit status
```

The dependency direction must remain:

```text
measurement
    ->
comparison
    ->
policy
```

Never:

```text
measurement
    ->
CI threshold logic
```

and never:

```text
Comparator
    ->
CI-specific behavior
```

---

# Recommended CI Boundary

The intended final benchmark split is:

```text
FAST / DETERMINISTIC

B0–B5 workloads
    |
    v
routine CI regression checks
```

versus:

```text
EXPENSIVE / NONDETERMINISTIC

B6 real coding-agent benchmark
    |
B7 repeated trials
    |
B8 comparison
    |
manual / scheduled / release evaluation
```

B8 itself remains usable for both categories.

B9 policy-driven gating should initially focus on the deterministic side.

---

# Exit Criteria

Phase B9 is complete when all of the following are true:

1. A deterministic benchmark result can be evaluated against an explicit regression policy.
2. Regression policy is separate from benchmark measurement.
3. Regression policy is separate from B8 comparison semantics.
4. B8 remains a descriptive comparison command.
5. A CI-oriented regression-check command is available.
6. Baseline and candidate inputs are explicit.
7. Regression checks reuse B8 delta semantics.
8. Per-metric thresholds are configurable.
9. Relative-increase thresholds are supported where appropriate.
10. Absolute-increase thresholds are supported where appropriate.
11. Zero baselines are handled safely.
12. Missing required metrics do not silently pass.
13. Unknown policy metrics are rejected.
14. Regression results have a canonical structured representation.
15. Human-readable regression output identifies the metric, baseline, candidate, change, threshold, and status.
16. Machine-readable JSON regression output is available.
17. Regression-evaluation JSON is versioned if it forms a separate serialized contract.
18. CI receives a non-zero exit status when a configured regression is detected.
19. Infrastructure failure remains distinguishable from a valid regression result.
20. A small deterministic regression suite exists or is explicitly selected.
21. The default CI suite uses reproducible B0–B5 workloads.
22. Full B6/B7 LLM-agent benchmarks are not required for normal CI gating.
23. CI-specific orchestration remains thin and does not duplicate benchmark semantics.
24. Per-workload regressions remain observable rather than being hidden by suite-wide averages.
25. Developers can reproduce the regression check locally using the same policy.
26. Baselines are not automatically rewritten by the benchmark tool.
27. Thresholds are not automatically adjusted.
28. Existing B0–B8 functionality remains intact.

The completed Phase B9 should establish the development loop:

```text
change AgentCap
      |
      v
run deterministic benchmark
      |
      v
compare with explicit baseline
      |
      v
apply explicit regression policy
      |
      +--> no configured regression
      |        -> CI continues
      |
      +--> configured regression
               -> CI fails with
                  measurable evidence
```

while keeping the more expensive real-agent evaluation path separate:

```text
significant AgentCap change
      |
      v
B6/B7 real coding-agent evaluation
      |
      v
B8 comparison
      |
      v
human analysis of
task success / context cost /
recovery / latency
```

This preserves the distinction between fast deterministic development feedback and the higher-level question of whether AgentCap actually improves complete coding-agent workflows.
