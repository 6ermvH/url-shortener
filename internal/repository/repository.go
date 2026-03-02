package repository

import "context"

type Repository interface {
	GetLongUrlByShortUrl(ctx context.Context, shortUrl string) (string, error)
	UpsertUrlMapping(ctx context.Context, shortUrl, longUrl string) error
}
