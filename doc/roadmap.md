# AgentCap — Roadmap and Progress Summary

## 1. Project Overview

Project name: AgentCap  
CLI command: `acap`  
Implementation language: Go

AgentCap is a local command-output reduction and progressive-disclosure layer for AI coding agents.

The core problem is that coding agents frequently consume large amounts of context from shell/tool output that contains:

- repeated information
- successful build/test noise
- large search results
- large Git diffs
- repeated diagnostics
- unchanged results from previously executed commands

AgentCap intercepts command execution, stores the complete result locally, and returns only the information likely to be useful for the agent's next action.

The product direction is not simply:

```text
large output
   ->
truncate/filter
   ->
smaller output
```

Instead:

```text
execute
   |
   v
capture full result
   |
   v
structure / reduce
   |
   v
compact capsule
   |
   v
agent
```

with recoverability:

```text
capsule
   |
   +--> targeted drill-down
   |
   +--> full raw result
```

and session-aware behavior:

```text
previous result
       +
current result
       |
       v
only newly useful information
```

---

# 2. Core Product Principles

AgentCap has evolved around several key principles.

## Preserve correctness before compression

Priority order:

```text
1. command semantics
2. actionable information
3. errors and diagnostics
4. raw-result recoverability
5. state correctness
6. low latency
7. token/output reduction
```

A larger correct result is preferable to an aggressively compressed but misleading result.

---

## Full raw results remain recoverable

Compression is a presentation optimization.

It is not destructive storage.

Every captured command result can be represented as:

```text
raw stdout
raw stderr
metadata
structured indexes
capsule
presentation
```

The current result never depends on replaying previous deltas to reconstruct it.

---

## Progressive disclosure

The desired interaction model is:

```text
compact result
   |
   v
targeted detail
   |
   v
broader detail
   |
   v
raw only when necessary
```

Typical retrieval commands:

```bash
acap show <id>
acap show <id> --file ...
acap show <id> --hunk ...
acap show <id> --errors
acap show <id> --warnings
acap show <id> --test ...
acap show <id> --match ...
acap show <id> --lines ...
acap raw <id>
```

`acap raw` is the last-resort path rather than the normal workflow.

---

## Stateful compression

AgentCap tracks what it has already returned within an agent session.

Repeated execution should not repeatedly consume equivalent context.

The central principle is:

> Do not resend information the coding agent has already seen when the new command result can be represented safely as unchanged or as a small delta.

---

## Local-first architecture

AgentCap remains local.

No remote service is required.

No LLM is used internally for output summarization.

Parsing and comparison are deterministic.

---

# 3. Storage Architecture

Storage was finalized during Phase 2.1.

Project state is stored under:

```text
<project-root>/.acap/
```

Canonical layout:

```text
.acap/
├── store.db
└── objects/
```

`store.db` is SQLite.

The division of responsibilities is:

```text
SQLite
  -> structured metadata
  -> capsules
  -> indexes
  -> statistics
  -> schema versions
  -> relationships

objects/
  -> raw stdout
  -> raw stderr
  -> large opaque output payloads
```

The filename is intentionally:

```text
store.db
```

because the surrounding `.acap` directory already provides the AgentCap namespace.

---

# 4. Project Root Resolution

AgentCap storage is project-local rather than based directly on the invocation directory.

Root discovery follows approximately:

```text
1. nearest ancestor containing .acap/
2. nearest Git repository root
3. current working directory
```

This means:

```text
/repo/.acap/store.db
```

is reused when commands are run from:

```text
/repo/
/repo/src/
/repo/src/parser/
/repo/tests/
```

while exact execution `cwd` is still retained independently.

This distinction is important:

```text
project root
    !=
execution cwd
```

---

# 5. Current CLI Model

The core manual interface includes:

```bash
acap run <command...>
acap show <result-id>
acap raw <result-id>
acap stats
acap clean
```

Additional drill-down selectors are available depending on result type.

Agent integration added in Phase 6 allows normal shell commands to benefit from AgentCap without manually writing:

