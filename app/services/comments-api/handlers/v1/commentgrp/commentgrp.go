// Package commentgrp maintains the group of handlers for comment access.
package commentgrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dudakovict/social-network/business/core/comment"
	"github.com/dudakovict/social-network/business/sys/auth"
	v1Web "github.com/dudakovict/social-network/business/web/v1"
	"github.com/dudakovict/social-network/foundation/web"
)

// Handlers manages the set of comment enpoints.
type Handlers struct {
	Core comment.Core
	Auth *auth.Auth
}

// Create adds a new comment to the system.
func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	var nc comment.NewComment
	if err := web.Decode(r, &nc); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	c, err := h.Core.Create(ctx, nc, v.Now)
	if err != nil {
		return fmt.Errorf("comment[%+v]: %w", &c, err)
	}

	return web.Respond(ctx, w, c, http.StatusCreated)
}

// Update updates a comment in the system.
func (h Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	var upd comment.UpdateComment
	if err := web.Decode(r, &upd); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	commentID := web.Param(r, "id")

	c, err := h.Core.QueryByID(ctx, commentID)
	if err != nil {
		return v1Web.NewRequestError(err, http.StatusBadRequest)
	}

	// If you are not an admin and looking to retrieve someone other than yourself.
	if !claims.Authorized(auth.RoleAdmin) && claims.Subject != c.UserID {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	if err := h.Core.Update(ctx, commentID, upd, v.Now); err != nil {
		switch {
		case errors.Is(err, comment.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, comment.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s] Comment[%+v]: %w", commentID, &upd, err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Delete removes a comment from the system.
func (h Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	commentID := web.Param(r, "id")

	c, err := h.Core.QueryByID(ctx, commentID)
	if err != nil {
		return v1Web.NewRequestError(err, http.StatusBadRequest)
	}

	// If you are not an admin and looking to delete someone other than yourself.
	if !claims.Authorized(auth.RoleAdmin) && claims.Subject != c.UserID {
		return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
	}

	if err := h.Core.Delete(ctx, commentID); err != nil {
		switch {
		case errors.Is(err, comment.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, comment.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s]: %w", commentID, err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Query returns a list of comments with paging.
func (h Handlers) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	page := web.Param(r, "page")
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return v1Web.NewRequestError(fmt.Errorf("invalid page format [%s]", page), http.StatusBadRequest)
	}
	rows := web.Param(r, "rows")
	rowsPerPage, err := strconv.Atoi(rows)
	if err != nil {
		return v1Web.NewRequestError(fmt.Errorf("invalid rows format [%s]", rows), http.StatusBadRequest)
	}

	comments, err := h.Core.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return fmt.Errorf("unable to query for comments: %w", err)
	}

	return web.Respond(ctx, w, comments, http.StatusOK)
}

// QueryByID returns a comment by its ID.
func (h Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	/*
		claims, err := auth.GetClaims(ctx)
		if err != nil {
			return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
		}
	*/

	commentID := web.Param(r, "id")

	/*
		// If you are not an admin and looking to retrieve someone other than yourself.
		if !claims.Authorized(auth.RoleAdmin) && claims.Subject != postID {
			return v1Web.NewRequestError(auth.ErrForbidden, http.StatusForbidden)
		}
	*/

	c, err := h.Core.QueryByID(ctx, commentID)
	if err != nil {
		switch {
		case errors.Is(err, comment.ErrInvalidID):
			return v1Web.NewRequestError(err, http.StatusBadRequest)
		case errors.Is(err, comment.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s]: %w", commentID, err)
		}
	}

	return web.Respond(ctx, w, c, http.StatusOK)
}
