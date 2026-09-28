# AgentCap Benchmark CLI Roadmap

## Goal

新しい CLI を、たとえば次の形で用意します。

```bash
acap bench ...
```

目的は単なる reducer の圧縮率測定ではなく、

```text
raw workflow cost
vs
AgentCap stateless cost
vs
AgentCap stateful cost
vs
integrated workflow cost
```

を**同じ workload に対して再現可能に比較すること**です。

最終的な主 KPI は既存 roadmap と同じく、

```text
Total agent-visible tokens/bytes
over the complete coding workflow
```

とします。`show` / `raw` による情報回収もコストに含めます。:chatgpt-content-reference{index="2"}

---

## Phase B0 — Benchmark Measurement Foundation

**目的:** benchmark CLI を作る前に、何をどう測定するかを固定する。

まず内部に共通の measurement model を導入します。

```go
type BenchmarkMeasurement struct {
    RawBytes           int64
    AgentVisibleBytes  int64
    Duration           time.Duration

    CommandCount       int
    ShowCount          int
    RawCount           int

    AcapProcessingTime time.Duration
}
```

実際の型名は既存 architecture に合わせます。

重要なのは、少なくとも次を別々に記録できることです。

```text
raw command output
capsule output
stateful/delta output
show output
raw retrieval output
AgentCap processing overhead
```

この Phase では benchmark runner を作り込みません。

**Exit criteria:**

```text
同じ command execution について

raw cost
AgentCap-visible cost
recovery cost
latency

を機械的に取得できる
```

---

## Phase B1 — `acap bench command`

最初のユーザー向け benchmark CLI を作ります。

```bash
acap bench command -- <command...>
```

例:

```bash
acap bench command -- git diff
acap bench command -- go test ./...
acap bench command -- cargo test
```

出力イメージ:

```text
Benchmark: git diff

raw:
  bytes:          184,231

agentcap:
  visible_bytes:   12,842
  reduction:        93.0%
  processing:        4.8ms

result:
  id: abc123
```

ここでは **stateless な単一 command benchmark** に限定します。

### 重要な制約

benchmark のために command を2回実行してはいけません。

Phase 6 の重要 invariant、

> Every intercepted command must execute exactly once.

を benchmark にも維持します。:chatgpt-content-reference{index="3"}

1 execution から raw と compressed の双方を測定します。

### Exit criteria

```bash
acap bench command -- <command>
```

だけで raw / reduced size / reduction / AgentCap overhead を取得できる。

---

# Phase B2 — Repeated-Command / Stateful Benchmark

次に AgentCap の大きな差別化である stateful compression を測定します。

例えば:

```bash
acap bench session
```

または workload ファイルを使って:

```bash
acap bench run benchmark.yaml
```

概念的には、

```yaml
steps:
  - run: git diff
  - run: go test ./...
  - run: git diff
  - run: go test ./...
```

のような workload を実行します。

測定するのは:

```text
raw total
stateless total
stateful total
```

です。

出力例:

```text
Commands: 12

                     bytes
raw                428,120
stateless           71,840
stateful            28,410

stateful vs raw       -93.4%
stateful vs stateless -60.4%

unchanged results: 3
delta results:     5
full results:      4
```

これにより roadmap の

```text
raw cost
vs
stateless AgentCap cost
vs
stateful AgentCap cost
```

を直接測定できます。:chatgpt-content-reference{index="4"}

### Exit criteria

複数 command からなる workflow に対して、

```text
raw
stateless
stateful
```

の total agent-visible bytes を比較できる。

---

# Phase B3 — Workload Specification

B2 の ad-hoc benchmark を再現可能な benchmark suite にします。

例えば:

```text
benchmarks/
├── git/
│   ├── repeated-diff.yaml
│   └── status-edit-status.yaml
├── build/
│   ├── compile-fix.yaml
│   └── warnings.yaml
└── test/
    ├── fail-fix-pass.yaml
    └── repeated-pass.yaml
```

workload は command だけではなく、状態遷移を表現できる必要があります。

例えば:

```text
compile
  -> 7 errors

apply fixture change

compile
  -> 3 errors

apply fixture change

compile
  -> PASS
```

ここは AgentCap にとって重要です。

単に同じ静的 command を100回実行する benchmark では、Phase 5 の

```text
7 errors
 -> 3 errors
 -> 1 error
 -> PASS
```

という価値を正しく測れません。:chatgpt-content-reference{index="5"}

### Safety

fixture repository / temporary workspace を利用し、ユーザーの実 repository を benchmark のために変更しない設計にします。

---

# Phase B4 — Benchmark Result Format

人間向け出力と machine-readable output を分離します。

```bash
acap bench run benchmarks/git/repeated-diff.yaml
```

通常は compact table。

そして:

```bash
acap bench run benchmarks/git/repeated-diff.yaml --json
```

で structured result。

例えば:

```json
{
  "schema_version": 1,
  "workload": "git/repeated-diff",
  "commands": 8,
  "raw_bytes": 240182,
  "stateless_bytes": 41820,
  "stateful_bytes": 16201,
  "show_bytes": 0,
  "raw_retrieval_bytes": 0,
  "processing_ns": 4812000
}
```

これにより CI や外部 analysis tool から利用できます。

### 重要

benchmark result schema は versioning します。

