#!/usr/bin/env sh
# Run the V1-05 native Feature lifecycle gate without a host npm installation.
set -eu
script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd -P)
exec python3 "$script_dir/verify-v1-05-host.py" "$@"
