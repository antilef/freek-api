package domain

import "time"

type Member struct {
	ID          int64     `json:"id"`
	HouseholdID int64     `json:"household_id"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
