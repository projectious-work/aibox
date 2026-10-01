#!/usr/bin/env sh
# Host launcher: the Python verifier bootstraps a pinned upstream CLI without npm.
set -eu
script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd -P)
exec python3 "$script_dir/verify-v1-04-host.py" "$@"
