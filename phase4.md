# AgentCap Phase 4 — Git-Aware Command Compression

## Project

Project name: AgentCap  
CLI command: `acap`  
Language: Go

AgentCap is a command-output reduction layer for AI coding agents.

Previous phases established:

- Phase 0: reliable command execution
- Phase 1: command-aware output reduction
- Phase 2: persistent results, capsules, and drill-down
- Phase 2.1: project-local SQLite storage under `.acap/store.db`
- Phase 3: session-aware deduplication and delta results

Phase 4 introduces Git-aware command handling.

The goal is:

> Represent Git output in a compact structure that preserves the information a coding agent needs to understand repository changes, while avoiding repeated transmission of large diffs and status output.

Phase 4 must remain deterministic and local.

Do not invoke an LLM for Git summarization.

---

# 1. Phase 4 Goal

Implement specialized Git reducers and delta logic for:

```text
git status
git diff
git diff --stat
git show
git log
git branch
```

Focus especially on:

```text
git status
git diff
```

These are the highest-priority commands for coding-agent workflows.

Phase 4 must provide:

- Git command classification
- structured parsing of supported Git output
- compact Git capsules
- file-level diff summaries
- hunk-level progressive disclosure
- Git-aware drill-down
- session-aware Git deltas
- conservative fallback for unsupported Git options/output
- exact preservation of raw output through existing result storage

---

# 2. Core Product Principle

Git output is not arbitrary text.

For example, a coding agent usually does not need to receive a 5,000-line diff immediately.

Instead:

```text
git diff
   |
   v
structured parse
   |
   v
file-level summary
   |
   v
compact capsule
```

Then, only if needed:

```text
acap show <id> --file src/foo.go
```

or:

```text
acap show <id> --hunk 3
```

should expose detailed diff content.

The desired behavior is:

> Show change topology first, change content second.

---

# 3. Supported Git Commands

Phase 4 should explicitly support these command families:

```text
git status
git diff
git show
git log
git branch
```

Support common variants where practical.

Examples:

```bash
git status
git status --short
git status --porcelain
git status --porcelain=v1

git diff
git diff --cached
git diff --staged
git diff HEAD
git diff HEAD~1
git diff --stat

git show
git show <commit>
git show --stat <commit>

git log
git log --oneline
git log -n 20

git branch
git branch -a
```

Do not attempt to support every Git subcommand.

Unsupported Git commands should fall back to the generic reducer.

---

# 4. Do Not Rewrite User Commands

AgentCap must execute the exact Git command requested.

Do NOT silently replace:

```bash
git status
```

with:

```bash
git status --porcelain
```

or replace:

```bash
git diff
```

with an internal alternative.

The command must remain authoritative.

AgentCap may parse its output after execution.

Do not change Git semantics for easier parsing.

---

# 5. Git Detection

Classify a command as Git-aware when the executable is Git.

Examples:

```text
git status
/usr/bin/git status
```

should be recognized.

Avoid matching unrelated executables merely because their name contains `git`.

Use the resolved executable basename where appropriate.

---

# 6. Git Reducer Architecture

Introduce Git-specific reducers without polluting the generic reducer layer.

Possible structure:

```text
internal/
  reduce/
    git/
      status.go
      diff.go
      show.go
      log.go
      branch.go
```

or equivalent.

Each reducer should expose structured intermediate data rather than immediately generating strings.

For example:

```go
type GitStatusResult struct {
    Branch    string
    Ahead     int
    Behind    int
    Files     []GitStatusFile
}
```

This structure can then be:

- rendered
- indexed
- stored
- compared in Phase 3 delta logic

---

# 7. Structured Git Data

Prefer structured internal representations.

For example:

```go
type GitFileChange struct {
    Path       string
    OldPath    string
    Status     string
    Additions  int
    Deletions  int
    Binary     bool
}
```

For diffs:

```go
type GitDiff struct {
    Files []GitDiffFile
}

type GitDiffFile struct {
    OldPath   string
    NewPath   string
    Status    string
    Binary    bool
    Hunks     []GitDiffHunk
}
```

Exact type names are not mandatory.

The important requirement is:

> Do not make rendering logic the only representation of parsed Git state.

---

# 8. `git status`

Implement a specialized reducer for `git status`.

Capture useful information such as:

