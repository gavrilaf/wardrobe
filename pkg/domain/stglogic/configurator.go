package stglogic

import (
	"context"

	"github.com/gavrilaf/wardrobe/pkg/domain/dto"
	"github.com/gavrilaf/wardrobe/pkg/fs"
)

type Configurator interface {
	PrepareStorage(ctx context.Context) error

	GetBucket(infoObject dto.InfoObject, file dto.File) string
	GenerateFileName(infoObject dto.InfoObject, file dto.File) string
}

func NewConfigurator(fs fs.Storage) Configurator {
	return &configurator{
		fs: fs,
	}
}

//

type configurator struct {
	fs fs.Storage
}

func (c *configurator) PrepareStorage(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (c *configurator) GetBucket(infoObject dto.InfoObject, file dto.File) string {
	//TODO implement me
	panic("implement me")
}

func (c *configurator) GenerateFileName(infoObject dto.InfoObject, file dto.File) string {
	//TODO implement me
	panic("implement me")
}
