package shortener

import (
	"context"
	"net/url"
	"sync/atomic"

	"github.com/sreway/shorturl/config"

	"go.uber.org/zap"

	entity "github.com/sreway/shorturl/internal/domain/url"
	"github.com/sreway/shorturl/internal/usecases/adapters/storage"
	"github.com/sreway/shorturl/pkg/tools/base62"
	"github.com/sreway/shorturl/pkg/tools/logger"
)

type UseCase struct {
	cfg     *config.ShortURL
	counter uint64
	storage storage.URL
	logger  *zap.Logger
}

func (uc *UseCase) CreateURL(ctx context.Context, rawURL string) (*entity.URL, error) {
	longURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		uc.logger.Error("parse long url",
			zap.Error(err),
			zap.String("longURL", rawURL),
		)
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
		uc.logger.Error("store url",
			zap.Error(err),
			zap.Uint64("id", id),
			zap.String("longURL", longURL.String()),
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
		uc.logger.Error("decode short url",
			zap.Error(err),
		)
		return nil, ErrDecodeURL
	}

	shortURL := &url.URL{
		Scheme: uc.cfg.BaseURL.Scheme,
		Host:   uc.cfg.BaseURL.Host,
	}

	shortURL.Path = urlID

	longURL, err := uc.storage.Get(ctx, id)
	if err != nil {
		uc.logger.Error("get storage url",
			zap.Error(err),
			zap.Uint64("id", id),
			zap.String("shortURL", shortURL.String()),
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
	l := logger.GetLogger()

	return &UseCase{
		counter: cfg.Counter,
		cfg:     cfg,
		storage: s,
		logger:  l.With(zap.String("service", "shortener")),
	}
}
