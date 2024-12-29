package store

import (
	"context"
	"database/sql"
)

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) Create(ctx context.Context, comment *Comments) error {
	query := `
		INSERT INTO comments (post_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *CommentStore) GetByPostID(ctx context.Context, postID int64) ([]Comments, error) {
	qry := `select c.id ,c.post_id ,c.content, c.created_at ,u.username ,u.id 
				from comments c join users u 
				on u.id = c.user_id 
				where c.post_id = $1
				order by c.created_at  desc;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, qry, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comments{}

	for rows.Next() {
		var c Comments
		c.User = User{}
		err := rows.Scan(&c.ID, &c.PostID, &c.Content, &c.CreatedAt,
			&c.User.Username, &c.User.ID)
		if err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}
	return comments, nil
}
