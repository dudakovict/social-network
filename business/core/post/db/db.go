// Package db contains post related CRUD functionality.
package db

import (
	"context"
	"fmt"

	"github.com/dudakovict/social-network/business/sys/database"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Store manages the set of API's for post access.
type Store struct {
	log          *zap.SugaredLogger
	tr           database.Transactor
	db           sqlx.ExtContext
	isWithinTran bool
}

// NewStore constructs a data for api access.
func NewStore(log *zap.SugaredLogger, db *sqlx.DB) Store {
	return Store{
		log: log,
		tr:  db,
		db:  db,
	}
}

// WithinTran runs passed function and do commit/rollback at the end.
func (s Store) WithinTran(ctx context.Context, fn func(sqlx.ExtContext) error) error {
	if s.isWithinTran {
		return fn(s.db)
	}
	return database.WithinTran(ctx, s.log, s.tr, fn)
}

// Tran return new Store with transaction in it.
func (s Store) Tran(tx sqlx.ExtContext) Store {
	return Store{
		log:          s.log,
		tr:           s.tr,
		db:           tx,
		isWithinTran: true,
	}
}

// Create inserts a new post into the database.
func (s Store) Create(ctx context.Context, p Post) error {
	const q = `
	INSERT INTO posts
		(post_id, title, description, date_created, date_updated)
	VALUES
		(:post_id, :title, :description, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, p); err != nil {
		return fmt.Errorf("inserting post: %w", err)
	}

	return nil
}

// Update replaces a post document in the database.
func (s Store) Update(ctx context.Context, p Post) error {
	const q = `
	UPDATE
		posts
	SET 
		"title" = :title,
		"description" = :description,
		"date_updated" = :date_updated
	WHERE
		post_id = :post_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, p); err != nil {
		return fmt.Errorf("updating postID[%s]: %w", p.ID, err)
	}

	return nil
}

// Delete removes a post from the database.
func (s Store) Delete(ctx context.Context, postID string) error {
	data := struct {
		PostID string `db:"post_id"`
	}{
		PostID: postID,
	}

	const q = `
	DELETE FROM
		posts
	WHERE
		post_id = :post_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("deleting postID[%s]: %w", postID, err)
	}

	return nil
}

// Query retrieves a list of existing posts from the database.
func (s Store) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]Post, error) {
	data := struct {
		Offset      int `db:"offset"`
		RowsPerPage int `db:"rows_per_page"`
	}{
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
	}

	const q = `
	SELECT
		*
	FROM
		posts
	ORDER BY
		post_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var ps []Post
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &ps); err != nil {
		return nil, fmt.Errorf("selecting posts: %w", err)
	}

	return ps, nil
}

// QueryByID gets the specified post from the database.
func (s Store) QueryByID(ctx context.Context, postID string) (Post, error) {
	data := struct {
		PostID string `db:"post_id"`
	}{
		PostID: postID,
	}

	const q = `
	SELECT
		*
	FROM
		posts
	WHERE 
		post_id = :post_id`

	var p Post
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &p); err != nil {
		return Post{}, fmt.Errorf("selecting postID[%q]: %w", postID, err)
	}

	return p, nil
}
