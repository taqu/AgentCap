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
| output sizes | — | `raw_bytes` (reference), `stateless_bytes`, `stateful_bytes` |
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

The comparison JSON has its own `schema_version` (currently 2; version 2 added
the single-result output-size metrics), independent of the benchmark result
schema:

```json
{
  "schema_version": 2,
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
is applied only by `acap bench check` (below).

## Regression checks

`acap bench compare` is **descriptive**: it reports what changed and always
exits 0 for a valid comparison. `acap bench check` is **policy-driven**: it
applies explicit thresholds to the same comparison and sets the exit status.

```text
acap bench suite --out <dir> benchmarks/regression.yaml
acap bench check --baseline <path> --candidate <path> --policy benchmarks/regression.yaml [--json]
```

### The deterministic suite

[`regression.yaml`](regression.yaml) is the single source of both the CI
workload list and the thresholds. It selects small, fast, deterministic
workloads covering git diff reduction, repeated/unchanged/delta stateful
output, build and test error transitions, and show/raw recovery.
`acap bench suite` runs each listed workload exactly as `acap bench run` does
(fresh temporary workspace and benchmark session) and writes one result per
workload, named `<workload with / replaced by _>.json`, into an output
directory that must be new or empty. Coding-agent workloads are rejected.

### Policy

```yaml
version: 1
workloads:
  - workloads/git/repeated-diff.yaml   # relative to this file
metrics:
  stateful_bytes:
    max_relative_increase: 0.10        # fraction: 0.10 = +10%
    max_absolute_increase: 32          # in the metric's unit (bytes, ns, count)
```

- Every rule applies to every suite workload separately. Workloads are never
  summed or averaged, so one workload's improvement cannot hide another's
  regression.
- The increase is the B8 `delta = candidate - baseline`; the relative change is
  the B8 `relative_delta`. The checker performs no delta arithmetic of its own.
- A decrease or no change always passes. An increase **equal** to a limit
  passes.
- With both limits, a check is a regression only when the increase exceeds
  both — the absolute limit is a noise floor for small values (10 B → 20 B is
  +100% but only +10 B).
- A zero baseline makes the relative change undefined: the absolute limit
  decides if configured, otherwise the check is `not_evaluable`.
- Policy metric names are the B8 comparison metric names: `stateless_bytes`,
  `stateful_bytes`, `total_visible_bytes`, `show_bytes`,
  `raw_retrieval_bytes`, `command_count`, `show_count`, `raw_retrieval_count`,
  `processing_ns`, `wall_time_ns`, and the repeated-result
  `median_*`/`total_*` metrics. Unknown names, unknown keys, duplicate
  metrics, negative or non-numeric limits, and rules without a limit are
  rejected, so a typo cannot silently drop a check.
- `raw_bytes` cannot be gated. Raw output measures the fixture and tools, not
  AgentCap; a raw change is reported as a reference note only.
- Deterministic workload results do not record wall time, so the checked-in
  policy gates AgentCap `processing_ns` with a wide tolerance instead of
  wall time. Byte metrics are exact and use tight limits.

Thresholds are configuration: the tool never adjusts them and never replaces a
baseline.

### Outcomes and exit status

Each check is `pass`, `regression`, or `not_evaluable`; each workload is
`evaluated`, `missing` (no baseline or candidate result), or `incompatible`.
The evaluation **fails** when any check regresses or is not evaluable, when a
suite workload is missing or incompatible, or when nothing was evaluated.
Missing measurements are never treated as zero.

| exit | meaning | JSON `status` |
|---|---|---|
| 0 | every check passed | `passed` |
| 1 | regression or not-evaluable check | `failed` |
| 2 | evaluation could not be completed (invalid policy, missing or invalid result files) | `error` |

The evaluation JSON has its own `schema_version` (currently 1). Each check
embeds the B8 comparison metric plus its limits and status:

```json
{
  "schema_version": 1,
  "status": "failed",
  "passed": false,
  "summary": {"workloads": 4, "failed_workloads": 1, "checks": 20, "passed": 19, "regressions": 1, "not_evaluable": 0},
  "workloads": [
    {
      "workload": "git/repeated-diff",
      "status": "evaluated",
      "raw_bytes": {"name": "raw_bytes", "unit": "bytes", "baseline": 716, "candidate": 716, "delta": 0, "relative_delta": 0, "baseline_status": "measured", "candidate_status": "measured"},
      "checks": [
        {"name": "stateful_bytes", "unit": "bytes", "baseline": 203, "candidate": 290, "baseline_status": "measured", "candidate_status": "measured",
         "delta": 87, "relative_delta": 0.4286, "max_relative_increase": 0.1, "max_absolute_increase": 32, "status": "regression"}
      ]
    }
  ]
}
```

### CI

[`.github/workflows/benchmark-regression.yml`](../.github/workflows/benchmark-regression.yml)
runs on pull requests. It builds `acap` from the pull request's base revision
(the baseline) and from the pull request (the candidate), runs
`acap bench suite` with each binary against the **candidate's** suite file,
workloads, and fixtures, and then runs `acap bench check`. Using the same
inputs on both sides isolates AgentCap code changes from benchmark input
changes, and running both on the same host keeps latency comparable. Each side
uses its own `ACAP_ROOT` store. The baseline, candidate, and evaluation JSON
are uploaded as the `benchmark-regression-results` artifact. The workflow
contains no benchmark semantics of its own.

If the base revision predates `acap bench suite`, the workflow emits a warning
that the regression check was **not performed** rather than reporting a pass.

### Reproducing a CI failure locally

Download the artifact and rerun the same policy:

```text
acap bench check --baseline benchmark-results/baseline --candidate benchmark-results/candidate --policy benchmarks/regression.yaml
```

Or regenerate both sides:

```text
git worktree add ../acap-base main
(cd ../acap-base && go build -o ../acap-base.exe ./cmd/acap)
go build -o acap-candidate.exe ./cmd/acap
ACAP_ROOT=$(mktemp -d) ../acap-base.exe bench suite --out results/baseline benchmarks/regression.yaml
ACAP_ROOT=$(mktemp -d) ./acap-candidate.exe bench suite --out results/candidate benchmarks/regression.yaml
./acap-candidate.exe bench check --baseline results/baseline --candidate results/candidate --policy benchmarks/regression.yaml
```

### Coding-agent benchmarks are not a CI gate

`acap bench agent` (B6/B7) depends on external LLM providers, costs money,
is rate-limited, and is nondeterministic, so it is intentionally excluded from
routine CI. It remains the way to evaluate what the deterministic suite cannot:
agent task success, real show/raw recovery behavior, and complete coding-agent
efficiency. Run it manually (for example before a release or after a
significant AgentCap change) with `--repeat`, and compare configurations with
`acap bench compare`.