```bash
acap run ...
```

for every command.

---

# 6. Phase 0 — Execution Core

Status:

```text
COMPLETE
```

Purpose:

> Establish a reliable command execution foundation before implementing compression.

Implemented concepts:

- `acap run <command...>`
- direct argv execution
- no implicit shell reconstruction
- stdin forwarding
- stdout preservation
- stderr preservation
- exit-code propagation
- cwd inheritance
- environment inheritance
- context cancellation
- signal handling
- large-output-safe streaming
- command-not-found handling
- execution timing
- tests for execution semantics

Major design rule:

```text
AgentCap must not change command semantics.
```

Phase 0 deliberately contained no compression.

---

# 7. Phase 1 — Core Command Output Compression

Status:

```text
COMPLETE
```

Purpose:

> Add deterministic command-aware reducers while keeping unknown commands usable.

Primary reducer targets:

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

Architecture:

```text
known command
   -> specialized reducer

unknown command
   -> generic reducer
```

Important behavior:

- small output remains mostly unchanged
- large output is summarized structurally
- stderr is preserved conservatively
- ANSI/progress noise can be removed from agent-facing presentation
- generic fallback remains conservative
- raw size vs returned size is measurable

The Phase 1 output introduced the basic AgentCap capsule concept.

---

# 8. Phase 2 — Result Store and Progressive Disclosure

Status:

```text
COMPLETE
```

Purpose:

> Make compression recoverable rather than destructive.

Introduced:

- compact result IDs
- persistent result storage
- complete raw stdout retention
- complete raw stderr retention
- stored command metadata
- stored capsules
- `acap show`
- `acap raw`
- line-range retrieval
- text matching
- path-oriented drill-down
- cache cleanup
- result lifecycle tests

Core model:

```text
capture
   ->
capsule
   ->
drill-down
```

The important shift was:

> AgentCap no longer discarded omitted information.

---

# 9. Phase 2.1 — Project-Local SQLite Storage

Status:

```text
COMPLETE
```

Purpose:

> Finalize storage architecture before stateful result relationships were introduced.

Major changes:

```text
global ~/.cache/agentcap style storage
                  ->
project-local .acap storage
```

Canonical database:

```text
.acap/store.db
```

Implemented architecture:

```text
SQLite = structured state
objects = large raw streams
```

Phase 2.1 also established:

- project root discovery
- SQLite schema versioning
- migration mechanism
- atomic database/result operations
- raw object storage
- concurrent access handling
- project-local stats
- project-local cleanup
- corruption handling
- isolation between repositories

`.acap` is treated as disposable local state rather than durable user data.

It should normally be ignored by Git.

AgentCap does not silently modify `.gitignore`.

---

# 10. Phase 3 — Session-Aware Compression

Status:

```text
COMPLETE
```

Purpose:

> Stop repeatedly transmitting information already shown to the coding agent.

Introduced:

- AgentCap sessions
- result history
- stable command identity
- baseline selection
- stdout/stderr hashes
- exact-result deduplication
- delta presentation
- structural delta hooks
- session-level statistics

Command identity is conservative and includes at least:

```text
argv
cwd
```

Baseline selection uses:

```text
same session
same command identity
latest equivalent result
```

Exact unchanged detection requires strong evidence such as:

```text
same exit status
same stdout
same stderr
```

rather than capsule equality alone.

Typical unchanged presentation:

```text
@acap <new-id> unchanged from <previous-id>
exit=0
```

---

# 11. Phase 3 Delta Principle

A delta is only a presentation optimization.

Results remain independently complete.

For:

```text
A full
|
B delta from A
|
C delta from B
```

the stored representation is still:

```text
A = complete
B = complete
C = complete
```

C never requires reconstructing:

```text
A + delta B + delta C
```

This keeps storage and recovery robust.

---

# 12. Phase 3 Stateful KPI

Phase 3 introduced an important distinction:

```text
raw cost
vs
stateless AgentCap cost
vs
stateful AgentCap cost
```

