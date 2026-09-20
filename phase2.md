# AgentCap Phase 2 — Result Store, Capsules, and Drill-Down

## Project

Project name: AgentCap  
CLI command: `acap`  
Language: Go

AgentCap is a command-output reduction layer for AI coding agents.

Phase 0 established reliable command execution.

Phase 1 added command-aware output reduction.

Phase 2 must introduce persistent command results and progressive disclosure.

The core concept of this phase is:

> Capture the complete command result locally, return a compact capsule to the agent, and allow the agent to retrieve more detail on demand without re-running the command.

This phase is what turns AgentCap from a simple output compressor into a recoverable command-result layer.

---

# 1. Phase 2 Goal

Implement persistent result capture and drill-down.

Phase 2 must provide:

- unique result IDs
- persistent local result storage
- complete raw stdout storage
- complete raw stderr storage
- command metadata storage
- reduced capsule storage
- `acap show <id>`
- `acap raw <id>`
- basic drill-down selectors
- result inspection without re-running commands
- cleanup / retention basics
- tests for persistence and retrieval

The primary objective is:

> Compression must become reversible at the command-result level.

The default agent-facing result may be compact, but the full original command output must remain locally available.

---

# 2. Core Model

The Phase 2 execution flow should become:

```text
command
   |
   v
executor
   |
   +------> raw stdout
   |
   +------> raw stderr
   |
   +------> metadata
   |
   v
reducer
   |
   v
capsule
   |
   +------> returned to agent
   |
   v
persistent result store
```

Later:

```text
acap show <id>
        |
        v
stored structured result
        |
        v
selected detail
```

Or:

```text
acap raw <id>
        |
        v
original stored output
```

Do not require command re-execution for drill-down.

---

# 3. Result Identity

Every captured command result must receive a stable result ID.

Example:

```text
@acap 8f31c2
```

The ID should be:

- short enough for an AI agent to reuse easily
- collision-resistant enough for local usage
- deterministic only if useful, but determinism is not required
- safe for use as a directory or file name

A reasonable format is:

```text
8f31c2
```

or:

```text
r_8f31c2
```

Do not use long UUIDs in normal output unless necessary.

The user-facing identifier should remain compact.

Internally, a longer identifier may be stored if desired.

---

# 4. Result Store

Introduce a persistent result store.

A reasonable default location on Unix-like systems is:

```text
~/.cache/agentcap/
```

Prefer using the platform-appropriate user cache directory via Go APIs rather than hard-coding the path.

A possible layout is:

```text
~/.cache/agentcap/
└── results/
    └── 8f31c2/
        ├── meta.json
        ├── stdout
        ├── stderr
        └── capsule
```

The exact format may differ.

The result store must preserve enough data to retrieve the original result later.

---

# 5. Stored Metadata

Each result should contain metadata similar to:

```json
{
  "id": "8f31c2",
  "command": ["rg", "Workspace", "."],
  "exit_code": 0,
  "started_at": "...",
  "duration_ms": 87,
  "stdout_bytes": 48291,
  "stderr_bytes": 0,
  "reducer": "rg",
  "created_at": "..."
}
```

This exact schema is not mandatory.

At minimum preserve:

- result ID
- argv
- exit code
- execution duration
- creation time
- stdout size
- stderr size
- selected reducer
- raw result locations or equivalent references

Keep the metadata schema simple and versionable.

---

# 6. Capture Architecture

Phase 1 may have streamed or buffered output differently.

Phase 2 must capture complete output without introducing unbounded memory usage.

Do NOT implement:

```text
entire stdout -> RAM
entire stderr -> RAM
```

as the only strategy for arbitrarily large output.

Prefer streaming capture.

For example:

```text
child stdout
   |
   +--> result store
   |
   +--> reducer input
```

and similarly for stderr.

If reducers currently require complete buffers, it is acceptable to use temporary files as the canonical captured representation.

The design should allow large command output to be stored safely.

---

# 7. Preserve Full Raw Output

