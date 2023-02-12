package shortener

import (
	"context"
	"net/url"
	"sync/atomic"

	"go.uber.org/zap"

	entity "github.com/sreway/shorturl/internal/domain/url"
	"github.com/sreway/shorturl/internal/usecases/adapters/storage"
	"github.com/sreway/shorturl/pkg/tools/base62"
	"github.com/sreway/shorturl/pkg/tools/logger"
)

type UseCase struct {
	storage storage.URL
	baseURL *url.URL
	counter uint64
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
		Scheme: uc.baseURL.Scheme,
		Host:   uc.baseURL.Host,
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
		Scheme: uc.baseURL.Scheme,
		Host:   uc.baseURL.Host,
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

func New(s storage.URL, baseURL *url.URL, counter uint64) *UseCase {
	l := logger.GetLogger()

	return &UseCase{
		counter: counter,
		baseURL: baseURL,
		storage: s,
		logger:  l.With(zap.String("service", "shortener")),
	}
}