- branch name
- detached HEAD state
- upstream relation
- ahead count
- behind count
- staged changes
- unstaged changes
- untracked files
- conflicts

Example input may be verbose:

```text
On branch feature/foo
Your branch is ahead of 'origin/feature/foo' by 2 commits.

Changes to be committed:
  modified: src/foo.go

Changes not staged for commit:
  modified: src/bar.go

Untracked files:
  testdata/new.txt
```

Return something like:

```text
@acap 8f31c2 git-status
branch=feature/foo ahead=2

staged=1
M src/foo.go

unstaged=1
M src/bar.go

untracked=1
? testdata/new.txt
```

Keep it compact.

---

# 9. Status Priority

Conflicts must be surfaced prominently.

Example:

```text
@acap a18f9c git-status
branch=feature/foo

conflicts=2
UU src/parser.go
AA src/new.go

unstaged=3
...
```

Do not bury conflicts beneath ordinary modified files.

Conflict state is highly actionable.

---

# 10. Small `git status`

If status is already tiny:

```text
nothing to commit, working tree clean
```

return something compact such as:

```text
@acap 81ac44 git-status clean
branch=main
```

Do not add unnecessary verbose explanation.

---

# 11. `git status --short`

Preserve the concise nature of short/porcelain output.

Do not expand it into verbose prose.

AgentCap may aggregate counts while preserving entries.

Example:

```text
@acap d81bc0 git-status files=3
M  src/foo.go
 M src/bar.go
?? testdata/new.txt
```

---

# 12. Rename and Copy Handling

Handle Git rename/copy paths correctly.

Example:

```text
old.go -> new.go
```

should not be interpreted as two independent files.

Represent it compactly:

```text
R old.go -> new.go
```

Preserve both paths.

---

# 13. `git diff` Priority

`git diff` is the most important Phase 4 feature.

Large diffs should not be returned in full by default.

A coding agent should first receive:

- number of changed files
- total additions/deletions
- changed paths
- per-file additions/deletions
- file status
- important metadata such as rename/binary changes

Example:

```text
@acap 4a21f9 git-diff
files=11 +482 -193

src/workspace.go    M +201 -74
src/parser.go       M +84  -31
src/lexer.go        M +12  -8
tests/workspace_test.go M +185 -80
...
```

---

# 14. Diff File Ordering

Preserve Git's original file ordering unless there is a strong reason not to.

Do not sort alphabetically by default if that changes the order presented by Git.

Deterministic original order is generally the best representation.

---

# 15. File-Level Diff Metadata

For each changed file, capture where possible:

```text
path
old path
status
additions
deletions
binary
mode changes
rename/copy information
```

Avoid parsing only the `+++` and `---` lines.

Git diff headers contain useful semantics.

---

# 16. Hunk Parsing

Parse unified diff hunks.

For example:

```text
@@ -120,8 +120,14 @@ func Resolve(...)
```

Represent internally:

```go
type GitDiffHunk struct {
    Index       int
    OldStart    int
    OldLines    int
    NewStart    int
    NewLines    int
    Header      string
    RawStart    ...
    RawEnd      ...
}
```

Exact fields may differ.

The key requirement is to support precise drill-down later.

---

# 17. Hunk-Level Summary

Do not necessarily return every diff line.

A file summary may include hunk locations.

Example:

```text
src/workspace.go M +201 -74 hunks=4
  h1 @@ -81,7 +81,19 @@ NewWorkspace
  h2 @@ -412,18 +424,41 @@ Resolve
  h3 @@ -702,9 +737,18 @@ lookup
  h4 @@ -911,11 +955,21 @@ updateIndex
```

This lets the coding agent decide which portion to inspect.

Keep hunk headers only when useful.

---

# 18. `acap show --file`

Extend Phase 2 drill-down for Git results.

Support:

```bash
acap show <result-id> --file src/workspace.go
```

For a Git diff result, this should show the complete diff for that file.

It must NOT re-run Git.

It must retrieve the data from the stored result.

Example:

```text
@acap 4a21f9 file=src/workspace.go
+201 -74 hunks=4

diff --git ...
...
```

---

# 19. Hunk Drill-Down

Support:

```bash
acap show <id> --file src/workspace.go --hunk 2
```

or an equivalent compact interface.

This should return only the selected hunk plus enough header/context information to understand it.

