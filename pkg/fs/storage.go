package fs

import (
	"context"
)

//go:generate mockery --name Storage --outpkg fsmocks --output ./fsmocks --dir .
type Storage interface {
	Ping() error

	IsBucketExists(ctx context.Context, name string) (bool, error)
	CreateBucket(ctx context.Context, name string) error

	CreateFile(ctx context.Context, f File) error
	GetFile(ctx context.Context, bucketName, fileName string) (File, error)
}
