package config

import (
	"net/url"

	"github.com/caarlos0/env/v7"
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
		GetCompressTypes() []string
		GetCompressLevel() int
		GetCookie() *cookie
	}

	ShortURL interface {
		GetBaseURL() *url.URL
	}

	Storage interface {
		GetPostgres() *postgres
		GetCache() *cache
	}

	Postgres interface {
		GetDSN() string
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
		Scheme        string   `env:"SERVER_SCHEME" envDefault:"http"`
		Address       string   `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`
		CompressTypes []string `env:"HTTP_COMPRESS_TYPES" envDefault:"text/plain,application/json" envSeparator:","`
		CompressLevel int      `env:"HTTP_COMPRESS_LEVEL" envDefault:"5"`
		cookie        *cookie
		SecretKey     string `env:"HTTP_SECRET_KEY" envDefault:"secret"`
	}

	cookie struct {
		SignID    string `env:"COOKIE_SIGN_ID" envDefault:"user_id"`
		SecretKey string `env:"COOKIE_SECRET_KEY" envDefault:"secret_key"`
	}

	shortURL struct {
		BaseURL *url.URL `env:"BASE_URL" envDefault:"http://127.0.0.1:8080"`
	}

	storage struct {
		cache    *cache
		postgres *postgres
	}

	cache struct {
		FilePath string `env:"FILE_STORAGE_PATH" envDefault:"./storage.json"`
	}

	postgres struct {
		DSN string `env:"DATABASE_DSN"`
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

func (h *http) GetCompressTypes() []string {
	return h.CompressTypes
}

func (h *http) GetCookie() *cookie {
	return h.cookie
}

func (h *http) GetCompressLevel() int {
	return h.CompressLevel
}

func (s *shortURL) GetBaseURL() *url.URL {
	return s.BaseURL
}

func (store *storage) Cache() *cache {
	return store.cache
}

func (store *storage) Postgres() *postgres {
	return store.postgres
}

func (c *cache) GetFilePath() string {
	return c.FilePath
}

func (s *postgres) GetDSN() string {
	return s.DSN
}

func NewConfig() (*config, error) {
	cfg := new(config)
	cfg.shortURL = new(shortURL)
	cfg.http = new(http)
	cfg.http.cookie = new(cookie)
	cfg.storage = new(storage)
	cfg.storage.cache = new(cache)
	cfg.storage.postgres = new(postgres)

	if err := env.Parse(cfg.http); err != nil {
		return nil, err
	}

	if err := env.Parse(cfg.shortURL); err != nil {
		return nil, err
	}

	if err := env.Parse(cfg.storage.cache); err != nil {
		return nil, err
	}

	if err := env.Parse(cfg.storage.postgres); err != nil {
		return nil, err
	}

	if err := env.Parse(cfg.http.cookie); err != nil {
		return nil, err
	}

	return cfg, nil
}
