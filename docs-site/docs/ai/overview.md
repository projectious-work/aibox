---
sidebar_position: 1
title: "Overview"
---

# AI Providers

aibox integrates AI coding agents into your development container. Agents are installed per-project via the addon system — not baked into the base image.

## Configuration

Set your providers in `aibox.toml`:

```toml
[ai]
providers = ["claude", "aider"]
```

This automatically installs the corresponding CLI tools, mounts their config directories, creates Zellij layout tabs, and adds VS Code terminal profiles.

## Supported Providers

| Provider | CLI Tool | Install Method | Config Directory |
|----------|----------|---------------|-----------------|
| [Claude](claude.md) | `claude` | curl installer | `.claude/` |
| [Aider](aider.md) | `aider` | uv tool | `.aider/` |
| [Gemini](gemini.md) | `gemini` | npm | `.gemini/` |
| [Mistral](mistral.md) | `mistral` | pip | — |

## How It Works

When you set `providers = ["claude"]` and run `aibox sync`:

1. The `ai-claude` addon is automatically enabled
2. The generated Dockerfile includes the Claude CLI install command
3. The generated docker-compose.yml mounts `.aibox-home/.claude/` into the container
4. Zellij layouts include a Claude tab/pane
5. VS Code devcontainer.json gets a Claude terminal profile

## Multiple Providers

You can use multiple providers simultaneously:

```toml
[ai]
providers = ["claude", "aider"]
```

In the **cowork** Zellij layout, multiple providers are stacked side-by-side for parallel AI collaboration. In **dev** and **focus** layouts, each provider gets its own tab.

## API Keys

API keys are stored in each provider's config directory inside `.aibox-home/`. Since these directories are mounted from the host, keys persist across container rebuilds.

- **Claude:** Authenticated via `claude` CLI login (stored in `.claude/`)
- **Aider:** Set `OPENAI_API_KEY` or `ANTHROPIC_API_KEY` in `[container.environment]`
- **Gemini:** Set `GOOGLE_API_KEY` in `[container.environment]`
- **Mistral:** Set `MISTRAL_API_KEY` in `[container.environment]`
