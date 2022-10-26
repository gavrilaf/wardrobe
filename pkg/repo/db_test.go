//go:build integration

package repo_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"

	"github.com/gavrilaf/wardrobe/pkg/repo"
)

var (
	dbName       = "wardrobe-test"
	db           *repo.DB
	kyivLocation *time.Location
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	kyivLocation, err = time.LoadLocation("Europe/Kiev") // server timezone
	if err != nil {
		log.Fatalf("Could not load timezone: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("postgres", "15.0", []string{"POSTGRES_USER=postgres", "POSTGRES_PASSWORD=secret", "POSTGRES_DB=" + dbName})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	connString := fmt.Sprintf("postgres://postgres:secret@%s/%s?sslmode=disable", resource.GetHostPort("5432/tcp"), dbName)

	db, err = repo.NewDB(ctx, connString, 2)
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	if _, err = pgx.Connect(ctx, connString); err != nil {
		log.Fatalf("Could not connect with pgx: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestPing(t *testing.T) {
	err := db.Ping()
	assert.NoError(t, err)
}
