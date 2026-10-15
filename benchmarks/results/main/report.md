# Claude Code coding workflow evaluation

Attempted 120 conditions; recorded 120 trials across 10 workloads. Successful: 120; unsuccessful: 0.

Mean paired token reduction among jointly successful enabled-mode pairs: -4.1%. Mean paired time overhead: 2.3%.

Compared modes: disabled, integrated, stateless, stateful. Each enabled run is paired with the matching disabled baseline by workload and repetition.

Tokens include uncached input, cache creation, cache reads, and output, using provider-reported final usage. Bytes are UTF-8 textual tool-result content; they are never converted to tokens. AgentCap shell bytes include drill-down. Resource aggregates include failures; successful-only aggregates and matched successful pairs are separate in aggregate.json.

| Level | Mode | Success | Mean / median tokens | Mean / median seconds | Mean tools / Bash | Show + raw | Repeated Bash | Build/test | Failed validations |
|---|---|---|---|---|---|---|---|---|---|
| L1 | disabled | 6/6 | 75452.3 / 77422.0 | 34.1 / 35.7 | 6.0 / 2.7 | 0.0 | 0.0 | 1.0 | 0.0 |
| L1 | integrated | 6/6 | 73642.8 / 75594.0 | 30.2 / 30.1 | 6.2 / 2.2 | 0.0 | 0.0 | 1.0 | 0.0 |
| L1 | stateless | 6/6 | 78944.2 / 77898.5 | 37.0 / 37.2 | 6.5 / 2.8 | 0.0 | 0.0 | 1.0 | 0.0 |
| L1 | stateful | 6/6 | 74895.0 / 77888.0 | 33.5 / 34.3 | 6.2 / 2.3 | 0.0 | 0.0 | 1.0 | 0.0 |
| L2 | disabled | 6/6 | 94861.2 / 103462.5 | 42.5 / 45.2 | 9.3 / 3.3 | 0.0 | 0.0 | 1.0 | 0.0 |
| L2 | integrated | 6/6 | 107660.8 / 108385.0 | 46.0 / 46.3 | 9.7 / 3.7 | 0.0 | 0.0 | 1.0 | 0.0 |
| L2 | stateless | 6/6 | 93261.2 / 98299.0 | 42.8 / 41.9 | 9.0 / 3.0 | 0.0 | 0.0 | 1.0 | 0.0 |
| L2 | stateful | 6/6 | 92427.2 / 91701.0 | 42.0 / 40.0 | 8.8 / 2.8 | 0.0 | 0.0 | 1.0 | 0.0 |
| L3 | disabled | 6/6 | 107744.2 / 103825.0 | 56.7 / 58.6 | 10.7 / 2.0 | 0.0 | 0.0 | 1.0 | 0.0 |
| L3 | integrated | 6/6 | 120581.3 / 119043.0 | 56.9 / 55.0 | 11.0 / 2.3 | 0.0 | 0.0 | 1.0 | 0.0 |
| L3 | stateless | 6/6 | 117862.0 / 117951.5 | 58.1 / 56.5 | 11.0 / 2.0 | 0.0 | 0.0 | 1.0 | 0.0 |
| L3 | stateful | 6/6 | 110972.5 / 109474.5 | 59.1 / 58.5 | 11.0 / 2.3 | 0.0 | 0.0 | 1.0 | 0.0 |
| L4 | disabled | 6/6 | 156357.8 / 147187.0 | 70.9 / 71.9 | 13.5 / 2.3 | 0.0 | 0.0 | 1.3 | 0.3 |
| L4 | integrated | 6/6 | 165006.0 / 154906.0 | 73.0 / 67.0 | 14.0 / 3.0 | 0.0 | 0.0 | 1.5 | 0.5 |
| L4 | stateless | 6/6 | 156674.2 / 151847.5 | 74.7 / 78.7 | 13.5 / 2.5 | 0.0 | 0.0 | 1.5 | 0.5 |
| L4 | stateful | 6/6 | 156428.7 / 151854.0 | 68.5 / 67.1 | 13.2 / 2.2 | 0.0 | 0.0 | 1.2 | 0.2 |
| L5 | disabled | 6/6 | 322544.8 / 331267.5 | 126.0 / 125.1 | 22.3 / 5.2 | 0.0 | 0.2 | 3.0 | 0.2 |
| L5 | integrated | 6/6 | 307233.3 / 314122.0 | 120.3 / 120.3 | 21.8 / 4.7 | 0.0 | 0.0 | 2.3 | 0.0 |
| L5 | stateless | 6/6 | 314547.3 / 314094.5 | 116.1 / 112.5 | 21.5 / 4.5 | 0.0 | 0.0 | 2.5 | 0.0 |
| L5 | stateful | 6/6 | 349513.8 / 350989.0 | 131.3 / 132.8 | 22.8 / 5.3 | 0.0 | 0.0 | 3.3 | 0.0 |

Paired means below use only jointly successful trials. Negative token/output reduction means increased consumption.

