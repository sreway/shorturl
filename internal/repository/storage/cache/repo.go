package cache

import (
	"context"
	"net/url"
	"os"
	"sync"
	"sync/atomic"

	"golang.org/x/exp/slog"
)

type repo struct {
	data    map[uint64]*url.URL
	counter uint64
	file    *os.File
	fileUse bool
	logger  *slog.Logger
	mu      sync.RWMutex
}

func (r *repo) Add(ctx context.Context, longURL *url.URL) (uint64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_ = ctx
	id := atomic.AddUint64(&r.counter, 1)
	r.data[id] = longURL
	return id, nil
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

func (r *repo) Close() error {
	if !r.fileUse {
		return nil
	}

	err := r.fileClose()
	if err != nil {
		r.logger.Error("failed close url repository file", err)
	}

	r.logger.Info("success close url repository file")
	return nil
}

func New(opts ...Option) *repo {
	log := slog.New(slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("repository", "cache")}))

	r := &repo{
		data:   map[uint64]*url.URL{},
		logger: log,
	}

	for _, opt := range opts {
		err := opt(r)
		if err != nil {
			r.logger.Error("failed apply option", err)
		}
	}

	return r
}
