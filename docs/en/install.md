# Install

This project currently ships as a Go CLI.

## From Source

```sh
git clone https://github.com/toshikimiyagawa/codex-approval-gate.git
cd codex-approval-gate
go test ./...
go install ./cmd/codex-approval-gate
```

## With `go install`

After the module is published with a stable version, install with:

```sh
go install github.com/toshikimiyagawa/codex-approval-gate/cmd/codex-approval-gate@latest
```

Until release tags exist, prefer installing from a local checkout.

## Config Location

The CLI reads config from `--config` when provided.

If `--config` is omitted, it searches:

1. `./approval-gate.toml`
2. OS user config directory:
   - macOS: `~/Library/Application Support/codex-approval-gate/config.toml`
   - Linux: `~/.config/codex-approval-gate/config.toml`

## Minimal Config

```toml
[provider]
type = "openai"
base_url = "http://localhost:1234"
model = "local-model"
api_key_env = ""
timeout = "10s"

[codex]
output_mode = "codex"

[audit]
path = "./audit.jsonl"
include_raw_input = false
```

Use `api_key_env` when your local OpenAI-compatible endpoint requires a token.

## Verify

```sh
codex-approval-gate codex --config examples/approval-gate.toml < examples/codex/hooks.json
```

Provider failures, invalid provider JSON, and ambiguous responses return `ask`.
