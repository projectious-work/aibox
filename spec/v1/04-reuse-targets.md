# 4. Reuse assessment, targets and dependencies

All target and dependency selections below are **proposed and unverified for
v1** unless marked accepted direction. Source reading is not compatibility,
license/supply-chain clearance or execution evidence.

## Reuse decisions

| Capability | Existing implementation to prefer | Residual aibox work / acceptance |
|---|---|---|
| Dev Container interpretation, build/up/exec, Features/Templates | [Dev Container CLI](https://github.com/devcontainers/cli) — accepted foundation | Public-command delegation; AC-BUILD/AC-LIFE; no internal API import or compiler. |
| Stop/remove and Compose services | Runtime CLI / native Compose | Verify exact identity and differences between stop/remove; no daemon API recreation. |
| MCP protocol | [Official Go SDK](https://github.com/modelcontextprotocol/go-sdk) | Typed bounded tools/resources and stdio; AC-CLI/AC-GUIDE. No custom protocol engine. |
| Basic language/tool installs | [Maintainer Features](https://github.com/devcontainers/features), native package managers | Evaluate exact options/platform fit; fill proven gaps with small Features. |
| Workspace terminal | tmux, PowerKit, Yazi, Vim, Starship, lnav, native plugins | Retain assets/hooks/palettes; AC-UX/AC-PREVIEW. |
| Source diff | Git, delta, LazyGit | Curated configuration and convenient launch; AC-REVIEW. |
| PR workflow | Optional [gh-dash](https://github.com/dlvhdr/gh-dash) plus gh | GitHub-specific review integration; native forge web fallback; not a new forge abstraction. |
| Harness sidebar | Evaluate [tmux-agent-status](https://github.com/samleeney/tmux-agent-status) first | Session scoping, theme/layout integration, hook compatibility and telemetry gaps; AC-SIDEBAR. |
| Rich monitoring alternative | [agent-dashboard](https://github.com/bjornjee/agent-dashboard) | Not selected: orchestration/workflow/remote control exceed monitoring scope; evaluate only a separable observational mode. |
| Preview/render | Existing Poppler, resvg/librsvg, Ghostscript, Rich, ffmpeg, csvkit, SQLite, bat, chafa/timg/mutool/entr | Preserve small glue where native plugins do not cover existing behavior. |
| Processkit integration | Product-owned supported installer/config/migration APIs | G02: API availability and pinned compatibility need agreement; do not copy process logic. |
| Config schema/JSONC | Upstream schema and maintained Go parser/validator | G01/G04: select dependencies after conformance test; no new JSONC parser. |
| Image supply chain | Existing SBOM/scanning/signing tools | Exact-artifact evidence, local release workflow; AC-RELEASE. |

### Tool-family disposition

Every one of the 98 recipe tool entries in [addons.md](ledger/addons.md) has
its source default/pin/option declaration. All retain version and explicit
disablement where supported. This matrix assigns a reuse route; it does not
assert a matching Feature exists for every tool.

| Recipes | Proposed route | Gaps to prove |
|---|---|---|
| python, node, go, rust | Upstream language Features first; native managers for companion tools | uv/pip/poetry/pdm, pnpm/yarn/bun, clippy/rustfmt/audit/cross target controls; no hidden compulsory extras. |
| latex, typst | Suitable Feature or small package-backed aibox Feature | All TeX collections, SVG conversion, architecture coverage and optionality. |
| go-quality, go-release, release, supply-chain | Existing upstream release binaries/packages with verified provenance | Per-tool version/disable controls, pinned checksums, language-neutral composition. |
| 10 AI recipes | Upstream harness installers or trustworthy Features | Eleven harness identities include host-only Cursor; install differs from integration; preserve auth/config/hook contracts. |
| Six docs recipes | Existing framework packages/build tools | Respect consuming repository lockfile; installed convenience tooling does not own app dependencies. |
| git-ui | gh and LazyGit native packages | Optional installation, existing binary absence semantics, new gh-dash selection separate. |
| preview-archive, preview-enhanced, data-preview | Native packages plus retained preview glue | Pixel/text rendering, read-only data access, all keybindings and watch flows. |
| diagramming, data-visualization, mermaid | D2/Graphviz, Vega/Vega-Lite, Mermaid/Puppeteer | Compatible browser/runtime pins; opt-in heavy dependencies, no implicit docs dependency. |
| browser-testing | Playwright + axe and matched browser revisions | Chromium default, Firefox/WebKit optional; repo package lock remains authoritative. |
| infrastructure, kubernetes, cloud-* and cloudflare | Existing upstream CLI packages/Features | Pinning, checksums and amd64/arm64; installing a tool grants no credentials/host authority. |
| audio-voice | PulseAudio client/SoX/ALSA packages | Native host service readiness separate; no in-container host setup. |

Base-image packages and runtime assets are separately enumerated in
[runtime-assets.json](ledger/runtime-assets.json). Do not consider a tool
covered merely because it is absent from the addon catalog. For each selected
dependency record upstream owner, license, version/digest, provenance source,
platform assets, known vulnerabilities, maintenance evidence and update owner.
Unverified dependencies are blocked, not silently accepted as “standard”.

## Target matrix

**R-TARGETS:** preserve the established v0 platform/user workflows; qualify
each advertised combination with evidence. The installation docs list four
host binary targets. Runtime mentions and host-gate evidence differ; this
draft does not transform mentions into a fully supported Cartesian product.

| ID | Host / architecture | Runtime route | v0 evidence / proposed v1 obligation |
|---|---|---|---|
| T01 | macOS arm64 | OrbStack Docker-compatible runtime + Compose | Documented host use and release-host workflow; primary qualification target. |
| T02 | macOS arm64 | Docker Desktop + Compose | Documented prerequisite/audio route; qualify lifecycle, mounts, UID, audio and previews. |
| T03 | macOS amd64 | Compatible Docker runtime + Compose | Published binary target; exact OS/runtime support versions must be evidenced before support claim. |
| T04 | Linux amd64 | Docker Engine + Compose | Published binary/install route; native ownership/networking and lifecycle gates required. |
| T05 | Linux arm64 | Docker Engine + Compose | Published binary target; native image/tool architecture gates required. |
| T06 | Linux/macOS as applicable | Podman with compatible Compose provider | v0 runtime-specific behavior exists; qualification required, not disposable optional parity. Resolve exact supported cells in G05. |
| T07 | Debian-based Linux amd64/arm64 container | Local Go MCP/runtime helpers, no host runtime credentials | Mandatory local customization/doctor/preview/harness tests on both architectures. |
| T08 | SSH access to an existing prepared machine | Operator executes supported native tooling there | Evaluate existing v0 scope; no new universal remote provisioning API implied. |
| T09 | Kubernetes, new cloud backends, Windows native, Alpine | No initial new support claim | Older rewrite ideas are not v0 parity. Any addition requires explicit scope and verification. |

Version floors/ceilings for OS/runtime/Compose/CLI remain G05. Before baseline
acceptance, publish an explicit matrix of tested versions and evidence, with
owner-approved treatment of any formerly supported unavailable target. No
runtime compatibility claim follows solely from Docker-API compatibility.
Nested rootless Podman as an optional **workspace tool** is distinct from
Podman as the **host lifecycle runtime**; neither grants host socket access.

## Dependency and compatibility matrix

**R-DEPENDENCIES:** select an exact tested dependency set before implementation
baseline/release gates; ranges alone are not reproducibility. Commit lockfiles,
image/Feature digests and checksums; retain supported range separately.

| Dependency | Role / where needed | Selection state / constraints |
|---|---|---|
| Go toolchain | Build CLI/MCP core for darwin/linux amd64/arm64 | Language accepted; supported stable version chosen with SDK minimum and recorded build provenance. No Go toolchain needed by end user of binary. |
| Official Go MCP SDK | Go module | Preferred; pin tested tagged release in go.mod/go.sum, protocol/client compatibility explicit. |
| Dev Container CLI | External operator dependency | Foundation accepted; pin tested release and its runtime requirements; no fallback to private internals. |
| CLI's Node/runtime distribution | Operator-side CLI prerequisite | Package/install method determines runtime; do not assume Go binary makes the upstream CLI dependency disappear. |
| Docker/Compose or qualified Podman/provider | Operator lifecycle/build | Explicit per-target versions; validate before side effects; never install silently. |
| OCI base image + package inputs | Workspace content | Debian baseline, two native architectures; immutable digest and package manifest. |
| Features/Templates | Standard content distribution | Explicit references/digests; published artifacts and native CLI use tested independently. |
| tmux/PowerKit/Yazi/Vim/Starship/etc. | Workspace UX | Source pins/assets in ledger; retain behavior and test actual version combinations. |
| AI harness binaries | Selected optional workspace content | Existing recipe defaults are historical evidence, not a perpetual pin; native hook/API capability matrix required. |
| processkit | Optional integration | v0 docs identify v0.28.8 for v0.35.0; v1 compatibility/interface selection is G02, no unconditional dependency on an unreleased rewrite. |
| Sidebar and PR TUI | New optional UX | Candidates only; G06/G07 qualify maintenance/license, privileges, telemetry, rendering and integration. |
| Scanners/signers/release tools | Maintainer workflow | Existing toolset; exact candidate-bound execution evidence, no mandatory user runtime cost. |

Research anchors checked 2026-09-25 (inspection snapshots, **not release pins**):

- Dev Container CLI: `5dc7533314b5ba7ec3875c30143dfe1aec644870`.
- Official Go SDK: `e07f0c9d5abf509ac1e47abf27cfa539eeda64a5`.
- tmux-agent-status: `546b6ca51c75b415db0b3ce06910703f875d3aae`.
- gh-dash: `1b14dd961ea47e2ad53eae5fed7fc942219bcc01`.

The official CLI currently distinguishes implemented build/up/exec from
unimplemented stop/down in its README. Qualification must inspect the pinned
release; the aibox core may fill this bounded gap through native tooling.
[CLI reference source](https://github.com/devcontainers/cli/blob/5dc7533314b5ba7ec3875c30143dfe1aec644870/README.md)

The SDK provides Go server/client primitives, including stdio. Choose the
protocol version with actual clients rather than assuming the latest protocol
is universally supported.
[SDK reference source](https://github.com/modelcontextprotocol/go-sdk/blob/e07f0c9d5abf509ac1e47abf27cfa539eeda64a5/README.md)