Do not require the agent to retrieve the entire file diff to inspect one hunk.

---

# 20. `--hunk` Identity

Hunk IDs may be result-local and file-local.

For example:

```text
file=src/workspace.go hunk=1
file=src/workspace.go hunk=2
```

No globally stable hunk identity is required.

A later `git diff` execution may produce different hunk numbering.

Document that hunk numbers are scoped to the stored result.

---

# 21. Context Preservation

When showing a hunk, preserve the original unified diff context.

Do not strip all unchanged context lines.

The coding agent may need surrounding code to understand the change.

However, do not add unrelated sections from other hunks.

---

# 22. Binary Files

Do not attempt to interpret binary diffs.

Represent them as metadata.

Example:

```text
assets/model.bin binary changed
```

If Git provides size or mode information, preserve it where useful.

Do not dump binary data.

---

# 23. New and Deleted Files

Represent clearly:

```text
A src/new.go +142
D src/old.go -87
```

Do not force these into generic modified-file semantics.

For a newly added file, additions may equal the full file length.

That is acceptable.

---

# 24. Renames

For renames:

```text
R src/old.go -> src/new.go
```

include similarity percentage where Git provides it and where useful.

Example:

```text
R90 src/old.go -> src/new.go
```

Do not expand rename-only changes into misleading delete/add pairs.

---

# 25. Mode Changes

Preserve mode changes compactly.

Example:

```text
M scripts/build.sh mode 100644 -> 100755
```

A mode-only change is significant even if there are no text lines changed.

Do not report it as unchanged.

---

# 26. `git diff --stat`

This command is already compact.

Avoid over-compressing it.

AgentCap may normalize it to the standard Git diff capsule format, but should not make the output less informative.

If output is already below the normal compression threshold, returning nearly unchanged output is acceptable.

---

# 27. `git diff --cached` / `--staged`

Treat staged diffs separately from unstaged diffs.

The command identity already differs because argv differs.

Do not merge them.

The capsule should optionally indicate scope:

```text
git-diff staged
```

Example:

```text
@acap a51fc2 git-diff staged
files=4 +88 -21
...
```

---

# 28. Arbitrary Revision Diffs

Commands such as:

```bash
git diff HEAD~1 HEAD
git diff abc123 def456
```

should use the Git diff reducer if the output is a normal unified Git diff.

Do not require the diff to represent the working tree.

The capsule should not assume staged/unstaged semantics unless known from argv.

---

# 29. Path-Limited Diff

For:

```bash
git diff -- src/parser.go
```

respect the user-specified path scope.

Do not compare the result with an unrelated full-repository diff in Phase 3.

The command key already includes argv and prevents that baseline error.

---

# 30. Combined and Merge Diffs

Git may produce combined diff formats for merges.

If the parser does not reliably support them:

```text
fall back to generic reducer
```

or a conservative Git summary.

Do not incorrectly parse complex combined diff formats as ordinary unified diffs.

Correct fallback is better than wrong structure.

---

# 31. `git show`

Support common `git show` output.

A normal `git show` combines:

```text
commit metadata
+
patch
```

Return a capsule such as:

```text
@acap 4ca112 git-show
commit=8fd31a2
author=...
date=...
subject=Fix workspace resolution

files=3 +81 -29

src/workspace.go M +55 -18
src/resolver.go  M +18 -9
tests/...        M +8  -2
```

Do not include long commit messages repeatedly unless useful.

---

# 32. `git show --stat`

Treat similarly to compact diff stats.

Preserve:

- commit ID
- subject
- changed-file summary

Avoid unnecessary patch parsing when no patch is present.

---

# 33. Commit Message Handling

Preserve:

- subject line
- optionally a short body preview

Do not discard the commit subject.

For long commit bodies, use conservative truncation with an omitted marker.

Do not attempt to summarize commit messages semantically.

---

# 34. `git log`

Support common log output.

For standard or oneline logs, preserve:

- commit ID
- date where available
- author where available
- subject

Example:

```text
@acap 71bc32 git-log commits=20

8fd31a2 Fix workspace resolution
29ae114 Add parser fixtures
91ac7cd Refactor index lookup
...
```

If the user explicitly asks for a small `-n` value, avoid additional compression.

---

# 35. Large `git log`

