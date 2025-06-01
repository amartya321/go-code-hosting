package repository

import "github.com/amartya321/go-code-hosting/internal/model"

type RepoRepository interface {
	Create(repo *model.Repository) (*model.Repository, error)
	GetByID(repoID int) (*model.Repository, error)
	ListByOwner(ownerID int) ([]*model.Repository, error)
	Update(repo *model.Repository) (*model.Repository, error)
	Delete(repoID int) error
}
