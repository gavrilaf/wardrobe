package minio

import (
	"context"
	"fmt"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/gavrilaf/wardrobe/pkg/fs"
)

func New(endpoint, user, password string) (fs.Storage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(user, password, ""),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create minio, %w", err)
	}

	return &storage{client: client}, nil
}

type storage struct {
	client *minio.Client
}

//

func (s *storage) Ping() error {
	url := fmt.Sprintf("%s/minio/health/live", s.client.EndpointURL().String())
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}
	return nil
}

func (s *storage) IsBucketExists(ctx context.Context, name string) (bool, error) {
	return s.client.BucketExists(ctx, name)
}

func (s *storage) CreateBucket(ctx context.Context, name string) error {
	return s.client.MakeBucket(ctx, name, minio.MakeBucketOptions{})
}

func (s *storage) CreateFile(ctx context.Context, f fs.File) error {
	opts := minio.PutObjectOptions{
		ContentType: f.ContentType,
	}

	_, err := s.client.PutObject(ctx, f.Bucket, f.Name, f.Reader, f.Size, opts)
	return err
}

func (s *storage) GetFile(ctx context.Context, bucketName, fileName string) (fs.File, error) {
	obj, err := s.client.GetObject(ctx, bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return fs.File{}, fmt.Errorf("minio GetObject failed (%s, %s), %w", bucketName, fileName, err)
	}

	stat, err := obj.Stat() // TODO: optimization point (??)
	if err != nil {
		return fs.File{}, fmt.Errorf("minio object stat failed (%s, %s), %w", bucketName, fileName, err)
	}

	return fs.File{
		Bucket:      bucketName,
		Name:        fileName,
		ContentType: stat.ContentType,
		Size:        stat.Size,
		Reader:      obj,
	}, nil
}
