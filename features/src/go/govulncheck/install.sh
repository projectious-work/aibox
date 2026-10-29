#!/bin/sh
# go/govulncheck: built with the installed go/go toolchain by go install
# (https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck); the Go checksum database verifies modules.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

af_go_install golang.org/x/vuln/cmd/govulncheck@v$AF_VERSION

af_finish
