# AgentCap Phase 5 — Build, Compiler, and Test-Aware Compression

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
- Phase 4: Git-aware compression and structural Git deltas

Phase 5 introduces specialized handling for build systems, compilers, and test runners.

The central goal is:

> Preserve actionable failures and diagnostics while aggressively reducing successful, repetitive, and low-value build/test output.

This phase should make AgentCap especially useful during the most common coding-agent loop:

```text
edit
build
inspect error
edit
build
test
inspect failure
edit
test
```

Phase 5 must remain local, deterministic, and conservative.

Do not use an LLM to summarize compiler or test output.

---

# 1. Phase 5 Goal

Implement specialized reducers for the initial build/test toolchain set:

```text
Go:
  go build
  go test

C / C++:
  gcc
  g++
  clang
  clang++

Rust:
  cargo build
  cargo check
  cargo test
```

Where practical, also recognize build orchestrators whose output mainly contains compiler diagnostics:

```text
make
cmake --build
ninja
```

However, these orchestrators are secondary.

The priority order is:

```text
1. go test
2. go build
3. gcc / g++
4. clang / clang++
5. cargo check / build
6. cargo test
7. make / ninja / cmake --build
```

Do not attempt to support every language or build system in Phase 5.

---

# 2. Core Product Principle

Build and test output contains highly uneven information value.

For example:

```text
hundreds of successful compilation lines
thousands of passing tests
progress output
repeated template diagnostics
repeated stack traces
one actual root-cause compiler error
```

A coding agent usually needs:

```text
what failed?
where?
why?
what changed since the previous run?
```

Therefore Phase 5 should transform:

```text
large execution log
        |
        v
structured diagnostics
        |
        v
compact failure/success capsule
```

while preserving the complete raw output in the Phase 2.1 object store.

---

# 3. Do Not Rewrite Commands

AgentCap must execute the exact requested command.

Do not silently transform:

```bash
go test ./...
```

into:

```bash
go test -json ./...
```

Do not rewrite:

```bash
cargo test
```

into another Cargo invocation.

Do not add compiler flags merely to simplify parsing.

The requested command remains authoritative.

AgentCap parses the captured output afterward.

Additional hidden build/test subprocesses are not allowed unless there is an exceptional, documented reason.

---

# 4. Build/Test Classification

Add structured classification for supported command families.

Examples:

```text
go test ./...
go build ./...

gcc ...
g++ ...
clang ...
clang++ ...

cargo check
cargo build
cargo test
```

Handle resolved executable paths where practical:

```text
/usr/bin/gcc
/usr/local/bin/clang++
```

Do not classify an unrelated executable merely because its filename contains `gcc`, `go`, or `cargo`.

---

# 5. Architecture

Introduce build/test-specific reducers outside the generic reducer package.

A reasonable structure:

```text
internal/
  reduce/
    build/
      common.go
      go_build.go
      gcc.go
      clang.go
      cargo_build.go

    test/
      common.go
      go_test.go
      cargo_test.go
```

Or equivalent.

Prefer structured intermediate representations over rendering directly from raw text.

---

# 6. Common Diagnostic Model

Introduce a normalized internal diagnostic representation.

For example:

```go
type Diagnostic struct {
    Tool       string
    Severity   Severity
    Message    string

    File       string
    Line       int
    Column     int

    Code       string

    Notes      []DiagnosticNote

    RawStart   int64
    RawEnd     int64
}
```

Exact fields may differ.

Useful severity categories include:

```text
error
warning
note
help
info
```

Do not force all tools into fields they do not provide.

---

# 7. Test Result Model

Use a separate structured representation for tests.

For example:

```go
type TestRun struct {
    Passed  int
    Failed  int
    Skipped int

    Packages []TestPackage
    Failures []TestFailure
}
```

A failure may contain:

```go
type TestFailure struct {
    Package  string
    Test     string
    Message  string
    File     string
    Line     int
    Output   string
}
```

Do not mix compiler diagnostics and test failures into one opaque type.

---

# 8. Raw Output Remains Authoritative

All build/test output must continue to be stored exactly through the existing result store.

The specialized parser is an optimization layer.

These must remain available:

```bash
acap show <id>
acap raw <id>
```

If parsing fails, the raw command result must still be complete and accessible.

---

# 9. Successful Build Output

Successful builds should normally produce extremely compact output.

Example:

```bash
acap run go build ./...
```

may return:

```text
@acap 91bf20 go-build PASS
duration=4.2s
```

If there were warnings:

```text
@acap 91bf20 go-build PASS warnings=3

W src/foo.go:81 unused ...
W ...
```

Do not preserve hundreds of routine successful build lines unless specifically requested.

---

# 10. Failed Build Output

A failed build should emphasize actionable diagnostics.

Example:

```text
@acap a81c20 clang-build FAIL
errors=3 warnings=12

E src/parser.cpp:481:17
  no matching function for call to 'parse'

E src/workspace.cpp:918:12
  use of undeclared identifier 'ctx'

E src/index.cpp:201:9
  cannot convert ...

warnings=12 grouped
```

