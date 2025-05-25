package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/amartya321/go-code-hosting/internal/handler"
	"github.com/amartya321/go-code-hosting/internal/handler/service"
	"github.com/amartya321/go-code-hosting/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"

	"github.com/golang-migrate/migrate/v4/database/sqlite"
)

func main() {

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable must be set")
	}

	// Open or create the SQLite database
	dbPath := "users.db"
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Apply database migrations
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		log.Fatalf("migrate driver error: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"sqlite3", driver,
	)
	if err != nil {
		log.Fatalf("failed to initialize migrate: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("database migration failed: %v", err)
	}
	log.Println("✅ Database migrations applied")

	// Initialize storage, services, and handlers
	store := storage.NewSQLiteUserRepositoryFromDB(db)

	r := chi.NewRouter()

	// Create the in-memory store and user service and user handler
	//	store := storage.NewInMemoryUserRepository()
	userService := service.NewUserService(store)
	authService := service.NewAuthService(jwtSecret)
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(userService, authService)

	// Register the handler
	handler.RegisterRoutes(r, userHandler, authHandler)

	// Start the HTTP server
	log.Println("Starting server on :8080...")
	serverErr := http.ListenAndServe(":8080", r)
	if serverErr != nil {
		log.Fatal(serverErr)
	}
}