This is a strict requirement.

If AgentCap returns:

```text
@acap 8f31c2 matches=431 files=52
...
```

the complete original output must remain recoverable.

AgentCap must not store only the reduced result.

The following must be possible:

```bash
acap raw 8f31c2
```

and should reproduce the original captured output as faithfully as practical.

---

# 8. `acap raw`

Implement:

```bash
acap raw <result-id>
```

Default behavior should output the complete stored command output.

Consider preserving stdout and stderr semantics.

Possible forms:

```bash
acap raw 8f31c2
acap raw 8f31c2 --stdout
acap raw 8f31c2 --stderr
```

At minimum:

```bash
acap raw <id>
```

must work.

If stdout and stderr were captured separately, do not invent a false original interleaving.

A simple documented default is acceptable.

For example:

```text
stdout
followed by stderr
```

or output stdout by default with explicit `--stderr`.

Prefer explicit semantics over pretending exact cross-stream ordering was preserved.

---

# 9. `acap show`

Implement:

```bash
acap show <result-id>
```

`show` must operate on the stored result.

It must NOT re-run the command.

Default behavior should display the stored capsule or a useful expanded version of it.

Example:

```bash
$ acap show 8f31c2

@acap 8f31c2
command: rg Workspace .
exit=0
matches=317 files=42

files:
src/workspace.go 81
tests/workspace_test.go 92
src/resolver.go 47
...
```

Keep the output optimized for agent consumption.

---

# 10. Progressive Disclosure

Introduce progressive disclosure.

The model should not need to jump directly from:

```text
compact capsule
```

to:

```text
entire raw output
```

Provide useful intermediate retrieval.

For example:

```text
Level 0
summary

Level 1
grouped details

Level 2
selected file/path/group

Level 3
raw
```

Do not over-formalize these levels in the public API yet.

The key requirement is that useful intermediate detail retrieval exists.

---

# 11. Drill-Down Selectors

Implement a small set of generic selectors.

Recommended initial options:

```bash
acap show <id> --match <text>
acap show <id> --from <line>
acap show <id> --to <line>
acap show <id> --lines <start>:<end>
```

Where applicable, also support:

```bash
acap show <id> --path <path>
```

Do not implement a complex query language.

Keep selectors simple and composable.

---

# 12. Line Range Retrieval

Line-based retrieval is especially important.

Example:

```bash
acap show 8f31c2 --lines 100:160
```

This should retrieve the relevant range from the stored raw output.

Avoid reading the entire file into memory merely to return a small range if efficient streaming is straightforward.

Line numbers should be clearly defined.

Use 1-based line numbering unless there is a compelling reason otherwise.

Document the behavior.

---

# 13. Text Matching

Implement:

```bash
acap show <id> --match Workspace
```

This should search the stored output.

Useful behavior:

```text
@acap 8f31c2 match="Workspace" hits=23

41: type Workspace struct {
89: func NewWorkspace(...)
123: ...
```

Avoid dumping thousands of matches.

Use a reasonable default limit.

Report omitted match count where applicable.

---

# 14. Path-Based Drill-Down

Where the stored command output contains recognizable path-oriented structure, allow:

```bash
acap show <id> --path src/workspace.go
```

This is particularly useful for:

- `rg`
- `grep`
- `find`
- `tree`

Do not implement language semantics.

Path filtering should be based on Phase 1 reducer/index information or conservative textual matching.

If path-specific drill-down cannot be reliably supported for a result type, return a clear error rather than guessing.

---

# 15. Capsule Format

Each result should store the exact reduced capsule produced for the original invocation.

Example:

```text
@acap 8f31c2 rg matches=317 files=42

files:
src/workspace.go 81
tests/workspace_test.go 92
...

sample:
...
```

The result ID must appear prominently.

The agent should be able to copy it directly into:

```bash
acap show 8f31c2
```

or:

```bash
acap raw 8f31c2
```

---

# 16. Result ID in Normal `acap run` Output

