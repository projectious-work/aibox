# V1-05 native tool catalog

V1-05 is in progress on `feat/v1-05-native-tool-catalog`. The phase maps the
frozen v0 tool ledger to project-owned Dev Container installation choices.
Qualified upstream Features and native image or Dockerfile packages are used
directly. `customizations.aibox` remains optional workspace UX metadata and
does not install tools. An aibox-owned installer Feature requires a specific
gap justification and owner review.

## Current work

- The frozen ledger has 98 rows and 41 recipes. Source and OCI review found
  four upstream Feature mappings, 68 gaps, and 26 community candidates needing
  further qualification. All 98 frozen source hashes match; 38 current v0
  checkout entries differ from that frozen baseline. The durable audit and OCI
  evidence are under `spec/v1/evidence/V1-05/`.
- `spec/v1/ledger/v1-05-tool-selection.json` accounts for each row with its
  current classification, source, native proposal, and decision state.
  `spec/v1/scripts/validate-v1-05-catalog.py` checks full census coverage and
  the active example against the reviewed mapping.
- The customized example now selects the upstream Hugo Feature at an immutable
  OCI manifest digest with the frozen `0.165.0` extended version. Its source,
  options and intended architectures are recorded in
  `spec/v1/fixtures/v1-05-native-features.json`. Minimal remains feature-free.
- Go CLI and local MCP inspection now report the selected Feature references
  without exposing their options. The paired inspection test covers the
  customized example. Source-level validation passes; live Hugo installation
  and absence checks remain pending.
- `scripts/verify-v1-05-host.sh` is the Python-only host gate for the active
  Hugo selection. It checks enabled and disabled containers with the pinned
  Dev Container CLI and records source, architecture, step logs and scoped
  cleanup. It must be run by the host operator on both Linux architectures.
- `dev-notes/V1-05-gap-review.md` lists all 94 unresolved rows, each proposed
  native route, conditional upstream alternatives, and two suggested removal
  candidates. Owner decisions are pending. No gap installer or removal has been
  implemented.

## Acceptance still pending

V1-05 remains in progress until owner dispositions for the unresolved rows,
chosen install units, two-architecture installation checks, and final v1 docs
validation pass. The four upstream mappings have source and OCI evidence only;
their install/version/disable behavior has not yet passed live tests. No v1
Feature artifact or release has been published.
