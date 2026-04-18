# IMPLEMENTATION STATUS

## purpose
This document separates the target architecture from the currently implemented state of `SEAL_HAT_LLM`.

---

## implemented now

### documentation and governance
- architecture, runbooks, incident taxonomy, recovery playbooks, eval taxonomy, and harness-aware training docs exist
- specialist slot model is documented and instantiated for the first specialist
- markdown preservation convention for blocked in-session updates is documented

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

### SQL scaffold
- base schema for specialists, slots, memory records, embeddings, postmortems, and routing audit
- hierarchical memory schema for regions and clusters
- starter memory functions and coarse-to-fine search function
- seed and verify scripts

### Python HAT tooling
- slot-aware prompt construction
- governed dataset building
- repo slot loading
- optional Postgres corpus ingestion
- governance-pressure negative example generation
- HF/LoRA-style export helpers
- PEFT LoRA training scaffold

---

## partially implemented / scaffolded

### runtime orchestration
- parent/specialist orchestration is scaffolded conceptually but not yet fully implemented as a live model-host abstraction
- tool broker integration is still conceptual
- slot DB/filesystem synchronization is only a snapshot scaffold

### lifecycle and health
- health update hooks exist in Go but depend on SQL contracts that must remain in sync
- degraded/suspended transitions are conceptually defined and partially wired, but not yet enforced by a full state machine

### retrieval and memory operations
- retrieval wrapper exists, but SQL/runtime contracts need active verification as both evolve
- candidate staging paths exist in Go, but the DB layer must stay aligned with those expectations

### training
- dataset generation is useful now
- actual fine-tuning/training is still scaffold-level and framework-dependent

---

## known weak spots
- Go runtime and SQL schema/functions must be kept in sync manually right now
- compile/test status is not yet treated as a strong gate
- CI is newly added and should be treated as the beginning, not the end, of validation
- architecture docs are ahead of production readiness

---

## intentionally deferred
- full multi-specialist orchestration
- mature eval execution engine
- production tool broker
- full benchmark harness
- distributed or large-scale training workflows

---

## next recommended milestone
The next milestone should be:

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

---

## summary
This repository is no longer just an idea, but it is not yet a production runtime.
It is a strong, explicit, architecture-first scaffold with enough implementation to support focused reconciliation and end-to-end stabilization.