Preserve:

- source path
- line
- column
- primary message
- diagnostic code where available
- important notes
- template/include context only when needed

---

# 11. Diagnostic Ordering

Preserve the original meaningful diagnostic order.

The first compiler error is often particularly important.

Do not sort diagnostics alphabetically by file unless the tool naturally groups them that way.

Do not move warnings ahead of errors.

Recommended output order:

```text
fatal errors
errors
important notes associated with errors
warnings
```

while preserving relative order within the original diagnostic stream where practical.

---

# 12. Error Priority

Errors are more important than warnings.

When output is large:

```text
all errors should generally remain visible
```

up to a reasonable safety limit.

Warnings may be summarized or grouped more aggressively.

Do not drop errors merely to hit an arbitrary compression ratio.

---

# 13. Diagnostic Limits

Introduce conservative limits for pathological builds.

For example:

```text
errors=824 shown=50 omitted=774
warnings=2201 shown=20 omitted=2181
```

This is preferable to returning megabytes of diagnostics.

The limit should be configurable internally and easy to tune.

The result must clearly state omitted counts.

---

# 14. Warning Grouping

Repeated warnings should be grouped where safe.

Example:

```text
warning: unused variable 'x'
```

appearing in 42 files may become:

```text
W unused variable 'x' x42
sample:
  src/a.cpp:81
  src/b.cpp:91
  src/c.cpp:14
```

Do not merge warnings whose messages only superficially resemble one another.

Use deterministic normalized diagnostic identity.

---

# 15. Error Deduplication

Compilers, especially C++, may emit repeated diagnostic chains.

AgentCap may collapse identical diagnostics when:

```text
same severity
same primary message
same logical source location
same diagnostic code where available
```

Be conservative.

Do not collapse separate errors merely because they share the same first sentence.

---

# 16. Diagnostic Chains

Preserve meaningful relationships such as:

```text
error
  note
  note
```

or:

```text
error
  instantiated from ...
  required from ...
```

Do not detach notes from their parent diagnostic when the relationship is obvious.

Internally model a primary diagnostic plus associated notes where practical.

---

# 17. C++ Template Errors

C++ template errors are a major token source.

Phase 5 should reduce repeated instantiation context while preserving the useful root cause.

For example, instead of returning 300 lines:

```text
required from ...
required from ...
required from ...
...
no matching function ...
```

return something like:

```text
E src/foo.cpp:91
  no matching function for call to make_node(...)

instantiation chain: 8 frames
  include/foo.hpp:181
  src/parser.cpp:410
  ...
```

Allow the agent to drill down if it needs the full diagnostic chain.

Do not attempt semantic interpretation of C++ templates.

---

# 18. Include Stack Compression

GCC/Clang may print long include stacks.

Compress repeated include prefixes.

Example:

```text
included-from chain=14
  src/main.cpp
  include/workspace.hpp
  include/parser.hpp
  ...
```

Preserve enough of the chain to identify relevant project headers.

Do not blindly remove all include context.

---

# 19. System Header Noise

Diagnostics originating in system or standard-library headers may be noisy.

Do not automatically suppress them.

However, when the compiler clearly points back to a project call site, prefer presenting the project-facing diagnostic first.

Do not classify paths as unimportant merely because they contain `/usr/include`.

Compiler output remains authoritative.

---

# 20. GCC Reducer

Support common GCC diagnostic formats such as:

```text
file.cpp:12:8: error: ...
file.cpp:14:5: warning: ...
```

Parse:

```text
file
line
column
severity
message
```

Also preserve:

```text
note:
fatal error:
```

where present.

Do not depend on color output being enabled or disabled.

Use ANSI-normalized parser input while keeping raw bytes unchanged.

---

# 21. Clang Reducer

Support standard Clang diagnostics.

Clang often includes:

```text
source line
caret indicator
fix-it hints
notes
```

Preserve the source/caret snippet when it materially helps explain the diagnostic.

Do not return every repeated caret block if many diagnostics share the same source region.

---

# 22. Fix-It Hints

If Clang provides a concise fix-it or `did you mean` message, preserve it.

Example:

```text
help: did you mean 'resolveNode'?
```

This is highly actionable for a coding agent.

Do not invent fixes.

Only surface compiler-provided suggestions.

---

# 23. Go Build Reducer

Support:

```bash
go build
go build ./...
```

Go compiler diagnostics are relatively compact, but repository-wide builds may repeat package context.

Example:

```text
@acap 8bc210 go-build FAIL
packages_failed=2

E internal/parser/parser.go:81:14
  undefined: parseNode

E internal/index/index.go:144:9
  cannot use result (...) as ...
```

Preserve package names where useful.

---

# 24. Go Test Reducer

`go test` is one of the highest-priority Phase 5 targets.

Support common output from:

```bash
go test
go test ./...
go test -v ./...
```

Identify:

```text
packages tested
packages passed
packages failed
test failures
panic/fatal failures
```

Example:

