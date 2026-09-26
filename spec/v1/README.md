# aibox v1 — Dev Container CLI-based product specification

Version: **0.2 implementation draft, 2026-09-26**. Review target: `v1.x-dev`.
Product boundary accepted; this specification is **not yet accepted**.
No product implementation, release or merge is authorized by this draft.

## Accepted inputs

The owner confirmed all five boundary review points, selected a thin wrapper
with **MCP first**, and selected **Go** for the application core and adapters.
The CLI is a human-facing projection of the same operations. Ordinary config
files are edited by humans/agents; no configuration mutation API is planned.

Decision provenance in the company coordination repository:

- `DEC-20260924_1508-PluckySage-rebase-aibox-v1-on-dev-container`: Dev Container
  CLI foundation, all v0 user-facing capability parity and reuse-first.
- `DEC-20260924_1541-SunnyTower-permit-container-local-customization-without-granting`:
  local customization allowed, host/lifecycle authority forbidden inside.
- `DEC-20260925_0635-SprySpring-refine-aibox-boundary-around-file-based`:
  configuration hypothesis and new sidebar/review workflows.
- `DEC-20260925_0652-HopefulTower-approve-aibox-v1-product-boundary-and`:
  boundary review points 1–5 confirmed, Go and thin MCP-first wrapper selected.

The [boundary review input](boundary.md) and [command audit](command-disposition.md)
are preserved for context. Their earlier tentative language does not reopen
these decisions. The numbered chapters below form the proposed implementation
contract.
The discarded specification from PR #455/#462 and existing old v1 code are
not the implementation baseline for this rewrite.

## Read in this order

1. [Product boundary and requirements](01-product.md)
2. [Configuration, persistence and migration](02-configuration.md)
3. [CLI, MCP and authority contracts](03-interfaces.md)
4. [Reuse assessment, target and dependency matrices](04-reuse-targets.md)
5. [Feature ledger and acceptance criteria](05-acceptance.md)
6. [Standards, trace and review gates](06-delivery.md)
7. [Go architecture and package boundaries](07-architecture.md)
8. [Native artifacts, Features and ownership](08-native-artifacts.md)
9. [Lifecycle algorithms and failure recovery](09-lifecycle.md)
10. [CLI, MCP, doctor and progressive guidance](10-interfaces-and-guidance.md)
11. [Workspace UX and full v0 parity](11-workspace-experience.md)
12. [Migration, named environments and persistence](12-migration-and-state.md)
13. [Security and trust boundary](13-security.md)
14. [Verification, target qualification and release](14-verification-and-release.md)

The [roadmap](roadmap.yaml) is the canonical phase graph: nineteen
content-specific phases in five capability groups. Each phase links the
chapters governing its implementation. `planned` is not a shipped claim.

### Detailed v0 evidence

Baseline: v0.35.0 release-line commit
`9061bf76d2a1f6ac3c7093d8c605bcc4de7eddf6` (latest published release checked
2026-09-25). The PR base is the separate `v1.x-dev` commit
`d9a8e22b4ba4385f1dc0134f58f5405417116d85`.

- [Configuration ledger](ledger/configuration.md), [machine-readable rows](ledger/configuration.json),
  [type/enum/alias declarations](ledger/configuration-types.json),
  [default-function and Default-implementation evidence](ledger/defaults.json),
  [proposed closed aibox customization schema](customization.schema.json).
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
  v0 implementation and the proposed v1 owner.
- [Versioned operation result schema](operation-result.schema.json) is the
  shared CLI/MCP envelope proposal; [request schema](operation-request.schema.json)
  closes common fields and operation-required selectors. Actual client checks
  and generated per-tool SDK conformance remain G04.

The census is exhaustive for the declared source scopes and includes custom
deserializer exceptions described in chapter 2. It is not a claim that all
v0 behavior is inferred from declarations, that every environment identifier
is public, or that proposed v1 mappings work. The behavioral ledger supplies
the user-facing contract; unresolved mappings are explicit review blockers.
Default implementations and validation logic remain authoritative in the
linked immutable source; migration tests must check effective values, not
merely parse the type inventory.

Reproduce the census from repository root (Git history containing the baseline
is required; no product runtime or credentials are needed):

```sh
node spec/v1/scripts/inventory.mjs --check
node spec/v1/scripts/validate.mjs
```

Omit `--check` only when intentionally regenerating inventory evidence.
The Node scripts are specification maintenance tooling, not an exception to
the Go product-language decision. They do not run hooks or container commands.

## Normative convention

`R-*` requirements and `AC-*` acceptance groups are proposed normative
contracts. A ledger row inherits its referenced acceptance group and the
all-features rule. `candidate`, `unverified` and `blocked` are not support
claims. No requirement is satisfied by this documentation PR alone.
Any later v0 release triggers a ledger delta review before v1 replacement;
retirement/deferment requires owner approval, not a rewritten inventory.
