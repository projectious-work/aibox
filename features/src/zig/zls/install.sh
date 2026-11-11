#!/bin/sh
# zig/zls: zls from the producer's release artifact
# (https://zigtools.org/zls/install/),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init
af_apt_require xz-utils minisign

minisign_key=RWR+9B91GBZ0zOjh6Lr17+zKf5BoSuFvrx2xSeDE57uIYvnKBGmMjOex
arch=$(af_arch_name x86_64 aarch64)
case "$AF_VERSION/$AF_ARCH" in
  0.17.0/amd64) sha256=dfa4514e1868803fdf2ae56832c88b9d400094d70ea6bb136b57bcc91af10258 ;;
  0.17.0/arm64) sha256=0d314f3cab30da770328e1b8c84ab89bb308e56dbc26a85ac82429ff02b44873 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/zigtools/zls/releases/download/$AF_VERSION/zls-$arch-linux.tar.xz"

af_fetch "$url" "$AF_TMP/artifact.tar.xz" "$sha256"
af_download "$url.minisig" "$AF_TMP/artifact.tar.xz.minisig"
af_helper_install xz-utils minisign
af_verify_minisign "$AF_TMP/artifact.tar.xz" "$AF_TMP/artifact.tar.xz.minisig" "$minisign_key"
af_extract "$AF_TMP/artifact.tar.xz" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/zls" zls

af_finish
