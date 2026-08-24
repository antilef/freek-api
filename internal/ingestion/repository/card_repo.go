// Package repository defines persistence interfaces for the ingestion context.
package repository

import (
	"context"
	"freek/internal/ingestion/domain"
)

// CardRepository persists and queries household credit cards.
type CardRepository interface {
	Create(ctx context.Context, card *domain.CreditCard) error
	GetByID(ctx context.Context, id int64) (*domain.CreditCard, error)
	ListByHousehold(ctx context.Context, householdID int64) ([]domain.CreditCard, error)
	ListByMember(ctx context.Context, householdID int64, memberID int64) ([]domain.CreditCard, error)
	Update(ctx context.Context, card *domain.CreditCard) error
	Delete(ctx context.Context, id int64) error
	HasExpenses(ctx context.Context, cardID int64) (bool, error)
}
