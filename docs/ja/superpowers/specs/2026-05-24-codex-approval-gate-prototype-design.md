# codex-approval-gate プロトタイプ設計

日付: 2026-05-24

## 目的

`codex-approval-gate` の最初の最小プロトタイプを作る。これは Codex の `PermissionRequest` hook に対する approval gate として動作する、Go 製の単一 CLI である。

このプロトタイプは stdin から Codex hook JSON を受け取り、OpenAI-compatible な local LLM endpoint に判定を依頼し、stdout に Codex hook response を出力する。必要に応じて JSONL audit log も記録する。失敗、無効な応答、曖昧な結果は必ず `ask` を返す。

## 初期コマンド

最初に対応するコマンドは次の通り。

```sh
codex-approval-gate codex --config approval-gate.toml
```

プロトタイプでは `--config` を必須にする。デフォルト設定ファイルの探索は後続の iteration に回す。

## パッケージ構成

```text
cmd/codex-approval-gate/
internal/
  adapters/codex/
  audit/
  config/
  judge/
  policy/
  providers/openai/
examples/
  codex/hooks.json
  approval-gate.toml
docs/
  en/
    codex-hook.md
  ja/
    codex-hook.md
```

パッケージ境界は意図的に小さく保つ。

- `adapters/codex` は Codex hook の入力と出力の形を扱う。
- `config` は TOML 設定を読み込み、必須項目を検証する。
- `judge` は provider の出力を `allow`、`deny`、`ask` のいずれかに正規化する。
- `providers/openai` は OpenAI-compatible な `/v1/chat/completions` endpoint を呼び出す。
- `audit` は設定されている場合に、判定ごとに 1 つの JSON object を書き込む。
- `policy` は prompt と rule の発展に備えた置き場とし、プロトタイプでは薄く保つ。

## Codex Adapter

入力 adapter は hook JSON を寛容に decode する。tool name、command、cwd、reason などの一般的な field が存在する場合は抽出しつつ、audit と prompt のために raw input を保持する。

出力 adapter は 2 つの response format に対応する。

デフォルトの Codex hook format:

```json
{
  "hookSpecificOutput": {
    "hookEventName": "PermissionRequest",
    "permissionDecision": "ask"
  }
}
```

simple compatibility format:

```json
{"decision":"ask"}
```

出力 mode は config で選ぶ。デフォルトは Codex hook format とする。

## Judge の挙動

judge は設定済み provider に短い prompt を送る。provider には decision と任意の reason を含む structured JSON を返すよう指示する。

有効な decision は次の 3 つ。

- `allow`
- `deny`
- `ask`

provider error、timeout、不正な JSON、decision の欠落、未知の decision、曖昧な text の場合、judge は必ず `ask` を返す。この fail-to-human の挙動を public contract とする。

## OpenAI-Compatible Provider

provider は最小限の chat completions request を実装する。

- `POST {base_url}/v1/chat/completions`
- config の model
- judge/policy 由来の system message と user message
- 指定された environment variable に API key がある場合はそれを使う
- config の request timeout

プロトタイプでは streaming、tool calls、retry、embeddings、provider 固有の拡張は不要。

## Config

TOML file はプロトタイプ用に次の field に対応する。

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

audit path が空の場合、audit logging は無効になる。

## Audit Log

有効な場合、audit log は JSONL として書き込む。各行には次の情報を含める。

- timestamp
- decision
- reason がある場合は reason
- provider type と model
- 要約された request field
- raw input hash
- `include_raw_input = true` の場合のみ raw input

audit write に失敗した場合、危険な action を許可してはならない。audit logging が有効で log を書き込めない場合、command は `ask` を返す。

## テスト

実装は fixture-based test から始める。初期ケースとして次を必須にする。

- 有効な Codex hook input と provider の `allow` により、Codex hook response の `allow` が返る
- provider error では `ask` が返る
- malformed または未知の provider decision では `ask` が返る
- simple output mode では `{"decision":"ask"}` が返る
- audit enabled では JSONL record が 1 行書き込まれる

fixture は、それを使う package test の近くに置く。CLI integration test では temporary config file と `httptest` の OpenAI-compatible endpoint を使える。

## 対象外

プロトタイプでは、完全な policy rule language、ccgate compatibility、Claude-specific adapter、plugin packaging、install script、chat 以外の provider protocol、background daemon は扱わない。

