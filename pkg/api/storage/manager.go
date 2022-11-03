package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/gavrilaf/wardrobe/pkg/domain/dto"
	"github.com/gavrilaf/wardrobe/pkg/fs"
	"github.com/gavrilaf/wardrobe/pkg/repo"
	"github.com/gavrilaf/wardrobe/pkg/utils/log"
)

type Manager interface {
	CreateInfoObject(ctx context.Context, obj dto.InfoObject) (int, error)
	AddFile(ctx context.Context, infoObjID int, fileMeta dto.File, r io.Reader) (int, error)
	FinilizeInfoObject(ctx context.Context, id int) error

	GetInfoObject(ctx context.Context, id int) (dto.InfoObject, error)
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

func (m *manager) CreateInfoObject(ctx context.Context, obj dto.InfoObject) (int, error) {
	var objectID int

	txErr := m.tx.RunWithTx(ctx, func(ctx context.Context) error {
		dbObj, err := obj.ToDBType()
		if err != nil {
			return fmt.Errorf("failed to create db info object, %w", err)
		}

		objectID, err = m.db.CreateInfoObject(ctx, dbObj)
		if err != nil {
			return fmt.Errorf("failed to create info object %v, %w", obj, err)
		}

		for _, tag := range obj.Tags {
			tagLower := strings.ToLower(tag)

			err = m.db.LinkTag(ctx, objectID, tagLower)
			if err != nil {
				return fmt.Errorf("failed to link tag %s to the info object %d, %w", tagLower, objectID, err)
			}
		}

		return nil
	})

	if txErr != nil {
		return 0, txErr
	}

	log.FromContext(ctx).Infof("info object created (%d, %s)", objectID, obj.Name)

	return objectID, nil
}

func (m *manager) AddFile(ctx context.Context, infoObjID int, fileMeta dto.File, r io.Reader) (int, error) {
	var fileID int

	txErr := m.tx.RunWithTx(ctx, func(ctx context.Context) error {
		obj, err := m.db.GetById(ctx, infoObjID)
		if err != nil {
			return fmt.Errorf("failed to get info object from db %d, %w", infoObjID, err)
		}

		if obj.Uploaded != nil {
			return fmt.Errorf("info object %d finalized, failed to add file, %w", infoObjID, err)
		}

		bucket := "some-bucket"

		dbFileMeta := fileMeta.ToDbType()
		dbFileMeta.InfoObjectID = infoObjID

		fileID, err = m.db.AddFile(ctx, dbFileMeta)
		if err != nil {
			return fmt.Errorf("failed to add file meta %v, %w", dbFileMeta, err)
		}

		file := fs.File{
			Bucket:      bucket,
			Name:        fileMeta.Name,
			ContentType: fileMeta.ContentType,
			Size:        fileMeta.Size,
			Reader:      io.NopCloser(r),
		}
		err = m.fs.CreateFile(ctx, file)
		if err != nil {
			return fmt.Errorf("failed to upload file to the storage (%d, %v), %w", infoObjID, fileMeta, err)
		}

		return nil
	})

	if txErr != nil {
		return 0, txErr
	}

	log.FromContext(ctx).Infof("added file %d to the info object %d, (%s, %s, %d)", fileID, infoObjID,
		fileMeta.Name, fileMeta.ContentType, fileMeta.Size)

	return fileID, nil
}

func (m *manager) FinilizeInfoObject(ctx context.Context, id int) error {
	if err := m.db.MarkUploaded(ctx, id); err != nil {
		return fmt.Errorf("failed to finilize info object %d", id)
	}

	log.FromContext(ctx).Infof("info object %d finalized", id)
	return nil
}

func (m *manager) GetInfoObject(ctx context.Context, id int) (dto.InfoObject, error) {
	dbObj, err := m.db.GetById(ctx, id)
	if err != nil {
		return dto.InfoObject{}, fmt.Errorf("failed to get info object  %d, %w", id, err)
	}

	dbFiles, err := m.db.GetFiles(ctx, id)
	if err != nil {
		return dto.InfoObject{}, fmt.Errorf("failed to get info object files %d, %w", id, err)
	}

	tags, err := m.db.GetTags(ctx, id)
	if err != nil {
		return dto.InfoObject{}, fmt.Errorf("failed to get info object tags %d, %w", id, err)
	}

	dtoFiles := make([]dto.File, 0, len(dbFiles))
	for _, dbFile := range dbFiles {
		dtoFiles = append(dtoFiles, dto.FileFromDBType(dbFile))
	}

	dtoObj := dto.InfoObjectFromDBType(dbObj)
	dtoObj.Files = dtoFiles
	dtoObj.Tags = tags

	return dtoObj, nil
}
