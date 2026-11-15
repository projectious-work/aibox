# Optional tool catalog: requirements

Normative. Requirements use *EARS* patterns. This version restates the V1-05
contract, frozen as contract commit
`3d93aa7e32f5cd49416da2087348baf4cb8c3d42` (amendment 4) and implemented and
accepted at `651506fa`, without changing its meaning. Its sources are
chapter 8 of the v1 product description, sections "Feature composition
contract", "Tool Feature interface" and "Census traceability", and the
`AC-TOOLS` group of chapter 5. Where a requirement and the catalog differ, the
catalog wins and the difference is a specification defect (TOOL-001).

"The catalog" is `catalog/v1-05-features.json`; "the schema" is
`schemas/feature-catalog.schema.json`. Field names in `code` are catalog
fields. A Feature is "aibox-owned" when its catalog `owner` is `aibox`.
Verification levels: **J** = a journey in `journeys.md` (container
acceptance or specification check); **C** = component test only, with the
reason given.

## Catalog

| ID | Requirement | Level |
|---|---|---|
| TOOL-001 | The catalog, valid against the schema, shall be the interface catalog for retained and added tools: reference base, supported and unsupported targets, error codes, every group with its shape, every Feature with its ID, owner, options, dependencies, environment, tools, installation provenance and observable checks, every removed census row, every compatibility note and the smoke tests. `notes` fields are informative and not requirements. | J |
| TOOL-002 | Every one of the 98 `ADDON:*` census rows shall map to a catalog Feature tool or to a removal with an owner decision, through the selection table `ledger/v0-parity/v1-05-tool-selection.json`; no row shall be dropped. | J |
| TOOL-003 | For each retained tool, the selection table and the catalog shall record the install unit, publisher, immutable inputs, license and provenance. | J |

## Identity and shapes

| ID | Requirement | Level |
|---|---|---|
| TOOL-010 | An aibox-owned Feature with catalog ID `<parent>/<name>` shall be published as `<registry>/<parent>/<name>:<featureVersion>` with the tags `<major>` and `<major>.<minor>`; its `devcontainer-feature.json` shall have `id` `<name>` and `version` `featureVersion`. | J |
| TOOL-011 | The sources of an aibox-owned Feature shall live in `features/src/<catalog-id>/` with `devcontainer-feature.json`, `install.sh` and helper files, and nothing outside that directory shall be packaged. | J |
| TOOL-012 | An upstream Feature shall be referenced at the catalog's `upstream.ref` digest with the catalog's `upstream.options`, and aibox shall not republish it. | J |
| TOOL-013 | A `single` group shall have one tool Feature; an `atomic` group shall have one Feature that installs all members and no member or bundle Feature; a `semantic` group shall have one Feature per member and a bundle Feature `<group-path>/all`. | J |
| TOOL-014 | A bundle Feature shall list every member in `dependsOn` at its default options, shall install nothing itself, and shall have no options. | J |

## Options

| ID | Requirement | Level |
|---|---|---|
| TOOL-020 | Where the catalog entry has `options.version`, the Feature shall declare a string `version` option with exactly that `enum` and `default`; otherwise it shall have no `version` option. | J |
| TOOL-021 | If the `version` option has a value outside its `enum`, then the Feature shall fail with `AIBOX_FEATURE_UNSUPPORTED_VERSION` before installation. | J |
| TOOL-022 | No option shall accept `latest`, a range or a URL, and no Feature shall have an `enabled` option. | J |
| TOOL-023 | The specification checks shall reject an example that sets an option key its Feature does not declare. | J |
| TOOL-024 | A tool with method `debian-package` shall install the candidate version of each named package from the base's configured repositories with `--no-install-recommends`, and shall not add or change APT sources. | J |

## Dependencies

| ID | Requirement | Level |
|---|---|---|
| TOOL-030 | Each `dependsOn` entry naming an aibox-owned Feature shall become a standard `dependsOn` reference to `<registry>/<id>:<major>` with an empty option object; an entry naming an upstream Feature shall become its `upstream.ref` with its `upstream.options`. | J |
| TOOL-031 | Each `requires` entry shall become an `installsAfter` reference to `<registry>/<id>`. | J |
| TOOL-032 | If a Feature that `requires` a runtime is installed while that runtime is not installed, then the Feature shall fail with `AIBOX_FEATURE_DEPENDENCY_MISSING`, naming the missing Feature ID in the detail; otherwise it shall use whichever `enum` version of the runtime is installed. | J |
| TOOL-033 | If a `versionCoupling` is not met (rule `equal`: the coupled Feature's installed version equals this Feature's `version`; rule `table`: it is allowed by the table for this `version`), or the coupled Feature is not installed, then the Feature shall fail with `AIBOX_FEATURE_COUPLING_MISMATCH`. | C: the catalog marks the code `component` |
| TOOL-034 | If a tool is installed a second time with a different version, then the second installation shall fail with `AIBOX_FEATURE_VERSION_CONFLICT`; an identical repeat shall succeed without change. | J |

## Targets and errors

