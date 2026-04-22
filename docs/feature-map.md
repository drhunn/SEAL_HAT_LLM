# Feature map

This document defines the repo's **feature families** and the **subfeatures** that live inside them.

Use this to keep feature branches narrow and to stop broad hand-waving from turning into branch sprawl.

## Rule

Use:
- **feature** for a real system lane
- **subfeature** for a bounded deliverable inside that lane

Branches stay at the **feature** level:
- `feature/<feature-name>`

Subfeatures are named in the feature spec and acceptance criteria.
They do not need their own long branch naming scheme unless the feature has clearly split into separate features.

---

## 1. Runtime task system

This is the core execution lane.

### Subfeatures
- task contract
- routing policy
- execution planning
- parent fallback behavior
- multimodal execution

---

## 2. Model unit runtime

This is the one-harness-per-model lane.

### Subfeatures
- unit spec
- unit registry
- local runtime bootstrap
- health and lifecycle

---

## 3. Storage and persistence

This is the storage lane. Do not smear it across runtime branches.

### Subfeatures
- shared DSN mode
- embedded Postgres bootstrap
- embedded Postgres lifecycle hardening
- route episode persistence
- specialist artifact persistence
- promotion and rollback

---

## 4. Task intake and inter-unit transport

This is where parent/specialist handoff belongs.

### Subfeatures
- file-backed task inbox
- specialist-side task RPC server
- parent-side task RPC dispatch
- parent result orchestration
- distributed task observability

---

## 5. Slot and bundle system

This is a separate feature family, not a runtime footnote.

### Subfeatures
- slot loading
- canonical bundle compilation
- bundle persistence
- slot sync

---

## 6. Governance and growth

This is the research core.

### Subfeatures
- SEAL proposal generation
- DEN growth staging
- admissibility and review
- evaluation and postmortem integration
- promotion decisions

---

## 7. Memory and sharing

This needs its own boundary because it can go bad fast.

### Subfeatures
- model-local memory
- cross-model sharing
- retrieval policy

---

## 8. Tool plane

This is later work, but still a real feature family.

### Subfeatures
- in-process tool usage
- shared RPC/IPC tool plane
- tool governance

---

## 9. Python HAT and training

Separate subsystem. Keep it that way.

### Subfeatures
- dataset building
- corpus ingestion
- training scaffold
- model artifact packaging

---

## Recommended feature priorities

If the goal is to move the repo from scaffold to functioning specialist system, prioritize these features first:

1. `task-intake-and-inter-unit-transport`
   - especially:
     - specialist-side task RPC server
     - parent-side task RPC dispatch
     - parent result orchestration
2. `runtime-task-system`
   - especially routing/execution policy cleanup
3. `storage-and-persistence`
   - especially embedded Postgres lifecycle hardening
4. `governance-and-growth`
   - especially making the evidence/change loop less theatrical and more real

---

## Branch naming convention

Use feature branches like:
- `feature/runtime-task-system`
- `feature/model-unit-runtime`
- `feature/storage-and-persistence`
- `feature/task-intake-and-inter-unit-transport`
- `feature/slot-and-bundle-system`
- `feature/governance-and-growth`
- `feature/memory-and-sharing`
- `feature/tool-plane`
- `feature/python-hat-and-training`

Inside the feature spec, name the included subfeatures explicitly.

Example:
- branch: `feature/task-intake-and-inter-unit-transport`
- spec includes subfeatures:
  - specialist-side task RPC server
  - parent-side task RPC dispatch
  - parent result orchestration

That keeps branches readable and keeps subfeature scope explicit without inventing a second branch taxonomy.

---

## What a feature spec should contain

Every feature spec should say:
- which feature family it belongs to
- which subfeatures are included in this slice
- what is explicitly out of scope
- the acceptance criteria
- the tests required for the feature to count as complete
- known risks or follow-on subfeatures

If the branch cannot name the feature family and the included subfeatures, it is probably too vague.
