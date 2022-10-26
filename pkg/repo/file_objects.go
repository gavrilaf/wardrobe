package repo

import (
	"context"

	"github.com/gavrilaf/wardrobe/pkg/repo/dbtypes"
)

//go:generate mockery --name FileObjects --outpkg repomocks --output ./repomocks --dir .
type FileObjects interface {
	Create(ctx context.Context, name, contentType string) (int, error)

	GetById(ctx context.Context, id int) (dbtypes.FO, error)
}