Every stored result should expose its ID.

Example:

```bash
$ acap run rg "Workspace" .

@acap 8f31c2 rg matches=317 files=42
...
```

Even when the original output is small and reduction is minimal, a stored result may still receive an ID.

Prefer consistency.

If storage is disabled or unavailable, fail clearly or use a documented fallback.

Do not silently claim a result is recoverable if it was not stored.

---

# 17. Store Before Returning

Prefer this order:

```text
execute
capture
persist
reduce
store capsule
return capsule
```

or another ordering that guarantees that a result ID returned to the agent refers to a valid stored result.

Do not emit:

```text
@acap abc123
```

and then fail to persist `abc123`.

Atomic or near-atomic result creation is desirable.

---

# 18. Atomic Result Creation

Avoid partially valid stored results.

A reasonable approach:

```text
results/.tmp-XYZ/
   meta
   stdout
   stderr
   capsule

rename -> results/8f31c2/
```

Use atomic rename where supported.

If execution fails normally with a non-zero child exit code, that is still a valid result and should be stored.

Only internal AgentCap storage failures should prevent a result from being committed.

---

# 19. Non-Zero Exit Results

Failed commands are often the most useful command results.

For example:

```bash
acap run go test ./...
```

returning exit code `1` must still produce a stored result.

The result should contain:

- stdout
- stderr
- exit code
- capsule
- metadata

Then the agent must be able to run:

```bash
acap show <id>
acap raw <id>
```

even though the original command failed.

---

# 20. Storage Failure

Define clear behavior when the result store cannot be written.

Examples:

- disk full
- permission denied
- corrupt cache directory

Do not silently discard the full result while returning a capsule that claims drill-down is available.

Prefer:

- clear stderr diagnostic
- non-zero AgentCap internal failure if reliable storage is required
- or an explicitly documented degraded mode

Choose one behavior and test it.

Reliability is more important than pretending success.

---

# 21. Result Metadata Inspection

Optionally support:

```bash
acap show <id> --meta
```

Example:

```text
id=8f31c2
command=rg Workspace .
exit=0
duration=87ms
stdout=48291B
stderr=0B
reducer=rg
created=...
```

Keep this compact.

Do not expose internal implementation paths unless debug mode is enabled.

---

# 22. Internal Result Model

Introduce a clear result abstraction.

For example:

```go
type StoredResult struct {
    ID        string
    Command   []string
    ExitCode  int
    Duration  time.Duration
    CreatedAt time.Time

    Reducer string

    StdoutPath  string
    StderrPath  string
    CapsulePath string
}
```

The exact representation is flexible.

Avoid loading complete raw data into this struct.

Store references to data rather than duplicating large byte slices.

---

# 23. Storage Interface

Introduce a small storage interface.

For example:

```go
type Store interface {
    Create(ctx context.Context, meta Metadata) (*Writer, error)
    Open(ctx context.Context, id string) (*StoredResult, error)
}
```

Or something simpler.

The purpose is:

- keep filesystem persistence isolated
- make tests easier
- avoid coupling CLI logic directly to cache directory layout

Do not build a generic database abstraction.

Filesystem storage is sufficient.

---

# 24. Reducer Integration

Phase 1 reducers should continue to operate.

Do not rewrite reducer behavior unnecessarily.

Instead, integrate them into the new result pipeline.

Conceptually:

```text
captured result
      |
      +--> persistent raw storage
      |
      v
reducer
      |
      v
capsule
      |
      +--> persistent capsule
      |
      +--> stdout to agent
```

If a reducer fails internally, consider falling back to the generic reducer.

Do not lose the raw result.

---

# 25. Optional Result Indexes

Reducers may generate lightweight indexes to improve drill-down.

For example, the `rg` reducer could record:

```text
file -> raw line ranges
```

The `find` reducer could record:

```text
top-level root -> raw line ranges
```

These indexes are useful but should remain lightweight.

A possible stored file:

```text
index.json
```

Do not build a general semantic indexing framework.

