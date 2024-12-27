package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	qry := `INSERT INTO posts (content, title, user_id, tags)
					VALUES ($1,$2,$3,$4) RETURNING id, created_at, updated_at`

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

func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	qry := `select p.id, p.user_id, p.title, p.content, p.tags, p.created_at, p.updated_at  
			from posts p 
			where p.id = $1`
	var post Post
	err := s.db.QueryRowContext(ctx, qry, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}