```text
@acap f12a91 go-test FAIL
packages=42 pass=40 fail=2
tests_failed=5

internal/parser
  TestParseInvalidToken
  TestParseTemplate
  TestParseNested

internal/index
  TestLookupMissing
  TestLookupAmbiguous
```

---

# 25. Go Test Failure Detail

Include concise failure details.

Example:

```text
TestParseInvalidToken
  parser_test.go:418
  expected TokenIdent, got TokenError
```

If failure output is long:

```text
failure_output_lines=184 shown=30
```

Preserve the beginning and relevant assertion/error region where detectable.

Raw output remains available.

---

# 26. Go Test Success

Successful test runs should be extremely compact.

Example:

```text
@acap 81c912 go-test PASS
packages=42 duration=8.4s
```

If test counts cannot be safely derived from normal output, do not invent them.

Prefer:

```text
packages=42
```

over a fabricated test count.

---

# 27. `go test -v`

Verbose Go tests can produce enormous amounts of successful log output.

On successful tests:

- remove routine `=== RUN`
- remove routine `--- PASS`
- collapse successful per-test output

Preserve:

- failing tests
- output associated with failure
- package summaries
- panics
- unexpected logs that the parser cannot safely classify

Be conservative with arbitrary user logging.

---

# 28. Panic Handling

Panics are high-priority diagnostics.

Preserve:

```text
panic message
relevant stack frames
test/package identity
```

Avoid dumping hundreds of runtime/internal frames when they add little value.

For Go stack traces, prefer project frames when safely identifiable, but retain enough context to avoid misleading the agent.

---

# 29. Stack Trace Compression

For any supported test runner that emits long stack traces:

```text
preserve top/root relevant frames
collapse obvious repetitive/internal frames
report omitted frame count
```

Example:

```text
stack:
  internal/parser.Parse(...)
  internal/parser.TestParse(...)
  ...
runtime/internal frames omitted=18
```

Only categorize frames as internal when the format is reliably recognized.

---

# 30. Cargo Build / Check

Support:

```bash
cargo check
cargo build
```

Rust diagnostics commonly include:

```text
error[E....]
warning:
 --> src/foo.rs:...
 help:
 note:
```

Preserve:

- error code
- file
- line/column
- primary message
- help
- relevant notes

Example:

```text
@acap a9120d cargo-check FAIL
errors=4 warnings=11

E[E0382] src/foo.rs:81:17
  borrow of moved value: `value`

E[E0277] src/bar.rs:91:9
  trait bound `Foo: Bar` is not satisfied
```

---

# 31. Rust Error Codes

Preserve Rust compiler error codes such as:

```text
E0382
E0277
```

These are compact and highly useful.

Do not remove them during normalization.

---

# 32. Rust Help and Notes

Rust compiler diagnostics often include useful `help:` suggestions.

Keep concise compiler-generated help.

Example:

```text
help: consider borrowing here
```

Do not promote every note to a separate top-level error.

Associate notes/help with the primary diagnostic.

---

# 33. Cargo Test

Support common:

```bash
cargo test
cargo test --workspace
```

Extract:

```text
test result
passed
failed
ignored
measured
filtered out
failing test names
failure details
```

Example:

```text
@acap f1c20a cargo-test FAIL
pass=1241 fail=3 ignored=12

failures:
parser::tests::invalid_token
workspace::tests::resolve_missing
index::tests::lookup_duplicate
```

---

# 34. Cargo Test Failure Detail

Preserve per-test failure sections.

Example:

```text
workspace::tests::resolve_missing
  src/workspace.rs:418
  assertion failed: result.is_some()
```

Avoid repeating the entire Cargo summary after every failure.

---

# 35. Multiple Test Binaries

Cargo may execute multiple test binaries.

Preserve enough grouping to distinguish them.

Example:

```text
lib agentcap_core
  240 pass
  2 fail

tests integration
  81 pass
  1 fail
```

Do not assume all test names are globally unique.

---

# 36. Build Orchestrators

Secondary support may include:

```text
make
ninja
cmake --build
```

These tools often interleave:

```text
build progress
compiler invocation
compiler diagnostics
linker output
```

The initial goal is not to fully model the build graph.

Instead:

```text
remove obvious progress noise
surface recognized compiler/linker diagnostics
preserve failed command information
```

If parsing confidence is low, use generic reduction.

---

# 37. Make

For `make`, preserve:

```text
failed target
failed command where available
compiler diagnostics
final make error
```

Suppress repetitive successful target lines where safe.

Example:

```text
@acap 1cb822 make FAIL
target=all

E src/parser.cpp:481 ...
E src/workspace.cpp:918 ...

make: *** [parser.o] Error 1
```

---

# 38. Ninja

Ninja progress lines such as:

```text
[127/842] Building CXX object ...
```

may be heavily reduced.

On success:

```text
@acap ... ninja PASS steps=842
```

where step count is reliably available.

On failure, preserve:

```text
FAILED:
failed command
diagnostics
```

---

# 39. CMake Build

Support only:

```bash
cmake --build ...
```

if practical.

Do not attempt to parse CMake configure/generate semantics deeply in Phase 5.

Treat the output primarily as a build orchestrator stream.

