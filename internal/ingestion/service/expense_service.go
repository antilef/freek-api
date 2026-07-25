package service

import (
	"context"
	"fmt"
	"freek/internal/ingestion/domain"
	"freek/internal/ingestion/repository"
	"freek/pkg/dates"
	"freek/pkg/money"
	"time"
)

type ExpenseService struct {
	repo repository.ExpenseRepository
}

func NewExpenseService(repo repository.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo}
}

func (s *ExpenseService) Create(ctx context.Context, expense *domain.Expense) error {
	if expense.AmountCLP <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	if expense.Description == "" {
		return fmt.Errorf("description is required")
	}
	if expense.PaymentType == domain.PaymentTypeCredit && expense.CreditCardID == nil {
		return fmt.Errorf("credit_card_id is required for credit expenses")
	}
	if expense.Installments < 1 {
		expense.Installments = 1
	}

	if expense.Status == "" {
		expense.Status = "active"
	}

	return s.repo.Create(ctx, expense)
}

func (s *ExpenseService) GetByID(ctx context.Context, id int64) (*domain.Expense, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ExpenseService) List(ctx context.Context, filter repository.ExpenseFilter) ([]domain.Expense, error) {
	return s.repo.List(ctx, filter)
}

func (s *ExpenseService) Update(ctx context.Context, expense *domain.Expense) error {
	return s.repo.Update(ctx, expense)
}

func (s *ExpenseService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *ExpenseService) CalculateMonthlyAmount(expense *domain.Expense) money.CLP {
	total := money.CLP(expense.AmountCLP)
	if expense.InterestRate > 0 {
		total = money.CLP(float64(total) * (1 + expense.InterestRate))
	}
	return money.DivideRoundUp(total, expense.Installments)
}

func (s *ExpenseService) CalculateFirstSettlementDate(expense *domain.Expense, closingDay int) time.Time {
	return dates.NextClosingDate(expense.PurchaseDate, closingDay)
}
