# IMPLEMENTATION STATUS

## purpose
This document separates the target architecture from the currently implemented state of `SEAL_HAT_LLM`.

---

## implemented now

### documentation and governance
- architecture, runbooks, incident taxonomy, recovery playbooks, eval taxonomy, lineage strategy, context-window strategy, multimodal extension, and harness-aware training docs exist
- specialist slot model is documented and instantiated for the first specialist
- markdown preservation convention for blocked in-session updates is documented
- MM-ELLS naming is now reflected in the core docs

### Go runtime scaffold
- config loading
- Postgres connection pooling
- slot file loading
- runtime startup path
- harness service scaffold
- postmortem creation path
- eval staging path
- retrieval wrapper scaffold
- routing audit write path scaffold
- lifecycle service scaffold
- standalone `cmd/verify` command
- basic Go unit tests for retrieval helpers and recovery planner logic

### SQL scaffold
- base schema for specialists, slots, memory records, embeddings, postmortems, eval cases, self-edit candidates, and routing audit
- hierarchical memory schema for regions and clusters
- starter memory functions and coarse-to-fine search function
- multimodal-memory scaffold tables
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
- Python compile checks and unit tests are wired into CI
- `make verify` now exists for schema/runtime smoke checking

---

## partially implemented / scaffolded

### runtime orchestration
- parent/specialist orchestration is scaffolded conceptually but not yet fully implemented as a live model-host abstraction
- tool broker integration is still conceptual
- slot DB/filesystem synchronization is only a snapshot scaffold

### lifecycle and health
- health update hooks exist in Go and SQL, but they still need ongoing reconciliation as the schema and runtime evolve
- degraded/suspended transitions are conceptually defined and partially wired, but not yet enforced by a full state machine

### retrieval and memory operations
- retrieval wrapper exists and now has a better SQL contract target, but SQL/runtime contracts still need active verification as both evolve
- candidate staging paths exist in Go, and the DB layer is closer to matching them, but this remains an area to watch
- context-overflow summarization to Postgres is now architectural policy, but not yet fully implemented as a live orchestration workflow

### training
- dataset generation is useful now
- training support is more real than before, but still scaffold-level and framework-dependent rather than benchmarked production training

### multimodal direction
- multimodal architecture and memory scaffold direction are now documented
- the live runtime and HAT tooling are still primarily text-first in implementation

---

## known weak spots
- Go runtime and SQL schema/functions must still be kept in sync deliberately
- CI is useful, but it is still a floor rather than proof of full runtime integration
- architecture docs are ahead of production readiness
- multimodal support is planned earlier now, but still not operationally complete

---

## intentionally deferred
- full multi-specialist orchestration
- mature eval execution engine
- production tool broker
- full benchmark harness
- distributed or large-scale training workflows
- live multimodal model integration

---

## next recommended milestone
The next milestone should still be:

**single-specialist loop works end to end**

That means:
1. config loads
2. slots load
3. DB connects
4. retrieval works
5. postmortem writes
6. eval case writes
7. health update works
8. verify command passes
9. context-overflow policy starts being implemented concretely

---

## summary
This repository is no longer just an idea, but it is not yet a production runtime.
It is a strong, explicit, architecture-first scaffold with enough implementation to support focused reconciliation, end-to-end stabilization, and earlier multimodal-aware planning.
