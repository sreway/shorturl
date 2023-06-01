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
		GetSwagger() *swagger
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
	// Swagger describes the implementation of th Swagger configuration.
	Swagger interface {
		GetTitle() string
		GetDescription() string
		GetHost() string
		GetBasePath() string
		GetSchemes() []string
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
		swagger       *swagger
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

	swagger struct {
		Title       string   `json:"title" env:"SWAGGER_TITLE" envDefault:"Shortener API"`
		Description string   `json:"description" env:"SWAGGER_DESCRIPTION"`
		Host        string   `json:"host" env:"SWAGGER_HOST" envDefault:"127.0.0.1:8080"`
		BasePath    string   `json:"basePath" env:"SWAGGER_BASE_PATH"`
		Schemes     []string `json:"schemes" env:"SWAGGER_SCHEMES" envSeparator:":" envDefault:"http"`
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

// GetSwagger implements getting Swagger configuration.
func (h *http) GetSwagger() *swagger {
	return h.swagger
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

// GetTitle implements getting Swagger title.
func (s *swagger) GetTitle() string {
	return s.Title
}

// GetDescription implements getting Swagger description.
func (s *swagger) GetDescription() string {
	return s.Description
}

// GetHost implements getting Swagger host.
func (s *swagger) GetHost() string {
	return s.Host
}

// GetBasePath implements getting Swagger base path.
func (s *swagger) GetBasePath() string {
	return s.BasePath
}

// GetSchemes implements getting Swagger schemes.
func (s *swagger) GetSchemes() []string {
	return s.Schemes
}

// NewConfig implements the creation of the application configuration.
func NewConfig() (*config, error) {
	cfg := new(config)
	cfg.shortURL = new(shortURL)
	cfg.http = new(http)
	cfg.http.cookie = new(cookie)
	cfg.http.swagger = new(swagger)
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

	if err := env.Parse(cfg.http.swagger); err != nil {
		return nil, err
	}

	return cfg, nil
}
