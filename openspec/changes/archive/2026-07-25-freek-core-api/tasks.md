## 1. Project Setup

- [x] 1.1 Initialize Go module (`go mod init freek`)
- [x] 1.2 Add dependencies: gorilla/mux, lib/pq or pgx, golang-migrate
- [x] 1.3 Create directory structure: cmd/server, internal/{ingestion,calculation,shared}, pkg
- [x] 1.4 Create cmd/server/main.go entry point with basic HTTP server

## 2. Database & Migrations

- [x] 2.1 Create migration 001_households.sql (households table)
- [x] 2.2 Create migration 002_members.sql (members table with household_id FK)
- [x] 2.3 Create migration 003_categories.sql (categories table with predefined seed data)
- [x] 2.4 Create migration 004_credit_cards.sql (credit_cards table with closing_day, member_id FK)
- [x] 2.5 Create migration 005_expenses.sql (expenses table with all fields: amount_clp, payment_type, installments, interest_rate, notes, member_id FK, credit_card_id FK, category_id FK)
- [x] 2.6 Configure migration runner in main.go

## 3. Shared Utilities

- [x] 3.1 Create pkg/money/money.go with CLP type (int64), String() with dots formatting, parsing
- [x] 3.2 Create pkg/dates/periods.go with period calculation (start_day to end_day across months)
- [x] 3.3 Create pkg/dates/closing.go with next closing date calculation (purchase_date → first settlement)

## 4. Domain Entities

- [x] 4.1 Create internal/ingestion/domain/household.go (Household struct)
- [x] 4.2 Create internal/ingestion/domain/member.go (Member struct)
- [x] 4.3 Create internal/ingestion/domain/credit_card.go (CreditCard struct with closing_day)
- [x] 4.4 Create internal/ingestion/domain/expense.go (Expense struct with payment_type, installments, interest_rate, notes)
- [x] 4.5 Create internal/ingestion/domain/category.go (Category struct with is_predefined flag)
- [x] 4.6 Create internal/calculation/domain/projection.go (ProjectionMonth, ProjectionResult structs)
- [x] 4.7 Create internal/calculation/domain/settlement.go (SettlementPeriod struct)

## 5. Repositories - Ingestion

- [x] 5.1 Create internal/ingestion/repository/household_repo.go (interface)
- [x] 5.2 Create internal/ingestion/repository/member_repo.go (interface)
- [x] 5.3 Create internal/ingestion/repository/card_repo.go (interface)
- [x] 5.4 Create internal/ingestion/repository/expense_repo.go (interface)
- [x] 5.5 Create internal/ingestion/repository/category_repo.go (interface)
- [x] 5.6 Implement internal/ingestion/repository/impl/postgres/household_repo.go
- [x] 5.7 Implement internal/ingestion/repository/impl/postgres/member_repo.go
- [x] 5.8 Implement internal/ingestion/repository/impl/postgres/card_repo.go
- [x] 5.9 Implement internal/ingestion/repository/impl/postgres/expense_repo.go
- [x] 5.10 Implement internal/ingestion/repository/impl/postgres/category_repo.go

## 6. Services - Ingestion

- [x] 6.1 Create internal/ingestion/service/household_service.go (Create, Get, AddMember, ListMembers, RemoveMember)
- [x] 6.2 Create internal/ingestion/service/card_service.go (Create, List, Get, Update, Delete with closing_day validation)
- [x] 6.3 Create internal/ingestion/service/expense_service.go (Create with installment calculation, Update, Delete, List with filters)
- [x] 6.4 Create internal/ingestion/service/category_service.go (List, Create custom, Update, Delete with predefined protection)

## 7. Handlers - Ingestion

- [x] 7.1 Create internal/ingestion/handler/household_handler.go (POST /households, GET /households/{id}, POST /households/{id}/members, GET /households/{id}/members, DELETE /households/{id}/members/{member_id})
- [x] 7.2 Create internal/ingestion/handler/card_handler.go (POST /households/{id}/cards, GET /households/{id}/cards, GET /households/{id}/cards/{card_id}, PUT /households/{id}/cards/{card_id}, DELETE /households/{id}/cards/{card_id})
- [x] 7.3 Create internal/ingestion/handler/expense_handler.go (POST /households/{id}/expenses, GET /households/{id}/expenses with filters, PUT /households/{id}/expenses/{expense_id}, DELETE /households/{id}/expenses/{expense_id})
- [x] 7.4 Create internal/ingestion/handler/category_handler.go (GET /households/{id}/categories, POST /households/{id}/categories, PUT /households/{id}/categories/{category_id}, DELETE /households/{id}/categories/{category_id})

## 8. Repositories - Calculation

- [x] 8.1 Create internal/calculation/repository/expense_read_repo.go (interface for read-only queries)
- [x] 8.2 Implement internal/calculation/repository/impl/postgres/expense_read_repo.go (queries for projection, history, summaries)

## 9. Services - Calculation

- [x] 9.1 Create internal/calculation/service/projection_service.go (ProjectByCard, ProjectAll, calculateMonthlyInstallment with closing_day logic)
- [x] 9.2 Create internal/calculation/service/history_service.go (QueryByPeriod, QueryBySettlement, SearchByText)
- [x] 9.3 Create internal/calculation/service/summary_service.go (SummaryByCategory, SummaryByMember, ComparePeriods)

## 10. Handlers - Calculation

- [x] 10.1 Create internal/calculation/handler/projection_handler.go (GET /households/{id}/projections with card_id, member_id, months params)
- [x] 10.2 Create internal/calculation/handler/history_handler.go (GET /households/{id}/history with from, to, last, period, card_id, settlement, q params)
- [x] 10.3 Create internal/calculation/handler/summary_handler.go (GET /households/{id}/history/summary with group_by param, GET /households/{id}/history/compare)

## 11. API Router & Integration

- [x] 11.1 Create internal/shared/household.go (HouseholdContext middleware to extract household_id)
- [x] 11.2 Wire up Gorilla/mux router with all routes grouped by domain (ingestion vs calculation)
- [x] 11.3 Add JSON request/response encoding middleware
- [x] 11.4 Add error handling middleware (validation errors, not found, conflict)
- [x] 11.5 Connect database, repositories, services, and handlers in main.go
- [x] 11.6 Test full flow: create household → add member → create card → create expense → query projection
