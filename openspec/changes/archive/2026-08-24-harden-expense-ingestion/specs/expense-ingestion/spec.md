## MODIFIED Requirements

### Requirement: Create expense
The system SHALL allow creating an expense with amount, description, category, payment type, and optional details. The system SHALL store the raw amounts as provided (`amount_clp`, `installments`, `interest_rate`); derived values (total with interest, monthly installment amount) SHALL be calculated at read time, not persisted at ingestion.

#### Scenario: Cash/debit expense
- **WHEN** user sends POST /api/v1/households/{id}/expenses with amount 50000, description "Supermercado", payment_type "cash", category_id
- **THEN** system creates expense and returns 201 with expense data

#### Scenario: Credit card expense (single payment)
- **WHEN** user sends POST /api/v1/households/{id}/expenses with amount 80000, payment_type "credit", credit_card_id, installments 1
- **THEN** system creates expense with installments=1 and returns 201

#### Scenario: Credit card expense with installments
- **WHEN** user sends POST /api/v1/households/{id}/expenses with amount 120000, payment_type "credit", credit_card_id, installments 6, interest_rate 0
- **THEN** system creates expense storing amount_clp=120000, installments=6, interest_rate=0 and returns 201

#### Scenario: Credit card expense with interest
- **WHEN** user sends POST /api/v1/households/{id}/expenses with amount 100000, payment_type "credit", credit_card_id, installments 3, interest_rate 0.15
- **THEN** system stores amount_clp=100000, installments=3, interest_rate=0.15 as-is (no derived values persisted) and returns 201

#### Scenario: Missing required fields
- **WHEN** user sends POST /api/v1/households/{id}/expenses without amount or description
- **THEN** system returns 400 with validation errors

#### Scenario: Credit expense without card
- **WHEN** user sends POST /api/v1/households/{id}/expenses with payment_type "credit" but no credit_card_id
- **THEN** system returns 400 with error "credit_card_id is required for credit expenses"

### Requirement: Update expense
The system SHALL allow updating expense details, applying the same validations as creation, and returning appropriate HTTP status codes.

#### Scenario: Successful update
- **WHEN** user sends PUT /api/v1/households/{id}/expenses/{expense_id} with new amount
- **THEN** system updates expense and returns 200 with updated data

#### Scenario: Update with invalid amount
- **WHEN** user sends PUT /api/v1/households/{id}/expenses/{expense_id} with amount 0 or negative
- **THEN** system returns 400 with error "amount must be greater than 0" (no database error surfaces)

#### Scenario: Update credit expense without card
- **WHEN** user sends PUT /api/v1/households/{id}/expenses/{expense_id} changing payment_type to "credit" without credit_card_id
- **THEN** system returns 400 with error "credit_card_id is required for credit expenses"

#### Scenario: Update nonexistent expense
- **WHEN** user sends PUT /api/v1/households/{id}/expenses/{expense_id} for an expense that does not exist
- **THEN** system returns 404

## ADDED Requirements

### Requirement: Derived amounts are consistent
All read-side calculations of total-with-interest and monthly installment amounts SHALL use a single formula: total = amount * (1 + interest_rate), monthly installment = round-up division of total by installments (`pkg/money.DivideRoundUp`). No read path SHALL truncate.

#### Scenario: Monthly amount rounding
- **WHEN** any read endpoint derives monthly amounts for an expense of total 115000 across 3 installments
- **THEN** each installment is reported as 38334 (rounded up)

#### Scenario: Consistency between read paths
- **WHEN** history, summary, and projection endpoints derive values for the same expense
- **THEN** all report identical monthly amounts

### Requirement: Delete is permanent
The system SHALL delete expenses permanently. Soft-delete is not supported in this version; deleted expenses SHALL NOT be recoverable through the API.

#### Scenario: Successful deletion
- **WHEN** user sends DELETE /api/v1/households/{id}/expenses/{expense_id}
- **THEN** system removes the expense row permanently and returns 204

#### Scenario: Deleted expense is gone from all queries
- **WHEN** an expense has been deleted
- **THEN** it does not appear in list, history, summary, or projection results

#### Scenario: Delete nonexistent expense
- **WHEN** user sends DELETE /api/v1/households/{id}/expenses/{expense_id} for an expense that does not exist
- **THEN** system returns 404

### Requirement: Purchase date defaults to today
The system SHALL default purchase_date to the current date when the client omits it.

#### Scenario: Create without purchase_date
- **WHEN** user sends POST /api/v1/households/{id}/expenses without purchase_date
- **THEN** system stores the current date and returns 201 with purchase_date set to today

#### Scenario: Create with explicit purchase_date
- **WHEN** user sends POST /api/v1/households/{id}/expenses with purchase_date "2026-07-14"
- **THEN** system stores exactly that date
