package storage

import (
	"context"
	"net/url"
)

//go:generate mockgen -source=./internal/usecases/adapters/storage/interfaces.go -destination=./internal/usecases/adapters/storage/mocks/mock_url.go -package=storageMock
type URL interface {
	Add(ctx context.Context, id, userID [16]byte, value *url.URL) error
	Get(ctx context.Context, id [16]byte) (value url.URL, userID [16]byte, err error)
	GetByUserID(ctx context.Context, userID [16]byte) (map[[16]byte]url.URL, error)
	Close() error
}
