package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"freek/internal/ingestion/domain"
	"freek/internal/ingestion/repository"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	dbURL := "postgres://localhost:5432/freek_test?sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Skip("Skipping: PostgreSQL not available:", err)
	}
	if err := db.Ping(); err != nil {
		t.Skip("Skipping: Cannot connect to PostgreSQL:", err)
	}

	// Clean test data
	db.Exec("DELETE FROM expenses")
	db.Exec("DELETE FROM credit_cards")
	db.Exec("DELETE FROM members")
	db.Exec("DELETE FROM categories WHERE household_id IS NOT NULL")
	db.Exec("DELETE FROM households")

	return db
}

func TestHouseholdRepo_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewHouseholdRepo(db)

	household := &domain.Household{Name: "Test Family"}
	err := repo.Create(context.Background(), household)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if household.ID == 0 {
		t.Error("expected ID to be set")
	}
	if household.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestHouseholdRepo_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewHouseholdRepo(db)

	household := &domain.Household{Name: "Test Family"}
	repo.Create(context.Background(), household)

	found, err := repo.GetByID(context.Background(), household.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if found == nil {
		t.Fatal("expected household to be found")
	}
	if found.Name != "Test Family" {
		t.Errorf("name = %s, want 'Test Family'", found.Name)
	}
}

func TestHouseholdRepo_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewHouseholdRepo(db)

	found, err := repo.GetByID(context.Background(), 999999)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if found != nil {
		t.Error("expected nil for non-existent household")
	}
}

func TestMemberRepo_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	householdRepo := NewHouseholdRepo(db)
	memberRepo := NewMemberRepo(db)

	household := &domain.Household{Name: "Test Family"}
	householdRepo.Create(context.Background(), household)

	member := &domain.Member{HouseholdID: household.ID, Name: "Juan"}
	err := memberRepo.Create(context.Background(), member)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if member.ID == 0 {
		t.Error("expected ID to be set")
	}
}

func TestMemberRepo_ListByHousehold(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	householdRepo := NewHouseholdRepo(db)
	memberRepo := NewMemberRepo(db)

	household := &domain.Household{Name: "Test Family"}
	householdRepo.Create(context.Background(), household)

	memberRepo.Create(context.Background(), &domain.Member{HouseholdID: household.ID, Name: "Juan"})
	memberRepo.Create(context.Background(), &domain.Member{HouseholdID: household.ID, Name: "María"})

	members, err := memberRepo.ListByHousehold(context.Background(), household.ID)
	if err != nil {
		t.Fatalf("ListByHousehold failed: %v", err)
	}

	if len(members) != 2 {
		t.Errorf("len(members) = %d, want 2", len(members))
	}
}

func TestMemberRepo_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	householdRepo := NewHouseholdRepo(db)
	memberRepo := NewMemberRepo(db)

	household := &domain.Household{Name: "Test Family"}
	householdRepo.Create(context.Background(), household)

	member := &domain.Member{HouseholdID: household.ID, Name: "Juan"}
	memberRepo.Create(context.Background(), member)

	err := memberRepo.Delete(context.Background(), member.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	found, _ := memberRepo.GetByID(context.Background(), member.ID)
	if found != nil {
		t.Error("expected member to be deleted")
	}
}

func TestCardRepo_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	householdRepo := NewHouseholdRepo(db)
	cardRepo := NewCardRepo(db)

	household := &domain.Household{Name: "Test Family"}
	householdRepo.Create(context.Background(), household)

	card := &domain.CreditCard{
		HouseholdID: household.ID,
		Name:        "CMR",
		ClosingDay:  15,
	}
	err := cardRepo.Create(context.Background(), card)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if card.ID == 0 {
		t.Error("expected ID to be set")
	}
}

func TestCardRepo_ListByHousehold(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	householdRepo := NewHouseholdRepo(db)
	cardRepo := NewCardRepo(db)

	household := &domain.Household{Name: "Test Family"}
	householdRepo.Create(context.Background(), household)

	cardRepo.Create(context.Background(), &domain.CreditCard{
		HouseholdID: household.ID,
		Name:        "CMR",
		ClosingDay:  15,
	})
	cardRepo.Create(context.Background(), &domain.CreditCard{
		HouseholdID: household.ID,
		Name:        "Scotiabank",
		ClosingDay:  20,
	})

	cards, err := cardRepo.ListByHousehold(context.Background(), household.ID)
	if err != nil {
		t.Fatalf("ListByHousehold failed: %v", err)
	}

	if len(cards) != 2 {
		t.Errorf("len(cards) = %d, want 2", len(cards))
	}
}

