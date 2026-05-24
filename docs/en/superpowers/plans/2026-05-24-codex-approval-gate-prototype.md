# codex-approval-gate Prototype Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the first Go CLI prototype for `codex-approval-gate codex`.

**Architecture:** Implement a narrow vertical slice: Codex adapter, TOML config, judge, OpenAI-compatible provider, audit writer, and CLI wiring. The default response format is Codex hook output; simple output remains available through config.

**Tech Stack:** Go, standard library, `github.com/BurntSushi/toml`, `httptest`, fixture-based tests.

---

## File Structure

- `go.mod`: module metadata and TOML dependency.
- `cmd/codex-approval-gate/main.go`: CLI entrypoint and `codex` subcommand.
- `internal/adapters/codex/codex.go`: permissive input decoding and response encoding.
- `internal/adapters/codex/testdata/*.json`: input/output fixtures.
- `internal/config/config.go`: TOML loading and validation.
- `internal/judge/judge.go`: provider decision normalization.
- `internal/providers/openai/client.go`: OpenAI-compatible chat completions client.
- `internal/audit/audit.go`: JSONL audit record writer.
- `examples/approval-gate.toml`: minimal config example.
- `examples/codex/hooks.json`: Codex hook input example.
- `README.md`: quick usage.
- `docs/en/codex-hook.md` and `docs/ja/codex-hook.md`: hook behavior docs.

### Task 1: Go Module And Codex Adapter

**Files:**
- Create: `go.mod`
- Create: `internal/adapters/codex/codex.go`
- Create: `internal/adapters/codex/codex_test.go`
- Create: `internal/adapters/codex/testdata/permission_request.json`

- [ ] **Step 1: Write failing adapter tests**

Test that fixture input decodes, default response emits Codex hook output, and simple mode emits `{"decision":"ask"}`.

- [ ] **Step 2: Run adapter tests and verify RED**

Run: `go test ./internal/adapters/codex`

Expected: FAIL because the module and package do not exist yet.

- [ ] **Step 3: Implement minimal adapter**

Create `go.mod`, decode JSON into raw maps plus summary fields, and encode `allow` / `deny` / `ask`.

- [ ] **Step 4: Run adapter tests and verify GREEN**

Run: `go test ./internal/adapters/codex`

Expected: PASS.

### Task 2: Config Loader

**Files:**
- Create: `internal/config/config.go`
- Create: `internal/config/config_test.go`

- [ ] **Step 1: Write failing config tests**

Test valid TOML, missing provider fields, default `codex` output mode, and timeout parsing.

- [ ] **Step 2: Run config tests and verify RED**

Run: `go test ./internal/config`

Expected: FAIL because config package is missing.

- [ ] **Step 3: Implement config loader**

Use `github.com/BurntSushi/toml`, require provider type/base URL/model, default timeout to 10 seconds, and default output mode to `codex`.

- [ ] **Step 4: Run config tests and verify GREEN**

Run: `go test ./internal/config`

Expected: PASS.

### Task 3: Judge And Provider Interface

**Files:**
- Create: `internal/judge/judge.go`
- Create: `internal/judge/judge_test.go`

- [ ] **Step 1: Write failing judge tests**

Test provider `allow`, provider error, malformed provider output, and unknown decision normalization.

- [ ] **Step 2: Run judge tests and verify RED**

Run: `go test ./internal/judge`

Expected: FAIL because judge package is missing.

- [ ] **Step 3: Implement minimal judge**

Define request/response types, provider interface, prompt construction, and fail-to-`ask` normalization.

- [ ] **Step 4: Run judge tests and verify GREEN**

Run: `go test ./internal/judge`

Expected: PASS.

### Task 4: OpenAI-Compatible Provider

**Files:**
- Create: `internal/providers/openai/client.go`
- Create: `internal/providers/openai/client_test.go`

- [ ] **Step 1: Write failing provider tests**

Use `httptest` to verify request path, model, authorization header, parsed JSON decision, and non-2xx error.

- [ ] **Step 2: Run provider tests and verify RED**

Run: `go test ./internal/providers/openai`

Expected: FAIL because provider package is missing.

- [ ] **Step 3: Implement minimal client**

POST to `/v1/chat/completions`, send model/messages, parse the first choice message content as judge JSON, and return errors for HTTP/provider failures.

- [ ] **Step 4: Run provider tests and verify GREEN**

Run: `go test ./internal/providers/openai`

Expected: PASS.

### Task 5: Audit Writer

**Files:**
- Create: `internal/audit/audit.go`
- Create: `internal/audit/audit_test.go`

- [ ] **Step 1: Write failing audit tests**

Test disabled audit does nothing, enabled audit writes one JSONL record, and raw input is included only when configured.

- [ ] **Step 2: Run audit tests and verify RED**

Run: `go test ./internal/audit`

Expected: FAIL because audit package is missing.

- [ ] **Step 3: Implement audit writer**

Append one JSON object per decision, include timestamp, decision, reason, provider/model, summary, raw input hash, and optional raw input.

- [ ] **Step 4: Run audit tests and verify GREEN**

Run: `go test ./internal/audit`

Expected: PASS.

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

Use temporary config, fixture stdin, and `httptest` provider. Test allow output, provider error fallback to ask, simple output mode, and audit write.

- [ ] **Step 2: Run CLI tests and verify RED**

Run: `go test ./cmd/codex-approval-gate`

Expected: FAIL because CLI wiring is missing.

- [ ] **Step 3: Implement CLI wiring**

Parse `codex --config`, load config, read stdin, call provider through judge, write audit, and encode response to stdout. If audit writing fails, output `ask`.

- [ ] **Step 4: Add examples and bilingual docs**

Add example config, hook fixture, README quickstart, and matching English/Japanese hook docs.

- [ ] **Step 5: Run full verification**

Run: `go test ./...`

Expected: PASS.

- [ ] **Step 6: Commit**

Run:

```sh
git add .
git commit -m "Implement initial codex approval gate prototype"
```

