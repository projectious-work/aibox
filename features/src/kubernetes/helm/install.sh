#!/bin/sh
# kubernetes/helm: helm from the producer's release artifact
# (https://helm.sh/docs/intro/install/),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name amd64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  4.2.4/amd64) sha256=c306b46f719b0a4da32d0f78ee21bf90ce8d602f15b22ab753f0674d1670a7f3 ;;
  4.2.4/arm64) sha256=564de2191b881e9f71b5606b25345821ea1682f06ab90499d3ab22b530176da1 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://get.helm.sh/helm-v$AF_VERSION-linux-$arch.tar.gz"

af_fetch "$url" "$AF_TMP/artifact.tar.gz" "$sha256"
af_extract "$AF_TMP/artifact.tar.gz" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/linux-$arch/helm" helm

af_finish
