## 1. Workflow

- [x] 1.1 Add `-race` to `go test` in the integration job
- [x] 1.2 Resolve current `postgres:18` digest (`docker pull` + `inspect` or registry API) and pin service image as `postgres:18@sha256:<digest>` with a comment on how to bump

## 2. Dependabot

- [x] 2.1 Create `.github/dependabot.yml`: `github-actions` ecosystem, weekly schedule, grouped updates, assignees/reviewers left to repo defaults

## 3. Coverage

- [x] 3.1 In unit-tests job: run `go test -race ./... -coverprofile=coverage.out` and add a step printing total coverage (`go tool cover -func=coverage.out | tail -1`)

## 4. Verification

- [x] 4.1 `actionlint` clean; push and confirm all three jobs green with `-race` in integration and coverage line visible in unit-tests logs
- [x] 4.2 Confirm Dependabot config is picked up (repo Settings → Dependabot, or first grouped PR within the cycle)
