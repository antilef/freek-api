package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	calcHandler "freek/internal/calculation/handler"
	calcRepoImpl "freek/internal/calculation/repository/impl/postgres"
	calcService "freek/internal/calculation/service"
	ingHandler "freek/internal/ingestion/handler"
	ingRepoImpl "freek/internal/ingestion/repository/impl/postgres"
	ingService "freek/internal/ingestion/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://localhost:5432/freek?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("Connected to database")

	migrationsPath := "file://migrations"
	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		log.Fatal("Failed to create migrate instance:", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to run migrations:", err)
	}
	log.Println("Migrations applied successfully")

	// Repositories
	householdRepo := ingRepoImpl.NewHouseholdRepo(db)
	memberRepo := ingRepoImpl.NewMemberRepo(db)
	cardRepo := ingRepoImpl.NewCardRepo(db)
	expenseRepo := ingRepoImpl.NewExpenseRepo(db)
	categoryRepo := ingRepoImpl.NewCategoryRepo(db)
	expenseReadRepo := calcRepoImpl.NewExpenseReadRepo(db)

	// Services
	householdService := ingService.NewHouseholdService(householdRepo)
	memberService := ingService.NewMemberService(memberRepo)
	cardService := ingService.NewCardService(cardRepo)
	expenseService := ingService.NewExpenseService(expenseRepo)
	categoryService := ingService.NewCategoryService(categoryRepo)
	projectionService := calcService.NewProjectionService(expenseReadRepo, nil)
	historyService := calcService.NewHistoryService(expenseReadRepo)
	summaryService := calcService.NewSummaryService(expenseReadRepo)

	// Handlers
	householdHandler := ingHandler.NewHouseholdHandler(householdService, memberService)
	cardHandler := ingHandler.NewCardHandler(cardService)
	expenseHandler := ingHandler.NewExpenseHandler(expenseService)
	categoryHandler := ingHandler.NewCategoryHandler(categoryService)
	projectionHandler := calcHandler.NewProjectionHandler(projectionService)
	historyHandler := calcHandler.NewHistoryHandler(historyService)
	summaryHandler := calcHandler.NewSummaryHandler(summaryService)

	// Router
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods("GET")

	// Household routes
	api.HandleFunc("/households", householdHandler.Create).Methods("POST")
	api.HandleFunc("/households/{id}", householdHandler.GetByID).Methods("GET")
	api.HandleFunc("/households/{id}/members", householdHandler.AddMember).Methods("POST")
	api.HandleFunc("/households/{id}/members", householdHandler.ListMembers).Methods("GET")
	api.HandleFunc("/households/{id}/members/{member_id}", householdHandler.RemoveMember).Methods("DELETE")

	// Credit card routes
	api.HandleFunc("/households/{id}/cards", cardHandler.Create).Methods("POST")
	api.HandleFunc("/households/{id}/cards", cardHandler.List).Methods("GET")
	api.HandleFunc("/households/{id}/cards/{card_id}", cardHandler.GetByID).Methods("GET")
	api.HandleFunc("/households/{id}/cards/{card_id}", cardHandler.Update).Methods("PUT")
	api.HandleFunc("/households/{id}/cards/{card_id}", cardHandler.Delete).Methods("DELETE")

	// Expense routes
	api.HandleFunc("/households/{id}/expenses", expenseHandler.Create).Methods("POST")
	api.HandleFunc("/households/{id}/expenses", expenseHandler.List).Methods("GET")
	api.HandleFunc("/households/{id}/expenses/{expense_id}", expenseHandler.Update).Methods("PUT")
	api.HandleFunc("/households/{id}/expenses/{expense_id}", expenseHandler.Delete).Methods("DELETE")

	// Category routes
	api.HandleFunc("/households/{id}/categories", categoryHandler.List).Methods("GET")
	api.HandleFunc("/households/{id}/categories", categoryHandler.Create).Methods("POST")
	api.HandleFunc("/households/{id}/categories/{category_id}", categoryHandler.Update).Methods("PUT")
	api.HandleFunc("/households/{id}/categories/{category_id}", categoryHandler.Delete).Methods("DELETE")

	// Calculation routes
	api.HandleFunc("/households/{id}/projections", projectionHandler.Project).Methods("GET")
	api.HandleFunc("/households/{id}/history", historyHandler.Query).Methods("GET")
	api.HandleFunc("/households/{id}/history/summary", summaryHandler.Summary).Methods("GET")
	api.HandleFunc("/households/{id}/history/compare", summaryHandler.Compare).Methods("GET")

	log.Printf("Freek API starting on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
