package repo

import "context"

//go:generate mockery --name Tags --outpkg repomocks --output ./repomocks --dir .
type Tags interface {
	GetOrCreate(ctx context.Context, tag string) (int, error)

	LinkTag(ctx context.Context, foID, tagID int) error
}