---

# 40. Linker Errors

Support common linker failures where possible.

Examples:

```text
undefined reference to ...
duplicate symbol ...
ld: error: ...
```

Represent them as high-priority errors.

Example:

```text
E linker
  undefined reference to `Workspace::Resolve(...)'
  referenced by src/main.o
```

Do not misclassify linker errors as compiler warnings.

---

# 41. Build Command Failure Without Diagnostics

Sometimes a build command exits non-zero without recognized diagnostics.

In that case return:

```text
@acap ... build FAIL exit=1

unparsed output:
[first useful lines]
...
[last useful lines]
```

Do not return:

```text
FAIL errors=0
```

as though the parser understood the cause.

Clearly distinguish parsed diagnostics from unparsed failure output.

---

# 42. Success Detection

A command should only be marked:

```text
PASS
```

when its actual process exit code indicates success.

Never infer success solely from parsed text.

Likewise, non-zero exit is always failure at the command level even if no parser error was detected.

Exit status remains authoritative.

---

# 43. Failure Detection

Do not infer command failure only from seeing the word:

```text
error
```

Some tools may print expected errors during successful tests.

Use:

```text
process exit code
+
tool-specific structure
```

to determine command result state.

---

# 44. Phase 3 Session Integration

Build/test reducers must integrate with session-aware compression.

The baseline rules remain:

```text
same session
same cwd
same argv
latest equivalent previous result
```

Do not weaken these rules for build/test commands.

---

# 45. Build Delta Principle

Repeated build runs should emphasize diagnostic changes.

Example:

First:

```text
errors=7
```

After edits:

```text
errors=3
```

AgentCap should ideally return:

```text
@acap ... delta from ...
build FAIL
errors 7 -> 3

resolved=5
new=1

new:
E src/index.cpp:201 ...

remaining=2
```

rather than re-sending every unchanged error.

---

# 46. Diagnostic Identity for Delta

Use a conservative normalized diagnostic identity.

Potential components:

```text
tool
severity
diagnostic code
file
line
column
normalized primary message
```

Be careful with line-number shifts.

For Phase 5, exact or near-exact structural identity is sufficient.

Do not implement semantic diagnostic matching with embeddings.

---

# 47. Line Number Changes

Editing code can shift an otherwise identical diagnostic from:

```text
foo.cpp:81
```

to:

```text
foo.cpp:84
```

Do not aggressively assume these are the same error unless the matching rule is highly reliable.

It is acceptable for Phase 5 to report it as:

```text
old diagnostic resolved
new diagnostic added
```

Correctness is more important than maximizing deduplication.

---

# 48. Resolved Diagnostics

If a previous build had:

```text
E src/foo.cpp:81 ...
```

and the new equivalent build no longer contains it, the delta may report:

```text
resolved:
src/foo.cpp:81 ...
```

This is useful feedback to the coding agent.

Keep resolved diagnostics compact.

---

# 49. New Diagnostics

New errors should be surfaced prominently.

Example:

```text
new_errors=2

E src/foo.cpp:...
E src/bar.cpp:...
```

Do not bury them beneath a count of resolved warnings.

---

# 50. Error Count Transition

Always surface transitions such as:

```text
errors 8 -> 3
errors 0 -> 1
exit 0 -> 1
exit 1 -> 0
```

These are highly useful signals.

---

# 51. Success After Failure

A build/test transition from failure to success can be extremely compact.

Example:

```text
@acap a81c92 delta from c910ab
go-test FAIL -> PASS
resolved_failures=5
duration=7.8s
```

There is no need to resend prior failure detail unless explicitly requested.

---

# 52. Repeated Identical Failure

If exactly the same build/test result repeats:

```text
@acap ... unchanged from ...
exit=1
```

Phase 3 exact deduplication should apply.

Do not generate a verbose diagnostic capsule again.

---

# 53. Test Delta

Test-specific deltas should compare:

```text
failed test identities
pass/fail summary
```

Example:

Previous:

```text
5 failures
```

Current:

```text
2 failures
```

Return:

```text
@acap ... go-test delta
failed 5 -> 2

resolved:
TestParseA
TestParseB
TestLookupC

remaining:
TestResolveX
TestResolveY
```

---

# 54. Newly Failing Tests

If a previously passing or absent test now fails:

```text
new:
TestWorkspaceCache
```

include its failure detail.

New regressions are high priority.

---

# 55. Test Output Identity

Use the test framework's native test name as the primary identity.

Where necessary include:

```text
package/module/test binary
```

to avoid collisions.

Example:

```text
internal/parser::TestResolve
```

rather than only:

```text
TestResolve
```

when multiple packages may contain the same name.

---

# 56. Flaky Output

Do not attempt to determine whether a test is flaky.

AgentCap may report:

```text
previously failed -> now passed
```

but must not infer causality or test reliability.

---

# 57. Drill-Down: Diagnostics

Extend `acap show` with useful selectors.

For example:

```bash
acap show <id> --errors
acap show <id> --warnings
```

Optional:

```bash
acap show <id> --diagnostic 3
```

This should retrieve the selected stored diagnostic context without re-running the build.

---

# 58. Drill-Down: File

Reuse the existing file/path selector where appropriate:

```bash
acap show <id> --file src/parser.cpp
```

For a build result, return diagnostics associated with that file.

Example:

```text
@acap ... file=src/parser.cpp
errors=3 warnings=2

