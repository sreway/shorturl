package http

type (
	shortURLRequest struct {
		URL string `json:"url"`
	}
	batchURLRequest struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}
)
