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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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
	qry := `select p.id, p.user_id, p.title, p.content, p.tags, p.created_at, p.updated_at, version  
			from posts p 
			where p.id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var post Post
	err := s.db.QueryRowContext(ctx, qry, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
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
func (s *PostStore) Delete(ctx context.Context, postID int64) error {
	qry := `Delete from posts where id=$1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	respone, err := s.db.ExecContext(ctx, qry, postID)
	if err != nil {
		return err
	}

	rows, err := respone.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrorNotFound
	}
	return nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	qry := `update posts 
			SET	title = $1, content = $2, version = version + 1
			where id = $3 AND version=$4 RETURNING version`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, qry,
		post.Title,
		post.Content,
		post.ID,
		post.Version).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrorNotFound
		default:
			return err
		}
	}
	return nil
}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]FeedItem, error) {
	qry := `select p.id , p.user_id , p.title , p."content" , p.created_at ,p."version", p.tags , u.username ,COUNT(c.id) as comments_count
			from posts p 
			left join comments c on c.post_id  = p.id 
			left join users u on p.user_id =u.id 
			join followers f on f.follower_id =p.user_id or p.user_id = $1
			where
			 f.user_id = $1 or p.user_id =$1 AND
			 (p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%') AND
			(p.tags @> $5 OR $5 = '{}')
			group by p.id, u.username 
			order by p.created_at ` + fq.Sort + `
			limit $2
			offset $3;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, qry, userID, fq.Limit, fq.Offset, fq.Search, pq.Array(fq.Tags))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feed []FeedItem
	for rows.Next() {
		var f FeedItem
		err := rows.Scan(
			&f.ID,
			&f.UserID,
			&f.Title,
			&f.Content,
			&f.CreatedAt,
			&f.Version,
			pq.Array(&f.Tags),
			&f.User.Username,
			&f.CommentCount,
		)
		if err != nil {
			return nil, err
		}

		feed = append(feed, f)
	}
	return feed, nil
}
