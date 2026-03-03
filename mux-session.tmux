#!/usr/bin/env bash

CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

run_plugin() {
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

tmux bind-key f run-shell -b "#{plug_current_dir}/scripts/run.sh"
