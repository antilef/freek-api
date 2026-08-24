## Why

The pipeline is functional but its infrastructure has known soft spots: no race
detector on integration tests, a floating `postgres:18` image tag, no
automation for action updates, and no visibility into test coverage. Each one
is a slow-burn risk rather than an active bug.

## What Changes

- Run the integration job with `-race` (parity with unit-tests).
- Pin the PostgreSQL service image by digest instead of the floating `18` tag.
- Add GitHub Dependabot config for GitHub Actions version bumps
  (weekly, grouped).
- Add coverage reporting to unit tests (`go test -cover` summary in job logs;
  optional threshold later — not gating in this change).

## Capabilities

### New Capabilities

- `ci-infra`: Pipeline infrastructure hardening: race-detector parity across
  jobs, digest-pinned service images, automated dependency update automation,
  and coverage visibility.

### Modified Capabilities

None.

## Impact

- `.github/workflows/ci.yml`: `-race` flag, image digest, coverage step.
- `.github/dependabot.yml`: new file.
- No application code changes. Digest pinning introduces manual bump duty
  (mitigated by Dependabot for actions; Postgres bumps become explicit PRs).
