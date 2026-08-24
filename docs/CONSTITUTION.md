# Constitution

Non-negotiable principles for the freek codebase. Every contribution — human or
LLM-generated — must comply with all of them. If a change violates one of these
principles, the design is wrong, not the principle.

## 1. Household is the tenancy boundary
Every table, query, repository method and endpoint is scoped by `household_id`.
A repository method without a household scope does not exist. Never return data
across households, even in internal helpers.

## 2. Layered flow, always
`handler → service → repository (interface)`. Handlers only decode/encode HTTP
and map errors to status codes. Business rules live exclusively in services.
Repositories know nothing about HTTP. Domain packages never import handlers or
repositories.

## 3. Dependencies point inward via interfaces
Services depend on repository **interfaces** declared in their own context;
Postgres implementations live under `impl/postgres` and are wired by hand in
`cmd/server/main.go`. No global state, no service locators, no init-time magic.

## 4. Ingestion writes, calculation reads
Mutations belong to `internal/ingestion`; queries, projections and summaries to
`internal/calculation`. Calculation code never mutates persisted state, and the
two contexts do not import each other's internals.

## 5. Money is integer CLP, dates are explicit
All amounts are `int64` CLP with no decimals; use `pkg/money` for formatting,
interest and installment math — never raw float arithmetic on money. Calendar
logic (closing days, settlements, periods) lives only in `pkg/dates`; never
hand-roll month math in services or SQL.

## 6. Validate at the service boundary
Services are the guardians of invariants: positive amounts, required credit
card for credit expenses, closing day 1–28, installments ≥ 1, immutable
predefined categories. A handler that bypasses validation, or a repository that
silently "fixes" bad input, is a bug.

## 7. Explicit SQL, no ORMs or query builders
Queries are hand-written SQL inside Postgres implementations. Schema changes
require a numbered migration pair (`NNN_name.up.sql` / `.down.sql`) in
`migrations/`; never alter a schema from application code. Deletes follow the
established semantics: cascade with the household, `SET NULL` for optional
references.

## 8. Errors are values, wrapped with context
Functions return `error`; no panics across package boundaries, no ignored
errors (`_`). Errors carry enough context to identify the failing operation.
Handlers translate them once into JSON `{"error": ...}` with 400 for
validation and 500 for unexpected failures.

## 9. Idiomatic Go, minimal surface
Follow Effective Go: exported symbols need doc comments, accept `context.Context`
first where I/O happens, prefer small interfaces defined at the consumer,
table-driven tests next to the code. No comments restating the obvious, no
premature abstractions, no dependencies beyond the existing stack
(stdlib, gorilla/mux, lib/pq, golang-migrate) without explicit approval.

## 10. Specs before code
Behavioral changes start as an OpenSpec change (proposal → design → tasks) and
are implemented task by task. Code without tests, or tests requiring a database
without `./up-db.sh`, do not merge. Run `go test ./...` and keep it green.
