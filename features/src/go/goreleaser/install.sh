#!/bin/sh
# go/goreleaser: goreleaser from the producer's release artifact
# (https://goreleaser.com/install/),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name x86_64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  2.17.1/amd64) sha256=a99bbc7ae0d8d897b07c4c497a9b62f222558804715ef219d1af05a7e417bc80 ;;
  2.17.1/arm64) sha256=702f03769ac8bcb0e47839c82243cc614ae995633599a98c63062e13ea85f829 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/goreleaser/goreleaser/releases/download/v$AF_VERSION/goreleaser_Linux_$arch.tar.gz"

af_fetch "$url" "$AF_TMP/artifact.tar.gz" "$sha256"
af_extract "$AF_TMP/artifact.tar.gz" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/goreleaser" goreleaser

af_finish
