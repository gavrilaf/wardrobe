package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgxutil"

	"github.com/gavrilaf/wardrobe/pkg/repo/dbtypes"
)

//go:generate mockery --name InfoObjects --outpkg repomocks --output ./repomocks --dir .
type InfoObjects interface {
	CreateInfoObject(ctx context.Context, obj dbtypes.InfoObject) (int, error)
	AddFile(ctx context.Context, file dbtypes.File) (int, error)
	MarkUploaded(ctx context.Context, id int) error

	LinkTag(ctx context.Context, objID int, tag string) error

	GetById(ctx context.Context, id int) (dbtypes.InfoObject, error)

	GetFiles(ctx context.Context, id int) ([]dbtypes.File, error)
	GetTags(ctx context.Context, id int) ([]string, error)

	GetFile(ctx context.Context, fileID int) (dbtypes.File, error)

	GetStat(ctx context.Context) (dbtypes.Stat, error)
}

func (db *DB) CreateInfoObject(ctx context.Context, obj dbtypes.InfoObject) (int, error) {
	query := "INSERT INTO info_objects(name, author, source, published) VALUES($1, $2, $3, $4) RETURNING id"
	return pgxutil.SelectValue[int](ctx, db.Doer(ctx), query, obj.Name, obj.Author, obj.Source, obj.Published)
}

func (db *DB) AddFile(ctx context.Context, file dbtypes.File) (int, error) {
	query := "INSERT INTO files(info_object_id, bucket, name, content_type, size) VALUES($1, $2, $3, $4, $5) RETURNING id"
	return pgxutil.SelectValue[int](ctx, db.Doer(ctx), query, file.InfoObjectID, file.Bucket, file.Name, file.ContentType, file.Size)
}

func (db *DB) MarkUploaded(ctx context.Context, id int) error {
	query := "UPDATE info_objects SET finalized = CURRENT_TIMESTAMP WHERE id = $1"
	_, err := db.Doer(ctx).Exec(ctx, query, id)
	return err
}

func (db *DB) LinkTag(ctx context.Context, objID int, tag string) error {
	return db.RunWithTx(ctx, func(ctx context.Context) error {
		tagID, err := db.getOrCreateTag(ctx, tag)
		if err != nil {
			return fmt.Errorf("failed to get tag %s, %w", tag, err)
		}

		query := "INSERT INTO info_objects_tags(info_object_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING"
		_, err = db.Doer(ctx).Exec(ctx, query, objID, tagID)
		if err != nil {
			return fmt.Errorf("failed to link tag %s to object %d, %w", tag, objID, err)
		}

		return nil
	})
}

func (db *DB) GetById(ctx context.Context, id int) (dbtypes.InfoObject, error) {
	query := "SELECT * FROM info_objects WHERE id = $1"

	var obj dbtypes.InfoObject
	err := pgxutil.SelectStruct(ctx, db.Doer(ctx), &obj, query, id)

	return obj, err
}

func (db *DB) GetFiles(ctx context.Context, id int) ([]dbtypes.File, error) {
	query := "SELECT * FROM files WHERE info_object_id = $1 ORDER BY uploaded"

	var files []dbtypes.File
	err := pgxutil.SelectAllStruct(ctx, db.Doer(ctx), &files, query, id)

	return files, err
}

func (db *DB) GetTags(ctx context.Context, id int) ([]string, error) {
	query := "SELECT value " +
		"FROM tags t JOIN info_objects_tags iot ON t.id = iot.tag_id " +
		"WHERE iot.info_object_id = $1 " +
		"ORDER BY value"

	return pgxutil.SelectColumn[string](ctx, db.Doer(ctx), query, id)
}

func (db *DB) GetFile(ctx context.Context, fileID int) (dbtypes.File, error) {
	query := "SELECT * FROM files WHERE id = $1"

	var file dbtypes.File
	err := pgxutil.SelectStruct(ctx, db.Doer(ctx), &file, query, fileID)

	return file, err
}

func (db *DB) GetStat(ctx context.Context) (dbtypes.Stat, error) {
	doer := db.Doer(ctx)

	totalObjects, err := pgxutil.SelectValue[int64](ctx, doer, "SELECT COUNT(*) FROM info_objects")
	if err != nil {
		return dbtypes.Stat{}, fmt.Errorf("failed to calculate objects quantity, %w", err)
	}

	totalFiles, err := pgxutil.SelectValue[int64](ctx, doer, "SELECT COUNT(*) FROM files")
	if err != nil {
		return dbtypes.Stat{}, fmt.Errorf("failed to calculate files quantity, %w", err)
	}

	totalSize, err := pgxutil.SelectValue[int64](ctx, doer, "SELECT SUM(size) FROM files")
	if err != nil {
		return dbtypes.Stat{}, fmt.Errorf("failed to calculate total files size, %w", err)
	}

	return dbtypes.Stat{
		ObjectsCount: totalObjects,
		FilesCount:   totalFiles,
		TotalSize:    totalSize,
	}, nil
}

// private

func (db *DB) getOrCreateTag(ctx context.Context, tag string) (int, error) {
	query := "WITH new_row AS ( INSERT INTO tags(value) SELECT $1 WHERE NOT EXISTS ( " +
		"SELECT id FROM tags WHERE value = $1 ) " +
		"RETURNING id ) " +
		"( SELECT id FROM tags WHERE value = $1 " +
		"UNION ALL " +
		"SELECT id FROM new_row " +
		") LIMIT 1"

	return pgxutil.SelectValue[int](ctx, db.Doer(ctx), query, tag)
}
