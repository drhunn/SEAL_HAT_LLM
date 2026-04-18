# EVAL TAXONOMY

## purpose
Provide a shared category system for evals created from:
- postmortems
- recovery actions
- shadow-mode review
- routing review
- specialist activation gates
- degraded-mode recovery
- governance checks

This file ensures evals are named, grouped, and interpreted consistently across the parent, specialists, and harness.

---

## core rule
Every eval should have:
- a category
- a target model scope
- a risk level
- a pass criterion
- a failure interpretation
- a recommended follow-up path

Do not create one-off evals with unclear meaning if they can fit an existing category.

---

## top-level eval families

### 1. scope evals
Used to verify lane discipline.

Examples:
- specialist stays within scope
- specialist answers only in-lane portions
- specialist escalates out-of-lane work
- parent delegates clear in-lane work

Typical trigger:
- scope failure
- specialist suitability review
- shadow-mode activation gate

---

### 2. routing evals
Used to verify parent routing and fallback behavior.

Examples:
- correct specialist chosen
- parent retains constitutional/governance work
- multi-specialist review triggered when needed
- fallback to parent occurs on uncertainty

Typical trigger:
- routing failure
- parent postmortem
- degraded-mode review

---

### 3. memory evals
Used to verify memory-first behavior and durable-memory handling.

Examples:
- specialist searches memory before answering
- relevant durable records are used
- `MEMORY.md` is not treated as the primary store
- promoted memory outranks weak/noisy records

Typical trigger:
- memory failure
- retrieval tuning
- summary projection review

---

### 4. contradiction evals
Used to verify conflict detection and staged resolution.

Examples:
- contradiction is staged instead of silently replacing active truth
- contradicted records are penalized in retrieval
- model surfaces unresolved conflict appropriately
- harness review is required before durable replacement

Typical trigger:
- contradiction handling failure
- memory corruption concern
- high-impact evidence conflict

---

### 5. tool-use evals
Used to verify tool choice, tool order, and tool validation.

Examples:
- correct tool selected
- required tool used when available
- tool result checked before relying on it
- unauthorized tool path not used
- fallback behavior works when a tool is unavailable

Typical trigger:
- tool failure
- runtime policy review
- high-impact task tooling review

---

### 6. reasoning evals
Used to verify synthesis quality when scope, memory, and tools are otherwise available.

Examples:
- fact vs inference separation
- strong tradeoff analysis
- consistent internal logic
- uncertainty correctly stated
- high-signal technical synthesis

Typical trigger:
- reasoning failure
- weak architectural recommendation
- repeated in-lane weak answers

---

### 7. validation evals
Used to verify that important outputs are checked before acceptance.

Examples:
- rollback path identified for durable change
- high-impact claims validated
- tests/checks run where appropriate
- confidence reduced when evidence is incomplete

Typical trigger:
- validation failure
- unsafe proposal
- weak engineering discipline

---

### 8. escalation evals
Used to verify handoff and escalation behavior.

Examples:
- specialist escalates constitutional issue
- specialist escalates cross-domain ambiguity
- parent is informed when required
- model stops instead of improvising beyond lane

Typical trigger:
- escalation failure
- scope drift
- governance-sensitive ambiguity

---

### 9. governance evals
Used to verify hard boundary compliance.

Examples:
- no constitutional self-edit attempt
- no self-grant of tool authority
- no bypass of required approval
- no markdown-only permission assumption
- parent remains subject to postmortem rules

Typical trigger:
- governance failure
- runtime enforcement review
- suspension review

---

### 10. postmortem evals
Used to verify postmortem triggering and postmortem quality.

Examples:
- meaningful failure triggers postmortem
- postmortem contains root cause
- postmortem identifies missed step
- postmortem proposes remediation
- skipped postmortem is treated as failure

Typical trigger:
- postmortem compliance failure
- shallow postmortems
- recurrence without learning

---

### 11. slot-integrity evals
Used to verify correct slot usage and projection discipline.

Examples:
- constitutional content stays in constitutional slots
- operational changes target the right slots
- `MEMORY.md` stays compact
- raw logs do not leak into summary slots
- projection output remains summary-like

Typical trigger:
- slot integrity failure
- projection review
- harness slot validation work

---

### 12. heartbeat/freshness evals
Used to verify the research/update loop.

Examples:
- specialist tracks approved sources only
- novelty threshold filters low-value noise
- heartbeat updates produce candidates rather than direct truth
- stale knowledge triggers appropriate refresh behavior

Typical trigger:
- heartbeat/freshness failure
- stale specialist guidance
- noisy update cycles

---

### 13. activation-gate evals
Used during specialist creation and shadow review.

Examples:
- scope adherence in shadow mode
- safe tool usage in shadow mode
- postmortem generation in shadow mode
- no scope expansion attempt
- no constitutional self-edit attempt

Typical trigger:
- new specialist activation
- shadow-mode review
- recovery-to-active review

---

### 14. degraded-recovery evals
Used when a specialist is in degraded mode or returning from suspension/shadow.

Examples:
- prior failure pattern no longer reproduces
- narrowed routing works safely
- postmortem recurrence falls
- targeted recovery action holds under test

Typical trigger:
- degraded mode
- recovery plan execution
- pre-reactivation review

---

### 15. specialist-suitability evals
Used when the role or lane itself may be wrong.

Examples:
- specialist performs well across its supposed core lane
- overlap with sibling specialists is acceptable
- lane is not too broad or incoherent
- specialist still justifies existence as a separate child

Typical trigger:
- specialist suitability failure
- repeated core-lane failures
- split/merge/retire review

---

## eval dimensions
Every eval should carry these dimensions.

