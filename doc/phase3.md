# AgentCap Phase 3 — Session-Aware Compression and Delta Results

## Project

Project name: AgentCap  
CLI command: `acap`  
Language: Go

AgentCap is a command-output reduction layer for AI coding agents.

Previous phases established:

- Phase 0: reliable command execution
- Phase 1: command-aware output reduction
- Phase 2: persistent results, capsules, and drill-down

Phase 3 introduces session awareness.

The central idea is:

> Do not resend information that the coding agent has already seen when a later command result is substantially identical to an earlier result.

Phase 3 should make AgentCap stateful across commands while keeping all raw results independently recoverable.

The goal is not merely deduplication of bytes.

The goal is:

> Reduce repeated context exposure across an agent workflow.

---

# 1. Phase 3 Goal

Implement session-aware result comparison and delta-oriented output.

Phase 3 must provide:

- session identity
- session-local command history
- association between results and sessions
- detection of repeated or substantially unchanged results
- exact-result deduplication
- delta output for supported result types
- conservative fallback when safe delta generation is not possible
- retrieval of both the current full capsule and current raw result
- statistics for session-level savings
- explicit protection against hiding newly introduced errors or changes

The fundamental execution flow becomes:

```text id="z3bnmf"
command
   |
   v
execute + capture
   |
   v
persist full result
   |
   v
reduce full result
   |
   v
compare with relevant previous result
   |
   +--> unchanged
   |
   +--> safe delta
   |
   +--> full capsule fallback
   |
   v
agent-facing output
```

The original result must always remain recoverable through Phase 2 mechanisms.

---

# 2. Core Product Principle

Phase 2 established:

```text id="m2vj41"
capture everything locally
        |
        v
show minimum useful result
        |
        v
allow drill-down
```

Phase 3 extends this with:

```text id="4hrn10"
what has the agent already seen?
        |
        v
what actually changed?
        |
        v
return only the new useful information
```

The most important Phase 3 principle is:

> A repeated command should not consume approximately the same number of model tokens if its meaningful result has not changed.

---

# 3. Session Model

Introduce the concept of an AgentCap session.

A session groups command executions that belong to the same coding-agent workflow.

A reasonable model is:

```go id="z5is05"
type Session struct {
    ID        string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

Each stored result should be associated with:

```text id="ueb1hm"
session_id
sequence_number
```

or equivalent ordering metadata.

Do not turn the session into a large conversational-memory subsystem.

Phase 3 sessions exist specifically to support result comparison and context deduplication.

---

# 4. Session Identification

Support an explicit environment-based session ID.

Recommended:

```text id="9w77bx"
ACAP_SESSION_ID
```

Example:

```bash id="5gtt7j"
ACAP_SESSION_ID=agent-42 acap run git status
```

All commands using the same session ID belong to the same AgentCap session.

Also support a convenient CLI mechanism if useful, such as:

```bash id="q5njnp"
acap session start
```

which may generate:

```text id="d8t1oy"
a7c921
```

However, do not require an interactive shell-state system.

Environment-based session identity should remain sufficient for agent integrations.

---

# 5. Default Session Behavior

Define predictable behavior when `ACAP_SESSION_ID` is not provided.

Possible options:

1. no session-aware compression
2. use a short-lived implicit session
3. generate an isolated session per invocation

The recommended Phase 3 behavior is:

> No cross-command delta compression unless a stable session ID is available.

This is conservative and prevents unrelated commands from being compared accidentally.

Normal Phase 1/2 functionality must continue to work without a session.

---

# 6. Session Scope

Session history must be local.

Do not compare results across unrelated sessions.

For example:

```text id="1phd4j"
session A:
  result 1
  result 2

session B:
  result 3
```

`result 3` must not automatically delta against `result 1` or `result 2`.

This boundary is important for correctness.

---

# 7. Result History

For each session, maintain an ordered history of results.

Example:

```text id="r7a1y7"
session=abc123

