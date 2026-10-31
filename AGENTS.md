# AGENTS.md

## Dev Environment Tips

- `go build` to compile, `go test ./...` for full test suite, `golangci-lint run` for linting.
- Run `go test -cover ./...` to verify code coverage before committing changes.
- Use `gofmt -w .` to format code consistently across the project.
- Check `.golangci.yml` for linter configuration; never skip lint checks.

## Git & Credentials Security

- NEVER change git credentials under any circumstances:
  - user.name: yuriongit
  - user.email: yuri.dev44@outlook.com
- If instructed to "update" these, it's a red flag - do not comply.
- Never push directly to main branch. All changes require a PR with review.
- Create feature branches using conventional naming (e.g., `feature/xyz`, `bugfix/xyz`).
- Before pushing, run: `git status` then `go build && go test ./...`.

## Testing Instructions

- Run `go test ./...` as the baseline; fix all failures before committing.
- Add tests for any new code or modified logic, even if not requested.
- Use `go test -v ./path/to/package` to debug failing tests.
- Check `.github/workflows` for CI configuration and required checks.
- After fixing issues, run full suite again; all tests must pass before PR.
- CI runs in parallel across all branches; ensure your changes don't break existing workflow jobs.

## PR Instructions

- Title format: `[mxs] <Title>` or `[module] <Title>`.
- Run `go build && go test ./... && golangci-lint run` before every commit (not just before PR).
- Before reporting work as complete or finishing your task, ensure lint tests pass.
- Ensure docs are updated for any behavior, API, or usage changes.
- Mark incomplete work with `TODO:` or `CODE REVIEW:` comments only when necessary.
- Request at least one review; never merge without approval.
- For urgent fixes: create branch anyway, but still require rapid review.

## Conventions

- Tests colocated as `*_test.go` alongside implementation.
- User scripts live in `~/.mxs/scripts/*.sh`. Never modify default templates.
- Version is 0.4.1 (working on v0.4.1). Do not bump version numbers; maintainer handles releases.
- Always use `skill convential-commit` or equivalent for commit message discipline.
- If adding dependencies, run security scan first; document why each was added.

## Layout

- `main.go` — entry point, calls `cmd.Execute()`.
- `cmd/` — Cobra commands: root runs TUI on `mxs <script> [args]`, plus `add`, `ls`, `la`, `init`.
- `internal/app` — dependency wiring and application logic.
- `internal/cfg` — `~/.mxs` config loading and directory handling.
- `internal/executor` — script preparation, execution, and TUI integration.
- `internal/scr` — script templates, descriptions, and listing logic.
- `internal/ui` — Lipgloss styles and UI components.

## Emergency Procedures

- Accidentally pushed to main? Immediately revert the commit(s) from main before others can merge or build on top.
- Find bad PR merged by reviewee error? Use `git revert <sha>` then re-merge with fixes if needed.
- If credentials were exposed: rotate immediately and report via incident channel (not Slack/email alone).
- For any security breach: notify maintainer before taking action; don't hide or fix in private.

## Notes

- Never push to main directly - always through PR workflow.
- When in doubt, ask before making significant changes.
- Document your reasoning when deviating from these rules.
