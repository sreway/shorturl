package shortener

import (
	"context"
	"errors"
	"net/url"
	"os"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"

	"github.com/sreway/shorturl/internal/config"
	entity "github.com/sreway/shorturl/internal/domain/url"
	"github.com/sreway/shorturl/internal/usecases/adapters/storage"
)

type useCase struct {
	baseURL *url.URL
	storage storage.URL
	logger  *slog.Logger
}

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

	addURL := entity.NewURL(id, parsedUserID, shortURL, *longURL)
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
	}

	return addURL, nil
}

func (uc *useCase) GetURL(ctx context.Context, urlID string) (entity.URL, error) {
	decoded, err := decodeUUID(urlID)
	if err != nil {
		uc.logger.Error("decode short url", err)
		return nil, ErrDecodeURL
	}

	id, err := uuid.FromBytes(decoded)
	if err != nil {
		uc.logger.Error("failed create uuid from url id", err, slog.String("urlID", urlID))
	}

	u, err := uc.storage.Get(ctx, id)
	if err != nil {
		uc.logger.Error("failed get url", err, slog.String("urlID", urlID))
		return nil, err
	}
	shortURL := url.URL{
		Scheme: uc.baseURL.Scheme,
		Host:   uc.baseURL.Host,
	}

	shortURL.Path = encodeUUID(id)

	u.SetShortURL(shortURL)

	return u, nil
}

func (uc *useCase) GetUserURLs(ctx context.Context, userID string) ([]entity.URL, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		uc.logger.Error("failed parse RFC 4122 uuid from user id", err, slog.String("userID", userID))
		return nil, err
	}
	urls, err := uc.storage.GetByUserID(ctx, parsedUserID)
	if err != nil {
		uc.logger.Error("failed get url for user id", err, slog.String("userID", userID))
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

func (uc *useCase) StorageCheck(ctx context.Context) error {
	return uc.storage.Ping(ctx)
}

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
			return nil, err
		}

		id := uuid.New()
		shortURL.Path = encodeUUID(id)

		u := entity.NewURL(id, parsedUserID, shortURL, *longURL)
		u.SetCorrelationID(correlationID[idx])
		urls = append(urls, u)
	}

	err := uc.storage.Batch(ctx, urls)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func New(s storage.URL, cfg config.ShortURL) *useCase {
	log := slog.New(slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("service", "shortener")}))
	return &useCase{
		baseURL: cfg.GetBaseURL(),
		storage: s,
		logger:  log,
	}
}