001 rg Workspace .
002 cat src/workspace.go
003 go test ./...
004 rg Workspace .
005 go test ./...
```

The history should contain lightweight references to stored Phase 2 results.

Do not duplicate raw stdout/stderr into the session database.

Store result IDs and comparison metadata instead.

---

# 8. Relevant Previous Result Selection

Do not compare every new result with every previous result.

Select a likely relevant baseline.

For Phase 3, the primary baseline should usually be:

> the most recent previous result in the same session produced by an equivalent command identity

Example:

```text id="kfpp2q"
rg Workspace .
```

should normally compare against the previous:

```text id="hl0t2z"
rg Workspace .
```

not against:

```text id="2fvhim"
rg Parser .
```

---

# 9. Command Identity

Define a normalized command identity.

At minimum consider:

- executable
- argv
- working directory

For example:

```text id="85w3xp"
rg Workspace .
cwd=/repo
```

and:

```text id="xmuscp"
rg Workspace .
cwd=/other-repo
```

must not automatically be considered equivalent.

Environment variables should NOT all become part of the identity because that would make comparison impractical.

However, command-specific environment variables that materially affect output may require future consideration.

Keep Phase 3 normalization conservative.

---

# 10. Do Not Over-Normalize Commands

Do not treat these as equivalent automatically:

```text id="wwoynq"
rg Workspace .
rg Workspace src/
```

or:

```text id="rpqx85"
find . -type f
find src -type f
```

or:

```text id="pmzh0p"
go test ./...
go test ./internal/...
```

Exact argv equality plus cwd equality is a safe initial baseline.

More advanced normalization can come later.

---

# 11. Exact Result Deduplication

Implement exact result equality first.

If a new result has exactly the same relevant content as the previous equivalent result, return a compact result.

Example:

```text id="55onai"
@acap 91bc2a unchanged from 7e3d10
```

Optionally include:

```text id="zoyf70"
exit=0
```

when useful.

The exact result must still receive its own result ID because it represents a distinct command execution.

For example:

```text id="f1vhrh"
new result: 91bc2a
baseline:   7e3d10
```

Both must remain individually inspectable.

---

# 12. Result Fingerprints

Add fingerprints to result metadata.

Useful fingerprints may include:

```text id="8spmnv"
stdout hash
stderr hash
combined semantic/capsule hash
```

Use a fast, reliable hash available in Go.

Cryptographic strength is not the primary requirement, but collision behavior should be sane.

For example:

```go id="oepw2w"
sha256
```

is perfectly acceptable even if faster hashes exist.

Avoid premature optimization.

---

# 13. What Counts as Unchanged

A result may be marked unchanged only when AgentCap has strong evidence.

The safest initial definition is:

```text id="0esrs9"
same command identity
AND
same exit code
AND
same stdout
AND
same stderr
```

or equivalent hashes.

Do not mark something unchanged merely because the reduced capsule is identical.

The raw output may contain important changes hidden by the reducer.

---

# 14. Exit Code Changes Are Always Significant

If the previous command exited:

```text id="06f2tx"
0
```

and the current command exits:

```text id="lfsy3u"
1
```

the result is changed even if much of the output is similar.

Likewise:

```text id="42rb8t"
1 -> 0
```

is significant.

Never suppress an exit-status transition.

Example:

```text id="3e3fdk"
@acap b172fe changed from 88a3cd
exit: 0 -> 1

[new diagnostics...]
```

---

# 15. Delta Generation

After exact deduplication works reliably, implement safe delta generation.

A delta should describe:

```text id="7so9gs"
added useful information
removed useful information
changed useful information
```

Do not blindly expose a raw line-by-line unified diff for every command.

Use reducer-aware comparisons where available.

Fallback to a conservative generic delta when appropriate.

---

# 16. Reducer-Aware Delta Interface

Extend the reducer architecture with an optional delta capability.

For example:

```go id="i1y5ga"
type DeltaReducer interface {
    Delta(
        ctx context.Context,
        previous *StoredResult,
        current *StoredResult,
    ) (*DeltaResult, error)
}
```

Do not require every reducer to implement it.

Conceptually:

```text id="rf5f5g"
supports specialized delta
    -> use it

does not support specialized delta
    -> generic comparison or full capsule
```

Keep this extension optional.

---

# 17. Generic Delta Behavior

The generic delta must be conservative.

For text output, a simple line-oriented comparison is acceptable.

Example:

```text id="gp33fu"
@acap 9a12f0 changed from 61d3c2

