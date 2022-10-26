package minio

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/gavrilaf/wardrobe/pkg/fs"
	"github.com/gavrilaf/wardrobe/pkg/utils/log"
)

func New(endpoint, user, password, bucket string) (fs.Storage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(user, password, ""),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect minio, %w", err)
	}

	return &storage{
		bucketName: bucket,
		client:     client,
	}, nil
}

type storage struct {
	bucketName string
	client     *minio.Client
}

func (s *storage) Prepare(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket %s existing, %w", s.bucketName, err)
	}

	if !exists {
		err = s.client.MakeBucket(ctx, s.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket %s, %w", s.bucketName, err)
		}

		log.FromContext(ctx).Infof("bucket %s created", s.bucketName)
	}

	return nil
}

func (s *storage) CreateObject(ctx context.Context, fo fs.Object) error {
	opts := minio.PutObjectOptions{
		ContentType: fo.ContentType,
	}

	info, err := s.client.PutObject(ctx, s.bucketName, fo.Name, fo.Reader, fo.Size, opts)
	if err != nil {
		return fmt.Errorf("failed to upload object %s, %w", fo.Name, err)
	}

	log.FromContext(ctx).Infof("object uploaded (%s, %s, %s, %d)", info.Bucket, info.Location, info.Key, info.Size)

	return nil
}

func (s *storage) GetObject(ctx context.Context, name string) (fs.Object, error) {
	//TODO implement me
	panic("implement me")
}

func (s *storage) DeleteObject(ctx context.Context, name string) error {
	//TODO implement me
	panic("implement me")
}
