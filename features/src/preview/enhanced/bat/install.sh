#!/bin/sh
# preview/enhanced/bat: bat from the producer's release artifact
# (https://github.com/sharkdp/bat#installation),
# pinned by SHA256 per version and architecture, recorded from the
# producer's GitHub asset digests. Installs the bat command.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name x86_64 aarch64)
case "$AF_VERSION/$AF_ARCH" in
  0.26.1/amd64) sha256=0dcd8ac79732c0d5b136f11f4ee00e581440e16a44eab5b3105b611bbf2cf191 ;;
  0.26.1/arm64) sha256=6369242c584065f195fb20cb36fbd7cb63ae690605bbe89868a7596b596c2c23 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/sharkdp/bat/releases/download/v$AF_VERSION/bat-v$AF_VERSION-$arch-unknown-linux-musl.tar.gz"

af_fetch "$url" "$AF_TMP/artifact.tar.gz" "$sha256"
af_extract "$AF_TMP/artifact.tar.gz" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/bat-v$AF_VERSION-$arch-unknown-linux-musl/bat" bat

af_finish