Do not introduce SQLite unless there is a demonstrated need.

Plain files are sufficient for Phase 2.

---

# 26. Generic Results

Drill-down must also work for unsupported commands handled by the generic reducer.

At minimum, generic results should support:

```bash
acap show <id>
acap show <id> --match <text>
acap show <id> --lines X:Y
acap raw <id>
```

This guarantees that progressive disclosure is not limited only to known command types.

---

# 27. Result Retrieval Must Be Cheap

`acap show` should be substantially cheaper than re-running the command.

Do not invoke:

- the original command
- external parsing tools
- an LLM
- remote services

during result retrieval.

The result should be derived entirely from local stored data.

---

# 28. Cleanup

Introduce basic cache cleanup.

At minimum implement one of:

```bash
acap clean
```

or:

```bash
acap cache clean
```

A simple implementation that removes stored results is sufficient.

Optionally support:

```bash
acap clean --older-than 7d
```

but do not spend excessive effort on retention policy in Phase 2.

---

# 29. Automatic Retention

Avoid unbounded cache growth.

Implement a simple retention strategy.

Possible approach:

```text
max cache size
or
max result count
or
max age
```

Use conservative defaults.

For example:

```text
max age: 7 days
```

or a bounded total cache size.

The exact policy is less important than ensuring the cache does not grow forever.

Make cleanup behavior predictable.

---

# 30. Do Not Delete Active Results

Cleanup must not remove a result while it is currently being written.

Use temporary directories or lock-free atomic creation patterns so cleanup can distinguish complete results from incomplete ones.

Avoid complex locking if atomic filesystem operations are sufficient.

---

# 31. Cache Corruption

Handle malformed stored results gracefully.

For example:

```bash
acap show deadbe
```

where metadata exists but stdout is missing.

Expected behavior:

- no panic
- clear error
- non-zero exit
- identify the damaged result if useful

Do not attempt to reconstruct missing data by re-running the command.

---

# 32. Missing Result

For:

```bash
acap show doesnotexist
```

or:

```bash
acap raw doesnotexist
```

return:

- concise stderr message
- non-zero exit code

Do not search unrelated directories or attempt fuzzy ID guessing unless explicitly implemented later.

Exact ID lookup is sufficient.

---

# 33. ID Prefix Matching

Optionally allow unique prefixes.

For example, if the stored ID is:

```text
8f31c2
```

then:

```bash
acap show 8f3
```

may work if that prefix is unique.

If multiple results match, return an ambiguity error.

This is optional.

Do not prioritize it over correctness.

---

# 34. Token Reduction Accounting

Phase 1 already tracks reduction statistics.

Extend statistics so retrieval can be measured separately.

Useful counters:

```text
run raw bytes
run returned bytes
show calls
raw calls
show returned bytes
raw returned bytes
```

This will later allow calculation of:

> total session bytes/tokens consumed after drill-down

This metric is more important than initial compression ratio alone.

Do not build full session analytics yet.

---

# 35. Important Future Metric

Design statistics so the following can eventually be measured:

```text
initial result bytes
+
all drill-down bytes
=
total bytes exposed to agent
```

Compared with:

```text
full raw result bytes
```

This enables evaluation of actual progressive-disclosure efficiency.

Phase 2 does not need full cross-command session tracking yet.

---

# 36. Output Examples

## Initial command

```bash
$ acap run rg "Workspace" .

@acap 8f31c2 rg matches=317 files=42

files:
src/workspace.go 81
tests/workspace_test.go 92
src/resolver.go 47

sample:
src/workspace.go:41:type Workspace struct {
src/workspace.go:89:func NewWorkspace(...)
tests/workspace_test.go:18:...

omitted=311
```

## Show

```bash
$ acap show 8f31c2

@acap 8f31c2 rg matches=317 files=42
...
```

## Match

```bash
$ acap show 8f31c2 --match Resolve

@acap 8f31c2 match="Resolve" hits=14

src/workspace.go:421:func (w *Workspace) Resolve(...)
...
```

