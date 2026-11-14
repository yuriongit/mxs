#!/usr/bin/env bash
set -euo pipefail

# Records every demo tape in its own container, in parallel.
# Run from the repo root. Set CPUS to cap each container (e.g. CPUS=1.5).

image=mxs-demos
ctx=$(mktemp -d)
trap 'rm -rf "$ctx"' EXIT

CGO_ENABLED=0 GOOS=linux go build -o "$ctx/mxs" .
docker build -q -f demos/Dockerfile -t "$image" "$ctx"

pids=()
tapes=()
for tape in ./demos/*/*.tape; do
  docker run --rm \
    --user "$(id -u):$(id -g)" \
    -e HOME=/tmp \
    ${CPUS:+--cpus "$CPUS"} \
    -v "$PWD/demos:/vhs/demos" \
    "$image" "$tape" &
  pids+=($!)
  tapes+=("$tape")
done

fail=0
for i in "${!pids[@]}"; do
  if ! wait "${pids[$i]}"; then
    echo "failed: ${tapes[$i]}" >&2
    fail=1
  fi
done
exit "$fail"
