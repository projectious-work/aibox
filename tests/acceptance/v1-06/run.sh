#!/usr/bin/env sh
# V1-06 acceptance suite (settings-and-logging@v1).
#   run.sh --aibox PATH --evidence DIR
# Exit 0: every journey passed; 1: a journey failed or was not run;
# 2: usage or infrastructure error. See runner.py for details.
set -eu
suite_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd -P)
command -v python3 >/dev/null 2>&1 || { echo "run.sh: python3 is required" >&2; exit 2; }
exec python3 -I "$suite_dir/runner.py" "$@"
