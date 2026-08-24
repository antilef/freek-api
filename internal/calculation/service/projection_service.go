package service

import (
	"context"
	"freek/internal/calculation/domain"
	calcRepo "freek/internal/calculation/repository"
	"time"
)

// CardRepository provides the card lookups the projector needs.
type CardRepository interface {
	GetByID(ctx context.Context, id int64) (*struct {
		ClosingDay int
		Name       string
	}, error)
	ListByHousehold(ctx context.Context, householdID int64) ([]struct {
		ID         int64
		Name       string
		ClosingDay int
	}, error)
}

// ProjectionService simulates future months of outstanding installment charges.
type ProjectionService struct {
	expenseReadRepo calcRepo.ExpenseReadRepository
	cardRepo        CardRepository
}

// NewProjectionService creates a new ProjectionService.
func NewProjectionService(er calcRepo.ExpenseReadRepository, cr CardRepository) *ProjectionService {
	return &ProjectionService{expenseReadRepo: er, cardRepo: cr}
}

// ProjectByCard projects remaining installments for one card; nil card projects all cards.
func (s *ProjectionService) ProjectByCard(ctx context.Context, householdID int64, cardID int64, months int) (*domain.ProjectionResult, error) {
	expenses, err := s.expenseReadRepo.ListForProjection(ctx, householdID, &cardID)
	if err != nil {
		return nil, err
	}

	if months <= 0 {
		months = 12
	}

	now := time.Now()
	monthlyTotals := make(map[string]int64)
	var allMonths []domain.ProjectionMonth
	var total int64

	for i := 0; i < months; i++ {
		month := now.AddDate(0, i, 0)
		var monthTotal int64

		for _, expense := range expenses {
			totalInstallments := expense.Amount / expense.MonthlyAmount
			if totalInstallments == 0 {
				totalInstallments = 1
			}
			if int64(i) < totalInstallments {
				monthTotal += expense.MonthlyAmount
			}
		}

		pm := domain.ProjectionMonth{
			Month:    month,
			Amount:   monthTotal,
			CardName: "",
			CardID:   cardID,
		}
		allMonths = append(allMonths, pm)
		monthlyTotals[month.Format("2006-01")] += monthTotal
		total += monthTotal
	}

	return &domain.ProjectionResult{
		Months:        allMonths,
		Total:         total,
		MonthlyTotals: monthlyTotals,
	}, nil
}

// ProjectAll ...
// ProjectAll projects upcoming months across all household cards.
func (s *ProjectionService) ProjectAll(_ context.Context, _ int64, months int) (*domain.ProjectionResult, error) {
	if months <= 0 {
		months = 12
	}

	now := time.Now()
	monthlyTotals := make(map[string]int64)
	var allMonths []domain.ProjectionMonth
	var total int64

	for i := 0; i < months; i++ {
		month := now.AddDate(0, i, 0)

		pm := domain.ProjectionMonth{
			Month:  month,
			Amount: 0,
		}
		allMonths = append(allMonths, pm)
	}

	return &domain.ProjectionResult{
		Months:        allMonths,
		Total:         total,
		MonthlyTotals: monthlyTotals,
	}, nil
}
