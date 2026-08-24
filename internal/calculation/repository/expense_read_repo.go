// Package repository defines calculation read-side persistence interfaces.
package repository

import (
	"context"
	"freek/internal/calculation/domain"
	"time"
)

// ExpenseReadFilter ...
type ExpenseReadFilter struct {
	HouseholdID  int64
	CreditCardID *int64
	MemberID     *int64
	From         *time.Time
	To           *time.Time
	Settlement   *string
	Search       *string
}

// SummaryItem ...
type SummaryItem struct {
	CategoryID   *int64
	CategoryName string
	MemberID     *int64
	MemberName   string
	Total        int64
}

// ExpenseReadRepository ...
type ExpenseReadRepository interface {
	ListForProjection(ctx context.Context, householdID int64, cardID *int64) ([]domain.SettlementExpense, error)
	ListByPeriod(ctx context.Context, filter ExpenseReadFilter) ([]domain.SettlementExpense, error)
	SummaryByCategory(ctx context.Context, filter ExpenseReadFilter) ([]SummaryItem, error)
	SummaryByMember(ctx context.Context, filter ExpenseReadFilter) ([]SummaryItem, error)
}