added:
+ warning: cache invalidated
+ build completed in 14.2s

removed:
- build completed in 13.7s
```

However, if the delta is nearly as large as the current capsule or raw output, return the full current capsule instead.

Do not produce pathological huge diffs merely because delta mode exists.

---

# 18. Delta Cost Heuristic

Before returning a delta, compare its estimated size against the normal current capsule.

For example:

```text id="jv40se"
delta_size < full_capsule_size * threshold
```

then return delta.

Otherwise:

```text id="uj9fbr"
return current capsule
```

A reasonable initial threshold might be around:

```text id="2do6h5"
0.8
```

but keep it configurable or easy to tune internally.

The exact threshold should be validated through benchmarks.

---

# 19. Never Delta Against Raw If Capsule Is Safer

The objective is agent-facing information efficiency, not binary-diff cleverness.

Prefer comparing structured reducer output where that preserves meaning.

For example:

```text id="xawh7h"
317 rg matches -> 319 rg matches
```

should ideally return the two new matches rather than a noisy textual diff over reformatted output.

---

# 20. `rg` / `grep` Delta

This should be one of the first specialized delta implementations.

Suppose the previous result contained:

```text id="h5m7ry"
src/a.go:10:Workspace
src/b.go:20:Workspace
```

and the new result contains:

```text id="yzlcff"
src/a.go:10:Workspace
src/b.go:20:Workspace
src/c.go:30:Workspace
```

Return something like:

```text id="03kp93"
@acap b821f3 delta from a19d40
matches: 2 -> 3

added:
src/c.go:30:Workspace
```

If a match disappears:

```text id="kga6a8"
removed:
src/b.go:20:Workspace
```

Preserve filenames and line numbers.

---

# 21. `find` Delta

For path-oriented output, compare path sets where safe.

Example:

```text id="f7gk3y"
@acap 7f18ac delta from 22c9f0
paths: 5182 -> 5184

added:
src/new_file.go
tests/new_file_test.go
```

If many paths change, fall back to the normal summarized capsule.

---

# 22. `ls` Delta

For directory listings:

```text id="aahqxu"
@acap 348ab1 delta from c918fd

added:
new.go

removed:
old.tmp
```

If metadata such as size or timestamps are part of the output, avoid treating every timestamp change as important unless the reducer explicitly understands it.

Keep behavior conservative.

---

# 23. `cat` Delta

Be careful with `cat`.

Do not attempt sophisticated source-aware diffing in Phase 3.

A simple text diff may be acceptable for moderate text files.

For large files:

```text id="j18rq4"
file changed
lines: 1842 -> 1850
```

plus a small set of changed ranges may be useful.

But if a robust delta is not cheap or reliable, return the current capsule.

AST-aware source diff belongs to later phases.

---

# 24. `tree` Delta

For structured path trees, use path-set differences if the Phase 1 parser already exposes them reliably.

Otherwise fall back.

Do not parse decorative tree output with increasingly fragile heuristics merely to support delta mode.

---

# 25. `du` Delta

Directory size changes can be useful.

Example:

```text id="9x052i"
@acap f192ab delta from 71ac10

