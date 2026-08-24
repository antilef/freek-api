package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"freek/internal/ingestion/domain"
	"freek/internal/ingestion/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

type testCardRepo struct {
	cards []domain.CreditCard
}

func (r *testCardRepo) Create(_ context.Context, card *domain.CreditCard) error {
	card.ID = int64(len(r.cards) + 1)
	r.cards = append(r.cards, *card)
	return nil
}

func (r *testCardRepo) GetByID(_ context.Context, id int64) (*domain.CreditCard, error) {
	for _, c := range r.cards {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, nil
}

func (r *testCardRepo) ListByHousehold(_ context.Context, _ int64) ([]domain.CreditCard, error) {
	return r.cards, nil
}

func (r *testCardRepo) ListByMember(_ context.Context, _ int64, _ int64) ([]domain.CreditCard, error) {
	return r.cards, nil
}

func (r *testCardRepo) Update(_ context.Context, _ *domain.CreditCard) error {
	return nil
}

func (r *testCardRepo) Delete(_ context.Context, _ int64) error {
	return nil
}

func (r *testCardRepo) HasExpenses(_ context.Context, _ int64) (bool, error) {
	return false, nil
}

func TestCardHandler_Create(t *testing.T) {
	repo := &testCardRepo{}
	svc := service.NewCardService(repo)
	handler := NewCardHandler(svc)

	body := `{"name":"CMR","closing_day":15}`
	req := httptest.NewRequestWithContext(context.Background(), "POST", "/api/v1/households/1/cards", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	handler.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}

	var card domain.CreditCard
	json.NewDecoder(w.Body).Decode(&card)

	if card.Name != "CMR" {
		t.Errorf("name = %s, want 'CMR'", card.Name)
	}
	if card.ClosingDay != 15 {
		t.Errorf("closing_day = %d, want 15", card.ClosingDay)
	}
}

func TestCardHandler_Create_InvalidClosingDay(t *testing.T) {
	repo := &testCardRepo{}
	svc := service.NewCardService(repo)
	handler := NewCardHandler(svc)

	body := `{"name":"CMR","closing_day":29}`
	req := httptest.NewRequestWithContext(context.Background(), "POST", "/api/v1/households/1/cards", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	handler.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCardHandler_List(t *testing.T) {
	repo := &testCardRepo{}
	svc := service.NewCardService(repo)
	handler := NewCardHandler(svc)

	req := httptest.NewRequestWithContext(context.Background(), "GET", "/api/v1/households/1/cards", nil)
	w := httptest.NewRecorder()

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var cards []domain.CreditCard
	json.NewDecoder(w.Body).Decode(&cards)

	if len(cards) != 0 {
		t.Errorf("len(cards) = %d, want 0", len(cards))
	}
}
