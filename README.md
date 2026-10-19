# MXS

**Current Version:** `v0.2.0-beta` (under development)

A simple & colored CLI to manage and execute your Bash scripts.

[mxs-cli.vercel.app](https://mxs-cli.vercel.app)

## Features

1. Script execution
2. Run a script using it's name (no file extension needed)
3. List your scripts with `mxs ls` (only `.sh` files, shown without the extension)
4. Global configuration directory
5. Execution status reports with exit codes
6. Structured and colored output, with separate opening and closing messages and
   logs in between
7. Argument support
8. Automated installers for Linux and macOS

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
mxs demo World
# mxs demo "Your FirstName"
```

## Preview

To see a preview of MXS, visit [mxs-cli.vercel.app/preview](https://mxs-cli.vercel.app/preview)

## Documents

- [/docs/architecture.md](./docs/architecture.md)
- [/docs/planned.md](./docs/planned.md)

## License

MIT
