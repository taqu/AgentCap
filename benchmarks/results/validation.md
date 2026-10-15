# Implementation validation

- All ten fixtures fail their initial Go tests and pass tests plus vet after
  reference repairs applied only in disposable workspaces. Detailed command
  output is recorded in `../fixture-validation.json`.
- `go test ./internal/agentbench ./internal/workload ./internal/benchcompare`
  passed, including native tool-result parsing, permission parity, missing usage,
  auxiliary model usage and shell error-output accounting tests.
- The existing four-workload deterministic suite completed successfully. Its
  canonical result files are in `deterministic/`.
- `go test ./...` passed all packages except the unchanged
  `internal/integration/antigravity` on Windows. Failures concerned WSL/Linux cwd
  and environment handling, shell output expectations and file permissions.
  The suite was rerun outside the sandbox after sandboxed Bash access failures;
  those Antigravity failures remained. No Antigravity source was changed.
- `git diff --check` passed.
