#!/usr/bin/env bash
set -euo pipefail

export HOST_UID=$(id -u)
export HOST_GID=$(id -g)

docker compose build
docker compose run --rm goreleaser

python3 scripts/gen_changelog.py
