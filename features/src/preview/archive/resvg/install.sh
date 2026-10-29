#!/bin/sh
# preview/archive/resvg: resvg from the producer's release artifact
# (https://github.com/linebender/resvg#install),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name x86_64 aarch64)
case "$AF_VERSION/$AF_ARCH" in
  0.49.0/amd64) sha256=216637e75e909e53e3098b0875cd5f7da2fdc8b4359771b4a29b539afff9d3f5 ;;
  0.49.0/arm64) sha256=00356cb6984836d7638fe10597a01ec4776b0c356e164538f165eca067f0d982 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/linebender/resvg/releases/download/v$AF_VERSION/resvg-linux-$arch.tar.gz"

af_fetch "$url" "$AF_TMP/artifact.tar.gz" "$sha256"
af_extract "$AF_TMP/artifact.tar.gz" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/resvg" resvg

af_finish
