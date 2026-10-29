---
apiVersion: processkit.projectious.work/v2
kind: DecisionRecord
metadata:
  id: DEC-20261002_0046-SmoothLotus-package-coupled-addon-groups-atomically-and
  created: '2026-10-02T00:46:19+00:00'
  updated: '2026-10-02T00:46:54+00:00'
spec:
  title: Package coupled addon groups atomically and semantic groups as selectable
    Features
  state: accepted
  decision: Classify each retained v0 addon group by cohesion. A coupled group whose
    members are only meaningful together becomes one first-party Feature with no subtool
    references and no all bundle; audio-voice is the example. A semantic collection
    of independently useful tools exposes individual Feature references plus a versioned
    all bundle using standard dependsOn; data-preview is the example. Each install.sh
    follows producer-documented installation methods. No public index publication
    yet. Supported container bases are declared and tested per Feature, not guaranteed
    by VS Code or a mandatory Microsoft base image.
  context: The owner refined the V1-05 grouped Feature design after the prior decision
    had treated every retained tool as an independent Feature. v0 audio-voice combines
    SoX, its PulseAudio plugin, PulseAudio utilities and ALSA compatibility; data-preview
    is a semantic collection of independently useful tools.
  rationale: The boundary follows useful user choices and actual dependency relationships
    rather than mirroring every v0 ledger row as a separate package. It avoids exposing
    broken partial installs while retaining independent selection where tools are
    standalone. Feature support varies by base image, so compatibility must be proven
    per artifact.
  alternatives:
  - option: Publish one Feature per v0 tool plus all bundle for every group
    outcome: Rejected for tightly coupled groups because partial selections would
      not be useful or supportable.
  - option: Publish one Feature per entire v0 addon group
    outcome: Rejected for semantic collections because it would remove useful independent
      selection.
  - option: Require Microsoft Dev Container base images
    outcome: Rejected as a general Feature rule; VS Code accepts user images and each
      Feature has its own base compatibility.
  consequences: Revise V1-05 spec and roadmap, classify all v0 groups, retain traceability
    for every ledger row, and test each atomic or independent Feature against declared
    bases and architectures. The accepted grouped first-party Feature decision is
    superseded where it demanded one Feature per individual tool. The exact support
    matrix and per-group classification still require evidence and owner retain/remove
    dispositions.
  decided_at: '2026-10-02T00:46:19+00:00'
  supersedes: DEC-20261002_0037-PatientDawn-use-grouped-first-party-features-for
---
