package service

import (
	"context"
	"freek/internal/ingestion/domain"
	"freek/internal/ingestion/repository"
)

// HouseholdService ...
type HouseholdService struct {
	repo repository.HouseholdRepository
}

// NewHouseholdService creates a new HouseholdService.
func NewHouseholdService(repo repository.HouseholdRepository) *HouseholdService {
	return &HouseholdService{repo: repo}
}

// Create ...
func (s *HouseholdService) Create(ctx context.Context, household *domain.Household) error {
	return s.repo.Create(ctx, household)
}

// GetByID ...
func (s *HouseholdService) GetByID(ctx context.Context, id int64) (*domain.Household, error) {
	return s.repo.GetByID(ctx, id)
}
