# Install

この project は現時点では Go CLI として提供する。

## Source から install

```sh
git clone https://github.com/toshikimiyagawa/codex-approval-gate.git
cd codex-approval-gate
go test ./...
go install ./cmd/codex-approval-gate
```

## `go install` を使う

module が stable version として公開された後は、次の command で install できる。

```sh
go install github.com/toshikimiyagawa/codex-approval-gate/cmd/codex-approval-gate@latest
```

release tag が存在するまでは、local checkout から install する。

## Config Location

`--config` が指定されている場合、CLI はその config を読む。

`--config` が省略された場合、次の順に探す。

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

local OpenAI-compatible endpoint が token を要求する場合は `api_key_env` を使う。

## Verify

```sh
codex-approval-gate codex --config examples/approval-gate.toml < examples/codex/hooks.json
```

provider failure、不正な provider JSON、曖昧な response は `ask` を返す。
