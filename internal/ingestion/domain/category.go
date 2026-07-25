package domain

import "time"

type Category struct {
	ID           int64     `json:"id"`
	HouseholdID  *int64    `json:"household_id,omitempty"`
	Name         string    `json:"name"`
	Icon         *string   `json:"icon,omitempty"`
	IsPredefined bool      `json:"is_predefined"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
