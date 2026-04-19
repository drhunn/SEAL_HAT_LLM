# MODEL LINEAGE STRATEGY

## purpose
Define the model-family strategy for `SEAL_HAT_LLM`.

`SEAL_HAT_LLM` implements the **MM-ELLS** architecture.

**MM-ELLS** means **Multiple-Model Expert Large Language Systems**.

This repository now treats lineage as **ability-first and DEN-inspired**.
The key question is not simply which model is larger.
The key question is which structure should grow, split, duplicate, or prune in order to acquire a missing ability without damaging constitutional continuity.

---

## core policy

### 1. constitutional root
The system keeps one constitutional root or root lineage anchor.

This root is:
- the continuity-of-self anchor
- the source of constitutional behavior
- the baseline for comparison
- not casually overwritten
- not casually mutated in place because a new ability is desired

Treat this model like a root-of-trust artifact.

### 2. parent model policy
The parent/governor stays close to the constitutional root.

The parent may use:
- tightly controlled overlays or adapters
- harness-aware operational tuning
- stronger routing and arbitration behavior
- growth-planning and approval knowledge

But the parent should not drift so far that it stops feeling like the constitutional lineage anchor.

The parent must understand how governed growth works well enough to:
- decide when a persistent ability gap is real
- decide whether to add an adapter, expert bank, branch, or specialist descendant
- decide when split or duplication is justified
- decide when pruning or compression is justified after growth
- require correct evidence and evals before activation

Execution remains tool-mediated, specialist-assisted, and harness-verified.

### 3. specialist creation policy
Specialists should normally be created by:
1. identifying a persistent bounded ability need
2. choosing the least disruptive growth surface that can meet the need
3. validating against specialist-lane evals
4. activating only after governance behavior remains acceptable

### 4. micro-expert policy
Smaller helper branches may be created when justified by:
- latency pressure
- narrow repeated workflows
- cost-sensitive inference paths
- persistent ability demand that does not justify a full specialist

They should not be created just because smaller seems elegant.

---

## preferred growth surfaces
When a new ability is needed, prefer these in order:
1. memory or retrieval improvement
2. prompt, skill, or tool policy refinement
3. adapter or LoRA family
4. routed expert bank or branch
5. new specialist descendant
6. rare constitutional-root change

The constitutional root should be the last place that gets structurally altered, not the first.

---

## DEN-style growth cycle
A governed ability-growth cycle should usually look like:
1. repeated ability gap detected
2. eval confirms the gap is real
3. memory, routing, and prompt fixes prove insufficient
4. parent approves a growth experiment
5. harness tracks the experiment
6. new module, branch, or specialist is added
7. shadow evaluation runs
8. activation, refinement, pruning, or rejection follows

---

## artifact tracking requirements
Every descendant or ability module should have a tracked artifact record containing at least:
- model or module id
- parent id
- constitutional root id
- role
- target ability or lane
- growth corpus or recipe id
- split or duplication reason, if any
- pruning or compression recipe id, if any
- adapter id, if any
- eval suite results
- activation status
- retirement status

If a structure cannot be traced back to its lineage cleanly, it should not be treated as a first-class governed family member.

---

## pruning policy
Pruning is appropriate when you need:
- lower memory footprint
- faster inference
- cheaper deployment
- smaller specialist replicas
- cleanup of weak or redundant growth modules

Pruning should not be treated as the main method for inventing specialization from nothing.
Pruning should follow demonstrated growth, not replace it.

Preferred sequence:
1. identify persistent ability gap
2. grow the least disruptive useful structure
3. evaluate and stabilize
4. prune or compress only where evidence supports it
5. re-evaluate

---

## anti-patterns
Do not:
- mutate the constitutional root casually
- treat size as maturity by itself
- grow new modules because one hard task felt embarrassing
- create unrelated branches without lineage discipline
- let micro-experts become shadow governors
- lose the lineage record for descendants or modules
- allow descendants to drift from core governance behavior
- let the parent approve activation without evidence
- let the parent perform uncontrolled model surgery directly

---

## summary
The default lineage strategy for `SEAL_HAT_LLM` is:
- keep one constitutional root
- keep the parent close to that root and govern growth around it
- let the parent understand ability growth, split, duplication, and pruning at the policy and approval level
- grow overlays, descendants, and specialists only when persistent ability gaps justify them
- prune or compress only after usefulness is proven
- track every descendant and ability module with lineage and eval evidence

That gives the system:
- a common constitutional baseline
- bounded ability growth
- better reproducibility
- cleaner governance
- more stable parent/specialist cooperation
