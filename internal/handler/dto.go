package handler

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortURL string `json:"shortUrl"`
}

type resolveResponse struct {
	OriginalURL string `json:"originalUrl"`
}

type errorResponse struct {
	Error string `json:"error"`
}
