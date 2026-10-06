#!/usr/bin/env bash
set -euo pipefail

NAME="${1:-}"

if [[ -z "$NAME" ]]; then
  echo "✗ Error: Name is required" >&2
  echo "Usage: xs demo <your-first-name>" >&2
  exit 1
fi

echo "• Starting demo script"

sleep 1

echo "Hello, $NAME!"

sleep 1

echo "Bye now, $NAME :("

sleep 1

echo "✓ Demo script completed successfully"
