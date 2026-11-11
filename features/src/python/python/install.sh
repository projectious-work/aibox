#!/bin/sh
# python/python: a uv-managed CPython (python-build-standalone) installed
# with the python/uv Feature (https://docs.astral.sh/uv/guides/install-python/).
# uv verifies each download against the SHA256 in its embedded metadata.
# pip is part of the managed distribution.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

uv="$AF_ROOT/usr/local/bin/uv"
[ -x "$uv" ] || af_die "uv is missing although python/uv is a dependsOn dependency"
install_dir="$AF_ROOT/opt/aibox/python"

export UV_PYTHON_INSTALL_DIR="$install_dir"
export UV_CACHE_DIR="$AF_TMP/uv-cache"
export UV_NO_CONFIG=1
"$uv" python install --no-bin "$AF_VERSION" \
  || af_fail DOWNLOAD_FAILED "source=python-build-standalone version=$AF_VERSION"
python=$("$uv" python find --managed-python --no-project "$AF_VERSION") \
  || af_die "uv cannot find the installed Python $AF_VERSION"
"$python" --version | grep -qx "Python $AF_VERSION" || af_die "$python is not Python $AF_VERSION"

install -d "$AF_ROOT/usr/local/bin"
ln -sfn "$python" "$AF_ROOT/usr/local/bin/python3"
ln -sfn "$python" "$AF_ROOT/usr/local/bin/python"

af_finish
