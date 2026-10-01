# V1-04 upstream Dev Container CLI archive

`devcontainer-cli-0.89.0.tgz` is the published `@devcontainers/cli` 0.89.0
archive from `https://registry.npmjs.org/@devcontainers/cli/-/cli-0.89.0.tgz`.
Its SHA256 is
`49c7d71d40058f89e1fd8b019a193ed4215b7fc773c0f6273f7032a46cd33f4b`.
The package includes its MIT `LICENSE.txt` and `ThirdPartyNotices.txt`.

The V1-04 host verifier checks this digest and package metadata, then invokes
the upstream CLI through a temporary Node executable. It does not call npm or
contact the npm registry. Node 20.19.5 archives come from `nodejs.org` and
are checked against the fixed platform SHA256 values in
`scripts/verify-v1-04-host.py`. The temporary CLI and Node files are removed
after the run; the host's tool installation is unchanged.

Updating this archive requires updating the version, digest, Node
compatibility checks, offline verifier, and host evidence contract together.
