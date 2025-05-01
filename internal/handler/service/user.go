package service

import (
	"github.com/amartya321/go-code-hosting/internal/model"
	"github.com/amartya321/go-code-hosting/internal/storage"
)

type UserService struct {
	repo storage.UserRepository
}

func NewUserService(repo storage.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(username, email string) (model.User, error) {
	user := model.NewUser(username, email)
	err := s.repo.Create(user)
	return user, err
}

func (s *UserService) ListUsers() []model.User {
	return s.repo.List()
}
