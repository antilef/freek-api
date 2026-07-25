## ADDED Requirements

### Requirement: Predefined categories exist
The system SHALL have a set of predefined categories available for all households.

#### Scenario: Default categories
- **WHEN** system starts
- **THEN** predefined categories exist: Alimentación, Transporte, Salud, Educación, Entretenimiento, Hogar, Ropa, Servicios, Otros

### Requirement: List categories
The system SHALL allow listing all categories for a household (predefined + custom).

#### Scenario: Household with only predefined
- **WHEN** user sends GET /api/v1/households/{id}/categories
- **THEN** system returns 200 with predefined categories

#### Scenario: Household with custom categories
- **WHEN** user sends GET /api/v1/households/{id}/categories and household has custom categories
- **THEN** system returns 200 with predefined + custom categories

### Requirement: Create custom category
The system SHALL allow creating custom categories per household.

#### Scenario: Successful creation
- **WHEN** user sends POST /api/v1/households/{id}/categories with name "Mascotas" and icon "paw"
- **THEN** system creates category and returns 201 with category data

#### Scenario: Duplicate name
- **WHEN** user sends POST /api/v1/households/{id}/categories with name "Alimentación" (predefined exists)
- **THEN** system returns 409 with error "category already exists"

### Requirement: Update category
The system SHALL allow updating custom categories. Predefined categories can only change icon.

#### Scenario: Update custom category
- **WHEN** user sends PUT /api/v1/households/{id}/categories/{category_id} with new name
- **THEN** system updates category and returns 200

#### Scenario: Update predefined category icon
- **WHEN** user sends PUT /api/v1/households/{id}/categories/{category_id} for predefined category with new icon
- **THEN** system updates icon and returns 200

#### Scenario: Update predefined category name
- **WHEN** user sends PUT /api/v1/households/{id}/categories/{category_id} for predefined category with new name
- **THEN** system returns 403 with error "cannot rename predefined category"

### Requirement: Delete custom category
The system SHALL allow deleting custom categories if no expenses use them.

#### Scenario: No expenses
- **WHEN** user sends DELETE /api/v1/households/{id}/categories/{category_id} and no expenses use it
- **THEN** system deletes category and returns 204

#### Scenario: Has expenses
- **WHEN** user sends DELETE /api/v1/households/{id}/categories/{category_id} and expenses use it
- **THEN** system returns 409 with error "category has expenses, cannot delete"

#### Scenario: Predefined category
- **WHEN** user sends DELETE /api/v1/households/{id}/categories/{category_id} for predefined category
- **THEN** system returns 403 with error "cannot delete predefined category"
