#!/bin/sh
# node/bun: bun from the producer's release artifact
# (https://bun.sh/docs/installation),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name x64 aarch64)
case "$AF_VERSION/$AF_ARCH" in
  1.4.0/amd64) sha256=2d03fb5fb83ac8b567aca0a281b2ce1a1a19d488f56c2968d88c3f25e92fe452 ;;
  1.4.0/arm64) sha256=4b1a332ee861983eb93bcfe6f770fff94e3e31b2c388bdaea3c8ed35e58eed0e ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/oven-sh/bun/releases/download/bun-v$AF_VERSION/bun-linux-$arch.zip"

af_fetch "$url" "$AF_TMP/artifact.zip" "$sha256"
af_extract "$AF_TMP/artifact.zip" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/bun-linux-$arch/bun" bun

af_finish
