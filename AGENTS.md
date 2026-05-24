# Repository Guidelines

## Documentation Languages

- Keep repository documentation under `docs/en/` and `docs/ja/`.
- The directory structure under `docs/en/` and `docs/ja/` must match.
- When adding or moving a documentation file, add or update both the English and Japanese versions in the matching relative path.
- If a document starts in one language, add the translated counterpart before considering the documentation change complete.
- Do not leave canonical documentation directly under `docs/` except for the `en/` and `ja/` language roots.

## Pull Requests

- Write pull request descriptions in Japanese.
- Keep PR descriptions concise and include a summary plus the test plan.
- Prefer PRs grouped by a single user-visible objective, not by every small commit.
- Use small commits inside a PR for reviewable steps, but keep the PR as the unit that answers "what does this enable?"
- Split PRs when changes have different review concerns, different rollback needs, or unrelated risk profiles.
- For larger feature work, write or update a short spec/plan before implementation and keep tasks in that plan small enough to verify independently.

## Plan Checkboxes

- Use plan checkboxes when a spec/plan exists for a feature-sized change.
- Mark a task checkbox when that task's implementation and verification are complete, not only at the end of the whole PR.
- Update checkboxes before pausing, handing off, opening a PR, or reporting that a planned task is complete.
- For small changes without a plan, do not create checkbox-only process overhead; summarize completed work and verification in the PR description instead.
- If a plan is intentionally not kept in sync, note that in the PR description so unfinished checkboxes are not ambiguous.
