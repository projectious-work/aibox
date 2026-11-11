#!/bin/sh
# ai/continue: the npm package @continuedev/cli installed globally with the node/node Feature's npm
# (https://www.npmjs.com/package/@continuedev/cli), pinned by the registry's sha512 integrity value.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

case "$AF_VERSION" in
  1.5.47) integrity=sha512-gtpewV3RoIOD9dyTtKIBi1SY0VOHRu3Ehe7C/mmnswm+j34MPyrcQhQaWj/m+jdfGO4fNIKdrgGIlLso1ULDFw== ;;
  *) af_die "no pinned integrity for $AF_VERSION" ;;
esac
af_npm_fetch @continuedev/cli "$AF_VERSION" "$integrity"
af_npm_install

af_finish
