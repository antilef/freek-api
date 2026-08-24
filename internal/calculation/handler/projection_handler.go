package handler

import (
	"encoding/json"
	"freek/internal/calculation/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ProjectionHandler handles HTTP requests for credit charge projections.
type ProjectionHandler struct {
	projectionService *service.ProjectionService
}

// NewProjectionHandler creates a new ProjectionHandler.
func NewProjectionHandler(ps *service.ProjectionService) *ProjectionHandler {
	return &ProjectionHandler{projectionService: ps}
}

// Project handles GET requests simulating upcoming installment charges per card or across all cards.
func (h *ProjectionHandler) Project(w http.ResponseWriter, r *http.Request) {
	householdID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid household id"}`, http.StatusBadRequest)
		return
	}

	months := 12
	if monthsStr := r.URL.Query().Get("months"); monthsStr != "" {
		if m, err := strconv.Atoi(monthsStr); err == nil && m > 0 {
			months = m
		}
	}

	var result *interface{}

	cardIDStr := r.URL.Query().Get("card_id")
	if cardIDStr != "" {
		cardID, err := strconv.ParseInt(cardIDStr, 10, 64)
		if err != nil {
			http.Error(w, `{"error":"invalid card_id"}`, http.StatusBadRequest)
			return
		}
		proj, err := h.projectionService.ProjectByCard(r.Context(), householdID, cardID, months)
		if err != nil {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		r2 := interface{}(proj)
		result = &r2
	} else {
		proj, err := h.projectionService.ProjectAll(r.Context(), householdID, months)
		if err != nil {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		r2 := interface{}(proj)
		result = &r2
	}

	w.Header().Set("Content-Type", "application/json")
	// Response headers are already committed; nothing useful to do on
	// write failure (Constitution #11).
	json.NewEncoder(w).Encode(result) //nolint:errcheck // response already committed
}
