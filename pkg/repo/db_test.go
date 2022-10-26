//go:build integration

package repo_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

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
		log.Panicf("could not connect to docker: %s", err)
	}

	kyivLocation, err = time.LoadLocation("Europe/Kiev") // server timezone
	if err != nil {
		log.Panicf("could not load timezone: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("postgres", "15.0", []string{"POSTGRES_USER=postgres", "POSTGRES_PASSWORD=secret", "POSTGRES_DB=" + dbName})
	if err != nil {
		log.Panicf("could not start resource: %s", err)
	}

	defer func() {
		if err = pool.Purge(resource); err != nil {
			log.Printf("could not purge resource: %s\n", err)
		}
	}()

	connString := fmt.Sprintf("postgres://postgres:secret@%s/%s?sslmode=disable", resource.GetHostPort("5432/tcp"), dbName)

	db, err = repo.NewDB(ctx, connString, 2)
	if err != nil {
		log.Panicf("could not start database: %s", err)
	}

	if err = db.Migrate(ctx, "../../migration"); err != nil {
		log.Panicf("migration failed: %s", err)
	}

	code := m.Run()

	log.Printf("Tests finished with code: %d\n", code)
}

func TestPing(t *testing.T) {
	err := db.Ping()
	assert.NoError(t, err)
}
