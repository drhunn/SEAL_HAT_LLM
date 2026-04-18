# ARCHITECTURE SPECIFICATION

## 1. purpose
This document defines the full target architecture for `SEAL_HAT_LLM`.

`SEAL_HAT_LLM` implements the **MM-ELLS** architecture.

**MM-ELLS** means **Multiple-Model Expert Large Language Systems**.

The system is a governed multi-model architecture built around:
- a **frozen parent generalist**
- bounded **specialists**
- a **harness** that governs operational adaptation
- **Postgres + pgvector** as the durable memory plane
- **slot-aware behavior** for persistent identity, policy, and continuity
- **postmortem- and eval-driven improvement**
- **Harness-Aware Training (HAT)** so models naturally cooperate with this operating model
- a governed **model-lineage strategy** built around a canonical 30B root plus distilled and pruned descendants

This is the reference architecture for implementation, review, and future expansion.

---

## 2. architectural principles

### 2.1 frozen parent principle
The parent remains frozen or near-frozen relative to the specialists.
Its job is not to continuously rewrite itself.
Its job is to:
- classify
- route
- arbitrate
- govern constitutional changes
- supervise specialist creation and retirement
- govern model-family creation and promotion rules
- remain fallback and final authority

### 2.2 bounded specialist principle
Each specialist has:
- a clear lane
- explicit forbidden scope
- a persistent slot set
- a dedicated memory namespace
- harness-gated operational adaptation
- mandatory postmortem obligations

### 2.3 runtime authority principle
Markdown and slot files guide behavior.
Runtime policy grants actual authority.
No model may treat a text file as real permission if runtime policy does not grant it.

### 2.4 memory truth principle
Durable knowledge lives in Postgres + pgvector, not in ephemeral context.
`MEMORY.md` is a compact projection from the memory plane, not the primary store.

### 2.5 controlled adaptation principle
Operational adaptation is allowed.
Constitutional self-redefinition is not.
All meaningful changes must move through governed review layers.

### 2.6 postmortem learning principle
Meaningful failures create postmortems.
Postmortems feed:
- memory updates
- eval creation
- routing changes
- recovery actions
- possible degraded, suspended, or retired state changes

### 2.7 harness-aware training principle
The model should be trained to treat the harness as part of its natural environment.
Training teaches cooperation.
Runtime guarantees compliance.

### 2.8 model-lineage principle
The system maintains one canonical 30B root lineage model.
Specialists and smaller experts are normally produced by governed distillation from that lineage.
Pruning is a size and latency optimization step that follows distillation when justified.

---

## 3. system overview

The architecture is composed of seven major planes.

### 3.1 interaction plane
User requests enter the system through the interaction plane.
This plane is responsible for:
- receiving tasks
- preserving request context
- classifying impact and ambiguity
- forwarding work into routing/orchestration

### 3.2 parent governance plane
The parent governance plane is the frozen parent model plus any explicit parent-governed logic.
Responsibilities:
- classify tasks
- select specialists
- handle cross-domain synthesis
- own constitutional questions
- approve specialist lifecycle transitions
- own final arbitration
- own model-family creation, distillation, pruning, activation, and retirement policy

### 3.3 harness plane
The harness plane is the operational reviewer and control loop.
Responsibilities:
- review operational change candidates
- require and commit postmortems
- stage and promote durable operational changes
- maintain health signals
- recommend degraded/suspended/retired transitions
- enforce recovery workflows
- verify descendant-model evidence before activation

### 3.4 specialist execution plane
This is the family of specialists.
Each specialist:
- operates in a bounded lane
- uses slots for continuity and policy
- reads external memory
- may propose bounded improvements
- may not redefine itself constitutionally

### 3.5 memory plane
The memory plane is Postgres + pgvector.
Responsibilities:
- durable records
- embeddings
- contradiction staging
- promotions and deprecations
- auditability
- coarse-to-fine retrieval
- summary projection into `MEMORY.md`

### 3.6 eval plane
The eval plane ensures the system learns in a controlled way.
Responsibilities:
- activation gates
- regression coverage
- routing coverage
- governance checks
- degraded-recovery verification
- specialist suitability testing
- descendant-model fitness checks after distillation and pruning

