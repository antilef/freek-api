package handler

import (
	"encoding/json"
	"freek/internal/ingestion/domain"
	"freek/internal/ingestion/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// HouseholdHandler handles HTTP requests for household and member endpoints.
type HouseholdHandler struct {
	householdService *service.HouseholdService
	memberService    *service.MemberService
}

// NewHouseholdHandler creates a new HouseholdHandler.
func NewHouseholdHandler(hs *service.HouseholdService, ms *service.MemberService) *HouseholdHandler {
	return &HouseholdHandler{householdService: hs, memberService: ms}
}

// Create handles POST requests to create a household.
func (h *HouseholdHandler) Create(w http.ResponseWriter, r *http.Request) {
	var household domain.Household
	if err := json.NewDecoder(r.Body).Decode(&household); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if household.Name == "" {
		http.Error(w, `{"error":"name is required"}`, http.StatusBadRequest)
		return
	}

	if err := h.householdService.Create(r.Context(), &household); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// Response already committed; nothing useful to do on write failure.
	json.NewEncoder(w).Encode(household) //nolint:errcheck // response already committed
}

// GetByID handles GET requests for a single household.
func (h *HouseholdHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	household, err := h.householdService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	if household == nil {
		http.Error(w, `{"error":"household not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// Response already committed; nothing useful to do on write failure.
	json.NewEncoder(w).Encode(household) //nolint:errcheck // response already committed
}

// AddMember handles POST requests to add a member to a household.
func (h *HouseholdHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	householdID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid household id"}`, http.StatusBadRequest)
		return
	}

	var member domain.Member
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	member.HouseholdID = householdID

	if err := h.memberService.Create(r.Context(), &member); err != nil {
		if err.Error() == "member already exists" {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// Response already committed; nothing useful to do on write failure.
	json.NewEncoder(w).Encode(member) //nolint:errcheck // response already committed
}

// ListMembers handles GET requests returning a household's members.
func (h *HouseholdHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	householdID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid household id"}`, http.StatusBadRequest)
		return
	}

	members, err := h.memberService.ListByHousehold(r.Context(), householdID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	if members == nil {
		members = []domain.Member{}
	}

	w.Header().Set("Content-Type", "application/json")
	// Response already committed; nothing useful to do on write failure.
	json.NewEncoder(w).Encode(members) //nolint:errcheck // response already committed
}

// RemoveMember handles DELETE requests to remove a member from a household.
func (h *HouseholdHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	memberID, err := strconv.ParseInt(mux.Vars(r)["member_id"], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid member id"}`, http.StatusBadRequest)
		return
	}

	if err := h.memberService.Delete(r.Context(), memberID); err != nil {
		if err.Error() == "member has expenses, cannot delete" {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
