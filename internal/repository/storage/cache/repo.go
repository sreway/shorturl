package cache

import (
	"context"
	"net/url"
	"os"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"

	entity "github.com/sreway/shorturl/internal/domain/url"
	"github.com/sreway/shorturl/internal/usecases/shortener"
)

type repo struct {
	data    map[uuid.UUID]*item
	file    *os.File
	fileUse bool
	logger  *slog.Logger
	mu      sync.RWMutex
}

func (r *repo) Add(ctx context.Context, id, userID [16]byte, value *url.URL) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_ = ctx

	r.data[id] = &item{
		UserID: userID,
		Value:  value,
	}
	return nil
}

func (r *repo) Get(_ context.Context, id [16]byte) (entity.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	i, ok := r.data[id]
	if !ok {
		return nil, shortener.ErrNotFound
	}
	return entity.NewURL(id, i.UserID, nil, i.Value), nil
}

func (r *repo) GetByUserID(_ context.Context, userID [16]byte) ([]entity.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := []entity.URL{}

	for k, v := range r.data {
		if v.UserID == userID {
			result = append(result, entity.NewURL(k, v.UserID, nil, v.Value))
		}
	}

	return result, nil
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

func (r *repo) Ping(_ context.Context) error {
	return ErrInvalidStorageType
}

func New(opts ...Option) *repo {
	log := slog.New(slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("repository", "cache")}))

	r := &repo{
		data:   map[uuid.UUID]*item{},
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
