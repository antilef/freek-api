package service

import (
	"context"
	"fmt"
	"freek/internal/ingestion/domain"
	"freek/internal/ingestion/repository"
)

// CategoryService ...
type CategoryService struct {
	repo repository.CategoryRepository
}

// NewCategoryService creates a new CategoryService.
func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// ListByHousehold ...
func (s *CategoryService) ListByHousehold(ctx context.Context, householdID int64) ([]domain.Category, error) {
	return s.repo.ListByHousehold(ctx, householdID)
}

// Create ...
func (s *CategoryService) Create(ctx context.Context, category *domain.Category) error {
	exists, err := s.repo.ExistsByName(ctx, *category.HouseholdID, category.Name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("category already exists")
	}
	return s.repo.Create(ctx, category)
}

// Update ...
func (s *CategoryService) Update(ctx context.Context, category *domain.Category) error {
	existing, err := s.repo.GetByID(ctx, category.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("category not found")
	}
	if existing.IsPredefined && category.Name != existing.Name {
		return fmt.Errorf("cannot rename predefined category")
	}
	return s.repo.Update(ctx, category)
}

// Delete ...
func (s *CategoryService) Delete(ctx context.Context, categoryID int64) error {
	existing, err := s.repo.GetByID(ctx, categoryID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("category not found")
	}
	if existing.IsPredefined {
		return fmt.Errorf("cannot delete predefined category")
	}

	hasExpenses, err := s.repo.HasExpenses(ctx, categoryID)
	if err != nil {
		return err
	}
	if hasExpenses {
		return fmt.Errorf("category has expenses, cannot delete")
	}
	return s.repo.Delete(ctx, categoryID)
}
