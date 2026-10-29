#!/usr/bin/env sh
# One full AC-TOOLS host run for one architecture:
#   tests/acceptance/v1-05/run-host.sh debian-13-amd64|debian-13-arm64 [runner options]
# Starts a registry on localhost:5000, publishes every aibox-owned Feature with
# scripts/v1-05-features.sh, runs every AC-TOOLS item on this host's target and
# writes evidence.json, summary.txt and logs to a new evidence directory
# (default: $HOME/aibox-v1-05-evidence/<target>-<UTC time>). The target must match
# the Docker server architecture; the runner refuses a mismatch.
set -eu
usage() { echo "usage: $0 debian-13-amd64|debian-13-arm64 [--output ABS_DIR] [runner options]" >&2; exit 2; }
[ $# -ge 1 ] || usage
target=$1; shift
case "$target" in debian-13-amd64|debian-13-arm64) ;; *) usage ;; esac
suite_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd -P)
case " $* " in
  *" --output "*|*" --output="*) exec python3 "$suite_dir/runner.py" --target "$target" "$@" ;;
esac
output="${HOME:?}/aibox-v1-05-evidence/$target-$(date -u +%Y%m%dT%H%M%SZ)"
exec python3 "$suite_dir/runner.py" --target "$target" --output "$output" "$@"
