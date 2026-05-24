# codex-approval-gate

`codex-approval-gate` は、Codex `PermissionRequest` hook を OpenAI-compatible な local LLM approval gate に渡すための小さな Go CLI prototype です。

最初の command は次の通りです。

```sh
codex-approval-gate codex --config examples/approval-gate.toml
```

stdin から Codex hook JSON を読み取り、stdout に approval response を出力します。有効な decision は `allow`、`deny`、`ask` です。provider failure、invalid provider output、ambiguous decision は必ず `ask` を返します。

## 試す

```sh
go run ./cmd/codex-approval-gate codex --config examples/approval-gate.toml < examples/codex/hooks.json
```

設定された provider は OpenAI-compatible な `/v1/chat/completions` endpoint を公開している必要があります。

## Documentation

- English: [docs/en/codex-hook.md](docs/en/codex-hook.md)
- English policy: [docs/en/policy.md](docs/en/policy.md)
- English smoke test: [docs/en/smoke-test.md](docs/en/smoke-test.md)
- Japanese: [docs/ja/codex-hook.md](docs/ja/codex-hook.md)
- Japanese policy: [docs/ja/policy.md](docs/ja/policy.md)
- Japanese smoke test: [docs/ja/smoke-test.md](docs/ja/smoke-test.md)
