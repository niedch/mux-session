#!/bin/sh
set -e

# Repository and binary details
REPO="niedch/mux-session"
BINARY_NAME="mux-session"
INSTALL_DIR="${HOME}/.local/bin"

# Determine OS and architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

# Get the latest release tag
LATEST_TAG=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_TAG" ]; then
    echo "Could not fetch the latest release tag. Aborting."
    exit 1
fi

# Construct the download URL
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_TAG}/${BINARY_NAME}_${OS}_${ARCH}.tar.gz"

# Download and extract the binary
echo "Downloading ${BINARY_NAME} from ${DOWNLOAD_URL}..."
TMPDIR=$(mktemp -d)
curl -sL "${DOWNLOAD_URL}" -o "${TMPDIR}/${BINARY_NAME}.tar.gz"
tar xzf "${TMPDIR}/${BINARY_NAME}.tar.gz" -C "${TMPDIR}"

# Install the binary
echo "Installing ${BINARY_NAME} to ${INSTALL_DIR}..."
mkdir -p "${INSTALL_DIR}"
mv "${TMPDIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
rm -rf "${TMPDIR}"

echo "${BINARY_NAME} installed successfully!"
