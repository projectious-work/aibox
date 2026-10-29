#!/bin/sh
# ai/copilot: copilot from the producer's release artifact
# (https://github.com/github/copilot-cli#installation),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name x64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  1.0.105/amd64) sha256=aba1f06f27bc0447656647f76e0a4cf58b4b2a341c75ff04a750982b6aeb6bd2 ;;
  1.0.105/arm64) sha256=e7b488cc04ac3e2ed800e1bf0aace33ada42dcf167529b3d2f48a61e1535a847 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/github/copilot-cli/releases/download/v$AF_VERSION/copilot-linux-$arch.tar.gz"

af_fetch "$url" "$AF_TMP/artifact.tar.gz" "$sha256"
af_extract "$AF_TMP/artifact.tar.gz" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/copilot" copilot

af_finish
