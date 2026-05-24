# Codex Setup

この page は `codex-approval-gate` を Codex `PermissionRequest` hook に接続する手順を説明する。

## 1. CLI を install する

project が pre-release の間は local checkout から install する。

```sh
go install ./cmd/codex-approval-gate
```

install 先 directory が `PATH` に入っていることを確認する。

## 2. Local LLM Endpoint を設定する

Codex を実行する directory に `approval-gate.toml` を作るか、[install.md](install.md) に記載された user config file を作る。

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

[policy]
read_only_prefixes = []
risky_prefixes = []
risky_substrings = []
```

endpoint は `POST /v1/chat/completions` に対応している必要がある。first choice message content は次のような JSON にする。

```json
{"decision":"ask","reason":"not enough context"}
```

## 3. Codex Hook を設定する

hook JSON を stdin で受け取る command として CLI を設定する。

設定例:

```json
{
  "hooks": {
    "PermissionRequest": {
      "command": "codex-approval-gate codex"
    }
  }
}
```

config discovery を使わない場合は明示的な path を渡す。

```json
{
  "hooks": {
    "PermissionRequest": {
      "command": "codex-approval-gate codex --config /absolute/path/approval-gate.toml"
    }
  }
}
```

実際の hook file の場所と形式は、使用している Codex configuration format に合わせて調整する。

## 4. Recommended Audit Path

実際の hook usage では absolute audit log path を使う。

```toml
[audit]
path = "/Users/you/Library/Logs/codex-approval-gate/audit.jsonl"
include_raw_input = false
```

`include_raw_input` は、command、path、secret の可能性がある情報を audit log に保存してよい場合だけ有効にする。

## Troubleshooting

- Config not found: `--config` を渡すか、`./approval-gate.toml` を作る。
- Provider unreachable: hook response は `ask` に fallback する。
- Provider が non-JSON content を返す: hook response は `ask` に fallback する。
- Local HTTPS endpoint が self-signed certificate を使う: 現時点では plain local HTTP を使うか、trusted proxy を前段に置く。
