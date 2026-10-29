#!/bin/sh
# browser-testing/axe-playwright: the npm packages @axe-core/playwright and
# axe-core installed globally next to the browser-testing/playwright Feature
# (https://github.com/dequelabs/axe-core-npm), pinned by sha512 integrity.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

case "$AF_VERSION" in
  4.13.0)
    integrity=sha512-6YLx+kxXu5GJceG4ozFg+33a2EMTdjYwWGloJ3sb9Kta5pp+ZNS53uxGVog5JetIY8s++P5UrtX+cri+u0VAVg==
    core_integrity=sha512-UzGt8zg7Ny8djbYMhxl2zuEevVa7r2gJjYY5Lwr1xM7+XU2nd6CkIWFTVcCIbAP63vSz71NaVyyuSk9lHKcy0A==
    ;;
  *) af_die "no pinned integrity for $AF_VERSION" ;;
esac
af_npm_fetch @axe-core/playwright "$AF_VERSION" "$integrity"
af_npm_fetch axe-core "$AF_VERSION" "$core_integrity"
af_npm_install

af_finish
