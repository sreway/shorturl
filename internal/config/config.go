// Package config implements and describes the application configuration.
package config

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/caarlos0/env/v7"
)

// Config describes the implementation of the application configuration.
type Config interface {
	GetHTTP() *http
	GetShortURL() *shortURL
	GetStorage() *storage
	GetGRPC() *grpc
}

// HTTP describes the implementation of the http server configuration.
type HTTP interface {
	GetScheme() string
	GetAddress() string
	GetCompressTypes() []string
	GetCompressLevel() int
	GetCookie() *cookie
	GetSwagger() *swagger
	GetTLS() *tls
	GetTrustedSubnet() *net.IPNet
}

// GRPC describes the implementation of the grpc server configuration.
type GRPC interface {
	GetTLS() *tls
	UseTLS() bool
	Enabled() bool
	GetAddress() string
}

// ShortURL describes the implementation of the URL shortening service configuration.
type ShortURL interface {
	GetBaseURL() *url.URL
	GetCheckTaskInterval() time.Duration
	GetMaxTaskQueue() int
}

// Storage describes the implementation of the application storage configuration.
type Storage interface {
	GetPostgres() *postgres
	GetCache() *cache
}

// Postgres describes the implementation of the PostgreSQL storage configuration.
type Postgres interface {
	GetDSN() string
	GetMigrateURL() string
}

// Cache describes the implementation of the in-memory storage configuration.
type Cache interface {
	GetFilePath() string
}

// Swagger describes the implementation of th Swagger configuration.
type Swagger interface {
	GetTitle() string
	GetDescription() string
	GetHost() string
	GetBasePath() string
	GetSchemes() []string
}

// config implements application configuration.
type config struct {
	HTTP     *http     `json:"http"`
	GRPC     *grpc     `json:"grpc"`
	ShortURL *shortURL `json:"short_url"`
	Storage  *storage  `json:"storage"`
}

// http implements http server configuration.
type http struct {
	Scheme        string   `json:"scheme" env:"SERVER_SCHEME"`
	Address       string   `json:"server_address" env:"SERVER_ADDRESS"`
	CompressTypes []string `json:"compress_types" env:"HTTP_COMPRESS_TYPES" envSeparator:","`
	CompressLevel int      `json:"compress_level" env:"HTTP_COMPRESS_LEVEL"`
	EnableHTTPS   bool     `json:"enable_https" env:"ENABLE_HTTPS"`
	Cookie        *cookie  `json:"cookie"`
	TLS           *tls     `json:"tls"`
	Swagger       *swagger `json:"swagger"`
	TrustedSubnet *subnet  `json:"trusted_subnet" env:"TRUSTED_SUBNET"`
}

// grpc implements grpc server configuration.
type grpc struct {
	Enable    bool   `json:"enable"`
	Address   string `json:"server_address" env:"SERVER_ADDRESS"`
	EnableTLS bool   `json:"enable_tls"`
	TLS       *tls   `json:"tls"`
}

// subnet describes ip subnet type.
type subnet net.IPNet

// UnmarshalText implements custom unmarshal subnet data.
func (s *subnet) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*s = subnet(net.IPNet{})
		return nil
	}

	_, p, err := net.ParseCIDR(string(text))
	if err != nil {
		return err
	}

	*s = subnet(*p)

	return nil
}

// cookie implements http server cookies configuration.
type cookie struct {
	SignID    string `json:"sign_id" env:"COOKIE_SIGN_ID"`
	SecretKey string `json:"secret_key" env:"COOKIE_SECRET_KEY"`
}

// tls implements http/grpc server tls configuration.
type tls struct {
	CertPath string `json:"cert_path" env:"TLS_CERT_PATH"`
	KeyPath  string `json:"key_path" env:"TLS_KET_PATH"`
}

// shortURL implements shortener configuration.
type shortURL struct {
	BaseURL           *url.URL      `json:"base_url" env:"BASE_URL"`
	CheckTaskInterval time.Duration `json:"check_task_interval" env:"CHECK_TASK_INTERVAL"`
	MaxTaskQueue      int           `json:"max_task_queue" env:"MAX_TASK_QUEUE"`
}

// storage implements storage configuration.
type storage struct {
	Cache    *cache    `json:"cache"`
	Postgres *postgres `json:"postgres"`
}

// cache implements in-memory storage configuration.
type cache struct {
	FilePath string `json:"file_path" env:"FILE_STORAGE_PATH"`
}

// postgres implements postgres configuration.
type postgres struct {
	DSN        string `json:"dsn" env:"DATABASE_DSN"`
	MigrateURL string `json:"migrate_url" env:"MIGRATE_URL"`
}