For hundreds or thousands of commits, return a bounded set.

Useful strategy:

```text
commits=842 shown=30 omitted=812
```

Include the most recent entries in the same order Git returned them.

Do not reorder by author or subject.

---

# 36. `git branch`

Support:

```bash
git branch
git branch -a
git branch -r
```

Preserve:

- current branch
- branch names
- remote/local distinction where visible

For huge branch lists, group or truncate conservatively.

Example:

```text
@acap 21cd90 git-branch current=feature/foo
local=12 remote=41

local:
* feature/foo
  main
  release/1.4
...
```

---

# 37. Detached HEAD

Represent detached HEAD explicitly.

Example:

```text
branch=DETACHED
head=8fd31a2
```

Do not invent a branch name.

---

# 38. Git Error Handling

Git failures must preserve diagnostics.

Examples:

```text
fatal: not a git repository
fatal: ambiguous argument
fatal: bad revision
```

Do not replace these with an empty Git capsule.

If the command exits non-zero and parsing confidence is low, prefer the normal stderr-preserving fallback.

---

# 39. Raw Output Remains Authoritative

All Phase 4 results must continue to preserve the exact raw stdout/stderr through Phase 2.1 storage.

Git parsing is an optimization layer.

The raw result remains authoritative.

These commands must always remain available:

```bash
acap show <id>
acap raw <id>
```

---

# 40. Structured Git Index Storage

For Git results, store lightweight parsed indexes if useful for drill-down.

For example:

```text
result_git_files
result_git_hunks
```

may be represented either:

- in SQLite tables
- as serialized structured data associated with the result

Prefer SQLite if it simplifies reliable file/hunk lookup.

Do not put enormous raw diff lines into relational tables.

---

# 41. Suggested SQLite Extensions

A possible schema:

```sql
CREATE TABLE git_diff_files (
    result_id TEXT NOT NULL,
    file_index INTEGER NOT NULL,
    old_path TEXT,
    new_path TEXT,
    status TEXT NOT NULL,
    additions INTEGER,
    deletions INTEGER,
    binary INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (result_id, file_index)
);
```

Potential hunk metadata:

```sql
CREATE TABLE git_diff_hunks (
    result_id TEXT NOT NULL,
    file_index INTEGER NOT NULL,
    hunk_index INTEGER NOT NULL,
    old_start INTEGER,
    old_lines INTEGER,
    new_start INTEGER,
    new_lines INTEGER,
    header TEXT,
    raw_start INTEGER,
    raw_end INTEGER,
    PRIMARY KEY (result_id, file_index, hunk_index)
);
```

Exact schema is flexible.

Use migrations from Phase 2.1.

---

# 42. Do Not Store Every Diff Line in SQLite

Keep raw patch text in the existing raw object store.

SQLite should index:

```text
where is the file?
where is the hunk?
what metadata belongs to it?
```

The raw object store should answer:

```text
what are the exact bytes?
```

Maintain the Phase 2.1 principle:

```text
SQLite = structured state
objects = large opaque content
```

---

# 43. Byte Offsets vs Line Offsets

For Git drill-down, byte offsets into the raw stdout object may be more efficient than rescanning large diffs repeatedly.

Consider storing:

```text
raw_start_byte
raw_end_byte
```

for files and hunks.

This enables:

```text
seek -> stream selected range
```

without loading the entire raw diff.

Line offsets are also acceptable if already supported.

Choose the simpler robust implementation.

---

# 44. Phase 3 Integration

Git reducers should participate in session-aware compression.

The baseline rules remain:

```text
same session
same cwd
same argv
latest equivalent previous result
```

Do not change Phase 3 command identity rules merely for Git.

---

# 45. Git Status Delta

Implement a specialized delta for `git status`.

Example previous state:

```text
M src/a.go
M src/b.go
```

Current:

```text
M src/a.go
M src/b.go
?? src/c.go
```

Return:

```text
@acap c91ac2 delta from a71f20 git-status

added:
? src/c.go
```

If a file disappears because it became clean:

```text
resolved:
M src/b.go
```

or equivalent compact representation.

---

# 46. Git Status State Transitions

Track meaningful transitions.

Examples:

```text
untracked -> staged
unstaged -> staged
conflicted -> modified
modified -> clean
```

Represent them if the structured parser can do so reliably.

Example:

