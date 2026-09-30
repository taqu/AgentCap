# Reproducible benchmark workloads

Run one trial with:

```text
acap bench run benchmarks/workloads/git/repeated-diff.yaml
acap bench run --verbose benchmarks/workloads/build/compile-fix.yaml
acap bench run benchmarks/workloads/git/repeated-diff.yaml --json
```

Every invocation validates the complete definition, copies the fixture into a
fresh `acap-bench-*` temporary directory, and creates a fresh B2 benchmark
session. `run` steps are measured by the existing B2 API. Mutations only alter
the temporary workspace and do not create results, contribute bytes, or enter
session history. The workspace is removed on success or failure by default;
`--keep-workspace` retains it and prints its path.

Workload command results are stored in the project from which `acap` is
invoked, rather than in the disposable workspace. Every result ID shown by
`--verbose` remains available to `acap show <id>` and `acap raw <id>`.

## Layout

```text
benchmarks/
  workloads/       versioned YAML definitions
  fixtures/        self-contained initial workspace trees and mutation states
```

A workload must be below `benchmarks/workloads`, and its `fixture` must resolve
below the matching `benchmarks/fixtures`. Fixture trees and mutation paths may
not contain symlinks. Workloads are executable repository content, like scripts
or tests; path checks are defense in depth, not a command sandbox.

## Version 1 schema

Every file requires `version: 1`, a stable logical `name`, a relative
`fixture`, and one or more `steps`. Unknown fields, unknown versions, ambiguous
steps, empty `argv`, absolute paths, and escaping `..` paths are rejected before
commands run.

```yaml
version: 1
name: example/change-and-check
fixture: ../../fixtures/example
git:
  init: true
steps:
  - run:
      argv: ["git", "diff"]
      cwd: "."
      expect:
        exit: 0
  - copy:
      from: .states/next.txt
      to: data.txt
  - write:
      path: note.txt
      content: "deterministic content\n"
  - mkdir:
      path: output
  - remove:
      path: obsolete.txt
  - show:
      command: 1
  - raw:
      command: 1
      stream: stdout
```

Supported steps are:

- `run`: structured `argv`, optional workspace-relative `cwd`, and optional
  `expect.exit`. It invokes no shell and executes exactly once.
- `copy`: copy a file or directory from the fixture tree to the workspace.
- `write`: replace or create a file whose parent already exists.
- `remove`: remove a workspace child.
- `mkdir`: create a workspace directory and missing parents.
- `show`: retrieve the normal stored capsule for a preceding, 1-based `run`
  number through the same path as default `acap show`.
- `raw`: retrieve `stdout` (the default) or `stderr` for a preceding, 1-based
  `run` number through the same path as `acap raw`.

Only `run` enters B2 command measurement and command counts. `show` and `raw`
measure recovery bytes and calls without rerunning the referenced command;
mutation steps remain unmeasured. A target command's non-zero status is ordinary
workload data unless `expect.exit` is present and does not match. A missing
executable, unsafe path, failed mutation, retrieval error, or failed assertion
is a workload infrastructure failure identifying the step number and type.

When `git.init` is true, setup initializes and commits the copied fixture before
the first step. It uses repository-local identity, disables signing and hooks,
uses fixed author/committer dates, ignores system/global Git configuration, and
does not use a network remote. Git setup is harness work, not a measured step.

The bundled workloads are:

- `git/repeated-diff`: changed diff, changed diff, then unchanged diff.
- `build/compile-fix`: two compile failures with fewer errors, then success.
- `test/fail-fix-pass`: failing test, passing test, then a repeated pass.
- `recovery/show-and-raw`: one captured Git diff followed by two `show`
  retrievals and one raw retrieval.

They use Git and the Go standard toolchain only and require no network access.

## Stable result JSON

`--json` writes one JSON object and no human headers or command output to
stdout. Schema version 4 retains progressive-disclosure recovery measurements
and adds optional coding-agent trial metadata:

