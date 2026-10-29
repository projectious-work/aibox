#!/bin/sh
# release/hadolint: hadolint from the producer's release artifact
# (https://github.com/hadolint/hadolint#install),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name x86_64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  2.15.1/amd64) sha256=c7187db94eeeeca956519a6af171adc31453941a1e777961f6e680f697c8c507 ;;
  2.15.1/arm64) sha256=f6198ef8090f404dbb771abfee086eb8c48ac177f30da7fd3510aca35b344b5d ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/hadolint/hadolint/releases/download/v$AF_VERSION/hadolint-linux-$arch"

af_fetch "$url" "$AF_TMP/artifact" "$sha256"
af_install_bin "$AF_TMP/artifact" hadolint

af_finish
