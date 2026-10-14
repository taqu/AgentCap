# AgentCap Phase 1 — Core Command Output Compression

## Project

Project name: AgentCap  
CLI command: `acap`  
Language: Go

AgentCap is a command-output reduction layer designed for AI coding agents.

The long-term goal is to intercept command/tool execution, retain the useful information from command results, and significantly reduce the number of tokens returned to the coding agent.

Phase 1 should focus on a small, robust foundation rather than broad command coverage.

The implementation must remain simple, deterministic, fast, and safe.

---

# 1. Phase 1 Goal

Implement the first usable version of AgentCap with:

- command execution through `acap run`
- command classification
- command-specific output reducers
- generic fallback reduction
- token/size reduction statistics
- support for:
  - `ls`
  - `find`
  - `grep`
  - `rg`
  - `cat`
  - `head`
  - `tail`
  - `tree`
  - `du`
  - `wc`

The primary goal is:

> Reduce command output sent to an AI coding agent while preserving the information most likely to be useful for the agent's next action.

Phase 1 is NOT intended to provide semantic code analysis, persistent result drill-down, session-level deduplication, Git-aware compression, or build/test parsing.

Those belong to later phases.

---

# 2. Core Design Principle

AgentCap must behave as a transparent command execution layer.

Conceptually:

```text
command argv
    |
    v
executor
    |
    +--> stdout
    +--> stderr
    +--> exit code
    |
    v
command classifier
    |
    v
reducer
    |
    v
compact agent-facing output
```

Running a command through AgentCap must not silently change the meaning of the command.

For example:

```bash
acap run rg "Workspace" .
```

should execute the real `rg` binary with the supplied arguments.

AgentCap may transform the output returned to the caller, but it must preserve important execution metadata such as the exit code.

---

# 3. CLI

Implement at least:

```bash
acap run <command> [args...]
acap stats
```

Examples:

```bash
acap run ls -la
acap run find . -type f
acap run rg "Workspace" .
acap run cat src/main.go
```

Do not implement shell parsing inside `acap run`.

The CLI should treat everything after `run` as an argv array and execute it directly.

For example:

```bash
acap run rg "foo bar" src
```

must result in an argv equivalent to:

```text
["rg", "foo bar", "src"]
```

Avoid executing the command via:

```bash
sh -c
```

or:

```bash
bash -c
```

unless absolutely necessary.

Direct process execution is preferred.

---

# 4. Execution Layer

Create a reusable execution abstraction.

A reasonable internal representation is:

```go
type ExecutionResult struct {
    Command   []string
    ExitCode  int
    Stdout    []byte
    Stderr    []byte
    Duration  time.Duration
}
```

This exact type is not mandatory.

The design should, however, clearly preserve:

- argv
- stdout
- stderr
- exit code
- execution duration

Also handle:

- command not found
- interrupted processes
- non-zero exit codes
- empty stdout
- empty stderr
- very large output

Do not treat a non-zero exit code as an internal AgentCap failure.

For example:

```bash
rg nonexistent-pattern .
```

may legitimately return a non-zero exit code.

AgentCap should preserve that behavior.

---

# 5. Exit Code Preservation

This is a strict requirement.

If the wrapped command exits with code N, then:

```bash
acap run ...
```

should normally exit with the same code N.

Compression must never turn a failed command into an apparently successful one.

Internal AgentCap failures should use clearly distinguishable failure behavior.

---

# 6. Reducer Architecture

Do not implement all reduction logic in one large function.

Use a reducer abstraction.

For example:

```go
type Reducer interface {
    Match(cmd Command) bool
    Reduce(ctx context.Context, result *ExecutionResult) (*ReducedResult, error)
}
```

The exact API may differ.

The important architectural property is:

```text
known command
    -> specialized reducer

unknown command
    -> generic reducer
```

Possible implementation layout:

```text
cmd/
  acap/
    main.go

internal/
  exec/
  classify/
  reduce/
    reducer.go
    generic.go
    ls.go
    find.go
    grep.go
    cat.go
    tree.go
    du.go
    wc.go
  render/
  stats/
```

Use idiomatic Go package boundaries.