The target metric is not merely:

```text
initial compression %
```

but:

```text
total agent-visible bytes/tokens
over an entire workflow
```

including subsequent drill-down.

---

# 13. Phase 4 — Git-Aware Compression

Status:

```text
COMPLETE
```

Purpose:

> Treat Git output as structured repository-change information instead of arbitrary terminal text.

Primary supported families:

```text
git status
git diff
git diff --stat
git show
git log
git branch
```

Highest priority:

```text
git status
git diff
```

---

# 14. Git Diff Model

Large Git diffs are represented initially as change topology.

Example concept:

```text
files=11 +482 -193

src/workspace.go       M +201 -74 hunks=4
src/parser.go          M +84  -31 hunks=2
tests/workspace_test.go M +185 -80 hunks=5
...
```

The model is:

```text
change topology first
change contents on demand
```

Drill-down includes:

```bash
acap show <id> --file src/workspace.go
acap show <id> --file src/workspace.go --hunk 2
```

No Git command re-execution is required.

---

# 15. Git Structured State

Phase 4 introduced structured concepts such as:

```text
changed file
old/new path
status
additions
deletions
binary state
mode change
rename/copy
hunks
raw ranges
```

Special cases include:

- new files
- deleted files
- renames
- binary changes
- executable/mode changes
- conflicts
- detached HEAD

Unsupported or ambiguous Git formats fall back conservatively.

---

# 16. Git Stateful Delta

Repeated:

```bash
git diff
```

does not necessarily resend the complete diff topology.

The desired behavior is:

```text
first git diff
 -> complete compact topology

edit

second git diff
 -> only newly changed Git structure

third git diff with no changes
 -> unchanged marker
```

This is effectively a structured "diff of the repository diff".

---

# 17. Phase 5 — Build / Compiler / Test-Aware Compression

Status:

```text
COMPLETE
```

Purpose:

> Preserve actionable failures while aggressively reducing successful and repetitive build/test output.

Primary supported toolchains:

```text
Go:
  go build
  go test

C/C++:
  gcc
  g++
  clang
  clang++

Rust:
  cargo check
  cargo build
  cargo test
```

Secondary orchestrator handling includes, where practical:

```text
make
ninja
cmake --build
```

---

# 18. Diagnostic Model

Build failures are represented structurally.

Concepts include:

```text
tool
severity
diagnostic code
file
line
column
message
notes
help
raw range
```

High-value information is prioritized:

```text
errors
associated notes
compiler suggestions
warnings
```

Repeated warning/error chains can be grouped conservatively.

---

# 19. Compiler Noise Reduction

Phase 5 targets common token-heavy compiler output such as:

```text
C++ template instantiation chains
include stacks
repeated warnings
caret/source blocks
linker errors
Rust note/help chains
```

Important diagnostics remain recoverable through drill-down.

The parser does not attempt semantic code analysis.

It only structures tool-provided diagnostics.

---

# 20. Test Result Model

Test output is represented around:

```text
pass
fail
skip/ignored
packages/suites
failing tests
failure detail
panic/stack traces
```

Successful test runs become extremely compact.

Example concept:

```text
@acap <id> go-test PASS
packages=42 duration=8.1s
```

Failure output emphasizes failing tests rather than thousands of passing tests.

---

# 21. Build/Test Stateful Delta

The major Phase 5 value comes from transitions such as:

```text
7 errors
   ->
3 errors
   ->
1 error
   ->
PASS
```

Rather than sending every diagnostic set repeatedly, AgentCap can represent:

```text
resolved
new
remaining
```

Likewise for tests:

```text
5 failed
   ->
2 failed
   ->
PASS
```

New failures and regressions are prioritized.

---

# 22. Build/Test Drill-Down

Supported result-specific retrieval includes concepts such as:

```bash
acap show <id> --errors
acap show <id> --warnings
acap show <id> --file src/parser.cpp
acap show <id> --test TestParseInvalidToken
```

These operate entirely on stored command results.

