package db

import (
	"time"
)

// Comment represent the structure we need for moving data
// between the app and the database.
type Comment struct {
	ID          string    `db:"comment_id"`
	Description string    `db:"description"`
	UserID      string    `db:"user_id"`
	PostID      string    `db:"post_id"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

type Post struct {
	ID          string    `db:"post_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	UserID      string    `db:"user_id"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

type PostComment struct {
	ID                 string    `db:"post_id"`
	Title              string    `db:"title"`
	Description        string    `db:"description"`
	UserID             string    `db:"user_id"`
	DateCreated        time.Time `db:"date_created"`
	DateUpdated        time.Time `db:"date_updated"`
	CommentID          string    `db:"comment_id"`
	CommentDescription string    `db:"comment_description"`
	CommentUserID      string    `db:"comment_user_id"`
	CommentDateCreated time.Time `db:"comment_date_created"`
	CommentDateUpdated time.Time `db:"comment_date_updated"`
}
