# Sandboxed connectivity attempt

Both smoke conditions were actually attempted. Claude exhausted API retries and
reported `API Error: Unable to connect to API (ConnectionRefused)` before any
coding work. These are infrastructure failures, not evidence about AgentCap
efficiency. Synthetic zero usage in the original JSON is not a model token
measurement. Logs are preserved without rewriting the observations.

The unrestricted smoke rerun is in `../smoke-unrestricted`; the main evaluation
is in `../main`. This failed connectivity attempt is excluded from the main
comparison.
