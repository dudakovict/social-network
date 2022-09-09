// Package handlers manages the different versions of the API.
package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/dudakovict/social-network/app/services/posts-api/handlers/debug/checkgrp"
	v1PostGrp "github.com/dudakovict/social-network/app/services/posts-api/handlers/v1/postgrp"
	v1TestGrp "github.com/dudakovict/social-network/app/services/posts-api/handlers/v1/testgrp"
	postCore "github.com/dudakovict/social-network/business/core/post"
	"github.com/dudakovict/social-network/business/sys/auth"
	"github.com/dudakovict/social-network/business/sys/nats"
	"github.com/dudakovict/social-network/business/web/v1/mid"
	"github.com/dudakovict/social-network/foundation/web"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// DebugStandardLibraryMux registers all the debug routes from the standard library
// into a new mux bypassing the use of the DefaultServerMux. Using the
// DefaultServerMux would be a security risk since a dependency could inject a
// handler into our service without us knowing it.
func DebugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Register all the standard library debug endpoints.
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

// DebugMux registers all the debug standard library routes and then custom
// debug application routes for the service. This bypassing the use of the
// DefaultServerMux. Using the DefaultServerMux would be a security risk since
// a dependency could inject a handler into our service without us knowing it.
func DebugMux(build string, log *zap.SugaredLogger, db *sqlx.DB) http.Handler {
	mux := DebugStandardLibraryMux()

	// Register debug check endpoints.
	cgh := checkgrp.Handlers{
		Build: build,
		Log:   log,
		DB:    db,
	}
	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	return mux
}

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
	Auth     *auth.Auth
	DB       *sqlx.DB
	NATS     *nats.NATS
}

// APIMux constructs an http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig) *web.App {

	// Construct the web.App which holds all routes.
	app := web.NewApp(
		cfg.Shutdown,
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
		mid.Metrics(),
		mid.Panics(),
	)

	// Load the routes for the different versions of the API.
	v1(app, cfg)

	return app
}

// v1 binds all the version 1 routes.
func v1(app *web.App, cfg APIMuxConfig) {
	const version = "v1"

	tgh := v1TestGrp.Handlers{
		Log: cfg.Log,
	}
	app.Handle(http.MethodGet, version, "/test", tgh.Test)
	app.Handle(http.MethodGet, version, "/testauth", tgh.Test, mid.Authenticate(cfg.Auth), mid.Authorize("ADMIN"))

	// Register post management and authentication endpoints.
	pgh := v1PostGrp.Handlers{
		Core: postCore.NewCore(cfg.Log, cfg.DB, cfg.NATS),
		Auth: cfg.Auth,
	}
	app.Handle(http.MethodGet, version, "/posts/:page/:rows", pgh.Query, mid.Authenticate(cfg.Auth))
	app.Handle(http.MethodGet, version, "/posts/:id", pgh.QueryByID, mid.Authenticate(cfg.Auth))
	app.Handle(http.MethodPost, version, "/posts", pgh.Create, mid.Authenticate(cfg.Auth))
	app.Handle(http.MethodPut, version, "/posts/:id", pgh.Update, mid.Authenticate(cfg.Auth))
	app.Handle(http.MethodDelete, version, "/posts/:id", pgh.Delete, mid.Authenticate(cfg.Auth))
}
