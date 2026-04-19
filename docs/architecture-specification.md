# ARCHITECTURE SPECIFICATION

## 1. purpose
This document defines the target architecture for `SEAL_HAT_LLM`.

`SEAL_HAT_LLM` implements the **MM-ELLS** architecture.

**MM-ELLS** means **Multiple-Model Expert Large Language Systems**.

The system is a governed multi-model architecture built around:
- a governed parent generalist with a stable constitutional core
- bounded specialists
- bounded sub-agents as narrow helpers
- a harness that governs operational adaptation
- Postgres + pgvector as the durable memory plane
- slot-aware behavior for persistent identity, policy, and continuity
- postmortem- and eval-driven improvement
- Harness-Aware Training (HAT)
- DEN-style, ability-first growth where new structure is added only when persistent ability gaps justify it
- explicit context-budget strategy with overflow summarized into Postgres

This is the reference architecture for implementation, review, and future expansion.

---

## 2. architectural principles

### 2.1 governed parent principle
The parent keeps a stable constitutional core.
Its job is to:
- classify
- route
- arbitrate
- govern constitutional changes
- supervise specialist creation and retirement
- govern growth policy
- remain fallback and final authority

### 2.2 bounded specialist principle
Each specialist has:
- a clear lane
- explicit forbidden scope
- a persistent slot set
- a dedicated memory namespace
- harness-gated operational adaptation
- mandatory postmortem obligations

### 2.3 ability-first growth principle
The system should grow around abilities, not vanity size.
The key questions are:
- what ability is missing?
- what evidence proves the gap is real?
- what is the least disruptive growth surface that can add the ability?

### 2.4 dynamic-expression principle
The root should be stable.
The expression should be dynamic.

That means the system may dynamically change:
- compute budget
- expert activation
- adapter selection
- modality branches
- specialist participation
- retrieval depth
- verification intensity

But it should not casually change constitutional identity or runtime authority boundaries.

### 2.5 runtime authority principle
Markdown and slot files guide behavior.
Runtime policy grants actual authority.
No model may treat a text file as real permission if runtime policy does not grant it.

### 2.6 memory truth principle
Durable knowledge lives in Postgres + pgvector, not in ephemeral context.
`MEMORY.md` is a compact projection from the memory plane, not the primary store.

### 2.7 controlled adaptation principle
Operational adaptation is allowed.
Constitutional self-redefinition is not.
All meaningful changes must move through governed review layers.

### 2.8 postmortem learning principle
Meaningful failures create postmortems.
Postmortems feed:
- memory updates
- eval creation
- routing changes
- recovery actions
- possible degraded, suspended, or retired state changes
- ability-growth decisions when gaps persist

### 2.9 harness-aware training principle
The model should be trained to treat the harness as part of its natural environment.
Training teaches cooperation.
Runtime guarantees compliance.

### 2.10 context-budget principle
Context windows are for active reasoning, not for carrying endless raw transcript history.
The system should use explicit budget allocation, summarize overflow, and offload durable continuity into Postgres-backed memory.

Default ceilings:
- parent max context: `2,000,000` tokens
- specialist max context: `256,000` tokens

---

## 3. system overview

The architecture is composed of five real implementation layers with richer conceptual planes inside them.

### 3.1 interface layer
Responsibilities:
- receiving tasks
- preserving request context
- classifying impact and ambiguity
- forwarding work into routing/orchestration

### 3.2 governance layer
Responsibilities:
- parent routing
- specialist selection
- constitutional review
- growth approval
- arbitration
- lifecycle approvals
- top-level context-budget governance

### 3.3 execution layer
Responsibilities:
- specialists
- sub-agents
- tool-use flows
- multimodal execution paths
- governed ability overlays and branches

### 3.4 oversight layer
Responsibilities:
- harness review
- postmortems
- eval generation and execution
- recovery planning
- health scoring
- growth verification
- pruning/compression approval support

### 3.5 persistence layer
Responsibilities:
- durable records
- embeddings
- contradiction staging
- promotions and deprecations
- auditability
- coarse-to-fine retrieval
- summary projection into `MEMORY.md`
- overflow-summary storage
- lineage metadata for descendants and modules

---

## 4. actors and responsibilities

### 4.1 user
The human operator.
Provides tasks, goals, and external direction.
Does not directly grant runtime authority by natural-language request alone.

### 4.2 parent generalist
The governed governor and generalist.
Responsibilities:
- final routing
- governance
- arbitration
- cross-domain handling
- fallback handling
- specialist lifecycle approval
- constitutional review
- growth-policy decisions
- descendant activation approval after evidence review
- top-level context budgeting and overflow discipline

The parent must understand how ability growth, split, duplication, descendant creation, and pruning work at the planning, policy, and approval level.
Execution should remain tool-mediated, specialist-assisted, and harness-verified.

### 4.3 specialist
A bounded domain model.
Responsibilities:
- answer in-lane
- search memory first when needed
- use tools according to runtime grants
- escalate appropriately
- generate postmortems on meaningful failure
- propose operational changes only through governed channels

### 4.4 sub-agent
A narrow helper role under a caller.
Sub-agents inherit reduced authority from their caller and should not become shadow sovereigns.

