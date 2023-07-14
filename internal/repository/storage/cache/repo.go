// Package cache implements a repository for storing short URLs in the in-memory storage.
package cache

import (
	"context"
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

// Add implements saving short URL.
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

// Get implements getting short URL.
func (r *repo) Get(_ context.Context, id uuid.UUID) (entity.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	i, ok := r.data[id]
	if !ok {
		return nil, entity.ErrNotFound
	}

	u := entity.NewURL(id, i.UserID)
	u.SetLongURL(i.Value)
	u.SetDeleted(i.Deleted)
	return u, nil
}

// GetByUserID implements getting short URLs for user ID.
func (r *repo) GetByUserID(_ context.Context, userID uuid.UUID) ([]entity.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := []entity.URL{}

	for k, v := range r.data {
		if v.UserID == userID {
			u := entity.NewURL(k, v.UserID)
			u.SetLongURL(v.Value)
			u.SetDeleted(v.Deleted)
			result = append(result, u)
		}
	}

	return result, nil
}

// Close implements closing the connection to the file storage.
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

// Ping implements health check storage.
func (r *repo) Ping(_ context.Context) error {
	return ErrInvalidStorageType
}

// Batch implements saving multiple short URLs.
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

// BatchDelete implements the deletion multiple short URLs.
func (r *repo) BatchDelete(_ context.Context, urls []entity.URL) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, item := range urls {
		v, ok := r.data[item.ID()]
		if !ok {
			r.logger.Error("url not found", entity.ErrNotFound, slog.String("func", "BatchDelete"))
		}

		v.Deleted = true
		r.data[item.ID()] = v
	}

	return nil
}

// GetUserCount implements the getting user count.
func (r *repo) GetUserCount(_ context.Context) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	users := map[uuid.UUID]struct{}{}
	for _, v := range r.data {
		users[v.UserID] = struct{}{}
	}
	return len(users), nil
}

// GetURLCount implements the getting url count.
func (r *repo) GetURLCount(_ context.Context) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.data), nil
}

// New implements the creation of storage.
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
