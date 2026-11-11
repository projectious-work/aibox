#!/bin/sh
# security/crypto/sops: sops from the producer's release artifact
# (https://github.com/getsops/sops#download),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name amd64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  3.13.3/amd64) sha256=e5bec3346a873ae91d871550f3e698c1aad962aff462a080e40f25fde17fef6b ;;
  3.13.3/arm64) sha256=53b0abacd38ef1b12a66d6c100956691b9cefce018d91f81e73ddf7438b94d77 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/getsops/sops/releases/download/v$AF_VERSION/sops-v$AF_VERSION.linux.$arch"

af_fetch "$url" "$AF_TMP/artifact" "$sha256"
af_install_bin "$AF_TMP/artifact" sops

af_finish
