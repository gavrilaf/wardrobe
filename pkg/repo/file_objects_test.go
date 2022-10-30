//go:build integration

package repo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gavrilaf/wardrobe/pkg/repo"
	"github.com/gavrilaf/wardrobe/pkg/repo/dbtypes"
)

func TestFileObjects(t *testing.T) {
	var fileObjects repo.FileObjects = db
	ctx := context.TODO()

	var (
		fo = dbtypes.FO{
			Name:        "name",
			ContentType: "text/plain",
			Author:      "author",
			Source:      "tg",
			Bucket:      "test-bucket",
			FileName:    "file-name",
		}
		foID int
		err  error
	)

	t.Run("create file object", func(t *testing.T) {
		foID, err = fileObjects.Create(ctx, fo)
		assert.NoError(t, err)

		assert.NotZero(t, foID)
	})

	t.Run("read file object", func(t *testing.T) {
		fo2, err := fileObjects.GetById(ctx, foID)
		assert.NoError(t, err)

		assert.NotZero(t, fo2.Created)

		fo.ID = foID
		fo.Created = fo2.Created

		assert.Equal(t, fo, fo2)
	})
}
