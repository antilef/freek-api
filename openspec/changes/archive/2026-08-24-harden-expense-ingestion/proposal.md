## Why

The ingestion side of the API was archived as complete, but an audit found four gaps between the specs and the implementation: derived interest values are computed inconsistently across the codebase (SQL truncation vs Go round-up), `Update` bypasses all validation and returns misleading HTTP status codes, soft-delete semantics are ambiguous (a `status` column exists but nothing uses it), and omitting `purchase_date` stores year-1 dates instead of defaulting to today.

## What Changes

- **Interest doctrine: derive at read time**: ingestion continues to store raw facts (`amount_clp`, `installments`, `interest_rate`); the spec now states this explicitly. Derived values (`total_with_interest`, monthly installment) are computed at read time.
- **Consistent derivation formula**: all read-side calculations SHALL use round-up division (`pkg/money.DivideRoundUp`), replacing the SQL truncation in the calculation repo.
- **Update validation**: `PUT /expenses/{id}` applies the same validations as create and returns `400` with a clear error message on invalid input (instead of a raw database error as `500`), and `404` for missing expenses.
- **Delete is hard delete**: explicitly specified that deletion permanently removes the expense; soft-delete is not supported in this version. The unused `status` column remains but no behavior claims rest on it.
- **purchase_date defaulting**: when omitted, `purchase_date` defaults to the current date instead of Go's zero time.

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `expense-ingestion`: spec corrected to state raw-fact storage at ingestion (no persisted calculated fields); update endpoint gains validation requirements with proper 400/404 status codes; delete is explicitly hard delete; purchase_date defaults to current date when omitted; derivation formula is standardized to round-up division on the read side.

## Impact

- `internal/ingestion/service/expense_service.go`: shared Create/Update validation, purchase_date defaulting
- `internal/ingestion/handler/expense_handler.go`: status code mapping
- `internal/calculation/repository/impl/postgres/expense_read_repo.go`: replace SQL truncation with round-up derivation consistent with `pkg/money`
- No schema changes; no migrations
