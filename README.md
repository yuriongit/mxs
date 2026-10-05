# Xs

Xs, pronounced as _"**ex • es**"_, is a command-line tool for script execution and
AI-assisted script management.

## Features

1. Global configuration directory that holds scripts
2. Script execution
3. Colored and structured output

## Configuration

Example paragraph that explains Xs's configuration.

- Configuration directory & options: [/.xs](./docs/.xs.md)

## Infrastructure

| Layer   | Tool                 |
| ------- | -------------------- |
| Main    | Go, Cobra, BubbleTea |
| Tooling | golangci-lint, Go    |
| UI      | Bubbles, Lipgloss    |
| Other   | Name                 |

## Requirements

- Go 1.27.1 or later

## Quick Start

Clone the repo & build the binary:

```bash
git clone https://github.com/yuriongit/xs.git
cd xs
go build
go install
```

Start:

```bash
xs --help
```

Init Xs & Run demo script:

```bash
xs init
xs demo
```

## Docs

- [/docs/architecture.md](./docs/architecture.md)
- [/docs/planned.md](./docs/planned.md)

## Images

- [/.github/images](./images)

## License

MIT