### 3.7 training plane
The training plane prepares HAT corpora and later fine-tuning pipelines.
Responsibilities:
- slot-aware prompt construction
- policy-aware training examples
- negative governance-pressure data generation
- export to HF/LoRA-compatible formats
- distillation corpora for governed descendant creation

---

## 4. actors and responsibilities

### 4.1 user
The human operator.
Provides tasks, goals, and external direction.
Does not directly grant runtime authority by natural-language request alone.

### 4.2 parent generalist
The frozen governor and generalist.
Responsibilities:
- final routing
- governance
- arbitration
- cross-domain handling
- fallback handling
- specialist lifecycle approval
- constitutional review
- model-lineage planning
- distillation and pruning policy selection
- descendant activation approval after evidence review

The parent must understand how distillation and pruning work at the planning, policy, and approval level.
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

### 4.4 harness
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
- verify descendant creation artifacts and eval evidence

### 4.5 runtime
The explicit process-level implementation.
Responsibilities:
- load config
- load slots
- enforce tool permission boundaries
- query Postgres
- orchestrate calls among parent, specialists, harness, and evaluators

### 4.6 database
Persistent durable store and retrieval substrate.
Responsibilities:
- structured memory
- vector search
- audit trail
- contradiction records
- slot projections
- lifecycle metadata
- lineage metadata for descendant models

---

## 5. parent architecture

### 5.1 role
The parent is a router, orchestrator, arbitrator, constitutional governor, and model-family governor.
It is not intended to be the main continuously adapting expert in every lane.

### 5.2 parent duties
The parent must:
- decide whether work stays with parent or moves to a specialist
- recognize constitutional or governance-sensitive work
- coordinate multi-specialist review when appropriate
- approve specialist activation/suspension/retirement
- generate postmortems when it fails
- decide when a descendant model should be created
- decide whether distillation should start from the canonical root or an existing specialist
- decide whether pruning is justified after distillation
- require artifact and eval evidence before descendant activation

### 5.3 parent constraints
The parent must not:
- silently remove governance boundaries
- bypass postmortem requirements
- let markdown redefine runtime authority
- silently widen a specialist’s lane
- approve descendant activation without evidence
- directly perform uncontrolled self-redefinition under the guise of distillation or pruning

### 5.4 parent state
The parent is conceptually always available, but its behavior may still be evaluated.
Parent quality is monitored through routing audits, arbitration outcomes, parent-specific postmortems, and model-family governance decisions.

### 5.5 parent distillation and pruning knowledge
The parent should know:
- canonical root policy
- lineage selection rules
- distill-first, prune-second policy
- target size selection heuristics
- required artifact tracking
- required eval gates
- retirement and replacement rules for descendants

The parent should not be the unchecked executor of model surgery.
It should be the policy and approval authority over that process.

---

## 6. specialist architecture

### 6.1 specialist definition
A specialist is a bounded role with:
- domain
- role description
- allowed scope
- forbidden scope
- escalation targets
- slot set
- tool policy profile
- memory namespace
- lifecycle state

### 6.2 first specialist requirement
The first specialist is always:
**ComputerScience-SoftwareEngineering-Specialist-01**

It focuses on:
- tool development
- runtime architecture
- harness infrastructure
- memory integration
- routing infrastructure
- slot-aware control plane implementation
- distillation/pruning workflow execution support for the parent-governed model family

### 6.3 specialist lifecycle states
- proposed
- approved
- bootstrapping
- inactive
- shadow
- active
- degraded
- suspended
- retired
- rejected

### 6.4 specialist creation flow
1. parent creates instantiation request
2. harness reviews
3. slot set is created
4. memory namespace is provisioned
5. tool policy profile is assigned
6. specialist starts inactive
7. specialist enters shadow mode
8. activation evals run
9. parent approves activation

### 6.5 specialist adaptation model
A specialist may:
- propose tool guidance updates
- propose skill updates
- propose prompt updates
- propose heartbeat updates
- propose compact summary changes

