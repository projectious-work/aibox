#!/bin/sh
# typst: typst from the producer's release artifact
# (https://github.com/typst/typst#installation),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init
af_apt_require xz-utils

arch=$(af_arch_name x86_64 aarch64)
case "$AF_VERSION/$AF_ARCH" in
  0.15.0/amd64) sha256=59b207df01be2dab9f13e80f73d04d7ff8273ffd46b3dd1b9eef5c60f3eeabea ;;
  0.15.0/arm64) sha256=cdf50ffc7b8ba759ed02200632eda3d78eb8b99aacb6611f4f75684990647620 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/typst/typst/releases/download/v$AF_VERSION/typst-$arch-unknown-linux-musl.tar.xz"

af_fetch "$url" "$AF_TMP/artifact.tar.xz" "$sha256"
af_helper_install xz-utils
af_extract "$AF_TMP/artifact.tar.xz" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/typst-$arch-unknown-linux-musl/typst" typst

af_finish
