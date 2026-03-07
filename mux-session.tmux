#!/usr/bin/env bash

CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

run_plugin() {
  echo "Running plugin"
  if [[ -n "$TMUX_VERSION" ]]; then
    major="$(echo "$TMUX_VERSION" | cut -d. -f1)"
    minor="$(echo "$TMUX_VERSION" | cut -d. -f2)"
    if [[ "$major" -ge 3 ]] && [[ "$minor" -ge 2 ]]; then
      tmux display-popup -E -w 100% -h 100% "$CURRENT_DIR/scripts/run.sh"
    else
      tmux run-shell -b "$CURRENT_DIR/scripts/run.sh"
    fi
  else
    tmux run-shell -b "$CURRENT_DIR/scripts/run.sh"
  fi
}

MUX_SESSION_KEY="${MUX_SESSION_KEY:-M}"

tmux bind-key "$MUX_SESSION_KEY" run-shell "tmux neww $CURRENT_DIR/scripts/run.sh"