...
```

Do not confuse this with Phase 4 Git file drill-down; route behavior according to result type.

---

# 59. Drill-Down: Test

Support a test selector where practical:

```bash
acap show <id> --test TestParseInvalidToken
```

For Cargo, allow fully qualified names:

```bash
acap show <id> --test parser::tests::invalid_token
```

Return stored failure detail only.

Do not re-run the test.

---

# 60. Drill-Down Must Use Stored Results

All Phase 5 drill-down operations must use the captured/stored command result.

Do not implement:

```text
--test X -> rerun test X
```

or:

```text
--file X -> rerun compiler
```

AgentCap retrieval and command execution must remain separate concepts.

---

# 61. SQLite Structured Metadata

Persist lightweight parsed build/test metadata where useful.

Possible tables:

```text
diagnostics
test_runs
test_failures
```

Do not store enormous raw diagnostic bodies redundantly if they already exist in raw object storage.

---

# 62. Suggested Diagnostic Schema

A possible schema:

```sql
CREATE TABLE diagnostics (
    result_id TEXT NOT NULL,
    diagnostic_index INTEGER NOT NULL,

    tool TEXT NOT NULL,
    severity TEXT NOT NULL,
    code TEXT,

    file TEXT,
    line INTEGER,
    column_no INTEGER,

    message TEXT NOT NULL,

    raw_start INTEGER,
    raw_end INTEGER,

    PRIMARY KEY (result_id, diagnostic_index)
);
```

Exact schema may differ.

Use the Phase 2.1 migration mechanism.

---

# 63. Suggested Test Failure Schema

For example:

```sql
CREATE TABLE test_failures (
    result_id TEXT NOT NULL,
    failure_index INTEGER NOT NULL,

    suite TEXT,
    test_name TEXT NOT NULL,
    file TEXT,
    line INTEGER,

    raw_start INTEGER,
    raw_end INTEGER,

    PRIMARY KEY (result_id, failure_index)
);
```

Do not over-normalize the database.

Only persist what materially improves drill-down and Phase 3 comparison.

---

# 64. Raw Offsets

Where possible, store byte offsets or equivalent references into raw stdout/stderr.

This enables efficient retrieval of:

```text
one diagnostic
one failure
one stack trace
```

without loading the full result.

Use whichever offset model is already established by Phase 4.

Avoid creating a second incompatible indexing approach.

---

# 65. stdout vs stderr

Compiler and test tools vary in their stream usage.

Some compilers primarily emit diagnostics to stderr.

Some test runners mix stdout and stderr.

Do not assume:

```text
stderr = failure
stdout = success
```

Parse both streams where necessary.

Keep their raw storage separate.

---

# 66. Cross-Stream Diagnostics

If meaningful diagnostic context spans stdout/stderr, do not invent exact cross-stream ordering unless AgentCap already captures it reliably.

It is acceptable to parse streams independently and present:

```text
stderr diagnostics
stdout test summary
```

The raw streams remain authoritative.

---

# 67. ANSI and Progress Cleanup

Use normalized parser input.

Examples of low-value noise:

```text
spinner frames
progress bars
repeated percentages
ANSI colors
terminal cursor updates
```

Do not modify the stored raw objects.

Maintain:

```text
raw = exact capture
normalized = parser input
```

---

# 68. Parser Confidence and Fallback

If a tool-specific parser cannot confidently interpret output:

```text
use generic/current reducer behavior
```

Do not emit fabricated counts or structured fields.

Examples:

```text
unknown compiler diagnostic format
custom test harness
localized output format
unsupported Cargo output variant
```

Conservative fallback is mandatory.

---

# 69. Locale Variability

Compiler output may be localized.

Do not assume all installations emit English diagnostics unless the tool behavior guarantees it.

Do not silently alter locale environment variables merely to simplify parsing.

If the parser depends on English phrases and encounters unknown output:

```text
fallback
```

rather than misparse.

---

# 70. No Hidden Environment Mutation

Do not silently set:

```text
LC_ALL=C
LANG=C
RUSTFLAGS=...
CFLAGS=...
```

to simplify parsing.

That could alter command behavior.

The command must run under the environment requested by the user/agent.

---

# 71. No JSON Mode Rewriting

Some tools offer machine-readable modes.

Examples include:

```text
go test -json
cargo --message-format=json
```

Do not automatically enable them.

If the user explicitly invoked a structured-output mode, AgentCap may parse it.

But AgentCap must not alter argv to obtain structured output.

---

# 72. User-Requested JSON Output

If the actual command outputs JSON, preserve raw output and consider using the appropriate structured parser if straightforward.

Do not pretty-print or expand JSON unnecessarily.

AgentCap's objective remains lower agent-visible token cost.

---

# 73. Large Build Logs

A build may produce hundreds of megabytes.

Avoid:

```text
entire log -> []byte
```

where practical.

Reuse the streaming object-store architecture.

Parsers should ideally operate:

```text
streaming
or
bounded-buffer
or
temporary-file-backed
```

for large results.

Document any in-memory parser limits.

---

# 74. Maximum Parser Memory

Introduce a safe parser strategy for extreme results.

For example:

```text
parse first N MB + indexed relevant sections
```

or a streaming state machine.

If the parser limit is exceeded:

```text
fall back conservatively
```

Do not crash AgentCap because a compiler emitted pathological output.

---

# 75. Truncated Diagnostic Parsing

If AgentCap cannot fully parse an extreme result, clearly indicate that parsing was partial.

Example:

```text
diagnostics>=500 shown=50
parser_limit_reached=true
```

Do not claim an exact total if the entire result was not processed.

---

# 76. Benchmark Fixtures

Create deterministic fixtures such as:

```text
testdata/build/
  gcc-single-error.txt
  gcc-many-errors.txt
  gcc-template-error.txt
  clang-error.txt
  clang-warning.txt
  linker-errors.txt

  go-build-error.txt
  go-test-pass.txt
  go-test-fail.txt
  go-test-verbose.txt
  go-test-panic.txt

  cargo-check-error.txt
  cargo-test-pass.txt
  cargo-test-fail.txt
