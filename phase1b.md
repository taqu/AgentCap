# Task: Implement AgentCap Benchmark CLI Phase B0 + B1

Implement the first benchmark infrastructure for AgentCap.

This task combines:

- **Phase B0 — Benchmark Measurement Foundation**
- **Phase B1 — Single-Command Benchmark CLI**

The goal is to establish a trustworthy measurement layer and expose the first user-facing benchmark command:

```bash
acap bench command -- <command...>
```

This phase measures a **single command execution** and compares its raw output cost with the AgentCap presentation cost.

Do not implement workflow benchmarks, repeated trials, benchmark fixture formats, coding-agent orchestration, or statistical comparison in this phase.

---

# 1. Core Goal

A user should be able to run:

```bash
acap bench command -- git diff
```

or:

```bash
acap bench command -- go test ./...
```

and receive a report showing, at minimum:

```text
raw output size
AgentCap-visible output size
reduction
AgentCap processing overhead
command execution duration
exit status
result ID, when applicable
```

The benchmark must use AgentCap's real execution, capture, reduction, and storage path.

Do not create a fake or simplified reducer path specifically for benchmarking.

---

# 2. Critical Correctness Invariant: Execute Exactly Once

This is the most important requirement.

The benchmark target command MUST execute exactly once.

Never implement benchmarking as:

```text
run command normally
    +
run command through AgentCap
```

That is incorrect.

Commands may have side effects:

```text
rm
mv
git add
git commit
go generate
database migrations
custom scripts
```

The correct model is:

```text
                 command
                    |
                    v
              execute once
                    |
                    v
             capture raw result
                    |
             +------+------+
             |             |
             v             v
      raw measurement   AgentCap processing
                             |
                             v
                       visible result
                             |
                             v
                    reduced measurement
```

Both benchmark sides must be derived from the **same captured execution result**.

Add tests that make accidental double execution detectable.

For example, use a test command or helper that increments/writes a value when executed and assert that the side effect occurs exactly once.

---

# 3. Inspect Existing Architecture First

Before implementing anything, inspect the current AgentCap codebase.

Identify and reuse the existing paths for:

- command execution
- stdout/stderr capture
- execution duration
- exit-code handling
- result persistence
- reducer/classifier selection
- capsule generation
- presentation rendering
- result IDs
- project-root resolution
- `.acap` storage
- command metadata
- existing statistics

Do not duplicate these systems inside the benchmark package.

The benchmark layer should primarily **observe and measure existing behavior**.

If refactoring is required to expose measurement points cleanly, keep the refactor small and behavior-preserving.

---

# 4. Phase B0 — Measurement Model

Introduce an internal benchmark measurement model.

Use names appropriate to the existing Go codebase rather than blindly copying the example below.

Conceptually, the model should contain data such as:

```go
type BenchmarkMeasurement struct {
    RawStdoutBytes     int64
    RawStderrBytes     int64
    RawBytes           int64

    AgentVisibleBytes  int64

    ExecutionDuration  time.Duration
    ProcessingDuration time.Duration

    ExitCode           int

    ResultID           string
}
```

Add fields only when they have clear semantics and can be measured reliably.

Avoid speculative metrics that are not needed by B0/B1.

---

# 5. Define Byte Semantics Precisely

Do not leave metrics such as `raw_bytes` or `agent_visible_bytes` ambiguous.

Define them in code comments and tests.

At minimum:

## Raw bytes

`raw_bytes` should represent the complete captured command output before AgentCap reduction.

Prefer:

```text
raw_bytes =
    raw stdout bytes
    +
    raw stderr bytes
```

If the existing execution model requires different semantics, document and test them explicitly.

Do not silently count terminal rendering metadata, database metadata, or filesystem storage overhead as raw command output.

## Agent-visible bytes

`agent_visible_bytes` should represent the actual AgentCap command-result presentation that would be returned to the coding agent for this execution.

Measure the rendered output, not an intermediate structured object.

For example, do not measure only:

```text
serialized capsule structure
```

if the actual agent receives:

```text
header
+
capsule
+
result ID
+
metadata
```

Measure what the agent actually sees.

This distinction is important because the long-term KPI is agent-visible context cost.

---

# 6. stdout and stderr

Preserve the existing AgentCap semantics for stdout and stderr.

