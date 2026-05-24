# Codex Setup

This page describes how to wire `codex-approval-gate` into a Codex `PermissionRequest` hook.

## 1. Install The CLI

Install from a local checkout while the project is pre-release:

```sh
go install ./cmd/codex-approval-gate
```

Ensure the install directory is on `PATH`.

## 2. Configure The Local LLM Endpoint

Create `approval-gate.toml` in the directory where Codex will run, or create the user config file documented in [install.md](install.md).

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

The endpoint must support `POST /v1/chat/completions`. The first choice message content must be JSON:

```json
{"decision":"ask","reason":"not enough context"}
```

## 3. Configure The Codex Hook

Use the CLI as the command that receives hook JSON on stdin.

Example shape:

```json
{
  "hooks": {
    "PermissionRequest": {
      "command": "codex-approval-gate codex"
    }
  }
}
```

If you do not use config discovery, pass an explicit path:

```json
{
  "hooks": {
    "PermissionRequest": {
      "command": "codex-approval-gate codex --config /absolute/path/approval-gate.toml"
    }
  }
}
```

Adjust the exact hook file location and shape to the Codex configuration format you are using.

## 4. Recommended Audit Path

Use an absolute audit log path for real hook usage:

```toml
[audit]
path = "/Users/you/Library/Logs/codex-approval-gate/audit.jsonl"
include_raw_input = false
```

Enable `include_raw_input` only if you are comfortable storing commands, paths, and possible secrets in the audit log.

## Troubleshooting

- Config not found: pass `--config` or create `./approval-gate.toml`.
- Provider unreachable: the hook response falls back to `ask`.
- Provider returns non-JSON content: the hook response falls back to `ask`.
- Local HTTPS endpoint uses a self-signed certificate: prefer plain local HTTP for now, or put a trusted proxy in front.
