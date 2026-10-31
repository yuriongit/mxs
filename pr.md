## Pull Request Type

- [ ] Bugfix
- [ ] Feature
- [ ] Code style update (formatting, renaming)
- [ ] Refactoring (no functional changes, no API changes)
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