| ID | Requirement | Level |
|---|---|---|
| TOOL-040 | Before any change to the image, `install.sh` shall read `/etc/os-release` and `dpkg --print-architecture`, and if the combination is not in the Feature's `targets`, it shall fail with `AIBOX_FEATURE_UNSUPPORTED_TARGET`. | J |
| TOOL-041 | When a failure listed in the catalog's `errors` occurs, the Feature shall write one line `aibox-feature error <CODE> feature=<catalog-id> <detail>` to standard error and exit non-zero. The catalog ID is the Feature under test, except that for `AIBOX_FEATURE_UNSUPPORTED_TARGET` it may be any Feature in its transitive `dependsOn` or `requires` closure. | J |
| TOOL-042 | If an error marked `beforeInstall` occurs, then no package, file or user change shall remain. | J |
| TOOL-043 | Errors marked `component` in the catalog (`AIBOX_FEATURE_COUPLING_MISMATCH`, `AIBOX_FEATURE_DOWNLOAD_FAILED`, `AIBOX_FEATURE_VERIFICATION_FAILED`, `AIBOX_FEATURE_PACKAGE_UNAVAILABLE`) shall be verified by component tests, and the component tests shall cover the guards of `AIBOX_FEATURE_UNSUPPORTED_TARGET`, `AIBOX_FEATURE_UNSUPPORTED_VERSION` and `AIBOX_FEATURE_DEPENDENCY_MISSING` for every aibox-owned Feature they apply to. | C: black-box tests cannot inject these faults |
| TOOL-044 | aibox shall not claim support for Arch, Nix or any target not listed in the catalog's `targets`. | J |

## Presence and absence

| ID | Requirement | Level |
|---|---|---|
| TOOL-050 | When a Feature is selected at a version, every check and every `commands` entry of its tools shall pass in the running container, as the fixture's `remoteUser` through `devcontainer exec` and as `root` through `docker exec --user root`, by the check kinds `version`, `command`, `package`, `texlive` and `path` of chapter 8. | J |
| TOOL-051 | Where no Feature is selected on the reference base, no `commands` entry outside `baseline.provides` shall resolve, no `package` check shall be installed, no `path` check shall exist and no `command` check shall pass. | J |
| TOOL-052 | When a semantic member is selected alone, with the Features its `dependsOn` closure `requires`, the observables of every sibling outside its transitive `dependsOn` and `requires` closure shall be absent, except those exempted by an installed tool's `alsoProvides` or by `baseline.provides` as chapter 8 defines. | J |
| TOOL-053 | When a Feature with `dependsOn` is selected, the Features it depends on shall pass their presence checks without explicit selection. | J |
| TOOL-054 | When a bundle is selected, every member shall pass its presence checks. | J |
| TOOL-055 | Disabling a tool shall be done by omitting its Feature; `customizations.aibox` shall not install tools. | J |

## Side effects

| ID | Requirement | Level |
|---|---|---|
| TOOL-060 | The published `devcontainer-feature.json` of an aibox-owned Feature shall have none of `onCreateCommand`, `updateContentCommand`, `postCreateCommand`, `postStartCommand`, `postAttachCommand`, `mounts`, `privileged`, `capAdd`, `securityOpt`, `init` and `entrypoint`, and its `containerEnv` shall be exactly the catalog's `containerEnv`, or absent when the catalog has none. | J |
| TOOL-061 | A build shall create, change or remove no file under the workspace folder other than `.devcontainer/devcontainer-lock.json`. | J |
| TOOL-062 | A tool that would update itself shall be configured not to, and no tool shall replace a project lockfile. | J |

## Smoke tests and stacks

| ID | Requirement | Level |
|---|---|---|
| TOOL-070 | For each `smoke` entry, with the fixture copied under its base name into the workspace root, `argv` run there as the remote user shall exit 0 and create `produces`, and where `producesMagicHex` is given the file shall start with those bytes. | J |
| TOOL-071 | Each known-risk stack of chapter 5 (`STACK-PYTHON`, `STACK-NODE`, `STACK-DOCS`, `STACK-NATIVE`, `STACK-SECURITY`) shall build and pass the presence checks of its Features. | J |

## Publishing and inspection

| ID | Requirement | Level |
|---|---|---|
| TOOL-080 | `scripts/v1-05-features.sh publish --registry <host>/<prefix>` shall publish every aibox-owned Feature to `<host>/<prefix>/<catalog-id>` with at least the tags `<featureVersion>`, `<major>` and `<major>.<minor>`, rewrite aibox `dependsOn` and `installsAfter` references to that prefix, leave upstream references unchanged, and exit 0 on success and non-zero with a message on failure. | J |
| TOOL-081 | `aibox inspect` and the MCP tool `inspect_workspace` shall report a configuration entry with key `devcontainer.features`, value the sorted list of selected Feature references, `layer` `project`, `source` `.devcontainer/devcontainer.json`, `domain` `native`, `redacted` false and `invocationOnly` false, shall never echo Feature options, and shall omit the entry when no Feature is selected. | J |
| TOOL-082 | Acceptance shall pass on each target in the catalog's `targets`. | J |
