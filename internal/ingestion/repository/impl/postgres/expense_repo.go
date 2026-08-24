package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"freek/internal/ingestion/domain"
	"freek/internal/ingestion/repository"
	"time"
)

type ExpenseRepo struct {
	db *sql.DB
}

func NewExpenseRepo(db *sql.DB) *ExpenseRepo {
	return &ExpenseRepo{db: db}
}

func (r *ExpenseRepo) Create(ctx context.Context, expense *domain.Expense) error {
	query := `INSERT INTO expenses (household_id, amount_clp, description, notes, category_id, payment_type,
		credit_card_id, member_id, purchase_date, installments, interest_rate, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		expense.HouseholdID, expense.AmountCLP, expense.Description, expense.Notes,
		expense.CategoryID, expense.PaymentType, expense.CreditCardID, expense.MemberID,
		expense.PurchaseDate, expense.Installments, expense.InterestRate, expense.Status).
		Scan(&expense.ID, &expense.CreatedAt, &expense.UpdatedAt)
}

func (r *ExpenseRepo) GetByID(ctx context.Context, id int64) (*domain.Expense, error) {
	query := `SELECT id, household_id, amount_clp, description, notes, category_id, payment_type,
		credit_card_id, member_id, purchase_date, installments, interest_rate, status, created_at, updated_at
		FROM expenses WHERE id = $1`
	e := &domain.Expense{}
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&e.ID, &e.HouseholdID, &e.AmountCLP, &e.Description, &e.Notes, &e.CategoryID,
			&e.PaymentType, &e.CreditCardID, &e.MemberID, &e.PurchaseDate,
			&e.Installments, &e.InterestRate, &e.Status, &e.CreatedAt, &e.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *ExpenseRepo) List(ctx context.Context, filter repository.ExpenseFilter) ([]domain.Expense, error) {
	query := `SELECT id, household_id, amount_clp, description, notes, category_id, payment_type,
		credit_card_id, member_id, purchase_date, installments, interest_rate, status, created_at, updated_at
		FROM expenses WHERE household_id = $1`
	args := []interface{}{filter.HouseholdID}
	argIdx := 2

	if filter.CreditCardID != nil {
		query += fmt.Sprintf(" AND credit_card_id = $%d", argIdx)
		args = append(args, *filter.CreditCardID)
		argIdx++
	}
	if filter.MemberID != nil {
		query += fmt.Sprintf(" AND member_id = $%d", argIdx)
		args = append(args, *filter.MemberID)
		argIdx++
	}
	if filter.CategoryID != nil {
		query += fmt.Sprintf(" AND category_id = $%d", argIdx)
		args = append(args, *filter.CategoryID)
		argIdx++
	}
	if filter.From != nil {
		query += fmt.Sprintf(" AND purchase_date >= $%d", argIdx)
		args = append(args, *filter.From)
		argIdx++
	}
	if filter.To != nil {
		query += fmt.Sprintf(" AND purchase_date <= $%d", argIdx)
		args = append(args, *filter.To)
		argIdx++
	}
	if filter.LastDays != nil {
		query += fmt.Sprintf(" AND purchase_date >= $%d", argIdx)
		args = append(args, time.Now().AddDate(0, 0, -*filter.LastDays))
		argIdx++
	}
	if filter.Search != nil {
		query += fmt.Sprintf(" AND (description ILIKE $%d OR notes ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+*filter.Search+"%")
	}

	query += " ORDER BY purchase_date DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []domain.Expense
	for rows.Next() {
		var e domain.Expense
		if err := rows.Scan(&e.ID, &e.HouseholdID, &e.AmountCLP, &e.Description, &e.Notes, &e.CategoryID,
			&e.PaymentType, &e.CreditCardID, &e.MemberID, &e.PurchaseDate,
			&e.Installments, &e.InterestRate, &e.Status, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}
	return expenses, nil
}

func (r *ExpenseRepo) Update(ctx context.Context, expense *domain.Expense) error {
	query := `UPDATE expenses SET amount_clp = $1, description = $2, notes = $3, category_id = $4,
		payment_type = $5, credit_card_id = $6, member_id = $7, purchase_date = $8,
		installments = $9, interest_rate = $10, status = $11, updated_at = NOW()
		WHERE id = $12`
	_, err := r.db.ExecContext(ctx, query,
		expense.AmountCLP, expense.Description, expense.Notes, expense.CategoryID,
		expense.PaymentType, expense.CreditCardID, expense.MemberID, expense.PurchaseDate,
		expense.Installments, expense.InterestRate, expense.Status, expense.ID)
	return err
}

func (r *ExpenseRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM expenses WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
