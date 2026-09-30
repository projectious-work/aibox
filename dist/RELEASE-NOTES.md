# aibox v0.35.2 — 2026-09-29

**Summary:** This patch corrects the macOS release-host Yazi smoke probe so
the v0.35.1 tool and Markdown-preview updates can complete host validation and
publication. No project configuration change is required.

## Fixed

- Verify Yazi startup in the generated tmux work pane, where a PTY is present,
  instead of running interactive `yazi --debug` through a non-TTY container
  command and failing with `Not a tty`.
- Guard the release-host script against reintroducing a non-TTY Yazi debug
  invocation.
- Replay ANSI color escapes that span multiple screencast events before
  checking terminal theme invariants, avoiding false status-row failures.

## Compatibility

- Minimum processkit version remains v0.28.8.
- This patch does not change aibox CLI behavior or project configuration.

## Upgrade notes

Use v0.35.2 for the macOS host release phase. The v0.35.1 host gate failed
before macOS binaries or images were published; do not reuse its failed run.

[v0.35.2]: https://github.com/projectious-work/aibox/compare/v0.35.1...v0.35.2