| Level | Mode | Successful pairs | Token reduction % | Tool-output reduction % | Time overhead % |
|---|---|---|---|---|---|
| L1 | disabled | 0/0 | unavailable | unavailable | unavailable |
| L1 | integrated | 6/6 | 2.3 | -25.5 | -9.4 |
| L1 | stateless | 6/6 | -4.8 | -35.7 | 11.3 |
| L1 | stateful | 6/6 | 0.3 | 4.9 | 0.4 |
| L2 | disabled | 0/0 | unavailable | unavailable | unavailable |
| L2 | integrated | 6/6 | -16.5 | 2.4 | 11.5 |
| L2 | stateless | 6/6 | -0.5 | -27.5 | 4.0 |
| L2 | stateful | 6/6 | 1.5 | -40.4 | 0.8 |
| L3 | disabled | 0/0 | unavailable | unavailable | unavailable |
| L3 | integrated | 6/6 | -12.7 | -1.1 | 1.1 |
| L3 | stateless | 6/6 | -10.5 | 0.9 | 4.2 |
| L3 | stateful | 6/6 | -3.8 | 0.0 | 5.7 |
| L4 | disabled | 0/0 | unavailable | unavailable | unavailable |
| L4 | integrated | 6/6 | -6.9 | 15.7 | 3.3 |
| L4 | stateless | 6/6 | -3.0 | 1.3 | 5.5 |
| L4 | stateful | 6/6 | -2.1 | 2.0 | -3.2 |
| L5 | disabled | 0/0 | unavailable | unavailable | unavailable |
| L5 | integrated | 6/6 | 3.8 | -19.6 | -2.9 |
| L5 | stateless | 6/6 | 1.2 | -18.1 | -5.3 |
| L5 | stateful | 6/6 | -10.1 | -21.0 | 7.6 |

