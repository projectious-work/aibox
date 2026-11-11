---
apiVersion: processkit.projectious.work/v2
kind: DecisionRecord
metadata:
  id: DEC-20261002_0054-AmberCompass-qualify-feature-compatibility-per-tool-version
  created: '2026-10-02T00:54:23+00:00'
spec:
  title: Qualify Feature compatibility per tool version and container base
  state: accepted
  decision: For each aibox Feature, define the supported matrix by tool version, container
    distribution and release, and architecture. The installer selects an applicable
    producer-documented canonical installation method for each supported matrix cell.
    If no such method exists, an official distribution package may be used after verifying
    exact requested tool version availability in the selected repository/release/architecture,
    provenance, dependencies, successful build and runtime behavior. A package name
    appearing in an APT index alone does not qualify all Debian-derived distributions.
    Unsupported combinations fail clearly; VS Code and aibox do not silently substitute
    versions or bases.
  context: The owner accepted that first-party Features own their compatibility decisions,
    as Dev Container Features do in VS Code. A supported tool version can require
    a different installation method or different supported bases from another version.
    When the tool producer offers no documented route for a base, the owner permits
    a tested package-manager fallback from that distribution.
  rationale: Compatibility is a property of the Feature's actual installer and selected
    version, not of VS Code or a mandatory Microsoft base image. Different tool releases
    may need different installers. Distribution repositories differ by release, architecture,
    sources and package version; exact-cell verification is needed for reliable version
    selection.
  alternatives:
  - option: Require a single base image for all Features
    outcome: Rejected because standard Dev Container Features compose with user-selected
      images and can declare narrower support.
  - option: Treat any APT package index entry as support for every Debian derivative
    outcome: Rejected because repository contents, package versions and dependencies
      vary by distribution release and architecture.
  - option: Use arbitrary undocumented installer fallback
    outcome: Rejected because provenance and behavior would not meet the producer-installation
      principle.
  consequences: V1-05 must record per-version installation routes, exact distro/release/architecture
    support, package-manager fallback evidence and explicit unsupported outcomes for
    each Feature. Build tests must exercise supported routes and expected failures.
    The prior blanket rule that absence of a producer-documented method makes a base
    unsupported is superseded by the verified distribution-package fallback.
  decided_at: '2026-10-02T00:54:23+00:00'
---
