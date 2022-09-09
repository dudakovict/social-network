// Package comment provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package comment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dudakovict/social-network/business/core/comment/db"
	"github.com/dudakovict/social-network/business/sys/database"
	"github.com/dudakovict/social-network/business/sys/nats"
	"github.com/dudakovict/social-network/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("comment not found")
	ErrInvalidID             = errors.New("ID is not in its proper form")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

// Core manages the set of API's for comment access.
type Core struct {
	store db.Store
	n     *nats.NATS
}

// NewCore constructs a core for comment api access.
func NewCore(log *zap.SugaredLogger, sqlxDB *sqlx.DB, n *nats.NATS) Core {
	c := Core{
		store: db.NewStore(log, sqlxDB),
		n:     n,
	}

	/*
		err := n.Subscribe("post-created", func(m *stan.Msg) {
			buf := bytes.NewReader(m.Data)
			dec := gob.NewDecoder(buf)

			var dbP db.Post

			err := dec.Decode(&dbP)
			if err != nil {
				fmt.Errorf("decoding: %w", err)
			}

			log.Infof("[%+v]:", dbP)

			if err := c.store.CreatePost(context.Background(), dbP); err != nil {
				fmt.Errorf("create: %w", err)
			}
		})
	*/
	l := Listener{
		log:    log,
		store:  c.store,
		client: c.n.Client,
	}

	l.PostCreated()
	l.PostUpdated()
	l.PostDeleted()

	return c
}

// Create inserts a new comment into the database.
func (c Core) Create(ctx context.Context, nc NewComment, now time.Time) (Comment, error) {
	if err := validate.Check(nc); err != nil {
		return Comment{}, fmt.Errorf("validating data: %w", err)
	}

	dbC := db.Comment{
		ID:          validate.GenerateID(),
		Description: nc.Description,
		PostID:      nc.PostID,
		UserID:      nc.UserID,
		DateCreated: now,
		DateUpdated: now,
	}

	// This provides an example of how to execute a transaction if required.
	tran := func(tx sqlx.ExtContext) error {
		if err := c.store.Tran(tx).Create(ctx, dbC); err != nil {
			return fmt.Errorf("create: %w", err)
		}
		return nil
	}

	if err := c.store.WithinTran(ctx, tran); err != nil {
		return Comment{}, fmt.Errorf("tran: %w", err)
	}

	// if err := c.store.Create(ctx, dbP); err != nil {
	// 	return Post{}, fmt.Errorf("create: %w", err)
	// }

	return toComment(dbC), nil
}

// Update replaces a comment document in the database.
func (c Core) Update(ctx context.Context, commentID string, uc UpdateComment, now time.Time) error {
	if err := validate.CheckID(commentID); err != nil {
		return ErrInvalidID
	}

	if err := validate.Check(uc); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	dbC, err := c.store.QueryByID(ctx, commentID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating comment commentID[%s]: %w", commentID, err)
	}

	if uc.Description != nil {
		dbC.Description = *uc.Description
	}
	dbC.DateUpdated = now

	if err := c.store.Update(ctx, dbC); err != nil {
		return fmt.Errorf("udpate: %w", err)
	}

	return nil
}

// Delete removes a comment from the database.
func (c Core) Delete(ctx context.Context, commentID string) error {
	if err := validate.CheckID(commentID); err != nil {
		return ErrInvalidID
	}

	if err := c.store.Delete(ctx, commentID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing comments from the database.
func (c Core) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]Comment, error) {
	dbComments, err := c.store.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query: %w", err)
	}

	return toCommentSlice(dbComments), nil
}

// QueryByID gets the specified comment from the database.
func (c Core) QueryByID(ctx context.Context, commentID string) (Comment, error) {
	if err := validate.CheckID(commentID); err != nil {
		return Comment{}, ErrInvalidID
	}

	dbC, err := c.store.QueryByID(ctx, commentID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Comment{}, ErrNotFound
		}
		return Comment{}, fmt.Errorf("query: %w", err)
	}

	return toComment(dbC), nil
}

func (c Core) QueryByUserID(ctx context.Context, userID string) ([]Comment, error) {
	if err := validate.CheckID(userID); err != nil {
		return nil, ErrInvalidID
	}

	dbComments, err := c.store.QueryByUserID(ctx, userID)

	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return toCommentSlice(dbComments), nil
}

func (c Core) QueryByPostID(ctx context.Context, postID string) ([]Comment, error) {
	if err := validate.CheckID(postID); err != nil {
		return nil, ErrInvalidID
	}

	dbComments, err := c.store.QueryByPostID(ctx, postID)

	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return toCommentSlice(dbComments), nil
}

func (c Core) QueryPostsByPostID(ctx context.Context, postID string) (interface{}, error) {
	if err := validate.CheckID(postID); err != nil {
		return nil, ErrInvalidID
	}

	dbPosts, err := c.store.QueryPostsByPostID(ctx, postID)

	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return toPostCommentSlice(dbPosts), nil
}
