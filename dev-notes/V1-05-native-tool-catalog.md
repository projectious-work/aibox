# V1-05 native tool catalog

V1-05 contract is frozen at contract commit
`d1c7adb487bd4cfa7449550ffcd3da05153537ed` on `contract/v1-05-tool-catalog`
(amendment 3, DEC-20261031_1750-NeatDaisy). It replaces amendment 2,
`d3f6e08eaf232eb6aa7c77cf6c6d2fcea10d1c3c` (DEC-20261029_2134-NobleBlossom),
amendment 1, `59c35acc1cf230c81658a5763205b1f91a0b1092` (DEC-20261029_1725-AmberBeam), and
the first contract commit `34d06da0cf42d73ff05971c7235ae1324d7e768a`
(DEC-20261029_1515-ShinySpire). The acceptance-test and implementation tracks
continue from that commit.

## Baseline

- Specification base: `v1.x-dev` at `9cdcbcbe` (R-CYCLE from PR #496).
- Phase: V1-05, restore the complete optional tool catalog.
- Chapters: 04 reuse targets, 05 acceptance (`AC-TOOLS`), 08 native
  artifacts (Feature composition contract and Tool Feature interface), 19
  configuration realization.
- Requirements: R-REUSE, R-DEPENDENCIES, R-CYCLE, R-REVIEW.
- Census baseline: v0 `9061bf76d2a1f6ac3c7093d8c605bcc4de7eddf6`, 98 addon
  tool rows in 41 recipes.
- Decisions:
  - DEC-20261002_0046-SmoothLotus and DEC-20261002_0054-AmberCompass:
    grouped Features, per-cell qualification.
  - DEC-20261026_1412-FastForge: frozen test-ready contract and parallel
    tracks.
  - DEC-20261029_0802-DeepOtter: trixie amd64/arm64 only, exact pins,
    upstream Features used directly, separate contract branch.
  - DEC-20261029_0835-HumbleBridge: row dispositions and removals.
  - DEC-20261029_1034-ActiveJay: additions, uv-managed Python, security
    groups.
  - DEC-20261029_1409-DaringTrail: review resolutions for runtime
    dependencies, Debian versions, version counts and approved deviations.

## Contract contents

| Artifact | Role |
|---|---|
| [`spec/v1/catalog/v1-05-features.json`](../spec/v1/catalog/v1-05-features.json) | Interface catalog: 39 groups, 120 Features, reference base, removals, compatibility notes, error codes, targets, smoke tests |
| [`spec/v1/feature-catalog.schema.json`](../spec/v1/feature-catalog.schema.json) | Catalog schema, checked by `validate.mjs` |
| [Chapter 8, Tool Feature interface](../spec/v1/08-native-artifacts.md#tool-feature-interface) | Identity, shapes, options, dependencies, targets, errors, presence, absence, side effects, publishing for tests, inspection |
| [Chapter 5, AC-TOOLS items](../spec/v1/05-acceptance.md#ac-tools--complete-optional-tool-matrix) | Acceptance item IDs and known-risk stacks |
| [Chapter 19](../spec/v1/19-configuration-realization.md) | Grouped Feature realization and owner-approval rule |
| `spec/v1/fixtures/v1-05-base/`, `spec/v1/fixtures/v1-05-targets/` | Reference base (snapshot `20261028T000000Z`) and unsupported-target fixtures |
| [`spec/v1/ledger/v1-05-tool-selection.json`](../spec/v1/ledger/v1-05-tool-selection.json) | Census table: each of the 98 rows retained with its Feature or removed with its decision |
| `spec/v1/fixtures/v1-05-smoke/` | Smoke-test inputs |
| `spec/v1/scripts/validate-v1-05-catalog.py` | Cross-reference checks between census, catalog and table |

Census result: 93 rows retained, 5 removed (mistral, starlight, docusaurus,
mkdocs, rich). Catalog additions without a census row: vibe, typescript (tsc,
biome, oxlint), zig (zig, zls), preview-enhanced (glow, bat), security/crypto (age, sops, gnupg,
openssh-client, minisign, openssl) and semgrep.

Version and install-route data were researched on 2026-10-29 from producer
release data, registries and the trixie, trixie-updates and trixie-security
package indexes for both architectures. Version-output regular expressions use
a version-token rule rather than exact output formats, so an output layout
change does not become a false contract failure. The
[gap review](V1-05-gap-review.md) is the pre-decision input and is superseded
by the catalog.

## Inventory and trace

| Requirement | Planned work | Evidence | Verification |
|---|---|---|---|
| R-REUSE | Upstream Features used directly; owned Features reuse producer installers | Catalog `upstream`, `method`, `docsUrl`, `verification` | Specification check; `AC-TOOLS/present` for upstream Features |
| R-DEPENDENCIES | Exact pins, checksums or signatures, licences, Debian versions from the base snapshot | Catalog tools; selection table | Specification check; component tests for download and verification errors; `AC-TOOLS/present`, `AC-TOOLS/version` |
| R-CYCLE | Frozen contract, separate tracks, integration loop | This logbook | Plan review; contract commit record |
| R-REVIEW | Independent plan review and implementation verification | Reviewer record below | Review reports |
| AC-TOOLS | Every `AC-TOOLS/*` item in chapter 5 | Acceptance test results on both targets | Black-box acceptance tests |
| Chapter 8 census obligations 1–5 | Census traceability, compatibility notes, isolation and side-effect checks | Catalog and selection table | `AC-TOOLS/trace`, `AC-TOOLS/isolated`, `AC-TOOLS/side-effects` |

## Plan

Implementation track waves, in dependency order. Paths are owned by the
implementation track unless noted.

| Wave | Work | Files | Requirements |
|---|---|---|---|
| 1 | Shared install conventions: target and option guards, error line, download and verification helpers with component tests for the component-verified codes; publish script | `features/src/**/install.sh` helpers, `features/test/`, `scripts/v1-05-features.sh` | R-DEPENDENCIES, AC-TOOLS/unsupported-*, AC-TOOLS/requires |
| 2 | Foundation Features: `python/uv`, `python/python`, `node/node`, `go/go`, `rust/rustc`, `zig/zig`, `latex/core`, `browser-testing/playwright` | `features/src/<id>/` | AC-TOOLS/present, version |
| 3 | Remaining tool Features, then atomic Features | `features/src/<id>/` | AC-TOOLS/present, version, isolated, depends, side-effects |
| 4 | Bundles; starter APT suites; inspection parity for `devcontainer.features` | `features/src/*/all/`, `templates/src/`, `spec/v1/examples/` rendering, `internal/inspection/` | AC-TOOLS/bundle, conflict, inspect |
| 5 | User documentation for selection, versions, removals and compatibility notes; v1 site build | `docs-site/content/docs/tool-catalog.md` | Roadmap docs obligation |

The acceptance-test track works on `tests/v1-05-tool-catalog` from the
contract commit, owns `tests/acceptance/v1-05/`, and covers every
`AC-TOOLS/*` item on both targets. The implementation track works on
`feat/v1-05-tool-features` and may reuse PR #493's inspection code where it
matches the contract.

## Review

Designated reviewer: TeamMember `sage` (agent), model `claude-opus-5-5`,
provider Anthropic. `route_task` proposed no TeamMember for this review.

| Round | Commit | Result |
|---|---|---|
| 1 | `a3850809` | Not test-ready: 11 blocking findings. Owner decisions in DEC-20261029_1409-DaringTrail; remaining findings fixed as contract defects. |
| 2 | `25c5a472` | Not test-ready: 10 of 11 round-1 findings resolved; 3 new blocking findings (Rust `PATH` syntax, sibling commands provided by rustup proxies and the TeX Live scheme, error-line Feature ID for dependencies). Final round; closed by the integration owner below. |

### Integration-owner closure of round 2

R-REVIEW allows two rounds, so the integration owner closed the round 2
findings against the baseline:

- Rust `containerEnv.PATH` uses `${PATH}`, the Feature `containerEnv` form.
- `rust/clippy`, `rust/rustfmt` and the `latex/*` collections list no
  `commands`; their `command` and `texlive` checks identify the component,
  because rustup proxies and the TeX Live scheme can provide the executables.
- An `AIBOX_FEATURE_UNSUPPORTED_TARGET` line may name any Feature in the
  dependency closure of the Feature under test.
- Non-blocking: the ledger audit field is renamed `audit-classification` and
  marked historical; `ai/claude` sets `DISABLE_UPDATES=1`, which the Claude
  Code environment variable reference documents as blocking all updates;
  Debian installs use `--no-install-recommends`; the starter suite change is
  worded as V1-05 work; `requires` counts as a dependency for isolation.
- Debian package versions were already resolved from snapshot
  `20261028T000000Z` for both architectures, not from live indexes.
- The pinned CLI `0.89.0` selects plain HTTP for a `localhost` registry
  (`src/spec-configuration/containerCollectionsOCI.ts`, the protocol
  selection for `localhost`) and its Feature option type declares `enum`
  without enforcing it. If integration shows otherwise, the failure is a
  contract defect under R-CYCLE step 8.

## Contract amendment 1

Both tracks reported contract questions after the freeze. Classified as
contract defects under R-CYCLE step 8 and fixed in amendment 1:

| Source | Defect | Change |
|---|---|---|
| Tests Q1 | Runtimes required by automatically installed Features were not selected | Selections add `requires` of the whole `dependsOn` closure |
| Tests Q3, implementation 1 | Absence could not observe Rust components; `texlive` absence would conflict with the TeX Live scheme | `command` checks are absence observables; `version` and `texlive` checks are not |
| Tests Q4, Q5 | Dependency closure and `alsoProvides` scope were unclear | Transitive closure; exemptions cover only what they name; new `feature` kind exempts `rust/x86-64-cross` on amd64 |
| Tests Q11, Q12 | Git status needs a repository the base fixture lacks; metadata keys were open | Workspace folder unchanged; exact metadata key list |
| Tests Q6, Q15 | Publish prerequisites and registry image unstated | `devcontainer` and `node` from `PATH`; registry chosen by tests and recorded |
| Tests Q2, Q7, Q8, Q9, Q13 | Test naming, bundle presence, conflict value, installed package | ID form, stack IDs, bundles excluded from `present`, first non-default value, `dpkg` status rule |
| Implementation 2 | The pinned CLI also publishes `latest` | Tags are a minimum; extra tags are never selected |
| Implementation 3 | `versionCoupling` table direction | Key is this Feature's version; missing coupled Feature fails |
| Implementation 4 | Upstream Features in `dependsOn` | Use `upstream.ref` and `upstream.options` |

Amendment review (sage, `claude-opus-5-5`, Anthropic) of `16a8b73d` found one
blocking defect: the Dev Container CLI writes
`.devcontainer/devcontainer-lock.json` into the workspace, so the workspace
rule failed by construction. The rule now exempts that file. Non-blocking:
per-item test-ID slots, `AC-TOOLS/absent` wording aligned with chapter 8,
metadata read from the test registry, and the `requires` closure repeated to a
fixed point.
Round 2 of the amendment review approved `4ec6dc8e` with no blocking findings;
its cosmetic test-ID note is applied.

Tests Q10 (MCP inspection), Q14 (requires target) and the implementation's
error-code mapping for apt update and tool-internal downloads follow the
existing contract and need no change.

## Contract amendment 2

| Source | Defect | Change |
|---|---|---|
| Implementation (wave 3) | `rich-cli` 1.8.1 prints version 1.8.0, so `AC-TOOLS/present preview/enhanced/rich` cannot pass; with uv-managed Python the rich library cannot serve aibox's Yazi preview glue | Owner removed rich and added `preview/enhanced/glow` 3.0.0 and `preview/enhanced/bat` 0.26.1 (DEC-20261029_2107-QuietSeal); census row removed; `COMPAT-RICH` updated |
| Tests A1 | Whether an exempt command name also exempts a `command` check was unstated | It exempts only the matching `commands` entry |

glow 3.0.0 and bat 0.26.1 were checked on 2026-10-29: archive checksums agree
with the GitHub asset digests and glow's `checksums.txt`, and both arm64
binaries print the catalog version; amd64 binaries are checked by hash only
until the host run. glow drops front matter, so previews show it through bat.
The v0 base installed bat; if a later base revision restores it, `bat` must be
added to `baseline.provides` (V1-13 hand-off).

Amendment review (sage, `claude-opus-5-5`, Anthropic) approved `48f54282` with
no blocking findings. Its notes are applied: the exempt-command sentence covers
`baseline.provides` and `alsoProvides`, chapter 4 marks glow as new, and
`COMPAT-V0-REQUIRES` no longer lists the removed rich row.

## Contract amendment 3

The nested Podman light pass on `debian-13-arm64` found two smoke commands that
cannot pass with any correct implementation (DEC-20261031_1743-MightyRaven):

| Source | Defect | Change |
|---|---|---|
| Light pass | mdbook resolves `--dest-dir` from the working directory, so `../mdbook-book` points outside the workspace | `SMOKE-MDBOOK` uses `--dest-dir mdbook-book` |
| Re-run after the node fix (`rerun2-arm64-211050`); the light pass failed earlier on the same test because of the missing `libatomic1` | `browser-testing/chromium` installs with `--no-shell`, but `playwright screenshot` defaults to the headless shell | `SMOKE-CHROMIUM` adds `--channel chromium`; the Feature keeps `--no-shell`; `COMPAT-CHROMIUM-CHANNEL` documents the channel for users |

Each corrected command passed on `debian-13-arm64` in a temporary catalog edit
under rootless Podman, from the dirty tree at `3ca64b86` that also carried the
Feature fixes (`mdbook-try-*` and `chromium-try-*`, in the local test cache only).
Neither has run on `debian-13-amd64` or Docker. The defects are not attributed
to the engine because of how the tools resolve these paths, not because of a
Docker run.

Amendment review (sage, `claude-opus-5-5`, Anthropic) approved `f01b7111` with
no blocking findings. Its notes are applied: the compatibility note's examples
are complete and name the `use` option and `chromium.launch`, and this section
states the evidence precisely. Integration must regenerate the catalog reference
when it moves onto this commit, and changes the page's wording so that a note
may record behaviour unchanged from v0.

## Contract amendment 4

The owner asked for shorter AC-TOOLS runs (DEC-20261031_1856-CuriousBeacon). A
full target planned 466 builds, 322 of them expected failures that mostly
exercise the shared guard code.

| Item | Full scope | Sampled scope |
|---|---|---|
| `requires` | 19 | 3 |
| `unsupported-version` | 65 | 9 |
| `unsupported-target` | 234 | 24 |
| Container builds per target | 466 | 184 |

`AC-TOOLS/guards` adds the component tests to each run, so every aibox Feature
still has offline evidence for these error codes. What sampled scope gives up
is real-container failure evidence for each Feature, which full scope provides
on the owner's command.

Feature install steps never hit the layer cache on Podman: 0 of 268 real build
logs show a hit, because the Dev Container CLI's steps carry a digest that differs
per build. Reusing `latex/core` layers is therefore not available through the
pinned CLI on Podman and is not implemented, which departs from
DEC-20261031_1856-CuriousBeacon. A small Podman 5.4.2 build did reuse a
`RUN --mount` layer built from other mounted content, so the runner fails a build
that takes a Feature step from the cache. Parallel Podman builds use separate
stores so that removing one build's image cannot remove layers another has just
matched; that race is reasoned, not observed.

## Integration loop

All runs used rootless Podman 5.4.2 on `debian-13-arm64`. Evidence outside the
repository stays in the test cache; the passing runs are copied to
[`spec/v1/evidence/V1-05-fbd178d6/`](../spec/v1/evidence/V1-05-fbd178d6/).

| Run | Source | Result | Classification |
|---|---|---|---|
| Light pass and reruns, 2026-10-30 | `3ca64b86` (dirty) | Two smoke commands could not pass | Contract defects: amendment 3 |
| `full-arm64-20261031T182046Z` | `60567029` | 872 of 875 passed; `go/goreleaser` `present`, `isolated` and `side-effects` failed | Infrastructure failure: go.dev dropped an HTTP/2 stream (curl exit 92) while `go/go` downloaded `go1.27.0.linux-arm64.tar.gz`. `curl --retry` does not retry that error. Hardened with `--retry-all-errors` (`8031b4ea`); regression test `fbd178d6` |
| `par4-*` and `shared4-*`, `--jobs 4` | `b73bd368`, `df2cac38` | 101 to 443 passed per run; failures were `devcontainer up` errors, `EAGAIN` and incomplete probes | Infrastructure failure in the acceptance harness: workers set `TMPDIR`, and with `XDG_RUNTIME_DIR` unset rootless Podman keeps runtime state there, so `podman rm` could not reach containers the Dev Container CLI started. Each leaked about 27 tasks until the cgroup's 1200-task limit was reached. Fixed by the integration owner in `71124b37` (BACK-20261101_1120-SmoothRiver), with unit tests: workers pin the engine's runtime directory, and a run stops before its next build when a container cannot be removed or tasks reach 80% of `pids.max` |
| `leakfix4-full-arm64-20261102T181158Z`, `--scope full --jobs 4` | `71124b37` | 876 of 876 passed, no cleanup failures, 51 minutes; peak 429 tasks | Passed |
| `leakfix4-sampled-arm64-20261102T181158Z`, `--jobs 4` | `71124b37` | 594 of 594 passed, no cleanup failures, 49 minutes | Passed |
| `guards-arm64-fbd178d6`, `--item guards` | `fbd178d6` | 1 of 1 passed | Diagnostic: `fbd178d6` changes only `features/test/`, which `AC-TOOLS/guards` runs |

`fbd178d6` changes no Feature, script or harness file relative to `71124b37`,
so the arm64 runs at `71124b37` stand for it.

## Host handoff: debian-13-amd64

Chapter 14 requires a human-authorized operator for host gates. Acceptance needs
a sampled run on `debian-13-amd64`; it is also the first run on Docker.

- **Commit:** `fbd178d6` on `integration/v1-05-tool-catalog`, from a clean
  checkout. Later commits on the branch change documentation only.
- **Host:** an amd64 machine with a working `docker` CLI and daemon (daemon
  architecture `amd64`), Python 3.10 or later, Git, port 5000 free on
  localhost, and about 80 GB of free image space (the arm64 full run with
  `--jobs 4` peaked at about 70 GB).
- **Network:** Docker Hub, ghcr.io, snapshot.debian.org, the producer download
  sites in the catalog, and nodejs.org.
- **Credentials:** none. If Docker Hub rate-limits anonymous pulls, `docker
  login` with any account first.
- **Command:**

  ```sh
  git fetch origin && git checkout --detach fbd178d6
  tests/acceptance/v1-05/run-host.sh debian-13-amd64
  ```

  Add `--scope full` for the optional full run. `--jobs` defaults to at most 4
  from the host's CPUs and memory. Parallel builds have not yet run on Docker;
  if failures look like resource or engine errors, rerun with `--jobs 1`.
- **Effects:** starts a registry container on `localhost:5000`, builds the
  Go binary in `golang:1.27.0`, and builds about 184 workspace images (466
  with `--scope full`). It removes the run's containers, images, workspaces
  and registry at the end; pulled base images stay in the Docker cache.
- **Duration:** about an hour on 8 CPUs.
- **Return:** the evidence directory the run prints,
  `$HOME/aibox-v1-05-evidence/debian-13-amd64-<UTC time>`, as a whole. The
  exit status is 0 only when every test passed and cleanup succeeded.

## Logbook

| Date | Entry |
|---|---|
| 2026-10-29 | Contract drafted on `contract/v1-05-tool-catalog` (`a3850809`). |
| 2026-10-29 | Review round 1 findings resolved (`25c5a472`). |
| 2026-10-29 | Review round 2 findings closed by the integration owner (`34d06da0`). |
| 2026-10-29 | Owner froze the contract at `34d06da0` (DEC-20261029_1515-ShinySpire). |
| 2026-10-29 | Tracks started: tests `799d234b`..`b646e40a`, implementation waves 1–2 `2cfa66e5`..`977a7eed`. Amendment 1 drafted from their contract questions; review pending. |
| 2026-10-29 | Owner approved amendment 1 as contract commit `59c35acc1cf230c81658a5763205b1f91a0b1092` (DEC-20261029_1725-AmberBeam); both tracks move onto it. |
| 2026-10-29 | Tracks reported: tests `85d77f02`, implementation wave 3 `5b38f55d`. Amendment 2 drafted (rich replaced by glow and bat); review pending. |
| 2026-10-29 | Implementation wave 4 reported (`13997c2a`). Owner approved amendment 2 as contract commit `d3f6e08eaf232eb6aa7c77cf6c6d2fcea10d1c3c` (DEC-20261029_2134-NobleBlossom); both tracks move onto it. |
| 2026-10-30 | Nested Podman light pass on `debian-13-arm64`; amendment 3 drafted from its two smoke defects. |
| 2026-10-31 | Owner approved amendment 3 (DEC-20261031_1743-MightyRaven); sage approved `f01b7111` with no blocking findings; frozen as contract commit `d1c7adb487bd4cfa7449550ffcd3da05153537ed` (DEC-20261031_1750-NeatDaisy). |
| 2026-10-31 | Full arm64 run at `60567029`: 872 of 875 passed; the 3 `go/goreleaser` failures were a dropped download (infrastructure), hardened in `8031b4ea`. Parallel runs exhausted the cgroup task limit (BACK-20261101_1120-SmoothRiver). |
| 2026-11-02 | Runner leak fixed (`71124b37`). Full (876/876) and sampled (594/594) arm64 runs with `--jobs 4` passed. |
| 2026-11-03 | Regression test for interrupted downloads (`fbd178d6`); `AC-TOOLS/guards` passed on arm64. Owner chose to run `debian-13-amd64` on an amd64 Docker host; handoff above. |
| 2026-11-05 | arc42 maintainer architecture documentation added under `docs-site/content/docs/architecture/` (R-CODE). |