The benchmark must not alter:

- stream capture
- error preservation
- reducer behavior
- exit-code behavior

Record stdout/stderr sizes separately if this is straightforward and useful.

At minimum ensure:

```text
raw_bytes = stdout + stderr
```

is correctly measured.

Test commands that produce:

- stdout only
- stderr only
- both stdout and stderr
- no output

---

# 7. Execution Duration vs AgentCap Processing Duration

Keep these measurements separate.

## Execution duration

Measure the actual target command execution.

Conceptually:

```text
process start
    ->
process completion
```

This includes the command's own runtime.

## AgentCap processing duration

Measure AgentCap's post-execution work relevant to producing the compact result.

This may include, depending on the existing architecture:

```text
classification
reduction
structuring
capsule generation
rendering
```

Be precise about the boundary.

Do not claim that this metric represents complete AgentCap overhead if it excludes meaningful work such as storage.

If storage is excluded, either:

- call the metric something precise such as `reduction_duration`, or
- document the exclusion clearly.

Prefer accurate naming over an impressive-looking number.

---

# 8. Do Not Benchmark Storage Size Yet

This phase is about **agent-visible command-result cost**, not `.acap` disk usage.

Do not treat:

```text
SQLite bytes
object-store bytes
database growth
```

as command-output benchmark metrics.

Storage benchmarking can be added separately later.

---

# 9. Implement `acap bench`

Introduce the benchmark command family without disturbing existing CLI behavior.

Target hierarchy:

```text
acap bench
    |
    +-- command
```

For this phase, only `command` needs to exist.

Expected syntax:

```bash
acap bench command -- <command...>
```

Examples:

```bash
acap bench command -- git status

acap bench command -- git diff

acap bench command -- go test ./...

acap bench command -- cargo check
```

Follow the existing CLI framework and conventions.

Do not introduce a second CLI parser or unrelated command-dispatch architecture.

---

# 10. Preserve argv Semantics

The command following `--` must use the same direct argv execution semantics as normal AgentCap execution.

Do not reconstruct it into a shell command string.

For example:

```bash
acap bench command -- printf '%s\n' 'hello world'
```

must preserve arguments correctly according to the existing AgentCap execution model.

Do not introduce implicit:

```bash
sh -c ...
```

unless the existing command execution path explicitly requires it.

---

# 11. Benchmark Execution Flow

The implementation should conceptually follow:

```text
acap bench command -- <command>
              |
              v
       existing executor
              |
              v
       capture complete result
              |
       +------+------+
       |             |
       v             v
 measure raw      existing
   bytes          AgentCap
                  pipeline
                      |
                      v
                 presentation
                      |
                      v
             measure visible bytes
                      |
                      v
                store result
                      |
                      v
              benchmark report
```

Adapt ordering where required by the existing implementation.

The important properties are:

1. exactly one target execution
2. real AgentCap reduction
3. normal result storage
4. measurements derived from that execution
5. benchmark reporting does not change the reduced result itself

---

# 12. Human-Readable Output

Default output should be concise and terminal-friendly.

A reasonable conceptual format is:

```text
Benchmark: git diff

command:
  exit:              0
  duration:          42.3ms

raw:
  stdout:           181240 B
  stderr:                0 B
  total:            181240 B

agentcap:
  visible:           12842 B
  reduction:          92.9%
  processing:          4.8ms

result:
  id: abc123
```

Do not copy this format mechanically if it conflicts with existing AgentCap output conventions.

Optimize for easy human comparison.

Avoid excessive decoration.

---

# 13. Reduction Calculation

When raw output is non-zero, calculate:

```text
reduction =
    1 - (agent_visible_bytes / raw_bytes)
```

and display it as a percentage.

Be careful about the interpretation.

This is:

```text
single-command output reduction
```

It is NOT:

```text
workflow token reduction
```

and NOT:

```text
coding-agent efficiency improvement
```

Do not label it as either.

---

# 14. Zero-Output Commands

Handle:

```text
raw_bytes = 0
```

explicitly.

Do not divide by zero.

For example, reduction may be rendered as:

```text
n/a
```

Do not invent a misleading 0% or 100% reduction unless there is a strong existing convention for it.

Add a test.

---

# 15. Commands That Fail

