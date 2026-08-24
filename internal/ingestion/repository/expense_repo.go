package repository

import (
	"context"
	"freek/internal/ingestion/domain"
	"time"
)

// ExpenseFilter selects which household expenses List returns.
type ExpenseFilter struct {
	HouseholdID  int64
	CreditCardID *int64
	MemberID     *int64
	CategoryID   *int64
	From         *time.Time
	To           *time.Time
	LastDays     *int
	Search       *string
}

// ExpenseRepository persists and queries household expenses.
type ExpenseRepository interface {
	Create(ctx context.Context, expense *domain.Expense) error
	GetByID(ctx context.Context, id int64) (*domain.Expense, error)
	List(ctx context.Context, filter ExpenseFilter) ([]domain.Expense, error)
	Update(ctx context.Context, expense *domain.Expense) error
	Delete(ctx context.Context, id int64) error
}
