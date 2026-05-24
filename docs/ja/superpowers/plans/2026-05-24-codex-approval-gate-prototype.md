# codex-approval-gate プロトタイプ実装計画

> **Agentic worker 向け:** REQUIRED SUB-SKILL: この計画を task ごとに実装するには superpowers:subagent-driven-development（推奨）または superpowers:executing-plans を使うこと。進捗管理には checkbox（`- [ ]`）形式を使う。

**Goal:** `codex-approval-gate codex` の最初の Go CLI プロトタイプを作る。

**Architecture:** Codex adapter、TOML config、judge、OpenAI-compatible provider、audit writer、CLI wiring の細い vertical slice を実装する。デフォルト response format は Codex hook output とし、config で simple output も使えるようにする。

**Tech Stack:** Go、standard library、`github.com/BurntSushi/toml`、`httptest`、fixture-based tests。

---

## File Structure

- `go.mod`: module metadata と TOML dependency。
- `cmd/codex-approval-gate/main.go`: CLI entrypoint と `codex` subcommand。
- `internal/adapters/codex/codex.go`: 寛容な input decoding と response encoding。
- `internal/adapters/codex/testdata/*.json`: input/output fixture。
- `internal/config/config.go`: TOML loading と validation。
- `internal/judge/judge.go`: provider decision normalization。
- `internal/providers/openai/client.go`: OpenAI-compatible chat completions client。
- `internal/audit/audit.go`: JSONL audit record writer。
- `examples/approval-gate.toml`: 最小 config example。
- `examples/codex/hooks.json`: Codex hook input example。
- `README.md`: quick usage。
- `docs/en/codex-hook.md` と `docs/ja/codex-hook.md`: hook behavior docs。

### Task 1: Go Module And Codex Adapter

**Files:**
- Create: `go.mod`
- Create: `internal/adapters/codex/codex.go`
- Create: `internal/adapters/codex/codex_test.go`
- Create: `internal/adapters/codex/testdata/permission_request.json`

- [ ] **Step 1: Write failing adapter tests**

fixture input が decode できること、default response が Codex hook output を出すこと、simple mode が `{"decision":"ask"}` を出すことをテストする。

- [ ] **Step 2: Run adapter tests and verify RED**

Run: `go test ./internal/adapters/codex`

Expected: module と package がまだ存在しないため FAIL。

- [ ] **Step 3: Implement minimal adapter**

`go.mod` を作り、JSON を raw map と summary field に decode し、`allow` / `deny` / `ask` を encode する。

- [ ] **Step 4: Run adapter tests and verify GREEN**

Run: `go test ./internal/adapters/codex`

Expected: PASS。

### Task 2: Config Loader

**Files:**
- Create: `internal/config/config.go`
- Create: `internal/config/config_test.go`

- [ ] **Step 1: Write failing config tests**

valid TOML、provider field の不足、`codex` output mode の default、timeout parsing をテストする。

- [ ] **Step 2: Run config tests and verify RED**

Run: `go test ./internal/config`

Expected: config package がまだ存在しないため FAIL。

- [ ] **Step 3: Implement config loader**

`github.com/BurntSushi/toml` を使い、provider type/base URL/model を必須にし、timeout は 10 秒、output mode は `codex` を default にする。

- [ ] **Step 4: Run config tests and verify GREEN**

Run: `go test ./internal/config`

Expected: PASS。

### Task 3: Judge And Provider Interface

**Files:**
- Create: `internal/judge/judge.go`
- Create: `internal/judge/judge_test.go`

- [ ] **Step 1: Write failing judge tests**

provider の `allow`、provider error、malformed provider output、unknown decision normalization をテストする。

- [ ] **Step 2: Run judge tests and verify RED**

Run: `go test ./internal/judge`

Expected: judge package がまだ存在しないため FAIL。

- [ ] **Step 3: Implement minimal judge**

request/response type、provider interface、prompt construction、fail-to-`ask` normalization を定義する。

- [ ] **Step 4: Run judge tests and verify GREEN**

Run: `go test ./internal/judge`

Expected: PASS。

### Task 4: OpenAI-Compatible Provider

**Files:**
- Create: `internal/providers/openai/client.go`
- Create: `internal/providers/openai/client_test.go`

- [ ] **Step 1: Write failing provider tests**

`httptest` を使い、request path、model、authorization header、JSON decision parsing、non-2xx error を検証する。

- [ ] **Step 2: Run provider tests and verify RED**

Run: `go test ./internal/providers/openai`

Expected: provider package がまだ存在しないため FAIL。

- [ ] **Step 3: Implement minimal client**

`/v1/chat/completions` に POST し、model/messages を送り、first choice message content を judge JSON として parse し、HTTP/provider failure では error を返す。

- [ ] **Step 4: Run provider tests and verify GREEN**

Run: `go test ./internal/providers/openai`

Expected: PASS。

### Task 5: Audit Writer

**Files:**
- Create: `internal/audit/audit.go`
- Create: `internal/audit/audit_test.go`

- [ ] **Step 1: Write failing audit tests**

disabled audit が何もしないこと、enabled audit が JSONL record を 1 行書くこと、raw input が設定時だけ含まれることをテストする。

- [ ] **Step 2: Run audit tests and verify RED**

Run: `go test ./internal/audit`

Expected: audit package がまだ存在しないため FAIL。

- [ ] **Step 3: Implement audit writer**

decision ごとに JSON object を 1 行 append し、timestamp、decision、reason、provider/model、summary、raw input hash、任意の raw input を含める。

- [ ] **Step 4: Run audit tests and verify GREEN**

Run: `go test ./internal/audit`

Expected: PASS。

### Task 6: CLI Wiring And Documentation

**Files:**
- Create: `cmd/codex-approval-gate/main.go`
- Create: `cmd/codex-approval-gate/main_test.go`
- Create: `examples/approval-gate.toml`
- Create: `examples/codex/hooks.json`
- Create: `README.md`
- Create: `docs/en/codex-hook.md`
- Create: `docs/ja/codex-hook.md`

- [ ] **Step 1: Write failing CLI integration tests**

temporary config、fixture stdin、`httptest` provider を使う。allow output、provider error fallback to ask、simple output mode、audit write をテストする。

- [ ] **Step 2: Run CLI tests and verify RED**

Run: `go test ./cmd/codex-approval-gate`

Expected: CLI wiring がまだ存在しないため FAIL。

- [ ] **Step 3: Implement CLI wiring**

`codex --config` を parse し、config を load し、stdin を読み、judge 経由で provider を呼び、audit を書き、response を stdout に encode する。audit write が失敗した場合は `ask` を出力する。

- [ ] **Step 4: Add examples and bilingual docs**

example config、hook fixture、README quickstart、対応する英語/日本語 hook docs を追加する。

- [ ] **Step 5: Run full verification**

Run: `go test ./...`

Expected: PASS。

- [ ] **Step 6: Commit**

Run:

```sh
git add .
git commit -m "Implement initial codex approval gate prototype"
```

