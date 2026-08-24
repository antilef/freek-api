package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"freek/internal/ingestion/domain"
	"freek/internal/ingestion/service"
	"github.com/gorilla/mux"
)

type testHouseholdRepo struct {
	households []domain.Household
}

func (r *testHouseholdRepo) Create(ctx context.Context, h *domain.Household) error {
	h.ID = int64(len(r.households) + 1)
	r.households = append(r.households, *h)
	return nil
}

func (r *testHouseholdRepo) GetByID(ctx context.Context, id int64) (*domain.Household, error) {
	for _, h := range r.households {
		if h.ID == id {
			return &h, nil
		}
	}
	return nil, nil
}

func (r *testHouseholdRepo) Update(ctx context.Context, h *domain.Household) error {
	return nil
}

func (r *testHouseholdRepo) Delete(ctx context.Context, id int64) error {
	return nil
}

func TestHouseholdHandler_Create(t *testing.T) {
	repo := &testHouseholdRepo{}
	svc := service.NewHouseholdService(repo)
	handler := NewHouseholdHandler(svc, nil)

	body := `{"name":"Familia Pérez"}`
	req := httptest.NewRequest("POST", "/api/v1/households", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}

	var household domain.Household
	json.NewDecoder(w.Body).Decode(&household)

	if household.Name != "Familia Pérez" {
		t.Errorf("name = %s, want 'Familia Pérez'", household.Name)
	}
}

func TestHouseholdHandler_Create_MissingName(t *testing.T) {
	repo := &testHouseholdRepo{}
	svc := service.NewHouseholdService(repo)
	handler := NewHouseholdHandler(svc, nil)

	body := `{}`
	req := httptest.NewRequest("POST", "/api/v1/households", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHouseholdHandler_GetByID(t *testing.T) {
	repo := &testHouseholdRepo{}
	svc := service.NewHouseholdService(repo)
	handler := NewHouseholdHandler(svc, nil)

	req := httptest.NewRequest("GET", "/api/v1/households/1", nil)
	w := httptest.NewRecorder()

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	handler.GetByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}
