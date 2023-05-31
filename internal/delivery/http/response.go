package http

type (
	shortURLResponse struct {
		Result string `json:"result"`
	}
	userURLResponse struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}
	batchURLResponse struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}
)
