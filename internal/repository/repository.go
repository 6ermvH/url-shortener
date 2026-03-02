package repository

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("url not found")
	ErrConflict = errors.New("short url already taken")
)

type Repository interface {
	GetByShort(ctx context.Context, short string) (URLMapping, error)
	Save(ctx context.Context, urlMap URLMapping) error
}
