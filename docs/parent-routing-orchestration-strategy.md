# PARENT ROUTING AND ORCHESTRATION STRATEGY

## purpose
Define how the parent generalist should evolve from a fixed routing policy into a **learning routing/orchestration layer** under harness and slot governance.

This document is a research-plan document.
It does not claim that the current runtime already implements all of this.

---

## core goal
The parent should get better at deciding:
- whether to handle a task itself
- whether to route to a specialist
- which specialist is the best target
- whether retrieval should happen first
- whether multiple specialists or multimodal fusion are required
- whether escalation, deferment, or refusal is correct

The parent should improve at routing and orchestration **without** gaining the authority to widen governance boundaries on its own.

The specialists the parent routes to are supposed to be **explicit smaller models derived from the base model**, not internal MoE shards.
That means routing quality is not just about correctness.
It is also about sending repeated narrow work to the smaller descendant model that can do it faster, cheaper, and with better lane-specific resolution than the parent.

Those specialists are intended to run as **model units**, not as loose executor names.
A routing target should ultimately mean:
- model artifact
- local harness
- local embedded Postgres
- local slots and lifecycle state

---

## non-goals
This strategy does **not** allow the learned parent to:
- redefine constitutional slots
- widen tool permissions on its own
- bypass review requirements
- create or retire specialists directly
- perform DEN-style structural changes without SEAL justification and harness approval

The parent may learn **how to route and orchestrate**.
It may not learn **what rules exist**.

---

## current starting point
The current runtime already has:
- explicit routing inputs
- a deterministic executor-selection policy
- execution planning
- routing and execution telemetry helpers
- task processing
- a conservative task inbox
- postmortem and eval staging after incidents

That is enough to start collecting route episodes and benchmarking routing quality.

It is **not** yet the final architecture.
The current code still routes to executor names inside one runtime scaffold rather than to fully isolated model units with their own harnesses and embedded databases.

---

## learning surfaces

### 1. route selection
The parent should learn:
- parent vs specialist
- specialist A vs specialist B
- single specialist vs fusion
- route confidence

### 2. orchestration sequencing
The parent should learn:
- retrieval before execution
- specialist chaining
- reroute after failure
- escalation after repeated degradation

### 3. fallback behavior
The parent should learn when:
- text-only fallback is acceptable
- a task should be rejected as out of lane
- a task should be deferred for review

### 4. tool-use planning
Later versions should let the parent learn tool-use planning across a shared tool plane while still respecting local harness policy.
That means learning when:
- a local model unit should call a shared tool executable
- the tool result is sufficient for the task
- specialist routing should happen before or after tool use
- a tool result should be rejected or reviewed by the local harness

---

## required data model
The runtime should emit a **route episode** for each task.

Minimum fields:
- task id
- task summary
- task class
- primary modality
- secondary modalities
- available specialists
- chosen route
- confidence
- fallback flag
- review flag
- retrieval-first flag
- execution result
- final eval result
- incident/postmortem flags
- latency
- cost if available
- whether rerouting occurred

Without a route-episode record, the parent cannot improve in any serious way.

A fuller later route-episode schema should also capture:
- target model-unit id
- target harness id or version
- local store snapshot or revision ref if relevant
- shared tool invocation refs, if tools were used
- whether the route crossed model-unit boundaries

---

## labels and metrics

### route labels
At minimum, classify episodes as:
- correct route
- acceptable route
- unnecessary handoff
- missed specialist handoff
- wrong specialist
- wrong fusion decision
- wrong escalation decision
- failed despite correct routing
- failed due to bad routing

### metrics
Track at least:
- route correctness
- handoff precision
- handoff recall
- unnecessary handoff rate
- missed specialist rate
- wrong fusion rate
- wrong escalation rate
- retrieval-before-routing miss rate
- overall task success after routing
- latency / cost per route class

