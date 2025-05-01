package handler

import (
	"net/http"

	"github.com/amartya321/go-code-hosting/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r *chi.Mux, userHandler *UserHandler) {

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Go Code Hosting Service 🚀"))
	})

	r.Route("/users", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware) // 👈 protect /users with auth

		r.Post("/", userHandler.CreateUser)
		r.Get("/", userHandler.ListUsers)
	})

}
