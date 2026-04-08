#!/bin/sh
# xdg-open shim for headless aibox containers.
#
# A real `xdg-open` tries to launch a browser via X/Wayland/desktop
# integration, none of which exists inside a dev container. Without
# this shim, every CLI tool that calls xdg-open (gh auth login, git
# credential helpers, opencode, claude code device-flow, ...) errors
# out with a confusing "Failed opening a web browser" message.
#
# This shim prints the URL with clear "open this on your host machine"
# instructions and exits 0, so the calling tool thinks the browser
# opened successfully and proceeds with whatever polling/waiting it
# was going to do (typically OAuth device-flow polling). The user has
# the URL in plain text and can paste it into the host browser.
#
# Exit 0 is deliberate: any non-zero status would make the calling
# tool fail or fall back to an even worse path. We are pretending the
# browser opened.
#
# See DEC-031.
set -e

if [ "$#" -lt 1 ]; then
    echo "xdg-open: missing URL argument" >&2
    exit 1
fi

URL="$1"

# Plain rendering — no box framing. URLs vary in length and a unicode
# box that doesn't pad to the URL width looks broken. A simple
# blank-line / heading / indented-URL / blank-line block stands out
# from any preceding tool output without requiring layout math.
printf '\n'
printf '\033[1;33m  ┃ Headless container — no browser available here.\033[0m\n'
printf '\033[1;33m  ┃ Open this URL in your host machine'"'"'s browser:\033[0m\n'
printf '\n'
printf '      \033[1;36m%s\033[0m\n' "$URL"
printf '\n'
printf '\033[2m  (the calling tool is now waiting for you to authorize in the browser)\033[0m\n'
printf '\n'

exit 0
