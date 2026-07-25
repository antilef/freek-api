package service

import (
	"context"
	"freek/internal/calculation/repository"
	"time"
)

type SummaryService struct {
	expenseReadRepo repository.ExpenseReadRepository
}

func NewSummaryService(er repository.ExpenseReadRepository) *SummaryService {
	return &SummaryService{expenseReadRepo: er}
}

func (s *SummaryService) SummaryByCategory(ctx context.Context, filter repository.ExpenseReadFilter) ([]repository.SummaryItem, error) {
	return s.expenseReadRepo.SummaryByCategory(ctx, filter)
}

func (s *SummaryService) SummaryByMember(ctx context.Context, filter repository.ExpenseReadFilter) ([]repository.SummaryItem, error) {
	return s.expenseReadRepo.SummaryByMember(ctx, filter)
}

func (s *SummaryService) ComparePeriods(ctx context.Context, householdID int64, period1From, period1To, period2From, period2To time.Time) (map[string]interface{}, error) {
	filter1 := repository.ExpenseReadFilter{
		HouseholdID: householdID,
		From:        &period1From,
		To:          &period1To,
	}
	filter2 := repository.ExpenseReadFilter{
		HouseholdID: householdID,
		From:        &period2From,
		To:          &period2To,
	}

	items1, err := s.expenseReadRepo.SummaryByCategory(ctx, filter1)
	if err != nil {
		return nil, err
	}

	items2, err := s.expenseReadRepo.SummaryByCategory(ctx, filter2)
	if err != nil {
		return nil, err
	}

	var total1, total2 int64
	for _, item := range items1 {
		total1 += item.Total
	}
	for _, item := range items2 {
		total2 += item.Total
	}

	return map[string]interface{}{
		"period1": map[string]interface{}{
			"from":  period1From,
			"to":    period1To,
			"total": total1,
			"items": items1,
		},
		"period2": map[string]interface{}{
			"from":  period2From,
			"to":    period2To,
			"total": total2,
			"items": items2,
		},
		"difference": total2 - total1,
	}, nil
}
