package storage

import "github.com/amartya321/go-code-hosting/internal/model"

type UserRepository interface {
	Create(user model.User) error
	List() []model.User
}

type InMemoryUserRepository struct {
	user []model.User
}

//Kind of like a constructor for the InMemoryUserRepository struct
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{user: []model.User{}}
}

func (s *InMemoryUserRepository) Create(user model.User) error {
	s.user = append(s.user, user)
	return nil
}

func (s *InMemoryUserRepository) List() []model.User {
	return s.user
}
