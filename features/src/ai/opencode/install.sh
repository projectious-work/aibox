#!/bin/sh
# ai/opencode: opencode from the producer's release artifact
# (https://opencode.ai/docs/),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name x64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  1.19.5/amd64) sha256=588310a336fa3b9df36b3f0225b26f679d453efaad0cf3d66c3651ef3ec1647a ;;
  1.19.5/arm64) sha256=fae5be262d693b7c5e562856cc473e59e741772936db0c0600ac0d5ed5bd088b ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/anomalyco/opencode/releases/download/v$AF_VERSION/opencode-linux-$arch.tar.gz"

af_fetch "$url" "$AF_TMP/artifact.tar.gz" "$sha256"
af_extract "$AF_TMP/artifact.tar.gz" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/opencode" opencode

af_finish