```text
src/foo.go unstaged -> staged
```

Do not reduce all transitions to simple added/removed path entries when richer state is cheap and reliable.

---

# 47. Git Diff Delta

This is a major Phase 4 feature.

Suppose an agent runs:

```bash
git diff
```

then edits one file and runs:

```bash
git diff
```

again.

Do NOT return the entire updated diff if most of it was already seen.

Return the difference between the previous and current Git diff structure.

Example:

```text
@acap 9ab21c delta from 81ac34 git-diff

changed:
src/workspace.go
  +14 -3
  new hunk @@ -511,6 +511,17 @@ Resolve

unchanged_files=10
```

This is substantially more useful than raw textual diff-of-diffs.

---

# 48. Diff-of-Diff Semantics

Avoid naive unified-diffing of two Git patch outputs when structured comparison is available.

Compare:

```text
file identity
hunks
changed lines
```

instead.

The output should describe:

> What changed in the repository diff since the agent last inspected it?

not:

> What textual characters changed in Git's patch formatting?

---

# 49. Stable File Identity

For ordinary modifications, file path is sufficient.

For rename cases, use old/new path information to avoid treating a rename as arbitrary deletion/addition where possible.

Do not attempt sophisticated Git object identity tracking in Phase 4.

Structured paths and status are sufficient.

---

# 50. Hunk Comparison

Use conservative hunk comparison.

Useful inputs:

```text
old/new ranges
hunk header
hunk content hash
```

If an existing hunk remains byte-identical, it can be treated as unchanged.

If it changes materially, include it in the current delta.

Do not attempt semantic code equivalence.

---

# 51. Diff Delta Fallback

If structural comparison becomes ambiguous:

```text
return the current full Git capsule
```

Examples:

- parser failure
- unsupported combined diff
- extensive rename ambiguity
- extremely large structural delta
- previous structured metadata unavailable

Phase 3's fallback principle still applies.

---

# 52. Diff Delta Cost

Compare the generated Git delta size against the full Git capsule.

If:

```text
delta >= useful fraction of full capsule
```

return the full current capsule instead.

Use the existing Phase 3 size heuristic.

Do not create a separate incompatible heuristic unless necessary.

---

# 53. `git log` Session Delta

A simple useful delta is possible.

Example:

First:

```text
A
B
C
```

Later:

```text
D
A
B
C
```

Return:

```text
@acap ... git-log delta
new_commits=1

D Add workspace cache
```

Keep this conservative.

Do not attempt commit graph reasoning.

---

# 54. `git branch` Session Delta

Support simple added/removed/current-branch changes.

Example:

```text
added:
feature/new-index

current:
main -> feature/new-index
```

Avoid elaborate branch ancestry analysis.

---

# 55. `git show` Delta

Dedicated session delta for `git show` is low priority.

Exact repeated output deduplication from Phase 3 is sufficient initially.

Do not spend excessive effort on specialized `git show` deltas unless implementation is straightforward.

---

# 56. Git Diff Capsule Example

A target output:

```text
@acap d831ac git-diff
scope=unstaged
files=8 +214 -73

M src/workspace.go +91 -21 hunks=3
M src/parser.go +44 -18 hunks=2
M src/index.go +22 -9 hunks=2
A src/cache.go +57 hunks=1
D src/old_cache.go -25 hunks=1
R95 src/foo.go -> src/bar.go +0 -0

binary=1
assets/test.bin
```

Keep the format compact and predictable.

---

# 57. File Drill-Down Example

```bash
acap show d831ac --file src/workspace.go
```

Output:

```text
@acap d831ac file=src/workspace.go
M +91 -21 hunks=3

@@ -81,7 +81,19 @@ NewWorkspace
...

@@ -412,18 +424,41 @@ Resolve
...

@@ -911,11 +955,21 @@ updateIndex
...
```

This is allowed to be significantly larger than the capsule because it is explicitly requested.

---

# 58. Hunk Drill-Down Example

```bash
acap show d831ac --file src/workspace.go --hunk 2
```

Output:

```text
@acap d831ac file=src/workspace.go hunk=2

@@ -412,18 +424,41 @@ Resolve
...
```

Do not include unrelated hunks.

---

# 59. File Matching

When resolving:

```bash
--file src/foo.go
```

prefer exact normalized path matching.

