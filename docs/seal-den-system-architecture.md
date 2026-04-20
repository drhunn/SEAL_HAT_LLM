# SEAL AND DEN SYSTEM ARCHITECTURE

## purpose
Define how `SEAL_HAT_LLM` should implement the relationship between **SEAL** as the governed adaptation layer and **DEN** as the governed structural expansion layer.

This document is intentionally architectural.
It defines where each responsibility belongs in the runtime, how proposals should move through the harness, and what packages and data structures should exist in the repository.

---

## core framing
In this repository:

- **SEAL** decides **when** the system should adapt and **how** adaptation should be governed
- **DEN** decides **where** capacity should change and **when** structural expansion is justified because existing capacity is insufficient
- the **harness** acts as the approving adult and experiment gatekeeper
- the **parent** remains the constitutional governor
- **Postgres + pgvector** remain the durable memory, audit, and lineage plane

This means the system should not jump directly from failure to growth.
It should move through a governed sequence:

1. observe failure or repeated weakness
2. stage evidence
3. let SEAL decide whether adaptation is justified
4. let the harness and parent approve the adaptation surface
5. let DEN choose the smallest structural change that can address the persistent gap
6. run a bounded experiment
7. verify, promote, or roll back

---

## architectural split

### 1. execution plane
The execution plane handles live work:
- task intake
- slot bundle loading
- memory retrieval
- routing
- execution
- result generation
- telemetry emission

This plane should never directly self-modify durable structure.
It should only emit signals and proposals.

### 2. SEAL plane
The SEAL plane is the **adaptation governor**.
It consumes signals such as:
- postmortems
- eval failures
- repeated routing misses
- repeated retrieval misses
- tool-use failures
- contradiction staging in memory
- context overflow patterns
- specialist underperformance clusters

SEAL answers:
- is the gap real?
- is the gap persistent?
- is it local or systemic?
- can it be solved without structural growth?
- what is the safest adaptation surface?
- what approval level is required?

SEAL emits **adaptation proposals**, not raw edits.

### 3. DEN plane
The DEN plane is the **structural expansion mechanism**.
It runs only after SEAL has justified adaptation and governance has approved the proposal.

DEN answers:
- where is capacity insufficient?
- can the problem be solved with an operational patch?
- should an existing specialist absorb the work?
- should a specialist be split?
- should a new specialist be created?
- should an adapter, branch, or expert bank be added?
- what should remain frozen?
- how can the change be rolled back safely?

DEN emits **growth plans** and bounded experiments.

### 4. oversight plane
The oversight plane owns:
- approval checks
- constitutional vs operational boundaries
- promotion and rollback decisions
- policy checks before durable changes
- audit and lineage updates

### 5. lineage plane
The lineage plane tracks:
- parent and specialist ancestry
- derived branches
- split specialists
- adapters and expert banks
- retirement, rollback, and pruning history

---

## control loop
The intended control loop is:

```text
request
  ->
runtime execution plane
  ->
telemetry + postmortem + eval + routing audit
  ->
SEAL review plane
  ->
adaptation proposal
  ->
approval gate
  ->
DEN growth plane
  ->
governed experiment
  ->
verify / promote / rollback
```

The most important design rule is:

**SEAL must prefer the smallest reversible fix.**

Preferred adaptation order:
1. slot patch
2. prompt / skill / tool patch
3. retrieval / memory patch
4. adapter tuning
5. new specialist
6. specialist split
7. larger structural expansion

The system should not use DEN as an excuse to add capacity casually.

---

## repository package layout
The existing repository already contains runtime, routing, memory, growth, slots, and verification scaffolding.
The architecture should grow by adding a few narrowly scoped packages.

Recommended layout:

```text
cmd/
  harness/
  verify/
  sealctl/
  denctl/

internal/
  config/
  db/
  execution/
  growth/
  memory/
  modality/
  modelhost/
  routing/
  runtime/
  slots/
  slotsync/

  telemetry/
  seal/
  den/
  oversight/
  lineage/
  evalrun/
```

