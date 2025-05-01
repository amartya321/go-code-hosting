package handler

import (
	"encoding/json"
	"net/http"

	"github.com/amartya321/go-code-hosting/internal/handler/service"
)

// func handleCreateUser(w http.ResponseWriter, r *http.Request) {
// 	var input struct {
// 		Username string `json:"username"`
// 		Email    string `json:"email"`
// 	}
// 	err := json.NewDecoder(r.Body).Decode(&input)
// 	if err != nil {
// 		http.Error(w, "invalid input", http.StatusBadRequest)
// 		return
// 	}
// 	user := model.NewUser(input.Username, input.Email)
// 	users = append(users, user)

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(user)
// }

// func handleListUsers(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(users)
// }

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.svc.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users := h.svc.ListUsers()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
