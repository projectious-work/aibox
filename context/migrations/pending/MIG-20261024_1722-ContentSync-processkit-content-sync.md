---
apiVersion: processkit.projectious.work/v1
kind: Migration
metadata:
  id: MIG-20261024_1722-ContentSync-processkit-content-sync
  created: 2026-10-24T17:22:22Z
spec:
  source: processkit
  source_url: "https://github.com/projectious-work/processkit.git"
  from_version: v0.28.1
  to_version: v0.28.8
  state: pending
  generated_by: aibox apply
  generated_at: 2026-10-24T17:22:22Z
  summary: 0 changed upstream, 1 conflicts, 9 new, 0 removed, 0 stale-removed (10 groups affected)
  affected_groups:
    - AGENTS
    - context/artifacts
    - lib
    - skills/data-ai
    - skills/design
    - skills/devops
    - skills/documents
    - skills/engineering
    - skills/processkit
    - skills/product
  affected_files:
    - { path: AGENTS.md, classification: conflict }
    - { path: context/artifacts/ART-20260503_1424-ModelSpec-alibaba-qwen3-6-flash.md, classification: new-upstream }
    - { path: context/artifacts/ART-20260503_1424-ModelSpec-alibaba-qwen3-7-max.md, classification: new-upstream }
    - { path: context/artifacts/ART-20260503_1424-ModelSpec-alibaba-qwen3-8-max-preview.md, classification: new-upstream }
    - { path: context/artifacts/ART-20260503_1424-ModelSpec-google-gemini-3-6-flash.md, classification: new-upstream }
    - { path: context/artifacts/ART-20260503_1424-ModelSpec-minimax-minimax-m3.md, classification: new-upstream }
    - { path: context/artifacts/ART-20260503_1424-ModelSpec-moonshot-kimi-k3.md, classification: new-upstream }
    - { path: context/artifacts/ART-20260503_1424-ModelSpec-openai-gpt-5-3-codex-spark.md, classification: new-upstream }
    - { path: context/artifacts/ART-20260503_1424-ModelSpec-subquadratic-subq-1-1-small.md, classification: new-upstream }
    - { path: context/artifacts/ART-20260503_1424-ModelSpec-xiaomi-mimo-7b.md, classification: new-upstream }
    - { path: context/skills/_lib/processkit/paths.py, classification: changed-locally-only }
    - { path: context/skills/_lib/processkit/schema.py, classification: changed-locally-only }
    - { path: context/skills/data-ai/ai-fundamentals/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/data-ai/ai-fundamentals/references/math-foundations.md, classification: changed-locally-only }
    - { path: context/skills/data-ai/ai-fundamentals/references/ml-concepts.md, classification: changed-locally-only }
    - { path: context/skills/data-ai/data-pipeline/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/data-ai/data-quality/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/data-ai/data-science/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/data-ai/data-science/references/statistical-methods.md, classification: changed-locally-only }
    - { path: context/skills/data-ai/data-science/references/tidy-data-principles.md, classification: changed-locally-only }
    - { path: context/skills/data-ai/data-science/references/visualization-guidelines.md, classification: changed-locally-only }
    - { path: context/skills/data-ai/embedding-vectordb/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/data-ai/llm-evaluation/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/data-ai/ml-pipeline/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/data-ai/ml-pipeline/references/pipeline-stages.md, classification: changed-locally-only }
    - { path: context/skills/data-ai/pandas-polars/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/data-ai/pandas-polars/references/api-comparison.md, classification: changed-locally-only }
    - { path: context/skills/data-ai/rag-engineering/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/data-ai/rag-engineering/references/chunking-strategies.md, classification: changed-locally-only }
    - { path: context/skills/data-ai/rag-engineering/references/evaluation.md, classification: changed-locally-only }
    - { path: context/skills/data-ai/rag-engineering/references/retrieval-patterns.md, classification: changed-locally-only }
    - { path: context/skills/design/excalidraw/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/design/excalidraw/references/json-schema.md, classification: changed-locally-only }
    - { path: context/skills/design/frontend-design/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/design/frontend-design/references/accessibility-checklist.md, classification: changed-locally-only }
    - { path: context/skills/design/mobile-app-design/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/design/mobile-app-design/references/platform-guidelines.md, classification: changed-locally-only }
    - { path: context/skills/design/seo-optimization/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/design/seo-optimization/references/technical-seo-checklist.md, classification: changed-locally-only }
    - { path: context/skills/devops/alerting-oncall/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/devops/ci-cd-setup/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/devops/container-orchestration/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/devops/container-orchestration/references/compose-patterns.md, classification: changed-locally-only }
    - { path: context/skills/devops/distributed-tracing/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/devops/dns-networking/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/devops/dns-networking/references/protocol-reference.md, classification: changed-locally-only }
    - { path: context/skills/devops/dns-networking/references/troubleshooting-tools.md, classification: changed-locally-only }
    - { path: context/skills/devops/dockerfile-review/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/devops/incident-response/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/devops/kubernetes-basics/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/devops/kubernetes-basics/references/cluster-architecture.md, classification: changed-locally-only }
    - { path: context/skills/devops/kubernetes-basics/references/resource-cheatsheet.md, classification: changed-locally-only }
    - { path: context/skills/devops/kubernetes-basics/references/troubleshooting.md, classification: changed-locally-only }
    - { path: context/skills/devops/linux-administration/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/devops/linux-administration/references/commands-cheatsheet.md, classification: changed-locally-only }
    - { path: context/skills/devops/linux-administration/references/systemd-reference.md, classification: changed-locally-only }
    - { path: context/skills/devops/logging-strategy/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/devops/logging-strategy/references/structured-logging.md, classification: changed-locally-only }
    - { path: context/skills/devops/metrics-management/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/devops/metrics-management/assets/metric-spec.yaml, classification: changed-locally-only }
    - { path: context/skills/devops/metrics-monitoring/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/devops/metrics-monitoring/references/metric-types.md, classification: changed-locally-only }
    - { path: context/skills/devops/postmortem-writing/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/devops/terraform-basics/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/documents/data-storytelling/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/documents/data-visualization/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/documents/data-visualization/references/chart-selection.md, classification: changed-locally-only }
    - { path: context/skills/documents/docx-authoring/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/documents/latex-authoring/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/documents/latex-authoring/references/math-reference.md, classification: changed-locally-only }
    - { path: context/skills/documents/latex-authoring/references/packages.md, classification: changed-locally-only }
    - { path: context/skills/documents/latex-authoring/references/tikz-reference.md, classification: changed-locally-only }
    - { path: context/skills/documents/pdf-workflow/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/documents/pptx-authoring/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/documents/xlsx-modeling/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/api-design/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/api-design/references/openapi-patterns.md, classification: changed-locally-only }
    - { path: context/skills/engineering/api-design/references/rest-conventions.md, classification: changed-locally-only }
    - { path: context/skills/engineering/auth-patterns/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/auth-patterns/references/jwt-reference.md, classification: changed-locally-only }
    - { path: context/skills/engineering/auth-patterns/references/oauth-flows.md, classification: changed-locally-only }
    - { path: context/skills/engineering/caching-strategies/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/code-generation/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/concurrency-patterns/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/concurrency-patterns/references/patterns-catalog.md, classification: changed-locally-only }
    - { path: context/skills/engineering/database-migration/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/database-modeling/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/database-modeling/references/modeling-patterns.md, classification: changed-locally-only }
    - { path: context/skills/engineering/debugging/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/dependency-audit/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/dependency-management/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/domain-driven-design/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/domain-driven-design/references/ddd-building-blocks.md, classification: changed-locally-only }
    - { path: context/skills/engineering/error-handling/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/event-driven-architecture/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/event-driven-architecture/references/messaging-patterns.md, classification: changed-locally-only }
    - { path: context/skills/engineering/fastapi-patterns/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/fastapi-patterns/references/endpoint-patterns.md, classification: changed-locally-only }
    - { path: context/skills/engineering/flutter-development/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/flutter-development/references/widget-catalog.md, classification: changed-locally-only }
    - { path: context/skills/engineering/go-conventions/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/go-conventions/references/go-patterns.md, classification: changed-locally-only }
    - { path: context/skills/engineering/graphql-patterns/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/grpc-protobuf/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/grpc-protobuf/references/proto-conventions.md, classification: changed-locally-only }
    - { path: context/skills/engineering/integration-testing/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/integration-testing/references/test-fixtures.md, classification: changed-locally-only }
    - { path: context/skills/engineering/java-patterns/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/load-testing/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/microservice-creation/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/nosql-patterns/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/performance-profiling/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/performance-profiling/references/profiling-tools.md, classification: changed-locally-only }
    - { path: context/skills/engineering/pixijs-gamedev/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/pixijs-gamedev/references/api-cheatsheet.md, classification: changed-locally-only }
    - { path: context/skills/engineering/python-best-practices/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/refactoring/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/refactoring/references/code-smells.md, classification: changed-locally-only }
    - { path: context/skills/engineering/refactoring/references/gof-patterns.md, classification: changed-locally-only }
    - { path: context/skills/engineering/reflex-python/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/reflex-python/references/component-reference.md, classification: changed-locally-only }
    - { path: context/skills/engineering/rust-conventions/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/secret-management/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/secure-coding/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/secure-coding/references/owasp-checklist.md, classification: changed-locally-only }
    - { path: context/skills/engineering/shell-scripting/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/shell-scripting/references/bash-patterns.md, classification: changed-locally-only }
    - { path: context/skills/engineering/software-architecture/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/software-architecture/references/patterns.md, classification: changed-locally-only }
    - { path: context/skills/engineering/software-modularization/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/sql-patterns/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/sql-patterns/references/query-patterns.md, classification: changed-locally-only }
    - { path: context/skills/engineering/sql-patterns/references/schema-design.md, classification: changed-locally-only }
    - { path: context/skills/engineering/sql-style-guide/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/system-design/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/system-design/references/estimation-cheatsheet.md, classification: changed-locally-only }
    - { path: context/skills/engineering/tailwind/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/tailwind/references/cheatsheet.md, classification: changed-locally-only }
    - { path: context/skills/engineering/tdd-workflow/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/testing-strategy/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/threat-modeling/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/typescript-patterns/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/engineering/webhook-integration/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/processkit/repository-portfolio-review/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/product/email-drafter/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/product/estimation-planning/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/product/onboarding-guide/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/product/prd-writing/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/product/sprint-retrospective/SKILL.md, classification: changed-locally-only }
    - { path: context/skills/product/user-research/SKILL.md, classification: changed-locally-only }
