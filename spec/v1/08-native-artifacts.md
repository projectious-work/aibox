# 8. Native project artifacts, Features and file ownership

This chapter specifies the concrete v1 project shape (R-CONFIG, R-OWNERSHIP,
R-INDEPENDENCE, R-REUSE). It is not a license to implement a new configuration
language. The project entry point is `.devcontainer/devcontainer.json` (JSONC
as accepted by the pinned Dev Container CLI); referenced native files retain
their own syntax and owner. aibox-specific intent is limited to the closed
`customizations.aibox` schema. The v0 `aibox.toml` is migration input only.
Maintained example definitions follow the chapter 2 commented-option
convention. Each aibox schema option appears in place with a purpose and
valid-values comment; inactive options remain commented. The list is generated
from the closed aibox schema, so a schema addition or changed enum requires
an updated explanation and regenerated examples in the same change. Native
examples cover only aibox-relevant fields, not every Dev Container CLI or
third-party extension property. The minimal starter keeps the entire aibox
extension commented out and remains usable without aibox; curated examples
uncomment only their selected settings. Comments never change the resolved
native definition.

## Starter distribution and upstream interoperability

Ship two versioned Dev Container Templates: `minimal` (Debian, shell, tmux,
Yazi, one selectable harness) and `curated` (same base plus documented optional
Feature choices and UX assets). A Template copies initial files; it does not
subsequently own or overwrite user edits. The normal user journey is
`aibox build`, `aibox up` and `aibox attach`. Separately, both templates MUST
pass direct upstream `devcontainer build`, `devcontainer up` and
`devcontainer exec` interoperability tests without the aibox host binary.
Product-specific local UX rendering may require
the aibox runtime Feature, but the container remains a valid Dev Container if
that Feature is removed; its special UX options then produce a doctor warning
instead of silently pretending to be active.

The initial project tree is:

```text
.devcontainer/
  devcontainer.json       # project-owned, sole Dev Container entry point
  Dockerfile              # optional project-owned native build input
  compose.yaml            # optional project-owned service definition
  compose.local.yaml      # optional ignored user-owned override
AGENTS.md                 # consuming repository or processkit-owned, never silently replaced
```

The Feature lockfile is owned by the Dev Container CLI; aibox MUST NOT introduce a
parallel Feature lock. The pinned CLI's build/up create it by default, and
release builds enforce it with the supported frozen-lockfile option. Existing
v0 `.aibox-home/` may remain during migration/rollback, but
new v1 projects use one documented persistent home mount. The default starter
MUST NOT create `.aibox/` or `.aibox-home/`. Do not copy or mount a whole host
home implicitly. Templates ignore any explicitly selected project-local bind
directory and private environment files. Operation records live in the process XDG
state directory; migration backups use an explicit destination.

The default image creates non-root user `aibox` with home `/home/aibox`.
The starter declares `remoteUser: "aibox"` and a project-scoped named volume:

```jsonc
"mounts": [
  "source=aibox-home-${devcontainerId},target=/home/aibox,type=volume"
]
```

V1-04 verifies the pinned CLI's volume identity and rebuild behavior. A user
may instead specify a native bind mount from an explicit host directory,
including an ignored project-local `.aibox-home`, through `mounts` or Compose
volumes. Changing `containerUser`/`remoteUser` and the Dockerfile user also
requires aligning the home and mount target with the actual container user's
home. `doctor` checks UID/GID, ownership and writability without silently
chowning host data. A mount may mask image-baked home contents; required
defaults belong in versioned Feature/image assets or idempotent first-use
initialization. Further cache/data persistence uses native mounts/Compose
volumes; aibox adds no volume schema or reconciler. Bind source paths resolve
on the container daemon host, which matters for remote runtimes.

For example, a user can replace—not add to—the named home volume with:

```jsonc
"mounts": [
  { "source": "${localWorkspaceFolder}/.aibox-home", "target": "/home/aibox", "type": "bind" }
]
```

For custom user `dev` whose Dockerfile home is `/home/dev`, the same native
definition instead uses `remoteUser: "dev"` and targets `/home/dev`. If
`containerUser` differs, its own process home and permissions are reviewed
separately; the two users need not be identical. The operator must ensure the
bind source exists on the daemon host and that `dev` can write it. A host bind has
different ownership, portability and cleanup behavior from a named volume.

