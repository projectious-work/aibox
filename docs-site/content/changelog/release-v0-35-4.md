---
title: "v0.35.4 — current Claude Code on every apply"
description: "Resolves Claude Code latest to a concrete version so apply rebuilds with the newest CLI, refreshes curated tool pins, and patches Rust TLS dependencies."
date: 2026-10-24
author: "projectious.work"
tags: [release]
---

v0.35.4 fixes `aibox apply` leaving an old Claude Code CLI in the container.
With `version = "latest"`, the generated Dockerfile previously contained the
same unpinned `apt-get install claude-code` line on every apply, so Docker
reused the cached layer. aibox now resolves `latest` against Anthropic's signed
apt repository and pins the concrete package version (for example
`claude-code=2.1.310-1`), recording it in `aibox.lock`. A new upstream release
changes the Dockerfile and triggers a rebuild of that layer.

`aibox up` now recreates the container when the environment from
`.aibox-local.toml` changed since the container was created, so rotated
credentials such as `GH_TOKEN` take effect after `aibox apply && aibox up`.
Compose reads `env_file` only when it creates a container; aibox now labels the
container with a digest of the generated env file and compares it on `up`.

Curated tool pins are refreshed: fzf 0.74.5, delta 0.20.1, starship 1.27.0,
uv 0.13.5, Hugo 0.168.0, MkDocs Material 9.8.0, Zensical 0.0.72, Go 1.27.3,
Rust 1.99.1, pnpm 12.15.0, Bun 1.4.4, PDM 2.29.3, OpenTofu 1.13.2,
Ansible 14.5.0, Packer 1.16.2, Helm 4.3.1, Kustomize 5.8.3, and Tau 0.4.8.
Previous pins remain selectable.

Rust dependencies are updated, including rustls 0.23.46, rustls-webpki
0.103.16, and webpki-roots 1.0.11, which address RUSTSEC-2026-0346 through
-0349, RUSTSEC-2026-0381, and RUSTSEC-2026-0382.

No project configuration change is required. Run `aibox apply` to pick up the
current Claude Code CLI.

[Full v0.35.4 release notes](https://github.com/projectious-work/aibox/releases/tag/v0.35.4)
