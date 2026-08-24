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

type testCategoryRepo struct {
	categories []domain.Category
}

func (r *testCategoryRepo) ListByHousehold(_ context.Context, _ int64) ([]domain.Category, error) {
	return r.categories, nil
}

func (r *testCategoryRepo) Create(_ context.Context, category *domain.Category) error {
	category.ID = int64(len(r.categories) + 1)
	r.categories = append(r.categories, *category)
	return nil
}

func (r *testCategoryRepo) Update(_ context.Context, _ *domain.Category) error {
	return nil
}

func (r *testCategoryRepo) Delete(_ context.Context, _ int64) error {
	return nil
}

func (r *testCategoryRepo) GetByID(_ context.Context, id int64) (*domain.Category, error) {
	for _, c := range r.categories {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, nil
}

func (r *testCategoryRepo) ExistsByName(_ context.Context, _ int64, _ string) (bool, error) {
	return false, nil
}

func (r *testCategoryRepo) HasExpenses(_ context.Context, _ int64) (bool, error) {
	return false, nil
}

func TestCategoryHandler_List(t *testing.T) {
	repo := &testCategoryRepo{}
	svc := service.NewCategoryService(repo)
	handler := NewCategoryHandler(svc)

	req := httptest.NewRequestWithContext(context.Background(), "GET", "/api/v1/households/1/categories", nil)
	w := httptest.NewRecorder()

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var categories []domain.Category
	json.NewDecoder(w.Body).Decode(&categories)

	if len(categories) != 0 {
		t.Errorf("len(categories) = %d, want 0", len(categories))
	}
}

func TestCategoryHandler_Create(t *testing.T) {
	repo := &testCategoryRepo{}
	svc := service.NewCategoryService(repo)
	handler := NewCategoryHandler(svc)

	body := `{"name":"Mascotas","icon":"paw"}`
	req := httptest.NewRequestWithContext(context.Background(), "POST", "/api/v1/households/1/categories", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)

	handler.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}

	var category domain.Category
	json.NewDecoder(w.Body).Decode(&category)

	if category.Name != "Mascotas" {
		t.Errorf("name = %s, want 'Mascotas'", category.Name)
	}
	if category.IsPredefined {
		t.Error("expected is_predefined to be false")
	}
}
