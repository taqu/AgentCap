# AgentCap Phase 2.1 — Project-Local SQLite Storage

## Project

Project name: AgentCap  
CLI command: `acap`  
Language: Go

AgentCap is a command-output reduction layer for AI coding agents.

Previous phases established:

- Phase 0: reliable command execution
- Phase 1: command-aware output reduction
- Phase 2: persistent results, capsules, and drill-down

Phase 2.1 must finalize the local storage architecture before session-aware compression is introduced in Phase 3.

The main change is:

> Move AgentCap result persistence to a project-local `.acap` directory backed by SQLite.

The canonical database path should be:

```text
<project-root>/.acap/store.db
```

Because the database already lives inside `.acap`, do not use redundant names such as:

```text
agentcap.db
agentcap-store.db
acap.db
```

Use:

```text
store.db
```

The filename should describe its role, not repeat the product name.

---

# 1. Phase 2.1 Goal

Implement a project-local storage architecture with:

- project-root discovery
- `.acap/` storage directory
- `.acap/store.db`
- SQLite-backed metadata and result indexing
- schema versioning
- raw stdout/stderr object storage
- result lifecycle integration
- cleanup support
- migration away from the previous global filesystem result store where practical
- tests for storage isolation and repository discovery

The primary goal is:

> Make AgentCap state naturally scoped to the project being analyzed.

Phase 2.1 must establish the storage foundation required by future session-aware result history without implementing Phase 3 session logic.

---

# 2. Canonical Storage Layout

Use the following project-local layout:

```text
<project-root>/
├── .acap/
│   ├── store.db
│   └── objects/
│       ├── ...
│       └── ...
├── .git/
├── src/
└── ...
```

The canonical database path is:

```text
.acap/store.db
```

The `objects/` directory stores large raw command outputs or other blob-like data that should not necessarily live directly inside SQLite.

Do not create multiple databases for individual AgentCap subsystems.

Use one local SQLite database as the primary metadata store.

---

# 3. Why Project-Local Storage

AgentCap results are inherently project-specific.

Commands such as:

```bash
rg Workspace .
find . -type f
cat src/workspace.go
go test ./...
git diff
```

derive their meaning from the repository or project in which they execute.

Therefore, result state should follow the project rather than the user's machine globally.

The desired model is:

```text
repo-a/
  .acap/store.db

repo-b/
  .acap/store.db
```

Results from repository A must not accidentally become candidates for repository B.

This isolation will become especially important in Phase 3.

---

# 4. Do Not Use the Current Directory Directly

Do not create `.acap` in every working directory from which AgentCap is invoked.

For example:

```text
/repo/
  src/
    parser/
```

If the user runs:

```bash
cd /repo/src/parser
acap run rg Workspace ..
```

AgentCap should NOT create:

```text
/repo/src/parser/.acap/
```

if `/repo` is the detected project root.

It should use:

```text
/repo/.acap/
```

instead.

Otherwise AgentCap state will fragment across arbitrary subdirectories.

---

# 5. Project Root Discovery

Implement deterministic project-root discovery.

Use the following priority order:

```text
1. nearest ancestor containing .acap/
2. nearest ancestor representing the repository root
3. current working directory
```

For the repository-root check, initially support Git repositories.

Search upward for:

```text
.git
```

This may be a directory or other valid Git repository marker.

Do not require invoking the `git` executable simply to discover the root if filesystem traversal is sufficient.

---

# 6. `.acap` Takes Priority

If an ancestor contains:

```text
.acap/
```

that ancestor should normally be treated as the AgentCap project root.

Example:

```text
/work/
  .acap/
  src/
    parser/
      current-directory
```

Even if another marker exists higher in the tree, use:

```text
/work/
```

as the AgentCap root.

This makes AgentCap project boundaries explicit and stable.

---

# 7. Git Root Fallback

If no `.acap` ancestor exists, search upward for the nearest Git project marker.

Example:

```text
/repo/
  .git/
  src/
    parser/
```

Running AgentCap from:

```text
/repo/src/parser
```

should resolve:

```text
project root = /repo
storage      = /repo/.acap
database     = /repo/.acap/store.db
```

---

# 8. Non-Git Projects

AgentCap must not require Git.

If no `.acap` ancestor and no repository marker can be found:

```text
project root = current working directory
```

