#!/bin/sh
# preview/enhanced/glow: glow from the producer's release artifact
# (https://github.com/charmbracelet/glow#installation),
# pinned by SHA256 per version and architecture, recorded from the
# producer's checksums.txt and GitHub asset digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name x86_64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  3.0.0/amd64) sha256=13e05e4b2acc18d2aee44291aefe6325b077ec321b631a0cfa780e8e3bc33f78 ;;
  3.0.0/arm64) sha256=810c39f4691feb75e675a5f5f54b9fa091354e7599d9cf9c9fc9b97e47a99759 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/charmbracelet/glow/releases/download/v$AF_VERSION/glow_${AF_VERSION}_Linux_$arch.tar.gz"

af_fetch "$url" "$AF_TMP/artifact.tar.gz" "$sha256"
af_extract "$AF_TMP/artifact.tar.gz" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/glow_${AF_VERSION}_Linux_$arch/glow" glow

af_finish