// swagger implements swagger configuration.
type swagger struct {
	Title       string   `json:"title" env:"SWAGGER_TITLE"`
	Description string   `json:"description" env:"SWAGGER_DESCRIPTION"`
	Host        string   `json:"host" env:"SWAGGER_HOST"`
	BasePath    string   `json:"basePath" env:"SWAGGER_BASE_PATH"`
	Schemes     []string `json:"schemes" env:"SWAGGER_SCHEMES" envSeparator:":"`
}

// GetHTTP implements getting http server configuration.
func (c *config) GetHTTP() *http {
	return c.HTTP
}

// GetShortURL implements getting URL shortening service configuration.
func (c *config) GetShortURL() *shortURL {
	return c.ShortURL
}

// GetStorage implements getting storage configuration.
func (c *config) GetStorage() *storage {
	return c.Storage
}

// GetGRPC implements getting grpc configuration.
func (c *config) GetGRPC() *grpc {
	return c.GRPC
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
	return h.Cookie
}

// GetTLS implements getting http server tls configuration.
func (h *http) GetTLS() *tls {
	return h.TLS
}

// GetCompressLevel implements getting http server compression level.
func (h *http) GetCompressLevel() int {
	return h.CompressLevel
}

// GetSwagger implements getting Swagger configuration.
func (h *http) GetSwagger() *swagger {
	return h.Swagger
}

// GetTrustedSubnet implements getting trusted subnet.
func (h *http) GetTrustedSubnet() *net.IPNet {
	return (*net.IPNet)(h.TrustedSubnet)
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

// GetCache implements getting in-memory storage configuration.
func (store *storage) GetCache() *cache {
	return store.Cache
}

// GetPostgres implements getting PostgreSQL storage configuration.
func (store *storage) GetPostgres() *postgres {
	return store.Postgres
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

// Enabled implements getting information about the need to use grpc service.
func (g *grpc) Enabled() bool {
	return g.Enable
}

// GetTLS implements  getting tls configuration
func (g *grpc) GetTLS() *tls {
	return g.TLS
}

// UseTLS implements getting information about the need to use tls configuration.
func (g *grpc) UseTLS() bool {
	return g.EnableTLS
}

// GetAddress implements getting grpc server address.
func (g *grpc) GetAddress() string {
	return g.Address
}

// NewConfig implements the creation of the application configuration.
func NewConfig() (*config, error) {
	cfg := defaultConfig()
	jsonConfigPath := os.Getenv("CONFIG")
	if len(jsonConfigPath) > 0 {
		fileObj, err := os.Open(jsonConfigPath)
		if err != nil {
			return nil, err
		}
		defer fileObj.Close()
		if err = json.NewDecoder(fileObj).Decode(cfg); err != nil {
			return nil, err
		}
	}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	if cfg.HTTP.EnableHTTPS {
		cfg.HTTP.Scheme = "https"
	} else {
		cfg.HTTP.Scheme = "http"
	}

	cfg.ShortURL.BaseURL = &url.URL{Scheme: cfg.HTTP.Scheme, Host: cfg.HTTP.Address}
	cfg.HTTP.Swagger.Host = fmt.Sprintf("%s://%s", cfg.HTTP.Scheme, cfg.HTTP.Address)
	cfg.HTTP.Swagger.Schemes = append(cfg.HTTP.Swagger.Schemes, cfg.HTTP.Scheme)
	return cfg, nil
}

// defaultConfig implements create application configuration with default values.
func defaultConfig() *config {
	return &config{
		HTTP: &http{
			Scheme:  "http",
			Address: "127.0.0.1:8080",
			CompressTypes: []string{
				"text/plain", "application/json",
			},
			CompressLevel: 5,
			EnableHTTPS:   false,
			TLS: &tls{
				CertPath: "./certs/server.crt",
				KeyPath:  "./certs/server.key",
			},
			Cookie: &cookie{
				SignID:    "user_id",
				SecretKey: "secret_key",
			},
			Swagger: &swagger{
				Title: "Shortener API",
			},
		},
		GRPC: &grpc{
			Address:   "127.0.0.1:8080",
			EnableTLS: false,
			TLS: &tls{
				CertPath: "./certs/server.crt",
				KeyPath:  "./certs/server.key",
			},
		},
		Storage: &storage{
			Cache: &cache{
				FilePath: "./storage.json",
			},
			Postgres: &postgres{
				MigrateURL: "file://migrations/postgres",
			},
		},
		ShortURL: &shortURL{
			CheckTaskInterval: 5 * time.Second,
			MaxTaskQueue:      100,
		},
	}
}
