## ADDED Requirements

### Requirement: Create expense
The system SHALL allow creating an expense with amount, description, category, payment type, and optional details.

#### Scenario: Cash/debit expense
- **WHEN** user sends POST /api/v1/households/{id}/expenses with amount 50000, description "Supermercado", payment_type "cash", category_id
- **THEN** system creates expense and returns 201 with expense data

#### Scenario: Credit card expense (single payment)
- **WHEN** user sends POST /api/v1/households/{id}/expenses with amount 80000, payment_type "credit", credit_card_id, installments 1
- **THEN** system creates expense with installments=1 and returns 201

#### Scenario: Credit card expense with installments
- **WHEN** user sends POST /api/v1/households/{id}/expenses with amount 120000, payment_type "credit", credit_card_id, installments 6, interest_rate 0
- **THEN** system creates expense with 6 installments of 20000 each and returns 201

#### Scenario: Credit card expense with interest
- **WHEN** user sends POST /api/v1/households/{id}/expenses with amount 100000, payment_type "credit", installments 3, interest_rate 0.15
- **THEN** system calculates total 115000, creates expense with installments of approximately 38333 each and returns 201

#### Scenario: Missing required fields
- **WHEN** user sends POST /api/v1/households/{id}/expenses without amount or description
- **THEN** system returns 400 with validation errors

#### Scenario: Credit expense without card
- **WHEN** user sends POST /api/v1/households/{id}/expenses with payment_type "credit" but no credit_card_id
- **THEN** system returns 400 with error "credit_card_id is required for credit expenses"

### Requirement: Add notes to expense
The system SHALL allow adding free-text notes to any expense.

#### Scenario: Create with notes
- **WHEN** user sends POST /api/v1/households/{id}/expenses with notes "Compré un vino para la junta de amigos"
- **THEN** system saves notes and returns them in expense data

#### Scenario: Update notes
- **WHEN** user sends PUT /api/v1/households/{id}/expenses/{expense_id} with new notes
- **THEN** system updates notes and returns 200 with updated expense

### Requirement: Assign member to expense
The system SHALL allow assigning a member (titular/executor) to an expense for filtering.

#### Scenario: Expense with member
- **WHEN** user sends POST /api/v1/households/{id}/expenses with member_id
- **THEN** system associates expense with member and returns 201

#### Scenario: Expense without member
- **WHEN** user sends POST /api/v1/households/{id}/expenses without member_id
- **THEN** system creates expense with null member and returns 201

### Requirement: List expenses
The system SHALL allow listing expenses with optional filters.

#### Scenario: List all expenses
- **WHEN** user sends GET /api/v1/households/{id}/expenses
- **THEN** system returns 200 with all expenses ordered by purchase_date descending

#### Scenario: Filter by member
- **WHEN** user sends GET /api/v1/households/{id}/expenses?member_id=123
- **THEN** system returns only expenses where member_id matches

#### Scenario: Filter by card
- **WHEN** user sends GET /api/v1/households/{id}/expenses?card_id=456
- **THEN** system returns only expenses where credit_card_id matches

#### Scenario: Filter by category
- **WHEN** user sends GET /api/v1/households/{id}/expenses?category_id=789
- **THEN** system returns only expenses where category_id matches

#### Scenario: Filter by date range
- **WHEN** user sends GET /api/v1/households/{id}/expenses?from=2026-07-14&to=2026-07-19
- **THEN** system returns expenses where purchase_date is between from and to (inclusive)

#### Scenario: Filter by last N days
- **WHEN** user sends GET /api/v1/households/{id}/expenses?last=7d
- **THEN** system returns expenses from the last 7 days

#### Scenario: Combined filters
- **WHEN** user sends GET /api/v1/households/{id}/expenses?member_id=123&card_id=456&from=2026-07-01
- **THEN** system returns expenses matching all filters

### Requirement: Update expense
The system SHALL allow updating expense details.

#### Scenario: Successful update
- **WHEN** user sends PUT /api/v1/households/{id}/expenses/{expense_id} with new amount
- **THEN** system updates expense and returns 200 with updated data

### Requirement: Delete expense
The system SHALL allow deleting an expense.

#### Scenario: Successful deletion
- **WHEN** user sends DELETE /api/v1/households/{id}/expenses/{expense_id}
- **THEN** system deletes expense and returns 204