後から metric を追加しても既存 benchmark データを解釈できるようにします。

---

# Phase B5 — Recovery-Cost Benchmark

ここで初めて progressive disclosure のコストを本格的に評価します。

AgentCap の KPI は capsule size だけでは不十分です。

ロードマップにもある通り、

```text
initial command results
+
all acap show calls
+
all acap raw calls
+
integration instruction overhead
```

まで含める必要があります。:chatgpt-content-reference{index="6"}

そこで benchmark measurement に、

```text
initial_visible_bytes
show_bytes
raw_retrieval_bytes
total_visible_bytes
```

を明示的に持たせます。

結果:

```text
Initial AgentCap output       18,420
show retrieval                4,810
raw retrieval                     0
                              ------
total agent-visible           23,230

raw baseline                 381,920

effective reduction            93.9%
```

特に、

```text
95% initial reduction
+
frequent acap raw
```

を「成功」と判定しないことが重要です。これは現在の roadmap の failure metric と直接対応します。:chatgpt-content-reference{index="7"}

---

# Phase B6 — Coding-Agent Workflow Benchmark

ここで benchmark CLI を実際の coding agent evaluation に拡張します。

対象:

```text
bug fixing
test failure repair
small feature
refactoring
compile-error repair
Git review
repository search
```

既存 roadmap の Phase 7 evaluation とここで合流します。

ただし `acap` 自身が coding agent を実装する必要はありません。

構造としては:

```text
Benchmark Runner
      |
      +---- workload
      |
      +---- agent adapter
      |        |
      |        +-- Claude Code
      |        +-- Codex
      |
      +---- AgentCap measurement
      |
      +---- result
```

Phase 6 と同じように agent-specific logic を adapter に閉じ込めます。

---

## B6 の比較モード

ここで初めて正式に:

```text
A. disabled
B. stateless
C. stateful
D. integrated
```

を比較します。

既存 roadmap の benchmark matrix をそのまま実行可能な形にする Phase です。:chatgpt-content-reference{index="8"}

例えば概念上:

```bash
acap bench agent \
    --workload bugfix-001 \
    --mode disabled

acap bench agent \
    --workload bugfix-001 \
    --mode integrated
```

ただし実際の CLI syntax は B6 実装時に agent adapter の制約を見て確定します。

---

# Phase B7 — Repetition and Statistical Comparison

LLM agent benchmark は非決定的なので、一回の結果だけで判断しないようにします。

```bash
acap bench agent \
    --workload bugfix-001 \
    --mode integrated \
    --repeat 5
```

などの repeated trials を導入します。

集計対象:

```text
success count
median agent-visible bytes
median wall time
median command count
show/raw rate
AgentCap overhead
```

個々の trial も保存します。

平均だけに潰さないことが重要です。

---

# Phase B8 — `acap bench compare`

保存した benchmark runs を比較する command を追加します。

例えば:

```bash
acap bench compare baseline.json candidate.json
```

または benchmark directory 単位で比較します。

出力イメージ:

```text
                         disabled     stateful
Task success               5/5          5/5
Median visible bytes      412k          91k
Median commands             31           32
Median wall time          84.2s        86.1s
Raw fallback                 -          0.2/run
```

ここでは「compression ratio が良いから勝ち」と自動判定しません。

**task success / information recovery / context cost / latency の各軸をそのまま表示**します。

---

# Phase B9 — CI Regression Benchmark

最後に benchmark CLI を development feedback loop に組み込みます。

用途:

```text
Reducer変更
     |
     v
benchmark
     |
     +--> output regression?
     +--> latency regression?
     +--> recovery regression?
```

例えば CI で lightweight deterministic benchmark を実行し、

```text
raw bytes
capsule bytes
stateful bytes
processing latency
```

の大きな regression を検出できるようにします。

LLM を使う B6/B7 の full agent benchmark は、通常 CI とは分離した方がよいです。

---

# 推奨実装順

全体をまとめると、

```text
B0  Measurement foundation
 |
B1  Single-command benchmark
 |
B2  Stateful workflow benchmark
 |
B3  Reproducible workload format
 |
B4  Stable JSON/result format
 |
B5  Progressive-disclosure recovery cost
 |
+-------------------------------+
| ここまで AgentCap-only        |
+-------------------------------+
 |
B6  Real coding-agent benchmark
 |
B7  Repeated trials / aggregation
 |
B8  Benchmark comparison
 |
B9  CI regression benchmark
```

という境界を置くのがよいです。

**B0〜B5だけでもかなり価値があります。** Agent integration や LLM の非決定性を持ち込む前に、AgentCap の「本当に workflow 全体で output を減らしているのか」を deterministic に測定できるからです。

その後 B6 で初めて、「出力量は減ったが coding agent の task performance を悪化させていないか」という、現在の roadmap が最終的に答えたい問いに進めます。ロードマップ自身も、Phase 6 後の中心課題を「機能追加」ではなく、実際の agent efficiency / task performance の検証としています。:chatgpt-content-reference{index="9"}

個人的には、次の実装単位は **B0+B1 を一つの Phase として coding agent に渡す**のがちょうどいいです。B1 は既存の execution/capture/reducer/store をかなり再利用でき、benchmark infrastructure の設計を早い段階で実データを使って検証できます。
