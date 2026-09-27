# 5. Feature ledger and acceptance criteria

Every row below is required for the replacement release and requires v1
verification. Baseline references are the immutable documentation/asset/type indexes
in `ledger/`; existing tests are inputs to the new fixtures, not proof of v1.
The [behavioral source map](feature-trace.md) binds each F row to its
primary implementation and v1 implementation component; it complements rather than
replaces the source declaration ledgers. No F row closes on an inventory count.

## Behavioral feature ledger

| ID | Observable v0 capability / new requirement | Acceptance group |
|---|---|---|
| F01 | Initialize new repository and adopt existing repository without losing files | AC-CONFIG, AC-MIG |
| F02 | Minimal/harness-only workspace; optional processkit packages/source/version/fork | AC-PK |
| F03 | Build, cache bypass, rebuild, start, attach, stop and remove | AC-BUILD, AC-LIFE |
| F04 | Runtime status/detail and structured table/JSON/YAML information | AC-CLI, AC-DOCTOR |
| F05 | Emergency harness entry without tmux/Yazi/status dependency | AC-LIFE |
| F06 | Named environment save/list/detail/switch/delete, including process context | AC-MIG |
| F07 | Scoped backup, soft reset/recovery, context reset guidance | AC-MIG |
| F08 | Scoped cache/container/worktree reclamation, preview and dirty-data protection | AC-MIG, AC-SEC |
| F09 | Installer/update/check/dry-run/uninstall/completion/help/version outcomes | AC-CLI, AC-RELEASE |
| F10 | Image selection, native custom build, Compose overrides, named paths and local overrides | AC-CONFIG, AC-BUILD |
| F11 | Container user/hostname/UID, environment, extra volumes, network keepalive | AC-CONFIG, AC-LIFE |
| F12 | Persistent private home/auth/cache/SSH/tool state and disable/preserve/purge choices | AC-MIG, AC-SEC |
| F13 | Every addon/tool enabled/version choice and transitive required capability | AC-TOOLS |
| F14 | Harness enable without install, ordered launch, multiple panes and native config/hooks | AC-HARNESS |
| F15 | Shared/custom and personal MCP servers; managed-entry preservation; gateway choice | AC-HARNESS, AC-PK |
| F16 | Execution policy axes, per-harness overrides, allow/deny/ask semantics | AC-SEC, AC-HARNESS |
| F17 | AGENTS canonical/pointer behavior and provider-neutral repository context | AC-PK |
| F18 | tmux layouts ai/dev/focus/cowork, optional LazyGit window, ordered harnesses | AC-UX |
| F19 | Pane/window/session navigation, split/resize/zoom/frame/popup/scrollback/search/chooser keys | AC-UX |
| F20 | Configurable prefix/session name, layout/theme chooser, live switch and state reset | AC-UX |
| F21 | Header/status extended/plain/disabled, element order/toggles/labels/separators/spacing/refresh | AC-UX |
| F22 | Resource/OOM/process/MCP/migration metrics, log pane, cloud/forge/Git/network/provider status | AC-DOCTOR, AC-UX |
| F48 | Existing opt-in provider quota polling/admin-usage rollups, cache/timeout controls and explicit credential references | AC-DOCTOR, AC-SEC |
| F23 | Terminal/tab header status incl. Ghostty; templates, symbols, identity, priority, bounds and restoration | AC-TITLE |
| F24 | Working/question/done/error/idle signaling, native/partial/manual fallback, completion TTL | AC-TITLE |
| F25 | Optional OSC-9/bell notifications, transition-only behavior and unsupported-terminal fallback | AC-TITLE |
| F26 | All 37 theme families/76 variants; dark/light/auto and per-family overrides | AC-UX |
| F27 | Semantic palette/decoration mapping, emphasis levels/overrides, contrast, monochrome/no-color | AC-UX |
| F28 | Consistent theme consumers: tmux/PowerKit, Vim, Yazi, Starship, LazyGit, bat/delta, fzf/eza, less/man, lnav and supported AI TUIs | AC-UX |
| F29 | Eight prompt presets, legacy alias and custom native prompt configuration | AC-UX |
| F30 | Yazi file navigation/open in Vim/fullscreen popup, async preview refresh | AC-PREVIEW |
| F31 | Raster/SVG/EPS/PDF/Markdown/text/video/archive preview renderers and terminal fallbacks | AC-PREVIEW |
| F32 | SQLite read-only preview, CSV/TSV and XLS/XLSX formatting | AC-PREVIEW |
| F33 | Independent preview scrolling, non-wrapping horizontal view, persistent line numbers/wrap | AC-PREVIEW |
| F34 | Recursive directory-size tree with hierarchy/totals; selectable read-only Vim preview | AC-PREVIEW |
| F35 | All supported Yazi copy actions and Vim selection through tmux/host clipboard channel | AC-PREVIEW, AC-SEC |
| F36 | Standalone image/PDF/SVG/EPS viewers and watch/reload; Excalidraw text fallback only | AC-PREVIEW |
| F37 | Vim user settings, persistent undo, rg integration and Git/XDG credential-helper behavior | AC-UX, AC-HARNESS |
| F38 | LaTeX multiple documents/engines/options/output/cache, build/watch and browser PDF preview | AC-PREVIEW |
| F39 | Optional audio output/input, tool install control, server selection and host diagnostics/setup | AC-AUDIO |
| F40 | Browser testing, Chromium default/optional Firefox/WebKit, axe accessibility checks | AC-TOOLS |
| F41 | Graphics D2/Graphviz/Vega/Vega-Lite/Mermaid, optional browser-heavy rendering | AC-TOOLS |
| F42 | Infrastructure/cloud/Kubernetes, security/release/quality and documentation tool bundles | AC-TOOLS |
| F43 | Configuration/integrity/security diagnostics; unknown keys/migrations/remediation | AC-DOCTOR, AC-MIG |
| F44 | Human and machine-readable catalog/provenance/workspace/harness projections | AC-CLI, AC-GUIDE |
| F45 | Linux/macOS binary targets, amd64/arm64 image/tool compatibility | AC-TARGET |
| F46 | Human-dev/headless selection with truthful current limits | AC-TARGET |
| F47 | Existing theme gallery, guides, cheatsheet, compatibility and upgrade docs | AC-GUIDE, AC-DOCS, AC-RELEASE |
| N01 | Optional left tmux pane for every current-session harness and available usage/quota/context data | AC-SIDEBAR |
| N02 | Convenient local diff plus PR review with existing tools | AC-REVIEW |

