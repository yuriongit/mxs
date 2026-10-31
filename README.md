# MXS

**Current Version:** `v0.4.1 (under development)`

A simple & colored CLI to manage and execute your Bash scripts.

[mxs-cli.vercel.app](https://mxs-cli.vercel.app)

## Features

1. Script execution: Run a script using its name (no file extension needed)
2. Script creation with the `add` command
3. List your scripts with the `ls` command (`.sh` files, shown without the extension)
4. Verbose listing with the `la` command, including each script's description
5. Global configuration directory: `~/.mxs`
6. Structured execution UI: live spinner, streamed logs, and execution time
7. Script descriptions shown alongside your script's output
8. Argument support
9. Automated installers for Linux and macOS

## Configuration

MXS' configuration directory, `~/.mxs`, is the home for a user's scripts.

For more info, view: [/docs/configuration-directory.md](./docs/configuration-directory.md)
(coming soon)

## Infrastructure

| Layer   | Tool                 |
| ------- | -------------------- |
| Main    | Go, Cobra, BubbleTea |
| Tooling | golangci-lint, Go    |
| UI      | Bubbles, Lipgloss    |

## Installation

To install MXS, visit [mxs-cli.vercel.app/install](https://mxs-cli.vercel.app/install).

Installers are available for Linux and macOS.

## Build From Source

### Requirements

For clarification, these requirements _only_ apply if you are building from source.

| Tool  | Type/Version |
| ----- | ------------ |
| Shell | Bash         |
| Go    | 1.27.1+      |

---

Clone the repo:

```bash
git clone https://github.com/yuriongit/mxs.git
cd mxs
```

Build and install the binary:

```bash
go build
go install
```

Initialize MXS:

```bash
mxs init
```

Run the demo:

```bash
mxs demo purple
# mxs demo "your favorite color"
```

## Preview

To see a preview of MXS, visit [mxs-cli.vercel.app/preview](https://mxs-cli.vercel.app/preview)

## Documents

- [Product Requirements Document](./docs/MVP-PRD.md)
- [Architecture](./docs/architecture.md) (coming soon)
- [Configuration Directory](./docs/configuration-directory.md) (coming soon)

## License

MIT
