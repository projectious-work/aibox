# Security and software-supply-chain standard

## Purpose and scope

This standard defines the minimum security and dependency practices for
projectious.work software. Projects add stronger controls according to their
trust boundary, data sensitivity, deployment model, and release profile.

Terms in *italics* are [Semantic Anchors](https://github.com/llm-coding/Semantic-Anchors/tree/6554641a11fdf99fdbdf09cc0602268e3db88966); see
[README](README.md#conventions).

## Security model

Every released application documents:

- protected assets and sensitive data;
- trust boundaries and privileged actors;
- credentials and authorization model;
- external processes, services, plugins, and network destinations;
- persistence, retention, and deletion responsibilities;
- likely abuse and failure paths, modeled with *STRIDE*, and for every web or
  HTTP interface checked against the current *OWASP Top 10*; and
- supported security-reporting process.

Security-sensitive behavior fails closed. Error messages and diagnostics do
not expose credentials, personal data, private content, or unsafe recovery
instructions.

## Input, filesystem, and process safety

Applications validate untrusted structured input, versions, sizes, paths,
archives, URLs, and encodings wherever untrusted data enters, and again
before use in a filesystem or process operation (*Defensive Programming
according to McConnell*; see the
[architecture standard](software-architecture-and-code-design.md#code-conventions)).
File operations account for traversal,
symlinks, special files, permissions, partial writes, and races.

Programs invoking external tools preserve argument boundaries and avoid a shell
unless shell semantics are explicitly required and reviewed. They control
working directories, forwarded environment variables, stdin, timeouts,
cancellation, and captured output. Supported dependency versions are tested.

Host gates resolve implementations by reviewed capability probes rather than
product names. Fixed package-manager symlinks are acceptable only when both
the invocation path and canonical target are validated and evidenced. Missing
prerequisites produce advisory instructions; trusted gates do not elevate or
modify system configuration.

## Secrets and sensitive data

- Secrets are never committed or placed in ordinary configuration.
- Redaction occurs before buffering, logging, telemetry, evidence export, test
  artifacts, or error rendering.
- Test fixtures use reserved domains, synthetic identities, and placeholders.
- Retention is minimized and documented.
- Debug or raw modes require explicit action and remain subject to known-secret
  protection.
- Repositories and release artifacts undergo secret scanning.

## Dependency policy

Every direct dependency has a documented purpose. Selection and upgrade review
considers maintenance, license, provenance, transitive weight, vulnerabilities,
platform support, and whether a standard-library or small local solution is
clearer.

Lockfiles are committed where supported. Dependency updates include relevant
tests, advisory review, and changelog or migration review. Vendored source,
fonts, media, binaries, and generated bundles retain required notices and
provenance.

## Release supply chain

Applicable release profiles require:

- vulnerability analysis of reachable code and dependency manifests;
- static security analysis appropriate to the language;
- license and provenance review;
- SBOM generation for final artifacts;
- scanning of built artifacts, not only source;
- checksums and signatures or attestations;
- build metadata identifying the source commit and toolchain; and
- independent post-publication verification.

Signing identities, publishing credentials, trust roots, and trusted caches
are scoped to the minimum release stage, subject, destination, and lifetime
that require them. Production credentials are not exposed to candidate
preparation or credential-free validation. Authorization evidence identifies
the exact immutable subject and intended distribution target; it cannot be
replayed for a changed candidate or broader destination.

Vulnerability evidence preserves complete raw scanner output and a normalized
disposition. Fix-available findings at the configured blocking severity fail
the gate. No-fix findings at that severity become explicit
publication-required warnings rather than being silently ignored or making
release permanently impossible. Scanner failure and reviewed exceptions are
distinct outcomes.

Content-addressed caches MAY be reused only when their identity binds the
candidate/source and relevant toolchain inputs and no validation command is
skipped. The selected cache policy is evidenced and a force-fresh mode is
available. Downloaded package/build caches remain outside repository scanner
roots or use an explicit reviewed exclusion.

External runtime prerequisites are inventoried separately from packaged
dependencies. A source-provided container definition does not imply that the
project publishes or supports a container image.

## Security regressions and response

Security fixes add permanent negative tests or fixtures. Public communication
uses coordinated disclosure and safe release notes. Published vulnerable tags
are never moved; corrected versions are released normally.

## Agent briefing

> Identify the trust boundary before coding; model abuse with *STRIDE* and
> check web interfaces against the *OWASP Top 10*. Validate untrusted input
> where it enters and before filesystem or process use. Preserve argument and
> path boundaries, redact before data leaves the core, review every direct
> dependency, and verify the final published artifact, not only the build.
