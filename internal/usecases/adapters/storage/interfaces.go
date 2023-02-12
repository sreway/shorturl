package storage

import (
	"context"
	"net/url"
)

type URL interface {
	Add(ctx context.Context, id uint64, longURL *url.URL) error
	Get(ctx context.Context, id uint64) (longURL *url.URL, err error)
}
