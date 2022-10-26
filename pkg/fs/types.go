package fs

import (
	"context"
	"io"
)

type Object struct {
	Name        string
	ContentType string
	Size        int64
	Reader      io.Reader
}

//go:generate mockery --name Storage --outpkg fsmocks --output ./fsmocks --dir .
type Storage interface {
	Prepare(ctx context.Context) error

	CreateObject(ctx context.Context, fo Object) error
	GetObject(ctx context.Context, name string) (Object, error)
	DeleteObject(ctx context.Context, name string) error
}
