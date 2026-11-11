#!/bin/sh
# rust/clippy: the clippy component of the installed rust/rustc toolchain through
# rustup (https://rust-lang.github.io/rustup/concepts/components.html); rustup
# verifies it against the release's channel manifest.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

af_rust_env
"$AF_RUSTUP" component add clippy || af_fail DOWNLOAD_FAILED "source=rustup component=clippy"
chmod -R a+w "$RUSTUP_HOME" "$CARGO_HOME"

af_finish
