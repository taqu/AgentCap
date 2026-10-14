# Phase B5 — Progressive-Disclosure Recovery-Cost Benchmark

## Objective

Extend the benchmark system so it measures the **complete agent-visible cost of progressive disclosure**, not only the initial compressed command output.

AgentCap may substantially reduce the first response while requiring the agent to recover omitted information later through commands such as:

```bash
acap show ...
acap raw ...
```

Those recovery operations are part of the actual workflow cost and must therefore be included in benchmark measurements.

The primary concept introduced in Phase B5 is:

```text
initial AgentCap-visible output
+
show retrieval output
+
raw retrieval output
=
total agent-visible output
```

The benchmark must make it possible to distinguish a genuinely efficient compressed workflow from one that initially looks small but repeatedly falls back to information recovery.

Do not implement real coding-agent integration yet. That belongs to Phase B6.

---

# 1. Preserve the Existing Measurement Pipeline

Build B5 on top of the benchmark infrastructure implemented in B0–B4.

Do not create a parallel benchmark execution path.

The intended architecture remains conceptually:

```text
workload
   |
   v
benchmark runner
   |
   v
AgentCap execution / retrieval
   |
   v
measurements
   |
   v
BenchmarkResult
   |
   +--> human output
   |
   +--> JSON output
```

Extend the existing measurement/result model where necessary.

Do not redesign the reducer, session model, or workload runner merely to support B5.

---

# 2. Define the Recovery-Cost Model

The benchmark must distinguish at least these values:

```text
initial_visible_bytes
show_bytes
raw_retrieval_bytes
total_visible_bytes
```

Their meanings are:

### `initial_visible_bytes`

Bytes presented to the agent by the initial benchmarked command results.

This is the compressed/stateful output produced before any explicit progressive-disclosure recovery operation.

### `show_bytes`

Total agent-visible bytes returned by benchmarked `acap show` retrieval operations.

### `raw_retrieval_bytes`

Total agent-visible bytes returned by benchmarked raw-output retrieval operations.

### `total_visible_bytes`

The total output actually exposed to the agent during the measured workflow:

```text
total_visible_bytes =
    initial_visible_bytes
  + show_bytes
  + raw_retrieval_bytes
```

If the existing architecture has additional already-measured agent-visible recovery output that clearly belongs in this total, include it consistently and document the rule.

Do not count internal storage bytes or internal intermediate representations as agent-visible bytes.

---

# 3. Keep Initial Compression and Effective Compression Separate

Do not report only the reduction achieved by the initial capsule.

The benchmark must make it possible to distinguish:

```text
initial reduction
```

from:

```text
effective reduction after recovery
```

For example:

```text
raw baseline                  381,920
initial AgentCap output        18,420
show retrieval                 4,810
raw retrieval                      0
                              -------
total agent-visible            23,230
```

The meaningful workflow-level comparison is:

```text
total agent-visible
vs
raw baseline
```

not merely:

```text
initial AgentCap output
vs
raw baseline
```

A benchmark result with excellent initial compression but expensive recovery must expose that fact directly.

---

# 4. Extend the Stable Benchmark Result Schema

Extend the B4 machine-readable result schema with the B5 recovery metrics.

Use names consistent with the current implementation.

Conceptually, schema results should expose:

```json
{
  "schema_version": 2,
  "workload": "example",
  "commands": 8,

  "raw_bytes": 381920,

  "initial_visible_bytes": 18420,
  "show_bytes": 4810,
  "raw_retrieval_bytes": 0,
  "total_visible_bytes": 23230,

  "processing_ns": 4812000
}
```

Retain existing B4 fields where they remain meaningful.

Do not remove useful existing metrics merely because B5 introduces a more complete total.

---

# 5. Handle Schema Versioning Correctly

B4 introduced an explicitly versioned benchmark result schema.

Determine whether adding the B5 fields is compatible with the project's chosen schema-versioning policy.

If the B4 contract treats field additions as a schema change, increment the version.

For example:

```go
const BenchmarkResultSchemaVersion = 2
```

If the repository explicitly defines additive optional fields as compatible with schema version 1, follow that established policy instead.

Do not silently change the serialized contract while ignoring the versioning rules introduced in B4.

Add or update tests accordingly.

---

# 6. Recovery Operations Must Be Explicit Workload Actions

B5 must measure actual recovery behavior.

Do not estimate recovery cost from stored raw output size.