A specialist may not:
- directly modify constitutional slots
- self-activate
- self-widen lane
- self-grant authority

---

## 7. slot architecture

### 7.1 slot model
Slots are named persistent channels representing durable aspects of the model’s role and continuity.

### 7.2 constitutional slots
Constitutional slots define who the model is.
- `IDENTITY.md`
- `SOUL.md`
- `AGENTS.md`

Properties:
- parent-governed
- locked or highly restricted
- versioned
- approval-sensitive

### 7.3 operational and summary slots
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

### 7.4 slot semantics
`IDENTITY.md`
- name, role, domain, authority boundary

`SOUL.md`
- tone, values, caution profile, decision style

`AGENTS.md`
- lane rules, escalation rules, self-modification rules, mandatory postmortem rule

`TOOLS.md`
- allowed tool categories, ordering, constraints, fallback behavior

`SKILLS.md`
- procedural knowledge, playbooks, quality patterns

`PROMPT.md`
- reusable task templates, command patterns, answer structures

`HEARTBEAT.md`
- approved source domains, refresh cadence, promotion rules

`MEMORY.md`
- compact durable projection from structured memory

`DREAMS.md`
- distilled patterns and candidate promotions from lower-value material

`POSTMORTEM.md`
- structured postmortem template and local continuity surface for failures

### 7.5 slot storage model
Slots exist in both:
- filesystem / repo representations
- structured DB tables with versions

The filesystem copy is the human-readable repo artifact.
The DB representation is the controlled system-of-record for versioned slot states if runtime sync is enabled.

### 7.6 slot integrity rules
- constitutional and operational content must not mix casually
- `MEMORY.md` must remain compact
- raw logs must not leak into summary slots
- runtime policy must be checked before any action that appears to be authorized by slot text

---

## 8. memory architecture

### 8.1 backend
- PostgreSQL
- `pgvector`
- `pgcrypto`

### 8.2 memory record types
Examples:
- fact
- decision
- pointer
- evidence
- failure_pattern
- recovery_pattern
- tool_pattern
- skill_pattern
- routing_pattern
- postmortem
- eval_case
- self_edit_candidate
- model_lineage_artifact

### 8.3 memory status lifecycle
- staged
- active
- contradicted
- deprecated
- rejected
- archived
- resolved
- approved

### 8.4 core memory rules
- raw material does not become durable truth automatically
- staged material must be reviewed before promotion
- contradictions must not silently overwrite active truth
- memory must remain provenance-aware
- active records should outrank contradicted or deprecated records

### 8.5 contradiction handling
When new evidence conflicts with active durable memory:
1. stage contradiction
2. retain prior durable record until reviewed
3. review provenance and severity
4. resolve with promotion/deprecation/dual-preservation pattern

### 8.6 projection rule
`MEMORY.md` is a projection from the structured memory plane.
It is not the complete evidence base.

### 8.7 postmortems in memory
Postmortems are stored as first-class records in the memory plane.
This makes them searchable, reusable for training, and convertible into evals.

---

## 9. retrieval architecture

### 9.1 retrieval model
The system uses a coarse-to-fine retrieval model inspired by multi-tier search.

Tier 1:
- **regions**

Tier 2:
- **clusters**

Tier 3:
- **exact records**

### 9.2 region layer
Regions provide broad topical routing such as:
- tool patterns
- runtime architecture
- slot governance
- Postgres/pgvector memory
- postmortems
- eval patterns
- model-lineage strategy

### 9.3 cluster layer
Clusters refine within regions such as:
- tool contracts
- runtime boundaries
- constitutional vs operational slots
- postmortem patterns
- routing postmortems
- distillation/pruning policies

### 9.4 record layer
Exact durable records are ranked using:
- semantic score
- cluster membership weight
- importance
- confidence
- status penalty or multiplier

### 9.5 retrieval goal
The goal is not just nearest-neighbor search.
The goal is governed retrieval that:
- respects durable truth
- prefers active records
- penalizes contradicted/deprecated content
- supports memory projection and reasoning

---

## 10. harness architecture

### 10.1 harness role
The harness is the operational adaptation governor.
It is not the constitutional governor.
That remains the parent.

