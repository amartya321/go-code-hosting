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

	r.Route("/users", func(r chi.Router) {
		// Public: create
		r.Post("/", userHandler.CreateUser)

		// Protected: anything below
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTMiddleware)

			r.Get("/", userHandler.ListUsers)
			r.Get("/{id}", userHandler.GetUserByID)
			// r.Put("/{id}",userHandler.UpdateUser)
			// r.Delete("/{id}", userHandler.DeleteUser)
		})
	})

}
