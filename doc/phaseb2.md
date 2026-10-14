# Task: Implement AgentCap Benchmark CLI Phase B2 — Stateful Session Benchmarking

Implement Phase B2 of the AgentCap benchmark CLI.

Phase B0+B1 established measurement of a single command execution through:

```bash
acap bench command -- <command...>
```

Phase B2 extends this foundation to measure the value of AgentCap's existing session-aware compression across a sequence of command executions.

The primary goal is to make this comparison measurable:

```text
raw workflow cost
        vs
stateless AgentCap cost
        vs
stateful AgentCap cost
```

Do not implement benchmark workload files, fixture mutation, coding-agent orchestration, repeated statistical trials, or CI regression thresholds in this phase.

---

# 1. Inspect the Existing B0+B1 Implementation First

Before changing code, inspect the implementation produced by B0+B1.

Identify:

- the benchmark measurement model
- how raw bytes are measured
- how agent-visible bytes are measured
- execution-duration measurement
- AgentCap processing-duration measurement
- how `acap bench command` invokes the normal execution pipeline
- how benchmark executions are stored
- how benchmark sessions currently interact with normal AgentCap sessions
- any JSON output already implemented
- tests covering exactly-once execution

Reuse the existing measurement infrastructure.

Do not create a second measurement model for B2.

If B0+B1 differs from the original design instructions, treat the actual implementation as the starting point and preserve compatible behavior.

---

# 2. Core Goal

Provide a way to benchmark multiple commands within one controlled AgentCap benchmark session.

The benchmark must be able to report:

```text
raw total
stateless AgentCap total
stateful AgentCap total
```

for the commands executed during that session.

It should also report useful stateful behavior such as:

```text
full results
delta results
unchanged results
```

The important metric is the cumulative amount of command-result information that would be visible to an agent over the sequence.

---

# 3. Reuse the Existing Session and Delta Engine

AgentCap already has session-aware behavior.

Do not implement a benchmark-specific deduplication or delta algorithm.

The benchmark must exercise the same:

```text
command identity
baseline selection
unchanged detection
structural delta
full fallback
```

used by normal AgentCap operation.

Conceptually:

```text
command #1
    |
    v
full AgentCap presentation
    |
    v
session state

command #2
    |
    v
existing baseline selector
    |
    v
existing delta engine
    |
    +--> unchanged
    |
    +--> delta
    |
    +--> full fallback
```

The benchmark layer observes these results and records their sizes.

It must not decide independently whether two results are unchanged.

---

# 4. Preserve the Existing Baseline Rules

Use AgentCap's normal baseline selection rules.

Do not weaken command identity merely to make benchmark results look better.

In particular, preserve the existing conservative semantics around:

```text
same session
same command identity
latest equivalent result
```

and whatever actual command identity fields are implemented.

Do not compare unrelated commands merely because their output happens to be similar.

---

# 5. Three Costs Must Be Distinct

B2 must distinguish three different measurements.

## Raw cost

For each execution:

```text
raw cost =
    complete captured stdout
    +
    complete captured stderr
```

Workflow raw cost is:

```text
sum(raw cost of every command execution)
```

## Stateless AgentCap cost

This represents what AgentCap would have returned if each command result were reduced normally but no previous command result had been visible.

Conceptually:

```text
full reduced presentation #1
+
full reduced presentation #2
+
full reduced presentation #3
...
```

This is NOT raw output.

It is AgentCap reduction without session/delta savings.

## Stateful AgentCap cost

This is the actual presentation after the existing session layer and delta engine have been applied.

Conceptually:

```text
full reduced presentation #1
+
delta presentation #2
+
unchanged presentation #3
...
```

This distinction is central to Phase B2.

---

# 6. Do Not Execute Commands Multiple Times to Obtain the Three Costs

This is a critical correctness requirement.

Never implement:

```text
execute for raw
+
execute for stateless
+
execute for stateful
```

Each target command must execute exactly once.

Instead:

```text
                     execute once
                          |
                          v
                    captured result
                          |
          +---------------+---------------+
          |               |               |
          v               v               v
      raw bytes      full reduced      session/delta
                     presentation      presentation
                          |               |
                          v               v
                   stateless bytes   stateful bytes
```

