#!/bin/sh
# kubernetes/kustomize: kustomize from the producer's release artifact
# (https://kubectl.docs.kubernetes.io/installation/kustomize/binaries/),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name amd64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  5.8.1/amd64) sha256=029a7f0f4e1932c52a0476cf02a0fd855c0bb85694b82c338fc648dcb53a819d ;;
  5.8.1/arm64) sha256=0953ea3e476f66d6ddfcd911d750f5167b9365aa9491b2326398e289fef2c142 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv$AF_VERSION/kustomize_v${AF_VERSION}_linux_$arch.tar.gz"

af_fetch "$url" "$AF_TMP/artifact.tar.gz" "$sha256"
af_extract "$AF_TMP/artifact.tar.gz" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/kustomize" kustomize

af_finish
