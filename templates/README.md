# Native workspace Templates

`src/minimal` and `src/curated` are individually versioned Dev Container
Template sources (`1.0.0`). They follow the [upstream Template format](https://github.com/devcontainers/spec/blob/main/docs/specs/devcontainer-templates.md).
Their source distribution is checked in; no OCI Template or aibox runtime
Feature is published by this implementation. No GitHub workflow is required.

Apply local sources into an existing project that has no `.devcontainer`:

```sh
node templates/render.mjs --template minimal --output /absolute/path/to/project
node templates/render.mjs --template curated --output /absolute/path/to/another-project --harness none
```

The helper copies initial `.devcontainer` files with standard Template option
substitution. It refuses existing configuration. The files then belong to the
project; editing or rebuilding never re-applies a Template automatically.
For an OCI distribution, use the pinned upstream CLI's `templates publish templates/src` command
with an approved registry/namespace as a separate publication task. The CLI
packages sources during publication; 0.89.0 has no standalone Template
`package` command. `devcontainer templates apply` consumes published OCI refs;
there is no fabricated registry reference for these source Templates.

Both Templates build Debian with bash, git, tmux and Yazi. The `harness` option
accepts `codex` (default: checksummed native Codex 0.107.0) or `none`. This is
installation only; no harness is started and no credentials are transferred.
Authenticate explicitly inside the container when wanted. Complete optional
tool/harness selection and runtime UX integration belong to later roadmap phases.

The default user is `aibox`, with `containerUser` and `remoteUser` both `aibox`.
Its named volume is `aibox-home-${devcontainerId}` mounted at `/home/aibox`.
This retains user configuration and private login/cache state across rebuilds.
No host home is copied or mounted, and no `.aibox/` or `.aibox-home/` is created.
The curated Template adds `/etc/tmux.conf` with mouse support, bounded history
and a small neutral status line. tmux loads user `~/.tmux.conf` afterwards;
user overrides remain authoritative even when the home volume masks image files.
Minimal has no active aibox extension. Curated is also directly usable without
aibox UX metadata or an aibox runtime Feature.

Native lifecycle verification uses the pinned upstream Dev Container CLI:

```sh
devcontainer build --workspace-folder /absolute/path/to/project
devcontainer up --workspace-folder /absolute/path/to/project
devcontainer exec --workspace-folder /absolute/path/to/project bash -lc 'id; tmux -V; yazi --version; codex --version'
```

Use `tmux -V; yazi --version` when the harness is `none`. See the V1-04 host
verification script for direct build/up/exec, named-volume identity and rebuild
persistence evidence. Its Python bootstrap uses a checked-in upstream CLI
archive and a checksummed temporary Node runtime; it needs no npm installation
or npm registry access on the host. Source/metadata checks do not establish
live runtime qualification.

## Optional native Features

Feature installation is a native project edit. For example, the published
[upstream Node Feature](https://github.com/devcontainers/features/tree/main/src/node)
can be selected using this inspected immutable manifest (Feature 1.7.1):

```jsonc
"features": {
  "ghcr.io/devcontainers/features/node@sha256:8c0de46939b61958041700ee89e3493f3b2e4131a06dc46b4d9423427d06e5f6": {
    "version": "22.14.0",
    "nodeGypDependencies": false,
    "pnpmVersion": "none",
    "nvmVersion": "0.40.3",
    "installYarnUsingApt": false
  }
}
```

The manifest and selected option names were checked against the published OCI
artifact, and Node/nvm versions exist upstream. This optional composition is
not part of the default Template or V1-04 lifecycle gate; tool catalog and
optional Feature qualification happen in V1-05. The upstream CLI owns Feature
lockfiles; aibox adds no alternative Feature resolver or lockfile. Rebuild when
changing Feature installation/version options.

For a custom user, change Dockerfile `ARG USERNAME`, `containerUser`,
`remoteUser` and the home mount target together. The rendered
`spec/v1/examples/custom-user` uses `dev` and `/home/dev`. The
`spec/v1/examples/bind-home` example replaces the volume with an explicitly
selected `${localWorkspaceFolder}/.aibox-home` bind and ignores it in Git;
create the source on the daemon host first and check writable ownership.
Bind mounts have different portability and cleanup behavior. Neither example
copies a whole host home or silently changes host ownership.

`provenance.json` records base/index platform digests, immutable Debian
snapshot/package versions and checksummed tool assets for linux/amd64 and
linux/arm64. APT verifies signed repository metadata. The initial snapshot
transport is HTTP so the minimal base can bootstrap CA certificates; downloaded
binaries use HTTPS plus SHA256 verification. Maintainers own pin updates and
requalification. These platform pins are available upstream; live verification
covers the architectures recorded in its evidence, not every advertised asset.
