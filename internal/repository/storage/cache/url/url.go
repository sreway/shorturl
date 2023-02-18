package url

import (
	"context"
	"net/url"
	"sync"
)

type repo struct {
	data map[uint64]*url.URL
	mu   sync.RWMutex
}

func (r *repo) Add(ctx context.Context, id uint64, longURL *url.URL) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_ = ctx

	r.data[id] = longURL
	return nil
}

func (r *repo) Get(ctx context.Context, id uint64) (*url.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_ = ctx

	v, ok := r.data[id]
	if !ok {
		return nil, ErrNotFound
	}

	return v, nil
}

func New() *repo {
	return &repo{
		data: map[uint64]*url.URL{},
	}
}
