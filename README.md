# codex-approval-gate

`codex-approval-gate` is a small Go CLI prototype for routing Codex `PermissionRequest` hooks through an OpenAI-compatible local LLM approval gate.

The first command is:

```sh
codex-approval-gate codex --config examples/approval-gate.toml
```

It reads Codex hook JSON from stdin and writes an approval response to stdout. Valid decisions are `allow`, `deny`, and `ask`. Provider failures, invalid provider output, and ambiguous decisions always return `ask`.

## Try It

```sh
go run ./cmd/codex-approval-gate codex --config examples/approval-gate.toml < examples/codex/hooks.json
```

The configured provider must expose an OpenAI-compatible `/v1/chat/completions` endpoint.

## Documentation

- English: [docs/en/codex-hook.md](docs/en/codex-hook.md)
- Japanese: [docs/ja/codex-hook.md](docs/ja/codex-hook.md)
