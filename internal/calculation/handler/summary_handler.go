package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"freek/internal/calculation/repository"
	"freek/internal/calculation/service"
	"github.com/gorilla/mux"
)

type SummaryHandler struct {
	summaryService *service.SummaryService
}

func NewSummaryHandler(ss *service.SummaryService) *SummaryHandler {
	return &SummaryHandler{summaryService: ss}
}

func (h *SummaryHandler) Summary(w http.ResponseWriter, r *http.Request) {
	householdID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid household id"}`, http.StatusBadRequest)
		return
	}

	filter := repository.ExpenseReadFilter{HouseholdID: householdID}

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

	groupBy := r.URL.Query().Get("group_by")

	var items []repository.SummaryItem
	if groupBy == "member" {
		items, err = h.summaryService.SummaryByMember(r.Context(), filter)
	} else {
		items, err = h.summaryService.SummaryByCategory(r.Context(), filter)
	}

	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *SummaryHandler) Compare(w http.ResponseWriter, r *http.Request) {
	householdID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid household id"}`, http.StatusBadRequest)
		return
	}

	period1FromStr := r.URL.Query().Get("period1_from")
	period1ToStr := r.URL.Query().Get("period1_to")
	period2FromStr := r.URL.Query().Get("period2_from")
	period2ToStr := r.URL.Query().Get("period2_to")

	period1From, err := time.Parse("2006-01-02", period1FromStr)
	if err != nil {
		http.Error(w, `{"error":"invalid period1_from"}`, http.StatusBadRequest)
		return
	}
	period1To, err := time.Parse("2006-01-02", period1ToStr)
	if err != nil {
		http.Error(w, `{"error":"invalid period1_to"}`, http.StatusBadRequest)
		return
	}
	period2From, err := time.Parse("2006-01-02", period2FromStr)
	if err != nil {
		http.Error(w, `{"error":"invalid period2_from"}`, http.StatusBadRequest)
		return
	}
	period2To, err := time.Parse("2006-01-02", period2ToStr)
	if err != nil {
		http.Error(w, `{"error":"invalid period2_to"}`, http.StatusBadRequest)
		return
	}

	result, err := h.summaryService.ComparePeriods(r.Context(), householdID, period1From, period1To, period2From, period2To)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
