// Package usecases describes the services used in the application.
package usecases

import (
	"context"

	"github.com/sreway/shorturl/internal/domain/stats"

	"github.com/sreway/shorturl/internal/domain/url"
)

// Shortener describes the implementation of the URL shortening service.
//
//go:generate mockgen -source=./internal/usecases/interfaces.go -destination=./internal/usecases/mock/mock_usecases.go -package=usecaseMock
type Shortener interface {
	CreateURL(ctx context.Context, rawURL string, userID string) (url.URL, error)
	BatchURL(ctx context.Context, correlationID, rawURL []string, userID string) ([]url.URL, error)
	GetURL(ctx context.Context, urlID string) (url.URL, error)
	GetUserURLs(ctx context.Context, userID string) ([]url.URL, error)
	DeleteURL(ctx context.Context, userID string, urlID []string) error
	StorageCheck(ctx context.Context) error
	GetStats(ctx context.Context) (stats.Collection, error)
}
