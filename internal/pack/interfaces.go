package pack

import "context"

type Repository interface {
	GetPackSizes(ctx context.Context) ([]int, error)
	InsertPackSize(ctx context.Context, size int) error
	DeletePackSize(ctx context.Context, size int) error
}
