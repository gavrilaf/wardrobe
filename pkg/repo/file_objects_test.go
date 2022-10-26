//go:build integration

package repo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gavrilaf/wardrobe/pkg/repo"
)

func TestFileObjects(t *testing.T) {
	var fileObjects repo.FileObjects = db

	ctx := context.TODO()

	var (
		foID        int
		err         error
		name        = "test-1"
		contentType = "text/plain"
	)

	t.Run("create file object", func(t *testing.T) {
		foID, err = fileObjects.Create(ctx, name, contentType)
		assert.NoError(t, err)

		assert.NotZero(t, foID)
	})

	t.Run("read file object", func(t *testing.T) {
		fo, err := fileObjects.GetById(ctx, foID)
		assert.NoError(t, err)

		assert.Equal(t, foID, fo.ID)
		assert.Equal(t, name, fo.Name)
		assert.Equal(t, contentType, fo.ContentType)
		assert.Zero(t, fo.Size)
		assert.NotZero(t, fo.Created)
		assert.Nil(t, fo.Uploaded)
	})

}
