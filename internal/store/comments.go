package store

import (
	"context"
	"database/sql"
)

type CommentStore struct {
	db *sql.DB
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
