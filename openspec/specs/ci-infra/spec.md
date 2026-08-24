# CI Infrastructure

## Purpose
Pipeline infrastructure guarantees: race detection across all test jobs,
reproducible digest-pinned service images, automated action updates, and
test coverage visibility.

## Requirements

### Requirement: Race detection in all test jobs
Every pipeline job that runs Go tests SHALL enable the race detector
(`go test -race`). A data race detected in either unit or integration tests
SHALL fail the corresponding job.

#### Scenario: Race detected in integration tests fails the job
- **WHEN** the integration job runs the test suite and the race detector reports a race
- **THEN** the integration job fails

#### Scenario: Clean code passes both jobs
- **WHEN** no races are present
- **THEN** unit-tests and integration jobs pass with `-race` enabled

### Requirement: Pinned service images
Service container images in the pipeline SHALL be pinned by immutable digest.
Floating tags alone SHALL NOT be used for services.

#### Scenario: Service uses digest-pinned image
- **WHEN** the integration job provisions PostgreSQL
- **THEN** the image reference includes a `@sha256:` digest

### Requirement: Automated action updates
The repository SHALL configure GitHub Dependabot for the `github-actions`
ecosystem with a weekly schedule and grouped updates, so that newer versions
of pinned actions arrive as reviewable pull requests.

#### Scenario: New action version available
- **WHEN** an action referenced in workflows publishes a newer version
- **THEN** Dependabot opens a grouped update pull request within its weekly cycle

### Requirement: Coverage visibility
The unit-tests job SHALL produce and print a total test coverage figure in the
job logs on every run. No coverage threshold SHALL gate merges as part of this
requirement.

#### Scenario: Coverage reported per run
- **WHEN** the unit-tests job completes
- **THEN** the job logs include the total statement coverage percentage
