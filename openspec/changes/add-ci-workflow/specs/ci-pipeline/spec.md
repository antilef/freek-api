## ADDED Requirements

### Requirement: CI runs on push and pull requests
The project SHALL run an automated CI pipeline on GitHub Actions for every
push to `main` and every pull request targeting `main`. Superseded runs on the
same reference SHALL be cancelled.

#### Scenario: Push to main triggers pipeline
- **WHEN** a commit is pushed to `main`
- **THEN** the CI workflow starts and executes lint, unit test, and integration test jobs

#### Scenario: Pull request triggers pipeline
- **WHEN** a pull request targeting `main` is opened or updated
- **THEN** the CI workflow starts and its result gates the PR checks

### Requirement: Lint job
The pipeline SHALL run golangci-lint using a configuration file committed to
the repository, with the Go toolchain version taken from `go.mod`. The job
SHALL fail on any linter finding.

#### Scenario: Lint passes on clean code
- **WHEN** the codebase has no findings for the enabled linters
- **THEN** the lint job passes

#### Scenario: Lint failure blocks pipeline
- **WHEN** golangci-lint reports at least one finding
- **THEN** the lint job fails and the pipeline is marked red

### Requirement: Unit tests job
The pipeline SHALL compile and run the full Go test suite with the race
detector enabled, without requiring external services.

#### Scenario: Unit tests pass
- **WHEN** all tests pass (including tests that skip due to no database)
- **THEN** the unit test job passes

#### Scenario: Failing test blocks pipeline
- **WHEN** any test fails
- **THEN** the unit test job fails and the pipeline is marked red

### Requirement: Integration tests run against real PostgreSQL
The pipeline SHALL provision a PostgreSQL 18 service container reachable at
`localhost:5432` with database `freek_test`, apply all migrations from
`migrations/` using golang-migrate CLI, wait for database readiness before
migrating, and then run the full test suite. If the database or a migration
fails, the job SHALL fail.

#### Scenario: Database ready and migrated
- **WHEN** the integration job starts
- **THEN** it waits until PostgreSQL accepts connections, applies all up-migrations successfully, and only then runs tests

#### Scenario: Database unreachable fails the job
- **WHEN** PostgreSQL never becomes ready within the configured timeout
- **THEN** the integration job fails instead of running skipped-green tests

#### Scenario: Failed migration fails the job
- **WHEN** any migration script returns an error
- **THEN** the integration job fails before running tests

### Requirement: No silent skips in CI
Repository integration tests SHALL NOT skip due to database unavailability
when running in a CI environment. When environment variable `CI` is set to
`1`, a failed database connection in test setup SHALL abort the test with a
fatal error.

#### Scenario: Test setup hard-fails under CI
- **WHEN** the repository test suite runs with `CI=1` and PostgreSQL is not reachable
- **THEN** the affected test aborts with a fatal setup error and the suite fails

#### Scenario: Local development keeps skipping
- **WHEN** the repository test suite runs locally without PostgreSQL and without `CI=1`
- **THEN** database-dependent tests are skipped as today
