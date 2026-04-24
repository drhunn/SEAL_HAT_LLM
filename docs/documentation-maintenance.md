# Documentation maintenance policy

Documentation is part of the change, not cleanup after the fact.

## Rule

Every slice that changes runtime behavior, storage behavior, protocol behavior, workflow behavior, or public project structure must update the relevant documentation in the same slice.

Do not merge code that leaves the docs describing an older system.

## Required documentation updates by change type

### Runtime behavior

Update at least one of:
- `docs/implementation-status.md`
- feature-specific docs under `docs/features/`
- the affected runtime design note

### SQL, storage, or persistence behavior

Update at least one of:
- `docs/sql-contracts.md`
- `docs/schema-runtime-reconciliation.md`
- `docs/implementation-status.md`

### Task RPC, dispatch, or inter-unit behavior

Update at least one of:
- `docs/task-dispatch.md`
- `docs/implementation-status.md`
- the relevant feature docs under `docs/features/`

### Branch workflow, CI, or process behavior

Update at least one of:
- `docs/branch-workflow.md`
- `docs/feature-map.md`
- `.github/workflows/ci.yml` comments if the behavior is encoded there

### Feature completion criteria

Update:
- `docs/features/<feature>/acceptance-criteria.md`
- `docs/features/<feature>/test-plan.md`

## Slice checklist

Each slice must answer:

1. What code changed?
2. What user-visible, runtime, storage, or workflow behavior changed?
3. Which docs now need to change?
4. Which docs were intentionally not changed, and why?
5. What tests or verify commands prove the docs are still true?

## Admission rule

A slice is not complete if:
- the implementation changed but the status docs did not
- the SQL contract changed but `docs/sql-contracts.md` did not
- protocol behavior changed but protocol docs did not
- CI/workflow behavior changed but workflow docs did not
- acceptance criteria still contain placeholders for a feature being claimed as implemented

## Style rule

Update the smallest relevant documentation file.
Do not rewrite the entire docs tree to describe one change.
Do not hide behavior changes in vague architecture prose.
Use concrete statements about what exists, what is partial, and what is not implemented.
