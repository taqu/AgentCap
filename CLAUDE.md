# AgentCap

AgentCap is active in this project. Shell command output is automatically compressed.

Compact results include an ID like `@acap 81bc2f`. Use targeted drill-down before fetching raw output:

- `acap show <id>` — full stored capsule
- `acap show <id> --errors` — build errors only
- `acap show <id> --warnings` — build warnings only
- `acap show <id> --test <name>` — specific test failure
- `acap show <id> --file <path>` — specific file diff
- `acap show <id> --hunk <N>` — specific diff hunk
- `acap show <id> --match <text>` — search stored output
- `acap show <id> --lines X:Y` — line range
- `acap raw <id>` — full raw output (use as last resort)

Prefer `acap show` with a selector over `acap raw`.

<!-- agentcap -->
## AgentCap

AgentCap is active in this project. Shell command output is automatically compressed.

Compact results include an ID like `@acap 81bc2f`. Use targeted drill-down before fetching raw output:

- `acap show <id>` — full stored capsule
- `acap show <id> --errors` — build errors only
- `acap show <id> --warnings` — build warnings only
- `acap show <id> --test <name>` — specific test failure
- `acap show <id> --file <path>` — specific file diff
- `acap show <id> --hunk <N>` — specific diff hunk
- `acap show <id> --match <text>` — search stored output
- `acap show <id> --lines X:Y` — line range
- `acap raw <id>` — full raw output (use as last resort)

Prefer `acap show` with a selector over `acap raw`.
