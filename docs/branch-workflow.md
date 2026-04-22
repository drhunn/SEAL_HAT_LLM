# BRANCH WORKFLOW

## purpose

Define a simple staged branch flow for `SEAL_HAT_LLM` so work moves through design, feature specification, development, test, and release in a predictable way.

---

## branch model

Long-lived branches:

- `main`
- `Test`
- `Development`
- `Design/Architecture`

Feature branches:

- `feature/<feature-name>`

There is not one shared `feature` branch.
There is one branch per feature.

---

## flow

```text
Design/Architecture -> feature/<feature-name> -> Development -> Test -> main
```

---

## branch responsibilities

### `Design/Architecture`
This is the architecture and design branch.

Use it for:
- design changes
- architecture notes
- interface contracts
- schema proposals
- ADRs
- workflow changes
- feature planning inputs

This branch is where broader system direction is worked out before a feature branch is cut.

### `feature/<feature-name>`
This is a per-feature branch created for one specific feature.

Use it for:
- feature-specific specification
- feature-specific acceptance criteria
- minimal scaffolding
- skeleton files and initial wiring
- the tests required for the feature to be considered complete
- identifying dependencies and open design questions for that feature

A feature branch must produce:
- the feature spec
- the completion criteria
- the first scaffolding for the feature
- the tests that define done for that feature

It should stay narrow.
It is not the place to absorb unrelated work.

Recommended naming examples:
- `feature/research-specialist-template`
- `feature/branch-workflow-enforcement`
- `feature/seal-den-persistence`

### `Development`
This is the active integration branch.

Use it for:
- implementation work promoted out of a feature branch
- integration across packages
- iterative fixes during active development
- documentation updates tied to implemented behavior

This is the branch where the feature becomes real code rather than just a spec and scaffold.

### `Test`
This is the stabilization and verification branch.

Use it for:
- verification
- regression checking
- integration validation
- release hardening
- final bug fixes needed before promotion to `main`

Only work needed to get the candidate safely into `main` should land here.

### `main`
This is the stable branch.

Use it for:
- promoted, tested, review-ready work
- release-quality state

Do not treat `main` as a scratch branch.

---

## promotion rules

### `Design/Architecture` -> `feature/<feature-name>`
Cut a feature branch when:
- the feature is scoped
- the architecture impact is understood well enough to start
- the intended direction is documented

### `feature/<feature-name>` -> `Development`
Promote when:
- the feature spec exists
- the acceptance criteria exist
- the basic scaffolding exists
- the tests required for feature completion exist
- the work has a bounded target
- open questions are documented

### `Development` -> `Test`
Promote when:
- implementation is materially complete
- tests or validation steps exist for the change
- docs are updated enough to test the feature correctly
- known risks are documented

### `Test` -> `main`
Promote when:
- the change is validated
- regressions are addressed or accepted explicitly
- the branch is fit for stable use

---

## operating rules

- Prefer one feature branch per feature.
- Do not mix unrelated features in one feature branch.
- Send architectural churn back to `Design/Architecture` instead of hiding it in implementation branches.
- Keep `Development` as the integration surface, not the planning surface.
- Keep `Test` as the validation surface, not the place for broad new feature work.
- Keep `main` clean.

---

## pull request direction

Preferred PR direction:

- `Design/Architecture` -> `feature/<feature-name>` by branch cut, not PR
- `feature/<feature-name>` -> `Development`
- `Development` -> `Test`
- `Test` -> `main`

When a feature reveals a broader architectural issue, update `Design/Architecture` first or in parallel, then continue the feature branch from that clarified design state.

---

## notes

This workflow is intentionally simple.
It is meant to reduce branch sprawl while preserving a clear path from design to stable release.

The long-lived branches represent stages.
Feature branches represent individual units of work.
