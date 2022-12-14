// Package testgrp contains all the test handlers.
package testgrp

import (
	"context"
	"errors"
	"math/rand"
	"net/http"

	webv1 "github.com/dudakovict/social-network/business/web/v1"
	"github.com/dudakovict/social-network/foundation/web"
	"go.uber.org/zap"
)

// Handlers manages the set of check enpoints.
type Handlers struct {
	Log *zap.SugaredLogger
}

// Test handler is for development.
func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		return webv1.NewRequestError(errors.New("trusted error"), http.StatusBadRequest)
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