Do not over-engineer the abstraction.

---

# 7. Output Format

Phase 1 output should be optimized for an LLM rather than for decorative terminal presentation.

Prefer compact plain text.

Avoid verbose JSON unless explicitly requested in the future.

A typical result can look like:

```text
@acap ls entries=1042 files=891 dirs=151

dirs:
src/
tests/
docs/
vendor/

largest:
build.log 18.2MB
data.bin 7.3MB

omitted=1018
```

Keep labels short and predictable.

Avoid:

- ASCII boxes
- Unicode decoration
- unnecessary explanatory prose
- repeating information already obvious from context

The format does not need to become a stable public protocol in Phase 1.

However, keep it deterministic enough for automated testing.

---

# 8. Generic Reducer

Every unsupported command must still work.

Implement a generic reducer.

The generic reducer should be conservative.

Useful generic operations include:

- strip ANSI escape sequences
- remove or collapse terminal progress updates
- collapse long runs of identical lines
- detect extremely large output
- keep a useful prefix
- keep a useful suffix
- report how much output was omitted

Example:

```text
@acap generic lines=8421 shown=120 omitted=8301

[first lines...]

...

[last lines...]
```

Do NOT aggressively summarize arbitrary text using heuristics that may destroy errors or diagnostics.

For unknown commands, preserving information is more important than maximizing compression.

---

# 9. `ls` Reducer

Support common `ls` output.

For large directory listings, avoid returning every entry.

Extract useful information such as:

- total entry count
- file count where practical
- directory count where practical
- visible top-level directories
- optionally largest files if size information is available
- omitted count

Example:

```text
@acap ls entries=1042 files=891 dirs=151

dirs:
src/
tests/
docs/
vendor/

largest:
build.log 18.2MB
data.bin 7.3MB

omitted=1018
```

For small output, do not compress unnecessarily.

A short `ls` result should usually be returned almost unchanged.

---

# 10. `find` Reducer

`find` can easily produce thousands of paths.

For large path lists, summarize them structurally.

Useful aggregation:

- total paths
- probable file count
- probable directory count
- counts grouped by top-level directory
- counts grouped by extension
- a small representative sample

Example:

```text
@acap find paths=5182

by-root:
vendor 2811
tests 1256
src 824
docs 291

extensions:
.go 381
.cpp 179
.h 204

sample:
src/main.go
src/workspace.go
tests/workspace_test.go

omitted=5179
```

Do not assume all `find` output contains filesystem paths.

If the output cannot safely be interpreted as a path list, fall back to generic reduction.

---

# 11. `grep` / `rg` Reducer

Treat `grep` and `rg` as particularly important.

Coding agents frequently generate large search outputs.

For normal file-and-line search results, group matches by file.

Example input:

```text
src/workspace.go:41:type Workspace struct {
src/workspace.go:89:func NewWorkspace(...)
src/workspace.go:123:func ...
tests/workspace_test.go:18:...
...
```

Possible reduced output:

```text
@acap rg matches=317 files=42

files:
src/workspace.go 81
tests/workspace_test.go 92
src/resolver.go 47
...

sample:
src/workspace.go:41:type Workspace struct {
src/workspace.go:89:func NewWorkspace(...)
tests/workspace_test.go:18:...

omitted=...
```

Important requirements:

- preserve filenames
- preserve line numbers when present
- preserve representative matching lines
- do not remove all context from an error-like match
- keep output ordering deterministic

Small search results should remain mostly untouched.

---

# 12. `cat` Reducer

Do not implement AST parsing in Phase 1.

For small files, return the content unchanged.

For large files, use simple structural reduction.

Example:

```text
@acap cat file=src/workspace.go lines=1842 bytes=61243

[first 80 lines]

...

[last 80 lines]

omitted_lines=1682
```

If the file appears to be binary, do not dump binary content into the model context.

Return compact metadata instead.

Example:

```text
@acap cat binary bytes=482193
output omitted
```

Do not attempt semantic source-code summarization yet.

---

# 13. `head` and `tail`

These commands already constrain output.

Avoid applying aggressive additional compression.

Normally:

- strip irrelevant ANSI/progress noise if present
- otherwise preserve output

AgentCap should not make a naturally small command less useful merely to claim a higher compression ratio.

---

# 14. `tree` Reducer

Large tree output should be collapsed by depth and/or subtree size.

Preserve important top-level structure.

Example:

```text
@acap tree dirs=316 files=5182

.
├── cmd/
│   └── acap/
├── internal/
│   ├── exec/
│   ├── reduce/
│   └── stats/
├── tests/
└── ...

depth_truncated=true
```

Exact visual formatting is optional.

Prefer simple output if parsing tree's human-readable formatting becomes brittle.

---

# 15. `du` Reducer

Large `du` output should focus on size-heavy paths.

Useful result:

```text
@acap du entries=1834 total=4.8GB

largest:
2.9GB vendor/
1.1GB build/
420MB .git/
180MB testdata/

omitted=1830
```

Retain enough information to identify disk-heavy areas.

---

# 16. `wc` Reducer

`wc` output is usually already compact.

Preserve it nearly verbatim.

Only normalize obvious formatting noise if necessary.

Do not reduce useful output simply because a reducer exists.

---

# 17. Compression Thresholds

Do not compress every output.

Introduce configurable or internal thresholds.

For example:

```text
small output
    -> return essentially unchanged

medium output
    -> light reduction

large output
    -> command-specific summarization
```

The exact thresholds may be adjusted based on tests.

Prefer byte/line-based thresholds in Phase 1.

Do not introduce a tokenizer dependency merely to decide whether output is large.

---

# 18. Statistics

Implement basic reduction statistics.

At minimum measure:

```text
commands processed
raw bytes
returned bytes
bytes saved
reduction percentage
```

If a lightweight token estimator is straightforward, also expose estimated token counts.

Do not make token estimation a critical dependency.

`acap stats` could produce something like:

```text
commands: 142
raw: 8.4MB
returned: 1.7MB
saved: 6.7MB
reduction: 79.8%
```

Statistics should preferably survive multiple command executions.

A small local stats file is acceptable.

Keep persistence simple.

---

# 19. Performance Requirements

AgentCap exists to help interactive coding agents.

It must not introduce noticeable latency for ordinary commands.

Avoid:

- spawning unnecessary subprocesses
- parsing output multiple times
- reading the same data repeatedly
- complex regex chains over huge buffers where streaming is practical

Phase 1 does not need sophisticated streaming compression, but the architecture should not make it impossible later.

---

# 20. Memory Safety for Large Output

Do not assume command output is always small.

Prefer a design that can eventually spill large output to disk.

If Phase 1 initially stores stdout/stderr in memory, introduce a clearly defined maximum and a safe fallback.

AgentCap itself should not consume hundreds of megabytes because a child process emits a huge build log.

Do not silently truncate without indicating that truncation occurred.

---

# 21. stderr Handling

Do not casually discard stderr.

Errors are often the most important part of command output for a coding agent.

Reducers should preserve stderr conservatively.

If stdout is large but stderr contains only a few useful lines, those stderr lines should remain visible.

The renderer may distinguish them compactly if useful.

---

# 22. ANSI and Progress Noise

Remove terminal-only noise that has little value to an LLM.

Examples:

- ANSI color sequences
- carriage-return progress updates
- repeated progress percentages
- spinner frames

Do not remove text merely because it came from stderr.

---

# 23. Safety / Transparency

AgentCap must not:

- modify files unless the wrapped command itself does so
- rewrite the command being executed
- silently retry commands
- invent command output
- infer nonexistent filesystem state
- turn errors into success
- suppress every diagnostic merely to reduce output

Compression must remain deterministic and inspectable.

---

# 24. Testing

Add unit tests for every reducer.

Tests should cover:

- small output that should remain unchanged
- large output
- malformed/unexpected output
- empty output
- Unicode filenames/text
- ANSI sequences
- non-zero command exits
- stdout and stderr together

Also add integration tests for `acap run`.

At minimum verify:

```text
acap run printf ...
acap run command-that-exits-nonzero
acap run rg ...
```

when the relevant external tools are available.

