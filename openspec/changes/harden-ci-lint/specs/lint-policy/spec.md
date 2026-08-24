## ADDED Requirements

### Requirement: Enforced linters
The lint gate SHALL run `gofumpt` (formatter), `govet`, `staticcheck`, `unused`,
`ineffassign`, `errcheck`, `noctx`, `rowserrcheck`, `revive`, and `nolintlint`.
The codebase SHALL pass all of them.

#### Scenario: Clean repository passes
- **WHEN** the lint job runs on a commit with no findings
- **THEN** the lint job passes

#### Scenario: Unchecked rows.Err fails lint
- **WHEN** code iterates `database/sql` rows without checking `rows.Err()`
- **THEN** lint fails with a `rowserrcheck` finding

### Requirement: Truncated query results fail loudly
Read-model repository queries SHALL check `rows.Err()` after iteration and
propagate the error. A mid-iteration failure SHALL NOT be reported as a
successful partial result.

#### Scenario: Read repo propagates iteration errors
- **WHEN** any read-repo method finishes scanning rows
- **THEN** it returns any accumulated iteration error instead of partial data

### Requirement: Justified exclusions only
Every linter exclusion SHALL be either (a) a path-scoped rule in
`.golangci.yml` with a written rationale comment, or (b) an inline
`//nolint:<linter> // <reason>` directive. Bare `//nolint` directives SHALL
fail lint via `nolintlint`.

#### Scenario: Bare nolint is rejected
- **WHEN** a commit adds `//nolint` or `//nolint:errcheck` without a reason comment
- **THEN** the lint job fails

#### Scenario: Documented exclusions are accepted
- **WHEN** handler response encoding uses `//nolint:errcheck // response already committed`
- **THEN** lint passes for that line

### Requirement: Surgical exclusion scope
Exclusions SHALL cover exactly two accepted patterns: error-returning calls in
`_test.go` fixture/seed code (`errcheck`), and JSON response encoding after
headers are committed (`errcheck`). No other patterns are excluded; disabling
a linter globally to clear findings SHALL NOT be possible through config alone
without an explicit policy change.

#### Scenario: New unchecked error outside exclusions fails
- **WHEN** production code ignores an error return outside the two accepted patterns
- **THEN** lint fails with the corresponding finding

### Requirement: Exclusion policy documented for contributors
AGENTS.md SHALL document the exclusion policy for LLM-assisted development:
exclusions are last resort, always carry reasons, and new recurring patterns
require an OpenSpec change to amend `.golangci.yml`.

#### Scenario: Policy discoverable
- **WHEN** an LLM or contributor reads AGENTS.md
- **THEN** it finds explicit instructions on when and how to write justified exclusions
