#!/usr/bin/env bash
# Installs MXS on Linux (x86_64).
set -euo pipefail

REPO="yuriongit/mxs"
VERSION="0.2.0"
BIN="mxs_${VERSION}_linux_amd64.tar.gz"
URL="https://github.com/${REPO}/releases/download/v${VERSION}/${BIN}"
INSTALL_DIR="/usr/local/bin"
DEST="${INSTALL_DIR}/mxs"

# 1. OS & Architecture Check
OS="$(uname -s)"
ARCH="$(uname -m)"

if [ "$OS" != "Linux" ] || [ "$ARCH" != "x86_64" ]; then
  echo "Error: This installer is intended for Linux x86_64 only." >&2
  echo "Detected OS: $OS, Architecture: $ARCH" >&2
  exit 1
fi

# 2. Download Helper Function (Tries curl first, falls back to wget)
download_file() {
  local url="$1"
  local dest="$2"

  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url" -o "$dest"
  elif command -v wget >/dev/null 2>&1; then
    wget -qO "$dest" "$url"
  else
    echo "Error: Neither 'curl' nor 'wget' is installed on this system." >&2
    echo "Please install curl or wget and try again." >&2
    exit 1
  fi
}

# 3. Create Temp File
TMP="$(mktemp)"
trap 'rm -f "$TMP"' EXIT

# 4. Download & Set Permissions
echo "Downloading MXS ${VERSION}..."
download_file "$URL" "$TMP"
chmod +x "$TMP"

# 5. Install to /usr/local/bin (Elevate with sudo only if needed)
echo "Installing binary to ${DEST}..."
if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP" "$DEST"
else
  sudo mv "$TMP" "$DEST"
fi

echo "Installation complete!"
echo "Run 'mxs help' or 'mxs init' to get started."
