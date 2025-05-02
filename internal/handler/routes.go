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

	r.Post("/login", authHandler.Login)

	r.Post("/users", userHandler.CreateUser)

	r.With(middleware.JWTMiddleware).
		Get("/users", userHandler.ListUsers)

}
