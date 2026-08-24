package service

import (
	"context"
	"errors"
	"fmt"
	"freek/internal/ingestion/domain"
	"freek/internal/ingestion/repository"
	"freek/pkg/dates"
	"freek/pkg/money"
	"time"
)

// ErrExpenseNotFound is returned when an expense does not exist in the household.
var ErrExpenseNotFound = errors.New("expense not found")

// ExpenseService enforces expense business rules over the repository.
type ExpenseService struct {
	repo repository.ExpenseRepository
}

// NewExpenseService creates a new ExpenseService.
func NewExpenseService(repo repository.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo}
}

func (s *ExpenseService) validate(expense *domain.Expense) error {
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
	return nil
}

// Create validates the expense and persists it, defaulting purchase_date to today.
func (s *ExpenseService) Create(ctx context.Context, expense *domain.Expense) error {
	if err := s.validate(expense); err != nil {
		return err
	}

	if expense.PurchaseDate.IsZero() {
		now := time.Now()
		expense.PurchaseDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	}

	return s.repo.Create(ctx, expense)
}

// GetByID returns an expense by ID within the household scope.
func (s *ExpenseService) GetByID(ctx context.Context, id int64) (*domain.Expense, error) {
	return s.repo.GetByID(ctx, id)
}

// List returns expenses matching the filter within the household scope.
func (s *ExpenseService) List(ctx context.Context, filter repository.ExpenseFilter) ([]domain.Expense, error) {
	return s.repo.List(ctx, filter)
}

// Update validates changes and persists them; not-found yields ErrExpenseNotFound.
func (s *ExpenseService) Update(ctx context.Context, expense *domain.Expense) error {
	existing, err := s.repo.GetByID(ctx, expense.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrExpenseNotFound
	}

	if err := s.validate(expense); err != nil {
		return err
	}

	return s.repo.Update(ctx, expense)
}

// Delete permanently removes an expense; not-found yields ErrExpenseNotFound.
func (s *ExpenseService) Delete(ctx context.Context, id int64) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrExpenseNotFound
	}

	return s.repo.Delete(ctx, id)
}

// CalculateMonthlyAmount derives the round-up monthly installment at read time.
func (s *ExpenseService) CalculateMonthlyAmount(expense *domain.Expense) money.CLP {
	total := money.TotalWithInterest(expense.AmountCLP, expense.InterestRate)
	return money.MonthlyInstallment(total, expense.Installments)
}

// CalculateFirstSettlementDate returns the settlement date charged with the first installment.
func (s *ExpenseService) CalculateFirstSettlementDate(expense *domain.Expense, closingDay int) time.Time {
	return dates.NextClosingDate(expense.PurchaseDate, closingDay)
}
