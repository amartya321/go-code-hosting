package handler

import (
	"encoding/json"
	"net/http"

	"github.com/amartya321/go-code-hosting/internal/handler/service"
	"github.com/amartya321/go-code-hosting/internal/middleware"
	"github.com/go-chi/chi/v5"
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

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "In URI Id is not valid", http.StatusBadRequest)
		return
	}
	user, err := h.svc.GetUserByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	//first need to check if the token Id and the user Id matches or not
	authId, ok := middleware.UserIDFromContext(r.Context())
	if !ok || authId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	if id != authId {
		http.Error(w, "Unauthorized to update this user", http.StatusForbidden)
		return
	}

	var input service.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	updatedUser, err := h.svc.UpdateUser(id, input)
	if err != nil {
		http.Error(w, "Could not update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)

}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	authId, ok := middleware.UserIDFromContext(r.Context())
	if !ok || authId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	if id != authId {
		http.Error(w, "Unauthorized to update this user", http.StatusForbidden)
		return
	}

	if err := h.svc.DeleteUser(id); err != nil {
		http.Error(w, "Could not delete user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
