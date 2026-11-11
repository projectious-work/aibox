---
apiVersion: processkit.projectious.work/v2
kind: DecisionRecord
metadata:
  id: DEC-20261002_0037-PatientDawn-use-grouped-first-party-features-for
  created: '2026-10-02T00:37:01+00:00'
  updated: '2026-10-02T00:46:54+00:00'
spec:
  title: Use grouped first-party Features for retained v0 tools
  state: superseded
  decision: For each retained individual v0 ledger tool, provide a separately selectable
    aibox-owned Dev Container Feature in a category namespace. Provide group-level
    all Features using standard dependsOn relationships. A Feature installer follows
    the producer's documented installation method, using qualified native packages
    or release artifacts as that documentation directs. Keep installation in standard
    devcontainer.json Features, not customizations.aibox. Do not publish the collection
    to the public Feature index yet.
  context: V1-05 inventories 98 individual tools across v0 addon recipes. Source review
    qualified only four upstream Features. The owner accepted individually selectable
    first-party Features organized in groups such as ai/tau, with group all bundles,
    and required each installer to follow the tool producer's documented installation
    method. Per-tool retain/remove decisions and unsupported-distribution policy remain
    under review.
  rationale: This preserves independent selection and familiar v0 catalog groups while
    exposing standard Dev Container Feature artifacts usable without aibox. Producer
    documentation keeps installation provenance explicit. Standard dependsOn supplies
    bundles without a custom addon resolver.
  alternatives:
  - option: Use only direct upstream Features or native Dockerfile packages for most
      tools
    outcome: Superseded as the default product catalog surface; reusable first-party
      Feature interfaces are now desired for retained v0 tools.
  - option: One monolithic Feature with boolean switches for every group
    outcome: Rejected because it obscures independently selected tools and makes version/disable
      behavior harder to verify.
  - option: Put tool enablement under customizations.aibox
    outcome: Rejected because upstream Dev Container tools do not interpret it as
      installation intent.
  consequences: Revise V1-05 spec, roadmap, selection ledger and evidence plan. Feature
    IDs and bundles need per-tool version/disable, provenance and architecture checks.
    Existing upstream Feature reuse may occur inside first-party Features when appropriate.
    Audio bundle boundaries and unsupported distro behavior still require a design
    answer.
  decided_at: '2026-10-02T00:37:01+00:00'
  superseded_by: DEC-20261002_0046-SmoothLotus-package-coupled-addon-groups-atomically-and
---
