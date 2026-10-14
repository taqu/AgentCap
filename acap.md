# Task: Expand Agent Benchmark Workloads and Execute Claude Code Evaluation

## Objective

Extend the existing Agent Benchmark framework with realistic coding workloads of varying difficulty and complexity, then execute benchmark experiments using Claude Code.

The primary objective is to evaluate whether AgentCap improves coding-agent efficiency across different workload characteristics.

Evaluate:

- Total token consumption.
- Agent-visible tool-output volume.
- Task completion success.
- End-to-end execution time.
- Tool and shell command counts.
- AgentCap drill-down behavior.
- Session-aware compression effectiveness.
- Information loss or additional work caused by compression.

The benchmark must measure entire coding workflows, not isolated command-output compression ratios.

## 1. Inspect the Existing Benchmark

Before making changes:

1. Inspect the current benchmark implementation.
2. Identify existing workload definitions and test infrastructure.
3. Review the Claude Code adapter and supported execution modes.
4. Identify existing metric collection and report generation.
5. Determine how benchmark workspaces are initialized, reset, and validated.

Reuse the existing architecture wherever possible.

Do not create a parallel benchmark system.

If Claude Code support is incomplete, implement only the minimum additional functionality necessary to execute the workloads reliably.

## 2. Add Five Workload Complexity Levels

Implement representative workloads across the following complexity levels.

### Level 1 — Simple

Characteristics:

- Single-file modification.
- Small codebase.
- Minimal investigation.
- Few tool calls.
- Short execution time.

Example tasks:

- Fix a simple conditional bug.
- Correct an incorrect return value.
- Add a small unit test.
- Repair a straightforward compilation error.

Purpose:

Establish baseline agent overhead and identify whether AgentCap introduces unnecessary costs for small tasks.

### Level 2 — Moderate

Characteristics:

- Two to five files involved.
- Some repository exploration.
- Multiple build/test commands.
- Moderate diagnostic output.

Example tasks:

- Fix a bug spanning several functions.
- Repair failing unit tests.
- Implement a small feature requiring tests.
- Correct a configuration-related issue.

Purpose:

Measure compression benefits during ordinary coding workflows.

### Level 3 — Complex

Characteristics:

- Multiple modules.
- Nontrivial repository exploration.
- Several dependent changes.
- Repeated build/test iterations.
- Significant Git diff and diagnostic output.

Example tasks:

- Implement a feature across multiple modules.
- Fix a bug involving several interacting components.
- Refactor an interface while updating callers.
- Resolve multiple related test failures.

Purpose:

Evaluate how AgentCap handles repeated diagnostics, Git changes, and incremental agent workflows.

### Level 4 — High Complexity

Characteristics:

- Large repository or substantial code surface.
- Many related files.
- Significant search and investigation.
- Multiple implementation stages.
- Extensive build/test output.
- Repeated inspection of Git changes.

Example tasks:

- Refactor a subsystem without changing external behavior.
- Resolve a cross-module regression.
- Migrate an internal API across multiple packages.
- Repair a complex failure involving several components.

Purpose:

Measure the benefits of structural reduction, progressive disclosure, and session-aware deltas.

### Level 5 — Extended Workflow

Characteristics:

- Long-running, multi-stage coding task.
- Repeated investigation, implementation, validation, and correction.
- Multiple failing-to-passing transitions.
- Frequent repeated commands.
- Large cumulative tool output.
- Complex result history.

Example tasks:

- Implement a feature, discover integration failures, and repair them.
- Refactor multiple modules while maintaining test compatibility.
- Resolve several interacting regressions.
- Complete a staged implementation requiring repeated verification.

Purpose:

Stress-test AgentCap's ability to avoid retransmitting information already observed by the agent.

Level 5 should test workflow length and accumulated context usage, not merely repository size.

## 3. Workload Design Requirements

Each workload must include:

- Unique workload identifier.
- Complexity level.
- Repository or fixture identifier.
- Starting repository revision or fixture version.
- Clear task instructions.
- Expected observable behavior.
- Objective success criteria.
- Automated validation commands.
- Execution timeout.
- Estimated task characteristics.

Prefer deterministic, self-contained workloads.

