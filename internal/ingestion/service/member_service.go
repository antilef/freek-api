package service

import (
	"context"
	"fmt"
	"freek/internal/ingestion/domain"
	"freek/internal/ingestion/repository"
)

type MemberService struct {
	repo repository.MemberRepository
}

func NewMemberService(repo repository.MemberRepository) *MemberService {
	return &MemberService{repo: repo}
}

func (s *MemberService) Create(ctx context.Context, member *domain.Member) error {
	exists, err := s.repo.ExistsByName(ctx, member.HouseholdID, member.Name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("member already exists")
	}
	return s.repo.Create(ctx, member)
}

func (s *MemberService) ListByHousehold(ctx context.Context, householdID int64) ([]domain.Member, error) {
	return s.repo.ListByHousehold(ctx, householdID)
}

func (s *MemberService) Delete(ctx context.Context, memberID int64) error {
	hasExpenses, err := s.repo.HasExpenses(ctx, memberID)
	if err != nil {
		return err
	}
	if hasExpenses {
		return fmt.Errorf("member has expenses, cannot delete")
	}
	return s.repo.Delete(ctx, memberID)
}