For example:

```bash
cd /tmp/my-project
acap run find . -type f
```

may create:

```text
/tmp/my-project/.acap/
```

This behavior should be predictable and documented.

---

# 9. Explicit Root Override

Consider supporting an explicit override.

For example:

```text
ACAP_ROOT=/repo
```

or a CLI option such as:

```bash
acap --root /repo run ...
```

This is useful for:

- non-Git projects
- monorepos
- nested repositories
- testing
- agent integrations

If implemented, explicit configuration must take precedence over automatic discovery.

Keep this feature simple.

---

# 10. Storage Directory Creation

Create:

```text
<project-root>/.acap/
```

only when AgentCap actually needs persistent storage.

Do not create `.acap` merely for:

```bash
acap --help
```

or other non-storage operations.

Use appropriately restrictive permissions where supported.

AgentCap output may contain:

- source code
- compiler diagnostics
- paths
- test output
- environment-derived information
- secrets accidentally printed by child commands

The storage directory should therefore not be world-readable by default.

---

# 11. Git Ignore Behavior

`.acap/` should normally not be committed to source control.

However, AgentCap must NOT silently edit:

```text
.gitignore
```

during normal command execution.

Do not mutate repository configuration merely because AgentCap ran.

Documentation should recommend:

```gitignore
.acap/
```

Optionally, a future or explicit initialization command may offer guidance.

Automatic `.gitignore` modification is outside the core Phase 2.1 requirement.

---

# 12. SQLite Database

Use SQLite for structured storage.

Canonical path:

```text
.acap/store.db
```

SQLite should store structured metadata such as:

- result IDs
- argv
- cwd
- exit status
- timestamps
- durations
- reducer identity
- capsule
- output sizes
- output hashes where available
- object references
- statistics
- schema version information

Do not introduce an external database server.

AgentCap must remain a single local application.

---

# 13. SQLite Driver

Choose a Go SQLite implementation appropriate for a single-binary CLI.

Prefer:

- reliable
- actively maintained
- well-tested
- simple deployment

Avoid introducing a runtime dependency on a separately installed SQLite command-line executable.

The database must be accessed through Go code.

Consider portability implications before choosing a CGO-only dependency.

If a pure-Go driver is appropriate and mature enough for the project, prefer it for easy distribution.

Document the chosen driver and why it was selected.

---

# 14. Database Responsibilities

SQLite is the authoritative structured metadata store.

Do not use a mixture of:

```text
meta.json
history.json
index.json
store.db
```

for overlapping metadata unless there is a clear reason.

Prefer:

```text
SQLite:
  structured state

objects/:
  large raw byte streams
```

This avoids multiple competing metadata sources.

---

# 15. Raw Output Storage

Do not automatically store every large stdout/stderr payload as a SQLite BLOB.

AgentCap may eventually process very large outputs.

Examples:

```text
100 MB build logs
large generated source listings
large test output
large search results
```

Storing all of these directly inside `store.db` can cause unnecessary database growth and maintenance overhead.

Use external object files for raw output.

---

# 16. Object Store

Use:

```text
.acap/objects/
```

for raw output payloads.

A possible structure:

```text
.acap/
├── store.db
└── objects/
    ├── 2a/
    │   └── 2a743... 
    ├── 91/
    │   └── 91ac8...
    └── ...
```

Content-addressed paths are recommended but not mandatory.

A flat object directory is acceptable initially if simpler.

Do not create deeply nested structures without a demonstrated reason.

---

# 17. Content Addressing

Consider using the output hash as the raw object identity.

For example:

```text
sha256(raw stdout)
```

can produce:

```text
objects/2a/2a743fe...
```

Advantages:

- exact duplicate payloads can share storage
- object integrity can be checked
- Phase 3 already benefits from output hashes
- cleanup can eventually use references

However, do not implement complex garbage collection merely to support content addressing.

Correctness first.

---

# 18. stdout and stderr Objects

Store stdout and stderr independently.

For example, a result row may reference:

```text
stdout_object
stderr_object
```

Do not merge streams merely for storage convenience.

Phase 2 semantics must remain intact.

The result should continue to preserve:

- stdout separately
- stderr separately
- exit status independently

---

# 19. Small vs Large Payloads

Keep the policy simple.

The recommended initial design is:

