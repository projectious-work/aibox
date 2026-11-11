# Specification layout and lifecycle standard

## Purpose and scope

This standard defines where specifications live, what one specification
version contains, and how it is frozen and changed. It applies to every
projectious.work repository that implements against specifications. The
[spec-driven development cycle](spec-driven-development-cycle.md) defines how
work proceeds against them.

Terms in *italics* are [Semantic Anchors](https://github.com/llm-coding/Semantic-Anchors/tree/6554641a11fdf99fdbdf09cc0602268e3db88966); see
[README](README.md#conventions).

## What gets a specification

- Every **feature** gets a specification. It is the feature's contract.
- A **fix** cites the specification version that already states the correct
  behavior. It adds a new version only when the specification was wrong or
  silent about the case.
- A **roadmap phase** has no specification of its own. It references the
  specifications of its features (`specs` in `roadmap.yaml`, see the
  [roadmap standard](product-roadmap-and-development-evidence.md)).
- Product-level descriptions (vision, architecture, glossary) are
  documentation, not contracts. They live in the repository's documentation
  and may be cited from specifications.

## Names

| Artefact | Pattern | Example |
|---|---|---|
| Directory | `spec/<slug>/v<N>/` | `spec/payments-api/v2/` |
| Branch | `spec/[<issue>-]<slug>-v<N>` | `spec/123-payments-api-v2` |
| Tag | `spec/<slug>/v<N>` | `spec/payments-api/v2` |
| Reference | `<slug>@v<N>` | `payments-api@v2` |

- `<slug>` follows the branch-name rules of the
  [git standard](git-branching-and-releases.md#branch-names), must not end in
  `-v<digits>`, is unique in the repository, and never changes or is reused.
- Ticket numbers appear only in branch names and in `SPEC.toml`; tickets
  change, a frozen directory name cannot.
- `N` starts at 1 and increases by 1. There are no minor or patch versions.

## Contents of a version

Every version directory is self-sufficient: reading the latest frozen version
alone gives everything needed to implement and test the feature.

```text
spec/<slug>/v<N>/
  SPEC.toml         required  slug, version, status, supersedes, roadmap,
                              issues, decisions
  requirements.md   required  normative *EARS* requirements with stable IDs
  interfaces.*      required  machine-readable interface catalog (TOML, YAML,
                              or JSON), validated with the specification
  journeys.md       required  acceptance journeys citing requirement IDs
  overview.md       required  non-normative context, rationale, boundaries
  CHANGES.md        N > 1     delta to v<N-1>
  FREEZE.sha256     frozen    hash manifest of every other file
  schemas/          optional  schemas referenced by the interface catalog
  fixtures/         optional  test data and golden output
  examples/         optional  maintained, runnable examples
  diagrams/         optional  diagram sources; rendered images only if needed
  design/           optional  user-facing alternatives and mockups
  artifacts.toml    optional  large binaries by reference
```

`SPEC.toml`:

```toml
slug = "payments-api"
version = 2
status = "frozen"            # draft | frozen
supersedes = "payments-api@v1"   # omitted for v1
roadmap = ["RM-31"]
issues = ["123"]
decisions = ["ADR-7"]
```

- **Requirement IDs** (for example `PAY-012`) are unique per slug and never
  reused across versions; a removed requirement keeps its ID retired.
- **Interface catalog:** everything the feature makes observable, such as
  operations, artifacts, options and defaults, result fields, error codes,
  supported targets, and the observable check for each effect, including
  absence when something is disabled.
- **`CHANGES.md`** lists added, changed, and retired requirement IDs, every
  file that differs from `v<N-1>`, and one line of rationale. Every file it
  does not list is byte-identical to the previous version.
- **`FREEZE.sha256`** uses `sha256sum` format, one line per file in the
  directory except itself, sorted by path.
- **Large binaries** (over 1 MB) are never copied into a version. They are
  referenced from `artifacts.toml` with URI, SHA-256, size, and media type.
- Code, test implementations, and evidence ledgers are not part of a
  specification; they live in the repository's source, test, and evidence
  locations.

## Lifecycle

1. **Draft.** Create `spec/<slug>/v<N>/` with `status = "draft"` on branch
   `spec/[<issue>-]<slug>-v<N>`. A draft is edited freely, and only on its
   `spec/` branch.
2. **Review.** The specification must be **test-ready**: an author who reads
   only this version can derive every acceptance test without reading code and
   without inventing an identifier, option, default, result field, error code,
   or supported target. The review follows the
   [spec-driven cycle](spec-driven-development-cycle.md#4-review-the-plan-against-the-complete-baseline);
   on its small-change path, the owner's request replaces the review, but the
   version must still be test-ready.
3. **Freeze.** After owner approval, the last commit on the branch sets
   `status = "frozen"` and adds `FREEZE.sha256`. The branch is squash-merged
   into `main` with the title `docs(spec): freeze <slug>@v<N>`, and the squash
   commit is tagged `spec/<slug>/v<N>` by the agent or person who merged it.
   The frozen version is the contract.
4. **Change.** A frozen version is never edited, moved, or deleted. Every
   change, including a typo, creates `v<N+1>` as a full copy with its
   `CHANGES.md`, through its own `spec/` branch and the same review and
   freeze. All versions stay in the tree.
5. **Implement.** Implementation pull requests (`feat/`, `fix/`, `chore/`)
   carry one `Spec: <slug>@v<N>` line per specification, naming a frozen
   version. They never change files under `spec/`.

### Work in flight when a new version freezes

The integration owner reads `CHANGES.md` for every open branch that cites the
previous version. If no requirement the branch implements changed, the
`Spec:` line moves to the new version without further review. Otherwise the
affected work stops and the
[baseline-change rule](spec-driven-development-cycle.md#7-control-discoveries-and-baseline-changes)
applies. Released work keeps citing the version it was released against.

## Reading rule for agents

Agents read the latest frozen version of a specification, or the draft they
are working on. They read older versions only to compare, for example with
`git diff --no-index spec/<slug>/v1 spec/<slug>/v2`.

## Validation

`python3 .standards/standards.py check-specs --base main` (in the standards
repository: `python3 scripts/standards.py check-specs --base main`) runs
locally before every merge, and without `--base` on `main` after tagging. It
fails on:

- directory, slug, or version-number violations, missing required files, or
  an invalid `SPEC.toml`;
- a frozen version whose files do not match `FREEZE.sha256`, or whose tag is
  missing or differs from its content (with `--base`, only for versions
  already frozen there, since the tag follows the merge);
- a draft or frozen status that contradicts the presence of `FREEZE.sha256`;
- `CHANGES.md` missing, or not listing exactly the files that differ from the
  previous version;
- an interface catalog that is not `interfaces.toml`, `.yaml`, `.yml`, or
  `.json`;
- with `--base <ref>`: any change to a frozen version, or any deletion of a
  version directory, compared with `<ref>`.

## Agent briefing

> Every feature has a specification; fixes cite it and add a version only when
> it was wrong. Versions live in self-sufficient `spec/<slug>/vN/` directories
> (`SPEC.toml`, *EARS* requirements, interface catalog, journeys). Freeze =
> owner approval, `FREEZE.sha256`, squash merge, tag `spec/<slug>/vN`. Frozen
> versions are never edited or deleted; any change is `vN+1` via its own `spec/`
> branch. Read only the latest frozen version; cite it as `Spec: <slug>@vN`.
