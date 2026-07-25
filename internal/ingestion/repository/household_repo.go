package repository

import (
	"context"
	"freek/internal/ingestion/domain"
)

type HouseholdRepository interface {
	Create(ctx context.Context, household *domain.Household) error
	GetByID(ctx context.Context, id int64) (*domain.Household, error)
}
