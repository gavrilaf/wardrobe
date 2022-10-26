//go:build integration

package repo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gavrilaf/wardrobe/pkg/repo"
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
