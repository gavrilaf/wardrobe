package minio

import (
	"context"
	
	"github.com/gavrilaf/wardrobe/pkg/fs"
)

func New() fs.Storage {
	return &storage{}
}

type storage struct {
}

func (s storage) CreateObject(ctx context.Context, fo fs.Object) error {
	//TODO implement me
	panic("implement me")
}

func (s storage) GetObject(ctx context.Context, name string) (fs.Object, error) {
	//TODO implement me
	panic("implement me")
}

func (s storage) DeleteObject(ctx context.Context, name string) error {
	//TODO implement me
	panic("implement me")
}
