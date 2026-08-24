## Context

CI ships a deliberately permissive `.golangci.yml` (`govet`, `staticcheck`,
`unused`, `ineffassign`; `errcheck` disabled). A full sweep with candidate
linters shows: 8 files needing formatting, 3 `rows.Err` gaps, 6 `noctx`
findings, 32 `errcheck` findings (~25 in test fixtures, 4 handler
`json.Encode`, few misc), ~50 missing doc comments. The repo is LLM-assisted:
once `//nolint` patterns exist, models will replicate them without judgment.
Constitution #11 (new) codifies the exclusion policy; this change makes it
mechanical.

## Goals / Non-Goals

**Goals:**

- CI enforces formatting, error-checking hygiene, context-aware I/O, and doc
  comments — with zero unexplained exclusions.
- Every accepted pattern (test seeds, response encode) has one documented,
  surgical exclusion path.
- All existing findings fixed; pipeline stays green.

**Non-Goals:**

- Enabling additional linters not validated here (e.g., cyclop, funlen).
- Coverage thresholds or infra changes (see `harden-ci-infra`).
- Refactoring beyond what linters demand.

## Decisions

### D1: nolintlint is the enforcement backbone

Enable `nolintlint` with `require-explanation` and `require-specific`. Any
bare `//nolint` fails CI regardless of author (human or LLM). This converts
Constitution #11 from documentation into mechanism.

*Alternative:* rely on code review to catch bare directives — unreliable with
LLM-generated volume.

### D2: Exclusions live in two places only

1. **Config-level** (`exclude-rules` in `.golangci.yml`) for *path-scoped*
   patterns: `errcheck` excluded in `_test.go` files entirely (seed calls are
   setup noise; genuine test bugs surface at assertion time). One rule,
   commented, covers ~25 findings without thousands of directives.
2. **Inline** (`//nolint:errcheck // reason`) for the 4 handler
   `Encode` sites — they are specific call sites embodying Constitution #11's
   "response already committed" case, and should be visible where they happen.

*Alternative:* inline nolint on every test seed line (noise, 25+ directives);
global errcheck disable (status quo, hides real findings like `rows.Close`).

### D3: Formatters via golangci-lint v2 `fmt`

Use the built-in formatter gate (`gofumpt`) instead of a separate
`gofmt -l` step — single tool, single config file, same action in CI.

### D4: revive with defaults, no custom config

Default revive rules cover exported-symbol docs plus low-noise checks. No
custom rule tuning until defaults produce sustained false positives.

### D5: rows.Err fixes are bug fixes, not style

Each read-repo loop gains `if err := rows.Err(); err != nil { return nil, err }`
after the loop. Behavior change is intentional: truncated result sets now fail
loudly instead of returning partial data.

## Risks / Trade-offs

- [revive default set may flag new patterns later] → defaults are stable;
  additions go through a change per Constitution #11.
- [Large diff touches many files] → mechanical edits only; full test suite +
  lint must pass before commit.
- [nolintlint friction annoys contributors] → that is the point; explanation
  cost is one comment line.

## Migration Plan

Single PR: config + all fixes together so CI is never red between commits.
Rollback reverts to permissive config.

## Open Questions

None.
