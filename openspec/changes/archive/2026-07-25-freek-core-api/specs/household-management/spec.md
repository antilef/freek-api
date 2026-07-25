## ADDED Requirements

### Requirement: Create household
The system SHALL allow creating a household with a name.

#### Scenario: Successful creation
- **WHEN** user sends POST /api/v1/households with name "Familia Pérez"
- **THEN** system creates household and returns 201 with household data including id and name

#### Scenario: Missing name
- **WHEN** user sends POST /api/v1/households without name
- **THEN** system returns 400 with error "name is required"

### Requirement: Get household
The system SHALL allow retrieving household details by id.

#### Scenario: Household exists
- **WHEN** user sends GET /api/v1/households/{id}
- **THEN** system returns 200 with household data including members count

#### Scenario: Household not found
- **WHEN** user sends GET /api/v1/households/{id} for non-existent id
- **THEN** system returns 404 with error "household not found"

### Requirement: Add member to household
The system SHALL allow adding members to a household. Members are labels for filtering expenses.

#### Scenario: Successful addition
- **WHEN** user sends POST /api/v1/households/{id}/members with name "Juan"
- **THEN** system creates member and returns 201 with member data

#### Scenario: Duplicate member name
- **WHEN** user sends POST /api/v1/households/{id}/members with name "Juan" when "Juan" already exists
- **THEN** system returns 409 with error "member already exists"

### Requirement: List members
The system SHALL allow listing all members of a household.

#### Scenario: Household has members
- **WHEN** user sends GET /api/v1/households/{id}/members
- **THEN** system returns 200 with array of members

#### Scenario: Household has no members
- **WHEN** user sends GET /api/v1/households/{id}/members for empty household
- **THEN** system returns 200 with empty array

### Requirement: Remove member
The system SHALL allow removing a member from a household.

#### Scenario: Member exists and has no expenses
- **WHEN** user sends DELETE /api/v1/households/{id}/members/{member_id}
- **THEN** system deletes member and returns 204

#### Scenario: Member has expenses
- **WHEN** user sends DELETE /api/v1/households/{id}/members/{member_id} and member has associated expenses
- **THEN** system returns 409 with error "member has expenses, cannot delete"