Because the specialists are intended to be smaller derived models, the routing benchmark should also measure:
- parent latency vs specialist latency
- parent cost vs specialist cost
- lane-specific quality difference between parent and specialist
- cases where the parent should keep work because the specialist does not actually outperform it yet
- local tool-call cost vs specialist-handoff cost in the later shared-tool architecture

---

## training and deployment plan

### phase 1: shadow mode
Keep the current deterministic router as the live path.
Add a shadow parent router that proposes:
- chosen executor or executor set
- confidence
- rationale
- retrieval-first recommendation
- review recommendation

Do not let the shadow router control execution yet.
Just log its choices and compare them against actual outcomes.

### phase 2: offline learning
Train on route episodes using one or more of:
- ranking
- classification
- bandit-style policy improvement
- harness-aware SFT on route decision traces

Start small.
The goal is not to make the parent mystical.
The goal is to make it measurably better than the fixed baseline.

### phase 3: gated online use
Allow the learned parent to control routing only when:
- confidence is high
- harness policy passes
- fallback exists
- the route is auditable

All low-confidence cases should fall back to the deterministic baseline.

### phase 4: orchestration learning
After basic routing works, let the parent learn:
- multi-step plans
- specialist chaining
- retrieval-before-routing
- reroute-after-failure
- fusion planning
- shared tool-use sequencing across model units

---

## relation to SEAL and DEN

### parent
Learns how to route and orchestrate better.

### SEAL
Watches routing/orchestration outcomes and decides whether repeated failures justify adaptation.
Examples:
- routing policy keeps missing the same task family
- fusion is overused or underused
- the current specialist layout is inadequate
- the parent cannot route a recurring task class correctly
- the parent keeps sending work to a specialist that does not actually outperform it
- the current model-unit layout is wrong even if the routing policy itself is reasonable

### DEN
Acts only after SEAL justification and harness approval.
Examples:
- split one specialist into two
- create a new specialist
- add a branch or adapter surface
- restructure specialist coverage
- produce a smaller distilled/pruned descendant for a recurring task family
- produce a full model-unit bundle with local harness and embedded Postgres bootstrap

So the parent’s learning loop is a **runtime capability improvement**.
DEN is a **structural change mechanism**.
They are not the same thing.

---

## governance rules
The following must remain hard-governed:
- slot mutability boundaries
- tool permissions
- constitutional protections
- review requirements
- promotion and rollback policy
- structural DEN changes
- cross-model sharing and replication policy

The parent can optimize within those boundaries.
It cannot rewrite them.

---

## tool-plane note
Later versions are intended to move tools into a shared RPC/IPC tool plane where multiple harnesses can use the same tool executable.

That does **not** mean shared authority.
Each local harness must still decide:
- whether a tool may be invoked
- whether the result may be trusted
- whether the result may be written into local memory
- whether the result triggers escalation or review

A shared tool executable is shared capability, not shared governance.

---

## minimum deliverables
A serious first research milestone should include:
1. route episode schema
2. a persisted route-episode store
3. a fixed routing benchmark task set
4. a baseline deterministic router report
5. a shadow learned-router report
6. side-by-side routing metrics
7. incident and eval comparison by route class

Without those, “learning orchestration” is just a slogan.

---

## immediate next steps
1. define the route-episode schema
2. persist route episodes for every runtime task
3. define the routing benchmark corpus
4. run the deterministic baseline
5. add a shadow learned router
6. compare outcomes before allowing any learned routing in the live path

---

## summary
The parent should become a **learning routing/orchestration layer**.
It should improve through measured route episodes, benchmarked comparisons, and harness-gated deployment.
It should route work toward explicit smaller derived specialist models when those models actually provide better lane-specific speed, cost, or quality.
It should ultimately route across **model units** with local harnesses and embedded databases rather than flat executor names, and later coordinate shared tool use through a reusable RPC/IPC tool plane.
SEAL should decide when routing/orchestration failure justifies deeper change.
DEN should perform only the approved structural changes.
The harness and slot governance remain the hard boundary the whole time.
