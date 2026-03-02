package service

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/6ermvH/url-shortener/internal/handler"
	"github.com/6ermvH/url-shortener/internal/repository"
	"github.com/6ermvH/url-shortener/pkg/base63"
)

const (
	shortURLLen = 10
	alphabet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

type Service struct {
	repo     repository.Repository
	encoding *base63.Encoding
}

func New(repo repository.Repository) *Service {
	return &Service{
		repo:     repo,
		encoding: base63.NewEncoding(alphabet, shortURLLen),
	}
}

func (s *Service) Shorten(
	ctx context.Context,
	req handler.ShortenRequest,
) (handler.ShortenResponse, error) {
	short := s.generateShort(req.URL, 0)

	err := s.repo.Save(ctx, repository.URLMapping{
		ShortURL:    short,
		OriginalURL: req.URL,
	})
	if err != nil {
		return handler.ShortenResponse{}, fmt.Errorf("save url mapping: %w", err)
	}

	return handler.ShortenResponse{ShortURL: short}, nil
}

func (s *Service) Resolve(ctx context.Context, short string) (handler.ResolveResponse, error) {
	m, err := s.repo.GetByShort(ctx, short)
	if err != nil {
		return handler.ResolveResponse{}, fmt.Errorf("get by short: %w", err)
	}

	return handler.ResolveResponse{OriginalURL: m.OriginalURL}, nil
}

func (s *Service) generateShort(url string, attempt int) string {
	input := url
	if attempt > 0 {
		input = fmt.Sprintf("%s#%d", url, attempt)
	}

	hash := sha256.Sum256([]byte(input))

	return s.encoding.EncodeToString(hash[:8])
}
