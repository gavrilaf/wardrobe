package repo

import (
	"context"

	"github.com/gavrilaf/wardrobe/pkg/repo/dbtypes"
)

type FO interface {
	Add(ctx context.Context, name, contentType string, size int64) (int, error)

	GetByName(ctx context.Context, name string) (dbtypes.FO, error)
}
