# V1-06 layered settings and diagnostic sinks

V1-06 implements the frozen specification `settings-and-logging@v1`
([`spec/settings-and-logging/v1/`](../spec/settings-and-logging/v1/), tag
`spec/settings-and-logging/v1`, contract commit
`7c2d4b89566d1869d720cd8b23a5afaae90ce68f`). Implementation pull requests
carry `Spec: settings-and-logging@v1` and never change `spec/`.

Status: plan reviewed; tracks starting. The cycle is the one in chapter 21 of the v1 product
description (development cycle steps 3 to 10).

## Roles

| Role | Actor |
|---|---|
| Owner | Repository owner; approved the contract on 2026-11-11 |
| Integration owner | `claude-opus-5-5` (Anthropic), the primary session |
| Designated reviewer (plan review and verification) | `sage`: a separate agent session, model `claude-opus-5-5` (Anthropic), that gets only the contract, this note and the repository, not the author's reasoning |
| Acceptance-test track | A fresh agent session with no conversation history, on `feat/settings-and-logging-tests` |
| Implementation track | A separate agent session on `feat/settings-and-logging-impl` |

Both tracks branch from the phase branch `feat/settings-and-logging` at the
commit that adds this note. That commit holds the frozen contract unchanged.
The phase branch squash-merges into `main` when every acceptance test passes
and verification is closed.

## Inventory and trace

Every normative requirement of the contract, with the work that delivers it
and how it is verified. **J** is a black-box journey from `journeys.md` on the
tests track; **C** is a component test on the implementation track, for the
reason the contract gives. Requirement levels are the contract's.

