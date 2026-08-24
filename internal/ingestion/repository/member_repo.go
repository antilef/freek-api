package repository

import (
	"context"
	"freek/internal/ingestion/domain"
)

// MemberRepository persists and queries household members.
type MemberRepository interface {
	Create(ctx context.Context, member *domain.Member) error
	GetByID(ctx context.Context, id int64) (*domain.Member, error)
	ListByHousehold(ctx context.Context, householdID int64) ([]domain.Member, error)
	Delete(ctx context.Context, id int64) error
	ExistsByName(ctx context.Context, householdID int64, name string) (bool, error)
	HasExpenses(ctx context.Context, memberID int64) (bool, error)
}
