#!/bin/sh
# preview/archive/poppler: Debian package(s) poppler-utils at the candidate version of the
# base's configured repositories (method debian-package).
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

af_debian_install poppler-utils

af_finish
