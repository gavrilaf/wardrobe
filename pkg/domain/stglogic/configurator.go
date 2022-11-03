package stglogic

import (
	"context"
	"fmt"
	"path"

	"github.com/gavrilaf/wardrobe/pkg/utils/log"

	"github.com/gavrilaf/wardrobe/pkg/fs"
	"github.com/gavrilaf/wardrobe/pkg/repo/dbtypes"
	"github.com/gavrilaf/wardrobe/pkg/utils/idgen"
)

type Configurator interface {
	PrepareStorage(ctx context.Context) error

	GetBucket(obj dbtypes.InfoObject, file dbtypes.File) string
	GenerateFileName(obj dbtypes.InfoObject, file dbtypes.File) (string, error)
}

func NewConfigurator(fs fs.Storage, snf idgen.Snowflake) Configurator {
	return &configurator{
		fs:  fs,
		snf: snf,
	}
}

//

const bucketName = "wardrobe-info-objects"

type configurator struct {
	fs  fs.Storage
	snf idgen.Snowflake
}

func (c *configurator) PrepareStorage(ctx context.Context) error {
	exists, err := c.fs.IsBucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket %s existing, %w", bucketName, err)
	}

	if !exists {
		err = c.fs.CreateBucket(ctx, bucketName)
		if err != nil {
			return fmt.Errorf("failed to create bucket %s, %w", bucketName, err)
		}

		log.FromContext(ctx).Infof("bucket %s created", bucketName)
	} else {
		log.FromContext(ctx).Infof("bucket %s already exists", bucketName)
	}

	return nil
}

func (c *configurator) GetBucket(obj dbtypes.InfoObject, file dbtypes.File) string {
	// simplest logic for now, only one static bucket
	return bucketName
}

func (c *configurator) GenerateFileName(obj dbtypes.InfoObject, file dbtypes.File) (string, error) {
	fn := file.Name

	ext := path.Ext(fn)
	if len(ext) > 0 {
		fn = fn[:len(fn)-len(ext)] // remove extension
	}

	pt := obj.Published

	nextID, err := c.snf.NextID()
	if err != nil {
		return "", fmt.Errorf("failed to generate file name (%d, %s), %w", obj.ID, file.Name, err)
	}

	newFn := fmt.Sprintf("%d-%02d-%02d-%02d-%02d-%s-%d",
		pt.Year(), pt.Month(), pt.Day(), pt.Hour(), pt.Minute(), fn, nextID)

	if len(ext) > 0 {
		newFn += ext
	}

	return newFn, nil
}
