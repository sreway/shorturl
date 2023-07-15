// Package shortener implements a URL shortening service.
package shortener

import (
	"context"
	"errors"
	"net/url"
	"os"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"

	"github.com/sreway/shorturl/internal/config"
	"github.com/sreway/shorturl/internal/domain/stats"
	entity "github.com/sreway/shorturl/internal/domain/url"
	"github.com/sreway/shorturl/internal/usecases/adapters/storage"
)

type (
	useCase struct {
		baseURL   *url.URL
		storage   storage.URL
		logger    *slog.Logger
		taskQueue chan task
	}
)

// CreateURL implements the creation of a short URL.
func (uc *useCase) CreateURL(ctx context.Context, rawURL string, userID string) (entity.URL, error) {
	longURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		uc.logger.Error("parse long url", err, slog.String("longURL", rawURL))
		return nil, ErrParseURL
	}

	shortURL := url.URL{
		Scheme: uc.baseURL.Scheme,
		Host:   uc.baseURL.Host,
	}

	parsedUserID, err := uuid.ParseBytes([]byte(userID))
	if err != nil {
		uc.logger.Error("failed parse RFC 4122 uuid from user id", err, slog.String("userID", userID))
		return nil, err
	}

	id := uuid.New()

	shortURL.Path = encodeUUID(id)

	addURL := entity.NewURL(id, parsedUserID)
	addURL.SetShortURL(shortURL)
	addURL.SetLongURL(*longURL)

	err = uc.storage.Add(ctx, addURL)
	if err != nil && !errors.Is(err, entity.ErrAlreadyExist) {
		uc.logger.Error("store url", err, slog.String("id", id.String()),
			slog.String("longURL", longURL.String()),
		)
		return nil, err
	}
	if errors.Is(err, entity.ErrAlreadyExist) {
		uc.logger.Error("store url", err, slog.String("id", id.String()),
			slog.String("longURL", longURL.String()),
		)

		var errURL *entity.ErrURL
		if errors.As(err, &errURL) {
			shortURL.Path = encodeUUID(errURL.ID())
			addURL.SetShortURL(shortURL)
			return addURL, err
		}
		return nil, err
	}

	return addURL, nil
}

// GetURL implements getting short URL.
func (uc *useCase) GetURL(ctx context.Context, urlID string) (entity.URL, error) {
	decoded, err := decodeUUID(urlID)
	if err != nil {
		uc.logger.Error("decode short url", err)
		return nil, ErrDecodeURL
	}

	id, err := uuid.FromBytes(decoded)
	if err != nil {
		uc.logger.Error("failed create uuid from url id", err, slog.String("urlID", urlID))
		return nil, ErrParseUUID
	}

	u, err := uc.storage.Get(ctx, id)
	if err != nil {
		uc.logger.Error("failed get url", err, slog.String("urlID", urlID))
		return nil, err
	}

	if u.Deleted() {
		return nil, entity.ErrDeleted
	}

	shortURL := url.URL{
		Scheme: uc.baseURL.Scheme,
		Host:   uc.baseURL.Host,
	}

	shortURL.Path = encodeUUID(id)

	u.SetShortURL(shortURL)

	return u, nil
}

// GetUserURLs implements getting short URLs for user ID.
func (uc *useCase) GetUserURLs(ctx context.Context, userID string) ([]entity.URL, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		uc.logger.Error("failed parse RFC 4122 uuid from user id", err, slog.String("userID", userID))
		return nil, ErrParseUUID
	}
	urls, err := uc.storage.GetByUserID(ctx, parsedUserID)
	if err != nil {
		uc.logger.Error("failed get url for user id", err, slog.String("userID", userID))
		return nil, err
	}

	for idx, i := range urls {
		shortURL := url.URL{
			Scheme: uc.baseURL.Scheme,
			Host:   uc.baseURL.Host,
		}
		shortURL.Path = encodeUUID(i.ID())
		urls[idx].SetShortURL(shortURL)
	}

	return urls, nil
}

// StorageCheck implements storage health check.
func (uc *useCase) StorageCheck(ctx context.Context) error {
	return uc.storage.Ping(ctx)
}

// BatchURL implements the creation of several short URLs.
func (uc *useCase) BatchURL(ctx context.Context, correlationID, rawURL []string, userID string) ([]entity.URL, error) {
	urls := []entity.URL{}

	for idx, item := range rawURL {
		longURL, err := url.ParseRequestURI(item)
		if err != nil {
			uc.logger.Error("parse long url", err, slog.String("BatchURL", item))
			return nil, ErrParseURL
		}

		shortURL := url.URL{
			Scheme: uc.baseURL.Scheme,
			Host:   uc.baseURL.Host,
		}

		parsedUserID, err := uuid.ParseBytes([]byte(userID))
		if err != nil {
			uc.logger.Error("failed parse RFC 4122 uuid from user id", err, slog.String("userID", userID))
			return nil, ErrParseUUID
		}

		id := uuid.New()
		shortURL.Path = encodeUUID(id)

		u := entity.NewURL(id, parsedUserID)
		u.SetShortURL(shortURL)
		u.SetLongURL(*longURL)
		u.SetCorrelationID(correlationID[idx])
		urls = append(urls, u)
	}

	err := uc.storage.Batch(ctx, urls)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

// DeleteURL implements the deletion multiple short URLs.
func (uc *useCase) DeleteURL(_ context.Context, userID string, urlID []string) error {
	urls := []entity.URL{}
	parsedUserID, err := uuid.ParseBytes([]byte(userID))
	if err != nil {
		uc.logger.Error("failed parse RFC 4122 uuid from user id", err, slog.String("userID", userID))
		return ErrParseUUID
	}

	for _, i := range urlID {
		decoded, err := decodeUUID(i)
		if err != nil {
			uc.logger.Error("decode short url", err)
			return ErrDecodeURL
		}

		id, err := uuid.FromBytes(decoded)
		if err != nil {
			uc.logger.Error("failed create uuid from url id", err, slog.String("urlID", i))
			return ErrParseUUID
		}

		u := entity.NewURL(id, parsedUserID)
		u.SetDeleted(true)
		urls = append(urls, u)
	}

	if len(uc.taskQueue) == cap(uc.taskQueue) {
		return ErrTaskBufferFull
	}

	uc.taskQueue <- *NewTask(deleteAction, urls)

	return nil
}

// GetStats implements getting stats of the short URLs service.
func (uc *useCase) GetStats(ctx context.Context) (stats.Collection, error) {
	userCount, err := uc.storage.GetUserCount(ctx)
	if err != nil {
		uc.logger.Error("failed get stats user count", err)
		return nil, err
	}

	urlCount, err := uc.storage.GetURLCount(ctx)
	if err != nil {
		uc.logger.Error("failed get stats url count", err)
		return nil, err
	}

	userStats := stats.NewUserStats(stats.UserCount(userCount))
	urlStats := stats.NewURLStats(stats.URLCount(urlCount))
	collection := stats.NewCollectionStats(userStats, urlStats)
	return collection, nil
}

// New implements the creation of a URL shortening service.
func New(s storage.URL, cfg config.ShortURL) *useCase {
	log := slog.New(slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("service", "shortener")}))
	taskQueue := make(chan task, cfg.GetMaxTaskQueue())
	return &useCase{
		baseURL:   cfg.GetBaseURL(),
		storage:   s,
		logger:    log,
		taskQueue: taskQueue,
	}
}