---

# Migration MIG-20261024_1722-ContentSync-processkit-content-sync

From `v0.28.1` to `v0.28.8` (source: `https://github.com/projectious-work/processkit.git`).

0 changed upstream, 1 conflicts, 9 new, 0 removed, 0 stale-removed (10 groups affected)

## Counts

- unchanged: 572
- changed-locally-only: 140
- changed-upstream-only: 0
- conflict: 1
- new-upstream: 9
- removed-upstream: 0
- removed-upstream-stale: 0

## Changes by group

### AGENTS

**conflict**

- `AGENTS.md` → `AGENTS.md`

### context/artifacts

**new-upstream**

- `context/artifacts/ART-20260503_1424-ModelSpec-subquadratic-subq-1-1-small.md` → `context/artifacts/ART-20260503_1424-ModelSpec-subquadratic-subq-1-1-small.md`
- `context/artifacts/ART-20260503_1424-ModelSpec-alibaba-qwen3-7-max.md` → `context/artifacts/ART-20260503_1424-ModelSpec-alibaba-qwen3-7-max.md`
- `context/artifacts/ART-20260503_1424-ModelSpec-alibaba-qwen3-6-flash.md` → `context/artifacts/ART-20260503_1424-ModelSpec-alibaba-qwen3-6-flash.md`
- `context/artifacts/ART-20260503_1424-ModelSpec-moonshot-kimi-k3.md` → `context/artifacts/ART-20260503_1424-ModelSpec-moonshot-kimi-k3.md`
- `context/artifacts/ART-20260503_1424-ModelSpec-minimax-minimax-m3.md` → `context/artifacts/ART-20260503_1424-ModelSpec-minimax-minimax-m3.md`
- `context/artifacts/ART-20260503_1424-ModelSpec-xiaomi-mimo-7b.md` → `context/artifacts/ART-20260503_1424-ModelSpec-xiaomi-mimo-7b.md`
- `context/artifacts/ART-20260503_1424-ModelSpec-google-gemini-3-6-flash.md` → `context/artifacts/ART-20260503_1424-ModelSpec-google-gemini-3-6-flash.md`
- `context/artifacts/ART-20260503_1424-ModelSpec-openai-gpt-5-3-codex-spark.md` → `context/artifacts/ART-20260503_1424-ModelSpec-openai-gpt-5-3-codex-spark.md`
- `context/artifacts/ART-20260503_1424-ModelSpec-alibaba-qwen3-8-max-preview.md` → `context/artifacts/ART-20260503_1424-ModelSpec-alibaba-qwen3-8-max-preview.md`

