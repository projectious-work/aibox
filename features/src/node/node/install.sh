#!/bin/sh
# node/node: Node.js with npm and npx from the producer's Linux binary
# archive (https://nodejs.org/en/download), pinned by the SHA256 values in
# https://nodejs.org/dist/v<version>/SHASUMS256.txt.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init
af_apt_require xz-utils

arch=$(af_arch_name x64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  22.23.4/amd64) sha256=1805bcd35e6815559bd2eb97db438d3fdb0eb3a5e6a977b7c1badca76337e9ce ;;
  22.23.4/arm64) sha256=6287aa46b6c375d816e63ae66d024fbbfd70404825aabdfdf898a970871d2f4b ;;
  24.22.2/amd64) sha256=48cc0e21e048790e54c0cf5749c2725005313d652651c40b7a21feb9697aaf75 ;;
  24.22.2/arm64) sha256=eb187837a2230d4047161ac0be008584e16b46f726e9c9490ee35dc6b17d9ba8 ;;
  26.13.0/amd64) sha256=d28a96af6c38596f6d9d97c0ee531c13d238d757742387d1928de822d959e997 ;;
  26.13.0/arm64) sha256=406fe3fbde1c054adff8b297c16fdb5af6659db327c1ac8acd46b56ef1af9a92 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
name="node-v$AF_VERSION-linux-$arch"

af_fetch "https://nodejs.org/dist/v$AF_VERSION/$name.tar.xz" "$AF_TMP/$name.tar.xz" "$sha256"
af_helper_install xz-utils
tar -xJf "$AF_TMP/$name.tar.xz" -C "$AF_TMP" --no-same-owner
install -d "$AF_ROOT/usr/local" "$AF_ROOT/usr/local/etc" "$AF_ROOT/usr/local/share/licenses/node"
cp -R "$AF_TMP/$name/bin" "$AF_TMP/$name/include" "$AF_TMP/$name/lib" "$AF_TMP/$name/share" "$AF_ROOT/usr/local/"
install -m 0644 "$AF_TMP/$name/LICENSE" "$AF_ROOT/usr/local/share/licenses/node/LICENSE"
# npm must not check for or announce its own updates.
printf 'update-notifier=false\nfund=false\n' > "$AF_ROOT/usr/local/etc/npmrc"

af_finish
