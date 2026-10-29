#!/bin/sh
# security/scan/gitleaks: gitleaks from the producer's release artifact
# (https://github.com/gitleaks/gitleaks#installing),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name x64 arm64)
case "$AF_VERSION/$AF_ARCH" in
  8.30.1/amd64) sha256=551f6fc83ea457d62a0d98237cbad105af8d557003051f41f3e7ca7b3f2470eb ;;
  8.30.1/arm64) sha256=e4a487ee7ccd7d3a7f7ec08657610aa3606637dab924210b3aee62570fb4b080 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/gitleaks/gitleaks/releases/download/v$AF_VERSION/gitleaks_${AF_VERSION}_linux_$arch.tar.gz"

af_fetch "$url" "$AF_TMP/artifact.tar.gz" "$sha256"
af_extract "$AF_TMP/artifact.tar.gz" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/gitleaks" gitleaks

af_finish