```text
capsule            -> SQLite
metadata           -> SQLite
raw stdout/stderr  -> object files
```

Do this even for small raw output unless there is a compelling implementation reason not to.

A single storage rule is easier to reason about than:

```text
< 64 KB -> BLOB
>= 64 KB -> file
```

Hybrid size thresholds can be introduced later if benchmarks justify them.

---

# 20. Suggested Database Schema

A minimal initial schema may contain:

```sql
CREATE TABLE schema_info (
    version INTEGER NOT NULL
);

CREATE TABLE results (
    id TEXT PRIMARY KEY,
    created_at TEXT NOT NULL,
    started_at TEXT,
    cwd TEXT NOT NULL,
    argv_json TEXT NOT NULL,
    exit_code INTEGER NOT NULL,
    duration_ns INTEGER NOT NULL,

    reducer TEXT,

    stdout_object TEXT,
    stderr_object TEXT,

    stdout_bytes INTEGER NOT NULL,
    stderr_bytes INTEGER NOT NULL,

    stdout_hash TEXT,
    stderr_hash TEXT,

    capsule TEXT NOT NULL
);
```

The exact schema may differ.

Use appropriate SQLite types and constraints.

Avoid designing the full Phase 3 schema now.

---

# 21. Statistics Storage

Existing AgentCap statistics may also move into SQLite.

A simple schema is sufficient.

For example:

```sql
CREATE TABLE stats (
    key TEXT PRIMARY KEY,
    value INTEGER NOT NULL
);
```

Or derive statistics from result rows where practical.

Do not build a large analytics subsystem.

The goal is to remove unnecessary scattered global state.

---

# 22. Schema Versioning

Introduce explicit schema versioning now.

This is important because Phase 3 and later phases will likely add:

- session IDs
- baseline references
- presentation type
- command keys
- additional indexes

The database must know which schema it contains.

A simple:

```text
schema version = 1
```

is sufficient.

---

# 23. Migrations

Implement a small migration mechanism.

For example:

```text
version 1 -> version 2
version 2 -> version 3
```

Do not introduce a full migration framework unless necessary.

Simple ordered Go migration functions are sufficient.

Example:

```go
func migrateV1ToV2(tx *sql.Tx) error
```

Every migration should run transactionally where SQLite permits it.

---

# 24. New Database Initialization

When:

```text
.acap/store.db
```

does not exist:

1. create `.acap`
2. initialize SQLite
3. create the current schema
4. record the schema version
5. create `objects/` if needed

Initialization must be safe to repeat.

Do not recreate or destroy a valid existing database.

---

# 25. Unsupported Future Schema

If AgentCap opens a database with a schema version newer than the binary understands:

```text
database version > supported version
```

fail clearly.

Do not attempt to downgrade automatically.

Example:

```text
AgentCap store schema is newer than this acap binary supports
```

Keep the diagnostic concise.

---

# 26. SQLite Configuration

Use SQLite pragmas intentionally.

Evaluate options such as:

```text
foreign_keys
busy_timeout
journal_mode
synchronous
```

Do not blindly copy high-performance configurations from server applications.

AgentCap is a local CLI with relatively small concurrent write volume.

Correctness and crash tolerance are more important than extreme transaction throughput.

WAL mode may be useful, especially for future concurrent agent commands, but justify its use and test cleanup implications.

---

# 27. Transaction Boundaries

Result creation should be atomic from the database perspective.

A partially stored result must not appear as a valid completed result.

Conceptually:

```text
capture command
write raw objects
begin transaction
insert result metadata
commit
```

or another ordering with equivalent correctness.

If metadata commit fails, do not expose the result ID as successfully persisted.

---

# 28. Object Write Atomicity

Write new object files safely.

Prefer:

```text
temporary file
     |
     v
fsync/close where appropriate
     |
     v
atomic rename
```

Do not leave partially written objects under their final names.

If content addressing is used and the object already exists, reuse it safely.

---

# 29. Database Is Project-Local State

Do not store absolute machine-global result relationships outside the project store unless necessary.

Everything needed to:

```bash
acap show <id>
acap raw <id>
```

should normally be available from:

```text
<project-root>/.acap/
```

This makes the AgentCap state self-contained.

---

# 30. `acap show` Root Resolution

Commands such as:

```bash
acap show 8f31c2
```

