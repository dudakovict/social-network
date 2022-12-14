// Package dbtest contains supporting code for running tests that hit the DB.
package dbtest

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"testing"
	"time"

	"github.com/dudakovict/social-network/business/data/comment/dbschema"
	"github.com/dudakovict/social-network/business/sys/auth"
	"github.com/dudakovict/social-network/business/sys/database"
	"github.com/dudakovict/social-network/business/sys/nats"
	"github.com/dudakovict/social-network/foundation/docker"
	"github.com/dudakovict/social-network/foundation/keystore"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Success and failure markers.
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

// StartDB starts a database instance.
func StartDB() (*docker.Container, error) {
	image := "postgres:13-alpine"
	port := "5432"
	args := []string{"-e", "POSTGRES_PASSWORD=postgres"}

	return docker.StartContainer(image, port, args...)
}

// StopDB stops a running database instance.
func StopDB(c *docker.Container) {
	docker.StopContainer(c.ID)
}

func StartNATS() (*docker.Container, error) {
	image := "nats-streaming:0.17.0"
	port := "4222"
	args := []string{"-p", "4222", "-m", "8222", "-hbi", "5s", "-hbt", "5s", "-hbf", "2", "-SD", "-cid", "social-network"}

	return docker.StartContainer(image, port, args...)
}

func StopNATS(c *docker.Container) {
	docker.StopContainer(c.ID)
}

// NewUnit creates a test database inside a Docker container. It creates the
// required table structure but the database is otherwise empty. It returns
// the database to use as well as a function to call at the end of the test.
func NewUnit(t *testing.T, nc *docker.Container, dbc *docker.Container, dbName string) (*zap.SugaredLogger, *sqlx.DB, *nats.NATS, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbM, err := database.Open(database.Config{
		User:       "postgres",
		Password:   "postgres",
		Host:       dbc.Host,
		Name:       "postgres",
		DisableTLS: true,
	})
	if err != nil {
		t.Fatalf("Opening database connection: %v", err)
	}

	t.Log("Waiting for database to be ready ...")

	if err := database.StatusCheck(ctx, dbM); err != nil {
		t.Fatalf("status check database: %v", err)
	}

	t.Log("Database ready")

	if _, err := dbM.ExecContext(context.Background(), "CREATE DATABASE "+dbName); err != nil {
		t.Fatalf("creating database %s: %v", dbName, err)
	}
	dbM.Close()

	// =========================================================================

	db, err := database.Open(database.Config{
		User:       "postgres",
		Password:   "postgres",
		Host:       dbc.Host,
		Name:       dbName,
		DisableTLS: true,
	})
	if err != nil {
		t.Fatalf("Opening database connection: %v", err)
	}

	t.Log("Migrate and seed database ...")

	if err := dbschema.Migrate(ctx, db); err != nil {
		docker.DumpContainerLogs(t, dbc.ID)
		t.Fatalf("Migrating error: %s", err)
	}

	if err := dbschema.Seed(ctx, db); err != nil {
		docker.DumpContainerLogs(t, dbc.ID)
		t.Fatalf("Seeding error: %s", err)
	}

	t.Log("Ready for testing ...")

	var buf bytes.Buffer
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	writer := bufio.NewWriter(&buf)
	log := zap.New(
		zapcore.NewCore(encoder, zapcore.AddSync(writer), zapcore.DebugLevel)).
		Sugar()

	t.Log("Opening NATS connection ...")

	n, err := nats.Connect(nats.Config{
		ClusterID: "social-network",
		ClientID:  "posts",
		Host:      nc.Host,
	})

	if err != nil {
		t.Fatalf("Connecting to NATS: %v", err)
	}

	t.Log("NATS ready ...")

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		n.Client.Close()
		db.Close()

		log.Sync()

		writer.Flush()
		fmt.Println("******************** LOGS ********************")
		fmt.Print(buf.String())
		fmt.Println("******************** LOGS ********************")
	}

	return log, db, n, teardown
}

// Test owns state for running and shutting down tests.
type Test struct {
	DB       *sqlx.DB
	Log      *zap.SugaredLogger
	Auth     *auth.Auth
	NATS     *nats.NATS
	Teardown func()

	t *testing.T
}

// NewIntegration creates a database, seeds it, constructs an authenticator.
func NewIntegration(t *testing.T, nc *docker.Container, dbc *docker.Container, dbName string) *Test {
	log, db, n, teardown := NewUnit(t, nc, dbc, dbName)

	// Create RSA keys to enable authentication in our service.
	keyID := "4754d86b-7a6d-4df5-9c65-224741361492"
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	// Build an authenticator using this private key and id for the key store.
	auth, err := auth.New(keyID, keystore.NewMap(map[string]*rsa.PrivateKey{keyID: privateKey}))
	if err != nil {
		t.Fatal(err)
	}

	test := Test{
		DB:       db,
		Log:      log,
		Auth:     auth,
		NATS:     n,
		t:        t,
		Teardown: teardown,
	}

	return &test
}

// StringPointer is a helper to get a *string from a string. It is in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func StringPointer(s string) *string {
	return &s
}

// IntPointer is a helper to get a *int from a int. It is in the tests package
// because we normally don't want to deal with pointers to basic types but it's
// useful in some tests.
func IntPointer(i int) *int {
	return &i
}