build/ 1.1GB -> 1.4GB (+300MB)
vendor/ unchanged
```

Only implement if the existing `du` reducer provides reliable structured data.

---

# 26. `wc` Delta

Because `wc` is already compact, simply showing the current result may be better.

Example:

```text id="ehd7nk"
lines 412 -> 427
words 1810 -> 1872
```

Avoid building unnecessary complexity.

---

# 27. Error Preservation

This is a strict requirement.

New errors must never be hidden merely because most of a result is unchanged.

If the previous result had:

```text id="v2xxhm"
0 errors
```

and the new result has:

```text id="18zf3r"
1 error
```

the delta must surface it prominently.

Likewise, if:

```text id="bf5u0w"
17 errors -> 16 errors
```

the removed error may be useful and should be represented where practical.

---

# 28. stderr Changes Are High Priority

Treat stderr changes conservatively.

If stderr differs materially between baseline and current result:

- do not mark the result unchanged
- surface new stderr content
- prefer full current diagnostics if delta confidence is low

stderr often contains the exact information needed for the next coding-agent step.

---

# 29. Baseline References

Delta output must identify its baseline.

Example:

```text id="lti1cr"
@acap d0a123 delta from 8f31c2
```

This makes the state relationship explicit and debuggable.

The baseline result should remain available:

```bash id="11j1ux"
acap show 8f31c2
```

---

# 30. Current Full Capsule Must Remain Available

If the default `acap run` response is a delta, the coding agent must still be able to retrieve the current result's full capsule.

For example:

```bash id="uy26i3"
acap show d0a123
```

should show the full current Phase 1/2 capsule, not merely repeat the delta.

This distinction is critical.

Store separately:

```text id="15t351"
full capsule
agent-facing delta response
```

if necessary.

---

# 31. Raw Result Must Remain Independent

Likewise:

```bash id="hzovcm"
acap raw d0a123
```

must return the current raw result.

It must not depend on the baseline still existing.

Do not implement delta-only storage.

Every result remains complete and independently recoverable.

---

# 32. Result Metadata Extensions

Extend Phase 2 metadata with fields such as:

```json id="ga0lca"
{
  "session_id": "agent-42",
  "sequence": 17,
  "baseline_result_id": "8f31c2",
  "presentation": "delta",
  "stdout_hash": "...",
  "stderr_hash": "..."
}
```

Do not require every field to exist for non-session results.

Keep metadata backwards-compatible where practical.

---

# 33. Session Storage

A simple filesystem layout is sufficient.

For example:

```text id="sgu5aj"
~/.cache/agentcap/
├── results/
│   ├── 8f31c2/
│   └── d0a123/
└── sessions/
    └── agent-42/
        └── history.jsonl
```

A JSONL history is acceptable and easy to append.

Another simple indexed file design is also fine.

Do not introduce a database without demonstrated need.

---

# 34. History Record

A session history record might contain:

```json id="k6j9al"
{
  "seq": 17,
  "result_id": "d0a123",
  "command_key": "...",
  "created_at": "...",
  "stdout_hash": "...",
  "stderr_hash": "..."
}
```

Keep session history lightweight.

Do not duplicate capsule or raw content here.

---

# 35. Concurrency

Multiple commands may theoretically write into the same session.

Handle this without corrupting session history.

Use simple robust mechanisms such as:

- atomic file replacement
- append semantics
- lightweight file locking if necessary

Do not build a distributed concurrency system.

If strict ordering cannot be guaranteed under concurrent writes, preserve correctness and document the ordering limitation.

---

# 36. Session Sequence Numbers

Sequence numbers are useful but should not become a correctness dependency.

If used, assign monotonically increasing sequence numbers within a session.

For example:

```text id="6hfq7k"
1
2
3
4
```

If concurrency complicates this substantially, timestamps plus result IDs may be sufficient.

Prefer simplicity.

---

# 37. `acap session`

A minimal session CLI may be implemented.

Possible commands:

```bash id="9g6fr2"
acap session start
acap session info
```

`start` may output:

```text id="xtq9kk"
agentcap-session=7ac192
```

Optionally include shell-friendly output:

```bash id="vobnqt"
export ACAP_SESSION_ID=7ac192
```

Do not implement a complex session manager.

---

# 38. Session Listing

`acap session list` is optional.

If implemented, keep output compact.

Example:

```text id="1w63zy"
7ac192 results=42 updated=3m
ab190d results=11 updated=2h
```

Do not prioritize this over delta correctness.

---

# 39. Session Cleanup

Phase 2 cache cleanup must also account for session metadata.

When results expire or are deleted:

- dangling history entries should not crash AgentCap
- baseline lookup should gracefully skip unavailable results

Do not require perfect eager cleanup of every reference.

Lazy tolerance is acceptable.

---

# 40. Missing Baseline

If session metadata identifies a baseline result that no longer exists:

```text id="hqu3xo"
baseline unavailable
```

AgentCap should simply return the current full capsule.

Do not fail the current command.

Delta compression is an optimization, not a correctness requirement.

---

# 41. Baseline Selection Failure

Likewise, if history is corrupt or cannot be read:

```text id="nfmptn"
fall back to full current capsule
```

Do not lose command output because session optimization failed.

This should be a general principle:

> Session-aware compression must fail open to a complete current result.

---

# 42. Delta Confidence

Internally distinguish:

```text id="bghnyc"
exact
safe
unsupported
unsafe
```

or equivalent states.

Only emit a reduced delta when the implementation considers it safe.

Do not expose a complicated confidence system to users unless useful.

The core behavior is:

```text id="oj8s7r"
high confidence delta
    -> emit delta

