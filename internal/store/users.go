package store

import (
	"context"
	"database/sql"
)

type User struct {
    ID        int64  `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Username  string `json:"username"`
    Email     string `json:"email"`
    Password  []byte `json:"-"`
    CreatedAt string `json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func (u *UserStore) Create(ctx context.Context, user *User) error {
	query := `
	INSERT INTO users (first_name, last_name, username, password, email)
	VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at
	`
	if err := u.db.QueryRowContext(ctx, query, user.FirstName, user.LastName, user.Username, user.Password, user.Email).Scan(
		&user.ID,
		&user.CreatedAt,
	); err != nil {
		return err
	}

	return nil
}
