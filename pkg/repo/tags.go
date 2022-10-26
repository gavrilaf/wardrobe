package repo

import "context"

type Tags interface {
	GetOrCreate(ctx context.Context, tag string) (int, error)

	LinkTag(ctx context.Context, foID, tagID int) error
}
