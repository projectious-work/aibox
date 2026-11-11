---
apiVersion: processkit.projectious.work/v2
kind: DecisionRecord
metadata:
  id: DEC-20261030_2049-StrongOasis-feature-adaptation-guidance-templates-and-recipe
  created: '2026-10-30T20:49:44+00:00'
spec:
  title: Feature adaptation guidance, templates and recipe test move to V1-09 and
    V1-16
  state: accepted
  decision: Feature adaptation stays out of V1-05. V1-09 registers read-only MCP template
    resources alongside guide resources. V1-16 delivers a `features` guide and an
    Adapt a Feature docs page, the `feature`, `minimal` and `curated` templates, and
    an AC-GUIDE acceptance test that copies an aibox Feature, adds a version and builds
    it as a local Feature with the pinned Dev Container CLI without the aibox binary.
    The template store returns content and never writes to a project.
  context: V1-05 Features are standard Dev Container Features. Users can select listed
    versions, but any other change means copying a Feature, which was neither documented
    nor tested. The owner asked for these additions in later phases where they fit,
    and asked whether aibox MCP can deliver descriptions and templates for creating
    and changing Features.
  rationale: Chapter 10 already gives V1-09 the MCP resource surface and V1-16 the
    guide content, and AC-GUIDE already requires recipes to run against a release.
    The V1-05 contract stays frozen and strict version pinning stays in place, with
    copying as the supported way to adapt a Feature.
  consequences: Chapter 10, AC-GUIDE and the V1-09 and V1-16 roadmap entries change;
    V1-16 gains an explicit dependency on V1-05. A user who copies a Feature owns
    its pins and security updates.
  deciders:
  - ACTOR-20260410_2209-SnappyFrog-bernhard
  decided_at: '2026-10-30T20:49:44+00:00'
---