Do not silently select a different file because of fuzzy similarity.

If multiple paths are ambiguous, return a concise ambiguity error.

Optional suffix matching may be supported only when unique.

---

# 60. Paths With Spaces

Git paths may contain spaces and escaped characters.

Handle them correctly.

Do not parse paths by naive whitespace splitting.

Use Git diff/status syntax carefully.

Add fixtures containing:

```text
src/file with spaces.go
```

and unusual Unicode filenames.

---

# 61. Quoted Git Paths

Git may quote unusual paths depending on configuration.

Do not assume every path is raw UTF-8 with no escaping.

Support common Git path quoting rules where necessary.

If parsing confidence is low, fall back rather than produce incorrect path metadata.

---

# 62. Git Config Variability

User Git configuration can affect output.

Examples include:

```text
color.ui
core.quotePath
log formatting
pager settings
```

AgentCap should handle ANSI cleanup through existing infrastructure where applicable.

Do not mutate user Git configuration.

Do not set global Git configuration.

---

# 63. Pager Behavior

Commands such as:

```bash
git log
git show
```

may use pagers in interactive environments.

AgentCap should avoid hanging on a pager where existing execution architecture permits safe control.

If AgentCap already has established non-interactive environment handling, reuse it.

Do not globally modify Git config.

If environment changes such as:

```text
GIT_PAGER=cat
```

would alter normal command semantics, do not introduce them silently without careful justification.

Prefer handling actual captured output from the requested invocation.

---

# 64. ANSI Handling

Phase 1 already supports ANSI cleanup where appropriate.

Git reducers may operate on normalized text for parsing while raw stored output remains unchanged.

This distinction is important:

```text
raw object = exact captured bytes
parser input = normalized representation if required
```

Do not overwrite raw data with normalized output.

---

# 65. Maximum Diff Size

Very large Git diffs must not cause excessive memory usage.

Avoid requiring:

```text
entire 500 MB diff -> one giant Go string
```

if the current architecture supports streaming or temporary-file parsing.

A practical initial implementation may parse moderately sized diffs in memory with a clear limit and fall back for extreme cases.

Do not crash or exhaust memory.

---

# 66. Parser Failure Behavior

If Git-specific parsing fails:

```text
raw result remains stored
generic/current fallback remains available
```

Do not fail the child Git command.

Example debug behavior:

```text
git reducer parse failed: unsupported combined diff
fallback=generic
```

Normal user output should remain concise.

---

# 67. Parsing Must Be Deterministic

Do not use:

- LLMs
- embeddings
- fuzzy semantic models

for Git parsing.

Use:

- Git syntax
- deterministic state machines
- line parsing
- structured metadata

The same input should produce the same capsule.

---

# 68. Testing Fixtures

Create deterministic fixtures for:

```text
git-status-clean.txt
git-status-modified.txt
git-status-conflict.txt
git-status-renames.txt

git-diff-small.txt
git-diff-large.txt
git-diff-add-delete.txt
git-diff-rename.txt
git-diff-binary.txt
git-diff-mode-change.txt
git-diff-path-spaces.txt

git-show.txt
git-log.txt
git-branch.txt
```

Prefer fixtures generated from known temporary Git repositories where practical.

---

# 69. Integration Test Repositories

Create temporary Git repositories during tests.

Use actual Git commands where Git is available.

Test workflows such as:

```text
init repository
commit baseline
modify file
acap run git status
acap run git diff
stage file
acap run git diff --cached
```

Skip gracefully if Git is unavailable.

Do not make all unit tests depend on an external Git installation.

---

# 70. Status Tests

Test:

- clean repository
- staged modification
- unstaged modification
- untracked file
- deleted file
- rename
- conflict
- detached HEAD where practical
- ahead/behind parsing where practical
- filenames with spaces
- Unicode paths

---

# 71. Diff Tests

Test:

- one modified file
- many modified files
- added file
- deleted file
- rename
- binary file
- mode change
- multiple hunks
- zero-context diff if requested
- paths with spaces
- unusual hunk headers
- empty diff

---

# 72. Drill-Down Tests

Verify:

```bash
acap show <id> --file ...
```

returns only the correct file diff.

Verify:

```bash
acap show <id> --file ... --hunk ...
```

returns only the requested hunk.

