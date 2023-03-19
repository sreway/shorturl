package usecases

import (
	"context"

	"github.com/sreway/shorturl/internal/domain/url"
)

//go:generate mockgen -source=./internal/usecases/interfaces.go -destination=./internal/usecases/mocks/mock_usecases.go -package=usecaseMock
type Shortener interface {
	CreateURL(ctx context.Context, rawURL string, userID string) (url.URL, error)
	GetURL(ctx context.Context, urlID string) (url.URL, error)
	GetUserURLs(ctx context.Context, userID string) ([]url.URL, error)
}
