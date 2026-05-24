# codex-approval-gate Prototype Design

Date: 2026-05-24

## Goal

Build the first minimal prototype of `codex-approval-gate`, a single Go CLI that acts as an approval gate for Codex `PermissionRequest` hooks.

The prototype accepts Codex hook JSON on stdin, asks an OpenAI-compatible local LLM endpoint for a decision, writes a Codex hook response to stdout, and optionally records a JSONL audit log. Any failure, invalid response, or ambiguous result must return `ask`.

## Initial Command

The first supported command is:

```sh
codex-approval-gate codex --config approval-gate.toml
```

For the prototype, `--config` is required. Default config discovery is deferred until a later iteration.

## Package Layout

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
  codex-hook.md
```

The package boundaries are intentionally small:

- `adapters/codex` handles Codex hook input and output shapes.
- `config` loads TOML configuration and validates required fields.
- `judge` normalizes provider output into `allow`, `deny`, or `ask`.
- `providers/openai` calls an OpenAI-compatible `/v1/chat/completions` endpoint.
- `audit` writes one JSON object per decision when configured.
- `policy` is a placeholder for prompt and rule evolution, kept thin in the prototype.

## Codex Adapter

The input adapter should decode hook JSON permissively. It should extract common fields such as tool name, command, cwd, and reason when present, while preserving the raw input for audit and prompting.

The output adapter supports two response formats:

Default Codex hook format:

```json
{
  "hookSpecificOutput": {
    "hookEventName": "PermissionRequest",
    "permissionDecision": "ask"
  }
}
```

Simple compatibility format:

```json
{"decision":"ask"}
```

The config selects the output mode. The default is the Codex hook format.

## Judge Behavior

The judge sends a compact prompt to the configured provider. The provider should be instructed to return structured JSON with a decision and optional reason.

Valid decisions are:

- `allow`
- `deny`
- `ask`

The judge must return `ask` for provider errors, timeouts, malformed JSON, missing decisions, unknown decisions, and ambiguous text. This fail-to-human behavior is part of the public contract.

## OpenAI-Compatible Provider

The provider implements the minimum chat completions request:

- `POST {base_url}/v1/chat/completions`
- model from config
- system and user messages from judge/policy
- API key from a configured environment variable when present
- request timeout from config

The prototype does not need streaming, tool calls, retries, embeddings, or provider-specific extensions.

## Config

The TOML file should support the prototype fields:

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

If audit path is empty, audit logging is disabled.

## Audit Log

When enabled, the audit log writes JSONL. Each line should include:

- timestamp
- decision
- reason, when available
- provider type and model
- summarized request fields
- raw input hash
- raw input only when `include_raw_input = true`

Audit write failures should not allow a risky action. The command should return `ask` if audit logging is enabled but the log cannot be written.

## Tests

Implementation starts with fixture-based tests. Required initial cases:

- valid Codex hook input plus provider `allow` returns Codex hook response `allow`
- provider error returns `ask`
- malformed or unknown provider decision returns `ask`
- simple output mode returns `{"decision":"ask"}`
- audit enabled writes one JSONL record

Fixtures should live near the package tests that use them. CLI integration tests can use temporary config files and an `httptest` OpenAI-compatible endpoint.

## Out of Scope

The prototype will not include a full policy rule language, ccgate compatibility, Claude-specific adapters, plugin packaging, install scripts, non-chat provider protocols, or background daemons.
