# Optional tool catalog: acceptance journeys

Normative. `interfaces.json` and the catalog prevail over this file (see
`requirements.md`, "Precedence and reading rules"). Container journeys are
realized by the AC-TOOLS harness in `tests/acceptance/v1-05/`; the
specification-check journey by `scripts/spec/validate.mjs` and
`scripts/spec/validate-v1-05-catalog.py`; the guards journey by the component
tests in `features/test/`.

## Common setup

- **Base.** Every build starts from the catalog `baseline` fixture
  (`fixtures/v1-05-base`) and adds Features to its `features` object.
- **Closure.** Except for `AC-TOOLS/requires`, every selection also includes,
  at their defaults, the Features that any Feature in the selection's
  transitive `dependsOn` closure `requires`, repeated until no Feature is
  added.
- **Registry.** Each run starts an OCI Distribution registry of its choice,
  recorded in the evidence, on `localhost:5000`, which the pinned Dev
  Container CLI contacts over plain HTTP. aibox-owned Features are published
  with TOOL-090 to `localhost:5000/aibox-features` and selected by their
  `<major>` tag, so that automatic dependencies resolve to the same instance.
  Upstream Features are selected by digest.
- **Checks** run as TOOL-050 to TOOL-053 define; absence as TOOL-060 to
  TOOL-063 define.
- **Errors.** An error journey passes when the build fails and some line of
  its combined output contains `aibox-feature error <CODE> feature=<catalog-id>`,
  regardless of any prefix the build adds; the catalog ID is as TOOL-041
  defines.
- **Build sharing.** Presence checks may share one build across Features that
  do not conflict; absence, isolation, bundle and error journeys use exactly
  the selection their item defines.
- **Targets.** A run covers the host's native architecture; acceptance needs
  one passing run on each target in `interfaces.json` `targets`.
- **Owner filter.** Items marked "aibox" apply only to aibox-owned Features.

## Test IDs

Each container test is named
`AC-TOOLS/<item> [<feature-id> [<tool>]] [<target-id>] [<key>=<value> ...]`,
with the slots per item in `interfaces.json` `acceptance.testIdSlots`, for
example `AC-TOOLS/present python/uv uv debian-13-arm64` or
`AC-TOOLS/stack STACK-PYTHON debian-13-amd64`. Test IDs are the same in both
scopes.

## Scope

An acceptance run uses **sampled** scope; **full** scope runs on the owner's
command.

- **Sampled.** A sampled item builds the first eligible Feature in catalog
  order for each value of its key. The key is the install method: a bundle's
  install method is `bundle`, and any other Feature's is the `method` of its
  first tool. `AC-TOOLS/requires` is keyed by the distinct `requires` set.
  Because the first Feature in install order fails, a sampled
  `unsupported-target` build may exercise a dependency's guard rather than
  the sampled Feature's; J-TOOLS-16 covers each Feature's own guard.
- **Full.** Builds every eligible Feature: `requires` per Feature with
  `requires`, `unsupported-version` per Feature with `version`, and
  `unsupported-target` per Feature × fixture.

## Journeys

| ID | Item | Per | Covers | Expected |
|---|---|---|---|---|
| J-TOOLS-01 | `AC-TOOLS/present` | non-bundle Feature × target | TOOL-012, TOOL-015, TOOL-020, TOOL-024, TOOL-044, TOOL-050 to TOOL-053 | All checks and commands pass for the default version; aibox Features resolve from their `<major>` tag |
| J-TOOLS-02 | `AC-TOOLS/version` | Feature × non-default `enum` value × target | TOOL-015, TOOL-020, TOOL-050, TOOL-051, TOOL-053 | All checks pass for that version |
| J-TOOLS-03 | `AC-TOOLS/absent` | target | TOOL-060, TOOL-062, TOOL-063, TOOL-064 | The reference base with no Feature has no absence observable outside `baseline.provides` |
| J-TOOLS-04 | `AC-TOOLS/isolated` | semantic member × target | TOOL-019, TOOL-061, TOOL-062, TOOL-063, TOOL-064 | The member alone leaves its other siblings absent |
| J-TOOLS-05 | `AC-TOOLS/bundle` | bundle × target | TOOL-014, TOOL-066 | The bundle passes the presence checks of every member |
| J-TOOLS-06 | `AC-TOOLS/depends` | Feature with `dependsOn` × target | TOOL-030, TOOL-065 | The `dependsOn` Features pass presence without explicit selection |
| J-TOOLS-07 | `AC-TOOLS/requires` (aibox) | distinct `requires` set (sampled) | TOOL-032, TOOL-041 | Selected without a Feature it `requires`, the build fails with `AIBOX_FEATURE_DEPENDENCY_MISSING` naming it |
| J-TOOLS-08 | `AC-TOOLS/unsupported-version` (aibox) | install method among Features with `version` (sampled) | TOOL-021, TOOL-041 | `version` `0.0.0-unsupported` fails with `AIBOX_FEATURE_UNSUPPORTED_VERSION` |
| J-TOOLS-09 | `AC-TOOLS/unsupported-target` (aibox) | install method × `unsupportedTargets` fixture (sampled) | TOOL-040, TOOL-041 | Each fixture fails with `AIBOX_FEATURE_UNSUPPORTED_TARGET` |
| J-TOOLS-10 | `AC-TOOLS/conflict` | member with more than one `enum` value | TOOL-034 (conflict) | A bundle plus its member at the first non-default value fails with `AIBOX_FEATURE_VERSION_CONFLICT` |
| J-TOOLS-11 | `AC-TOOLS/side-effects` (aibox) | Feature | TOOL-070, TOOL-071 | The published metadata, read from the Feature's tarball layer in the test registry, meets TOOL-070, and the workspace folder is unchanged after the build except for the CLI lockfile |
| J-TOOLS-12 | `AC-TOOLS/smoke` | smoke entry × target | TOOL-080 | Each smoke entry creates its file |
| J-TOOLS-13 | `AC-TOOLS/stack` | stack × target | TOOL-081 | Each known-risk stack builds and passes presence checks |
| J-TOOLS-14 | `AC-TOOLS/inspect` | example `customized` and `minimal` | TOOL-091 | CLI and MCP report `devcontainer.features` identically, without options |
| J-TOOLS-15 | `AC-TOOLS/trace` | census row | TOOL-002 | Every census row maps to a catalog Feature tool or a removal with a decision |
| J-TOOLS-16 | `AC-TOOLS/guards` | run | TOOL-043 and the C entries of `requirements.md` | `scripts/v1-05-features.sh test` (`node --test features/test/*.test.mjs`) passes with none skipped, including a test of each of the three sampled items' error codes for every aibox Feature they apply to |
| J-TOOLS-17 | Specification check | run | TOOL-001, TOOL-003 (fields), TOOL-004 (shapes), TOOL-005, TOOL-013, TOOL-022 (version values), TOOL-023 | `scripts/spec/validate.mjs` validates the catalog against the schema and `scripts/spec/validate-v1-05-catalog.py` passes its cross-reference checks |

TOOL-090's success path is exercised by every container journey, because each
publishes the Features it builds; its failure paths are component-tested.