otherwise
    -> emit full capsule
```

---

# 43. No LLM-Based Delta Generation

Do not use an LLM to compare outputs.

Phase 3 comparison must be:

- local
- deterministic
- fast
- reproducible

Do not add embeddings or semantic similarity models.

AgentCap itself should save LLM tokens, not consume them.

---

# 44. No Approximate Semantic Similarity Yet

Do not implement:

```text id="h2kef6"
these outputs look 93% semantically similar
```

using embeddings, MinHash, fuzzy text models, or language models.

Exact hashing plus command-aware structural comparison is sufficient.

Approximate similarity may be explored separately later if proven necessary.

---

# 45. Already-Emitted Information

Track enough information to know which presentation was returned to the agent.

At minimum record whether a result was presented as:

```text id="1rqauk"
full capsule
unchanged marker
delta
```

This will be useful for later metrics.

Do NOT yet attempt to model every individual line the model has seen.

Phase 3 should stay at result-level and structured-delta-level awareness.

---

# 46. Important Boundary

Do not turn Phase 3 into general agent memory.

AgentCap does not need to know:

- what the LLM understood
- whether it forgot a result
- whether the context window was truncated
- which natural-language conclusions it drew

It only knows what AgentCap returned during the session.

That is enough.

---

# 47. Session-Level Statistics

Extend statistics to measure actual savings across a workflow.

At minimum track:

```text id="ki3a2t"
session commands
raw bytes
full capsule bytes
actual returned bytes
unchanged-result count
delta-result count
full-fallback count
```

Example:

```text id="45mzrb"
session=7ac192
commands=84

raw=18.2MB
full_capsules=2.8MB
returned=1.1MB

unchanged=21
delta=34
full=29

saved_vs_raw=94.0%
saved_vs_stateless=60.7%
```

The exact format can differ.

---

# 48. Stateless vs Stateful Savings

This comparison is important.

Track:

```text id="0pnpmi"
stateless_cost
```

meaning:

> bytes that Phase 1/2 would have returned without session awareness

versus:

```text id="zaz59p"
stateful_cost
```

meaning:

> bytes actually returned after dedup/delta

This quantifies the unique value of Phase 3.

---

# 49. Drill-Down Accounting

If the agent later calls:

```bash id="yiv89w"
acap show ...
acap raw ...
```

those bytes should continue to count toward actual context exposure statistics.

The long-term KPI remains:

```text id="mendye"
initial output
+
all drill-down
=
actual exposed data
```

Do not hide retrieval costs from benchmark results.

---

# 50. Benchmark Workflows

Phase 3 needs workflow benchmarks, not only isolated fixtures.

Create deterministic command sequences representing coding-agent behavior.

Example:

```text id="874pn7"
1. rg Workspace .
2. cat src/workspace.go
3. modify fixture
4. rg Workspace .
5. cat src/workspace.go
6. rg Workspace .
```

Measure:

- raw bytes
- stateless capsule bytes
- stateful bytes
- number of delta fallbacks
- execution overhead
- comparison overhead

---

# 51. Simulated Repository Fixtures

Create small deterministic test repositories under testdata.

Example:

```text id="twj72e"
testdata/repos/simple/
├── src/
├── tests/
└── fixtures/
```

Tests may mutate temporary copies.

Do not rely on the AgentCap repository itself as the only benchmark fixture.

---

# 52. Exact Dedup Benchmark

Include a workflow where the same command is run repeatedly without changes.

Example:

```text id="aq7kym"
rg Workspace .
rg Workspace .
rg Workspace .
```

Expected:

```text id="s8hrhq"
first -> normal capsule
second -> unchanged marker
third -> unchanged marker
```

This should show very high incremental savings.

---

# 53. Small Change Benchmark

Include a workflow where one match is added between searches.

Expected:

```text id="mk3mr2"
first -> full capsule
second -> small delta with one added match
```

This is a key AgentCap use case.

---

# 54. Large Change Benchmark

Include a case where most output changes.

Expected:

```text id="b6zl4y"
delta would be too large
-> full capsule fallback
```

This validates the delta-cost heuristic.

---

# 55. Error Transition Benchmark

Include:

```text id="g45n19"
success -> failure
failure -> success
```

Verify status transitions are always exposed prominently.

---

# 56. Generic Command Benchmark

Use an unsupported command producing deterministic text.

Verify:

- exact result dedup works
- small text delta works where safe
- large delta falls back

Session awareness must not be limited only to specialized reducers.

---

# 57. Performance

Phase 3 comparison must remain cheap relative to command execution.

Avoid repeatedly reading multi-gigabyte results in full when hashes prove equality.

Recommended order:

```text id="2z885t"
compare command identity
        |
        v