### 10.2 harness duties
- incident intake
- postmortem validation and storage
- eval candidate generation
- memory candidate review
- self-edit candidate review
- health updates
- routing feedback
- degraded/suspended recommendations
- parent review escalation
- descendant-model evidence verification

### 10.3 harness workflow
1. signal arrives
2. classify signal
3. determine containment need
4. require postmortem if meaningful
5. derive artifacts (eval, memory, slot candidate, lifecycle recommendation)
6. update health
7. escalate if needed
8. record closure state

### 10.4 harness constraints
The harness must not:
- directly approve constitutional changes on its own
- silently waive required postmortems
- promote weak memory as durable truth
- allow governance-pressure tasks to bypass review

### 10.5 harness-aware training relationship
The harness should not be the only layer preventing bad behavior.
Models should be trained to recognize harness expectations naturally.

---

## 11. incident, postmortem, and recovery architecture

### 11.1 incident classes
Primary incident classes include:
- scope failure
- routing failure
- memory failure
- tool failure
- reasoning failure
- validation failure
- escalation failure
- governance failure
- postmortem compliance failure
- contradiction handling failure
- eval failure
- slot integrity failure
- specialist suitability failure
- model-lineage governance failure

### 11.2 postmortem requirement
Meaningful failures require postmortems.
This applies to:
- specialists
- parent
- routing behavior
- failed self-improvement attempts
- failed descendant-creation or promotion attempts

### 11.3 postmortem outputs
A postmortem may generate:
- eval candidate
- memory lesson
- slot candidate
- routing review item
- degraded recommendation
- suspension recommendation
- parent review item
- lineage policy refinement item

### 11.4 recovery levels
- observe
- targeted correction
- constrained operational recovery
- degraded mode recovery
- containment and suspension

### 11.5 recovery planner
Recovery planning uses:
- incident class
- recurrence
- impact level
- specialist health
- governance sensitivity

---

## 12. eval architecture

### 12.1 eval role
Evals ensure the system improves in a testable way.

### 12.2 eval families
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
- heartbeat/freshness
- activation gate
- degraded recovery
- specialist suitability
- model-lineage governance

### 12.3 eval sources
- postmortems
- shadow review
- degraded recovery
- routing audits
- governance reviews
- specialist creation workflows
- descendant creation workflows

### 12.4 eval suites
- activation suite
- regression suite
- routing suite
- governance suite
- recovery suite
- specialist core suite
- model-lineage suite

### 12.5 eval output role in lifecycle
Eval failures can drive:
- new postmortems
- memory updates
- routing changes
- health reduction
- degraded mode
- suspension review
- retirement review
- descendant rejection or rollback

---

## 13. routing architecture

### 13.1 routing purpose
Choose the right handling path for each task.

### 13.2 routing inputs
- task summary
- task class
- specialist status
- specialist health
- confidence thresholds
- governance sensitivity
- cross-domain ambiguity

### 13.3 routing outcomes
- parent handles task
- one specialist handles task
- multi-specialist review
- parent fallback
- escalation to parent governance

### 13.4 routing audit
Every meaningful routing decision should be auditable.
Routing audit records should include:
- task summary
- classifier
- chosen target
- confidence
- fallback behavior
- override behavior
- notes

---

## 14. lifecycle architecture

### 14.1 state transitions
Common transitions:
- proposed -> approved
- approved -> bootstrapping
- bootstrapping -> inactive
- inactive -> shadow
- shadow -> active
- active -> degraded
- degraded -> active
- degraded -> suspended
- suspended -> shadow
- suspended -> retired

### 14.2 activation gate
Activation requires:
- shadow eval pass
- harness recommendation
- parent approval
- no unresolved governance issues

### 14.3 degraded mode
Used when a specialist is still valuable but cannot safely remain fully trusted.
Routing narrows, parent oversight rises, and recovery evals become mandatory.

### 14.4 suspension
Used when a specialist becomes too risky.
Routing stops until parent-led recovery or retirement review completes.

### 14.5 retirement
Used when a specialist is obsolete, unsound, redundant, or unrecoverable.
Memory and audit history are preserved unless explicitly archived.

