#!/bin/sh
# Installs MXS on macOS (Apple Silicon).
# Downloads one binary from GitHub Releases and puts it in /usr/local/bin.
set -eu

REPO="yuriongit/mxs"
VERSION="v0.1.0-beta"
BIN="mxs-${VERSION}-darwin-arm64"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${BIN}"

if [ "$(uname -s)" != "Darwin" ] || [ "$(uname -m)" != "arm64" ]; then
  echo "This script is for macOS on Apple Silicon (arm64) only." >&2
  exit 1
fi

TMP="$(mktemp)"
trap 'rm -f "$TMP"' EXIT

echo "Downloading $URL"
curl -fsSL "$URL" -o "$TMP"
chmod +x "$TMP"

echo "Installing to /usr/local/bin/mxs (may ask for your password)"
sudo mkdir -p /usr/local/bin
sudo mv "$TMP" /usr/local/bin/mxs

echo "Installation complete"
echo "Run: 'mxs help' for help or 'mxs init' to initialize MXS"
