#!/bin/sh
# Installs MXS on Linux (x86_64).
# Downloads one binary from GitHub Releases and puts it in /usr/local/bin.
set -eu

REPO="yuriongit/mxs"
VERSION="v0.1.0-beta"
BIN="mxs-${VERSION}-linux-x86_64"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${BIN}"

if [ "$(uname -s)" != "Linux" ] || [ "$(uname -m)" != "x86_64" ]; then
  echo "This script is for Linux x86_64 only." >&2
  exit 1
fi

TMP="$(mktemp)"
trap 'rm -f "$TMP"' EXIT

echo "Downloading $URL"
curl -fsSL "$URL" -o "$TMP"
chmod +x "$TMP"

echo "Installing to /usr/local/bin/mxs (may ask for your password)"
sudo mv "$TMP" /usr/local/bin/mxs

echo "Installation complete"
echo "Run: 'mxs help' for help or 'mxs init' to initialize MXS"
