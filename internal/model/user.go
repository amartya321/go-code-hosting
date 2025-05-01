package model

import "github.com/google/uuid"

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

func NewUser(userName, email string) User {
	return User{
		ID:       uuid.NewString(),
		Username: userName,
		Email:    email,
	}
}
