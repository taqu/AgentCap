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
stdout. Schema version 2 adds progressive-disclosure recovery measurements:

```json
{
  "schema_version": 2,
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
