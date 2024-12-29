package store

import (
	"context"
	"database/sql"
)

type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, followerID, userID int64) error {
	qry := `insert into followers (user_id, follower_id) 
					values ($1,$2)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, qry, userID, followerID)
	return err
}

func (s *FollowerStore) Unfollow(ctx context.Context, followerID, userID int64) error {
	qry := `Delete from followers 
				WHERe user_id =$1 and follower_id =$2`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, qry, userID, followerID)
	return err
}