Verify invalid file and hunk selectors fail cleanly.

---

# 73. Session Delta Tests

Create workflow tests.

Example:

```text
git diff #1
modify one line
git diff #2
```

Verify:

- #2 uses #1 as baseline
- only newly changed Git structure is presented where appropriate

Then:

```text
git diff #3
```

with no repository changes.

Verify:

```text
unchanged
```

through Phase 3 exact deduplication.

---

# 74. Status Delta Tests

Example sequence:

```text
clean
modify file
stage file
modify second file
```

Verify each status delta exposes only the new state transition while never hiding conflicts/errors.

---

# 75. Large Change Test

Modify many files so the structural delta approaches the full capsule size.

Verify AgentCap falls back to the full current capsule.

Do not force delta presentation.

---

# 76. Benchmarking

Benchmark representative operations:

```text
parse 10-file diff
parse 100-file diff
parse large multi-hunk diff
render file summary
file drill-down
hunk drill-down
Git delta comparison
```

Measure:

```text
parse time
allocations
capsule bytes
raw bytes
delta bytes
```

---

# 77. Phase 4 KPI

Measure:

```text
raw Git output
vs
Phase 4 full Git capsule
vs
Phase 4 + Phase 3 stateful delta
```

Useful metrics:

```text
initial reduction ratio
stateful reduction ratio
drill-down bytes
total agent-visible bytes
```

Do not report only raw compression percentage.

---

# 78. Realistic Git Workflow Benchmark

Create a deterministic workflow approximating coding-agent behavior:

```text
1. git status
2. git diff
3. inspect one file
4. modify fixture
5. git diff
6. git status
7. stage file
8. git diff --cached
9. git status
```

Measure total agent-visible output.

This workflow is more meaningful than isolated command benchmarks.

---

# 79. Storage Integration

All Git-specific metadata must be associated with existing result IDs.

Do not create a second result identity system.

The authoritative relationship remains:

```text
result ID
  |
  +--> metadata in .acap/store.db
  +--> capsule
  +--> Git structured metadata
  +--> raw stdout/stderr objects
```

---

# 80. Schema Migration

Use the existing Phase 2.1 migration mechanism.

Do not delete and recreate `.acap/store.db`.

Add Git-specific schema incrementally.

For example:

```text
schema vN -> vN+1
```

Tests must cover migration from the previous Phase 3 schema.

---

# 81. Do Not Add Git Repository Indexing

Phase 4 is command-result parsing.

Do NOT create a persistent model of the entire Git repository.

Do not index:

- all commits
- all trees
- all blobs
- repository history graph
- blame information

Only store enough structured information for executed command results.

---

# 82. Do Not Call `git` Again for Parsing

After the user's Git command completes, avoid running additional Git commands merely to interpret the result.

For example, do not automatically run:

```bash
git diff --numstat
git status --porcelain
git rev-parse
```

after every command unless there is an exceptionally strong reason.

Phase 4 should primarily parse the captured command output.

Additional hidden subprocesses:

- add latency
- can observe a changed repository state
- complicate semantics
- make results harder to reproduce

Prefer single-command capture.

---

# 83. No Source-Code AST Parsing

Do not parse changed source files with AST tools.

Phase 4 Git semantics are:

```text
file
patch
hunk
line
```

not:

```text
function
class
symbol
call graph
```

Semantic source integration belongs to later phases.

---

# 84. No Build/Test Handling

Do not add compiler or test parsing in Phase 4.

Commands such as:

```text
go test
cargo test
clang
gcc
```

belong to Phase 5.

Keep Phase 4 focused.

---

# 85. No Agent-Specific Hooks

Do not add:

- Claude Code hooks
- Codex hooks
- Gemini CLI hooks
- shell interception
- MCP

Those remain Phase 6 work.

Phase 4 should work entirely through:

```bash
acap run git ...
```

---

# 86. Suggested Implementation Order

Implement in this order:

1. introduce Git command classification
2. implement common Git parsing utilities
3. implement `git status`
4. add status fixtures/tests
5. implement basic `git diff` file parsing
6. parse additions/deletions and statuses
7. implement hunk indexing
8. implement compact Git diff capsule
9. implement `acap show --file`
10. implement `--hunk`
11. persist Git file/hunk indexes in SQLite
12. implement Git status session delta
13. implement Git diff structural delta
14. implement binary/rename/mode handling
15. implement `git show`
16. implement `git log`
17. implement `git branch`
18. add large-output limits/fallbacks
19. add integration tests
20. add workflow benchmarks
21. update documentation

