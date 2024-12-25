package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	qry := `INSERT INTO posts (content, title, user_id, tags)
					VALUES ($1,$2,$3,$4,) RETURNING id, created_at, updated_at`

	err := s.db.QueryRowContext(ctx, qry,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}