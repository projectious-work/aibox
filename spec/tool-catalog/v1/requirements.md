# Optional tool catalog: requirements

Normative. Requirements use *EARS* patterns.

## Precedence and reading rules

- `catalog/v1-05-features.json` ("the catalog", valid against
  `schemas/feature-catalog.schema.json`, "the schema") and `interfaces.json`
  prevail over this file and over `journeys.md`. Where they differ, the
  catalog and `interfaces.json` win, and the difference is a specification
  defect that is fixed in `tool-catalog@v2`.
- The chapters of the v1 product description cited in `overview.md` are
  provenance only. They do not add to or override this version.
- `notes` fields of the catalog are informative research records, not
  requirements.
- Field names in `code` are catalog fields unless `interfaces.json` is named.
  A Feature is **aibox-owned** when its catalog `owner` is `aibox`, and
  **upstream** otherwise.

**Verified by** names the journey of `journeys.md` (J) or the component test
file (C) that verifies the requirement. "Not separately verified" marks an
obligation that is part of the contract but that no accepted V1-05 check
exercises on its own; it is reviewed against the catalog and the sources.

## Catalog and census

| ID | Requirement | Verified by |
|---|---|---|
| TOOL-001 | The catalog shall list the reference base, the supported and unsupported targets, the error codes, every group with its shape, every Feature with its ID, owner, options, dependencies, environment, tools, installation provenance and observable checks, every removed census row, every compatibility note and the smoke tests, and shall validate against the schema. | J-TOOLS-17 |
| TOOL-002 | Every one of the 98 `ADDON:<recipe>/<tool>` census rows of v0 revision `interfaces.json` `external.censusBaselineRevision` shall map to exactly one catalog Feature tool or to a removal with an owner decision, and no census row shall be dropped. Several rows may map to the same Feature. | J-TOOLS-15 |
| TOOL-003 | The selection table `interfaces.json` `external.selectionTable`, keyed by `ADDON:<recipe>/<tool>`, shall record for each row `v0-source`, `v0-default`, `v0-versions`, `v1-ref`, `v1-option`, `install-owner`, `dependency`, `license`, `digest`, `platforms`, `test-id` and `decision`, and for each retained tool its install unit, publisher, immutable inputs, license, and maintenance and security evidence. | J-TOOLS-17 (fields); content not separately verified |
| TOOL-004 | Each retained v0 addon group shall be classified by cohesion as `single`, `atomic` or `semantic`, and the selection table shall document the classification and the member-to-Feature mapping. | J-TOOLS-17 (group shapes); the table's documentation is not separately verified |
| TOOL-005 | A row must not claim independent version or disable support that its Feature does not provide. | J-TOOLS-17 (`v1-option` equals the Feature's option) |
| TOOL-006 | If a v0 project used per-tool switches within a group that is now `atomic`, then the conversion shall report those states as an explicit compatibility decision and shall not claim that the single Feature preserves them. | Not separately verified (conversion is V1-20) |
| TOOL-007 | A new aibox-owned Feature outside the catalog shall require an owner decision; a missing upstream Feature alone shall not be permission to add one. | Not separately verified (process rule) |
| TOOL-008 | aibox shall not submit these Features to the public Feature index, and shall record exact release references only after the artifacts exist. | Not separately verified |
| TOOL-009 | The base image shall enumerate its packages and runtime glue from `ledger/v0-parity/base-build.json`, with build inputs and packages pinned by digest or version and with update owner and provenance; tools outside the catalog shall not disappear, and a heavy optional tool shall not be forced into the base because an optional Feature could not remove it later. | Not separately verified (base toolchain is verified by `native-starters@v1`) |

## Identity and shapes

| ID | Requirement | Verified by |
|---|---|---|
| TOOL-010 | An aibox-owned Feature with catalog ID `<parent>/<name>` shall be published as `<registry>/<parent>/<name>:<featureVersion>` with the moving tags `<major>` and `<major>.<minor>`, and its `devcontainer-feature.json` shall have `id` `<name>` and `version` `featureVersion`. | C: `features/test/metadata.test.mjs`, `features/test/publish.test.mjs`; J-TOOLS-01 selects the `<major>` tag |
| TOOL-011 | The sources of an aibox-owned Feature shall live in `features/src/<catalog-id>/` with `devcontainer-feature.json`, `install.sh` and any helper files, and nothing outside that directory shall be packaged. | C: `features/test/metadata.test.mjs` |
| TOOL-012 | An upstream Feature shall be referenced at the catalog's `upstream.ref` digest with the catalog's `upstream.options`, and aibox shall not republish it. | J-TOOLS-01 |
| TOOL-013 | A `single` group shall have one tool Feature; an `atomic` group shall have one Feature that installs all members and no member or bundle Feature; a `semantic` group shall have one Feature per selectable member and a bundle Feature `<group-path>/all`. An `all` bundle shall exist only where its members also have individual Features. | J-TOOLS-17; C: `features/test/bundles.test.mjs` |
| TOOL-014 | A bundle Feature shall list every member in `dependsOn` at its default options, shall install nothing itself, and shall have no options. | C: `features/test/bundles.test.mjs`; J-TOOLS-05 |

## Installation method

| ID | Requirement | Verified by |
|---|---|---|
| TOOL-015 | Each Feature shall own a support matrix by selected tool version, container distro and release, and architecture, and its `install.sh` shall select, for each matrix cell, an applicable canonical method from the tool producer's installation documentation; different tool versions may use different methods. | J-TOOLS-01, J-TOOLS-02 (each declared cell builds) |
| TOOL-016 | Where the producer has no applicable method, an official distro package shall be used only after checking that the exact requested tool version is available from the selected base's configured repository on that architecture, that its provenance and dependencies are acceptable, and that installation and runtime tests pass; an APT index entry for one release or architecture shall not qualify other Debian-based images. | Not separately verified (recorded per tool in the catalog and the selection table) |
| TOOL-017 | A Feature shall pin and verify every release or package input it consumes; the Feature's own version shall not be taken to pin the installed tool version. | C: `features/test/foundation.test.mjs`, `features/test/tools.test.mjs` (`AIBOX_FEATURE_VERIFICATION_FAILED`) |
| TOOL-018 | A Feature shall reuse a qualified upstream Feature, native package manager or official distribution inside it where that meets this contract, and shall not implement an addon resolver or package-manager API; `customizations.aibox` shall carry UX preferences only and shall not install tools. | Not separately verified |
| TOOL-019 | Hard prerequisites shall be native dependencies or documented prerequisites, and an optional selection shall not acquire unexpected extras; a disabled heavy component shall not become an implicit dependency. | J-TOOLS-04 (siblings stay absent) |

## Options

| ID | Requirement | Verified by |
|---|---|---|
| TOOL-020 | Where the catalog entry has `options.version`, the Feature shall declare a string `version` option with exactly that `enum` and `default`; otherwise it shall have no `version` option. The option is the tool's own version; a package revision such as an APT `-0` suffix is an installation detail. | C: `features/test/metadata.test.mjs`; J-TOOLS-01, J-TOOLS-02 |
| TOOL-021 | If the `version` option of an aibox-owned Feature has a value outside its `enum`, then the Feature shall fail with `AIBOX_FEATURE_UNSUPPORTED_VERSION` before installation. | J-TOOLS-08; C: `features/test/conventions.test.mjs` |
| TOOL-022 | No option shall accept `latest`, a range or a URL, and no Feature shall have an `enabled` option. | J-TOOLS-17 rejects ranges and URLs in `version` values (`exactVersion` pattern); `latest` and `enabled` are not separately verified |
| TOOL-023 | The specification checks shall reject an example that sets an option key its Feature does not declare. | J-TOOLS-17 (example Feature options must equal the catalog's) |
| TOOL-024 | A tool with method `debian-package` shall install the candidate version of each named package from the base's configured repositories with `--no-install-recommends`, and its Feature shall not add or change APT sources. The base's snapshot pin, not the Feature, fixes that version, and the `version` of a `package` check is the version expected on the reference base. | C: `features/test/tools.test.mjs`; J-TOOLS-01 (package checks) |
| TOOL-025 | The minimal and curated starters shall configure `trixie`, `trixie-updates` and `trixie-security` from one snapshot. | Verified by the starter checks of `native-starters@v1` |

## Dependencies

| ID | Requirement | Verified by |
|---|---|---|
| TOOL-030 | Each `dependsOn` entry naming an aibox-owned Feature shall become a standard `dependsOn` reference to `<registry>/<id>:<major>` with an empty option object, and an entry naming an upstream Feature shall become its `upstream.ref` with its `upstream.options`, so that the Dev Container CLI installs it automatically. | C: `features/test/metadata.test.mjs`; J-TOOLS-06 |
| TOOL-031 | Each `requires` entry, used only for multi-version runtimes, shall become an `installsAfter` reference to `<registry>/<id>`. | C: `features/test/metadata.test.mjs` |
| TOOL-032 | If an aibox-owned Feature that `requires` a runtime is installed while that runtime is not installed, then the Feature shall fail before installation with `AIBOX_FEATURE_DEPENDENCY_MISSING`, naming the missing Feature ID in the detail. Otherwise it shall use whichever `enum` version of the runtime is installed. | J-TOOLS-07; C: `features/test/conventions.test.mjs` |
| TOOL-033 | If an aibox-owned Feature has a `versionCoupling` and the coupled Feature is not installed, or (rule `equal`) is installed at a version other than this Feature's `version`, or (rule `table`) is installed at a version the table does not allow for this Feature's `version`, then the Feature shall fail with `AIBOX_FEATURE_COUPLING_MISMATCH`. | C: `features/test/conventions.test.mjs` (the catalog marks the code `component`) |
| TOOL-034 | If a tool of an aibox-owned Feature is installed a second time with a different version, for example through a bundle and an explicit member selection, then the second installation shall fail with `AIBOX_FEATURE_VERSION_CONFLICT`; an identical repeat shall succeed without change. | J-TOOLS-10 (conflict); C: `features/test/conventions.test.mjs` (both) |

## Targets and errors

| ID | Requirement | Verified by |
|---|---|---|
| TOOL-040 | Before any change to the image, the `install.sh` of an aibox-owned Feature shall read `/etc/os-release` and `dpkg --print-architecture`, and if the combination is not in the Feature's `targets`, it shall fail with `AIBOX_FEATURE_UNSUPPORTED_TARGET`. The catalog's `unsupportedTargets` are the fixtures for this behaviour. | J-TOOLS-09; C: `features/test/conventions.test.mjs` |
| TOOL-041 | When a failure listed in the catalog's `errors` occurs in an aibox-owned Feature, the Feature shall write one line `aibox-feature error <CODE> feature=<catalog-id> <detail>` to standard error and exit non-zero, so that the Dev Container build fails. For `AIBOX_FEATURE_UNSUPPORTED_TARGET` the catalog ID may be the Feature under test or any Feature it transitively `dependsOn` or `requires`, because the first Feature in install order fails; for every other code it is the Feature under test. | J-TOOLS-07, J-TOOLS-08, J-TOOLS-09; C: `features/test/*.test.mjs` |
| TOOL-042 | If an error marked `beforeInstall` occurs, then no package, file or user change shall remain. | C: `features/test/conventions.test.mjs`, `foundation.test.mjs`, `tools.test.mjs` ("before any change") |
| TOOL-043 | Errors marked `component` in the catalog shall be verified by component tests, and the component tests shall cover the `AIBOX_FEATURE_UNSUPPORTED_TARGET`, `AIBOX_FEATURE_UNSUPPORTED_VERSION` and `AIBOX_FEATURE_DEPENDENCY_MISSING` guards of every aibox-owned Feature they apply to. Download and checksum failures, a missing exact package version and wrong-architecture failures shall be tested where the selected installation mechanism performs those operations. | J-TOOLS-16 |
| TOOL-044 | Upstream Features have no aibox error contract; for them, acceptance shall verify only presence, absence, isolation, bundles and smoke tests on the supported targets. | J-TOOLS-01, J-TOOLS-03 to J-TOOLS-05, J-TOOLS-12 |
| TOOL-045 | aibox shall support only the targets in the catalog's `targets`, and shall not infer Arch or Nix support from a generic Linux claim. | Not separately verified (no Arch or Nix fixture exists) |
| TOOL-046 | A passing Feature build shall not be taken to establish a client editor's own compatibility with that base. | Not separately verified |

## Presence and checks

| ID | Requirement | Verified by |
|---|---|---|
| TOOL-050 | When a Feature is selected, every check of each of its tools and every `commands` entry shall pass in the running container, evaluated as the fixture's `remoteUser` through `devcontainer exec` and again as `root` through `docker exec --user root <containerId>`. A `commands` entry passes when it resolves through `command -v`. | J-TOOLS-01, J-TOOLS-02 |
| TOOL-051 | A check shall pass by its `kind` as `interfaces.json` `checkKinds` defines: `version` when `argv` exits 0 and its output matches `regex` after `{version}` is replaced by the regular-expression-escaped selected version; `command` when `argv` exits 0 and its output matches `regex`, which has no placeholder; `package` when `dpkg-query -W -f '${Version}' <package>` exits 0 and prints exactly `version`; `texlive` when `tlmgr info --only-installed --data name <package>` prints `package`; `path` when `test -d <path>` succeeds. | J-TOOLS-01, J-TOOLS-02 |
| TOOL-052 | Regular expressions shall be ECMAScript, searched anywhere in the combined standard output and standard error, with `^` and `$` matching at line boundaries. | J-TOOLS-01 |
| TOOL-053 | The selected version of a `version` check shall be the Feature's `version` option value, otherwise the tool's `version`. `versionFrom` names the Feature whose version the tool follows and shall not change how its checks run. A tool's `packages` lists what it installs and is informative. | J-TOOLS-01, J-TOOLS-02 |

## Absence

| ID | Requirement | Verified by |
|---|---|---|
| TOOL-060 | While no Feature is selected on the reference base, no `commands` entry other than those in `baseline.provides` shall resolve through `command -v`, no `package` check shall be installed, no `path` check shall exist, and no `command` check shall pass. A package is installed only when `dpkg-query -W -f '${db:Status-Status}' <package>` prints `installed`. | J-TOOLS-03 |
| TOOL-061 | When a semantic member is selected alone, together with the Features its `dependsOn` closure `requires`, the absence observables of TOOL-060 of every sibling shall be absent unless the sibling is in that member's transitive `dependsOn` or `requires` closure. | J-TOOLS-04 |
| TOOL-062 | An installed tool's `alsoProvides` entry for the target shall exempt only what it names: a `command` entry exempts that command, a `package` entry exempts that package, and a `feature` entry exempts every observable of that Feature. A command name in `baseline.provides`, or in an installed tool's `alsoProvides` `command` entry, shall exempt only the matching `commands` entry; a `command` check whose `argv` runs that command remains an absence observable. | J-TOOLS-03, J-TOOLS-04 |
| TOOL-063 | `version` and `texlive` checks shall not be absence observables, and a tool whose observables are all exempt or in `baseline.provides` shall have no absence test. | J-TOOLS-03, J-TOOLS-04 |
| TOOL-064 | Disabling a tool shall mean omitting its Feature, and an explicit disable shall result in absence where v0 promises it, including absence from the base image; there shall be no hidden addon solver and no unexpected extras. | J-TOOLS-03, J-TOOLS-04 (absence); the absence of a solver is not separately verified |
| TOOL-065 | Where a Feature with `dependsOn` is selected, the Features it depends on shall pass their presence checks without explicit selection. | J-TOOLS-06 |
| TOOL-066 | Where a bundle is selected, every member shall pass its presence checks. | J-TOOLS-05 |

## Side effects

| ID | Requirement | Verified by |
|---|---|---|
| TOOL-070 | The published `devcontainer-feature.json` of an aibox-owned Feature, read from the Feature's tarball layer in the test registry, shall have none of `onCreateCommand`, `updateContentCommand`, `postCreateCommand`, `postStartCommand`, `postAttachCommand`, `mounts`, `privileged`, `capAdd`, `securityOpt`, `init` and `entrypoint`, and its `containerEnv` shall be exactly the catalog's `containerEnv`, or absent when the catalog has none. `privileges` is explanatory text only. | J-TOOLS-11; C: `features/test/metadata.test.mjs` |
| TOOL-071 | A build shall create, change or remove no file under the workspace folder other than the Dev Container CLI's own `.devcontainer/devcontainer-lock.json`. | J-TOOLS-11 |
| TOOL-072 | A tool that would update itself shall be configured not to, and tools shall never replace project lockfiles; native project lockfiles stay authoritative for application dependencies. | Not separately verified (recorded per tool in the catalog) |

## Smoke tests and stacks

| ID | Requirement | Verified by |
|---|---|---|
| TOOL-080 | For each `smoke` entry, the fixture file or directory shall be copied under its base name into the root of the workspace folder, and `argv`, run there as the remote user, shall exit 0 and create `produces` relative to the workspace folder; where `producesMagicHex` is given, the file's first bytes shall equal those hexadecimal bytes. | J-TOOLS-12 |
| TOOL-081 | Each known-risk stack in `interfaces.json` `acceptance.stacks` shall build with its listed Features, plus `node/node`, `go/go` and `rust/rustc` where a member requires them, and shall pass the presence checks of its Features. | J-TOOLS-13 |

## Publishing and inspection

| ID | Requirement | Verified by |
|---|---|---|
| TOOL-090 | `scripts/v1-05-features.sh publish --registry <host>/<prefix>`, run from the repository root with `docker`, `devcontainer` and `node` found on `PATH`, shall publish every aibox-owned Feature to `<host>/<prefix>/<catalog-id>` with at least the tags `<featureVersion>`, `<major>` and `<major>.<minor>`, shall rewrite aibox `dependsOn` and `installsAfter` references to that prefix, shall leave upstream references unchanged, and shall exit 0 on success and non-zero with a message on failure. Tags the pinned CLI adds, such as `latest`, are allowed but never selected by tests. | C: `features/test/publish.test.mjs`; success path in every container journey |
| TOOL-091 | `aibox inspect` and the MCP tool `inspect_workspace` shall report a configuration entry with key `devcontainer.features`, value the sorted list of selected Feature references, `layer` `project`, `source` `.devcontainer/devcontainer.json`, `domain` `native`, `redacted` false and `invocationOnly` false; Feature options shall never be echoed, and the entry shall be absent when no Feature is selected. | J-TOOLS-14 |

Retired IDs: none. IDs TOOL-054 to TOOL-059, TOOL-067 to TOOL-069,
TOOL-073 to TOOL-079 and TOOL-082 to TOOL-089 are unused.
