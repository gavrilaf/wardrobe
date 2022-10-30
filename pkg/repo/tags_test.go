//go:build integration

package repo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gavrilaf/wardrobe/pkg/repo"
	"github.com/gavrilaf/wardrobe/pkg/repo/dbtypes"
)

func TestTags(t *testing.T) {
	var tagsDB repo.Tags = db

	ctx := context.TODO()

	var (
		tag1ID int
		err    error
	)

	t.Run("create tag", func(t *testing.T) {
		tag1ID, err = tagsDB.GetOrCreateTag(ctx, "tag1")
		assert.NoError(t, err)

		assert.NotZero(t, tag1ID)
	})

	t.Run("get tag by id", func(t *testing.T) {
		s, err := tagsDB.GetTag(ctx, tag1ID)
		assert.NoError(t, err)

		assert.Equal(t, "tag1", s)
	})

	t.Run("creating the same tag should return existing", func(t *testing.T) {
		id, err := tagsDB.GetOrCreateTag(ctx, "tag1")
		assert.NoError(t, err)

		assert.Equal(t, tag1ID, id)
	})
}

func TestFileObjectTags(t *testing.T) {
	ctx := context.TODO()

	var tagsDB repo.Tags = db
	var fileObjects repo.FileObjects = db

	fo := dbtypes.FO{Name: "fo-tags"}

	t.Run("link and read file object tags", func(t *testing.T) {
		foID, err := fileObjects.Create(ctx, fo)
		assert.NoError(t, err)

		tags, err := tagsDB.GetFileObjectTags(ctx, foID)
		assert.NoError(t, err)
		assert.Empty(t, tags)

		tag1ID, err := tagsDB.GetOrCreateTag(ctx, "fo-tag1")
		assert.NoError(t, err)

		err = tagsDB.LinkTag(ctx, foID, tag1ID)
		assert.NoError(t, err)

		tags, err = tagsDB.GetFileObjectTags(ctx, foID)
		assert.NoError(t, err)

		assert.Equal(t, []string{"fo-tag1"}, tags)

		tag2ID, err := tagsDB.GetOrCreateTag(ctx, "fo-tag2")
		assert.NoError(t, err)

		err = tagsDB.LinkTag(ctx, foID, tag2ID)
		assert.NoError(t, err)

		tags, err = tagsDB.GetFileObjectTags(ctx, foID)
		assert.NoError(t, err)

		assert.Equal(t, []string{"fo-tag1", "fo-tag2"}, tags)
	})
}
