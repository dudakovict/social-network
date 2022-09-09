// Package db contains comment related CRUD functionality.
package db

import (
	"context"
	"fmt"

	"github.com/dudakovict/social-network/business/sys/database"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Store manages the set of API's for comment access.
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

// Create inserts a new comment into the database.
func (s Store) Create(ctx context.Context, c Comment) error {
	const q = `
	INSERT INTO comments
		(comment_id, description, post_id, user_id, date_created, date_updated)
	VALUES
		(:comment_id, :description, :post_id, :user_id, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, c); err != nil {
		return fmt.Errorf("inserting comment: %w", err)
	}

	return nil
}

// Update replaces a comment document in the database.
func (s Store) Update(ctx context.Context, c Comment) error {
	const q = `
	UPDATE
		comments
	SET 
		"description" = :description,
		"date_updated" = :date_updated
	WHERE
		comment_id = :comment_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, c); err != nil {
		return fmt.Errorf("updating commentID[%s]: %w", c.ID, err)
	}

	return nil
}

// Delete removes a comment from the database.
func (s Store) Delete(ctx context.Context, commentID string) error {
	data := struct {
		CommentID string `db:"comment_id"`
	}{
		CommentID: commentID,
	}

	const q = `
	DELETE FROM
		comments
	WHERE
		comment_id = :comment_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("deleting commentID[%s]: %w", commentID, err)
	}

	return nil
}

// Query retrieves a list of existing comments from the database.
func (s Store) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]Comment, error) {
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
		comments
	ORDER BY
		comment_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var comms []Comment
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &comms); err != nil {
		return nil, fmt.Errorf("selecting comments: %w", err)
	}

	return comms, nil
}

// QueryByID gets the specified comment from the database.
func (s Store) QueryByID(ctx context.Context, commentID string) (Comment, error) {
	data := struct {
		CommentID string `db:"comment_id"`
	}{
		CommentID: commentID,
	}

	const q = `
	SELECT
		*
	FROM
		comments
	WHERE 
		comment_id = :comment_id`

	var c Comment
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &c); err != nil {
		return Comment{}, fmt.Errorf("selecting commentID[%q]: %w", commentID, err)
	}

	return c, nil
}

func (s Store) QueryByUserID(ctx context.Context, userID string) ([]Comment, error) {
	data := struct {
		UserID string `db:"user_id"`
	}{
		UserID: userID,
	}

	const q = `
	SELECT
		*
	FROM
		comments
	WHERE
		user_id = :user_id`

	var comms []Comment
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &comms); err != nil {
		return nil, fmt.Errorf("selecting comments userID[%s]: %w", userID, err)
	}

	return comms, nil
}

func (s Store) QueryByPostID(ctx context.Context, postID string) ([]Comment, error) {
	data := struct {
		PostID string `db:"post_id"`
	}{
		PostID: postID,
	}

	const q = `
	SELECT
		*
	FROM
		comments
	WHERE
		post_id = :post_id`

	var comms []Comment
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &comms); err != nil {
		return nil, fmt.Errorf("selecting comments postID[%s]: %w", postID, err)
	}

	return comms, nil
}

func (s Store) CreatePost(ctx context.Context, p Post) error {
	const q = `
	INSERT INTO posts
		(post_id, title, description, user_id, date_created, date_updated)
	VALUES
		(:post_id, :title, :description, :user_id, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, p); err != nil {
		return fmt.Errorf("inserting post: %w", err)
	}

	return nil
}

func (s Store) UpdatePost(ctx context.Context, p Post) error {
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

func (s Store) DeletePost(ctx context.Context, postID string) error {
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
		return fmt.Errorf("deleting commentID[%s]: %w", postID, err)
	}

	return nil
}

func (s Store) QueryPostsByPostID(ctx context.Context, postID string) ([]PostComment, error) {
	data := struct {
		PostID string `db:"post_id"`
	}{
		PostID: postID,
	}

	const q = `
	SELECT
		p.*, c.comment_id, c.description as comment_description, c.user_id as comment_user_id, c.date_created as comment_date_created, c.date_updated as comment_date_updated
	FROM
		posts as p
	LEFT JOIN
		comments as c
	ON
		p.post_id = c.post_id
	WHERE
		p.post_id = :post_id`

	var pscomms []PostComment
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &pscomms); err != nil {
		return nil, fmt.Errorf("selecting comments postID[%s]: %w", postID, err)
	}

	return pscomms, nil
}
