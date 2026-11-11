#!/bin/sh
# docs/zensical: zensical as a uv tool (https://zensical.org/docs/get-started/),
# with a uv-managed CPython. uv verifies every wheel against the package
# index hashes; dependency resolution is frozen at the catalog date.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

af_uv_tool_install --python 3.14.7 zensical==$AF_VERSION

af_finish
