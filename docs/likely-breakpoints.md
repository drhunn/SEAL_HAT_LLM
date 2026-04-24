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
- route episode writes
- artifact lifecycle writes

Guardrail:
- treat migrations as admission-critical
- avoid silent schema edits
- keep bootstrap order explicit
- keep `docs/sql-contracts.md` aligned with runtime write paths

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
- proposal, growth-plan, lineage, route episode, artifact event, and bundle persistence

Guardrail:
- keep identifier strategy documented and stable
- change it only with coordinated Go + SQL updates

## 5. Duplicated policy drift

Failure mode:
- routing and execution packages encode overlapping policy in separate places.

What breaks:
- the system routes to one executor and plans for another
- verify becomes noisy for the wrong reason
- remote dispatch fires for the wrong target or fails to fire for a mapped target

Guardrail:
- keep one source of truth for task-to-executor policy
- keep routing/execution target contracts in `unitref`
- avoid making routing/execution import the full `unit` package again

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
- keep `docs/acceptance-criteria.md` aligned with actual admission gates

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

## 10. Task RPC socket flakiness

Failure mode:
- Unix socket tests pass locally but flake under CI load.
- Raw socket reads/writes block without deadlines.
- Server goroutines are not waited on during cleanup.

What breaks:
- task RPC tests
- remote dispatch integration tests
- CI trust

Guardrail:
- use deadlines for raw socket tests
- wait for server shutdown in test cleanup
- run `go test ./internal/taskrpc -count=20` before treating the transport as stable

## 11. Static remote dispatch mapping drift

Failure mode:
- `task_dispatch.remote_unit_sockets` points to a stale socket path or a unit ID that routing never produces.

What breaks:
- parent-side remote dispatch
- route episode accounting for remote tasks

Guardrail:
- keep mappings explicit
- fail loudly on missing mappings for intended remote dispatch
- do not add silent fallback unless it is a named, tested policy

## 12. Skipped DB-backed tests

Failure mode:
- DB-backed proof tests are skipped because `TEST_DATABASE_DSN` is unset, and the branch is treated as validated anyway.

What breaks:
- route episode proof
- end-to-end remote dispatch proof
- artifact lifecycle proof

Guardrail:
- require a DB-backed validation pass before merge
- make skipped DB tests visible in review notes

## 13. Artifact lifecycle theater

Failure mode:
- state transitions exist, but model artifacts are not actually packaged, signed, or rolled out.

What breaks:
- promotion semantics
- rollback expectations
- operator trust

Guardrail:
- treat lifecycle state enforcement as state accounting only
- do not claim bundle production or deployment until artifact packaging exists

## Summary

If this repo breaks, the most likely causes are still boring ones:
- schema/runtime mismatch
- migration drift
- policy duplication
- socket-test flakiness
- DB-backed tests being skipped
- docs outrunning code
- soft checks being mistaken for real gates

That is where review energy should go first.
