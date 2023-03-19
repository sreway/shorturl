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

func (r *repo) Get(ctx context.Context, id [16]byte) (value url.URL, userID [16]byte, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_ = ctx

	i, ok := r.data[id]
	if !ok {
		return url.URL{}, [16]byte{}, ErrNotFound
	}

	return *i.Value, i.UserID, nil
}

func (r *repo) GetByUserID(ctx context.Context, userID [16]byte) ([]entity.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	_ = ctx
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