must resolve the current project store using the same project-root discovery logic.

Do not search every AgentCap database on the machine for an ID.

If the current project's store does not contain the result, return a clear missing-result error.

This preserves project isolation.

---

# 31. `acap raw` Root Resolution

Likewise:

```bash
acap raw 8f31c2
```

must use:

```text
current project -> .acap/store.db
```

It must not perform a machine-wide result search.

Explicit cross-project access can be considered later if ever needed.

---

# 32. Current Working Directory Metadata

Every result should continue to store the exact working directory used for command execution.

This is separate from the project root.

For example:

```text
project_root=/repo
cwd=/repo/src/parser
```

Both may matter.

Do not replace the result cwd with the detected root.

---

# 33. Project Root Metadata

It is acceptable to store the detected project root in the database or derive it from the database location.

Avoid repeatedly storing redundant values in every row unless useful.

The important distinction is:

```text
storage scope = project root
execution scope = cwd
```

---

# 34. Previous Global Store

Previous Phase 2 implementations may store data under something such as:

```text
~/.cache/agentcap/
```

Phase 2.1 should stop using the global result store as the primary storage location.

New command results must go to:

```text
<project-root>/.acap/
```

---

# 35. Migration From Existing Global Results

Do NOT attempt to automatically assign arbitrary historical global results to projects unless their original project path is known reliably.

Possible safe behavior:

- leave old global results untouched
- document them as legacy data
- optionally provide an explicit migration tool if metadata contains an unambiguous cwd

Automatic migration is not required for Phase 2.1.

Do not guess result ownership.

---

# 36. Global Cache Use After Phase 2.1

Do not use:

```text
~/.cache/agentcap/
```

for normal project command results.

A global AgentCap directory may remain appropriate in the future for:

- global configuration
- temporary anonymous state
- application metadata

But Phase 2.1 should avoid adding new global state unless needed.

The primary result store is project-local.

---

# 37. Cache Semantics

Treat `.acap` as disposable local state.

Important principle:

> `.acap` is a cache and execution-history store, not durable user data.

If `.acap` is deleted:

- source code must remain unaffected
- AgentCap must recreate it when necessary
- old result IDs simply become unavailable

Do not require users to back up `.acap`.

---

# 38. Git Clean and Repository Cleanup

Be aware that commands such as:

```bash
git clean -fdx
```

may remove `.acap`.

This is acceptable under the disposable-state model.

AgentCap must tolerate the store disappearing between invocations.

On the next execution, recreate it as needed.

Do not assume the database persists forever.

---

# 39. Do Not Depend on Stored State for Command Correctness

The wrapped command must still execute correctly if:

```text
.acap/
```

does not exist.

Storage enables:

- result retrieval
- compression history
- future session logic

It must not become necessary for basic child-process execution semantics.

If persistence fails, follow the storage-failure behavior established by Phase 2.

---

# 40. Store Interface

Preserve or introduce a clean storage abstraction.

For example:

```go
type Store interface {
    CreateResult(ctx context.Context, result *Result) error
    GetResult(ctx context.Context, id string) (*StoredResult, error)
    DeleteResult(ctx context.Context, id string) error
    Close() error
}
```

The exact API is flexible.

Do not expose raw SQL throughout unrelated packages.

SQLite-specific logic should remain concentrated in the storage package.

---

# 41. Suggested Package Layout

A possible layout:

```text
internal/
├── project/
│   └── root.go
├── executor/
├── reduce/
├── result/
├── store/
│   ├── store.go
│   ├── sqlite.go
│   ├── schema.go
│   ├── migrate.go
│   └── objects.go
├── query/
├── render/
└── stats/
```

Keep project-root discovery separate from database mechanics.

Do not over-engineer repository abstractions.

---

# 42. Store Initialization API

A useful flow is:

```text
cwd
 |
 v
ResolveProjectRoot()
 |
 v
OpenProjectStore(root)
 |
 +--> .acap/store.db
 |
 +--> .acap/objects/
```

This keeps location policy out of individual CLI commands.

`run`, `show`, `raw`, `clean`, and `stats` should reuse the same root/store resolution logic.

---

# 43. Indexes

Add SQLite indexes only where justified.

Phase 2.1 likely needs very few.

Potential examples:

```sql
CREATE INDEX idx_results_created_at
ON results(created_at);
```

