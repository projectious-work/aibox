#!/bin/sh
# browser-testing/playwright: the Playwright test runner and CLI as a global
# npm package (https://playwright.dev/docs/intro), pinned by the registry's
# sha512 integrity value. Browsers are separate Features.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

case "$AF_VERSION" in
  1.62.1) integrity=sha512-DTcUc8qii+cpHvtOwggMtBRMjKZHXYWdw8syRYu2vtzuq4Wxphqq4NfCs5Zt44L6mA8rfDfj+PHnxFc/FeK6mQ== ;;
  *) af_die "no pinned integrity for $AF_VERSION" ;;
esac
package="$AF_TMP/playwright-test-$AF_VERSION.tgz"

af_download "https://registry.npmjs.org/@playwright/test/-/test-$AF_VERSION.tgz" "$package"
af_verify_integrity "$package" "$integrity"
# npm verifies the exact playwright and playwright-core dependencies against
# the registry's integrity values.
PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1 npm_config_cache="$AF_TMP/npm-cache" npm_config_logs_dir="$AF_TMP/npm-logs" \
  npm install --global --prefix "$AF_ROOT/usr/local" --no-audit --no-fund --no-update-notifier "$package" \
  || af_fail DOWNLOAD_FAILED "package=@playwright/test@$AF_VERSION"

af_finish