### package responsibilities

#### `internal/telemetry`
Responsible for normalized runtime signals.

Suggested files:
- `types.go`
- `events.go`
- `collector.go`

#### `internal/seal`
Responsible for adaptation governance.

Suggested files:
- `types.go`
- `service.go`
- `signals.go`
- `clustering.go`
- `scoring.go`
- `policy.go`
- `proposals.go`

#### `internal/den`
Responsible for structural growth planning.

Suggested files:
- `types.go`
- `service.go`
- `capacity.go`
- `planner.go`
- `freeze.go`
- `experiments.go`

#### `internal/oversight`
Responsible for promotion, rollback, and policy enforcement.

Suggested files:
- `types.go`
- `service.go`
- `approvals.go`
- `promotion.go`
- `rollback.go`
- `policy_checks.go`

#### `internal/lineage`
Responsible for model and specialist ancestry tracking.

Suggested files:
- `types.go`
- `service.go`
- `graph.go`

---

## recommended interfaces

### SEAL interface shape
SEAL should work with three core concepts:
- `Signal`
- `GapCluster`
- `AdaptationProposal`

The minimum service interface should look like:

```go
package seal

type Service interface {
    ReviewSpecialist(ctx context.Context, specialistID string) ([]AdaptationProposal, error)
}
```

SEAL proposal surfaces should include:
- `slot_patch`
- `prompt_patch`
- `retrieval_patch`
- `tool_patch`
- `adapter_tuning`
- `specialist_split`
- `new_specialist`
- `expert_expansion`

### DEN interface shape
DEN should accept approved proposals and turn them into a `GrowthPlan`.

The minimum service interface should look like:

```go
package den

type Service interface {
    PlanGrowth(ctx context.Context, proposalID string) (*GrowthPlan, error)
}
```

DEN growth surfaces should include:
- `operational_patch`
- `adapter`
- `split_specialist`
- `new_specialist`
- `expert_bank`

Each plan should include:
- experiment name
- freeze plan
- rollback plan
- affected specialist or lineage node

### oversight interface shape
Oversight should approve proposals and finalize outcomes.

```go
package oversight

type Service interface {
    ApproveProposal(ctx context.Context, proposalID string) error
    PromoteExperiment(ctx context.Context, experimentID string) (*PromotionDecision, error)
    RollbackExperiment(ctx context.Context, experimentID string, reason string) error
}
```

---

## database additions
The following tables should be added by new migrations rather than by rewriting the original schema in place.

### `adaptation_signals`
Stores normalized runtime evidence.

Fields should include:
- `id`
- `specialist_id`
- `category`
- `severity`
- `surface`
- `task_class`
- `summary`
- `evidence_refs`
- `occurred_at`
- `created_at`

### `gap_clusters`
Stores SEAL’s grouped interpretation of repeated failures.

Fields should include:
- `id`
- `specialist_id`
- `category`
- `surface`
- `count`
- `persistence_score`
- `severity_score`
- `reversibility_score`
- `summaries`
- `evidence_refs`
- `created_at`

### `adaptation_proposals`
Stores SEAL output.

Fields should include:
- `id`
- `specialist_id`
- `cluster_id`
- `surface`
- `reason`
- `requested_by`
- `risk_level`
- `requires_harness`
- `requires_parent`
- `rollback_required`
- `status`
- `created_at`
- `approved_at`

### `growth_plans`
Stores DEN output.

Fields should include:
- `id`
- `proposal_id`
- `specialist_id`
- `surface`
- `reason`
- `experiment_name`
- `freeze_plan`
- `rollback_plan`
- `status`
- `created_at`

### `growth_experiments`
Stores actual experiment lifecycle state.

Fields should include:
- `id`
- `growth_plan_id`
- `specialist_id`
- `surface`
- `status`
- `artifact_refs`
- `metrics`
- `created_at`
- `completed_at`

