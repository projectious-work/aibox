# V1-04 native workspace starters

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
| Template rendering and overwrite refusal | The implementation agent reports the source/helper checks passed. Attach candidate-specific offline verification output after the final revision is fixed. |
| Rendered example/schema checks | Pending final candidate verification. |
| Pinned Dev Container CLI configuration parsing | Pending final candidate verification. |
| Direct upstream `build`, `up` and `exec` | **Pending host check.** Run `./scripts/verify-v1-04-host.sh --output /path/to/evidence`; Python verifies a checked-in Dev Container CLI `0.89.0` archive and downloads a checksummed temporary Node runtime. No npm command or npm registry access is needed on the host. The script records lifecycle and home-volume persistence results. |
| Go CLI/MCP parity | The maintained V1-03 `inspect_workspace` operation remains read-only; V1-04 adds no new operation. Capture parity against the rendered starter in the phase evidence. |
| Hugo documentation build | Passed locally with Hugo Extended v0.167.0: 35 pages and 142 static files. |

No v1 release or published Template reference exists. Template registry
publication and direct lifecycle runtime qualification are not implied by local
source templates or schema parsing. The roadmap points to this note as the
in-progress record; replace that reference with a checked-in evidence record
when the candidate revision and local gates are finalized. Append the host
evidence path and tested runtime after the host gate completes.

## Remaining evidence

Before setting the roadmap item to `done`, run the host gate above and attach
its `evidence.json` and logs, finish rendered-example and pinned-CLI checks,
and record successful CLI/MCP inspection parity. This dev environment has no
container runtime, so the direct lifecycle result must come from the supported
host check. A future `shipped` status still requires published release and
Template artifact evidence.
