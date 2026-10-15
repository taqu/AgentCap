# Claude coding workload evaluation

This experiment uses the existing `acap bench agent` trial runner and
`acap bench compare` reports. No compression algorithms are modified.
The ten definitions live in `workloads/coding`; `coding-matrix.json` freezes
their complexity level, expected behavior, fixture version, estimated scope,
model, timeout, modes and three repetitions. Each fixture is a separate Go
module with only standard-library dependencies.

Run from the repository root with Python 3, Git, Go and authenticated Claude
Code installed:

```powershell
go build -o .tmp/acap-bench.exe ./cmd/acap
python benchmarks/validate_coding.py
python benchmarks/evaluate_coding.py --out benchmarks/results/smoke --smoke
python benchmarks/evaluate_coding.py --out benchmarks/results/main
python benchmarks/evaluate_coding.py --out benchmarks/results/main --append --additional-modes stateless stateful
python benchmarks/evaluate_coding.py --out benchmarks/results/main --refresh-comparisons --report-only
```

Output directories must be empty; choose a new directory for a rerun.
`--report-only --out <existing-directory>` regenerates analysis without
executing an agent. `--binary` selects an already-built AgentCap executable.
`generate_coding.py` regenerates v1 inputs; never regenerate or edit definitions
mid-experiment. The matrix manifest records fixture and definition SHA-256
hashes, executable SHA-256, repository revision, tool versions and every invocation.

The existing runner copies fixtures into fresh temporary workspaces, initializes
a deterministic Git commit, starts a fresh Claude conversation and isolated
AgentCap store/session, validates directly after the agent exits, then removes
the disposable workspace. Development source is never the trial workspace.
Validation checks original tests and go.mod against the initial Git commit,
runs uncached tests and runs vet. The reference-validation script separately
proves initial failure and reference success; its repairs are never copied to
agent workspaces. A final agent message is not evidence of success.

Claude modes use identical permission bypass, model, built-in tool availability,
project settings scope and MCP isolation. Integrated mode installs the existing
AgentCap hook plus recovery instructions. Disabled mode installs neither.
External user settings are excluded; Claude still uses the user's configured
authentication. Runs are sequential. The Baseline/Integrated pair order
alternates by workload and repetition. Optional Stateless/Stateful trials reuse
the same frozen baseline result for the matching workload and repetition; each
enabled-mode trial still gets a fresh workspace, conversation and state store.
The optional modes append their own 60 conditions. The smoke pair is separate
from the main results. It may incur provider charges or subscription quota.

Each result JSON is the canonical existing benchmark result. Each `.log` retains
verbose JSONL stdout and stderr. Paired comparison JSON comes from the existing
comparison command. `aggregate.json` and `report.md` extend that reporting with
complexity groups, means, medians, paired token reduction and time overhead.
Unconditional and successful-only summaries remain separate. Both members must
succeed for a pair to support efficiency conclusions.

`wall_time_ns` measures complete agent execution, including configuration, CLI
startup and provider/tool work. It excludes benchmark-only setup and validation.
Manifest `elapsed` covers the full harness invocation including those steps.

Token totals explicitly sum provider input, output, cache-creation input and
cache-read input across `modelUsage` (including auxiliary models), falling back
to final `usage` when per-model usage is absent. The two sources are not added
together. `usage_available` distinguishes missing/error usage from zero;
older consumers must not infer availability from absent numeric fields.
`tool_visible_bytes` measures UTF-8 text from all native tool-result messages;
`total_visible_bytes` measures the existing shell-output/retrieval boundary.
These different byte boundaries must not be mixed in a comparison.
Nontext tool payloads and hidden prompts are excluded. Shell counts represent
Bash tool invocations, including compound scripts. Tool counts represent native
tool invocations. AgentCap raw/reduced/stateful bytes, retrievals, processing
time and unchanged/delta/fallback counts come from isolated store-stat deltas.
Repeated commands count identical Bash command strings. Build/test/vet calls
and failed intermediate validations count matching Bash calls and native error
tool results. Repair iterations remain unavailable rather than inferred from
command counts. Per-run command sequences are retained in aggregate.json.

The existing Claude adapter conservatively bypasses commands mentioning `cd`,
other shell state mutations, background calls and retrievals. This behavior is
part of the measured integration. `intercepted_commands` is separate from all
Bash invocations. AgentCap raw/reduced/stateful measurements cover intercepted
commands only; all-tool bytes also cover bypassed output. Do not call zero
intercepted output a 100% workflow compression gain.

L4 concurrency fixtures include a concurrent test but the default verifier does
not use the race detector, which requires additional platform toolchains.
L5 requires meaningful staged layer validation; actual command repetition and
workflow length are agent decisions. Scope levels need empirical calibration
before generalizing beyond these fixtures. Do not infer a benefit from reduced
bytes alone, from a failed task, or from the smoke pair alone.
