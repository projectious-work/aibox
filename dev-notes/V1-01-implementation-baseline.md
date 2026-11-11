# V1-01 — implementation baseline

- **Roadmap item:** [V1-01](../spec/v1/roadmap.yaml), implementation baseline.
- **Scope and status:** Merged baseline work with independent conformance and
  release acceptance still pending; roadmap status is `in_progress`.
- **Reader:** v1 implementers and independent conformance reviewers.
- **Accepted baseline:** [v1 specification](../spec/v1/README.md) after PR #463;
  v0 parity snapshot `9061bf76d2a1f6ac3c7093d8c605bcc4de7eddf6`.

## Behavior and boundaries

PR [#464](https://github.com/projectious-work/aibox/pull/464) added the
reproducible v0 configuration, addon, command, asset, environment and
documentation inventories, schema/contract checks, and native configuration
examples. This is evidence for planning conversion; it does not prove a v1
binary, runtime parity, or a supported target.

The baseline retained source links, hashes and exact counts rather than
extrapolating behavior from configuration schemas. The prior v1 Docsy build
is historical pre-migration evidence, not the current documentation stack.
These choices avoid presenting unimplemented behavior as available.

## Validation and consequences

The source commit, commands, counts and limitations are in
[baseline evidence](../spec/v1/baseline-evidence.md). `node
spec/v1/scripts/baseline.mjs` checks the ledger against the frozen v0 source;
the specification and contract validators check internal consistency. No
container lifecycle acceptance test or published-v1 verification is claimed.

Security and compatibility consequence: the inventory is read-only and does
not import v0 host privileges into v1. Conversion behavior and target support
still require the later V1-20 migration and V1-22 target-qualification phases.

Documentation changed: the [specification index](../spec/v1/README.md),
ledgers, examples and baseline evidence serve maintainers and reviewers; no
user-facing v1 operation is documented as usable. V1-01 needs independent
conformance review and release acceptance before `shipped` can be claimed.
Follow-up runtime parity proof is owned by V1-04 through V1-22 in the roadmap.
