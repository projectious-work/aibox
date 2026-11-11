# Application configuration standard

## Purpose and scope

This standard defines predictable, inspectable, and safe configuration for
projectious.work applications. It applies to applications that combine two or
more configuration sources. Libraries with no autonomous runtime use their
host application's configuration contract instead.

Terms in *italics* are [Semantic Anchors](https://github.com/llm-coding/Semantic-Anchors/tree/6554641a11fdf99fdbdf09cc0602268e3db88966); see
[README](README.md#conventions).

Services and workers follow *Twelve-Factor App* factor III: configuration comes
from the environment, and no per-environment configuration files are
committed. Layered configuration files are the company exception for CLI
applications.

## Precedence

When present (CLIs; services normally use only defaults, environment, and
flags), configuration layers resolve from lower to higher precedence:

```text
built-in defaults
→ system configuration
→ user configuration
→ project or repository policy
→ environment variables
→ command-line or invocation overrides
```

A project MAY omit layers but MUST NOT reorder layers silently. Invocation
overrides affect only that invocation unless an explicit write command is
used.

## Authority is separate from precedence

Higher precedence does not grant unlimited authority. Every setting belongs to
an authority class, such as presentation, operation, security, credential,
executable discovery, storage, or policy. A committed project file MUST NOT,
without explicit authorization:

- select or replace executables;
- introduce credentials or secret values;
- weaken system or user security policy;
- redirect logs, caches, or durable state outside permitted roots; or
- enable a new network or output destination.

Rejected settings produce actionable diagnostics rather than being ignored.

## File and value behavior

- Configuration file locations and path anchors are documented.
- Project policy and mutable operational state use separate locations.
- Unknown keys and unsupported versions follow an explicit fail or
  compatibility-warning policy.
- Relative paths resolve against the configuration file or another documented
  anchor, never an accidental current directory.
- Environment variables use a project-specific uppercase prefix.
- Empty, absent, null, and default values have distinct documented semantics.
- Configuration schemas are versioned when the file is a public contract.
- Deprecated keys identify their replacement and removal horizon.

Secrets are injected, referenced, or obtained from an approved secret store.
No configuration file at any layer contains a secret value.

## Effective configuration and provenance

Applications with layered configuration MUST expose a diagnostic view of the
effective configuration. It reports, in deterministic key order:

- every supported key exactly once;
- its redacted effective value;
- the winning source and source layer;
- loaded and absent file layers;
- overridden sources when useful; and
- rejected settings with reasons.

Sensitive paths, credentials, tokens, personal data, and secret-shaped values
undergo the normal redaction policy. Diagnostic output MUST be available in a stable
machine-readable form for applications that expose machine interfaces.

## Testing

Tests cover precedence permutations, authority restrictions, missing and
malformed files, unknown keys, deprecated values, path anchoring, environment
isolation, redaction, and deterministic provenance output. Tests MUST NOT read
the developer's real home configuration or ambient credentials.

## Agent briefing

> Configuration precedence answers which source wins; authority answers
> whether that source may set the key at all. Preserve both, expose redacted
> provenance, and never use committed project configuration as an implicit
> privilege escalation path. Services and workers are environment-first
> (*Twelve-Factor App* III); layered files are for CLIs.
