package models

// ShortenRequest represents the request to shorten a URL
type ShortenRequest struct {
	LongURL string `json:"long_url" validate:"required,url"`
}

// ShortenResponse represents the response after shortening a URL
type ShortenResponse struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

// ErrorResponse represents the response after an error
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// StatsResponse represents service statistics
type StatsResponse struct {
	TotalURLsShortened int `json:"total_urls_shortened"`
	TotalRedirects     int `json:"total_redirects"`
}
