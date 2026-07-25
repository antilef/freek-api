## ADDED Requirements

### Requirement: Create credit card
The system SHALL allow creating a credit card for a household with name, closing day, and optional member as titular.

#### Scenario: Successful creation
- **WHEN** user sends POST /api/v1/households/{id}/cards with name "CMR", closing_day 15, member_id for titular
- **THEN** system creates card and returns 201 with card data

#### Scenario: Card without titular (shared)
- **WHEN** user sends POST /api/v1/households/{id}/cards with name "Itaú", closing_day 10, no member_id
- **THEN** system creates card as shared and returns 201

#### Scenario: Invalid closing day
- **WHEN** user sends POST /api/v1/households/{id}/cards with closing_day outside 1-28
- **THEN** system returns 400 with error "closing_day must be between 1 and 28"

### Requirement: List credit cards
The system SHALL allow listing all credit cards of a household.

#### Scenario: Household has cards
- **WHEN** user sends GET /api/v1/households/{id}/cards
- **THEN** system returns 200 with array of cards including titular member name

#### Scenario: Filter by titular
- **WHEN** user sends GET /api/v1/households/{id}/cards?member_id=123
- **THEN** system returns only cards where member_id matches

### Requirement: Get credit card details
The system SHALL allow retrieving a specific credit card with its balance summary.

#### Scenario: Card exists
- **WHEN** user sends GET /api/v1/households/{id}/cards/{card_id}
- **THEN** system returns 200 with card data and summary of active installments

### Requirement: Update credit card
The system SHALL allow updating card name, closing day, or titular.

#### Scenario: Successful update
- **WHEN** user sends PUT /api/v1/households/{id}/cards/{card_id} with new closing_day
- **THEN** system updates card and returns 200 with updated data

### Requirement: Delete credit card
The system SHALL allow deleting a credit card if it has no associated expenses.

#### Scenario: No expenses
- **WHEN** user sends DELETE /api/v1/households/{id}/cards/{card_id} and card has no expenses
- **THEN** system deletes card and returns 204

#### Scenario: Has expenses
- **WHEN** user sends DELETE /api/v1/households/{id}/cards/{card_id} and card has expenses
- **THEN** system returns 409 with error "card has expenses, cannot delete"
