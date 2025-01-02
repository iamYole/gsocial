package store

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Post struct {
	ID        int64      `json:"id"`
	Content   string     `json:"content"`
	Title     string     `json:"title"`
	UserID    int64      `json:"user_id"`
	Tags      []string   `json:"tags"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Version   int        `json:"version"`
	Comments  []Comments `json:"comments"`
	User      User       `json:"user"`
}

type Password struct {
	pword *string
	hash  []byte
}

func (p *Password) Set(pword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pword), bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	p.pword = &pword
	p.hash = hash

	return nil
}

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  Password  `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}

type Comments struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `json:"user"`
}

type Follower struct {
	UserID     int64     `json:"user_id"`
	FollowerID int64     `json:"follower_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type FeedItem struct {
	Post
	CommentCount int `json:"comments_count"`
}