### lib

**changed-locally-only**

- `context/skills/_lib/processkit/paths.py` → `context/skills/_lib/processkit/paths.py`
- `context/skills/_lib/processkit/schema.py` → `context/skills/_lib/processkit/schema.py`

### skills/data-ai

**changed-locally-only**

- `context/skills/data-ai/pandas-polars/references/api-comparison.md` → `context/skills/data-ai/pandas-polars/references/api-comparison.md`
- `context/skills/data-ai/pandas-polars/SKILL.md` → `context/skills/data-ai/pandas-polars/SKILL.md`
- `context/skills/data-ai/llm-evaluation/SKILL.md` → `context/skills/data-ai/llm-evaluation/SKILL.md`
- `context/skills/data-ai/data-science/references/tidy-data-principles.md` → `context/skills/data-ai/data-science/references/tidy-data-principles.md`
- `context/skills/data-ai/data-science/references/visualization-guidelines.md` → `context/skills/data-ai/data-science/references/visualization-guidelines.md`
- `context/skills/data-ai/data-science/references/statistical-methods.md` → `context/skills/data-ai/data-science/references/statistical-methods.md`
- `context/skills/data-ai/data-science/SKILL.md` → `context/skills/data-ai/data-science/SKILL.md`
- `context/skills/data-ai/data-quality/SKILL.md` → `context/skills/data-ai/data-quality/SKILL.md`
- `context/skills/data-ai/rag-engineering/references/chunking-strategies.md` → `context/skills/data-ai/rag-engineering/references/chunking-strategies.md`
- `context/skills/data-ai/rag-engineering/references/retrieval-patterns.md` → `context/skills/data-ai/rag-engineering/references/retrieval-patterns.md`
- `context/skills/data-ai/rag-engineering/references/evaluation.md` → `context/skills/data-ai/rag-engineering/references/evaluation.md`
- `context/skills/data-ai/rag-engineering/SKILL.md` → `context/skills/data-ai/rag-engineering/SKILL.md`
- `context/skills/data-ai/embedding-vectordb/SKILL.md` → `context/skills/data-ai/embedding-vectordb/SKILL.md`
- `context/skills/data-ai/data-pipeline/SKILL.md` → `context/skills/data-ai/data-pipeline/SKILL.md`
- `context/skills/data-ai/ai-fundamentals/references/math-foundations.md` → `context/skills/data-ai/ai-fundamentals/references/math-foundations.md`
- `context/skills/data-ai/ai-fundamentals/references/ml-concepts.md` → `context/skills/data-ai/ai-fundamentals/references/ml-concepts.md`
- `context/skills/data-ai/ai-fundamentals/SKILL.md` → `context/skills/data-ai/ai-fundamentals/SKILL.md`
- `context/skills/data-ai/ml-pipeline/references/pipeline-stages.md` → `context/skills/data-ai/ml-pipeline/references/pipeline-stages.md`
- `context/skills/data-ai/ml-pipeline/SKILL.md` → `context/skills/data-ai/ml-pipeline/SKILL.md`

