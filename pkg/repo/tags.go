package repo

import (
	"context"
	"github.com/jackc/pgxutil"
)

//go:generate mockery --name Tags --outpkg repomocks --output ./repomocks --dir .
type Tags interface {
	GetOrCreateTag(ctx context.Context, tag string) (int, error)
	GetTag(ctx context.Context, id int) (string, error)
	LinkTag(ctx context.Context, foID, tagID int) error
}

func (db *DB) GetOrCreateTag(ctx context.Context, tag string) (int, error) {
	query := "WITH new_row AS ( INSERT INTO tags(value) SELECT $1 WHERE NOT EXISTS ( " +
		"SELECT id FROM tags WHERE value = $1 ) " +
		"RETURNING id ) " +
		"( SELECT id FROM tags WHERE value = $1 " +
		"UNION ALL " +
		"SELECT id FROM new_row " +
		") LIMIT 1"

	return pgxutil.SelectValue[int](ctx, db.Doer(ctx), query, tag)
}

func (db *DB) GetTag(ctx context.Context, id int) (string, error) {
	query := "SELECT value FROM tags WHERE id = $1"

	return pgxutil.SelectValue[string](ctx, db.Doer(ctx), query, id)
}

func (db *DB) LinkTag(ctx context.Context, foID, tagID int) error {
	return nil
}