Use Go, C/C++, or Rust repositories where practical because AgentCap already supports relevant build and test output.

Include workloads that naturally trigger:

- `rg`, `grep`, `find`, or similar searches.
- File inspection.
- Git status and diff operations.
- Build commands.
- Test commands.
- Repeated diagnostics.
- Repeated commands with unchanged results.

Do not artificially force the agent to execute unnecessary commands merely to improve compression statistics.

Do not include the intended solution in the agent-facing workload instructions.

## 4. Workload Coverage

Provide at least two workloads per complexity level, for a total of at least ten workloads.

Prefer different workload types within each level, such as:

| Level | Workload A | Workload B |
|---|---|---|
| L1 | Simple bug fix | Single-file test repair |
| L2 | Multi-file bug fix | Small feature |
| L3 | Cross-module feature | Interface refactoring |
| L4 | Subsystem refactoring | Complex regression |
| L5 | Multi-stage feature | Extended repair workflow |

Difficulty should increase through actual engineering complexity, investigation requirements, dependency scope, and validation effort.

Do not rely solely on artificial output volume or arbitrary delays to define complexity.

## 5. Benchmark Execution Modes

Run each workload using:

**Mode A — Baseline**

Claude Code without AgentCap.

**Mode B — AgentCap Enabled**

Claude Code with the existing AgentCap integration enabled.

If stateless and stateful modes are already available in the benchmark framework, include them as additional comparison modes where practical.

Keep all other conditions equivalent:

- Same workload.
- Same starting repository revision.
- Same Claude model.
- Same agent instructions, except integration-specific instructions.
- Same permissions and tool access.
- Same timeout and resource constraints.
- Same validation procedure.

Each run must use a fresh workspace and isolated AgentCap session.

## 6. Reproducibility and Fairness

Implement fair paired comparisons.

Requirements:

1. Reset the workspace before every run.
2. Do not reuse Claude conversations between runs.
3. Do not reuse AgentCap state between independent benchmark runs.
4. Prevent output artifacts from previous runs from affecting later runs.
5. Record model identifiers, CLI versions, and relevant configuration.
6. Record workload revisions and random seeds where applicable.
7. Avoid changing task definitions after observing results.
8. Use identical success criteria for both modes.

If resources permit, run each workload and mode multiple times to account for model variability.

Target at least three repetitions per condition for the main evaluation.

Alternate or randomize execution order where practical to reduce ordering bias.

A benchmark execution must not modify or destroy the development repository.

## 7. Metrics Collection

Collect the following metrics per run.

### Primary Metrics

- Task success.
- Total model input tokens.
- Total model output tokens.
- Total reported model token usage.
- Total agent-visible tool-output bytes.
- End-to-end wall-clock duration.

### AgentCap-Specific Metrics

- Raw command-output bytes.
- Reduced command-output bytes.
- Stateful returned bytes.
- Number of unchanged-result presentations.
- Number of delta presentations.
- Number of full fallbacks.
- Number of `acap show` calls.
- Number of `acap raw` calls.
- Bytes returned by drill-down.
- AgentCap processing overhead.

### Workflow Metrics

- Total tool calls.
- Total shell commands.
- Repeated command count.
- Build/test invocation count.
- Number of failed validation attempts.
- Number of repair iterations.

Use provider-reported usage data for model tokens.

Account explicitly for cache-related token categories if Claude Code exposes them.

Do not silently treat bytes as tokens.

If precise metrics are unavailable, report them as unavailable and explain the limitation.

## 8. Benchmark Execution

After implementing and validating the new workloads:

1. Execute the benchmark using Claude Code.
2. Run both Baseline and AgentCap Enabled modes.
3. Verify that all configured workloads are attempted.
4. Validate task success using automated checks.
5. Save detailed results for every run.
6. Generate an aggregate comparison report.

Start with a small smoke test to verify the execution infrastructure.

Then execute the workload matrix.

Do not mark a workload successful based only on Claude Code's final response.

If actual benchmark execution is blocked by authentication, permissions, budget, or missing dependencies, report the blocker precisely.

Do not fabricate benchmark results.

## 9. Analyze Results by Complexity

Aggregate results separately for L1 through L5.

For each level, calculate:

