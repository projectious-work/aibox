# aibox v0.35.1 — 2026-09-29

**Summary:** This patch restores rendered Markdown previews for Yazi users and
refreshes the v0 tool catalog to current upstream releases. Update the v0 CLI,
apply the project configuration, and rebuild the container to receive the new tools.

## Changed

- Refresh base-image tools and curated addon defaults, including Yazi 26.9.1,
  Go 1.27.1, Rust 1.98.1, Hugo 0.167.0, Playwright 1.63.0, pnpm 12.8.1,
  Mermaid CLI 12.0.0, kubectl 1.37.1, and Tau 0.4.6.
- Update the v0 example configuration, generated development Dockerfile,
  documentation, and verified binary digests for affected tools.

## Fixed

- Render Markdown through Debian's system Python, where the opt-in
  `preview-enhanced` addon installs `python3-rich`.
- Invalidate stale plain-text preview caches and handle non-`less` pagers
  without passing them `less`-specific flags.
- Keep the tmux attention test isolated from the active Codex session so its
  fixture result is reproducible.
- Refresh the Rust lockfile, including rustls 0.23.45, to clear the release
  audit advisory.

## Compatibility

- Minimum processkit version remains v0.28.8.
- The Rich renderer remains opt-in through `preview-enhanced`; the fallback
  preview remains available when that addon is not selected.

## Upgrade notes

After upgrading to v0.35.1, run `aibox apply` to refresh the managed Yazi
plugin and full-pane helper. Rebuild the container against the v0.35.1 image
once the host image-publishing phase is complete. Review projects that rely on
pnpm 11, Mermaid CLI 11, or kubectl 1.36 before adopting the new defaults;
explicit version pins remain available. Mermaid CLI 12 requires Node 22.13+
(the default Node 26 satisfies this).

[v0.35.1]: https://github.com/projectious-work/aibox/compare/v0.35.0...v0.35.1
