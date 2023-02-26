package storage

import (
	"context"
	"net/url"
)

//go:generate mockgen -source=./internal/usecases/adapters/storage/interfaces.go -destination=./internal/usecases/adapters/storage/mocks/mock_url.go -package=storageMock
type URL interface {
	Add(ctx context.Context, longURL *url.URL) (uint64, error)
	Get(ctx context.Context, id uint64) (longURL *url.URL, err error)
	Close() error
}
