// Package config implements and describes the application configuration.
package config

import (
	"net/url"
	"time"

	"github.com/caarlos0/env/v7"
)

type (
	// Config describes the implementation of the application configuration.
	Config interface {
		HTTP() *http
		ShortURL() *shortURL
		Storage() *storage
	}
	// HTTP describes the implementation of the http server configuration.
	HTTP interface {
		GetScheme() string
		GetAddress() string
		GetCompressTypes() []string
		GetCompressLevel() int
		GetCookie() *cookie
	}
	// ShortURL describes the implementation of the URL shortening service configuration.
	ShortURL interface {
		GetBaseURL() *url.URL
		GetCheckTaskInterval() time.Duration
		GetMaxTaskQueue() int
	}
	// Storage describes the implementation of the application storage configuration.
	Storage interface {
		GetPostgres() *postgres
		GetCache() *cache
	}
	// Postgres describes the implementation of the PostgreSQL storage configuration.
	Postgres interface {
		GetDSN() string
		GetMigrateURL() string
	}
	// Cache describes the implementation of the in-memory storage configuration.
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
		BaseURL           *url.URL      `env:"BASE_URL" envDefault:"http://127.0.0.1:8080"`
		CheckTaskInterval time.Duration `env:"CHECK_TASK_INTERVAL" envDefault:"5s"`
		MaxTaskQueue      int           `env:"MAX_TASK_QUEUE" envDefault:"100"`
	}

	storage struct {
		cache    *cache
		postgres *postgres
	}

	cache struct {
		FilePath string `env:"FILE_STORAGE_PATH" envDefault:"./storage.json"`
	}

	postgres struct {
		DSN        string `env:"DATABASE_DSN"`
		MigrateURL string `env:"MIGRATE_URL" envDefault:"file://migrations/postgres"`
	}
)

// HTTP implements getting http server configuration.
func (c *config) HTTP() *http {
	return c.http
}

// ShortURL implements getting URL shortening service configuration.
func (c *config) ShortURL() *shortURL {
	return c.shortURL
}

// Storage implements getting storage configuration.
func (c *config) Storage() *storage {
	return c.storage
}

// GetScheme implements getting http server scheme (http/https).
func (h *http) GetScheme() string {
	return h.Scheme
}

// GetAddress implements getting http server address.
func (h *http) GetAddress() string {
	return h.Address
}

// GetCompressTypes implements getting http server compression types.
func (h *http) GetCompressTypes() []string {
	return h.CompressTypes
}

// GetCookie implements getting http server cookie configuration.
func (h *http) GetCookie() *cookie {
	return h.cookie
}

// GetCompressLevel implements getting http server compression level.
func (h *http) GetCompressLevel() int {
	return h.CompressLevel
}

// GetBaseURL implements getting the base URL for the URL shortening service.
func (s *shortURL) GetBaseURL() *url.URL {
	return s.BaseURL
}

// GetMaxTaskQueue implements getting the task limit for the URL shortening service.
func (s *shortURL) GetMaxTaskQueue() int {
	return s.MaxTaskQueue
}

// GetCheckTaskInterval implements getting the task verification interval for the URL shortening service.
func (s *shortURL) GetCheckTaskInterval() time.Duration {
	return s.CheckTaskInterval
}

// Cache implements getting in-memory storage configuration.
func (store *storage) Cache() *cache {
	return store.cache
}

// Postgres implements getting PostgreSQL storage configuration.
func (store *storage) Postgres() *postgres {
	return store.postgres
}

// GetFilePath implements getting the file path for the in-memory storage.
func (c *cache) GetFilePath() string {
	return c.FilePath
}

// GetDSN implements getting the DSN URL for PostgreSQL storage.
func (p *postgres) GetDSN() string {
	return p.DSN
}

// GetMigrateURL implements getting migration location for PostgreSQL storage.
func (p *postgres) GetMigrateURL() string {
	return p.MigrateURL
}

// NewConfig implements the creation of the application configuration.
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
