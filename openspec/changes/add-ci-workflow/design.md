## Context

The repository is a Go 1.26 modular monolith (see `docs/ARCHITECTURE.md`) with
unit tests (`pkg/money`, `pkg/dates`, services, handlers) and PostgreSQL
integration tests in `internal/ingestion/repository/impl/postgres`. The
integration tests connect to a hardcoded DSN
(`postgres://localhost:5432/freek_test?sslmode=disable`) and **skip silently**
when the database is unreachable. Locally, `./up-db.sh` provisions a
PostgreSQL 18 container and applies migrations with `psql`. There is no lint
configuration and no CI today.

Constraints from `docs/CONSTITUTION.md`: no new Go dependencies without
explicit approval; idiomatic Go; tests must run green via `go test ./...`.

## Goals / Non-Goals

**Goals:**

- Every push to `main` and every PR is verified automatically.
- Integration tests run against a real PostgreSQL 18 instance whose schema is
  created by golang-migrate — the same mechanism production uses.
- A green pipeline guarantees integration tests actually ran (no silent skips).
- Lint gate exists from day one, initially permissive, ready to tighten.

**Non-Goals:**

- Tightening the linter ruleset (follow-up once CI is stable).
- Deployment, releases, artifact publishing.
- Test parallelization/sharding or coverage reporting.

## Decisions

### D1: PostgreSQL as a GitHub Actions service container

`services:` block with `postgres:18`, matching the local dev image. The
service maps port 5432 on the runner, so the hardcoded test DSN works without
any code change.

*Alternatives:* manual `docker run` (needs Docker-in-Docker setup for no gain);
SQLite or testcontainers (new dependency, violates Constitution #9).

### D2: Migrations applied with golang-migrate CLI (option B)

Install the pinned `migrate` CLI and run
`migrate -path migrations -database "$DSN" up`. This exercises the same tool
the server uses at startup (`cmd/server/main.go`), unlike `psql` which would
test a second migration path.

*Alternatives:* raw `psql -f` per file (works, but duplicates the mechanism);
a Go bootstrap step (overkill today).

### D3: Anti-skip guard at two layers

1. **Workflow**: a step pings the database (`pg_isready`) and fails the job if
   unreachable before tests start.
2. **Test code**: when env var `CI=1`, `setupTestDB` calls `t.Fatal` instead
   of `t.Skip` on connection errors. Belt and suspenders: even if the workflow
   guard is removed later, a DB outage can never produce a false green.

*Alternatives:* only the workflow guard (fragile — survives workflow edits);
parsing `go test -v` output for skips (brittle).

### D4: Three jobs, not one

`lint`, `unit-tests`, `integration` run in parallel for faster feedback and
independent red/green signals (e.g., a lint regression doesn't hide a test
failure). Cost is minor runner startup overhead, acceptable for this repo size.

### D5: Permissive golangci-lint initially

Committed `.golangci.yml` enabling low-noise linters (`govt`-style checks:
`govet`, `staticcheck`, `unused`, `ineffassign`) and explicitly disabling
`errcheck` — the codebase currently ignores errors like `db.Exec` cleanup in
tests and unchecked `json.NewEncoder` writes. Ruleset will be tightened in a
follow-up change.

*Alternatives:* strict config from day one (red CI immediately, blocks this
change's own verification); only `go vet` (no path to progressive tightening).

### D6: Toolchain pinning

- Go version read from `go.mod` via `actions/setup-go@v5`
  (`go-version-file`), so CI always matches the declared toolchain.
- Third-party actions referenced by immutable SHA, not floating tags.
- `concurrency` group cancels superseded runs on the same ref.

## Risks / Trade-offs

- [Postgres service slow to become ready] → health options / explicit
  `pg_isready` wait loop with timeout before migrating.
- [golangci-lint flags something unexpected] → config starts permissive;
  failures are fixable lint issues, not behavior changes.
- [`CI=1` heuristic couples tests to GitHub Actions naming] → it is the
  de-facto standard variable set by Actions; documented in test file.
- [Migration drift between migrate CLI and app startup] → both use the same
  library and path; risk accepted as negligible.

## Migration Plan

Additive only: new workflow + lint config + one test-file edit. Rollback is
reverting the commit. No data or API impact.

## Open Questions

None — decisions confirmed during exploration (option B for migrations,
three jobs, fail-on-missing-DB, permissive lint).
