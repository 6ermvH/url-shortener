package repository

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("url not found")

type Repository interface {
	GetByShort(ctx context.Context, short string) (URLMapping, error)
	Save(ctx context.Context, urlMap URLMapping) error
}