### `lineage_nodes`
Stores any durable node in the system lineage.

Node types should eventually include:
- parent model
- specialist
- adapter
- branch
- expert bank

### `lineage_edges`
Stores relationships such as:
- spawned
- split into
- tuned into
- derived from
- retired into

### `slot_bundle_versions`
Stores compiled specialist slot bundles and source maps.

Fields should include:
- `id`
- `specialist_id`
- `bundle_toml`
- `source_map`
- `version_label`
- `created_by`
- `created_at`

### `promotion_decisions`
Stores harness or parent outcomes.

Fields should include:
- `id`
- `experiment_id`
- `decision`
- `reason`
- `approved_by`
- `created_at`

---

## runtime integration points

### execution
Execution should emit telemetry such as:
- tool failure
- unsupported output risk
- degraded fallback usage
- cross-modal grounding failure
- repeated execution degradation

### routing
Routing should emit telemetry such as:
- wrong specialist chosen
- no viable route
- repeated reroute
- text-only fallback overuse

### memory
Memory should emit telemetry such as:
- empty retrieval
- contradiction hit
- pointer miss
- context overflow triggered

### growth
The existing `growth` package should either:
- remain the experiment execution layer, or
- be narrowed over time while `seal` and `den` own decision-making

The preferred split is:
- `seal` decides whether to adapt
- `den` decides where to expand
- `growth` executes bounded experiments

---

## approval ladder

### operational-only changes
These include:
- `PROMPT.md`
- `SKILLS.md`
- `TOOLS.md`
- `HEARTBEAT.md`
- retrieval or memory heuristics

These should be specialist-proposed, SEAL-scored, and harness-approved.

### structural but non-constitutional growth
These include:
- new specialist creation
- specialist split
- adapter creation
- expert-bank creation
- branch creation

These should require:
- SEAL justification
- DEN planning
- harness verification
- parent approval

### constitutional changes
These include:
- identity changes
- authority widening
- governance rule changes
- slot mutability changes

These remain parent-governed only.

---

## freeze policy
The DEN layer should obey a conservative freeze policy:

- freeze the parent core by default
- freeze mature specialist cores by default
- allow plasticity only in:
  - new adapters
  - new specialists
  - explicitly immature modules
  - bounded experimental branches

The system should follow the rule:

**freeze what is mature; train what is missing.**

This keeps growth local and governed instead of allowing whole-system drift.

---

## implementation milestones

### milestone 1
**single-specialist adaptive loop works end to end**

Deliverables:
- adaptation signals table
- SEAL proposal generation
- harness approval for operational proposals
- compiled slot bundle persistence/versioning
- `verify` checks for signal -> proposal -> bundle-version flow

### milestone 2
**non-neural governed DEN growth exists**

Deliverables:
- DEN can choose among:
  - operational patch
  - retrieval patch
  - new specialist
  - specialist split
- lineage tables exist
- growth experiments are tracked
- rollback decisions are recorded

### milestone 3
**neural growth surfaces become real**

Deliverables:
- adapter experiments
- freeze-plan enforcement
- promotion decisions based on evals
- governed expert-bank experiments

---

## immediate next step
The next recommended implementation step is still **milestone 1 only**.

That means:
- add telemetry tables and writes
- add SEAL proposal generation
- persist compiled slot bundles
- extend `cmd/verify` to exercise the new loop

This gives the repository a real adaptive architecture before it commits to larger structural expansion.

---

## summary
The system architecture should treat:
- **SEAL** as the closed-loop adaptation governor
- **DEN** as the constrained capacity allocator and expansion mechanism
- the **harness** as the approving adult
- the **parent** as the constitutional sovereign
- **Postgres + pgvector** as the durable evidence, memory, and lineage plane
- **specialists** as bounded execution surfaces

The system should not become more capable by accident.
It should become more capable through governed evidence, bounded experiments, and reversible growth.
