package service

type ShortenInput struct {
	URL string
}

type ShortenResult struct {
	ShortURL string
}

type ResolveResult struct {
	OriginalURL string
}
