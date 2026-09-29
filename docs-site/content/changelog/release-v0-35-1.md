---
title: "v0.35.1 — Yazi Markdown preview and refreshed tools"
description: "Restores rendered Markdown previews and refreshes v0 base-image and addon tool versions."
date: 2026-09-29
author: "projectious.work"
tags: [release]
---

The v0.35.1 patch restores rendered Markdown in Yazi when a project also has
an application Python that shadows Debian's system Python. The preview plugin
now uses the interpreter that receives the opt-in `preview-enhanced` addon's
Rich package, and it discards stale plain-text preview caches.

The full-pane preview also respects custom pagers such as `cat` without passing
them `less`-only flags. The release refreshes base-image and addon tool defaults,
including Yazi, Go, Rust, Hugo, Playwright, pnpm, Kubernetes tooling, and Tau,
and updates the CLI lockfile to clear the rustls audit advisory.

After upgrading, run `aibox apply` and rebuild the container. Projects using
pnpm 11, Mermaid CLI 11, or kubectl 1.36 should review the new defaults and
retain explicit older pins where needed.

[Full v0.35.1 release notes](https://github.com/projectious-work/aibox/releases/tag/v0.35.1)
