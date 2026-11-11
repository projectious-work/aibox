#!/bin/sh
# node/yarn: the npm package @yarnpkg/cli-dist installed globally with the node/node Feature's npm
# (https://yarnpkg.com/getting-started/install), pinned by the registry's sha512 integrity value.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

case "$AF_VERSION" in
  4.17.0) integrity=sha512-quDpJzauzaCeutv4S7MObzjKcExJU0iKPzwl2LN66AP2MSVC/yTDuwnTyR3NwaA6euvSMBhk8MzX2GBDvVXxsA== ;;
  *) af_die "no pinned integrity for $AF_VERSION" ;;
esac
af_npm_fetch @yarnpkg/cli-dist "$AF_VERSION" "$integrity"
af_npm_install

af_finish