This may help cleanup.

Do NOT yet create Phase 3-specific indexes for:

```text
session_id
command_key
baseline_id
```

unless the columns are already required.

Avoid speculative schema design.

---

# 44. Cleanup

Update:

```bash
acap clean
```

to operate on the current project's:

```text
.acap/store.db
.acap/objects/
```

Do not clean stores belonging to other repositories.

Project-local behavior should be explicit.

---

# 45. Object Garbage Collection

When a result is deleted, raw output objects may become unreferenced.

Implement one of these simple approaches:

Option A:

```text
one object per result stream
```

Delete it with the result.

Option B:

```text
content-addressed shared objects
```

Delete only when no database row references the object.

If reference counting or garbage collection becomes complicated, prefer Option A initially.

Do not over-engineer deduplicated object storage in Phase 2.1.

---

# 46. `acap clean` Safety

Cleanup must never delete files outside:

```text
<project-root>/.acap/
```

Be extremely conservative with filesystem deletion.

Resolve and validate cleanup paths.

Never derive recursive delete paths directly from untrusted database strings without validation.

---

# 47. SQLite Corruption

Handle database corruption gracefully.

Do not panic.

Return a clear diagnostic such as:

```text
AgentCap store is unreadable: .acap/store.db
```

Do not silently delete the database and lose stored results.

Automatic destructive repair is out of scope.

---

# 48. Missing Object Files

If a result row exists but its raw stdout/stderr object is missing:

- `acap show` should still work if the capsule exists
- `acap raw` should return a clear error for the missing stream
- no panic
- do not re-run the original command automatically

Treat the result as partially corrupted.

---

# 49. Concurrent Access

Design the SQLite layer so multiple AgentCap processes can safely access the same project store.

This is important even before Phase 3 because coding agents may invoke commands concurrently.

At minimum:

- concurrent readers must be safe
- normal short writes should not corrupt the database
- busy errors should have reasonable handling

Use a small busy timeout if appropriate.

Do not build a custom daemon solely to serialize access.

---

# 50. No Long-Lived Database Daemon

AgentCap should remain:

```text
CLI process
  -> open SQLite
  -> perform operation
  -> close
```

Do not introduce:

- background service
- socket server
- local daemon
- database broker

Phase 2.1 should preserve the single-binary CLI model.

---

# 51. Result IDs

Do not change result IDs merely because storage changed.

Existing Phase 2 result ID semantics should remain valid.

IDs should stay:

- compact
- project-local
- collision-resistant enough for local use

SQLite should treat the result ID as the primary logical key.

---

# 52. Result ID Scope

After Phase 2.1, result IDs should be understood as project-local identifiers.

For example:

```text
repo-a: 8f31c2
repo-b: 8f31c2
```

could theoretically both exist.

This is acceptable.

Do not require globally unique IDs across all projects unless already naturally provided by the current implementation.

---

# 53. Capsule Storage

Store the complete Phase 1/2 capsule inside SQLite.

Capsules are generally small and frequently accessed.

Do not store them as separate object files unless benchmarks show a reason.

Example:

```text
results.capsule
```

should allow:

```bash
acap show <id>
```

to retrieve the normal capsule without opening a raw object file.

---

# 54. Raw Retrieval

`acap raw <id>` should:

1. resolve project root
2. open `.acap/store.db`
3. find the result
4. resolve the stdout/stderr object references
5. stream raw content

Do not load the entire raw object into memory.

Large result retrieval should remain streaming.

---

# 55. Drill-Down Queries

Existing Phase 2 features such as:

```bash
acap show <id> --lines X:Y
acap show <id> --match foo
acap show <id> --path path
```

must continue to work.

It is acceptable for these operations to stream/search raw object files referenced by SQLite.

Do not regress Phase 2 progressive disclosure behavior.

---

# 56. Statistics

`acap stats` should become project-local by default.

Running:

```bash
cd repo-a
acap stats
```

should report AgentCap activity for repo A.

Running:

```bash
cd repo-b
acap stats
```

should report repo B.

Do not silently combine statistics from unrelated projects.

A future explicit global stats feature may be added separately.

---

# 57. Debug Output

Extend debug mode to report storage resolution.

Useful information:

```text
project_root=/repo
store=/repo/.acap/store.db
result_id=8f31c2
```

