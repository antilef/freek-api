## Context

The CI pipeline (see archived `add-ci-workflow`) runs three jobs: lint,
unit-tests (`-race`), integration (PostgreSQL 18 service container + migrate
CLI). Known soft spots: race detector absent from the job that touches real
SQL concurrency paths; `postgres:18` tag is mutable; no automation keeps
pinned action SHAs current; coverage is invisible.

## Goals / Non-Goals

**Goals:**

- Same rigor in every test job: race detection everywhere tests run.
- Reproducible service images: same Postgres build on every run.
- Dependency updates arrive as reviewable PRs instead of silent drift.
- Coverage visible per run, establishing a baseline before any threshold.

**Non-Goals:**

- Gating on coverage percentage (future change once baseline exists).
- Branch protection rules (GitHub repo settings, documented separately).
- Renovate/other tooling — Dependabot suffices for this dependency surface.

## Decisions

### D1: `-race` in the integration job

Race conditions around DB pools/transactions are exactly what `-race`
catches. Cost is ~2x runtime of `go test`, acceptable at current suite size.
Unit and integration jobs become symmetric.

### D2: Pin postgres image by digest

Resolve `postgres:18` to its current digest and use
`image: postgres:18@sha256:<digest>`. Guarantees byte-identical service across
runs and protects against upstream tag mutation. Bumps become explicit.

*Alternative:* keep floating tag (reproducibility hole); pin full version
(e.g., 18.2) — still mutable, weaker guarantee.

### D3: Dependabot, weekly, grouped

`.github/dependabot.yml` with `package-ecosystem: github-actions`,
weekly schedule, grouped into one PR. Digest bumps arrive as ready-to-review
PRs with changelog links.

### D4: Coverage as report, not gate

Add `go test -cover ./...` output (and `go test -coverprofile` with a total
line in job logs) to unit-tests. No threshold: premature gates get gamed or
block progress. Baseline first; threshold is a follow-up change if wanted.

## Risks / Trade-offs

- [-race slows integration job] → minutes today; revisit only if it grows.
- [Digest pinning hides new Postgres point releases] → Dependabot-style PRs or
  manual periodic bump; explicitness is the goal.
- [Coverage number invites cargo-cult thresholds] → explicitly out of scope;
  any future threshold goes through a change.

## Migration Plan

Single commit touching workflow + new dependabot config. Rollback reverts.
No application impact.

## Open Questions

None.
