#!/bin/sh
# release/shellcheck: shellcheck from the producer's release artifact
# (https://github.com/koalaman/shellcheck#installing),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name x86_64 aarch64)
case "$AF_VERSION/$AF_ARCH" in
  0.11.0/amd64) sha256=b7af85e41cc99489dcc21d66c6d5f3685138f06d34651e6d34b42ec6d54fe6f6 ;;
  0.11.0/arm64) sha256=68a8133197a50beb8803f8d42f9908d1af1c5540d4bb05fdfca8c1fa47decefc ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/koalaman/shellcheck/releases/download/v$AF_VERSION/shellcheck-v$AF_VERSION.linux.$arch.tar.gz"

af_fetch "$url" "$AF_TMP/artifact.tar.gz" "$sha256"
af_extract "$AF_TMP/artifact.tar.gz" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/shellcheck-v$AF_VERSION/shellcheck" shellcheck

af_finish
