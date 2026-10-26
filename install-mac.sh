#!/usr/bin/env bash
# Installs MXS on macOS (Apple Silicon arm64).
set -euo pipefail

REPO="yuriongit/mxs"
VERSION="0.3.0"
BIN="mxs_${VERSION}_darwin_arm64.tar.gz"
URL="https://github.com/${REPO}/releases/download/v${VERSION}/${BIN}"
INSTALL_DIR="/usr/local/bin"
DEST="${INSTALL_DIR}/mxs"

# 1. OS & Architecture Check
OS="$(uname -s)"
ARCH="$(uname -m)"

if [ "$OS" != "Darwin" ] || [ "$ARCH" != "arm64" ]; then
  echo "Error: This installer is intended for macOS on Apple Silicon (arm64) only." >&2
  echo "Detected OS: $OS, Architecture: $ARCH" >&2
  exit 1
fi

# 2. Check for required tools
if ! command -v tar >/dev/null 2>&1; then
  echo "Error: 'tar' is required to extract MXS." >&2
  exit 1
fi

# Download Helper Function (Tries curl first, falls back to wget)
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

# 3. Create Temp Directory & Archive File
TMP_DIR="$(mktemp -d)"
TMP_TAR="${TMP_DIR}/archive.tar.gz"

# Clean up temporary directory on script exit
trap 'rm -rf "$TMP_DIR"' EXIT

# 4. Download & Extract Archive
echo "Downloading MXS ${VERSION} for macOS (arm64)..."
download_file "$URL" "$TMP_TAR"

echo "Extracting archive..."
tar -xzf "$TMP_TAR" -C "$TMP_DIR"

if [ ! -f "${TMP_DIR}/mxs" ]; then
  echo "Error: Expected 'mxs' binary was not found inside the archive." >&2
  exit 1
fi

chmod +x "${TMP_DIR}/mxs"

# 5. Install binary to /usr/local/bin
echo "Installing binary to ${DEST}..."

# Ensure target directory exists
if [ ! -d "$INSTALL_DIR" ]; then
  if [ -w "$(dirname "$INSTALL_DIR")" ]; then
    mkdir -p "$INSTALL_DIR"
  else
    sudo mkdir -p "$INSTALL_DIR"
  fi
fi

# Move binary with appropriate permissions
if [ -w "$INSTALL_DIR" ]; then
  mv "${TMP_DIR}/mxs" "$DEST"
else
  sudo mv "${TMP_DIR}/mxs" "$DEST"
fi

echo "Installation complete!"
echo "Run 'mxs help' or 'mxs init' to get started."