Do not make the entire test suite depend on optional third-party commands.

---

# 25. Benchmark Fixtures

Create deterministic fixtures representing large outputs.

Examples:

```text
testdata/
  ls-large.txt
  find-large.txt
  rg-large.txt
  cat-large.txt
  tree-large.txt
  du-large.txt
```

Use these fixtures for both correctness tests and benchmarks.

Benchmark:

- reduction latency
- allocations where useful
- raw size
- reduced size

Do not optimize prematurely, but make regressions measurable.

---

# 26. Observability

Add an opt-in debug mode.

For example:

```bash
ACAP_DEBUG=1 acap run ...
```

or:

```bash
acap --debug run ...
```

Debug output should explain:

- detected command
- selected reducer
- raw size
- reduced size
- execution duration
- reduction duration

Debug information must go to stderr and must not contaminate normal compact output.

---

# 27. Configuration

Keep Phase 1 configuration minimal.

Do not build a large configuration framework.

It is acceptable to hard-code reasonable defaults initially.

If configuration is necessary, limit it to parameters such as:

```text
compression threshold
maximum preview lines
maximum output size
stats enabled/disabled
```

Do not add YAML/TOML configuration unless the implementation clearly requires it.

---

# 28. Explicit Non-Goals

Do NOT implement the following in Phase 1:

- MCP server
- Claude Code integration
- Codex-specific integration
- shell hooks
- command alias rewriting
- persistent command-result IDs
- `acap show`
- `acap raw`
- result drill-down
- cross-command deduplication
- session-aware state
- Git semantic parsing
- compiler-specific diagnostics
- test framework parsing
- AST parsing
- Tree-sitter
- source-code symbol extraction
- embeddings
- LLM-based summarization
- remote services
- telemetry
- cloud storage
- plugin ecosystem

Do not introduce infrastructure for these unless it is naturally required by the Phase 1 architecture.

Keep the codebase small.

---

# 29. Suggested Implementation Order

Implement in this order:

1. CLI skeleton
2. process execution and exit-code propagation
3. `ExecutionResult`
4. ANSI/progress cleanup utilities
5. reducer interface
6. generic reducer
7. rendering
8. statistics
9. `rg` / `grep` reducer
10. `find` reducer
11. `ls` reducer
12. `cat` reducer
13. `head` / `tail`
14. `tree`
15. `du`
16. `wc`
17. integration tests
18. benchmarks
19. documentation

Do not begin with all reducers simultaneously.

Make the execution and reducer framework reliable first.

---

# 30. Definition of Done

Phase 1 is complete when all of the following are true:

1. `acap` builds as a single Go binary.

2. The following works reliably:

```bash
acap run <command> [args...]
```

3. The child command's exit status is preserved.

4. Unsupported commands still run through the generic reducer.

5. Specialized reducers exist for:

```text
ls
find
grep
rg
cat
head
tail
tree
du
wc
```

6. Small outputs are not unnecessarily compressed.

7. Large outputs are substantially reduced while preserving actionable information.

8. stderr and failure diagnostics remain visible.

9. ANSI/progress noise is removed where appropriate.

10. `acap stats` reports at least raw bytes, returned bytes, and reduction percentage.

11. Each reducer has meaningful tests.

12. Large-output fixtures and basic benchmarks exist.

13. No MCP, AST, Git-specific, session-level, or LLM-based functionality has been added.

---

# 31. Engineering Priorities

When tradeoffs are necessary, use this priority order:

1. Correctness
2. Preserve actionable information
3. Preserve command semantics and exit status
4. Conservative behavior
5. Low latency
6. Token/output reduction
7. Breadth of command support

A reducer that saves fewer tokens but reliably preserves useful diagnostics is better than an aggressive reducer that sometimes hides information the coding agent needs.

---

# 32. Final Deliverable

At the end of the implementation:

- run the full Go test suite
- run relevant benchmarks
- verify representative commands manually
- report the final project structure
- summarize implemented reducers
- report known limitations
- report benchmarked reduction ratios for the included fixtures

Do not proceed into Phase 2 work.

Stop once the Phase 1 Definition of Done is satisfied.
