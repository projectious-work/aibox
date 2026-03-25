---
sidebar_position: 2
title: "Prompt Presets"
---

# Starship Prompt Presets

aibox includes 6 [Starship](https://starship.rs) prompt presets that work with any theme. Set a preset in `aibox.toml`:

```toml
[appearance]
prompt = "default"
```

## Available Presets

### default

Full-featured prompt with directory, git, language versions, and command duration. Uses Nerd Font symbols.

### plain

Same information as default but uses ASCII characters only — no Nerd Font required.

### minimal

Directory and git branch only. Two-line with a simple `>` character. For distraction-free work.

### nerd-font

Rich prompt with Nerd Font icons for languages, git status, and system info. Requires a [Nerd Font](https://www.nerdfonts.com/) installed on the host terminal.

### pastel

Soft colors with rounded segment separators. Nerd Font recommended.

### bracketed

Each prompt segment wrapped in square brackets. Clean, structured appearance without special fonts.

## Changing Presets

1. Edit `aibox.toml`:
   ```toml
   [appearance]
   prompt = "minimal"
   ```

2. Run sync:
   ```bash
   aibox sync
   ```

The Starship config is regenerated at `.aibox-home/.config/starship.toml`. Colors are derived from the active theme.
