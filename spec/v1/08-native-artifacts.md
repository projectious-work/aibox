# 8. Native project artifacts, Features and file ownership

This chapter specifies the concrete v1 project shape (R-CONFIG, R-OWNERSHIP,
R-INDEPENDENCE, R-REUSE). It is not a license to implement a new configuration
language. The project entry point is `.devcontainer/devcontainer.json` (JSONC
as accepted by the pinned Dev Container CLI); referenced native files retain
their own syntax and owner. aibox-specific intent is limited to the closed
`customizations.aibox` schema. The v0 `aibox.toml` is migration input only.

## Starter distribution and direct-upstream journey

Ship two versioned Dev Container Templates: `minimal` (Debian, shell, tmux,
Yazi, one selectable harness) and `curated` (same base plus documented optional
Feature choices and UX assets). A Template copies initial files; it does not
subsequently own or overwrite user edits. Both templates MUST pass a direct
upstream `devcontainer build`, `devcontainer up` and `devcontainer exec` test
without the aibox host binary. Product-specific local UX rendering may require
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
.aibox/
  overrides/              # ignored user-owned native tmux/theme/Yazi/Starship overrides
  state/                  # ignored local receipts and migration journals, not auth content
  backups/                # ignored, explicit migration/recovery backups
AGENTS.md                 # consuming repository or processkit-owned, never silently replaced
```

The Feature lockfile is owned by the Dev Container CLI; aibox MUST NOT introduce a
parallel Feature lock. The pinned CLI's build/up create it by default, and
release builds enforce it with the supported frozen-lockfile option. Existing
v0 `.aibox-home/` may remain during migration/rollback, but
new v1 projects use one documented persistent home mount. Do not copy a whole
host home. Templates include `.gitignore` entries for local overrides,
private environment files, receipts, backups and persistent auth state.

The `customizations.aibox` schema in this specification is a proposed public
contract, not a replacement for the upstream schema. `schemaVersion: "1"`
closes accepted keys. The 115 legacy UX leaf paths map under `workspace` with
identical semantic names; aliases and contradictory legacy forms are handled
only by the converter. `harnesses.order` and `harnesses.launch.<name>.enabled`
express launch intent, not installation. `latex` and `diagnostics` contain only
aibox-owned UX/diagnostic preferences. New `workspace.sidebar` and
`workspace.review` are opt-in. Unknown versions/keys and null fail closed;
there is no permissive opaque extension bag.

## One owner for each setting

| Setting class | Authoritative input | Interpreter | Durable user override |
|---|---|---|---|
| Workspace name, image/build, user, ports, mounts, lifecycle | Standard `devcontainer.json`, native Dockerfile/Compose | Dev Container CLI and native tools | Edited native file; operator policy may deny unsafe requests |
| Install/disable/version of optional tool or harness | Feature presence and documented Feature options | Feature installer and upstream CLI | Project edit of Feature reference/options |
| Theme family/mode/variant, prompt/layout, tmux/status/title UX | `customizations.aibox.workspace` | Bounded local aibox UX renderer | Explicit native file in ignored `.aibox/overrides/` |
| Native tool advanced settings | Native tmux, PowerKit, Yazi, Vim, Starship, LazyGit and harness files | Owning tool | User-owned native file; never copied back into namespace |
| Optional audio client | Audio Feature and approved native env/mount declarations | Package manager, runtime and client tools | Host audio service remains operator-owned |
| Processkit source/version/packages | Versioned processkit Feature input/lock and supported processkit interface | processkit installer | Its own documented user config/context |
| Host operation policy and credentials | Operator-local policy/secret store outside project root | Go policy adapter and native provider tool | Not repository-controlled |
| User authentication/home/cache | Scoped persistent home or native provider store | Harness/provider | Not committed or embedded in image |

Every `CFG:*` migration row must resolve to exactly one of these authorities,
with a concrete destination path/property, transform, default behavior and
fixture. A row that says merely “Feature options” or “native integration” is
not implementation-ready. This is a specification completion gate, not work to
silently improvise during implementation. The machine-readable field ledger
is inventory evidence; its generic targets are not binding implementation
instructions.

## Feature composition contract

Use the Dev Container Feature mechanism for installable capabilities. A
Feature is a packaged installer with declared options, dependencies and
entry point; it is not an aibox package-manager API. Standard or
maintainer-owned Features are preferred when they meet the complete tool
contract. An aibox-owned Feature is permitted only for a proven gap or
aibox-specific integration, and must reuse existing package managers/tool
distribution rather than implement another resolver. A thin aibox Feature MAY
compose several upstream Features when conditional v0 bundle behavior cannot
be expressed otherwise; it must not hide an unbounded addon registry.

Every one of the 98 `ADDON:*` tool entries has these migration obligations:

1. Record one selected v1 install unit/reference and its publisher, immutable
   digest, license, current maintenance evidence and security review.
2. Map v0 absent/default-enabled/explicit-enabled/explicit-disabled/version
   states to Feature presence/options; explicit disable MUST result in absence
   of the tool where v0 promises it, including from the base image.
3. Map v0 recipe `requires` into explicit Feature dependencies or a documented
   prerequisite; prove there is no hidden addon solver or unexpected extras.
4. Test default and non-default supported versions on both container
   architectures, including a failed download/checksum and wrong architecture.
5. Prove that native project lockfiles remain authoritative for application
   dependencies; bundled convenience tools never rewrite them.

The selection table must be a reviewed source file keyed by `ADDON:<recipe>/<tool>`;
its fields are `v0-source`, `v0-default`, `v0-versions`, `v1-ref`, `v1-option`,
`install-owner`, `dependency`, `license`, `digest`, `platforms`, `test-id`, and
`decision`. The current `ledger/addons.json` is only its input. Phase V1-04
cannot be accepted without the filled table and actual Feature manifests.
This explicit design gate avoids the false certainty of assigning all 98 tools
to guessed public Features today.

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
rejected rather than overwritten.

Refresh is transactional at the file-set level as far as the filesystem
supports it: validate all inputs; render to a private staging directory;
check output paths and permissions; atomically replace managed files; retain
the previous generation until tmux/tool reload succeeds; roll back on failure.
The result reports changed files and any `rebuild_required` setting separately.
It must not run `postCreateCommand`, install packages, regenerate host files,
restart the container or recreate a tmux session implicitly. A live theme
change may reload tmux/PowerKit and selected tools; a setting requiring image
rebuild remains pending for the external operator.

## Native configuration example

This example is schematic until selected Feature refs/digests and pins are
recorded in the V1-04 table. It demonstrates ownership and syntax, not an
assertion that a registry artifact already exists:

```jsonc
{
  "name": "my-project",
  "build": { "dockerfile": "Dockerfile" },
  "remoteUser": "aibox",
  "features": {
    "ghcr.io/devcontainers/features/go:1": {}
  },
  "customizations": {
    "aibox": {
      "schemaVersion": "1",
      "workspace": { "theme": "gruvbox", "mode": "dark", "layout": "dev" },
      "harnesses": { "order": ["codex"], "launch": { "codex": { "enabled": true } } }
    }
  }
}
```

The Feature reference above is illustrative, not a tested release pin. Before
this is promoted to a maintained example, replace it
with the selected, digest-bound release and prove direct-upstream execution.

## Artifact acceptance

- Validate minimal, multi-harness, disabled-tool, version override, audio,
  LaTeX, local override and Compose examples with the upstream schema and
  aibox extension schema; reject unknown Feature options and aibox keys.
- Diff generated files against explicit ownership declarations; user files
  survive repeated refresh, rebuild and version upgrade byte-for-byte.
- Compare direct Dev Container CLI and wrapped journeys from a clean checkout;
  the wrapper may add diagnostics but not change native workspace meaning.
- Secret fixtures and private overrides remain outside Git, image layers,
  machine output, published docs and release artifacts.
