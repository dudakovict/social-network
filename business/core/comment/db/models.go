package db

import (
	"time"
)

// Comment represent the structure we need for moving data
// between the app and the database.
type Comment struct {
	ID          string    `db:"comment_id"`
	Description string    `db:"description"`
	PostID      string    `db:"post_id"`
	UserID      string    `db:"user_id"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}
