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
- modality-aware routing and execution planning
- model-host abstraction and simulated host registry wiring
- basic Go unit tests for retrieval helpers, routing, execution planning, and recovery planner logic

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
- `make verify` exists for schema/runtime smoke checking

---

## partially implemented / scaffolded

### runtime orchestration
- parent/specialist orchestration is scaffolded conceptually but not yet a full live production loop
- sub-agent orchestration remains mostly architectural and policy-level rather than deeply implemented
- tool broker integration is still conceptual
- slot DB/filesystem synchronization is only a snapshot scaffold

### lifecycle and health
- health update hooks exist in Go and SQL, but they still need ongoing reconciliation as the schema and runtime evolve
- degraded/suspended transitions are conceptually defined and partially wired, but not yet enforced by a full state machine

### retrieval and memory operations
- retrieval wrapper exists, but SQL/runtime contracts still need active verification as both evolve
- candidate staging paths exist in Go, and the DB layer is closer to matching them, but this remains an area to watch
- context-overflow summarization to Postgres is now architectural policy, but not yet fully implemented as a live orchestration workflow
- multimodal execution persistence helpers now exist, but the runtime still needs deeper integration beyond smoke-path usage

### training and growth
- dataset generation is useful now
- training support is more real than before, but still scaffold-level and framework-dependent rather than benchmarked production training
- DEN-style ability-first growth is now documented, but governed dynamic growth is still more architectural than operational

### multimodal direction
- multimodal architecture and memory scaffold direction are documented
- the live runtime now has modality-aware routing and execution planning, but it is still not a full real multimodal production stack

---

## known weak spots
- Go runtime and SQL schema/functions must still be kept in sync deliberately
- CI is useful, but it is still a floor rather than proof of full runtime integration
- architecture docs are ahead of production readiness
- multimodal support is no longer just imagined, but still not operationally complete
- DEN-style dynamic growth policy is ahead of runtime machinery that can safely execute it

---

## intentionally deferred
- full multi-specialist orchestration
- mature eval execution engine
- production tool broker
- full benchmark harness
- distributed or large-scale training workflows
- live multimodal model integration beyond the current scaffolding
- full governed dynamic-growth machinery for branches, expert banks, and structured pruning loops

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
10. the first ability-growth path is executed in a governed way rather than only described in docs

---

## summary
This repository is no longer just an idea, but it is not yet a production runtime.
It is a strong, explicit, architecture-first scaffold with enough implementation to support focused reconciliation, end-to-end stabilization, earlier multimodal-aware planning, and a DEN-style ability-first growth direction.
