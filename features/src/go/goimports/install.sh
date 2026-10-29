#!/bin/sh
# go/goimports: built with the installed go/go toolchain by go install
# (https://pkg.go.dev/golang.org/x/tools/cmd/goimports); the Go checksum database verifies modules.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

af_go_install golang.org/x/tools/cmd/goimports@v$AF_VERSION

af_finish
