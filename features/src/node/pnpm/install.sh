#!/bin/sh
# node/pnpm: the npm package pnpm installed globally with the node/node Feature's npm
# (https://pnpm.io/installation), pinned by the registry's sha512 integrity value.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

case "$AF_VERSION" in
  11.22.0) integrity=sha512-H/hwxMYTPf2I+yr8Rt0T1H8JyXlLQ4xv20fKmMrzvBY4HuC+k6CRuOOCTPAfiJ9G19niCRD7C+GrD7W6qA3WIQ== ;;
  *) af_die "no pinned integrity for $AF_VERSION" ;;
esac
af_npm_fetch pnpm "$AF_VERSION" "$integrity"
af_npm_install

af_finish
