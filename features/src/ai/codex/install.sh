#!/bin/sh
# ai/codex: the npm package @openai/codex installed globally with the node/node Feature's npm
# (https://developers.openai.com/codex/cli), pinned by the registry's sha512 integrity value.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

case "$AF_VERSION" in
  0.165.0) integrity=sha512-aiwtQVp2KBhJld966CuU5bVfQ46VNB4mNTl6giX3WGZdt4euI3tv2RUqJpnupPNA8DXHg2z+RoQ0Zfw87dy8Lw== ;;
  *) af_die "no pinned integrity for $AF_VERSION" ;;
esac
af_npm_fetch @openai/codex "$AF_VERSION" "$integrity"
af_npm_install

af_finish
