package storage

import "github.com/amartya321/go-code-hosting/internal/model"

type UserRepository interface {
	Create(users model.User) error
	List() []model.User
	FindByUserName(username string) (*model.User, error)
	FindByUserId(username string) (*model.User, error)
}

type InMemoryUserRepository struct {
	users []model.User // instead of `users []model.User`

}

//Kind of like a constructor for the InMemoryUserRepository struct
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{users: []model.User{}}
}

func (s *InMemoryUserRepository) Create(users model.User) error {
	s.users = append(s.users, users)
	return nil
}

func (s *InMemoryUserRepository) List() []model.User {
	return s.users
}

func (s *InMemoryUserRepository) FindByUserName(username string) (*model.User, error) {
	for _, user := range s.users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, nil
}
