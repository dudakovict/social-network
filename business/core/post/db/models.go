package db

import (
	"time"
)

// Post represent the structure we need for moving data
// between the app and the database.
type Post struct {
	ID          string    `db:"post_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}