The workload should be able to explicitly perform recovery operations.

Use the existing B3 workload model and extend it minimally if necessary.

Conceptually, a workload might contain:

```yaml
steps:
  - run: go test ./...

  - show: <result reference>

  - raw: <result reference>
```

The exact syntax must follow the architecture already implemented in B3.

Do not introduce the example syntax literally if a better step/action abstraction already exists.

The important property is that benchmark workloads can reproducibly specify:

```text
command
    ->
optional show
    ->
optional raw retrieval
```

and the bytes returned by those operations are measured.

---

# 7. Use Real AgentCap Retrieval Paths

Recovery benchmark steps must exercise the normal AgentCap retrieval implementation.

Do not simulate:

```text
show_bytes
raw_retrieval_bytes
```

by reading files directly from the benchmark code.

The benchmark should invoke the same underlying application/service path that normal:

```bash
acap show ...
acap raw ...
```

operations use.

CLI subprocess execution is not required if the codebase provides an appropriate internal API.

Prefer reuse of the normal internal execution path over spawning `acap` recursively.

The benchmark must measure real AgentCap behavior without unnecessarily coupling itself to CLI text parsing.

---

# 8. Do Not Re-Execute Original Commands During Recovery

Preserve the core execution invariant:

```text
Every intercepted command executes exactly once.
```

A `show` or `raw` recovery operation must retrieve information captured from the original execution.

It must not rerun the underlying command.

For example:

```text
go test ./...
     |
     | execute once
     v
captured result
     |
     +--> initial capsule
     |
     +--> acap show
     |
     +--> acap raw
```

Never:

```text
go test ./...       # initial result

go test ./...       # show

go test ./...       # raw
```

Add tests protecting this invariant.

---

# 9. Measure Agent-Visible Retrieval Output, Not Stored Payload Size

For recovery operations, measure the bytes actually presented to the agent.

For example:

```text
stored raw output:          150 KB
acap show output:             4 KB
```

must contribute:

```text
show_bytes += 4 KB
```

not 150 KB.

Likewise, if `acap raw` presents the complete raw result, measure the actual output emitted by that retrieval operation.

The benchmark KPI concerns:

```text
agent-visible context cost
```

not storage utilization.

---

# 10. Count Repeated Recovery Operations

If the workload performs multiple retrieval operations, count all of them.

For example:

```text
command A
show A
command B
show B
raw B
show A
```

must include the output from every retrieval operation.

Do not deduplicate retrieval bytes merely because the same result is requested more than once.

If the agent sees the same 5 KB twice, the workflow consumed approximately 10 KB of agent-visible output.

This is important because repeated recovery is itself a behavior B5 is intended to expose.

---

# 11. Track Recovery Counts

In addition to byte totals, track operation counts where practical:

```text
show_count
raw_retrieval_count
```

These are useful diagnostics because two workloads can have similar byte totals but very different fallback behavior.

Conceptually:

```json
{
  "show_count": 3,
  "raw_retrieval_count": 1
}
```

Use naming consistent with the current result schema.

These counts should reflect actual retrieval operations performed during the workload.

Do not count internal implementation calls that are invisible at the benchmark/workflow level.

---

# 12. Human-Readable Output

Extend the B4 human-readable result so recovery cost is obvious.

Prefer output conceptually similar to:

```text
Benchmark: test/fail-fix-pass

Agent-visible output

initial                    18,420
show                        4,810
raw retrieval                   0
                           ------
total                      23,230

raw baseline              381,920

effective reduction          93.9%

Recovery
show calls                       2
raw retrievals                   0

processing                    4.8ms
```

Exact formatting should follow existing CLI conventions.

Keep the output compact.

The important requirement is that a human can immediately see:

```text
initial cost
recovery cost
total cost
raw baseline
```

without manually calculating them.

---

# 13. Machine-Readable Output

`--json` must expose the recovery metrics as numeric fields.

For example:

```bash
acap bench run benchmark.yaml --json
```

should provide enough structured information for external tools to calculate or inspect:

```text
initial compression
recovery overhead
effective compression
show/raw frequency
```

Do not serialize byte counts or durations as human-formatted strings.

Good:

```json
"show_bytes": 4810
```

Bad:

```json
"show_bytes": "4.8 KB"
```

Likewise, preserve B4's requirement that stdout in JSON mode contains only the machine-readable result.

---

