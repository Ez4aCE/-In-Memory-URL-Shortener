package models

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortCode string `json:"short_code"`
}

type URLData struct {
	URL    string `json:"url"`
	Clicks int    `json:"clicks"`
}
