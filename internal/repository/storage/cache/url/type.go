package url

import (
	"context"
	"net/url"
	"sync"
)

type Repo struct {
	data map[uint64]*url.URL
	mu   sync.RWMutex
}

func (r *Repo) Add(ctx context.Context, id uint64, longURL *url.URL) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_ = ctx

	r.data[id] = longURL
	return nil
}

func (r *Repo) Get(ctx context.Context, id uint64) (*url.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_ = ctx

	v, ok := r.data[id]
	if !ok {
		return nil, ErrNotFound
	}

	return v, nil
}

func New() *Repo {
	return &Repo{
		data: map[uint64]*url.URL{},
	}
}
