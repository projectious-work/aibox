#!/bin/sh
# rust/rustc: rustc and cargo through rustup
# (https://rust-lang.github.io/rustup/installation/other.html). rustup-init
# 1.29.2 is pinned by its published SHA256; rustup verifies each toolchain
# component against the SHA256 in the release's channel manifest.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init
af_apt_require gcc libc6-dev

rustup_version=1.29.2
triple=$(af_arch_name x86_64 aarch64)-unknown-linux-gnu
case "$AF_ARCH" in
  amd64) sha256=5be191e8c70a5007074254cb74d78dbfdd25cd5a61c27de7b8e5410122ec24f4 ;;
  arm64) sha256=a9c168dbea81a84671c6b7a5c1795cb0eb699020110d107ddfca427bd797e697 ;;
  *) af_die "no pinned digest for $AF_ARCH" ;;
esac

af_fetch "https://static.rust-lang.org/rustup/archive/$rustup_version/$triple/rustup-init" "$AF_TMP/rustup-init" "$sha256"
chmod 0755 "$AF_TMP/rustup-init"
# rustc links through the system C toolchain.
af_apt_install gcc libc6-dev

# The same locations as the Feature's containerEnv.
export RUSTUP_HOME="$AF_ROOT/usr/local/rustup"
export CARGO_HOME="$AF_ROOT/usr/local/cargo"
"$AF_TMP/rustup-init" -y --no-modify-path --profile minimal \
  --default-host "$triple" --default-toolchain "$AF_VERSION" \
  || af_fail DOWNLOAD_FAILED "url=https://static.rust-lang.org/dist/channel-rust-$AF_VERSION.toml"
"$CARGO_HOME/bin/rustup" set auto-self-update disable
# As in the producer's container image, the remote user may add crates,
# components and targets.
chmod -R a+w "$RUSTUP_HOME" "$CARGO_HOME"

af_finish
