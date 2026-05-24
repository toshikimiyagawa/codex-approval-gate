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
