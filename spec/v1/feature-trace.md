# v0 behavior-to-v1 trace (reviewed source map)

The paths below are relative to the immutable v0.35.0 source commit
`9061bf76d2a1f6ac3c7093d8c605bcc4de7eddf6`. They identify the primary
behavioral source, not just a declaration. `05-acceptance.md` defines the
observable obligation. The v1 implementation component identifies responsibility; it does not
claim implementation or runtime equivalence. Each F row must acquire an
executable fixture and target evidence before parity can close.

| ID | Primary v0 behavior source | Proposed v1 owner / retained asset |
|---|---|---|
| F01 | `cli/src/main.rs`, `content_init.rs` | Native Template/starter and bounded conversion utility; never regenerate user files blindly. |
| F02 | `cli/src/content_install.rs`, `context.rs` | Optional processkit Feature/integration; standalone harness workspace is first-class. |
| F03 | `cli/src/main.rs`, `container.rs`, `runtime.rs` | Dev Container CLI for build/up/exec; exact native stop/remove gap adapter. |
| F04 | `cli/src/output.rs`, `runtime_resources.rs` | Shared Go result core, scoped runtime or local inspection. |
| F05 | `cli/src/harness_commands.rs`, `main.rs` | Bounded native exec/recovery entry with no tmux dependency. |
| F06 | `cli/src/env.rs` | Separate tested native/delegated snapshot utility; not merely multiple JSON files. |
| F07 | `cli/src/reset.rs`, `main.rs` | Explicit backup/recovery utility plus processkit-owned context recovery. |
| F08 | `cli/src/prune.rs` | Native scoped cleanup with ownership proofs and consent. |
| F09 | `cli/src/update.rs`, `main.rs` | Distribution tooling plus generated CLI completion/help. |
| F10 | `cli/src/generate.rs`, `config.rs` | Native `devcontainer.json`, Dockerfile and Compose inputs. |
| F11 | `cli/src/config.rs`, `generate.rs` | Native user/env/mount/lifecycle settings and local keepalive helper. |
| F12 | `cli/src/runtime_home.rs`, `seed.rs` | Explicit persistent scoped home/mounts and native harness auth. |
| F13 | `cli/src/addons.rs`, `addon_loader.rs`, `addon_registry.rs`, `addons/` | Pinned Features/packages; per-tool route in tool selection ledger. |
| F14 | `cli/src/harness_commands.rs`, `tmux/layouts.rs` | Feature install plus native harness files and retained pane launch glue. |
| F15 | `cli/src/mcp_registration.rs`, `runtime_sync.rs` | Native harness MCP files; processkit owns its own entries. |
| F16 | `cli/src/preauth.rs`, `config.rs` | Native harness policy plus independent operator authority. |
| F17 | `cli/src/content_install.rs`, `compliance.rs` | Processkit-owned canonical `AGENTS.md` and provider pointers. |
| F18 | `cli/src/tmux/layouts.rs`, `images/base-debian/config/bin/aibox-tmux-session.sh` | tmux layout assets and bounded local session helper. |
| F19 | `images/base-debian/config/tmux/tmux.conf`, `cli/src/tmux/mod.rs` | Retain tmux/PowerKit keymap and native settings. |
| F20 | `cli/src/tmux/sync.rs`, `theme_cmd.rs` | Explicit local override and atomic managed refresh. |
| F21 | `cli/src/tmux/status.rs`, `images/base-debian/config/tmux/powerkit-plugins/` | Retain status renderer/plugins; native tmux/PowerKit files. |
| F22 | `cli/src/runtime_resources.rs`, `images/base-debian/config/tmux/powerkit-plugins/` | Local status collectors with truthful unavailable state. |
| F23 | `cli/src/hook_registration.rs`, `cli/src/templates/aibox-agent-signal.sh` | Retain/specialize title hooks and terminal escape channel. |
| F24 | `cli/src/hook_registration.rs`, `cli/src/templates/aibox-codex-notify.sh` | Harness-native signal adapters plus manual fallback. |
| F25 | `cli/src/hook_registration.rs`, `cli/src/templates/aibox-agent-signal.sh` | Opt-in native terminal bell/OSC notification path. |
| F26 | `cli/src/themes.rs`, `cli/src/theme_cmd.rs` | Preserve palette data and every family/variant; native consumers. |
| F27 | `cli/src/themes.rs` | Preserve semantic palette renderer and accessible fallbacks. |
| F28 | `cli/src/themes.rs`, `cli/src/seed.rs` | Retain per-tool theme adapters until native replacements pass equivalence. |
| F29 | `cli/src/config.rs`, `cli/src/seed.rs` | Starship native config plus all eight presets/alias. |
| F30 | `images/base-debian/config/yazi/init.lua`, `keymap.toml` | Retain Yazi navigation, popup and async preview glue. |
| F31 | `images/base-debian/config/bin/aibox-preview.sh`, `config/yazi/plugins/rich-preview.yazi/` | Existing renderers and bounded retained preview glue. |
| F32 | `images/base-debian/config/yazi/plugins/sqlite-preview.yazi/`, `tabular-preview.yazi/` | Retain read-only data preview plugins. |
| F33 | `images/base-debian/config/yazi/plugins/preview-options.yazi/`, `keymap.toml` | Retain stateful preview option/keymap behavior. |
| F34 | `images/base-debian/config/yazi/plugins/rich-preview.yazi/`, `keymap.toml` | Retain directory tree and selectable read-only preview. |
| F35 | `images/base-debian/config/yazi/keymap.toml`, `cli/src/seed.rs` | Native clipboard channel and existing copy actions. |
| F36 | `images/base-debian/config/bin/aibox-preview.sh` | Retain standalone viewers/watch and renderer fallbacks. |
| F37 | `cli/src/seed.rs`, `images/base-debian/config/` | Native Vim/Git/XDG files; preserve user overrides. |
| F38 | `cli/src/latex.rs`, `images/base-debian/config/bin/aibox-latex-preview.py` | Container-local build/watch and read-only Compose sidecar. |
| F39 | `cli/src/audio.rs`, `config.rs` | Optional audio Feature plus separately authorized host setup. |
| F40 | `addons/tools/browser-testing.yaml` | Playwright/axe with matched native browser packages. |
| F41 | `addons/tools/diagramming.yaml`, `data-visualization.yaml`, `mermaid.yaml` | Existing renderer packages; no new graphics engine. |
| F42 | `addons/tools/`, `addons/languages/` | Maintained Features/packages, per-tool evidence required. |
| F43 | `cli/src/doctor.rs`, `integrity.rs` | Read-only Go doctor with stable findings and scoped checks. |
| F44 | `cli/src/workspace_manifest.rs`, `provider_backend.rs`, `output.rs` | Versioned Go results and read-only catalog resources. |
| F45 | `scripts/build-macos.sh`, `cli/src/runtime.rs` | Tested darwin/linux binary and image qualification cells. |
| F46 | `cli/src/config.rs`, `cli/src/harness_commands.rs` | Profile selection with truthful headless limits. |
| F47 | `docs-site/content/docs/` | Version-aligned user docs and generated MCP guidance. |
| F48 | `cli/src/tmux/status.rs`, `images/base-debian/config/tmux/powerkit-plugins/modelstatus_provider.sh` | Retain opt-in provider polling without conflating quota/account/session metrics. |

New N01/N02 have no v0 source. Their candidate and acceptance routes are in
`04-reuse-targets.md` and `05-acceptance.md`; do not count them as parity.