The exhaustive parameter/tool/command records specialize these rows. A test
covering one theme, addon or default field does not close the whole group.
Every current choice needs either evidence or an explicit blocked result.

## Acceptance groups

### AC-CONTRACTS — executable contracts and quality limits

Apply chapter 18 and the request/result/policy/receipt schemas to all operations.
Validate with Draft 2020-12 including date-time formats and external references;
the same positive/negative corpus must pass Go, CLI and MCP SDK adapters.
Cover absent containers and first-build failure, invalid input without digest,
irrelevant fields, wrong contexts, data payloads, partial effects, cancellation,
timeouts, expired grants, receipt crash recovery and a new client inspecting
an old operation. Verify Compose exclusive/shared membership and no volume
deletion. Measure every R-LIMITS bound using the declared reference environment;
publish observed values, not just pass labels. Reject oversize inputs before
delegation, and bound/truncate diagnostic output without truncating JSON syntax.

### AC-CONFIG — configuration equivalence

Given a default v0 fixture and then each non-default field/alias fixture,
conversion preserves effective intent in standard/native destinations.
Cover omitted/false/empty/explicit values, arrays/order, dynamic map keys,
aliases, malformed/unknown fields, contradictory compatibility forms, JSONC,
Unicode/spaces in paths and source precedence. For each `CFG:*` row record
fixture, expected effective value, destination, negative case and evidence.
The customization schema must declare each of its 115 v0 UX leaf
paths and reject unknown keys, unknown enumerated choices and null; validate
the representative example in chapter 2. For standard fields and Feature
options, verify the *resolved* native Dev Container configuration, not only
the source JSON text. A field remains blocked until its exact destination and
effective-value fixture are recorded; the path-pattern census is not proof.
No config mutation command appears in CLI/MCP. Pure validation causes no
filesystem/network/process effects beyond explicitly requested reads.

### AC-BUILD — reproducible upstream delegation

Given identical pinned source/base/Features and tool versions, build with
aibox and directly with upstream CLI yields equivalent declared capabilities.
Record actual digests/provenance without promising byte-identical builds when
the dependency chain is not reproducible. Test cache/no-cache, failed fetch,
checksum mismatch, unavailable architecture and cancellation. The running
workspace remains intact after build-only failure. Changed build inputs
invalidate old release evidence. No secret appears in layers or logs.

