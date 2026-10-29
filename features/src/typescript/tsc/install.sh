#!/bin/sh
# typescript/tsc: the npm package typescript installed globally with the node/node Feature's npm
# (https://www.typescriptlang.org/download/), pinned by the registry's sha512 integrity value.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

case "$AF_VERSION" in
  7.0.2) integrity=sha512-8FYau96o3NKOhbjKi/qNvG/W5jhzxkbdm5sj9AbZ/5T5sWqn3hJgLfGx27sRKZWTvyzCP8dLRBTf5tBTSRVUNA== ;;
  *) af_die "no pinned integrity for $AF_VERSION" ;;
esac
af_npm_fetch typescript "$AF_VERSION" "$integrity"
af_npm_install

af_finish
