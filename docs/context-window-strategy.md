# CONTEXT WINDOW STRATEGY

## purpose
Define the context-window policy for `SEAL_HAT_LLM` and the MM-ELLS architecture.

This document establishes two core rules:
- the parent may use a very large context window
- specialists remain much tighter and offload overflow to Postgres as structured summaries

---

## hard maximums

### parent
- **hard max window:** `2,000,000 tokens`

### specialists
- **hard max window:** `256,000 tokens`

These are hard architectural ceilings, not a license to treat the context window as infinite rolling transcript storage.

---

## core principle

**Context window is for active reasoning.**
**Postgres is for durable summarized continuity.**

The system should prefer:
- active task context in-window
- retrieved evidence in-window
- compact slot and policy state in-window
- summarized historical continuity in Postgres

It should avoid carrying oversized raw histories forward turn after turn when summary + retrieval will do the job better.

---

## default percentage budget

### parent default budget
- **system + constitutional slots:** `10–15%`
- **operational slots + summaries:** `10–20%`
- **retrieved memory + evidence:** `25–40%`
- **active task / working context:** `20–30%`
- **response headroom:** `10–15%`

### parent budget at 2,000,000 tokens
- **system + constitutional slots:** `200k–300k`
- **operational slots + summaries:** `200k–400k`
- **retrieved memory + evidence:** `500k–800k`
- **active task / working context:** `400k–600k`
- **response headroom:** `200k–300k`

### specialist default budget
- **system + constitutional slots:** `10–15%`
- **operational slots + summaries:** `10–20%`
- **retrieved memory + evidence:** `30–40%`
- **active task / working context:** `20–30%`
- **response headroom:** `10–15%`

### specialist budget at 256,000 tokens
- **system + constitutional slots:** `25.6k–38.4k`
- **operational slots + summaries:** `25.6k–51.2k`
- **retrieved memory + evidence:** `76.8k–102.4k`
- **active task / working context:** `51.2k–76.8k`
- **response headroom:** `25.6k–38.4k`

---

## overflow policy
When context becomes too large, the system should:
1. compress and deduplicate
2. summarize older working context
3. offload those summaries into Postgres
4. keep provenance and retrieval metadata
5. retrieve the summary later only when it is actually needed

This is the default overflow handling model for both parent and specialists.

---

## offload targets
Good candidates for offloading include:
- older chain-of-work summaries
- prior task-state summaries
- tool trace summaries
- intermediate synthesis notes
- arbitration summaries
- postmortem summaries
- long evidence bundles summarized with pointers to raw assets
- cross-specialist handoff summaries

Bad candidates for blind offloading include:
- current active instructions still needed for the task
- unresolved contradictions that still need live attention
- response headroom
- core constitutional and policy context

---

## required metadata for offloaded summaries
Offloaded summaries should not be stored as meaningless text blobs.
They should carry at least:
- summary text
- source references or asset pointers
- timestamp
- task/session linkage
- specialist id or parent id
- importance
- confidence
- retrieval tags
- modality metadata where relevant

That makes offloaded context retrievable and governable instead of merely archived.

---

## pressure thresholds

### parent thresholds
- **below 70%:** normal
- **above 70%:** begin compression and deduplication
- **above 85%:** summarize older traces to Postgres
- **above 92%:** allow only high-value raw inserts
- **above 97%:** emergency trim before continuation

### specialist thresholds
- **below 70%:** normal
- **above 70%:** summarize older local context
- **above 85%:** offload aggressively to Postgres
- **above 92%:** require summary substitution for raw history
- **above 97%:** emergency trim before continuation

---

## specialist discipline
Specialists have smaller windows for a reason:
- narrower lane
- faster inference
- lower cost
- less drift
- stronger retrieval discipline

A specialist should not depend on huge rolling histories when durable memory and retrieval can carry continuity more cleanly.

---

## parent discipline
The parent has a much larger window because it may need to support:
- routing
- orchestration
- arbitration
- governance
- cross-specialist synthesis
- multi-stage evidence handling

But even the parent should not treat 2M tokens as a substitute for structured memory.
The large window is a convenience and capability layer, not the primary memory system.

---

## design rule
The architecture must not require giant windows to remain coherent.
Large windows are valuable, but durable continuity still belongs in Postgres-backed summarized memory.

---

## summary
Default policy:
- **parent max:** `2,000,000 tokens`
- **specialist max:** `256,000 tokens`
- **budgeting:** proportional percentage budgets
- **overflow handling:** summarize and offload to Postgres
- **memory philosophy:** retrieval-first, not transcript-hoarding
