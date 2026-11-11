#!/bin/sh
# security/crypto/minisign: minisign from the producer's release artifact
# (https://jedisct1.github.io/minisign/),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name x86_64 aarch64)
case "$AF_VERSION/$AF_ARCH" in
  0.12/amd64) sha256=9a599b48ba6eb7b1e80f12f36b94ceca7c00b7a5173c95c3efc88d9822957e73 ;;
  0.12/arm64) sha256=9a599b48ba6eb7b1e80f12f36b94ceca7c00b7a5173c95c3efc88d9822957e73 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/jedisct1/minisign/releases/download/$AF_VERSION/minisign-$AF_VERSION-linux.tar.gz"

af_fetch "$url" "$AF_TMP/artifact.tar.gz" "$sha256"
af_extract "$AF_TMP/artifact.tar.gz" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/minisign-linux/$arch/minisign" minisign

af_finish
