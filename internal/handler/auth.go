package handler

import (
	"encoding/json"
	"net/http"

	"github.com/amartya321/go-code-hosting/internal/handler/service"
)

type AuthHandler struct {
	userSvc *service.UserService
	authSvc *service.AuthService
}

func NewAuthHandler(userSvc *service.UserService, authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{userSvc: userSvc, authSvc: authSvc}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	user, err := h.userSvc.Authenticate(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}
	token, err := h.authSvc.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})

}
