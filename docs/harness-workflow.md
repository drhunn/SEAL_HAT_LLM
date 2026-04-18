# HARNESS WORKFLOW

## purpose
Define the concrete operational workflow for how the harness processes incidents, postmortems, evals, slot candidates, memory candidates, health updates, and state transitions.

This is the main day-to-day review loop for the system.

---

## core responsibilities
The harness is responsible for:
- reviewing operational changes
- promoting or rejecting durable memory
- committing postmortems
- generating and reviewing eval candidates
- reviewing self-edit candidates
- updating specialist health signals
- recommending degraded, suspension, recovery, or activation transitions
- escalating constitutional and specialist-design issues to the parent

The harness is not the constitutional governor.
The parent remains the constitutional governor.

---

## main workflow
1. task executes or eval runs
2. signal is produced
3. harness classifies the signal
4. harness decides whether containment is needed
5. harness opens or validates postmortem if required
6. harness identifies derived artifacts
7. harness routes artifacts for review
8. harness updates health signals
9. harness recommends routing or lifecycle changes if needed
10. harness records closure state

---

## signal sources
Signals may come from:
- runtime task traces
- routing audit
- memory contradictions
- failed evals
- specialist self-report
- parent report
- human review
- shadow-mode review
- degraded-mode monitoring

---

## signal classes
The harness should first classify a signal as one of:
- normal success
- low-value noise
- postmortem-required incident
- candidate memory promotion
- candidate eval creation
- candidate operational slot update
- contradiction requiring staging or review
- routing review item
- lifecycle review item
- parent review item

---

## containment rules
Immediate containment is required when:
- constitutional self-edit attempt occurs
- tool authority boundary is violated
- repeated required postmortems are skipped
- severe routing failure affects high-impact tasks
- contradiction mishandling threatens durable truth
- specialist is clearly unsafe to continue routing to

Containment actions include:
- fallback to parent
- stop routing to specialist
- freeze candidate promotion
- force degraded state recommendation
- force suspension recommendation

---

## postmortem workflow
When a meaningful failure is detected:
1. confirm postmortem is required
2. create or validate structured postmortem
3. classify incident using `incident-classification.md`
4. identify root cause layer
5. identify missed step
6. mark preventability
7. identify remediation targets
8. create derived artifacts as appropriate

Reject shallow postmortems.
A shallow postmortem does not count as compliance.

---

## derived artifact workflow
A postmortem may create:
- durable memory candidate
- contradiction record
- eval case candidate
- operational slot candidate
- routing review item
- parent review item
- lifecycle review item

The harness should create only the smallest useful artifact set.
Do not generate noise artifacts just because a postmortem exists.

---

## durable memory review
When reviewing a memory candidate, the harness should ask:
- is this stable enough to become durable?
- is provenance sufficient?
- is it lane-appropriate?
- is this actually a contradiction case?
- should this stay staged instead of active?

Promote only when the evidence and usefulness justify it.

---

## eval workflow
When reviewing an eval candidate, the harness should ask:
- is this failure preventable or likely to recur?
- does an equivalent eval already exist?
- which canonical category from `eval-taxonomy.md` fits best?
- should this be regression, activation-gate, degraded-recovery, or governance coverage?

Use canonical eval names whenever possible.

---

## operational slot candidate workflow
When reviewing a candidate for:
- `TOOLS.md`
- `SKILLS.md`
- `PROMPT.md`
- `HEARTBEAT.md`
- `MEMORY.md`

The harness should ask:
- is this actually an operational issue?
- is the proposed change narrow?
- is it reversible?
- is it already covered elsewhere?
- should this be memory instead of slot guidance?
- does it need a paired eval?

Do not use operational slot edits to patch constitutional or specialist-design problems.

---

## contradiction workflow
When new evidence conflicts with active durable memory:
1. stage contradiction
2. do not auto-replace active record
3. review provenance and authority
4. decide whether to resolve, deprecate, or preserve both with conflict marker
5. update retrieval penalties or memory candidates if needed

---

## health update workflow
Health signals should consider:
- recent postmortem rate
- recent eval results
- recent routing failures
- contradiction load
- repeated skipped validation or escalation
- repeated governance pressure

After each meaningful incident cluster, refresh specialist health.

---

## lifecycle recommendation workflow
The harness may recommend:
- remain active
- remain in shadow
- activate from shadow
- degrade
- recover from degraded
- suspend
- retire review

### recommend degrade when
- repeated meaningful failures exist
- specialist still seems recoverable
- high-impact routing should narrow

### recommend suspend when
- governance-threatening behavior appears
- repeated severe failures continue
- routing is no longer safe

### recommend parent review when
- constitutional implications exist
- lane definition appears flawed
- specialist suitability is in doubt
- suspension or retirement is likely

---

## parent handoff contract
When escalating to the parent, include:
- incident summary
- primary and secondary failure classes
- severity and preventability
- root cause layer
- recommended action
- whether containment is already in place
- relevant postmortems, evals, and memory refs

---

## closure states
A harness review should end in one of:
- no action
- postmortem accepted
- postmortem rejected for rewrite
- memory candidate staged
- memory candidate promoted
- eval candidate created
- slot candidate staged
- routing review item created
- degraded recommended
- suspension recommended
- parent review opened
- issue closed with rationale

---

## anti-patterns
Do not:
- promote raw logs to durable truth
- accept shallow postmortems
- generate duplicate evals endlessly
- use slot edits to avoid parent governance review
- leave repeated incidents without health updates
- let contradictions silently overwrite durable memory

---

## summary rule
The harness exists to make adaptation controlled.

That means:
- incidents become postmortems
- postmortems become targeted follow-up actions
- durable changes are reviewed
- health signals are updated
- routing and lifecycle decisions are evidence-based
- constitutional matters go to the parent
