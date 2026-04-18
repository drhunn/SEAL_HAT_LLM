# KNOWN GAPS

## purpose
Track explicit known gaps so the repository does not look more complete than it is.

---

## resolved gap
The earlier standalone Go verification-command gap has been addressed.
The repository now includes:
- `cmd/verify/main.go`
- `make verify`

That closes the earlier note about having no dedicated verification entrypoint.

---

## current meaningful gaps

### 1. schema/runtime reconciliation is still an active concern
The Go runtime and SQL layer are both present, but they must continue to be actively reconciled as the scaffold evolves.

### 2. end-to-end compile/integration confidence is still limited
The repository has stronger structure now, but it is still not the same thing as a fully proven end-to-end runtime.

### 3. tool broker and live model execution are still scaffold-level
The architecture supports them, but they are not yet fully implemented.

### 4. eval execution is still lighter than the eval architecture
Eval generation concepts are strong; full execution and gating are still less mature.

### 5. slot DB/filesystem synchronization is still basic
Filesystem snapshotting exists, but true bidirectional synchronization and conflict handling are not yet mature.

### 6. training workflows are scaffolded, not fully benchmarked
The Python HAT layer can prepare strong datasets and training inputs, but production-quality benchmarked training workflows remain future work.

---

## companion docs
Use these documents together:
- `docs/implementation-status.md`
- `docs/likely-breakpoints.md`
- `docs/schema-runtime-reconciliation.md`
- `docs/sql-contracts.md`

---

## summary
The repo has moved beyond a pure idea stage.
Its biggest remaining risks are no longer missing concepts, but keeping the implementation layers aligned and verifying that the intended control loop really works end to end.