### skills/design

**changed-locally-only**

- `context/skills/design/mobile-app-design/references/platform-guidelines.md` → `context/skills/design/mobile-app-design/references/platform-guidelines.md`
- `context/skills/design/mobile-app-design/SKILL.md` → `context/skills/design/mobile-app-design/SKILL.md`
- `context/skills/design/seo-optimization/references/technical-seo-checklist.md` → `context/skills/design/seo-optimization/references/technical-seo-checklist.md`
- `context/skills/design/seo-optimization/SKILL.md` → `context/skills/design/seo-optimization/SKILL.md`
- `context/skills/design/excalidraw/references/json-schema.md` → `context/skills/design/excalidraw/references/json-schema.md`
- `context/skills/design/excalidraw/SKILL.md` → `context/skills/design/excalidraw/SKILL.md`
- `context/skills/design/frontend-design/references/accessibility-checklist.md` → `context/skills/design/frontend-design/references/accessibility-checklist.md`
- `context/skills/design/frontend-design/SKILL.md` → `context/skills/design/frontend-design/SKILL.md`

### skills/devops

**changed-locally-only**

- `context/skills/devops/linux-administration/references/commands-cheatsheet.md` → `context/skills/devops/linux-administration/references/commands-cheatsheet.md`
- `context/skills/devops/linux-administration/references/systemd-reference.md` → `context/skills/devops/linux-administration/references/systemd-reference.md`
- `context/skills/devops/linux-administration/SKILL.md` → `context/skills/devops/linux-administration/SKILL.md`
- `context/skills/devops/logging-strategy/references/structured-logging.md` → `context/skills/devops/logging-strategy/references/structured-logging.md`
- `context/skills/devops/logging-strategy/SKILL.md` → `context/skills/devops/logging-strategy/SKILL.md`
- `context/skills/devops/incident-response/SKILL.md` → `context/skills/devops/incident-response/SKILL.md`
- `context/skills/devops/metrics-monitoring/references/metric-types.md` → `context/skills/devops/metrics-monitoring/references/metric-types.md`
- `context/skills/devops/metrics-monitoring/SKILL.md` → `context/skills/devops/metrics-monitoring/SKILL.md`
- `context/skills/devops/terraform-basics/SKILL.md` → `context/skills/devops/terraform-basics/SKILL.md`
- `context/skills/devops/container-orchestration/references/compose-patterns.md` → `context/skills/devops/container-orchestration/references/compose-patterns.md`
- `context/skills/devops/container-orchestration/SKILL.md` → `context/skills/devops/container-orchestration/SKILL.md`
- `context/skills/devops/kubernetes-basics/references/cluster-architecture.md` → `context/skills/devops/kubernetes-basics/references/cluster-architecture.md`
- `context/skills/devops/kubernetes-basics/references/troubleshooting.md` → `context/skills/devops/kubernetes-basics/references/troubleshooting.md`
- `context/skills/devops/kubernetes-basics/references/resource-cheatsheet.md` → `context/skills/devops/kubernetes-basics/references/resource-cheatsheet.md`
- `context/skills/devops/kubernetes-basics/SKILL.md` → `context/skills/devops/kubernetes-basics/SKILL.md`
- `context/skills/devops/metrics-management/SKILL.md` → `context/skills/devops/metrics-management/SKILL.md`
- `context/skills/devops/metrics-management/assets/metric-spec.yaml` → `context/skills/devops/metrics-management/assets/metric-spec.yaml`
- `context/skills/devops/dns-networking/references/troubleshooting-tools.md` → `context/skills/devops/dns-networking/references/troubleshooting-tools.md`
- `context/skills/devops/dns-networking/references/protocol-reference.md` → `context/skills/devops/dns-networking/references/protocol-reference.md`
- `context/skills/devops/dns-networking/SKILL.md` → `context/skills/devops/dns-networking/SKILL.md`
- `context/skills/devops/postmortem-writing/SKILL.md` → `context/skills/devops/postmortem-writing/SKILL.md`
- `context/skills/devops/dockerfile-review/SKILL.md` → `context/skills/devops/dockerfile-review/SKILL.md`
- `context/skills/devops/distributed-tracing/SKILL.md` → `context/skills/devops/distributed-tracing/SKILL.md`
- `context/skills/devops/ci-cd-setup/SKILL.md` → `context/skills/devops/ci-cd-setup/SKILL.md`
- `context/skills/devops/alerting-oncall/SKILL.md` → `context/skills/devops/alerting-oncall/SKILL.md`

