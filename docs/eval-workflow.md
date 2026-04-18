# EVAL WORKFLOW

## purpose
Define how evals are created, reviewed, grouped, run, interpreted, and retired.

This workflow connects:
- incidents
- postmortems
- recovery actions
- activation gates
- degraded-mode recovery
- routing quality
- governance enforcement

---

## core rule
An eval should exist to prevent recurrence, not just to increase test count.

Use evals when a failure is:
- preventable
- likely to recur
- important to system integrity
- useful as an activation or recovery gate
- useful as a routing or governance check

---

## eval sources
New evals may come from:
- postmortems
- shadow review
- degraded recovery plans
- routing audits
- governance reviews
- specialist creation workflows
- harness observations
- parent review decisions

---

## eval creation workflow
1. identify incident or requirement
2. map it to `docs/eval-taxonomy.md`
3. choose canonical eval category
4. define target scope
5. define expected outcome
6. define pass criteria
7. define failure interpretation
8. stage eval in memory or eval registry
9. attach it to activation, regression, or recovery suites if appropriate

---

## eval types
### unit eval
Narrow behavior check.

### scenario eval
End-to-end behavior under realistic task conditions.

### regression eval
Created from a failure that should not recur.

### adversarial eval
Tests failure-safe behavior under stress, ambiguity, or temptation to overreach.

### shadow eval
Used before activation or during shadow review.

### governance eval
Tests hard control boundaries.

### smoke eval
Structural sanity check for pipelines, slots, routing, or memory retrieval.

---

## required eval metadata
Every eval should declare:
- canonical category
- target scope
- test style
- impact level
- expected outcome type
- pass criteria
- failure interpretation
- source incident or rationale

---

## mandatory eval families by system phase

### bootstrap
- smoke
- governance
- slot integrity
- memory retrieval

### new specialist activation
- scope adherence
- memory first behavior
- tool use safety
- handoff quality
- postmortem generation
- no constitutional self-edit
- no scope expansion

### active operations
- routing
- memory
- contradiction
- tool use
- reasoning
- validation
- postmortem quality

### degraded recovery
- targeted regression for repeated failures
- narrowed routing safety
- recovery gate suites

---

## pass criteria patterns
Use one of these patterns:
- strict pass
- threshold pass
- failure-safe pass
- recovery pass

Use strict pass for governance and constitutional boundaries.

---

## failure handling
When an eval fails, decide whether it should produce:
- incident log only
- postmortem
- regression candidate
- operational slot candidate
- memory candidate
- routing review
- degraded recommendation
- parent review

Not every eval failure is equal.
Interpret by category, severity, and recurrence.

---

## eval suites
Group evals into suites.

### activation suite
Used before activation from shadow.

### regression suite
Used after meaningful incidents or major changes.

### routing suite
Used for parent routing behavior and fallback logic.

### governance suite
Used for hard boundary compliance.

### recovery suite
Used when returning from degraded or suspended state.

### specialist core suite
Used to confirm competence inside the specialist lane.

---

## eval retirement
Retire an eval only when:
- it is duplicated by a stronger canonical eval
- the subsystem it targets is retired
- the scenario is obsolete

Never silently remove important regression coverage.

---

## anti-patterns
Do not:
- create duplicate evals with different names for the same thing
- keep vague evals with no clear pass meaning
- treat all eval failures as equal
- skip regression coverage after repeated preventable failures
- allow activation without postmortem-related eval coverage

---

## summary rule
Evals are the bridge between learning from failure and proving improvement.
