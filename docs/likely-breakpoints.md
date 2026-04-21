# Likely breakpoints

This file records where SEAL_HAT_LLM is most likely to break as the repo evolves.

## 1. Go runtime ↔ SQL function drift

Failure mode:
- Go expects function names, arguments, or return columns that no longer match the database.

What breaks:
- retrieval
- health updates
- summary projection
- verify

Guardrail:
- update `docs/sql-contracts.md` whenever a runtime-facing SQL function changes
- run strict verify after SQL changes

## 2. Go runtime ↔ SQL schema drift

Failure mode:
- Go writes columns or tables that do not exist in the target database.
- Migrations exist in the repo, but were never applied where the runtime is pointed.

What breaks:
- bundle persistence
- telemetry writes
- proposal / growth-plan writes
- lifecycle writes

Guardrail:
- treat migrations as admission-critical
- avoid silent schema edits
- keep bootstrap order explicit

## 3. Namespace drift

Failure mode:
- new SQL files stop using `agent_core` consistently.

What breaks:
- the runtime appears to work locally but writes to the wrong schema or reads from the wrong place.

Guardrail:
- keep `agent_core` as the canonical runtime schema
- qualify names or set `search_path` consistently

## 4. Identifier drift

Failure mode:
- one side assumes text IDs while the other side silently switches to DB-generated UUIDs or a different text format.

What breaks:
- proposal, growth-plan, lineage, and bundle persistence

Guardrail:
- keep identifier strategy documented and stable
- change it only with coordinated Go + SQL updates

## 5. Duplicated policy drift

Failure mode:
- routing and execution packages encode overlapping policy in separate places.

What breaks:
- the system routes to one executor and plans for another
- verify becomes noisy for the wrong reason

Guardrail:
- keep one source of truth for task-to-executor policy

## 6. Docs ahead of code

Failure mode:
- architecture notes are read as if they describe a completed runtime.

What breaks:
- planning quality
- review quality
- roadmap discipline

Guardrail:
- use `implementation-status.md` as the reality anchor
- do not advertise architecture as completed behavior

## 7. Preserved `.md` code ambiguity

Failure mode:
- archived or preserved `.md` code is mistaken for a live entrypoint.

What breaks:
- maintenance
- debugging
- onboarding

Guardrail:
- label preserved code clearly
- keep one canonical executable path

## 8. Verify false confidence

Failure mode:
- soft verify is treated like a real gate, or strict verify becomes too broad to diagnose cleanly.

What breaks:
- admission quality
- debugging speed
- trust in CI

Guardrail:
- use strict verify in CI
- keep verify narrow enough that failures are interpretable

## 9. First-specialist scope creep

Failure mode:
- the first specialist becomes a dumping ground for every technical task in the repo.

What breaks:
- lane clarity
- future specialist design
- governance boundaries

Guardrail:
- keep lane boundaries explicit
- split responsibilities only when the codebase earns it

## Summary

If this repo breaks, the most likely causes are still boring ones:
- schema/runtime mismatch
- migration drift
- policy duplication
- docs outrunning code
- soft checks being mistaken for real gates

That is where review energy should go first.
