#!/bin/sh
# security/scan/cosign: cosign from the producer's release artifact
# (https://docs.sigstore.dev/cosign/system_config/installation/),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name amd64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  3.1.2/amd64) sha256=f7622ed3cf22e55e1ae6377c080979ff77a22da9981c11df222a2e444991e7cf ;;
  3.1.2/arm64) sha256=90e7ae0b5dfd60f20816b52c012addf7fc055ebcc7bea4ce81c428ca8518c302 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/sigstore/cosign/releases/download/v$AF_VERSION/cosign-linux-$arch"

af_fetch "$url" "$AF_TMP/artifact" "$sha256"
af_install_bin "$AF_TMP/artifact" cosign

af_finish
