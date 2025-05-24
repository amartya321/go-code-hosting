package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

type ctxKey string

const ctxKeyUserID ctxKey = "userID"

func JWTMiddleware(next http.Handler) http.Handler {
	secret := []byte(os.Getenv("JWT_SECRET"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Extract the “sub” claim
		claims := token.Claims.(*jwt.MapClaims)
		sub, ok := claims["sub"].(string)
		if !ok || sub == "" {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Put it into context for handlers downstream
		ctx := context.WithValue(r.Context(), ctxKeyUserID, sub)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