All three measurements must be derived from the same captured execution.

Preserve and extend the B1 exactly-once tests.

---

# 7. Stateless Presentation Must Be Real

Do not estimate stateless cost from:

```text
raw_bytes * expected_ratio
```

or from historical statistics.

Generate or obtain the actual normal full AgentCap presentation for the captured result before session delta optimization.

The intended measurement point is approximately:

```text
capture
   |
   v
classifier/reducer
   |
   v
full capsule
   |
   v
full presentation  <---- stateless measurement
   |
   v
session layer
   |
   v
delta engine
   |
   v
stateful presentation <---- stateful measurement
```

Adapt this to the actual implementation.

If the current architecture discards or hides the pre-delta rendered representation, perform the smallest reasonable refactor to expose it.

Do not duplicate the reducer.

---

# 8. Measurement Model

Extend the B0 measurement model rather than replacing it.

Conceptually, each command measurement should now be able to represent:

```go
type BenchmarkMeasurement struct {
    RawStdoutBytes        int64
    RawStderrBytes        int64
    RawBytes              int64

    StatelessVisibleBytes int64
    StatefulVisibleBytes  int64

    ExecutionDuration     time.Duration
    ProcessingDuration    time.Duration

    PresentationKind      PresentationKind

    ExitCode              int
    ResultID              string
}
```

Use existing project naming conventions.

`PresentationKind` should reuse existing AgentCap concepts if they already exist.

Possible logical values are:

```text
full
delta
unchanged
```

Do not create a second incompatible classification if the session layer already exposes equivalent state.

---

# 9. Session Aggregate

Introduce a reusable aggregate measurement for a benchmark session.

Conceptually:

```go
type BenchmarkSessionMeasurement struct {
    CommandCount int

    RawBytes      int64
    StatelessBytes int64
    StatefulBytes  int64

    FullCount      int
    DeltaCount     int
    UnchangedCount int

    ExecutionDuration  time.Duration
    ProcessingDuration time.Duration
}
```

Again, adapt names and fields to the actual codebase.

The aggregate should be computed from individual command measurements rather than maintained through fragile duplicated counters where possible.

This structure should be reusable by future B3/B6 benchmark runners.

Do not implement those runners yet.

---

# 10. CLI Design

Add a minimal user-facing way to create and use a benchmark session.

First inspect the existing AgentCap session CLI and integration model.

Prefer reusing existing session concepts rather than introducing an unrelated benchmark-session subsystem.

Choose the smallest CLI that fits the current architecture.

A reasonable design may be something like:

```bash
acap bench session start
```

which returns a benchmark/session identifier, followed by:

```bash
acap bench command --session <id> -- git diff
acap bench command --session <id> -- go test ./...
acap bench command --session <id> -- git diff
```

and finally:

```bash
acap bench session show <id>
```

However, do NOT mechanically implement this syntax if the existing AgentCap session interface provides a cleaner solution.

Alternative designs are acceptable if they:

1. clearly define the benchmark session boundary
2. allow multiple independently invoked CLI processes to participate
3. reuse normal AgentCap session semantics
4. allow aggregate results to be inspected
5. do not require commands to be rerun

Document the chosen CLI design in the final report.

---

# 11. Benchmark Session Must Survive Separate CLI Invocations

Do not make B2 useful only inside one in-memory Go process.

A realistic benchmark sequence may be:

```text
shell invocation #1
shell invocation #2
shell invocation #3
```

The benchmark session therefore needs to use AgentCap's existing persistent session/state infrastructure where appropriate.

Do not create an in-memory-only session accumulator that loses state when `acap` exits.

---

# 12. Session Isolation

Benchmark sessions must not accidentally inherit unrelated previous AgentCap state.

For example, a benchmark should not become artificially smaller because the user previously ran:

```bash
acap run git diff
```

outside the benchmark.

A newly started benchmark session should have a clear baseline boundary.

Likewise, two benchmark sessions should not share delta baselines unless the normal session architecture explicitly requires it.

Add tests for isolation.

---

# 13. Do Not Disable Normal Result Storage

