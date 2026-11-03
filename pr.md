## Pull Request Type

- [ ] Bugfix
- [ ] Feature
- [ ] Code style update (formatting, renaming)
- [x] Refactoring (no functional changes, no API changes)
- [ ] Build related changes
- [ ] Continuous Integration/Continuous Delivery
- [x] Documentation content changes
- [ ] Other (please describe):

## Summary

Harden AGENTS.md with security guardrails for git credentials, testing workflow requirements, PR instructions, and development environment tips.

## Changes

- Update AGENTS.md with credential security warnings and prohibition on changing git user details.
- Add comprehensive testing instructions: running test suites before commits and CI requirements.
- Define PR workflow including title format, lint check requirements, review process, and merge guidelines.
- Document code conventions for tests, script locations, version management, and commit message discipline.
- Clarify project layout with paths for main.go, cmd/, internal/app, internal/cfg, internal/executor, internal/scr, and internal/ui.
- Add emergency procedures for accidental main pushes, bad PR reverts, credential rotation, and security breach protocols.

## Bugfixes

### #77 — Preserve interleaved stdout and stderr output

- **Fix executor output race:** Merge the concurrent `stdout` and `stderr`
  readers into a single goroutine backed by a mutex-serialized line writer, so
  lines are forwarded in stream order without concurrent channel writes.
- **Simplify pipe handling:** Replace manual `StdoutPipe`/`StderrPipe` scanning
  with `cmd.Stdout`/`cmd.Stderr` + `cmd.Run()`, letting the runtime own pipe
  lifecycle and cleanly surface non-zero exit errors.
- Add tests for partial-line buffering and streamed line ordering, including
  interleaved output.
