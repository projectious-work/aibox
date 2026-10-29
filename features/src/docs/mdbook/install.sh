#!/bin/sh
# docs/mdbook: mdbook from the producer's release artifact
# (https://rust-lang.github.io/mdBook/guide/installation.html),
# pinned by SHA256 per version and architecture, recorded from the
# producer's published digests.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

arch=$(af_arch_name x86_64 aarch64)
case "$AF_VERSION/$AF_ARCH" in
  0.5.4/amd64) sha256=5222beabd3e37dc5be0d18ff99b79058469354db5c220153a1b92db5ba12be89 ;;
  0.5.4/arm64) sha256=753e5c5c363ee8a56972344dcf91466f005a51db84a7aeffe427ae3ef83d6d44 ;;
  *) af_die "no pinned digest for $AF_VERSION/$AF_ARCH" ;;
esac
url="https://github.com/rust-lang/mdBook/releases/download/v$AF_VERSION/mdbook-v$AF_VERSION-$arch-unknown-linux-musl.tar.gz"

af_fetch "$url" "$AF_TMP/artifact.tar.gz" "$sha256"
af_extract "$AF_TMP/artifact.tar.gz" "$AF_TMP/x"
af_install_bin "$AF_TMP/x/mdbook" mdbook

af_finish
