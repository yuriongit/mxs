# AGENTS.md

Guidance for agents working on MXS.

## Project

Go CLI to manage and run Bash scripts. Cobra + BubbleTea + Lipgloss.
Module: `github.com/yuriongit/mxs`. Go 1.27.1.

## Commands

- Build: `go build`
- Test: `go test ./...`
- Lint: `golangci-lint run` (config: `.golangci.yml`, v2)
- Format: `gofmt -w .`

## Layout

- `main.go` — entry, calls `cmd.Execute()`
- `cmd/` — Cobra commands; `root` runs the TUI on `mxs <script> [args]`, plus `add`, `ls`, `la`, `init`
- `internal/app` — dependency wiring
- `internal/cfg` — `~/.mxs` config and directory handling
- `internal/executor` — script preparation/execution + TUI
- `internal/scr` — script template, description, listing
- `internal/ui` — Lipgloss styles and components

## Conventions

- Tests colocated as `*_test.go`.
- User scripts live in `~/.mxs/scripts/*.sh`.
- Version is 0.4.1 under development. Do not bump version numbers; the maintainer handles releases.
