#!/bin/sh
# browser-testing/firefox: the firefox build of the installed Playwright version
# with its system dependencies (playwright install --with-deps,
# https://playwright.dev/docs/browsers). Playwright publishes no checksums for
# browser builds (catalog verification: none).
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

playwright="$AF_ROOT/usr/local/bin/playwright"
[ -x "$playwright" ] || af_die "playwright is missing although browser-testing/playwright is a dependsOn dependency"
browsers="$AF_ROOT/opt/aibox/ms-playwright"
# Playwright installs the browser's system dependencies with apt-get.
af_apt_prepare || af_fail PACKAGE_UNAVAILABLE "packages=playwright-dependencies reason=apt-get-update-failed"
PLAYWRIGHT_BROWSERS_PATH="$browsers" "$playwright" install --with-deps --no-shell firefox \
  || af_fail DOWNLOAD_FAILED "source=playwright-cdn browser=firefox"
[ -d "$browsers/firefox-1538" ] || af_die "Playwright did not install firefox revision 1538"
chmod -R a+rX "$browsers"

af_finish
