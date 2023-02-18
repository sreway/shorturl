package shortener

import (
	"context"
	"net/url"
	"os"
	"sync/atomic"

	"github.com/sreway/shorturl/internal/config"

	"golang.org/x/exp/slog"

	entity "github.com/sreway/shorturl/internal/domain/url"
	"github.com/sreway/shorturl/internal/usecases/adapters/storage"
)

type useCase struct {
	baseURL *url.URL
	counter uint64
	storage storage.URL
	logger  *slog.Logger
}

func (uc *useCase) CreateURL(ctx context.Context, rawURL string) (entity.URL, error) {
	longURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		uc.logger.Error("parse long url", err, slog.String("longURL", rawURL))
		return nil, ErrParseURL
	}

	id := atomic.AddUint64(&uc.counter, 1)

	shortURL := &url.URL{
		Scheme: uc.baseURL.Scheme,
		Host:   uc.baseURL.Host,
	}

	shortURL.Path = uintEncode(id)

	err = uc.storage.Add(ctx, id, longURL)
	if err != nil {
		uc.logger.Error("store url", err, slog.Uint64("id", id),
			slog.String("longURL", longURL.String()),
		)
		return nil, err
	}

	return entity.NewURL(id, shortURL, longURL), nil
}

func (uc *useCase) GetURL(ctx context.Context, urlID string) (entity.URL, error) {
	id, err := uintDecode(urlID)
	if err != nil {
		uc.logger.Error("decode short url", err)
		return nil, ErrDecodeURL
	}

	shortURL := &url.URL{
		Scheme: uc.baseURL.Scheme,
		Host:   uc.baseURL.Host,
	}

	shortURL.Path = urlID

	longURL, err := uc.storage.Get(ctx, id)
	if err != nil {
		uc.logger.Error("get storage url", err, slog.Uint64("id", id),
			slog.String("shortURL", shortURL.String()),
		)
		return nil, err
	}
	return entity.NewURL(id, shortURL, longURL), nil
}

func New(s storage.URL, cfg config.ShortURL) *useCase {
	log := slog.New(slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("service", "shortener")}))
	return &useCase{
		counter: cfg.GetCounter(),
		baseURL: cfg.GetBaseURL(),
		storage: s,
		logger:  log,
	}
}
