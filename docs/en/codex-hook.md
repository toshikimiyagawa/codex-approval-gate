# Codex Hook Prototype

`codex-approval-gate codex` reads a Codex `PermissionRequest` hook payload from stdin and writes a hook response to stdout.

## Command

```sh
codex-approval-gate codex --config approval-gate.toml
```

The `--config` flag is required in the prototype.
If omitted, the CLI looks for `./approval-gate.toml` and then `codex-approval-gate/config.toml` under the OS user config directory.

## Decisions

The approval gate returns one of:

- `allow`
- `deny`
- `ask`

Any provider error, timeout, invalid JSON, unknown decision, or ambiguous result returns `ask`.

## Default Output

```json
{
  "hookSpecificOutput": {
    "hookEventName": "PermissionRequest",
    "permissionDecision": "ask"
  }
}
```

## Simple Output

Set `codex.output_mode = "simple"` to write:

```json
{"decision":"ask"}
```

## Audit Log

When `audit.path` is set, the command appends one JSON object per decision to that JSONL file. Raw input is omitted by default and included only when `audit.include_raw_input = true`.
