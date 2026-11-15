# 11. Workspace experience and v0 parity

This chapter is the implementer's contract for R-PARITY, R-LOCAL and
R-OWNERSHIP. The source-of-truth inventory is `feature-trace.md` plus the
configuration, addon and runtime-asset ledgers. All F01–F48 behaviors remain
required even where not repeated below. A missing row is a failing parity
test, not implicit permission to retire it. New N01/N02 are additional.

## Harnesses and durable user state (V1-10)

Support the existing eleven harness identities, including host-only Cursor.
The Feature selection installs an in-container binary only for applicable
harnesses; the aibox `harnesses` namespace selects enablement, order and
session launch. Never equate installed with enabled. Preserve each provider's
own login/config/state directory on a scoped persistent home mount, not in an
image layer, generated file or project Git. Disabled harnesses retain their
state. A harness launcher checks binary presence and configuration and reports
a scoped warning; it must not silently install, authenticate or grant host
access. Per-harness hooks and MCP configuration are merged only through each
provider's supported file format, with ownership markers and a preserved
personal section. Processkit-owned MCP configs and personal MCP entries must
not be overwritten by aibox refresh. Cursor integration is limited to
documented host editor behavior; no fictional in-container CLI is required.

The container entrypoint restores a normal shell when no harness is selected
or a selected harness fails. Recovery entry bypasses tmux and all status/
preview plugins. Home persistence is explicit and portable across rebuild;
cross-major state migration gets a backup and provider compatibility check.
No agent inside the container receives a runtime socket, unrestricted host
filesystem mount, operator MCP endpoint, host-management credential or capability to
start/stop/rebuild the container. Local agent edits to tmux/theme/Yazi files
are allowed as the container user and activate only through local refresh.
Explicit native workspace/home mounts and scoped credential transfers in
chapter 17 are allowed; they never confer host-management authority.

## tmux, status, titles and prompt (V1-11–09)

Recreate the four v0 layouts, pane ordering, window naming, attach/resume,
navigation keys, copy mode, scroll behavior and all documented mode/status
options from `configuration-types.json`. Golden tests compare a generated
session plan (windows, panes, commands, working directories) for each layout
with v0 behavioral fixtures; interactive smoke tests verify focus and key
bindings. Header bars expose current project, active harness, branch/status,
and attention state using a bounded, sanitized rendering contract. User
customizations cover placement, colors, separators, left/right segments,
refresh interval and compact mode, including the v0 option names. The header
must remain legible in narrow terminals, low-color terminals and light/dark
themes. A provider hook that cannot supply state shows `unknown`, not a fake
idle indicator.

Terminal title forwarding uses terminal-standard escape sequences emitted by
tmux with sanitization and length bounds; prove behavior in Ghostty and a
generic xterm-compatible terminal, and document terminals that disable
title changes. The title includes a stable project/session prefix and a
short current-window/attention state, with no file paths, tokens or prompts.
When title forwarding is unsupported it degrades to internal tmux header
only. Refresh reloads tmux configuration without killing existing panes and
reports settings requiring a new session.

Preserve all 37 theme families, 76 variants, eight Starship presets,
semantic-color mappings, icons, font/powerline fallbacks and accessibility
modes enumerated in the v0 ledger. The renderer maps one aibox theme intent
to native tmux, PowerKit, Yazi, Starship and editor settings; exact native
files remain overridable. Snapshot tests cover every family/variant and
light/dark/auto modes; representative interactive tests check contrast,
status, selection and preview visibility. `auto` follows available terminal
signals with a documented fallback and never queries host appearance from
inside the container.

## Preview, editor, LaTeX, audio and processkit (V1-13–12)

Preserve Yazi keybindings, file-type rules and fallback rendering for text,
Markdown, images, PDF, SVG, video/audio, archives, tabular/JSON, SQLite and
other formats in the runtime-asset ledger. Use the listed native tools where
licensed/qualified; retain a small preview glue script only where a native
plugin does not preserve behavior. Preview runs read-only against the chosen
file, with time/size/output bounds and no automatic execution of document
macros or archive entries. Missing optional renderer yields a useful text
fallback and doctor finding. Vim, clipboard behavior, independent viewers
and watch commands remain available per the F rows; do not force a desktop
clipboard on headless targets.

LaTeX supports multiple named documents with explicit source, build, watch
and preview paths. Build/watch run inside the workspace; the preview sidecar
only serves completed PDFs read-only over a loopback-bound port. Watchers
debounce changes, retain the last good PDF on error, and stop with the
container. The generated per-document commands and URLs are deterministic;
remote access requires user-controlled port forwarding. Compilation cannot
be triggered via the preview HTTP service.

Audio remains opt-in. The Feature provides client packages, while the
definition declares a platform-qualified native environment/mount bridge.
Host PulseAudio/PipeWire/OrbStack/Desktop service setup remains a human
host action. A local agent can test input/output devices and report status,
but cannot create host services or alter host audio policy. With audio off,
no bridge mount or audio packages are present. A missing host bridge degrades
the workspace rather than failing unrelated tools.

Processkit is optional. Integrate through its supported installer/config
interfaces, pinned source and package tier; never copy its entity engine,
write its `context/` directly or regenerate a user's AGENTS instructions
unasked. Preserve source forks, package selection, user skills and personal
MCP entries during refresh and migration. A workspace without processkit
retains its harness and terminal functions. Compatibility tests use a real
pinned processkit fixture and a no-processkit fixture.

## N01: observational harness sidebar (V1-17)

A left tmux pane is opt-in, collapsible, resizable and scoped to the current
tmux session. It lists harness process/session identity and state (`working`,
`idle`, `waiting_for_user`, `stopped`, `unknown`) with observation time and
source. Polling is bounded and stops when hidden. Prefer a qualified existing
tmux status tool; an aibox adapter is allowed only for missing per-provider
signals. Never infer a question solely from low CPU or a stopped process.
Context-window usage, current-session usage and 5-hour/weekly limits appear
only when a documented provider-local API or hook exposes them reliably.
Absent/unavailable/shared-account quota data displays `unavailable` and
source/freshness, never an invented percentage. The sidebar is observational:
no account scraping, credential collection, remote control or provider
automation. Test simultaneous harnesses, shared quotas, stale hooks, pane
close/reopen and narrow layouts.

## N02: local diff and PR review (V1-18)

The default diff path uses Git plus delta; optional LazyGit provides staging,
navigation and side-by-side diff where supported. A qualified `gh-dash` plus
`gh` offers GitHub PR browse/review if the user opts in and authenticates.
Existing forge web UI remains the fallback for inline comments or actions
the TUI cannot implement correctly. Read-only local diff works offline and
without a forge token. Network writes (review submission, comment, push,
merge) always use the native `gh`/forge tool with its own confirmation and
credentials; neither the local agent nor operator MCP gets a generic PR-write
capability. Compare two commits, staged/unstaged changes, binary diff
fallback, conflict state and large changes in acceptance tests. Document
screen-reader/keyboard limitations of selected TUIs rather than claiming
unqualified accessibility.
