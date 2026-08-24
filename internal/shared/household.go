// Package shared provides cross-context HTTP middleware.
package shared

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

// contextKey is the private key type for context values.
type contextKey string

// HouseholdIDKey is the context key carrying the scoped household ID.
const HouseholdIDKey contextKey = "household_id"

// HouseholdContext extracts the household scope from the query string or path
// and stores it in the request context.
func HouseholdContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var householdID int64

		idStr := r.URL.Query().Get("household_id")
		if idStr != "" {
			var err error
			householdID, err = strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				http.Error(w, `{"error":"invalid household_id"}`, http.StatusBadRequest)
				return
			}
		} else {
			path := r.URL.Path
			parts := strings.Split(path, "/")
			for i, part := range parts {
				if part == "households" && i+1 < len(parts) {
					id, err := strconv.ParseInt(parts[i+1], 10, 64)
					if err == nil {
						householdID = id
						break
					}
				}
			}
		}

		if householdID == 0 {
			http.Error(w, `{"error":"household_id required"}`, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), HouseholdIDKey, householdID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetHouseholdID returns the scoped household ID from ctx, or 0 when absent.
func GetHouseholdID(ctx context.Context) int64 {
	if id, ok := ctx.Value(HouseholdIDKey).(int64); ok {
		return id
	}
	return 0
}
