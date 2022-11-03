//go:build integration

package repo_test

import (
	"context"
	"github.com/gavrilaf/wardrobe/pkg/utils/timex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gavrilaf/wardrobe/pkg/repo"
	"github.com/gavrilaf/wardrobe/pkg/repo/dbtypes"
)

func TestInfoObjects(t *testing.T) {
	ctx := context.TODO()

	var infoObjects repo.InfoObjects = db

	var (
		infoObject = dbtypes.InfoObject{
			Name:      "name",
			Author:    "author",
			Source:    "tg",
			Published: timex.Date(2022, time.January, 10),
		}

		file1 = dbtypes.File{
			Bucket:      "bucket1",
			Name:        "file1",
			ContentType: "text/plain",
			Size:        123,
		}

		file2 = dbtypes.File{
			Bucket:      "bucket2",
			Name:        "file2",
			ContentType: "image/jpg",
			Size:        1024,
		}

		tag1 = "tag1"
		tag2 = "very long tag we must check it"

		infoObjectID, file1ID, file2ID int
		err                            error
	)

	t.Run("create info object", func(t *testing.T) {
		infoObjectID, err = infoObjects.CreateInfoObject(ctx, infoObject)
		assert.NoError(t, err)

		assert.NotZero(t, infoObjectID)
	})

	t.Run("get info object", func(t *testing.T) {
		readObj, err := infoObjects.GetById(ctx, infoObjectID)
		assert.NoError(t, err)

		assert.NotZero(t, readObj.Created)
		assert.Nil(t, readObj.Finalized)

		expected := infoObject
		expected.ID = infoObjectID
		expected.Created = readObj.Created

		assert.Equal(t, expected, expected)
	})

	t.Run("add files", func(t *testing.T) {
		file1.InfoObjectID = infoObjectID
		file1ID, err = infoObjects.AddFile(ctx, file1)
		assert.NoError(t, err)
		assert.NotZero(t, file1ID)

		file2.InfoObjectID = infoObjectID
		file2ID, err = infoObjects.AddFile(ctx, file2)
		assert.NoError(t, err)
		assert.NotZero(t, file2ID)
	})

	t.Run("mark info object uploaded", func(t *testing.T) {
		err = infoObjects.MarkUploaded(ctx, infoObjectID)
		assert.NoError(t, err)

		t.Run("uploaded should not be empty", func(t *testing.T) {
			readObj, err := infoObjects.GetById(ctx, infoObjectID)
			assert.NoError(t, err)
			assert.NotZero(t, readObj.Finalized)
		})
	})

	t.Run("get files", func(t *testing.T) {
		files, err := infoObjects.GetFiles(ctx, infoObjectID)
		assert.NoError(t, err)
		assert.Len(t, files, 2)

		assert.NotZero(t, files[0].Uploaded)
		assert.NotZero(t, files[1].Uploaded)

		expected := []dbtypes.File{file1, file2}
		expected[0].ID = file1ID
		expected[0].Uploaded = files[0].Uploaded
		expected[1].ID = file2ID
		expected[1].Uploaded = files[1].Uploaded

		assert.Equal(t, expected, files)
	})

	t.Run("link tags", func(t *testing.T) {
		err = infoObjects.LinkTag(ctx, infoObjectID, tag1)
		assert.NoError(t, err)

		err = infoObjects.LinkTag(ctx, infoObjectID, tag2)
		assert.NoError(t, err)

		t.Run("link tags again", func(t *testing.T) {
			err = infoObjects.LinkTag(ctx, infoObjectID, tag1)
			assert.NoError(t, err)
		})
	})

	t.Run("get tags", func(t *testing.T) {
		tags, err := infoObjects.GetTags(ctx, infoObjectID)
		assert.NoError(t, err)

		expected := []string{tag1, tag2}
		assert.Equal(t, expected, tags)
	})
}