| Requirements | Delivered by (wave, package) | Verification |
|---|---|---|
| CFG-001 to CFG-009 (layers, precedence, defaults, no writes, resolution before the operation, CLI and MCP alike, failure order, `version` and `help` untouched) | W1 `internal/settings` resolver; W4 wiring in `cmd/aibox` and `internal/mcpserver` | J: J-CFG-01, 02, 03, 04, 05, 15 (h), 17, 23 (a). C: resolver unit tests per key and layer pair |
| CFG-010 to CFG-019, CFG-024 (system and user files, XDG rule, absent files, read and schema defects, duplicate keys, null, `--settings-file`, no project or cwd settings file, no home) | W1 `internal/settings` file loader with a duplicate-key-aware JSON decoder and the process-settings schema | J: J-CFG-06, 07, 08, 09, 23 (b) |
| CFG-020 to CFG-022 (project layer from `customizations.aibox.workspace` only) | W1 `internal/settings`, reading the selected `devcontainer.json` through `internal/project` | J: J-CFG-09, 10 |
| CFG-030 to CFG-044 except CFG-037 (retired) (env files, `AIBOX_ENV_FILE`, grammar, process environment mapping, registered names, nearest name, conversion, limit of files) | W1 `internal/settings` env-file parser and environment mapper; registered names from `fixtures/registered-environment.json` compiled in | J: J-CFG-11, 12, 13, 14. CFG-043: C (no child process in V1-06) |
| CFG-050 to CFG-057 (flags, ranges, `--project` anchoring, no chdir, `mcp serve` and `plain`, fixed MCP settings, usage errors, `mcp serve` bootstrap failure) | W4 `cmd/aibox` command-line parser (replaces `flag`, which accepts repeated flags and cannot report `usage_error` by the contract's rules); `internal/mcpserver` | J: J-CFG-15, 16, 18, 19. CFG-053: C |
| CFG-060 to CFG-068 (effective-configuration view, `overridden`, `settingsSources`, redaction in the view, CLI and MCP parity, no hooks or engine, view only on request, schema validity) | W3 `internal/contract` result fields `settingsSources` and `overridden`; W4 `internal/usecase` inspect with `effectiveConfig`; MCP input `effectiveConfig` | J: J-CFG-01, 17, 20, 21, 22, 24 |
| OUT-001 to OUT-008 (stdout only the result, JSON, YAML, human text, colour rule, pre-operation error envelope, bootstrap diagnostic, exit codes) | W3 `internal/output`: JSON, block YAML without a new dependency (D7), human summary, colour rule, bootstrap diagnostic | J: J-OUT-01 to 05, J-CFG-15 (g), 18 |
| LOG-001 to LOG-008 (event lists and order, schema, request ID and sequence, process correlation ID, level filter, `settings.rejected`, nothing on stdout, cancellation) | W2 `internal/logging` events, levels, plain and JSON Lines formats; W4 emission points in the use case, CLI and MCP | J: J-LOG-01 to 04, every rejection journey. LOG-008: C |
| LOG-010 to LOG-020 (sinks, file anchoring, `path_escape`, modes, append and rotation, retention, concurrent writers, no foreign files, insecure mode) | W2 `internal/logging` stderr sink and file sink with lock file, rotation and retention | J: J-LOG-03 to 10 |
| LOG-030 to LOG-034 (degraded sink, `logging.required`, closed stderr, bounded queue, flush limit) | W2 `internal/logging` queue and failure handling; W3 warning `log_sink_degraded` | J: J-LOG-10, 11. C: LOG-030 write failure, LOG-031 mutating part, LOG-033, LOG-034 |
| RED-001 to RED-005 (known secrets, secret shapes, home prefix, split writes, no way to disable) | W2 `internal/redact`, applied before any sink or queue and in the view | J: J-CFG-22, J-LOG-12. RED-004: C |
| EXE-001 to EXE-003 (read-class time limit, timeout outcome, limit in the view) | W4 `internal/usecase` context deadline from `execution.timeoutSeconds` | J: J-CFG-01 (EXE-003). EXE-001, EXE-002: C |

No requirement is `not applicable`, `blocked` or `deferred` beyond what the
contract itself places in later phases (overview, "Out of scope").

## Plan

### Waves

The implementation track works test-first (R-TESTSTYLE) in this order. A wave
starts when the interfaces it consumes are stable.

1. **W1 settings core** (`internal/settings`): the key catalog of
   `interfaces.json` as Go data (15 keys, 7 layers, `env` names, defaults per
   mode, ranges and values), file loading and schema checks, env-file grammar,
   environment mapping with registered names and nearest-name suggestion,
   conversion, failure order, and the resolved value with its source and the
   overridden values. Pure functions over an injected file system view and
   environment, so tests need no real home.
   - The home directory comes from `HOME` only, with no `os/user` fallback
     (CFG-024, J-CFG-23 (b)); the binary stays static.
   - Settings files are checked by a hand-written validator that maps each
     defect to the reason `settings.schemaDefects` assigns; a generic
     JSON Schema validator does not give those reasons.
     `github.com/google/jsonschema-go` stays an indirect dependency.
   - The resolver never changes the process environment (no `os.Setenv`), and
     keys from an env file never enter the environment given to the
     process-runner port. That is the component test for CFG-043.
   - `internal/settings` owns the path rules shared with W2: the XDG state
     directory, `logging.fileSink.stateDirectory`, the permitted root and the
     anchoring of `logging.file`. W2 calls them; nothing writes them twice.
2. **W2 logging and redaction** (`internal/logging`, `internal/redact`):
   event model, level filter, plain and JSON Lines encoders, bounded queue
   with drop and direct-write rules, stderr sink that tolerates a closed
   reader, file sink (anchoring through W1, permitted root, modes, append,
   lock file, rotation, retention), and redaction of known secrets, secret
   shapes and the home prefix.
   - Events are validated against `log-event.schema.json` in tests, not at
     run time.
   - The lock is `<file>.lock` at mode 0600, held with `flock` as in
     `internal/operation/lock_unix.go`. It is not counted by and never deleted
     by retention. No descriptor of the log file is kept between events; a
     write waits for the lock at most as long as LOG-034's `flushSeconds`
     allows, and otherwise counts as a sink failure. A component test runs 16
     writers contending for the lock and shows that all finish within that
     limit and none reaches the failure path (LOG-018).
3. **W3 result and output** (`internal/contract`, `internal/output`): result
   fields `settingsSources` and `overridden`, warning `log_sink_degraded`,
   pre-operation error envelope (`invalid_request`), the operation-result
   schema of this version, and the JSON, YAML and human encoders with the
   colour rule and the bootstrap diagnostic.
4. **W4 wiring** (`cmd/aibox`, `internal/usecase`, `internal/inspection`,
   `internal/mcpserver`):
   - A command-line parser that follows `cli.flags`, `cli.syntax` and
     `settings.selectors`. The parser reports only `usage_error` (CFG-056,
     CFG-044); flag values are checked by the resolver at the invocation layer,
     so a defective user file is reported before `--log-level loud`
     (J-CFG-15 (h)). From a rejected command line the parser still extracts
     `--log-format` and `--log-level` for `logging.rejectionFormat`.
   - Settings resolution before the operation; `inspect --effective-config`;
     MCP input `effectiveConfig`; settings fixed at `mcp serve` start; event
     emission in the orders of `logging.eventOrder`; exit codes.
   - Time limits and cancellation: a context is passed into
     `internal/inspection` (today it takes none), with the deadline from
     `execution.timeoutSeconds` and SIGINT and SIGTERM handling (EXE-001,
     EXE-002, LOG-008 as component tests).
   - The mutating part of LOG-031 is a component test with a test double for
     an operation of class mutating behind the use-case port. It asserts
     `operation_record_unavailable` (category `internal`, exit 5) before any
     side effect. No placeholder operation is registered (R-CODE).
5. **W5 documentation and architecture**:
   - The v1 site: settings hierarchy, environment and flag precedence,
     logging sinks, redaction and troubleshooting reference, including
     capturing the stderr of `mcp serve` under systemd and launchd (overview,
     "Out of scope").
   - arc42 pages under `docs-site/content/docs/architecture/`: building
     blocks (`internal/settings`, `internal/logging`, `internal/redact`,
     `internal/output`), concepts (logging and redaction), decisions, risks.
   - ADRs under `docs/adr/`: replacing the standard `flag` parser, and the
     hand-written YAML encoder (D7).

Data that comes from the contract (the registered environment names, any
schema copy) is compiled into Go, because `go:embed` cannot reach `spec/`. A
drift test reads `spec/settings-and-logging/v1/fixtures/registered-environment.json`
and `schemas/` when the tests run and fails on any difference.

### File ownership

| Paths | Owner |
|---|---|
| `internal/**`, `cmd/**`, `go.mod`, `go.sum` | Implementation track |
| `docs-site/content/**` pages for settings and logging, `docs-site/content/docs/architecture/**`, `docs/adr/**` | Implementation track (W5) |
| `tests/acceptance/v1-06/**` | Acceptance-test track |
| `scripts/verify-v1-06.sh`, `dev-notes/V1-06-settings-and-logging.md`, the V1-06 entry of `roadmap.yaml` | Integration owner |
| `spec/settings-and-logging/v1/**` and the rest of `spec/` | Nobody during the phase, (the roadmap is now `roadmap.yaml` at the repository root, outside `spec/`) |

The tests track does not read `internal/**`, `cmd/**` or any implementation
branch. The implementation track does not read `tests/acceptance/v1-06/**` or
the tests branch.

### Frozen shared interfaces

Both tracks rely only on what the contract defines: everything in
`interfaces.json`, the schemas and the fixtures of
`spec/settings-and-logging/v1/`. That includes the CLI flags, syntax and exit
codes, the MCP tool `inspect_workspace` and its input, result documents on
stdout, events on stderr, the system and user settings paths and the XDG
rules, the env-file grammar, the `AIBOX_*` names, the project members of
`devcontainer.json`, and the directories `<state>/aibox/logs/` and
`<state>/aibox/operations/`. The tests track builds the executable from the
commit under test and uses nothing else from the implementation.

### Acceptance harness

The tests track chooses the harness language and layout under these fixed
rules:

- **Entry point:** `tests/acceptance/v1-06/run.sh --aibox PATH --evidence DIR`.
  Exit 0 when every journey passed, 1 when any journey failed or was not run,
  2 for a usage or infrastructure error. `DIR` receives one result record per
  journey (ID, covered requirements, outcome, observed output on failure) and
  a summary.
- **Build:** the binary under test is built static with `CGO_ENABLED=0 go build
  -o PATH ./cmd/aibox` in `localhost/aibox-docs-build:go1.27.0`, because the
  development container has no Go.
- **Run:** the suite runs in the development container, which has Python 3,
  `pty` and rootless Podman.
- **Root journeys are required.** J-CFG-02 (pairs with `system`), J-CFG-03,
  J-CFG-06 (a) and J-CFG-08 run in a disposable container from
  `docker.io/library/debian@sha256:a72b4433de1cf7d6bae0c4cee74d0547f58601a643b3375f070c51600d63b6b2`
  (Debian 13 slim, present locally, no network needed). The test copies the
  static binary in, writes `/etc/aibox/settings.json` as root, creates an
  ordinary user with `useradd`, and runs aibox as that user with
  `podman exec --user`. If Podman is unavailable, that is an infrastructure
  failure: one evidenced report to the owner, and the gate stays blocked. It
  is never a skip.
- **YAML 1.2:** J-OUT-02 uses `ruamel.yaml` 0.19.1, vendored as its wheel
  under `tests/acceptance/v1-06/` with its SHA-256 recorded and checked. The
  gate runs offline.
- **J-CFG-24:** the frozen journey names `spec/v1/examples/customized`. #504
  moved the example to `examples/customized` before the tracks started
  (logbook, 2026-11-15); the test runs the journey's command with that path
  and nothing else changed. The frozen text stays a historical reference
  ([path map](../docs/engineering/spec-layout-migration.md)).
- **Vendored wheel:** the tests track fetches the `ruamel.yaml` wheel once,
  records its source URL and SHA-256, and treats a missing or mismatching
  wheel as an infrastructure error (exit 2), never as a skip.

### Offline gate

`scripts/verify-v1-06.sh` (integration owner) runs, in the docs-build image
where Go is needed:

- `gofmt -l cmd internal` (empty), `go vet ./...`, `go test ./...`,
  `go test -race ./...`, `go mod verify`;
- `sha256sum -c` of `spec/settings-and-logging/v1/FREEZE.sha256`, and an empty
  `git diff 7c2d4b89 -- spec/settings-and-logging/v1`;
- the static build and `tests/acceptance/v1-06/run.sh` (in the development
  container);
- the J-CFG-24 demo command, recorded as evidence;
- the v1 Hugo build and `git diff --check`.

The demo line in `scripts/verify-v1-05.sh` stays as it is.

### Risks

| Risk | Handling |
|---|---|
| `flag` from the standard library accepts a repeated flag and prints its own usage text, which breaks CFG-056 and OUT-007 | W4 replaces it with a small parser driven by the flag table (ADR) |
| Concurrent appends and rotation (LOG-018) across processes | Lock file with `flock`, rotation under the lock; J-LOG-08 with 16 processes |
| YAML without a dependency (D7) | Block YAML whose scalars are JSON scalars; the tests track checks it with the vendored YAML 1.2 parser |
| Root journeys depend on rootless Podman in the development container | Required, never skipped; unavailability blocks the gate and is reported to the owner |

## Logbook

| Date | Entry |
|---|---|
| 2026-11-15 | Phase started after V1-05 closed (#523). Plan written by the integration owner; plan review by `sage` requested. |
| 2026-11-15 | Plan review round 1 by `sage` on `9e7be163`: ready with changes; blocking B1 (root journeys could be skipped), B2 (J-CFG-24 path guidance pointed at files that do not exist yet), B3 (no entry point, build and run environment or YAML parser for the suite; ownership of `verify-v1-06.sh`); non-blocking N1 to N12. All applied in the following commit: root journeys required in a pinned Debian image, frozen J-CFG-24 path with a logbook route for #504, fixed entry point and environment, vendored `ruamel.yaml`, the gate script with the integration owner, and the wave details for N2 to N12. |
| 2026-11-15 | Plan review round 2 by `sage` on `88f37621`: ready; B1 to B3 and N1 to N12 closed. New non-blocking N13 (lock-wait failure path against LOG-018) and N14 (vendored wheel provenance) applied in the following commit. |
| 2026-11-15 | Decision of the integration owner: #504 (PR #524, `70a750f4`) moved `spec/v1/examples/` to `examples/` and `spec/v1/roadmap.yaml` to `roadmap.yaml` before the tracks started. J-CFG-24 runs with `examples/customized`; the frozen contract is unchanged and its path is a historical reference. The phase branch merged `main` at `70a750f4`. |
| 2026-11-15 | Acceptance-test track delivered `4bc6f711` on `feat/settings-and-logging-tests`: 41 journeys, 43 harness unit tests; against the pre-V1-06 binary 3 journeys pass and 38 fail, none skipped. Root journeys run in the pinned Debian image under rootless Podman with the engine's default network, because `--network=none` cannot start in the development container; nothing in the container uses the network. |
| 2026-11-15 | Contract readings by the integration owner for the points the tests track raised. None changes the contract; each is the literal reading. (1) The system settings file in root journeys is `root:root` mode 0644: the common setup's 0600 applies to files the journey's user owns, and a root-owned 0600 file would be `unreadable`. (2) J-CFG-07 `schemaVersion` `2` is the string `"2"`; by `settings.schemaDefects` any value other than `"1"` is `unsupported_version`. (3) `--context` keeps the operation core's behaviour: optional, default `local`; journeys that omit it run as written. (4) J-LOG-08 and J-LOG-10 inspect project P. (5) J-LOG-07: each run is `inspect --project P --format json`; "only the events of the last run" means every line carries that run's `requestId`; the fifth refill starts with `E`; with `retentionFiles` 1 the age step has no `r.log.2` and is recorded as such. (6) J-LOG-09 asserts only its own two files and successful runs. (7) J-CFG-19: a rejection by the input schema is either a JSON-RPC error or a tool result with `isError` true. (8) LOG-003: `sequence` numbers the events that are emitted for the request; a filtered event takes no number. (9) J-OUT-02 ignores the top-level `requestId` and `data.observedAt`. (10) J-LOG-06 compares `/tmp/x.log` before and after. |
| 2026-11-15 | Implementation track delivered `2c510d67` on `feat/settings-and-logging-impl` (waves W1 to W5; gofmt, vet, tests, race, `go mod verify`, static build, Hugo build passed on the track). Its 15 contract readings are listed in its report and reviewed at verification. |
| 2026-11-15 | Integration run 1: both tracks merged on the phase branch (`ae557a89`), static binary built in the docs-build image, `tests/acceptance/v1-06/run.sh`: 40 of 41 journeys passed. J-OUT-02, case "View `--format yaml`", failed: the YAML value differs from the JSON form only in `entry(output.format).value` (`yaml` against `json`) besides `requestId` and `observedAt`. Classification: **test defect**. CFG-061 requires that entry to carry the value of `--format`, so the two calls can never be equal; the journey's ignore list is incomplete. Integration-owner decision against the contract: the View comparison also ignores `entry(output.format).value`. Returned to the tests track. |
| 2026-11-15 | Integration run 2 after the test fix `7f78de6b`: 41 of 41 journeys passed, none skipped. The offline gate `scripts/verify-v1-06.sh` (integration owner, `004c7c98`; the harness unit tests moved to the development container, because the docs-build image's Python 3.13 cannot import the vendored wheel) passed at `004c7c98` in `localhost/aibox-docs-build:go1.27.0`: frozen contract unchanged, gofmt, vet, tests, race, `go mod verify`, static build, darwin/arm64 build, roadmap demo (`logging.level` `debug` from layer `user`), specification and example checks, Hugo build (57 pages), `git diff --check`, no v0 change. On the gate's binary (SHA-256 `579d9ce2…`) the 44 harness unit tests and all 41 journeys passed in the development container. Evidence: [`evidence/v1/V1-06-004c7c98/`](../evidence/v1/V1-06-004c7c98/). Next: independent verification (step 9). |