| Workload | Rep | Mode | Success / status | Model tokens | All-tool B | Bash calls | AgentCap raw → stateful B | Processing ms | Δ / same | Show / raw | Seconds |
|---|---|---|---|---|---|---|---|---|---|---|---|
| coding/l1-clamp | 1 | disabled | True / completed | 77327.0 | 818 | — / 3 | n/a | n/a | n/a | 0 / 0 | 35.3 |
| coding/l1-clamp | 1 | integrated | True / completed | 78823.0 | 830 | 0 / 3 | 0 → 0 | 0.00 | 0 / 0 | 0 / 0 | 36.9 |
| coding/l1-duration | 1 | integrated | True / completed | 79198.0 | 1006 | 1 / 3 | 149 → 37 | 11.70 | 0 / 0 | 0 / 0 | 37.4 |
| coding/l1-duration | 1 | disabled | True / completed | 70997.0 | 4535 | — / 2 | n/a | n/a | n/a | 0 / 0 | 28.8 |
| coding/l2-config | 1 | disabled | True / completed | 74066.0 | 5530 | — / 2 | n/a | n/a | n/a | 0 / 0 | 30.8 |
| coding/l2-config | 1 | integrated | True / completed | 108806.0 | 1386 | 1 / 4 | 201 → 37 | 15.58 | 0 / 0 | 0 / 0 | 48.3 |
| coding/l2-pagination | 1 | integrated | True / completed | 108405.0 | 1290 | 1 / 4 | 202 → 37 | 12.23 | 0 / 0 | 0 / 0 | 53.4 |
| coding/l2-pagination | 1 | disabled | True / completed | 107068.0 | 1448 | — / 4 | n/a | n/a | n/a | 0 / 0 | 46.9 |
| coding/l3-ledger | 1 | disabled | True / completed | 133726.0 | 8223 | — / 2 | n/a | n/a | n/a | 0 / 0 | 55.0 |
| coding/l3-ledger | 1 | integrated | True / completed | 132056.0 | 8091 | 1 / 2 | 327 → 37 | 10.86 | 0 / 0 | 0 / 0 | 49.7 |
| coding/l3-planner | 1 | integrated | True / completed | 103535.0 | 6844 | 1 / 2 | 233 → 37 | 11.87 | 0 / 0 | 0 / 0 | 58.6 |
| coding/l3-planner | 1 | disabled | True / completed | 104901.0 | 6790 | — / 2 | n/a | n/a | n/a | 0 / 0 | 63.7 |
| coding/l4-store | 1 | disabled | True / completed | 100842.0 | 8360 | — / 2 | n/a | n/a | n/a | 0 / 0 | 57.3 |
| coding/l4-store | 1 | integrated | True / completed | 122339.0 | 8230 | 1 / 2 | 326 → 37 | 11.58 | 0 / 0 | 0 / 0 | 64.7 |
| coding/l4-context-migration | 1 | integrated | True / completed | 260129.0 | 5367 | 3 / 6 | 686 → 598 | 40.32 | 0 / 0 | 0 / 0 | 104.2 |
| coding/l4-context-migration | 1 | disabled | True / completed | 200601.0 | 9664 | — / 2 | n/a | n/a | n/a | 0 / 0 | 80.3 |
| coding/l5-ledger-workflow | 1 | disabled | True / completed | 264403.0 | 12118 | — / 5 | n/a | n/a | n/a | 0 / 0 | 95.2 |
| coding/l5-ledger-workflow | 1 | integrated | True / completed | 275101.0 | 12713 | 3 / 5 | 713 → 752 | 32.39 | 0 / 0 | 0 / 0 | 105.4 |
| coding/l5-planner-workflow | 1 | integrated | True / completed | 327437.0 | 11975 | 3 / 5 | 540 → 579 | 31.02 | 0 / 0 | 0 / 0 | 134.9 |
| coding/l5-planner-workflow | 1 | disabled | True / completed | 357656.0 | 12072 | — / 5 | n/a | n/a | n/a | 0 / 0 | 129.7 |
| coding/l1-clamp | 2 | integrated | True / completed | 66362.0 | 1050 | 0 / 1 | 0 → 0 | 0.00 | 0 / 0 | 0 / 0 | 23.2 |
| coding/l1-clamp | 2 | disabled | True / completed | 70949.0 | 4286 | — / 2 | n/a | n/a | n/a | 0 / 0 | 26.4 |
| coding/l1-duration | 2 | disabled | True / completed | 77936.0 | 1104 | — / 3 | n/a | n/a | n/a | 0 / 0 | 36.8 |
| coding/l1-duration | 2 | integrated | True / completed | 79260.0 | 1122 | 0 / 3 | 0 → 0 | 0.00 | 0 / 0 | 0 / 0 | 31.4 |
| coding/l2-config | 2 | integrated | True / completed | 104587.0 | 5567 | 1 / 3 | 201 → 37 | 13.02 | 0 / 0 | 0 / 0 | 39.6 |
| coding/l2-config | 2 | disabled | True / completed | 100267.0 | 5691 | — / 4 | n/a | n/a | n/a | 0 / 0 | 46.9 |
| coding/l2-pagination | 2 | disabled | True / completed | 106658.0 | 1445 | — / 4 | n/a | n/a | n/a | 0 / 0 | 43.4 |
| coding/l2-pagination | 2 | integrated | True / completed | 110737.0 | 2644 | 1 / 4 | 1241 → 1221 | 9.91 | 0 / 0 | 0 / 0 | 41.2 |
| coding/l3-ledger | 2 | integrated | True / completed | 118712.0 | 8100 | 1 / 2 | 327 → 37 | 10.43 | 0 / 0 | 0 / 0 | 49.7 |
| coding/l3-ledger | 2 | disabled | True / completed | 98135.0 | 7675 | — / 2 | n/a | n/a | n/a | 0 / 0 | 41.6 |
| coding/l3-planner | 2 | disabled | True / completed | 104063.0 | 6790 | — / 2 | n/a | n/a | n/a | 0 / 0 | 62.0 |
| coding/l3-planner | 2 | integrated | True / completed | 132005.0 | 7019 | 1 / 4 | 233 → 37 | 10.46 | 0 / 0 | 0 / 0 | 71.1 |
| coding/l4-store | 2 | integrated | True / completed | 104333.0 | 8745 | 1 / 2 | 490 → 503 | 12.46 | 0 / 0 | 0 / 0 | 61.0 |
| coding/l4-store | 2 | disabled | True / completed | 121549.0 | 8506 | — / 2 | n/a | n/a | n/a | 0 / 0 | 61.1 |
| coding/l4-context-migration | 2 | disabled | True / completed | 172825.0 | 9560 | — / 2 | n/a | n/a | n/a | 0 / 0 | 79.9 |
| coding/l4-context-migration | 2 | integrated | True / completed | 193423.0 | 4724 | 1 / 3 | 326 → 37 | 10.19 | 0 / 0 | 0 / 0 | 68.5 |
| coding/l5-ledger-workflow | 2 | integrated | True / completed | 317804.0 | 12250 | 3 / 5 | 803 → 842 | 42.67 | 0 / 0 | 0 / 0 | 101.1 |
| coding/l5-ledger-workflow | 2 | disabled | True / completed | 296491.0 | 5553 | — / 6 | n/a | n/a | n/a | 0 / 0 | 120.6 |
| coding/l5-planner-workflow | 2 | disabled | True / completed | 353893.0 | 12005 | — / 5 | n/a | n/a | n/a | 0 / 0 | 140.1 |
| coding/l5-planner-workflow | 2 | integrated | True / completed | 279124.0 | 12173 | 1 / 4 | 394 → 407 | 16.06 | 0 / 0 | 0 / 0 | 129.5 |
| coding/l1-clamp | 3 | disabled | True / completed | 77517.0 | 822 | — / 3 | n/a | n/a | n/a | 0 / 0 | 36.1 |
| coding/l1-clamp | 3 | integrated | True / completed | 65849.0 | 715 | 0 / 1 | 0 → 0 | 0.00 | 0 / 0 | 0 / 0 | 23.8 |
| coding/l1-duration | 3 | integrated | True / completed | 72365.0 | 4631 | 1 / 2 | 149 → 37 | 9.83 | 0 / 0 | 0 / 0 | 28.8 |
| coding/l1-duration | 3 | disabled | True / completed | 77988.0 | 1112 | — / 3 | n/a | n/a | n/a | 0 / 0 | 41.3 |
| coding/l2-config | 3 | disabled | True / completed | 74420.0 | 5530 | — / 2 | n/a | n/a | n/a | 0 / 0 | 35.9 |
| coding/l2-config | 3 | integrated | True / completed | 105065.0 | 5606 | 1 / 3 | 201 → 37 | 14.33 | 0 / 0 | 0 / 0 | 44.4 |
| coding/l2-pagination | 3 | integrated | True / completed | 108365.0 | 1293 | 1 / 4 | 202 → 37 | 12.58 | 0 / 0 | 0 / 0 | 49.2 |
| coding/l2-pagination | 3 | disabled | True / completed | 106688.0 | 1448 | — / 4 | n/a | n/a | n/a | 0 / 0 | 51.2 |
| coding/l3-ledger | 3 | disabled | True / completed | 102053.0 | 8223 | — / 2 | n/a | n/a | n/a | 0 / 0 | 55.2 |
| coding/l3-ledger | 3 | integrated | True / completed | 117806.0 | 8100 | 1 / 2 | 327 → 37 | 13.89 | 0 / 0 | 0 / 0 | 51.3 |
| coding/l3-planner | 3 | integrated | True / completed | 119374.0 | 6844 | 1 / 2 | 233 → 37 | 14.32 | 0 / 0 | 0 / 0 | 61.0 |
| coding/l3-planner | 3 | disabled | True / completed | 103587.0 | 6832 | — / 2 | n/a | n/a | n/a | 0 / 0 | 62.7 |
| coding/l4-store | 3 | disabled | True / completed | 120202.0 | 8504 | — / 3 | n/a | n/a | n/a | 0 / 0 | 63.8 |
| coding/l4-store | 3 | integrated | True / completed | 138917.0 | 8736 | 2 / 3 | 506 → 532 | 25.52 | 0 / 0 | 0 / 0 | 65.5 |
| coding/l4-context-migration | 3 | integrated | True / completed | 170895.0 | 9529 | 1 / 2 | 326 → 37 | 12.18 | 0 / 0 | 0 / 0 | 74.1 |
| coding/l4-context-migration | 3 | disabled | True / completed | 222128.0 | 9855 | — / 3 | n/a | n/a | n/a | 0 / 0 | 83.2 |
| coding/l5-ledger-workflow | 3 | disabled | True / completed | 308642.0 | 12118 | — / 5 | n/a | n/a | n/a | 0 / 0 | 102.4 |
| coding/l5-ledger-workflow | 3 | integrated | True / completed | 333494.0 | 11519 | 3 / 5 | 624 → 111 | 44.37 | 0 / 0 | 0 / 0 | 111.2 |
| coding/l5-planner-workflow | 3 | integrated | True / completed | 310440.0 | 12020 | 1 / 4 | 394 → 407 | 17.47 | 0 / 0 | 0 / 0 | 139.8 |
| coding/l5-planner-workflow | 3 | disabled | True / completed | 354184.0 | 12486 | — / 5 | n/a | n/a | n/a | 0 / 0 | 168.0 |
| coding/l1-clamp | 1 | stateless | True / completed | 77358.0 | 708 | 1 / 3 | 149 → 37 | 12.89 | 0 / 0 | 0 / 0 | 36.3 |
| coding/l1-clamp | 1 | stateful | True / completed | 65829.0 | 1580 | 0 / 1 | 0 → 0 | 0.00 | 0 / 0 | 0 / 0 | 25.5 |
| coding/l1-duration | 1 | stateful | True / completed | 77923.0 | 1000 | 1 / 3 | 149 → 37 | 12.81 | 0 / 0 | 0 / 0 | 35.7 |
| coding/l1-duration | 1 | stateless | True / completed | 77880.0 | 998 | 1 / 3 | 149 → 37 | 11.03 | 0 / 0 | 0 / 0 | 38.9 |
| coding/l2-config | 1 | stateless | True / completed | 74071.0 | 5438 | 1 / 2 | 201 → 37 | 12.10 | 0 / 0 | 0 / 0 | 34.9 |
| coding/l2-config | 1 | stateful | True / completed | 89788.0 | 5476 | 1 / 2 | 201 → 37 | 12.26 | 0 / 0 | 0 / 0 | 38.1 |
| coding/l2-pagination | 1 | stateful | True / completed | 106587.0 | 1277 | 1 / 4 | 202 → 37 | 13.62 | 0 / 0 | 0 / 0 | 51.1 |
| coding/l2-pagination | 1 | stateless | True / completed | 107005.0 | 1283 | 1 / 4 | 202 → 37 | 19.05 | 0 / 0 | 0 / 0 | 48.8 |
| coding/l3-ledger | 1 | stateless | True / completed | 118866.0 | 7896 | 1 / 2 | 327 → 37 | 10.51 | 0 / 0 | 0 / 0 | 60.5 |
| coding/l3-ledger | 1 | stateful | True / completed | 116196.0 | 8002 | 1 / 2 | 327 → 37 | 13.87 | 0 / 0 | 0 / 0 | 51.6 |
| coding/l3-planner | 1 | stateful | True / completed | 128960.0 | 6847 | 1 / 4 | 233 → 37 | 12.03 | 0 / 0 | 0 / 0 | 75.5 |
| coding/l3-planner | 1 | stateless | True / completed | 113663.0 | 6746 | 1 / 2 | 233 → 37 | 14.90 | 0 / 0 | 0 / 0 | 52.6 |
| coding/l4-store | 1 | stateless | True / completed | 117929.0 | 8133 | 1 / 2 | 326 → 37 | 12.21 | 0 / 0 | 0 / 0 | 57.4 |
| coding/l4-store | 1 | stateful | True / completed | 117749.0 | 8181 | 1 / 2 | 326 → 37 | 15.15 | 0 / 0 | 0 / 0 | 59.9 |
| coding/l4-context-migration | 1 | stateful | True / completed | 199551.0 | 9379 | 1 / 2 | 326 → 37 | 12.37 | 0 / 0 | 0 / 0 | 75.7 |
| coding/l4-context-migration | 1 | stateless | True / completed | 166750.0 | 9485 | 1 / 2 | 326 → 37 | 11.58 | 0 / 0 | 0 / 0 | 78.5 |
| coding/l5-ledger-workflow | 1 | stateless | True / completed | 315878.0 | 11540 | 3 / 5 | 844 → 111 | 39.88 | 0 / 0 | 0 / 0 | 113.9 |
| coding/l5-ledger-workflow | 1 | stateful | True / completed | 360120.0 | 11786 | 4 / 6 | 744 → 148 | 52.22 | 0 / 0 | 0 / 0 | 131.5 |
| coding/l5-planner-workflow | 1 | stateful | True / completed | 354677.0 | 12321 | 3 / 5 | 859 → 898 | 39.47 | 0 / 0 | 0 / 0 | 134.1 |
| coding/l5-planner-workflow | 1 | stateless | True / completed | 288251.0 | 11910 | 1 / 2 | 619 → 632 | 11.15 | 0 / 0 | 0 / 0 | 104.7 |
| coding/l1-clamp | 2 | stateful | True / completed | 85449.0 | 4109 | 1 / 2 | 149 → 37 | 12.81 | 0 / 0 | 0 / 0 | 32.9 |
| coding/l1-clamp | 2 | stateless | True / completed | 77405.0 | 818 | 0 / 3 | 0 → 0 | 0.00 | 0 / 0 | 0 / 0 | 37.4 |
| coding/l1-duration | 2 | stateless | True / completed | 78011.0 | 1112 | 0 / 3 | 0 → 0 | 0.00 | 0 / 0 | 0 / 0 | 39.8 |
| coding/l1-duration | 2 | stateful | True / completed | 77853.0 | 996 | 1 / 3 | 149 → 37 | 14.06 | 0 / 0 | 0 / 0 | 38.9 |
| coding/l2-config | 2 | stateful | True / completed | 74352.0 | 5476 | 1 / 2 | 201 → 37 | 12.65 | 0 / 0 | 0 / 0 | 33.0 |
| coding/l2-config | 2 | stateless | True / completed | 74152.0 | 5438 | 1 / 2 | 201 → 37 | 19.26 | 0 / 0 | 0 / 0 | 35.1 |
| coding/l2-pagination | 2 | stateless | True / completed | 106751.0 | 1280 | 1 / 4 | 202 → 37 | 13.62 | 0 / 0 | 0 / 0 | 48.5 |
| coding/l2-pagination | 2 | stateful | True / completed | 115839.0 | 5504 | 1 / 4 | 202 → 37 | 13.60 | 0 / 0 | 0 / 0 | 51.4 |
| coding/l3-ledger | 2 | stateful | True / completed | 99582.0 | 7945 | 1 / 2 | 327 → 37 | 11.74 | 0 / 0 | 0 / 0 | 57.9 |
| coding/l3-ledger | 2 | stateless | True / completed | 117062.0 | 7994 | 1 / 2 | 327 → 37 | 13.71 | 0 / 0 | 0 / 0 | 53.7 |
| coding/l3-planner | 2 | stateless | True / completed | 121678.0 | 6703 | 1 / 2 | 233 → 37 | 10.92 | 0 / 0 | 0 / 0 | 69.0 |
| coding/l3-planner | 2 | stateful | True / completed | 102753.0 | 6746 | 1 / 2 | 233 → 37 | 13.74 | 0 / 0 | 0 / 0 | 60.4 |
| coding/l4-store | 2 | stateful | True / completed | 135351.0 | 8803 | 2 / 3 | 670 → 696 | 25.40 | 0 / 0 | 0 / 0 | 59.5 |
| coding/l4-store | 2 | stateless | True / completed | 117747.0 | 8181 | 1 / 2 | 326 → 37 | 13.04 | 0 / 0 | 0 / 0 | 60.9 |
| coding/l4-context-migration | 2 | stateless | True / completed | 209620.0 | 9991 | 2 / 3 | 506 → 532 | 27.26 | 0 / 0 | 0 / 0 | 89.2 |
| coding/l4-context-migration | 2 | stateful | True / completed | 168357.0 | 9485 | 1 / 2 | 326 → 37 | 12.32 | 0 / 0 | 0 / 0 | 85.2 |
| coding/l5-ledger-workflow | 2 | stateful | True / completed | 367792.0 | 12137 | 4 / 6 | 833 → 885 | 88.83 | 0 / 0 | 0 / 0 | 142.6 |
| coding/l5-ledger-workflow | 2 | stateless | True / completed | 320155.0 | 12229 | 3 / 5 | 803 → 842 | 39.74 | 0 / 0 | 0 / 0 | 114.7 |
| coding/l5-planner-workflow | 2 | stateless | True / completed | 312311.0 | 12177 | 2 / 5 | 695 → 721 | 30.13 | 0 / 0 | 0 / 0 | 151.2 |
| coding/l5-planner-workflow | 2 | stateful | True / completed | 321358.0 | 12529 | 2 / 4 | 695 → 721 | 27.82 | 0 / 0 | 0 / 0 | 126.5 |
| coding/l1-clamp | 3 | stateless | True / completed | 85094.0 | 4078 | 1 / 2 | 149 → 37 | 13.33 | 0 / 0 | 0 / 0 | 32.5 |
| coding/l1-clamp | 3 | stateful | True / completed | 64330.0 | 571 | 1 / 2 | 149 → 37 | 11.18 | 0 / 0 | 0 / 0 | 30.9 |
| coding/l1-duration | 3 | stateful | True / completed | 77986.0 | 1110 | 0 / 3 | 0 → 0 | 0.00 | 0 / 0 | 0 / 0 | 37.2 |
| coding/l1-duration | 3 | stateless | True / completed | 77917.0 | 998 | 1 / 3 | 149 → 37 | 9.89 | 0 / 0 | 0 / 0 | 37.1 |
| coding/l2-config | 3 | stateless | True / completed | 107741.0 | 1379 | 1 / 4 | 201 → 37 | 11.57 | 0 / 0 | 0 / 0 | 54.3 |
| coding/l2-config | 3 | stateful | True / completed | 74383.0 | 5476 | 1 / 2 | 201 → 37 | 11.13 | 0 / 0 | 0 / 0 | 36.5 |
| coding/l2-pagination | 3 | stateful | True / completed | 93614.0 | 1144 | 1 / 3 | 202 → 37 | 11.76 | 0 / 0 | 0 / 0 | 42.0 |
| coding/l2-pagination | 3 | stateless | True / completed | 89847.0 | 5343 | 1 / 2 | 202 → 37 | 11.91 | 0 / 0 | 0 / 0 | 35.4 |
| coding/l3-ledger | 3 | stateless | True / completed | 117470.0 | 8043 | 1 / 2 | 327 → 37 | 43.93 | 0 / 0 | 0 / 0 | 55.3 |
| coding/l3-ledger | 3 | stateful | True / completed | 116984.0 | 8002 | 1 / 2 | 327 → 37 | 13.73 | 0 / 0 | 0 / 0 | 50.2 |
| coding/l3-planner | 3 | stateful | True / completed | 101360.0 | 6942 | 0 / 2 | 0 → 0 | 0.00 | 0 / 0 | 0 / 0 | 59.1 |
| coding/l3-planner | 3 | stateless | True / completed | 118433.0 | 6746 | 1 / 2 | 233 → 37 | 12.55 | 0 / 0 | 0 / 0 | 57.6 |
| coding/l4-store | 3 | stateless | True / completed | 136945.0 | 8337 | 2 / 3 | 506 → 230 | 20.01 | 0 / 0 | 0 / 0 | 78.8 |
| coding/l4-store | 3 | stateful | True / completed | 117463.0 | 8133 | 1 / 2 | 326 → 37 | 11.20 | 0 / 0 | 0 / 0 | 56.6 |
| coding/l4-context-migration | 3 | stateful | True / completed | 200101.0 | 9326 | 1 / 2 | 326 → 37 | 10.86 | 0 / 0 | 0 / 0 | 74.3 |
| coding/l4-context-migration | 3 | stateless | True / completed | 191054.0 | 9689 | 2 / 3 | 506 → 230 | 23.33 | 0 / 0 | 0 / 0 | 83.6 |
| coding/l5-ledger-workflow | 3 | stateless | True / completed | 312011.0 | 11608 | 3 / 5 | 844 → 111 | 30.43 | 0 / 0 | 0 / 0 | 101.2 |
| coding/l5-ledger-workflow | 3 | stateful | True / completed | 345835.0 | 12381 | 4 / 6 | 833 → 885 | 41.45 | 0 / 0 | 0 / 0 | 116.7 |
| coding/l5-planner-workflow | 3 | stateful | True / completed | 347301.0 | 12705 | 3 / 5 | 859 → 898 | 32.23 | 0 / 0 | 0 / 0 | 136.7 |
| coding/l5-planner-workflow | 3 | stateless | True / completed | 338678.0 | 12178 | 3 / 5 | 540 → 579 | 30.49 | 0 / 0 | 0 / 0 | 111.1 |

