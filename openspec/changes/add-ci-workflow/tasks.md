## 1. Test hardening (anti-skip)

- [x] 1.1 In `internal/ingestion/repository/impl/postgres/repo_test.go`: when env `CI=1`, replace `t.Skip` with `t.Fatal` on DB open/ping failures in `setupTestDB`; keep skipping behavior locally

## 2. Lint configuration

- [x] 2.1 Add `.golangci.yml` with permissive ruleset: enable `govet`, `staticcheck`, `unused`, `ineffassign`; disable `errcheck`
- [x] 2.2 Run golangci-lint locally and fix any findings until clean

## 3. Workflow

- [x] 3.1 Create `.github/workflows/ci.yml` with triggers (`push` to `main`, `pull_request` to `main`) and a `concurrency` group that cancels superseded runs
- [x] 3.2 Add `lint` job: actions/checkout + actions/setup-go (`go-version-file: go.mod`) + golangci-lint action, all pinned by SHA
- [x] 3.3 Add `unit-tests` job: checkout + setup-go + `go test -race ./...`
- [x] 3.4 Add `integration` job with PostgreSQL 18 service container (`POSTGRES_DB=freek_test`, port 5432), `pg_isready` wait loop with timeout, migrations via pinned golang-migrate CLI against the test DSN, then `go test ./...`

## 4. Verification

- [x] 4.1 Run `go test ./...` locally with `CI=1` and database up → green; stop the database and confirm tests now fail (not skip) under `CI=1`, then restart it
- [ ] 4.2 Validate workflow syntax (e.g., `actionlint` or YAML lint) and push a branch to observe the pipeline go green on GitHub Actions
