# MXS

A simple & colored CLI to manage and execute your Bash scripts.

[mxs-cli.vercel.app](https://mxs-cli.vercel.app)

## Features

1. Script execution
2. Run a script using it's name (no file extension needed)
3. Global configuration directory
4. Execution status reports with exit codes
5. Structured and colored output
6. Argument support

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

## Build From Source

### Requirements

| Tool  | Type/Version |
| ----- | ------------ |
| Shell | Bash         |
| Go    | 1.27.1+      |

Clone the repo:

```bash
git clone https://github.com/yuriongit/mxs.git
cd mxs
```

Build the binary:

```bash
go build
```

Install the binary:

```bash
go install
```

Start:

```bash
mxs help
```

Initialize MXS:

```bash
mxs init
```

Run the demo:

```bash
mxs demo "Your FirstName"
```

## Images

To see a preview of MXS, view the [/docs/preview.md](./docs/preview.md) document.

_Images directory - [/.github/images](./github/images)_

## Documents

- [/docs/architecture.md](./docs/architecture.md)
- [/docs/planned.md](./docs/planned.md)

## License

MIT
