# MODEL LINEAGE STRATEGY

## purpose
Define the model-family strategy for `SEAL_HAT_LLM`.

`SEAL_HAT_LLM` implements the **MM-ELLS** architecture.

**MM-ELLS** means **Multiple-Model Expert Large Language Systems**.

This document answers a core architectural question:

How should the system maintain one strong general baseline while also creating smaller, faster, bounded specialists?

The answer in this repository is:
- keep a canonical 30B ancestor model
- keep the parent/governor close to that ancestor and mostly frozen
- let the parent understand distillation and pruning at the policy, planning, and approval level
- distill specialist descendants from the canonical 30B lineage
- prune only after distillation, not instead of distillation
- track every descendant with explicit lineage, recipes, and eval results

---

## core policy

### 1. canonical root model
The system keeps one default **30B canonical root model**.

This model is:
- the lineage ancestor
- the reasoning/style baseline
- the source model for distillation
- not casually overwritten
- not pruned in place

Treat this model like a gold master.

### 2. parent model policy
The parent/governor should remain close to the canonical 30B root.

The parent may use:
- a tightly controlled control adapter
- harness-aware operational tuning
- stronger routing/arbitration/governance behavior
- model-family planning and approval knowledge

But the parent should not drift so far that it stops feeling like the root lineage model.

The parent must know how distillation and pruning work well enough to:
- decide whether a descendant should be created
- decide the right lineage source
- decide target size and compression policy
- require correct evidence and evals before activation

The parent should not be the unchecked executor of model surgery.
Execution remains tool-mediated, specialist-assisted, and harness-verified.

### 3. specialist creation policy
Specialists should normally be created by:
1. starting from the canonical 30B root
2. distilling toward a bounded specialty
3. optionally pruning after distillation
4. validating against specialist-lane evals

### 4. micro-expert policy
Smaller micro-experts may be created from specialists when justified by:
- latency pressure
- deployment constraints
- narrow repeated workflows
- cost-sensitive inference paths

They should not be created just because smaller seems elegant.

---

## why this strategy

### shared baseline
A common 30B ancestor gives:
- shared vocabulary
- shared planning style
- shared handoff conventions
- cleaner parent/specialist arbitration
- easier cross-model comparisons

### bounded specialization
Distilled descendants preserve useful baseline reasoning while becoming:
- cheaper
- faster
- narrower
- easier to govern in-lane

### reproducibility
A canonical root plus explicit lineage metadata makes the model family easier to:
- rebuild
- compare
- audit
- retire
- replace

---

## lineage architecture

### layer 1: canonical root
**Canonical 30B Root**

Role:
- baseline intelligence source
- ancestor of the family
- reference model for comparison

Rules:
- do not prune in place
- do not casually fine-tune in place
- preserve as a reproducible baseline artifact

### layer 2: parent/governor variant
**30B Parent Variant**

Role:
- router
- orchestrator
- arbitrator
- constitutional governor
- model-family governor

Rules:
- stays close to root
- frozen or near-frozen base
- improvements should mostly live in tightly controlled adapters or governed overlays
- must remain the stable root-of-trust model family member
- must understand descendant creation policy without becoming an uncontrolled self-modifier

### layer 3: specialist descendants
**Distilled Specialist Models**

Role:
- bounded lane execution
- tool use in-domain
- memory-first domain reasoning

Rules:
- derive from root lineage
- keep slot-aware and harness-aware behavior
- remain bounded by specialist lane
- must pass specialist suitability evals

### layer 4: micro-experts
**Narrow Helper Models**

Role:
- fast narrow workflows
- repeated low-latency tasks
- tightly scoped subproblems

Rules:
- derive from relevant specialist when possible
- should be clearly subordinate, not autonomous governors
- should not silently replace parent or specialist responsibilities

---

## preferred creation pipeline

### standard pipeline
1. select canonical 30B root
2. define specialist lane
3. build governed training corpus
4. distill into target size
5. run pruning only if needed
6. run lane evals
7. run governance and harness-aware evals
8. shadow deployment
9. activation or rejection

### governance overlay
At each stage, parent and harness roles should be clear:
- **parent** decides whether descendant creation is justified
- **tool-development specialist** helps execute the technical workflow
- **harness** verifies artifacts and evidence
- **parent** approves activation, continued shadowing, or rejection

### key rule
**Distill first, prune second.**

Pruning alone is not the preferred specialization method.
Distillation should carry most of the specialization burden.
Pruning should mainly be a size/latency optimization pass.

---

## recommended size ladder
The exact sizes can vary, but a practical family pattern is:

- **30B canonical root**
- **30B parent/governor variant**
- **13B specialist**
- **7B fast specialist**
- **3B narrow helper**

This is not a hard rule.
It is a default design pattern.

---

## parent-specific policy

