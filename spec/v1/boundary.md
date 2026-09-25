> Historical review input, copied for self-contained specification review. Owner confirmed boundary review points 1–5 and selected Go / thin MCP-first wrapper on 2026-09-25 (DEC-HopefulTower). Earlier open-choice language below records the review history; README.md and numbered chapters define the current draft. This snapshot does not authorize implementation.

# aibox v1.x product boundary — draft for owner review

Status: **proposed; not an accepted specification or implementation authorization**.
Prepared: 2026-09-24; revised after owner review: 2026-09-25.
Product implementation belongs in the aibox repository;
this document belongs here as a cross-product boundary proposal.

## 1. Proposed product definition

**aibox is a curated, reproducible AI development workspace distribution built
on the Dev Container CLI and standard Dev Container artifacts.** It supplies
the tools, terminal experience, harness integrations, customizations, and
tested combinations that make an environment pleasant and dependable to use.

Its users are developers working interactively with AI harnesses and operators
preparing environments in which an agent runs. A useful first proof is opening
an existing repository in a workspace with the selected tools, persistent
authentication, tmux layout, theme, and working previews, then rebuilding it
without losing personal state or project changes.

The product can remain substantial even if it has no host-side executable of
its own. Versioned content, integration, compatibility testing, migration,
documentation, and support are product responsibilities.

### Owner requirements already established

The current instruction, recorded as
`DEC-20260924_1508-PluckySage-rebase-aibox-v1-on-dev-container`, establishes:

1. Use the existing Dev Container CLI as the foundation for v1.x.
2. Preserve **all current v0.x user-facing features**, including optional
   tools, Yazi previews, themes, and customizations.
3. Prefer established open-source tools; minimize functionality recreated by
   aibox. A smaller implementation is desirable, a reduced user experience is
   not an authorized trade-off.
4. Prefer option B: a thin, MCP-first operational interface with a human CLI
   over the same core. This is the owner's current leaning; exact commands
   and packaging remain proposals. Configuration is edited in files by humans
   or agents, not through configuration-mutating CLI/MCP commands.
5. Review the product boundary first, then develop the specification under
   company standards.

The September 25 review adds the `devcontainer.json` working assumption,
explicit tmux header/title parity and two new capabilities: an optional
harness-status sidebar and a convenient diff/PR-review workflow. These new
capabilities are distinguished from existing v0 parity below.

Everything more specific below is a recommendation for review. In particular,
this draft does not select an implementation language or approve the old Go
engine, a new daemon, a plugin framework, or a universal deployment model.

## 2. Recovered direction and conflicts

The evidence reviewed comprises indexed decisions, discussion metadata,
handover records, company architecture and portfolio documents, all thirteen
company engineering standards, the prior specification and its revert, and
the current v0 release documentation and selected implementation assets.
This is a reconstruction from available records, not a claim to possess every
past conversation transcript. Some discussion tools expose metadata without
the original discussion body.

| Earlier direction | Proposed treatment now |
|---|---|
| Prefer established standards and mature tools (`QuietHarbor`, April 14) | Retain; make reuse demonstrable for each proposed component. |
| Preserve the rich existing workspace (`NimbleLark`, August 21) | Retain and strengthen to the owner's explicit all-features requirement. |
| Separate infrastructure, workspace, agent loop, and process memory | Retain the responsibility split in section 3. |
| Standard Dev Container artifacts (`AbleDew`, August 24) | Retain; the CLI itself is now the foundation, rather than merely a backend of a new engine. |
| Optional product-owned handoffs (`AbleRaven`) and direct Kaits orchestration (`LuckyMoss`) | Retain; products remain usable independently. |
| Orthogonal capabilities (`SoftSummit`) | Retain as a way to reason about requirements; it does not justify a new capability language or solver. |
| Go engine, compiled runtime adapters, `aiboxctl`, a custom template/deployment compiler (`LuckyLark`, `RoyalSpark`, `CalmWolf`, `ShinyLily`) | Reopen. These are historical proposals, not inherited implementation constraints. |
| Universal secrets model, saved plans, broad remote/Kubernetes/confidential-computing targets | Reassess individually. Preserve existing user capabilities; avoid importing speculative scope from the discarded rewrite. |
| DevPod excluded as the production foundation (`FluentFlute`) | No proposal here to reverse that choice. Its old assessment is not a fresh evaluation of today's project. |
| July stabilization briefing | Historical; not the new v1 baseline. |

