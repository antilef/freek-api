package handler

import (
	"encoding/json"
	"errors"
	"freek/internal/ingestion/domain"
	"freek/internal/ingestion/repository"
	"freek/internal/ingestion/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// ExpenseHandler handles HTTP requests for expense endpoints.
type ExpenseHandler struct {
	expenseService *service.ExpenseService
}

// NewExpenseHandler creates a new ExpenseHandler.
func NewExpenseHandler(es *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{expenseService: es}
}

// Create handles POST requests to register a new expense.
func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	householdID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid household id"}`, http.StatusBadRequest)
		return
	}

	var expense domain.Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	expense.HouseholdID = householdID

	if err := h.expenseService.Create(r.Context(), &expense); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// Response headers are already committed; nothing useful to do on
	// write failure (Constitution #11).
	json.NewEncoder(w).Encode(expense) //nolint:errcheck // response already committed
}

// List handles GET requests returning filtered household expenses.
func (h *ExpenseHandler) List(w http.ResponseWriter, r *http.Request) {
	householdID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid household id"}`, http.StatusBadRequest)
		return
	}

	filter := repository.ExpenseFilter{HouseholdID: householdID}

	if cardIDStr := r.URL.Query().Get("card_id"); cardIDStr != "" {
		cardID, err := strconv.ParseInt(cardIDStr, 10, 64)
		if err == nil {
			filter.CreditCardID = &cardID
		}
	}
	if memberIDStr := r.URL.Query().Get("member_id"); memberIDStr != "" {
		memberID, err := strconv.ParseInt(memberIDStr, 10, 64)
		if err == nil {
			filter.MemberID = &memberID
		}
	}
	if categoryIDStr := r.URL.Query().Get("category_id"); categoryIDStr != "" {
		categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
		if err == nil {
			filter.CategoryID = &categoryID
		}
	}
	if fromStr := r.URL.Query().Get("from"); fromStr != "" {
		from, err := time.Parse("2006-01-02", fromStr)
		if err == nil {
			filter.From = &from
		}
	}
	if toStr := r.URL.Query().Get("to"); toStr != "" {
		to, err := time.Parse("2006-01-02", toStr)
		if err == nil {
			filter.To = &to
		}
	}
	if lastStr := r.URL.Query().Get("last"); lastStr != "" {
		days, err := strconv.Atoi(lastStr[:len(lastStr)-1])
		if err == nil && lastStr[len(lastStr)-1] == 'd' {
			filter.LastDays = &days
		}
	}
	if q := r.URL.Query().Get("q"); q != "" {
		filter.Search = &q
	}

	expenses, err := h.expenseService.List(r.Context(), filter)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	if expenses == nil {
		expenses = []domain.Expense{}
	}

	w.Header().Set("Content-Type", "application/json")
	// Response headers are already committed; nothing useful to do on
	// write failure (Constitution #11).
	json.NewEncoder(w).Encode(expenses) //nolint:errcheck // response already committed
}

// Update handles PUT requests to modify an expense.
func (h *ExpenseHandler) Update(w http.ResponseWriter, r *http.Request) {
	expenseID, err := strconv.ParseInt(mux.Vars(r)["expense_id"], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid expense id"}`, http.StatusBadRequest)
		return
	}

	var expense domain.Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	expense.ID = expenseID

	if err := h.expenseService.Update(r.Context(), &expense); err != nil {
		if errors.Is(err, service.ErrExpenseNotFound) {
			http.Error(w, `{"error":"expense not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// Response headers are already committed; nothing useful to do on
	// write failure (Constitution #11).
	json.NewEncoder(w).Encode(expense) //nolint:errcheck // response already committed
}

// Delete handles DELETE requests to permanently remove an expense.
func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	expenseID, err := strconv.ParseInt(mux.Vars(r)["expense_id"], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid expense id"}`, http.StatusBadRequest)
		return
	}

	if err := h.expenseService.Delete(r.Context(), expenseID); err != nil {
		if errors.Is(err, service.ErrExpenseNotFound) {
			http.Error(w, `{"error":"expense not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