Each benchmarked command remains a normal independently complete AgentCap result.

Preserve:

```text
raw stdout
raw stderr
metadata
structured indexes
capsule
```

as required by the existing store.

The stateful presentation is only a presentation optimization.

Do not store only the delta.

After a benchmark execution:

```bash
acap show <result-id>
acap raw <result-id>
```

must continue to work.

---

# 14. Human-Readable Session Output

Provide a compact summary.

Conceptually:

```text
Benchmark Session: <id>

commands: 12

output:
  raw:                  428120 B
  stateless:             71840 B
  stateful:              28410 B

reduction:
  stateless vs raw:       83.2%
  stateful vs raw:        93.4%
  stateful vs stateless:  60.5%

presentations:
  full:                       4
  delta:                      5
  unchanged:                  3

time:
  command execution:        ...
  AgentCap processing:      ...
```

Use existing AgentCap formatting conventions where appropriate.

Keep the output compact.

The purpose is comparison, not verbose diagnostics.

---

# 15. Reduction Metrics

Calculate and label the comparisons precisely.

## Stateless reduction vs raw

```text
1 - stateless_bytes / raw_bytes
```

## Stateful reduction vs raw

```text
1 - stateful_bytes / raw_bytes
```

## Additional stateful savings

```text
1 - stateful_bytes / stateless_bytes
```

The last metric is especially useful because it isolates the benefit of session-aware compression from the benefit of ordinary reducer compression.

Use a clear label such as:

```text
stateful vs stateless
```

rather than claiming it is total AgentCap improvement.

Handle zero denominators explicitly as `n/a`.

---

# 16. Per-Command Detail

Provide a way to inspect the commands contributing to a session benchmark if this can be implemented cleanly.

For example, a session detail view might show:

```text
#  command             raw      stateless   stateful   kind
1  git diff            84 KB     11 KB       11 KB     full
2  go test ./...       31 KB      4 KB        4 KB     full
3  git diff            87 KB     12 KB        2 KB     delta
4  go test ./...       31 KB      4 KB       40 B      unchanged
```

Do not make this table mandatory if it would substantially complicate B2.

The aggregate metrics are the primary requirement.

---

# 17. Processing Duration

Preserve the exact timing definition established by B0+B1.

Do not silently change the meaning of `processing_duration`.

If B2 introduces measurable delta/session processing overhead, decide whether:

```text
processing duration
```

already includes it.

If it does not, either:

- extend the existing timing boundary consistently, or
- add a clearly named separate metric.

Do not mix incompatible timing definitions between B1 and B2.

Document the final timing boundary.

---

# 18. B1 Compatibility

Existing:

```bash
acap bench command -- <command...>
```

must continue to work.

Do not force all B1 users to explicitly create a benchmark session unless there is a compelling architectural reason.

Its existing single-command output and semantics should remain compatible where practical.

If B1 is intentionally stateless, keep that behavior.

B2 session-aware behavior should be explicitly requested through the benchmark-session mechanism.

---

# 19. JSON Output

If B1 already implements `--json`, extend it consistently.

Do not invent an unrelated JSON representation.

A session summary may conceptually contain:

```json
{
  "session_id": "...",
  "command_count": 12,
  "raw_bytes": 428120,
  "stateless_bytes": 71840,
  "stateful_bytes": 28410,
  "full_count": 4,
  "delta_count": 5,
  "unchanged_count": 3
}
```

Preserve existing schema conventions.

If B1 intentionally deferred JSON until B4, continue to defer it.

Do not let JSON schema work expand B2 unnecessarily.

---

# 20. No Token Estimation

Continue using bytes as the canonical deterministic measurement.

Do not add model-specific tokenizers.

B2 should measure:

```text
raw bytes
stateless agent-visible bytes
stateful agent-visible bytes
```

Token-aware evaluation belongs to a later agent benchmark phase.

---

# 21. Important Test Scenarios

Add deterministic tests covering at least the following.

## A. Repeated identical result

Run an equivalent command twice in the same benchmark session with identical output.

Expected conceptual behavior:

```text
first  -> full
second -> unchanged
```

Verify:

```text
stateful_bytes < stateless_bytes
```

for the second result where appropriate.

---

## B. Changed result

Use a controlled command or fixture where output changes between executions.

Expected:

```text
first  -> full
second -> delta or full fallback
```

according to the existing AgentCap delta implementation.

Do not require delta if the normal reducer correctly chooses full fallback.

Verify that the benchmark reports the actual presentation kind.

---

## C. Different command identity

Run two different commands that happen to produce identical output.

They must not be treated as unchanged merely because their output matches.

---

## D. Session isolation

Run the same command in two different benchmark sessions.

The second session must not use the first session as its baseline.

---

## E. Cross-process persistence

Create a benchmark session in one CLI invocation and execute subsequent commands through separate invocations.

Verify that stateful behavior still works.

---

## F. Exactly-once execution

Extend the B1 side-effect test to session benchmarking.

Every target command must still execute exactly once.

---

## G. Failed commands

Repeated failing compiler/test commands should remain benchmarkable.

Verify that non-zero exit status does not break aggregation.

---

## H. Empty output

Ensure zero-output commands do not break reduction calculations or session aggregation.

---

## I. Result recoverability

Verify that every benchmark result can still be retrieved independently using the existing `show` and `raw` paths.

---

# 22. Test Stateful Savings Without Faking Them

Do not unit-test B2 by manually assigning:

```text
kind = unchanged
stateful_bytes = 20
```

for all important integration tests.

At least some tests must exercise the actual:

```text
executor
-> reducer
-> session
-> baseline selector
-> delta engine
-> renderer
```

path.

Unit tests for aggregation arithmetic are useful, but they are not sufficient.

---

# 23. Benchmark Metadata

Store only the benchmark metadata necessary to reconstruct or aggregate B2 results.

Avoid prematurely building a full benchmark database.

If existing result/session metadata can represent:

```text
benchmark session membership
stateless visible bytes
stateful visible bytes
raw bytes
presentation kind
```

reuse or minimally extend it.

Do not create a parallel benchmark result store unless clearly necessary.

The benchmark system should remain subordinate to the existing `.acap` storage architecture.

---

# 24. Do Not Confuse AgentCap Session with Benchmark Result History

Keep the concepts clear.

An AgentCap session exists to determine what information the agent has already seen.

A benchmark session exists to define which executions should be measured together.

They may map 1:1 in B2, and that is likely desirable, but do not accidentally couple them in a way that prevents future benchmark runners from controlling measurement boundaries.

If they are mapped 1:1, document that design explicitly.

---

# 25. No Artificial Repetition

Do not add a flag such as:

```bash
--repeat 100
```

in B2.

Executing the same command repeatedly automatically would:

- risk side effects
- distort realistic workflow behavior
- overlap with later repeated-trial benchmarking

B2 measures commands that the benchmark/user explicitly executes.

---

# 26. No Workload Specification Yet

Do not implement:

```text
benchmark.yaml
steps:
  - run: ...
  - edit: ...
```

in this phase.

That belongs to Phase B3.

For B2, users or tests manually perform the sequence.

This keeps the session measurement foundation independent from the future workload format.

---

# 27. No Coding-Agent Orchestration

Do not invoke:

- Claude Code
- Codex
- Gemini CLI
- OpenCode
- any LLM API

from the B2 benchmark implementation.

B2 measures AgentCap's deterministic stateful compression layer.

Real coding-agent workflow evaluation belongs to a later phase.

---

# 28. Architecture Boundary

Keep the architecture conceptually:

```text
                    AgentCap Core
                         |
        +----------------+----------------+
        |                |                |
        v                v                v
     Executor          Reducer         Storage
                         |
                         v
                   Full Presentation
                         |
                         v
                    Session Layer
                         |
                         v
                     Delta Engine
                         |
                         v
                Stateful Presentation
                         |
             +-----------+-----------+
             |                       |
             v                       v
          Agent                  Benchmark
                                 Measurement
```

Benchmarking should observe the existing pipeline.

Avoid spreading:

```go
if benchmarkMode {
    ...
}
```

through reducers, session logic, or storage.

---

# 29. Future Compatibility