`customizations.aibox` is not an upstream Feature or a generic customization
facility. It is product-specific data in `devcontainer.json`, read
only by software that elects to interpret it. Dev Container Features **can**
consume their own typed options and render assets packaged in the Feature at
image build. A Feature lifecycle hook may read files that a Dev Container
Template placed in the workspace after the workspace is available. Build-time
`install.sh` cannot assume access to those workspace files. Chapter 19 defines
115 legacy UX leaf paths under `workspace` as an optional runtime compatibility interface, not an infrastructure
override system. Native-only users may omit it. A Feature packages the local
interpreter and assets, without duplicating those preferences as build options.
Changing install/version Feature options requires a rebuild.
Never duplicate an authoritative value across Feature options and this
namespace. Moving a field requires schema, ledger, converter, example and
fixture updates together. The current `schemaVersion: "1"` closes the
specified keys; aliases and contradictory legacy forms are converter-only.
`harnesses.order` and `harnesses.launch.<name>.enabled`
express launch intent, not installation. `latex` and `diagnostics` contain only
aibox-owned UX/diagnostic preferences. New `workspace.sidebar` and
`workspace.review` are opt-in. Unknown versions/keys and null fail closed;
there is no permissive opaque extension bag.

## One owner for each setting

| Setting class | Authoritative input | Interpreter | Durable user override |
|---|---|---|---|
| Workspace name, image/build, user, ports, mounts, lifecycle | Standard `devcontainer.json`, native Dockerfile/Compose | Dev Container CLI and native tools | Edited native file; operator policy may deny unsafe requests |
| Install/disable/version of optional tool or harness | Native Feature/image/Dockerfile/package inputs | Native installer and upstream CLI | Project edit of native selection/build inputs |
| Theme family/mode/variant, prompt/layout, tmux/status/title UX | Versioned distribution defaults and optional `customizations.aibox.workspace`; no duplicate Feature preference options | Bounded local UX renderer; native-only operation also supported | Tool-native user files in persistent home |
| Native tool advanced settings | Native tmux, PowerKit, Yazi, Vim, Starship, LazyGit and harness files | Owning tool | User-owned native file; never copied back into namespace |
| Optional audio client | Audio Feature and approved native env/mount declarations | Package manager, runtime and client tools | Host audio service remains operator-owned |
| Processkit source/version/packages | Versioned processkit Feature input/lock and supported processkit interface | processkit installer | Its own documented user config/context |
| Host operation policy and credentials | Operator-local policy/secret store outside project root | Go policy adapter and native provider tool | Not repository-controlled |
| User authentication/home/cache | Scoped persistent home or native provider store | Harness/provider | Not committed or embedded in image |

Every `CFG:*` migration row must resolve to exactly one of these authorities,
with a concrete destination path/property, transform, default behavior and
fixture under chapter 19's conversion rules. The machine-readable field ledger
is inventory evidence; its generic targets are not binding implementation
instructions. Exact per-tool references and per-harness translations must be recorded in
the implementation manifest before implementing those adapters; unsupported
conversion is never silently reported as parity.

## Feature composition contract

Classify each retained v0 addon group by the usefulness of selecting its
members separately. A tightly coupled group is one aibox-owned Feature with
no subtool Feature routes or `all` bundle: `audio-voice` is the example, under
`ghcr.io/projectious-work/aibox-features/audio/voice`. A semantic collection
of independently useful tools has one Feature per selectable tool and a
versioned `all` bundle using standard `dependsOn`: `data-preview` is the
example, under `.../data/preview/<tool>` and `.../data/preview/all`.
Only publish an `all` bundle where its members also have individual routes;
avoid using a collection's metadata artifact as a Feature. Document each
group's classification and the member-to-Feature mapping in the selection
table. The 98 v0 rows remain individually traceable even when one Feature
installs several coupled components. Do not submit these collections to the
public Feature index yet. Record exact release references only after artifacts
exist. If a v0 project used per-tool switches within a now-atomic group, the
converter must report those states as an explicit compatibility decision;
it must not silently claim that the single Feature preserves them.

Each Feature owns a support matrix by selected tool version, container distro
and release, and architecture. Its `install.sh` selects an applicable canonical
method from the tool producer's installation documentation for that matrix
cell; different tool versions may need different methods. If the producer has
no applicable method, an official distro package is an allowed fallback only
after checking that the exact requested tool version is available from the
selected base's configured repository on that architecture, that its provenance
and dependencies are acceptable, and that installation and runtime tests pass.
An APT index entry for one release or architecture does not qualify all
Debian-based images. Pin and verify consumed release or package inputs; the
Feature's own version does not pin the installed tool version. Reuse a
qualified upstream Feature, native package manager or official distribution
*inside* an aibox Feature where it meets that contract. A Feature is a packaged
installer with declared options, dependencies and entry point; it is not an
aibox package-manager API. `customizations.aibox` carries UX preferences only
and must not install tools. An unsupported version/base/architecture combination
fails clearly at build time before installation; Arch and Nix are not inferred
as supported from a generic Linux claim. Follow
DEC-20261002_0054-AmberCompass-qualify-feature-compatibility-per-tool-version.
Feature compatibility is distinct from any client editor's own container
requirements; a passing Feature build does not establish VS Code Server or
extension compatibility on that base.

