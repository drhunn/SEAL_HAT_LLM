# Documentation guide

This directory contains both real runtime guidance and forward-looking architecture notes.
Do not treat every document here as proof that the code already does the thing it describes.

## Read first

These files are the current reality anchors:
- `implementation-status.md` — what is implemented, what is scaffolded, and what is still missing
- `likely-breakpoints.md` — where the repo is most likely to fail as it evolves
- `sql-contracts.md` — the contract between SQL schema/functions and the Go runtime
- `hat-llm.md` — the current shape of the Python dataset/training-support layer

## Core architecture notes

These documents explain the intended shape of the system:
- `layering-strategy.md`
- `sub-agent-strategy.md`
- `developmental-model.md`
- `context-window-strategy.md`
- `multimodal-architecture.md`
- `seal-den-system-architecture.md` — parent/specialist roles, per-model harness/db target shape, and oversight/DEN structure
- `parent-routing-orchestration-strategy.md` — learned parent routing over specialist model units and later shared tool use
- `den-freeze-policy.md`
- `dynamic-architecture-strategy.md`
- `ability-first-growth-strategy.md`
- `model-lineage-strategy.md` — descendant model artifacts, specialist lineage, and bundle-level model-unit packaging

## Contracts and reconciliation

Use these when changing code that touches storage or runtime behavior:
- `sql-contracts.md`
- `schema-runtime-reconciliation.md`
- `implementation-status.md`
- `likely-breakpoints.md`

## Rule of thumb

When docs and code disagree:
1. trust the code and the failing verify path first
2. update the docs in the same change set that fixes the mismatch
3. do not sell architecture as implementation
