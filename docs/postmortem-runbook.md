# POSTMORTEM RUNBOOK

## core rule
A meaningful failure requires a postmortem.

This applies to:
- the parent
- all specialists
- routing failures
- failed self-improvement attempts
- skipped required postmortems

## workflow
1. detect incident
2. contain if needed
3. create structured postmortem
4. classify the incident
5. identify root cause layer
6. identify the missed step
7. mark preventability
8. choose remediation target
9. create derived artifacts
10. route for harness or parent approval
11. verify closure

## minimum required fields
- incident id
- task summary
- expected behavior
- actual behavior
- what went wrong
- primary failure class
- secondary classes
- severity
- preventability
- impact scope
- root cause
- missed evidence or step
- correction
- recommended remediation targets
- requires harness review
- requires parent review
- regression test needed

## possible outputs
- eval case
- memory artifact
- operational slot candidate
- routing review item
- degraded recommendation
- suspension recommendation
- parent review item

## meta-postmortem rule
If the system fails in how it handles a failure, create a second-order postmortem.