```

Avoid giant committed fixtures.

Generate large synthetic variants in tests/benchmarks when needed.

---

# 77. Integration Tests: Go

Where Go is available, create temporary modules.

Test:

```text
successful build
compile failure
successful test
single failed test
multiple failed tests
panic
verbose test output
multiple packages
```

Since AgentCap itself is written in Go, these integration tests are high value.

---

# 78. Integration Tests: C/C++

Where GCC/Clang are available, compile small deterministic source files.

Test:

```text
syntax error
type error
warning
multiple errors
template failure
linker failure
successful build
```

Skip tool-specific tests cleanly when the compiler is unavailable.

---

# 79. Integration Tests: Rust

If Cargo is available, create a small temporary crate.

Test:

```text
cargo check pass
cargo check failure
cargo test pass
cargo test failure
```

Do not make Rust availability mandatory for the entire Go test suite.

---

# 80. Session Workflow Tests

Create realistic coding-agent workflows.

Example:

```text
1. go test ./...
   -> 5 failures

2. modify fixture

3. go test ./...
   -> 2 failures

4. modify fixture

5. go test ./...
   -> PASS

6. go test ./...
   -> unchanged PASS
```

Verify output transitions at every step.

---

# 81. Compiler Workflow Test

Example:

```text
build #1
  errors=7

fix several errors

build #2
  errors=3

fix remaining errors

build #3
  PASS
```

Measure total returned bytes.

This should demonstrate the benefit of structural diagnostic deltas.

---

# 82. New Regression Test

Test:

```text
build #1 -> 2 errors
build #2 -> 1 old error resolved, 1 new error introduced
```

The delta must prominently show the new error.

Do not output only:

```text
errors 2 -> 2
```

because that hides the actual state change.

---

# 83. Diagnostic Reordering Test

Some builds may output the same diagnostics in a slightly different order.

Phase 5 may normalize diagnostic identity enough to avoid reporting everything as changed where safe.

However, do not prioritize order-insensitive matching over correctness.

Exact structural matching is sufficient initially.

---

# 84. Phase 5 Statistics

Extend statistics with build/test-relevant measurements.

Potential metrics:

```text
build/test commands
raw diagnostic bytes
stateless capsule bytes
stateful returned bytes
diagnostics parsed
diagnostics omitted
repeated diagnostics collapsed
test failures parsed
drill-down bytes
```

Keep user-facing `acap stats` concise.

Detailed metrics may remain debug/benchmark-only.

---

# 85. Phase 5 KPI

Measure:

```text
raw build/test output
vs
structured Phase 5 capsule
vs
Phase 3 + Phase 5 stateful output
```

The most important metric remains:

> Total agent-visible bytes/tokens across a realistic edit-build-test workflow.

Do not optimize merely for one large log's initial compression percentage.

---

# 86. Information Recovery Metric

Track how often a coding-agent benchmark must call:

```text
acap raw
acap show --diagnostic
acap show --test
```

after receiving a Phase 5 capsule.

Extremely high compression is not useful if every result requires immediate raw retrieval.

Use this metric when tuning reducers.

---

# 87. Success Capsule Target

Successful common commands should generally fit into a few lines.

Examples:

```text
@acap 91aa12 go-build PASS duration=3.8s
```

```text
@acap 18bc92 go-test PASS packages=42 duration=8.1s
```

```text
@acap c91f21 cargo-test PASS pass=1284 ignored=7
```

Do not embellish successful output.

---

# 88. Failure Capsule Target

Failures should answer:

```text
what failed?
how many?
where?
what is the first/root actionable error?
what changed from the previous equivalent run?
```

before showing less useful context.

---

# 89. Do Not Infer Root Cause

AgentCap may prioritize the first diagnostic or preserve compiler grouping.

It must not claim:

```text
this is the root cause
```

unless the underlying tool explicitly says so.

Use neutral labels such as:

```text
first error
primary diagnostic
```

Do not perform semantic diagnosis.

---

# 90. No Automatic Fixing

Do NOT:

- edit source files
- apply compiler fix-its
- rerun commands
- rerun only failing tests
- invoke formatters
- install dependencies

Phase 5 only executes the requested command and transforms its captured result.

---

# 91. No AST Integration

Do not add:

```text
Tree-sitter
AST parsing
symbol resolution
semantic code analysis
```

to explain compiler errors.

Phase 5 diagnostics remain tool-output driven.

---

# 92. No Agent-Specific Integration

Do not add:

```text
Claude Code hooks
Codex hooks
Gemini CLI hooks
OpenCode hooks
MCP
shell interception
```

Those belong to Phase 6.

Phase 5 must be testable through:

```bash
acap run ...
```

alone.

---

# 93. No Additional Languages Yet

Do not add specialized reducers for:

```text
pytest
Jest
Vitest
Java/Javac
Maven
Gradle
dotnet
Swift
Ruby
PHP
```

during Phase 5 unless there is a trivial generic reuse and no scope increase.

The initial supported language/toolchain set is intentionally narrow.

---

# 94. SQLite Migration

Use the existing Phase 2.1 schema migration mechanism.

Do not recreate:

```text
.acap/store.db
```

to add Phase 5 metadata.

Add the minimum necessary tables/indexes.

Test migration from the Phase 4 schema.

---

# 95. Result Type

Add a reliable result classification such as:

```text
go_build
go_test
gcc
clang
cargo_build
cargo_test
```

or equivalent.

Do not overload reducer-name strings inconsistently across the codebase.

The result type should make drill-down routing predictable.

---

# 96. Specialized Delta Interface

Reuse the Phase 3 delta abstraction.

Do not create a second independent delta engine for build/test results.

Conceptually:

```text
DeltaReducer
  |
  +-- Git
  +-- diagnostics
  +-- test results
