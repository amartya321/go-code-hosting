package service

import (
	"fmt"
	"log"

	"github.com/amartya321/go-code-hosting/internal/model"
	"github.com/amartya321/go-code-hosting/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo storage.UserRepository
}

func NewUserService(repo storage.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(username, email, passwrod string) (model.User, error) {
	user := model.NewUser(username, email)

	hash, err := bcrypt.GenerateFromPassword([]byte(passwrod), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("UserService.CreateUser(%q) password Generation failed: %v", username, err)
		return user, err

	}

	user.PasswordHash = string(hash)

	if err := s.repo.Create(user); err != nil {
		log.Printf("UserService.CreateUser(%q) Create failed: %v", username, err)
		return user, err
	}

	return user, nil
}

func (s *UserService) ListUsers() []model.User {
	return s.repo.List()
}

// Authenticate checks the username/password and returns the user if valid
func (s *UserService) Authenticate(username, passwrod string) (*model.User, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		log.Printf("UserService.Authenticate(%q) FindByUsername failed: %v", username, err)
		return nil, err
	}
	if user == nil {
		log.Printf("UserService.Authenticate(%q) user not found", username)
		return nil, fmt.Errorf("user not found")
	}

	// 2) Compare the provided password against the stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(passwrod)); err != nil {
		log.Printf("UserService.Authenticate(%q) password comparison failed: %v", username, err)
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil

}
