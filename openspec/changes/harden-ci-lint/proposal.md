## Why

CI is green but the lint gate is permissive by design (errcheck disabled, no
formatters). A full lint sweep revealed 3 real correctness gaps
(`rows.Err` unchecked in all read-model queries — silently truncated results),
~90 style/robustness findings, and no policy governing linter exclusions.
Without a policy, every future `//nolint` (human or LLM-generated) becomes an
unreviewed escape hatch.

## What Changes

- Enable additional golangci-lint checks: `gofumpt` formatter, `revive`,
  `errcheck`, `noctx`, `rowserrcheck`, and `nolintlint` (requires a written
  reason on every `//nolint` directive).
- Fix all resulting findings:
  - check `rows.Err()` after iteration in all 3 read-repo queries;
  - adopt context-aware calls (`PingContext`, `NewRequestWithContext`);
  - add doc comments to exported symbols (Constitution #9);
  - format all files with gofumpt.
- Define surgical exclusions for accepted patterns only:
  - test-fixture seed calls in `_test.go` files (`errcheck`);
  - `json.NewEncoder(w).Encode(...)` after headers are sent (`errcheck`),
    per new Constitution principle #11.
- Document the exclusion policy for LLM and human contributors in AGENTS.md.
- Tighten `.golangci.yml` so CI fails on any new unexplained exclusion.

## Capabilities

### New Capabilities

- `lint-policy`: Enforced lint ruleset and exclusion policy: which linters run,
  what is excluded surgically and why, and mechanical enforcement of justified
  exclusions via `nolintlint`.

### Modified Capabilities

None.

## Impact

- `.golangci.yml`: expanded linters + commented exclude-rules.
- ~15 source/test files touched: formatting, doc comments, `rows.Err()` checks,
  context-aware calls. No behavior changes except earlier failure on truncated
  query iteration (a bug fix).
- `AGENTS.md`: new "Linter exclusions" section.
- CI workflow unchanged (same golangci-lint action; stricter config picked up
  from repo).
