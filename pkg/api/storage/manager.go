package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/gavrilaf/wardrobe/pkg/domain/dto"
	"github.com/gavrilaf/wardrobe/pkg/domain/stglogic"
	"github.com/gavrilaf/wardrobe/pkg/fs"
	"github.com/gavrilaf/wardrobe/pkg/repo"
	"github.com/gavrilaf/wardrobe/pkg/utils/log"
)

type Manager interface {
	CreateInfoObject(ctx context.Context, obj dto.InfoObject) (int, error)
	FinalizeInfoObject(ctx context.Context, id int) error
	GetInfoObject(ctx context.Context, id int) (dto.InfoObject, error)

	AddFile(ctx context.Context, objID int, fileMeta dto.File, r io.Reader) (int, error)
	GetFile(ctx context.Context, fileID int) (fs.File, error)

	GetStat(ctx context.Context) (dto.Stat, error)
}

type Config struct {
	Tx              repo.TxRunner
	InfoObjects     repo.InfoObjects
	FS              fs.Storage
	StgConfigurator stglogic.Configurator
}

func NewManager(cfg Config) Manager {
	return &manager{
		tx:  cfg.Tx,
		db:  cfg.InfoObjects,
		fs:  cfg.FS,
		cnf: cfg.StgConfigurator,
	}
}

//

type manager struct {
	tx  repo.TxRunner
	db  repo.InfoObjects
	fs  fs.Storage
	cnf stglogic.Configurator
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

func (m *manager) AddFile(ctx context.Context, objID int, fileMeta dto.File, r io.Reader) (int, error) {
	var fileID int

	txErr := m.tx.RunWithTx(ctx, func(ctx context.Context) error {
		obj, err := m.db.GetById(ctx, objID)
		if err != nil {
			return fmt.Errorf("failed to get info object from db %d, %w", objID, err)
		}

		if obj.Finalized != nil {
			return fmt.Errorf("info object %d finalized, failed to add file", objID)
		}

		dbFileMeta := fileMeta.ToDbType()

		bucket := m.cnf.GetBucket(obj, dbFileMeta)
		fileName, err := m.cnf.GenerateFileName(obj, dbFileMeta)
		if err != nil {
			return fmt.Errorf("failed to generate file name (%v, %v), %w", obj, dbFileMeta, err)
		}

		dbFileMeta.InfoObjectID = objID
		dbFileMeta.Bucket = bucket
		dbFileMeta.Name = fileName

		fileID, err = m.db.AddFile(ctx, dbFileMeta)
		if err != nil {
			return fmt.Errorf("failed to add file meta %v, %w", dbFileMeta, err)
		}

		file := fs.File{
			Bucket:      bucket,
			Name:        fileName,
			ContentType: fileMeta.ContentType,
			Size:        fileMeta.Size,
			Reader:      io.NopCloser(r),
		}
		err = m.fs.CreateFile(ctx, file)
		if err != nil {
			return fmt.Errorf("failed to upload file to the storage (%d, %v), %w", objID, fileMeta, err)
		}

		return nil
	})

	if txErr != nil {
		return 0, txErr
	}

	log.FromContext(ctx).Infof("added file %d to the info object %d, (%s, %s, %d)", fileID, objID,
		fileMeta.Name, fileMeta.ContentType, fileMeta.Size)

	return fileID, nil
}

func (m *manager) FinalizeInfoObject(ctx context.Context, id int) error {
	if err := m.db.MarkUploaded(ctx, id); err != nil {
		return fmt.Errorf("failed to finalize info object %d", id)
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

func (m *manager) GetFile(ctx context.Context, fileID int) (fs.File, error) {
	fileMeta, err := m.db.GetFile(ctx, fileID)
	if err != nil {
		return fs.File{}, fmt.Errorf("failed to read file %d meta, %w", fileID, err)
	}

	file, err := m.fs.GetFile(ctx, fileMeta.Bucket, fileMeta.Name)
	if err != nil {
		return fs.File{}, fmt.Errorf("failed to read file %v from storage, %w", fileMeta, err)
	}

	return file, err
}

func (m *manager) GetStat(ctx context.Context) (dto.Stat, error) {
	dbStat, err := m.db.GetStat(ctx)
	if err != nil {
		return dto.Stat{}, fmt.Errorf("failed to calculate storage stat, %w", err)
	}

	return dto.StatFromDDType(dbStat), err
}