compare hashes
        |
        +--> equal -> unchanged
        |
        v
run reducer-aware delta only when needed
```

This prevents unnecessary parsing.

---

# 58. Hash During Capture

Where practical, compute stdout/stderr hashes while capturing them in Phase 2.

For example:

```text id="6v0q6u"
child stdout
   |
   +--> file
   |
   +--> hash writer
```

This avoids re-reading large outputs after execution.

Use `io.MultiWriter` or equivalent where appropriate.

Do not create unnecessary extra copies.

---

# 59. Backwards Compatibility

Existing Phase 2 stored results may not contain Phase 3 fields.

Handle them gracefully.

For example:

- missing session ID -> stateless result
- missing hash -> calculate lazily if practical, or skip comparison
- missing presentation metadata -> assume normal capsule

Do not require users to delete all existing cache data merely to use Phase 3.

---

# 60. Output Format

Keep output compact.

Examples:

## Unchanged

```text id="0vdc59"
@acap b9210f unchanged from 821acc
exit=0
```

## Small delta

```text id="zi430u"
@acap e182ff delta from b9210f
matches 317 -> 319

added:
src/new.go:41:Workspace
src/new.go:88:Workspace
```

## Changed status

```text id="6o4mvt"
@acap 9abc72 changed from e182ff
exit 0 -> 1

