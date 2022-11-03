//go:build integration

package minio_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"testing"

	"github.com/cenkalti/backoff"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"

	"github.com/gavrilaf/wardrobe/pkg/fs"
	"github.com/gavrilaf/wardrobe/pkg/fs/minio"
)

var storage fs.Storage

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Panicf("could not connect to docker: %s", err)
	}

	options := &dockertest.RunOptions{
		Repository: "minio/minio",
		Tag:        "latest",
		Cmd:        []string{"server", "/data"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"9000/tcp": {{HostPort: "9000"}},
		},
		Env: []string{"MINIO_ROOT_USER=minio", "MINIO_ROOT_PASSWORD=miniopsw"},
	}

	resource, err := pool.RunWithOptions(options)
	if err != nil {
		log.Panicf("could not start resource: %s", err)
	}

	defer func() {
		if err = pool.Purge(resource); err != nil {
			log.Printf("could not purge resource: %s\n", err)
		}
	}()

	endpoint := fmt.Sprintf("localhost:%s", resource.GetPort("9000/tcp"))

	storage, err = minio.New(endpoint, "minio", "miniopsw")
	if err != nil {
		log.Panicf("could not create minio: %s", err)
	}

	code := m.Run()

	log.Printf("Tests finished with code: %d\n", code)
}

func TestMinioStorage(t *testing.T) {
	ctx := context.TODO()

	bucketName := "test-bucket"
	fileName := "test-file"

	t.Run("minio is online", func(t *testing.T) {
		err := backoff.Retry(func() error {
			return storage.Ping()
		}, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3))

		assert.NoError(t, err)
	})

	t.Run("bucket should not exist", func(t *testing.T) {
		exist, err := storage.IsBucketExists(ctx, bucketName)
		assert.NoError(t, err)
		assert.False(t, exist)
	})

	t.Run("create bucket", func(t *testing.T) {
		err := storage.CreateBucket(ctx, bucketName)
		assert.NoError(t, err)
	})

	t.Run("bucket should exist", func(t *testing.T) {
		exist, err := storage.IsBucketExists(ctx, bucketName)
		assert.NoError(t, err)
		assert.True(t, exist)
	})

	var content = "This is content example"

	t.Run("create object", func(t *testing.T) {
		err := storage.CreateFile(ctx, fs.File{
			Bucket:      bucketName,
			Name:        fileName,
			ContentType: "text/plain",
			Size:        int64(len(content)),
			Reader:      io.NopCloser(strings.NewReader(content)),
		})

		assert.NoError(t, err)
	})

	t.Run("read object", func(t *testing.T) {
		fp, err := storage.GetFile(ctx, bucketName, fileName)
		assert.NoError(t, err)

		assert.Equal(t, bucketName, fp.Bucket)
		assert.Equal(t, fileName, fp.Name)
		assert.Equal(t, "text/plain", fp.ContentType)
		assert.Equal(t, int64(len(content)), fp.Size)

		s, err := io.ReadAll(fp.Reader)
		assert.NoError(t, err)

		assert.Equal(t, content, string(s))

		assert.NoError(t, fp.Reader.Close())
	})
}
