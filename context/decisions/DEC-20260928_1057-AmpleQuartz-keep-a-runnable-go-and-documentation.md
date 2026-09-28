---
apiVersion: processkit.projectious.work/v2
kind: DecisionRecord
metadata:
  id: DEC-20260928_1057-AmpleQuartz-keep-a-runnable-go-and-documentation
  created: '2026-09-28T10:57:41+00:00'
spec:
  title: Keep a runnable Go and documentation demo after every v1 implementation phase
  state: accepted
  decision: 'Beginning with V1-03, every v1 roadmap phase preserves a runnable aibox
    Go executable and a reproducible, feature-specific customer demonstration paired
    with version-accurate built documentation. V1-03 introduces the first safe read-only
    executable; later phases extend the same binary incrementally. A local demo or
    merged phase is not a shipped release: publication and shipped status still require
    the company release and conformance gates.'
  context: The owner observed that the current V1-03 Go implementation is testable
    but has no executable, while the roadmap deferred the thin CLI to V1-09. This
    blocks customer demonstration and delays feedback despite the desired release-often
    development model.
  rationale: Working software and documentation on every increment expose integration
    and usability gaps early, permit small verified prerelease candidates, and avoid
    a long run of invisible package-only work. Separating demo readiness from publication
    preserves the company evidence and approval requirements.
  alternatives:
  - option: Keep V1-03 through V1-08 as internal packages and first expose a CLI in
      V1-09
    rejected_because: Leaves multiple phases without a customer-runnable product and
      postpones binary-level feedback.
  - option: Mark every successful local demo as shipped
    rejected_because: Contradicts company status semantics and release artifact verification.
  consequences: Revise V1-03 through V1-23 roadmap deliverables, schema and validator;
    add a reusable demo evidence contract to architecture, verification and documentation
    chapters; reconcile alpha/beta gates with incremental supported-scope claims.
    Existing V1-03 work remains in progress until a runnable binary and associated
    evidence exist.
  decided_at: '2026-09-28T10:57:41+00:00'
---
