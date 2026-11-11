---
title: Architecture decisions
description: "Index of the aibox v1 decision records that shape the architecture."
weight: 40
---

# Architecture decisions

Architecture decisions are recorded as processkit DecisionRecords
(`DEC-*`) in
[`context/decisions/`](https://github.com/projectious-work/aibox/tree/v1.x-dev/context/decisions).
Each record follows the ADR structure by Michael Nygard: context, decision,
consequences. The record is authoritative; the one-line summaries below only
help you find it. To change a decision, record a new one that supersedes it.

This index lists the decisions that the specification and the phase dev
notes cite. A later phase adds the decisions it relies on.

## Tool Feature architecture

| Decision | Summary |
|---|---|
| [DEC-20261002_0046-SmoothLotus](https://github.com/projectious-work/aibox/blob/v1.x-dev/context/decisions/DEC-20261002_0046-SmoothLotus-package-coupled-addon-groups-atomically-and.md) | Package a tightly coupled tool group as one atomic Feature; give a semantic group one Feature per tool plus an `all` bundle. |
| [DEC-20261002_0054-AmberCompass](https://github.com/projectious-work/aibox/blob/v1.x-dev/context/decisions/DEC-20261002_0054-AmberCompass-qualify-feature-compatibility-per-tool-version.md) | Qualify Feature compatibility per tool version, distro release and architecture cell. |
| DEC-20261001_1954-SmoothTrout | Earlier restriction on routine aibox-owned installer Features; superseded by the grouped Feature direction (chapter 8). Not in this repository's decision log. |
| [DEC-20261029_0802-DeepOtter](https://github.com/projectious-work/aibox/blob/v1.x-dev/context/decisions/DEC-20261029_0802-DeepOtter-v1-05-contract-scope-trixie-two.md) | V1-05 scope: Debian trixie on `amd64` and `arm64` only, exact pins, upstream Features used directly, separate contract branch. |
| [DEC-20261029_0835-HumbleBridge](https://github.com/projectious-work/aibox/blob/v1.x-dev/context/decisions/DEC-20261029_0835-HumbleBridge-v1-05-owner-row-dispositions-first.md) | Owner disposition of every v0 tool row, including removals. |
| [DEC-20261029_1034-ActiveJay](https://github.com/projectious-work/aibox/blob/v1.x-dev/context/decisions/DEC-20261029_1034-ActiveJay-v1-05-catalog-additions-vibe-uv.md) | Catalog additions, uv-managed Python and the security tool groups. |
| [DEC-20261029_1409-DaringTrail](https://github.com/projectious-work/aibox/blob/v1.x-dev/context/decisions/DEC-20261029_1409-DaringTrail-v1-05-contract-review-resolutions-runtime.md) | Review resolutions for runtime dependencies, Debian versions, version counts and approved deviations. |
| [DEC-20261029_2107-QuietSeal](https://github.com/projectious-work/aibox/blob/v1.x-dev/context/decisions/DEC-20261029_2107-QuietSeal-replace-rich-cli-with-glow-and.md) | Replace `rich-cli` with `glow` and `bat` for previews. |

## Delivery process

| Decision | Summary |
|---|---|
| [DEC-20261026_1412-FastForge](https://github.com/projectious-work/aibox/blob/v1.x-dev/context/decisions/DEC-20261026_1412-FastForge-v1-phases-start-from-a-frozen.md) | Each v1 phase starts from a frozen, test-ready contract with parallel acceptance-test and implementation tracks. |
| [DEC-20261029_1515-ShinySpire](https://github.com/projectious-work/aibox/blob/v1.x-dev/context/decisions/DEC-20261029_1515-ShinySpire-freeze-the-v1-05-phase-contract.md) | Freeze the first V1-05 contract. |
| [DEC-20261029_1725-AmberBeam](https://github.com/projectious-work/aibox/blob/v1.x-dev/context/decisions/DEC-20261029_1725-AmberBeam-freeze-v1-05-contract-amendment-1.md) | Freeze V1-05 contract amendment 1 (dependency closure, absence checks, publish prerequisites). |
| [DEC-20261029_2134-NobleBlossom](https://github.com/projectious-work/aibox/blob/v1.x-dev/context/decisions/DEC-20261029_2134-NobleBlossom-freeze-v1-05-contract-amendment-2.md) | Freeze V1-05 contract amendment 2 (rich replaced by glow and bat). |
| [DEC-20261031_1743-MightyRaven](https://github.com/projectious-work/aibox/blob/v1.x-dev/context/decisions/DEC-20261031_1743-MightyRaven-approve-v1-05-contract-amendment-3.md) | Approve amendment 3, which corrects the mdbook and Chromium smoke commands. |
| [DEC-20261031_1750-NeatDaisy](https://github.com/projectious-work/aibox/blob/v1.x-dev/context/decisions/DEC-20261031_1750-NeatDaisy-freeze-v1-05-contract-amendment-3.md) | Freeze V1-05 contract amendment 3 as the current contract commit. |
| [DEC-20261031_1856-CuriousBeacon](https://github.com/projectious-work/aibox/blob/v1.x-dev/context/decisions/DEC-20261031_1856-CuriousBeacon-reduce-ac-tools-run-time-through.md) | Shorten AC-TOOLS runs with a sampled scope for error items; full scope runs on the owner's command. Its layer-reuse part is not implemented (see [risks](risks.md)). |

## Decisions recorded in the specification

The Go core, the delegation of container semantics to the Dev Container CLI,
and the package boundaries were accepted as part of the v1 specification
([chapter 7](https://github.com/projectious-work/aibox/blob/v1.x-dev/spec/v1/07-architecture.md))
rather than as separate `DEC-*` records. The dev notes cite no decision
record for them.
