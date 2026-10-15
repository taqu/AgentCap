# Claude Code coding workflow evaluation

Attempted 2 conditions; recorded 2 trials across 1 workloads. Successful: 2; unsuccessful: 0.

Mean paired token reduction among jointly successful runs: -1.7%. Mean paired time overhead: -5.4%.

Tokens include uncached input, cache creation, cache reads, and output, using provider-reported final usage. Bytes are UTF-8 textual tool-result content; they are never converted to tokens. AgentCap shell bytes include drill-down. Resource aggregates include failures; successful-only aggregates and matched successful pairs are separate in aggregate.json.

| Level | Mode | Success | Mean / median tokens | Mean / median seconds | Mean commands | Mean drill-downs |
|---|---|---|---|---|---|---|
| L1 | disabled | 1/1 | 76965.0 / 76965.0 | 38.0 / 38.0 | 3.0 | 0.0 |
| L1 | integrated | 1/1 | 78300.0 / 78300.0 | 35.9 / 35.9 | 1.0 | 0.0 |
| L2 | disabled | 0/0 | unavailable / unavailable | unavailable / unavailable | unavailable | unavailable |
| L2 | integrated | 0/0 | unavailable / unavailable | unavailable / unavailable | unavailable | unavailable |
| L3 | disabled | 0/0 | unavailable / unavailable | unavailable / unavailable | unavailable | unavailable |
| L3 | integrated | 0/0 | unavailable / unavailable | unavailable / unavailable | unavailable | unavailable |
| L4 | disabled | 0/0 | unavailable / unavailable | unavailable / unavailable | unavailable | unavailable |
| L4 | integrated | 0/0 | unavailable / unavailable | unavailable / unavailable | unavailable | unavailable |
| L5 | disabled | 0/0 | unavailable / unavailable | unavailable / unavailable | unavailable | unavailable |
| L5 | integrated | 0/0 | unavailable / unavailable | unavailable / unavailable | unavailable | unavailable |

| Workload | Rep | Mode | Success / status | Tokens | Tool bytes | Seconds | Show + raw |
|---|---|---|---|---|---|---|---|
| coding/l1-clamp | 1 | disabled | True / completed | 76965.0 | 818 | 38.0 | 0 |
| coding/l1-clamp | 1 | integrated | True / completed | 78300.0 | 720 | 35.9 | 0 |

## Evidence and limitations

Raw JSONL and stderr are retained in each .log file. Canonical acap bench compare reports are retained for each complete pair. Inspect command sequences there before attributing a token regression to compression. Tool-output bytes exclude images and internal system prompts. Cache token accounting is explicit, but billed cost is not inferred. Repair iterations and failed agent validation attempts require manual transcript review; they are not inferred from total command counts. Shell commands count Bash invocations, not individual commands in compound scripts.

The complexity levels describe intended engineering scope, not empirically calibrated difficulty. L5 requests meaningful layer-by-layer validation, without artificial repetition or output padding. Workflow length and repeated commands remain agent decisions. Three repetitions are the default; a smoke run supports infrastructure conclusions only.

## Recommendations

First resolve any recorded execution or interception failures before interpreting efficiency. Use successful matched pairs to assess token and time regressions; report task success separately. Extend repetitions before claiming a complexity trend. Compression algorithms were left unchanged.