### AC-LIFE — lifecycle and exact identity

In a disposable authorized target, start a missing/stopped/running environment;
attach interactively and through recovery; stop twice; remove twice; rebuild
with persistent user data. Verify defined no-change outcomes. Run unrelated
same-name resources in another context and verify zero changes there. Simulate
runtime disappearance, permission denial, interrupted replacement and concurrent
commands. Report actual partial state and safe recovery. Removing/recreating
an old name cannot make an old request target the new runtime.

### AC-CLI — shared interface conformance

For each operation, identical logical inputs through CLI/MCP produce equivalent
core results/errors/effects. Test stdout/stderr separation, JSON/YAML projection,
exit codes, invalid arguments, timeout/cancellation, child error propagation,
noninteractive execution and help/completion. All 14 v0 declaration groups,
66 atomic actions and 74 arguments (including global flags) have a documented
disposition; every old
action has a usable replacement, not simply
“removed”. Test scripts against published
binary, not only Go functions. Verify protocol negotiation with selected clients.

### AC-CONFIG-SOURCES — layered process settings and provenance

Test built-in, system, user, eligible project, multiple explicitly selected
env files, inherited process environment and CLI/MCP invocation overrides in
every pairwise precedence order and representative full-stack cases. Confirm
missing optional files are reported, while malformed/unreadable files,
duplicate/unknown keys, unsupported versions and unauthorized sources fail
before side effects. Exercise `AIBOX_<SECTION>__<KEY>` parsing, null/empty/
false, private env-file permissions, relative-path anchoring, ambient-env
isolation and a deterministic redacted effective-config view. Native
Dev Container fields remain under upstream semantics; no env/CLI aibox key
rewrites Features, image, mounts or hooks. Tests use temporary homes and
prove a project/env/MCP request cannot alter operator policy, executable
selection or allowed log egress. CLI and MCP resolve equivalent settings for
their shared use-case inputs; MCP callers cannot change server logging.

### AC-LOGGING — sink, redaction and evidence separation

Validate the closed log-event schema and each declared stderr, rotating-file
and service-manager capture route on CLI, operator MCP and local MCP. Machine
stdout contains only its result/protocol while progress and diagnostics use
stderr; `NO_COLOR` and non-TTY behavior are deterministic. Test levels,
filters, event correlation, initiating/executing actor attribution, child
source attribution, partial/cancelled outcomes, concurrent writes, rotation,
retention, backpressure and unavailable/closed sinks. A required sink or
receipt failure refuses mutation before effects; mid-operation failure
reports degraded observability without replaying work. Known secrets,
private paths and chunked child output never reach any sink, buffer or test
artifact. Rotating logs cannot erase durable operation receipts/evidence.

### AC-SEC — authority and malicious-input tests

From the in-container agent, attempt operator mode, forged tool calls,
host-socket access, path traversal/symlink escape, command injection, altered
mounts/hooks/privilege, project-supplied “allow” and replay against a new target.
All host/lifecycle operations remain unavailable or denied with no host side
effects. Authorized local tmux edit/reload succeeds. A changed input digest
invalidates operator approval. Secret fixtures are absent from returned errors,
logs, evidence and builds. Validate names and scope again before deletion.
Test host trust review separately from schema validity and local user trust.

### AC-SECRETS — credential transfer and containment

Run chapter 17's file bind, agent socket, native environment, SOPS/age,
OpenBao/provider and BuildKit canary matrix on supported targets. Verify the
credential-free starter, custom user/home, consent, least-scope mount,
redaction, broker failure, cleanup and absence of host authority. A failed
mode is reported precisely; no silent downgrade or credential copy is allowed.

### AC-TOOLS — complete optional tool matrix

For every `ADDON:*` tool: installation/version smoke test on both image
architectures, minimal enabled/disabled cases, declared dependencies, supported
version alternatives, absence when disabled, conflicting selections and failed
downloads. Verify base toolchain separately. Test one real browser/graphics/doc
build per relevant tool family, not only `--version`. Verify no disabled heavy
component becomes an implicit dependency. Record selected package/Feature and
license/provenance clearance per row. Use pairwise interactions plus explicit
known-risk full stacks; justify combinations rather than claiming exhaustive
Cartesian-product testing.

### AC-HARNESS — installation, state and integration

