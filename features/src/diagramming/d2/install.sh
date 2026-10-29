#!/bin/sh
# diagramming/d2: d2 from the producer's release artifact
# (https://d2lang.com/tour/install/),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name amd64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  0.7.1/amd64) sha256=eb172adf59f38d1e5a70ab177591356754ffaf9bebb84e0ca8b767dfb421dad7 ;;
  0.7.1/arm64) sha256=ce3a0b985a8f91335a826c254b3a88736fd81afcdd08b58f6c749d2add6864b0 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/d2lang/d2/releases/download/v$AF_VERSION/d2-v$AF_VERSION-linux-$arch.tar.gz"

af_fetch "$url" "$AF_TMP/artifact.tar.gz" "$sha256"
af_extract "$AF_TMP/artifact.tar.gz" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/d2-v$AF_VERSION/bin/d2" d2

af_finish
