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
```

Supported steps are:

- `run`: structured `argv`, optional workspace-relative `cwd`, and optional
  `expect.exit`. It invokes no shell and executes exactly once.
- `copy`: copy a file or directory from the fixture tree to the workspace.
- `write`: replace or create a file whose parent already exists.
- `remove`: remove a workspace child.
- `mkdir`: create a workspace directory and missing parents.

Only `run` is measured. A target command's non-zero status is ordinary workload
data unless `expect.exit` is present and does not match. A missing executable,
unsafe path, failed mutation, or failed assertion is a workload infrastructure
failure identifying the step number and type.

When `git.init` is true, setup initializes and commits the copied fixture before
the first step. It uses repository-local identity, disables signing and hooks,
uses fixed author/committer dates, ignores system/global Git configuration, and
does not use a network remote. Git setup is harness work, not a measured step.

The bundled workloads are:

- `git/repeated-diff`: changed diff, changed diff, then unchanged diff.
- `build/compile-fix`: two compile failures with fewer errors, then success.
- `test/fail-fix-pass`: failing test, passing test, then a repeated pass.

They use Git and the Go standard toolchain only and require no network access.

## Stable result JSON

`--json` writes one JSON object and no human headers or command output to
stdout. Schema version 1 contains these stable fields:

```json
{
  "schema_version": 1,
  "workload": "git/repeated-diff",
  "commands": 3,
  "raw_bytes": 716,
  "stateless_bytes": 285,
  "stateful_bytes": 203,
  "show_bytes": 0,
  "raw_retrieval_bytes": 0,
  "processing_ns": 151000000
}
```

Counts, byte measurements, and duration are numeric. `processing_ns` preserves
the existing B2 AgentCap-processing definition. `show_bytes` and
`raw_retrieval_bytes` are zero because recovery operations are not measured
until a later benchmark phase. The schema version is independent of workload
definition version 1 and the internal store schema. Consumers should select
their interpretation using `schema_version`; later schema versions may add or
change fields.