For each supported harness identity, test enabled/install/version/order choices,
host-only integration, startup, native config, preserved user entries,
authentication persistence and per-harness execution restrictions. Hook state
coverage is native/partial/wrapper/manual per capability, never claimed uniform.
An unsupported permission control fails visibly rather than silently broadening
access. Optional MCP/processkit absence must not prevent harness-only operation.
Credentials used for live checks are owner-provided in the correct authority
context, not copied into fixtures or logs.

### AC-PK — optional processkit

Fresh harness-only environment has no mandatory processkit files/services/hooks.
Enabled integration verifies pinned source/version/hash, compatible fork/source,
content/user override preservation, custom MCP entries and processkit-owned
migration/doctor workflows. Missing supported API is a blocking integration
result. Removing aibox's entity commands does not remove access to those user
workflows. Include source changes and same-version manifest changes in tests.

### AC-UX — terminal and customization parity

Parameterize all four layouts, harness counts/order, optional panes, tmux
prefix/session settings, header modes and every theme/prompt choice. Verify
key actions, visibility/contrast, resize/narrow panes, no-color/plain-font,
light/dark/auto, emphasis overrides and native AI-TUI fallback tiers. Use both
structural assertions and deterministic visual baselines across the palette
matrix; visual comparison of representative screenshots supplements automated tests.
User changes survive refresh/rebuild. Failed reload leaves working config;
session reset requires consent. Plugin version conflicts fail actionably.
Prove each retained UX field's single owner (Feature option, native file or
aibox runtime namespace), and direct-upstream build without optional aibox UX
integration. A Feature option change may require rebuild; a live native edit
must not be misreported as a rebuild-only setting.

### AC-TITLE — attention and outer terminal headers

Replay lifecycle fixtures for each native/partial/manual integration; test
question/error priority across background panes, permission reply clearing,
done TTL, unknown/idle states, model/effort identity and message bounds.
Observe titles in Ghostty and another supported terminal, with no competing
shell title writer. Restore original title on exit; sanitize escape characters.
Notifications occur only on configured transitions. Unsupported notifications
degrade without breaking titles or harness operation. No host API permission
is required for ordinary terminal output.

### AC-PREVIEW — files, data and document workflows

Maintain fixtures for every F31/F32 format plus malformed/large/untrusted files,
spaces/Unicode, stale cache and no-renderer fallback. Preview is read-only and
bounded; never execute document content. Verify J/K versus j/k, horizontal
view, wrap/line-number persistence, selectable Vim preview, directory totals,
Yazi/Vim clipboard and watch refresh inside tmux. Test PDF/SVG/EPS image/text
fallback paths. LaTeX multi-document build/watch stays inside the dev-container;
sidecar serves completed PDFs read-only on loopback by default. Reject public
binding without explicit operator authority; no source/secret directory leakage.

### AC-AUDIO — optional voice path

Disabled audio installs/opens nothing unexpected. Enabled client tooling honors
install selection and server settings. Test playback and explicitly authorized
input on qualified host routes; distinguish missing host service from client
misconfiguration. Host setup is external and separately approved; local agent
cannot create host services or expand network exposure. No silent microphone
enablement, broad unauthenticated listener or stored transcript/credential leak.

### AC-MIG — persistence, named state, recovery and cleanup

Fixtures include custom themes/hooks, personal env/MCP data, extra mounts,
disabled harness auth, dirty/untracked worktrees, multiple named environments,
processkit shared/non-shared context and unknown config fields. Preview makes
no changes; explicit conversion has scoped backup and a complete mapping.
Interrupt each stage and verify old workspace recovery. Switch named states
without cross-contamination. Ordinary stop/remove/rebuild retains user data;
purge/cleanup enumerates exact resources and requires confirmation. Backup
verification includes permissions and restoration, not just file existence.
Test default `/home/aibox` named-volume retention, explicit host-directory
bind, custom Dockerfile/user/home plus aligned mount target, first-use
initialization on empty/masking volumes, and remote-daemon path semantics.
The default starter creates no `.aibox/` or `.aibox-home/`.

### AC-DOCTOR — actionable diagnostics

Inject each finding class and assert stable code, severity, provenance,
remediation link and correct machine status/exit. Missing data is unknown;
local host checks are not-authorized/skipped rather than healthy. Verify no
auto-fix/migration side effects and no hook execution during validation.
Read-only inspection includes bounded resource/status/log data with redaction.
Provider API rate limits, organization usage and subscription five-hour/weekly
allowances are different measurements. Preserve existing opt-in status controls
without relabeling one as another. Billable probes/admin-key access require
explicit consent, bounded polling and disclosed cost; no automatic activation
to populate the new sidebar. Do not copy stale comments about price into the
product contract; qualify actual source semantics and current authorization.

