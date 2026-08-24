# Architecture

Freek is a household expense-tracking API written in Go. It follows a
**layered, modular monolith** design with two bounded contexts — **ingestion**
(the write side) and **calculation** (the read side) — over a single
PostgreSQL database. The split resembles lightweight CQRS: writes go through
`internal/ingestion`, reads and analytics through `internal/calculation`.

```
                    ┌──────────────────────────────────────────┐
 HTTP (gorilla/mux) │            cmd/server/main.go            │  wiring + routes
                    └───────┬───────────────────────┬──────────┘
                            │                       │
        ┌───────────────────▼───────┐   ┌───────────▼───────────────────┐
        │      ingestion (write)    │   │     calculation (read)         │
        │  handler → service → repo │   │  handler → service → read repo │
        └───────────────┬───────────┘   └───────────┬───────────────────┘
                        │      shared (household scoping)
                        ▼                           ▼
                   PostgreSQL ◄──── pkg/money, pkg/dates (pure domain helpers)
```

## Tech stack

- **Go**, standard library `net/http`
- **gorilla/mux** for routing (`/api/v1/...` prefix)
- **PostgreSQL** via `database/sql` + `lib/pq`
- **golang-migrate** for schema migrations (`migrations/`, applied on startup)
- No external frameworks; manual dependency injection in `main.go`

## Module layout

```
cmd/server/main.go          Entry point: config, DB connect, migrations, DI, routes
internal/
  ingestion/                Write-side bounded context
    domain/                 Entities: Household, Member, Category, CreditCard, Expense
    repository/             Interfaces (ports) per aggregate
    repository/impl/postgres/
    service/                Business rules and validation (one service per entity)
    handler/                HTTP handlers: decode → validate → service → encode
  calculation/              Read-side bounded context (no writes)
    domain/                 Read models: ProjectionResult, SettlementPeriod, SummaryItem
    repository/             ExpenseReadRepository interface + filters/DTOs
    repository/impl/postgres/
    service/                ProjectionService, HistoryService, SummaryService
    handler/
  shared/                   Cross-context middleware (household scoping)
pkg/
  money/                    CLP integer type, interest/installment math
  dates/                    Closing dates, settlement cycles, custom periods
migrations/                 Numbered SQL migrations (up/down pairs)
docs/, openspec/            Documentation and spec-driven change workflow
```

Dependency rule: `handler → service → repository(interface)`. Implementations
live under `impl/postgres` and are injected in `main.go`. Domain packages never
import handlers or repositories.

## Bounded contexts

### Ingestion (write side)

Owns the persisted entities and all mutations:

- CRUD for households, members, cards, expenses and categories.
- Business validation lives in services, e.g.
  - credit expenses require a `credit_card_id`;
  - card `closing_day` must be 1–28;
  - predefined categories are immutable;
  - amounts must be positive integer CLP; installments coerced to ≥ 1.

### Calculation (read side)

Stateless queries over expenses; never mutates data:

- **ProjectionService**: simulates future months of outstanding credit
  installments per card or across all cards (default horizon 12 months).
- **HistoryService**: lists expenses grouped by settlement period (`YYYY-MM`).
- **SummaryService**: totals by category/member over date ranges and
  two-period comparison.

The read side has its own repository interface (`ExpenseReadRepository`) and
read models, so its SQL shape can evolve independently of write entities.

### Shared

`HouseholdContext` middleware extracts the household scope from the path
(`/households/{id}/...`). Every query in both contexts is implicitly filtered
by `household_id`; FKs use `ON DELETE CASCADE` so deleting a household removes
all its data.

## Data model

Tables (see `docs/GLOSARY.md` for full field reference):

| Table | Notes |
|---|---|
| `households` | Multi-tenancy root |
| `members` | Unique `(household_id, name)`; no roles |
| `categories` | Nullable `household_id` ⇒ global predefined category |
| `credit_cards` | `closing_day` CHECK 1–28 |
| `expenses` | `amount_clp BIGINT > 0`; enum `payment_type ('cash','debit','credit')`; `installments ≥ 1`; `interest_rate NUMERIC(5,4) ≥ 0` |

Delete semantics: children cascade with their household; optional references
(`category_id`, `member_id`, `credit_card_id`) use `ON DELETE SET NULL`.

Migrations run automatically at startup (`migrate.Up()`) and are also applied
by the test bootstrap script (`./up-db.sh` against `freek_test` database).

## Money and dates

Money handling and calendar math are extracted into pure packages with no I/O,
unit-tested independently:

- `pkg/money`: `CLP int64` (no decimals). `TotalWithInterest`,
  `DivideRoundUp`, `MonthlyInstallment` implement the installment math used by
  both contexts.
- `pkg/dates`: statement closing dates (`NextClosingDate`, `ClosingDates`),
  settlement cycle boundaries (`SettlementForMonth`), and configurable monthly
  periods that may span month boundaries (`PeriodConfig`).

## API surface

All routes live under `/api/v1` (see `cmd/server/main.go`):

- `POST /households`, `GET /households/{id}`
- Members: `GET|POST /households/{id}/members`, `DELETE .../members/{member_id}`
- Cards: `GET|POST /households/{id}/cards`, `GET|PUT|DELETE .../cards/{card_id}`
- Expenses: `GET|POST /households/{id}/expenses`, `PUT|DELETE .../expenses/{expense_id}`
- Categories: `GET|POST /households/{id}/categories`, `PUT|DELETE .../categories/{category_id}`
- Reads:
  - `GET /households/{id}/projections?months=&card_id=`
  - `GET /households/{id}/history?settlement=YYYY-MM`
  - `GET /households/{id}/history/summary?from=&to=&group_by=`
  - `GET /households/{id}/history/compare?from1=&to1=&from2=&to2=`
- `GET /health`

Handlers return JSON errors as `{"error": "..."}`; validation failures are 400,
unexpected failures 500.

## Testing

Table-driven unit tests per layer (`*_test.go` next to the code): services,
handlers and pure packages (`pkg/money`, `pkg/dates`). Integration tests for
Postgres repositories require the local database started with `./up-db.sh`;
run everything with `go test ./...`.
