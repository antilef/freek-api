package service

import (
	"context"
	"fmt"
	"freek/internal/ingestion/domain"
	"freek/internal/ingestion/repository"
	"freek/pkg/dates"
	"freek/pkg/money"
	"time"
)

type CardService struct {
	repo repository.CardRepository
}

func NewCardService(repo repository.CardRepository) *CardService {
	return &CardService{repo: repo}
}

func (s *CardService) Create(ctx context.Context, card *domain.CreditCard) error {
	if card.ClosingDay < 1 || card.ClosingDay > 28 {
		return fmt.Errorf("closing_day must be between 1 and 28")
	}
	return s.repo.Create(ctx, card)
}

func (s *CardService) ListByHousehold(ctx context.Context, householdID int64) ([]domain.CreditCard, error) {
	return s.repo.ListByHousehold(ctx, householdID)
}

func (s *CardService) ListByMember(ctx context.Context, householdID int64, memberID int64) ([]domain.CreditCard, error) {
	return s.repo.ListByMember(ctx, householdID, memberID)
}

func (s *CardService) GetByID(ctx context.Context, id int64) (*domain.CreditCard, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CardService) Update(ctx context.Context, card *domain.CreditCard) error {
	if card.ClosingDay < 1 || card.ClosingDay > 28 {
		return fmt.Errorf("closing_day must be between 1 and 28")
	}
	return s.repo.Update(ctx, card)
}

func (s *CardService) Delete(ctx context.Context, cardID int64) error {
	hasExpenses, err := s.repo.HasExpenses(ctx, cardID)
	if err != nil {
		return err
	}
	if hasExpenses {
		return fmt.Errorf("card has expenses, cannot delete")
	}
	return s.repo.Delete(ctx, cardID)
}

func (s *CardService) NextClosingDate(purchaseDate time.Time, closingDay int) time.Time {
	return dates.NextClosingDate(purchaseDate, closingDay)
}

func (s *CardService) CalculateMonthlyAmount(amountCLP int64, installments int, interestRate float64) money.CLP {
	total := money.CLP(amountCLP)
	if interestRate > 0 {
		total = money.CLP(float64(total) * (1 + interestRate))
	}
	return money.DivideRoundUp(total, installments)
}
