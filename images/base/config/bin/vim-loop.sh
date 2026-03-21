#!/bin/bash
# vim-loop.sh — Run vim in a loop so the editor pane never dies.
# When vim exits (:q), it restarts with an empty buffer.
# Exit the loop with :cq (exit with error code) or Ctrl+C.

while true; do
    vim "$@"
    exit_code=$?
    # :cq exits with code 1 — use this to truly quit
    if [ "$exit_code" -ne 0 ]; then
        break
    fi
    # Normal :q — restart vim
done
