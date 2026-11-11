#!/bin/sh
# audio/voice: Debian package(s) sox libsox-fmt-pulse pulseaudio-utils libasound2-plugins at the candidate version of the
# base's configured repositories (method debian-package).
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

af_debian_install sox libsox-fmt-pulse pulseaudio-utils libasound2-plugins

af_finish
