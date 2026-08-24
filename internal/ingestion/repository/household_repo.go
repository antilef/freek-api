package repository

import (
	"context"
	"freek/internal/ingestion/domain"
)

// HouseholdRepository persists and queries households.
type HouseholdRepository interface {
	Create(ctx context.Context, household *domain.Household) error
	GetByID(ctx context.Context, id int64) (*domain.Household, error)
}
