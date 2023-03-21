package storage

import (
	"context"
	"net/url"

	entity "github.com/sreway/shorturl/internal/domain/url"
)

//go:generate mockgen -source=./internal/usecases/adapters/storage/interfaces.go -destination=./internal/usecases/adapters/storage/mocks/mock_url.go -package=storageMock
type URL interface {
	Add(ctx context.Context, id, userID [16]byte, value *url.URL) error
	Get(ctx context.Context, id [16]byte) (entity.URL, error)
	GetByUserID(ctx context.Context, userID [16]byte) ([]entity.URL, error)
	Close() error
	Ping(ctx context.Context) error
}
