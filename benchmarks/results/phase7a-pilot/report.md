# AgentCap Phase 7 benchmark report

- Batches: phase7a-pilot
- AgentCap versions: 2f2c593b4639+dirty
- Agents/models: antigravity: gemini-3.8-flash-medium (1.2.14), claude: claude-sonnet-5-5 (2.1.286)
- Tasks: 3 · Runs: 12 included, 0 excluded · Modes: disabled, integrated
- Primary metric: total agent-visible command-result bytes per workflow (capsules + deltas + `acap show` + `acap raw` + instruction overhead), counted once from the agent's own tool-result stream. Tokens marked *est* are bytes/4; host tokens are reported separately.

## antigravity

| mode | runs | pass | median context (IQR) | tokens est | median ctx, passing | savings vs OFF (all / passing) | median wall s | median shell cmds | repeated | host input tok (median) | native tool bytes (median) |
|---|---|---|---|---|---|---|---|---|---|---|---|
| OFF | 3 | 2/3 (67%) | 16.4 KiB (12.2 KiB–33.3 KiB) | 4193 | 12.2 KiB | – / – | 73 | 9 | 2 | 421996 | 174 B |
| FULL | 3 | 2/3 (67%) | 12.1 KiB (9192 B–22.0 KiB) | 3095 | 9192 B | -26% / -26% | 80 | 12 | 1 | 502604 | 174 B |

**AgentCap behavior (FULL):** 41 results (full 40 / delta 0 / unchanged 1) · intercepted 44, bypassed 0, adapter failures 0 · `show` 0 calls (0 B, in 0 runs) · `raw` 0 calls (0 B, in 0 runs; immediate 0) · raw fallback rate 0% of results · targeted drill-down rate 0% of results

**Latency (FULL, per intercepted command):** adapter 0.14 ms · AgentCap core 1.65 ms

**OFF isolation:** 33 hook decisions all bypassed (0 intercepted).

### antigravity — per task (medians; OFF → FULL)

| task | pass OFF | pass FULL | context OFF | context FULL | Δ | cmds OFF→FULL | wall s OFF→FULL | show/raw FULL |
|---|---|---|---|---|---|---|---|---|
| cpp-compile-repair | 1/1 | 1/1 | 8127 B | 6005 B | -26% | 9→12 | 73→80 | 0/0 |
| go-review-diff | 1/1 | 1/1 | 16.4 KiB | 12.1 KiB | -26% | 8→10 | 59→52 | 0/0 |
| go-search-large | 0/1 | 0/1 | 50.2 KiB | 32.0 KiB | -36% | 16→22 | 90→130 | 0/0 |

### antigravity — command families (total bytes over all runs)

| family | OFF bytes | OFF calls | FULL bytes | FULL calls | Δ |
|---|---|---|---|---|---|
| search | 17.5 KiB | 5 | 26.0 KiB | 9 | +48% |
| git | 23.2 KiB | 12 | 3947 B | 15 | -83% |
| test | 15.0 KiB | 8 | 1536 B | 8 | -90% |
| listing | 9873 B | 4 | 3278 B | 4 | -67% |
| generic | 0 B | 0 | 10.5 KiB | 4 | – |
| file-read | 7129 B | 1 | 2830 B | 1 | -60% |
| build | 2163 B | 3 | 2167 B | 3 | +0% |

## claude

| mode | runs | pass | median context (IQR) | tokens est | median ctx, passing | savings vs OFF (all / passing) | median wall s | median shell cmds | repeated | host input tok (median) | native tool bytes (median) |
|---|---|---|---|---|---|---|---|---|---|---|---|
| OFF | 3 | 3/3 (100%) | 4582 B (4206 B–10170 B) | 1146 | 4582 B | – / – | 11 | 3 | 0 | 60241 | 0 B |
| FULL | 3 | 3/3 (100%) | 5645 B (4821 B–8042 B) | 1411 | 5645 B | +23% / +23% | 13 | 3 | 0 | 78105 | 0 B |

**AgentCap behavior (FULL):** 12 results (full 12 / delta 0 / unchanged 0) · intercepted 12, bypassed 0, adapter failures 0 · `show` 0 calls (0 B, in 0 runs) · `raw` 0 calls (0 B, in 0 runs; immediate 0) · raw fallback rate 0% of results · targeted drill-down rate 0% of results

**Latency (FULL, per intercepted command):** adapter 0.11 ms · AgentCap core 1.08 ms

**OFF isolation:** 9 hook decisions all bypassed (0 intercepted).

### claude — per task (medians; OFF → FULL)

| task | pass OFF | pass FULL | context OFF | context FULL | Δ | cmds OFF→FULL | wall s OFF→FULL | show/raw FULL |
|---|---|---|---|---|---|---|---|---|
| cpp-compile-repair | 1/1 | 1/1 | 3831 B | 3997 B | +4% | 3→3 | 14→12 | 0/0 |
| go-review-diff | 1/1 | 1/1 | 4582 B | 5645 B | +23% | 1→3 | 7→13 | 0/0 |
| go-search-large | 1/1 | 1/1 | 15.4 KiB | 10.2 KiB | -34% | 5→6 | 11→16 | 0/0 |

### claude — command families (total bytes over all runs)

| family | OFF bytes | OFF calls | FULL bytes | FULL calls | Δ |
|---|---|---|---|---|---|
| git | 17.3 KiB | 2 | 5645 B | 3 | -68% |
| file-read | 3739 B | 5 | 6832 B | 5 | +83% |
| test | 2675 B | 1 | 2725 B | 1 | +2% |
| listing | 0 B | 0 | 3003 B | 1 | – |
| search | 52 B | 1 | 1877 B | 2 | +3510% |

## Per task category (all agents, medians)

| category | runs OFF/FULL | pass OFF | pass FULL | context OFF | context FULL | Δ |
|---|---|---|---|---|---|---|
| compile-repair | 2/2 | 100% | 100% | 5979 B | 5001 B | -16% |
| git-review | 2/2 | 100% | 100% | 10.4 KiB | 9012 B | -16% |
| search | 2/2 | 50% | 50% | 32.8 KiB | 21.1 KiB | -36% |

## Anomalies and excluded runs (for manual review)

- `20260930T193858.012-44610` antigravity FULL go-search-large — outcome FAIL