### AC-SIDEBAR — new optional session overview

Enable a left pane in each supported layout; collapse/resize/disable without
losing other panes. Display every harness in the current session and no unrelated
session. Distinguish working, idle, exited/stopped, question and error based
on reliable events; indicate unknown/stale states. Show per-session usage and
context numerator/denominator only when known. Show five-hour/weekly account
allowances/reset times where supported, labeled account-wide and deduplicated.
Test multiple panes sharing one account, unavailable endpoints, expired auth,
rate limiting, stale data and event storms. Bounded collection must not make
the shell/harness unusable. Exact refresh/cost/performance budgets and per-harness
telemetry sources are G06. Monitoring alone must not send prompts, approve
permissions, manage host containers or install a new workflow governance layer.

### AC-REVIEW — new diff/PR workflow

Inspect working-tree/staged/branch diffs with Git/delta/LazyGit, including large
and binary changes. From the optional PR TUI, find a PR, read description and
checks, navigate diff, safely check out with dirty-worktree protection, and
perform explicitly authorized comment/review submission using native tooling.
Document web fallback for unsupported inline review features. Test no-auth,
read-only credentials, network error and non-GitHub repository behavior.
No implicit approval, push or merge and no aibox-built diff engine.

### AC-TARGET — target qualification

Execute the relevant lifecycle/UX/audio/browser/preview gates on each supported
target cell with exact versions/architecture recorded. Native host-only tests
are invoked by an authorized host operator through bounded handoffs; unavailable hosts remain blocked.
Headless selection runs without requiring interactive tmux; do not advertise
new PID-1/agent-runner guarantees without an additional contract and tests.

### AC-GUIDE — progressive disclosure and docs

Discover/read each declared resource with no arbitrary path escape, credentials,
private process context or host runtime requirement. Run recipe steps against
the specified release and assert expected results; unsupported version/topic
is explicit. Human docs, Markdown, MCP and discovery files agree. Measure that
the initial index is bounded and a theme query retrieves the specific recipe,
not the full manual. Client compatibility is recorded per protocol version.

### AC-RELEASE — delivery and exact evidence

Validate Go binary targets, OCI/Feature/Template artifacts, checksums, SBOM,
scanning, signing/provenance and clean-room published-artifact installation.
Use candidate-bound evidence with no unexplained skips; independently verify
published artifacts. Local docs build/deploy, no project-authored GitHub
Actions. Release notes, compatibility matrix, roadmap and migration guide
describe implemented behavior only. Version-line artifacts and documentation must agree.

### AC-DOCS — shared Hugo site and versioned prerelease truth

For **every** roadmap phase, verify its `docs` deliverable was updated with
the implementation, name changed documentation paths and audiences, run the
local v1 Hugo build from the phase candidate commit, and check links, schemas
and changed executable examples. A build failure or stale/misleading page
blocks that phase's completion; a development note alone does not satisfy
the gate. Check v0 current content and Releases labels remain unchanged.
V1-01 records the existing docs-build baseline; V1-20 establishes the shared
Hugo route before any implementation phase. Later phases rerun the build,
including internal phases and phases that do not publish a public preview.

Build v0 current, v1 preview, an alpha archive and a beta archive through the
same Hugo/brand-theme module build route. Verify a pinned shared theme,
line-specific content and metadata, reproducible local output and no Docsy
runtime or submodule requirement. The Releases dropdown on both lines shows
v0.x as the current product release and alpha/beta only as v1 previews; it
includes published archives, keyboard/active-state semantics and a no-JS
fallback. Publishing either line preserves the other line's files and
manifest entries. Broken cross-line page paths fall back to a useful landing
page, not a 404. Validate each prerelease candidate's getting started,
configuration, migration, authority and feature pages against its exact
binary/Feature behavior and release notes. Documentation, public Markdown,
discovery files and MCP guidance agree; no stale v0-only or reverted-v1
instructions remain in current v1 pages. Stable-v1/current labels are applied only when publishing the stable release.

## Evidence record and closure

Each record names requirement/feature/field/tool IDs, case and fixture, target
versions, source/candidate digest, executor/time, command, expected/actual result,
artifact references and pass/fail/blocked status. No secret payloads. A group
closes only when all applicable rows have evidence or accepted non-applicability.
Verify the complete ledger; counts alone never prove parity.
