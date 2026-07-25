package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"freek/internal/ingestion/domain"
	"freek/internal/ingestion/service"
	"github.com/gorilla/mux"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(cs *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: cs}
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	householdID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid household id"}`, http.StatusBadRequest)
		return
	}

	categories, err := h.categoryService.ListByHousehold(r.Context(), householdID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	if categories == nil {
		categories = []domain.Category{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	householdID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid household id"}`, http.StatusBadRequest)
		return
	}

	var category domain.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	category.HouseholdID = &householdID
	category.IsPredefined = false

	if err := h.categoryService.Create(r.Context(), &category); err != nil {
		if err.Error() == "category already exists" {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.ParseInt(mux.Vars(r)["category_id"], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid category id"}`, http.StatusBadRequest)
		return
	}

	var category domain.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	category.ID = categoryID

	if err := h.categoryService.Update(r.Context(), &category); err != nil {
		if err.Error() == "cannot rename predefined category" {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusForbidden)
			return
		}
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.ParseInt(mux.Vars(r)["category_id"], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid category id"}`, http.StatusBadRequest)
		return
	}

	if err := h.categoryService.Delete(r.Context(), categoryID); err != nil {
		if err.Error() == "cannot delete predefined category" {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusForbidden)
			return
		}
		if err.Error() == "category has expenses, cannot delete" {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
