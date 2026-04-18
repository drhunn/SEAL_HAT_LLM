# LAYERING STRATEGY

## purpose
Define the recommended implementation layering model for `SEAL_HAT_LLM` and MM-ELLS.

The architecture may be described in richer conceptual planes, but the implementation should stay compact enough to build, test, and maintain.

---

## recommended implementation layering
Use **5 real implementation layers**.

1. **Interface layer**
2. **Governance layer**
3. **Execution layer**
4. **Oversight layer**
5. **Persistence layer**

This is the preferred buildable runtime structure.

---

## layer 1: interface layer
Responsibilities:
- user/task intake
- request shaping
- attachment intake
- session entry state
- task envelopes and metadata

This layer should not become the governor.
It prepares work for the governance layer.

---

## layer 2: governance layer
Responsibilities:
- parent routing
- arbitration
- specialist selection
- lifecycle approvals
- model-family planning
- context-budget governance
- constitutional review

This is where the parent and parent-governed logic live.

Sub-roles can exist here, but they should remain governance helpers rather than independent sovereigns.

---

## layer 3: execution layer
Responsibilities:
- specialists
- tool-use flows
- task completion
- handoffs
- narrow task workers

This is where domain work gets done.
The execution layer should not own constitutional authority.

---

## layer 4: oversight layer
Responsibilities:
- harness logic
- postmortems
- eval generation and later execution
- recovery planning
- health scoring
- degraded/suspended recommendations

This layer exists so the system can improve without dissolving its own boundaries.

---

## layer 5: persistence layer
Responsibilities:
- Postgres
- pgvector
- lineage records
- audit trail
- durable memory
- retrieval
- summary offload storage
- multimodal asset references

This layer is not just storage. It is the durable continuity substrate.

---

## why five layers
Five layers are enough separation to keep the system clean without turning it into procedural theater.

### too few layers
If there are too few layers, the architecture tends to blur:
- governance into execution
- harness into runtime glue
- memory into an afterthought

### too many layers
If there are too many layers, the architecture tends to bloat:
- latency rises
- ownership gets fuzzy
- integration becomes harder than the actual problem
- the structure starts becoming the project

---

## conceptual planes versus implementation layers
MM-ELLS can still use richer conceptual descriptions such as:
- interaction
- parent governance
- harness
- specialist execution
- memory
- eval
- training

But those should compress into the 5 real implementation layers when building the actual runtime.

That keeps the docs expressive and the implementation sane.

---

## sublayers
Sublayers are allowed inside a layer.
For example:

### governance sublayers
- router
- arbitrator
- lineage governor
- context-budget governor

### execution sublayers
- specialists
- tool runners
- handoff manager
- ephemeral worker agents

### oversight sublayers
- postmortem engine
- eval engine
- recovery planner
- health scorer

These should stay sublayers, not become a forest of new top-level layers.

---

## implementation rule
Keep the top-level layering shallow.
Use sublayers only when they clearly improve clarity or testing.

If a new layer exists only because one module felt important, it probably should not be a top-level layer.

---

## summary
The recommended implementation structure for MM-ELLS is:
- **5 real layers**
- richer conceptual planes allowed in docs
- sublayers allowed internally
- no uncontrolled layer explosion

That gives the architecture enough structure to be serious without letting the structure become the system.
