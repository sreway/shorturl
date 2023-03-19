package shortener

import (
	"context"
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

	shortURL := &url.URL{
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

	err = uc.storage.Add(ctx, id, parsedUserID, longURL)
	if err != nil {
		uc.logger.Error("store url", err, slog.String("id", id.String()),
			slog.String("longURL", longURL.String()),
		)
		return nil, err
	}

	u := entity.NewURL(id, parsedUserID, shortURL, longURL)

	return u, nil
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

	longURL, userID, err := uc.storage.Get(ctx, id)
	if err != nil {
		uc.logger.Error("failed get url", err, slog.String("urlID", urlID))
	}
	shortURL := &url.URL{
		Scheme: uc.baseURL.Scheme,
		Host:   uc.baseURL.Host,
	}

	shortURL.Path = encodeUUID(id)

	return entity.NewURL(id, userID, shortURL, &longURL), nil
}

func (uc *useCase) GetUserURLs(ctx context.Context, userID string) ([]entity.URL, error) {
	parsedUserID, err := uuid.ParseBytes([]byte(userID))
	if err != nil {
		uc.logger.Error("failed parse RFC 4122 uuid from user id", err, slog.String("userID", userID))
		return nil, err
	}

	urls, err := uc.storage.GetByUserID(ctx, parsedUserID)
	if err != nil {
		uc.logger.Error("failed get url for user id", err, slog.String("userID", userID))
	}

	for idx, i := range urls {
		shortURL := &url.URL{
			Scheme: uc.baseURL.Scheme,
			Host:   uc.baseURL.Host,
		}
		shortURL.Path = encodeUUID(i.ID())
		urls[idx].SetShortURL(shortURL)
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
