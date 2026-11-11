#!/bin/sh
# typescript/biome: biome from the producer's release artifact
# (https://biomejs.dev/guides/manual-installation/),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name x64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  2.6.4/amd64) sha256=eb9afd5af3c04bd242be6db39fb812a025b42469af8c6a621535036dbd9b0c63 ;;
  2.6.4/arm64) sha256=11a2e6770b3d7c8ce8cbec4652a516584b179dc30febecce0ca852d87074bee8 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/biomejs/biome/releases/download/%40biomejs%2Fbiome%40$AF_VERSION/biome-linux-$arch"

af_fetch "$url" "$AF_TMP/artifact" "$sha256"
af_install_bin "$AF_TMP/artifact" biome

af_finish