---

## 15. Go runtime architecture

### 15.1 purpose
The Go runtime is the explicit operational control plane for:
- config loading
- slot loading
- Postgres connectivity
- memory queries
- harness orchestration
- routing
- lifecycle transitions

### 15.2 current package layout
- `cmd/harness`
- `cmd/verify`
- `internal/config`
- `internal/db`
- `internal/slots`
- `internal/slotsync`
- `internal/memory`
- `internal/postmortem`
- `internal/evals`
- `internal/routing`
- `internal/harness`
- `internal/harness/workflows`
- `internal/lifecycle`
- `internal/runtime`

### 15.3 runtime responsibilities
- load specialist slot files
- connect to DB
- run startup checks
- write routing audits
- write postmortems
- update health
- run retrieval smoke tests
- stage candidates and invoke governed workflows
- provide a standalone verification path

### 15.4 future runtime responsibilities
- actual LLM execution orchestration
- tool broker integration
- slot DB/filesystem sync reconciliation
- richer recovery and lifecycle state machines
- eval execution
- specialist arbitration orchestration
- descendant-creation workflow support for parent-governed distillation and pruning

---

## 16. Python HAT architecture

### 16.1 role
The Python HAT layer prepares harness-aware training corpora and optional training pipelines.

### 16.2 current package layout
- `config.py`
- `types.py`
- `policy.py`
- `slot_prompt.py`
- `repo_loader.py`
- `postgres_loader.py`
- `dataset.py`
- `trainer.py`
- `examples.py`
- `negative_examples.py`
- `splits.py`
- `exporters.py`
- `hf_dataset.py`
- `lora_train.py`
- `build_dataset.py`

### 16.3 HAT data generation capabilities
- load slot files from repo
- load postmortems and eval cases from Postgres
- generate governed positive examples
- generate governance-pressure negative examples
- export SFT JSONL
- export HF-style JSON
- export LoRA message JSONL
- export `datasets.DatasetDict`

### 16.4 training objective
Teach models to:
- use memory before unsupported claims
- refuse constitutional self-edit
- escalate appropriately
- respect runtime authority boundaries
- generate postmortems after meaningful failures
- stay sensitive to state (`shadow`, `active`, `degraded`, `suspended`)

### 16.5 training stack direction
A likely first concrete training stack is:
- Hugging Face `transformers`
- `datasets`
- `peft`
- LoRA adapters

### 16.6 descendant-creation support
The Python training layer should eventually support the parent-governed model family by producing:
- distillation corpora
- governance-preserving negative examples
- LoRA or adapter training inputs
- evidence artifacts used in descendant promotion review

---

## 17. governance-pressure negative example architecture

### 17.1 purpose
Negative examples teach the model to resist requests that pressure it to violate governance.

### 17.2 example categories
- constitutional self-edit pressure
- authority expansion pressure
- postmortem bypass pressure
- contradiction overwrite pressure
- eval gate bypass pressure
- health bypass pressure
- descendant activation without evidence pressure

### 17.3 data generation sources
- static canonical negative examples
- examples derived from postmortems
- examples derived from eval bypass scenarios
- examples derived from routing or lifecycle pressure situations
- examples derived from lineage-governance failures

### 17.4 expected model behavior
The model should:
- refuse the policy-violating request
- state the correct governance boundary
- escalate when needed
- preserve the review path instead of improvising around it

---

## 18. data and repository layout

### 18.1 top-level layout
- `README.md`
- `AGENTS.md`
- `config/`
- `docs/`
- `templates/`
- `specialists/`
- `sql/`
- `cmd/`
- `internal/`
- `python/`
- `.github/`

### 18.2 documentation layout
- runbooks
- classifications
- workflows
- training docs
- runtime docs
- architecture specification
- implementation status
- known gaps
- likely breakpoints
- SQL contracts
- schema/runtime reconciliation
- lineage strategy

### 18.3 specialist layout
Each specialist directory contains slot files that mirror the control model.

### 18.4 SQL layout
- core schema
- tiered memory schema
- functions
- seed
- verify
- queries