A command returning a non-zero exit status is still a valid benchmark execution.

For example:

```bash
acap bench command -- go test ./...
```

may fail because tests fail.

The benchmark should still:

- capture the complete result
- run normal AgentCap reduction
- store the result
- print benchmark measurements
- preserve/report the target command's exit status

Do not treat a non-zero target exit code as an internal benchmark failure.

Follow existing `acap run` exit-code semantics where practical.

---

# 16. Result Storage and Recoverability

Benchmark executions should preserve AgentCap's normal recoverability guarantees.

If normal command execution produces a stored result ID, benchmark execution should do so as well.

The result should remain inspectable using existing commands such as:

```bash
acap show <id>
acap raw <id>
```

Do not create benchmark-only result storage that bypasses `.acap/store.db` or the object store.

A user should be able to benchmark a command and then inspect exactly what was captured.

---

# 17. Do Not Introduce Stateful Comparison Yet

Phase B1 is intentionally single-command and effectively stateless from the benchmark user's perspective.

Do NOT implement:

```text
raw vs stateless vs stateful
```

comparison yet.

Do not add:

```text
baseline selection
repeated-command benchmark logic
workflow totals
delta benchmark scenarios
```

Those belong to Phase B2.

If the existing AgentCap execution path automatically has session behavior, ensure the benchmark has clearly defined behavior and does not accidentally produce misleading stateful measurements.

Prefer an isolated/stateless benchmark context for B1 unless the architecture makes another approach clearly safer.

Document the chosen behavior in code and tests.

---

# 18. Do Not Implement Token Estimation Yet

B0/B1 should use bytes as the canonical metric.

Do not add tokenizer dependencies merely to report approximate tokens.

Reasons:

- token counts depend on the model/tokenizer
- AgentCap itself is model-independent
- bytes are deterministic
- later benchmark phases can add tokenizer-aware analysis separately

The initial metric should therefore be:

```text
agent-visible bytes
```

not guessed tokens.

---

# 19. Optional `--json`

Add `--json` in B1 only if it fits naturally into the existing CLI/output architecture without materially expanding the task.

If implemented:

```bash
acap bench command --json -- git diff
```

should produce a stable structured representation such as:

```json
{
  "command": ["git", "diff"],
  "exit_code": 0,
  "raw_stdout_bytes": 181240,
  "raw_stderr_bytes": 0,
  "raw_bytes": 181240,
  "agent_visible_bytes": 12842,
  "reduction_ratio": 0.9291,
  "execution_duration_ns": 42300000,
  "processing_duration_ns": 4800000,
  "result_id": "abc123"
}
```

If implementing JSON would require premature benchmark-schema design, defer it to Phase B4.

Do not compromise B0/B1 architecture merely to add `--json`.

---

# 20. Measurement Must Not Affect Presentation

Avoid circular measurement behavior.

For example, the following is incorrect:

```text
AgentCap presentation
+
benchmark statistics
=
agent_visible_bytes
```

The benchmark report itself is not part of the command result that a coding agent would normally receive.

Therefore:

```text
agent_visible_bytes
```

must measure the normal AgentCap result presentation only.

The additional:

```text
Benchmark: ...
raw: ...
agentcap: ...
```

report is benchmark UI and should not be included in that value.

Add a test or clear code structure that makes this distinction obvious.

---

# 21. Tests

Add focused tests for the measurement foundation and CLI.

At minimum cover:

### Execute exactly once

Use an observable side effect and prove that:

```bash
acap bench command -- <side-effect-command>
```

executes the target once.

This is mandatory.

### Raw byte counting

Test:

```text
stdout only
stderr only
stdout + stderr
empty output
```

### Agent-visible byte counting

Verify that the measurement corresponds to the rendered AgentCap result.

### Zero raw output

Verify no divide-by-zero behavior.

### Successful command

Verify normal benchmark reporting.

### Failed command

Verify measurements are still generated and exit status is handled correctly.

### argv preservation

Verify arguments containing spaces and special-but-non-shell arguments are preserved correctly.

### Storage

Verify the generated result can be retrieved through the existing result store.

### Existing behavior

Run existing execution/reducer/storage tests to ensure the benchmark implementation does not change normal `acap run` behavior.

---

# 22. Prefer Reusable Instrumentation

Design B0 so later phases can reuse it.

