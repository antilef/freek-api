package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"freek/internal/ingestion/domain"
	"freek/internal/ingestion/service"
	"github.com/gorilla/mux"
)

type CardHandler struct {
	cardService *service.CardService
}

func NewCardHandler(cs *service.CardService) *CardHandler {
	return &CardHandler{cardService: cs}
}

func (h *CardHandler) Create(w http.ResponseWriter, r *http.Request) {
	householdID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid household id"}`, http.StatusBadRequest)
		return
	}

	var card domain.CreditCard
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	card.HouseholdID = householdID

	if err := h.cardService.Create(r.Context(), &card); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(card)
}

func (h *CardHandler) List(w http.ResponseWriter, r *http.Request) {
	householdID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid household id"}`, http.StatusBadRequest)
		return
	}

	var cards []domain.CreditCard
	memberIDStr := r.URL.Query().Get("member_id")
	if memberIDStr != "" {
		memberID, parseErr := strconv.ParseInt(memberIDStr, 10, 64)
		if parseErr != nil {
			http.Error(w, `{"error":"invalid member_id"}`, http.StatusBadRequest)
			return
		}
		cards, err = h.cardService.ListByMember(r.Context(), householdID, memberID)
	} else {
		cards, err = h.cardService.ListByHousehold(r.Context(), householdID)
	}

	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	if cards == nil {
		cards = []domain.CreditCard{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}

func (h *CardHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	cardID, err := strconv.ParseInt(mux.Vars(r)["card_id"], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid card id"}`, http.StatusBadRequest)
		return
	}

	card, err := h.cardService.GetByID(r.Context(), cardID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	if card == nil {
		http.Error(w, `{"error":"card not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(card)
}

func (h *CardHandler) Update(w http.ResponseWriter, r *http.Request) {
	cardID, err := strconv.ParseInt(mux.Vars(r)["card_id"], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid card id"}`, http.StatusBadRequest)
		return
	}

	var card domain.CreditCard
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	card.ID = cardID

	if err := h.cardService.Update(r.Context(), &card); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(card)
}

func (h *CardHandler) Delete(w http.ResponseWriter, r *http.Request) {
	cardID, err := strconv.ParseInt(mux.Vars(r)["card_id"], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid card id"}`, http.StatusBadRequest)
		return
	}

	if err := h.cardService.Delete(r.Context(), cardID); err != nil {
		if err.Error() == "card has expenses, cannot delete" {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