func TestCardRepo_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	householdRepo := NewHouseholdRepo(db)
	cardRepo := NewCardRepo(db)

	household := &domain.Household{Name: "Test Family"}
	householdRepo.Create(context.Background(), household)

	card := &domain.CreditCard{
		HouseholdID: household.ID,
		Name:        "CMR",
		ClosingDay:  15,
	}
	cardRepo.Create(context.Background(), card)

	card.Name = "CMR Actualizado"
	card.ClosingDay = 20
	err := cardRepo.Update(context.Background(), card)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	found, _ := cardRepo.GetByID(context.Background(), card.ID)
	if found.Name != "CMR Actualizado" {
		t.Errorf("name = %s, want 'CMR Actualizado'", found.Name)
	}
	if found.ClosingDay != 20 {
		t.Errorf("closing_day = %d, want 20", found.ClosingDay)
	}
}

func TestExpenseRepo_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	householdRepo := NewHouseholdRepo(db)
	expenseRepo := NewExpenseRepo(db)

	household := &domain.Household{Name: "Test Family"}
	householdRepo.Create(context.Background(), household)

	now := time.Now()
	expense := &domain.Expense{
		HouseholdID:  household.ID,
		AmountCLP:    50000,
		Description:  "Supermercado",
		PaymentType:  domain.PaymentTypeCash,
		PurchaseDate: now,
		Installments: 1,
		Status:       "active",
	}
	err := expenseRepo.Create(context.Background(), expense)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if expense.ID == 0 {
		t.Error("expected ID to be set")
	}
}

func TestExpenseRepo_List(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	householdRepo := NewHouseholdRepo(db)
	expenseRepo := NewExpenseRepo(db)

	household := &domain.Household{Name: "Test Family"}
	householdRepo.Create(context.Background(), household)

	now := time.Now()
	expenseRepo.Create(context.Background(), &domain.Expense{
		HouseholdID:  household.ID,
		AmountCLP:    50000,
		Description:  "Supermercado",
		PaymentType:  domain.PaymentTypeCash,
		PurchaseDate: now,
		Installments: 1,
		Status:       "active",
	})
	expenseRepo.Create(context.Background(), &domain.Expense{
		HouseholdID:  household.ID,
		AmountCLP:    30000,
		Description:  "Bencina",
		PaymentType:  domain.PaymentTypeDebit,
		PurchaseDate: now,
		Installments: 1,
		Status:       "active",
	})

	expenses, err := expenseRepo.List(context.Background(), repository.ExpenseFilter{HouseholdID: household.ID})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(expenses) != 2 {
		t.Errorf("len(expenses) = %d, want 2", len(expenses))
	}
}

func TestCategoryRepo_ListByHousehold(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	householdRepo := NewHouseholdRepo(db)
	categoryRepo := NewCategoryRepo(db)

	household := &domain.Household{Name: "Test Family"}
	householdRepo.Create(context.Background(), household)

	categories, err := categoryRepo.ListByHousehold(context.Background(), household.ID)
	if err != nil {
		t.Fatalf("ListByHousehold failed: %v", err)
	}

	// Should have predefined categories
	if len(categories) < 9 {
		t.Errorf("len(categories) = %d, want >= 9 (predefined)", len(categories))
	}
}

func TestCategoryRepo_Create_Custom(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	householdRepo := NewHouseholdRepo(db)
	categoryRepo := NewCategoryRepo(db)

	household := &domain.Household{Name: "Test Family"}
	householdRepo.Create(context.Background(), household)

	category := &domain.Category{
		HouseholdID:  &household.ID,
		Name:         "Mascotas",
		Icon:         strPtr("paw"),
		IsPredefined: false,
	}
	err := categoryRepo.Create(context.Background(), category)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if category.ID == 0 {
		t.Error("expected ID to be set")
	}
}

func strPtr(s string) *string {
	return &s
}
