# Smoke Test

Use this flow after changing the Codex adapter, provider client, or prompt policy.

## 1. Build Or Run The CLI

```sh
go run ./cmd/codex-approval-gate codex --config examples/approval-gate.toml < examples/codex/hooks.json
```

The command expects an OpenAI-compatible endpoint at the configured `provider.base_url`.

## 2. Configure A Local Endpoint

Edit `examples/approval-gate.toml` or copy it to a local path:

```toml
[provider]
type = "openai"
base_url = "http://localhost:1234"
model = "local-model"
timeout = "10s"
```

The provider must implement `POST /v1/chat/completions` and return a first choice message whose content is JSON:

```json
{"decision":"ask","reason":"not enough context"}
```

## 3. Capture Real Hook Payloads

When validating against Codex itself, capture the stdin payload before sending it to the gate. Keep captured examples out of git if they include private paths, commands, or secrets.

After redaction, add representative payloads under:

```text
internal/adapters/codex/testdata/
```

Then add adapter tests that assert the extracted summary fields:

- tool name
- command
- cwd
- reason, when present

## 4. Verify

```sh
CGO_ENABLED=0 go test -count=1 ./...
```

`CGO_ENABLED=0` avoids a macOS test binary loader issue observed in this project environment.
