package storage

import (
	"context"
	"fmt"
	dto2 "github.com/gavrilaf/wardrobe/pkg/domain/dto"
	"io"
	"strings"
	"time"

	"github.com/gavrilaf/wardrobe/pkg/api/dto"
	"github.com/gavrilaf/wardrobe/pkg/fs"
	"github.com/gavrilaf/wardrobe/pkg/repo"
	"github.com/gavrilaf/wardrobe/pkg/utils/log"
)

type Manager interface {
	CreateObject(ctx context.Context, fo dto.FO) (int, error)
	UploadContent(ctx context.Context, id int, r io.Reader, size int64) error

	GetObject(ctx context.Context, id int) (dto.FO, error)
}

type Config struct {
	Tx          repo.TxRunner
	InfoObjects repo.InfoObjects
	FS          fs.Storage
}

func NewManager(cfg Config) Manager {
	return &manager{
		tx: cfg.Tx,
		db: cfg.InfoObjects,
		fs: cfg.FS,
	}
}

//

type manager struct {
	tx repo.TxRunner
	db repo.InfoObjects
	fs fs.Storage
}

func (m *manager) CreateObject(ctx context.Context, obj dto) (int, error) {
	var (
		objectID int
		err      error
	)

	err = m.tx.RunWithTx(ctx, func(ctx context.Context) error {
		foDb := fo.ToDBType()

		foDb.FileName = makeFileName(fo)

		foID, err := m.fileObjects.Create(ctx, foDb)
		if err != nil {
			return fmt.Errorf("failed to add file meta to the db (%s, %s), %w", fo.Name, fo.ContentType, err)
		}

		for _, tag := range fo.Tags {
			tt := strings.ToLower(tag)

			tagID, err := m.tags.GetOrCreateTag(ctx, tt)
			if err != nil {
				return fmt.Errorf("failed to get or create tag (%s, %s), %w", fo.Name, tt, err)
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

	err = m.tx.RunWithTx(ctx, func(ctx context.Context) error {

		err = m.fileObjects.MarkAsUploaded(ctx, id, size)
		if err != nil {
			return fmt.Errorf("failed to mark object as uploaded %d, %w", id, err)
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

		return nil
	})

	log.FromContext(ctx).Infof("object uploaded, (%d, %s, %s, %d)", id, foDb.Name, foDb.ContentType, size)
	return nil
}

func (m *manager) GetObject(ctx context.Context, id int) (dto.FO, error) {
	foDb, err := m.fileObjects.GetById(ctx, id)
	if err != nil {
		return dto.FO{}, fmt.Errorf("failed to retrieve file object from db %d, %w", id, err)
	}

	tags, err := m.tags.GetFileObjectTags(ctx, id)
	if err != nil {
		return dto.FO{}, fmt.Errorf("failed to retrieve file object tags from db %d, %w", id, err)
	}

	fo := dto2.MakeFOFromDBType(foDb)
	fo.Tags = tags

	return fo, nil
}

// move this logic

func makeFileName(fo dto.FO) string {
	return fmt.Sprintf("%s_%s_%d", fo.Source, fo.Name, time.Now().UnixMilli())
}
