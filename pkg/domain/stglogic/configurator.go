package stglogic

import (
	"context"
	"fmt"
	"path"

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

type configurator struct {
	fs  fs.Storage
	snf idgen.Snowflake
}

func (c *configurator) PrepareStorage(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (c *configurator) GetBucket(obj dbtypes.InfoObject, file dbtypes.File) string {
	//TODO implement me
	panic("implement me")
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
