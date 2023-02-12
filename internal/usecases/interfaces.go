package usecases

import (
	"context"

	entity "github.com/sreway/shorturl/internal/domain/url"
)

type Shortener interface {
	CreateURL(ctx context.Context, rawURL string) (*entity.URL, error)
	GetURL(ctx context.Context, urlID string) (*entity.URL, error)
}