- Task success rate.
- Average and median total token usage.
- Average and median execution time.
- Average tool-output reduction.
- Average command count.
- Average AgentCap drill-down frequency.

Calculate token reduction as:

`Token Reduction (%) = 100 × (Baseline Tokens - AgentCap Tokens) / Baseline Tokens`

Calculate execution overhead as:

`Execution Overhead (%) = 100 × (AgentCap Time - Baseline Time) / Baseline Time`

Compare equivalent workloads using paired results.

Do not interpret failed or incomplete tasks as ordinary successful runs.

Report success rates separately and clearly distinguish unconditional resource usage from resource usage conditional on successful completion.

## 10. Identify AgentCap Failure Modes

Investigate cases where:

- AgentCap increases total token consumption.
- Claude Code repeatedly requests raw output.
- Compression causes additional investigation.
- Task completion fails only with AgentCap enabled.
- Delta presentations appear to confuse the agent.
- Additional tool calls eliminate initial compression gains.
- Processing overhead materially increases runtime.

For each significant regression, identify the workload, command sequence, and available supporting evidence.

Do not assume that lower tool-output volume automatically implies better benchmark performance.

## 11. Generate Benchmark Reports

Generate machine-readable results and a human-readable Markdown report using the existing report architecture.

The report should contain:

### Executive Summary

- Number of workloads executed.
- Successful and failed runs.
- Overall token reduction.
- Overall execution-time impact.
- Key conclusions.

### Results by Complexity

Compare Baseline and AgentCap Enabled for L1–L5.

### Per-Workload Results

Include task success, model token usage, output volume, execution time, and AgentCap retrieval behavior.

### Analysis

Explain:

- Which complexity levels benefit most.
- Whether benefits increase with workflow length.
- Whether stateful compression reduces repeated context.
- Whether excessive drill-down negates compression benefits.
- Whether task quality or success rate changes.

### Recommendations

Identify concrete areas for improvement, prioritized by measured impact.

Avoid adding new AgentCap functionality based solely on speculation.

## 12. Implementation Constraints

- Reuse the existing Agent Benchmark architecture.
- Preserve compatibility with existing workloads and agents.
- Keep workloads deterministic and independently reproducible.
- Avoid unnecessary third-party dependencies.
- Keep benchmark fixtures separate from production source.
- Do not weaken validation to increase success rates.
- Do not modify AgentCap's compression algorithms during the benchmark.
- Avoid unrelated refactoring.

This task is focused on workload expansion, benchmark execution, and evaluation.

## 13. Deliverables

Provide:

1. At least ten workloads covering five complexity levels.
2. Automated workload setup and validation.
3. Any necessary Claude Code benchmark integration updates.
4. Benchmark execution configuration.
5. Raw run results and execution logs.
6. Aggregate benchmark results.
7. A Markdown analysis report.
8. Documentation explaining how to reproduce the experiment.

## 14. Final Acceptance Criteria

The task is complete when:

- All five complexity levels are represented.
- At least two workloads exist per level.
- Workloads have objective automated validation.
- Baseline and AgentCap-enabled evaluations are supported.
- Benchmark runs are isolated and reproducible.
- Actual Claude Code benchmark runs have been executed.
- Results include task success and resource consumption.
- A comparative report has been generated.
- Existing benchmark functionality remains operational.

If the complete evaluation cannot be executed, clearly distinguish completed implementation work, executed runs, and remaining blockers.

## Execution Instructions

Proceed autonomously through the following phases:

**Phase 1 — Inspection:** Understand the existing benchmark infrastructure.

**Phase 2 — Workload Implementation:** Add the five complexity levels and their workloads.

**Phase 3 — Validation:** Verify workload initialization, correctness checks, and isolation.

**Phase 4 — Smoke Test:** Run a small baseline/AgentCap pair.

**Phase 5 — Full Evaluation:** Execute the benchmark matrix using Claude Code.

**Phase 6 — Analysis:** Generate results, comparisons, and recommendations.

Do not stop after implementing the workloads. Execute the benchmark and analyze the results.

Do not claim completion unless execution evidence exists.

**Primary principle: Measure whether AgentCap makes Claude Code more efficient at completing real coding tasks, not simply whether it produces smaller command outputs.**
