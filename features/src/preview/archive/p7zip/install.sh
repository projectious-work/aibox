#!/bin/sh
# preview/archive/p7zip: Debian package(s) 7zip at the candidate version of the
# base's configured repositories (method debian-package).
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

af_debian_install 7zip

af_finish