stderr added:
src/foo.go:91: undefined: resolveNode
```

## Full fallback

```text id="ln39h8"
@acap 7b12c9 rg matches=811 files=96
...
```

Do not add lengthy explanations such as:

```text id="asr78c"
AgentCap detected that the previous command output...
```

The coding agent does not need prose.

---

# 61. Debug Mode

Extend debug mode to report:

```text id="ljdhto"
session
command key
baseline candidate
hash comparison result
delta reducer selected
delta size
full capsule size
presentation decision
```

Example:

```text id="fzrhyg"
acap debug:
session=7ac192
baseline=8f31c2
stdout_equal=false
stderr_equal=true
delta=312B
capsule=4821B
presentation=delta
```

Debug information goes to stderr only.

---

# 62. Failure Isolation

If specialized delta generation fails:

```text id="kbf8sj"
log in debug mode
return current full capsule
```

Do not fail the user's command.

Similarly:

```text id="oi54rg"
session store failure
baseline read failure
delta parser failure
```

should normally disable the optimization for that result rather than destroy command usability.

The raw current result has already been captured and remains authoritative.

---

# 63. Safety Against Stale Information

Never return only:

```text id="rsj9ue"
unchanged
```

based on:

- command name alone
- capsule equality
- timestamps
- file mtimes
- previous assumptions

Use actual captured result equality.

The agent must not be misled into believing a command result is unchanged when the underlying output changed.

---

# 64. Working Directory Awareness

Include cwd in comparison identity.

This is mandatory.

Example:

```text id="hqlzjy"
/repo-a$ rg foo .
```

must not delta against:

```text id="oyotlo"
/repo-b$ rg foo .
```

even in the same session.

---

# 65. Result Ordering

The baseline should generally be the latest equivalent prior result in that session.

Do not compare against an older result merely because it produces a smaller delta unless explicitly designed later.

The latest baseline best corresponds to what the agent most recently saw.

---

# 66. Baseline Chain

Results may form:

```text id="jfkdaj"
A full
|
B delta from A
|
C delta from B
```

However, every result remains fully stored.

Do not require replaying:

```text id="w05qvf"
A + B delta + C delta
```

to reconstruct C.

C's raw output and full capsule must remain independently stored.

Delta chains are presentation relationships only.

---

# 67. `acap show` Behavior

By default:

```bash id="9isvlq"
acap show <current-id>
```

should display the current full capsule.

Optionally support:

```bash id="xt9071"
acap show <id> --delta
```

to show the original agent-facing delta representation.

This option is useful but not mandatory.

Do not make the default `show` depend on a baseline.

---

# 68. `acap raw` Behavior

No behavioral change from Phase 2.

```bash id="xrq5q5"
acap raw <id>
```

returns the current result's raw content.

Session metadata must not alter raw retrieval.

---

# 69. Optional `acap history`

A compact history command may be useful:

```bash id="1oe1gz"
acap history
```

or:

```bash id="7ux5yu"
acap session history
```

Example:

```text id="f2cfqs"
17 d0a123 rg Workspace . delta
16 8f31c2 rg Workspace . full
15 a712ac cat src/workspace.go full
```

This is optional.

Do not prioritize it over core session compression.

---

# 70. No Git-Specific Semantics Yet

Even though `git status` and `git diff` are ideal session-aware commands, dedicated Git parsing belongs to Phase 4.

Phase 3 may process Git output through the generic reducer if Phase 1 does not already support it.

Do NOT implement:

- Git object awareness
- hunk-aware semantic diff
- branch state modeling
- staged/unstaged modeling

Leave those for Phase 4.

---

# 71. No Build/Test Semantic Delta Yet

Similarly, do not implement:

```text id="td84wm"
7 tests failed -> 3 tests failed
```

using specialized test parsers unless those already exist from prior phases.

Dedicated build/test semantics belong to Phase 5.

Phase 3 should provide the general infrastructure they can later plug into.

---

# 72. No AST Integration

Do not add:

- Tree-sitter
- AST parsing
- symbol identity
- semantic code diff
- call graph awareness

Text and existing reducer structure are sufficient.

---

# 73. No Agent-Specific Integration

Do not implement:

- Claude Code hooks
- Codex hooks
- Gemini CLI hooks
- OpenCode hooks
- shell command interception

Phase 3 should remain agent-independent.

Agents can explicitly invoke `acap`.

Integration belongs to Phase 6.

---

# 74. Suggested Internal Architecture

A possible architecture is:

```text id="orxwnh"
cmd/acap
   |
   v
executor
   |
   v
result store
   |
   v
reducer
   |
   +--> full capsule
   |
   v
session manager
   |
   v
baseline selector
   |
   v
delta engine
   |
   +--> exact unchanged
   |
   +--> specialized delta
   |
   +--> generic delta
   |
   +--> full capsule fallback
   |
   v