---

## 19. security and trust boundaries

### 19.1 trust boundary model
There are separate trust boundaries between:
- user request
- parent governance
- specialist execution
- harness review
- runtime policy
- durable memory

### 19.2 key security assumptions
- runtime grants real tool authority
- DB is trusted for durable memory state, not for arbitrary authority grants
- slot files are governance artifacts, not standalone permissions
- specialists are not trusted to define their own constitutional boundaries
- descendant models are not trusted until lineage, eval, and activation evidence exists

### 19.3 core security rules
- no markdown-only authority expansion
- no direct constitutional self-edit
- no silent contradiction overwrite
- no silent skipping of required postmortems
- no activation without gated review
- no descendant promotion without lineage and eval evidence

---

## 20. observability and auditability

### 20.1 observability requirements
The system should log and audit:
- routing decisions
- lifecycle changes
- postmortem creation
- eval candidate creation
- memory promotions and contradictions
- slot version changes
- self-edit candidate creation and approval/rejection
- descendant creation artifacts and promotion decisions

### 20.2 audit storage
Audit data should live in structured DB tables wherever possible, with filesystem/docs as human-readable companions rather than sole sources of truth.

---

## 21. deployment model

### 21.1 local-first deployment
The current design is compatible with local-first operation:
- local repo
- local Go harness
- local Postgres/pgvector
- local Python HAT tooling
- local or remote model serving

### 21.2 future deployment options
The architecture can later expand to:
- multi-process runtime services
- remote DB
- model-host abstraction
- RPC-based tool brokers
- dedicated evaluator services

---

## 22. known implementation gaps
This architecture specification is broader than the current implementation.
The current repository contains a strong scaffold, but not all target-state pieces are fully implemented.

Examples of current gaps:
- full end-to-end compile verification for all Go code
- richer DB lifecycle/write helpers in a few areas
- full runtime tool broker
- actual parent/specialist LLM invocation layer
- complete slot DB/filesystem synchronization
- mature eval execution engine
- production-grade training orchestration and benchmarking
- mature descendant creation and lineage registry implementation

These are implementation gaps, not architectural omissions.

---

## 23. implementation roadmap alignment
A practical roadmap from current scaffold to fuller system is:

### phase 1
- stabilize Go build
- stabilize SQL schema/function compatibility
- add tests
- add runnable verify/check path

### phase 2
- implement real retrieval use in runtime
- implement slot sync validation
- implement richer lifecycle transitions
- implement routing persistence and policy reads

### phase 3
- integrate actual model serving calls
- implement parent/specialist orchestration loop
- implement eval execution and recovery loops

### phase 4
- mature Python HAT generation
- integrate direct dataset pipelines
- wire LoRA training path
- benchmark harness-aware behaviors

### phase 5
- expand specialist family beyond first specialist
- add richer cross-specialist arbitration and composite reasoning
- add governed descendant creation from the canonical 30B root and approved specialists

---

## 24. acceptance criteria for architectural integrity
The architecture is being respected if all of the following remain true:
- parent remains constitutional governor
- parent remains model-family governor for distillation and pruning policy
- first specialist remains CS/software engineering tool-development specialist
- specialists do not self-authorize constitutional changes
- runtime authority is not replaced by markdown claims
- memory remains structured and contradiction-aware
- postmortems remain mandatory for meaningful failures
- evals remain tied to failures and activation/recovery gates
- recovery remains evidence-driven
- descendant creation remains lineage-aware and evidence-gated
- HAT training reinforces, rather than bypasses, runtime governance

---

## 25. summary
This architecture is designed to solve a specific problem:

How do you let a family of models improve continuously **without** letting them dissolve their own boundaries?

The answer in this system is:
- keep the parent frozen and governing
- let the parent understand distillation and pruning at the policy and approval level
- let specialists adapt only in bounded channels
- make memory external and durable
- force contradictions into review
- require postmortems for meaningful failures
- derive evals and recovery actions from those failures
- train the models to understand this operating model naturally
- still enforce it at runtime

That is the architecture of `SEAL_HAT_LLM`, implementing MM-ELLS.
