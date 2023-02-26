package config

import (
	"net/url"

	"github.com/caarlos0/env/v6"
)

type (
	Config interface {
		HTTP() *http
		ShortURL() *shortURL
		Storage() *storage
	}

	HTTP interface {
		GetScheme() string
		GetAddress() string
	}

	ShortURL interface {
		GetBaseURL() *url.URL
		GetCounter() uint64
	}

	Cache interface {
		GetFilePath() string
	}

	config struct {
		http     *http
		shortURL *shortURL
		storage  *storage
	}

	http struct {
		Scheme  string `env:"SERVER_SCHEME" envDefault:"http"`
		Address string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
	}

	shortURL struct {
		BaseURL *url.URL `env:"BASE_URL" envDefault:"http://127.0.0.1:8080"`
		Counter uint64   `env:"COUNTER" envDefault:"1000000000"`
	}

	storage struct {
		cache *cache
	}

	cache struct {
		FilePath string `env:"FILE_STORAGE_PATH" envDefault:"./storage.json"`
	}
)

func (c *config) HTTP() *http {
	return c.http
}

func (c *config) ShortURL() *shortURL {
	return c.shortURL
}

func (c *config) Storage() *storage {
	return c.storage
}

func (h *http) GetScheme() string {
	return h.Scheme
}

func (h *http) GetAddress() string {
	return h.Address
}

func (s *shortURL) GetBaseURL() *url.URL {
	return s.BaseURL
}

func (s *shortURL) GetCounter() uint64 {
	return s.Counter
}

func (store *storage) Cache() *cache {
	return store.cache
}

func (c *cache) GetFilePath() string {
	return c.FilePath
}

func NewConfig() (*config, error) {
	cfg := new(config)
	cfg.shortURL = new(shortURL)
	cfg.http = new(http)
	cfg.storage = new(storage)
	cfg.storage.cache = new(cache)

	if err := env.Parse(cfg.http); err != nil {
		return nil, err
	}

	if err := env.Parse(cfg.shortURL); err != nil {
		return nil, err
	}

	if err := env.Parse(cfg.storage.cache); err != nil {
		return nil, err
	}

	return cfg, nil
}

func NewHTTPConfig(scheme, address string) *http {
	return &http{
		scheme,
		address,
	}
}

func NewShortURLConfig(u *url.URL, counter uint64) *shortURL {
	return &shortURL{
		u,
		counter,
	}
}
