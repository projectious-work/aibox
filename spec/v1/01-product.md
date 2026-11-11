# 1. Product and architecture

## Product contract

**R-PRODUCT:** aibox is a curated, reproducible AI development workspace
distribution built on the existing Dev Container CLI. Its value is the
integrated tools, terminal UX, customization, safe maintenance, compatibility
and guidance. It remains usable without processkit, ainfra, Airunner or Kaits.

**R-PARITY:** preserve every supported user-facing v0.35.0 outcome and choice,
including optional capabilities. Preserve useful content/integration code
where possible. A different command spelling or implementation is allowed;
unannounced feature loss is not. Documented broken/unsafe details are corrected,
not reproduced. New N01/N02 requirements are additional scope, not alleged v0
features. Stable v1 replacement is gated on parity and their accepted criteria.

**R-CORE:** implement one Go application core with thin CLI and MCP adapters.
Delegate Dev Container processing/build/up/exec semantics to the public CLI;
use native runtime/Compose operations only for verified gaps. No independent
Dev Container compiler, scheduler, builder, dependency solver, desired-state
database, general plugin SDK, secret broker or embedded AI planner.

Component responsibilities (package names are illustrative):

| Component | Responsibility | Forbidden expansion |
|---|---|---|
| Go operation core | Typed inputs/results, authority checks, exact resource identity, delegation, cancellation and diagnostics | Reimplement upstream lifecycle/config merge semantics |
| CLI adapter | Human rendering, machine output, interactive attachment | Independent operation logic/config mutation |
| MCP adapter | Bounded tools, schemas, stdio protocol using official Go SDK | Arbitrary host shell or invented JSON-RPC implementation |
| Knowledge resources | Versioned task recipes, catalog, schemas, migration docs | Internal process context exposure or second doc source |
| Workspace content | Features/templates, palettes, layouts, hooks and preview integration | Host-management credentials or privileged bridge |
| Recovery/migration utility | Reviewed v0 conversion, backup/diff and rollback assistance | A permanent second configuration engine |

**R-BOUNDARY:** ainfra prepares infrastructure; aibox consumes an existing
authorized target, optionally through a product-owned result. Harness/Airunner
owns the agent loop; Kaits owns cross-product orchestration; processkit owns
repository process memory, entities and migrations. aibox can integrate these
products without reproducing their implementation or requiring their future
unreleased versions. If processkit lacks a supported integration interface,
that is an explicit integration gap, not permission to copy its entity engine.

## Two actors, enforced authority

**R-AUTHORITY:** the repository agent runs inside the dev-container. It may
edit local tmux/themes/prompts/tool configuration and invoke container-local
validation/refresh for its user. The external management agent acts separately
on the human's behalf and may invoke authorized lifecycle operations on
identified containers. In-container code never receives host runtime sockets,
host credentials, an operator endpoint or a general host command bridge.

Scope is enforced by process placement, OS permissions, mounted resources and
accessible endpoints. Tool hiding and an `operator` flag alone are insufficient.
A shared repository file is untrusted input to the manager: review current
hooks, mounts, requested privilege, paths and executable sources before acting.
Editing a file does not approve executing it. Bind approval to the actual
resolved inputs; revalidate if those inputs change before execution.

Normal repository/home mounts and preconfigured audio/clipboard/title channels
are purpose-limited integrations. They grant no access to arbitrary host
files, host shell, service setup, or container lifecycle. Do not claim the
container is a hardened hostile-code sandbox merely because it is a container.

**R-LOCAL:** local configuration that requires rebuild produces an actionable
`rebuild_required` result. The local agent cannot perform or automatically
trigger that rebuild. Local tmux reload is permitted; destructive session
recreation requires separate explicit consent and must not be confused with
container restart. Multiple harness panes do not imply separate identities.

## Quality targets

**R-REUSE:** justify each aibox-owned component by a named uncovered user
requirement and a reuse assessment. Prefer native files and maintained tools.
Do not rewrite a working preview helper merely to make every asset Go.

**R-OBSERVABILITY:** expose truthful state, provenance and limitations. Missing
measurements are unknown, not zero. Error handling names the failed step,
affected resources, any completed effects and safe next action. Logs and
evidence redact secrets before persistence; normal stdout remains machine-safe.

**R-INDEPENDENCE:** a standard workspace can be built and started directly
with the supported upstream CLI without an aibox host daemon. aibox-specific
UX code may run through installed content/hooks; customizations are not
automatically implemented by upstream tools. Removal of the convenience
interface must not destroy the workspace definition or stored user data.
