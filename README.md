# XS

XS, pronounced as _"**ex • es | x • s**"_, is a command-line tool for script execution
and management.

_Soon, some AI-assisted managerial operations will be implemented._

## Features

1. Global configuration directory support with offered automatic setup
2. Script execution with reported success and errors
3. Structured colored output

## Configuration

XS's configuration directory, `~/.xs`, is currently the home for a user's scripts.

- `~/.xs` directory layout document (coming soon): [/docs/xs.md](./docs/xs.md)

## Infrastructure

| Layer   | Tool                 |
| ------- | -------------------- |
| Main    | Go, Cobra, BubbleTea |
| Tooling | golangci-lint, Go    |
| UI      | Bubbles, Lipgloss    |

## Requirements

- Go 1.27.1 or later

## Quick Start

Clone the repo:

```bash
git clone https://github.com/yuriongit/xs.git
cd xs
```

Build the binary:

```bash
go build
go install
```

Start:

```bash
xs help
```

Init XS and run the demo script:

```bash
xs init
xs demo "Your FirstName"
```

## Images

To see a preview of XS, view the [/docs/preview.md](./docs/preview.md) document.

_Images directory - [/.github/images](./github/images)_

## Docs

- [/docs/architecture.md](./docs/architecture.md)
- [/docs/planned.md](./docs/planned.md)

## License

MIT
