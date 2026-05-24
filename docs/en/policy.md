# Minimal Policy Precheck

The policy precheck summarizes obvious command risk before the local LLM makes a decision.

It does not automatically allow or deny requests. The result is added to the judge prompt as context. Provider failures and ambiguous results still return `ask`.

## Verdicts

- `safe`: command matches a narrow read-only pattern.
- `risky`: command matches a destructive or privilege-changing pattern.
- `unknown`: command does not match a known pattern.

## Initial Patterns

Read-only examples:

- `git status`
- `git diff`
- `ls`
- `pwd`
- `cat`

Risky examples:

- `rm`
- `sudo`
- `chmod`
- `chown`
- shell redirection with `>`

The policy is intentionally conservative. Unknown commands are left to the provider and still fall back to `ask` on failure.

## Config

The built-in patterns remain enabled. Add local patterns in TOML:

```toml
[policy]
read_only_prefixes = ["go test"]
risky_prefixes = ["docker system prune"]
risky_substrings = ["| pbcopy"]
```

Configured risky patterns are evaluated before read-only patterns. This keeps destructive local overrides from being treated as safe.
