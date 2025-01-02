package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"time"
)

type UserStore struct {
	db *sql.DB
}

func (u *UserStore) CreateForSeed(ctx context.Context, user *User) error {
	qry := `INSERT INTO users (username, password, email)
					VALUES ($1,$2,$3) RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := u.db.QueryRowContext(ctx, qry,
		user.Username,
		user.Password.hash,
		user.Email,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	qry := `INSERT INTO users (username, password, email)
					VALUES ($1,$2,$3) RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(ctx, qry,
		user.Username,
		user.Password.hash,
		user.Email,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate kay violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate kay violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
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

func (u *UserStore) CreateAndInvite(ctx context.Context, user *User, token string,
	invitationExp time.Duration) error {
	return withTx(u.db, ctx, func(tx *sql.Tx) error {
		//transaction wrapper
		//create user
		if err := u.Create(ctx, tx, user); err != nil {
			return err
		}
		//create user invite
		if err := u.createUserInvitation(ctx, tx, token, invitationExp, user.ID); err != nil {
			return err
		}
		return nil
	})

}
func (u *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration,
	userID int64) error {
	qry := `insert into user_invitations (token, user_id, expiry) values ($1, $2, $3);`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, qry, token, userID, time.Now().Add(exp))
	if err != nil {
		return err
	}
	return nil
}
func (u *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	qry := `select u.id, u.username, u.email, u.created_at, u.is_active 
			from users u  join user_invitations ui 
					on u.id = ui.user_id 
			where ui.token =$1 and ui.expiry > $2;`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	err := tx.QueryRowContext(ctx, qry, hashToken, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive,
	)
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
func (u *UserStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
	qry := `update users set username=$1, email=$2, is_active= $3
			where username=$4;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, qry, user.Username, user.Email, user.IsActive, user.Username)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) deleteUserInvitation(ctx context.Context, tx *sql.Tx, userID int64) error {
	qry := `delete from  user_invitations where user_id = $1;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, qry, userID)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) Activate(ctx context.Context, token string) error {
	return withTx(u.db, ctx, func(tx *sql.Tx) error {
		//find the user
		user, err := u.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		user.IsActive = true
		if err := u.update(ctx, tx, user); err != nil {
			return err
		}

		if err := u.deleteUserInvitation(ctx, tx, user.ID); err != nil {
			return err
		}
		return nil
	})
}