The owner still decides which v0 rows to retain or remove. The grouped Feature
direction supersedes DEC-20261001_1954-SmoothTrout's restriction on routine
aibox-owned installer Features; the group boundary is fixed by
DEC-20261002_0046-SmoothLotus-package-coupled-addon-groups-atomically-and. Native
image/Dockerfile/package units are implementation inputs for those Features,
not a parallel project-facing addon resolver.

Every one of the 98 `ADDON:*` tool entries has these migration obligations:

1. Record a selected qualified v1 native install unit/reference and its
   publisher, immutable inputs, license, maintenance/security evidence, or an
   explicit owner-approved removal/scope decision. Retain every census row.
2. Map v0 absent/default-enabled/explicit-enabled/explicit-disabled/version
   states to native selection/options; explicit disable MUST result in absence
   where v0 promises it, including from the base image. Record any reduced
   support explicitly for owner decision.
3. Map v0 recipe `requires` into native dependencies or a documented
   prerequisite; prove there is no hidden addon solver or unexpected extras.
4. Test default/non-default supported versions and disable on both container
   architectures, for every declared distro/release cell. Test download/checksum,
   missing exact package version and wrong-architecture failures where those
   operations are performed by the selected installation mechanism.
5. Prove that native project lockfiles remain authoritative for application
   dependencies; bundled convenience tools never rewrite them.

### Tool Feature interface

[`catalog/v1-05-features.json`](catalog/v1-05-features.json), validated by
[`feature-catalog.schema.json`](feature-catalog.schema.json), is the frozen
interface catalog for retained and added tools. It lists the reference base,
supported and unsupported targets, error codes, every group with its shape,
every Feature with its ID, owner, options, dependencies, environment, tools,
installation provenance and observable checks, every removed census row,
every compatibility note and the smoke tests. Prose in this section defines
how those fields are realized; where prose and catalog differ, the catalog
wins and the difference is a contract defect. `notes` fields are informative
research records, not requirements.

**Identity.** An aibox-owned Feature with catalog ID `<parent>/<name>` is
published as `<registry>/<parent>/<name>:<featureVersion>` with the moving
tags `<major>` and `<major>.<minor>`. Its `devcontainer-feature.json` has
`id` = `<name>` and `version` = `featureVersion`. Sources live in
`features/src/<catalog-id>/` with `devcontainer-feature.json`, `install.sh` and
any helper files; nothing outside that directory is packaged. Upstream
Features are referenced at the catalog's `upstream.ref` digest with the
catalog's `upstream.options`; aibox does not republish them.

**Shapes.** A `single` group has one tool Feature. An `atomic` group has one
Feature that installs all members and has no member or bundle Feature. A
`semantic` group has one Feature per member and a bundle Feature
`<group-path>/all` whose `dependsOn` lists every member at its default
options; the bundle installs nothing itself and has no options.

**Options.** A Feature whose catalog entry has `options.version` declares a
string `version` option with exactly that `enum` and `default`. Any other value
fails with `AIBOX_FEATURE_UNSUPPORTED_VERSION` before installation. The option
is the tool's own version; a package revision such as an APT `-0` suffix is an
installation detail. A Feature without `options.version` has no `version`
option. No option accepts `latest`, a range or a URL, and no Feature has an
`enabled` option. The specification checks reject an example that sets an
option key its Feature does not declare.

**Debian packages.** A tool with method `debian-package` installs the
candidate version of each named package from the base's configured
repositories with `--no-install-recommends`; the base's snapshot pin, not the
Feature, fixes that version. Such a Feature does not add or change APT
sources. The `version` of a
`package` check is the version expected on the reference base. V1-05 changes
the minimal and curated starters to configure `trixie`, `trixie-updates` and
`trixie-security` from one snapshot.

