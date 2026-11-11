#!/bin/sh
# security/scan/grype: grype from the producer's release artifact
# (https://github.com/anchore/grype#installation),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name amd64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  0.116.1/amd64) sha256=0122df7b655981abe547ad3d2190d65551dac6a2bfc80b4dc2a989b5d0587458 ;;
  0.116.1/arm64) sha256=a8d7504a149629324eb5f4ce3dc25dfd211bbfe047e64ee2bf7844b466c3d84d ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/anchore/grype/releases/download/v$AF_VERSION/grype_${AF_VERSION}_linux_$arch.tar.gz"

af_fetch "$url" "$AF_TMP/artifact.tar.gz" "$sha256"
af_extract "$AF_TMP/artifact.tar.gz" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/grype" grype

af_finish
