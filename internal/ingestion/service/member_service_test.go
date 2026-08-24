package service

import (
	"context"
	"freek/internal/ingestion/domain"
	"testing"
)

type mockMemberRepo struct {
	members []domain.Member
	nextID  int64
}

func newMockMemberRepo() *mockMemberRepo {
	return &mockMemberRepo{nextID: 1}
}

func (r *mockMemberRepo) Create(_ context.Context, member *domain.Member) error {
	member.ID = r.nextID
	r.nextID++
	r.members = append(r.members, *member)
	return nil
}

func (r *mockMemberRepo) GetByID(_ context.Context, id int64) (*domain.Member, error) {
	for _, m := range r.members {
		if m.ID == id {
			return &m, nil
		}
	}
	return nil, nil
}

func (r *mockMemberRepo) ListByHousehold(_ context.Context, _ int64) ([]domain.Member, error) {
	return r.members, nil
}

func (r *mockMemberRepo) Delete(_ context.Context, id int64) error {
	for i, m := range r.members {
		if m.ID == id {
			r.members = append(r.members[:i], r.members[i+1:]...)
			return nil
		}
	}
	return nil
}

func (r *mockMemberRepo) ExistsByName(_ context.Context, _ int64, name string) (bool, error) {
	for _, m := range r.members {
		if m.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func (r *mockMemberRepo) HasExpenses(_ context.Context, _ int64) (bool, error) {
	return false, nil
}

func TestMemberService_Create(t *testing.T) {
	repo := newMockMemberRepo()
	svc := NewMemberService(repo)

	member := &domain.Member{
		HouseholdID: 1,
		Name:        "Juan",
	}

	err := svc.Create(context.Background(), member)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if member.ID == 0 {
		t.Error("expected ID to be set")
	}
}

func TestMemberService_Create_DuplicateName(t *testing.T) {
	repo := newMockMemberRepo()
	svc := NewMemberService(repo)

	member1 := &domain.Member{HouseholdID: 1, Name: "Juan"}
	svc.Create(context.Background(), member1)

	member2 := &domain.Member{HouseholdID: 1, Name: "Juan"}
	err := svc.Create(context.Background(), member2)
	if err == nil {
		t.Error("expected error for duplicate name")
	}
}

func TestMemberService_ListByHousehold(t *testing.T) {
	repo := newMockMemberRepo()
	svc := NewMemberService(repo)

	svc.Create(context.Background(), &domain.Member{HouseholdID: 1, Name: "Juan"})
	svc.Create(context.Background(), &domain.Member{HouseholdID: 1, Name: "María"})

	members, err := svc.ListByHousehold(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(members) != 2 {
		t.Errorf("len(members) = %d, want 2", len(members))
	}
}

func TestMemberService_Delete_NoExpenses(t *testing.T) {
	repo := newMockMemberRepo()
	svc := NewMemberService(repo)

	member := &domain.Member{HouseholdID: 1, Name: "Juan"}
	svc.Create(context.Background(), member)

	err := svc.Delete(context.Background(), member.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
