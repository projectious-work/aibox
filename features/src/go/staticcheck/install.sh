#!/bin/sh
# go/staticcheck: built with the installed go/go toolchain by go install
# (https://pkg.go.dev/honnef.co/go/tools/cmd/staticcheck); the Go checksum database verifies modules.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

af_go_install honnef.co/go/tools/cmd/staticcheck@v$AF_VERSION

af_finish
