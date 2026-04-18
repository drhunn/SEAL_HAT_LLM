# OPS RUNBOOK

## steady-state rules
- The parent remains frozen.
- Specialists may improve only through bounded operational channels.
- Runtime policy controls actual tools.
- Durable memory promotion requires harness review.
- Constitutional changes require parent approval.
- Meaningful failures require postmortems.

## lifecycle states
- proposed
- approved
- bootstrapping
- inactive
- shadow
- active
- degraded
- suspended
- retired
- rejected

## adding a specialist
1. parent creates `config/specialist-instantiation.toml`
2. harness reviews
3. parent approves or rejects
4. template pack instantiated
5. constitutional slots created and locked
6. tool policy and memory namespace provisioned
7. specialist starts `inactive`
8. specialist enters `shadow`
9. shadow evals run
10. parent approves activation

## degraded mode
Use when repeated meaningful failures occur but recovery still looks possible.

In degraded mode:
- routing exposure narrows
- parent awareness increases
- high-impact tasks prefer parent involvement
- operational changes receive stricter review

## suspension
Suspend when:
- governance boundaries are threatened
- postmortem compliance repeatedly fails
- unsafe tool behavior repeats
- constitutional self-edit attempt occurs

## retirement
Retire when:
- the domain is obsolete
- specialist is redundant
- repeated failures show it should not continue
- parent explicitly decides to retire it

Retirement keeps history, audit, and postmortems unless explicitly archived.
