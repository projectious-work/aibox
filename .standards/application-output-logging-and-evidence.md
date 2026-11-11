# Application output, logging, and evidence standard

## Purpose and scope

This standard separates requested results, operational diagnostics, and
durable domain evidence. It applies to command-line applications, services,
workers, and user-facing applications according to their application profile.

Terms in *italics* are [Semantic Anchors](https://github.com/llm-coding/Semantic-Anchors/tree/6554641a11fdf99fdbdf09cc0602268e3db88966); see
[README](README.md#conventions).

## Three information channels

1. **Application output** communicates the result requested by a person or
   program.
2. **Operational logs** explain runtime decisions, dependencies, and failures.
3. **Domain or run evidence** records durable facts required for recovery,
   audit, reconciliation, or product behavior.

These channels MAY share an internal event model but MUST have distinct
contracts and retention.

Operational logs follow *Twelve-Factor App* factor XI: an application writes
log events as an unbuffered stream (stdout for services, workers, and remote
MCP servers; stderr for CLIs and stdio MCP servers, whose stdout carries only
the protocol), and the environment collects, routes, and retains them. Logs may be
rotated, sampled, or lost, so nothing may depend on them as proof. Whatever must
survive that, to support recovery or audit or to answer "what exactly
happened", is evidence: the application writes it deliberately to its own store
with its own retention. A log line becomes evidence only when a product
contract explicitly copies it into that store.

## Semantic events and rendering

Applications construct typed or structured semantic events before rendering.
Terminal escape sequences, localized prose, and sink formatting do not belong
in domain event objects.

Structured events SHOULD include where applicable:

- schema or event version;
- timestamp and severity;
- stable event name;
- component and operation;
- correlation, request, run, or trace identifier;
- outcome and error category;
- attributed child or external dependency; and
- redacted structured fields.

The actor identifies the human, service, or agent that performed or attested
the action. It MUST NOT be replaced by the affected subject, task, gate, or
artifact identifier. Delegated actions record both the initiating and
executing identities when that distinction matters to authority or audit.

Human renderers MAY use color, tables, progress, and other decoration. Plain
and machine renderers remain stable and deterministic.

Long-running trusted gates follow the
[host-phase operator contract](human-controlled-host-phase-execution.md#operator-contract)
for events, renderers, and progress. TUI state, widget state, and clipboard
contents are never evidence.

## Streams and formats

CLI applications write requested results to stdout and diagnostics to stderr.
Machine-readable output MUST contain only the declared envelope on stdout.
Non-interactive execution selects plain output unless explicitly overridden.
Color control respects `NO_COLOR` and project-specific settings.

Services and workers leave log routing and rotation to the platform; CLIs MAY
offer file sinks. Applications document supported log sinks, formats, levels,
rotation, retention, and failure behavior. Enabling a sink that can disclose data or
perform network egress requires an authorized configuration source.

## Redaction and data handling

Redaction occurs before events reach any renderer, buffer, file, telemetry
exporter, or test artifact. It covers credentials, secret-shaped values,
personal data, sensitive paths, and product-specific protected fields.
Chunked and concurrent output MUST not bypass redaction at boundaries.

Raw modes are exceptional. They require explicit invocation, document their
risk, and MUST NOT imply that known secrets are safe to print.

## Failure and correlation

Logs and evidence distinguish refusal, dependency failure, timeout,
cancellation, partial success, and complete success. Child-process or remote
events retain source attribution. Concurrent events carry correlation data and
must not be reordered into a false lifecycle narrative.

Immutable validation evidence and mutable publication journals have separate
contracts. Manifest-covered files stop changing before their digests are
finalized; later publication commands and results cannot mutate that manifest.

## Evidence integrity

A gate or milestone pass records enough structured evidence to reproduce or
independently assess the claim. Where applicable this includes the exact
subject or candidate, validator and version, invocation or evaluation profile,
actor, occurrence time, outcome, warnings or exceptions, and a stable evidence
locator or digest. A bare success sentence, HTTP status, or release URL is
useful context but is not sufficient proof when stronger machine-verifiable
evidence exists.

Evidence references remain resolvable outside the producer's workspace.
Producer-local absolute paths MAY appear as diagnostic metadata but MUST NOT be
the sole locator for an audit record, release artifact, or gate result.

## Testing

Tests separately assert output, diagnostics, structured logs, evidence,
redaction, exit or outcome status, rotation, concurrency, and sink failure.
Golden output is reviewed as a product contract; volatile values are
normalized deliberately rather than ignored broadly.

## Agent briefing

> Output answers the request. Logs are a disposable event stream that explains
> operation (*Twelve-Factor App* XI). Evidence is a durable record written on
> purpose that must survive log loss; a log line is never proof. Do not collapse
> these channels. Structure and redact events before rendering or export.
