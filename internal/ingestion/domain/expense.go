package domain

import "time"

type PaymentType string

const (
	PaymentTypeCash   PaymentType = "cash"
	PaymentTypeDebit  PaymentType = "debit"
	PaymentTypeCredit PaymentType = "credit"
)

type Expense struct {
	ID           int64      `json:"id"`
	HouseholdID  int64      `json:"household_id"`
	AmountCLP    int64      `json:"amount_clp"`
	Description  string     `json:"description"`
	Notes        *string    `json:"notes,omitempty"`
	CategoryID   *int64     `json:"category_id,omitempty"`
	PaymentType  PaymentType `json:"payment_type"`
	CreditCardID *int64     `json:"credit_card_id,omitempty"`
	MemberID     *int64     `json:"member_id,omitempty"`
	PurchaseDate time.Time  `json:"purchase_date"`
	Installments int        `json:"installments"`
	InterestRate float64    `json:"interest_rate"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
