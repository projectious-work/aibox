#!/bin/sh
# security/scan/syft: syft from the producer's release artifact
# (https://github.com/anchore/syft#installation),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name amd64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  1.50.0/amd64) sha256=bf7b29ff57f06da30918266a0e1c2885a8f99784798d1bdb1628886aa015d788 ;;
  1.50.0/arm64) sha256=887c57cbcc2d0e8c5c110a4571a3fc7150058b24d74f993ee4663516e5c8ce86 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/anchore/syft/releases/download/v$AF_VERSION/syft_${AF_VERSION}_linux_$arch.tar.gz"

af_fetch "$url" "$AF_TMP/artifact.tar.gz" "$sha256"
af_extract "$AF_TMP/artifact.tar.gz" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/syft" syft

af_finish
