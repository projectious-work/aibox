#!/usr/bin/env sh
# Host launcher for the V1-05 AC-TOOLS acceptance suite. Requires Python 3 and
# Docker; the runner bootstraps the pinned Dev Container CLI itself.
set -eu
suite_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd -P)
exec python3 "$suite_dir/runner.py" "$@"