**Dependencies.** Each `dependsOn` entry becomes a standard `dependsOn`
reference to `<registry>/<id>:<major>` with an empty option object; the Dev
Container CLI installs it automatically. Each `requires` entry, used only for
multi-version runtimes, becomes an `installsAfter` reference to
`<registry>/<id>`; the Feature checks before installation that the runtime is
installed and otherwise fails with `AIBOX_FEATURE_DEPENDENCY_MISSING`, naming
the missing Feature ID in the detail. It uses whichever enum version of the
runtime is installed. A `versionCoupling` of rule `equal` requires the coupled
Feature's installed version to equal this Feature's `version`; `table` uses
the listed mapping; a violation fails with `AIBOX_FEATURE_COUPLING_MISMATCH`.
When a tool is installed twice with different versions, for example through a
bundle and an explicit member selection with another version, the second
installation fails with `AIBOX_FEATURE_VERSION_CONFLICT`; an identical repeat
succeeds without change.

**Targets.** Before any change to the image, `install.sh` reads
`/etc/os-release` and `dpkg --print-architecture`. A combination not listed in
the Feature's `targets` fails with `AIBOX_FEATURE_UNSUPPORTED_TARGET`; the
catalog's `unsupportedTargets` are the fixtures for this behavior. Upstream
Features have no aibox error contract; for them, V1-05 verifies only presence,
absence, isolation, bundles and smoke tests on the supported targets.

**Errors.** Every failure listed in the catalog's `errors` writes one line
`aibox-feature error <CODE> feature=<catalog-id> <detail>` to standard error
and exits non-zero, so the Dev Container build fails. A test passes when the
build fails and some line of its combined output contains
`aibox-feature error <CODE> feature=<catalog-id>`, regardless of any prefix the
build adds. For `AIBOX_FEATURE_UNSUPPORTED_TARGET` the catalog ID may be the
Feature under test or any Feature it transitively `dependsOn` or `requires`,
because the first Feature in install order fails; for every other code it is
the Feature under test. `beforeInstall` errors leave no package, file or user change
behind. Errors marked `component` are verified by implementation-track
component tests because black-box tests cannot inject the fault.

**Reference base and absence.** Acceptance tests build the catalog `baseline`
fixture and add Features to its `features` object. A command or package in
`baseline.provides` exists in the base and is excluded from absence checks.
With no Feature selected, no other catalog `commands` entry resolves through
`command -v`, no other `package` check is installed and no `path` check
exists. For a semantic member selected alone, together with the Features it
`requires`, the commands, packages and paths of every sibling that is neither
a `dependsOn` or `requires` dependency nor listed in an installed tool's
`alsoProvides` for that target are absent. Disabling a tool means omitting its Feature.

**Checks.** Every check is evaluated in the running container, as the
fixture's `remoteUser` through `devcontainer exec`, and again as `root`
through `docker exec --user root <containerId>`. Every entry in `commands`
resolves through `command -v`. Regular expressions are ECMAScript, searched
anywhere in the combined standard output and error, with `^` and `$`
matching at line boundaries.

| Kind | Passes when |
|---|---|
| `version` | `argv` exits 0 and output matches `regex` after `{version}` is replaced by the regular-expression escaped selected version |
| `command` | `argv` exits 0 and output matches `regex`, which has no placeholder |
| `package` | `dpkg-query -W -f '${Version}' <package>` exits 0 and prints exactly `version` |
| `texlive` | `tlmgr info --only-installed --data name <package>` prints `package` |
| `path` | `test -d <path>` succeeds |

The selected version of a `version` check is the Feature's `version` option
value, otherwise the tool's `version`. `versionFrom` names the Feature whose
version the tool follows and does not change how its checks run. A tool's
`packages` lists what it installs and is informative; absence is judged only
through checks and commands.

**Side effects.** Owned Features declare no lifecycle commands, no mounts, no
`privileged` mode, capabilities or devices; `privileges` is explanatory text
only. Their `containerEnv` is exactly the catalog's `containerEnv`. They
write nothing under the workspace folder, so the workspace Git status is
unchanged by a build. A tool that would update itself is configured not to.
Tools never replace project lockfiles.

**Smoke tests.** For each `smoke` entry, the fixture file or directory is
copied under its base name into the root of the workspace folder, `argv` runs
there as the remote user, exits 0 and creates `produces` relative to the
workspace folder. When `producesMagicHex` is given, the file's first bytes
equal those hexadecimal bytes.