## Path

```bash
$ acap show 8f31c2 --path src/workspace.go

@acap 8f31c2 path=src/workspace.go matches=81
...
```

## Raw

```bash
$ acap raw 8f31c2
```

returns the complete stored result.

---

# 37. Performance Requirements

Result persistence should add low overhead.

Avoid:

- unnecessary copies of complete stdout
- multiple complete reads of large files
- repeated serialization of raw output
- loading raw output into memory during `show` unless required

Use buffered I/O where appropriate.

Measure storage and retrieval overhead in benchmarks.

---

# 38. Storage Security

Stored output may contain source code, environment-derived output, paths, compiler diagnostics, or secrets printed by child commands.

Therefore:

- store results only locally
- use user-private file permissions where practical
- do not upload anything
- do not introduce telemetry
- do not expose cache files to other users by default

On Unix-like systems, prefer directories/files that are not world-readable.

---

# 39. No Remote Components

Phase 2 must remain entirely local.

Do NOT add:

- remote APIs
- cloud storage
- synchronization
- remote databases
- telemetry
- external analytics

AgentCap must remain a single local Go executable.

---

# 40. Testing

Add tests for the full result lifecycle.

At minimum test:

## Result creation

```text
run command
-> result ID created
-> metadata exists
-> stdout exists
-> stderr exists
-> capsule exists
```

## Successful retrieval

```text
acap show <id>
acap raw <id>
```

## Non-zero child exit

Failed child commands must still produce retrievable results.

## Exact raw preservation

Stored stdout and stderr must match the child output.

## Large output

Capture large output without excessive memory usage.

## Line range

Verify:

```bash
acap show <id> --lines X:Y
```

returns correct lines.

## Match

Verify:

```bash
acap show <id> --match foo
```

returns relevant matches and respects limits.

## Missing ID

Verify clean failure.

## Corrupt result

Verify no panic.

## Cleanup

Verify completed old results can be removed safely.

## Atomic creation

Verify incomplete temporary results are not treated as valid results.

---

# 41. Benchmarks

Add benchmarks for:

- result creation
- writing large stdout
- writing large stderr
- reading a capsule
- line-range retrieval
- text matching
- raw retrieval

Use deterministic generated fixtures.

Measure at least:

```text
elapsed time
allocations where useful
bytes written
bytes read
```

Do not optimize microseconds at the expense of correctness.

---

# 42. Suggested Project Structure

A possible structure after Phase 2:

```text
agentcap/
├── cmd/
│   └── acap/
│       └── main.go
│
├── internal/
│   ├── executor/
│   ├── reduce/
│   ├── render/
│   ├── result/
│   │   ├── result.go
│   │   └── metadata.go
│   ├── store/
│   │   ├── store.go
│   │   ├── filesystem.go
│   │   └── cleanup.go
│   ├── query/
│   │   ├── lines.go
│   │   └── match.go
│   └── stats/
│
└── go.mod
```

This is only a suggestion.

Prefer simple idiomatic Go boundaries.

Do not create unnecessary abstraction layers.

---

# 43. CLI Scope for Phase 2

By the end of this phase, the primary CLI should be approximately:

```text
acap run <command...>
acap show <id>
acap raw <id>
acap stats
acap clean
```

Optional:

```text
acap show <id> --meta
acap show <id> --match <text>
acap show <id> --path <path>
acap show <id> --lines X:Y
```

Do not expand the CLI significantly beyond this.

---

# 44. Explicit Non-Goals

Do NOT implement the following in Phase 2:

- session-aware compression
- cross-command deduplication
- comparison against previous command results
- Git semantic diff handling
- Git-specific result history
- compiler-specific diagnostics
- test framework parsing beyond existing Phase 1 reducers
- AST parsing
- Tree-sitter
- symbol extraction
- semantic source navigation
- repository-wide indexing
- embeddings
- LLM summarization
- MCP
- Claude Code integration
- Codex-specific integration
- shell hooks
- automatic command rewriting
- agent-specific adapters
- remote storage
- telemetry
- SQLite unless clearly necessary
- plugin system

