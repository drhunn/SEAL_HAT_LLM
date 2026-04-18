# SUB-AGENT STRATEGY

## purpose
Define how sub-agents should exist inside `SEAL_HAT_LLM` and the MM-ELLS architecture.

Sub-agents are allowed, but they must remain bounded helpers rather than independent sovereigns.

---

## two kinds of sub-agents
The preferred model is to support **two kinds of sub-agents**.

### 1. ephemeral sub-agents
These are short-lived workers created for one task or one narrow step.

Examples:
- source reader
- retrieval triager
- tool plan writer
- diff checker
- citation checker
- failure reviewer
- postmortem drafter

Properties:
- task-scoped
- narrow context
- no durable identity required
- no independent lifecycle authority
- usually no direct durable memory writes

### 2. persistent sub-agents
These are longer-lived helper roles under the parent or under a specialist.

Examples:
- parent routing aide
- parent lineage-review aide
- specialist tool-planning aide
- specialist memory-triage aide
- harness postmortem reviewer

Properties:
- repeated helper role
- bounded lane under a caller
- reduced authority compared with the caller
- should be used sparingly

---

## authority model
Sub-agents are **workers**, not mini-governors.

A sub-agent should:
- inherit a reduced authority set from its caller
- have one clear caller
- have one narrow purpose
- have one output contract

A sub-agent should not:
- own constitutional authority
- own lifecycle authority
- self-authorize tool expansion
- become an untracked shadow hierarchy
- write durable truth directly unless explicitly allowed through governed paths

### authority reduction rule
- parent sub-agent < parent
- specialist sub-agent < specialist
- harness helper < harness

Equal authority should never be the default.

---

## hierarchy depth
Keep the hierarchy shallow.

Preferred depth:
- **Parent**
  - **Specialist**
    - **Sub-agent**

That is the normal stopping point.

A sub-agent spawning sub-agents should be rare and tightly controlled.
If the design starts looking like a bureaucracy simulator, it has already gone too far.

---

## memory model
Sub-agents should usually use:
- temporary working memory
- task-local context
- caller-provided memory slices

Sub-agents should usually not:
- write durable memory directly
- promote memory directly
- update constitutional or operational slots directly

Durable results should normally be written through:
- the caller
- the harness
- explicit reviewed memory/write paths

---

## when to use sub-agents
Sub-agents are useful for:
- parallel reading
- evidence triage
- comparing candidate answers
- tool-call planning
- retrieval filtering
- postmortem drafting
- eval drafting
- policy-compliance checking before promotion

They are especially useful when a narrow helper can reduce the load on the parent or on a specialist without needing full specialist status.

---

## when not to use sub-agents
Do not use a sub-agent when the role clearly needs:
- durable identity
- long-term memory
- independent lane ownership
- repeated independent routing
- independent promotion/suspension/retirement semantics

If it needs those things, it is not really a sub-agent anymore.
It is trying to become a specialist.

---

## boundary rule
Use this rule:

**If it needs identity, long-term memory, its own lane, and repeated independent routing, it is no longer a sub-agent. It is a specialist.**

That is the clean boundary.

---

## parent sub-agents
Good parent sub-agents include:
- routing aide
- arbitration aide
- lineage-review aide
- context-budget aide

These should help the parent make better decisions, but they should not become shadow governors with independent constitutional power.

---

## specialist sub-agents
Good specialist sub-agents include:
- retrieval aide
- tool-planning aide
- evidence summarizer
- diff checker
- local quality checker

These should help the specialist execute better, but they should remain clearly subordinate to the specialist.

---

## harness helpers
The harness can also have helper sub-agents such as:
- postmortem reviewer
- eval drafter
- evidence normalizer
- contradiction reviewer

These remain helpers under harness control, not independent oversight sovereigns.

---

## persistence policy
Default preference:
- **ephemeral sub-agents by default**
- **persistent sub-agents only when repeated value is obvious**

Persistent helper roles should be justified by repeated operational value, not by aesthetic preference.

---

## implementation guidance
The implementation should support sub-agents as:
- bounded helper roles
- temporary worker contexts
- narrow tool scopes
- caller-owned result handling

The implementation should avoid:
- untracked sub-agent sprawl
- recursive agent trees by default
- persistent hidden state outside governed memory systems

---

## relationship to layering
Sub-agents live mostly inside the **execution layer** and **governance layer** as sublayers, not as new top-level implementation layers.

That means:
- the 5-layer implementation model still holds
- sub-agents are internal helpers within those layers
- sub-agents do not justify exploding the top-level architecture

---

## summary
The recommended sub-agent model for MM-ELLS is:
- support **ephemeral sub-agents** by default
- support **persistent sub-agents** sparingly
- keep sub-agents subordinate and narrow
- keep hierarchy depth shallow
- route durable writes through governed paths
- promote a role to **specialist** when it clearly needs identity, memory, lane ownership, and repeated routing
