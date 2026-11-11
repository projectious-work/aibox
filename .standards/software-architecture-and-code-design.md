# Software architecture and code design standard

## Purpose and scope

This standard defines the default structure of projectious.work software:
how boundaries, modules, and code are shaped, and how architecture is
documented. It applies to applications with domain logic and external
integrations — CLIs, services, workers, and MCP servers. Libraries apply the
inner rules and leave composition to their host.

Terms in *italics* are [Semantic Anchors](https://github.com/llm-coding/Semantic-Anchors/tree/6554641a11fdf99fdbdf09cc0602268e3db88966); see
[README](README.md#conventions).

## Principles

- The core does not know its delivery mechanism or infrastructure.
- Each architectural frame has one job; frames are not mixed into vague
  hybrids.
- Module boundaries are as deliberate as service boundaries.
- Architecture documentation records what a reader cannot recover from the
  code, especially decisions and their reasons.
- The simplest design that meets the specification wins (*KISS Principle*);
  nothing is built for a need the specification does not state (*YAGNI (You
  Aren't Gonna Need It)*).

## Boundaries: Hexagonal Architecture

*Hexagonal Architecture* governs the system boundary:

- Domain and application form the core. Public functionality lives there.
- The core defines ports (interfaces) for everything external: CLI and MCP
  entry points, filesystems, processes, container engines, Git, network
  services, model providers, clocks, and randomness.
- Adapters implement ports. Dependencies point inward only; the core never
  imports an adapter.
- A composition root at the edge (the program's `main`, server bootstrap, or
  test harness) wires adapters to ports.

## Inside the core: Clean Architecture

*Clean Architecture* governs the inside of the core: business rules live in
entities and use cases; frameworks, databases, transport, and serialization
are details behind interfaces.

## Modules

- Cut the core along bounded contexts (*Domain-Driven Design according to
  Evans*, tactical patterns). Use one Ubiquitous Language: the same terms in
  specification, code, tests, CLI output, and documentation.
- Modules are *Deep Modules*: a small interface over substantial behavior.
  Do not split code into shallow pass-through layers.
- Module boundaries follow ownership, because the structure will mirror the
  owning teams (*Conway's Law*): one owner per module. A repository holds one
  concern.
- Modules are feature-centric and own their data and state. They communicate
  only through explicit module APIs and never reach into another module's
  internal types, tables, or files.
- Module and package dependencies are acyclic.
- Every piece of knowledge has one authoritative representation (*DRY (Don't
  Repeat Yourself)*). Code that only looks alike is not duplication.

## Code conventions

- The repository's formatter and linter configuration is authoritative. Code
  that does not pass them is not finished. A repository without a linter
  configuration enforces the language's standard formatter and static checks
  (for example `gofmt` and `go vet`) in its verification scripts; adopting a
  linter commits its configuration and names the tool in the project
  specification.
- Language idioms: *Effective Go* for Go, *Effective Python* for Python, the
  Rust API Guidelines for Rust. The repository's configuration wins over the
  idiom source where they differ.
- *SOLID Principles* apply at type and module level; code aims at the *CUPID
  Properties*. Functions keep a *Single Level of Abstraction Principle*: a use
  case reads as calls to ports and domain operations; details move into named
  functions.
- Name a recurring structure by its *GoF Design Patterns* name (an adapter is
  an Adapter). Introduce a pattern only when it removes a present problem.
- Core operations state preconditions, postconditions, and invariants
  (*Design by Contract*); a violation inside the core is a programming error
  and fails fast. Adapters validate everything that crosses a trust boundary
  (*Defensive Programming according to McConnell*) and never pass unchecked
  input inward.

## Testing consequences

The core is tested through its ports with in-memory fakes; adapters receive
their own integration tests against the real dependency. The test style
follows the
[software verification standard](software-verification-and-release-engineering.md)
(*TDD, Chicago School*).

## Architecture documentation

- Structure per *arc42 Architecture Documentation*. Fill only sections that
  carry information; at minimum context and scope (§3), building blocks (§5),
  crosscutting concepts (§8), architecture decisions (§9), and risks and
  technical debt (§11). Quality goals and requirements, where documented, use
  the *ISO/IEC 25010* characteristics.
- Diagrams per *C4-Diagrams*: context and container views for every product,
  component views where a container's structure is not obvious.
- Architecture decisions as *ADR according to Nygard* (context, decision,
  consequences), stored next to the code or in the project's decision store.
  Superseded decisions are marked, never deleted. A decision about technology
  or technique compares the options with the current *SOTA (State-of-the-Art)*
  and dates that comparison.
- Before removing a boundary, check, or rule whose purpose is unclear, find
  the reason it exists (*Chesterton's Fence*), usually in an ADR.

## Agent briefing

> Ports in the core, adapters at the edge, dependencies inward (*Hexagonal
> Architecture*); rules in entities and use cases (*Clean Architecture*);
> feature *Deep Modules* own their data and talk only through module APIs.
> *KISS Principle*, *YAGNI*; *Design by Contract* inside, defensive checks at
> trust boundaries. The formatter and linter decide style. Document with
> *arc42*, *C4-Diagrams*, and *ADR according to Nygard*.
