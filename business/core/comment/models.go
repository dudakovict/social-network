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
	UserID      string    `json:"user_id"`
	PostID      string    `json:"post_id"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}

type Post struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserID      string    `json:"user_id"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}

type PostComment struct {
	ID                 string    `json:"id"`
	Title              string    `json:"title"`
	Description        string    `json:"description"`
	UserID             string    `json:"user_id"`
	DateCreated        time.Time `json:"date_created"`
	DateUpdated        time.Time `json:"date_updated"`
	CommentID          string    `json:"comment_id"`
	CommentDescription string    `json:"comment_description"`
	CommentUserID      string    `json:"comment_user_id"`
	CommentDateCreated time.Time `json:"comment_date_created"`
	CommentDateUpdated time.Time `json:"comment_date_updated"`
}

// NewComment contains information needed to create a new Comment.
type NewComment struct {
	Description string `json:"description" validate:"required"`
	UserID      string `json:"user_id" validate:"required"`
	PostID      string `json:"post_id" validate:"required"`
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

func toPostComment(dbPcomm db.PostComment) PostComment {
	pcu := (*PostComment)(unsafe.Pointer(&dbPcomm))
	return *pcu
}

func toPostCommentSlice(dbPsComms []db.PostComment) []PostComment {
	pscomms := make([]PostComment, len(dbPsComms))
	for i, dbPcomm := range dbPsComms {
		pscomms[i] = toPostComment(dbPcomm)
	}
	return pscomms
}
