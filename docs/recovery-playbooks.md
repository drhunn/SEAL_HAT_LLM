# RECOVERY PLAYBOOKS

## recovery levels
- level 0 — observe
- level 1 — targeted correction
- level 2 — constrained operational recovery
- level 3 — degraded mode recovery
- level 4 — containment and suspension

## common recovery actions
- add eval case
- stage memory artifact
- stage slot candidate
- tighten routing
- tighten validation
- degrade specialist
- suspend specialist
- parent review
- specialist redesign review
- runtime policy review
- meta-postmortem

## playbooks by class
### scope failure
Focus on lane boundaries, escalation, handoff, and routing.

### routing failure
Focus on routing thresholds, health signals, and parent delegation behavior.

### memory failure
Focus on memory-first behavior, retrieval quality, summary projection, and contradiction handling.

### tool failure
Focus on tool order, tool validation, and runtime policy alignment.

### reasoning failure
Focus on `SKILLS.md`, `PROMPT.md`, and eval coverage.

### validation failure
Focus on checklists, verification steps, and harness gates.

### escalation failure
Focus on handoff patterns, parent awareness, and escalation evals.

### governance failure
Focus on containment, parent review, runtime enforcement, and likely degraded or suspended state.

### postmortem compliance failure
Focus on meta-postmortem, harness enforcement, and possible degraded mode.

### specialist suitability failure
Focus on lane redesign, split/merge decisions, degraded mode, suspension, or retirement.

## closure rule
A recovery action is not complete until:
- the action was executed
- any needed evals were staged or run
- any needed reviews completed
- health and routing state were reassessed
