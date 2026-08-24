## Why

The repository has no CI. Tests only run on developer machines, and the
integration tests silently skip (`t.Skip`) when PostgreSQL is unavailable —
so a green `go test ./...` locally does not prove the SQL layer works. We need
automated verification on every push and pull request, with a database that
is guaranteed to exist.

## What Changes

- Add a GitHub Actions workflow (`.github/workflows/ci.yml`) with three jobs:
  - **lint**: golangci-lint with an initially permissive configuration
    (`.golangci.yml` committed to the repo; tightened later).
  - **unit-tests**: `go test ./...` using the Go version declared in `go.mod`.
  - **integration**: same test suite against a PostgreSQL 18 service
    container, with migrations applied via golang-migrate CLI (the same
    mechanism the server uses at startup).
- The integration job must fail hard if the database is unreachable — silent
  skips are treated as failure.
- Make repository tests fail instead of skip when running in CI: when the
  `CI=1` environment variable is set, `setupTestDB` uses `t.Fatal` on DB
  connection errors instead of `t.Skip`.

## Capabilities

### New Capabilities

- `ci-pipeline`: Automated build verification on GitHub Actions: lint, unit
  tests, and integration tests with guaranteed PostgreSQL availability and
  no silent skips.

### Modified Capabilities

None. The `t.Fatal`-under-CI change in repository tests is test
infrastructure, not API-level behavior; its requirements live under
`ci-pipeline` ("No silent skips in CI").

## Impact

- New files: `.github/workflows/ci.yml`, `.golangci.yml`.
- Modified: `internal/ingestion/repository/impl/postgres/repo_test.go`
  (skip → fatal under CI).
- No application code, API, or schema changes.
- Dependencies added at CI level only (golangci-lint action, golang-migrate
  CLI); no new Go module dependencies.
