package postgres

import (
	"context"
	"database/sql"
	"freek/internal/ingestion/domain"
)

type CardRepo struct {
	db *sql.DB
}

func NewCardRepo(db *sql.DB) *CardRepo {
	return &CardRepo{db: db}
}

func (r *CardRepo) Create(ctx context.Context, card *domain.CreditCard) error {
	query := `INSERT INTO credit_cards (household_id, name, closing_day, member_id, credit_limit)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query, card.HouseholdID, card.Name, card.ClosingDay, card.MemberID, card.CreditLimit).
		Scan(&card.ID, &card.CreatedAt, &card.UpdatedAt)
}

func (r *CardRepo) GetByID(ctx context.Context, id int64) (*domain.CreditCard, error) {
	query := `SELECT id, household_id, name, closing_day, member_id, credit_limit, created_at, updated_at
		FROM credit_cards WHERE id = $1`
	c := &domain.CreditCard{}
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&c.ID, &c.HouseholdID, &c.Name, &c.ClosingDay, &c.MemberID, &c.CreditLimit, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *CardRepo) ListByHousehold(ctx context.Context, householdID int64) ([]domain.CreditCard, error) {
	query := `SELECT id, household_id, name, closing_day, member_id, credit_limit, created_at, updated_at
		FROM credit_cards WHERE household_id = $1 ORDER BY name`
	rows, err := r.db.QueryContext(ctx, query, householdID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []domain.CreditCard
	for rows.Next() {
		var c domain.CreditCard
		if err := rows.Scan(&c.ID, &c.HouseholdID, &c.Name, &c.ClosingDay, &c.MemberID, &c.CreditLimit, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cards, nil
}

func (r *CardRepo) ListByMember(ctx context.Context, householdID int64, memberID int64) ([]domain.CreditCard, error) {
	query := `SELECT id, household_id, name, closing_day, member_id, credit_limit, created_at, updated_at
		FROM credit_cards WHERE household_id = $1 AND member_id = $2 ORDER BY name`
	rows, err := r.db.QueryContext(ctx, query, householdID, memberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []domain.CreditCard
	for rows.Next() {
		var c domain.CreditCard
		if err := rows.Scan(&c.ID, &c.HouseholdID, &c.Name, &c.ClosingDay, &c.MemberID, &c.CreditLimit, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	return cards, nil
}

func (r *CardRepo) Update(ctx context.Context, card *domain.CreditCard) error {
	query := `UPDATE credit_cards SET name = $1, closing_day = $2, member_id = $3, credit_limit = $4, updated_at = NOW()
		WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, card.Name, card.ClosingDay, card.MemberID, card.CreditLimit, card.ID)
	return err
}

func (r *CardRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM credit_cards WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *CardRepo) HasExpenses(ctx context.Context, cardID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM expenses WHERE credit_card_id = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, cardID).Scan(&exists)
	return exists, err
}
