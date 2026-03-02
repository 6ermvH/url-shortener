package service

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"shortUrl"`
}

type ResolveResponse struct {
	OriginalURL string `json:"originalUrl"`
}
