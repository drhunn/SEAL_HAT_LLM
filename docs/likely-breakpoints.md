# LIKELY BREAKPOINTS

## purpose
Record the areas most likely to break as `SEAL_HAT_LLM` evolves.

---

## 1. Go runtime ↔ SQL function drift
The Go runtime already expects more from the SQL layer than the earliest scaffold provided.
Typical breakpoints:
- missing helper functions
- changed function signatures
- changed return columns
- changed enum names or status values

### mitigation
- update `docs/sql-contracts.md` whenever SQL functions change
- run `cmd/verify` after SQL changes
- keep retrieval return shapes synchronized

---

## 2. Go runtime ↔ SQL schema drift
Likely failure modes:
- Go writes columns that do not exist yet
- Go expects tables that were never created or renamed
- lifecycle writes assume richer schema than the DB currently has

### mitigation
- treat schema/runtime reconciliation as a tracked task
- prefer explicit migrations over silent schema edits

---

## 3. rename drift after repo rename
The repository has been renamed to `SEAL_HAT_LLM`, but name drift can still occur across:
- docs
- Go module path
- package comments
- old commit messages and preserved markdown code references

### mitigation
- keep repo/module/doc naming explicit
- allow Python package naming to stay stable only if that is intentional and documented

---

## 4. architecture docs ahead of implementation
The docs are detailed and useful, but they are broader than the live runtime.

### mitigation
- use `docs/implementation-status.md` as the reality check
- do not assume a documented workflow is already enforced in code

---

## 5. preserved `.md` code versus live code confusion
Some code is preserved as `.md` due to in-session connector limitations.

### mitigation
- always leave a note in the folder
- state whether the `.md` file is archival or canonical
- keep one executable entrypoint only

---

## 6. retrieval scoring contract drift
The retrieval wrapper and SQL function must agree on result shape and score semantics.

### mitigation
- keep retrieval smoke tests in `cmd/verify`
- document return columns in `docs/sql-contracts.md`

---

## 7. health and lifecycle semantics drift
Health score logic, degraded rules, and lifecycle transitions may drift between docs, Go code, and SQL helpers.

### mitigation
- make health update formulas explicit
- define allowed transitions clearly
- audit lifecycle writes

---

## 8. training/data format proliferation
The Python HAT layer can emit multiple formats.
That is useful, but it increases the risk of ambiguity about which format is canonical.

### mitigation
- keep `build_dataset.py` as the canonical dataset-building entrypoint
- document output purposes clearly

---

## 9. CI false confidence
Basic CI is helpful, but passing import/build checks does not mean the runtime is truly integrated.

### mitigation
- treat CI as a floor, not a proof of completeness
- add DB-backed smoke coverage later

---

## 10. first-specialist scope creep
The first specialist is correctly focused on CS/software engineering and tool development, but it could gradually become too broad.

### mitigation
- keep lane boundaries explicit
- do not use the first specialist as a dumping ground for everything technical

---

## summary
If the repo breaks, the most likely causes are:
- schema/runtime mismatch
- function contract drift
- naming drift
- docs outrunning implementation
- preserved-code ambiguity
