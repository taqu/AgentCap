# Infrastructure smoke pair

Both trials passed automated verification and demonstrated real Claude Code
execution. The enabled trial demonstrated AgentCap interception.

This smoke pair preceded the final accounting corrections: its token fields use
top-level Claude usage (excluding separately reported auxiliary model usage),
and its enabled command count covers intercepted commands. The main evaluation
counts every Bash invocation equally and includes all per-model usage. Preserve
this pair as infrastructure evidence; do not pool its metrics with main trials
or use it to conclude an efficiency benefit.
