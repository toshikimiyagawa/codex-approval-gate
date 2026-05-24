# Smoke Test

Codex adapter、provider client、prompt policy を変更した後は、この流れで確認する。

## 1. CLI を build または run する

```sh
go run ./cmd/codex-approval-gate codex --config examples/approval-gate.toml < examples/codex/hooks.json
```

この command は、`provider.base_url` に設定された OpenAI-compatible endpoint を期待する。

## 2. Local Endpoint を設定する

`examples/approval-gate.toml` を編集するか、local path に copy して使う。

```toml
[provider]
type = "openai"
base_url = "http://localhost:1234"
model = "local-model"
timeout = "10s"
```

provider は `POST /v1/chat/completions` を実装し、first choice message の content として次のような JSON を返す必要がある。

```json
{"decision":"ask","reason":"not enough context"}
```

## 3. 実 Hook Payload を capture する

Codex 本体と検証するときは、gate に渡す前の stdin payload を capture する。private path、command、secret が含まれる場合は git に入れない。

redact 後、代表的な payload を次の場所に追加する。

```text
internal/adapters/codex/testdata/
```

そのうえで、adapter test に次の summary field の assertion を追加する。

- tool name
- command
- cwd
- reason がある場合は reason

## 4. Verify

```sh
CGO_ENABLED=0 go test -count=1 ./...
```

`CGO_ENABLED=0` は、この project environment で確認された macOS test binary loader issue を避けるために使う。
