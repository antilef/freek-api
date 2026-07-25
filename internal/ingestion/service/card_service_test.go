package service

import (
	"context"
	"freek/internal/ingestion/domain"
	"testing"
)

type mockCardRepo struct {
	cards []domain.CreditCard
}

func newMockCardRepo() *mockCardRepo {
	return &mockCardRepo{}
}

func (r *mockCardRepo) Create(ctx context.Context, card *domain.CreditCard) error {
	card.ID = int64(len(r.cards) + 1)
	r.cards = append(r.cards, *card)
	return nil
}

func (r *mockCardRepo) GetByID(ctx context.Context, id int64) (*domain.CreditCard, error) {
	for _, c := range r.cards {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, nil
}

func (r *mockCardRepo) ListByHousehold(ctx context.Context, householdID int64) ([]domain.CreditCard, error) {
	return r.cards, nil
}

func (r *mockCardRepo) ListByMember(ctx context.Context, householdID int64, memberID int64) ([]domain.CreditCard, error) {
	return r.cards, nil
}

func (r *mockCardRepo) Update(ctx context.Context, card *domain.CreditCard) error {
	return nil
}

func (r *mockCardRepo) Delete(ctx context.Context, id int64) error {
	return nil
}

func (r *mockCardRepo) HasExpenses(ctx context.Context, cardID int64) (bool, error) {
	return false, nil
}

func TestCardService_Create_ValidClosingDay(t *testing.T) {
	repo := newMockCardRepo()
	svc := NewCardService(repo)

	card := &domain.CreditCard{
		HouseholdID: 1,
		Name:        "CMR",
		ClosingDay:  15,
	}

	err := svc.Create(context.Background(), card)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if card.ID == 0 {
		t.Error("expected ID to be set")
	}
}

func TestCardService_Create_InvalidClosingDay_Zero(t *testing.T) {
	repo := newMockCardRepo()
	svc := NewCardService(repo)

	card := &domain.CreditCard{
		HouseholdID: 1,
		Name:        "CMR",
		ClosingDay:  0,
	}

	err := svc.Create(context.Background(), card)
	if err == nil {
		t.Error("expected error for closing_day = 0")
	}
}

func TestCardService_Create_InvalidClosingDay_29(t *testing.T) {
	repo := newMockCardRepo()
	svc := NewCardService(repo)

	card := &domain.CreditCard{
		HouseholdID: 1,
		Name:        "CMR",
		ClosingDay:  29,
	}

	err := svc.Create(context.Background(), card)
	if err == nil {
		t.Error("expected error for closing_day = 29")
	}
}

func TestCardService_Create_SharedCard(t *testing.T) {
	repo := newMockCardRepo()
	svc := NewCardService(repo)

	card := &domain.CreditCard{
		HouseholdID: 1,
		Name:        "Itaú",
		ClosingDay:  10,
		MemberID:    nil,
	}

	err := svc.Create(context.Background(), card)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if card.MemberID != nil {
		t.Error("expected MemberID to be nil for shared card")
	}
}

func TestCardService_Delete_NoExpenses(t *testing.T) {
	repo := newMockCardRepo()
	svc := NewCardService(repo)

	card := &domain.CreditCard{
		HouseholdID: 1,
		Name:        "CMR",
		ClosingDay:  15,
	}
	svc.Create(context.Background(), card)

	err := svc.Delete(context.Background(), card.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
