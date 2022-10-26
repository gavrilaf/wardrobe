package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/gavrilaf/wardrobe/pkg/api/dto"
	"github.com/gavrilaf/wardrobe/pkg/fs"
	"github.com/gavrilaf/wardrobe/pkg/repo"
	"github.com/gavrilaf/wardrobe/pkg/utils/log"
)

type Manager interface {
	CreateObject(ctx context.Context, fo dto.FO) (int, error)
	UploadContent(ctx context.Context, id int, r io.Reader, size int64) error
}

type Config struct {
	Tx          repo.TxRunner
	FileObjects repo.FileObjects
	Tags        repo.Tags
	Stg         fs.Storage
}

func NewManager(cfg Config) Manager {
	return &manager{
		tx:          cfg.Tx,
		fileObjects: cfg.FileObjects,
		stg:         cfg.Stg,
	}
}

//

type manager struct {
	tx          repo.TxRunner
	fileObjects repo.FileObjects
	tags        repo.Tags
	stg         fs.Storage
}

func (m *manager) CreateObject(ctx context.Context, fo dto.FO) (int, error) {
	var (
		fileObjectID int
		err          error
	)

	err = m.tx.RunWithTx(ctx, func(ctx context.Context) error {
		foID, err := m.fileObjects.Create(ctx, fo.Name, fo.ContentType)
		if err != nil {
			return fmt.Errorf("failed to add file meta to the db (%s, %s), %w", fo.Name, fo.ContentType, err)
		}

		for _, tag := range fo.Tags {
			tagID, err := m.tags.GetOrCreateTag(ctx, tag)
			if err != nil {
				return fmt.Errorf("failed to get or create tag (%s, %s), %w", fo.Name, tag, err)
			}

			if err = m.tags.LinkTag(ctx, foID, tagID); err != nil {
				return fmt.Errorf("failed to link tag %d to the file object %d, %w", tagID, foID, err)
			}
		}

		fileObjectID = foID
		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to create file object, %w", err)
	}

	log.FromContext(ctx).Infof("object created (%d, %s, %s)", fileObjectID, fo.Name, fo.ContentType)

	return fileObjectID, err
}

func (m *manager) UploadContent(ctx context.Context, id int, r io.Reader, size int64) error {
	foDb, err := m.fileObjects.GetById(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to retrieve file object from db %d, %w", id, err)
	}

	err = m.stg.CreateObject(ctx, fs.Object{
		Name:        foDb.Name,
		ContentType: foDb.ContentType,
		Size:        size,
		Reader:      r,
	})
	if err != nil {
		return fmt.Errorf("failed to upload object to the storage (%d, %s), %w", id, foDb.Name, err)
	}

	log.FromContext(ctx).Infof("object uploaded, (%d, %s, %s, %d)", id, foDb.Name, foDb.ContentType, size)
	return nil
}
