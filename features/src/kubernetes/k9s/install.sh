#!/bin/sh
# kubernetes/k9s: k9s from the producer's release artifact
# (https://k9scli.io/topics/install/),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name amd64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  0.51.0/amd64) sha256=c3752ad51a5a4015a113819c4eeb6e55a4d0e4b8e652494797532f6fc8161dd7 ;;
  0.51.0/arm64) sha256=3ee05c82e5f9198928a4e86133608ba6a2c10a2244d6a7789e820f78319d640c ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/derailed/k9s/releases/download/v$AF_VERSION/k9s_Linux_$arch.tar.gz"

af_fetch "$url" "$AF_TMP/artifact.tar.gz" "$sha256"
af_extract "$AF_TMP/artifact.tar.gz" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/k9s" k9s

af_finish
