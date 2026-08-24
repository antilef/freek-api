// Package postgres implements calculation read repositories over PostgreSQL.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"freek/internal/calculation/domain"
	calcRepo "freek/internal/calculation/repository"
	"freek/pkg/money"
	"time"
)

// ExpenseReadRepo ...
type ExpenseReadRepo struct {
	db *sql.DB
}

// NewExpenseReadRepo ...
func NewExpenseReadRepo(db *sql.DB) *ExpenseReadRepo {
	return &ExpenseReadRepo{db: db}
}

// ListForProjection lists active credit expenses for projection.
func (r *ExpenseReadRepo) ListForProjection(ctx context.Context, householdID int64, cardID *int64) ([]domain.SettlementExpense, error) {
	query := `SELECT id, description, amount_clp, installments, interest_rate, purchase_date
		FROM expenses WHERE household_id = $1 AND payment_type = 'credit' AND status = 'active'`
	args := []interface{}{householdID}
	argIdx := 2

	if cardID != nil {
		query += fmt.Sprintf(" AND credit_card_id = $%d", argIdx)
		args = append(args, *cardID)
	}

	query += " ORDER BY purchase_date"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var expenses []domain.SettlementExpense
	for rows.Next() {
		var e domain.SettlementExpense
		var amountCLP int64
		var installments int
		var interestRate float64
		var purchaseDate time.Time
		if err := rows.Scan(&e.ExpenseID, &e.Description, &amountCLP, &installments, &interestRate, &purchaseDate); err != nil {
			return nil, err
		}
		total := money.TotalWithInterest(amountCLP, interestRate)
		e.MonthlyAmount = int64(money.MonthlyInstallment(total, installments))
		e.Amount = int64(total)
		expenses = append(expenses, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return expenses, nil
}

// ListByPeriod lists expenses matching the read filter.
func (r *ExpenseReadRepo) ListByPeriod(ctx context.Context, filter calcRepo.ExpenseReadFilter) ([]domain.SettlementExpense, error) {
	query := `SELECT e.id, e.description, e.amount_clp, e.payment_type, e.member_id, COALESCE(m.name, '') as member_name
		FROM expenses e
		LEFT JOIN members m ON e.member_id = m.id
		WHERE e.household_id = $1`
	args := []interface{}{filter.HouseholdID}
	argIdx := 2

	if filter.CreditCardID != nil {
		query += fmt.Sprintf(" AND e.credit_card_id = $%d", argIdx)
		args = append(args, *filter.CreditCardID)
		argIdx++
	}
	if filter.MemberID != nil {
		query += fmt.Sprintf(" AND e.member_id = $%d", argIdx)
		args = append(args, *filter.MemberID)
		argIdx++
	}
	if filter.From != nil {
		query += fmt.Sprintf(" AND e.purchase_date >= $%d", argIdx)
		args = append(args, *filter.From)
		argIdx++
	}
	if filter.To != nil {
		query += fmt.Sprintf(" AND e.purchase_date <= $%d", argIdx)
		args = append(args, *filter.To)
		argIdx++
	}
	if filter.Search != nil {
		query += fmt.Sprintf(" AND (e.description ILIKE $%d OR e.notes ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+*filter.Search+"%")
	}

	query += " ORDER BY e.purchase_date DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var expenses []domain.SettlementExpense
	for rows.Next() {
		var e domain.SettlementExpense
		var paymentType string
		var memberID sql.NullInt64
		var memberName string
		if err := rows.Scan(&e.ExpenseID, &e.Description, &e.Amount, &paymentType, &memberID, &memberName); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return expenses, nil
}

// SummaryByCategory totals expenses per category.
func (r *ExpenseReadRepo) SummaryByCategory(ctx context.Context, filter calcRepo.ExpenseReadFilter) ([]calcRepo.SummaryItem, error) {
	query := `SELECT e.category_id, COALESCE(c.name, 'Sin categoría') as category_name, SUM(e.amount_clp) as total
		FROM expenses e
		LEFT JOIN categories c ON e.category_id = c.id
		WHERE e.household_id = $1`
	args := []interface{}{filter.HouseholdID}
	argIdx := 2

	if filter.From != nil {
		query += fmt.Sprintf(" AND e.purchase_date >= $%d", argIdx)
		args = append(args, *filter.From)
		argIdx++
	}
	if filter.To != nil {
		query += fmt.Sprintf(" AND e.purchase_date <= $%d", argIdx)
		args = append(args, *filter.To)
	}

	query += " GROUP BY e.category_id, c.name ORDER BY total DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var items []calcRepo.SummaryItem
	for rows.Next() {
		var item calcRepo.SummaryItem
		if err := rows.Scan(&item.CategoryID, &item.CategoryName, &item.Total); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// SummaryByMember totals expenses per member.
func (r *ExpenseReadRepo) SummaryByMember(ctx context.Context, filter calcRepo.ExpenseReadFilter) ([]calcRepo.SummaryItem, error) {
	query := `SELECT e.member_id, COALESCE(m.name, 'Sin titular') as member_name, SUM(e.amount_clp) as total
		FROM expenses e
		LEFT JOIN members m ON e.member_id = m.id
		WHERE e.household_id = $1`
	args := []interface{}{filter.HouseholdID}
	argIdx := 2

	if filter.From != nil {
		query += fmt.Sprintf(" AND e.purchase_date >= $%d", argIdx)
		args = append(args, *filter.From)
		argIdx++
	}
	if filter.To != nil {
		query += fmt.Sprintf(" AND e.purchase_date <= $%d", argIdx)
		args = append(args, *filter.To)
	}

	query += " GROUP BY e.member_id, m.name ORDER BY total DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var items []calcRepo.SummaryItem
	for rows.Next() {
		var item calcRepo.SummaryItem
		if err := rows.Scan(&item.MemberID, &item.MemberName, &item.Total); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
