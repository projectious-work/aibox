# v0.35.2 — release-host Yazi probe fix

> Checks Yazi startup in a generated tmux pane so macOS release validation can complete.


The v0.35.2 patch fixes the macOS release-host smoke test. Its previous
non-interactive `yazi --debug` invocation failed with `Not a tty`, even after
the container image built and the generated tmux layout started successfully.
The probe now checks that Yazi runs in the generated work pane under a PTY.
The release visual checker also handles color escapes split between recording
events, so valid terminal themes no longer fail its status-row invariant.

No project configuration change is required. Use v0.35.2 when completing the
host release phase for the v0.35.1 tool and Markdown-preview updates.

[Full v0.35.2 release notes](https://github.com/projectious-work/aibox/releases/tag/v0.35.2)


---
Source: https://projectious-work.github.io/aibox/v0.x/v0.35.3/changelog/release-v0-35-2/index.md
