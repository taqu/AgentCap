#!/usr/bin/env bash
# Phase 7 baseline batch driver. Orchestrates `acap bench agent` only; all
# measurement, isolation and verification live in the benchmark CLI.
#
#   benchmarks/tools/run_phase7_batch.sh <batch-id> [reps] [agents...]
#
# Output: benchmarks/results/<batch-id>/{manifest.txt,runs/*.json,evidence/,logs/}
# Then:   acap bench report benchmarks/results/<batch-id>
set -u
shopt -s extglob
BATCH=${1:?batch id}
REPS=${2:-3}
shift 2 || true
AGENTS=("$@")
[ ${#AGENTS[@]} -eq 0 ] && AGENTS=(claude antigravity)

REPO=$(cd "$(dirname "$0")/../.." && pwd)
OUT="$REPO/benchmarks/results/$BATCH"
mkdir -p "$OUT/runs" "$OUT/evidence" "$OUT/logs" "$OUT/bin"
cd "$REPO"

declare -A MODEL=([claude]=claude-sonnet-5-5 [antigravity]=gemini-3.8-flash-medium [codex]=)
MODES=(disabled integrated)
TASKS=(benchmarks/workloads/phase7/${TASK_GLOB:-*}.yaml)
SEED=${SEED:-7}

# Freeze the AgentCap build for the whole batch (§71).
go build -o "$OUT/bin/acap" ./cmd/acap || exit 1
ACAP="$OUT/bin/acap"

{
  echo "batch: $BATCH"
  echo "started: $(date -u +%FT%TZ)"
  echo "agentcap_commit: $(git rev-parse HEAD)"
  echo "agentcap_worktree_diff_sha256: $(git diff HEAD -- . ':!benchmarks/results' | sha256sum | cut -c1-16)"
  echo "platform: $(uname -srm)"
  echo "go: $(go version)"
  echo "cc: $(cc --version | head -1)"
  echo "claude: $(claude --version 2>/dev/null) model=${MODEL[claude]} permission_mode=bypassPermissions setting_sources=project"
  echo "antigravity: $(agy --version 2>/dev/null) model=${MODEL[antigravity]} --dangerously-skip-permissions"
  echo "modes: ${MODES[*]}  reps: $REPS  seed: $SEED"
  echo "timeout: per-workload (15m), identical for OFF and FULL"
  echo "cache_policy: warm (shared Go build cache, warmed once before the batch for both modes)"
  echo "tasks:"; for t in "${TASKS[@]}"; do echo "  $t $(sha256sum "$t" | cut -c1-16)"; done
} > "$OUT/manifest.txt"

# Warmup (§73): identical for all modes — prime the Go build cache.
for f in benchmarks/fixtures/p7-go-*; do (cd "$f" && go build ./... >/dev/null 2>&1; go vet ./... >/dev/null 2>&1); done

run_agent() {
  local agent=$1
  # Randomized, interleaved order of (task, mode, rep) per agent (§21).
  local plan
  plan=$(for r in $(seq 1 "$REPS"); do for t in "${TASKS[@]}"; do for m in "${MODES[@]}"; do echo "$r $t $m"; done; done; done |
         python3 -c "import random,sys; l=sys.stdin.read().split('\n')[:-1]; random.Random($SEED+len('$agent')).shuffle(l); print('\n'.join(l))")
  echo "$plan" > "$OUT/logs/$agent.order"
  local i=0
  while read -r rep task mode; do
    i=$((i+1))
    local name; name=$(basename "$task" .yaml)
    local dst="$OUT/runs/${agent}__${name}__${mode}__r${rep}.json"
    [ -s "$dst" ] && continue   # resumable
    for attempt in 1 2; do
      "$ACAP" bench agent --workload "$task" --agent "$agent" --mode "$mode" --model "${MODEL[$agent]}" \
        --batch "$BATCH" --evidence-dir "$OUT/evidence" --json > "$dst.tmp" 2> "$OUT/logs/${agent}__${name}__${mode}__r${rep}.a$attempt.err"
      local rc=$?
      local outcome; outcome=$(python3 -c "import json,sys; print(json.load(open(sys.argv[1]))['workflow']['outcome'])" "$dst.tmp" 2>/dev/null || echo NONE)
      echo "$(date -u +%T) [$agent $i] $name $mode r$rep attempt=$attempt rc=$rc outcome=$outcome" >> "$OUT/logs/progress.log"
      if [ "$rc" -eq 0 ] && [ "$outcome" != INFRA_ERROR ] && [ "$outcome" != NONE ]; then
        mv "$dst.tmp" "$dst"; break
      fi
      # Keep the failed attempt for review; rerun once (§43).
      mv "$dst.tmp" "$OUT/logs/${agent}__${name}__${mode}__r${rep}.a$attempt.json.failed" 2>/dev/null
      sleep 20
    done
  done <<< "$plan"
}

for a in "${AGENTS[@]}"; do run_agent "$a" & done
wait
echo "finished: $(date -u +%FT%TZ)" >> "$OUT/manifest.txt"
"$ACAP" bench report "$OUT" > "$OUT/report.md"
echo "report: $OUT/report.md"
