# KNOWN GAPS

## in-session GitHub push limitation note
During this build-out, one intended file was not added successfully from the session:

- a standalone Go verification command (`cmd/verify` / `cmd/check` style entrypoint)

The repository was updated around that gap by:
- keeping the main Go harness runtime scaffold in place
- retaining the rest of the Go service wiring
- changing the `Makefile` to use `go test ./...` for the current check/test path instead of relying on the missing standalone verification command

## practical impact
Right now, the repository has:
- the Go harness entrypoint
- Postgres wiring
- slot loading
- memory/postmortem services
- routing/eval/recovery scaffolding

But it does **not yet** have a dedicated standalone Go verification CLI committed in-repo.

## recommended follow-up
Add one of the following in a later pass:
- `cmd/verify/main.go`
- `cmd/check/main.go`

That command should perform at least:
- config load check
- database connectivity check
- slot file load check
- health snapshot check
- optional retrieval smoke test

## status
This is a repository note so the gap is explicit rather than implicit.
