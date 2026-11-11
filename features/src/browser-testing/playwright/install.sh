#!/bin/sh
# browser-testing/playwright: the npm package @playwright/test installed globally with the node/node Feature's npm
# (https://playwright.dev/docs/intro), pinned by the registry's sha512 integrity value.
# Browsers are separate Features.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

case "$AF_VERSION" in
  1.62.1) integrity=sha512-DTcUc8qii+cpHvtOwggMtBRMjKZHXYWdw8syRYu2vtzuq4Wxphqq4NfCs5Zt44L6mA8rfDfj+PHnxFc/FeK6mQ== ;;
  *) af_die "no pinned integrity for $AF_VERSION" ;;
esac
af_npm_fetch @playwright/test "$AF_VERSION" "$integrity"
af_npm_install

af_finish