```json
{
  "schema_version": 4,
  "workload": "recovery/show-and-raw",
  "commands": 1,
  "raw_bytes": 221,
  "stateless_bytes": 95,
  "stateful_bytes": 95,
  "initial_visible_bytes": 95,
  "show_bytes": 176,
  "raw_retrieval_bytes": 221,
  "total_visible_bytes": 492,
  "show_count": 2,
  "raw_retrieval_count": 1,
  "processing_ns": 151000000
}
```

Counts, byte measurements, and duration are numeric. `initial_visible_bytes`
is the B2 stateful presentation cost. `total_visible_bytes` is computed once as
initial + every show + every raw retrieval and is the meaningful AgentCap-side
context cost when recovery occurs. Repeated retrievals are deliberately counted
again because those bytes are shown to the agent again. `processing_ns` includes
the existing command processing time plus measured AgentCap retrieval work, but
not general workload-runner overhead.

The result schema version is independent of workload definition version 1 and
the internal store schema. Consumers should select their interpretation using
`schema_version`; later schema versions may add or change fields.

## Coding-agent workloads

Coding-agent workloads use the same versioned fixture format, replacing
deterministic `steps` with task text and benchmark-only verification:

```yaml
version: 1
name: agent/go-bugfix
fixture: ../../fixtures/agent-go-bugfix
git:
  init: true
task: |
  Fix the implementation bug that causes the existing test to fail.
  Do not modify the test.
timeout: 10m
verify:
  - run:
      argv: ["go", "test", "./..."]
      expect: {exit: 0}
```

Run one trial with an explicit agent and mode, or repeat the same fixed
configuration sequentially:

```text
acap bench agent --workload benchmarks/workloads/agent/go-bugfix.yaml --agent codex --mode integrated
acap bench agent --workload benchmarks/workloads/agent/go-bugfix.yaml --agent codex --mode disabled --json
acap bench agent --workload benchmarks/workloads/agent/go-bugfix.yaml --agent codex --mode integrated --repeat 5 --json
```

The initial real adapter is Codex CLI. It uses non-interactive `codex exec`, an
ephemeral agent session, workspace-write sandboxing, automatic approval review,
captured JSONL logs, and an explicit timeout. Codex authentication must already
be configured. `--model` optionally selects a model; otherwise the installed
Codex default is used.

The four modes are:

- `disabled`: no AgentCap project hook. Command count and visible bytes come
  from Codex command-execution events, and the raw baseline equals visible
  command output.
- `stateless`: the AgentCap hook reduces each intercepted command without a
  session baseline.
- `stateful`: the hook uses the agent session for full/delta/unchanged behavior,
  but no AgentCap recovery guidance is added to the task workspace.
- `integrated`: stateful interception plus stable AgentCap `show`/`raw`
  instructions in the workspace. The agent decides whether to recover output.

Enabled-mode metrics are isolated store-stat deltas produced by the normal
AgentCap engine and retrieval paths. Disabled metrics come from agent JSONL
command events. The benchmark measures AgentCap-controlled command output; it
does not claim to measure model reasoning, hidden prompts, provider token use,
or native non-command tool payloads.

After the agent terminates, every `verify.run` executes directly through the
plain command executor. Its output determines `task_success` but is not included
in agent-visible bytes or command counts. A failed verification is a valid trial
with `task_success: false`; failure to start or configure the agent is an
infrastructure error. Timeout and non-zero agent exit are recorded as
`execution_status` values rather than fabricated successful results.

Schema version 4 agent results expose `agent`, `mode`,
`task_success`, `execution_status`, `agent_exit_code`, and `wall_time_ns`.
Complete agent wall time is separate from `processing_ns`, which records only
AgentCap command/retrieval processing.

### Repeated coding-agent trials

`--repeat N` requires `N >= 1` and runs the existing single-trial path N times
sequentially. Every trial receives a fresh fixture workspace, isolated AgentCap
store, and distinct stateful session scope. Task text, verifier, mode, timeout,
model selection, and integration instructions remain constant. `--repeat 1` is
the default and retains the compact single-trial B6 output.

For repeated runs, schema version 4 emits group metadata, requested and
completed trial counts, completion status, every indexed trial, and one
canonical `aggregate`. The aggregate reports objective success/failure counts,
timeout and agent-error counts, median total visible bytes, commands, complete
agent wall time, AgentCap processing time, and show/raw usage. Integer medians
use the midpoint of the two central values and round down for a half-unit.

