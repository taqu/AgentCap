# Phase B4 — Stable Benchmark Result Format

## Objective

Implement a stable benchmark result/output layer for `acap bench`.

Phase B4 must clearly separate:

1. human-readable benchmark output, and
2. machine-readable benchmark output.

The machine-readable format must be explicitly versioned so benchmark results can be stored, consumed by CI or external analysis tools, and interpreted after new metrics are added in future phases.

This phase is about the **result contract and presentation layer**.

Do not expand the benchmark execution model beyond what is already implemented by Phases B0–B3.

---

## Scope

Implement stable result formatting for benchmark workloads, including:

```bash
acap bench run <workload>
```

and:

```bash
acap bench run <workload> --json
```

The default output should be concise and optimized for humans.

`--json` must produce a stable structured representation intended for programs.

Reuse the benchmark measurements and workload execution infrastructure already implemented.

Do not introduce a second measurement path merely for formatting.

---

## 1. Introduce a Stable Benchmark Result Model

Create an explicit result type representing the externally meaningful result of a benchmark run.

Use names consistent with the existing codebase, but conceptually it should contain at least:

```go
type BenchmarkResult struct {
    SchemaVersion int    `json:"schema_version"`
    Workload      string `json:"workload"`
    Commands      int    `json:"commands"`

    RawBytes       int64 `json:"raw_bytes"`
    StatelessBytes int64 `json:"stateless_bytes"`
    StatefulBytes  int64 `json:"stateful_bytes"`

    ShowBytes         int64 `json:"show_bytes"`
    RawRetrievalBytes int64 `json:"raw_retrieval_bytes"`

    ProcessingNS int64 `json:"processing_ns"`
}
```

Adapt the exact Go structure to the current architecture rather than duplicating existing measurement types unnecessarily.

The important architectural distinction is:

```text
internal measurements
        |
        v
stable BenchmarkResult
        |
        +--> human formatter
        |
        +--> JSON formatter
```

Formatting code must not independently reconstruct benchmark semantics from command output.

There should be one canonical result model from which all presentation formats are produced.

---

## 2. Add Explicit Schema Versioning

Every machine-readable benchmark result must contain:

```json
"schema_version": 1
```

Define the current version centrally rather than scattering the literal value throughout the code.

For example:

```go
const BenchmarkResultSchemaVersion = 1
```

or the equivalent appropriate for the existing package structure.

The version belongs to the serialized benchmark result schema, not to workload specifications or internal persistence formats.

Do not reuse an unrelated existing schema version.

---

## 3. Implement `--json`

Support:

```bash
acap bench run benchmarks/git/repeated-diff.yaml --json
```

The output should conceptually resemble:

```json
{
  "schema_version": 1,
  "workload": "git/repeated-diff",
  "commands": 8,
  "raw_bytes": 240182,
  "stateless_bytes": 41820,
  "stateful_bytes": 16201,
  "show_bytes": 0,
  "raw_retrieval_bytes": 0,
  "processing_ns": 4812000
}
```

Use the actual workload identity and measurements produced by the current implementation.

### JSON requirements

The JSON output must:

- be valid JSON;
- contain no human-oriented headers or explanatory text;
- write benchmark data to stdout;
- use stable field names;
- use numeric values for byte counts, counts, and durations;
- include `schema_version`;
- remain suitable for shell pipelines and external tools.

For example, this must work:

```bash
acap bench run benchmark.yaml --json | jq '.stateful_bytes'
```

Do not emit logs, progress messages, or informational prose to stdout when `--json` is active.

If diagnostics are necessary, follow the project's existing stderr conventions.

---

## 4. Human-Readable Output

Without `--json`:

```bash
acap bench run benchmarks/git/repeated-diff.yaml
```

produce a compact human-readable summary.

Prefer a small table or similarly compact representation.

For example:

```text
Benchmark: git/repeated-diff
Commands: 8

                     bytes
raw                240,182
stateless           41,820
stateful            16,201

stateful vs raw       -93.3%
stateful vs stateless -61.3%

processing             4.8ms
```

Exact spacing and table implementation should follow existing CLI conventions.

The human format may derive presentation-only values such as percentage reductions from the canonical measurements.

Do not add presentation-only values to the stable JSON schema unless they are genuinely part of the benchmark result contract.

Prefer storing primitive measurements and deriving values such as:

```text
reduction percentage
formatted duration
formatted byte counts
```

at presentation time.

---

## 5. Keep Result Construction Separate From Formatting

Avoid code where benchmark execution directly prints results.

Prefer a structure conceptually similar to:

```go
result, err := runner.Run(...)
if err != nil {
    ...
}

if jsonOutput {
    return formatter.WriteJSON(out, result)
}

return formatter.WriteHuman(out, result)
```

The exact APIs should fit the existing architecture.

The important invariant is:

```text
benchmark execution != benchmark presentation
```

This separation will be required by later phases such as:

```text
B7 repeated trials
B8 benchmark comparison
B9 CI regression checks
```

Those phases must be able to consume benchmark results without parsing human CLI output.

---

## 6. Preserve Existing Measurement Semantics

Do not change the meaning of measurements introduced by previous phases.

In particular, preserve the distinction between:

```text
raw bytes
stateless AgentCap-visible bytes
stateful AgentCap-visible bytes
show bytes
raw retrieval bytes
AgentCap processing time
```

Do not merge these merely to simplify the output layer.

Likewise, preserve the existing invariant that benchmark instrumentation must not cause commands to execute more than once.

