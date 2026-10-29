# V1-05 AC-TOOLS acceptance suite: host run

Black-box acceptance tests for every `AC-TOOLS/*` item in
[chapter 5](../../../spec/v1/05-acceptance.md#ac-tools--complete-optional-tool-matrix),
planned from [the catalog](../../../spec/v1/catalog/v1-05-features.json) and
evaluated as [chapter 8](../../../spec/v1/08-native-artifacts.md#tool-feature-interface)
defines. Acceptance needs one full run on each supported target.

## Prerequisites

- **Docker**: a working `docker` CLI and daemon whose architecture is the
  target you run (`amd64` for `debian-13-amd64`, `arm64` for
  `debian-13-arm64`). Allow for a long run (about 466 container builds per
  target) and tens of gigabytes of image space. Port 5000 on localhost must be
  free.
- **Python 3.10 or later.**
- **Dev Container CLI 0.89.0 and Node.** The runner starts them itself by
  default. It unpacks `tools/vendor/devcontainer-cli-0.89.0.tgz` and downloads
  Node 20.19.5, checking both SHA-256 digests the same way as the V1-04 host
  verifier. It then puts them first on `PATH` as `devcontainer` and `node` for
  the publish script. To use your own installation instead, pass
  `--cli-from-path`. The runner checks that `devcontainer --version` prints
  `0.89.0`.
- **Registry image**: `registry:2.8.3` by default. Use `--registry-image` for
  another OCI Distribution image. The image and its digest are recorded in the
  evidence.
- **aibox binary for `AC-TOOLS/inspect`**: built from this checkout in
  `golang:<go.mod version>` unless you pass `--aibox PATH`.
- **Publish script**: `scripts/v1-05-features.sh` from the implementation
  branch. The run fails if it is missing.
- **Network access**: Docker Hub (base images), ghcr.io (upstream Features),
  snapshot.debian.org, the producer download sites in the catalog, and
  nodejs.org when the CLI is bootstrapped.

## Run

From a clean checkout of the integration branch, run one command per
architecture:

```sh
tests/acceptance/v1-05/run-host.sh debian-13-arm64   # on an arm64 Docker host
tests/acceptance/v1-05/run-host.sh debian-13-amd64   # on an amd64 Docker host
```

Each run does the following:

1. Starts a registry on `localhost:5000`.
2. Runs `scripts/v1-05-features.sh publish --registry localhost:5000/aibox-features`.
3. Builds each planned workspace from `spec/v1/fixtures/v1-05-base`, or from
   an unsupported-target fixture.
4. Runs the checks as the remote user (`devcontainer exec`) and as root
   (`docker exec --user root`).
5. Removes that run's containers, images, workspaces and registry.

The exit status is 0 only when every planned test passed and cleanup
succeeded. If the selected target does not match the Docker daemon, the run
is refused.

## Evidence

Each run writes a new directory, by default
`$HOME/aibox-v1-05-evidence/<target>-<UTC time>`; `--output ABS_DIR` sets
another location. The directory contains:

| File | Content |
|---|---|
| `evidence.json` | Source commit, host, target, Docker version and architecture, CLI and Node runtime, registry image and digest, publish result, and one result per test ID |
| `summary.txt` | Counts per item and every test that did not pass |
| `cli-runtime.json` | Pinned CLI and Node identity |
| `publish.log`, `host.log` | Publish and setup transcripts |
| `builds/<key>.log`, `tests/*.log` | Per-build Dev Container CLI, probe and inspection transcripts |

## Partial runs and planning

```sh
tests/acceptance/v1-05/run.sh --list --target debian-13-amd64      # every planned test ID
tests/acceptance/v1-05/run.sh --dry-run --item stack                # builds with their Feature selections
tests/acceptance/v1-05/run-host.sh debian-13-arm64 --item present,version --feature 'python/*'
tests/acceptance/v1-05/run.sh --item trace                          # offline, no Docker
```

`--item` takes AC-TOOLS items, `--feature` takes glob patterns for Feature or
smoke IDs, and `--test` takes glob patterns for full test IDs. A partial run
is diagnostic only; acceptance needs the full run.

## Unit tests

```sh
cd tests/acceptance/v1-05 && python3 -m unittest discover -s unit -t .
```
