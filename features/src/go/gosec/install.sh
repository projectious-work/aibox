#!/bin/sh
# go/gosec: gosec from the producer's release artifact
# (https://github.com/securego/gosec#install),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name amd64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  2.28.0/amd64) sha256=d7882e505b1ff345d458bf0e893eec8019bc849f861ad73a212869540dd505ff ;;
  2.28.0/arm64) sha256=63259681b6e4b9e7a24d4e187b485e75d3844d28d512b0c97dc831e51d374720 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/securego/gosec/releases/download/v$AF_VERSION/gosec_${AF_VERSION}_linux_$arch.tar.gz"

af_fetch "$url" "$AF_TMP/artifact.tar.gz" "$sha256"
af_extract "$AF_TMP/artifact.tar.gz" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/gosec" gosec

af_finish