Do not print full database internals during normal execution.

Debug output goes to stderr.

---

# 58. Testing: Root Discovery

Add tests for:

```text
.acap ancestor exists
.git ancestor exists
neither exists
nested directories
nested Git repositories
explicit root override
```

Use temporary directories.

Do not depend on the developer's real repository layout.

---

# 59. Testing: Project Isolation

Create two temporary project roots:

```text
project-a/
project-b/
```

Run/store results independently.

Verify:

```text
project-a result
```

cannot be retrieved while operating in:

```text
project-b
```

unless an explicit override points back to project A.

This is a key Phase 2.1 property.

---

# 60. Testing: SQLite Lifecycle

Test:

- first database creation
- reopening existing database
- schema initialization
- schema migration
- unsupported future schema
- transaction rollback
- missing database recreation
- corruption handling where practical

---

# 61. Testing: Raw Objects

Test:

- stdout object creation
- stderr object creation
- empty stdout
- empty stderr
- large objects
- exact byte preservation
- missing object
- failed object write
- duplicate object behavior if content addressing is implemented

---

# 62. Testing: Concurrent Access

Add a reasonable concurrency test.

For example, multiple goroutines/processes may insert independent result records into the same temporary project store.

Verify:

- no corruption
- all successful committed results remain readable
- transient SQLite locking is handled appropriately

Do not turn the test into a stress-testing framework.

---

# 63. Testing: Cleanup

Verify:

```bash
acap clean
```

only affects the active project's `.acap`.

Create neighboring directories with similarly named paths and ensure they are untouched.

Deletion safety matters more than cleanup performance.

---

# 64. Benchmarks

Benchmark at least:

```text
database open
result insert
capsule lookup
metadata lookup
large raw object write
large raw object streaming read
```

Use representative sizes such as:

```text
10 KB
1 MB
10 MB
```

A 100 MB fixture may be generated dynamically if useful, but do not commit giant fixture files.

---

# 65. Performance Goal

SQLite should not introduce significant latency for normal AgentCap usage.

Avoid:

- opening multiple database connections unnecessarily in one CLI invocation
- serializing the same metadata repeatedly
- reading raw object content for metadata-only operations
- running VACUUM automatically during normal commands

The store should remain lightweight.

---

# 66. No Automatic VACUUM on Every Cleanup

Do not run:

```sql
VACUUM;
```

after every result deletion.

That can be unexpectedly expensive.

If VACUUM support is useful, make it explicit or threshold-based later.

Phase 2.1 cleanup should prioritize responsiveness.

---

# 67. No Phase 3 Session Logic

This is a strict boundary.

Do NOT implement:

- sessions
- `ACAP_SESSION_ID` behavior
- command history chains
- previous-result lookup
- baseline selection
- unchanged-result detection
- delta generation
- cross-command deduplication
- session-aware statistics

The database may make these easier later, but Phase 2.1 must not implement them.

---

# 68. No Git-Specific Result Semantics

Git may be used only as a project-root marker.

Do NOT implement:

- `git status` parsing
- `git diff` parsing
- branch tracking
- commit tracking
- Git-aware result invalidation

Those remain future work.

---

# 69. No AST or Semantic Indexing

Do NOT add:

- Tree-sitter
- AST storage
- source symbols
- semantic indexes
- call graphs
- embeddings

`store.db` is not intended to become a general code intelligence database in this phase.

---

# 70. No SQLite Full-Text Search Yet

Do not add FTS tables merely for:

```bash
acap show --match
```

Existing streaming/text matching is sufficient.

SQLite FTS can be evaluated later if measurement shows a benefit.

Avoid speculative complexity.

---

# 71. No Remote Synchronization

`.acap/store.db` and `.acap/objects/` must remain local.

Do NOT implement:

- database uploads
- cloud backup
- synchronization
- telemetry
- remote cache
- shared team database

Project-local does not mean repository-shared.

`.acap` should normally remain ignored and machine-local.

---

# 72. Suggested Implementation Order

Implement in this order:

