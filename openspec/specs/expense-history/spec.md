# Expense History

## Purpose
Query, filter, and analyze past expenses.

## Requirements

### Requirement: Query expenses by flexible period
The system SHALL allow querying expenses by custom date ranges matching the user's payment period.

#### Scenario: Custom period
- **WHEN** user sends GET /api/v1/households/{id}/history?from=2026-07-05&to=2026-08-04
- **THEN** system returns all expenses within that date range

#### Scenario: Last N days
- **WHEN** user sends GET /api/v1/households/{id}/history?last=30d
- **THEN** system returns expenses from the last 30 days

#### Scenario: Current month by household period
- **WHEN** user sends GET /api/v1/households/{id}/history?period=current
- **THEN** system returns expenses within the household's configured period (start_day to end_day)

### Requirement: Query expenses by card settlement period
The system SHALL allow querying expenses that belong to a specific credit card settlement.

#### Scenario: Card settlement for specific month
- **WHEN** user sends GET /api/v1/households/{id}/history?card_id=456&settlement=2026-08
- **THEN** system returns expenses that will appear in the August settlement of that card

#### Scenario: Card settlement with closing day logic
- **WHEN** card closing day is 15 and user queries settlement 2026-08
- **THEN** system returns expenses purchased between July 16 and August 15

### Requirement: Expense summary by category
The system SHALL provide expense summaries grouped by category.

#### Scenario: Summary for period
- **WHEN** user sends GET /api/v1/households/{id}/history/summary?from=2026-07-01&to=2026-07-31
- **THEN** system returns total per category and grand total

### Requirement: Expense summary by member
The system SHALL provide expense summaries grouped by member (titular).

#### Scenario: Summary by member
- **WHEN** user sends GET /api/v1/households/{id}/history/summary?from=2026-07-01&to=2026-07-31&group_by=member
- **THEN** system returns total per member and grand total

### Requirement: Compare periods
The system SHALL allow comparing expenses between two periods.

#### Scenario: Compare two months
- **WHEN** user sends GET /api/v1/households/{id}/history/compare?period1_from=2026-06-01&period1_to=2026-06-30&period2_from=2026-07-01&period2_to=2026-07-31
- **THEN** system returns totals for each period and the difference

### Requirement: Search expenses by description
The system SHALL allow searching expenses by text in description or notes.

#### Scenario: Text search
- **WHEN** user sends GET /api/v1/households/{id}/history?q=vino
- **THEN** system returns expenses where description or notes contain "vino" (case-insensitive)