**Publishing and running tests.** `scripts/v1-05-features.sh publish
--registry <host>/<prefix>` publishes every aibox-owned Feature to
`<host>/<prefix>/<catalog-id>` with tags `<featureVersion>`, `<major>` and
`<major>.<minor>`, rewrites aibox `dependsOn` and `installsAfter` references
to that prefix, leaves upstream references unchanged, exits 0 on success and
non-zero with a message on failure. Its prerequisites are the pinned Dev
Container CLI and Docker. Acceptance tests run an OCI Distribution registry on
`localhost:5000`, which the pinned CLI contacts over plain HTTP, publish to
`localhost:5000/aibox-features`, and select aibox Features by their `<major>`
tag so automatic dependencies resolve to the same instance. Upstream Features
are selected by digest. A test run covers the host's native architecture;
acceptance requires one run on each supported target. Presence checks may
share one build across Features that do not conflict; absence, isolation,
bundle and error checks use exactly the selection their item defines.
Production publication and the public Feature index are out of V1-05 scope.

**Inspection.** `aibox inspect` and the MCP `inspect_workspace` tool report a
configuration entry with key `devcontainer.features`, value the sorted list of
selected Feature references, `layer` `project`, `source`
`.devcontainer/devcontainer.json`, `domain` `native`, `redacted` false and
`invocationOnly` false. Feature options are never echoed. The entry is absent
when no Feature is selected.

### Census traceability

The selection table is a source file keyed by `ADDON:<recipe>/<tool>`;
its fields are `v0-source`, `v0-default`, `v0-versions`, `v1-ref`, `v1-option`,
`install-owner`, `dependency`, `license`, `digest`, `platforms`, `test-id`, and
`decision`. The current `ledger/addons.json` is only its input. Phase V1-05
requires the complete table, owner dispositions for removals, group cohesion
classification, grouped Feature sources and references for retained tools,
producer installation provenance, and qualification evidence. The 98 census
rows remain individually traceable; semantic bundles are additional convenience
selections.

The base image must separately enumerate its packages and runtime glue from
`ledger/base-build.json`; tools outside the addon catalog do not disappear.
Build inputs and packages are pinned by digest/version, with update owner and
provenance. Do not force a heavy optional tool into the base merely because an
optional Feature cannot remove it later.

## Managed output and refresh

The renderer accepts an immutable, schema-validated aibox UX intent plus
versioned distribution assets and user-owned native overrides. It generates
only declared managed output files under the container-local home/XDG config
root. Each managed file has a recorded source digest and a header or manifest
claiming aibox ownership. It never writes `.devcontainer/devcontainer.json`,
Dockerfile, Compose, `AGENTS.md`, processkit context, user overrides or provider
auth. A path with a symlink, special file, wrong owner or unexpected content is
rejected rather than overwritten. The sole allowed managed generation-pointer
symlink and exact output layout are specified in chapter 19.

Refresh validates all inputs, renders a private complete generation, checks
paths/permissions, and atomically selects it through the managed pointer in
chapter 19. Retain the previous generation until reload succeeds; rollback
restores file selection and reports failed reload recovery. This does not
promise atomic rollback of tool-process side effects.
The result reports changed files and any `rebuild_required` setting separately.
It must not run `postCreateCommand`, install packages, regenerate host files,
restart the container or recreate a tmux session implicitly. A live theme
change may reload tmux/PowerKit and selected tools; a setting requiring image
rebuild remains pending for the external operator.

## Native configuration example

The maintained [minimal](examples/minimal/.devcontainer/devcontainer.json) and
[customized](examples/customized/.devcontainer/devcontainer.json) JSONC files
show the complete commented aibox option catalog and keep only their selected
native fields active. aibox preferences are configuration intent and require
the qualified runtime Feature for activation. An example must not invent a
Feature registry reference; a published reference includes its tested ID,
digest and actual options in the example and selection table together.

## Artifact acceptance

- Validate minimal, multi-harness, disabled-tool, version override, audio,
  LaTeX, local override and Compose examples with the upstream schema and
  aibox extension schema; reject unknown Feature options and aibox keys.
- Diff generated files against explicit ownership declarations; user files
  survive repeated refresh, rebuild and version upgrade byte-for-byte.
- Compare direct Dev Container CLI and wrapped journeys from a clean checkout;
  the wrapper may add diagnostics but not change native workspace meaning.
- Check each maintained example's active JSONC, exhaustive commented aibox
  option catalog and documentation against the closed customization schema.
  Test the aibox integration with the pinned CLI, not the upstream engine's
  entire option matrix or unrelated Features.
- Secret fixtures and private overrides remain outside Git, image layers,
  machine output, published docs and release artifacts.