1. implement project-root discovery
2. define canonical `.acap` paths
3. introduce SQLite driver
4. create schema version table
5. implement initial schema
6. implement database open/create lifecycle
7. implement raw object store
8. migrate result metadata persistence to SQLite
9. migrate capsule persistence to SQLite
10. migrate raw stdout/stderr persistence to objects
11. update `acap run`
12. update `acap show`
13. update `acap raw`
14. update `acap stats`
15. update `acap clean`
16. add schema migration framework
17. add corruption/error handling
18. add root-isolation tests
19. add concurrent-access tests
20. add storage benchmarks
21. update README/documentation

Do not begin Phase 3 work.

---

# 73. Definition of Done

Phase 2.1 is complete when all of the following are true:

1. AgentCap resolves a stable project root.

2. An existing `.acap` ancestor has priority.

3. Git repository root is used when no `.acap` ancestor exists.

4. Current working directory is used when no project marker exists.

5. Normal project results are stored under:

```text
<project-root>/.acap/
```

6. The SQLite database is named exactly:

```text
store.db
```

7. Structured result metadata is stored in SQLite.

8. Capsules are stored in SQLite.

9. Raw stdout and stderr remain independently recoverable.

10. Raw outputs are stored as external object files rather than requiring large SQLite BLOBs.

11. `acap show <id>` continues to work.

12. `acap raw <id>` continues to work.

13. Existing drill-down functionality continues to work.

14. `acap stats` reports project-local statistics.

15. `acap clean` only cleans the active project's storage.

16. Schema versioning exists.

17. Database migrations are supported through a small internal mechanism.

18. Opening a newer unsupported schema fails clearly.

19. Database corruption does not cause a panic.

20. Missing raw objects produce clear errors.

21. Concurrent normal access does not corrupt the store.

22. New command results are no longer primarily stored under `~/.cache/agentcap`.

23. `.gitignore` is not silently modified.

24. `.acap` is treated as disposable local state.

25. `go test ./...` passes.

26. `go build ./...` passes.

27. No Phase 3 session or delta logic has been implemented.

---

# 74. Engineering Priorities

When tradeoffs are necessary, use this order:

1. Result integrity
2. Project isolation
3. Exact raw-output recoverability
4. Safe filesystem behavior
5. SQLite transaction correctness
6. Low memory usage
7. Simple schema design
8. Low execution overhead
9. Future extensibility

Do not over-design the schema for hypothetical future features.

Phase 2.1 exists to provide a clean and reliable storage foundation.

---

# 75. Storage Design Principle

The target architecture is:

```text
project
  |
  v
.acap/
  |
  +--> store.db
  |      |
  |      +--> metadata
  |      +--> capsules
  |      +--> indexes
  |      +--> statistics
  |
  +--> objects/
         |
         +--> raw stdout
         +--> raw stderr
```

The conceptual separation is:

```text
SQLite = structured state
objects = large opaque byte streams
```

Keep this boundary simple.

---

# 76. Future Phase 3 Compatibility

While Phase 3 must not be implemented now, avoid making the schema impossible to extend with future fields such as:

```text
session_id
sequence
command_key
baseline_result_id
presentation
```

Use migrations rather than speculative nullable columns where possible.

The Phase 2.1 database should be easy to evolve, not pre-populated with future architecture.

---

# 77. Final Verification

Before finishing Phase 2.1, manually verify the following workflow.

From a repository root:

```bash
acap run rg "Workspace" .
```

Verify:

```text
.acap/store.db
.acap/objects/
```

exist.

Then:

```bash
acap show <result-id>
acap raw <result-id>
acap stats
```

Verify all operations use the project-local store.

Next:

```bash
cd src/some/nested/directory
acap show <same-result-id>
```

Verify AgentCap discovers the same project root and accesses the same database.

Then create or enter a different repository and verify:

```bash
acap show <first-project-result-id>
```

does not retrieve the first project's result.

Verify a non-Git directory:

```bash
mkdir /tmp/acap-test-project
cd /tmp/acap-test-project
acap run printf 'hello\n'
```

and confirm:

```text
/tmp/acap-test-project/.acap/store.db
```

is created.

Finally run:

```bash
gofmt
go test ./...
go build ./...
```

Report:

- selected SQLite driver and rationale
- project-root discovery rules
- final `.acap` layout
- SQLite schema
- schema migration mechanism
- raw object layout
- concurrency behavior
- cleanup behavior
- benchmark results
- known limitations

Do not proceed into Phase 3.

Stop once the Phase 2.1 Definition of Done is satisfied.