Do not begin with every Git command simultaneously.

Make `status` and `diff` excellent first.

---

# 87. Definition of Done

Phase 4 is complete when all of the following are true:

1. AgentCap detects supported Git commands reliably.

2. Unsupported Git commands safely use existing fallback behavior.

3. `git status` produces a compact structured capsule.

4. Clean status is represented compactly.

5. Staged, unstaged, untracked, deleted, renamed, and conflicted states are preserved.

6. Conflicts are surfaced prominently.

7. `git diff` produces a file-level summary rather than returning large patches by default.

8. Additions/deletions are reported per file where reliably derivable.

9. Added, deleted, renamed, binary, and mode-only changes are represented correctly.

10. Unified diff hunks are indexed.

11. `acap show <id> --file <path>` retrieves a stored file diff without re-running Git.

12. Hunk-level drill-down works.

13. Raw Git output remains independently recoverable.

14. Git parsing failure falls back safely.

15. Git status has specialized session-aware delta behavior.

16. Git diff has specialized structural delta behavior.

17. Exact repeated Git output still benefits from Phase 3 unchanged detection.

18. Large Git deltas fall back to the full Git capsule when appropriate.

19. `git show` has a useful commit + change summary.

20. `git log` has compact history output.

21. `git branch` has compact branch output.

22. Git structured metadata is integrated with `.acap/store.db`.

23. Schema changes use Phase 2.1 migrations.

24. Large diffs do not cause uncontrolled memory usage.

25. Tests cover representative Git edge cases.

26. Workflow benchmarks quantify Git token/output savings.

27. `go test ./...` passes.

28. `go build ./...` passes.

29. No Phase 5 build/test parsing has been implemented.

30. No Phase 6 agent-specific integration has been implemented.

---

# 88. Engineering Priorities

When tradeoffs are necessary, use this order:

1. Never misrepresent repository state
2. Preserve conflicts and errors
3. Exact raw-result recoverability
4. Correct file/status parsing
5. Reliable file/hunk drill-down
6. Conservative Git delta behavior
7. Lower total agent-visible output
8. Low latency
9. Support for additional Git variants

A generic fallback is always better than an incorrect Git interpretation.

---

# 89. Phase 4 Success Criterion

Phase 4 succeeds when a typical coding-agent loop:

```text
git status
git diff
edit
git diff
edit
git diff
```

no longer repeatedly sends the complete repository patch into the model context.

The expected flow becomes:

```text
first git diff
  -> compact file/hunk topology

agent requests one relevant file/hunk
  -> targeted detail

second git diff
  -> only newly changed diff structure

third unchanged git diff
  -> unchanged marker
```

This is the core Phase 4 value proposition.

---

# 90. Final Verification

Before finishing Phase 4, manually verify a repository workflow.

Start with a clean repository:

```bash
export ACAP_SESSION_ID=git-test

acap run git status
```

Modify several files:

```bash
acap run git status
acap run git diff
```

Verify the diff capsule contains file-level summaries rather than the full patch.

Use the returned result ID:

```bash
acap show <id> --file src/foo.go
acap show <id> --file src/foo.go --hunk 1
acap raw <id>
```

Verify:

- file drill-down is correct
- hunk drill-down is correct
- raw output exactly reflects the captured Git result

Modify only one additional region and run:

```bash
acap run git diff
```

Verify the session-aware output emphasizes only the newly changed Git structure.

Run again without changes:

```bash
acap run git diff
```

Verify Phase 3 unchanged detection applies.

Then test:

```bash
git add src/foo.go
acap run git status
acap run git diff --cached
```

Also test:

- new file
- deleted file
- rename
- binary file
- executable mode change
- conflict fixture
- filename with spaces

Finally run:

```bash
gofmt
go test ./...
go build ./...
```

Report:

- Git command coverage
- internal Git data model
- SQLite schema additions
- file/hunk indexing strategy
- Git delta strategy
- fallback rules
- large-diff behavior
- benchmark results
- known limitations

Do not proceed into Phase 5.

Stop once the Phase 4 Definition of Done is satisfied.
