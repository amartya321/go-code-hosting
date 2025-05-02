package main

import (
	"log"
	"net/http"

	"github.com/amartya321/go-code-hosting/internal/handler"
	"github.com/amartya321/go-code-hosting/internal/handler/service"
	"github.com/amartya321/go-code-hosting/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {

	store, err := storage.NewSQLiteUserRepository("users.db")
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	// Create the in-memory store and user service and user handler
	//	store := storage.NewInMemoryUserRepository()
	userService := service.NewUserService(store)
	authService := service.NewAuthService()
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
