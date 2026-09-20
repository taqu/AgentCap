## AgentCap Roadmap

| Phase | 目的 | 主な対象 |
|---|---|---|
| 0 | 最小実行基盤 | `acap run` |
| 1 | coreutils圧縮 | `ls/find/grep/rg/cat/tree/du/wc` |
| 2 | Capsule + drill-down | `acap show/raw` |
| 3 | stateful差分圧縮 | セッション履歴・重複除去 |
| 4 | Git対応 | `status/diff/log/show` |
| 5 | Build/Test対応 | Go/C/C++/Rust等 |
| 6 | Agent統合 | Claude Code / Codex / hooks / MCP |

### Phase 0 — Execution Core

まずは「コマンドを安全に代理実行できる」だけを完成させます。

```bash
acap run ls -la
acap run rg Workspace .
acap run cat src/main.go
```

内部は単純で構いません。

```text
argv
 ↓
executor
 ↓
stdout + stderr + exit code
 ↓
generic reducer
 ↓
stdout
```

この段階では高度な意味解析は不要です。

必要機能は、exit codeの完全保持、stdout/stderr分離、signal伝播、timeout、大容量出力のspill-to-disk、TTY判定、環境変数・cwd継承です。

特に重要なのは、

> **AgentCapを通したせいでコマンドの意味が変わらない**

ことです。

最初から、

```go
type ExecutionResult struct {
    Command   []string
    ExitCode  int
    StdoutRef string
    StderrRef string
    Duration  time.Duration
}
```

のような内部表現を持たせます。

**Phase 0の完了条件:** `acap run <cmd>` が普通のshell executionの代替として安定して使える。

---

## Phase 1 — Coreutils Compression

ここが最初の実用版です。

対象をかなり絞ります。

```text
ls
find
grep
rg
cat
head
tail
tree
du
wc
```

ただし、「単純truncate」は極力しません。

### `ls`

例えば1000ファイルなら、

```text
@acap id=a12f
entries=1042 files=891 dirs=151

dirs:
  src/
  tests/
  docs/
  vendor/
  ...

largest:
  build.log  18.2MB
  data.bin    7.3MB

1042 entries; full output available
```

### `find`

```text
@acap id=f8ca
files=5182 dirs=316

by-root:
 src       824
 tests    1256
 vendor   2811
 docs      291

ext:
 .go   381
 .cpp  179
 .h    204
```

### `rg / grep`

これはかなり重要です。

例えば、

```bash
rg "Workspace" .
```

の300ヒットを、

```text
@acap id=8dd1 matches=317 files=42

src/workspace.go       81
src/resolver.go        47
tests/workspace_test.go 92
...

top matches:
src/workspace.go:41:type Workspace struct {
src/workspace.go:89:func NewWorkspace(...)
...
```

にする。

### `cat`

最初はAST不要です。

長いファイルについて、

```text
@acap id=ac21
src/workspace.go
lines=1842 bytes=61K

preview:
1-80
...
1760-1842

omitted: 1681 lines
```

程度でよいです。

後々semantic source readingへ進化できます。

**Phase 1の最重要KPIは圧縮率ではありません。**

```text
raw bytes
raw estimated tokens
returned bytes
returned estimated tokens
```

を全部記録します。

例えば、

```text
saved_tokens = 82.4%
```

を `acap stats` で見られるようにする。

---

## Phase 2 — Capsule / Drill-down

ここからAgentCapらしくなります。

すべての実行結果にIDを付けます。

```bash
$ acap run find . -type f

@acap 82f1
files=5182
...
```

そして、

```bash
acap show 82f1
acap raw 82f1
```

を実装します。

さらに、

```bash
acap show 82f1 --path src
acap show 82f1 --match workspace
acap show 82f1 --from 100 --to 200
```

と掘れるようにします。

内部としては、

```text
Execution
  ├── Raw stdout
  ├── Raw stderr
  ├── Metadata
  ├── Capsule
  └── Index
```

を保存。

つまりAgentCapの本質である、

```text
capture
 ↓
capsule
 ↓
drill-down
```

が完成します。

ここまで来たら、初めてREADMEで

> Lossless progressive disclosure

を前面に出せます。

---

## Phase 3 — Stateful Compression

ここが普通の「output compressor」と差別化できるフェーズです。

Agentのセッションを持たせます。

```bash
acap session start
```

あるいは暗黙的に、

```text
ACAP_SESSION_ID
```

で管理。

記録するのは、

```text
command history
result IDs
output hashes
file fingerprints
already-emitted data
```

です。

例えば、

```bash
rg "Workspace" .
```

を2回実行して結果がほぼ同じなら、

```text
@acap d214
same as 93ef
changes:
 + src/new_workspace.go:81
 - src/old_workspace.go:57
```

だけ返します。

さらに、

```bash
git status
```

が同じなら、

```text
@acap unchanged since c19a
```

だけ。

ここで初めて、

> **モデルが既に見た情報を再送しない**

というAgentCap独自の思想が成立します。

名前をつけるなら、

**Session-aware Delta Compression**

あたりが分かりやすいです。

---

# Phase 4 — Git

coreutilsの次はGitです。

ここはCoding Agentで非常に呼び出し頻度が高いので、効果が大きい。

優先順は、

```text
git status
git diff
git diff --stat
git show
git log
git branch
```

です。

特に `git diff`。

巨大diffなら最初は、

```text
@acap diff=e71a

11 files
+482 -193

src/
 workspace.go    +201 -74
 parser.go        +84 -31

tests/
 workspace_test.go +185 -80
```

まで。

必要なら、

