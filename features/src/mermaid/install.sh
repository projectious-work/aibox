#!/bin/sh
# mermaid (atomic): mermaid-cli and Puppeteer as global npm packages
# (https://github.com/mermaid-js/mermaid-cli), pinned by sha512 integrity,
# with the chrome-headless-shell build Puppeteer pins, installed with its
# system dependencies into PUPPETEER_CACHE_DIR. Chrome for Testing publishes
# no checksums for that build.
set -eu
feature_dir=$(cd "${0%/*}" && pwd)
. "$feature_dir/lib/catalog.sh"
. "$feature_dir/lib/aibox-feature.sh"
af_init

af_npm_fetch puppeteer 25.11.0 sha512-n+psxZFH/qlJ3gwTZMxTLUZoE98FAOEjT/BSxHXvPmNCGAfN/kq0UnIZ8m9RB4rAVvJr7zRHUnKKg98GLMUyyQ==
af_npm_fetch @mermaid-js/mermaid-cli 11.16.0 sha512-0InK2nbVIMtzVzCugmdvPkAuvS6wRUqU6Utntff1n8c7lgfRZAdhKY6PSKvcIK9nFmuOUzAgB5+x/XWcroZ7Zg==
# mermaid-cli resolves the global puppeteer instead of a second peer copy.
af_npm_install --legacy-peer-deps
cache="$AF_ROOT/opt/aibox/puppeteer"
af_apt_prepare || af_fail PACKAGE_UNAVAILABLE "packages=chrome-headless-shell-dependencies reason=apt-get-update-failed"
PUPPETEER_CACHE_DIR="$cache" "$AF_ROOT/usr/local/bin/puppeteer" browsers install chrome-headless-shell --install-deps \
  || af_fail DOWNLOAD_FAILED "source=chrome-for-testing browser=chrome-headless-shell"
chmod -R a+rX "$cache"

af_finish