Design B2 so Phase B3 can later drive the same API with a reproducible workload.

Future architecture should be able to become:

```text
B3 workload runner
       |
       v
B2 benchmark session API
       |
       v
AgentCap execution/session pipeline
```

Likewise, later coding-agent benchmarks should be able to aggregate the same measurement structures.

Do not implement B3 now.

Just avoid designs that would force B2 to be rewritten.

---

# 30. Documentation

Update relevant CLI help and benchmark documentation.

Document:

- how to start/use a stateful benchmark session
- how to execute commands inside it
- how to inspect the aggregate
- what raw/stateless/stateful mean
- how reduction percentages are calculated
- that every target command executes once
- that stateful savings come from normal AgentCap session/delta behavior

Make clear that B2 measures deterministic AgentCap output reduction.

It does NOT measure:

```text
task success
model quality
coding-agent accuracy
token usage
```

yet.

---

# 31. Out of Scope

Do NOT implement the following:

- YAML/JSON workload definitions
- automatic workspace mutations
- fixture repositories
- benchmark task catalogs
- repeated benchmark trials
- statistical aggregation across trials
- coding-agent execution
- task-success scoring
- tokenizers
- LLM judges
- benchmark comparison across historical runs
- CI regression thresholds
- semantic source benchmarks
- adaptive reducer tuning
- Phase B3+ functionality

---

# 32. Validation

Before finishing:

1. Build AgentCap.
2. Run all existing tests.
3. Run all B0+B1 benchmark tests.
4. Run the new B2 tests.
5. Start a real benchmark session.
6. Execute a command that produces a full presentation.
7. Execute an equivalent command with unchanged output.
8. Verify the second execution uses normal unchanged behavior.
9. Execute a controlled changed-output scenario.
10. Verify delta/full-fallback behavior matches normal AgentCap behavior.
11. Inspect the session aggregate.
12. Verify raw/stateless/stateful totals manually for a small deterministic example.
13. Verify each target command executed exactly once.
14. Verify results using `acap show`.
15. Verify raw results using `acap raw`.
16. Verify a new benchmark session does not inherit the previous benchmark's baseline.
17. Confirm normal non-benchmark AgentCap behavior remains unchanged.

---

# 33. Acceptance Criteria

Phase B2 is complete when:

```text
[ ] Multiple command executions can be grouped into a benchmark session.

[ ] Benchmark sessions survive separate CLI invocations.

[ ] Benchmark sessions have isolated baseline boundaries.

[ ] Every target command executes exactly once.

[ ] Raw cost is measured from the complete captured result.

[ ] Stateless cost is measured from the real pre-delta AgentCap presentation.

[ ] Stateful cost is measured from the real session/delta presentation.

[ ] No command is rerun to obtain those three measurements.

[ ] Workflow totals can be calculated for raw/stateless/stateful bytes.

[ ] Full/delta/unchanged presentation counts are reported.

[ ] Existing AgentCap baseline and delta logic is reused.

[ ] Every stored result remains independently complete.

[ ] Existing `show` and `raw` recovery continues to work.

[ ] B1 single-command benchmarking remains compatible.

[ ] Session isolation is tested.

[ ] Cross-process session persistence is tested.

[ ] Exactly-once execution is tested.

[ ] No workload-runner or coding-agent functionality has leaked into B2.
```

---

# 34. Final Report

When implementation is complete, report:

1. files added or modified
2. CLI syntax chosen for benchmark sessions
3. how benchmark sessions map to AgentCap sessions
4. exact definition of `raw_bytes`
5. exact definition of `stateless_bytes`
6. exact definition of `stateful_bytes`
7. where the pre-delta/stateless measurement is taken
8. where the post-delta/stateful measurement is taken
9. how full/delta/unchanged are classified
10. how session isolation is guaranteed
11. how cross-process persistence works
12. tests added
13. exactly-once test strategy
14. example output from a real B2 benchmark
15. existing tests executed and results
16. any architectural issues discovered

If the implementation required a core refactor to expose the pre-delta presentation, describe that refactor and explain why it does not change normal AgentCap semantics.

Do not make unrelated refactors.