They never re-run the compiler or test command.

---

# 23. Phase 6 — Coding Agent Integration

Status:

```text
COMPLETE
```

Purpose:

> Move AgentCap from an explicitly invoked CLI wrapper into the normal coding-agent shell/tool execution path.

Before Phase 6:

```text
agent
 -> acap run git diff
```

After Phase 6:

```text
agent
 -> git diff
 -> AgentCap integration
 -> compact Git result
```

The model can continue using ordinary shell commands.

---

# 24. Integration Architecture

Phase 6 keeps core logic agent-independent.

Architecture:

```text
Coding Agent
     |
     v
Agent Adapter
     |
     v
AgentCap Integration Boundary
     |
     v
AgentCap Core
     |
     +-- execution
     +-- reduction
     +-- storage
     +-- session
     +-- delta
     +-- rendering
```

Agent-specific behavior is isolated in adapters.

The core does not contain pervasive logic such as:

```text
if claude ...
if codex ...
```

---

# 25. Primary Agent Integrations

Phase 6 targeted at least:

```text
Claude Code
OpenAI Codex
```

with the architecture intended to support additional adapters later, such as:

```text
Gemini CLI
OpenCode
generic hook-based integrations
```

These secondary integrations are not required to redefine the core.

---

# 26. Critical Integration Invariant

The most important Phase 6 correctness rule is:

> Every intercepted command must execute exactly once.

AgentCap must never accidentally create:

```text
agent executes command
+
AgentCap executes command again
```

This matters especially for commands with side effects.

Examples:

```text
rm
mv
git add
git commit
go generate
database migrations
custom scripts
```

---

# 27. Integration Session Mapping

A coding-agent task/session maps to an AgentCap session.

This allows transparent execution to reuse Phase 3 functionality:

```text
normal command #1
normal command #2
normal command #3
```

becomes internally:

```text
same AgentCap session
```

while separate agent tasks remain isolated.

---

# 28. Transparent Result Flow

A normal integrated flow now looks like:

```text
Coding Agent
    |
    | ordinary shell command
    v
AgentCap interception
    |
    v
execute once
    |
    v
capture raw output
    |
    v
store in .acap
    |
    v
specialized reducer
    |
    v
session delta if applicable
    |
    v
compact result
    |
    v
Coding Agent
```

Result IDs remain visible so progressive disclosure still works.

---

# 29. Integration Safety

Phase 6 design includes:

- recursion protection
- AgentCap self-command bypass
- explicit bypass/debug mechanism
- cwd preservation
- environment preservation
- exit-code preservation
- cancellation propagation
- configuration idempotency
- safe integration removal
- fail-open behavior when compression fails but command output is available
- no remote telemetry
- no mandatory daemon

---

# 30. Current State

The roadmap through Phase 6 is now implemented.

Current maturity can be summarized as:

```text
Execution foundation                  COMPLETE
Generic reduction                     COMPLETE
Progressive disclosure                COMPLETE
Project-local persistent storage      COMPLETE
SQLite storage architecture           COMPLETE
Session-aware compression             COMPLETE
Cross-command delta presentation      COMPLETE
Git-aware compression                 COMPLETE
Git structural delta                  COMPLETE
Compiler/build-aware compression      COMPLETE
Test-aware compression                COMPLETE
Diagnostic/test delta                 COMPLETE
Transparent coding-agent integration  COMPLETE
```

AgentCap has therefore moved beyond being a simple output-filter CLI.

It is now conceptually:

> A stateful command-result virtualization layer between coding agents and their execution environment.

---

# 31. Current Architecture

The accumulated architecture is approximately:

