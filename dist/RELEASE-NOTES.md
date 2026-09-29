# aibox v0.35.1 — 2026-09-29

**Summary:** This patch restores rendered Markdown previews for Yazi users when
an application Python shadows Debian's Rich-enabled interpreter. Update the
v0 CLI and apply the project configuration to refresh managed preview files.

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
once the host image-publishing phase is complete.

[v0.35.1]: https://github.com/projectious-work/aibox/compare/v0.35.0...v0.35.1
