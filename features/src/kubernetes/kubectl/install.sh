#!/bin/sh
# kubernetes/kubectl: kubectl from the producer's release artifact
# (https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name amd64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  1.36.4/amd64) sha256=8b8f088da2dab964f853b38464033b1be15ede2839eca751482357c45abdd05a ;;
  1.36.4/arm64) sha256=0ecf44450ee6063bf19dd166a103ee6df4a9034455c2abce626e6eea657d73fb ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://dl.k8s.io/release/v$AF_VERSION/bin/linux/$arch/kubectl"

af_fetch "$url" "$AF_TMP/artifact" "$sha256"
af_install_bin "$AF_TMP/artifact" kubectl

af_finish
