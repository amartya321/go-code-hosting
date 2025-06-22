package service

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"

	"regexp"

	"github.com/amartya321/go-code-hosting/internal/model"
	"github.com/amartya321/go-code-hosting/internal/storage/repository"
)

type RepoService struct {
	repoRepo repository.RepoRepository
	repoRoot string // the filesystem path where bare repos live, e.g. "./repos"
}

func NewRepoService(repoRepo repository.RepoRepository, repoRoot string) *RepoService {
	return &RepoService{
		repoRepo: repoRepo,
		repoRoot: repoRoot,
	}
}

// Predefined errors for RepoService.
var (
	// ErrRepoNotFound is returned when a repository does not exist.
	ErrRepoNotFound = errors.New("repository not found")

	// ErrRepoNameTaken is returned when a user already has a repository with the given name.
	ErrRepoNameTaken = errors.New("repository name already taken")

	// ErrForbidden is returned when a user tries to access or modify a repository they do not own.
	ErrForbidden = errors.New("forbidden")

	// ErrInvalidName is returned when a repository name does not match validation rules.
	ErrInvalidName = errors.New("invalid repository name")
)

// validateName ensures that the repository name is non-empty and contains only allowed characters.
// Here we allow letters, numbers, dashes, underscores, and dots. Adjust the regex as needed.
func validateName(name string) error {
	if name == "" {
		return ErrInvalidName
	}
	matched, _ := regexp.MatchString(`^[A-Za-z0-9._-]+$`, name)
	if !matched {
		return ErrInvalidName
	}
	return nil
}

func (s *RepoService) Create(ownerID int, name string, description string, isPrivate bool) (*model.Repository, error) {
	if err := validateName(name); err != nil {
		return nil, err
	}

	// 2. Ensure uniqueness: no existing repo with same name for this owner
	existingRepo, err := s.repoRepo.ListByOwner(ownerID)
	if err != nil {
		return nil, err
	}

	for _, repo := range existingRepo {
		if repo.Name == name {
			return nil, ErrRepoNameTaken
		}
	}

	repo := &model.Repository{
		OwnerID:     ownerID,
		Name:        name,
		Description: sql.NullString{String: description, Valid: description != ""},
		IsPrivate:   isPrivate,
	}

	created, err := s.repoRepo.Create(repo)
	if err != nil {
		return nil, fmt.Errorf("insert repo into database: %w", err)
	}

	onDisk := filepath.Join(s.repoRoot, fmt.Sprintf("%d_%s.git", ownerID, name))
	if err := os.MkdirAll(s.repoRoot, 0o755); err != nil {
		_ = s.repoRepo.Delete(created.ID)
		return nil, fmt.Errorf("create repos root directory: %w", err)
	}

	if _, err := git.PlainInit(onDisk, true); err != nil {
		_ = s.repoRepo.Delete(created.ID)
		return nil, fmt.Errorf("initialize bare git repo: %w", err)
	}

	return created, nil

}
