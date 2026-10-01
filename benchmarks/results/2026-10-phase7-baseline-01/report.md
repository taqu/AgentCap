# AgentCap Phase 7 benchmark report

- Batches: 2026-10-phase7-baseline-01
- AgentCap versions: 2f2c593b4639+dirty
- Agents/models: antigravity: gemini-3.8-flash-medium (1.2.14), claude: claude-sonnet-5-5 (2.1.286)
- Tasks: 8 · Runs: 50 included, 0 excluded · Modes: disabled, integrated
- Primary metric: total agent-visible command-result bytes per workflow (capsules + deltas + `acap show` + `acap raw` + instruction overhead), counted once from the agent's own tool-result stream. Tokens marked *est* are bytes/4; host tokens are reported separately.

## antigravity

| mode | runs | pass | median context (IQR) | tokens est | median ctx, passing | savings vs OFF (all / passing) | median wall s | median shell cmds | repeated | host input tok (median) | native tool bytes (median) |
|---|---|---|---|---|---|---|---|---|---|---|---|
| OFF | 2 | 1/2 (50%) | 1388 B (800 B–1976 B) | 347 | 2563 B | – / – | 25 | 4 | 0 | 97199 | 20 B |

**OFF isolation:** 8 hook decisions all bypassed (0 intercepted).

### antigravity — per task (medians; OFF → FULL)

| task | pass OFF | pass FULL | context OFF | context FULL | Δ | cmds OFF→FULL | wall s OFF→FULL | show/raw FULL |
|---|---|---|---|---|---|---|---|---|

### antigravity — command families (total bytes over all runs)

| family | OFF bytes | OFF calls | FULL bytes | FULL calls | Δ |
|---|---|---|---|---|---|
| test | 1656 B | 3 | 0 B | 0 | -100% |
| git | 890 B | 3 | 0 B | 0 | -100% |
| listing | 204 B | 1 | 0 B | 0 | -100% |
| file-read | 26 B | 1 | 0 B | 0 | -100% |

## claude

| mode | runs | pass | median context (IQR) | tokens est | median ctx, passing | savings vs OFF (all / passing) | median wall s | median shell cmds | repeated | host input tok (median) | native tool bytes (median) |
|---|---|---|---|---|---|---|---|---|---|---|---|
| OFF | 24 | 24/24 (100%) | 3844 B (2992 B–4582 B) | 961 | 3844 B | – / – | 11 | 3 | 0 | 60524 | 0 B |
| FULL | 24 | 24/24 (100%) | 3836 B (3126 B–5256 B) | 959 | 3836 B | -0% / -0% | 11 | 3 | 0 | 63339 | 54 B |

**AgentCap behavior (FULL):** 81 results (full 81 / delta 0 / unchanged 0) · intercepted 81, bypassed 6, adapter failures 0 · `show` 0 calls (0 B, in 0 runs) · `raw` 0 calls (0 B, in 0 runs; immediate 0) · raw fallback rate 0% of results · targeted drill-down rate 0% of results

**Latency (FULL, per intercepted command):** adapter 0.11 ms · AgentCap core 1.12 ms

**OFF isolation:** 79 hook decisions all bypassed (0 intercepted).

### claude — per task (medians; OFF → FULL)

| task | pass OFF | pass FULL | context OFF | context FULL | Δ | cmds OFF→FULL | wall s OFF→FULL | show/raw FULL |
|---|---|---|---|---|---|---|---|---|
| c-bugfix-ringbuf | 3/3 | 3/3 | 4972 B | 4878 B | -2% | 4→4 | 11→10 | 0/0 |
| cpp-compile-repair | 3/3 | 3/3 | 3537 B | 3919 B | +11% | 3→3 | 12→12 | 0/0 |
| go-compile-repair | 3/3 | 3/3 | 3891 B | 3708 B | -5% | 3→3 | 10→12 | 0/0 |
| go-feature-csv | 3/3 | 3/3 | 2802 B | 2867 B | +2% | 2→2 | 11→11 | 0/0 |
| go-refactor-mathx | 3/3 | 3/3 | 2999 B | 3607 B | +20% | 5→6 | 15→18 | 0/0 |
| go-review-diff | 3/3 | 3/3 | 4582 B | 5493 B | +20% | 1→3 | 8→10 | 0/0 |
| go-search-large | 3/3 | 3/3 | 7332 B | 12.0 KiB | +68% | 6→6 | 13→11 | 0/0 |
| go-test-repair | 3/3 | 3/3 | 3144 B | 3220 B | +2% | 3→3 | 8→11 | 0/0 |

### claude — command families (total bytes over all runs)

| family | OFF bytes | OFF calls | FULL bytes | FULL calls | Δ |
|---|---|---|---|---|---|
| file-read | 27.3 KiB | 32 | 27.7 KiB | 34 | +1% |
| git | 25.8 KiB | 5 | 15.3 KiB | 7 | -41% |
| search | 14.0 KiB | 9 | 26.4 KiB | 10 | +89% |
| test | 18.0 KiB | 16 | 17.6 KiB | 17 | -2% |
| build | 9226 B | 5 | 10.6 KiB | 6 | +18% |
| listing | 9174 B | 4 | 8255 B | 4 | -10% |
| generic | 2988 B | 8 | 3128 B | 9 | +5% |

## Per task category (all agents, medians)

| category | runs OFF/FULL | pass OFF | pass FULL | context OFF | context FULL | Δ |
|---|---|---|---|---|---|---|
| bug-fix | 3/3 | 100% | 100% | 4972 B | 4878 B | -2% |
| compile-repair | 6/6 | 100% | 100% | 3714 B | 3730 B | +0% |
| feature | 3/3 | 100% | 100% | 2802 B | 2867 B | +2% |
| git-review | 3/3 | 100% | 100% | 4582 B | 5493 B | +20% |
| refactor | 4/3 | 75% | 100% | 2958 B | 3607 B | +22% |
| search | 3/3 | 100% | 100% | 7332 B | 12.0 KiB | +68% |
| test-repair | 4/3 | 100% | 100% | 3058 B | 3220 B | +5% |

## Anomalies and excluded runs (for manual review)

None.