```

Keep the stateful presentation pipeline unified.

---

# 97. Delta Fallback

If diagnostic/test comparison is uncertain:

```text
return current full Phase 5 capsule
```

Do not emit a misleading reduced delta.

The same Phase 3 rule remains:

> Full current capsule is always safer than an incorrect delta.

---

# 98. Delta Cost Heuristic

Reuse the existing Phase 3 delta-size heuristic.

If:

```text
diagnostic delta size
```

is close to or larger than:

```text
full current build/test capsule
```

return the full current capsule.

Do not generate a giant "delta" for a completely changed failure set.

---

# 99. Debug Mode

Extend debug output with:

```text
detected tool
result type
parser selected
diagnostics parsed
test failures parsed
warnings grouped
parser fallback reason
baseline
delta type
full capsule size
delta size
presentation decision
```

Debug output goes to stderr.

Do not contaminate the normal agent-facing output.

---

# 100. Performance

Parsing should add little overhead compared with build/test execution.

Avoid:

```text
multiple complete rescans
large regex cascades
unnecessary SQLite writes per output line
```

Prefer:

```text
single parse pass
in-memory structured result for moderate logs
batched SQLite persistence
```

Do not insert every compiler line into SQLite individually.

---

# 101. Parsing Strategy

Prefer simple deterministic state machines over extremely complex regular expressions.

For example:

```text
detect diagnostic start
collect associated source snippet
collect notes/help
finalize diagnostic
```

This tends to be easier to maintain for GCC, Clang, and Rust output.

Use regex where appropriate for concise line formats.

---

# 102. Parser Isolation

Each parser should be independently testable.

For example:

```go
func ParseGCC(r io.Reader) (*BuildResult, error)
func ParseGoTest(r io.Reader) (*TestRun, error)
```

or equivalent.

Do not couple parsing to CLI rendering.

This is essential for reliable fixture-based testing.

---

# 103. Rendering Isolation

Structured data should be rendered separately.

Conceptually:

```text
raw output
   |
parser
   |
structured result
   |
renderer
   |
capsule
```

This allows:

```text
same structured result
 -> full capsule
 -> session delta
 -> drill-down
