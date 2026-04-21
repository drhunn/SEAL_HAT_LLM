# IMPLEMENTATION STATUS

## purpose
This document separates the target architecture from the currently implemented state of `SEAL_HAT_LLM`.

---

## implemented now

### documentation and governance
- architecture, runbooks, incident taxonomy, recovery playbooks, eval taxonomy, lineage strategy, developmental model, ability-first growth strategy, dynamic architecture strategy, context-window strategy, multimodal extension, and harness-aware training docs exist
- specialist slot model is documented and instantiated for the first specialist
- markdown preservation convention for blocked in-session updates is documented
- MM-ELLS naming is reflected in the core docs
- the documentation now frames growth as DEN-style, ability-first, and governed rather than size-first
- SEAL and DEN system architecture is now documented explicitly

### Go runtime scaffold
- config loading
- Postgres connection pooling
- slot file loading
- canonical specialist slot bundle compilation and TOML emission
- runtime startup path
- bounded runtime task processing path now exists through `runtime.Service.ProcessTask`
- file-backed runtime task inbox now exists and can claim JSON task files, process them, and move them to `processed/` or `failed/` without deleting them
- harness incident handling is now wired into the bounded task path for execution errors, parent-review requirements, and high-severity runtime signals
- harness service scaffold
- postmortem creation path
- eval staging path
- retrieval wrapper scaffold
- routing audit write path scaffold
- lifecycle service scaffold
- standalone `cmd/verify` command
- modality-aware routing and execution planning
- model-host abstraction and simulated host registry wiring
- telemetry collector, in-memory telemetry store, and retrieval/routing/execution signal helpers now exist
- SEAL proposal generation scaffold now exists
- DEN growth planning scaffold now exists
- oversight and lineage scaffolds now exist
- basic Go unit tests for retrieval helpers, routing, execution planning, recovery planner logic, slot compiler behavior, runtime task helpers, SEAL proposal generation, and DEN plan generation
- DB-backed runtime integration tests now cover bounded task success persistence and failure-triggered postmortem/eval writes when `TEST_DATABASE_DSN` is provided
- task inbox unit tests now cover file claiming, processed archiving, failed archiving, and error-note emission

### SQL scaffold
- base schema for specialists, slots, memory records, embeddings, postmortems, eval cases, self-edit candidates, and routing audit
- hierarchical memory schema for regions and clusters
- starter memory functions and coarse-to-fine search function
- multimodal-memory scaffold tables
- SEAL/DEN scaffold tables for adaptation signals, proposals, growth plans, lineage, promotion decisions, and slot bundle versions now exist
- seed and verify scripts

### Python HAT tooling
- slot-aware prompt construction
- governed dataset building
- repo slot loading
- optional Postgres corpus ingestion
- governance-pressure negative example generation
- HF/LoRA-style export helpers
- direct `DatasetDict` export support
- PEFT LoRA training scaffold
- Python unit tests for policy and dataset builder basics

### CI and checks
- basic GitHub Actions workflow exists
- Go build/test is wired into CI
- CI now injects `TEST_DATABASE_DSN` for DB-backed runtime integration coverage
- Python compile checks and unit tests are wired into CI
- `make verify` exists for schema/runtime smoke checking

---

## partially implemented / scaffolded

### runtime orchestration
- parent/specialist orchestration now has a bounded single-task execution path and a file-backed task inbox, but not yet a full live production loop with remote intake or queue semantics
- sub-agent orchestration remains mostly architectural and policy-level rather than deeply implemented
- tool broker integration is still conceptual
- slot DB/filesystem synchronization is now stronger because the canonical bundle compiler exists, but DB-backed bundle reconciliation is still not a full live sync workflow

### lifecycle and health
- health update hooks exist in Go and SQL, but they still need ongoing reconciliation as the schema and runtime evolve
- degraded/suspended transitions are conceptually defined and partially wired, but not yet enforced by a full state machine

### retrieval and memory operations
- retrieval wrapper exists, but SQL/runtime contracts still need active verification as both evolve
- candidate staging paths exist in Go, and the DB layer is closer to matching them, but this remains an area to watch
- context-overflow summarization to Postgres is now architectural policy, but not yet fully implemented as a live orchestration workflow
- multimodal execution persistence helpers now exist, and the bounded runtime task path now persists execution artifacts, but the runtime still needs deeper integration beyond the current single-task path
- `cmd/verify` now attempts durable bundle, signal, proposal, and growth-plan persistence, but still falls back safely when the new migration has not been applied yet

### training and growth
- dataset generation is useful now
- training support is more real than before, but still scaffold-level and framework-dependent rather than benchmarked production training
- DEN-style ability-first growth is no longer only documented: the repo now has telemetry, proposal, and planning scaffolds, and the runtime can stage a bounded task-follow-up growth record, but governed dynamic growth is still not a full production runtime
- the existing `growth` package should now be understood as the governed experiment-execution layer under SEAL and DEN decisions rather than a competing governor

### multimodal direction
- multimodal architecture and memory scaffold direction are documented
- the live runtime now has modality-aware routing and execution planning, but it is still not a full real multimodal production stack

---

## known weak spots
- Go runtime and SQL schema/functions must still be kept in sync deliberately
- the new SEAL/DEN migration must be applied before durable persistence paths fully work
- CI is useful, but it is still a floor rather than proof of full runtime integration
- architecture docs are ahead of production readiness
- multimodal support is no longer just imagined, but still not operationally complete
- DEN-style dynamic growth policy is now closer to code, but oversight/promotion/rollback are still not fully enforced end to end
- the file-backed inbox is intentionally conservative and local; it is not yet a remote intake API or a distributed queue

---

## intentionally deferred
- full multi-specialist orchestration
- mature eval execution engine
- production tool broker
- full benchmark harness
- distributed or large-scale training workflows
- live multimodal model integration beyond the current scaffolding
- full governed dynamic-growth machinery for branches, expert banks, and structured pruning loops
- remote task intake APIs and queue backends

---

## next recommended milestone
The next milestone should still be:

**single-specialist loop works end to end and persists durably**

That means:
1. config loads
2. slots load
3. canonical bundle compiles
4. DB connects
5. retrieval works
6. postmortem writes
7. eval case writes
8. health update works
9. verify command passes
10. bundle persistence works against the migrated DB
11. telemetry writes, SEAL proposal writes, and DEN plan writes work durably
12. context-overflow policy starts being implemented concretely
13. the first ability-growth path is executed in a governed way rather than only described in docs
14. the bounded runtime task path and file inbox graduate into a more explicit intake contract

---

## summary
This repository is no longer just an idea, but it is not yet a production runtime.
It is a strong, explicit, architecture-first scaffold with enough implementation to support focused reconciliation, end-to-end stabilization, canonical slot bundling, a bounded runtime task path, a conservative file-backed task inbox, DB-backed runtime integration coverage, early telemetry-driven SEAL/DEN experimentation, multimodal-aware planning, and a DEN-style ability-first growth direction.