### skills/documents

**changed-locally-only**

- `context/skills/documents/data-storytelling/SKILL.md` → `context/skills/documents/data-storytelling/SKILL.md`
- `context/skills/documents/data-visualization/references/chart-selection.md` → `context/skills/documents/data-visualization/references/chart-selection.md`
- `context/skills/documents/data-visualization/SKILL.md` → `context/skills/documents/data-visualization/SKILL.md`
- `context/skills/documents/docx-authoring/SKILL.md` → `context/skills/documents/docx-authoring/SKILL.md`
- `context/skills/documents/pdf-workflow/SKILL.md` → `context/skills/documents/pdf-workflow/SKILL.md`
- `context/skills/documents/xlsx-modeling/SKILL.md` → `context/skills/documents/xlsx-modeling/SKILL.md`
- `context/skills/documents/pptx-authoring/SKILL.md` → `context/skills/documents/pptx-authoring/SKILL.md`
- `context/skills/documents/latex-authoring/references/math-reference.md` → `context/skills/documents/latex-authoring/references/math-reference.md`
- `context/skills/documents/latex-authoring/references/packages.md` → `context/skills/documents/latex-authoring/references/packages.md`
- `context/skills/documents/latex-authoring/references/tikz-reference.md` → `context/skills/documents/latex-authoring/references/tikz-reference.md`
- `context/skills/documents/latex-authoring/SKILL.md` → `context/skills/documents/latex-authoring/SKILL.md`

