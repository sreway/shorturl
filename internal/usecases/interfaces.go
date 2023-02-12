package usecases

import (
	"context"

	entity "github.com/sreway/shorturl/internal/domain/url"
)

//go:generate mockgen -source=./internal/usecases/interfaces.go -destination=./internal/usecases/mocks/mock_usecases.go -package=usecaseMock
type Shortener interface {
	CreateURL(ctx context.Context, rawURL string) (*entity.URL, error)
	GetURL(ctx context.Context, urlID string) (*entity.URL, error)
}