```

without reparsing formatting strings.

---

# 104. Suggested Implementation Order

Implement in this order:

1. define common diagnostic model
2. define common test-result model
3. add build/test command classification
4. implement Go build parser
5. implement Go test parser
6. add Go fixtures/integration tests
7. implement GCC parser
8. implement Clang parser
9. implement linker diagnostic handling
10. add C/C++ fixtures/integration tests
11. implement Cargo build/check parser
12. implement Cargo test parser
13. add Rust fixtures/integration tests
14. persist structured diagnostics/test metadata
15. implement diagnostic drill-down
16. implement test drill-down
17. integrate diagnostic session deltas
18. integrate test session deltas
19. implement warning/error grouping
20. add large-log limits/fallbacks
21. optionally add basic make/ninja/cmake-build handling
22. add workflow benchmarks
23. update README/documentation

Do not begin with every supported build system at once.

Make Go support excellent first, then C/C++, then Rust.

---

# 105. Definition of Done

Phase 5 is complete when all of the following are true:

1. AgentCap reliably classifies supported Go, C/C++, and Rust build/test commands.

2. Unsupported build/test commands safely fall back to existing generic behavior.

3. Successful Go builds return compact success capsules.

4. Failed Go builds expose actionable compiler diagnostics.

5. `go test` distinguishes successful and failing packages.

6. Failed Go tests expose test names and useful failure context.

7. Go panics remain visible and useful.

8. GCC diagnostics are parsed into structured errors/warnings/notes.

9. Clang diagnostics are parsed into structured errors/warnings/notes.

10. C++ template/include diagnostic noise can be reduced conservatively.

11. Common linker failures are surfaced.

12. Cargo build/check diagnostics preserve Rust error codes and help/notes.

13. Cargo test results expose failing test names and failure detail.

14. Exit code remains authoritative for command success/failure.

15. All errors are preserved up to a clearly reported safety limit.

16. Repetitive warnings/diagnostics are grouped only when safely identical.

17. Complete raw stdout/stderr remain independently recoverable.

18. `acap show <id> --errors` works for supported build results.

19. `acap show <id> --warnings` works where applicable.

20. File-based diagnostic drill-down works.

21. Test-based drill-down works for supported test runners.

22. Drill-down never re-runs the original command.

23. Structured diagnostic/test metadata is integrated with `.acap/store.db`.

24. Schema changes use the existing migration system.

25. Phase 3 exact unchanged detection applies to repeated build/test results.

26. Specialized diagnostic delta behavior exists.

27. Specialized test-failure delta behavior exists.

28. Newly introduced errors/failures are surfaced prominently.

29. Resolved errors/failures can be represented compactly.

30. Failure-to-success transitions are represented compactly.

31. Large deltas fall back to the full Phase 5 capsule.

32. Large logs do not cause uncontrolled memory consumption.

33. Parser failures fall back safely.

34. Representative fixture tests exist.

35. Integration tests exist for available Go/C/C++/Rust toolchains.

36. Workflow benchmarks quantify token/output savings.

37. `go test ./...` passes.

38. `go build ./...` passes.

39. No Phase 6 agent-specific integration has been implemented.

---

# 106. Engineering Priorities

When tradeoffs are necessary, use this order:

1. Never hide a new error or failing test
2. Preserve exact raw output
3. Preserve exit-code semantics
4. Correct diagnostic/test identity
5. Useful failure context
6. Conservative session deltas
7. Lower total agent-visible output
8. Low parsing overhead
9. Broader toolchain coverage

A longer but correct failure capsule is preferable to a shorter capsule that hides the information required to fix the code.

---

# 107. Phase 5 Success Criterion

Phase 5 succeeds when a typical coding-agent loop such as:

```text
go test ./...
edit
go test ./...
edit
go test ./...
```

or:

```text
cmake --build build
edit
cmake --build build
```

does not repeatedly inject the same successful logs and unchanged diagnostics into the model context.

The desired behavior is:

```text
first failure
  -> structured actionable failures

next run
  -> new/resolved/remaining failures only

successful run
  -> compact PASS transition

unchanged successful run
  -> tiny unchanged result
```

This is the core Phase 5 value proposition.

---

# 108. Final Verification

Before finishing Phase 5, manually verify representative workflows.

## Go test workflow

Start with several failing tests:

```bash
export ACAP_SESSION_ID=phase5-go

acap run go test ./...
```

Fix some failures and run again:

```bash
acap run go test ./...
```

Verify:

- resolved failures are represented
- remaining failures are preserved
- new failures are clearly surfaced

Fix all tests:

```bash
acap run go test ./...
```

Verify the failure-to-success transition is compact.

Run once more without changes:

```bash
acap run go test ./...
```

Verify unchanged detection applies.

## C/C++ workflow

Compile a fixture containing multiple errors:

```bash
acap run clang++ ...
```

Verify structured diagnostics.

Fix most errors and compile again.

Verify the session delta focuses on changed diagnostics.

Also test:

- warning-only success
- template error
- linker error
- executable success

## Rust workflow

Run:

```bash
acap run cargo check
acap run cargo test
```

Verify Rust error codes, test names, help text, and failure transitions.

## Drill-Down

For returned result IDs verify:

```bash
acap show <id> --errors
acap show <id> --warnings
acap show <id> --file src/foo.cpp
acap show <id> --test <test-name>
acap raw <id>
```

Verify all retrieval uses stored result data.

Finally run:

```bash
gofmt
go test ./...
go build ./...
```

Report:

- supported build/test commands
- internal diagnostic model
- internal test-result model
- SQLite schema additions
- diagnostic identity rules
- session delta rules
- drill-down operations
- fallback behavior
- large-log behavior
- workflow benchmark results
- known limitations

Do not proceed into Phase 6.

Stop once the Phase 5 Definition of Done is satisfied.