# 14. Prefer Primitive Metrics Over Derived Metrics

Persist authoritative primitive measurements such as:

```text
raw_bytes
initial_visible_bytes
show_bytes
raw_retrieval_bytes
total_visible_bytes
show_count
raw_retrieval_count
processing_ns
```

Derived percentages can normally be calculated at presentation time.

For example:

```text
effective_reduction =
    1 - total_visible_bytes / raw_bytes
```

Avoid creating multiple persisted fields that express the same underlying measurement unless there is a clear compatibility or analysis requirement.

This reduces the chance that stored results become internally inconsistent.

---

# 15. Validate Accounting Invariants

Where applicable, enforce or test:

```text
total_visible_bytes
=
initial_visible_bytes
+ show_bytes
+ raw_retrieval_bytes
```

Do not allow formatters to independently calculate different totals.

Prefer computing the canonical total in one place.

Also preserve:

```text
show_count >= 0
raw_retrieval_count >= 0

show_bytes >= 0
raw_retrieval_bytes >= 0
total_visible_bytes >= 0
```

Use the project's normal approach to invariants rather than adding unnecessary runtime checks everywhere.

---

# 16. Baseline Semantics

Keep the raw baseline conceptually separate from recovery.

The raw baseline represents the workflow cost without AgentCap compression.

The AgentCap workflow represents:

```text
initial visible output
+
recovery output
```

Do not add AgentCap recovery operations to the raw baseline.

The comparison is conceptually:

```text
Raw workflow:
    normal raw command output

AgentCap workflow:
    compressed command output
    + show retrieval
    + raw retrieval
```

This is what allows the benchmark to answer whether progressive disclosure reduces total context consumption.

---

# 17. Avoid Artificial Success Criteria

Do not automatically classify a benchmark as successful merely because initial compression exceeds some threshold.

For example:

```text
95% initial reduction
+
frequent acap raw fallback
```

must remain visible as a potentially expensive workflow.

B5 should expose the underlying metrics.

Do not introduce a simplistic:

```text
PASS if compression > X%
```

policy.

Later comparison/regression phases may decide how metrics are interpreted.

B5's responsibility is accurate measurement.

---

# 18. Processing Overhead

Continue measuring AgentCap processing overhead using the existing B0–B4 semantics.

Recovery operations may introduce additional AgentCap processing.

If the existing measurement architecture can naturally attribute recovery processing time, include it consistently in the benchmark's AgentCap processing total.

Do not include unrelated benchmark-runner overhead merely because B5 adds more steps.

If precise attribution would require a broad timing-system redesign, preserve the existing processing metric semantics and document the limitation rather than expanding B5 substantially.

---

# 19. Workload Reproducibility

Recovery-cost benchmarks must remain deterministic and reproducible at the workload level.

A workload should explicitly determine whether and when recovery occurs.

B5 does **not** yet ask an LLM agent to decide:

```text
"Do I need acap show?"
```

or:

```text
"Should I fall back to acap raw?"
```

That behavior belongs to B6 and later agent evaluation.

B5 instead measures deterministic scenarios such as:

```text
command
show

command
show
raw
```

This gives us a controlled way to validate the accounting model before introducing agent nondeterminism.

---

# 20. Fixture Safety

Preserve B3's fixture/workspace safety rules.

Recovery-cost benchmarks must operate on benchmark fixtures or temporary workspaces where state changes are required.

Do not modify the user's real repository merely to exercise recovery behavior.

Recovery operations themselves should normally be read-only against previously captured results.

---

# 21. Tests

Add focused tests for recovery accounting.

At minimum, cover the following cases.

## No recovery

```text
command
```

Expected:

```text
show_count = 0
raw_retrieval_count = 0

show_bytes = 0
raw_retrieval_bytes = 0

total_visible_bytes = initial_visible_bytes
```

## Show recovery

```text
command
show
```

Verify:

```text
show_count = 1
show_bytes > 0

total_visible_bytes =
    initial_visible_bytes + show_bytes
```

## Raw recovery

```text
command
raw
```

Verify:

```text
raw_retrieval_count = 1
raw_retrieval_bytes > 0
```

and that the original command was not executed again.

## Multiple recovery operations

Test something conceptually equivalent to:

```text
command
show
show
raw
```

Verify that every agent-visible retrieval is counted.

Do not deduplicate repeated retrievals.

## Mixed workflow

