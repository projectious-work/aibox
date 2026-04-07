# aibox v0.14.4

CLI patch release. Three usability tweaks based on real-world use of v0.14.3:

1. **Yazi default ratio** changed from `[1, 3, 4]` to `[3, 5, 18]` — a much
   more preview-heavy split that makes the file preview the dominant pane.
2. **`ai` layout proportions** flipped from yazi 40% / AI 60% to **yazi 60%
   / AI 40%** — yazi gets the bigger half because file navigation is the
   primary surface in the AI-first layout.
3. **New `cowork-swap` layout** — a re-arrangement of `cowork` for users
   who prefer the editor on the right.

## 1. Yazi default ratio: `[3, 5, 18]`

The new default gives:

- 3/26 width to the parent column (≈ 12%)
- 5/26 width to the current directory column (≈ 19%)
- **18/26 width to the preview column** (≈ 69%)

This is a sharper, preview-heavy split. To get the previous behavior back,
edit `.aibox-home/.config/yazi/yazi.toml`:

```toml
[mgr]
ratio = [1, 3, 4]   # previous default
```

`aibox sync` will not overwrite a hand-edited yazi.toml — only newly
seeded files (and the migration from `[manager]` → `[mgr]`) touch it.

## 2. `ai` layout: yazi 60% / AI 40%

The previous v0.14.2/3 split was yazi 40% / AI 60%. The new split
emphasizes file navigation as the primary surface — the AI agent column
is still wide enough to read responses, and yazi's preview column gets
more room.

```
┌────────────────────────────────┬──────────────────────┐
│                                │                      │
│  yazi (60%)                    │  claude (40%)        │
│  parent | current | preview    │                      │
│  (with the new [3,5,18] ratio  │                      │
│   yazi's preview column is     │                      │
│   the visual centerpiece)      │                      │
│                                │                      │
├────────────────────────────────┴──────────────────────┤
│  status-bar                                           │
└───────────────────────────────────────────────────────┘
```

If you prefer the old proportions, you can override after `aibox sync` by
hand-editing `.aibox-home/.config/zellij/layouts/ai.kdl` — but that file
is regenerated on every sync, so the cleaner path is to wait for a future
release that adds per-layout pane-size config.

## 3. New layout: `cowork-swap`

A re-arrangement of `cowork` for users who prefer the editor on the right
(the bigger pane). The outer split is **40/60 instead of cowork's 50/50**,
and the editor and AI panes swap roles.

```
┌──────────────────┬────────────────────────────────────────┐
│  yazi   (40%)    │                                        │
│                  │                                        │
├──────────────────┤  vim editor (60%)                      │
│  AI agent (60%)  │                                        │
│                  │                                        │
├──────────────────┴────────────────────────────────────────┤
│  status-bar                                               │
└───────────────────────────────────────────────────────────┘

  Tab 1: cowork-swap   Tab 2: git   Tab 3: shell
```

When multiple AI providers are configured, they stack in the bottom-left
pane (same convention as `cowork`, `browse`, and `ai`).

When no AI provider is configured, the layout degenerates to a simple
yazi-left + vim-right shape (same as `dev`, with the `cowork-swap` tab
name preserved).

Select with:

```bash
aibox start --layout cowork-swap
```

or set as default in `aibox.toml`:

```toml
[customization]
layout = "cowork-swap"
```

## Tests

314 tests pass. Five new tests for `cowork-swap`, plus the `ai_layout_*`
tests updated for the new 60/40 proportions:

- `cowork_swap_layout_single_provider`
- `cowork_swap_layout_multiple_providers_stacked`
- `cowork_swap_layout_no_providers`
- `cowork_swap_layout_editor_right_of_files`
- `cowork_swap_layout_ai_below_files_in_left_stack`

## Container Images

- `ghcr.io/projectious-work/aibox:base-debian-v0.14.4`
- `ghcr.io/projectious-work/aibox:base-debian-latest`

The base image is rebuilt with the new yazi `[mgr]` ratio default
(`[3, 5, 18]`). All other content is identical to v0.14.3.

## CLI Binaries

- `aibox-v0.14.4-aarch64-unknown-linux-gnu.tar.gz`
- `aibox-v0.14.4-x86_64-unknown-linux-gnu.tar.gz`
- `aibox-v0.14.4-aarch64-apple-darwin.tar.gz` (added in Phase 2)
- `aibox-v0.14.4-x86_64-apple-darwin.tar.gz` (added in Phase 2)

## Upgrading

```bash
# Install the new CLI
curl -fsSL https://raw.githubusercontent.com/projectious-work/aibox/main/scripts/install.sh | bash

# Or on macOS, follow Phase 2 instructions on the host

# In each project, run sync — the cowork-swap.kdl layout will be added,
# and ai.kdl will be regenerated with the new 60/40 proportions:
cd <project>
aibox sync
```

The yazi.toml ratio is left alone if you've already hand-edited it; new
projects (created via `aibox init` after upgrading) will get the new
`[3, 5, 18]` default.
