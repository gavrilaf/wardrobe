package repo

import (
	"context"

	"github.com/jackc/pgxutil"

	"github.com/gavrilaf/wardrobe/pkg/repo/dbtypes"
)

//go:generate mockery --name FileObjects --outpkg repomocks --output ./repomocks --dir .
type FileObjects interface {
	Create(ctx context.Context, fo dbtypes.FO) (int, error)
	MarkAsUploaded(ctx context.Context, id int, size int64) error

	GetById(ctx context.Context, id int) (dbtypes.FO, error)
}

func (db *DB) Create(ctx context.Context, fo dbtypes.FO) (int, error) {
	query := "INSERT INTO file_objects(name, content_type, author, source, bucket, file_name) " +
		"VALUES($1, $2, $3, $4, $5, $6) " +
		"RETURNING id"

	return pgxutil.SelectValue[int](ctx, db.Doer(ctx), query, fo.Name, fo.ContentType, fo.Author, fo.Source, fo.Bucket, fo.FileName)
}

func (db *DB) MarkAsUploaded(ctx context.Context, id int, size int64) error {
	query := "UPDATE file_objects SET size = $2, uploaded = CURRENT_TIMESTAMP WHERE id = $1"

	_, err := db.Doer(ctx).Exec(ctx, query, id, size)
	return err
}

func (db *DB) GetById(ctx context.Context, id int) (dbtypes.FO, error) {
	query := "SELECT id, name, content_type, author, source, bucket, file_name, size, created, uploaded " +
		"FROM file_objects " +
		"WHERE id = $1"

	var fo dbtypes.FO
	err := pgxutil.SelectStruct(ctx, db.Doer(ctx), &fo, query, id)

	return fo, err
}
