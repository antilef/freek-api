package repository

import (
	"context"
	"freek/internal/ingestion/domain"
)

// CategoryRepository persists and queries global and household categories.
type CategoryRepository interface {
	ListByHousehold(ctx context.Context, householdID int64) ([]domain.Category, error)
	Create(ctx context.Context, category *domain.Category) error
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*domain.Category, error)
	ExistsByName(ctx context.Context, householdID int64, name string) (bool, error)
	HasExpenses(ctx context.Context, categoryID int64) (bool, error)
}