```text
                    Coding Agent
                         |
                         v
                Integration Adapter
                         |
                         v
                  AgentCap Core
                         |
          +--------------+--------------+
          |                             |
          v                             v
       Executor                     Project Root
          |                             |
          v                             v
     Raw Capture                  .acap/store.db
          |                             |
          +----------+------------------+
                     |
                     v
                 Classifier
                     |
                     v
                  Reducer
                     |
      +--------------+--------------+
      |              |              |
      v              v              v
   Generic          Git        Build / Test
      |              |              |
      +--------------+--------------+
                     |
                     v
               Full Capsule
                     |
                     v
               Session Layer
                     |
                     v
              Baseline Selector
                     |
                     v
                Delta Engine
                     |
          +----------+----------+
          |          |          |
          v          v          v
      unchanged    delta     full fallback
          |
          v
        Agent
```

---

# 32. Storage Architecture at Current State

```text
project/
├── .git/
├── .acap/
│   ├── store.db
│   └── objects/
├── src/
├── tests/
└── ...
```

SQLite holds structured state such as:

```text
results
result metadata
capsules
statistics
Git indexes
diagnostics
test failures
schema versions
session relationships
baseline relationships
```

Raw stdout/stderr remain in the object store.

---

# 33. What AgentCap Is Not

Even after Phase 6, AgentCap is intentionally not:

```text
a coding agent
an autonomous agent
a shell
an AST engine
a repository semantic database
an LLM summarizer
a remote service
a telemetry service
a build system
a Git replacement
```

Its responsibility remains narrow:

> Mediate command-result information between execution tools and AI coding agents.

---

# 34. Important Existing Boundaries

Several boundaries should remain unless future measurement strongly justifies changing them.

## No LLM internally

Reduction remains deterministic.

## No command rewriting

AgentCap observes and transforms results, not user intent.

## No automatic reruns

Drill-down reads stored results.

## No destructive compression

Raw output remains available.

## No session reconstruction dependency

Each result remains complete independently.

## No cross-project result mixing

Storage and session context remain project-scoped.

---

# 35. Main Differentiation

AgentCap's current differentiation is no longer simply:

```text
"compress terminal output"
```

It is the combination of:

```text
full local capture
+
command-aware structural reduction
+
progressive disclosure
+
project-local result persistence
+
session-aware deduplication
+
structured delta presentation
+
transparent agent integration
```

The particularly important distinction is:

> AgentCap optimizes how much command-result information the model must see across an entire workflow, rather than only compressing each command independently.

---

# 36. Current Core KPI

The primary metric should remain:

```text
Total agent-visible tokens/bytes
over a complete coding workflow
```

not:

```text
compression ratio for one command
```

The calculation conceptually includes:

```text
initial command results
+
all acap show calls
+
all acap raw calls
+
integration instruction overhead
=
total AgentCap context cost
```

Compare this against:

```text
raw agent workflow output
```

and:

```text
stateless compression workflow output
```

---

# 37. Secondary KPIs

Useful secondary metrics include:

```text
raw bytes
capsule bytes
stateful returned bytes
drill-down bytes
raw retrieval bytes
unchanged-result rate
delta-result rate
full-fallback rate
information-recovery call rate
command latency overhead
adapter overhead
```

For integration quality:

```text
commands intercepted
commands bypassed
integration failures
```

are also useful.

---

# 38. Important Failure Metric

A particularly important metric remains:

```text
how often does the coding agent immediately request raw output?
```

If AgentCap achieves:

```text
95% initial reduction
```

but agents routinely issue:

```bash
acap raw <id>
```

immediately afterward, the reduction is not useful.

The product should optimize:

```text
useful compression
```

rather than:

```text
maximum compression
```

---

# 39. Current Major Strengths

At the end of Phase 6, the strongest parts of the design are:

### Lossless progressive disclosure

Large results can remain outside model context until required.

### Session-aware state

Repeated command results can be represented incrementally.

### Structured Git handling

Coding-agent Git loops can avoid repeatedly transmitting entire patches.

### Structured build/test handling

Repeated failure sets can be represented as changes rather than full logs.

### Transparent integration

The model does not need to remember to prefix every command with `acap run`.

### Local project scope

State naturally follows the repository without becoming a global cross-project cache.

---

# 40. Key Risks Going Forward

The next discussion should probably focus less on adding more commands and more on validating the product model.

Important risks include:

