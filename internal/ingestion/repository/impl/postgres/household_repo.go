package postgres

import (
	"context"
	"database/sql"
	"freek/internal/ingestion/domain"
)

// HouseholdRepo implements HouseholdRepository over PostgreSQL.
type HouseholdRepo struct {
	db *sql.DB
}

// NewHouseholdRepo creates a new HouseholdRepo.
func NewHouseholdRepo(db *sql.DB) *HouseholdRepo {
	return &HouseholdRepo{db: db}
}

// Create ...
func (r *HouseholdRepo) Create(ctx context.Context, household *domain.Household) error {
	query := `INSERT INTO households (name) VALUES ($1) RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query, household.Name).Scan(&household.ID, &household.CreatedAt, &household.UpdatedAt)
}

// GetByID ...
func (r *HouseholdRepo) GetByID(ctx context.Context, id int64) (*domain.Household, error) {
	query := `SELECT id, name, created_at, updated_at FROM households WHERE id = $1`
	h := &domain.Household{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&h.ID, &h.Name, &h.CreatedAt, &h.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return h, nil
}
