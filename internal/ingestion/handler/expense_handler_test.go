package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"freek/internal/ingestion/domain"
	"freek/internal/ingestion/repository"
	"freek/internal/ingestion/service"
	"github.com/gorilla/mux"
)

type testExpenseRepo struct {
	expenses []domain.Expense
}

func (r *testExpenseRepo) Create(ctx context.Context, expense *domain.Expense) error {
	expense.ID = int64(len(r.expenses) + 1)
	r.expenses = append(r.expenses, *expense)
	return nil
}

func (r *testExpenseRepo) GetByID(ctx context.Context, id int64) (*domain.Expense, error) {
	for _, e := range r.expenses {
		if e.ID == id {
			return &e, nil
		}
	}
	return nil, nil
}

func (r *testExpenseRepo) List(ctx context.Context, filter repository.ExpenseFilter) ([]domain.Expense, error) {
	return r.expenses, nil
}

func (r *testExpenseRepo) Update(ctx context.Context, expense *domain.Expense) error {
	return nil
}

func (r *testExpenseRepo) Delete(ctx context.Context, id int64) error {
	return nil
}

func TestExpenseHandler_Create_CashExpense(t *testing.T) {
	repo := &testExpenseRepo{}
	svc := service.NewExpenseService(repo)
	handler := NewExpenseHandler(svc)

	body := `{"amount_clp":50000,"description":"Supermercado","payment_type":"cash","installments":1}`
	req := httptest.NewRequest("POST", "/api/v1/households/1/expenses", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	handler.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}

	var expense domain.Expense
	json.NewDecoder(w.Body).Decode(&expense)

	if expense.AmountCLP != 50000 {
		t.Errorf("amount_clp = %d, want 50000", expense.AmountCLP)
	}
}

func TestExpenseHandler_Create_CreditWithoutCard(t *testing.T) {
	repo := &testExpenseRepo{}
	svc := service.NewExpenseService(repo)
	handler := NewExpenseHandler(svc)

	body := `{"amount_clp":80000,"description":"Retail","payment_type":"credit","installments":3}`
	req := httptest.NewRequest("POST", "/api/v1/households/1/expenses", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	handler.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestExpenseHandler_Create_ZeroAmount(t *testing.T) {
	repo := &testExpenseRepo{}
	svc := service.NewExpenseService(repo)
	handler := NewExpenseHandler(svc)

	body := `{"amount_clp":0,"description":"Test","payment_type":"cash","installments":1}`
	req := httptest.NewRequest("POST", "/api/v1/households/1/expenses", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	handler.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestExpenseHandler_List(t *testing.T) {
	repo := &testExpenseRepo{}
	svc := service.NewExpenseService(repo)
	handler := NewExpenseHandler(svc)

	req := httptest.NewRequest("GET", "/api/v1/households/1/expenses", nil)
	w := httptest.NewRecorder()

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var expenses []domain.Expense
	json.NewDecoder(w.Body).Decode(&expenses)

	if len(expenses) != 0 {
		t.Errorf("len(expenses) = %d, want 0", len(expenses))
	}
}
