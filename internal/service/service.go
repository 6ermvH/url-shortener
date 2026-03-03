package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/url"

	"github.com/6ermvH/url-shortener/internal/repository"
	"github.com/6ermvH/url-shortener/pkg/base63"
)

const (
	shortURLLen = 10
	alphabet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

var (
	ErrNotFound   = errors.New("url not found")
	ErrEmptyURL   = errors.New("url is required")
	ErrInvalidURL = errors.New("invalid url")
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

func (s *Service) Shorten(ctx context.Context, input ShortenInput) (ShortenResult, error) {
	err := validateURL(input.URL)
	if err != nil {
		return ShortenResult{}, err
	}

	short := s.generateShort(input.URL)

	err = s.repo.Save(ctx, repository.URLMapping{
		ShortURL:    short,
		OriginalURL: input.URL,
	})
	if err != nil {
		return ShortenResult{}, fmt.Errorf("save url mapping: %w", err)
	}

	return ShortenResult{ShortURL: short}, nil
}

func (s *Service) Resolve(ctx context.Context, short string) (ResolveResult, error) {
	mapping, err := s.repo.GetByShort(ctx, short)
	if errors.Is(err, repository.ErrNotFound) {
		return ResolveResult{}, ErrNotFound
	}

	if err != nil {
		return ResolveResult{}, fmt.Errorf("get by short: %w", err)
	}

	return ResolveResult{OriginalURL: mapping.OriginalURL}, nil
}

func (s *Service) generateShort(rawURL string) string {
	hash := sha256.Sum256([]byte(rawURL))

	return s.encoding.EncodeToString(hash[:8])
}

func validateURL(rawURL string) error {
	if rawURL == "" {
		return ErrEmptyURL
	}

	u, err := url.ParseRequestURI(rawURL)
	if err != nil || u.Host == "" {
		return fmt.Errorf("%w: %q", ErrInvalidURL, rawURL)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("%w: scheme must be http or https", ErrInvalidURL)
	}

	return nil
}