B4 must only transform already collected measurements into result representations.

---

## 7. Forward Compatibility

Design schema version 1 so later phases can add metrics without requiring a redesign.

Future results may include fields such as:

```text
initial_visible_bytes
total_visible_bytes
show_count
raw_retrieval_count
task_success
wall_time
trial metadata
agent mode
```

Do not implement those future metrics unless they already exist naturally in the current measurement model.

However, avoid a design that would make adding them difficult.

In particular:

- do not encode metrics into positional arrays;
- do not encode numeric values as formatted strings;
- do not make human-readable output the persisted representation;
- do not require consumers to parse table text.

---

## 8. Do Not Implement B5 Yet

Phase B5 will make progressive-disclosure recovery cost a first-class benchmark concern.

Do not pull B5 into this phase.

It is acceptable for schema v1 to expose already available:

```json
"show_bytes": 0,
"raw_retrieval_bytes": 0
```

or their actual measured values if B0–B3 already provide them.

But do not introduce new recovery workflows, synthetic `show`/`raw` operations, or new progressive-disclosure benchmark behavior solely for B4.

B4 defines how existing results are represented.

---

## 9. Do Not Implement B6–B9

Explicitly exclude:

- coding-agent adapters;
- Claude Code/Codex integration;
- repeated LLM trials;
- statistical aggregation;
- benchmark-to-benchmark comparison;
- automatic regression thresholds;
- CI policy;
- task-success evaluation.

Those belong to later phases.

The B4 result model should make those features possible later, but must not implement them now.

---

## 10. Error Behavior

Preserve existing CLI error conventions.

For `--json`, benchmark execution failures must not result in a successful-looking benchmark result containing fabricated zero measurements.

A failed workload should continue to return an appropriate non-zero exit status.

Do not silently convert execution failures into:

```json
{
  "raw_bytes": 0,
  "stateful_bytes": 0
}
```

unless zero is genuinely the measured successful result.

Keep benchmark result data distinct from CLI/runtime errors.

---

## 11. Tests

Add focused tests for both the result model and CLI output.

At minimum, cover:

### JSON serialization

Verify that:

```text
schema_version
workload
commands
raw_bytes
stateless_bytes
stateful_bytes
show_bytes
raw_retrieval_bytes
processing_ns
```

are serialized with the intended stable names and numeric types.

### Schema version

Verify that newly produced results use:

```text
schema_version = 1
```

from the central schema-version definition.

### Valid JSON-only stdout

For:

```bash
acap bench run <fixture> --json
```

verify stdout can be parsed directly as JSON and does not contain human-readable prefixes or suffixes.

### Human output

Verify the normal command contains the essential benchmark measurements and remains reasonably compact.

Avoid brittle tests for exact whitespace unless the repository already uses golden-output tests for CLI formatting.

### Same underlying result

Where practical, test that human and JSON formatting operate on the same `BenchmarkResult` rather than independently calculating benchmark measurements.

### Execution invariant

Ensure introducing multiple formatters does not cause the workload or any command inside it to execute multiple times.

---

## 12. Documentation

Update benchmark CLI documentation/examples to show both forms:

```bash
acap bench run benchmarks/git/repeated-diff.yaml
```

and:

```bash
acap bench run benchmarks/git/repeated-diff.yaml --json
```

Document that the JSON result is versioned and intended for programmatic consumption.

Do not document schema v1 as permanently immutable.

The compatibility contract is:

```text
schema_version identifies how the serialized result should be interpreted.
```

---

## Implementation Guidance

Before changing code:

1. inspect the B0–B3 benchmark measurement/result structures;
2. inspect existing CLI formatting conventions;
3. inspect how other commands implement `--json`, if applicable;
4. reuse existing output abstractions where they fit;
5. avoid parallel benchmark-specific infrastructure unless necessary.

Prefer a small result/presentation layer over a broad refactor.

Do not modify reducer, session, capture, or workload semantics unless a minimal change is required to expose measurements that B4 needs.

If the current B0–B3 implementation differs from the roadmap, adapt B4 to the actual architecture rather than recreating the roadmap literally.

---

## Architectural Invariants

After Phase B4, the following should be true:

```text
workload execution
       |
       v
measurements
       |
       v
BenchmarkResult (schema v1)
       |
       +---- human-readable output
       |
       +---- JSON output
```

And:

```text
one workload execution
        ->
one canonical result
        ->
multiple possible representations
```

Never:

```text
human benchmark execution
+
separate JSON benchmark execution
```

or:

```text
human table
        ->
parse table
        ->
JSON
```

---

## Exit Criteria

Phase B4 is complete when all of the following are true:

1. `acap bench run <workload>` produces a concise human-readable benchmark summary.
2. `acap bench run <workload> --json` produces valid machine-readable JSON.
3. JSON contains an explicit `schema_version`.
4. Schema version 1 has stable, documented field names.
5. Numeric measurements are serialized as numeric values rather than formatted strings.
6. Human and JSON output originate from the same canonical benchmark result.
7. Output formatting does not execute benchmark commands again.
8. Existing B0–B3 benchmark behavior remains intact.
9. Tests cover JSON serialization, schema versioning, human output, and the single-execution invariant.
10. No B5–B9 functionality is introduced beyond what is necessary to keep the result model extensible.

The final implementation should make benchmark results reliable inputs for future CI, analysis, comparison, and agent-evaluation phases without coupling those consumers to the human CLI presentation.
