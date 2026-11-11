# Application profiles for engineering standards

## Purpose

Company standards describe outcomes that span several kinds of software. This
document defines applicability profiles so a CLI convention is not imposed on
a library, and a service requirement is not imposed on a documentation-only
repository.

Terms in *italics* are [Semantic Anchors](https://github.com/llm-coding/Semantic-Anchors/tree/6554641a11fdf99fdbdf09cc0602268e3db88966); see
[README](README.md#conventions).

A repository declares one primary profile and any secondary profiles in its
engineering documentation. It records justified exceptions rather than
silently ignoring an applicable requirement.

## Profiles

### CLI application

Applies configuration precedence, stdout/stderr separation, terminal and
machine rendering, exit codes, signal handling, cross-platform distribution,
and black-box command tests. Command grammar follows the
[command-line interface design standard](command-line-interface-design.md)
(*Unix Philosophy* and one recorded style).

### Service or worker

Applies API/protocol compatibility, structured logs, correlation and tracing,
health and readiness, cancellation, graceful shutdown, deployment
configuration, data retention, and service-level integration tests. Follows
*Twelve-Factor App* (see [Twelve-Factor alignment](#twelve-factor-alignment)).

### User-interface application

Applies accessible interaction, local and remote state boundaries, privacy,
error recovery, visual regression where valuable, API compatibility, and
user-journey tests, and the [user-interface design standard](user-interface-design.md).
CLI-only stream and exit-code rules do not apply unless the project also ships
a CLI.

### MCP server

Applies [agent-native product interfaces](agent-native-product-interfaces.md).
A local stdio server follows the CLI stream rules; a remote server also follows
the Service or worker profile.

### AI agent

Applies to an agent the company builds or configures: a product agent, a
sub-agent definition, a skill, or repository agent instructions. Follows the
[agent building standard](agent-building.md).

### Library or SDK

Applies API compatibility, dependency discipline, supported runtime versions,
examples as contracts, unit/property tests, and package publication. It
normally delegates logging and configuration policy to the host application.

### Schema, protocol, or process package

Applies independent contract versioning, metaschema validation, positive and
negative fixtures, compatibility and migration, generated-reference drift,
and consumer conformance tests.

### Infrastructure or template product

Applies declarative-contract validation, provider/tool compatibility,
credential isolation, plan/review/rollback or teardown evidence, disposable
end-to-end acceptance, and cost/cleanup controls.

### Documentation or website

Applies local reproducible builds, links, accessibility, generated-content
drift, public metadata, content/release consistency, and the
[user-interface design standard](user-interface-design.md). Runtime configuration
and application logging rules are normally not applicable.

## Twelve-Factor alignment

Services, workers, and remote MCP servers follow *Twelve-Factor App*. Other
profiles follow it where a factor is meaningful. The company deviations are
exactly these; a project records any further one in its applicability record.

| Factor | Company position |
|---|---|
| I Codebase | Followed: one repository per product; supported versions are `lts/*` branches ([branching](git-branching-and-releases.md)). |
| II Dependencies | Followed: declared and pinned ([security](security-and-software-supply-chain.md)). |
| III Config | **Deviation for CLIs:** layered files (system, user, project) plus environment and flags under the [configuration standard](application-configuration.md). Services and workers are environment-first. |
| IV Backing services | Followed by services, workers, and MCP servers that use them. |
| V Build, release, run | **Deviation:** no CI; the release is local and human-controlled. Stages stay separate: an immutable validated candidate, then bounded publication ([host-phase execution](human-controlled-host-phase-execution.md)). |
| VI Processes | Followed by services, workers, and MCP servers. CLIs keep state only in documented state locations. |
| VII Port binding | Followed by services and remote MCP servers. |
| VIII Concurrency | Followed by services and workers. |
| IX Disposability | Followed by services, workers, and MCP servers; for CLIs it means signal handling and cancellation. |
| X Dev/prod parity | **Deviation:** development containers are not the runtime. Host-only properties are proven on real hosts (host-gated release overlay). |
| XI Logs | Followed by services and workers. **Deviation for CLIs:** logs go to stderr, and file sinks are allowed. **Addition:** evidence is a separate channel ([output, logging, and evidence](application-output-logging-and-evidence.md)). |
| XII Admin processes | Followed by services: one-off commands use the same build and configuration. A CLI is its own admin surface. |

Libraries and SDKs delegate factors III to XII to their host application.
Documentation and websites apply no runtime factor.

## Conditional overlays

### Host-gated release

Applies when native-host, container-engine, signing, or publication authority
is intentionally unavailable in the development environment. It activates the
[human-controlled host-phase execution standard](human-controlled-host-phase-execution.md),
including the portable producer/verifier handover, lifecycle separation,
schema-versioned validation manifest, credential-free validation, bounded
publisher, capability adapters, and conformance cases.

The overlay is selected because of authority and trust boundaries, not because
of implementation language or a particular container product. Projects record
the host-only properties being proven and any reviewed non-applicable cases.

## Composite products

A product MAY declare several profiles. Requirements combine; one profile does
not weaken another. The repository documents component boundaries so reviewers
can identify which profile governs each artifact.

## Applicability record

The project records:

- selected profiles and components;
- supported platforms and distribution surfaces;
- applicable verification and release gates;
- stronger project-specific requirements; and
- non-applicable or deferred company requirements with rationale.

Profile selection is reviewed when the product boundary or distribution model
changes.

## Agent briefing

> Select standards by product behavior, not implementation language. Declare
> every applicable profile, combine their requirements, and make exceptions
> explicit and reviewable. Services, workers, and remote MCP servers follow
> *Twelve-Factor App* with the listed deviations: layered CLI configuration,
> local release, development containers, and evidence as a separate channel.
