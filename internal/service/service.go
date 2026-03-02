package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/url"

	"github.com/6ermvH/url-shortener/internal/handler"
	"github.com/6ermvH/url-shortener/internal/repository"
	"github.com/6ermvH/url-shortener/pkg/base63"
)

const (
	shortURLLen = 10
	alphabet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

var ErrNotFound = errors.New("url not found")

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
	err := validateURL(req.URL)
	if err != nil {
		return handler.ShortenResponse{}, err
	}

	for attempt := 0; ; attempt++ {
		short := s.generateShort(req.URL, attempt)

		err := s.repo.Save(ctx, repository.URLMapping{
			ShortURL:    short,
			OriginalURL: req.URL,
		})
		if errors.Is(err, repository.ErrConflict) {
			continue
		}

		if err != nil {
			return handler.ShortenResponse{}, fmt.Errorf("save url mapping: %w", err)
		}

		return handler.ShortenResponse{ShortURL: short}, nil
	}
}

func (s *Service) Resolve(ctx context.Context, short string) (handler.ResolveResponse, error) {
	m, err := s.repo.GetByShort(ctx, short)
	if errors.Is(err, repository.ErrNotFound) {
		return handler.ResolveResponse{}, ErrNotFound
	}

	if err != nil {
		return handler.ResolveResponse{}, fmt.Errorf("get by short: %w", err)
	}

	return handler.ResolveResponse{OriginalURL: m.OriginalURL}, nil
}

func (s *Service) generateShort(rawURL string, attempt int) string {
	input := rawURL
	if attempt > 0 {
		input = fmt.Sprintf("%s#%d", rawURL, attempt)
	}

	hash := sha256.Sum256([]byte(input))

	return s.encoding.EncodeToString(hash[:8])
}

func validateURL(rawURL string) error {
	if rawURL == "" {
		return errors.New("url is required")
	}

	u, err := url.ParseRequestURI(rawURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("invalid url: %q", rawURL)
	}

	return nil
}
