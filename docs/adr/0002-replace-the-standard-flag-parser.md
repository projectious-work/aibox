# 0002. Replace the standard flag parser

- **Status:** Accepted
- **Date:** 2026-11-15
- **Deciders:** implementation track of V1-06, under the plan reviewed by the
  integration owner (`dev-notes/V1-06-settings-and-logging.md`, risk table)

## Context

Up to V1-05, `cmd/aibox` parsed `inspect` and `mcp serve` with the standard
library's `flag` package. The settings-and-logging contract
(`spec/settings-and-logging/v1/`) fixes the command line exactly:

- `cli.syntax`: `--name value` and `--name=value`; an unknown flag, a
  repeated non-repeatable flag, a flag on a command that does not list it,
  and any positional argument are `usage_error` (CFG-056).
- `--env-file` is repeatable, at most 16 times (CFG-044).
- A rejection writes one `settings.rejected` event in the format that the
  command line asks for (`logging.rejectionFormat`), then either an
  `invalid_request` envelope (OUT-006) or one diagnostic line
  `aibox: invalid_input: <reason>: <text>` as the last stderr line (OUT-007).
- Flag values are not checked by the parser: a defective user settings file
  is reported before `--log-level loud` (CFG-008, J-CFG-15 (h)).

`flag` accepts a repeated flag silently (the last value wins), accepts
single-dash forms, stops at the first positional argument instead of
rejecting it, and prints its own usage text to its output, which would put a
second, uncontrolled message on stderr.

## Decision

Replace `flag` in `cmd/aibox` with a small parser (`cmd/aibox/cmdline.go`)
driven by one table that mirrors `cli.flags` and `settings.selectors`: name,
whether it takes a value, whether it may repeat, whether it sets a settings
key, and the commands that accept it. The parser:

- returns only `usage_error`, with a message for the diagnostic;
- hands settings flags to the resolver as raw `(flag, value)` pairs in
  command-line order, so value checks happen at the invocation layer;
- reads `--log-format`, `--log-level` and `--format` even from a rejected
  command line, for the rejection event and envelope;
- writes nothing itself.

A drift test compares the table with `interfaces.json`.

## Consequences

- Command-line behaviour is defined by the contract table, and a contract
  change that adds a flag fails the drift test until the table follows.
- aibox owns all of its stderr output; no library writes usage text.
- There is no generated `--help` per command; `aibox help` prints one usage
  text maintained by hand.
- Short options and flag abbreviations are not supported. Later commands
  (V1-09) extend the same table rather than adding a CLI framework; if the
  command set grows large, a framework can be reconsidered in a new ADR.
