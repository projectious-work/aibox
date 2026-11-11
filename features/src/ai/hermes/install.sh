#!/bin/sh
# ai/hermes: hermes as a uv tool (https://hermes-agent.nousresearch.com),
# with a uv-managed CPython. uv verifies every wheel against the package
# index hashes; dependency resolution is frozen at the catalog date.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

af_uv_tool_install --python 3.13.15 hermes-agent==$AF_VERSION

af_finish
