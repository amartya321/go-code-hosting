package handler

import (
	"net/http"

	"github.com/amartya321/go-code-hosting/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r *chi.Mux, userHandler *UserHandler, authHandler *AuthHandler) {

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Go Code Hosting Service 🚀"))
	})

	// Public: login to receive a JWT
	r.Post("/login", authHandler.Login)

	// Protected: any /users route now requires a valid JWT
	r.Route("/users", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)

		r.Post("/", userHandler.CreateUser)
		r.Get("/", userHandler.ListUsers)
	})

}
