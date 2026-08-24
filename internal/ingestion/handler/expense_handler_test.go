package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"freek/internal/ingestion/domain"
	"freek/internal/ingestion/repository"
	"freek/internal/ingestion/service"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

type testExpenseRepo struct {
	expenses []domain.Expense
}

func (r *testExpenseRepo) Create(_ context.Context, expense *domain.Expense) error {
	expense.ID = int64(len(r.expenses) + 1)
	r.expenses = append(r.expenses, *expense)
	return nil
}

func (r *testExpenseRepo) GetByID(_ context.Context, id int64) (*domain.Expense, error) {
	for _, e := range r.expenses {
		if e.ID == id {
			return &e, nil
		}
	}
	return nil, nil
}

func (r *testExpenseRepo) List(_ context.Context, _ repository.ExpenseFilter) ([]domain.Expense, error) {
	return r.expenses, nil
}

func (r *testExpenseRepo) Update(_ context.Context, _ *domain.Expense) error {
	return nil
}

func (r *testExpenseRepo) Delete(_ context.Context, _ int64) error {
	return nil
}

func TestExpenseHandler_Create_CashExpense(t *testing.T) {
	repo := &testExpenseRepo{}
	svc := service.NewExpenseService(repo)
	handler := NewExpenseHandler(svc)

	body := `{"amount_clp":50000,"description":"Supermercado","payment_type":"cash","installments":1}`
	req := httptest.NewRequestWithContext(context.Background(), "POST", "/api/v1/households/1/expenses", bytes.NewBufferString(body))
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
	req := httptest.NewRequestWithContext(context.Background(), "POST", "/api/v1/households/1/expenses", bytes.NewBufferString(body))
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
	req := httptest.NewRequestWithContext(context.Background(), "POST", "/api/v1/households/1/expenses", bytes.NewBufferString(body))
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

	req := httptest.NewRequestWithContext(context.Background(), "GET", "/api/v1/households/1/expenses", nil)
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

func TestExpenseHandler_Update_InvalidAmount(t *testing.T) {
	repo := &testExpenseRepo{}
	svc := service.NewExpenseService(repo)
	handler := NewExpenseHandler(svc)

	seed := &domain.Expense{
		ID:          1,
		HouseholdID: 1,
		AmountCLP:   50000,
		Description: "Test",
		PaymentType: domain.PaymentTypeCash,
	}
	repo.Create(context.Background(), seed)

	body := `{"amount_clp":0,"description":"Test","payment_type":"cash"}`
	req := httptest.NewRequestWithContext(context.Background(), "PUT", "/api/v1/households/1/expenses/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	vars := map[string]string{"id": "1", "expense_id": "1"}
	req = mux.SetURLVars(req, vars)

	handler.Update(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestExpenseHandler_Update_NotFound(t *testing.T) {
	repo := &testExpenseRepo{}
	svc := service.NewExpenseService(repo)
	handler := NewExpenseHandler(svc)

	body := `{"amount_clp":50000,"description":"Test","payment_type":"cash"}`
	req := httptest.NewRequestWithContext(context.Background(), "PUT", "/api/v1/households/1/expenses/999", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	vars := map[string]string{"id": "1", "expense_id": "999"}
	req = mux.SetURLVars(req, vars)

	handler.Update(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestExpenseHandler_Delete_NotFound(t *testing.T) {
	repo := &testExpenseRepo{}
	svc := service.NewExpenseService(repo)
	handler := NewExpenseHandler(svc)

	req := httptest.NewRequestWithContext(context.Background(), "DELETE", "/api/v1/households/1/expenses/999", nil)
	w := httptest.NewRecorder()

	vars := map[string]string{"id": "1", "expense_id": "999"}
	req = mux.SetURLVars(req, vars)

	handler.Delete(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestExpenseHandler_Create_DefaultPurchaseDate(t *testing.T) {
	repo := &testExpenseRepo{}
	svc := service.NewExpenseService(repo)
	handler := NewExpenseHandler(svc)

	body := `{"amount_clp":50000,"description":"Test","payment_type":"cash","installments":1}`
	req := httptest.NewRequestWithContext(context.Background(), "POST", "/api/v1/households/1/expenses", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	handler.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusCreated)
	}

	var expense domain.Expense
	json.NewDecoder(w.Body).Decode(&expense)

	now := time.Now()
	want := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	if !expense.PurchaseDate.Equal(want) {
		t.Errorf("purchase_date = %v, want %v (today)", expense.PurchaseDate, want)
	}
}
