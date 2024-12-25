package store

import (
	"context"
	"database/sql"
)

type UserStore struct {
	db *sql.DB
}

func (u *UserStore) Create(ctx context.Context, user *User) error {
	qry := `INSERT INTO users (username, password, email)
					VALUES ($1,$2,$3) RETURNING id, created_at`

	err := u.db.QueryRowContext(ctx, qry,
		user.ID,
		user.Username,
		user.Password,
		user.Email,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
