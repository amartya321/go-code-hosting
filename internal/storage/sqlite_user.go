package storage

import (
	"database/sql"

	"github.com/amartya321/go-code-hosting/internal/model"
	_ "modernc.org/sqlite" // pure Go SQLite driver
)

type SQLiteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(dbPath string) (*SQLiteUserRepository, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Create users table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL,
			email TEXT NOT NULL,
			password_hash TEXT NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	return &SQLiteUserRepository{db: db}, nil
}

func (r *SQLiteUserRepository) Create(user model.User) error {
	_, err := r.db.Exec(`INSERT INTO users (id, username, email) VALUES (?, ?, ?)`,
		user.ID, user.Username, user.Email)
	return err
}

func (r *SQLiteUserRepository) List() []model.User {
	rows, err := r.db.Query(`SELECT id, username, email FROM users`)
	if err != nil {
		return []model.User{}
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err == nil {
			users = append(users, user)
		}
	}
	return users
}

func (s *SQLiteUserRepository) FindByUsername(username string) (*model.User, error) {
	row := s.db.QueryRow(
		`SELECT id, username, email, password_hash 
		FROM users WHERE username = ?`, username,
	)
	var u model.User
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
