package cache

import (
	"context"
	"net/url"
	"os"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"

	entity "github.com/sreway/shorturl/internal/domain/url"
)

type repo struct {
	data    map[uuid.UUID]storageURL
	file    *os.File
	fileUse bool
	logger  *slog.Logger
	mu      sync.RWMutex
}

func (r *repo) Add(ctx context.Context, item entity.URL) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_ = ctx

	r.data[item.ID()] = storageURL{
		UserID: item.UserID(),
		Value:  item.LongValue(),
	}
	return nil
}

func (r *repo) Get(_ context.Context, id uuid.UUID) (entity.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	i, ok := r.data[id]
	if !ok {
		return nil, entity.ErrNotFound
	}
	return entity.NewURL(uuid.UUID{}, i.UserID, url.URL{}, i.Value), nil
}

func (r *repo) GetByUserID(_ context.Context, userID uuid.UUID) ([]entity.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := []entity.URL{}

	for k, v := range r.data {
		if v.UserID == userID {
			u := entity.NewURL(k, v.UserID, url.URL{}, v.Value)
			result = append(result, u)
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

func (r *repo) Batch(_ context.Context, urls []entity.URL) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, item := range urls {
		r.data[item.ID()] = storageURL{
			UserID: item.UserID(),
			Value:  item.LongValue(),
		}
	}
	return nil
}

func New(opts ...Option) *repo {
	log := slog.New(slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("repository", "cache")}))

	r := &repo{
		data:   map[uuid.UUID]storageURL{},
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
