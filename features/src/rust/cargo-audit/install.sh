#!/bin/sh
# rust/cargo-audit: built from crates.io with the installed rust/rustc
# toolchain (https://github.com/rustsec/rustsec/tree/main/cargo-audit); the
# crate's Cargo.lock (--locked) pins and checksums every dependency.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

af_rust_env
CARGO_TARGET_DIR="$AF_TMP/target" "$CARGO_HOME/bin/cargo" install --locked cargo-audit --version "$AF_VERSION" \
  || af_fail DOWNLOAD_FAILED "source=crates.io crate=cargo-audit@$AF_VERSION"
rm -rf "$CARGO_HOME/registry" "$CARGO_HOME/git"
chmod -R a+w "$RUSTUP_HOME" "$CARGO_HOME"

af_finish
