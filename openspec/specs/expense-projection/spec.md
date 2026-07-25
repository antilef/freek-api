# Expense Projection

## Purpose
Calculate and project credit card settlements for future months.

## Requirements

### Requirement: Project credit card settlements 12 months
The system SHALL calculate projected credit card settlements for the next 12 months.

#### Scenario: Single card with installments
- **WHEN** user sends GET /api/v1/households/{id}/projections?card_id=456&months=12
- **THEN** system returns monthly settlement amounts for each of the next 12 months for that card

#### Scenario: All cards projection
- **WHEN** user sends GET /api/v1/households/{id}/projections?months=12
- **THEN** system returns combined projection for all cards of the household

#### Scenario: Card with no active installments
- **WHEN** user sends GET /api/v1/households/{id}/projections?card_id=456
- **THEN** system returns 12 months with $0 amounts

### Requirement: Calculate monthly installment amount
The system SHALL correctly distribute installment amounts across billing cycles.

#### Scenario: Even installments without interest
- **WHEN** expense is 120000 CLP with 6 installments, 0% interest, card closing day 15, purchase date July 5
- **THEN** each month from July to December shows 20000 CLP

#### Scenario: Installments with interest
- **WHEN** expense is 100000 CLP with 3 installments, 15% interest, card closing day 15
- **THEN** total is 115000 CLP, each month shows approximately 38333 CLP

#### Scenario: Purchase after closing day
- **WHEN** purchase date is July 20 and card closing day is 15
- **THEN** first installment appears in August settlement (next closing)

### Requirement: Show projection summary
The system SHALL provide summary totals in projection response.

#### Scenario: Projection with data
- **WHEN** system returns projection for 12 months
- **THEN** response includes total per month and grand total of all months

### Requirement: Filter projection by member
The system SHALL allow filtering projection by member (titular).

#### Scenario: Filter by member
- **WHEN** user sends GET /api/v1/households/{id}/projections?member_id=123
- **THEN** system returns projection only for cards where member is titular

### Requirement: Handle multiple cards in projection
The system SHALL aggregate settlements from multiple cards into a household total.

#### Scenario: Two cards with different closing days
- **WHEN** household has CMR (closing 15) and Scotiabank (closing 20)
- **THEN** projection shows separate lines per card and a combined household total per month
