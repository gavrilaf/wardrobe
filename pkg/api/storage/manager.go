package storage

import (
	"context"
	"io"

	"github.com/gavrilaf/wardrobe/pkg/api/dto"
)

type Manager interface {
	CreateObject(ctx context.Context, fo dto.FO) (int, error)
	UploadContent(ctx context.Context, id int, r io.Reader) error
}

func NewManager() Manager {
	return &manager{}
}

//

type manager struct {
}

func (m *manager) CreateObject(ctx context.Context, fo dto.FO) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (m *manager) UploadContent(ctx context.Context, id int, r io.Reader) error {
	//TODO implement me
	panic("implement me")
}
