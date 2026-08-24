package postgres

import (
	"context"
	"database/sql"
	"freek/internal/ingestion/domain"
)

// CategoryRepo ...
type CategoryRepo struct {
	db *sql.DB
}

// NewCategoryRepo creates a new CategoryRepo.
func NewCategoryRepo(db *sql.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

// ListByHousehold ...
func (r *CategoryRepo) ListByHousehold(ctx context.Context, householdID int64) ([]domain.Category, error) {
	query := `SELECT id, household_id, name, icon, is_predefined, created_at, updated_at
		FROM categories WHERE household_id = $1 OR household_id IS NULL
		ORDER BY is_predefined DESC, name`
	rows, err := r.db.QueryContext(ctx, query, householdID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.HouseholdID, &c.Name, &c.Icon, &c.IsPredefined, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

// Create ...
func (r *CategoryRepo) Create(ctx context.Context, category *domain.Category) error {
	query := `INSERT INTO categories (household_id, name, icon, is_predefined)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query, category.HouseholdID, category.Name, category.Icon, category.IsPredefined).
		Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)
}

// Update ...
func (r *CategoryRepo) Update(ctx context.Context, category *domain.Category) error {
	var query string
	var args []interface{}

	if category.IsPredefined {
		query = `UPDATE categories SET icon = $1, updated_at = NOW() WHERE id = $2`
		args = []interface{}{category.Icon, category.ID}
	} else {
		query = `UPDATE categories SET name = $1, icon = $2, updated_at = NOW() WHERE id = $3`
		args = []interface{}{category.Name, category.Icon, category.ID}
	}

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

// Delete ...
func (r *CategoryRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetByID ...
func (r *CategoryRepo) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	query := `SELECT id, household_id, name, icon, is_predefined, created_at, updated_at
		FROM categories WHERE id = $1`
	c := &domain.Category{}
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&c.ID, &c.HouseholdID, &c.Name, &c.Icon, &c.IsPredefined, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

// ExistsByName ...
func (r *CategoryRepo) ExistsByName(ctx context.Context, householdID int64, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM categories WHERE (household_id = $1 OR household_id IS NULL) AND name = $2)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, householdID, name).Scan(&exists)
	return exists, err
}

// HasExpenses ...
func (r *CategoryRepo) HasExpenses(ctx context.Context, categoryID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM expenses WHERE category_id = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, categoryID).Scan(&exists)
	return exists, err
}
