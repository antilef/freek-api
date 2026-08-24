// Package service implements calculation read-side use cases.
package service

import (
	"context"
	"freek/internal/calculation/domain"
	calcRepo "freek/internal/calculation/repository"
	"freek/pkg/dates"
	"time"
)

// HistoryService queries expenses grouped by settlement period.
type HistoryService struct {
	expenseReadRepo calcRepo.ExpenseReadRepository
}

// NewHistoryService creates a new HistoryService.
func NewHistoryService(er calcRepo.ExpenseReadRepository) *HistoryService {
	return &HistoryService{expenseReadRepo: er}
}

// QueryByPeriod lists expenses within an explicit date range.
func (s *HistoryService) QueryByPeriod(ctx context.Context, filter calcRepo.ExpenseReadFilter) ([]domain.SettlementExpense, error) {
	return s.expenseReadRepo.ListByPeriod(ctx, filter)
}

// QueryBySettlement lists the expenses charged on a card's YYYY-MM settlement.
func (s *HistoryService) QueryBySettlement(ctx context.Context, householdID int64, cardID int64, settlement string, closingDay int) ([]domain.SettlementExpense, error) {
	t, err := time.Parse("2006-01", settlement)
	if err != nil {
		return nil, err
	}

	start, end := dates.SettlementForMonth(closingDay, t.Year(), t.Month())
	filter := calcRepo.ExpenseReadFilter{
		HouseholdID:  householdID,
		CreditCardID: &cardID,
		From:         &start,
		To:           &end,
	}

	return s.expenseReadRepo.ListByPeriod(ctx, filter)
}