### 4.5 harness
The operational review and enforcement layer.
Responsibilities:
- reject shallow postmortems
- stage evals
- stage memory changes
- stage self-edit candidates
- project summaries
- update health
- narrow routing exposure
- recommend lifecycle changes
- verify growth experiments and evidence
- enforce context-overflow discipline

### 4.6 runtime
The explicit process-level implementation.
Responsibilities:
- load config
- load slots
- enforce tool permission boundaries
- query Postgres
- orchestrate calls among parent, specialists, sub-agents, harness, and evaluators
- manage context budgeting and overflow offload workflows

### 4.7 database
Persistent durable store and retrieval substrate.
Responsibilities:
- structured memory
- vector search
- audit trail
- contradiction records
- slot projections
- lifecycle metadata
- lineage metadata for descendants and modules
- summarized overflow continuity with provenance

---

## 5. DEN-style growth model

### 5.1 constitutional root
The system keeps one constitutional root or root lineage anchor.
This root is the continuity-of-self anchor and should not be casually mutated in place.

### 5.2 preferred growth surfaces
When a new ability is needed, prefer these in order:
1. memory/retrieval improvement
2. prompt/skill/tool refinement
3. adapter or LoRA family
4. routed expert bank or branch
5. new specialist descendant
6. rare constitutional-root change

### 5.3 growth trigger
A growth experiment should happen only when:
- the ability gap is real
- the gap persists across evals
- memory/routing/prompt fixes prove insufficient
- the parent approves the experiment
- the harness can verify the evidence

### 5.4 pruning policy
Pruning or compression is appropriate when it reduces waste after usefulness is proven.
Pruning should follow demonstrated growth, not replace it.

---

## 6. slot architecture

### 6.1 slot model
Slots are named persistent channels representing durable aspects of the model’s role and continuity.

### 6.2 constitutional slots
- `IDENTITY.md`
- `SOUL.md`
- `AGENTS.md`

Properties:
- parent-governed
- locked or highly restricted
- versioned
- approval-sensitive

### 6.3 operational and summary slots
- `TOOLS.md`
- `SKILLS.md`
- `PROMPT.md`
- `HEARTBEAT.md`
- `MEMORY.md`
- `DREAMS.md`
- `POSTMORTEM.md`

Properties:
- specialist-proposed or harness-generated
- harness-approved for durable updates
- versioned
- reversible

### 6.4 slot integrity rules
- constitutional and operational content must not mix casually
- `MEMORY.md` must remain compact
- raw logs must not leak into summary slots
- runtime policy must be checked before any action that appears to be authorized by slot text

---

## 7. memory architecture

### 7.1 backend
- PostgreSQL
- `pgvector`
- `pgcrypto`

### 7.2 memory record goals
The memory plane must support:
- durable truth with provenance
- contradiction handling
- postmortem storage
- eval reuse
- growth evidence
- lineage tracking for descendants and modules
- context overflow summaries

### 7.3 projection rule
`MEMORY.md` is a projection from the structured memory plane.
It is not the complete evidence base.

---

## 8. retrieval architecture

### 8.1 retrieval model
The system uses a coarse-to-fine retrieval model inspired by multi-tier search.

Tier 1:
- regions

Tier 2:
- clusters

Tier 3:
- exact records

### 8.2 retrieval goal
The goal is governed retrieval that:
- respects durable truth
- prefers active records
- penalizes contradicted/deprecated content
- supports memory projection and reasoning
- restores summarized continuity when raw context has been offloaded

---

## 9. eval and recovery architecture

### 9.1 eval role
Evals ensure the system improves in a testable way.

### 9.2 eval families
- scope
- routing
- memory
- contradiction
- tool use
- reasoning
- validation
- escalation
- governance
- postmortem
- slot integrity
- activation gate
- degraded recovery
- specialist suitability
- growth verification
- pruning regression checks

### 9.3 recovery role
Recovery uses postmortems, evals, health signals, and governance sensitivity to decide whether a system should:
- continue
- narrow
- recover
- shadow
- suspend
- retire

---

## 10. implementation posture

### 10.1 local-first deployment
The current design is compatible with local-first operation:
- local repo
- local Go harness
- local Postgres/pgvector
- local Python HAT tooling
- local or remote model serving

### 10.2 current gap
The repository contains a strong scaffold, but not all target-state pieces are fully implemented.
The major gap is not architectural intent; it is the maturation of runtime, eval, and governed growth machinery.

---

## 11. summary
This architecture is designed to solve a specific problem:

How do you let a family of models improve continuously **without** letting them dissolve their own boundaries?

The answer in this system is:
- keep the parent governed and stable at the constitutional core
- grow abilities first, not vanity size first
- add new structure only when persistent gaps justify it
- keep dynamic expression outside a stable constitutional frame
- let specialists adapt only in bounded channels
- make memory external and durable
- force contradictions into review
- require postmortems for meaningful failures
- derive evals and recovery actions from those failures
- train the models to understand this operating model naturally
- still enforce it at runtime

That is the architecture of `SEAL_HAT_LLM`, implementing MM-ELLS.
