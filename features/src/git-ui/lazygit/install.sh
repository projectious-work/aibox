#!/bin/sh
# git-ui/lazygit: lazygit from the producer's release artifact
# (https://github.com/jesseduffield/lazygit#installation),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name x86_64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  0.66.2/amd64) sha256=158efc9ddfa1241fab5fb3526bb442c8dc6c4b1904d73213b9e4e93e03c52034 ;;
  0.66.2/arm64) sha256=79a3161eeb642f0473e23db6341f7baabd4a5933e921a9c1167c5d5c765c41a4 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/jesseduffield/lazygit/releases/download/v$AF_VERSION/lazygit_${AF_VERSION}_linux_$arch.tar.gz"

af_fetch "$url" "$AF_TMP/artifact.tar.gz" "$sha256"
af_extract "$AF_TMP/artifact.tar.gz" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/lazygit" lazygit

af_finish
