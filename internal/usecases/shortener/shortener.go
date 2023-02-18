package shortener

import (
	"context"
	"net/url"
	"os"
	"sync/atomic"

	"golang.org/x/exp/slog"

	"github.com/sreway/shorturl/config"

	entity "github.com/sreway/shorturl/internal/domain/url"
	"github.com/sreway/shorturl/internal/usecases/adapters/storage"
	"github.com/sreway/shorturl/pkg/tools/base62"
)

type UseCase struct {
	cfg     *config.ShortURL
	counter uint64
	storage storage.URL
	logger  *slog.Logger
}

func (uc *UseCase) CreateURL(ctx context.Context, rawURL string) (*entity.URL, error) {
	longURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		uc.logger.Error("parse long url", err, slog.String("longURL", rawURL))
		return nil, ErrParseURL
	}

	id := atomic.AddUint64(&uc.counter, 1)

	shortURL := &url.URL{
		Scheme: uc.cfg.BaseURL.Scheme,
		Host:   uc.cfg.BaseURL.Host,
	}

	shortURL.Path = base62.UIntEncode(id)

	err = uc.storage.Add(ctx, id, longURL)
	if err != nil {
		uc.logger.Error("store url", err, slog.Uint64("id", id),
			slog.String("longURL", longURL.String()),
		)
		return nil, err
	}

	return &entity.URL{
		ID:       id,
		LongURL:  longURL,
		ShortURL: shortURL,
	}, nil
}

func (uc *UseCase) GetURL(ctx context.Context, urlID string) (*entity.URL, error) {
	id, err := base62.UIntDecode(urlID)
	if err != nil {
		uc.logger.Error("decode short url", err)
		return nil, ErrDecodeURL
	}

	shortURL := &url.URL{
		Scheme: uc.cfg.BaseURL.Scheme,
		Host:   uc.cfg.BaseURL.Host,
	}

	shortURL.Path = urlID

	longURL, err := uc.storage.Get(ctx, id)
	if err != nil {
		uc.logger.Error("get storage url", err, slog.Uint64("id", id),
			slog.String("shortURL", shortURL.String()),
		)
		return nil, err
	}

	return &entity.URL{
		ID:       id,
		LongURL:  longURL,
		ShortURL: shortURL,
	}, nil
}

func New(s storage.URL, cfg *config.ShortURL) *UseCase {
	log := slog.New(slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("service", "shortener")}))
	return &UseCase{
		counter: cfg.Counter,
		cfg:     cfg,
		storage: s,
		logger:  log,
	}
}
