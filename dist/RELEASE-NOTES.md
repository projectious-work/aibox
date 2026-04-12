# aibox v0.17.17

Feature release. Two quality-of-life improvements to `aibox init` / `aibox sync`,
plus a dependency-declaration fix for the OpenAI addon.

## New features

### `aibox.toml` — inline addon documentation

`aibox.toml` now includes a block of inline comments above each
`[addons.<name>.tools]` section, generated from the addon definition:

```toml
# OpenAI Codex CLI
# codex: pin with codex = { version = "x.y.z" }
[addons.ai-openai.tools]
```

For addons with a curated version list, available versions are shown with the
default marked:

```toml
# Python
# python: 3.11 | 3.12 | 3.13 (default)
# uv: pin with uv = { version = "x.y.z" }
[addons.python.tools]
python = { version = "3.13" }
```

Tools that are disabled by default also get a hint comment so users know they
exist and can be enabled.

### `aibox init` — AI provider addons included in transitive `requires` expansion

AI provider addons were previously excluded from the transitive dependency pass
during `aibox init`, causing their `requires` entries to be silently ignored.
For example, `ai-openai` requires `node` (Codex CLI is an npm package), but
`node` was never pulled in automatically.

The init flow now combines user-selected addons with AI provider addons before
expanding `requires`, so:

```
aibox init --ai openai
```

correctly adds `node` (and any other transitively required addons) without
the user having to specify them manually.

## Bug fix

### `ai-openai.yaml` — declare `requires: [node]`

The `ai-openai` addon now formally declares `requires: [node]` in its addon
definition. This declaration is what enables the transitive expansion fix above.
Projects that already have `node` in their `aibox.toml` are unaffected.

## Housekeeping

- This project's own devcontainer (`aibox.lock`, `aibox.toml`,
  `.devcontainer/`) synced to aibox v0.17.16.
