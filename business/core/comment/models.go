package comment

import (
	"time"
	"unsafe"

	"github.com/dudakovict/social-network/business/core/comment/db"
)

// Comment represents an individual comment.
type Comment struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	PostID      string    `json:"post_id"`
	UserID      string    `json:"user_id"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}

// NewComment contains information needed to create a new Comment.
type NewComment struct {
	Description string `json:"description" validate:"required"`
	PostID      string `json:"post_id" validate:"required"`
	UserID      string `json:"user_id" validate:"required"`
}

// UpdateComment defines what information may be provided to modify an existing
// Comment. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types but we make exceptions around
// marshalling/unmarshalling.
type UpdateComment struct {
	Description *string `json:"description"`
}

// =============================================================================

func toComment(dbC db.Comment) Comment {
	cu := (*Comment)(unsafe.Pointer(&dbC))
	return *cu
}

func toCommentSlice(dbCs []db.Comment) []Comment {
	comments := make([]Comment, len(dbCs))
	for i, dbC := range dbCs {
		comments[i] = toComment(dbC)
	}
	return comments
}
