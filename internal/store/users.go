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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := u.db.QueryRowContext(ctx, qry,
		user.Username,
		user.Password,
		user.Email,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) GetByID(ctx context.Context, userID int64) (*User, error) {
	qry := `select u.id, u.username, u.email, u.password, u.created_at
			from users u 
			where u.id =$1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	err := u.db.QueryRowContext(ctx, qry, userID).Scan(&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}
