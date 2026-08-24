package handler

import (
	"encoding/json"
	"freek/internal/ingestion/domain"
	"freek/internal/ingestion/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CategoryHandler handles HTTP requests for category endpoints.
type CategoryHandler struct {
	categoryService *service.CategoryService
}

// NewCategoryHandler creates a new CategoryHandler.
func NewCategoryHandler(cs *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: cs}
}

// List handles GET requests returning global and household categories.
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
	// Response already committed; nothing useful to do on write failure.
	json.NewEncoder(w).Encode(categories) //nolint:errcheck // response already committed
}

// Create handles POST requests to create a household category.
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
	// Response already committed; nothing useful to do on write failure.
	json.NewEncoder(w).Encode(category) //nolint:errcheck // response already committed
}

// Update handles PUT requests to rename or re-icon a category.
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
	// Response already committed; nothing useful to do on write failure.
	json.NewEncoder(w).Encode(category) //nolint:errcheck // response already committed
}

// Delete handles DELETE requests to remove a category.
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
