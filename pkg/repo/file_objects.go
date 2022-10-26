package repo

import (
	"context"

	"github.com/jackc/pgxutil"

	"github.com/gavrilaf/wardrobe/pkg/repo/dbtypes"
)

//go:generate mockery --name FileObjects --outpkg repomocks --output ./repomocks --dir .
type FileObjects interface {
	Create(ctx context.Context, name, contentType string) (int, error)

	GetById(ctx context.Context, id int) (dbtypes.FO, error)
}

func (db *DB) Create(ctx context.Context, name, contentType string) (int, error) {
	query := "INSERT INTO file_objects(name, content_type) " +
		"VALUES($1, $2) " +
		"RETURNING id"

	return pgxutil.SelectValue[int](ctx, db.Doer(ctx), query, name, contentType)
}

func (db *DB) GetById(ctx context.Context, id int) (dbtypes.FO, error) {
	query := "SELECT id, name, content_type, size, created, uploaded " +
		"FROM file_objects " +
		"WHERE id = $1"

	var fo dbtypes.FO
	err := pgxutil.SelectStruct(ctx, db.Doer(ctx), &fo, query, id)

	return fo, err
}
