# aibox v1 — Software implementation specification

Version: **0.8, 2026-09-30**. Development line: `v1.x-dev`.

aibox combines native Dev Container artifacts with a Go operation core,
an MCP-first interface and a thin human CLI. Configuration is edited in files;
the application does not expose a configuration-mutation API.

The numbered chapters, schemas, feature ledger and roadmap define the
implementation contract. The [command disposition](command-disposition.md)
maps the v0 CLI to v1 operations and delegated workflows.

## Read in this order

1. [Product boundary and requirements](01-product.md)
2. [Configuration, persistence and migration](02-configuration.md)
3. [CLI, MCP and authority contracts](03-interfaces.md)
4. [Reuse assessment, target and dependency matrices](04-reuse-targets.md)
5. [Feature ledger and acceptance criteria](05-acceptance.md)
6. [Requirements trace and technical qualification](06-delivery.md)
7. [Go architecture and package boundaries](07-architecture.md)
8. [Native artifacts, Features and ownership](08-native-artifacts.md)
9. [Lifecycle algorithms and failure recovery](09-lifecycle.md)
10. [CLI, MCP, doctor and progressive guidance](10-interfaces-and-guidance.md)
11. [Workspace UX and full v0 parity](11-workspace-experience.md)
12. [Migration, named environments and persistence](12-migration-and-state.md)
13. [Security and trust boundary](13-security.md)
14. [Verification, target qualification and release](14-verification-and-release.md)
15. [Version-line documentation and prerelease publication](15-documentation.md)
16. [Layered CLI/MCP configuration and operational logging](16-configuration-and-logging.md)
17. [Secret and credential transfer into the Dev Container](17-secrets-and-credential-transfer.md)
18. [Executable operation, policy, recovery and quality contracts](18-operational-contracts.md)
19. [Configuration realization, installer and private-state rules](19-configuration-realization.md)
20. [Human and AI personas, user stories and solution paths](20-personas-and-user-stories.md)

The [roadmap](roadmap.yaml) is the canonical phase graph: twenty-three
content-specific phases in six capability groups. Each phase links the
chapters governing its implementation. Display order follows the intended
implementation sequence, with dependencies pointing only backward; IDs remain
stable references. The roadmap tracks delivery separately from these target
requirements. `implementationStatus: done` records locally verified completion;
only release-backed work with published conformance evidence becomes `shipped`.
From V1-03,
the roadmap's `demos` map requires a feature-specific run of the Go executable
and an observable outcome for every phase. The executable and matching Hugo
walkthrough remain demonstrable after each integrated increment; V1-03 must
deliver the first safe read-only binary and stdio MCP tool with equivalent
inspection results. Every executable roadmap phase proves CLI/MCP parity for
its operations; V1-09 expands the guarded registry and client coverage.
Verified checkpoints may publish scoped alpha/beta releases before full v0
parity. V1-19 hardens the cross-line archive/navigation matrix rather than
postponing the first matching prerelease documentation.

The [phase ID crosswalk](phase-id-crosswalk.md) maps identifiers used before
this sequential renumbering to their current IDs.

Implementation evidence and remaining gaps belong in the
[development logbook](../../dev-notes/README.md), not in the target contract.

### Detailed v0 evidence

Baseline: v0.35.0 release-line commit
`9061bf76d2a1f6ac3c7093d8c605bcc4de7eddf6`.

- [Configuration ledger](ledger/configuration.md), [machine-readable rows](ledger/configuration.json),
  [type/enum/alias declarations](ledger/configuration-types.json),
  [default-function and Default-implementation evidence](ledger/defaults.json),
  [closed aibox customization schema](customization.schema.json).
- [Tool/dependency ledger](ledger/addons.md), [full recipe metadata](ledger/addons.json).
- [Complete CLI declarations](ledger/commands.json), [atomic action rows](ledger/command-actions.json)
  and [argument rows](ledger/command-arguments.json),
  including argument annotations, aliases and resource enums; pair with the
  command disposition table rather than assuming every flag survives.
- [Runtime asset census](ledger/runtime-assets.json), [documentation coverage](ledger/documentation.json),
  [environment identifier census](ledger/environment.json), [counts](ledger/census.json).
- [Base image build/entrypoint declarations](ledger/base-build.json) retain
  package lists, tool pins and installation behavior outside the addon catalog.
- [Behavioral feature source map](feature-trace.md) links each F row to primary
  v0 implementation and its v1 implementation component.
- [Versioned operation result schema](operation-result.schema.json) is the
  shared CLI/MCP envelope; [request schema](operation-request.schema.json)
  closes common fields and operation-required selectors. Actual client checks
  and generated per-tool SDK conformance remain G04.
- [Process settings schema](process-settings.schema.json) and
  [log event schema](log-event.schema.json) define the separate bounded
  configuration and diagnostic contracts; chapter 16 fixes precedence,
  authority, sinks and failure handling.

The census is exhaustive for the declared source scopes and includes custom
deserializer exceptions described in chapter 2. It is not a claim that all
v0 behavior is inferred from declarations, that every environment identifier
is public, or that source inventory proves v1 behavior. The behavioral ledger supplies
the user-facing contract; unresolved mappings block the affected conversion.
Default implementations and validation logic remain authoritative in the
linked immutable source; migration tests must check effective values, not
merely parse the type inventory.

Reproduce the census from repository root (Git history containing the baseline
is required; no product runtime or credentials are needed):

```sh
node spec/v1/scripts/inventory.mjs --check
node spec/v1/scripts/contracts.mjs --check
node spec/v1/scripts/render-devcontainer-examples.mjs --check
node spec/v1/scripts/validate.mjs
uv run --script spec/v1/scripts/validate-contracts.py
```

Omit `--check` only when intentionally regenerating inventory evidence.
The Node/Python scripts are specification maintenance tooling, not an exception to
the Go product-language decision. They do not run hooks or container commands.

## Normative convention

`R-*` requirements and `AC-*` acceptance groups are normative contracts.
A ledger row inherits its referenced acceptance group and the all-features rule.
`candidate`, `unverified` and `blocked` are not support claims. Verify each
requirement against implemented behavior. Before replacing a later v0 release,
extend the ledger and parity tests to cover its changes.