```bash
acap show e71a --file src/workspace.go
```

で本文。

さらにPhase 3と組み合わせると、

```text
git diff #1
 ↓
agent sees it

code modification

git diff #2
 ↓
AgentCap computes delta against diff #1

only newly modified hunks
```

というかなり強い機能になります。

---

# Phase 5 — Build / Test / Compiler

ここでAgentCapの実用価値が一気に上がります。

最初は全部対応せず、

```text
go test
go build

gcc
g++
clang
clang++

cargo build
cargo test
```

程度。

Node.js系は後回しでもよいと思います。

あなたの狙いなら、

> **Go単一バイナリを置くだけ**

という体験を壊さないほうが大事です。

### Compiler

1000行のビルドログを、

```text
@acap build 2fc1 exit=1

errors=3 warnings=17

E src/parser.cpp:481
  no matching function for call to parse(...)

E src/workspace.cpp:918
  use of undeclared identifier 'ctx'

W src/foo.cpp:82
  unused variable 'result'

17 warnings grouped
```

にする。

同一テンプレートエラーなどはまとめる。

### Tests

```text
tests=1284
pass=1277
fail=7

failures:
 ParserTest.InvalidToken        3
 WorkspaceTest.Resolve         2
 IndexTest.Lookup              2

first failure:
 ...
```

とする。

失敗時は成功ログをかなり削る。

成功時ならさらに、

```text
PASS 1284 tests (18.4s)
```

だけでいい。

これはかなり削減効果が高いです。

---

# Phase 6 — Agent Integration

ここまではAgentCap単体で完成させます。

そのあと、

```text
Claude Code
Codex
Gemini CLI
OpenCode
その他agent
```

へ接続。

ただし専用統合ロジックをcoreに入れません。

```text
Agent
 ↓
adapter
 ↓
acap core
```

にします。

例えば、

```text
internal/
  agent/
    claude/
    codex/
```

ではなく、可能なら外部設定だけにする。

AgentCap coreは、

```text
stdin
argv
stdout
stderr
filesystem
```

しか知らない状態が理想です。

### 最後にMCP

MCPはここで初めて追加します。

例えば、

```text
acap.execute
acap.show
acap.search_result
acap.raw
```

というToolとして公開。

しかしMCPを**AgentCapそのもの**にはしません。

CLIが本体です。

---

# アーキテクチャも早期に固定した方がいい

おすすめはこれです。

```text
cmd/acap
   │
   ▼
Command Router
   │
   ▼
Executor
   │
   ├─────────── Raw Store
   │
   ▼
Classifier
   │
   ▼
Reducer
   │
   ▼
Capsule
   │
   ▼
Renderer
```

Reducerをinterface化します。

```go
type Reducer interface {
    Match(cmd Command) bool
    Reduce(ctx context.Context, r *ExecutionResult) (*Capsule, error)
}
```

そして、

```text
GenericReducer
LsReducer
FindReducer
RipgrepReducer
CatReducer
GitDiffReducer
GoTestReducer
CompilerReducer
```

と追加していく。

これはかなり重要です。

**「acapに特別対応されていないコマンドも必ず実行できる」**設計にします。

つまり、

```text
known command
→ specialized reducer

unknown command
→ generic reducer
```

です。

未知のCLIでも、

```text
first N important lines
last N lines
repeated-line collapse
ANSI removal
progress-line collapse
```

くらいはできる。

---

# Release milestone

プロジェクトとして切るなら、こんな感じが現実的です。

```text
v0.1
Execution core
Generic reducer
ls/find/rg
Token statistics

v0.2
Result store
acap show/raw
cat/tree/du
Progressive disclosure

v0.3
Session state
Hashing
Cross-command dedup
Delta results

v0.4
Git support
status/diff/show/log

v0.5
Build/Test
Go + GCC/Clang + Cargo

v0.6
Agent hooks
Claude Code / Codex

v0.7
Semantic source compression
symbol-aware file views

v0.8
AST integration

v0.9
MCP

v1.0
Stable output contract
Plugin/reducer API
Benchmark suite
```

ここで特に重要なのは、**ASTをv0.8くらいまで意図的に入れないこと**です。

AgentCapが成功するかどうかは、ASTの賢さより先に、

```text
$ curl ...
$ install acap
$ acap ...
```

だけで本当にトークンが減るか、

で判断できます。

---

## 最初の1本として作るべきMVP

さらに絞るなら、最初のMVPはこれだけで十分です。

```text
acap run
acap show
acap raw
acap stats

reducers:
  generic
  ls
  find
  rg
```

保存先：

```text
~/.cache/agentcap/
```

結果：

```text
~/.cache/agentcap/results/<id>/
    meta.json
    stdout
    stderr
    capsule
```

この時点でベンチマークを作ります。

```text
Linux kernel
LLVM
Kubernetes
Rust compiler
medium-sized Go repo
```

などでcoding-agentがよく使うコマンドを流し、

```text
Raw tokens
Returned tokens
Reduction %
Execution overhead
Information recovery calls
```

を計測する。

**最後の `Information recovery calls` が特に重要です。**

95%圧縮してもAgentが毎回、

```bash
acap raw ...
```

を呼ぶなら失敗です。

AgentCapで追うべき数字はむしろ、

> **Total session tokens saved after drill-down**

です。

これをプロジェクトの中心KPIにすると、単なる「grep出力を80%削減しました」という既存ツールとの比較から抜け出せます。

そしてこの順番なら、最初の開発ターゲットはかなり明確です。

**`v0.1 = Go単一バイナリ + generic/ls/find/rg + token accounting`**

ここから始めるのが一番いいと思います。
