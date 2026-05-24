# Codex Hook プロトタイプ

`codex-approval-gate codex` は stdin から Codex `PermissionRequest` hook payload を読み取り、stdout に hook response を書き出す。

## コマンド

```sh
codex-approval-gate codex --config approval-gate.toml
```

プロトタイプでは `--config` flag が必須。

## Decision

approval gate は次のいずれかを返す。

- `allow`
- `deny`
- `ask`

provider error、timeout、不正な JSON、未知の decision、曖昧な結果はすべて `ask` を返す。

## デフォルト出力

```json
{
  "hookSpecificOutput": {
    "hookEventName": "PermissionRequest",
    "permissionDecision": "ask"
  }
}
```

## Simple Output

`codex.output_mode = "simple"` を設定すると次を出力する。

```json
{"decision":"ask"}
```

## Audit Log

`audit.path` が設定されている場合、command は decision ごとに JSON object を 1 行 JSONL file に append する。raw input はデフォルトでは省略され、`audit.include_raw_input = true` の場合のみ含める。
