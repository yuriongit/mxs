# Demos

Terminal demo GIFs for MXS, recorded with [VHS](https://github.com/charmbracelet/vhs).

Each important command has its own directory with a `.tape` file

## Generate the GIFs

Requires Docker and Go. From the repo root:

```
./demos/generate_gifs.sh
```

Set `CPUS` to cap each container, e.g. `CPUS=1.5 ./demos/generate_gifs.sh`.

## How it works

The script builds a static Linux mxs binary and copies it into an image based on `ghcr.io/charmbracelet/vhs`. Each tape then runs in its own container, all in parallel, with `./demos` bind-mounted so the GIFs land on the host. Each container has its own filesystem, so `~/.mxs` inside it never touches a personal .mxs directory - every tape starts from a clean state.

The script exits non-zero and names the tape if any recording fails.