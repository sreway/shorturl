package http

import (
	"net/http"

	"github.com/go-chi/render"
)

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
	errResponse struct {
		Err            error  `json:"-"`
		HTTPStatusCode int    `json:"-"`
		ErrorText      string `json:"error,omitempty"`
	}
)

// Render renders a single payload and respond to the client request.
func (er *errResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, er.HTTPStatusCode)
	return nil
}
