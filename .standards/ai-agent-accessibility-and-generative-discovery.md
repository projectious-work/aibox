# AI-agent accessibility and generative discovery

## Purpose and scope

This standard defines how public projectious.work products make their released
knowledge and capabilities accessible to AI agents. It complements the
[open-source documentation strategy](open-source-documentation-strategy.md)
with three supported access paths:

Terms in *italics* are [Semantic Anchors](https://github.com/llm-coding/Semantic-Anchors/tree/6554641a11fdf99fdbdf09cc0602268e3db88966); see
[README](README.md#conventions).

1. canonical AI-readable product documentation;
2. public generative-discovery surfaces; and
3. a readily usable Model Context Protocol (MCP) server.

The standard applies to public software, libraries, themes, specifications,
design systems, and other released technical products. It governs product
consumption, not repository operation. Agent contribution instructions,
repository mutation, internal context, development governance, and provider
harness integration are outside this standard.

## Principles

1. **One product truth, several projections.** Human documentation, Markdown,
   generative-discovery files, and MCP resources describe the same released
   contract (*SSOT*). Generated projections are preferred to manually
   duplicated facts.
2. **Documentation remains independently useful.** An agent must not need an
   MCP client to install, understand, or use the product.
3. **MCP adds structured access** through read-only resources; tools follow
   [Optional MCP product tools](#optional-mcp-product-tools).
4. **Released facts only.** Maturity, versions, compatibility, limitations,
   security boundaries, and examples match the released implementation.
5. **Open formats and stable identifiers.** Prefer Markdown, semantic HTML,
   JSON Schema, and stable canonical URLs over JavaScript-only presentation or
   proprietary agent formats.
6. **Version alignment is a release property.** AI-facing surfaces are checked
   and published with the product version they describe.

## Required access paths

Every in-scope project MUST provide all three access paths below. They are
complementary: users can read the product documentation directly, discover it
through public generative surfaces, or connect an MCP client.

### AI-readable product documentation

The repository MUST contain a canonical, clutter-free Markdown product guide.
It MAY be split into linked pages when the product cannot be explained clearly
in one bounded document. It MUST cover, where applicable:

- product identity, purpose, intended users, maturity, and non-goals;
- current version and supported compatibility range;
- prerequisites, installation, and a minimal successful example;
- concepts needed to use the product correctly;
- configuration, inputs, outputs, defaults, and public machine contracts;
- supported components, APIs, commands, extension or override points;
- task-focused examples whose commands and expected results are current;
- limitations, trust boundaries, privacy behavior, and external dependencies;
- troubleshooting, diagnostics, and support routes;
- upgrades, deprecations, migrations, and release history; and
- licences, third-party attribution, and asset or dependency provenance.

The guide MUST distinguish normative reference from tutorials and examples.
Facts that also appear in schemas, package metadata, or generated reference
SHOULD be derived from, or checked against, those authoritative sources.

The published documentation site MUST expose the same material as semantic
HTML and SHOULD make the canonical Markdown source or an equivalent `.md`
representation directly retrievable. Essential information MUST remain usable
without client-side JavaScript.

### Public generative discovery

The published documentation root MUST provide `/llms.txt` as a concise index
to authoritative product material. It MUST:

- identify the product and its purpose;
- identify the current documented version;
- link to the canonical getting-started, reference, compatibility, security,
  release, licence, and provenance material that applies;
- prefer stable canonical URLs and concise descriptions; and
- exclude secrets, internal context, unpublished plans, duplicated navigation,
  and generated noise.

Projects SHOULD provide `/llms-full.txt` when a bounded, coherent compilation
improves offline or single-request consumption. It MUST NOT become an
unbounded dump of repository files, histories, generated assets, or internal
process records.

The documentation site MUST also provide:

- a sitemap containing canonical public pages;
- an explicit `robots.txt` crawl policy;
- canonical link metadata and stable, descriptive URLs;
- accurate titles, descriptions, headings, landmarks, and link text;
- structured data only when it matches visible content; and
- visible version, recency, ownership, and provenance signals.

`llms.txt` is a discovery aid, not access control, authorization, or a promise
of inclusion in a generative system. `robots.txt` controls compliant crawler
access; it does not provide product instructions or repository authority.

### MCP access

Every in-scope project MUST provide a documented, readily usable, read-only
MCP resource server. A shared projectious.work server implementation configured
by a small per-project manifest SHOULD be used for the baseline instead of
creating a bespoke protocol implementation in every repository.

The baseline server MUST expose resources for the applicable product material:

- product identity, version, status, and documentation index;
- getting started and minimal usage;
- configuration and machine-readable schemas;
- compatibility, limitations, and external dependencies;
- examples and public reference material;
- changelog, release, and migration guidance; and
- licence, attribution, and provenance.

The repository MUST document one copyable command that starts the server
locally over stdio. Runtime and package versions MUST be pinned or constrained,
and startup MUST NOT require secrets, repository write access, or an external
network connection after installation unless that dependency is an explicit
part of the product contract. A remotely hosted transport MAY additionally be
offered.

The resource server MUST NOT expose private company context, credentials,
unpublished plans, repository-operations records, maintainer-only operations,
or arbitrary
repository file access.

## Optional MCP product tools

Tools follow the *Principle of Least Privilege*: product access does not grant
permission to change a repository or external system.

A project SHOULD expose MCP tools when they add safe, meaningful product
functionality, such as validating configuration, generating a starter file,
querying supported capabilities, or performing a bounded product operation.
Merely wrapping document retrieval in a tool is insufficient; documents belong
in MCP resources.

Tool inputs and outputs MUST have explicit schemas and deterministic error
behavior. A tool that creates content SHOULD return the proposed content to the
client rather than writing to the user's filesystem unless direct mutation is
the product's explicit purpose.

Any tool with side effects MUST additionally define:

- the affected system, data, and reversibility;
- least-privilege authentication and authorization;
- product-verified approval and human-control expectations (client
  confirmation is not approval; see Agent-native product interfaces);
- secret, personal-data, logging, and retention behavior;
- timeouts, retries, idempotency, and partial-failure behavior; and
- security tests appropriate to its capability.

Side-effect-free tools (validation, returning generated content) need only the
schema and deterministic-error rules above. Consequential MCP tools MUST also
conform to
[Agent-native product interfaces](agent-native-product-interfaces.md), which
defines the shared application-core, CLI-parity, capability, authorization,
durable-operation, prompt/tool-safety, and recovery contract. The mandatory
read-only resource-server baseline in this document does not by itself make a
product operationally agent-native.

Remote protected servers MUST follow the current MCP authorization and
transport requirements and the remote-transport rules of Agent-native product
interfaces. A project MUST NOT add a write-capable tool merely to
satisfy this standard; a mutating tool requires an explicit product need,
authorization design, and human control.

## Project manifest and shared implementation

The company SHOULD maintain a versioned manifest schema and reusable Level 1
resource-server implementation. The manifest SHOULD contain references rather
than copied prose and SHOULD declare:

- product identifier, title, description, version, and canonical site;
- canonical Markdown documents and public schemas;
- resource names, media types, audiences, and version applicability;
- server package/version and the local startup command; and
- optional tools and their capability class.

The shared server SHOULD reject paths outside the declared public resource set,
resolve links predictably, report its own product and protocol versions, and
start without scanning the entire repository.

Until the shared implementation and manifest schema exist, a project MAY ship
an equivalent project-local read-only server. This temporary implementation
must still meet the verification and boundary requirements in this standard.

## Release and verification requirements

The local verification suite MUST check, as applicable:

- canonical documentation and public schemas are internally consistent;
- documentation examples and declared startup commands work;
- all required `/llms.txt` links resolve to canonical public material;
- `/llms-full.txt`, when present, is bounded and contains no excluded content;
- sitemap, canonical metadata, and crawler directives are valid;
- the MCP server starts from the documented command;
- MCP resource discovery returns the declared resources;
- every declared resource can be read and identifies the applicable version;
- undeclared files and private/internal paths cannot be read through MCP;
- tool schemas, errors, timeouts, and side-effect boundaries are tested; and
- generated projections are clean or reproducibly regenerated.

A release MUST NOT publish AI-facing material that describes unreleased
features as available. The GitHub release notes and documentation changelog
SHOULD state material changes to public schemas, MCP resources, tools, or
discovery paths.

The documentation site and local stdio MCP server are release surfaces. A
project that cannot provide one of them is not conformant and MUST document the
gap visibly until it is resolved.

## Relationship to other standards

- The [open-source documentation strategy](open-source-documentation-strategy.md)
  governs the complete human-facing documentation and community contract.
- [Compatibility and machine interfaces](compatibility-and-machine-interfaces.md)
  governs schemas, versioning, and machine-consumed contract changes.
- [Security and software supply chain](security-and-software-supply-chain.md)
  governs dependencies, provenance, secrets, and release integrity.
- [Agent-native product interfaces](agent-native-product-interfaces.md)
  governs operational MCP tools and equivalent CLI workflows.
- [Software verification and release engineering](software-verification-and-release-engineering.md)
  governs release evidence and validation depth.
- Repository operability for AI agents is deliberately outside this standard.

## References

- [The `llms.txt` proposal](https://llmstxt.org/)
- [Robots Exclusion Protocol, RFC 9309](https://www.rfc-editor.org/rfc/rfc9309.html)

## Agent briefing

> Every public product offers three AI access paths over one released truth: a
> canonical AI-readable Markdown product guide, public generative-discovery
> files, and a read-only MCP resource server. Documentation stays usable
> without MCP; MCP tools exist only for real product functionality; mutating
> tools need explicit authorization design. AI-facing surfaces are versioned
> and verified with the release.
