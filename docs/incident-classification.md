# INCIDENT CLASSIFICATION

## dimensions
Every meaningful incident should classify:
- primary failure class
- secondary contributing classes
- severity
- preventability
- impact scope
- required response level
- remediation targets

## primary failure classes
- scope failure
- routing failure
- memory failure
- tool failure
- reasoning failure
- validation failure
- escalation failure
- governance failure
- postmortem compliance failure
- contradiction handling failure
- eval failure
- slot integrity failure
- specialist suitability failure

## severity
- 1 = low
- 2 = moderate
- 3 = high
- 4 = critical

## preventability
- preventable
- partially preventable
- not clearly preventable

## impact scope
- local
- specialist-local
- routing-local
- memory-local
- governance-wide
- system-wide

## response levels
- level 0 = log only
- level 1 = postmortem required
- level 2 = harness review required
- level 3 = parent review required
- level 4 = immediate containment

## key rule
Classification is not paperwork.
It determines:
- whether the harness must act
- whether the parent must intervene
- whether routing should narrow
- whether evals should be added
- whether a specialist should degrade, suspend, or retire