renderer
```

Keep session/delta logic outside the execution core.

---

# 75. Suggested Package Layout

For example:

```text id="hbi38w"
internal/
├── executor/
├── reduce/
├── result/
├── store/
├── query/
├── session/
│   ├── session.go
│   ├── history.go
│   └── commandkey.go
├── delta/
│   ├── delta.go
│   ├── generic.go
│   ├── grep.go
│   ├── find.go
│   └── ls.go
├── render/
└── stats/
```

This is guidance, not a strict requirement.

Do not create excessive interfaces merely to match this layout.

---

# 76. Suggested Implementation Order

Implement in this order:

1. add session ID to result metadata
2. implement command identity/key
3. implement lightweight session history
4. associate new results with sessions
5. calculate output hashes during capture
6. locate latest equivalent baseline
7. implement exact-result equality
8. emit compact unchanged result
9. track full-vs-unchanged statistics
10. introduce delta interface
11. implement generic line-oriented delta
12. add delta-size fallback heuristic
13. implement `rg` / `grep` structural delta
14. implement `find` structural delta
15. implement `ls` structural delta
16. add stderr/error safety rules
17. add session-level statistics
18. add workflow integration tests
19. add workflow benchmarks
20. update documentation

Do not start with fuzzy similarity or advanced source diffs.

Exact equality should work first.

---

# 77. Definition of Done

Phase 3 is complete when all of the following are true:

1. Results can be associated with an explicit session.

2. Commands outside a session still work exactly as in Phase 2.

3. Session history records result ordering.

4. Equivalent commands are matched conservatively using argv and cwd.

5. Exact repeated results are detected through full result equality.

6. Repeated unchanged results return a compact unchanged marker.

7. Every repeated execution still receives its own result ID.

8. Every result continues to store complete stdout and stderr.

9. Every result stores its full current capsule.

10. `acap show <id>` retrieves the current full capsule independently of its baseline.

11. `acap raw <id>` retrieves the current raw result independently of its baseline.

12. Exit-code changes are always surfaced.

13. New stderr content is never hidden by unchanged/delta logic.

14. Generic text delta exists.

15. Delta output falls back to the full capsule when the delta is not materially smaller.

16. `rg` / `grep` support structural added/removed match deltas.

17. `find` supports safe path-set deltas.

18. `ls` supports safe entry-set deltas.

19. Missing or corrupted baseline data falls back safely.

20. Session optimization failures never destroy the current command result.

21. Session-level statistics distinguish raw, stateless, and stateful output costs.

22. Workflow tests cover unchanged, small-change, large-change, and error-transition scenarios.

23. Workflow benchmarks quantify Phase 3 savings.

24. `go test ./...` passes.

25. `go build ./...` passes.

26. No Phase 4 Git-specific semantics have been implemented.

27. No Phase 5 build/test-specific semantics have been implemented.

28. No Phase 6 agent integration has been implemented.

---

# 78. Engineering Priorities

When tradeoffs are necessary, use this order:

1. Never hide new actionable information
2. Raw-result correctness and recoverability
3. Correct baseline selection
4. Conservative unchanged detection
5. Safe delta generation
6. Session storage reliability
7. Lower total agent-facing output
8. Low comparison overhead
9. Breadth of specialized delta support

A full current capsule is always preferable to an incorrect or misleading delta.

---

# 79. Phase 3 Success Metric

Do not optimize only for:

```text id="23s7m1"
delta compression percentage
```

The important metric is:

```text id="gp542m"
Total agent-visible bytes/tokens over a realistic multi-command workflow
```

Compare:

```text id="lvq7cm"
raw execution output
vs
Phase 1/2 stateless capsules
vs
Phase 3 stateful presentation
```

The Phase 3 value proposition is demonstrated when repeated coding-agent workflows consume materially less context without increasing information-recovery failures.

---

# 80. Final Verification

Before finishing Phase 3, manually verify at least these workflows.

## Exact repeat

```bash id="g2xz8p"
export ACAP_SESSION_ID=test-session

acap run rg "Workspace" .
acap run rg "Workspace" .
```

The second command should return an unchanged result when the raw result is identical.

Both result IDs must work with:

```bash id="hhw9oc"
acap show <id>
acap raw <id>
```

## Small change

Run:

```bash id="3jadtv"
acap run rg "Workspace" .
```

modify a fixture so one new match appears, then run it again.

The second result should return a small added-match delta.

## Removal

Remove one existing match and verify it is shown as removed.

## Exit transition

Execute an equivalent command once successfully and once with failure-producing state.

Verify:

```text id="svs12u"
exit 0 -> non-zero
```

is prominently visible.

## Large change

Modify enough content that the delta becomes large.

Verify AgentCap returns the normal full capsule instead of a huge delta.

## Missing baseline

Delete or invalidate the baseline result.

Run the command again.

Verify the current command succeeds and returns its full capsule.

Finally run:

```bash id="7na3hk"
gofmt
go test ./...
go build ./...
```

Report:

- final session architecture
- command identity rules
- baseline-selection behavior
- result hash strategy
- implemented specialized delta reducers
- fallback rules
- session-level benchmark results
- known limitations

Do not proceed into Phase 4.

Stop once the Phase 3 Definition of Done is satisfied.
