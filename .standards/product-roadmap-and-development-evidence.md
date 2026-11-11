# Product roadmap and development evidence standard

## Purpose and scope

This standard defines how projectious.work software projects express future
direction and connect completed and released outcomes to implementation evidence.
It applies once a project discusses future work beyond the current release.

Terms in *italics* are [Semantic Anchors](https://github.com/llm-coding/Semantic-Anchors/tree/6554641a11fdf99fdbdf09cc0602268e3db88966); see
[README](README.md#conventions).

The canonical roadmap MUST be machine-readable YAML in `roadmap.yaml`, or at a
documented project-specific path. A rendered website or Markdown roadmap MAY be
generated from it, but MUST NOT become a competing source of truth.

## Principles

- Roadmaps communicate outcomes and commitment, not implementation inventory.
- Status is observed state. It is not a substitute for evidence.
- Dates are optional and appear only when the project will maintain them.
- Stable identifiers survive title, grouping, and sequencing changes.
- Planned behavior is never presented as released documentation.

## Required model

The roadmap schema MUST be versioned through a top-level `schemaVersion`
field. It MUST support groups (`id`, `title`) containing ordered items or
phases with these fields:

| Field | Requirement |
|---|---|
| `id` | Stable and unique within the project; never reused. |
| `title` | Short, outcome-oriented name. |
| `summary` | User or operator outcome, not only a component name. |
| `status` | One of `idea`, `planned`, `in_progress`, `done`, `released`, or `cancelled`. |
| `kind` | Optional. `feature` (default) or `phase` for an ordered delivery phase. |
| `target` | Optional release, milestone, or date. |
| `dependencies` | Optional stable IDs whose outcomes are prerequisites. |
| `issues` | Optional external tracking references. |
| `decisions` | Optional governing decision references. |
| `specs` | Frozen `<slug>@v<N>` versions defining the item; a phase lists its features' versions. Required before `in_progress`. |
| `devNote` | Required when a non-trivial item becomes `done`; retained when released. |
| `release` | Required for `released`; the published release tag (`vX.Y.Z`) containing the outcome. |
| `rationale` | Required for `cancelled` unless `decisions` gives the reason. |

Projects MAY extend this model with product-specific fields. Extensions MUST
not redefine the standard fields or statuses.

## Status semantics

- `idea`: exploratory and non-committal.
- `planned`: outcome and rough scope agreed, intended and sequenced, but
  implementation has not started.
- `in_progress`: active implementation has visible ownership and evidence.
- `done`: the item meets its *Definition of Done*, which is: agreed scope
  implemented and merged into `main` (or its `lts/*` line); verification passes
  without unexplained skips; specification, conformance matrix (where one
  exists, it MUST be accepted), and version-appropriate documentation agree; development evidence is
  complete. Items built through the
  [spec-driven cycle](spec-driven-development-cycle.md#10-complete-phase-evidence)
  also meet its §10 list. Publication and release-wide gates may still be
  pending.
- `released`: available in the named published release and supported by conformance
  evidence, documentation, and a development note when required.
- `cancelled`: intentionally not proceeding; rationale or a decision reference
  is present.

The normal progression is `planned` → `in_progress` → `done` → `released`.
Merging code alone is insufficient. Incomplete item-level acceptance remains
`in_progress`; a pending release alone does not.

Statuses record facts that already hold. `done` is set by a separate `chore/`
pull request after the item's merge commit is verified; `released` (with
`release: vX.Y.Z`) is set by a separate `chore/` pull request after
publication and post-publication verification. A prerelease does not make an
item `released`.

Changing an item to `released` additionally requires the applicable release gates,
publication, and release evidence. Several dependent items MAY become `done`
before a shared release gate and become `released` together in one named release.
An internal foundation MAY remain `done` when no separate release claim applies.

Implementation dependencies MAY be satisfied by `done` or `released` outcomes.
A dependency requiring released availability MUST explicitly require `released`
in the consuming item's acceptance criteria or a documented schema extension.
`target` expresses intent and MUST NOT substitute for actual release evidence.

Public roadmaps MUST distinguish completion from availability. Display unreleased
`done` outcomes as **Done — awaiting release**; internal-only foundations MAY use
**Done — internal foundation** when that scope is explicit. Neither label grants
a released support claim. Reopen an item as `in_progress` if its completion
evidence no longer supports the agreed scope.

### Relationship to the spec-driven development cycle

Projects that follow the [spec-driven development
cycle](spec-driven-development-cycle.md) map its events to statuses:

| Cycle event | Roadmap status |
|---|---|
| Specification versions frozen and plan approved | `planned` (if it was `idea`) |
| Implementation waves start | `in_progress` |
| Pull request merged | unchanged; stays `in_progress` |
| Phase evidence complete | `done` |
| Published with release evidence | `released` |

## Development notes

Non-trivial roadmap phases and material architectural changes maintain a note
at `dev-notes/<roadmap-id>-<slug>.md`. The note evolves during implementation
and records:

- roadmap item, scope, status, and intended reader;
- implemented behavior and important boundaries;
- decisions, alternatives, and deviations;
- tests and manual validation, including known gaps;
- security, compatibility, migration, and operational consequences;
- documentation and examples changed; and
- unresolved follow-up work linked to owned tracking items.

Development notes explain why and what was proven. They do not narrate commits
or replace stable user documentation. Tasks belong in the issue tracker;
roadmap items link to them through `issues`.

## Validation

Projects MUST validate roadmap YAML against the declared schema. Validation
MUST reject duplicate IDs, unknown statuses, missing done or released evidence
(including `released` without `release` and non-trivial `done` or `released`
without `devNote`), broken internal references such as unknown `dependencies`,
and references to absent development notes. Generated
roadmap pages MUST be reproducible locally and checked for drift.

Schema status enums, evidence conditions, validators, and renderers MUST support
`done` and `released` together. During migration, validators MAY read the former status
`shipped` as `released`; afterwards they reject it. Adoption MUST NOT
bulk-promote existing items: reconcile each status against implementation and
release evidence.

## Appendix: `roadmap.yaml` example (informative)

```yaml
schemaVersion: 1
groups:
  - id: onboarding
    title: First-run experience
    items:
      - id: RM-12
        title: Guided first upload
        summary: New users upload and classify their first document without help.
        status: released
        release: v2.4.0
        issues: [PROJ-101]
        specs: [guided-first-upload@v1]
        devNote: dev-notes/RM-12-guided-first-upload.md
      - id: RM-13
        title: Bulk import
        summary: Users import a whole folder in one step.
        status: in_progress
        specs: [bulk-import@v1]
        dependencies: [RM-12]
        target: 2.6.0
      - id: RM-20
        title: Offline mode
        summary: Users keep working without a connection.
        status: idea
```

## Agent briefing

> Keep `roadmap.yaml` the only roadmap source; generate pages, check drift.
> Status is observed state: never upgrade on existing code or a merged PR.
> `done` needs agreed scope, tests, docs, and a dev note if non-trivial;
> `released` also needs a named published `release`. Show unreleased work as
> Done — awaiting release. Validate against the schema; never bulk-promote.
