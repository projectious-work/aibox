# V1-04 native workspace starters

V1-04 implementation is done with local evidence; publication is pending.
The verified implementation commit is
`85cc0c9af67e74175b0efc1499d48193bb73a629`. The repeat host run used
`667d1fdd762afa940b8c26c40c6e7fe6e5cc2acd`, whose changes from the
implementation commit cover only evidence, documentation, and host-verifier
progress messages. The starter and Go implementation files are identical.

## Scope

V1-04 defines minimal and curated native Dev Container starters, persistent
home storage, custom-user and bind-mounted-home examples, and the ownership
boundary between native configuration, Features, and aibox workspace UX. It
retains the V1-03 read-only Go CLI and MCP inspection journey. It does not add
container lifecycle operations to aibox.

## Candidate implementation

- Source Templates live at `templates/src/minimal` and
  `templates/src/curated`. The local renderer applies Template options and
  refuses to overwrite an existing destination.
- Rendered examples live under `spec/v1/examples/{minimal,customized,
  custom-user,bind-home}`. The minimal example has no active aibox extension;
  curated is explicit about optional Features and UX metadata.
- Standard home persistence uses a project-scoped named volume mounted at the
  selected container user's home. The bind example uses ignored `.aibox-home`
  and requires its source directory on the daemon host. The custom-user example
  aligns the image user, remote user, home path and mount.
- User guidance is in the [native workspace starter guide](../docs-site/content/docs/native-starters.md).

## Verification record

| Check | Result |
|---|---|
| Template rendering and overwrite refusal | Passed in the [offline verification](../spec/v1/evidence/V1-04-85cc0c9a/offline-evidence.json) for source `85cc0c9af67e74175b0efc1499d48193bb73a629`. Both Template variants and `codex`/`none` options were checked. |
| Rendered example/schema checks | Passed for minimal, customized, custom-user, and bind-home in the same offline verification. |
| Pinned Dev Container CLI configuration parsing | Passed in the [first host run](../spec/v1/evidence/V1-04-85cc0c9a/host-run-1.json) for all four examples with CLI `0.89.0` on macOS arm64 and Docker. |
| Direct upstream `build`, `up` and `exec` | The [first](../spec/v1/evidence/V1-04-85cc0c9a/host-run-1.json) and [repeat](../spec/v1/evidence/V1-04-85cc0c9a/host-run-2.json) macOS arm64 Docker host runs each passed all 34 recorded steps and had zero cleanup failures. Both checked minimal, curated, custom-user, and bind-home; no-cache rebuild and home-marker persistence passed for the two Template variants. Each manifest records 43 checksummed artifacts. |
| Go CLI/MCP parity | Passed against all four starters; see the [MCP parity transcript](../spec/v1/evidence/V1-04-85cc0c9a/mcp-parity.txt). The maintained `inspect_workspace` operation remains read-only. |
| Hugo documentation build | Passed locally with Hugo Extended v0.167.0: 35 pages and 142 static files, including the [latest offline verification](../spec/v1/evidence/V1-04-85cc0c9a/offline-evidence-667d1fdd.json). |

No v1 release or published Template reference exists. Template registry
publication is not implied by local source Templates. Full logs for the two
passing runs remain in the owner-retained host evidence directories
`/tmp/aibox-v1-04-evidence-20261001-192444-13450` and
`/tmp/aibox-v1-04-repeat-20261001-213052-13450`; their checked-in manifests
record each log's digest. The host script uses a checked-in, checksummed Dev
Container CLI archive and a temporary checksummed Node runtime, without npm.

## Completion boundary

The local acceptance boundary is met by the two direct upstream host runs,
source and rendered Template checks, read-only Go CLI/MCP parity, and v1 Hugo
documentation build. The roadmap records `done` after this evidence is merged
to `v1.x-dev`. A future `shipped` status still requires published release and
Template artifact evidence.
