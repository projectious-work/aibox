#!/bin/bash
# open-in-editor.sh — Open a file in the vim editor pane from yazi/file manager.
#
# Uses zellij actions to focus the editor pane and send :e <file> to vim.
# Assumes the editor pane is to the right of the file manager pane.

file="$1"
[ -z "$file" ] && exit 1

# Focus the pane to the right (editor pane)
zellij action move-focus right

# Ensure vim is in normal mode (send Escape), then open the file
zellij action write 27
sleep 0.05
zellij action write-chars ":e ${file}"
zellij action write 13
