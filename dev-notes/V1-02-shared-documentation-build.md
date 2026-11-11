# V1-02 — shared documentation build

- **Roadmap item:** [V1-02](../spec/v1/roadmap.yaml), shared Hugo foundation.
- **Scope and status:** Merged documentation infrastructure with independent
  conformance and release acceptance pending; roadmap status is `in_progress`.
- **Reader:** documentation maintainers, release operators and reviewers.

## Behavior, decisions and boundaries

PR [#464](https://github.com/projectious-work/aibox/pull/464) moved the v1
preview to the pinned `brand-theme-hugo-vanilla` Hugo module. PR
[#465](https://github.com/projectious-work/aibox/pull/465) aligned v0 release
navigation. Both lines pin theme v0.3.4. The v0 root remains the current
product documentation; v1 is visibly a preview. Old Docsy source is retained
only as excluded migration history rather than forming a second live theme.

The accepted specification chose one shared build shell over maintaining
Docsy and the brand theme in parallel. V1-02 does not publish alpha/beta
snapshots; that is V1-19, with exact candidate and release checks. A local
build is not proof that published pages match a release artifact.

## Validation and consequences

[Documentation foundation evidence](../spec/v1/docs-foundation-evidence.md)
records the merge commits and local builds. On 2026-09-27, Hugo Extended
v0.165.0 built v1 (31 pages, 142 static files) and v0 (185 pages, 138 static
files). Candidate-specific alpha/beta archival and deployment checks remain
V1-19 work; local builds do not claim those gates passed.

Operational consequence: both lines require the pinned Go/Hugo/Node module
toolchain for local builds. No v1 CLI feature is implied by the preview site.
Documentation changed: the active [v1 site README](../docs-site/README.md),
[documentation contract](../spec/v1/15-documentation.md), navigation and
historical Docsy instructions. Follow-up publication and accessibility
verification remain owned by V1-19.
