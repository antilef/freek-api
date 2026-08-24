package postgres

import (
	"context"
	"database/sql"
	"freek/internal/ingestion/domain"
)

// MemberRepo implements MemberRepository over PostgreSQL.
type MemberRepo struct {
	db *sql.DB
}

// NewMemberRepo creates a new MemberRepo.
func NewMemberRepo(db *sql.DB) *MemberRepo {
	return &MemberRepo{db: db}
}

// Create inserts a member and populates generated fields.
func (r *MemberRepo) Create(ctx context.Context, member *domain.Member) error {
	query := `INSERT INTO members (household_id, name) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query, member.HouseholdID, member.Name).Scan(&member.ID, &member.CreatedAt, &member.UpdatedAt)
}

// GetByID returns a member by ID within the household scope.
func (r *MemberRepo) GetByID(ctx context.Context, id int64) (*domain.Member, error) {
	query := `SELECT id, household_id, name, created_at, updated_at FROM members WHERE id = $1`
	m := &domain.Member{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&m.ID, &m.HouseholdID, &m.Name, &m.CreatedAt, &m.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

// ListByHousehold returns all members of a household ordered by name.
func (r *MemberRepo) ListByHousehold(ctx context.Context, householdID int64) ([]domain.Member, error) {
	query := `SELECT id, household_id, name, created_at, updated_at FROM members WHERE household_id = $1 ORDER BY name`
	rows, err := r.db.QueryContext(ctx, query, householdID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var members []domain.Member
	for rows.Next() {
		var m domain.Member
		if err := rows.Scan(&m.ID, &m.HouseholdID, &m.Name, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return members, nil
}

// Delete removes a member row.
func (r *MemberRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM members WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// ExistsByName reports whether a household already has a member with that name.
func (r *MemberRepo) ExistsByName(ctx context.Context, householdID int64, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM members WHERE household_id = $1 AND name = $2)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, householdID, name).Scan(&exists)
	return exists, err
}

// HasExpenses reports whether the member is referenced by any expense.
func (r *MemberRepo) HasExpenses(ctx context.Context, memberID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM expenses WHERE member_id = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, memberID).Scan(&exists)
	return exists, err
}