### skills/engineering

**changed-locally-only**

- `context/skills/engineering/grpc-protobuf/references/proto-conventions.md` → `context/skills/engineering/grpc-protobuf/references/proto-conventions.md`
- `context/skills/engineering/grpc-protobuf/SKILL.md` → `context/skills/engineering/grpc-protobuf/SKILL.md`
- `context/skills/engineering/database-modeling/references/modeling-patterns.md` → `context/skills/engineering/database-modeling/references/modeling-patterns.md`
- `context/skills/engineering/database-modeling/SKILL.md` → `context/skills/engineering/database-modeling/SKILL.md`
- `context/skills/engineering/go-conventions/references/go-patterns.md` → `context/skills/engineering/go-conventions/references/go-patterns.md`
- `context/skills/engineering/go-conventions/SKILL.md` → `context/skills/engineering/go-conventions/SKILL.md`
- `context/skills/engineering/secret-management/SKILL.md` → `context/skills/engineering/secret-management/SKILL.md`
- `context/skills/engineering/error-handling/SKILL.md` → `context/skills/engineering/error-handling/SKILL.md`
- `context/skills/engineering/nosql-patterns/SKILL.md` → `context/skills/engineering/nosql-patterns/SKILL.md`
- `context/skills/engineering/python-best-practices/SKILL.md` → `context/skills/engineering/python-best-practices/SKILL.md`
- `context/skills/engineering/pixijs-gamedev/references/api-cheatsheet.md` → `context/skills/engineering/pixijs-gamedev/references/api-cheatsheet.md`
- `context/skills/engineering/pixijs-gamedev/SKILL.md` → `context/skills/engineering/pixijs-gamedev/SKILL.md`
- `context/skills/engineering/debugging/SKILL.md` → `context/skills/engineering/debugging/SKILL.md`
- `context/skills/engineering/auth-patterns/references/oauth-flows.md` → `context/skills/engineering/auth-patterns/references/oauth-flows.md`
- `context/skills/engineering/auth-patterns/references/jwt-reference.md` → `context/skills/engineering/auth-patterns/references/jwt-reference.md`
- `context/skills/engineering/auth-patterns/SKILL.md` → `context/skills/engineering/auth-patterns/SKILL.md`
- `context/skills/engineering/api-design/references/rest-conventions.md` → `context/skills/engineering/api-design/references/rest-conventions.md`
- `context/skills/engineering/api-design/references/openapi-patterns.md` → `context/skills/engineering/api-design/references/openapi-patterns.md`
- `context/skills/engineering/api-design/SKILL.md` → `context/skills/engineering/api-design/SKILL.md`
- `context/skills/engineering/integration-testing/references/test-fixtures.md` → `context/skills/engineering/integration-testing/references/test-fixtures.md`
- `context/skills/engineering/integration-testing/SKILL.md` → `context/skills/engineering/integration-testing/SKILL.md`
- `context/skills/engineering/fastapi-patterns/references/endpoint-patterns.md` → `context/skills/engineering/fastapi-patterns/references/endpoint-patterns.md`
- `context/skills/engineering/fastapi-patterns/SKILL.md` → `context/skills/engineering/fastapi-patterns/SKILL.md`
- `context/skills/engineering/dependency-audit/SKILL.md` → `context/skills/engineering/dependency-audit/SKILL.md`
- `context/skills/engineering/graphql-patterns/SKILL.md` → `context/skills/engineering/graphql-patterns/SKILL.md`
- `context/skills/engineering/threat-modeling/SKILL.md` → `context/skills/engineering/threat-modeling/SKILL.md`
- `context/skills/engineering/event-driven-architecture/references/messaging-patterns.md` → `context/skills/engineering/event-driven-architecture/references/messaging-patterns.md`
- `context/skills/engineering/event-driven-architecture/SKILL.md` → `context/skills/engineering/event-driven-architecture/SKILL.md`
- `context/skills/engineering/tailwind/references/cheatsheet.md` → `context/skills/engineering/tailwind/references/cheatsheet.md`
- `context/skills/engineering/tailwind/SKILL.md` → `context/skills/engineering/tailwind/SKILL.md`
- `context/skills/engineering/performance-profiling/references/profiling-tools.md` → `context/skills/engineering/performance-profiling/references/profiling-tools.md`
- `context/skills/engineering/performance-profiling/SKILL.md` → `context/skills/engineering/performance-profiling/SKILL.md`
- `context/skills/engineering/code-generation/SKILL.md` → `context/skills/engineering/code-generation/SKILL.md`
- `context/skills/engineering/caching-strategies/SKILL.md` → `context/skills/engineering/caching-strategies/SKILL.md`
- `context/skills/engineering/shell-scripting/references/bash-patterns.md` → `context/skills/engineering/shell-scripting/references/bash-patterns.md`
- `context/skills/engineering/shell-scripting/SKILL.md` → `context/skills/engineering/shell-scripting/SKILL.md`
- `context/skills/engineering/concurrency-patterns/references/patterns-catalog.md` → `context/skills/engineering/concurrency-patterns/references/patterns-catalog.md`
- `context/skills/engineering/concurrency-patterns/SKILL.md` → `context/skills/engineering/concurrency-patterns/SKILL.md`
- `context/skills/engineering/secure-coding/references/owasp-checklist.md` → `context/skills/engineering/secure-coding/references/owasp-checklist.md`
- `context/skills/engineering/secure-coding/SKILL.md` → `context/skills/engineering/secure-coding/SKILL.md`
- `context/skills/engineering/dependency-management/SKILL.md` → `context/skills/engineering/dependency-management/SKILL.md`
- `context/skills/engineering/testing-strategy/SKILL.md` → `context/skills/engineering/testing-strategy/SKILL.md`
- `context/skills/engineering/webhook-integration/SKILL.md` → `context/skills/engineering/webhook-integration/SKILL.md`
- `context/skills/engineering/tdd-workflow/SKILL.md` → `context/skills/engineering/tdd-workflow/SKILL.md`
- `context/skills/engineering/microservice-creation/SKILL.md` → `context/skills/engineering/microservice-creation/SKILL.md`
- `context/skills/engineering/software-architecture/references/patterns.md` → `context/skills/engineering/software-architecture/references/patterns.md`
- `context/skills/engineering/software-architecture/SKILL.md` → `context/skills/engineering/software-architecture/SKILL.md`
- `context/skills/engineering/system-design/references/estimation-cheatsheet.md` → `context/skills/engineering/system-design/references/estimation-cheatsheet.md`
- `context/skills/engineering/system-design/SKILL.md` → `context/skills/engineering/system-design/SKILL.md`
- `context/skills/engineering/java-patterns/SKILL.md` → `context/skills/engineering/java-patterns/SKILL.md`
- `context/skills/engineering/load-testing/SKILL.md` → `context/skills/engineering/load-testing/SKILL.md`
- `context/skills/engineering/reflex-python/references/component-reference.md` → `context/skills/engineering/reflex-python/references/component-reference.md`
- `context/skills/engineering/reflex-python/SKILL.md` → `context/skills/engineering/reflex-python/SKILL.md`
- `context/skills/engineering/refactoring/references/gof-patterns.md` → `context/skills/engineering/refactoring/references/gof-patterns.md`
- `context/skills/engineering/refactoring/references/code-smells.md` → `context/skills/engineering/refactoring/references/code-smells.md`
- `context/skills/engineering/refactoring/SKILL.md` → `context/skills/engineering/refactoring/SKILL.md`
- `context/skills/engineering/software-modularization/SKILL.md` → `context/skills/engineering/software-modularization/SKILL.md`
- `context/skills/engineering/sql-style-guide/SKILL.md` → `context/skills/engineering/sql-style-guide/SKILL.md`
- `context/skills/engineering/database-migration/SKILL.md` → `context/skills/engineering/database-migration/SKILL.md`
- `context/skills/engineering/flutter-development/references/widget-catalog.md` → `context/skills/engineering/flutter-development/references/widget-catalog.md`
- `context/skills/engineering/flutter-development/SKILL.md` → `context/skills/engineering/flutter-development/SKILL.md`
- `context/skills/engineering/rust-conventions/SKILL.md` → `context/skills/engineering/rust-conventions/SKILL.md`
- `context/skills/engineering/domain-driven-design/references/ddd-building-blocks.md` → `context/skills/engineering/domain-driven-design/references/ddd-building-blocks.md`
- `context/skills/engineering/domain-driven-design/SKILL.md` → `context/skills/engineering/domain-driven-design/SKILL.md`
- `context/skills/engineering/typescript-patterns/SKILL.md` → `context/skills/engineering/typescript-patterns/SKILL.md`
- `context/skills/engineering/sql-patterns/references/query-patterns.md` → `context/skills/engineering/sql-patterns/references/query-patterns.md`
- `context/skills/engineering/sql-patterns/references/schema-design.md` → `context/skills/engineering/sql-patterns/references/schema-design.md`
- `context/skills/engineering/sql-patterns/SKILL.md` → `context/skills/engineering/sql-patterns/SKILL.md`

### skills/processkit

**changed-locally-only**

- `context/skills/processkit/repository-portfolio-review/SKILL.md` → `context/skills/processkit/repository-portfolio-review/SKILL.md`

### skills/product

**changed-locally-only**

- `context/skills/product/estimation-planning/SKILL.md` → `context/skills/product/estimation-planning/SKILL.md`
- `context/skills/product/sprint-retrospective/SKILL.md` → `context/skills/product/sprint-retrospective/SKILL.md`
- `context/skills/product/user-research/SKILL.md` → `context/skills/product/user-research/SKILL.md`
- `context/skills/product/onboarding-guide/SKILL.md` → `context/skills/product/onboarding-guide/SKILL.md`
- `context/skills/product/email-drafter/SKILL.md` → `context/skills/product/email-drafter/SKILL.md`
- `context/skills/product/prd-writing/SKILL.md` → `context/skills/product/prd-writing/SKILL.md`

