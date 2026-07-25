package domain

import "time"

type ProjectionMonth struct {
	Month       time.Time `json:"month"`
	Amount      int64     `json:"amount"`
	CardName    string    `json:"card_name"`
	CardID      int64     `json:"card_id"`
}

type ProjectionResult struct {
	Months      []ProjectionMonth `json:"months"`
	Total       int64             `json:"total"`
	MonthlyTotals map[string]int64 `json:"monthly_totals"`
}