### target scope
- parent
- single specialist
- multi-specialist
- harness
- runtime policy
- system-wide

### impact level
- low
- medium
- high

### test style
- unit
- scenario
- regression
- adversarial
- shadow
- governance
- smoke

### expected outcome type
- must pass
- should pass
- must fail safely
- must escalate
- must fallback
- must trigger postmortem

---

## canonical eval categories
Use these canonical names where possible.

### scope
- `scope_adherence`
- `out_of_lane_escalation`
- `parent_delegation_quality`

### routing
- `single_specialist_routing`
- `cross_domain_routing`
- `fallback_to_parent`
- `multi_specialist_review_trigger`

### memory
- `memory_first_behavior`
- `durable_record_usage`
- `summary_projection_discipline`
- `memory_retrieval_quality`

### contradiction
- `contradiction_staging`
- `conflict_sensitive_answering`
- `contradicted_record_penalty`

### tools
- `tool_selection_quality`
- `tool_ordering_quality`
- `tool_output_validation`
- `unauthorized_tool_blocking`

### reasoning
- `fact_vs_inference_separation`
- `tradeoff_analysis_quality`
- `uncertainty_expression`
- `in_lane_synthesis_quality`

### validation
- `durable_change_validation`
- `rollback_path_presence`
- `high_impact_claim_validation`

### escalation
- `constitutional_issue_escalation`
- `cross_domain_ambiguity_escalation`
- `stop_instead_of_improvise`

### governance
- `constitutional_self_edit_blocking`
- `authority_expansion_blocking`
- `approval_gate_compliance`

### postmortem
- `postmortem_triggering`
- `postmortem_quality`
- `postmortem_follow_through`

### slot integrity
- `slot_boundary_integrity`
- `memory_md_compactness`
- `projection_output_quality`

### heartbeat
- `approved_source_discipline`
- `novelty_threshold_quality`
- `heartbeat_candidate_flow`

### activation / recovery
- `shadow_activation_gate`
- `degraded_recovery_gate`
- `reactivation_gate`

### suitability
- `core_lane_competence`
- `lane_coherence`
- `specialist_overlap_acceptability`

---

## pass criteria patterns
Use one of these patterns.

### strict pass
Used for governance, constitutional, and safety boundaries.

Interpretation:
- any failure is serious
- usually triggers immediate review or containment

### threshold pass
Used for broader scenario suites.

Interpretation:
- suite has minimum pass rate
- repeated misses still matter even above threshold

### failure-safe pass
Used when the desired behavior is safe refusal, fallback, or escalation.

Interpretation:
- success means the system did not overreach

### recovery pass
Used after degraded-mode remediation.

Interpretation:
- success means the specific failure pattern no longer reproduces acceptably

---

## failure interpretation rules
Each eval family should have a default interpretation.

### scope/routing failures
Usually imply:
- routing review
- possible degraded mode if repeated
- possible parent postmortem

### memory/contradiction failures
Usually imply:
- retrieval or memory review
- contradiction review
- possible skill/prompt candidate

### tool/validation failures
Usually imply:
- `TOOLS.md` or `SKILLS.md` candidate
- runtime policy review if authority confusion exists

### governance failures
Usually imply:
- parent review
- possible immediate containment
- degraded or suspended state likely

### postmortem eval failures
Usually imply:
- meta-postmortem
- stronger harness enforcement

### suitability failures
Usually imply:
- parent review
- redesign, split, merge, or retirement discussion

---

## mapping from incident class to default evals

### scope failure
Default evals:
- `scope_adherence`
- `out_of_lane_escalation`

### routing failure
Default evals:
- `single_specialist_routing`
- `fallback_to_parent`
- `multi_specialist_review_trigger`

### memory failure
Default evals:
- `memory_first_behavior`
- `durable_record_usage`
- `memory_retrieval_quality`

### contradiction handling failure
Default evals:
- `contradiction_staging`
- `conflict_sensitive_answering`

### tool failure
Default evals:
- `tool_selection_quality`
- `tool_output_validation`

### reasoning failure
Default evals:
- `fact_vs_inference_separation`
- `tradeoff_analysis_quality`

### validation failure
Default evals:
- `durable_change_validation`
- `rollback_path_presence`

### escalation failure
Default evals:
- `constitutional_issue_escalation`
- `cross_domain_ambiguity_escalation`

### governance failure
Default evals:
- `constitutional_self_edit_blocking`
- `approval_gate_compliance`

### postmortem compliance failure
Default evals:
- `postmortem_triggering`
- `postmortem_quality`

### slot integrity failure
Default evals:
- `slot_boundary_integrity`
- `memory_md_compactness`

### heartbeat/freshness failure
Default evals:
- `approved_source_discipline`
- `novelty_threshold_quality`

### specialist suitability failure
Default evals:
- `core_lane_competence`
- `lane_coherence`
- `specialist_overlap_acceptability`

---

## eval creation rules
Create a new eval when:
- the failure was preventable
- the failure is likely to recur
- the incident reveals an untested boundary
- the incident affects routing, memory, tools, escalation, or governance
- the incident suggests activation criteria are too weak

Do not create redundant evals when an existing category already captures the issue.

---

## eval retirement and cleanup
An eval may be retired only when:
- it is duplicated by a stronger canonical eval
- the specialist or subsystem it targets has been retired
- the scenario is obsolete

Do not silently remove evals created from important postmortems without recording why.

---

## summary rule
Every eval should answer:
- what category of behavior is being tested?
- what model or layer is under test?
- what counts as success or safe failure?
- what should happen if it fails?

That shared structure is the point of this taxonomy.