AgentCap intercepted 119/274 enabled-mode Bash invocations. Raw/reduced/stateful byte fields cover intercepted commands only. Unintercepted native outputs are included in the all-tool byte comparison above. A zero intercepted byte count is not evidence of complete output elimination.

## Regression evidence

- coding/l1-clamp integrated repetition 1: token usage increased 1.9%; time overhead 4.6%; 0 drill-downs. [Transcript](l1-clamp-1-integrated.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l1-clamp stateless repetition 1: token usage increased 0.0%; time overhead 3.0%; 0 drill-downs. [Transcript](l1-clamp-1-stateless.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l1-duration integrated repetition 1: token usage increased 11.6%; time overhead 30.1%; 0 drill-downs. [Transcript](l1-duration-1-integrated.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l1-duration stateless repetition 1: token usage increased 9.7%; time overhead 35.1%; 0 drill-downs. [Transcript](l1-duration-1-stateless.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l1-duration stateful repetition 1: token usage increased 9.8%; time overhead 24.0%; 0 drill-downs. [Transcript](l1-duration-1-stateful.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l2-config integrated repetition 1: token usage increased 46.9%; time overhead 56.7%; 0 drill-downs. [Transcript](l2-config-1-integrated.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l2-config stateless repetition 1: token usage increased 0.0%; time overhead 13.2%; 0 drill-downs. [Transcript](l2-config-1-stateless.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l2-config stateful repetition 1: token usage increased 21.2%; time overhead 23.8%; 0 drill-downs. [Transcript](l2-config-1-stateful.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l2-pagination integrated repetition 1: token usage increased 1.2%; time overhead 13.8%; 0 drill-downs. [Transcript](l2-pagination-1-integrated.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l3-planner stateless repetition 1: token usage increased 8.4%; time overhead -17.4%; 0 drill-downs. [Transcript](l3-planner-1-stateless.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l3-planner stateful repetition 1: token usage increased 22.9%; time overhead 18.5%; 0 drill-downs. [Transcript](l3-planner-1-stateful.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l4-store integrated repetition 1: token usage increased 21.3%; time overhead 12.9%; 0 drill-downs. [Transcript](l4-store-1-integrated.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l4-store stateless repetition 1: token usage increased 16.9%; time overhead 0.2%; 0 drill-downs. [Transcript](l4-store-1-stateless.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l4-store stateful repetition 1: token usage increased 16.8%; time overhead 4.6%; 0 drill-downs. [Transcript](l4-store-1-stateful.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l4-context-migration integrated repetition 1: token usage increased 29.7%; time overhead 29.7%; 0 drill-downs. [Transcript](l4-context-migration-1-integrated.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l5-ledger-workflow integrated repetition 1: token usage increased 4.0%; time overhead 10.7%; 0 drill-downs. [Transcript](l5-ledger-workflow-1-integrated.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l5-ledger-workflow stateless repetition 1: token usage increased 19.5%; time overhead 19.7%; 0 drill-downs. [Transcript](l5-ledger-workflow-1-stateless.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l5-ledger-workflow stateful repetition 1: token usage increased 36.2%; time overhead 38.2%; 0 drill-downs. [Transcript](l5-ledger-workflow-1-stateful.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l1-clamp stateless repetition 2: token usage increased 9.1%; time overhead 41.9%; 0 drill-downs. [Transcript](l1-clamp-2-stateless.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l1-clamp stateful repetition 2: token usage increased 20.4%; time overhead 24.6%; 0 drill-downs. [Transcript](l1-clamp-2-stateful.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l1-duration integrated repetition 2: token usage increased 1.7%; time overhead -14.7%; 0 drill-downs. [Transcript](l1-duration-2-integrated.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l1-duration stateless repetition 2: token usage increased 0.1%; time overhead 8.1%; 0 drill-downs. [Transcript](l1-duration-2-stateless.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l2-config integrated repetition 2: token usage increased 4.3%; time overhead -15.7%; 0 drill-downs. [Transcript](l2-config-2-integrated.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l2-pagination integrated repetition 2: token usage increased 3.8%; time overhead -5.2%; 0 drill-downs. [Transcript](l2-pagination-2-integrated.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l2-pagination stateless repetition 2: token usage increased 0.1%; time overhead 11.7%; 0 drill-downs. [Transcript](l2-pagination-2-stateless.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l2-pagination stateful repetition 2: token usage increased 8.6%; time overhead 18.5%; 0 drill-downs. [Transcript](l2-pagination-2-stateful.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l3-ledger integrated repetition 2: token usage increased 21.0%; time overhead 19.5%; 0 drill-downs. [Transcript](l3-ledger-2-integrated.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l3-ledger stateless repetition 2: token usage increased 19.3%; time overhead 29.2%; 0 drill-downs. [Transcript](l3-ledger-2-stateless.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l3-ledger stateful repetition 2: token usage increased 1.5%; time overhead 39.2%; 0 drill-downs. [Transcript](l3-ledger-2-stateful.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l3-planner integrated repetition 2: token usage increased 26.9%; time overhead 14.7%; 0 drill-downs. [Transcript](l3-planner-2-integrated.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l3-planner stateless repetition 2: token usage increased 16.9%; time overhead 11.3%; 0 drill-downs. [Transcript](l3-planner-2-stateless.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l4-store stateful repetition 2: token usage increased 11.4%; time overhead -2.6%; 0 drill-downs. [Transcript](l4-store-2-stateful.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l4-context-migration integrated repetition 2: token usage increased 11.9%; time overhead -14.2%; 0 drill-downs. [Transcript](l4-context-migration-2-integrated.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l4-context-migration stateless repetition 2: token usage increased 21.3%; time overhead 11.7%; 0 drill-downs. [Transcript](l4-context-migration-2-stateless.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l5-ledger-workflow integrated repetition 2: token usage increased 7.2%; time overhead -16.2%; 0 drill-downs. [Transcript](l5-ledger-workflow-2-integrated.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l5-ledger-workflow stateless repetition 2: token usage increased 8.0%; time overhead -4.8%; 0 drill-downs. [Transcript](l5-ledger-workflow-2-stateless.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l5-ledger-workflow stateful repetition 2: token usage increased 24.0%; time overhead 18.2%; 0 drill-downs. [Transcript](l5-ledger-workflow-2-stateful.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l1-clamp stateless repetition 3: token usage increased 9.8%; time overhead -10.1%; 0 drill-downs. [Transcript](l1-clamp-3-stateless.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l2-config integrated repetition 3: token usage increased 41.2%; time overhead 23.5%; 0 drill-downs. [Transcript](l2-config-3-integrated.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l2-config stateless repetition 3: token usage increased 44.8%; time overhead 51.1%; 0 drill-downs. [Transcript](l2-config-3-stateless.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l2-pagination integrated repetition 3: token usage increased 1.6%; time overhead -3.9%; 0 drill-downs. [Transcript](l2-pagination-3-integrated.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l3-ledger integrated repetition 3: token usage increased 15.4%; time overhead -6.9%; 0 drill-downs. [Transcript](l3-ledger-3-integrated.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l3-ledger stateless repetition 3: token usage increased 15.1%; time overhead 0.3%; 0 drill-downs. [Transcript](l3-ledger-3-stateless.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l3-ledger stateful repetition 3: token usage increased 14.6%; time overhead -9.0%; 0 drill-downs. [Transcript](l3-ledger-3-stateful.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l3-planner integrated repetition 3: token usage increased 15.2%; time overhead -2.8%; 0 drill-downs. [Transcript](l3-planner-3-integrated.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l3-planner stateless repetition 3: token usage increased 14.3%; time overhead -8.1%; 0 drill-downs. [Transcript](l3-planner-3-stateless.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l4-store integrated repetition 3: token usage increased 15.6%; time overhead 2.7%; 0 drill-downs. [Transcript](l4-store-3-integrated.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l4-store stateless repetition 3: token usage increased 13.9%; time overhead 23.5%; 0 drill-downs. [Transcript](l4-store-3-stateless.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l5-ledger-workflow integrated repetition 3: token usage increased 8.1%; time overhead 8.6%; 0 drill-downs. [Transcript](l5-ledger-workflow-3-integrated.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l5-ledger-workflow stateless repetition 3: token usage increased 1.1%; time overhead -1.1%; 0 drill-downs. [Transcript](l5-ledger-workflow-3-stateless.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.
- coding/l5-ledger-workflow stateful repetition 3: token usage increased 12.1%; time overhead 14.0%; 0 drill-downs. [Transcript](l5-ledger-workflow-3-stateful.log). Commands and failures are retained in aggregate.json. This association does not establish compression as the cause.

## Analysis

Enabled trials recorded 0 delta presentations, 0 unchanged presentations, and 0 show/raw calls. There is no observed drill-down cost in these trials.

Native tools other than Bash contributed 194552/593947 textual tool-result bytes (32.8%). They are outside the current shell integration. Interception coverage and native-output share limit the workflow benefit achievable by shell reduction.

The paired level table shows whether reductions increase with intended task complexity. The observed L5 command histories and repeated-command counts are in aggregate.json; intended complexity alone is not evidence of longer actual context history. Stateful effectiveness is limited by the number of observed unchanged/delta presentations, and cannot be inferred from the mode name. Quality is judged by the same protected tests and vet checks in both modes.

## Evidence and limitations

Raw JSONL and stderr are retained in each .log file. Canonical acap bench compare reports are retained for each complete pair. Tool-output bytes exclude images and internal system prompts. Cache token accounting is explicit, but billed cost is not inferred. Repair iterations are unavailable. Exact repeated command strings, build/test/vet Bash invocations and failed validation tool results are extracted into aggregate.json; native background notifications are not classified as validation attempts. Shell commands count Bash invocations, not individual commands in compound scripts.

The complexity levels describe intended engineering scope, not empirically calibrated difficulty. L5 requests meaningful layer-by-layer validation, without artificial repetition or output padding. Workflow length and repeated commands remain agent decisions. Three repetitions are the default; a smoke run supports infrastructure conclusions only.

## Recommendations

1. Use all-tool bytes and provider token usage for release decisions. Shell-only byte reduction can conceal bypassed and native tool costs measured above.
2. Investigate the recorded shell-state bypass command sequences before broadening adapter coverage. Preserve shell semantics; do not change compression algorithms based on these results alone.
3. Target follow-up experiments at the levels with token or time regressions in the paired table. Retain the same success criteria and use more repetitions before treating variability as a stable effect.
