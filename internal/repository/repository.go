package repository

import "context"

type Repository interface {
	GetByShort(ctx context.Context, short string) (URLMapping, error)
	Save(ctx context.Context, urlMap URLMapping) error
}
