# MXS

<!--/*
  Add custom configured colors: If there's any errors
  in the config, they must first be resolved. Until
  then, a script will not run.

  Add this functionality to the cnf package.

  These following colors will have a default color.

  Also, in the cnf file, add a themes section:

  theme: {
    colors: {
      errors: ["#LIGHT", "#DARK"],
      execution_progress: ["#LIGHT", "#DARK"],
      successful_execution: ["#LIGHT", "#DARK"]
    },
    spinner: {
      extraspace: false,
      type: "dots | monkey | meter",
      color: "#HEX_CODE"
    },
    verbosity: {
      supplied_args: true,
      stopwatch: true,
    },
    logs: {
      quotes: true
    }
  }
*/-->

A simple & colored CLI to manage and execute your Bash scripts.

[mxs-cli.vercel.app](https://mxs-cli.vercel.app)

## Features

1. Global configuration directory support with offered automatic setup
2. Script execution with reported success and errors
3. Structured and colored output
4. Arguments support

## Configuration

MXS's configuration directory, `~/.mxs`, is the home for a user's scripts.

- `~/.mxs` directory layout document (coming soon): [/docs/mxs.md](./docs/mxs.md)

## Infrastructure

| Layer   | Tool                 |
| ------- | -------------------- |
| Main    | Go, Cobra, BubbleTea |
| Tooling | golangci-lint, Go    |
| UI      | Bubbles, Lipgloss    |

## Quick Start

Currently, releases aren't available. To try out MXS, build it from source:

### Requirements

- Go: 1.27.1+
- Bash: 5.3.9+
- Architecture: x86_64

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