All valid executed trials participate in workflow-cost medians, including task
failures, timeouts, and non-zero agent exits. They are never replaced by zero or
discarded merely because verification failed. Infrastructure failure stops the
remaining repetitions and marks the partial group `incomplete`; requested and
completed counts make missing trials explicit. Individual observations remain
in `trials`, so the median never hides outliers.

B7 does not compare modes, rank configurations, calculate significance, or
apply CI thresholds. Those remain separate later phases.

## Comparing results

```text
acap bench compare [--json] <baseline.json> <candidate.json>
```

`acap bench compare` reads two JSON results previously written by
`acap bench run --json` or `acap bench agent --json`. It never executes a
workload, coding agent, verifier, or AgentCap processing, and never modifies
its inputs.

- The first input is the **baseline**, the second the **candidate**.
- Every absolute delta is `candidate - baseline`; a negative byte delta means
  the candidate used fewer bytes, a positive wall-time delta means it took
  longer.
- The relative delta is `(candidate - baseline) / baseline` and is `null`
  (omitted in the human view) when the baseline is zero.

Dimensions reported independently:

| dimension | repeated results | single results |
|---|---|---|
| task success | `success_count/trial_count` per side | `1/1` or `0/1` |
| agent-visible bytes | `median_total_visible_bytes` | `total_visible_bytes` |
| command count | `median_command_count` | `command_count` |
| wall time | `median_wall_time_ns` | `wall_time_ns` (agent results only) |
| show/raw recovery | `trials_with_show`, `trials_with_raw_retrieval`, `total_show_count`, `total_raw_retrieval_count` | `show_count`, `raw_retrieval_count`, `show_bytes`, `raw_retrieval_bytes` |
| AgentCap processing | `median_processing_ns` | `processing_ns` |

Repeated results are compared using their recorded B7 aggregates; medians are
not recomputed. Timeouts, agent errors, and cancellations are also reported as
per-side outcome counts. Outcome counts are shown as `count/total` side by side
and are never reduced to a delta. Summed totals (`total_show_count`,
`total_raw_retrieval_count`) get a delta only when both sides have the same
trial count.

Compatibility: both results must have the same workload name, the same agent
(or both be deterministic workload results), and the same result kind — a
single observation is never compared with a repeated-trial median. Modes may
differ; that is the usual purpose of a comparison. Repeated results may have
different trial counts; each side keeps its own trials and medians, and trial
counts are displayed. Supported inputs are single results with schema versions
2–4 and repeated results with schema version 4. Incompatible or unsupported
inputs fail with a non-zero exit and no stdout.

Missing values are not zero. Each side of each metric carries a status:
`measured`, `unavailable` (not recorded, for example no trials or no wall time
in deterministic workload results), or `not_applicable` (AgentCap processing
and show/raw recovery in `disabled` mode). The human view prints `n/a` and `-`
respectively, and no delta is computed involving such a value.

The comparison JSON has its own `schema_version` (currently 1), independent of
the benchmark result schema:

```json
{
  "schema_version": 1,
  "kind": "repeated",
  "baseline":  {"schema_version": 4, "workload": "go-bugfix", "agent": "codex", "mode": "disabled", "trial_count": 5, "requested_trial_count": 5, "run_status": "completed"},
  "candidate": {"schema_version": 4, "workload": "go-bugfix", "agent": "codex", "mode": "integrated", "trial_count": 5, "requested_trial_count": 5, "run_status": "completed"},
  "outcomes": [
    {"name": "task_success", "baseline": {"count": 5, "total": 5}, "candidate": {"count": 4, "total": 5}, "baseline_status": "measured", "candidate_status": "measured"}
  ],
  "metrics": [
    {"name": "median_total_visible_bytes", "unit": "bytes", "baseline": 412000, "candidate": 91000, "baseline_status": "measured", "candidate_status": "measured", "delta": -321000, "relative_delta": -0.7791}
  ]
}
```

The command deliberately produces no composite score, winner, or pass/fail
status, and its exit code does not depend on metric changes. Regression policy
is a separate later phase.
