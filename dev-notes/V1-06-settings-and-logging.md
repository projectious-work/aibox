# V1-06 layered settings and diagnostic sinks

V1-06 implements the frozen specification `settings-and-logging@v1`
([`spec/settings-and-logging/v1/`](../spec/settings-and-logging/v1/), tag
`spec/settings-and-logging/v1`, contract commit
`7c2d4b89566d1869d720cd8b23a5afaae90ce68f`). Implementation pull requests
carry `Spec: settings-and-logging@v1` and never change `spec/`.

Status: planning. The cycle is the one in chapter 21 of the v1 product
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
2. **W2 logging and redaction** (`internal/logging`, `internal/redact`):
   event model and validation against the log-event schema, level filter,
   plain and JSON Lines encoders, bounded queue with drop and direct-write
   rules, stderr sink that tolerates a closed reader, file sink (anchoring,
   permitted root, modes, append, lock file, rotation, retention), and
   redaction of known secrets, secret shapes and the home prefix.
3. **W3 result and output** (`internal/contract`, `internal/output`): result
   fields `settingsSources` and `overridden`, warning `log_sink_degraded`,
   pre-operation error envelope (`invalid_request`), the operation-result
   schema of this version, and the JSON, YAML and human encoders with the
   colour rule and the bootstrap diagnostic.
4. **W4 wiring** (`cmd/aibox`, `internal/usecase`, `internal/mcpserver`):
   a command-line parser that follows `cli.flags`, `cli.syntax` and
   `settings.selectors`; settings resolution before the operation;
   `inspect --effective-config`; MCP input `effectiveConfig`; settings fixed
   at `mcp serve` start; event emission in the orders of
   `logging.eventOrder`; time limits; exit codes.
5. **W5 documentation**: settings hierarchy, environment and flag precedence,
   logging sinks, redaction and troubleshooting reference on the v1 site,
   including capturing the stderr of `mcp serve` under systemd and launchd
   (overview, "Out of scope"); `scripts/verify-v1-06.sh` offline gate.

### File ownership

| Paths | Owner |
|---|---|
| `internal/**`, `cmd/**`, `go.mod`, `go.sum` | Implementation track |
| `tests/acceptance/v1-06/**` | Acceptance-test track |
| `docs-site/content/**` pages for settings and logging, `scripts/verify-v1-06.sh` | Implementation track (W5) |
| `dev-notes/V1-06-settings-and-logging.md`, `roadmap.yaml` (or `spec/v1/roadmap.yaml` until #504 moves it) | Integration owner |
| `spec/**` | Nobody: frozen |

### Frozen shared interfaces

The only interfaces both tracks rely on are the contract's public ones: the
CLI (`cli.flags`, `cli.syntax`, exit codes), the MCP tool `inspect_workspace`
and its input, stdout result documents, stderr events, files under
`<state>/aibox/logs/`, and the schemas under `spec/settings-and-logging/v1/schemas/`.
The tests track builds the executable with `go build ./cmd/aibox` and
uses nothing else from the implementation.

### Acceptance harness

The tests track owns its harness choice. Constraints from the contract and the
development container:

- Journeys run as a non-root user with a fresh `HOME` and a minimal
  environment (journeys, "Common setup").
- Journeys marked **root** (J-CFG-02 pairs with `system`, J-CFG-03, J-CFG-06 (a),
  J-CFG-08) need a disposable container in which the test writes
  `/etc/aibox/settings.json` as root and runs aibox as another user. The
  development container has rootless Podman; a Debian 13 image with the built
  binary copied in is enough.
- J-OUT-02 needs a YAML 1.2 parser; J-OUT-03 needs a pseudo-terminal.
- J-CFG-24 names `spec/v1/examples/customized`. Issue #504 moves the examples
  to `examples/customized`; the frozen path is a historical reference, and the
  test uses the example's current location
  (`docs/engineering/spec-layout-migration.md`).

### Risks

| Risk | Handling |
|---|---|
| `flag` from the standard library accepts a repeated flag and prints its own usage text, which breaks CFG-056 and OUT-007 | W4 replaces it with a small parser driven by the flag table |
| Concurrent appends and rotation (LOG-018) across processes | Lock file `<file>.lock` with `flock`, rotation under the lock; J-LOG-08 with 16 processes |
| YAML without a dependency (D7) | Block YAML whose scalars are JSON scalars; the tests track checks it with an independent YAML 1.2 parser |
| Root journeys depend on Podman in the development container | Recorded as infrastructure, not as a finding, if Podman is unavailable |

## Logbook

| Date | Entry |
|---|---|
| 2026-11-15 | Phase started after V1-05 closed (#523). Plan written by the integration owner; plan review by `sage` requested. |
