# aibox v0.35.4 — 2026-10-24

**Summary:** This patch makes `aibox apply` install the current Claude Code
CLI when the harness tracks `latest`. It also refreshes curated tool pins and
updates Rust dependencies to address the open RustSec advisories.

## Fixed

- Resolve Claude Code `version = "latest"` to the newest package in
  Anthropic's signed apt repository for the host architecture, and pin
  `claude-code=<version>` in the generated Dockerfile. Previously the
  unpinned install line never changed, so Docker reused a stale cached layer
  and `aibox apply` kept the old CLI. The resolved version is recorded in
  `aibox.lock`. If the lookup fails, the previous unpinned behavior is used.
- Recreate the container on `aibox up` when the environment from
  `.aibox-local.toml` changed since the container was created. Compose reads
  `env_file` only at creation, so a rotated token such as `GH_TOKEN` was
  previously ignored until a manual recreate. The generated compose service
  now carries an `aibox.local-env-hash` label (a digest, never the values),
  and `aibox apply` reports when the next `up` will recreate.
- Update rustls to 0.23.46, rustls-webpki to 0.103.16, and webpki-roots to
  1.0.11 (RUSTSEC-2026-0346, -0347, -0348, -0349, -0381, -0382), along with
  other compatible crate updates.

## Changed

- Base image: fzf 0.74.5, delta 0.20.1, starship 1.27.0, uv 0.13.5,
  Rust 1.99.1.
- Addon defaults: Hugo 0.168.0, MkDocs Material 9.8.0, Zensical 0.0.72,
  Go 1.27.3, Rust 1.99.1, uv 0.13.5, pnpm 12.15.0, Bun 1.4.4, PDM 2.29.3,
  OpenTofu 1.13.2, Ansible 14.5.0, Packer 1.16.2, Helm 4.3.1,
  Kustomize 5.8.3, Tau 0.4.8. Previous pins remain selectable.

## Compatibility

- Minimum processkit version remains v0.28.8.
- No project configuration change is required.
- The v1 development branch remains independent.

## Upgrade notes

Run `aibox apply` after upgrading to rebuild with the current Claude Code CLI.
Projects that pin an explicit Claude Code version, or that select the `stable`
channel, keep their existing behavior.

[v0.35.4]: https://github.com/projectious-work/aibox/compare/v0.35.3...v0.35.4
