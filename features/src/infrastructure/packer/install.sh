#!/bin/sh
# infrastructure/packer: packer from the producer's release artifact
# (https://developer.hashicorp.com/packer/install),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name amd64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  1.16.0/amd64) sha256=5edcd14ab59b535040c512dbecd6ec9ef976a000b073c19d93e4c431c948581e ;;
  1.16.0/arm64) sha256=cf18f03460d92265d49b56befff333e80641d845822799eab04357c39f75b5d7 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://releases.hashicorp.com/packer/$AF_VERSION/packer_${AF_VERSION}_linux_$arch.zip"

af_fetch "$url" "$AF_TMP/artifact.zip" "$sha256"
af_extract "$AF_TMP/artifact.zip" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/packer" packer

af_finish
