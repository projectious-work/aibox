---
sidebar_position: 3
title: "Local Config (.aibox-local.toml)"
---

# Local Config (.aibox-local.toml)

`.aibox-local.toml` is a personal, gitignored overlay that sits next to `aibox.toml` in the project root. It exists for secrets and per-developer settings that must never be committed to version control — API tokens, personal credential paths, and similar values that differ between contributors.

## Why it exists

`aibox.toml` is committed and shared across the team. That's the right place for project-wide settings: container name, addons, processkit version, shared environment variables, and so on. But tokens and personal bind mounts don't belong there. `.aibox-local.toml` gives every developer a private escape valve without requiring `.gitignore` discipline on every secret.

## Location and gitignore

`.aibox-local.toml` lives in the **project root**, next to `aibox.toml`:

```
my-project/
├── aibox.toml               ← committed, shared
├── .aibox-local.toml        ← gitignored, personal
├── .devcontainer/
└── context/
```

`aibox init` and `aibox sync` automatically add `.aibox-local.toml` to `.gitignore`. You do not need to do this manually.

## Supported sections

Only two sections are supported. Everything else must remain in `aibox.toml`.

### [container.environment]

Inject environment variables into the container. These are merged **on top of** any `[container.environment]` values in `aibox.toml`. If the same key appears in both files, the local value wins.

```toml
[container.environment]
GH_TOKEN            = "ghp_xxxxxxxxxxxxxxxxxxxx"
ANTHROPIC_API_KEY   = "sk-ant-api03-..."
OPENAI_API_KEY      = "sk-proj-..."
AWS_PROFILE         = "my-dev-profile"
```

### [[container.extra_volumes]]

Personal bind mounts appended **after** any volumes declared in `aibox.toml`. Each entry requires `source` (host path) and `target` (container path). `read_only` defaults to `false`.

```toml
[[container.extra_volumes]]
source = "~/.config/gh"
target = "/home/aibox/.config/gh"

[[container.extra_volumes]]
source = "~/.aws"
target = "/home/aibox/.aws"
read_only = true

[[container.extra_volumes]]
source = "~/.ssh/id_ed25519"
target = "/home/aibox/.ssh/id_ed25519"
read_only = true
```

## Merge behavior

| Section | Merge rule |
|---------|-----------|
| `[container.environment]` | Merged with `aibox.toml`; local values win on key conflicts |
| `[[container.extra_volumes]]` | Appended after `aibox.toml` volumes; no deduplication |

## Full example

A typical `.aibox-local.toml` for a developer working with Claude, GitHub, and AWS:

```toml
[container.environment]
ANTHROPIC_API_KEY = "sk-ant-api03-..."
GH_TOKEN          = "ghp_xxxxxxxxxxxxxxxxxxxx"
AWS_PROFILE       = "my-dev-profile"
AWS_REGION        = "eu-west-1"

[[container.extra_volumes]]
source = "~/.config/gh"
target = "/home/aibox/.config/gh"

[[container.extra_volumes]]
source = "~/.aws"
target = "/home/aibox/.aws"
read_only = true

[[container.extra_volumes]]
source = "~/.ssh/id_ed25519"
target = "/home/aibox/.ssh/id_ed25519"
read_only = true
```

## What is NOT supported

Everything outside of `[container.environment]` and `[[container.extra_volumes]]` is ignored. The following must remain in `aibox.toml`:

- Container name, hostname, user, `post_create_command`, `keepalive`
- `[addons]` — addon configuration
- `[processkit]` — content source and version pin
- `[skills]` — include/exclude lists
- `[ai]` — provider list
- `[customization]` — theme, prompt, layout
- `[audio]` — audio bridging

:::tip Applying changes
After editing `.aibox-local.toml`, run `aibox sync` (or `aibox sync --no-build` for a config-only refresh) to regenerate `.devcontainer/` files with the updated environment and volumes.
:::
