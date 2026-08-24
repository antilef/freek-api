package domain

import "time"

// SettlementPeriod ...
type SettlementPeriod struct {
	Period      time.Time           `json:"period"`
	TotalAmount int64               `json:"total_amount"`
	Expenses    []SettlementExpense `json:"expenses"`
}

// SettlementExpense ...
type SettlementExpense struct {
	ExpenseID     int64  `json:"expense_id"`
	Description   string `json:"description"`
	Amount        int64  `json:"amount"`
	MonthlyAmount int64  `json:"monthly_amount"`
}
