package storage

import (
	"context"

	"github.com/google/uuid"

	entity "github.com/sreway/shorturl/internal/domain/url"
)

//go:generate mockgen -source=./internal/usecases/adapters/storage/interfaces.go -destination=./internal/usecases/adapters/storage/mock/mock_url.go -package=storageMock
type URL interface {
	Add(ctx context.Context, url entity.URL) error
	Get(ctx context.Context, id uuid.UUID) (entity.URL, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]entity.URL, error)
	Ping(ctx context.Context) error
	Close() error
	Batch(ctx context.Context, urls []entity.URL) error
}
