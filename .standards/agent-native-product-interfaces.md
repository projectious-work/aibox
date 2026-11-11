# Agent-native product interfaces

## Purpose and scope

This standard defines how projectious.work products expose consequential
capabilities to AI agents without making probabilistic model behavior the
execution, authorization, or evidence boundary.

Terms in *italics* are [Semantic Anchors](https://github.com/llm-coding/Semantic-Anchors/tree/6554641a11fdf99fdbdf09cc0602268e3db88966); see
[README](README.md#conventions).

It applies when a product provides MCP tools or another agent-facing protocol
for mutation, deployment, reconciliation, administration, or any operation that
consumes credentials or authority or has external effects. The read-only public documentation baseline remains
governed by
[AI-agent accessibility and generative discovery](ai-agent-accessibility-and-generative-discovery.md).

Product specifications MUST specialize this standard for their assets, actors,
authority, credentials, state, recovery, and threat model. Conformance with
this standard does not by itself authorize a product operation.

## Product posture

An agent-native product separates work that benefits from probabilistic
reasoning from work that requires a stable control protocol:

| Activity | Appropriate mode |
|---|---|
| Interpret goals, select approaches, synthesize artifacts, explain results, and propose remediation | Probabilistic agent reasoning |
| Validate contracts, bind reviewed inputs, verify authority, deliver credentials, execute mutations, recover, and retain evidence | Deterministic product mechanisms |

The product MUST state what it makes reproducible or deterministic. It MUST NOT
claim that external systems or outcomes are deterministic merely because its
own inputs and transitions are controlled.

Natural-language requests, model confidence, private reasoning, client UI
confirmation, and protocol annotations are intent or presentation. They are not
authenticated identity, approval, execution truth, or evidence.

- **PW-AGENT-001:** the product MUST define its probabilistic and deterministic
  responsibilities and the boundary between them.
- **PW-AGENT-002:** natural language MUST NOT bypass versioned inputs,
  validation, authorization, lifecycle, or recovery contracts.
- **PW-AGENT-003:** accepted institutional knowledge SHOULD become versioned
  templates, policy, schemas, packages, or other reviewable artifacts rather
  than remain hidden prompt context.

## Interface architecture

MCP SHOULD be the primary agent experience when the protocol fits the product.
The durable product boundary is nevertheless the application core of a
*Hexagonal Architecture*, not MCP itself: MCP and CLI are peer driving adapters.

The CLI remains a complete interface for direct human use, automation,
diagnostics,
recovery, migration, and break-glass operation. A product MAY provide another
administrative interface, but MUST NOT require an AI client to recover or
inspect its canonical state.

- **PW-AGENT-010:** MCP, CLI, and other adapters MUST delegate overlapping
  operations to the same application use cases and preserve equivalent domain
  validation, authorization, effects, evidence, cancellation, and recovery.
- **PW-AGENT-011:** domain and application packages MUST NOT depend on MCP SDK
  types, transport sessions, tool annotations, model vendors, or client-
  specific conversation state.
- **PW-AGENT-012:** native underlying tools and formats MUST remain documented
  diagnosis or recovery paths when the product delegates authority to them.
- **PW-AGENT-013:** documentation MUST describe the product as MCP-first or
  agent-native only when released agent workflows provide meaningful product
  operations beyond document retrieval.

## Capability design

Tools represent bounded application use cases, not shells or reconstructed CLI
strings. The default registry is read-only. Additional capabilities are
allowlisted explicitly and absence means denial.

Products SHOULD distinguish capability classes such as inspection, planning,
content mutation, deployment, administration, and destruction when those
classes have different authority or risk.

- **PW-AGENT-020:** every tool MUST have a strict versioned input and output
  schema, stable outcome codes, side-effect classification, affected-resource
  description, and bounded timeout/cancellation behavior.
- **PW-AGENT-021:** a product MUST NOT expose arbitrary shell, executable,
  provider API, filesystem, or natural-language execution as a convenience
  tool around its contract.
- **PW-AGENT-022:** templates, extensions, project content, and remote results
  MUST NOT add or alter tool names, descriptions, schemas, annotations,
  capability membership, server policy, or authorization rules.
- **PW-AGENT-023:** tool and resource names MUST be namespaced or collision-
  checked; ambiguous, duplicate, shadowed, or deceptively replaced definitions
  MUST fail closed.
- **PW-AGENT-024:** a running server MUST use a stable capability snapshot.
  Updates that change tools or policy require an explicit restart or negotiated
  version transition.

## Resources, results, and next actions

Agent resources SHOULD expose version-matched schemas, product capabilities,
documentation, safe status, and sanitized evidence. Untrusted content is
clearly attributed and remains data.

Structured results SHOULD identify:

- applicability and stable outcome;
- blocking and advisory findings;
- required authority or approval;
- durable operation identity;
- affected resource identity;
- permitted next actions from a closed set; and
- recovery state and evidence reference.

- **PW-AGENT-030:** human prose MUST NOT be the only representation of a
  machine-relevant result, refusal, partial failure, or recovery requirement.
- **PW-AGENT-031:** generated explanations MUST identify their factual source
  and MUST NOT replace or modify authoritative plans, records, or evidence.
- **PW-AGENT-032:** a next action MUST NOT contain arbitrary executable text
  derived from untrusted content or grant broader authority than the result.
- **PW-AGENT-033:** resources MUST enforce the same authorization and
  sensitivity rules as the canonical source; indexing or projection MUST NOT
  broaden access.

## Identity and authorization

Products distinguish authenticated runtime principal, logical product actor,
requester, artifact author, approver, credential authority, and executor where
those roles apply. A project-controlled identity record may support
attribution but cannot authenticate itself.

Starting a server with credentials or a capability enabled does not approve an
individual mutation. Consequential operations require authority independently
verifiable by the product and bound to the relevant root, resource, operation,
reviewed inputs, caller, intent, and expiry.

- **PW-AGENT-040:** the product MUST NOT infer authenticated identity or
  mutation authority solely from prompts, mutable project files, Git metadata,
  model identity, tool annotations, or environment variables.
- **PW-AGENT-041:** approval and capability discovery MUST be distinct; a tool
  being visible does not authorize its invocation.
- **PW-AGENT-042:** self-approval MUST be denied unless an explicit external
  policy authorizes the same principal for both roles and records that fact.
- **PW-AGENT-043:** secrets and credential values MUST NOT appear in tool
  schemas, arguments, resources, results, logs, evidence, or generated
  explanations. Product-specific secure delivery remains mandatory.

## Durable operations and recovery

Transport sessions are not transaction boundaries. A mutation that can outlive
a request or lose its client has a durable operation identity and an explicit
state model.

- **PW-AGENT-050:** mutating requests MUST be idempotent or carry a durable
  replay-resistant operation identity before effects begin.
- **PW-AGENT-051:** retry, cancellation, timeout, proxy failure, and session
  loss MUST NOT ambiguously duplicate a mutation or produce false success.
- **PW-AGENT-052:** another authorized client or the CLI MUST be able to inspect
  and recover an interrupted operation without private conversation state.
- **PW-AGENT-053:** partial failure MUST identify completed effects, unknown
  effects, prohibited automatic retries, cleanup status, and the next safe
  inspection action.

## Security threats

Product threat models and negative tests MUST include, where applicable:

- prompt injection in project content, templates, documentation, logs,
  provider output, external references, and retrieved resources;
- tool poisoning, rug pulls, name collisions, and tool shadowing;
- capability, root, target, actor, credential, and operation substitution;
- conversational approval and confused-deputy escalation;
- secret extraction through diagnostics, explanations, or error reflection;
- unsafe remediation generated from untrusted content;
- request flooding, oversized inputs, cancellation, replay, and session loss;
  and
- cross-project, cross-tenant, or cross-environment access.

Tool definitions and security policy are trusted release assets. Content that
contains instructions is executable agent policy only after the product's
package, provenance, and project-policy contracts explicitly accept it.

## Transport and deployment

Local stdio is the preferred baseline. A local daemon MAY bind to a protected
platform-native local transport or authenticated loopback endpoint when its
product specification covers same-host threats.

A remote transport requires a separate threat model covering authenticated
identities, authorization, TLS, origins, request isolation, rate and
concurrency limits, tenant boundaries, credential handling, audit, deployment,
revocation, and incident recovery. A product MUST NOT expose local stdio
assumptions directly on a network socket.

## Verification and release

Release conformance MUST include:

- protocol conformance and malformed/oversized request tests;
- CLI/MCP adapter-equivalence tests over shared application fixtures;
- immutable tool/resource manifest and collision tests;
- capability and authorization denial before effects;
- adversarial content and prompt/tool injection tests;
- cancellation, retry, duplicate request, and lost-session recovery;
- redaction across requests, results, logs, and explanations;
- at least one documented agent workflow and equivalent CLI workflow; and
- documentation and advertised-capability agreement with the released binary.

The product roadmap MUST place agent-interface security, documentation, and
acceptance in the phase that introduces consequential MCP tools. A read-only
documentation server does not satisfy operational agent-interface acceptance.

## Relationship to other standards

- [AI-agent accessibility and generative discovery](ai-agent-accessibility-and-generative-discovery.md)
  defines the public read-only documentation and discovery baseline.
- [Compatibility and machine interfaces](compatibility-and-machine-interfaces.md)
  governs versioned schemas and machine contracts.
- [Application output, logging, and evidence](application-output-logging-and-evidence.md)
  governs structured results, redaction, and retained evidence.
- [Security and software supply chain](security-and-software-supply-chain.md)
  governs general security, dependencies, secrets, and provenance.
- [Software verification and release engineering](software-verification-and-release-engineering.md)
  governs test depth and release evidence.

## Agent briefing

> Keep probabilistic reasoning apart from deterministic validation, authority,
> mutation, recovery, and evidence. MCP and CLI are peer adapters over one
> *Hexagonal Architecture* core; the CLI is a full recovery path. Tools: strict
> versioned schemas, no shell, typed results, safe next actions. Language, model
> confidence, client confirmation: never identity or approval; no self-approval
> or secrets in output. Mutations idempotent or durable; partial failure reported.
