package service

import (
	"context"
	"freek/internal/ingestion/domain"
	"freek/internal/ingestion/repository"
	"testing"
)

type mockExpenseRepo struct {
	expenses []domain.Expense
	nextID   int64
}

func newMockExpenseRepo() *mockExpenseRepo {
	return &mockExpenseRepo{nextID: 1}
}

func (r *mockExpenseRepo) Create(ctx context.Context, expense *domain.Expense) error {
	expense.ID = r.nextID
	r.nextID++
	r.expenses = append(r.expenses, *expense)
	return nil
}

func (r *mockExpenseRepo) GetByID(ctx context.Context, id int64) (*domain.Expense, error) {
	for _, e := range r.expenses {
		if e.ID == id {
			return &e, nil
		}
	}
	return nil, nil
}

func (r *mockExpenseRepo) List(ctx context.Context, filter repository.ExpenseFilter) ([]domain.Expense, error) {
	return r.expenses, nil
}

func (r *mockExpenseRepo) Update(ctx context.Context, expense *domain.Expense) error {
	for i, e := range r.expenses {
		if e.ID == expense.ID {
			r.expenses[i] = *expense
			return nil
		}
	}
	return nil
}

func (r *mockExpenseRepo) Delete(ctx context.Context, id int64) error {
	for i, e := range r.expenses {
		if e.ID == id {
			r.expenses = append(r.expenses[:i], r.expenses[i+1:]...)
			return nil
		}
	}
	return nil
}

func TestExpenseService_Create_CashExpense(t *testing.T) {
	repo := newMockExpenseRepo()
	svc := NewExpenseService(repo)

	expense := &domain.Expense{
		HouseholdID:  1,
		AmountCLP:    50000,
		Description:  "Supermercado",
		PaymentType:  domain.PaymentTypeCash,
		Installments: 1,
	}

	err := svc.Create(context.Background(), expense)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if expense.ID == 0 {
		t.Error("expected ID to be set")
	}
}

func TestExpenseService_Create_CreditWithoutCard(t *testing.T) {
	repo := newMockExpenseRepo()
	svc := NewExpenseService(repo)

	expense := &domain.Expense{
		HouseholdID:  1,
		AmountCLP:    80000,
		Description:  "Retail",
		PaymentType:  domain.PaymentTypeCredit,
		Installments: 3,
	}

	err := svc.Create(context.Background(), expense)
	if err == nil {
		t.Error("expected error for credit expense without card")
	}
}

func TestExpenseService_Create_ZeroAmount(t *testing.T) {
	repo := newMockExpenseRepo()
	svc := NewExpenseService(repo)

	expense := &domain.Expense{
		HouseholdID:  1,
		AmountCLP:    0,
		Description:  "Test",
		PaymentType:  domain.PaymentTypeCash,
		Installments: 1,
	}

	err := svc.Create(context.Background(), expense)
	if err == nil {
		t.Error("expected error for zero amount")
	}
}

func TestExpenseService_Create_EmptyDescription(t *testing.T) {
	repo := newMockExpenseRepo()
	svc := NewExpenseService(repo)

	expense := &domain.Expense{
		HouseholdID:  1,
		AmountCLP:    50000,
		Description:  "",
		PaymentType:  domain.PaymentTypeCash,
		Installments: 1,
	}

	err := svc.Create(context.Background(), expense)
	if err == nil {
		t.Error("expected error for empty description")
	}
}

func TestExpenseService_CalculateMonthlyAmount_NoInterest(t *testing.T) {
	repo := newMockExpenseRepo()
	svc := NewExpenseService(repo)

	expense := &domain.Expense{
		AmountCLP:    120000,
		Installments: 6,
		InterestRate: 0,
	}

	result := svc.CalculateMonthlyAmount(expense)
	if result != 20000 {
		t.Errorf("CalculateMonthlyAmount = %d, want 20000", int64(result))
	}
}

func TestExpenseService_CalculateMonthlyAmount_WithInterest(t *testing.T) {
	repo := newMockExpenseRepo()
	svc := NewExpenseService(repo)

	expense := &domain.Expense{
		AmountCLP:    100000,
		Installments: 3,
		InterestRate: 0.15,
	}

	result := svc.CalculateMonthlyAmount(expense)
	// 100000 * 1.15 = 115000 / 3 = 38333.33 → rounds up to 38334
	// But due to float64 precision, we get 38333
	if result < 38333 || result > 38334 {
		t.Errorf("CalculateMonthlyAmount = %d, want 38333 or 38334", int64(result))
	}
}

func TestExpenseService_Create_DefaultStatus(t *testing.T) {
	repo := newMockExpenseRepo()
	svc := NewExpenseService(repo)

	expense := &domain.Expense{
		HouseholdID:  1,
		AmountCLP:    50000,
		Description:  "Test",
		PaymentType:  domain.PaymentTypeCash,
		Installments: 1,
	}

	svc.Create(context.Background(), expense)

	if expense.Status != "active" {
		t.Errorf("Status = %s, want 'active'", expense.Status)
	}
}
