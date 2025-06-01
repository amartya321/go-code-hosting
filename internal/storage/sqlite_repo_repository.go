package storage

import (
	"database/sql"

	"github.com/amartya321/go-code-hosting/internal/model"
	"github.com/amartya321/go-code-hosting/internal/storage/repository"
)

type sqliteRepoRepository struct {
	db *sql.DB
}

func NewSQLiteRepoRepository(db *sql.DB) repository.RepoRepository {
	return &sqliteRepoRepository{db: db}
}

func (r *sqliteRepoRepository) Create(repo *model.Repository) (*model.Repository, error) {
	result, error := r.db.Exec("INSERT INTO repositories (owner_id, name, description, is_private) VALUES (?, ?, ?, ?)", repo.OwnerID, repo.Name, repo.Description.String, repo.IsPrivate)
	if error != nil {
		return nil, error
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	repo.ID = int(id)
	return r.populateTimestamps(repo), nil

}

func (r *sqliteRepoRepository) GetByID(repoID int) (*model.Repository, error) {
	row := r.db.QueryRow(
		`SELECT id, owner_id, name, description, is_private, created_at, updated_at
           FROM repositories
          WHERE id = ?`, repoID,
	)
	var repo model.Repository
	err := row.Scan(&repo.ID, &repo.OwnerID, &repo.Name, &repo.Description, &repo.IsPrivate, &repo.CreatedAt, &repo.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No repository found
		}
		return nil, err // Other error
	}
	return &repo, nil

}

func (r *sqliteRepoRepository) ListByOwner(ownerID int) ([]*model.Repository, error) {
	rows, err := r.db.Query(
		`SELECT id, owner_id, name, description, is_private, created_at, updated_at
		   FROM repositories
		  WHERE owner_id = ?`, ownerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repos []*model.Repository
	for rows.Next() {
		var repo model.Repository
		if err := rows.Scan(&repo.ID, &repo.OwnerID, &repo.Name, &repo.Description, &repo.IsPrivate, &repo.CreatedAt, &repo.UpdatedAt); err != nil {
			return nil, err
		}
		repos = append(repos, &repo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return repos, nil
}

func (r *sqliteRepoRepository) Update(repo *model.Repository) (*model.Repository, error) {
	_, err := r.db.Exec("UPDATE repositories SET name = ?, description = ?, is_private = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", repo.Name, repo.Description.String, repo.IsPrivate, repo.ID)
	if err != nil {
		return nil, err
	}
	return r.populateTimestamps(repo), nil
}

func (r *sqliteRepoRepository) Delete(repoID int) error {
	result, err := r.db.Exec(`DELETE FROM repositories WHERE id = ?`, repoID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *sqliteRepoRepository) populateTimestamps(repo *model.Repository) *model.Repository {
	row := r.db.QueryRow("SELECT created_at, updated_at  FROM repositories WHERE id = ?", repo.ID)
	row.Scan(&repo.CreatedAt, &repo.UpdatedAt)

	return repo
}