Do not accidentally implement Phase 3 session state.

A stored command result is not yet a session.

---

# 45. Important Boundary: Result State vs Session State

Phase 2 stores individual results.

It does NOT reason about relationships between results.

Allowed:

```text
result A
result B
result C
```

Not yet allowed:

```text
B is nearly identical to A
C should only return its delta from B
the model has already seen these lines
```

Those concepts belong to Phase 3.

Keep this boundary clear.

---

# 46. Suggested Implementation Order

Implement in this order:

1. define result metadata model
2. implement result ID generation
3. implement filesystem store
4. implement atomic result creation
5. capture stdout/stderr into the store
6. integrate Phase 1 reducer output as stored capsule
7. include result ID in `acap run` output
8. implement `acap show <id>`
9. implement `acap raw <id>`
10. implement line-range retrieval
11. implement text matching
12. add path-based retrieval where reliable
13. integrate stats
14. implement cache cleanup
15. add corruption/error handling
16. add integration tests
17. add benchmarks
18. update README

Do not start with advanced indexing.

Make reliable capture and retrieval work first.

---

# 47. Definition of Done

Phase 2 is complete when all of the following are true.

1. Every `acap run` result receives a result ID.

2. Complete stdout is stored locally.

3. Complete stderr is stored locally.

4. Exit code and command metadata are stored.

5. The Phase 1 capsule is stored.

6. A returned result ID always refers to a valid persisted result.

7. `acap show <id>` works without re-running the command.

8. `acap raw <id>` retrieves the original captured result.

9. Non-zero child executions remain retrievable.

10. Large output does not require unbounded memory buffering.

11. `acap show <id> --lines X:Y` works.

12. `acap show <id> --match <text>` works.

13. Path-specific drill-down works where safely supported.

14. Missing and corrupted results fail cleanly.

15. Cache cleanup exists.

16. Stored files use appropriately private permissions where practical.

17. Statistics distinguish initial reduction from subsequent retrieval.

18. Automated tests cover the result lifecycle.

19. Benchmarks exist for large result storage and retrieval.

20. `go test ./...` passes.

21. `go build ./...` passes.

22. No Phase 3 session-aware logic has been added.

---

# 48. Engineering Priorities

When tradeoffs are necessary, use this order:

1. Exact recoverability of command results
2. Correct exit/error semantics
3. Storage reliability
4. Low memory usage
5. Simple drill-down behavior
6. Low latency
7. Compact agent-facing output
8. Advanced indexing

Do not sacrifice raw-result recoverability for higher compression.

---

# 49. Product Principle

Phase 1 introduced compression.

Phase 2 changes the product model.

AgentCap should now follow this principle:

> Do not permanently discard information just because it was not useful enough to place in the model context immediately.

Instead:

```text
capture everything locally
        |
        v
return the minimum useful capsule
        |
        v
allow cheap drill-down
```

This principle should guide implementation decisions throughout Phase 2.

---

# 50. Final Verification

Before finishing, manually verify flows such as:

```bash
acap run find . -type f
acap show <returned-id>
acap show <returned-id> --match internal
acap show <returned-id> --lines 1:50
acap raw <returned-id>
```

And:

```bash
acap run rg "Workspace" .
acap show <returned-id> --path src/workspace.go
```

Also test a failing command:

```bash
acap run sh -c 'echo failure >&2; exit 7'
```

Verify:

- exit code is still 7
- a result ID exists
- stderr is stored
- `show` works
- `raw` works

Finally run:

```bash
gofmt
go test ./...
go build ./...
```

Then report:

- final project structure
- result-store design
- result ID format
- supported drill-down operations
- storage location
- cleanup policy
- benchmark results
- known limitations

Do not proceed into Phase 3.

Stop once the Phase 2 Definition of Done is satisfied.
