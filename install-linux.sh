#!/usr/bin/env bash
# Installs MXS on Linux (x86_64).
set -euo pipefail

REPO="yuriongit/mxs"
VERSION="0.3.0"
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

# 2. Check for tar
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

# Clean up temporary directory and archive on script exit
trap 'rm -rf "$TMP_DIR"' EXIT

# 4. Download & Extract Archive
echo "Downloading MXS v${VERSION}..."
download_file "$URL" "$TMP_TAR"

echo "Extracting archive..."
tar -xzf "$TMP_TAR" -C "$TMP_DIR"

# Verify the binary exists inside the unpacked files
if [ ! -f "${TMP_DIR}/mxs" ]; then
  echo "Error: Expected 'mxs' binary was not found inside the archive." >&2
  exit 1
fi

chmod +x "${TMP_DIR}/mxs"

# 5. Install binary to /usr/local/bin
echo "Installing binary to ${DEST}..."
if [ -w "$INSTALL_DIR" ]; then
  mv "${TMP_DIR}/mxs" "$DEST"
else
  sudo mv "${TMP_DIR}/mxs" "$DEST"
fi

echo "Installation complete!"
echo "Run 'mxs help' or 'mxs init' to get started."
