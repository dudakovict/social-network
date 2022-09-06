package post

import (
	"time"
	"unsafe"

	"github.com/dudakovict/social-network/business/core/post/db"
)

// Post represents an individual post.
type Post struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserID      string    `json:"user_id"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}

// NewPost contains information needed to create a new Post.
type NewPost struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	UserID      string `json:"user_id" validate:"required"`
}

// UpdatePost defines what information may be provided to modify an existing
// Post. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types but we make exceptions around
// marshalling/unmarshalling.
type UpdatePost struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

// =============================================================================

func toPost(dbP db.Post) Post {
	pu := (*Post)(unsafe.Pointer(&dbP))
	return *pu
}

func toPostSlice(dbPs []db.Post) []Post {
	posts := make([]Post, len(dbPs))
	for i, dbP := range dbPs {
		posts[i] = toPost(dbP)
	}
	return posts
}
