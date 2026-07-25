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

type HistoryHandler struct {
	historyService *service.HistoryService
}

func NewHistoryHandler(hs *service.HistoryService) *HistoryHandler {
	return &HistoryHandler{historyService: hs}
}

func (h *HistoryHandler) Query(w http.ResponseWriter, r *http.Request) {
	householdID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid household id"}`, http.StatusBadRequest)
		return
	}

	filter := repository.ExpenseReadFilter{HouseholdID: householdID}

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
			t := time.Now().AddDate(0, 0, -days)
			filter.From = &t
		}
	}
	if q := r.URL.Query().Get("q"); q != "" {
		filter.Search = &q
	}

	expenses, err := h.historyService.QueryByPeriod(r.Context(), filter)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}
