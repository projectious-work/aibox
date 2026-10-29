#!/bin/sh
# zig/zig: the Zig compiler from the producer's release archive
# (https://ziglang.org/download/), pinned by the SHA256 in
# https://ziglang.org/download/index.json and verified with the producer's
# minisign key.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init
af_apt_require xz-utils minisign

minisign_key=RWSGOq2NVecA2UPNdBUZykf1CCb147pkmdtYxgb3Ti+JO/wCYvhbAb/U
arch=$(af_arch_name x86_64 aarch64)
case "$AF_VERSION/$AF_ARCH" in
  0.17.0/amd64) sha256=1cbe9df9f27e6b78d14ccbca43b6703a404ef79ef1c463de901d7f088d4e2026 ;;
  0.17.0/arm64) sha256=9e8d11661d4ae3bd57702a3832781e23ad151dde5798e16a5ccd503f65234ff8 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
name="zig-$arch-linux-$AF_VERSION"
url="https://ziglang.org/download/$AF_VERSION/$name.tar.xz"

af_fetch "$url" "$AF_TMP/$name.tar.xz" "$sha256"
af_download "$url.minisig" "$AF_TMP/$name.tar.xz.minisig"
af_helper_install xz-utils minisign
af_verify_minisign "$AF_TMP/$name.tar.xz" "$AF_TMP/$name.tar.xz.minisig" "$minisign_key"
tar -xJf "$AF_TMP/$name.tar.xz" -C "$AF_TMP" --no-same-owner
install -d "$AF_ROOT/opt/aibox" "$AF_ROOT/usr/local/bin"
rm -rf "$AF_ROOT/opt/aibox/zig"
mv "$AF_TMP/$name" "$AF_ROOT/opt/aibox/zig"
ln -sfn ../../../opt/aibox/zig/zig "$AF_ROOT/usr/local/bin/zig"

af_finish