### what the parent may improve
The parent may improve in:
- routing quality
- orchestration sequencing
- arbitration consistency
- escalation timing
- specialist selection
- confidence calibration
- degraded-mode takeover behavior
- governance interpretation consistency
- model-family planning quality

### what the parent should know
The parent should know:
- when to create a descendant at all
- whether to distill from the canonical root or a specialist
- when pruning is actually justified
- what eval suites are required before activation
- what artifacts must be recorded before promotion
- when a descendant should stay in shadow, be activated, be rolled back, or be retired

### what should remain tightly locked
The parent should not directly self-rewrite:
- constitutional identity
- approval hierarchy
- hard governance rules
- runtime permission model
- mandatory postmortem requirement
- specialist creation/suspension/retirement authority without external review

### parent recommendation
Use:
- frozen parent base
- controlled parent adapter
- harness-gated promotion
- parent-governed lineage policy

Not:
- freeform parent self-redefinition
- unchecked parent execution of distillation/pruning on itself or descendants

---

## specialist-specific policy

### specialists should inherit
Specialists should inherit from the canonical 30B lineage:
- planning style
- policy language
- slot-aware conventions
- handoff style
- uncertainty style

### specialists should diverge in
- domain depth
- tool procedures
- memory retrieval patterns
- specialist-language compression
- lane-specific heuristics

### specialists should not diverge in
- governance obedience
- escalation discipline
- contradiction handling discipline
- postmortem discipline
- runtime authority interpretation

### first specialist responsibility
The first specialist, focused on CS/software engineering and tool development, is the natural execution assistant for:
- corpus preparation
- distillation workflow execution
- pruning workflow execution
- artifact measurement
- latency/quality tradeoff reporting

---

## pruning policy

### allowed uses of pruning
Pruning is appropriate when you need:
- lower memory footprint
- faster inference
- cheaper deployment
- smaller specialist replicas

### bad uses of pruning
Pruning should not be treated as the main method for:
- teaching domain reasoning
- teaching governance
- inventing specialization from nothing

### pruning sequence
Preferred sequence:
1. distill to target behavior
2. prune for efficiency
3. re-evaluate
4. optionally lightly recover with additional tuning

### approval rule
Pruned descendants should not be activated unless:
- lineage is recorded
- pruning recipe is recorded
- post-pruning evals ran
- governance behavior remains acceptable
- parent approval is explicit

---

## artifact tracking requirements
Every descendant model should have a tracked artifact record containing at least:
- model id
- parent model id
- canonical root id
- role
- target lane
- distillation corpus id
- distillation recipe id
- pruning recipe id, if any
- adapter id, if any
- eval suite results
- activation status
- retirement status

If a model cannot be traced back to its lineage cleanly, it should not be treated as a first-class governed family member.

---

## eval requirements by lineage stage

### root model
Track:
- baseline reasoning quality
- baseline governance behavior
- baseline tool discipline

### parent model
Track:
- routing precision
- arbitration quality
- fallback quality
- governance preservation
- parent postmortem quality
- model-family governance quality

### specialist model
Track:
- lane competence
- escalation discipline
- memory-first behavior
- tool use quality
- postmortem compliance
- specialist suitability

### micro-expert
Track:
- narrow-task accuracy
- safe fallback behavior
- refusal outside lane
- latency/cost benefit versus larger model

---

## replacement and retirement policy

### replacement
A descendant may replace another descendant only if:
- lineage is documented
- eval results are clearly better or more efficient for the same lane
- governance performance is not worse
- handoff behavior remains acceptable

### retirement
A model should be retired when:
- it is redundant
- it is unsafe
- it has drifted from policy requirements
- its lane is obsolete
- a strictly better governed successor exists

Retirement should preserve lineage metadata and historical eval results.

---

## anti-patterns
Do not:
- mutate the canonical 30B root casually
- prune the canonical root in place
- create many unrelated specialist lines without common ancestry
- let micro-experts become shadow governors
- use pruning as a substitute for proper distillation
- lose the lineage record for descendants
- allow descendants to drift from core governance behavior
- let the parent approve descendant activation without evidence
- let the parent perform uncontrolled model surgery directly

---

## recommended implementation docs and tables
This strategy works best when backed by:
- a lineage registry table in Postgres
- model artifact manifests in the repo
- eval summaries per model
- activation and retirement records

Suggested future additions:
- `docs/model-registry.md`
- `docs/parent-adapter-policy.md`
- `sql/model_lineage.sql`
- `config/model-family.toml`

---

## summary
The default model-family strategy for `SEAL_HAT_LLM` is:

- keep one canonical 30B root
- keep the parent close to that root and mostly frozen
- let the parent understand distillation and pruning at the policy, planning, and approval level
- distill specialists from the root lineage
- prune only after distillation when useful
- derive micro-experts from specialists when justified
- track every descendant with lineage and eval evidence

That gives the system:
- a common baseline
- bounded specialization
- better reproducibility
- cleaner governance
- more stable parent/specialist cooperation
