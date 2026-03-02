package handler

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

type ResolveResponse struct {
	OriginalURL string `json:"original_url"`
}