Future phases will need to aggregate measurements approximately like:

```text
command #1 measurement
+
command #2 measurement
+
command #3 measurement
=
workflow measurement
```

Therefore avoid embedding all measurement logic directly inside the CLI handler.

Prefer a small internal API that can eventually be called by:

```text
bench command
bench workflow
bench agent
tests
```

without coupling core AgentCap behavior to benchmark CLI code.

However, do not build these future commands now.

---

# 23. Keep Benchmark Logic Outside Core Semantics

The architecture should remain conceptually:

```text
AgentCap Core
    |
    +-- execution
    +-- capture
    +-- reduction
    +-- storage
    +-- rendering

Benchmark Layer
    |
    +-- invokes/reuses core
    +-- observes measurements
    +-- reports results
```

Avoid:

```text
AgentCap Core
    |
    +-- if benchmark mode ...
```

spread throughout the codebase.

Small instrumentation hooks are acceptable where necessary, but benchmark-specific branching should remain localized.

---

# 24. No Benchmark-Specific Reduction

Do not tune reducers to make benchmark numbers look better.

`acap bench command` must exercise the same reducer behavior that normal AgentCap execution uses for the same captured result.

The benchmark is an observer of AgentCap behavior, not a special optimized mode.

---

# 25. Documentation

Update CLI help and relevant documentation.

At minimum:

```bash
acap bench --help
acap bench command --help
```

should explain that this benchmark measures a single command execution.

Document clearly that:

```text
reduction
```

means reduction in agent-visible bytes for that command, not end-to-end workflow savings.

If README already contains a benchmark section, update it only enough to document B1 accurately.

Do not document future B2+ functionality as implemented.

---

# 26. Out of Scope

Do NOT implement any of the following in this task:

- benchmark workload files
- YAML benchmark specifications
- multi-command workflows
- stateful benchmark comparison
- stateless vs stateful comparison
- repeated trials
- benchmark result databases
- benchmark history
- benchmark comparison commands
- coding-agent orchestration
- Claude/Codex benchmark runners
- task-success scoring
- LLM judging
- token estimation
- statistical analysis
- CI regression thresholds
- source-code semantic benchmarks
- Phase B2 or later functionality

Keep B0+B1 small and trustworthy.

---

# 27. Validation

Before finishing:

1. Build AgentCap.
2. Run all relevant existing tests.
3. Run the new benchmark tests.
4. Verify exactly-once execution explicitly.
5. Run at least one small successful command through `acap bench command`.
6. Run at least one command producing stderr.
7. Run at least one command returning non-zero.
8. Verify its stored result with `acap show`.
9. Verify its raw result with `acap raw`.
10. Confirm normal `acap run` behavior remains unchanged.

Also manually inspect the reported byte counts for at least one simple deterministic command.

---

# 28. Acceptance Criteria

Phase B0+B1 is complete when all of the following are true:

```text
[ ] A reusable benchmark measurement model exists.

[ ] Raw stdout/stderr/total bytes can be measured reliably.

[ ] Actual AgentCap-visible presentation bytes can be measured reliably.

[ ] Command execution duration is measured.

[ ] AgentCap processing duration has a precise documented boundary.

[ ] `acap bench command -- <command...>` exists.

[ ] The target command executes exactly once.

[ ] The benchmark uses the normal AgentCap reduction path.

[ ] The result remains stored and recoverable.

[ ] Non-zero command exits are benchmarkable.

[ ] Empty-output commands are handled correctly.

[ ] Reduction calculation is well-defined.

[ ] Benchmark-report bytes are not counted as agent-visible bytes.

[ ] No stateful/workflow benchmark functionality has leaked into B1.

[ ] Existing AgentCap behavior remains compatible.

[ ] Tests cover the critical measurement semantics and exactly-once invariant.
```

---

# 29. Final Report

When implementation is complete, report:

1. files added or modified
2. measurement model introduced
3. exact definitions of each reported metric
4. exact boundary used for processing duration
5. CLI syntax implemented
6. example benchmark output from a real test command
7. tests added
8. exactly-once execution test strategy
9. existing tests executed and their results
10. any architectural issues or ambiguities discovered

Also explicitly state whether `--json` was implemented or intentionally deferred.

Do not make unrelated refactors.