## Over-compression

Does AgentCap sometimes hide information that the coding agent actually needs?

## Excessive drill-down

Are agents compensating by calling `show` or `raw` too frequently?

## Session correctness

Are baselines always selected appropriately?

## Delta usability

Are deltas actually easier for models to reason about than complete capsules?

## Latency

Does interception and structured parsing materially slow interactive agent workflows?

## Storage growth

How quickly does `.acap` grow during long-running development?

## Integration brittleness

How stable are coding-agent hooks across tool versions?

---

# 41. Next Roadmap Decision

Phase 0–6 mainly built capability.

The next stage should likely be driven by measurement.

A useful fork is:

```text
Phase 7A: Benchmark / evaluation framework
```

versus:

```text
Phase 7B: Semantic source-code integration
```

The recommendation is to strongly consider doing benchmark/evaluation work first.

AgentCap now has enough functionality that additional complexity should be justified by measured agent behavior.

---

# 42. Candidate Phase 7 — Evaluation and Optimization

A logical next phase would be:

```text
Phase 7 — Agent Workflow Evaluation
```

Possible goals:

- establish reproducible coding-agent benchmarks
- compare AgentCap on/off
- measure total tokens
- measure task completion
- measure command count
- measure drill-down behavior
- measure latency
- measure accidental information loss
- identify reducers that are too aggressive
- identify commands responsible for most context usage

This would answer whether AgentCap is actually helping the agent rather than merely producing smaller terminal output.

---

# 43. Candidate Benchmark Design

A useful benchmark matrix could compare:

```text
A. AgentCap disabled

B. AgentCap stateless
   compression only

C. AgentCap stateful
   sessions + deltas

D. AgentCap fully integrated
   transparent Phase 6 mode
```

Across coding tasks such as:

```text
bug fixing
test failure repair
small feature implementation
refactoring
compile-error repair
Git review / diff inspection
large repository search
```

---

# 44. Candidate Metrics for Benchmarking

For each task measure:

```text
task success
total model input tokens
total tool-output tokens
wall-clock time
number of shell commands
number of repeated commands
number of acap show calls
number of acap raw calls
raw-output fallback rate
AgentCap processing latency
incorrect/missing-context incidents
```

This should provide much stronger evidence than standalone reducer benchmarks.

---

# 45. Candidate Phase 7 Alternative — Semantic Source Navigation

Another possible direction is deeper source-aware handling.

For example:

```text
cat large source file
   ->
symbol outline

acap show <id> --symbol Resolve
```

This could eventually integrate AST-based source structure.

However, this adds significant complexity.

It should probably be justified by benchmark evidence showing that source-file output is a major remaining context cost.

---

# 46. Potential AST Tool Relationship

If semantic source navigation becomes necessary, AgentCap should preferably not grow its own parser stack unnecessarily.

A clean long-term division would be:

```text
AgentCap
  -> command/result transport
  -> context compression
  -> progressive disclosure

AST Tool
  -> source structure
  -> symbols
  -> semantic navigation
  -> workspace analysis
```

AgentCap could consume semantic services rather than duplicate them.

This keeps the two projects conceptually separate.

---

# 47. Other Possible Future Areas

After measurement, potential future work includes:

```text
additional test runners
additional compilers
Docker / kubectl output
package managers
lint/static-analysis tools
coverage output
benchmark output
repository semantic source views
adaptive presentation budgets
automatic reducer tuning
```

These should not be added merely for breadth.

Prioritize commands that benchmark data shows are major context consumers.

---

# 48. Current Recommended Position

At the end of Phase 6:

```text
AgentCap core architecture is sufficiently complete for serious evaluation.
```

The project already has:

```text
execution
capture
storage
compression
recovery
state
delta
domain-aware parsing
transparent integration
```

The next major question is therefore no longer:

> What feature should AgentCap implement next?

It is:

> Does AgentCap materially improve coding-agent efficiency and task performance across realistic workflows, and where does it still waste context?

That should drive the next roadmap.
