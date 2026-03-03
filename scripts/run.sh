#!/bin/sh
set -e

if command -v mux-session >/dev/null 2>&1; then
    mux-session "$@"
else
    echo "mux-session not found, installing..."
    SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
    "$SCRIPT_DIR/../install.sh"
    mux-session "$@"
fi