Test multiple commands with different recovery behavior:

```text
command A
show A

command B

command C
show C
raw C
```

Verify aggregate accounting across the complete workload.

## JSON output

Verify that B5 metrics are serialized using stable field names and numeric types.

## Human output

Verify that initial, recovery, total, and baseline costs are visible.

Avoid unnecessarily brittle whitespace assertions.

## Single-execution invariant

Instrument a fixture command so the test can prove that:

```text
command + show + raw
```

still executes the underlying command exactly once.

This is a critical regression test.

---

# 22. Documentation

Update benchmark documentation with a recovery-cost example.

Explain the distinction between:

```text
initial_visible_bytes
```

and:

```text
total_visible_bytes
```

Document that:

```text
total_visible_bytes
```

is the meaningful AgentCap-side context-cost metric when progressive disclosure is involved.

Also document that repeated `show` and `raw` operations are intentionally counted repeatedly because the agent receives those bytes repeatedly.

---

# 23. Do Not Implement B6 Yet

Do not add:

```text
Claude Code adapter
Codex adapter
agent process control
agent decision tracing
task-success evaluation
agent-specific prompt integration
```

B5 remains an AgentCap-only deterministic benchmark phase.

The benchmark workload explicitly specifies recovery operations.

B6 will later determine whether real coding agents actually request the right information and how often they fall back to raw output.

---

# 24. Do Not Implement B7–B9

Also exclude:

```text
repeated LLM trials
median/statistical aggregation
benchmark comparison command
regression thresholds
CI benchmark policy
```

Keep B5 focused on making recovery cost measurable.

---

# Implementation Guidance

Before modifying code:

1. inspect the B0–B4 benchmark measurement and result types;
2. inspect the B3 workload action representation;
3. inspect the normal implementation paths used by `acap show` and `acap raw`;
4. identify how result/session IDs are passed between workload steps;
5. reuse those paths rather than duplicating retrieval logic;
6. inspect B4 schema-versioning policy before modifying JSON fields.

Prefer minimal extensions to existing abstractions.

If the current implementation already has a generic workload action model, extend that model rather than creating benchmark-specific special cases.

If `show` and `raw` share an internal retrieval abstraction, use it.

Do not parse the human-readable output of `acap show` or `acap raw` if structured internal APIs already exist.

---

# Architectural Invariants

After B5, the benchmark accounting should conceptually be:

```text
                    captured command
                          |
                          v
                  initial presentation
                          |
                 initial_visible_bytes
                          |
             +------------+-------------+
             |                          |
             v                          v
           show                        raw
             |                          |
             v                          v
        show_bytes            raw_retrieval_bytes
             |                          |
             +------------+-------------+
                          |
                          v
                total_visible_bytes
```

The fundamental accounting invariant is:

```text
total AgentCap workflow cost
=
initial presentation cost
+
all progressive-disclosure recovery cost
```

And recovery must always operate on previously captured results:

```text
one command execution
       |
       v
captured result
       |
       +--> initial
       +--> show
       +--> raw
```

Never re-execute commands merely to recover information.

---

# Exit Criteria

Phase B5 is complete when all of the following are true:

1. Benchmark workloads can deterministically exercise progressive-disclosure recovery.
2. Initial AgentCap-visible output is measured separately from recovery output.
3. `show` output contributes to `show_bytes`.
4. raw retrieval contributes to `raw_retrieval_bytes`.
5. Every repeated recovery operation contributes its actual agent-visible bytes.
6. `total_visible_bytes` represents the complete measured AgentCap-side context cost.
7. `show_count` and `raw_retrieval_count` are available as recovery diagnostics.
8. Human-readable output clearly separates initial, recovery, total, and raw-baseline costs.
9. JSON output exposes stable machine-readable recovery metrics.
10. Benchmark result schema versioning remains correct after the new fields are introduced.
11. `show` and `raw` use normal AgentCap retrieval paths rather than benchmark simulations.
12. Recovery never re-executes the original command.
13. Tests cover no-recovery, show, raw, repeated-recovery, and mixed-workflow cases.
14. Existing B0–B4 benchmark behavior remains intact.
15. No coding-agent integration or statistical benchmark functionality from B6–B9 is introduced.

The resulting benchmark must make this distinction observable:

```text
small initial capsule
```

is not necessarily equivalent to:

```text
small complete workflow cost
```

The benchmark should measure the latter.
