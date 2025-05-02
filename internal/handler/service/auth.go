package service

import (
	"time"

	"github.com/golang-jwt/jwt"
)

// AuthService handles JWT creation and validation
type AuthService struct {
	jwtSecret []byte
}

// NewAuthService makes an AuthService with the given secret
func NewAuthService(secret string) *AuthService {
	return &AuthService{jwtSecret: []byte(secret)}
}

// GenerateToken issues a JWT for the given userID, valid for 2 hours
func (s *AuthService) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(2 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
