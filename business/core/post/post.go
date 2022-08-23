// Package post provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package post

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dudakovict/social-network/business/core/post/db"
	"github.com/dudakovict/social-network/business/sys/database"
	"github.com/dudakovict/social-network/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("post not found")
	ErrInvalidID             = errors.New("ID is not in its proper form")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

// Core manages the set of API's for post access.
type Core struct {
	store db.Store
}

// NewCore constructs a core for post api access.
func NewCore(log *zap.SugaredLogger, sqlxDB *sqlx.DB) Core {
	return Core{
		store: db.NewStore(log, sqlxDB),
	}
}

// Create inserts a new post into the database.
func (c Core) Create(ctx context.Context, np NewPost, now time.Time) (Post, error) {
	if err := validate.Check(np); err != nil {
		return Post{}, fmt.Errorf("validating data: %w", err)
	}

	dbP := db.Post{
		ID:          validate.GenerateID(),
		Title:       np.Title,
		Description: np.Description,
		DateCreated: now,
		DateUpdated: now,
	}

	// This provides an example of how to execute a transaction if required.
	tran := func(tx sqlx.ExtContext) error {
		if err := c.store.Tran(tx).Create(ctx, dbP); err != nil {
			return fmt.Errorf("create: %w", err)
		}
		return nil
	}

	if err := c.store.WithinTran(ctx, tran); err != nil {
		return Post{}, fmt.Errorf("tran: %w", err)
	}

	// if err := c.store.Create(ctx, dbP); err != nil {
	// 	return Post{}, fmt.Errorf("create: %w", err)
	// }

	return toPost(dbP), nil
}

// Update replaces a post document in the database.
func (c Core) Update(ctx context.Context, postID string, up UpdatePost, now time.Time) error {
	if err := validate.CheckID(postID); err != nil {
		return ErrInvalidID
	}

	if err := validate.Check(up); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	dbP, err := c.store.QueryByID(ctx, postID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating post postID[%s]: %w", postID, err)
	}

	if up.Title != nil {
		dbP.Title = *up.Title
	}
	if up.Description != nil {
		dbP.Description = *up.Description
	}
	dbP.DateUpdated = now

	if err := c.store.Update(ctx, dbP); err != nil {
		return fmt.Errorf("udpate: %w", err)
	}

	return nil
}

// Delete removes a post from the database.
func (c Core) Delete(ctx context.Context, postID string) error {
	if err := validate.CheckID(postID); err != nil {
		return ErrInvalidID
	}

	if err := c.store.Delete(ctx, postID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing posts from the database.
func (c Core) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]Post, error) {
	dbPosts, err := c.store.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query: %w", err)
	}

	return toPostSlice(dbPosts), nil
}

// QueryByID gets the specified post from the database.
func (c Core) QueryByID(ctx context.Context, postID string) (Post, error) {
	if err := validate.CheckID(postID); err != nil {
		return Post{}, ErrInvalidID
	}

	dbP, err := c.store.QueryByID(ctx, postID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Post{}, ErrNotFound
		}
		return Post{}, fmt.Errorf("query: %w", err)
	}

	return toPost(dbP), nil
}
