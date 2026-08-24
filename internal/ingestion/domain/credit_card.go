package domain

import "time"

// CreditCard ...
type CreditCard struct {
	ID          int64     `json:"id"`
	HouseholdID int64     `json:"household_id"`
	Name        string    `json:"name"`
	ClosingDay  int       `json:"closing_day"`
	MemberID    *int64    `json:"member_id,omitempty"`
	CreditLimit *int64    `json:"credit_limit,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