The earlier [v1 specification PR #455](https://github.com/projectious-work/aibox/pull/455)
was merged and then [reverted by #462](https://github.com/projectious-work/aibox/pull/462).
Its content is useful history, not an approved implementation baseline for
this new direction. Existing decision records still need selective
reconciliation after boundary approval; this draft does not silently change
their states.

## 3. Ownership boundary

| Responsibility | Owner and aibox's relationship |
|---|---|
| Workspace experience | **aibox**: curate and integrate tools, configurations, themes, previews, layouts, harness startup and attention signals; test the delivered experience. |
| Dev Container interpretation and lifecycle execution | **Dev Container CLI**, with its supported runtime tools. aibox consumes their public interfaces. |
| Image building, container execution and multi-service composition | Existing OCI/build/container/Compose tools. aibox provides standard artifacts and tested recipes. |
| Infrastructure provisioning and target preparation | **ainfra** or independently managed infrastructure. Optional explicit handoff; ainfra is not a prerequisite for local aibox use. |
| Agent execution loop, retries and task execution policy | The chosen harness or **Airunner**. aibox makes it runnable; it does not implement another agent loop. |
| Multi-agent organization and cross-product orchestration | **Kaits** or another caller. No mandatory Kaits → Airunner → aibox → ainfra chain. |
| Repository memory, skills, process entities and their migrations | **processkit**. aibox offers optional installation/integration, using product-owned interfaces rather than reproducing processkit internals. |
| Credentials, model access and external service permissions | User/operator and the relevant provider/tool. aibox preserves explicitly selected auth state and safe integration, not an identity broker or secret vault. |
| Project code, dependencies, tests and application-specific policy | The consuming repository. Installed tools do not make aibox the owner of the application's lockfiles or tests. |

Multiple harness panes remain supported in a human workspace. They must not
be advertised as separate security identities merely because they occupy
different panes. Independently authorized agents require explicitly separate
isolation and authority arrangements.

### Proposed shipped surface

- Versioned Dev Container Templates, reusable Features where needed, standard
  image/Dockerfile/Compose assets, and documented native extension points.
- A curated workspace configuration and asset collection: palettes, prompt
  presets, layouts, preview integrations, harness hooks, and bounded runtime
  helpers. Existing useful integration code can be retained or adapted.
- Tested dependency selections, provenance, compatibility information,
  migration assistance, diagnostics, documentation, and release evidence.
- A proposed thin MCP-first interface for bounded operations and progressive
  access to how-to guidance, with a convenience CLI sharing the operational
  core. Standard workspace definitions remain directly usable by upstream tools.

Do not introduce a parallel addon package manager, dependency solver, build
system, container API, generic provider framework, agent scheduler, process
store, or secrets service. New code needs a named user requirement and a
documented gap in suitable existing tools.

## 4. Standards-first composition

Use `devcontainer.json` for supported workspace settings, standard Features
for installable capabilities, and standard Templates for initial project
files. Dockerfile and Compose remain native extension surfaces. Prefer a
suitable maintained upstream Feature over maintaining another installer.
An aibox-owned Feature is justified for missing coverage or genuinely
aibox-specific integration, with its own tests and maintenance commitment.

The official CLI supplies build, up, exec, configuration inspection, Feature
and Template tooling, and Feature update/lockfile operations. Its current
README still marks stop/down as unimplemented. Pin and test a supported CLI
release before relying on particular behavior.
[CLI source](https://github.com/devcontainers/cli/blob/main/README.md)

Templates provide initial files; do not assume they solve subsequent
three-way updates of user-modified projects. Features provide a standard
installation unit; do not assume they reproduce every existing per-tool
option. These are explicit specification gaps to resolve.
[Template specification](https://github.com/devcontainers/spec/blob/main/docs/specs/devcontainer-templates.md),
[Feature specification](https://github.com/devcontainers/spec/blob/main/docs/specs/devcontainer-features.md)

| Capability | Reuse proposal | Residual aibox responsibility |
|---|---|---|
| Build/start/execute | Dev Container CLI and native container tooling | Tested inputs, compatibility, helpful diagnostics and documented workflows. |
| Stop/remove/multiple environments | Native runtime/Compose operations where supported | Exact environment identification, safe scope, preservation of user state; prove workflow equivalence. |
| Tools and versions | Upstream Features and package managers | Curated selection, unsupported-tool recipes, options and compatibility tests. |
| Terminal and file browsing | tmux, PowerKit, Yazi and existing preview executables | Cohesive configuration, key bindings, palette mappings and small missing adapters. |
| Themes and prompts | Native tool themes/configuration and Starship | Preserve the shared semantic palette and supported customizations. Evaluate reuse before replacing existing renderers. |
| Updates and migration | Upstream version/lock mechanisms and native file/config tools | Project-safe migration from v0; managed/user ownership rules and recovery. A generic synchronizer is not assumed. |
| Secrets and authentication | Native provider authentication and established secret tools | Explicit mounts/references and redaction; no credentials embedded in images or committed configuration. |
| Diagnostics and evidence | Native inspection, scanners and test tools | Cross-component checks specific to the shipped workspace and trustworthy reports. |

**Working assumption: all v0 `aibox.toml` desired settings can be expressed
through the Dev Container configuration model, with `devcontainer.json` as
the entry point. No successor `aibox.toml` is planned.** Prove this through a
field-by-field migration ledger; any gap returns for review rather than
silently becoming a second configuration system.

Use standard fields first, Feature references/options for installation and
tool selection, and a namespaced `customizations.aibox` extension only for
aibox-specific UX intent that lacks standard semantics. The standard permits
tool-specific customizations; it does not make those properties executable
by the generic CLI. An aibox content/runtime component must interpret them.
This is a proposed mapping, not a claim that tmux, themes or audio are native
Dev Container properties.
[Standard customization namespaces](https://github.com/devcontainers/spec/blob/main/docs/specs/supporting-tools.md)

For example, audio enablement can select client tooling and declare required
environment/mount configuration. Provisioning the host audio service remains
an independently authorized host task; setting a JSON property grants no
authority to perform it. Credentials and operator authorization policy are
not project settings to embed in `devcontainer.json`.

Native Dockerfile/Compose and tool configuration files remain valid referenced
assets. Users and agents edit these files directly. The specification must
define precedence between declared defaults and durable local overrides,
especially when a container-local theme/tmux edit differs from the project
default. No silent overwrite or competing authoritative setting is allowed.

An upstream replacement must be evaluated for feature fit, maintenance,
license, provenance, security, portability and dependency weight. Reuse is
not satisfied by adding a large new framework to avoid a small bounded
integration. Conversely, old aibox code need not be rewritten merely because
the distribution architecture changes.

## 5. Full v0 parity is a release obligation

Inventory baseline: aibox v0.35.0 release-line source,
`9061bf76d2a1f6ac3c7093d8c605bcc4de7eddf6`.
[Source snapshot](https://github.com/projectious-work/aibox/tree/9061bf76d2a1f6ac3c7093d8c605bcc4de7eddf6)

This is a first capability inventory, not a completed source audit or proof
of v1 compatibility. The specification must expand every row into observable
requirements, all supported options, source references and executable or
owner-observed acceptance evidence. Refresh the baseline for v0 features
shipped while v1 is being built.

Parity means the user can achieve the same supported outcome, including
configuration choices, accessibility and state preservation. Internal
implementation and command spelling may change through a documented migration;
loss of convenience or configurability is not automatically equivalent.
Do not label dropped features “delegated” without a usable integrated workflow.
Existing defects and unsafe implementation details need not be copied.

| ID | Capability to preserve | Acceptance focus |
|---|---|---|
| P01 | New/existing-project setup; minimal/harness-only and optional processkit setups | No required processkit, ainfra or orchestration installation; understandable selected components. |
| P02 | Build/rebuild/no-cache, start/attach, stop, removal and named environments | Correct target, repeatable operations, isolation between environments; removal preserves intended project/home data. |
| P03 | Tool selection, disabling, version selection and catalog discovery | All current recipes and supported tool-level choices accounted for; optional heavy tools stay optional. |
| P04 | Persistent home, harness login/configuration, caches and user files | Rebuild/recreation, upgrade and rollback do not silently lose or expose state. |
| P05 | Local overrides, custom Dockerfile/Compose extensions, mounts, ports, environment, paths and resource settings | Native editing remains usable; precedence and managed-file ownership are explicit. |
| P06 | tmux workspace, keyboard navigation, pane control and `ai`, `dev`, `focus`, `cowork` layouts | Multiple ordered harness panes, optional lazygit, live layout switching and existing useful shortcuts. |
| P07 | Theme family/variant/mode selection, live changes, emphasis, accessibility and gallery | Every existing palette and supported consumer; dark/light/auto, contrast, monochrome and no-color behavior. |
| P08 | Starship presets and user prompt customization | All eight named presets, aliases and font-light choices; durable customization path. |
| P09 | Yazi previews and navigation | Media, documents, Markdown, data, archives, source text; preview scrolling independent of file navigation. |
| P10 | Yazi preview controls and clipboard workflows | Persistent line numbers/wrapping; horizontal view; selectable read-only Vim preview; recursive directory-size hierarchy/totals; all copy actions through supported host clipboard path. |
| P11 | Standalone image/PDF/SVG/EPS viewing and watch workflows | Terminal/tmux graphics fallbacks and document refresh; actual renderers tested. |
| P12 | LaTeX build/watch and read-only browser preview companion | Multi-document configuration, live refresh, loopback serving and remote-forwarding guidance; compilation remains in the development container. |
| P13 | Harness launch/config integration and attention/title signals | Supported native/partial/manual fallbacks, model/effort identity where available, multi-pane priority, notifications and terminal-title restoration. |
| P14 | Status and logs | Resource/cgroup/OOM/process visibility, available MCP/runtime diagnostics, configured status segments and readable lnav workflow. |
| P15 | Audio/voice and host audio setup/diagnostics | Explicit opt-in output/input bridge, persistence and host-specific instructions; no hidden authorization expansion. |
| P16 | Processkit content and harness projections | Optional versioned integration, migrations and customization preservation delegated through supported processkit interfaces. |
| P17 | Inspection, integrity/health/security checks and actionable diagnostics | Useful human and existing machine-readable workflows; clear remediation and no false success. |
| P18 | Backup, reset, cleanup, configuration migration and recovery | Preview/diff where applicable, precise scope, auth/cache preservation choices, recoverable changes and explicit destructive consent. |
| P19 | Installation, updates, uninstall, help and shell completion | Equally accessible supported workflow even if no aibox binary ships; distinguish runtime, content and dependency updates. |
| P20 | Platform/runtime compatibility and headless use | Carry forward supported v0 workflows; explicit tested host/runtime/architecture matrix and documented limitations. |
| P21 | tmux header/status bars and customizations | Extended/plain/disabled modes; configurable elements, order, labels, separators, spacing, refresh and theme integration; preserve supported pane/window/session presentation. |
| P22 | Forward harness status into the outer terminal/tab header, including Ghostty | Preserve customizable title format, state symbols, project/repository/window/harness/model identity, multi-pane aggregation and title restoration. Test Ghostty and terminal-neutral fallback; no Ghostty dependency. |

P21 and P22 make existing requirements explicit; they are not substitutes for
P06/P13/P14. Title forwarding uses the established terminal output channel,
not a host API or permission to change host settings.

The v0 installation docs list Linux/macOS on x86_64 and ARM64. This does not
prove every runtime/host combination works. Establish the exact existing
support matrix from release evidence; do not silently narrow it or advertise
new combinations because a standard nominally permits them. Headless metadata
in v0 is not evidence of all the stronger PID-1/signal contracts proposed in
the old v1 spec; classify those separately as prospective requirements.

Early previews may openly report incomplete parity. A stable v1 replacement
must satisfy the completed parity ledger; any proposed exception comes back
to the owner as a scope change. Documentation-only claims do not close a row.

### Tool inventory to carry into the specification

The source contains 41 addon recipes. Recipe names are inventory identifiers,
not a proposed new v1 package format:

- AI: `ai-aider`, `ai-claude`, `ai-codex`, `ai-continue`, `ai-copilot`,
  `ai-gemini`, `ai-hermes`, `ai-mistral`, `ai-opencode`, `ai-tau`.
- Languages/quality: `go`, `go-quality`, `latex`, `node`, `python`, `rust`, `typst`.
- Documentation: `docs-docusaurus`, `docs-hugo`, `docs-mdbook`, `docs-mkdocs`,
  `docs-starlight`, `docs-zensical`.
- Tools: `audio-voice`, `browser-testing`, `cloud-aws`, `cloud-azure`,
  `cloud-gcp`, `cloudflare`, `data-preview`, `data-visualization`, `diagramming`,
  `git-ui`, `go-release`, `infrastructure`, `kubernetes`, `mermaid`,
  `preview-archive`, `preview-enhanced`, `release`, `supply-chain`.

Include the base toolchain as well: tmux, Yazi, Vim, Git, ripgrep, fd, bat,
eza, zoxide, fzf, delta, Starship, lnav and shell/network/archive utilities.
Preserve optional browser engines and accessibility testing, graphics
renderers, cloud/Kubernetes/IaC tooling, release/security tooling, and
per-tool selection rather than treating each bundle as an indivisible checkbox.

Representative preview coverage includes raster images, SVG, EPS, PDF,
Markdown, SQLite read-only inspection, CSV/TSV, Excel, video, archives and
syntax-highlighted text. Current `.excalidraw` text fallback must not be
misrepresented as a rendered drawing preview.

### Customization inventory to carry into the specification

The theme documentation lists 37 families and 76 variants. Families:
andromeeda, aurora-x, ayu, borland, catppuccin, contrast, contrast-mono,
dracula, everforest, github, gruvbox, houston, kanagawa, laserwave, material,
min, mono, monokai, moonlight, night-owl, nord, norton, one-dark, plastic,
poimandres, phosphor, projectious, red, rose-pine, slack, snazzy, solarized,
synthwave-84, tokyo-night, vesper, vitesse, vscode.

Preserve semantic palette mappings across tmux/PowerKit, Vim, Yazi,
Starship, LazyGit, bat/delta, fzf/eza, less/man, lnav and supported AI TUIs.
Exact palettes and best-available native modes are different support levels;
document them honestly. Preserve contrast and active-pane/selection visibility,
not just family names.

Prompt presets: `default`, `plain`, `minimal`, `nerd-font`, `pastel`,
`powerline-pastel`, `bracketed`, `arrow`, including the documented
`pastel-powerline` alias. Existing managed-output overwrite behavior is a
limitation to improve through clear ownership, not a reason to lose customization.

## 6. Should aibox have its own wrapper CLI?

| Option | Advantages | Costs and risks |
|---|---|---|
| A. No host wrapper: distribution plus native commands | Least lifecycle duplication; familiar upstream interfaces; smallest host software surface; workspace remains directly inspectable. | Users may need several tools for setup, stop/remove, updates and recovery; safe environment selection and v0 convenience need a credible solution. Documentation alone may be insufficient. |
| B. Thin MCP-first interface plus convenience CLI | Gives management agents bounded lifecycle operations and targeted guidance; humans get consistent basic commands using the same core. | Requires authority separation and reliable upstream delegation; risk of expanding into configuration management or a second lifecycle engine. |
| C. Mandatory orchestration/compiler CLI | Maximum freedom over unified UX and custom abstractions. | Highest duplication and maintenance; another configuration/state model; strongest risk of recreating the discarded inner platform. No demonstrated need in this draft. |

**Current preference: B, MCP-first with a shared operational core.** This
matches the owner's September 25 leaning. The CLI is a convenient human
projection, not the reason to build a new engine. Keep standard artifacts
usable without it. The command-level recommendations are in the
[current CLI audit](command-disposition.md); they remain review proposals.

### Operations and progressive disclosure

Expose bounded lifecycle and inspection tools to the external management
agent. Offer local inspection, validation and narrowly scoped refresh tools
inside the container, with no host capabilities. Human CLI operations use
the same validation, authority checks and delegation code as MCP operations.

Publish versioned, task-sized how-to resources from canonical documentation:
for example, “edit themes” identifies authoritative files, allowed settings,
examples, validation, live reload versus rebuild requirements, and deeper
links. Start with a small index and retrieve only the needed recipe; no
embedded LLM or bespoke retrieval platform is assumed. Documentation belongs
in MCP resources; any search/discovery tool returns bounded references rather
than introducing a second documentation source.

The agent reads guidance, edits ordinary config files using its existing
editing capabilities, then requests validation or an authorized operation.
Do not add `set theme`, `enable addon`, generic `edit_config`, or equivalent
configuration-mutating tools. Local refresh may consume edited configuration
and regenerate explicitly managed outputs, but must preserve user overrides.

A wrapper, if approved, must:

- Operate on standard artifacts that remain usable through the upstream CLI.
- Delegate through supported public interfaces; avoid importing private CLI
  internals or reconstructing its configuration resolution and build logic.
- Own only demonstrable aibox-specific gaps. No generic adapter SDK, new
  lifecycle state database, or parallel lockfile for upstream-owned data.
- Preserve stdout/stderr, exit status, cancellation and argument safety;
  report actual delegated outcomes rather than inventing success semantics.
- Identify resources exactly for destructive operations and leave unrelated
  containers, volumes, images, worktrees and personal state untouched.
- Remain removable from normal build/start/exec workflows. Migration tooling
  may be a separate bounded utility rather than a permanent wrapper feature.

Evaluate these cases before deciding: new project; existing repository;
tool enable/disable/version change; start and tmux attach; safe stop/remove;
named environments; live theme/layout change; update with user modifications;
v0 migration/rollback; audio setup and failed-build diagnosis.

Runtime helpers for theme switching, previews or attention signals are a
separate question from a host wrapper. Keep them bounded to the current
workspace and prefer native commands/configuration. No commitment to an
`aiboxctl` binary, daemon or generalized runtime-control protocol is made.

### New requested workspace capabilities (beyond v0 parity)

**N01 — optional left-side tmux harness overview.** A configurable, collapsible
pane shows every harness in the current tmux session: identity, working,
idle/stopped, waiting/question and error state; usage and context-window
consumption; and five-hour/weekly allowance usage and reset times where the
provider exposes reliable data. Distinguish per-session token use from
account-wide shared allowances; never sum the same quota once per pane.
Unknown, unsupported and stale values are explicit, with source and freshness.
Do not infer questions from idleness or fabricate quota percentages.

Reuse a suitable OSS viewer before writing one. Initial shortlist:

- [tmux-agent-status](https://github.com/samleeney/tmux-agent-status) is the
  closest layout/state candidate: persistent sidebar and lifecycle-hook
  integrations, plus custom-agent integration points. Complete quota and
  context coverage is unproven; assess compatibility with aibox/PowerKit
  before selecting it.
- [agent-dashboard](https://github.com/bjornjee/agent-dashboard) advertises
  token/cost and rate-limit displays, but also orchestration, workflow gates,
  agent control and remote UI. Do not adopt those broader responsibilities
  merely to obtain a monitor. Consider only if a bounded observational use
  can be established without competing hooks or workflow ownership.

These are documentation-based candidates checked September 25, not tested
selections. Evaluate licenses, maintenance, credential use, current-session
filtering, narrow-pane usability, all supported harnesses and native data
sources. A capability matrix and gap report precede any custom collector/UI.

**N02 — convenient diff and PR review.** Retain LazyGit and delta for local
working-tree/branch/commit review. Recommend an optional GitHub review
workflow using [gh-dash](https://github.com/dlvhdr/gh-dash), which uses `gh`
and delta and provides PR/issue navigation and actions. It complements local
Git review rather than making aibox a review platform. Native forge web UI
remains the fallback for review details not covered by the selected TUI.
[LazyGit](https://github.com/jesseduffield/lazygit) remains useful independently
of GitHub. Other forges require their own supported integration; no universal
forge adapter is proposed. Verify diff navigation, checkout, comments and
review submission before promising the complete workflow. Network writes
still require the user's authorization.

N01/N02 are requested scope additions, not shipped v0 capabilities. The
specification must assign acceptance criteria and release milestones; no
dependency selection or implementation is approved by this shortlist.

## 7. Security and authority boundary

The repository agent inside the container and the external agent managing
containers for the human are **distinct actors**. The accepted clarification
is recorded in `DEC-20260924_1541-SunnyTower-permit-container-local-customization-without-granting`.

| Actor | Permitted scope | Excluded authority |
|---|---|---|
| In-container repository agent | User-authorized edits to local tmux/themes/prompts/tool settings; local validation and reload; repository work; task guidance. | Start/stop/rebuild/remove the enclosing container, host configuration/access, host runtime sockets, host command bridges and management credentials. |
| External management agent | Explicitly authorized configuration and lifecycle operations on identified environments for the human. | Implicit authority over unrelated environments, credentials or destructive actions. |

Enforce the separation through runtime capabilities and accessible endpoints,
not prompt instructions or hidden tool names. An in-container request cannot
switch into operator mode. It may prepare a config change requiring a rebuild,
but application requires a separate externally authorized action. Before
executing it, the manager rechecks the actual inputs, including new mounts,
hooks and privileges; repository edits are not trusted authorization.

Normal workspace files may be backed by explicit host mounts. That does not
grant arbitrary host filesystem access. Terminal output for titles/clipboard
and preconfigured audio channels are bounded user integrations, not general
host control. Establishing or broadening those host integrations is forbidden
to the in-container agent.

Standard artifacts are executable inputs, not a security sandbox. Dev
Container configuration can request host-side initialization commands,
mounts and runtime privileges. Trust review must precede execution, and
published templates should make authority requirements visible.
[Metadata reference](https://github.com/devcontainers/spec/blob/main/docs/specs/devcontainerjson-reference.md)

- Host lifecycle operations belong to the human-controlled host context or
  an independently authorized operator. No host shell bridge, host runtime
  socket or privileged companion is exposed to the development agent to
  bypass that boundary.
- Project configuration cannot silently grant credentials, privileged mounts,
  publication authority or broader host access. User/operator policy remains
  separate from repository-controlled preferences.
- Authentication persistence uses explicit, scoped state. Do not import an
  entire host home or copy credentials into build layers, evidence or source.
- Audio and clipboard integration have explicit, limited purposes and opt-in
  or documented controls. Neither creates general host command authority.
- Container isolation alone is not a claim of safe execution of hostile code.
  Document limitations and any harness sandbox compatibility concessions.
- Provisioning hosts, operating secret brokers, confidential-computing trust
  infrastructure and granting agent permissions stay outside aibox's core.

Remote/Kubernetes support from the old rewrite must not be assumed to follow
from adopting Dev Container CLI. Preserve proven v0 support; treat additional
targets as separately reviewed increments with existing-tool feasibility
and evidence. Do not build a universal target framework in anticipation.

## 8. Company standards applied to this boundary

The proposed primary profile is infrastructure/template distribution, with
documentation and package/schema obligations. Apply CLI and operational MCP
requirements if those interfaces are actually shipped; profiles are selected
by behavior, not used to force unnecessary executables into existence.

| Standard | Consequence for the later specification |
|---|---|
| Application profiles ([source](https://github.com/projectious-work/internal/blob/d096a1992ab91dea2265cf4df342e7a97380b0e4/docs/standards/application-profiles.md)) | Explicit applicability matrix for content, runtime helpers, optional host CLI and host-gated release. |
| Application configuration ([source](https://github.com/projectious-work/internal/blob/d096a1992ab91dea2265cf4df342e7a97380b0e4/docs/standards/application-configuration.md)) | Single ownership per setting, precedence/provenance, redacted inspection, separate preferences and authority. |
| Output, logging and evidence ([source](https://github.com/projectious-work/internal/blob/d096a1992ab91dea2265cf4df342e7a97380b0e4/docs/standards/application-output-logging-and-evidence.md)) | Separate command results, diagnostics and durable evidence; redact at source and bind evidence to exact artifacts. |
| Compatibility and machine interfaces ([source](https://github.com/projectious-work/internal/blob/d096a1992ab91dea2265cf4df342e7a97380b0e4/docs/standards/compatibility-and-machine-interfaces.md)) | Versioned contracts, dependency windows, safe migrations/rollback and documented machine-interface behavior. |
| Roadmap and development evidence ([source](https://github.com/projectious-work/internal/blob/d096a1992ab91dea2265cf4df342e7a97380b0e4/docs/standards/product-roadmap-and-development-evidence.md)) | Machine-readable roadmap and evidence-backed state; distinguish proposed, implemented, tested and released. |
| Spec-driven development ([source](https://github.com/projectious-work/internal/blob/d096a1992ab91dea2265cf4df342e7a97380b0e4/docs/standards/spec-driven-development-cycle.md)) | Accepted baseline before implementation; requirement/evidence traceability and independent plan/conformance reviews. |
| Security and supply chain ([source](https://github.com/projectious-work/internal/blob/d096a1992ab91dea2265cf4df342e7a97380b0e4/docs/standards/security-and-software-supply-chain.md)) | Threat model, dependency evaluation/pinning, provenance, SBOM/scanning/signing and no secret leakage. |
| Verification and release engineering ([source](https://github.com/projectious-work/internal/blob/d096a1992ab91dea2265cf4df342e7a97380b0e4/docs/standards/software-verification-and-release-engineering.md)) | Test integration and actual user workflows; visual/structural coverage, exact-candidate gates and published-artifact verification. |
| Human-controlled host phases ([source](https://github.com/projectious-work/internal/blob/d096a1992ab91dea2265cf4df342e7a97380b0e4/docs/standards/human-controlled-host-phase-execution.md)) | Prepare inside; owner invokes bounded host-only validation/publication with immutable handoff and explicit authority. |
| Host-gated conformance ([source](https://github.com/projectious-work/internal/blob/d096a1992ab91dea2265cf4df342e7a97380b0e4/docs/standards/host-gated-release-conformance.md)) | Capability probes, repeatability, scoped cleanup and verifiable host evidence for supported runtimes. |
| Branching and promotion ([source](https://github.com/projectious-work/internal/blob/d096a1992ab91dea2265cf4df342e7a97380b0e4/docs/standards/git-branching-and-release-promotion.md)) | v1 development/promotion line, fast-forward promotion pointers and explicit owner merge approval. |
| Open-source documentation ([source](https://github.com/projectious-work/internal/blob/d096a1992ab91dea2265cf4df342e7a97380b0e4/docs/standards/open-source-documentation-strategy.md)) | Implementation-aligned README/docs/roadmap/changelog/releases; local build/deploy; no project-authored GitHub Actions. |
| AI-readable discovery ([source](https://github.com/projectious-work/internal/blob/d096a1992ab91dea2265cf4df342e7a97380b0e4/docs/standards/ai-agent-accessibility-and-generative-discovery.md)) | Version-aligned public machine-readable docs and read-only MCP resources; prefer shared delivery machinery. |

The agent-native interface direction is accepted in `CuriousSpire`, while its
[normative document is still in internal PR #12](https://github.com/projectious-work/internal/pull/12)
at review time. Reconcile that publication status during specification.
Where aibox owns operational capabilities, use an interface-neutral core,
bounded authority and capability-equivalent CLI/MCP adapters where applicable.
This does not require wrapping every upstream command in a new MCP server.
Read-only product documentation and operational control are separate surfaces.

Follow the [repository positioning policy](https://github.com/projectious-work/internal/blob/d096a1992ab91dea2265cf4df342e7a97380b0e4/docs/portfolio/repository-positioning-policy.md):
describe v1 as a proposal until evidence changes that status. Its older
Rust-CLI positioning describes v0, not a language mandate for v1.

## 9. Review decisions and next specification step

Requested review of this draft:

1. Accept or amend the proposed **curated workspace distribution** product
   definition and ownership boundary.
2. Confirm the proposed parity interpretation: preserve all supported user
   outcomes and choices, with documented migration rather than mandatory old
   internal architecture or command spellings.
3. Review the proposed MCP-first option B and command dispositions in the
   interface audit, including read-only doctor and config-file-only editing.
4. Review the proposal to separate current-feature parity from additional
   remote-target/headless guarantees inherited only from the discarded spec.
5. Confirm the configuration-mapping assumption and acceptance scope for
   the new sidebar and diff/PR-review workflows; candidate tools remain open.

After boundary approval, the specification work should produce:

- An exhaustive v0 feature/configuration/command ledger with source and
  release-evidence references, including all tool options and customization
  assets; an agreed moving-baseline policy until v1 replacement.
- A reuse assessment mapping every requirement to a standard tool, retained
  aibox asset, or justified gap; focused feasibility experiments only when
  authorized, with explicit findings rather than assumed compatibility.
- A tested target/dependency matrix; configuration and file-ownership model;
  state/auth persistence, update, migration and rollback contracts.
- A wrapper decision backed by the workflow comparison; any CLI/MCP contracts
  limited to capabilities aibox genuinely owns.
- Acceptance criteria and evidence plans for parity, security, visual UX,
  diagnostics, failure/recovery, host gates and releases; a phased roadmap.
- Selective decision reconciliation and independent review, followed by owner
  acceptance before implementation starts.

No full specification, product implementation, branch promotion or release
is authorized by this draft.

## Source guide

- [Company stack boundaries](https://github.com/projectious-work/internal/blob/d096a1992ab91dea2265cf4df342e7a97380b0e4/docs/architecture/projectious-stack.md).
- Indexed decisions cited by their unique memorable ID components in section 2;
  especially `NimbleLark`, `AbleDew`, `TenderFlame`, `HardySpring`,
  `AbleRaven`, `LuckyMoss` and the new owner-direction record above.
- Discussion `DISC-20260821_0224-DaringButter-what-normative-product-boundary-engine-template`
  and handover `LOG-20260824_1610-KeenMoss-session-handover` supply historical context.
- [v0 documentation source](https://github.com/projectious-work/aibox/tree/9061bf76d2a1f6ac3c7093d8c605bcc4de7eddf6/docs-site/content/docs):
  CLI reference, configuration/local overrides, base image/runtime operations,
  previews/audio, themes/prompts/layouts/attention titles, addon catalog and compatibility history.
- [v0 addon recipes](https://github.com/projectious-work/aibox/tree/9061bf76d2a1f6ac3c7093d8c605bcc4de7eddf6/addons)
  and runtime assets complement those docs. The specification must complete
  source-level tracing; this draft has not executed v0 or v1 runtime parity tests.
- Upstream links in section 4 and section 7 were checked on 2026-09-24;
  their moving branches must be replaced by pinned references for acceptance.
