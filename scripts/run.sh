#!/usr/bin/env bash
set -e

MUX_INSTALL_DIR=/usr/local/bin

if command -v "$MUX_INSTALL_DIR/mux-session" >/dev/null 2>&1; then
    $MUX_INSTALL_DIR/mux-session
else
    echo "mux-session not found, installing..."
    SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
    "$SCRIPT_DIR/../install.sh"
    $MUX_INSTALL_DIR/mux-session
fi


