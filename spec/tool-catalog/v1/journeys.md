# Optional tool catalog: acceptance journeys

Normative. Each journey is one `AC-TOOLS` item of chapter 5 of the v1 product
description. The AC-TOOLS harness in `tests/acceptance/v1-05/` realizes them.
Test IDs, selections, sampled and full scope, and the known-risk stacks are
as chapter 5 defines them; this file names which requirements each item
covers.

## Common setup

- Every build starts from the catalog `baseline` fixture
  (`fixtures/v1-05-base`) and adds Features to its `features` object.
- Except for `AC-TOOLS/requires`, every selection also includes, at their
  defaults, the Features that any Feature in the selection's transitive
  `dependsOn` closure `requires`, repeated until no Feature is added.
- aibox-owned Features are published with TOOL-080 to an OCI Distribution
  registry on `localhost:5000` under `localhost:5000/aibox-features` and
  selected by their `<major>` tag; upstream Features are selected by digest.
- Checks run as in TOOL-050. An error passes when the build fails and a line
  of its combined output contains `aibox-feature error <CODE> feature=<catalog-id>`.
- A run covers the host's native architecture; acceptance needs one passing
  run on each target in the catalog's `targets` (TOOL-082).
- Items marked "aibox" apply only to aibox-owned Features.

## Journeys

| ID | AC-TOOLS item | Per | Covers | Expected |
|---|---|---|---|---|
| J-TOOLS-01 | `AC-TOOLS/present` | non-bundle Feature × target | TOOL-010, TOOL-012, TOOL-020, TOOL-024, TOOL-030, TOOL-031, TOOL-050 | All checks and commands pass for the default version |
| J-TOOLS-02 | `AC-TOOLS/version` | Feature × non-default `enum` value × target | TOOL-020, TOOL-050 | All checks pass for that version |
| J-TOOLS-03 | `AC-TOOLS/absent` | target | TOOL-051, TOOL-055 | The reference base with no Feature has no absence observable outside `baseline.provides` |
| J-TOOLS-04 | `AC-TOOLS/isolated` | semantic member × target | TOOL-013, TOOL-052 | The member alone leaves its other siblings absent |
| J-TOOLS-05 | `AC-TOOLS/bundle` | bundle × target | TOOL-013, TOOL-014, TOOL-054 | The bundle passes the presence checks of every member |
| J-TOOLS-06 | `AC-TOOLS/depends` | Feature with `dependsOn` × target | TOOL-030, TOOL-053 | The `dependsOn` Features pass presence without explicit selection |
| J-TOOLS-07 | `AC-TOOLS/requires` (aibox) | distinct `requires` set, sampled | TOOL-031, TOOL-032, TOOL-041, TOOL-042 | Selected without a Feature it `requires`, the build fails with `AIBOX_FEATURE_DEPENDENCY_MISSING` naming it |
| J-TOOLS-08 | `AC-TOOLS/unsupported-version` (aibox) | install method among Features with `version`, sampled | TOOL-021, TOOL-041, TOOL-042 | `version` `0.0.0-unsupported` fails with `AIBOX_FEATURE_UNSUPPORTED_VERSION` |
| J-TOOLS-09 | `AC-TOOLS/unsupported-target` (aibox) | install method × `unsupportedTargets` fixture, sampled | TOOL-040, TOOL-041, TOOL-042, TOOL-044 | Each fixture fails with `AIBOX_FEATURE_UNSUPPORTED_TARGET` |
| J-TOOLS-10 | `AC-TOOLS/conflict` | member with more than one `enum` value | TOOL-034 | A bundle plus its member at the first non-default value fails with `AIBOX_FEATURE_VERSION_CONFLICT` |
| J-TOOLS-11 | `AC-TOOLS/side-effects` (aibox) | Feature | TOOL-011, TOOL-060, TOOL-061, TOOL-062 | Published metadata meets TOOL-060 and the workspace folder is unchanged after the build |
| J-TOOLS-12 | `AC-TOOLS/smoke` | smoke entry × target | TOOL-070 | Each smoke entry creates its file |
| J-TOOLS-13 | `AC-TOOLS/stack` | stack × target | TOOL-071 | Each known-risk stack builds and passes presence checks |
| J-TOOLS-14 | `AC-TOOLS/inspect` | example `customized` and `minimal` | TOOL-081 | CLI and MCP report `devcontainer.features` identically, without options |
| J-TOOLS-15 | `AC-TOOLS/trace` | census row (specification check) | TOOL-001, TOOL-002, TOOL-003, TOOL-022, TOOL-023 | Every census row maps to a catalog Feature tool or a removal with a decision; the catalog validates against the schema and the cross-reference checks |
| J-TOOLS-16 | `AC-TOOLS/guards` | run | TOOL-033, TOOL-043 | The component tests pass with none skipped, including each sampled item's error code for every aibox-owned Feature it applies to |

TOOL-080 is exercised by every container journey, because each publishes the
Features it builds.
