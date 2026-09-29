---
title: "v0.35.1 — restored Yazi Markdown preview"
description: "Uses the Rich-enabled system Python for Markdown previews and keeps custom pagers working."
date: 2026-09-29
author: "projectious.work"
tags: [release]
---

The v0.35.1 patch restores rendered Markdown in Yazi when a project also has
an application Python that shadows Debian's system Python. The preview plugin
now uses the interpreter that receives the opt-in `preview-enhanced` addon's
Rich package, and it discards stale plain-text preview caches.

The full-pane preview also respects custom pagers such as `cat` without passing
them `less`-only flags. The release also refreshes the CLI dependency lockfile
to clear the rustls audit advisory. After upgrading, run `aibox apply` to
refresh the managed preview files.

[Full v0.35.1 release notes](https://github.com/projectious-work/aibox/releases/tag/v0.35.1)
